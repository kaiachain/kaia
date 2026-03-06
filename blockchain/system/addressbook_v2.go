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
	"math/big"
	"sort"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	abv2contracts "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// deriveAddressBookV2LogicAddr derives a collision-free logic contract address
// from the AddressBookV2 bytecode hash. If the derived address is already occupied,
// it increments a nonce and retries until an empty address is found.
func deriveAddressBookV2LogicAddr(state *state.StateDB) common.Address {
	for nonce := uint64(0); ; nonce++ {
		input := make([]byte, 0, len(AddressBookV2Code)+8)
		input = append(input, AddressBookV2Code...)
		input = append(input, new(big.Int).SetUint64(nonce).Bytes()...)
		addr := common.BytesToAddress(crypto.Keccak256(input)[12:])
		if state.Empty(addr) {
			return addr
		}
	}
}

// InstallAddressBookV2 deploys the AddressBookV2 UUPS proxy at AddressBookAddr (0x400),
// replacing the existing AddressBookV1. The logic contract is deployed at a
// collision-free address derived from the bytecode hash.
func InstallAddressBookV2(state *state.StateDB) error {
	logicAddr := deriveAddressBookV2LogicAddr(state)

	// Set ERC1967 proxy code at 0x400
	if err := state.SetCode(AddressBookAddr, ERC1967ProxyCode); err != nil {
		return err
	}
	// Point proxy's implementation slot to the logic contract
	state.SetState(AddressBookAddr, common.BytesToHash(ImplementationSlot), lpad32(logicAddr))
	// Clear slot 0 (_initialized) which legacy AddressBook may have set non-zero
	state.SetState(AddressBookAddr, lpad32(0), common.Hash{})

	// Deploy logic contract
	if err := state.SetCode(logicAddr, AddressBookV2Code); err != nil {
		return err
	}
	// Prevent direct initialization of the logic contract
	state.SetState(logicAddr, lpad32(0), lpad32([]byte{0xff}))

	return nil
}

// EncodeInitializeABv2 encodes the initialize() call for AddressBookV2 proxy.
// ABv2.initialize() takes no parameters — it reads all genesis data from ABv2DataContract via Registry.
func EncodeInitializeABv2(rules params.Rules) (common.Address, *types.Transaction, error) {
	data, err := AddressBookV2ABI.Pack("initialize")
	if err != nil {
		return common.Address{}, nil, err
	}
	from := params.SystemAddress
	gasLimit := params.UpperGasLimit
	intrinsicGas, err := types.IntrinsicGas(data, nil, nil, false, rules)
	if err != nil {
		return common.Address{}, nil, err
	}
	var (
		to  = AddressBookAddr
		msg = types.NewMessage(from, &to, 0, common.Big0, gasLimit, common.Big0, nil, nil, nil, data, false, intrinsicGas, nil, nil, nil, nil, nil)
	)
	return from, msg, nil
}

func EncodeWriteNodes(
	rules params.Rules,
	validators valset.NodeStateMap,
) (common.Address, *types.Transaction, error) {
	keys := make([]common.Address, 0, len(validators))
	for addr := range validators {
		keys = append(keys, addr)
	}
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].Bytes(), keys[j].Bytes()) < 0
	})

	var (
		nodeIds    = make([]common.Address, len(keys))
		newStates  = make([]uint8, len(keys))
		timeoutAts = make([]*big.Int, len(keys))
	)
	for idx, addr := range keys {
		vs := validators[addr]
		nodeIds[idx] = addr
		newStates[idx] = vs.State.ToUint8()
		var timeout time.Time
		switch vs.State {
		case valset.ValPaused:
			timeout = vs.PausedTimeout
		case valset.ValReady, valset.ValInactive:
			timeout = vs.IdleTimeout
		}
		if timeout.IsZero() {
			timeoutAts[idx] = big.NewInt(0)
		} else {
			timeoutAts[idx] = big.NewInt(timeout.Unix())
		}
	}

	data, err := AddressBookV2ABI.Pack("processSystemTransition", nodeIds, newStates, timeoutAts)
	if err != nil {
		return common.Address{}, nil, err
	}
	var (
		from     = params.SystemAddress
		gasLimit = params.UpperGasLimit
		to       = AddressBookAddr
	)
	intrinsicGas, err := types.IntrinsicGas(data, nil, nil, false, rules)
	if err != nil {
		return common.Address{}, nil, err
	}
	msg := types.NewMessage(
		from,         // from common.Address,
		&to,          // to *common.Address,
		0,            // nonce uint64,
		common.Big0,  // amount *big.Int,
		gasLimit,     // gasLimit uint64,
		common.Big0,  // gasPrice *big.Int
		nil,          // gasFeeCap *big.Int
		nil,          // gasTipCap *big.Int
		nil,          // blobGasFeeCap *big.Int
		data,         // data []byte
		false,        // checkNonce bool
		intrinsicGas, // intrinsicGas uint64
		nil,          // list AccessList
		nil,          // chainId *big.Int
		nil,          // blobHashes []common.Hash
		nil,          // sidecar *BlobTxSidecar
		nil,          // auth []SetCodeAuthorization
	)
	return from, msg, nil
}

