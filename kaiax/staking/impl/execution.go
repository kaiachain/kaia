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
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
)

func (s *StakingModule) PostInsertBlock(block *types.Block) error {
	isKaia := s.ChainConfig.IsKaiaForkEnabled(block.Number())
	if !isKaia {
		// Make sure the staking info for the new block is persisted.
		// The StakingInfo(sourceNum) will be persisted here, even if GetStakingInfo is never called elsewhere.
		if _, err := s.GetStakingInfo(block.NumberU64()); err != nil {
			return err
		}
	}
	return nil
}

func (s *StakingModule) RewindTo(newBlock *types.Block) {
	// Purge the staking info cache.
	s.stakingInfoCache.Purge()
}

func (s *StakingModule) RewindDelete(hash common.Hash, num uint64) {
	DeleteStakingInfo(s.ChainKv, num)
}
