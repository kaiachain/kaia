// Copyright 2026 The Kaia Authors
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
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
)

func (v *VRankModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   &VRankGovernanceAPI{v: v},
			Public:    true,
		},
	}
}

// VRankGovernanceAPI defines the governance namespace APIs for VRank scores.
type VRankGovernanceAPI struct {
	v *VRankModule
}

// GetPFS returns the Proposal Failure Scores for all validators at the given block number.
func (api *VRankGovernanceAPI) GetPFS(num *rpc.BlockNumber) (map[common.Address]uint64, error) {
	blockNum := api.resolveBlockNum(num)
	return api.v.GetPFS(blockNum)
}

// GetCFS returns the Consensus Failure Scores for all validators at the given block number.
func (api *VRankGovernanceAPI) GetCFS(num *rpc.BlockNumber) (map[common.Address]uint64, error) {
	blockNum := api.resolveBlockNum(num)
	return api.v.GetCFS(blockNum)
}

func (api *VRankGovernanceAPI) resolveBlockNum(num *rpc.BlockNumber) uint64 {
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		return api.v.Chain.CurrentHeader().Number.Uint64()
	}
	return num.Uint64()
}
