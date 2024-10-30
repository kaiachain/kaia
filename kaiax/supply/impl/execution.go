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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

// PostInsertBlock will try to advance the supplyCheckpoint by one block.
// If the new block is right after the last checkpoint, it will advance the supplyCheckpoint.
// Otherwise, it will leave the work to the catchup thread.
func (s *SupplyModule) PostInsertBlock(block *types.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newNum := block.NumberU64()
	if s.lastNum+1 == newNum && s.lastCheckpoint != nil {
		newCheckpoint, err := s.accumulateCheckpoint(s.lastNum, newNum, s.lastCheckpoint, true)
		if err != nil {
			return err
		}
		s.lastNum = newNum
		s.lastCheckpoint = newCheckpoint
	}
	return nil
}

func (s *SupplyModule) RewindTo(newBlock *types.Block) {
	// Soft reset to the nearest checkpoint interval less than the new block number,
	// so that the next accumulation will start below the rewound block.
	newLastNum := nearestCheckpointInterval(newBlock.NumberU64())
	WriteLastSupplyCheckpointNumber(s.ChainKv, newLastNum)
}

func (s *SupplyModule) RewindDelete(hash common.Hash, num uint64) {
	// TODO
}
