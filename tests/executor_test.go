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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/profile"
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

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, nodeAddr)

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
	}
	require.NoError(t, executor.ResetWithState(parentState, header))

	// Clone should fail ExecuteTransactions because it was never initialized.
	_, err = clone.ExecuteTransactions(nil)
	assert.ErrorIs(t, err, work.ErrExecutorNotInitialized)

	// Initialize clone separately — should succeed independently.
	cloneState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)
	require.NoError(t, clone.ResetWithState(cloneState, header))

	result, err := clone.ExecuteTransactions(nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), result.UsedGas)
}

// TestExecuteTransactions_MatchesInsertChain is the single canonical test for
// speculative execution correctness. It asserts that the speculative pipeline
// (ExecuteTransactions + FinalizeState) produces identical state root, receipt
// root, bloom, and gas-used as the proposer's MineABlock AND that InsertChain
// accepts the same block (confirming StateProcessor.Process also agrees).
func TestExecuteTransactions_MatchesInsertChain(t *testing.T) {
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
	parent := bcdata.bc.CurrentBlock()
	parentState, err := bcdata.bc.PrunableStateAt(parent.Root(), parent.NumberU64())
	require.NoError(t, err)

	executor := work.NewDefaultExecutor(bcdata.bc.Config(), bcdata.bc, common.Address{})
	require.NoError(t, executor.ResetWithState(parentState, block.Header()))

	specResult, err := executor.ExecuteTransactions(block.Transactions())
	require.NoError(t, err)

	specBlock, err := executor.FinalizeState(specResult)
	require.NoError(t, err)

	// --- 4. Compare speculative results with block header ---
	// State root: the most critical invariant.
	assert.Equal(t, block.Root(), specBlock.Root(),
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
