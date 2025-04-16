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
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/params"
)

const (
	bidChSize        = 100
	allowFutureBlock = 2

	MaxBidNum = 2048
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

	bidMsgCh   chan *auction.Bid
	newBidCh   chan *auction.Bid
	newBidFeed event.Feed
	wg         sync.WaitGroup

	running uint32
	stopped uint32
}

func NewBidPool(chainConfig *params.ChainConfig, chain backends.BlockChainForCaller) *BidPool {
	if chainConfig == nil || chain == nil {
		return nil
	}

	bp := &BidPool{
		ChainConfig:  chainConfig,
		Chain:        chain,
		bidMap:       make(map[common.Hash]*auction.Bid),
		bidTargetMap: make(map[uint64]map[common.Hash]*auction.Bid),
		bidWinnerMap: make(map[uint64]map[common.Address]common.Hash),
		bidMsgCh:     make(chan *auction.Bid, bidChSize),
		newBidCh:     make(chan *auction.Bid, bidChSize),
		running:      0, // not running yet
		stopped:      0, // not stopped
	}

	return bp
}

func (bp *BidPool) SubscribeNewBid(sink chan<- *auction.Bid) event.Subscription {
	// Do not prevent subscription before start
	// if atomic.LoadUint32(&bp.running) == 0 {
	// 	return nil
	// }
	return bp.newBidFeed.Subscribe(sink)
}

func (bp *BidPool) start() {
	// Start the bid pool.
	// running will be set 1 once it's ready.

	// If channels are closed, recreate them
	if atomic.LoadUint32(&bp.stopped) == 1 {
		bp.bidMsgCh = make(chan *auction.Bid, bidChSize)
		bp.newBidCh = make(chan *auction.Bid, bidChSize)
		atomic.StoreUint32(&bp.stopped, 0)
	}

	bp.wg.Add(2)
	go bp.handleBidMsg()
	go bp.handleNewBid()
}

func (bp *BidPool) stop() {
	// Stop the bid pool.
	atomic.CompareAndSwapUint32(&bp.running, 1, 0)
	bp.clearBidPool()

	// Only close channels if they haven't been closed before
	if atomic.CompareAndSwapUint32(&bp.stopped, 0, 1) {
		close(bp.bidMsgCh)
		close(bp.newBidCh)
	}
	bp.wg.Wait()
}

// removeOldBids removes the old bids for the given block number.
func (bp *BidPool) removeOldBids(num uint64, txHashMap map[common.Hash]struct{}) {
	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	// Remove the old bids.
	for bn := range bp.bidWinnerMap {
		if bn > num {
			break
		}

		for _, bh := range bp.bidWinnerMap[bn] {
			delete(bp.bidMap, bh)
		}
		delete(bp.bidTargetMap, bn)
		delete(bp.bidWinnerMap, bn)
	}

	// Remove the bid which target tx is in the txHashMap.
	toBlock := num + allowFutureBlock
	for blockNum := num + 1; blockNum <= toBlock; blockNum++ {
		targetMap := bp.bidTargetMap[blockNum]
		if targetMap == nil {
			continue
		}

		// Collect bids to remove first to avoid modifying map during iteration
		var bidsToRemove []*auction.Bid
		for _, bid := range targetMap {
			if _, ok := txHashMap[bid.TargetTxHash]; ok {
				bidsToRemove = append(bidsToRemove, bid)
			}
		}

		// Remove collected bids
		for _, bid := range bidsToRemove {
			delete(targetMap, bid.TargetTxHash)
			delete(bp.bidWinnerMap[blockNum], bid.Sender)
			delete(bp.bidMap, bid.Hash())
		}
	}
}

// clearBidPool clears the bid pool.
func (bp *BidPool) clearBidPool() {
	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	bp.bidMap = make(map[common.Hash]*auction.Bid)
	bp.bidTargetMap = make(map[uint64]map[common.Hash]*auction.Bid)
	bp.bidWinnerMap = make(map[uint64]map[common.Address]common.Hash)
}

// updateAuctionInfo updates the auction info if the auctioneer or auction entry point address is changed.
func (bp *BidPool) updateAuctionInfo(auctioneer common.Address, auctionEntryPoint common.Address) {
	bp.auctionInfoMu.Lock()
	defer bp.auctionInfoMu.Unlock()

	if bp.auctioneer == auctioneer && bp.auctionEntryPoint == auctionEntryPoint {
		return
	}

	// Clear the existing auction pool since the auctioneer or auction entry point address is changed.
	bp.clearBidPool()

	bp.auctioneer = auctioneer
	bp.auctionEntryPoint = auctionEntryPoint

	logger.Info("Update auction info", "auctioneer", auctioneer, "auctionEntryPoint", auctionEntryPoint)
}

// getTargetTxHashMap returns the target tx hash map for the given block number.
func (bp *BidPool) getTargetTxHashMap(num uint64) map[common.Hash]struct{} {
	bp.bidMu.RLock()
	defer bp.bidMu.RUnlock()

	txHashMap := make(map[common.Hash]struct{})
	for hash := range bp.bidTargetMap[num] {
		txHashMap[hash] = struct{}{}
	}
	return txHashMap
}

