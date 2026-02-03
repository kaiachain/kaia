// Modifications Copyright 2024 The Kaia Authors
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

package shared

import (
	"errors"
	"math/big"
	"sync"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var logger = log.NewModuleLogger(log.ConsensusIstanbul)

var (
	ErrExecutorNotInitialized = errors.New("executor not initialized, call Reset first")
	ErrStateRootMismatch      = errors.New("state root mismatch")
)

// DefaultExecutor implements the consensus.Executor interface.
// It provides transaction execution functionality that can be used by consensus engines.
type DefaultExecutor struct {
	mu sync.RWMutex

	config   *params.ChainConfig
	chain    consensus.ChainContext
	signer   types.Signer
	nodeAddr common.Address

	// Current execution state
	state    *state.StateDB
	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
	logs     []*types.Log
	usedGas  uint64

	// initialized indicates whether Reset has been called
	initialized bool
}

// NewDefaultExecutor creates a new DefaultExecutor instance.
func NewDefaultExecutor(config *params.ChainConfig, chain consensus.ChainContext, nodeAddr common.Address) *DefaultExecutor {
	return &DefaultExecutor{
		config:   config,
		chain:    chain,
		nodeAddr: nodeAddr,
		signer:   types.MakeSigner(config, nil),
	}
}

// Reset initializes the executor for a new block, setting up the state
// based on the parent block.
func (e *DefaultExecutor) Reset(parent *types.Block) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Get state from parent block
	statedb, err := e.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}

	// Create header for the new block
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).Add(num, common.Big1),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       parent.Time(), // Will be updated later
	}

	// Handle BaseFee for Magma fork
	if e.config.IsMagmaForkEnabled(header.Number) {
		header.BaseFee = parent.Header().BaseFee
	}

	// Handle BlobGasUsed for Osaka fork
	if e.config.IsOsakaForkEnabled(header.Number) {
		blobGasUsed := uint64(0)
		header.BlobGasUsed = &blobGasUsed
	}

	e.state = statedb
	e.header = header
	e.txs = nil
	e.receipts = nil
	e.logs = nil
	e.usedGas = 0
	e.signer = types.MakeSigner(e.config, header.Number)
	e.initialized = true

	return nil
}

// ResetWithState initializes the executor with a pre-existing state and header.
// This is useful when the caller has already prepared the state and header.
func (e *DefaultExecutor) ResetWithState(statedb *state.StateDB, header *types.Header) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.state = statedb
	e.header = header
	e.txs = nil
	e.receipts = nil
	e.logs = nil
	e.usedGas = 0
	e.signer = types.MakeSigner(e.config, header.Number)
	e.initialized = true

	return nil
}

// Execute executes a batch of transactions and returns the execution result.
func (e *DefaultExecutor) Execute(txs types.Transactions) (*consensus.ExecutionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, ErrExecutorNotInitialized
	}

	vmConfig := &vm.Config{}

	for _, tx := range txs {
		e.state.SetTxContext(tx.Hash(), common.Hash{}, len(e.txs))

		receipt, logs, err := e.applyTransaction(tx, vmConfig)
		if err != nil {
			// Skip failed transactions but continue with others
			logger.Trace("Transaction execution failed", "hash", tx.Hash(), "err", err)
			continue
		}

		e.txs = append(e.txs, tx)
		e.receipts = append(e.receipts, receipt)
		e.logs = append(e.logs, logs...)
		e.usedGas = e.header.GasUsed
	}

	return e.buildResult(), nil
}

// GetPendingState returns the current accumulated state during execution.
func (e *DefaultExecutor) GetPendingState() *state.StateDB {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.state
}

// GetPendingReceipts returns the receipts accumulated during execution.
func (e *DefaultExecutor) GetPendingReceipts() types.Receipts {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.receipts
}

// VerifyStateRoot verifies that the current state root matches the expected value.
func (e *DefaultExecutor) VerifyStateRoot(expected common.Hash) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.initialized {
		return ErrExecutorNotInitialized
	}

	root := e.state.IntermediateRoot(true)
	if root != expected {
		return ErrStateRootMismatch
	}
	return nil
}

// applyTransaction applies a single transaction to the current state.
func (e *DefaultExecutor) applyTransaction(tx *types.Transaction, vmConfig *vm.Config) (*types.Receipt, []*types.Log, error) {
	snap := e.state.Snapshot()

	receipt, _, err := e.chain.ApplyTransaction(e.config, &e.nodeAddr, e.state, e.header, tx, &e.header.GasUsed, vmConfig)
	if err != nil {
		e.state.RevertToSnapshot(snap)
		return nil, nil, err
	}

	return receipt, receipt.Logs, nil
}

// buildResult builds and returns the current execution result.
func (e *DefaultExecutor) buildResult() *consensus.ExecutionResult {
	stateRoot := e.state.IntermediateRoot(true)

	return &consensus.ExecutionResult{
		StateRoot: stateRoot,
		State:     e.state,
		Receipts:  e.receipts,
		Logs:      e.logs,
		UsedGas:   e.usedGas,
		Txs:       e.txs,
	}
}

// GetHeader returns the current header being built.
// This is useful for callers that need header information during execution.
func (e *DefaultExecutor) GetHeader() *types.Header {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.header
}

// GetTransactions returns the transactions that have been executed.
func (e *DefaultExecutor) GetTransactions() []*types.Transaction {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.txs
}
