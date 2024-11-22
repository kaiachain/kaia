package impl

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// council for block N is created based on the previous block's council.
// It is used to calculate committee or proposer. It's not display purpose.
// councilList(N-1) -> apply demoted rule -> council(N) -> block proposed and apply votes -> councilList(N)
type council struct {
	blockNumber uint64

	qualifiedValidators valset.AddressList // qualified is a subset of prev block's council list
	demotedValidators   valset.AddressList // demoted is a subset of prev block's council who are demoted as a member of committee

	councilAddressList valset.AddressList // total council node address list. the order is reserved.
}

func newCouncil(valCtx *valSetContext, councilAddressList []common.Address) (*council, error) {
	var (
		blockNumber    = valCtx.blockNumber
		rules          = valCtx.rules
		proposerPolicy = valCtx.prevBlockResult.proposerPolicy
		stakingAmount  = valCtx.prevBlockResult.consolidatedStakingAmount()

		isSingleMode = valCtx.prevBlockResult.pSet.GovernanceMode == "single"
		govNode      = valCtx.prevBlockResult.pSet.GoverningNode
		minStaking   = valCtx.prevBlockResult.pSet.MinimumStake.Uint64()
	)

	// create council struct for block N
	c := &council{
		blockNumber:        blockNumber,
		councilAddressList: make(valset.AddressList, len(councilAddressList)),
	}
	copy(c.councilAddressList, councilAddressList)

	// either proposer-policy is not weighted random or before istanbul HF,
	// do not filter out the demoted validators
	if !proposerPolicy.IsWeightedRandom() || !rules.IsIstanbul {
		c.qualifiedValidators = make(valset.AddressList, len(councilAddressList))
		copy(c.qualifiedValidators, councilAddressList)
		return c, nil
	}

	// Split the councilList(N-1) into qualified and demoted.
	// Qualified is a subset of the council who are qualified to be a committee member.
	// (1) If governance mode is single, always include the governing node.
	// (2) If no council members has enough KAIA, all members become qualified.
	for _, addr := range councilAddressList {
		staking, ok := stakingAmount[addr]
		if ok && (staking >= minStaking || (isSingleMode && addr == govNode)) {
			c.qualifiedValidators = append(c.qualifiedValidators, addr)
		} else {
			c.demotedValidators = append(c.demotedValidators, addr)
		}
	}

	// include all council members if case1 or case2
	//   case1. not a single mode && no qualified
	//   case2. single mode && len(qualified) is 1 && govNode is not qualified
	if len(c.qualifiedValidators) == 0 || (isSingleMode && len(c.qualifiedValidators) == 1 && stakingAmount[govNode] < minStaking) {
		c.demotedValidators = valset.AddressList{} // ensure demoted is empty
		c.qualifiedValidators = make(valset.AddressList, len(councilAddressList))
		copy(c.qualifiedValidators, councilAddressList)
	}

	sort.Sort(c.qualifiedValidators)
	sort.Sort(c.demotedValidators)

	return c, nil
}

// selectRandomCommittee composes a committee selecting validators randomly based on the seed value.
// It returns nil if the given committeeSize is bigger than validatorSize or proposer indexes are invalid.
func (c *council) selectRandomCommittee(valCtx *valSetContext, round uint64, proposers []common.Address) ([]common.Address, error) {
	var (
		validatorSize   = len(c.qualifiedValidators)
		validator       = c.qualifiedValidators
		num             = valCtx.blockNumber
		rules           = valCtx.rules
		committeeSize   = valCtx.prevBlockResult.pSet.CommitteeSize
		prevHeader      = valCtx.prevBlockResult.header
		pUpdateInterval = valCtx.prevBlockResult.pSet.ProposerUpdateInterval
		pUpdateBlock    = calcProposerBlockNumber(num, pUpdateInterval)
		author          = valCtx.prevBlockResult.author

		// pick current round's proposer and next proposer which address is different from current proposer
		proposer, proposerIdx                                 = pickWeightedRandomProposer(proposers, pUpdateBlock, num, round, c.qualifiedValidators, author)
		closestDifferentProposer, closestDifferentProposerIdx = pickWeightedRandomProposer(proposers, pUpdateBlock, num, round+1, c.qualifiedValidators, author)
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
		closestDifferentProposer, closestDifferentProposerIdx = pickWeightedRandomProposer(proposers, pUpdateBlock, num, round+i, c.qualifiedValidators, author)
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

	copied := make(valset.AddressList, len(c.qualifiedValidators))
	copy(copied, c.qualifiedValidators)

	seed := int64(binary.BigEndian.Uint64(prevBnRes.header.MixHash[:8]))
	rand.New(rand.NewSource(seed)).Shuffle(len(copied), copied.Swap)
	return copied[:prevBnRes.pSet.CommitteeSize], nil
}
