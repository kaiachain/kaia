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

package kaiax

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/networks/rpc"
)

// On top of below Module interfaces, every module must have these:
//
// Constructor: Initialize data structures such as maps and channels.
// Try not to perform any application-specific initializations here.
//   func NewAbc() (*Abc, error)
//
// Module-specific Init method: Pass in fields - config, other modules, db, etc.
//   type AbcOpts struct { ... }
//   func (*Abc) Init(*AbcOpts) error

// BaseModule must be implemented by every kaiax module
type BaseModule interface {
	Start() error
	Stop()
}

// A module can optionally implement below interfaces.

type JsonRpcModule interface {
	// Exposes the module-specific APIs like governance_ namespace.
	// Or you can implement a generic namespace like eth_ as a module.
	APIs() []rpc.API
}

// ConsensusModule consensus deals with states before block confirmation.
// Therefore these methods MUST NOT modify any persistent states.
// e.g. VerifyHeader shall not actually record the governance change.
type ConsensusModule interface {
	// Additional checks to perform during header verification
	VerifyHeader(*types.Header) error

	// Additional changes to the new header that is being created
	PrepareHeader(*types.Header) (*types.Header, error)

	// Additional state transitions after block txs have been executed
	FinalizeBlock() (*types.Block, error)
}

// ExecutionModule deals with execution of confirmed blocks.
// These methods MAY modify persistent states.
// e.g. PostInsertBlock() may record the governance change.
type ExecutionModule interface {
	// Additional actions to perform after inserting a block
	PostInsertBlock(*types.Block) error
}

// TxProcessModule intervenes how transactions are executed.
// Tx processing can happen on temporary states (e.g. eth_call),
// or run on states before confirmation (e.g. miner),
// therefore these methods MUST NOT modify any persistent states.
type TxProcessModule interface {
	// Additional actions to perform before EVM execution.
	// Optionally transform the tx to a regular Eth tx.
	// Optionally populate VM envs e.g. feepayer
	PreRunTx(*vm.EVM, *types.Transaction) (*types.Transaction, error)

	// Additional actions to perform after EVM execution.
	PostRunTx(*vm.EVM, *types.Transaction) error
}

// UnwindableModule is a module that can erase its persistent data
// associated to block numbers. If a module maintains block-number-dependent
// data, the module should implement UnwindableModule.
type UnwindableModule interface {
	// Actions to be taken when the head block is rewinded by one block, to `num-1`.
	// i.e. Delete data associated with the block number `num`.
	Unwind(num uint64) error
}

// TxPoolModule can intervene how the txpool handles transactions
// from the inception (e.g. AddLocal) to termination (e.g. drop).
type TxPoolModule interface {
	// Additional actions to be taken when a new tx arrives at txpool
	PreAddLocal(*types.Transaction) error
	PreAddRemote(*types.Transaction) error
}

// A module can freely add more methods.
// But try to follow the naming convention:
//
// Prefix 'Get' for read-only methods,
//   GetXyz(..)
//   e.g. GetConsensusInfo(...) (*ConsensusInfo, error)
// Prefix 'Handle' for methods that modifies any in-memory or persistent data
// But try not to add too many Handle method. Prefer above interfaces instead.
//   HandleXyz(...)
