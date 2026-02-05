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

package work

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

// clearSkippedTxs clears the skippedTxs map for test isolation.
func clearSkippedTxs() {
	skippedTxsMu.Lock()
	defer skippedTxsMu.Unlock()
	skippedTxs = make(map[common.Hash]time.Time)
}

func generateTestKey() *ecdsa.PrivateKey {
	key, _ := crypto.GenerateKey()
	return key
}

func createTestTransaction(nonce uint64, gasLimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(
		types.NewTransaction(nonce, common.HexToAddress("0xAAAA"), big.NewInt(100), gasLimit, big.NewInt(1), nil),
		types.LatestSignerForChainID(params.TestChainConfig.ChainID),
		key,
	)
	return tx
}

func setupTestTask(t *testing.T, key *ecdsa.PrivateKey) (*Task, *state.StateDB) {
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	statedb.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).SetUint64(params.KAIA))

	header := &types.Header{
		ParentHash:  common.Hash{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		BlockScore:  big.NewInt(1),
		Number:      big.NewInt(1),
		GasUsed:     0,
		Time:        big.NewInt(0),
		TimeFoS:     0,
		BaseFee:     big.NewInt(1),
	}

	task := NewTask(params.TestChainConfig, types.LatestSignerForChainID(params.TestChainConfig.ChainID), statedb, header)
	return task, statedb
}

// TestApplyTransactions_FirstTxExceedsGasLimit tests that the first transaction
// exceeding gas limit is added to skippedTxs and marked as unexecutable on retry.
func TestApplyTransactions_FirstTxExceedsGasLimit(t *testing.T) {
	clearSkippedTxs()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := generateTestKey()
	task, _ := setupTestTask(t, key)
	mockBC := mocks.NewMockBlockChain(ctrl)

	// Create a transaction
	tx := createTestTransaction(0, 50_000_000, key)

	// Mock ApplyTransaction to return ErrFirstTxLimitReached
	mockBC.EXPECT().
		ApplyTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(tx), gomock.Any(), gomock.Any()).
		Return(nil, nil, vm.ErrFirstTxLimitReached)

	// Create transaction set
	pending := map[common.Address]types.Transactions{
		crypto.PubkeyToAddress(key.PublicKey): {tx},
	}
	txs := types.NewTransactionsByPriceAndNonce(task.signer, pending, task.header.BaseFee)

	// Apply transactions with timeout
	done := make(chan struct{})
	go func() {
		task.ApplyTransactions(txs, mockBC, common.Address{}, nil)
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("ApplyTransactions timed out")
	}

	// Verify tx was added to skippedTxs
	skippedTxsMu.Lock()
	_, exists := skippedTxs[tx.Hash()]
	skippedTxsMu.Unlock()
	assert.True(t, exists, "tx should be in skippedTxs after exceeding gas limit")

	// Verify tx was not included in block
	assert.Equal(t, 0, len(task.Transactions()), "tx should not be included in block")
}

// TestApplyTransactions_SkippedTxMarkedUnexecutable tests that a previously skipped
// transaction is marked as unexecutable when encountered again.
func TestApplyTransactions_SkippedTxMarkedUnexecutable(t *testing.T) {
	clearSkippedTxs()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := generateTestKey()
	task, _ := setupTestTask(t, key)
	mockBC := mocks.NewMockBlockChain(ctrl)

	// Create a transaction
	tx := createTestTransaction(0, 50_000_000, key)

	// Pre-add tx to skippedTxs (simulating previous gas limit violation)
	addSkippedTx(tx.Hash())

	// No ApplyTransaction call expected - tx should be skipped before execution

	// Create transaction set
	pending := map[common.Address]types.Transactions{
		crypto.PubkeyToAddress(key.PublicKey): {tx},
	}
	txs := types.NewTransactionsByPriceAndNonce(task.signer, pending, task.header.BaseFee)

	// Apply transactions with timeout
	done := make(chan struct{})
	go func() {
		task.ApplyTransactions(txs, mockBC, common.Address{}, nil)
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("ApplyTransactions timed out - skipped tx might be causing infinite loop")
	}

	// Verify tx was marked as unexecutable
	assert.True(t, tx.IsMarkedUnexecutable(), "skipped tx should be marked as unexecutable")

	// Verify tx was not included in block
	assert.Equal(t, 0, len(task.Transactions()), "skipped tx should not be included in block")
}

