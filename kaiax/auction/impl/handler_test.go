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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var handlerTestBids = make([]*auction.Bid, 5)

func init() {
	// Initialize bids for each searcher
	for i, key := range []*ecdsa.PrivateKey{testSearcher1Key, testSearcher2Key, testSearcher3Key, testSearcher1Key, testSearcher1Key} {
		bid := &auction.Bid{}
		initBaseBid(bid, i, 2)

		// Set searcher address
		bid.Sender = crypto.PubkeyToAddress(key.PublicKey)

		// Generate searcher signature (EIP-712)
		digest := bid.GetHashTypedData(testChainConfig.ChainID, testAuctionEntryPoint)
		sig, _ := crypto.Sign(digest, key)
		// Convert V from 0/1 to 27/28
		sig[crypto.RecoveryIDOffset] += 27
		bid.SearcherSig = sig

		// Generate auctioneer signature
		searcherSig := bid.SearcherSig
		msg := fmt.Appendf(nil, "\x19Ethereum Signed Message:\n%d%s", len(searcherSig), searcherSig)
		digest = crypto.Keccak256(msg)
		sig, _ = crypto.Sign(digest, testAuctioneerKey)
		// Convert V from 0/1 to 27/28
		sig[crypto.RecoveryIDOffset] += 27
		bid.AuctioneerSig = sig

		handlerTestBids[i] = bid
	}
}

func TestHandler_HandleBid(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db         = database.NewMemoryDBManager()
		alloc      = testAllocStorage()
		config     = testRandaoForkChainConfig(big.NewInt(0))
		testPeerID = "testPeerID"
	)
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	auctionConfig := auction.AuctionConfig{
		Disable:        false,
		MaxBidPoolSize: 1024,
	}
	apiBackend := &MockBackend{}
	fakeDownloader := &downloader.FakeDownloader{}

	// Create a new auction module
	module := NewAuctionModule()
	require.NotNil(t, module)

	// Initialize the module with test configuration
	opts := &InitOpts{
		ChainConfig:   testChainConfig,
		AuctionConfig: &auctionConfig,
		Chain:         backend.BlockChain(),
		Backend:       apiBackend,
		Downloader:    fakeDownloader,
		NodeKey:       testNodeKey,
	}
	err := module.Init(opts)
	require.NoError(t, err)

	// Start the module
	err = module.Start()
	atomic.StoreUint32(&module.bidPool.running, 1)
	require.NoError(t, err)
	defer module.Stop()

	module.bidPool.auctioneer = testAuctioneer
	module.bidPool.auctionEntryPoint = testAuctionEntryPoint

	// Test handling a valid bid
	module.HandleBid(testPeerID, handlerTestBids[0])
	time.Sleep(10 * time.Millisecond)

	// Verify the bid was added to the pool
	assert.Equal(t, 1, len(module.bidPool.bidMap))

	// Test handling multiple bids
	for _, bid := range handlerTestBids[1:3] {
		module.HandleBid(testPeerID, bid)
	}
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 3, len(module.bidPool.bidMap))

	// Test handling a nil bid
	module.HandleBid(testPeerID, nil)
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 3, len(module.bidPool.bidMap)) // Stats should not change
}

