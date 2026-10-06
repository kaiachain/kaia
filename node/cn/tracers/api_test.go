// Modifications Copyright 2022 The klaytn Authors
// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/tracers/api_test.go (2022/08/08).
// Modified and improved for the klaytn development.
package tracers

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	kaiaapi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
)

var (
	errStateNotFound       = errors.New("state not found")
	errBlockNotFound       = errors.New("block not found")
	errTransactionNotFound = errors.New("transaction not found")
)

func emptyStructTraceResult(gas uint64) json.RawMessage {
	result, err := json.Marshal(&kaiaapi.ExecutionResult{
		Gas:         gas,
		Failed:      false,
		ReturnValue: "",
		StructLogs:  []kaiaapi.StructLogRes{},
	})
	if err != nil {
		panic(err)
	}
	return result
}

type testBackend struct {
	chainConfig *params.ChainConfig
	sealer      consensus.Sealer
	chaindb     database.DBManager
	chain       *blockchain.BlockChain

	refHook func() // Hook is invoked when the requested state is referenced
	relHook func() // Hook is invoked when the requested state is released
}

func newTestBackend(t *testing.T, n int, gspec *blockchain.Genesis, generator func(i int, b *blockchain.BlockGen)) *testBackend {
	backend := &testBackend{
		chainConfig: params.TestChainConfig,
		sealer:      faker.NewFaker(),
		chaindb:     database.NewMemoryDBManager(),
	}
	// Generate blocks for testing
	gspec.Config = backend.chainConfig
	var (
		gendb   = database.NewMemoryDBManager()
		genesis = gspec.MustCommit(gendb)
	)
	blocks, _ := blockchain.GenerateChain(backend.chainConfig, genesis, backend.sealer, gendb, n, generator)
	// Import the canonical chain
	gspec.MustCommit(backend.chaindb)
	cacheConfig := &blockchain.CacheConfig{
		CacheSize:           512,
		BlockInterval:       blockchain.DefaultBlockInterval,
		TriesInMemory:       blockchain.DefaultTriesInMemory,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:   512,
		ArchiveMode:         true, // Archive mode
	}
	chain, err := blockchain.NewBlockChain(backend.chaindb, cacheConfig, backend.chainConfig, backend.sealer, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	backend.chain = chain
	return backend
}

func (b *testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.chain.GetHeaderByHash(hash), nil
}

func (b *testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentHeader(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.chain.GetBlockByHash(hash), nil
}

func (b *testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	block := b.chain.GetBlockByNumber(uint64(number))
	if block == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", number)
	}
	return block, nil
}

func (b *testBackend) GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	tx, hash, blockNumber, index := b.chain.GetTxAndLookupInfoInCache(txHash)
	if tx == nil {
		return nil, common.Hash{}, 0, 0
	}
	return tx, hash, blockNumber, index
}

func (b *testBackend) RPCGasCap() *big.Int {
	return big.NewInt(int64(params.UpperGasLimit))
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	return b.chainConfig
}

func (b *testBackend) Sealer() consensus.Sealer {
	return b.sealer
}

func (b *testBackend) ChainDB() database.DBManager {
	return b.chaindb
}

