// Modifications Copyright 2026 The Kaia Authors
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

package tests

import (
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/profile"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/engine"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultExecutor_Clone verifies that Clone returns a fresh executor that
// shares immutable config but has independent mutable state.
func TestDefaultExecutor_Clone(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, nodeAddr, vm.Config{})

	// Clone before initialization — clone should also be uninitialized.
	clone := executor.Clone()
	require.NotNil(t, clone)

	// Initialize original — clone must remain uninitialized.
	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).Add(parent.Number(), common.Big1),
		Time:       new(big.Int).Add(parent.Time(), common.Big1),
		BaseFee:    parent.Header().BaseFee,
	}
	require.NoError(t, executor.ResetWithState(parentState, header))

	// Clone should fail ProcessBlock because it was never initialized.
	_, err = clone.ProcessBlock(nil)
	assert.ErrorIs(t, err, work.ErrExecutorNotInitialized)

	// Initialize clone separately — should succeed independently.
	cloneState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)
	require.NoError(t, clone.ResetWithState(cloneState, header))

	result, err := clone.ProcessBlock(nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), result.UsedGas)
}

// TestProcessBlock_MatchesInsertChain is the single canonical test for
// speculative execution correctness. It asserts that the speculative pipeline
// (ProcessBlock + FinalizeState) produces identical state root, receipt
// root, bloom, and gas-used as the proposer's MineABlock AND that InsertChain
// accepts the same block (confirming StateProcessor.Process also agrees).
func TestProcessBlock_MatchesInsertChain(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	// --- 1. Create value-transfer transactions ---
	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	var txs types.Transactions
	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]

	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)
	for i := range 5 {
		recipient := *bcdata.addrs[1+i%4]
		tx := types.NewTransaction(
			nonce,
			recipient,
			new(big.Int).SetUint64(1000),
			params.TxGas,
			gasPrice,
			nil,
		)
		signedTx, err := types.SignTx(tx, signer, sender)
		require.NoError(t, err)
		txs = append(txs, signedTx)
		nonce++
	}

	// --- 2. Mine a block (proposer path) ---
	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(txs, signer, prof, nil)
	require.NoError(t, err)
	require.Equal(t, len(txs), len(block.Transactions()), "all txs should be included")

	// --- 3. Speculative execution (validator path) ---
	// ProcessBlock delegates to Process() internally, which includes
	// FinalizeState. Do NOT call FinalizeState separately.
	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, common.Address{}, vm.Config{})
	require.NoError(t, executor.ResetWithState(parentState, block.Header()))

	specResult, err := executor.ProcessBlock(block.Transactions())
	require.NoError(t, err)

	// --- 4. Compare speculative results with block header ---
	// State root: the most critical invariant. Process() already called
	// FinalizeState which computes IntermediateRoot.
	specRoot := specResult.State.IntermediateRoot(true)
	assert.Equal(t, block.Root(), specRoot,
		"state root mismatch between proposer and speculative execution")

	// Receipt root.
	specReceiptRoot := types.DeriveReceiptsRoot(specResult.Receipts, block.Number())
	assert.Equal(t, block.ReceiptHash(), specReceiptRoot,
		"receipt root mismatch between proposer and speculative execution")

	// Bloom filter.
	specBloom := types.CreateBloom(specResult.Receipts)
	assert.Equal(t, block.Bloom(), specBloom,
		"bloom mismatch between proposer and speculative execution")

	// Gas used.
	assert.Equal(t, block.GasUsed(), specResult.UsedGas,
		"gas used mismatch between proposer and speculative execution")

	// Receipt count.
	assert.Equal(t, len(block.Transactions()), len(specResult.Receipts),
		"receipt count mismatch")

	// All receipts successful.
	for i, receipt := range specResult.Receipts {
		assert.Equal(t, uint(1), receipt.Status,
			"receipt %d should succeed", i)
	}

	// --- 5. InsertChain: confirms Process path also matches ---
	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err, "InsertChain should accept the block")
	assert.Equal(t, 0, n, "InsertChain should process the block at index 0")
}

