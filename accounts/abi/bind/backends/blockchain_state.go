// Modifications Copyright 2026 The Kaia Authors
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
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
)

// `StateBlockchainContractCaller` implements `ContractCaller` interface
type StateBlockchainContractCaller struct {
	bc    BlockChainForCaller
	state *state.StateDB
}

var _ bind.ContractCaller = (*StateBlockchainContractCaller)(nil)

func NewStateBlockchainContractBackend(bc BlockChainForCaller, state *state.StateDB) (*StateBlockchainContractCaller, error) {
	if state == nil {
		return nil, errors.New("statedb is nil")
	}
	return &StateBlockchainContractCaller{bc, state}, nil
}

func (b *StateBlockchainContractCaller) getBlock(num *big.Int) (*types.Block, error) {
	var block *types.Block
	if num == nil {
		block = b.bc.CurrentBlock()
	} else {
		header := b.bc.GetHeaderByNumber(num.Uint64())
		if header == nil {
			return nil, errBlockDoesNotExist
		}
		block = b.bc.GetBlock(header.Hash(), header.Number.Uint64())
	}
	if block == nil {
		return nil, errBlockDoesNotExist
	}
	return block, nil
}

// bind.ContractCaller defined methods

func (b *StateBlockchainContractCaller) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.state.GetCode(account), nil
}

func (b *StateBlockchainContractCaller) CallContract(ctx context.Context, call kaia.CallMsg, blockNumber *big.Int) ([]byte, error) {
	block, err := b.getBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	res, err := b.callContract(call, block, b.state)
	if err != nil {
		return nil, err
	}
	if len(res.Revert()) > 0 {
		return nil, blockchain.NewRevertError(res)
	}
	return res.Return(), res.Unwrap()
}

func (b *StateBlockchainContractCaller) callContract(call kaia.CallMsg, block *types.Block, state *state.StateDB) (*blockchain.ExecutionResult, error) {
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
