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

package work

import (
	"errors"
	"sync"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/builder"
)

var executorLogger = log.NewModuleLogger(log.ConsensusIstanbul)

var (
	ErrExecutorNotInitialized = errors.New("executor not initialized, call Reset first")
	ErrStateRootMismatch      = errors.New("state root mismatch")
)

// DefaultExecutor implements the consensus.Executor interface.
// It provides transaction execution functionality that can be used by consensus engines.
type DefaultExecutor struct {
	mu sync.RWMutex

	config   *params.ChainConfig
	chain    BlockChain
	signer   types.Signer
	nodeAddr common.Address

	// Current execution state
	state    *state.StateDB
	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
	logs     []*types.Log
	usedGas  uint64

	// Transaction bundling modules for gasless transactions
	txBundlingModules []builder.TxBundlingModule

	// initialized indicates whether Reset has been called
	initialized bool
}

// NewDefaultExecutor creates a new DefaultExecutor instance.
func NewDefaultExecutor(config *params.ChainConfig, chain BlockChain, nodeAddr common.Address) *DefaultExecutor {
	return &DefaultExecutor{
		config:   config,
		chain:    chain,
		nodeAddr: nodeAddr,
		signer:   types.MakeSigner(config, nil),
	}
}

// SetTxBundlingModules sets the transaction bundling modules for gasless transactions.
func (e *DefaultExecutor) SetTxBundlingModules(modules []builder.TxBundlingModule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.txBundlingModules = modules
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
// This uses work.Task.CommitTransactions for full transaction execution with bundle handling.
func (e *DefaultExecutor) Execute(txs *types.TransactionsByPriceAndNonce, mux *event.TypeMux) (*consensus.ExecutionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, ErrExecutorNotInitialized
	}

	// Create Task for transaction execution with bundle handling and time limits
	task := NewTask(e.config, e.signer, e.state, e.header)

	// Execute transactions using CommitTransactions (includes posting pending events)
	task.CommitTransactions(mux, txs, e.chain, e.nodeAddr, e.txBundlingModules)

	// Get results from task
	e.txs = task.Transactions()
	e.receipts = task.Receipts()
	e.usedGas = e.header.GasUsed

	return e.buildResult(), nil
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