func (b *testBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (*state.StateDB, StateReleaseFunc, error) {
	statedb, err := b.chain.StateAt(block.Root())
	if err != nil {
		return nil, nil, errStateNotFound
	}
	if b.refHook != nil {
		b.refHook()
	}
	release := func() {
		if b.relHook != nil {
			b.relHook()
		}
	}
	return statedb, release, nil
}

func (b *testBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (blockchain.Message, vm.BlockContext, vm.TxContext, *state.StateDB, StateReleaseFunc, error) {
	parent := b.chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.BlockContext{}, vm.TxContext{}, nil, nil, errBlockNotFound
	}
	statedb, release, err := b.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, vm.BlockContext{}, vm.TxContext{}, nil, nil, errStateNotFound
	}
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.BlockContext{}, vm.TxContext{}, statedb, release, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner(b.chainConfig, block.Number())
	for idx, tx := range block.Transactions() {
		msg, _ := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		txContext := blockchain.NewEVMTxContext(msg, block.Header(), b.chainConfig)
		blockContext := blockchain.NewEVMBlockContext(block.Header(), b.chain, nil)
		if idx == txIndex {
			return msg, blockContext, txContext, statedb, release, nil
		}
		vmenv := vm.NewEVM(blockContext, txContext, statedb, b.chainConfig, &vm.Config{Debug: true, EnableInternalTxTracing: true})
		if _, err := blockchain.ApplyMessage(vmenv, msg); err != nil {
			return nil, vm.BlockContext{}, vm.TxContext{}, nil, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		statedb.Finalise(true, true)
	}
	return nil, vm.BlockContext{}, vm.TxContext{}, nil, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}

func TestTraceChain(t *testing.T) {
	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.KAIA)},
		accounts[1].addr: {Balance: big.NewInt(params.KAIA)},
		accounts[2].addr: {Balance: big.NewInt(params.KAIA)},
	}}
	genBlocks := 50
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)

	var (
		ref   uint32 // total refs has made
		rel   uint32 // total rels has made
		nonce uint64
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		for j := 0; j < i+1; j++ {
			tx, _ := types.SignTx(types.NewTransaction(nonce, accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
			b.AddTx(tx)
			nonce += 1
		}
	})
	backend.refHook = func() { atomic.AddUint32(&ref, 1) }
	backend.relHook = func() { atomic.AddUint32(&rel, 1) }
	api := NewAPI(backend)

	single := `{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}`
	cases := []struct {
		start  uint64
		end    uint64
		config *TraceConfig
	}{
		{0, 50, nil},  // the entire chain range, blocks [1, 50]
		{10, 20, nil}, // the middle chain range, blocks [11, 20]
	}
	for _, c := range cases {
		ref, rel = 0, 0 // clean up the counters

		from, _ := api.blockByNumber(context.Background(), rpc.BlockNumber(c.start))
		to, _ := api.blockByNumber(context.Background(), rpc.BlockNumber(c.end))
		ret, err := api.traceChain(from, to, c.config, nil, nil)
		assert.NoError(t, err)

		for _, trace := range ret {
			for _, txTrace := range trace.Traces {
				blob, _ := json.Marshal(txTrace.Result)
				if string(blob) != single {
					t.Error("Unexpected tracing result")
				}
			}
		}
	}
}

func TestTraceCall(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(0)},
		accounts[1].addr: {Balance: big.NewInt(1000 * 10)},
		accounts[2].addr: {Balance: big.NewInt(0)},
	}}
	genBlocks := 10
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	api := NewAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[1] to account[0]
		//    value: 1000 kei
		//    fee:   0 kei
		tx, err := types.SignTx(types.NewTransaction(uint64(i), accounts[0].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[1].key)
		assert.NoError(t, err)
		b.AddTx(tx)
	}))

	testSuite := []struct {
		blockNumber rpc.BlockNumber
		call        kaiaapi.CallArgs
		config      *TraceConfig
		expectErr   error
		expect      interface{}
	}{
		// Standard JSON trace upon the genesis, plain transfer.
		{
			blockNumber: rpc.BlockNumber(0),
			call: kaiaapi.CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: hexutil.Big(*big.NewInt(1000)),
			},
			config:    nil,
			expectErr: errors.New("tracing failed: insufficient balance for transfer"),
			expect:    nil,
		},
		// Standard JSON trace upon the head, plain transfer.
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			call: kaiaapi.CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: hexutil.Big(*big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    emptyStructTraceResult(params.TxGas),
		},
		// Standard JSON trace upon the non-existent block, error expects
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			call: kaiaapi.CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: hexutil.Big(*big.NewInt(1000)),
			},
			config:    nil,
			expectErr: fmt.Errorf("the block does not exist (block number: %d)", genBlocks+1),
			expect:    nil,
		},
		// Standard JSON trace upon the latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: kaiaapi.CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: hexutil.Big(*big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    emptyStructTraceResult(params.TxGas),
		},
		// Standard JSON trace upon the pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			call: kaiaapi.CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: hexutil.Big(*big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    emptyStructTraceResult(params.TxGas),
		},
	}
	for _, testspec := range testSuite {
		result, err := api.TraceCall(context.Background(), testspec.call, rpc.BlockNumberOrHash{BlockNumber: &testspec.blockNumber}, testspec.config)
		assert.Equal(t, testspec.expectErr, err)
		assert.Equal(t, result, testspec.expect)
	}
}

