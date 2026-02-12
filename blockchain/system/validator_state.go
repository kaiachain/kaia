// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package system

import (
	"bytes"
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	contracts "github.com/kaiachain/kaia/contracts/contracts/system_contracts/vrank"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

var (
	ValidatorStateABI, _        = contracts.IValidatorStateMetaData.GetAbi()
	ValidatorStateNotRegistered = errors.New("ValidatorState contract not registered")
)

type ValidatorState struct {
	Addr          common.Address `abi:"addr"`
	State         uint8          `abi:"state"`
	IdleTimeout   *big.Int       `abi:"idleTimeout"`
	PausedTimeout *big.Int       `abi:"pausedTimeout"`
}

func ReadValidatorStateAddr(backend *backends.StateBlockchainContractBackend, num *big.Int) (common.Address, error) {
	validatorStateAddr, err := ReadActiveAddressFromRegistry(backend, ValidatorStateName, num)
	if err != nil {
		return common.Address{}, err
	}
	if validatorStateAddr == (common.Address{}) {
		return common.Address{}, ValidatorStateNotRegistered
	}
	return validatorStateAddr, nil
}

func EncodeWriteValidators(
	backend *backends.StateBlockchainContractBackend,
	rules params.Rules,
	num *big.Int,
	validatorStateAddr common.Address,
	validators valset.ValidatorChartMap,
) (common.Address, *types.Transaction, error) {
	// Generate order validator addresses to prevent invalid merkle roots
	// The order is important becasue of the solidity set type `EnumerableSet.AddressSet`
	keys := make([]common.Address, 0, len(validators))
	for addr := range validators {
		keys = append(keys, addr)
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].Bytes(), keys[j].Bytes()) < 0
	})

	// Prepare contract input
	valStates := make([]ValidatorState, len(validators))
	for idx, addr := range keys {
		var (
			state         = validators[addr].State.ToUint8()
			idleTimeout   = big.NewInt(validators[addr].IdleTimeout.Unix())
			pausedTimeout = big.NewInt(validators[addr].PausedTimeout.Unix())
		)
		valStates[idx] = ValidatorState{
			Addr:          addr,
			State:         state,
			IdleTimeout:   idleTimeout,
			PausedTimeout: pausedTimeout,
		}
	}
	data, err := ValidatorStateABI.Pack("setValidatorStates", valStates)
	if err != nil {
		return common.Address{}, nil, err
	}
	var (
		from     = params.SystemAddress
		gasLimit = params.UpperGasLimit
	)
	intrinsicGas, err := types.IntrinsicGas(data, nil, nil, false, rules)
	if err != nil {
		return common.Address{}, nil, err
	}
	msg := types.NewMessage(
		from,                // from common.Address,
		&validatorStateAddr, // to *common.Address,
		0,                   // nonce uint64,
		common.Big0,         // amount *big.Int,
		gasLimit,            // gasLimit uint64,
		common.Big0,         // gasPrice *big.Int
		nil,                 // gasFeeCap *big.Int
		nil,                 // gasTipCap *big.Int
		nil,                 // blobGasFeeCap *big.Int
		data,                // data []byte
		false,               // checkNonce bool
		intrinsicGas,        // intrinsicGas uint64
		nil,                 // list AccessList
		nil,                 // chainId *big.Int
		nil,                 // blobHashes []common.Hash
		nil,                 // sidecar *BlobTxSidecar
		nil,                 // auth []SetCodeAuthorization
	)
	return from, msg, nil
}

func ReadGetAllValidators(backend *backends.StateBlockchainContractBackend, validatorStateAddr common.Address, num *big.Int) (valset.ValidatorChartMap, error) {
	// Prepare caller
	caller, err := contracts.NewIValidatorStateCaller(validatorStateAddr, backend)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{BlockNumber: num}

	// Call contract
	valStates, err := caller.GetAllValidators(opts)
	if err != nil {
		return nil, err
	}

	// Parse the result
	validators := make(valset.ValidatorChartMap)
	for _, valState := range valStates {
		addr, state, pausedTimeout, idleTimeout := valState.Addr, valState.State, valState.PausedTimeout, valState.IdleTimeout
		validators[addr] = &valset.ValidatorChart{
			State:         valset.State(state),
			IdleTimeout:   time.Unix(idleTimeout.Int64(), 0),
			PausedTimeout: time.Unix(pausedTimeout.Int64(), 0),
		}
	}
	return validators, nil
}
