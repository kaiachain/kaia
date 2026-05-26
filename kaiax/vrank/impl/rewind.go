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

// RewindDelete removes score state for exactly one deleted block `num`.
// The blockchain host calls this once per deleted block during hard rewind.
func (v *VRankModule) RewindDelete(hash common.Hash, num uint64) {
	pfs := ReadCheckpointPFS(v.ChainKv, num)
	if pfs == nil {
		return
	}
	DeleteCheckpoint(v.ChainKv, num)

	lastCP, ok := ReadLastCheckpoint(v.ChainKv)
	if !ok || lastCP != num {
		return
	}

	cpInterval := v.scoreCheckpointInterval()
	for cpNum := num; cpNum >= cpInterval; cpNum -= cpInterval {
		prevCP := cpNum - cpInterval
		prevPFS := ReadCheckpointPFS(v.ChainKv, prevCP)
		if prevPFS != nil {
			WriteLastCheckpoint(v.ChainKv, prevCP)
			return
		}
	}
	DeleteLastCheckpoint(v.ChainKv)
}
