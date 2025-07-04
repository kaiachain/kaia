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
	"math/big"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/randao"
	"github.com/kaiachain/kaia/kaiax/staking"
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

	// Engine retrieves the header chain's consensus engine.
	Engine() Engine

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

// Engine is an algorithm agnostic consensus engine.
//
//go:generate mockgen -destination=./mocks/engine_mock.go -package=mocks github.com/kaiachain/kaia/consensus Engine
type Engine interface {
	// Author retrieves the Kaia address of the account that minted the given
	// block.
	Author(header *types.Header) (common.Address, error)

	// CanVerifyHeadersConcurrently returns true if concurrent header verification possible, otherwise returns false.
	CanVerifyHeadersConcurrently() bool

	// PreprocessHeaderVerification prepares header verification for heavy computation before synchronous header verification such as ecrecover.
	PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, header *types.Header) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, header *types.Header) error

	// Initialize runs any pre-transaction state modifications (e.g., EIP-2539)
	Initialize(chain ChainReader, header *types.Header, state *state.StateDB)

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new block for the given input block with the local miner's
	// seal place on top.
	Seal(chain ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error)

	// CalcBlockScore is the blockscore adjustment algorithm. It returns the blockscore
	// that a new block should have.
	CalcBlockScore(chain ChainReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainReader) []rpc.API

	// Protocol returns the protocol for this consensus
	Protocol() Protocol

	// GetConsensusInfo returns consensus information regarding the given block number.
	GetConsensusInfo(block *types.Block) (ConsensusInfo, error)

	PurgeCache()
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
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

// Istanbul is a consensus engine to avoid byzantine failure
type Istanbul interface {
	Engine

	// Start starts the engine
	Start(chain ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error

	// Stop stops the engine
	Stop() error

	// SetChain sets chain of the Istanbul backend
	SetChain(chain ChainReader)

	RegisterKaiaxModules(mGov gov.GovModule, mStaking staking.StakingModule, mValset valset.ValsetModule, mRandao randao.RandaoModule)

	kaiax.ConsensusModuleHost
	staking.StakingModuleHost
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
