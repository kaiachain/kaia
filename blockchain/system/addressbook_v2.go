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
	abv2contracts "github.com/kaiachain/kaia/contracts_permissionless/contracts/AddressBookV2"
	abv2data "github.com/kaiachain/kaia/contracts_permissionless/contracts/AddressBookV2/abv2data"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// ReadABv2Implementation reads the ABv2 logic contract address from ABv2DataContract.
// ABv2DataContract is resolved from Registry and stores the implementation as an immutable.
func ReadABv2Implementation(backend bind.ContractCaller, num *big.Int) (common.Address, error) {
	dataContractAddr, err := ReadActiveAddressFromRegistry(backend, "ABv2DataContract", num)
	if err != nil {
		return common.Address{}, err
	}
	caller, err := abv2data.NewABv2DataContractCaller(dataContractAddr, backend)
	if err != nil {
		return common.Address{}, err
	}
	return caller.Implementation(&bind.CallOpts{BlockNumber: num})
}

// InstallAddressBookV2 sets up the AddressBookV2 UUPS proxy at AddressBookAddr (0x400).
// logicAddr is the pre-deployed ABv2 implementation contract address
// (stored in ABv2DataContract.implementation()).
func InstallAddressBookV2(state *state.StateDB, logicAddr common.Address) error {
	// Set ERC1967 proxy code at 0x400
	if err := state.SetCode(AddressBookAddr, ERC1967ProxyV5Code); err != nil {
		return err
	}
	// Point proxy's implementation slot to the logic contract
	state.SetState(AddressBookAddr, common.BytesToHash(ImplementationSlot), lpad32(logicAddr))

	return nil
}

// encodeSystemMessage creates a system transaction message for the given target and data.
func encodeSystemMessage(rules params.Rules, to common.Address, data []byte) (common.Address, *types.Transaction, error) {
	var (
		from     = params.SystemAddress
		gasLimit = params.UpperGasLimit
	)
	intrinsicGas, err := types.IntrinsicGas(data, nil, nil, false, rules)
	if err != nil {
		return common.Address{}, nil, err
	}
	msg := types.NewMessage(from, &to, 0, common.Big0, gasLimit, common.Big0, nil, nil, nil, data, false, intrinsicGas, nil, nil, nil, nil, nil)
	return from, msg, nil
}

// EncodeInitializeABv2 encodes the initialize() call for AddressBookV2 proxy.
// ABv2.initialize() takes no parameters — it reads all genesis data from ABv2DataContract via Registry.
func EncodeInitializeABv2(rules params.Rules) (common.Address, *types.Transaction, error) {
	data, err := AddressBookV2ABI.Pack("initialize")
	if err != nil {
		return common.Address{}, nil, err
	}
	return encodeSystemMessage(rules, AddressBookAddr, data)
}

// EncodeWriteNodes encodes the processSystemTransition call with the given validator state changes.
func EncodeWriteNodes(
	rules params.Rules,
	validators valset.NodeStateMap,
) (common.Address, *types.Transaction, error) {
	nodeIds := make([]common.Address, 0, len(validators))
	for addr := range validators {
		nodeIds = append(nodeIds, addr)
	}
	// Sort for deterministic ABI encoding across all nodes
	sort.Slice(nodeIds, func(i, j int) bool {
		return bytes.Compare(nodeIds[i].Bytes(), nodeIds[j].Bytes()) < 0
	})

	var (
		newStates  = make([]uint8, len(nodeIds))
		timeoutAts = make([]*big.Int, len(nodeIds))
	)
	for idx, addr := range nodeIds {
		vs := validators[addr]
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
	return encodeSystemMessage(rules, AddressBookAddr, data)
}

// ReadNodeStates reads all validator states, timeouts, max counts, and exit threshold from ABv2 in a single MultiCall.
func ReadNodeStates(
	statedb *state.StateDB,
	chain backends.BlockChainForCaller,
	header *types.Header,
) (valset.NodeStateMap, time.Duration, time.Duration, uint64, uint64, uint64, error) {
	caller, err := NewMultiCallContractCaller(statedb, chain, header)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}
	opts := &bind.CallOpts{BlockNumber: header.Number}

	res, err := caller.MultiCallNodeStatesPermissionless(opts)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
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

	var (
		pauseTimeout      = time.Duration(res.PauseTimeout.Int64()) * time.Second
		idleTimeout       = time.Duration(res.IdleTimeout.Int64()) * time.Second
		maxValCount       = res.MaxValidatorCount.Uint64()
		maxReadyCandCount = res.MaxReadyCandidateCount.Uint64()
		exitThreshold     = res.ExitThreshold.Uint64()
	)
	return validators, pauseTimeout, idleTimeout, maxValCount, maxReadyCandCount, exitThreshold, nil
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
