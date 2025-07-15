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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/kaiax/gasless/mock"
	"github.com/stretchr/testify/assert"
)

var (
	tx1Sender = common.HexToAddress("0x1")
	tx2Sender = common.HexToAddress("0x2")
)

func TestFilterTxs(t *testing.T) {
	mAuction := prep(t)
	txs := make(map[common.Address]types.Transactions)

	for i := 0; i < 5; i++ {
		// tx.Time is set to current time
		tx := types.NewTransaction(uint64(i), tx1Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)
		tx2 := types.NewTransaction(uint64(i), tx2Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)

		time.Sleep(40 * time.Millisecond)
		txs[tx1Sender] = append(txs[tx1Sender], tx)
		txs[tx2Sender] = append(txs[tx2Sender], tx2)
	}

	// Not running
	mAuction.bidPool.running = 0
	mAuction.FilterTxs(txs)

	assert.Equal(t, 5, len(txs[tx1Sender]))
	assert.Equal(t, 5, len(txs[tx2Sender]))

	// Running
	mAuction.bidPool.running = 1
	mAuction.FilterTxs(txs)

	assert.Equal(t, 2, len(txs[tx1Sender]))
	assert.Equal(t, 2, len(txs[tx2Sender]))

	for i, tx := range txs[tx1Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}

	for i, tx := range txs[tx2Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}
}

func TestFilterTxs_TargetTx(t *testing.T) {
	mAuction := prep(t)
	txs := make(map[common.Address]types.Transactions)

	for i := 0; i < 5; i++ {
		// tx.Time is set to current time
		tx := types.NewTransaction(uint64(i), tx1Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)
		tx2 := types.NewTransaction(uint64(i), tx2Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)

		time.Sleep(40 * time.Millisecond)
		txs[tx1Sender] = append(txs[tx1Sender], tx)
		txs[tx2Sender] = append(txs[tx2Sender], tx2)
	}

	// Running
	mAuction.bidPool.bidTargetMap[1] = map[common.Hash]*auction.Bid{
		txs[tx1Sender][2].Hash(): nil,
		txs[tx2Sender][2].Hash(): nil,
	}

	mAuction.FilterTxs(txs)

	assert.Equal(t, 3, len(txs[tx1Sender]))
	assert.Equal(t, 3, len(txs[tx2Sender]))

	for i, tx := range txs[tx1Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}

	for i, tx := range txs[tx2Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}
}

func TestFilterTxs_GaslessTx(t *testing.T) {
	mAuction := prep(t)
	txs := make(map[common.Address]types.Transactions)

	for i := 0; i < 5; i++ {
		// tx.Time is set to current time
		tx := types.NewTransaction(uint64(i), tx1Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)
		tx2 := types.NewTransaction(uint64(i), tx2Sender, big.NewInt(1), 1000000, big.NewInt(1), nil)

		time.Sleep(40 * time.Millisecond)
		txs[tx1Sender] = append(txs[tx1Sender], tx)
		txs[tx2Sender] = append(txs[tx2Sender], tx2)
	}

	mockCtrl := gomock.NewController(t)
	mockGasless := mock.NewMockGaslessModule(mockCtrl)
	mockGasless.EXPECT().IsBundleTx(gomock.Any()).Return(true).AnyTimes()

	mAuction.RegisterGaslessModule(mockGasless)
	mAuction.FilterTxs(txs)

	assert.Equal(t, 5, len(txs[tx1Sender]))
	assert.Equal(t, 5, len(txs[tx2Sender]))

	for i, tx := range txs[tx1Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}

	for i, tx := range txs[tx2Sender] {
		assert.Equal(t, uint64(i), tx.Nonce())
	}
}