// TestApplyTransactions_FirstTxNormal tests that a normal first transaction
// is executed and included in the block.
func TestApplyTransactions_FirstTxNormal(t *testing.T) {
	clearSkippedTxs()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := generateTestKey()
	task, _ := setupTestTask(t, key)
	mockBC := mocks.NewMockBlockChain(ctrl)

	// Create a transaction
	tx := createTestTransaction(0, 21000, key)

	// Mock ApplyTransaction to return success
	receipt := &types.Receipt{Status: types.ReceiptStatusSuccessful}
	mockBC.EXPECT().
		ApplyTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(tx), gomock.Any(), gomock.Any()).
		Return(receipt, nil, nil)

	// Create transaction set
	pending := map[common.Address]types.Transactions{
		crypto.PubkeyToAddress(key.PublicKey): {tx},
	}
	txs := types.NewTransactionsByPriceAndNonce(task.signer, pending, task.header.BaseFee)

	// Apply transactions
	task.ApplyTransactions(txs, mockBC, common.Address{}, nil)

	// Verify tx was included in block
	assert.Equal(t, 1, len(task.Transactions()), "tx should be included in block")

	// Verify tx is not in skippedTxs
	skippedTxsMu.Lock()
	_, exists := skippedTxs[tx.Hash()]
	skippedTxsMu.Unlock()
	assert.False(t, exists, "normal tx should not be in skippedTxs")
}

// TestApplyTransactions_SecondTxHighGas tests that the second transaction
// is NOT subject to first tx gas limit (only 250ms time limit applies).
func TestApplyTransactions_SecondTxHighGas(t *testing.T) {
	clearSkippedTxs()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := generateTestKey()
	task, _ := setupTestTask(t, key)
	mockBC := mocks.NewMockBlockChain(ctrl)

	// Create two transactions
	tx1 := createTestTransaction(0, 21000, key)      // First tx - low gas
	tx2 := createTestTransaction(1, 50_000_000, key) // Second tx - high gas

	// Mock ApplyTransaction for first tx - success
	receipt1 := &types.Receipt{Status: types.ReceiptStatusSuccessful}
	mockBC.EXPECT().
		ApplyTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(tx1), gomock.Any(), gomock.Any()).
		Return(receipt1, nil, nil)

	// Mock ApplyTransaction for second tx - success (no gas limit check for non-first tx)
	receipt2 := &types.Receipt{Status: types.ReceiptStatusSuccessful}
	mockBC.EXPECT().
		ApplyTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(tx2), gomock.Any(), gomock.Any()).
		Return(receipt2, nil, nil)

	// Create transaction set
	pending := map[common.Address]types.Transactions{
		crypto.PubkeyToAddress(key.PublicKey): {tx1, tx2},
	}
	txs := types.NewTransactionsByPriceAndNonce(task.signer, pending, task.header.BaseFee)

	// Apply transactions
	task.ApplyTransactions(txs, mockBC, common.Address{}, nil)

	// Verify both txs were included in block
	assert.Equal(t, 2, len(task.Transactions()), "both txs should be included in block")

	// Verify tx2 is not in skippedTxs (second tx doesn't get added to skip list)
	skippedTxsMu.Lock()
	_, exists := skippedTxs[tx2.Hash()]
	skippedTxsMu.Unlock()
	assert.False(t, exists, "second tx should not be in skippedTxs even with high gas")
}
