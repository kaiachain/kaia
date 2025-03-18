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

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsApproveTx(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	g := NewGaslessModule()
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
	key, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig:   testChainConfig,
		GaslessConfig: testGaslessConfig,
		NodeKey:       key,
		Chain:         backend.BlockChain(),
		TxPool:        &testTxPool{},
	})
	require.NoError(t, err)

	privkey, _ := crypto.GenerateKey()
	correct := makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256})

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
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0xffff"), Amount: abi.MaxUint256}),
			false,
		},
		"invalid amount": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(0)}),
			false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsApproveTx(tc.tx)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestIsSwapTx(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	g := NewGaslessModule()
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
	key, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig:   testChainConfig,
		GaslessConfig: testGaslessConfig,
		NodeKey:       key,
		Chain:         backend.BlockChain(),
		TxPool:        &testTxPool{},
	})
	require.NoError(t, err)

	privkey, _ := crypto.GenerateKey()
	correct := makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)})

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
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xffff"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsSwapTx(tc.tx)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestIsExecutable(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	g := NewGaslessModule()
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
	sdb, _ := backend.BlockChain().State()
	key, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig:   testChainConfig,
		GaslessConfig: testGaslessConfig,
		NodeKey:       key,
		Chain:         backend.BlockChain(),
		TxPool:        &testTxPool{sdb},
	})
	require.NoError(t, err)

	privkey, _ := crypto.GenerateKey()
	testcases := map[string]struct {
		approve *types.Transaction
		swap    *types.Transaction
		ok      bool
	}{
		"correct gasless tx pair": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			true,
		},
		"correct single swap tx": {
			nil,
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(1), AmountRepay: big.NewInt(1021000), Deadline: big.NewInt(300)}),
			true,
		},
		"gasless tx pair with different sender address": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, nil, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			false,
		},
		"gasless tx pair with different token address": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xffff"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			false,
		},
		"gasless tx pair with non sequential nonce": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, privkey, 2, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			false,
		},
		"gasless tx pair with non head nonce": {
			makeApproveTx(t, privkey, 1, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, privkey, 2, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			false,
		},
		"single swap tx with non head nonce": {
			nil,
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1), Deadline: big.NewInt(300)}),
			false,
		},
		"single swap tx with invalid repay amount": {
			nil,
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1), Deadline: big.NewInt(300)}),
			false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ok := g.IsExecutable(tc.approve, tc.swap)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestGetLendTxGenerator(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
	nodekey, _ := crypto.GenerateKey()
	sdb, _ := backend.BlockChain().State()
	sdb.SetBalance(crypto.PubkeyToAddress(nodekey.PublicKey), new(big.Int).SetUint64(params.KAIA))

	testTxPoolConfig := blockchain.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""

	privkey, _ := crypto.GenerateKey()

	testcases := map[string]struct {
		approve *types.Transaction
		swap    *types.Transaction
	}{
		"correct gasless tx pair": {
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			makeSwapTx(t, privkey, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
		},
		"correct single swap tx": {
			nil,
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(1), AmountRepay: big.NewInt(1021000), Deadline: big.NewInt(300)}),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			g := NewGaslessModule()
			bc := &testBlockChain{sdb.Copy(), 10000000, new(event.Feed)}
			pool := blockchain.NewTxPool(testTxPoolConfig, testChainConfig, bc, &dummyGovModule{chainConfig: testChainConfig})
			err := g.Init(&InitOpts{
				ChainConfig:   testChainConfig,
				GaslessConfig: testGaslessConfig,
				NodeKey:       nodekey,
				Chain:         backend.BlockChain(),
				TxPool:        pool,
			})
			require.NoError(t, err)

			pool.RegisterTxPoolModule(g)

			ok := g.IsExecutable(tc.approve, tc.swap)
			require.True(t, ok)

			generator := g.GetLendTxGenerator(tc.approve, tc.swap)
			tx, err := generator.GetTx(0)
			require.NoError(t, err)

			// tx contents test
			require.Equal(t, crypto.PubkeyToAddress(privkey.PublicKey).Bytes(), tx.To().Bytes())
			lendAmount := tc.swap.Fee()
			if tc.approve != nil {
				lendAmount.Add(lendAmount, tc.approve.Fee())
			}
			require.Zero(t, lendAmount.Cmp(tx.Value()))

			// pool passing test
			pool.AddLocal(tx)
			pending, err := pool.Pending()
			require.NoError(t, err)
			flatten := flattenPoolTxs(pending)
			require.True(t, flatten[tx.Hash()])
		})
	}
}

func TestTxGeneratorHashUniqueness(t *testing.T) {
	hashSet := make(map[common.Hash]struct{})
	g := NewGaslessModule()
	for range 100 {
		generator := g.GetLendTxGenerator(nil, makeSwapTx(t, nil, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(1), AmountRepay: big.NewInt(1021000)}))
		_, ok := hashSet[generator.Id]
		assert.False(t, ok)
		hashSet[generator.Id] = struct{}{}
	}
}
