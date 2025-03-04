// Copyright 2024 The Kaia Authors
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
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func init() {
	conf := params.TestChainConfig.Copy()
	conf.IstanbulCompatibleBlock = big.NewInt(0)
	conf.LondonCompatibleBlock = big.NewInt(0)
	conf.EthTxTypeCompatibleBlock = big.NewInt(0)
	conf.MagmaCompatibleBlock = big.NewInt(0)
	fork.SetHardForkBlockNumberConfig(conf)
	blockchain.InitDeriveSha(conf)
}

func TestIsModuleTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	testcases := []struct {
		tx *types.Transaction
		ok bool
	}{
		{
			types.NewTransaction(0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1),
				hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240")),
			true,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		StateDB:     statedb,
	})
	for _, tc := range testcases {
		ok := g.IsModuleTx(tc.tx)
		require.Equal(t, tc.ok, ok)
	}
}

func TestIsReady(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)

	privkey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privkey.PublicKey)

	nodeKey, _ := crypto.GenerateKey()
	g := NewGaslessModule()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     nodeKey,
		StateDB:     statedb,
	})
	approveTx := func(nonce uint64) *types.Transaction {
		tx, err := makeTx(privkey, nonce, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240"))
		require.NoError(t, err)
		return tx
	}

	swapTx := func(nonce uint64) *types.Transaction {
		tx, err := makeTx(privkey, nonce, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000001ed688"))
		require.NoError(t, err)
		return tx
	}

	swapTx2 := func(nonce uint64) *types.Transaction {
		tx, err := makeTx(privkey, nonce, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000f9448"))
		require.NoError(t, err)
		return tx
	}

	other := func(nonce uint64) *types.Transaction {
		tx, err := makeTx(privkey, nonce, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)
		require.NoError(t, err)
		return tx
	}

	testcases := []struct {
		queue    map[uint64]*types.Transaction
		ready    types.Transactions
		i        uint64
		nonce    uint64
		expected bool
	}{
		{
			map[uint64]*types.Transaction{1: swapTx2(1)},
			types.Transactions{},
			1,
			1,
			true,
		},
		{
			map[uint64]*types.Transaction{1: approveTx(1), 2: swapTx(2), 3: other(3), 4: other(4)},
			types.Transactions{},
			1,
			1,
			true,
		},
		{
			map[uint64]*types.Transaction{2: swapTx(2), 3: other(3), 4: other(4)},
			types.Transactions{1: approveTx(1)},
			2,
			1,
			true,
		},
		{
			map[uint64]*types.Transaction{2: swapTx(2), 3: other(3), 4: other(4)},
			types.Transactions{},
			2,
			1,
			false,
		},
	}

	for _, tc := range testcases {
		g.StateDB.SetNonce(addr, tc.nonce)
		ok := g.IsReady(tc.queue, tc.i, tc.ready)
		require.Equal(t, tc.expected, ok)
	}
}

func TestPromoteGaslessTransactions(t *testing.T) {
	t.Parallel()

	type TxTypeTest int
	const (
		T       TxTypeTest = iota // regular tx
		A                         // approve tx
		SwithA                    // swap tx with approe tx
		SingleS                   // single swap tx
	)

	conf := params.TestChainConfig.Copy()
	conf.IstanbulCompatibleBlock = big.NewInt(0)
	conf.LondonCompatibleBlock = big.NewInt(0)
	conf.EthTxTypeCompatibleBlock = big.NewInt(0)
	conf.MagmaCompatibleBlock = big.NewInt(0)

	testTxPoolConfig := blockchain.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)

	g := NewGaslessModule()
	nodeKey, _ := crypto.GenerateKey()
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     nodeKey,
	})

	userKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	testcases := []struct {
		balance bool
		txs     []TxTypeTest
		pending []TxTypeTest
		queued  []TxTypeTest
	}{
		{
			false,
			[]TxTypeTest{T, A, SwithA, SingleS},
			[]TxTypeTest{},
			[]TxTypeTest{A, SwithA, SingleS},
		},
		{
			true,
			[]TxTypeTest{T, A, SwithA, SingleS},
			[]TxTypeTest{T},
			[]TxTypeTest{A, SwithA, SingleS},
		},
		{
			false,
			[]TxTypeTest{T, SingleS, A, SwithA},
			[]TxTypeTest{},
			[]TxTypeTest{SingleS, A, SwithA},
		},
		{
			true,
			[]TxTypeTest{T, SingleS, A, SwithA},
			[]TxTypeTest{T},
			[]TxTypeTest{SingleS, A, SwithA},
		},
		{
			false,
			[]TxTypeTest{A, T, SwithA, SingleS},
			[]TxTypeTest{},
			[]TxTypeTest{A, SwithA, SingleS},
		},
		{
			true,
			[]TxTypeTest{A, T, SwithA, SingleS},
			[]TxTypeTest{},
			[]TxTypeTest{A, T, SwithA, SingleS},
		},
		{
			false,
			[]TxTypeTest{A, SwithA, T, SingleS},
			[]TxTypeTest{A, SwithA},
			[]TxTypeTest{SingleS},
		},
		{
			true,
			[]TxTypeTest{A, SwithA, T, SingleS},
			[]TxTypeTest{A, SwithA, T},
			[]TxTypeTest{SingleS},
		},
		{
			false,
			[]TxTypeTest{A, SwithA, SingleS, T},
			[]TxTypeTest{A, SwithA},
			[]TxTypeTest{SingleS},
		},
		{
			true,
			[]TxTypeTest{A, SwithA, SingleS, T},
			[]TxTypeTest{A, SwithA},
			[]TxTypeTest{SingleS, T},
		},
		{
			false,
			[]TxTypeTest{SingleS, A, SwithA, T},
			[]TxTypeTest{SingleS},
			[]TxTypeTest{A, SwithA},
		},
		{
			true,
			[]TxTypeTest{SingleS, A, SwithA, T},
			[]TxTypeTest{SingleS},
			[]TxTypeTest{A, SwithA, T},
		},
		{
			false,
			[]TxTypeTest{SingleS, T, A, SwithA},
			[]TxTypeTest{SingleS},
			[]TxTypeTest{A, SwithA},
		},
		{
			true,
			[]TxTypeTest{SingleS, T, A, SwithA},
			[]TxTypeTest{SingleS, T},
			[]TxTypeTest{A, SwithA},
		},
	}

	for _, tc := range testcases {
		sdb := statedb.Copy()
		g.StateDB = sdb
		if tc.balance {
			// set some token
			sdb.SetBalance(crypto.PubkeyToAddress(userKey.PublicKey), new(big.Int).SetUint64(params.KAIA))
		}
		bc := &testBlockChain{sdb, 10000000, new(event.Feed)}
		pool := blockchain.NewTxPool(testTxPoolConfig, conf, bc, &dummyGovModule{chainConfig: conf}, []kaiax.TxPoolModule{g})

		txMap := map[TxTypeTest]*types.Transaction{}

		for i, ttype := range tc.txs {
			nonce := uint64(i)
			var tx *types.Transaction
			switch ttype {
			case T:
				tx, err = makeTx(userKey, nonce, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x"))
			case A:
				tx, err = makeApproveTx(userKey, nonce, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)})
			case SwithA:
				tx, err = makeSwapTx(userKey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)})
			case SingleS:
				tx, err = makeSwapTx(userKey, nonce, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000)})
			}
			require.NoError(t, err)
			txMap[ttype] = tx
			pool.AddLocal(tx)
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

