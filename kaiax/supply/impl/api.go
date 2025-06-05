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

package supply

import (
	"github.com/kaiachain/kaia/v2/kaiax/supply"
	"github.com/kaiachain/kaia/v2/networks/rpc"
)

func (s *SupplyModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewSupplyAPI(s),
			Public:    true,
		},
	}
}

type SupplyAPI struct {
	s *SupplyModule
}

func NewSupplyAPI(s *SupplyModule) *SupplyAPI {
	return &SupplyAPI{s: s}
}

// If showPartial == nil or *showPartial == false, the regular use case, this API either delivers the full result or fails.
// If showPartial == true, the advanced and debugging use case, this API delivers full or best effort partial result.
func (api *SupplyAPI) GetTotalSupply(num *rpc.BlockNumber, showPartial *bool) (*supply.TotalSupplyResponse, error) {
	var blockNum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNum = api.s.Chain.CurrentBlock().NumberU64()
	} else {
		blockNum = num.Uint64()
	}
	shouldShowPartial := showPartial != nil && *showPartial

	// 1. GetTotalSupply returns (nil, err) -> API returns (nil, err)
	// 2. GetTotalSupply returns (ts, err)
	// 2.1. if showPartial -> API returns (ts, nil)
	// 2.2. if !showPartial -> API returns (nil, err)
	// 3. GetTotalSupply returns (ts, nil) -> API returns (ts, nil)
	ts, err := api.s.GetTotalSupply(blockNum)
	if ts == nil {
		// 1. Essential components are missing.
		return nil, err
	}

	response := ts.ToResponse(blockNum, err)
	if err != nil {
		if shouldShowPartial {
			// 2.1. Deliver partial result
			return response, nil
		} else {
			// 2.2. Fail on partial (incomplete) result
			return nil, err
		}
	}

	// 3. Deliver full result
	return response, nil
}
