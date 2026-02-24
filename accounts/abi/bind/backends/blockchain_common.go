// Modifications Copyright 2026Kaia Authors
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
	"errors"
	"math/big"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/node/cn/filters"
)

type BlockchainContractCommonBackend struct {
	bc     BlockChainForCaller
	txPool TxPoolForCaller
	events *filters.EventSystem
}

func newBlockchainContractCommonBackend(bc BlockChainForCaller, tp TxPoolForCaller, es *filters.EventSystem) *BlockchainContractCommonBackend {
	return &BlockchainContractCommonBackend{bc, tp, es}
}

func (b *BlockchainContractCommonBackend) CodeAt(ctx context.Context, state *state.StateDB, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return state.GetCode(account), nil
}

// Executes a read-only function call with respect to the specified block's state, or latest state if not specified.
//
// Returns call result in []byte.
// Returns error when:
// - cannot find the corresponding block or stateDB
// - VM revert error
// - VM other errors (e.g. NotProgramAccount, OutOfGas)
// - Error outside VM
func (b *BlockchainContractCommonBackend) CallContract(ctx context.Context, state *state.StateDB, block *types.Block, call kaia.CallMsg, blockNumber *big.Int) ([]byte, error) {
	res, err := b.callContract(call, block, state)
	if err != nil {
		return nil, err
	}
	if len(res.Revert()) > 0 {
		return nil, blockchain.NewRevertError(res)
	}
	return res.Return(), res.Unwrap()
}

