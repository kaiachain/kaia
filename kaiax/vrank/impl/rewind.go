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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (v *VRankModule) RewindTo(newBlock *types.Block) {
	v.pfsCache.Purge()
	v.cpMatrixCache.Purge()
}

// RewindDelete removes scores for all blocks strictly above num.
// num must be the highest block that should survive (the new chain head after the rewind),
// NOT the block being deleted. The checkpoint at num itself is preserved.
func (v *VRankModule) RewindDelete(hash common.Hash, num uint64) {
	v.RemoveScoresAfter(num)
}

func (v *VRankModule) RemoveScoresAfter(blockNum uint64) {
	v.pfsCache.Purge()
	v.cpMatrixCache.Purge()

	head := v.Chain.CurrentHeader()
	if head == nil || head.Number == nil {
		return
	}
	for cpNum := calcCheckpointBlock(head.Number.Uint64()); ; {
		if cpNum <= blockNum {
			break
		}
		DeleteCheckpoint(v.ChainKv, cpNum)
		if cpNum < scoreCheckpointInterval {
			break
		}
		cpNum -= scoreCheckpointInterval
	}

	// Update lastCheckpoint pointer to the highest surviving checkpoint.
	newLastCP := calcCheckpointBlock(blockNum)
	pfs, _ := ReadCheckpoint(v.ChainKv, newLastCP)
	if pfs != nil {
		WriteLastCheckpoint(v.ChainKv, newLastCP)
	} else {
		DeleteLastCheckpoint(v.ChainKv)
	}
}