// TestInsertChain_SpeculativeCacheHit verifies that InsertChain uses a
// correctly populated speculative cache, skipping Process().
func TestInsertChain_SpeculativeCacheHit(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	// --- 1. Create and mine a block ---
	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	var txs types.Transactions
	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	for i := range 3 {
		tx := types.NewTransaction(nonce, *bcdata.addrs[1+i%4], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
		signedTx, err := types.SignTx(tx, signer, sender)
		require.NoError(t, err)
		txs = append(txs, signedTx)
		nonce++
	}

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(txs, signer, prof, nil)
	require.NoError(t, err)

	// --- 2. Speculatively execute (validator path) ---
	// Uses the same ProcessBlock workflow that KBFT-2 will use.
	// ProcessBlock delegates to Process() internally.
	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, common.Address{}, vm.Config{})
	require.NoError(t, executor.ResetWithState(parentState, block.Header()))

	specResult, err := executor.ProcessBlock(block.Transactions())
	require.NoError(t, err)

	// --- 3. Populate the speculative cache ---
	hitsBefore := bcdata.bc.SpeculativeCache().Hits()
	entry := bcdata.bc.SpeculativeCache().Reserve(block.Hash())
	entry.Complete(&blockchain.SpeculativeResult{
		State:            specResult.State,
		Receipts:         specResult.Receipts,
		Logs:             specResult.Logs,
		UsedGas:          specResult.UsedGas,
		InternalTxTraces: specResult.InternalTxTraces,
		ProcessStats: blockchain.ProcessStats{
			BeforeApplyTxs: specResult.BeforeApplyTxs,
			AfterApplyTxs:  specResult.AfterApplyTxs,
			AfterFinalize:  specResult.AfterFinalize,
		},
		Bloom:       types.CreateBloom(specResult.Receipts),
		ReceiptHash: types.DeriveReceiptsRoot(specResult.Receipts, block.Number()),
	}, nil)

	// --- 4. InsertChain should use the cached result ---
	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, hitsBefore+1, bcdata.bc.SpeculativeCache().Hits(), "cache should have been hit")
}

// TestInsertChain_SpeculativeCacheMiss verifies that InsertChain falls through
// to synchronous execution when the cache is empty.
func TestInsertChain_SpeculativeCacheMiss(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	tx := types.NewTransaction(nonce, *bcdata.addrs[1], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
	signedTx, err := types.SignTx(tx, signer, sender)
	require.NoError(t, err)

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(types.Transactions{signedTx}, signer, prof, nil)
	require.NoError(t, err)

	// Cache is empty — InsertChain must succeed via sync path.
	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, uint64(0), bcdata.bc.SpeculativeCache().Hits())
}

// TestInsertChain_SpeculativeCacheHashMismatch verifies that a cache entry
// for a different block hash is ignored and InsertChain falls through to sync.
func TestInsertChain_SpeculativeCacheHashMismatch(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	tx := types.NewTransaction(nonce, *bcdata.addrs[1], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
	signedTx, err := types.SignTx(tx, signer, sender)
	require.NoError(t, err)

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(types.Transactions{signedTx}, signer, prof, nil)
	require.NoError(t, err)

	// Populate cache with a DIFFERENT block hash.
	entry := bcdata.bc.SpeculativeCache().Reserve(common.HexToHash("0xdeadbeef"))
	entry.Complete(&blockchain.SpeculativeResult{UsedGas: 999}, nil)

	// InsertChain should ignore the mismatched entry and succeed via sync.
	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, uint64(0), bcdata.bc.SpeculativeCache().Hits())
}

