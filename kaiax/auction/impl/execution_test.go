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
	"context"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"gotest.tools/assert"
)

func testRandaoForkChainConfig(forkNum *big.Int) *params.ChainConfig {
	var config *params.ChainConfig

	config = &params.ChainConfig{
		ChainID: common.Big1,
		Governance: &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				UseGiniCoeff:          true,
				StakingUpdateInterval: 86400,
			},
			KIP71: params.GetDefaultKIP71Config(),
		},
	}
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = forkNum

	return config
}

func testAllocStorage() blockchain.GenesisAlloc {
	allocStorage := system.AllocRegistry(&params.RegistryConfig{
		Records: map[string]common.Address{
			system.AuctionEntryPointName: system.AuctionEntryPointAddrMock,
		},
		Owner: common.HexToAddress("0xffff"),
	})
	storage := make(map[common.Hash]common.Hash)
	storage[common.BytesToHash([]byte{0x00})] = common.BytesToHash([]byte{0x01})
	alloc := blockchain.GenesisAlloc{
		system.RegistryAddr: {
			Code:    system.RegistryMockCode,
			Balance: big.NewInt(0),
			Storage: allocStorage,
		},
		system.AuctionEntryPointAddrMock: {
			Code:    system.AuctionEntryPointMockCode,
			Balance: big.NewInt(0),
			Storage: storage,
		},
	}
	return alloc
}

type MockBackend struct{}

func (m *MockBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return nil
}

func TestUpdateAuctionInfo(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db     = database.NewMemoryDBManager()
		alloc  = testAllocStorage()
		config = testRandaoForkChainConfig(big.NewInt(0))
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	mAuction := NewAuctionModule()
	apiBackend := &MockBackend{}
	fakeDownloader := &downloader.FakeDownloader{}
	mAuction.Init(&InitOpts{
		ChainConfig: config,
		Chain:       backend.BlockChain(),
		Backend:     apiBackend,
		Downloader:  fakeDownloader,
	})

	// Not updated yet
	assert.Equal(t, mAuction.bidPool.auctioneer, common.Address{})
	assert.Equal(t, mAuction.bidPool.auctionEntryPoint, common.Address{})

	mAuction.PostInsertBlock(backend.BlockChain().CurrentBlock())

	assert.Equal(t, mAuction.bidPool.auctioneer, common.HexToAddress("0x01"))
	assert.Equal(t, mAuction.bidPool.auctionEntryPoint, system.AuctionEntryPointAddrMock)

	// Auction is running
	assert.Equal(t, atomic.LoadUint32(&mAuction.bidPool.running), uint32(1))

	// Update auction info again
	mAuction.PostInsertBlock(types.NewBlock(&types.Header{Number: big.NewInt(1)}, nil, nil))

	assert.Equal(t, mAuction.bidPool.auctioneer, common.Address{})
	assert.Equal(t, mAuction.bidPool.auctionEntryPoint, common.Address{})

	// Auction is not running
	assert.Equal(t, atomic.LoadUint32(&mAuction.bidPool.running), uint32(0))
}
