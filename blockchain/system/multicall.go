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

package system

import (
	"context"
	"math/big"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/contracts/contracts/system_contracts/multicall"
)

// ContractCallerForMultiCall is an implementation of ContractCaller only for MultiCall contract.
// The caller interacts with a multicall contract on a read only basis.
type ContractCallerForMultiCall struct {
	state  *state.StateDB               // the state that is under process
	chain  backends.BlockChainForCaller // chain containing the blockchain information
	header *types.Header                // the header of a new block that is under process
}

func (caller *ContractCallerForMultiCall) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return MultiCallCode, nil
}

// CallContract injects a multicall contract code into the state and executes the call.
func (caller *ContractCallerForMultiCall) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	gasPrice := big.NewInt(0) // execute call regardless of the balance of the sender
	gasLimit := uint64(1e8)   // enough gas limit to execute multicall contract functions
	intrinsicGas := uint64(0) // read operation doesn't require intrinsicGas

	// call.From: zero address will be assigned if nothing is specified
	// call.To: the target contract address will be assigned by `BoundContract`
	// call.Value: nil value is acceptable for `types.NewMessage`
	// call.Data: a proper value will be assigned by `BoundContract`
	// No need to handle access list here

	// Set the code of the multicall contract
	err := caller.state.SetCode(MultiCallAddr, MultiCallCode)
	if err != nil {
		return nil, err
	}

	msg := types.NewMessage(call.From, call.To, caller.state.GetNonce(call.From),
		call.Value, gasLimit, gasPrice, call.Data, false, intrinsicGas, nil)

	blockContext := blockchain.NewEVMBlockContext(caller.header, caller.chain, nil)
	txContext := blockchain.NewEVMTxContext(msg, caller.header, caller.chain.Config())
	txContext.GasPrice = gasPrice                                                                // set gasPrice again if baseFee is assigned
	evm := vm.NewEVM(blockContext, txContext, caller.state, caller.chain.Config(), &vm.Config{}) // no additional vm config required

	result, err := blockchain.ApplyMessage(evm, msg)
	return result.Return(), err
}

// NewMultiCallContractCaller creates a new instance of ContractCaller for MultiCall contract.
func NewMultiCallContractCaller(chain backends.BlockChainForCaller, header *types.Header) (*multicall.MultiCallContractCaller, error) {
	state, err := chain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	c := &ContractCallerForMultiCall{state, chain, header}
	return multicall.NewMultiCallContractCaller(MultiCallAddr, c)
}
