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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/auction"
	cn "github.com/kaiachain/kaia/node/cn/filters/mock"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	db := database.NewMemoryDBManager()
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	mockCtrl := gomock.NewController(t)
	mockBackend := cn.NewMockBackend(mockCtrl)
	apiBackend := &MockBackend{mockBackend}
	fakeDownloader := &downloader.FakeDownloader{}

	// Set up the auction config
	auctionInitConfig := &InitOpts{
		ChainConfig: testChainConfig,
		Backend:     apiBackend,
		Downloader:  fakeDownloader,
		NodeKey:     key,
	}

	tcs := map[string]struct {
		balance  *big.Int
		disabled bool
	}{
		"insufficient balance": {
			balance:  big.NewInt(0),
			disabled: true,
		},
		"sufficient balance": {
			balance:  AuctionLenderMinBal,
			disabled: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			alloc := blockchain.GenesisAlloc{
				addr: {
					Balance: tc.balance,
				},
			}
			backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)

			auctionInitConfig.AuctionConfig = auction.DefaultAuctionConfig()
			auctionInitConfig.Chain = backend.BlockChain()

			am := NewAuctionModule()
			am.Init(auctionInitConfig)

			assert.Equal(t, tc.disabled, am.IsDisabled())
		})
	}
}
