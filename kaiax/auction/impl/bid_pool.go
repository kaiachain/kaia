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

func (bp *BidPool) start() {
	atomic.CompareAndSwapUint32(&bp.running, 0, 1)
	bp.wg.Add(1)
	go bp.handleBidMsg()
}

func (bp *BidPool) stop() {
	atomic.CompareAndSwapUint32(&bp.running, 1, 0)
	bp.clearBidPool()
	close(bp.bidCh)
	bp.wg.Wait()
}

// removeOldBids removes the old bids for the given block number.
func (bp *BidPool) removeOldBids(num uint64) {
	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	// Clear the bid for the given block number.
	delete(bp.bidTargetMap, num)
	delete(bp.bidWinnerMap, num)
	for _, bid := range bp.allBids[num] {
		delete(bp.bidMap, bid.Hash())
	}
	delete(bp.allBids, num)
}

// clearBidPool clears the bid pool.
func (bp *BidPool) clearBidPool() {
	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	bp.bidTargetMap = make(map[uint64]map[common.Hash]*auction.Bid)
	bp.bidWinnerMap = make(map[uint64]map[common.Address]common.Hash)
	bp.bidMap = make(map[common.Hash]*auction.Bid)
	bp.allBids = make(map[uint64][]*auction.Bid)
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

// AddBid adds a bid to the bid pool.
func (bp *BidPool) AddBid(bid *auction.Bid) (common.Hash, error) {
	if atomic.LoadUint32(&bp.running) == 0 {
		return common.Hash{}, auction.ErrAuctionPaused
	}

	bp.bidMu.Lock()
	defer bp.bidMu.Unlock()

	if err := bp.validateBid(bid); err != nil {
		return common.Hash{}, err
	}

	if err := bp.insertBid(bid); err != nil {
		return common.Hash{}, err
	}

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
	bp.allBids[blockNumber] = append(bp.allBids[blockNumber], bid)

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
	if _, ok := bp.allBids[num]; !ok {
		bp.allBids[num] = make([]*auction.Bid, 0)
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
	if _, ok := bp.bidWinnerMap[blockNumber][sender]; ok {
		if !isSameBid(bid, bp.bidMap[bp.bidWinnerMap[blockNumber][sender]]) {
			return auction.ErrBidSenderExists
		}
	}

	// Check if the bid is valid.
	if err := bp.validateBidSigs(bid); err != nil {
		return err
	}

	return nil
}

func (bp *BidPool) stats() int {
	bp.bidMu.RLock()
	defer bp.bidMu.RUnlock()

	return len(bp.bidMap)
}

func (bp *BidPool) validateBidSigs(bid *auction.Bid) error {
	bp.auctionInfoMu.RLock()
	defer bp.auctionInfoMu.RUnlock()

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
	if atomic.LoadUint32(&bp.running) == 0 {
		return
	}
	bp.bidCh <- bid
}

func (bp *BidPool) handleBidMsg() {
	defer bp.wg.Done()

	for {
		select {
		case bid, ok := <-bp.bidCh:
			if !ok {
				return
			}
			bp.AddBid(bid)
		}
	}
}

func isSameBid(bid1 *auction.Bid, bid2 *auction.Bid) bool {
	return bid1.BlockNumber == bid2.BlockNumber && bid1.TargetTxHash == bid2.TargetTxHash
}
