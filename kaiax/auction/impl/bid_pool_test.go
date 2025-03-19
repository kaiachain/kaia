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
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testChainConfig = &params.ChainConfig{
		ChainID: big.NewInt(31337),
	}

	// Test private keys
	testAuctioneerKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	testSearcher1Key, _  = crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	testSearcher2Key, _  = crypto.HexToECDSA("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a")
	testSearcher3Key, _  = crypto.HexToECDSA("7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6")

	// Addresses derived from private keys
	testAuctioneer        = crypto.PubkeyToAddress(testAuctioneerKey.PublicKey)
	testAuctionEntryPoint = common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9")
	testSearcher1         = crypto.PubkeyToAddress(testSearcher1Key.PublicKey)
	testSearcher2         = crypto.PubkeyToAddress(testSearcher2Key.PublicKey)
	testSearcher3         = crypto.PubkeyToAddress(testSearcher3Key.PublicKey)

	testBids = make([]*auction.Bid, 5)
)

func init() {
	// Initialize bids for each searcher
	for i, key := range []*ecdsa.PrivateKey{testSearcher1Key, testSearcher2Key, testSearcher3Key, testSearcher1Key, testSearcher1Key} {
		bid := &auction.Bid{}
		initBaseBid(bid, i)

		if i == 3 {
			bid.TargetTxHash = common.HexToHash("0xf3c03c891206b24f5d2ff65b460df9b58c652279a3e0faed865dde4c46fe9da0")
		}

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

		testBids[i] = bid
	}
}

func initBaseBid(bid *auction.Bid, index int) {
	bid.TargetTxHash = common.HexToHash(fmt.Sprintf("0xf3c03c891206b24f5d2ff65b460df9b58c652279a3e0faed865dde4c46fe9da%d", index))
	bid.BlockNumber = 11
	bid.To = common.HexToAddress(fmt.Sprintf("0x5FC8d32690cc91D4c39d9d3abcBD16989F87570%d", index))
	bid.Nonce = uint64(index)
	bid.Bid = new(big.Int).Add(
		new(big.Int).SetBytes(common.Hex2Bytes("8ac7230489e80000")),
		new(big.Int).SetUint64(uint64(index)*1e18),
	)
	bid.CallGasLimit = 10000000
	bid.Data = common.Hex2Bytes("d09de08a")
}

func TestNewBidPool(t *testing.T) {
	testCases := []struct {
		name        string
		chainConfig *params.ChainConfig
		wantNil     bool
	}{
		{
			name:        "valid chain config",
			chainConfig: testChainConfig,
			wantNil:     false,
		},
		{
			name:        "nil chain config",
			chainConfig: nil,
			wantNil:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := NewBidPool(tc.chainConfig)
			if tc.wantNil {
				assert.Nil(t, pool)
			} else {
				assert.NotNil(t, pool)
				assert.Equal(t, tc.chainConfig, pool.ChainConfig)
				assert.NotNil(t, pool.bidTargetMap)
				assert.NotNil(t, pool.bidWinnerMap)
				assert.NotNil(t, pool.bidMap)
				assert.NotNil(t, pool.allBids)
				assert.Equal(t, uint32(0), pool.running)
			}
		})
	}
}

func TestBidPool_AddBid(t *testing.T) {
	pool := NewBidPool(testChainConfig)
	require.NotNil(t, pool)

	// Test adding bid when auction is paused
	_, err := pool.AddBid(testBids[0])
	assert.Equal(t, auction.ErrAuctionPaused, err)

	// Start the auction
	pool.running = 1
	pool.auctioneer = testAuctioneer
	pool.auctionEntryPoint = testAuctionEntryPoint

	// Test successful bid additions
	for _, bid := range testBids[:3] {
		hash, err := pool.AddBid(bid)
		require.NoError(t, err)
		assert.Equal(t, bid.Hash(), hash)

		// Verify bid was added correctly
		assert.Equal(t, bid, pool.bidMap[hash])
		assert.Equal(t, bid, pool.bidTargetMap[bid.BlockNumber][bid.TargetTxHash])
		assert.Contains(t, pool.bidWinnerMap[bid.BlockNumber], bid.Sender)
		assert.Contains(t, pool.allBids[bid.BlockNumber], bid)
	}

	// Test zero bid
	zeroBid := &auction.Bid{}
	*zeroBid = *testBids[0]
	zeroBid.Bid = big.NewInt(0)
	_, err = pool.AddBid(zeroBid)
	assert.Equal(t, auction.ErrZeroBid, err)

	// Test duplicate bid
	_, err = pool.AddBid(testBids[0])
	assert.Equal(t, auction.ErrBidAlreadyExists, err)

	// Test bid with same sender but different target
	duplicateSenderBid := &auction.Bid{}
	*duplicateSenderBid = *testBids[4]
	_, err = pool.AddBid(duplicateSenderBid)
	assert.Equal(t, auction.ErrBidSenderExists, err)

	// Test bid with higher amount for same target
	higherBid := &auction.Bid{}
	*higherBid = *testBids[3]
	_, err = pool.AddBid(higherBid)
	require.NoError(t, err)
	assert.Equal(t, higherBid, pool.bidTargetMap[testBids[3].BlockNumber][testBids[3].TargetTxHash])
	assert.Equal(t, higherBid.Hash(), pool.bidWinnerMap[testBids[3].BlockNumber][testBids[3].Sender])

	// Test bid with lower amount for same target
	lowerBid := &auction.Bid{}
	*lowerBid = *testBids[0]
	_, err = pool.AddBid(lowerBid)
	assert.Equal(t, auction.ErrLowBid, err)
}

