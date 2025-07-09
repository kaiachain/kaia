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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestFilterTxs(t *testing.T) {
	auctionModule := NewAuctionModule()

	txs := make(map[common.Address]types.Transactions)

	for i := 0; i < 5; i++ {
		// tx.Time is set to current time
		tx := types.NewTransaction(uint64(i), common.HexToAddress("0x1"), big.NewInt(1), 1000000, big.NewInt(1), nil)
		tx2 := types.NewTransaction(uint64(i), common.HexToAddress("0x2"), big.NewInt(1), 1000000, big.NewInt(1), nil)

		time.Sleep(40 * time.Millisecond)
		txs[common.HexToAddress("0x1")] = append(txs[common.HexToAddress("0x1")], tx)
		txs[common.HexToAddress("0x2")] = append(txs[common.HexToAddress("0x2")], tx2)
	}

	// [0, 40, 80, 120, 160] -> current: 200 -> [0, 40]
	auctionModule.FilterTxs(txs)

	assert.Equal(t, len(txs[common.HexToAddress("0x1")]), 2)
	assert.Equal(t, len(txs[common.HexToAddress("0x2")]), 2)

	for i, tx := range txs[common.HexToAddress("0x1")] {
		assert.Equal(t, tx.Nonce(), uint64(i))
	}

	for i, tx := range txs[common.HexToAddress("0x2")] {
		assert.Equal(t, tx.Nonce(), uint64(i))
	}
}
