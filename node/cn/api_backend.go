// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/api_backend.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package cn

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/bloombits"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node/cn/gasprice"
	"github.com/kaiachain/kaia/node/cn/tracers"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work"
)

// CNAPIBackend implements api.Backend for full nodes
type CNAPIBackend struct {
	cn  *CN
	gpo *gasprice.Oracle
}

// GetTxLookupInfoAndReceipt retrieves a tx and lookup info and receipt for a given transaction hash.
func (b *CNAPIBackend) GetTxLookupInfoAndReceipt(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return b.cn.blockchain.GetTxLookupInfoAndReceipt(txHash)
}

// GetTxAndLookupInfoInCache retrieves a tx and lookup info for a given transaction hash in cache.
func (b *CNAPIBackend) GetTxAndLookupInfoInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return b.cn.blockchain.GetTxAndLookupInfoInCache(txHash)
}

// GetBlockReceiptsInCache retrieves receipts for a given block hash in cache.
func (b *CNAPIBackend) GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	return b.cn.blockchain.GetBlockReceiptsInCache(blockHash)
}

// GetTxLookupInfoAndReceiptInCache retrieves a tx and lookup info and receipt for a given transaction hash in cache.
func (b *CNAPIBackend) GetTxLookupInfoAndReceiptInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return b.cn.blockchain.GetTxLookupInfoAndReceiptInCache(txHash)
}

func (b *CNAPIBackend) ChainConfig() *params.ChainConfig {
	return b.cn.chainConfig
}

func (b *CNAPIBackend) CurrentBlock() *types.Block {
	return b.cn.blockchain.CurrentBlock()
}

func doSetHead(bc work.BlockChain, cn consensus.Engine, gpo *gasprice.Oracle, targetBlkNum uint64) error {
	if err := bc.SetHead(targetBlkNum); err != nil {
		return err
	}
	// Initialize staking info cache, and governance cache
	cn.PurgeCache()
	gpo.PurgeCache()
	return nil
}

func (b *CNAPIBackend) SetHead(number uint64) error {
	b.cn.protocolManager.Downloader().Cancel()
	b.cn.protocolManager.SetSyncStop(true)
	defer b.cn.protocolManager.SetSyncStop(false)
	return doSetHead(b.cn.blockchain, b.cn.engine, b.gpo, number)
}

func (b *CNAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.cn.miner.PendingBlock()
		if block == nil {
			return nil, fmt.Errorf("pending block is not prepared yet")
		}
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.cn.blockchain.CurrentBlock().Header(), nil
	}
	header := b.cn.blockchain.GetHeaderByNumber(uint64(blockNr))
	if header == nil {
		return nil, fmt.Errorf("the header does not exist (block number: %d)", blockNr)
	}
	return header, nil
}

func (b *CNAPIBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := b.HeaderByHash(ctx, hash)
		if err != nil {
			return nil, err
		}
		return header, nil
	}
	return nil, fmt.Errorf("invalid arguments; neither block nor hash specified")
}

func (b *CNAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	if header := b.cn.blockchain.GetHeaderByHash(hash); header != nil {
		return header, nil
	}
	return nil, fmt.Errorf("the header does not exist (hash: %s)", hash.String())
}

func (b *CNAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.cn.miner.PendingBlock()
		if block == nil {
			return nil, fmt.Errorf("pending block is not prepared yet")
		}
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.cn.blockchain.CurrentBlock(), nil
	}
	block := b.cn.blockchain.GetBlockByNumber(uint64(blockNr))
	if block == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNr)
	}
	return block, nil
}

func (b *CNAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err := b.BlockByHash(ctx, hash)
		if err != nil {
			return nil, err
		}
		return block, nil
	}
	return nil, fmt.Errorf("invalid arguments; neither block nor hash specified")
}

func (b *CNAPIBackend) Pending() (*types.Block, types.Receipts, *state.StateDB) {
	return b.cn.miner.Pending()
}

func (b *CNAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, _, state := b.cn.miner.Pending()
		if block == nil || state == nil {
			return nil, nil, fmt.Errorf("pending block is not prepared yet")
		}
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.cn.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *CNAPIBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.cn.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, nil, fmt.Errorf("header for hash not found")
		}
		stateDb, err := b.cn.BlockChain().StateAt(header.Root)
		return stateDb, header, err
	}
	return nil, nil, fmt.Errorf("invalid arguments; neither block nor hash specified")
}

func (b *CNAPIBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	block := b.cn.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", hash.String())
	}
	return block, nil
}

// GetTxAndLookupInfo retrieves a tx and lookup info for a given transaction hash.
func (b *CNAPIBackend) GetTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return b.cn.blockchain.GetTxAndLookupInfo(hash)
}

// GetBlockReceipts retrieves the receipts for all transactions with given block hash.
func (b *CNAPIBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	return b.cn.blockchain.GetReceiptsByBlockHash(hash)
}