// TestInsertChain_SpeculativeCacheWait verifies that InsertChain blocks on an
// in-flight speculative entry instead of falling through to synchronous
// execution — the behaviour enabled by the non-zero TryGet timeout.
func TestInsertChain_SpeculativeCacheWait(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	tx := types.NewTransaction(nonce, *bcdata.addrs[1], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
	signedTx, err := types.SignTx(tx, signer, sender)
	require.NoError(t, err)

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(types.Transactions{signedTx}, signer, prof, nil)
	require.NoError(t, err)

	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, common.Address{}, vm.Config{})
	require.NoError(t, executor.ResetWithState(parentState, block.Header()))

	specResult, err := executor.ProcessBlock(block.Transactions())
	require.NoError(t, err)

	hitsBefore := bcdata.bc.SpeculativeCache().Hits()
	entry := bcdata.bc.SpeculativeCache().Reserve(block.Hash())
	go func() {
		time.Sleep(50 * time.Millisecond)
		entry.Complete(&blockchain.SpeculativeResult{
			State:            specResult.State,
			Receipts:         specResult.Receipts,
			Logs:             specResult.Logs,
			UsedGas:          specResult.UsedGas,
			InternalTxTraces: specResult.InternalTxTraces,
			ProcessStats: blockchain.ProcessStats{
				BeforeApplyTxs: specResult.BeforeApplyTxs,
				AfterApplyTxs:  specResult.AfterApplyTxs,
				AfterFinalize:  specResult.AfterFinalize,
			},
			Bloom:       types.CreateBloom(specResult.Receipts),
			ReceiptHash: types.DeriveReceiptsRoot(specResult.Receipts, block.Number()),
		}, nil)
	}()

	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, hitsBefore+1, bcdata.bc.SpeculativeCache().Hits(),
		"InsertChain should have waited for the in-flight speculative result")
}

// TestInsertChain_SpeculativeCacheTimeout verifies that InsertChain falls back
// to synchronous execution when the speculative entry stays in-flight past
// the timeout.
func TestInsertChain_SpeculativeCacheTimeout(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	tx := types.NewTransaction(nonce, *bcdata.addrs[1], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
	signedTx, err := types.SignTx(tx, signer, sender)
	require.NoError(t, err)

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(types.Transactions{signedTx}, signer, prof, nil)
	require.NoError(t, err)

	original := params.BlockGenerationTimeLimit
	params.BlockGenerationTimeLimit = 10 * time.Millisecond
	defer func() { params.BlockGenerationTimeLimit = original }()

	hitsBefore := bcdata.bc.SpeculativeCache().Hits()
	bcdata.bc.SpeculativeCache().Reserve(block.Hash())

	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, hitsBefore, bcdata.bc.SpeculativeCache().Hits(),
		"InsertChain should have timed out and fallen back to sync execution")
}