func TestHandler_RateLimit(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db          = database.NewMemoryDBManager()
		alloc       = testAllocStorage()
		config      = testRandaoForkChainConfig(big.NewInt(0))
		testPeerID1 = "testPeer1"
		testPeerID2 = "testPeer2"
	)
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	auctionConfig := auction.AuctionConfig{
		Disable:        false,
		MaxBidPoolSize: 10000,
	}
	apiBackend := &MockBackend{}
	fakeDownloader := &downloader.FakeDownloader{}

	// Create a new auction module
	module := NewAuctionModule()
	require.NotNil(t, module)

	// Initialize the module with test configuration
	opts := &InitOpts{
		ChainConfig:   testChainConfig,
		AuctionConfig: &auctionConfig,
		Chain:         backend.BlockChain(),
		Backend:       apiBackend,
		Downloader:    fakeDownloader,
		NodeKey:       testNodeKey,
	}
	err := module.Init(opts)
	require.NoError(t, err)

	// Start the module
	err = module.Start()
	atomic.StoreUint32(&module.bidPool.running, 1)
	require.NoError(t, err)
	defer module.Stop()

	module.bidPool.auctioneer = testAuctioneer
	module.bidPool.auctionEntryPoint = testAuctionEntryPoint

	// Commit a block to set current block number
	backend.Commit()
	currentBlock := backend.BlockChain().CurrentBlock().NumberU64()

	// Create many test bids with different target transactions to avoid duplicate rejection
	testBids := make([]*auction.Bid, 400)
	for i := 0; i < 400; i++ {
		key, _ := crypto.GenerateKey()
		bid := &auction.Bid{}

		initBaseBid(bid, i, currentBlock+1)
		bid.Sender = crypto.PubkeyToAddress(key.PublicKey)

		digest := bid.GetHashTypedData(testChainConfig.ChainID, testAuctionEntryPoint)
		sig, _ := crypto.Sign(digest, key)
		sig[crypto.RecoveryIDOffset] += 27
		bid.SearcherSig = sig

		searcherSig := bid.SearcherSig
		msg := fmt.Appendf(nil, "\x19Ethereum Signed Message:\n%d%s", len(searcherSig), searcherSig)
		digest = crypto.Keccak256(msg)
		sig, _ = crypto.Sign(digest, testAuctioneerKey)
		sig[crypto.RecoveryIDOffset] += 27
		bid.AuctioneerSig = sig

		testBids[i] = bid
	}

	// Test rate limiting for peer1
	// Send 350 bids rapidly (should be limited to 300)
	for i := 0; i < 350; i++ {
		module.HandleBid(testPeerID1, testBids[i])
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Should have accepted around 300 bids (rate limit)
	bidCount1 := len(module.bidPool.bidMap)
	assert.Equal(t, bidCount1, 300)

	// Test that peer2 is not affected by peer1's rate limit
	// Clear the bid pool first
	module.bidPool.clearBidPool()

	// Send bids from peer2
	for i := 0; i < 50; i++ {
		module.HandleBid(testPeerID2, testBids[i])
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Peer2 should have all its bids accepted
	bidCount2 := len(module.bidPool.bidMap)
	assert.Equal(t, 50, bidCount2)

	// Test rate limit recovery after 1 second
	module.bidPool.clearBidPool()
	time.Sleep(1100 * time.Millisecond) // Wait for rate limit to reset

	// Send another batch from peer1
	for i := 0; i < 100; i++ {
		module.HandleBid(testPeerID1, testBids[i])
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Should accept new bids after rate limit reset
	bidCount3 := len(module.bidPool.bidMap)
	assert.Equal(t, 100, bidCount3)
}

func TestHandler_SubscribeNewBid(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db         = database.NewMemoryDBManager()
		alloc      = testAllocStorage()
		config     = testRandaoForkChainConfig(big.NewInt(0))
		testPeerID = "testPeerID"
	)
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	auctionConfig := auction.AuctionConfig{
		Disable:        false,
		MaxBidPoolSize: 1024,
	}
	apiBackend := &MockBackend{}
	fakeDownloader := &downloader.FakeDownloader{}

	// Create a new auction module
	module := NewAuctionModule()
	require.NotNil(t, module)

	// Initialize the module with test configuration
	opts := &InitOpts{
		ChainConfig:   testChainConfig,
		AuctionConfig: &auctionConfig,
		Chain:         backend.BlockChain(),
		Backend:       apiBackend,
		Downloader:    fakeDownloader,
		NodeKey:       testNodeKey,
	}
	err := module.Init(opts)
	require.NoError(t, err)

	// Start the module
	err = module.Start()
	atomic.StoreUint32(&module.bidPool.running, 1)
	require.NoError(t, err)
	defer module.Stop()

	module.bidPool.auctioneer = testAuctioneer
	module.bidPool.auctionEntryPoint = testAuctionEntryPoint

	// Create a channel to receive new bids
	newBidCh := make(chan *auction.Bid, 10)

	// Subscribe to new bids
	sub := module.SubscribeNewBid(newBidCh)
	require.NotNil(t, sub)

	// Test receiving new bids
	go func() {
		module.HandleBid(testPeerID, handlerTestBids[0])
		module.HandleBid(testPeerID, handlerTestBids[1])
		module.HandleBid(testPeerID, handlerTestBids[2])
	}()

	// Verify we receive the bids in order
	for i := 0; i < 3; i++ {
		select {
		case bid := <-newBidCh:
			assert.Equal(t, handlerTestBids[i].Hash(), bid.Hash())
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for bid")
		}
	}

	// Test unsubscribing
	sub.Unsubscribe()

	// Send another bid
	module.HandleBid(testPeerID, handlerTestBids[3])

	// Verify we don't receive the bid after unsubscribing
	select {
	case bid := <-newBidCh:
		t.Fatalf("Received unexpected bid: %v", bid.Hash())
	case <-time.After(100 * time.Millisecond):
		// Expected timeout
	}
}