type tracedPrestateAccount struct {
	Balance *hexutil.Big                `json:"balance"`
	Nonce   uint64                      `json:"nonce"`
	Code    hexutil.Bytes               `json:"code"`
	Storage map[common.Hash]common.Hash `json:"storage"`
}

func assertTraceBalance(t *testing.T, account tracedPrestateAccount, want *big.Int) {
	t.Helper()
	if assert.NotNil(t, account.Balance) {
		assert.Equal(t, 0, (*big.Int)(account.Balance).Cmp(want))
	}
}

func traceCallPrestate(t *testing.T, api *API, call kaiaapi.CallArgs, config *TraceConfig) map[common.Address]tracedPrestateAccount {
	t.Helper()
	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), call, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return nil
	}
	raw, ok := result.(json.RawMessage)
	if !assert.True(t, ok, "expected json.RawMessage result, got %T", result) {
		return nil
	}
	var prestate map[common.Address]tracedPrestateAccount
	assert.NoError(t, json.Unmarshal(raw, &prestate))
	return prestate
}

func TestTraceCallPrestateTracerKeepsLogicalPrestate(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000aaaa")
	to := common.HexToAddress("0x000000000000000000000000000000000000bbbb")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	tracerName := "prestateTracer"
	gas := hexutil.Uint64(params.TxGas)
	gasPrice := hexutil.Big(*big.NewInt(1))
	call := kaiaapi.CallArgs{From: from, To: &to, Gas: &gas, GasPrice: &gasPrice}

	prestate := traceCallPrestate(t, api, call, &TraceConfig{Tracer: &tracerName})
	assertTraceBalance(t, prestate[from], big.NewInt(0))

	overrideBalance := hexutil.Big(*big.NewInt(0x42))
	overrideBalancePtr := &overrideBalance
	overrides := kaiaapi.EthStateOverride{
		from: {Balance: &overrideBalancePtr},
	}
	prestate = traceCallPrestate(t, api, call, &TraceConfig{Tracer: &tracerName, StateOverrides: &overrides})
	assertTraceBalance(t, prestate[from], big.NewInt(0x42))
}

func TestTraceCallPrestateTracerDiffModeKeepsLogicalPostBalance(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000a111")
	to := common.HexToAddress("0x000000000000000000000000000000000000b111")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	tracerName := "prestateTracer"
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(1))
	value := hexutil.Big(*big.NewInt(0x1234))
	overrideBalance := hexutil.Big(*big.NewInt(100000))
	overrideBalancePtr := &overrideBalance
	overrides := kaiaapi.EthStateOverride{
		from: {Balance: &overrideBalancePtr},
	}
	config := &TraceConfig{
		Tracer:         &tracerName,
		TracerConfig:   json.RawMessage(`{"diffMode":true}`),
		StateOverrides: &overrides,
	}

	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		To:       &to,
		Gas:      &gas,
		GasPrice: &gasPrice,
		Value:    value,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	raw, ok := result.(json.RawMessage)
	if !assert.True(t, ok, "expected json.RawMessage result, got %T", result) {
		return
	}

	var diff struct {
		Pre  map[common.Address]tracedPrestateAccount `json:"pre"`
		Post map[common.Address]tracedPrestateAccount `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(raw, &diff))

	assertTraceBalance(t, diff.Pre[from], (*big.Int)(&overrideBalance))
	wantPost := new(big.Int).Sub((*big.Int)(&overrideBalance), (*big.Int)(&value))
	assertTraceBalance(t, diff.Post[from], wantPost)
}

func TestTraceCallPrestateTracerDiffModeHidesMagmaFee(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000a121")
	rewardbase := common.HexToAddress("0x000000000000000000000000000000000000b121")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, func(i int, b *blockchain.BlockGen) {
		b.SetRewardbase(rewardbase)
	}))

	tracerName := "prestateTracer"
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(1))
	value := hexutil.Big(*big.NewInt(0x10000))
	overrideBalance := hexutil.Big(*big.NewInt(0x100000))
	overrideBalancePtr := &overrideBalance
	overrides := kaiaapi.EthStateOverride{
		from: {Balance: &overrideBalancePtr},
	}
	config := &TraceConfig{
		Tracer:         &tracerName,
		TracerConfig:   json.RawMessage(`{"diffMode":true}`),
		StateOverrides: &overrides,
	}

	// The call target is also the Magma rewardbase, so it receives both the
	// transferred value and the (half-burned) synthetic gas fee. The fee must be
	// hidden, leaving only the value in its post balance.
	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		To:       &rewardbase,
		Gas:      &gas,
		GasPrice: &gasPrice,
		Value:    value,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	raw, ok := result.(json.RawMessage)
	if !assert.True(t, ok, "expected json.RawMessage result, got %T", result) {
		return
	}

	var diff struct {
		Pre  map[common.Address]tracedPrestateAccount `json:"pre"`
		Post map[common.Address]tracedPrestateAccount `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(raw, &diff))

	assertTraceBalance(t, diff.Post[rewardbase], (*big.Int)(&value))
}

