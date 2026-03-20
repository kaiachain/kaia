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
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (v *VRankModule) PostInsertBlock(block *types.Block) error {
	blockNum := block.NumberU64()
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil
	}

	// GetPFS/GetCFS populate both caches; capture pfs directly to avoid a second cache lookup.
	pfs, err := v.GetPFS(blockNum)
	if err != nil {
		return err
	}
	if _, err := v.GetCFS(blockNum); err != nil {
		return err
	}

	if blockNum%scoreCheckpointInterval == 0 {
		// GetCFS populates cpMatrixCache (not a separate cfsCache), so fetch from there.
		cpRaw, hasCPS := v.cpMatrixCache.Get(blockNum)
		var cpMatrix map[common.Address]map[common.Address]uint64
		if hasCPS {
			cpMatrix = cloneCPMatrix(cpRaw.(map[common.Address]map[common.Address]uint64))
		} else {
			// unlikely to happen: perform exhaustive re-calculation.
			cpMatrix, err = v.newCPMatrix(blockNum)
			if err != nil {
				return err
			}
			cpMatrix, err = v.computeCPMatrix(calcEpochStart(blockNum), blockNum, cpMatrix)
			if err != nil {
				return err
			}
		}
		WriteCheckpoint(v.ChainKv, blockNum, pfs, cpMatrix)
		WriteLastCheckpoint(v.ChainKv, blockNum)
	}
	return nil
}
