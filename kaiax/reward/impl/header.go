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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/reward"
)

func (r *RewardModule) VerifyHeader(header *types.Header, _ *types.Header) error {
	if header.Number.Sign() == 0 {
		return nil
	}
	// Only WeightedRandom has a staking-registered reward address. Under RoundRobin
	// the proposer's Rewardbase is unconstrained, matching FinalizeState's behavior.
	if r.GovModule.GetParamSet(header.Number.Uint64()).ProposerPolicy != uint64(istanbul.WeightedRandom) {
		return nil
	}
	proposer, err := r.Chain.Sealer().Author(header)
	if err != nil {
		return err
	}
	expected := r.GetRewardAddress(header.Number.Uint64(), proposer)
	// No staking entry for this proposer: nothing to check against. FinalizeState
	// falls back to the node's own configured Rewardbase in this case, so do not
	// reject the header here.
	if expected == (common.Address{}) {
		return nil
	}
	if header.Rewardbase != expected {
		return reward.ErrInvalidRewardbase
	}
	return nil
}

func (r *RewardModule) PrepareHeader(header *types.Header) error {
	header.Rewardbase = r.Rewardbase
	return nil
}