func (b *BlockchainContractCommonBackend) callContract(call kaia.CallMsg, block *types.Block, state *state.StateDB) (*blockchain.ExecutionResult, error) {
	gasPrice := common.Big0
	if call.GasPrice != nil && (call.GasFeeCap != nil || call.GasTipCap != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	if !b.bc.Config().IsKaiaForkEnabled(block.Number()) { // before KIP-162
		if call.GasPrice != nil {
			gasPrice = call.GasPrice
		}
	} else { // after KIP-162
		if call.GasPrice != nil {
			gasPrice = call.GasPrice
		} else {
			if call.GasFeeCap == nil {
				call.GasFeeCap = big.NewInt(0)
			}
			if call.GasTipCap == nil {
				call.GasTipCap = big.NewInt(0)
			}
			gasPrice = math.BigMin(new(big.Int).Add(call.GasTipCap, block.Header().BaseFee), call.GasFeeCap)
		}
	}

	if call.Gas == 0 {
		call.Gas = uint64(3e8) // enough gas for ordinary contract calls
	}

	var accessList types.AccessList
	if call.AccessList != nil {
		accessList = *call.AccessList
	}
	intrinsicGas, err := types.IntrinsicGas(call.Data, accessList, nil, call.To == nil, b.bc.Config().Rules(block.Number()))
	if err != nil {
		return nil, err
	}

	msg := types.NewMessage(call.From, call.To, 0, call.Value, call.Gas, gasPrice, nil, nil, nil, call.Data,
		false, intrinsicGas, accessList, nil, nil, nil, nil)

	txContext := blockchain.NewEVMTxContext(msg, block.Header(), b.bc.Config())
	blockContext := blockchain.NewEVMBlockContext(block.Header(), b.bc, nil)

	// EVM demands the sender to have enough KAIA balance (gasPrice * gasLimit) in buyGas()
	// After KIP-71, gasPrice is nonzero baseFee, regardless of the msg.gasPrice (usually 0)
	// But our sender (usually 0x0) won't have enough balance. Instead we override gasPrice = 0 here
	txContext.GasPrice = big.NewInt(0)
	evm := vm.NewEVM(blockContext, txContext, state, b.bc.Config(), &vm.Config{})

	return blockchain.ApplyMessage(evm, msg)
}

// bind.ContractTransactor defined methods

func (b *BlockchainContractCommonBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	// TODO-Kaia this is not pending code but latest code
	state, err := b.bc.State()
	if err != nil {
		return nil, err
	}
	return state.GetCode(account), nil
}

func (b *BlockchainContractCommonBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if b.txPool != nil {
		return b.txPool.GetPendingNonce(account), nil
	}
	// TODO-Kaia this is not pending nonce but latest nonce
	state, err := b.bc.State()
	if err != nil {
		return 0, err
	}
	return state.GetNonce(account), nil
}

func (b *BlockchainContractCommonBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if b.bc.Config().IsMagmaForkEnabled(b.bc.CurrentBlock().Number()) {
		if b.txPool != nil {
			return new(big.Int).Mul(b.txPool.GasPrice(), big.NewInt(2)), nil
		} else {
			return new(big.Int).Mul(b.bc.CurrentBlock().Header().BaseFee, big.NewInt(2)), nil
		}
	} else {
		// This is used for sending txs, so it's ok to use the genesis value instead of the latest value from governance.
		return new(big.Int).SetUint64(b.bc.Config().UnitPrice), nil
	}
}

func (b *BlockchainContractCommonBackend) EstimateGas(ctx context.Context, call kaia.CallMsg) (uint64, error) {
	state, err := b.bc.State()
	if err != nil {
		return 0, err
	}
	balance := state.GetBalance(call.From) // from can't be nil

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *blockchain.ExecutionResult, error) {
		call.Gas = gas

		currentState, err := b.bc.State()
		if err != nil {
			return true, nil, nil
		}
		res, err := b.callContract(call, b.bc.CurrentBlock(), currentState)
		if err != nil {
			if errors.Is(err, blockchain.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return res.Failed(), res, nil
	}

	gasPrice := common.Big0
	if call.GasPrice != nil && (call.GasFeeCap != nil || call.GasTipCap != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if call.GasPrice != nil {
		gasPrice = call.GasPrice
	} else if call.GasFeeCap != nil {
		gasPrice = call.GasFeeCap
	}

	estimated, err := blockchain.DoEstimateGas(ctx, call.Gas, 0, call.Value, gasPrice, balance, executable)
	return uint64(estimated), err
}

func (b *BlockchainContractCommonBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if b.txPool == nil {
		return errors.New("tx pool not configured")
	}
	return b.txPool.AddLocal(tx)
}

func (b *BlockchainContractCommonBackend) ChainID(ctx context.Context) (*big.Int, error) {
	return b.bc.Config().ChainID, nil
}

// bind.ContractFilterer defined methods

func (b *BlockchainContractCommonBackend) FilterLogs(ctx context.Context, query kaia.FilterQuery) ([]types.Log, error) {
	// Convert the current block numbers into internal representations
	if query.FromBlock == nil {
		query.FromBlock = big.NewInt(b.bc.CurrentBlock().Number().Int64())
	}
	if query.ToBlock == nil {
		query.ToBlock = big.NewInt(b.bc.CurrentBlock().Number().Int64())
	}
	from := query.FromBlock.Int64()
	to := query.ToBlock.Int64()

	state, err := b.bc.State()
	if err != nil {
		return nil, err
	}
	bc, ok := b.bc.(*blockchain.BlockChain)
	if !ok {
		return nil, errors.New("BlockChainForCaller is not blockchain.BlockChain")
	}
	filter := filters.NewRangeFilter(&filterBackend{state.Database().TrieDB().DiskDB(), bc, nil}, from, to, query.Addresses, query.Topics)

	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]types.Log, len(logs))
	for i, log := range logs {
		res[i] = *log
	}
	return res, nil
}

func (b *BlockchainContractCommonBackend) SubscribeFilterLogs(ctx context.Context, query kaia.FilterQuery, ch chan<- types.Log) (kaia.Subscription, error) {
	// Subscribe to contract events
	sink := make(chan []*types.Log)

	if b.events == nil {
		return nil, errors.New("events system not configured")
	}
	sub, err := b.events.SubscribeLogs(query, sink)
	if err != nil {
		return nil, err
	}
	// Since we're getting logs in batches, we need to flatten them into a plain stream
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case logs := <-sink:
				for _, log := range logs {
					select {
					case ch <- *log:
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// bind.DeployBackend defined methods

func (b *BlockchainContractCommonBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	bc, ok := b.bc.(*blockchain.BlockChain)
	if !ok {
		return nil, errors.New("BlockChainForCaller is not blockchain.BlockChain")
	}
	receipt := bc.GetReceiptByTxHash(txHash)
	if receipt != nil {
		return receipt, nil
	}
	return nil, errors.New("receipt does not exist")
}

// sc.Backend requires BalanceAt and CurrentBlockNumber

func (b *BlockchainContractCommonBackend) BalanceAt(ctx context.Context, state *state.StateDB, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return state.GetBalance(account), nil
}

func (b *BlockchainContractCommonBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	return b.bc.CurrentBlock().NumberU64(), nil
}
