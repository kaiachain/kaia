package impl

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sort"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// council is used to calculate committee or proposer.
// councilList(N) -> apply demoted rule -> demoted,qualified -> block proposed and apply votes -> councilList(N+1)
type council struct {
	blockNumber uint64
	rules       params.Rules

	qualifiedValidators valset.AddressList // qualified is a list of validators which will be the member of the committee
	demotedValidators   valset.AddressList // demoted are subtract of qualified from council(N)

	councilAddressList valset.AddressList // total council node address list. the order is reserved.
}

func newCouncil(v *ValsetModule, num uint64) (*council, error) {
	c := &council{
		blockNumber: num,
		rules:       v.chain.Config().Rules(big.NewInt(int64(num))),
	}

	// council after processing block N-1
	cList, err := v.GetCouncilAddressList(num)
	if err != nil {
		return nil, err
	}
	c.councilAddressList = make(valset.AddressList, len(cList))
	copy(c.councilAddressList, cList)

	if num == 0 {
		c.qualifiedValidators = make(valset.AddressList, len(cList))
		c.demotedValidators = make(valset.AddressList, 0)
		copy(c.qualifiedValidators, cList)
		return c, nil
	}

	_, stakingAmounts, err := v.getStakingInfoWithStakingAmounts(num, cList)
	if err != nil {
		return nil, err
	}

	pSet, proposerPolicy, err := v.getPSetWithProposerPolicy(num)
	if err != nil {
		return nil, err
	}

	var (
		isSingleMode = pSet.GovernanceMode == "single"
		govNode      = pSet.GoverningNode
		minStaking   = pSet.MinimumStake.Uint64()
	)

	// either proposer-policy is not weighted random or before istanbul HF,
	// do not filter out the demoted validators
	if !proposerPolicy.IsWeightedRandom() || !c.rules.IsIstanbul {
		c.qualifiedValidators = make(valset.AddressList, len(cList))
		copy(c.qualifiedValidators, cList)
		return c, nil
	}

	// Split the council(N) into qualified and demoted.
	// Qualified is a subset of the council who are qualified to be a committee member.
	// (1) If governance mode is single, always include the governing node.
	// (2) If no council members has enough KAIA, all members become qualified.
	for _, addr := range cList {
		staking, ok := stakingAmounts[addr]
		if ok && (staking >= minStaking || (isSingleMode && addr == govNode)) {
			c.qualifiedValidators = append(c.qualifiedValidators, addr)
		} else {
			c.demotedValidators = append(c.demotedValidators, addr)
		}
	}

	// include all council members if case1 or case2
	//   case1. not a single mode && no qualified
	//   case2. single mode && len(qualified) is 1 && govNode is not qualified
	if len(c.qualifiedValidators) == 0 || (isSingleMode && len(c.qualifiedValidators) == 1 && stakingAmounts[govNode] < minStaking) {
		c.demotedValidators = valset.AddressList{} // ensure demoted is empty
		c.qualifiedValidators = make(valset.AddressList, len(cList))
		copy(c.qualifiedValidators, cList)
	}

	sort.Sort(c.qualifiedValidators)
	sort.Sort(c.demotedValidators)

	return c, nil
}

// selectRandomCommittee composes a committee selecting validators randomly based on the seed value.
// It returns nil if the given committeeSize is bigger than validatorSize or proposer indexes are invalid.
func (c *council) selectRandomCommittee(round uint64, pSet gov.ParamSet, proposers []common.Address, prevHeader *types.Header, prevAuthor common.Address) ([]common.Address, error) {
	var (
		pUpdateBlock = calcProposerBlockNumber(c.blockNumber, pSet.ProposerUpdateInterval)

		// pick current round's proposer and next proposer which address is different from current proposer
		proposer, proposerIdx                                 = pickWeightedRandomProposer(proposers, pUpdateBlock, c.blockNumber, round, c.qualifiedValidators, prevAuthor)
		closestDifferentProposer, closestDifferentProposerIdx = pickWeightedRandomProposer(proposers, pUpdateBlock, c.blockNumber, round+1, c.qualifiedValidators, prevAuthor)
	)

	// return early if the committee size is 1
	if pSet.CommitteeSize == 1 {
		return []common.Address{proposer}, nil
	}

	// closest next proposer who has different address with the proposer
	for i := uint64(2); i < pSet.ProposerUpdateInterval; i++ {
		if proposer != closestDifferentProposer {
			break
		}
		closestDifferentProposer, closestDifferentProposerIdx = pickWeightedRandomProposer(proposers, pUpdateBlock, c.blockNumber, round+i, c.qualifiedValidators, prevAuthor)
	}

	// ensure validator indexes are valid
	validatorSize := len(c.qualifiedValidators)
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
	if c.rules.IsIstanbul {
		seed += int64(round)
	}
	committee := make([]common.Address, pSet.CommitteeSize)
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
	for i := uint64(0); i < pSet.CommitteeSize-2; i++ {
		committee[i+2] = c.qualifiedValidators[indexs[i]]
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
func (c *council) selectRandaoCommittee(prevHeader *types.Header, pSet gov.ParamSet) ([]common.Address, error) {
	if prevHeader.MixHash == nil {
		return nil, errNilMixHash
	}

	copied := make(valset.AddressList, len(c.qualifiedValidators))
	copy(copied, c.qualifiedValidators)

	seed := int64(binary.BigEndian.Uint64(prevHeader.MixHash[:8]))
	rand.New(rand.NewSource(seed)).Shuffle(len(copied), copied.Swap)
	return copied[:pSet.CommitteeSize], nil
}