// TestKaiaBFT_SpeculativeExecution_ClonePath verifies that the kaiabft
// speculative execution path (Clone → ResetWithState → ProcessBlock →
// FinalizeState → populate cache) produces results that pass InsertChain's
// ValidateState, matching the proposer's block exactly.
func TestKaiaBFT_SpeculativeExecution_ClonePath(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	bcdata, err := NewBCDataWithConfigs(6, 4, Forks["Magma"], nil)
	require.NoError(t, err)
	defer bcdata.Shutdown()

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).Add(bcdata.bc.CurrentBlock().Header().BaseFee, big.NewInt(1))

	var txs types.Transactions
	sender := bcdata.privKeys[0]
	senderAddr := *bcdata.addrs[0]
	stateDb, err := bcdata.bc.State()
	require.NoError(t, err)
	nonce := stateDb.GetNonce(senderAddr)

	for i := range 5 {
		tx := types.NewTransaction(nonce, *bcdata.addrs[1+i%4], new(big.Int).SetUint64(1000), params.TxGas, gasPrice, nil)
		signedTx, err := types.SignTx(tx, signer, sender)
		require.NoError(t, err)
		txs = append(txs, signedTx)
		nonce++
	}

	prof := profile.NewProfiler()
	_, block, _, err := bcdata.MineABlock(txs, signer, prof, nil)
	require.NoError(t, err)

	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	origExecutor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, common.Address{}, vm.Config{})
	clonedExecutor := origExecutor.Clone()

	require.NoError(t, clonedExecutor.ResetWithState(parentState, block.Header()))

	specResult, err := clonedExecutor.ProcessBlock(block.Transactions())
	require.NoError(t, err)

	assert.Equal(t, block.Root(), specResult.State.IntermediateRoot(true), "state root mismatch")
	assert.Equal(t, block.GasUsed(), specResult.UsedGas, "gas used mismatch")

	hitsBefore := bcdata.bc.SpeculativeCache().Hits()
	entry := bcdata.bc.SpeculativeCache().Reserve(block.Hash())

	entry.Complete(&blockchain.SpeculativeResult{
		State:            specResult.State,
		Receipts:         specResult.Receipts,
		Logs:             specResult.Logs,
		UsedGas:          specResult.UsedGas,
		InternalTxTraces: specResult.InternalTxTraces,
		ProcessStats: blockchain.ProcessStats{
			BeforeApplyTxs: specResult.BeforeApplyTxs,
			AfterApplyTxs:  specResult.AfterApplyTxs,
			AfterFinalize:  specResult.AfterFinalize,
		},
		Bloom:       types.CreateBloom(specResult.Receipts),
		ReceiptHash: types.DeriveReceiptsRoot(specResult.Receipts, block.Number()),
	}, nil)

	n, err := bcdata.bc.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, hitsBefore+1, bcdata.bc.SpeculativeCache().Hits(), "speculative cache should have been hit")
}

// TestKaiaBFT_EngineFactory verifies that the consensus engine factory
// creates the correct engine type based on EngineType.
func TestKaiaBFT_EngineFactory(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	istanbulConfig := &istanbul.Config{
		Timeout:        10000,
		ProposerPolicy: istanbul.RoundRobin,
		Epoch:          30000,
		SubGroupSize:   21,
	}

	e1 := engine.NewEngine(&engine.Opts{
		IstanbulConfig: istanbulConfig,
		PrivateKey:     key,
		NodeType:       common.CONSENSUSNODE,
	})
	require.NotNil(t, e1)
	_, isHandler := e1.(consensus.Handler)
	assert.True(t, isHandler, "istanbul engine should implement Handler")

	e2 := engine.NewEngine(&engine.Opts{
		EngineType:     engine.KaiaBFT,
		IstanbulConfig: istanbulConfig,
		PrivateKey:     key,
		NodeType:       common.CONSENSUSNODE,
	})
	require.NotNil(t, e2)
	_, isHandler2 := e2.(consensus.Handler)
	assert.True(t, isHandler2, "kaiabft engine should implement Handler")
}

// TestKaiaBFT_BlockGenerationTimeLimit verifies that the kaiabft-specific
// block generation time limit constant exists and that the params global
// can be set to it. The actual CLI config path (which sets the global
// based on --consensus.engine) is tested by integration tests.
func TestKaiaBFT_BlockGenerationTimeLimit(t *testing.T) {
	// Verify the constants are defined with expected values.
	assert.Equal(t, 250*time.Millisecond, params.DefaultBlockGenerationTimeLimit)
	assert.Equal(t, 500*time.Millisecond, params.KaiaBFTBlockGenerationTimeLimit)

	// Verify the global can be set to the kaiabft value and restored.
	original := params.BlockGenerationTimeLimit
	defer func() { params.BlockGenerationTimeLimit = original }()

	params.BlockGenerationTimeLimit = params.KaiaBFTBlockGenerationTimeLimit
	assert.Equal(t, 500*time.Millisecond, params.BlockGenerationTimeLimit)
}
