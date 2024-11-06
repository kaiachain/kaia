package impl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

type subsetCouncilSlice []common.Address

func (sc subsetCouncilSlice) Len() int {
	return len(sc)
}

func (sc subsetCouncilSlice) Less(i, j int) bool {
	return strings.Compare(sc[i].String(), sc[j].String()) < 0
}

func (sc subsetCouncilSlice) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

func (sc subsetCouncilSlice) AddressStringList() []string {
	stringAddrs := make([]string, len(sc))
	for _, val := range sc {
		stringAddrs = append(stringAddrs, val.String())
	}
	return stringAddrs
}

func (sc subsetCouncilSlice) getIdxByAddress(addr common.Address) int {
	for i, val := range sc {
		if addr == val {
			return i
		}
	}
	// TODO-Kaia-Istanbul: Enable this log when non-committee nodes don't call `core.startNewRound()`
	// logger.Warn("failed to find an address in the validator list",
	// 	"address", addr, "validatorAddrs", valSet.validators.AddressStringList())
	return -1
}

// SortedAddressList retrieves the sorted address list of ValidatorSet in "ascending order".
// if public is false, sort it using bytes.Compare. It's for public purpose.
// - public-false usage: (getValidators/getDemotedValidators, defaultSet snap store, prepareExtra.validators)
// if public is true, sort it using strings.Compare. It's used for internal consensus purpose, especially for the source of committee.
// - public-true usage: (snap read/store/apply except defaultSet snap store, vrank log)
// TODO-kaia-valset: unify sorting.
func (sc subsetCouncilSlice) sortedAddressList(public bool) []common.Address {
	copiedSlice := make(subsetCouncilSlice, len(sc))
	copy(copiedSlice, sc)

	if public {
		// want reverse-sort: ascending order - bytes.Compare(ValidatorSet[i][:], ValidatorSet[j][:]) > 0
		sort.Slice(copiedSlice, func(i, j int) bool {
			return bytes.Compare(copiedSlice[i].Bytes(), copiedSlice[j].Bytes()) >= 0
		})
		sort.Sort(sort.Reverse(copiedSlice))
	} else {
		// want sort: descending order - strings.Compare(ValidatorSet[i].String(), ValidatorSet[j].String()) < 0
		sort.Sort(copiedSlice)
	}
	return copiedSlice
}

// council for block N is created based on the previous block's council.
// It is used to calculate committee or proposer. It's not display purpose.
// councilList(N-1) -> apply demoted rule -> council(N) -> block proposed and apply votes -> councilList(N)
type council struct {
	blockNumber uint64

	qualifiedValidators subsetCouncilSlice // qualified is a subset of prev block's council list
	demotedValidators   subsetCouncilSlice // demoted is a subset of prev block's council who are demoted as a member of committee

	councilAddressList subsetCouncilSlice // total council node address list. the order is reserved.
}

func newCouncil(db database.Database, valCtx *valSetContext) (*council, error) {
	var (
		blockNumber    = valCtx.blockNumber
		rules          = valCtx.rules
		proposerPolicy = valCtx.prevBlockResult.proposerPolicy
		stakingAmount  = valCtx.prevBlockResult.consolidatedStakingAmount()

		isSingleMode = valCtx.prevBlockResult.pSet.GovernanceMode == "single"
		govNode      = valCtx.prevBlockResult.pSet.GoverningNode
		minStaking   = valCtx.prevBlockResult.pSet.MinimumStake.Uint64()
	)

	// read council list of block N-1
	councilAddressList, err := ReadCouncilAddressListFromDb(db, blockNumber-1)
	if err != nil {
		return nil, err
	}

	// create council struct for block N
	c := &council{
		blockNumber:        blockNumber,
		councilAddressList: make(subsetCouncilSlice, len(councilAddressList)),
	}
	copy(c.councilAddressList, councilAddressList)

	// either proposer-policy is not weighted random or before istanbul HF,
	// do not filter out the demoted validators
	if !proposerPolicy.IsWeightedRandom() || !rules.IsIstanbul {
		c.qualifiedValidators = make(subsetCouncilSlice, len(councilAddressList))
		copy(c.qualifiedValidators, councilAddressList)
		return c, nil
	}

	// Split the councilList(N-1) into qualified and demoted.
	// Qualified is a subset of the council who are qualified to be a committee member.
	// (1) If governance mode is single, always include the governing node.
	// (2) If no council members has enough KAIA, all members become qualified.
	for addr, val := range stakingAmount {
		if val >= minStaking || (isSingleMode && addr == govNode) {
			c.qualifiedValidators = append(c.qualifiedValidators, addr)
		} else {
			c.demotedValidators = append(c.qualifiedValidators, addr)
		}
	}

	// include all council members if case1 or case2
	//   case1. not a single mode && no qualified
	//   case2. single mode && len(qualified) is 1 && govNode is not qualified
	if len(c.qualifiedValidators) == 0 || (isSingleMode && len(c.qualifiedValidators) == 1 && stakingAmount[govNode] < minStaking) {
		c.demotedValidators = subsetCouncilSlice{} // ensure demoted is empty
		c.qualifiedValidators = make(subsetCouncilSlice, len(councilAddressList))
		copy(c.qualifiedValidators, councilAddressList)
	}

	sort.Sort(c.qualifiedValidators)
	sort.Sort(c.demotedValidators)

	return c, nil
}

