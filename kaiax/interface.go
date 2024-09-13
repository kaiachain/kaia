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
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
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

// JsonRpcModule adds JSON-RPC APIs.
type JsonRpcModule interface {
	// Exposes the module-specific APIs like governance_ namespace.
	// Or expose a general-purpose APIs like kaia_ namespace, in which case
	// the module is likely an API-only module.
	APIs() []rpc.API
}

// ConsensusModule consensus deals with states before block confirmation.
// Therefore these methods MUST NOT modify any persistent states.
// e.g. VerifyHeader shall not actually record the governance change.
type ConsensusModule interface {
	// Additional checks to perform at the end of header verification.
	VerifyHeader(header *types.Header) error

	// Additional changes to the new header that is being created.
	// Headers are modified in-place.
	PrepareHeader(header *types.Header) error

	// Additional state transitions after block txs have been executed.
	// Header modifications and state transitions are applied in-place.
	FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error
}

// Any component or module that accomodate consensus modules.
type ConsensusModuleHost interface {
	RegisterConsensusModule(modules ...ConsensusModule)
}

// ExecutionModule deals with execution of confirmed blocks.
// These methods MAY modify persistent states.
// e.g. PostInsertBlock() may record the governance change.
type ExecutionModule interface {
	// Additional actions to perform after inserting a block.
	PostInsertBlock(block *types.Block) error
}

// Any component or module that accomodate execution modules.
type ExecutionModuleHost interface {
	RegisterExecutionModule(modules ...ExecutionModule)
}

// RewindableModule is a module that can erase its persistent data
// associated to block numbers. If a module maintains block-number-dependent
// data, the module should implement RewindableModule.
type RewindableModule interface {
	// RewindableModule must implement Start() and Stop() methods so it can be reset after rewind.
	// In other words, non-Rewindable modules must be indifferent to block rewinds.
	BaseModule

	// Actions to be taken when the header rewinds to given block.
	// The block is the new head block and the new current block.
	// e.g. Move the head pointer to `num`.
	//
	// Happens when "soft" rewind occurs such as --start-block-num or startup repair.
	// and also when "hard" rewind occurs such as debug_setHead.
	RewindTo(newBlock *types.Block)

	// Actions to be taken when the block data should be deleted at given block hash and number.
	// This method should deal with exactly one block.
	// e.g. Delete data associated with the block number `num`.
	//
	// Happens when "hard" rewind occurs such as debug_setHead.
	RewindDelete(hash common.Hash, num uint64)
}

// Any component or module that accomodate rewindable modules.
type RewindableModuleHost interface {
	RegisterRewindableModule(modules ...RewindableModule)
	StopRewindableModules()
	StartRewindableModules()
}

// TxProcessModule intervenes how transactions are executed.
// Tx processing can happen on temporary states (e.g. eth_call),
// or run on states before confirmation (e.g. miner),
// therefore these methods MUST NOT modify any persistent states.
type TxProcessModule interface {
	// Additional actions to perform before EVM execution.
	// Optionally transform the tx to a regular Eth tx.
	// Optionally populate VM envs e.g. feepayer
	PreRunTx(evm *vm.EVM, tx *types.Transaction) (*types.Transaction, error)

	// Additional actions to perform after EVM execution.
	PostRunTx(evm *vm.EVM, tx *types.Transaction) error
}

// Any component or module that accomodate tx process modules.
type TxProcessModuleHost interface {
	RegisterTxProcessModule(modules ...TxProcessModule)
}

// TxPoolModule can intervene how the txpool handles transactions
// from the inception (e.g. AddLocal) to termination (e.g. drop).
type TxPoolModule interface {
	// Additional actions to be taken when a new tx arrives at txpool
	PreAddLocal(*types.Transaction) error
	PreAddRemote(*types.Transaction) error
}

// Any component or module that accomodate txpool modules.
type TxPoolModuleHost interface {
	RegisterTxPoolModule(modules ...TxPoolModule)
}

// A module can freely add more methods.
// But try to follow the naming convention:
//
// Prefix 'Get' for read-only methods.
//   GetXyz(..)
//
// Prefix 'Handle' for other non-getters not defined in the interface above.
//   HandleXyz(...)
