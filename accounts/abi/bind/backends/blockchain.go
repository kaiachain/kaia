// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package backends

import (
	"context"
	"math/big"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/node/cn/filters"
)

// Maintain separate minimal interfaces of blockchain.BlockChain because ContractBackend are used
// in various situations. BlockChain instances are often passed down as different interfaces such as
// consensus.ChainReader, governance.blockChain, work.BlockChain.
type BlockChainForCaller interface {
	// Required by NewEVMContext
	consensus.ChainReader
}

// Maintain separate minimal interfaces of blockchain.TxPool because ContractBackend are used
// in various situations. TxPool instances are often passed down as work.TxPool.
type TxPoolForCaller interface {
	// Below is a subset of work.TxPool
	GetPendingNonce(addr common.Address) uint64
	AddLocal(tx *types.Transaction) error
	GasPrice() *big.Int
}

// BlockchainContractBackend implements bind.Contract* and bind.DeployBackend, based on
// a user-supplied blockchain.BlockChain instance.
// Its intended purpose is reading system contracts during block processing.
//
// Note that SimulatedBackend creates a new temporary BlockChain for testing,
// whereas BlockchainContractBackend uses an existing BlockChain with existing database.
type BlockchainContractBackend struct {
	bc              BlockChainForCaller
	bcCommonBackend *BlockchainContractCommonBackend
}

// This nil assignment ensures at compile time that BlockchainContractBackend implements bind.Contract* and bind.DeployBackend.
var (
	_ bind.ContractCaller     = (*BlockchainContractBackend)(nil)
	_ bind.ContractTransactor = (*BlockchainContractBackend)(nil)
	_ bind.ContractFilterer   = (*BlockchainContractBackend)(nil)
	_ bind.DeployBackend      = (*BlockchainContractBackend)(nil)
	_ bind.ContractBackend    = (*BlockchainContractBackend)(nil)
)

// `txPool` is required for bind.ContractTransactor methods and `events` is required for bind.ContractFilterer methods.
// If `tp=nil`, bind.ContractTransactor methods could return errors.
// If `es=nil`, bind.ContractFilterer methods could return errors.
func NewBlockchainContractBackend(bc BlockChainForCaller, tp TxPoolForCaller, es *filters.EventSystem) *BlockchainContractBackend {
	bcCommonBackend := newBlockchainContractCommonBackend(bc, tp, es)
	return &BlockchainContractBackend{bc, bcCommonBackend}
}

func (b *BlockchainContractBackend) getBlockAndState(num *big.Int) (*types.Block, *state.StateDB, error) {
	var block *types.Block
	if num == nil {
		block = b.bc.CurrentBlock()
	} else {
		header := b.bc.GetHeaderByNumber(num.Uint64())
		if header == nil {
			return nil, nil, errBlockDoesNotExist
		}
		block = b.bc.GetBlock(header.Hash(), header.Number.Uint64())
	}
	if block == nil {
		return nil, nil, errBlockDoesNotExist
	}

	state, err := b.bc.StateAt(block.Root())
	return block, state, err
}

// bind.ContractCaller defined methods

func (b *BlockchainContractBackend) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	_, state, err := b.getBlockAndState(blockNumber)
	if err != nil {
		return nil, err
	}
	return b.bcCommonBackend.CodeAt(ctx, state, account, blockNumber)
}

// Executes a read-only function call with respect to the specified block's state, or latest state if not specified.
//
// Returns call result in []byte.
// Returns error when:
// - cannot find the corresponding block or stateDB
// - VM revert error
// - VM other errors (e.g. NotProgramAccount, OutOfGas)
// - Error outside VM
func (b *BlockchainContractBackend) CallContract(ctx context.Context, call kaia.CallMsg, blockNumber *big.Int) ([]byte, error) {
	block, state, err := b.getBlockAndState(blockNumber)
	if err != nil {
		return nil, err
	}
	return b.bcCommonBackend.CallContract(ctx, state, block, call, blockNumber)
}

// bind.ContractTransactor defined methods

func (b *BlockchainContractBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return b.bcCommonBackend.PendingCodeAt(ctx, account)
}

func (b *BlockchainContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return b.bcCommonBackend.PendingNonceAt(ctx, account)
}

func (b *BlockchainContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return b.bcCommonBackend.SuggestGasPrice(ctx)
}

func (b *BlockchainContractBackend) EstimateGas(ctx context.Context, call kaia.CallMsg) (uint64, error) {
	return b.bcCommonBackend.EstimateGas(ctx, call)
}

func (b *BlockchainContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return b.bcCommonBackend.SendTransaction(ctx, tx)
}

func (b *BlockchainContractBackend) ChainID(ctx context.Context) (*big.Int, error) {
	return b.bcCommonBackend.ChainID(ctx)
}

// bind.ContractFilterer defined methods

func (b *BlockchainContractBackend) FilterLogs(ctx context.Context, query kaia.FilterQuery) ([]types.Log, error) {
	return b.bcCommonBackend.FilterLogs(ctx, query)
}

func (b *BlockchainContractBackend) SubscribeFilterLogs(ctx context.Context, query kaia.FilterQuery, ch chan<- types.Log) (kaia.Subscription, error) {
	return b.bcCommonBackend.SubscribeFilterLogs(ctx, query, ch)
}

// bind.DeployBackend defined methods

func (b *BlockchainContractBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return b.bcCommonBackend.TransactionReceipt(ctx, txHash)
}

// sc.Backend requires BalanceAt and CurrentBlockNumber

func (b *BlockchainContractBackend) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	_, state, err := b.getBlockAndState(blockNumber)
	if err != nil {
		return nil, err
	}
	return b.bcCommonBackend.BalanceAt(ctx, state, account, blockNumber)
}

func (b *BlockchainContractBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	return b.bcCommonBackend.CurrentBlockNumber(ctx)
}