func (c *council) getProposerFromProposers(proposers []common.Address, round uint64, proposerUpdateInterval uint64) (common.Address, int) {
	proposer := proposers[int(c.blockNumber+round)%int(proposerUpdateInterval)%len(proposers)]
	proposerIdx := c.qualifiedValidators.getIdxByAddress(proposer)
	return proposer, proposerIdx
}

// selectRandomCommittee composes a committee selecting validators randomly based on the seed value.
// It returns nil if the given committeeSize is bigger than validatorSize or proposer indexes are invalid.
func (c *council) selectRandomCommittee(valCtx *valSetContext, round uint64, proposers []common.Address) ([]common.Address, error) {
	var (
		validatorSize   = len(c.qualifiedValidators)
		validator       = c.qualifiedValidators
		rules           = valCtx.rules
		committeeSize   = valCtx.prevBlockResult.pSet.CommitteeSize
		prevHeader      = valCtx.prevBlockResult.header
		pUpdateInterval = valCtx.prevBlockResult.pSet.ProposerUpdateInterval

		// pick current round's proposer and next proposer which address is different from current proposer
		proposer, proposerIdx                                 = c.getProposerFromProposers(proposers, round, pUpdateInterval)
		closestDifferentProposer, closestDifferentProposerIdx = c.getProposerFromProposers(proposers, round+1, pUpdateInterval)
	)

	// return early if the committee size is 1
	if committeeSize == 1 {
		return []common.Address{proposer}, nil
	}

	// closest next proposer who has different address with the proposer
	for i := uint64(2); i < pUpdateInterval; i++ {
		if proposer != closestDifferentProposer {
			break
		}
		closestDifferentProposer, closestDifferentProposerIdx = c.getProposerFromProposers(proposers, round+i, pUpdateInterval)
	}

	// ensure validator indexes are valid
	if proposerIdx < 0 || closestDifferentProposerIdx < 0 || proposerIdx == closestDifferentProposerIdx ||
		validatorSize <= proposerIdx || validatorSize <= closestDifferentProposerIdx {
		return nil, fmt.Errorf("invalid indexes of validators. validatorSize: %d, proposerIdx:%d, nextProposerIdx:%d",
			validatorSize, proposerIdx, closestDifferentProposerIdx)
	}

	seed, err := convertHashToSeed(prevHeader.Hash())
	if err != nil {
		return nil, err
	}
	// shuffle the qualified validators except two proposers
	if rules.IsIstanbul {
		seed += int64(round)
	}
	committee := make([]common.Address, committeeSize)
	picker := rand.New(rand.NewSource(seed))
	pickSize := validatorSize - 2
	indexs := make([]int, pickSize)
	idx := 0
	for i := 0; i < validatorSize; i++ {
		if i != proposerIdx && i != closestDifferentProposerIdx {
			indexs[idx] = i
			idx++
		}
	}
	for i := 0; i < pickSize; i++ {
		randIndex := picker.Intn(pickSize)
		indexs[i], indexs[randIndex] = indexs[randIndex], indexs[i]
	}

	// first committee is the proposer and the second committee is the next proposer
	committee[0], committee[1] = proposer, closestDifferentProposer
	for i := uint64(0); i < committeeSize-2; i++ {
		committee[i+2] = validator[indexs[i]]
	}

	return committee, nil
}

// SelectRandaoCommittee composes a committee selecting validators randomly based on the mixHash.
// It is guaranteed that the committee include proposers for all rounds because
// the proposer is picked from the this committee. See weightedRandomProposer().
//
// def select_committee_KIP146(validators, committee_size, seed):
//
//	shuffled = shuffle_validators_KIP146(validators, seed)
//	return shuffled[:min(committee_size, len(validators))]
func (c *council) selectRandaoCommittee(prevBnRes *blockResult) ([]common.Address, error) {
	if prevBnRes.header.MixHash == nil {
		return nil, errNilMixHash
	}

	copied := make(subsetCouncilSlice, len(c.qualifiedValidators))
	copy(copied, c.qualifiedValidators)

	seed := int64(binary.BigEndian.Uint64(prevBnRes.header.MixHash[:8]))
	rand.New(rand.NewSource(seed)).Shuffle(len(copied), copied.Swap)
	return copied[:prevBnRes.pSet.CommitteeSize], nil
}
