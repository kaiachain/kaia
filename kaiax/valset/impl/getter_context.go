// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// blockContext is the common input for committee and proposer selection for a block.
// Because committee and proposer selection calls each other, we collect common inputs here.
type blockContext struct {
	num uint64

	// Computed inputs
	qualified    *valset.AddressSet
	prevHeader   *types.Header
	prevProposer common.Address
	// convenience fields
	rules params.Rules
	pset  gov.ParamSet
}

func (v *ValsetModule) getBlockContext(num uint64) (*blockContext, error) {
	qualified, err := v.getQualifiedValidators(num)
	if err != nil {
		return nil, err
	}
	prevHeader := v.Chain.GetHeaderByNumber(num - 1)
	if prevHeader == nil {
		return nil, errNoHeader
	}

	prevProposer := qualified.At(0)
	if num-1 > 0 {
		prevProposer, err = v.Chain.Engine().Author(prevHeader)
		if err != nil {
			return nil, err
		}
	}

	return &blockContext{
		num:          num,
		qualified:    qualified,
		prevHeader:   prevHeader,
		prevProposer: prevProposer,
		rules:        v.Chain.Config().Rules(new(big.Int).SetUint64(num)),
		pset:         v.GovModule.EffectiveParamSet(num),
	}, nil
}

func (v *ValsetModule) getCommittee(c *blockContext, round uint64) ([]common.Address, error) {
	if c.num == 0 {
		return c.qualified.List(), nil
	}
	if c.qualified.Len() <= int(c.pset.CommitteeSize) {
		return c.qualified.List(), nil
	}

	var useRandomCommittee bool
	switch istanbul.ProposerPolicy(c.pset.ProposerPolicy) {
	case istanbul.RoundRobin, istanbul.Sticky:
		useRandomCommittee = true
	case istanbul.WeightedRandom:
		useRandomCommittee = !c.rules.IsRandao
	default:
		return nil, errInvalidProposerPolicy
	}

	if useRandomCommittee {
		currProposer, err := v.getProposer(c, round)
		if err != nil {
			return nil, err
		}
		if c.pset.CommitteeSize == 1 {
			return []common.Address{currProposer}, nil
		}
		// From here, we can assume 1 < committeeSize < qualified.Len, so there must be a nextDistinctProposer.
		nextDistinctProposer, err := v.getNextDistinctProposer(c, currProposer, round)
		if err == errNoNextDistinctProposer {
			// This should not happen, but let the consensus to continue anyway.
			logger.Error("failed to find next distinct proposer", "num", c.num, "round", round)
			return c.qualified.List(), nil
		} else if err != nil {
			return nil, err
		}
		logger.Trace("SelectRandomCommittee", "number", c.num, "round", round, "curr", currProposer.Hex(), "next", nextDistinctProposer.Hex())
		return selectRandomCommittee(c.qualified, c.pset.CommitteeSize, c.prevHeader.Hash(), currProposer, nextDistinctProposer), nil
	} else {
		// In the legacy code weighted.go, the Randao mode also had the same [if committeeSize == 1 return {currProposer}] logic.
		// This logic is implicitly fulfilled: If committeeSize == 1, the committee shall return one address {x},
		// where x must be the current proposer, because the RandaoProposer is picked from the RandaoCommittee.
		return selectRandaoCommittee(c.qualified, c.pset.CommitteeSize, c.prevHeader.MixHash), nil
	}
}

func (v *ValsetModule) getNextDistinctProposer(c *blockContext, currProposer common.Address, round uint64) (common.Address, error) {
	switch istanbul.ProposerPolicy(c.pset.ProposerPolicy) {
	case istanbul.RoundRobin, istanbul.Sticky:
		// For RoundRobin and Sticky, next round proposer has to be different from the current.
		return v.getProposer(c, round+1)
	case istanbul.WeightedRandom:
		if c.rules.IsRandao { // If IsRandao this function is not called.
			return common.Address{}, errInvalidProposerPolicy
		}

		list, sourceNum, err := v.getProposerList(c)
		if err != nil {
			return common.Address{}, err
		}
		// scan one proposer update interval
		for i := uint64(1); i <= c.pset.ProposerUpdateInterval; i++ {
			nextProposer := selectWeightedRandomProposer(list, sourceNum, c.num, round+uint64(i))
			if currProposer != nextProposer {
				return nextProposer, nil
			}
		}
		return common.Address{}, errNoNextDistinctProposer
	default:
		return common.Address{}, errInvalidProposerPolicy
	}
}

