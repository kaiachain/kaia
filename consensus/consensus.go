// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from consensus/consensus.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package consensus

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// CurrentBlock revrieves the current block from the local chain.
	CurrentBlock() *types.Block

	// ValidateHeader validates a header according to chain rules.
	ValidateHeader(header *types.Header) error

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block

	// State retrieves statedb
	State() (*state.StateDB, error)

	// StateAt retrieves statedb on a particular point in time
	StateAt(root common.Hash) (*state.StateDB, error)
}

// ChainContext extends ChainReader with transaction execution capability.
// This interface is used by components that need to execute transactions.
type ChainContext interface {
	ChainReader

	// ApplyTransaction applies a transaction to the given state and returns the receipt.
	ApplyTransaction(config *params.ChainConfig, author *common.Address, statedb *state.StateDB,
		header *types.Header, tx *types.Transaction, usedGas *uint64, cfg *vm.Config) (*types.Receipt, *vm.InternalTxTrace, error)
}

type ChainReaderWithSealer interface {
	ChainReader
	Sealer() Sealer
}

// Sealer handles seal-related logic independent from consensus networking/runtime.
type Sealer interface {
	// Read information from header.Extra.
	Author(header *types.Header) (common.Address, error)
	Committers(header *types.Header) ([]common.Address, error)
	Vanity(header *types.Header) ([]byte, error)
	RawSeals(header *types.Header) ([]byte, [][]byte, error)
	Round(header *types.Header) (byte, error)
	Validators(header *types.Header) ([]common.Address, error)

	// Write information to header.Extra.
	WriteAuthorSeal(header *types.Header, seal []byte) error
	WriteCommittedSeals(header *types.Header, committedSeals [][]byte) error
	WriteRound(header *types.Header, round int64)
	WriteValidators(header *types.Header, validators []common.Address) error

	// Create seal bytes
	MakeAuthorSeal(header *types.Header) ([]byte, error)
	MakeCommittedSeal(header *types.Header) ([]byte, error)
	HeaderHash(header *types.Header) common.Hash
	SigHash(header *types.Header) common.Hash

	// Verify seals and derive consensus metadata.
	F(blockNum uint64, qualifiedlen, committeeSize int) int
	Quorum(blockNum uint64, qualifiedlen, committeeSize int) int
}

// Engine is an algorithm agnostic consensus engine.
//
//go:generate mockgen -destination=./mocks/engine_mock.go -package=mocks github.com/kaiachain/kaia/consensus Engine
type Engine interface {
	// Author retrieves the Kaia address of the account that minted the given
	// block.
	Author(header *types.Header) (common.Address, error)

	// Committers extracts committed-seal signer addresses from the given header.
	// Engines that do not use committed seals may return (nil, nil).
	Committers(header *types.Header) ([]common.Address, error)

	// RegisterKaiaxModules wires kaiax modules to the consensus engine.
	RegisterKaiaxModules(mGov gov.GovModule, mValset valset.ValsetModule)

	// Start starts any consensus-specific background lifecycle.
	// Engines without a background runtime should implement this as a no-op.
	Start(chain ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool, executor Executor) error

	// Stop stops any consensus-specific background lifecycle.
	// Engines without a background runtime should implement this as a no-op.
	Stop() error

	// PrepareExtra builds consensus-specific header.Extra content.
	PrepareExtra(header *types.Header, parent *types.Header) ([]byte, error)

	// SubmitTransactions submits transactions for execution and consensus.
	// Returns finalizeCh which receives the execution result when block is finalized (for DB write and broadcast).
	// onPrepared is called after transactions are executed and block is finalized but before consensus sealing.
	// This allows the caller to update pending block state for APIs.
	SubmitTransactions(txs *types.TransactionsByPriceAndNonce, state *state.StateDB, header *types.Header, mux *event.TypeMux, onPrepared func(*ExecutionResult)) (finalizeCh <-chan *ExecutionResult)

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainReader) []rpc.API

	// VerifySeals checks consensus-specific seals for the given header.
	VerifySeals(header *types.Header) error

	// Protocol returns the protocol for this consensus.
	Protocol() Protocol

	// GetConsensusInfo returns consensus information regarding the given block number.
	GetConsensusInfo(block *types.Block) (ConsensusInfo, error)

	PurgeCache()

	// SubscribeNewSequence subscribes to new sequence events.
	// Returns nil for non-Istanbul engines.
	SubscribeNewSequence() *event.TypeMuxSubscription
}

// Handler should be implemented is the consensus needs to handle and send peer's message
type Handler interface {
	// NewChainHead handles a new head block comes
	NewChainHead() error

	// HandleMsg handles a message from peer
	HandleMsg(address common.Address, data p2p.Msg) (bool, error)

	// SetBroadcaster sets the broadcaster to send message to peers
	SetBroadcaster(Broadcaster, common.ConnType)

	// RegisterConsensusMsgCode registers the channel of consensus msg.
	RegisterConsensusMsgCode(Peer)
}

type ConsensusInfo struct {
	// Proposer signs [sigHash] to make seal; Validators signs [block.Hash + msgCommit] to make committedSeal
	SigHash        common.Hash
	Proposer       common.Address
	OriginProposer *common.Address // the proposer of 0th round at the same block number
	Committee      []common.Address
	Committers     []common.Address
	Round          byte
}
