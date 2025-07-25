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
	"github.com/kaiachain/kaia/work/builder"
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

// Any component or module that accommodate consensus modules.
type ConsensusModuleHost interface {
	RegisterConsensusModule(modules ...ConsensusModule)

	// Temporal method to unregister consensus module.
	// See node/cn/backend.go#createSnapshot() for details.
	UnregisterConsensusModule(module ConsensusModule)
}

// ExecutionModule deals with execution of confirmed blocks.
// These methods MAY modify persistent states.
// e.g. PostInsertBlock() may record the governance change.
type ExecutionModule interface {
	// Additional actions to perform after inserting a block.
	PostInsertBlock(block *types.Block) error
}

// Any component or module that accommodate execution modules.
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

// Any component or module that accommodate rewindable modules.
type RewindableModuleHost interface {
	RegisterRewindableModule(modules ...RewindableModule)
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

// Any component or module that accommodate tx process modules.
type TxProcessModuleHost interface {
	RegisterTxProcessModule(modules ...TxProcessModule)
}

// TxPoolModule can intervene how the txpool handles transactions
// from the inception (e.g. AddLocal) to termination (e.g. drop).
//
//go:generate mockgen -destination=./mock/tx_pool_module.go -package=mock github.com/kaiachain/kaia/kaiax TxPoolModule
type TxPoolModule interface {
	// Additional actions to be taken when a new tx arrives at txpool
	PreAddTx(tx *types.Transaction, local bool) error

	// Additional checks to check if a given transaction should be handled by module.
	IsModuleTx(tx *types.Transaction) bool

	// Optional actions to check if sender balance is valid for module transaction.
	// This is mainly used on checking if module transaction be appended to queue.
	// If nil is returned, default check (balance > txFee) is performed. Otherwise, the returned function overrides default check.
	GetCheckBalance() func(tx *types.Transaction) error

	// Additional actions to check if a module transaction should be appended to pending
	IsReady(txs map[uint64]*types.Transaction, next uint64, ready types.Transactions) bool

	// Additional actions to perform before the txpool is reset.
	PreReset(oldHead, newHead *types.Header) (dropTxs []common.Hash)

	// Additional actions to perform after the txpool is reset.
	PostReset(oldHead, newHead *types.Header, queue, pending map[common.Address]types.Transactions)
}

// Any component or module that accommodate txpool modules.
type TxPoolModuleHost interface {
	RegisterTxPoolModule(modules ...TxPoolModule)
}

// TxBundlingModule can intervene how miner/proposer orders transactions in a block.
//
//go:generate mockgen -destination=./mock/tx_bundling_module.go -package=mock github.com/kaiachain/kaia/work/kaiax.TxBundlingModule
type TxBundlingModule interface {
	// The function finds transactions to be bundled.
	// New transactions can be injected.
	// returned bundles must not have conflict with `prevBundles`.
	// `txs` and `prevBundles` is read-only; it is only to check if there's conflict between new bundles.
	ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle

	// IsBundleTx returns true if the tx is a potential bundle tx.
	IsBundleTx(tx *types.Transaction) bool

	// GetMaxBundleTxsInPending returns the maximum number of transactions that can be bundled in pending.
	// This limitation works properly only when a module bundles only sequential txs by the same sender.
	GetMaxBundleTxsInPending() uint

	// GetMaxBundleTxsInQueue returns the maximum number of transactions that can be bundled in queue.
	// This limitation works properly only when a module bundles only sequential txs by the same sender.
	GetMaxBundleTxsInQueue() uint

	// FilterTxs filters transactions that are not bundled.
	FilterTxs(txs map[common.Address]types.Transactions)
}

// Any component or module that accomodate tx bundling modules.
type TxBundlingModuleHost interface {
	RegisterTxBundlingModule(modules ...TxBundlingModule)
}

// A module can freely add more methods.
// But try to follow the naming convention:
//
// Prefix 'Get' for read-only methods.
//   GetXyz(..)
//
// Prefix 'Handle' for other non-getters not defined in the interface above.
//   HandleXyz(...)