// ReadABv2Timeouts reads pauseTimeout and idleTimeout durations from the AddressBookV2 contract.
func ReadABv2Timeouts(
	backend bind.ContractCaller,
	num *big.Int,
) (pauseTimeout, idleTimeout time.Duration, err error) {
	caller, err := abv2contracts.NewAddressBookV2Caller(AddressBookAddr, backend)
	if err != nil {
		return 0, 0, err
	}
	opts := &bind.CallOpts{BlockNumber: num}
	pause, idle, err := caller.GetTimeouts(opts)
	if err != nil {
		return 0, 0, err
	}
	return time.Duration(pause.Int64()) * time.Second,
		time.Duration(idle.Int64()) * time.Second,
		nil
}

// ReadABv2MaxCounts reads maxValidatorCount and maxReadyCandidateCount from the AddressBookV2 contract.
func ReadABv2MaxCounts(
	backend bind.ContractCaller,
	num *big.Int,
) (maxValidatorCount, maxReadyCandidateCount uint64, err error) {
	caller, err := abv2contracts.NewAddressBookV2Caller(AddressBookAddr, backend)
	if err != nil {
		return 0, 0, err
	}
	opts := &bind.CallOpts{BlockNumber: num}
	valCount, candCount, err := caller.GetMaxCounts(opts)
	if err != nil {
		return 0, 0, err
	}
	return valCount.Uint64(), candCount.Uint64(), nil
}

func ReadGetAllValidators(
	statedb *state.StateDB,
	chain backends.BlockChainForCaller,
	header *types.Header,
) (valset.NodeStateMap, error) {
	caller, err := NewMultiCallContractCaller(statedb, chain, header)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{BlockNumber: header.Number}

	res, err := caller.MultiCallStakingInfoPermissionless(opts)
	if err != nil {
		return nil, err
	}

	validators := make(valset.NodeStateMap)
	for i, p := range res.Profiles {
		nodeState := valset.State(p.State)
		vs := &valset.ValidatorState{
			State:         nodeState,
			StakingAmount: new(big.Int).Div(res.StakingAmounts[i], big.NewInt(params.KAIA)).Uint64(),
		}
		switch nodeState {
		case valset.ValReady, valset.ValInactive:
			if p.TimeoutAt.Sign() > 0 {
				vs.IdleTimeout = time.Unix(p.TimeoutAt.Int64(), 0)
			}
		case valset.ValPaused:
			if p.TimeoutAt.Sign() > 0 {
				vs.PausedTimeout = time.Unix(p.TimeoutAt.Int64(), 0)
			}
		}
		validators[p.NodeId] = vs
	}
	return validators, nil
}

// ReadAddressBookV2BlsAll reads BLS public key info for all nodes from AddressBookV2.
// Used for permissionless blocks where ABv2 is the BLS source of truth.
func ReadAddressBookV2BlsAll(backend bind.ContractCaller, num *big.Int) (BlsPublicKeyInfos, error) {
	caller, err := abv2contracts.NewAddressBookV2Caller(AddressBookAddr, backend)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{BlockNumber: num}
	ret, err := caller.GetAllBlsInfo(opts)
	if err != nil {
		return nil, err
	}
	return buildBlsPublicKeyInfos(ret.NodeIdList, len(ret.PubkeyList), func(i int) ([]byte, []byte) {
		return ret.PubkeyList[i].PublicKey, ret.PubkeyList[i].Pop
	})
}