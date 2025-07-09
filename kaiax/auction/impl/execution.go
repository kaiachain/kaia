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
	"sync/atomic"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (a *AuctionModule) PostInsertBlock(block *types.Block) error {
	if !a.ChainConfig.IsRandaoForkEnabled(block.Number()) || a.Downloader.Synchronising() {
		atomic.CompareAndSwapUint32(&a.bidPool.running, 1, 0)
		return nil
	}

	if err := a.updateAuctionInfo(block.Number()); err != nil {
		logger.Error("failed to update auction info", "error", err)
		a.bidPool.clearBidPool()
		atomic.CompareAndSwapUint32(&a.bidPool.running, 1, 0)
	} else {
		atomic.CompareAndSwapUint32(&a.bidPool.running, 0, 1)
	}

	if atomic.LoadUint32(&a.bidPool.running) == 1 {
		a.bidPool.removeOldBids(block.Number().Uint64())
	}

	return nil
}

func (a *AuctionModule) RewindTo(newBlock *types.Block) {
	// Nothing to do.
}

func (a *AuctionModule) RewindDelete(hash common.Hash, num uint64) {
	// Nothing to do.
}
