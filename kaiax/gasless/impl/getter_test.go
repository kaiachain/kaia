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

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestIsApproveTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	privkey, _ := crypto.GenerateKey()
	// Legacy TestToken.approve(SwapRouter, 1000000)
	correct := makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)})
	testcases := map[string]struct {
		tx *types.Transaction
		ok bool
	}{
		"correct": {
			correct,
			true,
		},
		"invalid token address": {
			makeTx(t, privkey, 0, common.HexToAddress("0xffff"), big.NewInt(0), 1000000, big.NewInt(1), correct.Data()),
			false,
		},
		"invalid spender address": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0xffff"), Amount: big.NewInt(1000000)}),
			false,
		},
		"invalid amount": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(0)}),
			false,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		TxPool:      &testTxPool{},
	})
	require.NoError(t, err)
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsApproveTx(tc.tx)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestIsSwapTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	privkey, _ := crypto.GenerateKey()
	// Legacy TestRouter.swapForGas(Token, 10, 100, 2021000)
	correct := makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)})
	testcases := map[string]struct {
		tx *types.Transaction
		ok bool
	}{
		"correct": {
			correct,
			true,
		},
		"invalid swap router address": {
			makeTx(t, privkey, 0, common.HexToAddress("0xffff"), big.NewInt(0), 1000000, big.NewInt(1), correct.Data()),
			false,
		},
		"invalid token address": {
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xffff"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		TxPool:      &testTxPool{},
	})
	require.NoError(t, err)

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsSwapTx(tc.tx)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestIsExecutable(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	privkey, _ := crypto.GenerateKey()
	testcases := map[string]struct {
		approve *types.Transaction
		swap    *types.Transaction
		ok      bool
	}{
		"correct gasless tx pair": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			true,
		},
		"correct single swap tx": {
			nil,
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000)}),
			true,
		},
		"gasless tx pair with different sender address": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, nil, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
		"gasless tx pair with different token address": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xffff"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
		"gasless tx pair with invalid amount in": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(1000001), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
		"gasless tx pair with non sequential nonce": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, privkey, 2, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
		"gasless tx pair with non head nonce": {
			makeApproveTx(t, privkey, 1, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}),
			makeSwapTx(t, privkey, 2, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}),
			false,
		},
		"single swap tx with non head nonce": {
			nil,
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1)}),
			false,
		},
		"single swap tx with invalid repay amount": {
			nil,
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1)}),
			false,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	err := g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		TxPool:      &testTxPool{sdb},
	})
	require.NoError(t, err)

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsExecutable(tc.approve, tc.swap)
			require.Equal(t, tc.ok, ok)
		})
	}
}