// helper

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{Number: big.NewInt(0)}, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

type dummyGovModule struct {
	chainConfig *params.ChainConfig
}

func (d *dummyGovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	return gov.ParamSet{UnitPrice: d.chainConfig.UnitPrice}
}

func flattenPoolTxs(structured map[common.Address]types.Transactions) map[common.Hash]bool {
	flattened := map[common.Hash]bool{}
	for _, txs := range structured {
		for _, tx := range txs {
			flattened[tx.Hash()] = true
		}
	}
	return flattened
}

func makeApproveTx(privKey *ecdsa.PrivateKey, nonce uint64, approveArgs ApproveArgs) (*types.Transaction, error) {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
	}

	data := append([]byte{}, common.Hex2Bytes("095ea7b3")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Spender.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Amount.Bytes()), 32)...)
	approveTx, err := makeTx(privKey, nonce, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), data)
	if err != nil {
		return nil, err
	}

	return approveTx, nil
}

func makeSwapTx(privKey *ecdsa.PrivateKey, nonce uint64, swapArgs SwapArgs) (*types.Transaction, error) {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
	}

	data := append([]byte{}, common.Hex2Bytes("43bab9f7")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.Token.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountIn.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.MinAmountOut.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountRepay.Bytes()), 32)...)
	swapTx, err := makeTx(privKey, nonce, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), data)
	if err != nil {
		return nil, err
	}

	return swapTx, nil
}
