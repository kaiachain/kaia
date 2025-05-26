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

	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/networks/rpc"
)

func (s *StakingModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   newStakingAPI(s),
			Public:    true,
		},
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   newStakingAPI(s),
			Public:    true,
		},
	}
}

type stakingAPI struct {
	s *StakingModule
}

func newStakingAPI(s *StakingModule) *stakingAPI {
	return &stakingAPI{s}
}

func (api *stakingAPI) GetStakingInfo(num rpc.BlockNumber) (*staking.StakingInfoResponse, error) {
	if num == rpc.LatestBlockNumber {
		num = rpc.BlockNumber(api.s.Chain.CurrentBlock().NumberU64())
	}

	si, err := api.s.GetStakingInfo(num.Uint64())
	if err != nil {
		return nil, err
	}

	var useGini bool
	// Gini option deprecated since Kore, as All committee members have an equal chance
	if api.s.ChainConfig.IsKoreForkEnabled(new(big.Int).SetUint64(num.Uint64())) {
		useGini = false
	} else {
		useGini = api.s.useGiniCoeff
	}
	// Calculate Gini coefficient regardless of useGini flag
	return si.ToResponse(useGini, api.s.minimumStake.Uint64()), nil
}