func TestTraceCallPrestateTracerDiffModeCreate(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000cafe")
	created := crypto.CreateAddress(from, 0)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from:    {Balance: big.NewInt(5)},
		created: {Balance: big.NewInt(0x10)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	tracerName := "prestateTracer"
	gas := hexutil.Uint64(200000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	value := hexutil.Big(*big.NewInt(5))
	initCode := hexutil.Bytes(common.FromHex("0x303150600160005560006000526001601ff3"))
	config := &TraceConfig{Tracer: &tracerName, TracerConfig: json.RawMessage(`{"diffMode":true}`)}

	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		Gas:      &gas,
		GasPrice: &gasPrice,
		Value:    value,
		Data:     initCode,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	raw, ok := result.(json.RawMessage)
	if !assert.True(t, ok, "expected json.RawMessage result, got %T", result) {
		return
	}

	var diff struct {
		Pre  map[common.Address]tracedPrestateAccount `json:"pre"`
		Post map[common.Address]tracedPrestateAccount `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(raw, &diff))

	preCreated, ok := diff.Pre[created]
	if assert.True(t, ok, "created account missing from pre diff") {
		assertTraceBalance(t, preCreated, big.NewInt(0x10))
		assert.Equal(t, uint64(0), preCreated.Nonce)
		assert.NotContains(t, preCreated.Storage, common.Hash{})
	}
	postCreated, ok := diff.Post[created]
	if assert.True(t, ok, "created account missing from post diff") {
		assertTraceBalance(t, postCreated, big.NewInt(0x15))
		assert.Equal(t, uint64(1), postCreated.Nonce)
		assert.Equal(t, hexutil.Bytes{0}, postCreated.Code)
		assert.Equal(t, common.BigToHash(big.NewInt(1)), postCreated.Storage[common.Hash{}])
	}
}

func TestTraceCallPrestateTracerDiffModePrunesUnchangedStorage(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000d001")
	contract := common.HexToAddress("0x000000000000000000000000000000000000d002")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	tracerName := "prestateTracer"
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	code := hexutil.Bytes(common.FromHex("0x600054600160015500"))
	state := map[common.Hash]common.Hash{
		{}: common.BigToHash(big.NewInt(1)),
	}
	overrides := kaiaapi.EthStateOverride{
		contract: {Code: &code, State: &state},
	}
	config := &TraceConfig{Tracer: &tracerName, TracerConfig: json.RawMessage(`{"diffMode":true}`), StateOverrides: &overrides}

	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		To:       &contract,
		Gas:      &gas,
		GasPrice: &gasPrice,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	raw, ok := result.(json.RawMessage)
	if !assert.True(t, ok, "expected json.RawMessage result, got %T", result) {
		return
	}

	var diff struct {
		Pre  map[common.Address]tracedPrestateAccount `json:"pre"`
		Post map[common.Address]tracedPrestateAccount `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(raw, &diff))

	preContract := diff.Pre[contract]
	assert.NotContains(t, preContract.Storage, common.Hash{})
	assert.Equal(t, common.BigToHash(big.NewInt(1)), diff.Post[contract].Storage[common.BigToHash(big.NewInt(1))])
}

func TestTraceCallCallTracerWithLog(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x000000000000000000000000000000000000a222")
	contract := common.HexToAddress("0x000000000000000000000000000000000000b222")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	tracerName := "callTracer"
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	topic := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	code := hexutil.Bytes(common.FromHex("0x7f111111111111111111111111111111111111111111111111111111111111111160006000a100"))
	overrides := kaiaapi.EthStateOverride{
		contract: {Code: &code},
	}
	config := &TraceConfig{
		Tracer:         &tracerName,
		TracerConfig:   json.RawMessage(`{"withLog":true}`),
		StateOverrides: &overrides,
	}

	blockNumber := rpc.LatestBlockNumber
	result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		To:       &contract,
		Gas:      &gas,
		GasPrice: &gasPrice,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	frame, ok := result.(vm.CallFrame)
	if !assert.True(t, ok, "expected vm.CallFrame result, got %T", result) {
		return
	}
	if assert.Len(t, frame.Logs, 1) {
		assert.Equal(t, contract, frame.Logs[0].Address)
		assert.Equal(t, []common.Hash{topic}, frame.Logs[0].Topics)
		assert.Empty(t, frame.Logs[0].Data)
		assert.Equal(t, hexutil.Uint(0), frame.Logs[0].Index)
		assert.Equal(t, hexutil.Uint(0), frame.Logs[0].Position)
	}
	encoded, err := json.Marshal(frame)
	assert.NoError(t, err)
	assert.Contains(t, string(encoded), `"logs"`)
}

func TestTraceCallStructLoggerHonorsTimeout(t *testing.T) {
	from := common.HexToAddress("0x000000000000000000000000000000000000a333")
	contract := common.HexToAddress("0x000000000000000000000000000000000000b333")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	api := NewAPI(newTestBackend(t, 1, genesis, nil))

	gas := hexutil.Uint64(params.UpperGasLimit)
	gasPrice := hexutil.Big(*big.NewInt(0))
	code := hexutil.Bytes(common.FromHex("0x5b600056"))
	timeout := "1ns"
	overrides := kaiaapi.EthStateOverride{
		contract: {Code: &code},
	}
	config := &TraceConfig{
		LogConfig:      &vm.LogConfig{DisableMemory: true, DisableStack: true, DisableStorage: true, Limit: 1},
		Timeout:        &timeout,
		StateOverrides: &overrides,
	}

	blockNumber := rpc.LatestBlockNumber
	_, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
		From:     from,
		To:       &contract,
		Gas:      &gas,
		GasPrice: &gasPrice,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
	assert.ErrorContains(t, err, "tracing aborted")
}

func TestTraceCallStructLoggerSlotWait(t *testing.T) {
	from := common.HexToAddress("0xa666")
	to := common.HexToAddress("0xb666")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from: {Balance: big.NewInt(0)},
	}}
	backend := newTestBackend(t, 1, genesis, nil)
	t.Cleanup(backend.chain.Stop)
	api := NewAPI(backend)
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	blockNumber := rpc.LatestBlockNumber
	timeout := "1s"

	for _, tc := range []struct {
		name           string
		cancel         bool
		requestTimeout time.Duration
		wantErr        error
	}{
		{name: "execution timeout excludes queueing"},
		{name: "request cancellation", cancel: true, wantErr: context.Canceled},
		{name: "request deadline", requestTimeout: time.Second, wantErr: context.DeadlineExceeded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Keep the slot channel inside the fake-clock bubble. Do not run in parallel.
			slots := structTraceSlots
			defer func() { structTraceSlots = slots }()
			synctest.Test(t, func(t *testing.T) {
				structTraceSlots = make(chan struct{}, maxConcurrentStructTraces)
				acquired := 0
				defer func() {
					for range acquired {
						releaseStructTraceSlot()
					}
				}()
				for range maxConcurrentStructTraces {
					structTraceSlots <- struct{}{}
					acquired++
				}

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.requestTimeout > 0 {
					var cancelDeadline context.CancelFunc
					ctx, cancelDeadline = context.WithTimeout(ctx, tc.requestTimeout)
					defer cancelDeadline()
				}
				var (
					result interface{}
					err    error
				)
				done := make(chan struct{})
				go func() {
					result, err = api.TraceCall(ctx, kaiaapi.CallArgs{
						From: from, To: &to, Gas: &gas, GasPrice: &gasPrice,
					}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, &TraceConfig{Timeout: &timeout})
					close(done)
				}()

				synctest.Wait()
				time.Sleep(2 * time.Second)
				if tc.cancel {
					cancel()
				}
				synctest.Wait()
				if tc.wantErr == nil {
					select {
					case <-done:
						t.Fatalf("trace finished before acquiring a slot: %v", err)
					default:
					}
					releaseStructTraceSlot()
					acquired--
					synctest.Wait()
				}
				select {
				case <-done:
				default:
					t.Fatal("trace did not finish")
				}
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
					assert.ErrorContains(t, err, "tracing aborted")
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
				}
				assert.Len(t, structTraceSlots, acquired)
			})
		})
	}
}

