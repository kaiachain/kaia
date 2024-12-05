package impl

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sort"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

type committeeContext struct {
	qualified valset.AddressList
	num       uint64
	rules     params.Rules

	// pSet
	committeeSize          uint64
	proposerPolicy         ProposerPolicy
	proposerUpdateInterval uint64

	// num-1
	prevHeader *types.Header
	prevAuthor common.Address

	// num - (num % proposerUpdateInterval)
	pUpdateBlock uint64
	proposers    []common.Address
}

func newCommitteeContext(v *ValsetModule, num uint64) (*committeeContext, error) {
	if num < 1 {
		return nil, errGenesisNotCalculable
	}

	pSet, proposerPolicy, err := v.getPSetWithProposerPolicy(num)
	if err != nil {
		return nil, err
	}

	prevHeader := v.chain.GetHeaderByNumber(num - 1)
	if prevHeader == nil {
		return nil, errNilHeader
	}

	qualified, err := v.getQualifiedValidators(num)
	if err != nil {
		return nil, err
	}

	prevAuthor, err := v.chain.Engine().Author(prevHeader)
	if err != nil {
		return nil, err
	}

	var (
		rules        = v.chain.Config().Rules(big.NewInt(int64(num)))
		pUpdateBlock = uint64(0)
		proposers    = []common.Address(nil)
	)

	if proposerPolicy.IsWeightedRandom() && !rules.IsRandao {
		pUpdateBlock = calcProposerBlockNumber(num, pSet.ProposerUpdateInterval)
		proposers, err = v.getLegacyProposersList(pUpdateBlock)
		if err != nil {
			return nil, err
		}
	}

	return &committeeContext{
		qualified: qualified,
		num:       num,
		rules:     rules,

		committeeSize:          pSet.CommitteeSize,
		proposerPolicy:         proposerPolicy,
		proposerUpdateInterval: pSet.ProposerUpdateInterval,

		prevHeader: prevHeader,
		prevAuthor: prevAuthor,

		pUpdateBlock: pUpdateBlock,
		proposers:    proposers,
	}, nil
}

func (c *committeeContext) getCommittee(round uint64) ([]common.Address, error) {
	// if the committee size is bigger than qualified valSet size, return all qualified
	if c.committeeSize >= uint64(len(c.qualified)) {
		return c.qualified, nil
	}

	proposer, err := c.getProposer(round)
	if err != nil {
		return nil, err
	}

	// return early if the committee size is 1
	if c.committeeSize == 1 {
		return []common.Address{proposer}, nil
	}

	if !c.proposerPolicy.IsWeightedRandom() || !c.rules.IsRandao {
		repetition := uint64(2)
		if c.proposerPolicy.IsWeightedRandom() && c.rules.IsRandao {
			repetition = c.proposerUpdateInterval
		}

		// closest next proposer who has different address with the proposer
		// pick current round's proposer and next proposer which address is different from current proposer
		var nextDistinctProposer common.Address
		for i := uint64(1); i < repetition; i++ {
			nextDistinctProposer, err = c.getProposer(round + i)
			if err != nil {
				return nil, err
			}
			if proposer != nextDistinctProposer {
				break
			}
		}

		return c.selectRandomCommittee(round, proposer, nextDistinctProposer)
	}
	return c.selectRandaoCommittee()
}

// selectRandomCommittee composes a committee selecting validators randomly based on the seed value.
// It returns nil if the given committeeSize is bigger than validatorSize or proposer indexes are invalid.
func (c *committeeContext) selectRandomCommittee(round uint64, proposer, nextDistinctProposer common.Address) ([]common.Address, error) {
	proposerIdx := c.qualified.GetIdxByAddress(proposer)
	nextDistinctProposerIdx := c.qualified.GetIdxByAddress(nextDistinctProposer)

	// ensure validator indexes are valid
	validatorSize := len(c.qualified)
	if proposerIdx < 0 || nextDistinctProposerIdx < 0 || proposerIdx == nextDistinctProposerIdx ||
		validatorSize <= proposerIdx || validatorSize <= nextDistinctProposerIdx {
		return nil, fmt.Errorf("invalid indexes of validators. validatorSize: %d, proposerIdx:%d, nextDistinctProposerIdx:%d",
			validatorSize, proposerIdx, nextDistinctProposerIdx)
	}

	seed, err := convertHashToSeed(c.prevHeader.Hash())
	if err != nil {
		return nil, err
	}
	// shuffle the qualified validators except two proposers
	if c.rules.IsIstanbul {
		seed += int64(round)
	}
	committee := make([]common.Address, c.committeeSize)
	picker := rand.New(rand.NewSource(seed))
	pickSize := validatorSize - 2
	indexs := make([]int, pickSize)
	idx := 0
	for i := 0; i < validatorSize; i++ {
		if i != proposerIdx && i != nextDistinctProposerIdx {
			indexs[idx] = i
			idx++
		}
	}
	for i := 0; i < pickSize; i++ {
		randIndex := picker.Intn(pickSize)
		indexs[i], indexs[randIndex] = indexs[randIndex], indexs[i]
	}

	// first committee is the proposer and the second committee is the next proposer
	committee[0], committee[1] = proposer, nextDistinctProposer
	for i := uint64(0); i < c.committeeSize-2; i++ {
		committee[i+2] = c.qualified[indexs[i]]
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
func (c *committeeContext) selectRandaoCommittee() ([]common.Address, error) {
	if c.prevHeader.MixHash == nil {
		return nil, errNilMixHash
	}

	copied := make(valset.AddressList, len(c.qualified))
	copy(copied, c.qualified)

	seed := int64(binary.BigEndian.Uint64(c.prevHeader.MixHash[:8]))
	rand.New(rand.NewSource(seed)).Shuffle(len(copied), copied.Swap)
	return copied[:c.committeeSize], nil
}

func (c *committeeContext) getProposer(round uint64) (common.Address, error) {
	if c.proposerPolicy.IsDefaultSet() {
		// if the policy is round-robin or sticky, all the council members are qualified.
		// be cautious that the proposer may not be included in the committee list.
		copied := make(valset.AddressList, len(c.qualified))
		copy(copied, c.qualified)

		// sorting on council address list
		sort.Sort(copied)

		proposer, _ := pickRoundRobinProposer(c.qualified, c.proposerPolicy, c.prevAuthor, round)
		return proposer, nil
	}

	// before Randao, weightedrandom uses proposers to pick the proposer.
	if !c.rules.IsRandao {
		proposer, _ := pickWeightedRandomProposer(c.proposers, c.pUpdateBlock, c.num, round, c.qualified, c.prevAuthor)
		return proposer, nil
	}

	// after Randao, pick proposer from randao committee
	committee, err := c.selectRandaoCommittee()
	if err != nil {
		return common.Address{}, err
	}
	return committee[int(round)%len(c.qualified)], nil
}
