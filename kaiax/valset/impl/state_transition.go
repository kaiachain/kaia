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

package impl

import (
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

const (
	VRankEpoch           = 10
	ValPausedTimeout     = time.Hour * 8
	ActiveValidatorCount = 50
)

type ValidatorList struct {
	permlessMu   sync.RWMutex
	permlessVals valset.ValidatorStateMap
}

func newValidatorList(validatorChartMap valset.ValidatorStateMap) *ValidatorList {
	return &ValidatorList{permlessVals: validatorChartMap}
}

func convertToChartMap(nodeAddrs []common.Address) valset.ValidatorStateMap {
	validators := make(valset.ValidatorStateMap)
	for _, addr := range nodeAddrs {
		validators[addr] = &valset.ValidatorState{
			// assign `ValActive` state in the permissioned operation
			State: valset.ValActive,
		}
	}
	return validators
}

func (vs *ValidatorList) EqualState(other valset.ValidatorStateMap) bool {
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	return vs.permlessVals.EqualState(other)
}

func (vs *ValidatorList) String() string {
	if vs == nil {
		return ""
	}
	return vs.permlessVals.String()
}

func (vs *ValidatorList) Marshal() ([]byte, error) {
	if vs == nil {
		return nil, errors.New("permeless valset empty")
	}
	return json.Marshal(vs.permlessVals)
}

func (vs *ValidatorList) Copy() valset.CommonAddressSet {
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	return &ValidatorList{permlessVals: vs.permlessVals.Copy()}
}

// Council returns council where includes `ValActive`, `ValReady`, and `ValPaused`
func (vs *ValidatorList) Council() []common.Address {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return []common.Address{}
	}
	// return all state of validtaors execept for candidate states
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	var ret []common.Address
	for addr, val := range vs.permlessVals {
		switch val.State {
		case valset.ValPaused, valset.ValReady, valset.ValActive:
			ret = append(ret, addr)
		}
	}
	return ret
}

func (vs *ValidatorList) Len() int {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return 0
	}
	// return the length of all state of validtaors
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	return len(vs.permlessVals)
}

func (vs *ValidatorList) Contains(targetAddr common.Address) bool {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return false
	}
	// return all state of validtaors
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	_, exists := vs.permlessVals[targetAddr]
	return exists
}

func (vs *ValidatorList) Add(addr common.Address) {
	// Permissionless: NO-OP
}

func (vs *ValidatorList) Remove(targetAddr common.Address) bool {
	// Permissionless: NO-OP
	return false
}

func (vs *ValidatorList) Subtract(other *valset.AddressSet) *valset.AddressSet {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return valset.NewAddressSet([]common.Address{})
	}
	// do not read lock because of the manipluation on copied data
	copied := vs.Copy().(*ValidatorList).permlessVals
	for _, addr := range other.Council() {
		delete(copied, addr)
	}
	result := valset.NewAddressSet(nil)
	for addr := range copied {
		result.Add(addr)
	}
	return result
}

// GetDemoted returns all state of validators where the state is not `ValActive`
func (vs *ValidatorList) GetDemoted(
	_ valset.CommonAddressSet,
	_ map[common.Address]float64,
	_ uint64,
) *valset.AddressSet {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return valset.NewAddressSet([]common.Address{})
	}
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()

	demoted := valset.NewAddressSet(nil)
	for addr, val := range vs.permlessVals {
		if val.State != valset.ValActive {
			demoted.Add(addr)
		}
	}
	return demoted
}

func (v *ValsetModule) writeValidators(num uint64, validators valset.ValidatorStateMap) {
	if validators != nil {
		if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
			writeCouncilPermissionless(v.ChainKv, num, newValidatorList(validators))
		} else {
			writeCouncilPermissioned(v.ChainKv, num, newValidatorList(validators))
		}
		insertValidatorStateChangeBlockNum(v.ChainKv, num)
		v.validatorStateChangeBlockNumsCache = nil
	}
}

// TODO-Permissionless: Replace with KIP-227 implementation
func (v *ValsetModule) isPassVrankTest() bool {
	return true
}

func isVrankEpoch(num uint64) bool {
	return num%VRankEpoch == 0
}

func (v *ValsetModule) initialPromoteLegacyValidators(header *types.Header, vmenv *vm.EVM, state *state.StateDB) error {
	var (
		config    = v.Chain.Config()
		parentNum = new(big.Int).Sub(header.Number, common.Big1)
		nextNum   = new(big.Int).Add(header.Number, common.Big1)
	)
	backend := backends.NewStateBlockchainContractBackend(v.Chain, nil, nil, state)
	council, err := v.getCouncilPermissioned(header.Number.Uint64())
	if err != nil {
		logger.Error("Failed to get council", "number", header.Number.Uint64(), "err", err.Error())
		return err
	}
	validatorStateAddr, err := readValidatorStateAddr(backend, parentNum)
	if err != nil {
		// TODO-Permissionless: Change the log level to Debug
		logger.Error("Failed to fetch ValidatorState contract adress", "number", header.Number.Uint64(), "err", err.Error())
		return err
	}
	valStateMap := convertToChartMap(council.Council())
	v.writeValidators(nextNum.Uint64(), valStateMap)
	msg, from, err := prepareValidatorWrite(backend, config, state, nextNum, validatorStateAddr, valStateMap)
	if err == nil {
		blockchain.WriteValidators(msg, from, header, vmenv, state, config.Rules(header.Number))
	}
	return nil
}

