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

package consensus

import (
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

// ExecutionResult contains the results of executing a batch of transactions.
type ExecutionResult struct {
	// StateRoot is the root hash of the state trie after execution
	StateRoot common.Hash
	// State is the state database after execution
	State *state.StateDB
	// Receipts contains the receipts of all executed transactions
	Receipts types.Receipts
	// Logs contains all logs emitted during execution
	Logs []*types.Log
	// UsedGas is the total gas used by all transactions
	UsedGas uint64
	// Txs contains the executed transactions
	Txs []*types.Transaction
	// Block is the finalized block (set after Finalize is called)
	Block *types.Block

	// Timing information
	ExecuteTime  time.Duration // Time spent executing transactions
	FinalizeTime time.Duration // Time spent finalizing the block
	SealTime     time.Duration // Time spent in consensus (Seal)
}

// Transactions returns the executed transactions.
func (r *ExecutionResult) Transactions() []*types.Transaction {
	return r.Txs
}

// Executor is responsible for executing transactions and managing state transitions.
// It separates the transaction execution logic from the worker, allowing for cleaner
// architecture and better testability.
//
//go:generate mockgen -destination=./mocks/executor_mock.go -package=mocks github.com/kaiachain/kaia/consensus Executor
type Executor interface {
	// Execute executes a batch of transactions and returns the execution result.
	Execute(txs types.Transactions) (*ExecutionResult, error)
}
