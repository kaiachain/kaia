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
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/system"
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

	// Remove old bids unconditionally.
	a.bidPool.removeOldBids(block.Number().Uint64())

	return nil
}

func (a *AuctionModule) RewindTo(newBlock *types.Block) {
	// Nothing to do.
}

func (a *AuctionModule) RewindDelete(hash common.Hash, num uint64) {
	// Nothing to do.
}

// updateAuctionInfo updates the auctioneer address and auction entry point address for the given block number.
// It expects the `num` is after Randao fork.
func (a *AuctionModule) updateAuctionInfo(num *big.Int) error {
	header := a.Chain.GetHeaderByNumber(num.Uint64())
	if header == nil {
		return fmt.Errorf("failed to get header for block number %d", num.Uint64())
	}
	_, err := a.Chain.StateAt(header.Root)
	if err != nil {
		return fmt.Errorf("failed to get state for block number %d: %v", num.Uint64(), err)
	}

	backend := backends.NewBlockchainContractBackend(a.Chain, nil, nil)

	auctionEntryPointAddr, err := system.ReadActiveAddressFromRegistry(backend, system.AuctionEntryPointName, num)
	if err != nil {
		return fmt.Errorf("failed to read auction entry point address: %v", err)
	}

	auctioneer, err := system.ReadAuctioneer(backend, auctionEntryPointAddr, num)
	if err != nil {
		return fmt.Errorf("failed to read auctioneer address: %v", err)
	}

	a.bidPool.updateAuctionInfo(auctioneer, auctionEntryPointAddr)

	return nil
}