func (bp *BidPool) GetAuctionEntryPoint() common.Address {
	bp.auctionInfoMu.RLock()
	defer bp.auctionInfoMu.RUnlock()

	return bp.auctionEntryPoint
}

func (bp *BidPool) GetTargetTxMap(num uint64) map[common.Hash]*auction.Bid {
	bp.bidMu.RLock()
	defer bp.bidMu.RUnlock()

	return bp.bidTargetMap[num]
}

// AddBid adds a bid to the bid pool.
func (bp *BidPool) AddBid(bid *auction.Bid) (common.Hash, error) {
	if atomic.LoadUint32(&bp.running) == 0 {
		return common.Hash{}, auction.ErrAuctionPaused
	}

	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	if len(bp.bidMap) >= MaxBidNum {
		return common.Hash{}, auction.ErrBidPoolFull
	}

	if err := bp.validateBid(bid); err != nil {
		return common.Hash{}, err
	}

	if err := bp.insertBid(bid); err != nil {
		return common.Hash{}, err
	}

	bp.newBidCh <- bid

	return bid.Hash(), nil
}

func (bp *BidPool) insertBid(bid *auction.Bid) error {
	var (
		blockNumber  = bid.BlockNumber
		targetTxHash = bid.TargetTxHash
		sender       = bid.Sender
	)

	// If same block number, same target tx hash exists, replace it if it's better
	if existingTx, ok := bp.bidTargetMap[blockNumber][targetTxHash]; ok {
		// FCFS if the bid is the same.
		if existingTx.Bid.Cmp(bid.Bid) >= 0 {
			return auction.ErrLowBid
		}

		logger.Trace("Replace bid", "old", existingTx.Hash(), "new", bid.Hash())
		delete(bp.bidMap, existingTx.Hash())
		delete(bp.bidWinnerMap[blockNumber], existingTx.Sender)
	}

	hash := bid.Hash()

	bp.initializeBidMap(blockNumber)

	bp.bidMap[hash] = bid
	bp.bidTargetMap[blockNumber][targetTxHash] = bid
	bp.bidWinnerMap[blockNumber][sender] = hash

	logger.Trace("Add bid", "bid", hash)

	return nil
}

func (bp *BidPool) initializeBidMap(num uint64) {
	if _, ok := bp.bidTargetMap[num]; !ok {
		bp.bidTargetMap[num] = make(map[common.Hash]*auction.Bid)
	}
	if _, ok := bp.bidWinnerMap[num]; !ok {
		bp.bidWinnerMap[num] = make(map[common.Address]common.Hash)
	}
}

func (bp *BidPool) validateBid(bid *auction.Bid) error {
	blockNumber := bid.BlockNumber
	sender := bid.Sender

	curBlock := bp.Chain.CurrentBlock()
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
	if _, ok := bp.bidMap[bid.Hash()]; ok {
		return auction.ErrBidAlreadyExists
	}

	// Check if the sender is already in the winner list.
	if hash, ok := bp.bidWinnerMap[blockNumber][sender]; ok {
		if !isSameBid(bid, bp.bidMap[hash]) {
			return auction.ErrBidSenderExists
		}
	}

	// Check if the bid is valid.
	if err := bp.validateBidSigs(bid); err != nil {
		return err
	}

	return nil
}

func (bp *BidPool) validateBidSigs(bid *auction.Bid) error {
	bp.auctionInfoMu.RLock()
	defer bp.auctionInfoMu.RUnlock()

	if bid.SearcherSig == nil || len(bid.SearcherSig) != crypto.SignatureLength {
		return auction.ErrInvalidSearcherSig
	}
	if bid.AuctioneerSig == nil || len(bid.AuctioneerSig) != crypto.SignatureLength {
		return auction.ErrInvalidAuctioneerSig
	}

	// Verify the EIP712 signature.
	if err := bid.ValidateSearcherSig(bp.ChainConfig.ChainID, bp.auctionEntryPoint); err != nil {
		return err
	}

	// Verify the auctioneer signature.
	if err := bid.ValidateAuctioneerSig(bp.auctioneer); err != nil {
		return err
	}

	return nil
}

func (bp *BidPool) HandleBid(bid *auction.Bid) {
	if atomic.LoadUint32(&bp.running) == 0 || bid == nil {
		return
	}
	bp.bidMsgCh <- bid
}

func (bp *BidPool) handleBidMsg() {
	defer bp.wg.Done()

	for {
		select {
		case bid, ok := <-bp.bidMsgCh:
			if !ok {
				return
			}
			bp.AddBid(bid)
		}
	}
}

func (bp *BidPool) handleNewBid() {
	defer bp.wg.Done()

	for {
		select {
		case bid, ok := <-bp.newBidCh:
			if !ok {
				return
			}
			bp.newBidFeed.Send(bid)
		}
	}
}

func isSameBid(bid1 *auction.Bid, bid2 *auction.Bid) bool {
	return bid1.BlockNumber == bid2.BlockNumber && bid1.TargetTxHash == bid2.TargetTxHash
}
