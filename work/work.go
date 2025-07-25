// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from miner/miner.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package work

import (
	"fmt"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/snapshot"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/builder"
)

var logger = log.NewModuleLogger(log.Work)

// TxPool is an interface of blockchain.TxPool used by ProtocolManager and Backend.
//
//go:generate mockgen -destination=./mocks/txpool_mock.go -package=mocks github.com/kaiachain/kaia/work TxPool
type TxPool interface {
	// HandleTxMsg should add the given transactions to the pool.
	HandleTxMsg(types.Transactions)

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending() (map[common.Address]types.Transactions, error)

	CachedPendingTxsByCount(count int) types.Transactions

	// SubscribeNewTxsEvent should return an event subscription of
	// NewTxsEvent and send events to the given channel.
	SubscribeNewTxsEvent(chan<- blockchain.NewTxsEvent) event.Subscription

	GetPendingNonce(addr common.Address) uint64
	AddLocal(tx *types.Transaction) error
	GasPrice() *big.Int
	SetGasPrice(price *big.Int)
	Stop()
	Get(hash common.Hash) *types.Transaction
	Stats() (int, int)
	Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions)
	StartSpamThrottler(conf *blockchain.ThrottlerConfig) error
	StopSpamThrottler()

	kaiax.TxPoolModuleHost
}

// Backend wraps all methods required for mining.
type Backend interface {
	AccountManager() accounts.AccountManager
	BlockChain() BlockChain
	TxPool() TxPool
	ChainDB() database.DBManager
	ReBroadcastTxs(transactions types.Transactions)
}

// Miner creates blocks and searches for proof-of-work values.
type Miner struct {
	mux     *event.TypeMux
	worker  *worker
	mining  int32
	backend Backend
	engine  consensus.Engine

	canStart    int32 // can start indicates whether we can start the mining operation
	shouldStart int32 // should start indicates whether we should start after sync
}

func New(backend Backend, config *params.ChainConfig, mux *event.TypeMux, engine consensus.Engine, nodetype common.ConnType, nodeAddr common.Address, TxResendUseLegacy bool, govModule gov.GovModule) *Miner {
	miner := &Miner{
		backend:  backend,
		mux:      mux,
		engine:   engine,
		worker:   newWorker(config, engine, nodeAddr, backend, mux, nodetype, TxResendUseLegacy, govModule),
		canStart: 1,
	}
	// TODO-Kaia drop or missing tx
	miner.Register(NewCpuAgent(backend.BlockChain(), engine, nodetype))
	go miner.update()

	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (self *Miner) update() {
	events := self.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
out:
	for ev := range events.Chan() {
		switch ev.Data.(type) {
		case downloader.StartEvent:
			atomic.StoreInt32(&self.canStart, 0)
			if self.Mining() {
				self.Stop()
				atomic.StoreInt32(&self.shouldStart, 1)
				logger.Info("Mining aborted due to sync")
			}
		case downloader.DoneEvent, downloader.FailedEvent:
			shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

			atomic.StoreInt32(&self.canStart, 1)
			atomic.StoreInt32(&self.shouldStart, 0)
			if shouldStart {
				self.Start()
			}
			// unsubscribe. we're only interested in this event once
			events.Unsubscribe()
			// stop immediately and ignore all further pending events
			break out
		}
	}
}

func (self *Miner) Start() {
	atomic.StoreInt32(&self.shouldStart, 1)

	if atomic.LoadInt32(&self.canStart) == 0 {
		logger.Info("Network syncing, will start work afterwards")
		return
	}
	atomic.StoreInt32(&self.mining, 1)

	if self.worker.nodetype == common.CONSENSUSNODE {
		logger.Info("Starting mining operation")
	}
	self.worker.start()
	self.worker.commitNewWork()
}

func (self *Miner) Stop() {
	self.worker.stop()
	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.shouldStart, 0)
}

func (self *Miner) Register(agent Agent) {
	if self.Mining() {
		agent.Start()
	}
	self.worker.register(agent)
}

func (self *Miner) Unregister(agent Agent) {
	self.worker.unregister(agent)
}

func (self *Miner) Mining() bool {
	return atomic.LoadInt32(&self.mining) > 0
}

func (self *Miner) HashRate() (tot int64) {
	if pow, ok := self.engine.(consensus.PoW); ok {
		tot += int64(pow.Hashrate())
	}
	// do we care this might race? is it worth we're rewriting some
	// aspects of the worker/locking up agents so we can get an accurate
	// hashrate?
	for agent := range self.worker.agents {
		if _, ok := agent.(*CpuAgent); !ok {
			tot += agent.GetHashRate()
		}
	}
	return
}

func (self *Miner) SetExtra(extra []byte) error {
	// istanbul BFT
	maximumExtraDataSize := params.GetMaximumExtraDataSize()
	if uint64(len(extra)) > maximumExtraDataSize {
		return fmt.Errorf("Extra exceeds max length. %d > %v", len(extra), maximumExtraDataSize)
	}
	self.worker.setExtra(extra)
	return nil
}