// TODO-Permissionless: implement contract and align with it
func (v *ValsetModule) isInitializedABv2(validators valset.ValidatorStateMap) bool {
	return len(validators) > 0
}

func (v *ValsetModule) ProcessTransition(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	var (
		config    = v.Chain.Config()
		parentNum = new(big.Int).Sub(header.Number, common.Big1)
	)

	// [Promoting Legacy Validator #1]
	// initialize permissioned validators at `permless HF-1`
	// all validators are registered to `ValActive` state
	if config.IsPermissionlessForBlockParent(header.Number) {
		if err := v.initialPromoteLegacyValidators(header, vmenv, state); err != nil {
			logger.Error("Failed to promote legacy validators", "number", header.Number, "err", err.Error())
		}
	}

	// TODO-Permissionless: Store canddiate scores once ABV2 and VRank interface is finalized
	if config.IsPermissionlessForkEnabled(header.Number) {
		// 0. self-state transition(user tx) might have been executed at header.Number - 1
		// 1. read all validators from contrcat on every block
		backend := backends.NewStateBlockchainContractBackend(v.Chain, nil, nil, state)
		prevCouncil, err := v.getCouncilPermissionless(parentNum.Uint64())
		if err != nil {
			logger.Error("Failed to get council", "number", parentNum, "err", err.Error())
			return err
		}

		validatorStateAddr, err := readValidatorStateAddr(backend, parentNum)
		if err != nil {
			// TODO-Permissionless: Change the log level to Debug
			logger.Error("Failed to fetch ValidatorState contract adress", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}
		si, err := v.StakingModule.GetStakingInfoFromState(header.Number.Uint64(), state)
		if err != nil {
			return nil
		}
		validators, err := system.ReadGetAllValidators(backend, validatorStateAddr, si, parentNum)
		if err != nil {
			logger.Error("Failed to fetch all validators' state", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// [Promoting Legacy Validator #2]
		// if ABv2 is not initialized yet, promote legacy validators to the state of `ValActive`
		// this code is reachable when HF is passed && registry had not been set before HF
		if !v.isInitializedABv2(validators) {
			// the promote updates the state, thus the timing of promoting must be aligned with all nodes
			// which is selected as epoch number
			if isVrankEpoch(header.Number.Uint64()) {
				if err := v.initialPromoteLegacyValidators(header, vmenv, state); err != nil {
					logger.Error("Failed to promote legacy validators", "number", header.Number, "err", err.Error())
				}
			}
			return nil
		}

		// 2. check VRank violation
		newValidators, err := v.GetVrankViolationTransition(validators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process vrank violation", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// 3. timeout transition
		newValidators = v.GetTimeoutTransition(newValidators)

		// 4. epoch transition
		newValidators, err = v.GetEpochTransition(newValidators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process epoch transition", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		if !prevCouncil.EqualState(newValidators) {
			// if contract returns empty for validator query, it's a circumstance
			// where ValidatorState contract is deployed after permless HF
			newValidatorState := newValidators
			// no `Copy()` is required because no mutation on it
			if len(validators) == 0 {
				if prevCouncilVals, ok := prevCouncil.(*ValidatorList); ok {
					newValidatorState = prevCouncilVals.permlessVals
				}
			}
			// 5. write updated validators' state into checkpoint db
			v.writeValidators(header.Number.Uint64(), newValidatorState)
			// 6. write updated validators' state into contract
			msg, from, err := prepareValidatorWrite(backend, config, state, header.Number, validatorStateAddr, newValidatorState)
			if err == nil {
				blockchain.WriteValidators(msg, from, header, vmenv, state, config.Rules(header.Number))
			}
		}
	}
	return nil
}

func readValidatorStateAddr(backend *backends.StateBlockchainContractBackend, num *big.Int) (common.Address, error) {
	validatorStateAddr, err := system.ReadValidatorStateAddr(backend, num)
	if err != nil {
		return common.Address{}, err
	}
	return validatorStateAddr, nil
}

func prepareValidatorWrite(
	backend *backends.StateBlockchainContractBackend,
	config *params.ChainConfig,
	statedb *state.StateDB,
	num *big.Int,
	validatorStateAddr common.Address,
	validators valset.ValidatorStateMap,
) (*types.Transaction, common.Address, error) {
	from, msg, err := system.EncodeWriteValidators(
		backend,
		config.Rules(num),
		validatorStateAddr,
		validators,
	)
	if err != nil {
		logger.Error("Failed to encode WriteValidators", "number", num.Uint64(), "err", err.Error(), "validators", validators.String())
		return nil, common.Address{}, err
	}
	return msg, from, err
}