func TestBidPool_RemoveOldBids(t *testing.T) {
	pool := NewBidPool(testChainConfig)
	require.NotNil(t, pool)

	// Start the auction
	pool.running = 1
	pool.auctioneer = testAuctioneer
	pool.auctionEntryPoint = testAuctionEntryPoint

	for _, bid := range testBids[:3] {
		hash, err := pool.AddBid(bid)
		require.NoError(t, err)
		assert.Equal(t, bid.Hash(), hash)

		// Verify bid was added correctly
		assert.Equal(t, bid, pool.bidMap[hash])
		assert.Equal(t, bid, pool.bidTargetMap[bid.BlockNumber][bid.TargetTxHash])
		assert.Contains(t, pool.bidWinnerMap[bid.BlockNumber], bid.Sender)
		assert.Contains(t, pool.allBids[bid.BlockNumber], bid)
	}

	// Remove bids for block 11
	pool.removeOldBids(11)

	// Verify bids for block 11 were removed
	assert.Empty(t, pool.bidTargetMap[11])
	assert.Empty(t, pool.bidWinnerMap[11])
	assert.Empty(t, pool.allBids[11])
}

func TestBidPool_ClearBidPool(t *testing.T) {
	pool := NewBidPool(testChainConfig)
	require.NotNil(t, pool)

	// Start the auction
	pool.running = 1
	pool.auctioneer = testAuctioneer
	pool.auctionEntryPoint = testAuctionEntryPoint

	// Add some bids
	for _, bid := range testBids[:3] {
		hash, err := pool.AddBid(bid)
		require.NoError(t, err)
		assert.Equal(t, bid.Hash(), hash)

		// Verify bid was added correctly
		assert.Equal(t, bid, pool.bidMap[hash])
		assert.Equal(t, bid, pool.bidTargetMap[bid.BlockNumber][bid.TargetTxHash])
		assert.Contains(t, pool.bidWinnerMap[bid.BlockNumber], bid.Sender)
		assert.Contains(t, pool.allBids[bid.BlockNumber], bid)
	}

	// Clear the pool
	pool.clearBidPool()

	// Verify all maps are empty
	assert.Empty(t, pool.bidTargetMap)
	assert.Empty(t, pool.bidWinnerMap)
	assert.Empty(t, pool.bidMap)
	assert.Empty(t, pool.allBids)
}

func TestBidPool_UpdateAuctionInfo(t *testing.T) {
	pool := NewBidPool(testChainConfig)
	require.NotNil(t, pool)

	// Start the auction
	pool.running = 1
	pool.auctioneer = testAuctioneer
	pool.auctionEntryPoint = testAuctionEntryPoint

	// Add some bids
	for _, bid := range testBids[:3] {
		hash, err := pool.AddBid(bid)
		require.NoError(t, err)
		assert.Equal(t, bid.Hash(), hash)

		// Verify bid was added correctly
		assert.Equal(t, bid, pool.bidMap[hash])
		assert.Equal(t, bid, pool.bidTargetMap[bid.BlockNumber][bid.TargetTxHash])
		assert.Contains(t, pool.bidWinnerMap[bid.BlockNumber], bid.Sender)
		assert.Contains(t, pool.allBids[bid.BlockNumber], bid)
	}

	// Update auction info with same addresses
	pool.updateAuctionInfo(testAuctioneer, testAuctionEntryPoint)
	assert.NotEmpty(t, pool.bidTargetMap)
	assert.NotEmpty(t, pool.bidWinnerMap)
	assert.NotEmpty(t, pool.bidMap)
	assert.NotEmpty(t, pool.allBids)

	// Update auction info with different addresses
	newAuctioneer := common.HexToAddress("0x1234")
	newAuctionEntryPoint := common.HexToAddress("0x5678")
	pool.updateAuctionInfo(newAuctioneer, newAuctionEntryPoint)

	// Verify pool was cleared and addresses were updated
	assert.Empty(t, pool.bidTargetMap)
	assert.Empty(t, pool.bidWinnerMap)
	assert.Empty(t, pool.bidMap)
	assert.Empty(t, pool.allBids)
	assert.Equal(t, newAuctioneer, pool.auctioneer)
	assert.Equal(t, newAuctionEntryPoint, pool.auctionEntryPoint)
}