// Pending returns the currently pending block, corresponding receipts and associated state.
func (self *Miner) Pending() (*types.Block, types.Receipts, *state.StateDB) {
	return self.worker.pending()
}

// PendingBlock returns the currently pending block.
//
// Note, to access both the pending block and the pending state
// simultaneously, please use Pending(), as the pending state can
// change between multiple method calls
func (self *Miner) PendingBlock() *types.Block {
	return self.worker.pendingBlock()
}

// RegisterExecutionModule registers kaiax.ExecutionModule to underlying worker.
func (self *Miner) RegisterExecutionModule(modules ...kaiax.ExecutionModule) {
	self.worker.RegisterExecutionModule(modules...)
}

// RegisterTxBundlingModule registers kaiax.TxBundlingModule to underlying worker.
func (self *Miner) RegisterTxBundlingModule(txBundlingModules ...kaiax.TxBundlingModule) {
	modules := make([]builder.TxBundlingModule, len(txBundlingModules)) // TODO-Kaia: Remove this cast.
	for i, module := range txBundlingModules {
		modules[i] = module.(builder.TxBundlingModule)
	}
	self.worker.RegisterTxBundlingModule(modules...)
}

// BlockChain is an interface of blockchain.BlockChain used by ProtocolManager.
//
//go:generate mockgen -destination=./mocks/blockchain_mock.go -package=mocks github.com/kaiachain/kaia/work BlockChain
type BlockChain interface {
	Genesis() *types.Block

	CurrentBlock() *types.Block
	CurrentFastBlock() *types.Block
	HasBlock(hash common.Hash, number uint64) bool
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockByHash(hash common.Hash) *types.Block
	GetBlockByNumber(number uint64) *types.Block
	GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash

	CurrentHeader() *types.Header
	HasHeader(hash common.Hash, number uint64) bool
	GetHeader(hash common.Hash, number uint64) *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	GetHeaderByNumber(number uint64) *types.Header

	GetTd(hash common.Hash, number uint64) *big.Int
	GetTdByHash(hash common.Hash) *big.Int

	GetBodyRLP(hash common.Hash) rlp.RawValue

	GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts

	InsertChain(chain types.Blocks) (int, error)
	TrieNode(hash common.Hash) ([]byte, error)
	ContractCode(hash common.Hash) ([]byte, error)
	ContractCodeWithPrefix(hash common.Hash) ([]byte, error)
	Config() *params.ChainConfig
	State() (*state.StateDB, error)
	Rollback(chain []common.Hash)
	InsertReceiptChain(blockChain types.Blocks, receiptChain []types.Receipts) (int, error)
	InsertHeaderChain(chain []*types.Header, checkFreq int) (int, error)
	FastSyncCommitHead(hash common.Hash) error
	StateCache() state.Database

	SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription
	SetHead(head uint64) error
	Stop()

	SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription
	SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription
	IsParallelDBWrite() bool
	IsSenderTxHashIndexingEnabled() bool

	Processor() blockchain.Processor
	BadBlocks() ([]blockchain.BadBlockArgs, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	PrunableStateAt(root common.Hash, num uint64) (*state.StateDB, error)
	StateAtWithPersistent(root common.Hash) (*state.StateDB, error)
	StateAtWithGCLock(root common.Hash) (*state.StateDB, error)
	Export(w io.Writer) error
	ExportN(w io.Writer, first, last uint64) error
	Engine() consensus.Engine
	GetTxLookupInfoAndReceipt(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
	GetTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts
	GetTxLookupInfoAndReceiptInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
	GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetLogsByHash(hash common.Hash) [][]*types.Log
	ResetWithGenesisBlock(gb *types.Block) error
	Validator() blockchain.Validator
	HasBadBlock(hash common.Hash) bool
	WriteBlockWithState(block *types.Block, receipts []*types.Receipt, stateDB *state.StateDB) (blockchain.WriteResult, error)
	PostChainEvents(events []interface{}, logs []*types.Log)
	ApplyTransaction(config *params.ChainConfig, author *common.Address, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg *vm.Config) (*types.Receipt, *vm.InternalTxTrace, error)

	// State Migration
	PrepareStateMigration() error
	StartStateMigration(uint64, common.Hash) error
	StopStateMigration() error
	StateMigrationStatus() (bool, uint64, int, int, int, float64, error)

	// Warm up
	StartWarmUp(minLoad uint) error
	StartContractWarmUp(contractAddr common.Address, minLoad uint) error
	StopWarmUp() error

	// Collect state/storage trie statistics
	StartCollectingTrieStats(contractAddr common.Address) error
	GetContractStorageRoot(block *types.Block, db state.Database, contractAddr common.Address) (common.ExtHash, error)

	// Save trie node cache to this
	SaveTrieNodeCacheToDisk() error

	// KES
	BlockSubscriptionLoop(pool *blockchain.TxPool)
	CloseBlockSubscriptionLoop()

	// read-only mode
	CurrentBlockUpdateLoop(pool *blockchain.TxPool)

	// Snapshot
	Snapshots() *snapshot.Tree

	// kaiax module host
	kaiax.ExecutionModuleHost
	kaiax.RewindableModuleHost
}
