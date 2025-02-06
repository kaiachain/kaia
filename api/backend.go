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
// This file is derived from internal/ethapi/backend.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"context"
	"math/big"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
//
//go:generate mockgen -destination=./mocks/backend_mock.go -package=mock_api github.com/kaiachain/kaia/api Backend
type Backend interface {
	// General Kaia API
	Progress() kaia.SyncProgress
	ProtocolVersion() int
	SuggestPrice(ctx context.Context) (*big.Int, error)
	SuggestTipCap(ctx context.Context) (*big.Int, error)
	UpperBoundGasPrice(ctx context.Context) *big.Int
	LowerBoundGasPrice(ctx context.Context) *big.Int
	ChainDB() database.DBManager
	EventMux() *event.TypeMux
	AccountManager() accounts.AccountManager
	RPCEVMTimeout() time.Duration // global timeout for eth/kaia_call/estimateGas/estimateComputationCost
	RPCGasCap() *big.Int          // global gas cap for eth/kaia_call/estimateGas/estimateComputationCost
	RPCTxFeeCap() float64         // global tx fee cap in eth_signTransaction
	Engine() consensus.Engine
	FeeHistory(ctx context.Context, blockCount int, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error)

	// BlockChain API
	SetHead(number uint64) error
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error)
	BlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error)
	BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error)
	StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error)
	StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error)
	GetBlockReceipts(ctx context.Context, blockHash common.Hash) types.Receipts
	GetTxLookupInfoAndReceipt(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
	GetTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetEVM(ctx context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error)
	SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription
	IsParallelDBWrite() bool

	IsSenderTxHashIndexingEnabled() bool
	IsConsoleLogEnabled() bool

	// TxPool API
	SendTx(ctx context.Context, signedTx *types.Transaction) error
	GetPoolTransactions() (types.Transactions, error)
	GetPoolTransaction(txHash common.Hash) *types.Transaction
	GetPoolNonce(ctx context.Context, addr common.Address) uint64
	Stats() (pending int, queued int)
	TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions)
	SubscribeNewTxsEvent(chan<- blockchain.NewTxsEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block

	GetTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts
	GetTxLookupInfoAndReceiptInCache(Hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
}
