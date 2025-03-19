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

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/params"
)

const (
	bidChSize        = 100
	allowFutureBlock = 2
)

type BidPool struct {
	ChainConfig *params.ChainConfig
	Chain       backends.BlockChainForCaller

	auctionInfoMu     sync.RWMutex
	auctioneer        common.Address
	auctionEntryPoint common.Address

	bidMu        sync.RWMutex
	bidMap       map[common.Hash]*auction.Bid
	bidTargetMap map[uint64]map[common.Hash]*auction.Bid
	bidWinnerMap map[uint64]map[common.Address]common.Hash
	allBids      map[uint64][]*auction.Bid

	bidCh chan *auction.Bid
	wg    sync.WaitGroup

	running uint32
}

func NewBidPool(chainConfig *params.ChainConfig, chain backends.BlockChainForCaller) *BidPool {
	if chainConfig == nil || chain == nil {
		return nil
	}

	bp := &BidPool{
		ChainConfig:  chainConfig,
		Chain:        chain,
		bidTargetMap: make(map[uint64]map[common.Hash]*auction.Bid),
		bidWinnerMap: make(map[uint64]map[common.Address]common.Hash),
		bidMap:       make(map[common.Hash]*auction.Bid),
		allBids:      make(map[uint64][]*auction.Bid),
		bidCh:        make(chan *auction.Bid, bidChSize),
		running:      0, // not running yet
	}

	return bp
}

func (a *BidPool) start() {
	atomic.CompareAndSwapUint32(&a.running, 0, 1)
	a.wg.Add(1)
	go a.handleBidMsg()
}

func (a *BidPool) stop() {
	atomic.CompareAndSwapUint32(&a.running, 1, 0)
	a.clearBidPool()
	close(a.bidCh)
	a.wg.Wait()
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
	a.bidWinnerMap = make(map[uint64]map[common.Address]common.Hash)
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
		// FCFS if the bid is the same.
		if existingTx.Bid.Cmp(bid.Bid) >= 0 {
			return auction.ErrLowBid
		}

		logger.Trace("Replace bid", "old", existingTx.Hash(), "new", bid.Hash())
		delete(a.bidMap, existingTx.Hash())
		delete(a.bidWinnerMap[blockNumber], existingTx.Sender)
	}

	hash := bid.Hash()

	a.initializeBidMap(blockNumber)

	a.bidMap[hash] = bid
	a.bidTargetMap[blockNumber][targetTxHash] = bid
	a.bidWinnerMap[blockNumber][sender] = hash
	a.allBids[blockNumber] = append(a.allBids[blockNumber], bid)

	logger.Trace("Add bid", "bid", hash)

	return nil
}

func (a *BidPool) initializeBidMap(num uint64) {
	if _, ok := a.bidTargetMap[num]; !ok {
		a.bidTargetMap[num] = make(map[common.Hash]*auction.Bid)
	}
	if _, ok := a.bidWinnerMap[num]; !ok {
		a.bidWinnerMap[num] = make(map[common.Address]common.Hash)
	}
	if _, ok := a.allBids[num]; !ok {
		a.allBids[num] = make([]*auction.Bid, 0)
	}
}

func (a *BidPool) validateBid(bid *auction.Bid) error {
	blockNumber := bid.BlockNumber
	sender := bid.Sender

	curBlock := a.Chain.CurrentBlock()
	if curBlock == nil {
		return auction.ErrBlockNotFound
	}

	curNum := curBlock.NumberU64()
	if blockNumber <= curNum || blockNumber > curNum+allowFutureBlock {
		return auction.ErrInvalidBlockNumber
	}

	if bid.Bid.Sign() <= 0 {
		return auction.ErrZeroBid
	}

	// Check if the auction tx is already in the pool.
	if _, ok := a.bidMap[bid.Hash()]; ok {
		return auction.ErrBidAlreadyExists
	}

	// Check if the sender is already in the winner list.
	if _, ok := a.bidWinnerMap[blockNumber][sender]; ok {
		if !isSameBid(bid, a.bidMap[a.bidWinnerMap[blockNumber][sender]]) {
			return auction.ErrBidSenderExists
		}
	}

	// Check if the bid is valid.
	if err := a.validateBidSigs(bid); err != nil {
		return err
	}

	return nil
}

func (a *BidPool) stats() int {
	a.bidMu.RLock()
	defer a.bidMu.RUnlock()

	return len(a.bidMap)
}

func (a *BidPool) validateBidSigs(bid *auction.Bid) error {
	a.auctionInfoMu.RLock()
	defer a.auctionInfoMu.RUnlock()

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

func (a *BidPool) HandleBid(bid *auction.Bid) {
	if atomic.LoadUint32(&a.running) == 0 {
		return
	}
	a.bidCh <- bid
}

func (a *BidPool) handleBidMsg() {
	defer a.wg.Done()

	for {
		select {
		case bid, ok := <-a.bidCh:
			if !ok {
				return
			}
			a.AddBid(bid)
		}
	}
}

func isSameBid(bid1 *auction.Bid, bid2 *auction.Bid) bool {
	return bid1.BlockNumber == bid2.BlockNumber && bid1.TargetTxHash == bid2.TargetTxHash
}
