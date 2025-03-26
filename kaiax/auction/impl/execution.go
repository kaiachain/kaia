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
	"sync/atomic"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (a *AuctionModule) PostInsertBlock(block *types.Block) error {
	if a.Downloader.Synchronising() || !a.ChainConfig.IsRandaoForkEnabled(block.Number()) {
		atomic.CompareAndSwapUint32(&a.bidPool.running, 1, 0)
		return nil
	}

	if !a.updateAuctionInfo(block.Number()) {
		logger.Debug("stop auction since auctioneer or auction entry point is not set")
		atomic.CompareAndSwapUint32(&a.bidPool.running, 1, 0)
		return nil
	}

	atomic.CompareAndSwapUint32(&a.bidPool.running, 0, 1)

	txHashMap := make(map[common.Hash]struct{})
	for _, tx := range block.Transactions() {
		txHashMap[tx.Hash()] = struct{}{}
	}
	a.bidPool.removeOldBids(block.Number().Uint64(), txHashMap)

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
// It returns true if the non-zero auctioneer address and auction entry point address are set, otherwise false.
func (a *AuctionModule) updateAuctionInfo(num *big.Int) bool {
	auctioneer := common.Address{}
	auctionEntryPointAddr := common.Address{}

	defer func() {
		a.bidPool.updateAuctionInfo(auctioneer, auctionEntryPointAddr)
	}()

	header := a.Chain.GetHeaderByNumber(num.Uint64())
	if header == nil {
		return false
	}
	_, err := a.Chain.StateAt(header.Root)
	if err != nil {
		return false
	}

	backend := backends.NewBlockchainContractBackend(a.Chain, nil, nil)

	auctionEntryPointAddr, err = system.ReadActiveAddressFromRegistry(backend, system.AuctionEntryPointName, num)
	if err != nil {
		return false
	}

	if auctionEntryPointAddr == (common.Address{}) {
		return false
	}

	auctioneer, err = system.ReadAuctioneer(backend, auctionEntryPointAddr, num)
	if err != nil {
		return false
	}

	if auctioneer == (common.Address{}) || auctionEntryPointAddr == (common.Address{}) {
		return false
	}

	return true
}
