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
)

func (v *VRankModule) PostInsertBlock(block *types.Block) error {
	blockNum := block.NumberU64()
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil
	}

	pfs, err := v.GetPFS(blockNum)
	if err != nil {
		return err
	}
	cpMatrix, err := v.getCPMatrix(blockNum)
	if err != nil {
		return err
	}

	if blockNum%v.scoreCheckpointInterval() == 0 {
		WriteCheckpoint(v.ChainKv, blockNum, pfs, cpMatrix)
		WriteLastCheckpoint(v.ChainKv, blockNum)
	}
	return nil
}
