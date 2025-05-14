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
	"math/rand/v2"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestIsModuleTx(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	g := NewGaslessModule()
	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
	nodekey, _ := crypto.GenerateKey()
	err := g.Init(&InitOpts{
		ChainConfig:   testChainConfig,
		GaslessConfig: testGaslessConfig,
		NodeKey:       nodekey,
		Chain:         backend.BlockChain(),
		TxPool:        &testTxPool{},
	})
	require.NoError(t, err)

	privkey, _ := crypto.GenerateKey()
	testcases := []struct {
		tx *types.Transaction
		ok bool
	}{
		{
			makeApproveTx(t, privkey, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256}),
			true,
		},
		{
			makeSwapTx(t, privkey, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)}),
			true,
		},
		{
			makeTx(t, privkey, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil),
			false,
		},
		{
			nil,
			false,
		},
	}

	for _, tc := range testcases {
		ok := g.IsModuleTx(tc.tx)
		require.Equal(t, tc.ok, ok)
	}
}

func TestIsReady(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlTrace)

	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
	sdb, _ := backend.BlockChain().State()
	nodeKey, _ := crypto.GenerateKey()

	privkey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privkey.PublicKey)

	approveTx := func(nonce uint64) *types.Transaction {
		return makeApproveTx(t, privkey, nonce, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256})
	}

	swapTx := func(nonce uint64) *types.Transaction {
		return makeSwapTx(t, privkey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)})
	}

	singleSwapTx := func(nonce uint64) *types.Transaction {
		return makeSwapTx(t, privkey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000), Deadline: big.NewInt(300)})
	}

	other := func(nonce uint64) *types.Transaction {
		return makeTx(t, privkey, nonce, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)
	}

	testcases := map[string]struct {
		queue    map[uint64]*types.Transaction
		ready    types.Transactions
		i        uint64
		nonce    uint64
		expected bool
	}{
		// single swap test
		"correct single swap tx": {
			map[uint64]*types.Transaction{1: singleSwapTx(1), 2: other(2)},
			types.Transactions{},
			1,
			1,
			true,
		},
		"single swap tx with non-head nonce": {
			map[uint64]*types.Transaction{1: singleSwapTx(1), 2: other(2)},
			types.Transactions{other(0)},
			1,
			0,
			false,
		},
		// approve tx test
		"correct approve tx": {
			map[uint64]*types.Transaction{1: approveTx(1), 2: swapTx(2), 3: other(3)},
			types.Transactions{},
			1,
			1,
			true,
		},
		"approve tx without swap tx": {
			map[uint64]*types.Transaction{1: approveTx(1)},
			types.Transactions{},
			1,
			1,
			false,
		},
		"approve tx with non-sequentail swap tx": {
			map[uint64]*types.Transaction{1: approveTx(1), 2: other(2), 3: swapTx(3)},
			types.Transactions{},
			1,
			1,
			false,
		},
		"apporve tx with non-head nonce": {
			map[uint64]*types.Transaction{1: approveTx(1), 2: swapTx(2), 3: other(3)},
			types.Transactions{other(0)},
			1,
			0,
			false,
		},
		// swap test
		"correct swap tx": {
			map[uint64]*types.Transaction{2: swapTx(2), 3: other(3)},
			types.Transactions{approveTx(1)},
			2,
			1,
			true,
		},
		"swap tx without approve tx": {
			map[uint64]*types.Transaction{2: swapTx(2), 3: other(3)},
			types.Transactions{},
			2,
			1,
			false,
		},
		"gasless tx with non-sequential approve tx": {
			map[uint64]*types.Transaction{3: swapTx(3)},
			types.Transactions{approveTx(1), other(2)},
			3,
			1,
			false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			cdb := sdb.Copy()
			cdb.SetNonce(addr, tc.nonce)
			g := NewGaslessModule()
			err := g.Init(&InitOpts{
				ChainConfig:   testChainConfig,
				GaslessConfig: testGaslessConfig,
				NodeKey:       nodeKey,
				Chain:         backend.BlockChain(),
				TxPool:        &testTxPool{cdb},
			})
			require.NoError(t, err)

			ok := g.IsReady(tc.queue, tc.i, tc.ready)
			require.Equal(t, tc.expected, ok)
		})
	}
}