func selectRandomCommittee(qualified *valset.AddressSet, committeeSize uint64, prevBlockHash common.Hash, currProposer, nextDistinctProposer common.Address) []common.Address {
	qualified = qualified.Copy()
	qualified.Remove(currProposer)
	qualified.Remove(nextDistinctProposer)

	seed := valset.HashToSeedLegacy(prevBlockHash)
	shuffled := qualified.ShuffledListLegacy(seed)
	committee := []common.Address{currProposer, nextDistinctProposer}
	committee = append(committee, shuffled[:committeeSize-2]...)
	return committee
}

func selectRandaoCommittee(qualified *valset.AddressSet, committeeSize uint64, prevMixHash []byte) []common.Address {
	// Note: If committeeSize == 1, below code is equivalent to return []common.Address{currProposer}.
	// Because, the resulting committee will be one validator, and the only validator is also selected as the proposer.
	// Therefore no special handling for committeeSize == 1 is needed.
	mixHash := prevMixHash
	if len(prevMixHash) == 0 { // At exactly Randao hardfork block, or genesis block, prevMixHash is empty.
		mixHash = params.ZeroMixHash
	}
	seed := valset.HashToSeed(mixHash)
	shuffled := qualified.ShuffledList(seed)
	if len(shuffled) < int(committeeSize) {
		return shuffled
	} else {
		return shuffled[:committeeSize]
	}
}

func (v *ValsetModule) getProposer(c *blockContext, round uint64) (common.Address, error) {
	switch istanbul.ProposerPolicy(c.pset.ProposerPolicy) {
	case istanbul.RoundRobin:
		return selectRoundRobinProposer(c.qualified, c.prevProposer, round), nil
	case istanbul.Sticky:
		return selectStickyProposer(c.qualified, c.prevProposer, round), nil
	case istanbul.WeightedRandom:
		if c.rules.IsRandao {
			committee := selectRandaoCommittee(c.qualified, c.pset.CommitteeSize, c.prevHeader.MixHash)
			return selectRandaoProposer(committee, round), nil
		} else {
			list, sourceNum, err := v.getProposerList(c)
			if err != nil {
				return common.Address{}, err
			}
			return selectWeightedRandomProposer(list, sourceNum, c.num, round), nil
		}
	default:
		return common.Address{}, errInvalidProposerPolicy
	}
}

func selectRoundRobinProposer(qualified *valset.AddressSet, prevProposer common.Address, round uint64) common.Address {
	prevIdx := qualified.IndexOf(prevProposer)
	if prevIdx < 0 {
		prevIdx = 0
	}
	return qualified.At((prevIdx + int(round) + 1) % qualified.Len())
}

func selectStickyProposer(qualified *valset.AddressSet, prevProposer common.Address, round uint64) common.Address {
	prevIdx := qualified.IndexOf(prevProposer)
	if prevIdx < 0 {
		prevIdx = 0
	}
	return qualified.At((prevIdx + int(round)) % qualified.Len())
}

// listSourceNum is the block number at which the list is generated. The caller must ensure that `len(list) > 0` and `listSourceNum <= num - 1`
func selectWeightedRandomProposer(list []common.Address, listSourceNum, num uint64, round uint64) common.Address {
	idx := num + round - listSourceNum - 1
	return list[int(idx)%len(list)]
}

func selectRandaoProposer(committee []common.Address, round uint64) common.Address {
	return committee[round%uint64(len(committee))]
}
