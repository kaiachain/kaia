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

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/log"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var handlerTestBids = make([]*auction.Bid, 5)

func init() {
	// Initialize bids for each searcher
	for i, key := range []*ecdsa.PrivateKey{testSearcher1Key, testSearcher2Key, testSearcher3Key, testSearcher1Key, testSearcher1Key} {
		bid := &auction.Bid{}
		initBaseBid(bid, i)

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
		mockCtrl = gomock.NewController(t)
		chain    = chain_mock.NewMockBlockChain(mockCtrl)
		block    = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(10)})
	)
	defer mockCtrl.Finish()

	chain.EXPECT().CurrentBlock().Return(block).AnyTimes()

	auctionConfig := auction.AuctionConfig{
		Disable: false,
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
		Chain:         chain,
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
	module.HandleBid(handlerTestBids[0])
	time.Sleep(10 * time.Millisecond)

	// Verify the bid was added to the pool
	assert.Equal(t, 1, len(module.bidPool.bidMap))

	// Test handling multiple bids
	for _, bid := range handlerTestBids[1:3] {
		module.HandleBid(bid)
	}
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 3, len(module.bidPool.bidMap))

	// Test handling a nil bid
	module.HandleBid(nil)
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 3, len(module.bidPool.bidMap)) // Stats should not change
}

func TestHandler_SubscribeNewBid(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		mockCtrl = gomock.NewController(t)
		chain    = chain_mock.NewMockBlockChain(mockCtrl)
		block    = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(10)})
	)
	defer mockCtrl.Finish()

	chain.EXPECT().CurrentBlock().Return(block).AnyTimes()

	auctionConfig := auction.AuctionConfig{
		Disable: false,
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
		Chain:         chain,
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
		module.HandleBid(handlerTestBids[0])
		module.HandleBid(handlerTestBids[1])
		module.HandleBid(handlerTestBids[2])
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
	module.HandleBid(handlerTestBids[3])

	// Verify we don't receive the bid after unsubscribing
	select {
	case bid := <-newBidCh:
		t.Fatalf("Received unexpected bid: %v", bid.Hash())
	case <-time.After(100 * time.Millisecond):
		// Expected timeout
	}
}
