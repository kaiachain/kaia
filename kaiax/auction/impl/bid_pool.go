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
	"sync"
	"sync/atomic"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/params"
)

type BidPool struct {
	ChainConfig *params.ChainConfig

	auctionInfoMu     sync.RWMutex
	auctioneer        common.Address
	auctionEntryPoint common.Address

	bidMu        sync.RWMutex
	bidTargetMap map[uint64]map[common.Hash]*auction.Bid
	bidWinnerMap map[uint64]map[common.Address]struct{}
	bidMap       map[common.Hash]*auction.Bid
	allBids      map[uint64][]*auction.Bid

	running uint32
}

func NewBidPool(chainConfig *params.ChainConfig) *BidPool {
	if chainConfig == nil {
		return nil
	}

	return &BidPool{
		ChainConfig:  chainConfig,
		bidTargetMap: make(map[uint64]map[common.Hash]*auction.Bid),
		bidWinnerMap: make(map[uint64]map[common.Address]struct{}),
		bidMap:       make(map[common.Hash]*auction.Bid),
		allBids:      make(map[uint64][]*auction.Bid),
		running:      0, // not running yet
	}
}

// removeOldBids removes the old bids for the given block number.
func (a *BidPool) removeOldBids(num uint64) {
	a.bidMu.Lock()
	defer a.bidMu.Unlock()

	// Clear the bid for the given block number.
	delete(a.bidTargetMap, num)
	delete(a.bidWinnerMap, num)
	for _, bid := range a.allBids[num] {
		delete(a.bidMap, bid.Hash())
	}
	delete(a.allBids, num)
}

// clearBidPool clears the bid pool.
func (a *BidPool) clearBidPool() {
	a.bidMu.Lock()
	defer a.bidMu.Unlock()

	a.bidTargetMap = make(map[uint64]map[common.Hash]*auction.Bid)
	a.bidWinnerMap = make(map[uint64]map[common.Address]struct{})
	a.bidMap = make(map[common.Hash]*auction.Bid)
	a.allBids = make(map[uint64][]*auction.Bid)
}

// updateAuctionInfo updates the auction info if the auctioneer or auction entry point address is changed.
func (a *BidPool) updateAuctionInfo(auctioneer common.Address, auctionEntryPoint common.Address) {
	a.auctionInfoMu.Lock()
	defer a.auctionInfoMu.Unlock()

	if a.auctioneer == auctioneer && a.auctionEntryPoint == auctionEntryPoint {
		return
	}

	// Clear the existing auction pool since the auctioneer or auction entry point address is changed.
	a.clearBidPool()

	a.auctioneer = auctioneer
	a.auctionEntryPoint = auctionEntryPoint

	logger.Info("Update auction info", "auctioneer", auctioneer, "auctionEntryPoint", auctionEntryPoint)
}

// AddBid adds a bid to the bid pool.
func (a *BidPool) AddBid(bid *auction.Bid) (common.Hash, error) {
	if atomic.LoadUint32(&a.running) == 0 {
		return common.Hash{}, auction.ErrAuctionPaused
	}

	a.bidMu.Lock()
	defer a.bidMu.Unlock()

	if err := a.validateBid(bid); err != nil {
		return common.Hash{}, err
	}

	if err := a.validateBidSigs(bid); err != nil {
		return common.Hash{}, err
	}

	if err := a.insertBid(bid); err != nil {
		return common.Hash{}, err
	}

	return bid.Hash(), nil
}

func (a *BidPool) insertBid(bid *auction.Bid) error {
	var (
		blockNumber  = bid.BlockNumber
		targetTxHash = bid.TargetTxHash
		sender       = bid.Sender
	)

	// If same block number, same target tx hash exists, replace it if it's better
	if existingTx, ok := a.bidTargetMap[blockNumber][targetTxHash]; ok {
		if existingTx.Bid.Cmp(bid.Bid) > 0 {
			return auction.ErrLowBid
		}

		logger.Trace("Replace bid", "old", existingTx.Hash(), "new", bid.Hash())
		delete(a.bidMap, existingTx.Hash())
		delete(a.bidWinnerMap[blockNumber], existingTx.Sender)
	}

	hash := bid.Hash()

	a.bidMap[hash] = bid
	a.bidTargetMap[blockNumber][targetTxHash] = bid
	a.bidWinnerMap[blockNumber][sender] = struct{}{}
	if _, ok := a.allBids[blockNumber]; !ok {
		a.allBids[blockNumber] = make([]*auction.Bid, 0)
	}
	a.allBids[blockNumber] = append(a.allBids[blockNumber], bid)

	logger.Trace("Add bid", "bid", hash)

	return nil
}

func (a *BidPool) validateBid(bid *auction.Bid) error {
	blockNumber := bid.BlockNumber
	sender := bid.Sender

	// Check if the auction tx is already in the pool.
	if _, ok := a.bidMap[bid.Hash()]; ok {
		return auction.ErrBidAlreadyExists
	}

	// Check if the sender is already in the winner list.
	if _, ok := a.bidWinnerMap[blockNumber][sender]; ok {
		return auction.ErrBidSenderExists
	}

	return nil
}

func (a *BidPool) validateBidSigs(bid *auction.Bid) error {
	// Verify the EIP712 signature.
	if err := bid.ValidateSearcherSig(a.ChainConfig.ChainID, a.auctionEntryPoint); err != nil {
		return err
	}

	// Verify the auctioneer signature.
	if err := bid.ValidateAuctioneerSig(a.auctioneer); err != nil {
		return err
	}

	return nil
}
