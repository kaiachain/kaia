// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/types.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
)

// Validator is an interface which defines the standard for block validation. It
// is only responsible for validating block contents, as the header validation is
// done by the specific consensus engines.
type Validator interface {
	// RegisterHeaderModules registers header modules used during header validation.
	RegisterHeaderModules(modules ...kaiax.HeaderModule)

	// SetupKaiaxModules sets up Kaiax modules used during validation, such as the governance module.
	SetupKaiaxModules(mGov gov.GovModule)

	// Preprocess preprocesses the given headers concurrently.
	Preprocess(headers []*types.Header) (chan<- struct{}, <-chan error)

	// ValidateHeader validates or preprocesses the given header.
	ValidateHeader(header *types.Header) error

	// ValidateBody validates the given block's content.
	ValidateBody(block *types.Block) error

	// ValidateState validates the given statedb and optionally the receipts and
	// gas used.
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error
}

// Prefetcher is an interface for pre-caching transaction signatures and state.
type Prefetcher interface {
	// Prefetch processes the state changes according to the Kaia rules by running
	// the transaction messages using the statedb, but any changes are discarded. The
	// only goal is to pre-cache transaction signatures and state trie nodes.
	Prefetch(block *types.Block, stateDB *state.StateDB, cfg vm.Config, interrupt *uint32)
	PrefetchTx(block *types.Block, ti int, stateDB *state.StateDB, cfg vm.Config, interrupt *uint32)
}

// Processor is an interface for processing blocks using a given initial state.
type Processor interface {
	// InitializeState runs pre-transaction state modifications.
	InitializeState(header *types.Header, stateDB *state.StateDB)

	// RegisterBlockStateModule registers state transition modules.
	RegisterBlockStateModule(modules ...kaiax.BlockStateModule)

	// FinalizeState runs post-transaction state modifications and assembles final block.
	FinalizeState(header *types.Header, stateDB *state.StateDB, txs []*types.Transaction, receipts types.Receipts) (*types.Block, error)

	// Process processes the state changes according to the Kaia rules by running
	// the transaction messages using the statedb and applying any rewards to
	// the processor (coinbase).
	Process(block *types.Block, stateDB *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, []*vm.InternalTxTrace, ProcessStats, error)
}