func (b *CNAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	return b.cn.blockchain.GetLogsByHash(hash), nil
}

func (b *CNAPIBackend) GetEVM(ctx context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	vmError := func() error { return nil }

	txContext := blockchain.NewEVMTxContext(msg, header, b.ChainConfig())
	blockContext := blockchain.NewEVMBlockContext(header, b.cn.BlockChain(), nil)

	return vm.NewEVM(blockContext, txContext, state, b.cn.chainConfig, &vmCfg), vmError, nil
}

func (b *CNAPIBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return b.cn.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *CNAPIBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return b.cn.BlockChain().SubscribeChainEvent(ch)
}

func (b *CNAPIBackend) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return b.cn.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *CNAPIBackend) SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription {
	return b.cn.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *CNAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.cn.BlockChain().SubscribeLogsEvent(ch)
}

func (b *CNAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.cn.txPool.AddLocal(signedTx)
}

func (b *CNAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.cn.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *CNAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.cn.txPool.Get(hash)
}

func (b *CNAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) uint64 {
	return b.cn.txPool.GetPendingNonce(addr)
}

func (b *CNAPIBackend) Stats() (pending int, queued int) {
	return b.cn.txPool.Stats()
}

func (b *CNAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.cn.TxPool().Content()
}

func (b *CNAPIBackend) SubscribeNewTxsEvent(ch chan<- blockchain.NewTxsEvent) event.Subscription {
	return b.cn.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *CNAPIBackend) Progress() kaia.SyncProgress {
	return b.cn.Progress()
}

func (b *CNAPIBackend) ProtocolVersion() int {
	return b.cn.ProtocolVersion()
}

// SuggestPrice returns the baseFee * 2 if the current block is magma hard forked.
// If it's kaia hard forked, it returns the baseFee + SuggestTipCap.
// Other cases, it returns the unitPrice.
func (b *CNAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

// SuggestTipCap returns the 0 if the current block is magma hard forked.
// If it's kaia hard forked, it returns the SuggestTipCap based on fee history.
// Other cases, it returns the unitPrice.
func (b *CNAPIBackend) SuggestTipCap(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestTipCap(ctx)
}

func (b *CNAPIBackend) UpperBoundGasPrice(ctx context.Context) *big.Int {
	bignum := b.CurrentBlock().Number()
	pset := b.cn.govModule.GetParamSet(bignum.Uint64() + 1)
	if b.cn.chainConfig.IsMagmaForkEnabled(bignum) {
		return new(big.Int).SetUint64(pset.UpperBoundBaseFee)
	} else {
		return new(big.Int).SetUint64(pset.UnitPrice)
	}
}

func (b *CNAPIBackend) LowerBoundGasPrice(ctx context.Context) *big.Int {
	bignum := b.CurrentBlock().Number()
	pset := b.cn.govModule.GetParamSet(bignum.Uint64() + 1)
	if b.cn.chainConfig.IsMagmaForkEnabled(bignum) {
		return new(big.Int).SetUint64(pset.LowerBoundBaseFee)
	} else {
		return new(big.Int).SetUint64(pset.UnitPrice)
	}
}

func (b *CNAPIBackend) ChainDB() database.DBManager {
	return b.cn.ChainDB()
}

func (b *CNAPIBackend) EventMux() *event.TypeMux {
	return b.cn.EventMux()
}

func (b *CNAPIBackend) AccountManager() accounts.AccountManager {
	return b.cn.AccountManager()
}

func (b *CNAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.cn.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *CNAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.cn.bloomRequests)
	}
}

func (b *CNAPIBackend) IsParallelDBWrite() bool {
	return b.cn.BlockChain().IsParallelDBWrite()
}

func (b *CNAPIBackend) IsSenderTxHashIndexingEnabled() bool {
	return b.cn.BlockChain().IsSenderTxHashIndexingEnabled()
}

func (b *CNAPIBackend) IsConsoleLogEnabled() bool {
	return b.cn.config.UseConsoleLog
}

func (b *CNAPIBackend) RPCGasCap() *big.Int {
	return b.cn.config.RPCGasCap
}

func (b *CNAPIBackend) RPCEVMTimeout() time.Duration {
	return b.cn.config.RPCEVMTimeout
}

func (b *CNAPIBackend) RPCTxFeeCap() float64 {
	return b.cn.config.RPCTxFeeCap
}

func (b *CNAPIBackend) Engine() consensus.Engine {
	return b.cn.engine
}

func (b *CNAPIBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (*state.StateDB, tracers.StateReleaseFunc, error) {
	return b.cn.stateAtBlock(block, reexec, base, readOnly, preferDisk)
}

func (b *CNAPIBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (blockchain.Message, vm.BlockContext, vm.TxContext, *state.StateDB, tracers.StateReleaseFunc, error) {
	return b.cn.stateAtTransaction(block, txIndex, reexec, base, readOnly, preferDisk)
}

func (b *CNAPIBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return b.gpo.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}
