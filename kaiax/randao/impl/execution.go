// Copyright 2025 The Kaia Authors
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

	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
)

func (r *RandaoModule) PostInsertBlock(block *types.Block) error {
	isRandao := r.ChainConfig.IsRandaoForkEnabled(block.Number())
	if !isRandao || r.Downloader.Synchronising() {
		return nil
	}
	go func() {
		// Cache the bls pubkey for the next block.
		num := new(big.Int).Add(block.Number(), common.Big1)
		_, err := r.getAllCached(num)
		if err != nil {
			logger.Error("Failed to cache bls pubkey", "number", num.Uint64(), "error", err)
		}
	}()
	return nil
}

func (r *RandaoModule) RewindTo(newBlock *types.Block) {
	// Purge the bls pubkey cache.
	r.blsPubkeyCache.Purge()
}

func (r *RandaoModule) RewindDelete(hash common.Hash, num uint64) {
	// Nothing to do.
}