func TestPromoteGaslessTxsWithSingleSender(t *testing.T) {
	t.Parallel()

	type txTypeTest int
	const (
		T       txTypeTest = iota // regular tx
		A                         // approve tx
		SwithA                    // swap tx with approe tx
		SingleS                   // single swap tx
	)

	testTxPoolConfig := blockchain.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""

	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
	sdb, _ := backend.BlockChain().State()
	nodeKey, _ := crypto.GenerateKey()

	userKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	testcases := []struct {
		balance bool
		txs     []txTypeTest
		pending []txTypeTest
		queued  []txTypeTest
	}{
		{
			false,
			[]txTypeTest{A, T, SwithA, SingleS},
			[]txTypeTest{},
			[]txTypeTest{A, SwithA, SingleS},
		},
		{
			true,
			[]txTypeTest{A, T, SwithA, SingleS},
			[]txTypeTest{},
			[]txTypeTest{A, T, SwithA, SingleS},
		},
		{
			false,
			[]txTypeTest{A, SwithA, T, SingleS},
			[]txTypeTest{A, SwithA},
			[]txTypeTest{SingleS},
		},
		{
			true,
			[]txTypeTest{A, SwithA, T, SingleS},
			[]txTypeTest{A, SwithA, T},
			[]txTypeTest{SingleS},
		},
		{
			false,
			[]txTypeTest{A, SwithA, SingleS, T},
			[]txTypeTest{A, SwithA},
			[]txTypeTest{SingleS},
		},
		{
			true,
			[]txTypeTest{A, SwithA, SingleS, T},
			[]txTypeTest{A, SwithA},
			[]txTypeTest{SingleS, T},
		},
		{
			false,
			[]txTypeTest{SingleS, A, SwithA, T},
			[]txTypeTest{SingleS},
			[]txTypeTest{A, SwithA},
		},
		{
			true,
			[]txTypeTest{SingleS, A, SwithA, T},
			[]txTypeTest{SingleS},
			[]txTypeTest{A, SwithA, T},
		},
		{
			false,
			[]txTypeTest{SingleS, T, A, SwithA},
			[]txTypeTest{SingleS},
			[]txTypeTest{A, SwithA},
		},
		{
			true,
			[]txTypeTest{SingleS, T, A, SwithA},
			[]txTypeTest{SingleS, T},
			[]txTypeTest{A, SwithA},
		},
		{
			false,
			[]txTypeTest{T, A, SwithA, SingleS},
			[]txTypeTest{},
			[]txTypeTest{A, SwithA, SingleS},
		},
		{
			true,
			[]txTypeTest{T, A, SwithA, SingleS},
			[]txTypeTest{T},
			[]txTypeTest{A, SwithA, SingleS},
		},
		{
			false,
			[]txTypeTest{T, SingleS, A, SwithA},
			[]txTypeTest{},
			[]txTypeTest{SingleS, A, SwithA},
		},
		{
			true,
			[]txTypeTest{T, SingleS, A, SwithA},
			[]txTypeTest{T},
			[]txTypeTest{SingleS, A, SwithA},
		},
		{
			false,
			[]txTypeTest{SwithA, A, SingleS, T},
			[]txTypeTest{},
			[]txTypeTest{SwithA, A, SingleS},
		},
		{
			true,
			[]txTypeTest{SwithA, A, SingleS, T},
			[]txTypeTest{},
			[]txTypeTest{SwithA, A, SingleS, T},
		},
	}

	for _, tc := range testcases {
		cdb := sdb.Copy()
		if tc.balance {
			cdb.SetBalance(crypto.PubkeyToAddress(userKey.PublicKey), new(big.Int).SetUint64(params.KAIA))
		}
		bc := &testBlockChain{cdb, 10000000, new(event.Feed)}
		pool := blockchain.NewTxPool(testTxPoolConfig, testChainConfig, bc, &dummyGovModule{chainConfig: testChainConfig})
		g := NewGaslessModule()
		err := g.Init(&InitOpts{
			ChainConfig:   testChainConfig,
			GaslessConfig: testGaslessConfig,
			NodeKey:       nodeKey,
			Chain:         backend.BlockChain(),
			TxPool:        pool,
		})
		require.NoError(t, err)
		pool.RegisterTxPoolModule(g)
		txMap := map[txTypeTest]*types.Transaction{}

		for i, ttype := range tc.txs {
			nonce := uint64(i)
			var tx *types.Transaction
			switch ttype {
			case T:
				tx = makeTx(t, userKey, nonce, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x"))
			case A:
				tx = makeApproveTx(t, userKey, nonce, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256})
			case SwithA:
				tx = makeSwapTx(t, userKey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)})
			case SingleS:
				tx = makeSwapTx(t, userKey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000), Deadline: big.NewInt(300)})
			}
			txMap[ttype] = tx
			err = pool.AddLocal(tx)
			if err != nil {
				require.ErrorIs(t, err, blockchain.ErrInsufficientFundsFrom)
			}
		}

		pending, queued := pool.Content()
		pendingTxs := flattenPoolTxs(pending)
		queuedTxs := flattenPoolTxs(queued)

		require.Equal(t, len(tc.pending), len(pendingTxs))
		require.Equal(t, len(tc.queued), len(queuedTxs))

		for _, ttype := range tc.pending {
			_, ok := pendingTxs[txMap[ttype].Hash()]
			require.True(t, ok)
		}

		for _, ttype := range tc.queued {
			_, ok := queuedTxs[txMap[ttype].Hash()]
			require.True(t, ok)
		}

		pool.Stop()
	}
}