func TestResolveTraceTimeout(t *testing.T) {
	short := "1s"
	long := "30s"
	invalid := "invalid"
	for _, tc := range []struct {
		name    string
		config  *TraceConfig
		want    time.Duration
		wantErr bool
	}{
		{name: "default", want: defaultTraceTimeout},
		{name: "no timeout", config: &TraceConfig{}, want: defaultTraceTimeout},
		{name: "shorter", config: &TraceConfig{Timeout: &short}, want: time.Second},
		{name: "longer", config: &TraceConfig{Timeout: &long}, want: 30 * time.Second},
		{name: "invalid", config: &TraceConfig{Timeout: &invalid}, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveTraceTimeout(tc.config)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTraceCallStructLoggerMemoryConfig(t *testing.T) {
	from := common.HexToAddress("0xa444")
	contract := common.HexToAddress("0xb444")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from:     {Balance: big.NewInt(0)},
		contract: {Balance: big.NewInt(0), Code: common.FromHex("0x600160005200")},
	}}
	backend := newTestBackend(t, 1, genesis, nil)
	t.Cleanup(backend.chain.Stop)
	api := NewAPI(backend)
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	blockNumber := rpc.LatestBlockNumber

	for _, tc := range []struct {
		name   string
		config string
		memory bool
	}{
		{"nil", `null`, true},
		{"empty", `{}`, true},
		{"omitted", `{"disableStack":true}`, true},
		{"enabled", `{"disableMemory":false}`, true},
		{"disabled", `{"disableMemory":true}`, false},
		{"legacy logger timeout", `{"loggerTimeout":"1ns"}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var config *TraceConfig
			if err := json.Unmarshal([]byte(tc.config), &config); err != nil {
				t.Fatal(err)
			}
			result, err := api.TraceCall(context.Background(), kaiaapi.CallArgs{
				From: from, To: &contract, Gas: &gas, GasPrice: &gasPrice,
			}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, config)
			if err != nil {
				t.Fatal(err)
			}
			var decoded kaiaapi.ExecutionResult
			if err := json.Unmarshal(result.(json.RawMessage), &decoded); err != nil {
				t.Fatal(err)
			}
			if !assert.NotEmpty(t, decoded.StructLogs) {
				return
			}
			last := decoded.StructLogs[len(decoded.StructLogs)-1]
			assert.Equal(t, "STOP", last.Op)
			if tc.memory {
				assert.Equal(t, &[]string{fmt.Sprintf("%064x", 1)}, last.Memory)
			} else {
				assert.Nil(t, last.Memory)
			}
		})
	}
}

func TestTraceCallNamedTracerKeepsCancellationBehavior(t *testing.T) {
	from := common.HexToAddress("0xa555")
	contract := common.HexToAddress("0xb555")
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		from:     {Balance: big.NewInt(0)},
		contract: {Balance: big.NewInt(0), Code: common.FromHex("0x600160005200")},
	}}
	backend := newTestBackend(t, 1, genesis, nil)
	t.Cleanup(backend.chain.Stop)
	api := NewAPI(backend)
	gas := hexutil.Uint64(100000)
	gasPrice := hexutil.Big(*big.NewInt(0))
	blockNumber := rpc.LatestBlockNumber
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for _, name := range []string{"callTracer", "fastCallTracer", "prestateTracer", "noopTracer"} {
		t.Run(name, func(t *testing.T) {
			result, err := api.TraceCall(ctx, kaiaapi.CallArgs{
				From: from, To: &contract, Gas: &gas, GasPrice: &gasPrice,
			}, rpc.BlockNumberOrHash{BlockNumber: &blockNumber}, &TraceConfig{Tracer: &name})
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestStructTraceSlotWaitHonorsCancellation(t *testing.T) {
	acquired := 0
	defer func() {
		for range acquired {
			releaseStructTraceSlot()
		}
	}()
	for range maxConcurrentStructTraces {
		if err := acquireStructTraceSlot(context.Background()); err != nil {
			t.Fatal(err)
		}
		acquired++
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.ErrorIs(t, acquireStructTraceSlot(ctx), context.Canceled)
}

func TestTraceTransaction(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(2)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.KAIA)},
		accounts[1].addr: {Balance: big.NewInt(params.KAIA)},
	}}
	target := common.Hash{}
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	api := NewAPI(newTestBackend(t, 1, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 kei
		//    fee:   0 kei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(1), nil), signer, accounts[0].key)
		b.AddTx(tx)
		target = tx.Hash()
	}))
	result, err := api.TraceTransaction(context.Background(), target, nil)
	if err != nil {
		t.Errorf("Failed to trace transaction %v", err)
	}
	if !reflect.DeepEqual(result, emptyStructTraceResult(params.TxGas)) {
		t.Error("Transaction tracing result is different")
	}
}

func TestTraceBlock(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.KAIA)},
		accounts[1].addr: {Balance: big.NewInt(params.KAIA)},
		accounts[2].addr: {Balance: big.NewInt(params.KAIA)},
	}}
	genBlocks := 10
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	api := NewAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 kei
		//    fee:   0 kei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
	}))

	testSuite := []struct {
		blockNumber rpc.BlockNumber
		config      *TraceConfig
		expect      interface{}
		expectErr   error
	}{
		// Trace genesis block, expect error
		{
			blockNumber: rpc.BlockNumber(0),
			config:      nil,
			expect:      nil,
			expectErr:   errors.New("genesis is not traceable"),
		},
		// Trace head block
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{Result: emptyStructTraceResult(params.TxGas)},
			},
		},
		// Trace non-existent block
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			config:      nil,
			expectErr:   fmt.Errorf("the block does not exist (block number: %d)", genBlocks+1),
			expect:      nil,
		},
		// Trace latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{Result: emptyStructTraceResult(params.TxGas)},
			},
		},
		// Trace pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{Result: emptyStructTraceResult(params.TxGas)},
			},
		},
	}
	for _, testspec := range testSuite {
		result, err := api.TraceBlockByNumber(context.Background(), testspec.blockNumber, testspec.config)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, testspec.expectErr) {
				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error, get %v", err)
				continue
			}
			if len(result) != len(testspec.expect.([]*txTraceResult)) {
				t.Errorf("Result length mismatch, want %v, get %v", len(result), len(testspec.expect.([]*txTraceResult)))
			}
			for idx, r := range result {
				if !reflect.DeepEqual(r.Result, testspec.expect.([]*txTraceResult)[idx].Result) {
					t.Errorf("Result mismatch, want %v, get %v", testspec.expect, result)
				}
			}
		}
	}
}

type Account struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

type Accounts []Account

func (a Accounts) Len() int           { return len(a) }
func (a Accounts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Accounts) Less(i, j int) bool { return bytes.Compare(a[i].addr.Bytes(), a[j].addr.Bytes()) < 0 }

func newAccounts(n int) (accounts Accounts) {
	for range n {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	sort.Sort(accounts)
	return accounts
}