func TestPromoteGaslessTxsWithMultiSenders(t *testing.T) {
	t.Parallel()

	testTxPoolConfig := blockchain.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""

	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
	sdb, _ := backend.BlockChain().State()
	nodeKey, _ := crypto.GenerateKey()

	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	key4, _ := crypto.GenerateKey()

	// The number of variable means the number of sender. For exampele, A1 and S1 are Sender1's txs.
	A1 := makeApproveTx(t, key1, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256})
	S1 := makeSwapTx(t, key1, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)})

	A2 := makeApproveTx(t, key2, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: abi.MaxUint256})
	S2 := makeSwapTx(t, key2, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000), Deadline: big.NewInt(300)})

	S3 := makeSwapTx(t, nil, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000), Deadline: big.NewInt(300)})

	T4 := makeTx(t, key4, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)

	T5 := makeTx(t, nil, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)

	sdb.SetBalance(crypto.PubkeyToAddress(key2.PublicKey), new(big.Int).SetUint64(params.KAIA))
	sdb.SetBalance(crypto.PubkeyToAddress(key4.PublicKey), new(big.Int).SetUint64(params.KAIA))

	expected := []*types.Transaction{A1, S1, A2, S2, S3, T4}

	// send A1, S1, A2, S2, S3, T4, and T5 in random order and then check if pending has expected txs.
	for range make([]int, 1000) {
		cdb := sdb.Copy()
		bc := &testBlockChain{cdb, 10000000, new(event.Feed)}
		pool := blockchain.NewTxPool(testTxPoolConfig, testChainConfig, bc, &dummyGovModule{chainConfig: testChainConfig})
		g := NewGaslessModule()
		err := g.Init(&InitOpts{
			ChainConfig:   testChainConfig,
			GaslessConfig: testGaslessConfig,
			NodeKey:       nodeKey,
			Chain:         backend.BlockChain(),
			TxPool:        pool,
		})
		require.NoError(t, err)

		pool.RegisterTxPoolModule(g)

		txs := []*types.Transaction{A1, S1, A2, S2, S3, T4, T5}
		rand.Shuffle(len(txs), func(i, j int) {
			txs[i], txs[j] = txs[j], txs[i]
		})
		for i := range txs {
			err := pool.AddLocal(txs[i])
			if err != nil {
				require.ErrorIs(t, err, blockchain.ErrInsufficientFundsFrom)
			}
		}

		pending, err := pool.Pending()
		require.NoError(t, err)
		pendingTxs := flattenPoolTxs(pending)

		require.Equal(t, len(expected), len(pendingTxs))
		for i := range expected {
			_, ok := pendingTxs[expected[i].Hash()]
			require.True(t, ok)
		}

		pool.Stop()
	}
}
