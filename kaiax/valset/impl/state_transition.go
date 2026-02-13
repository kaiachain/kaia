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
	VRankEpoch           = 20
	ValPausedTimeout     = time.Hour * 8
	ActiveValidatorCount = 50
)

type ValidatorList struct {
	permlessMu       sync.RWMutex
	isLegacy         bool
	permissionedVals *valset.AddressSet
	permlessVals     valset.ValidatorChartMap
}

func convertToChartMap(nodeAddrs []common.Address) valset.ValidatorChartMap {
	validators := make(valset.ValidatorChartMap)
	for _, addr := range nodeAddrs {
		validators[addr] = &valset.ValidatorChart{
			// assign `ValActive` state in the permissioned operation
			State: valset.ValActive,
		}
	}
	return validators
}

func ConvertLegacyToValidatorList(nodeAddrs []common.Address) *ValidatorList {
	if len(nodeAddrs) == 0 {
		// zero council
		return nil
	}
	return &ValidatorList{
		isLegacy:         true,
		permissionedVals: valset.NewAddressSet(nodeAddrs),
		permlessVals:     convertToChartMap(nodeAddrs),
	}
}

func ConvertLegacySetToValidatorList(set *valset.AddressSet) *ValidatorList {
	if set == nil || set.Len() == 0 {
		return nil
	}
	return &ValidatorList{
		isLegacy:         true,
		permissionedVals: set,
		permlessVals:     convertToChartMap(set.List()),
	}
}

func NewValidatorList(validatorChartMap valset.ValidatorChartMap) *ValidatorList {
	return &ValidatorList{
		isLegacy:         false,
		permissionedVals: nil,
		permlessVals:     validatorChartMap,
	}
}

func (vs *ValidatorList) SetLegacyFalse() {
	vs.isLegacy = false
}

func (vs *ValidatorList) Copy() *ValidatorList {
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()

	var (
		permissionedVals *valset.AddressSet
		permlessVals     = vs.permlessVals.Copy()
	)
	if vs.permissionedVals != nil {
		permissionedVals = vs.permissionedVals.Copy()
	}
	return &ValidatorList{
		isLegacy:         vs.isLegacy,
		permissionedVals: permissionedVals,
		permlessVals:     permlessVals,
	}
}

func (vs *ValidatorList) List() []common.Address {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return []common.Address{}
	}
	if vs.isLegacy {
		return vs.permissionedVals.List()
	}
	// return all state of validtaors execept for candidate states
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	var ret []common.Address
	for addr, val := range vs.permlessVals {
		switch val.State {
		case valset.ValInactive, valset.ValPaused, valset.ValExiting, valset.ValReady, valset.ValActive:
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
	if vs.isLegacy {
		return vs.permissionedVals.Len()
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
	if vs.isLegacy {
		return vs.permissionedVals.Contains(targetAddr)
	}
	// return all state of validtaors
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	_, exists := vs.permlessVals[targetAddr]
	return exists
}

func (vs *ValidatorList) Add(addr common.Address) {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return
	}
	if vs.isLegacy {
		vs.permissionedVals.Add(addr)
		return
	}
	// In permissionless operation, `governance.vote('addvalidator', 0x...)` is not the expected usage. Thus, do nothing
}

func (vs *ValidatorList) Remove(targetAddr common.Address) bool {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return false
	}
	if vs.isLegacy {
		return vs.permissionedVals.Remove(targetAddr)
	}
	vs.permlessMu.Lock()
	defer vs.permlessMu.Unlock()

	val, exist := vs.permlessVals[targetAddr]
	if exist {
		val.State = valset.CandInactive
	}
	return exist
}

func (vs *ValidatorList) Subtract(other *valset.AddressSet) *valset.AddressSet {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return valset.NewAddressSet([]common.Address{})
	}
	if vs.isLegacy {
		return vs.permissionedVals.Subtract(other)
	}

	// do not read lock because of the manipluation on copied data
	copied := vs.Copy().permlessVals
	for _, addr := range other.List() {
		delete(copied, addr)
	}
	result := valset.NewAddressSet(nil)
	for addr := range copied {
		result.Add(addr)
	}
	return result
}

// GetDemoted returns all state of validators where the state is not `ValActive`
func (vs *ValidatorList) GetDemoted() *valset.AddressSet {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return valset.NewAddressSet([]common.Address{})
	}
	demoted := valset.NewAddressSet(nil)
	if vs.isLegacy {
		return demoted
	}
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()

	for addr, val := range vs.permlessVals {
		if val.State != valset.ValActive {
			demoted.Add(addr)
		}
	}
	return demoted
}

func (v *ValsetModule) writeValidators(num uint64, validators valset.ValidatorChartMap) {
	if validators != nil {
		writeCouncil(v.ChainKv, num, NewValidatorList(validators))
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

func (v *ValsetModule) ProcessTransition(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	config := v.Chain.Config()
	if config.IsPermissionlessForkEnabled(header.Number) {
		// 0. self-state transition(user tx) might have been executed at header.Number - 1
		// 1. read all validators from contrcat on every block
		var (
			parentNum = new(big.Int).Sub(header.Number, common.Big1)
			backend   = backends.NewStateBlockchainContractBackend(v.Chain, nil, nil, state)
		)

		prevCouncil, err := v.getCouncil(parentNum.Uint64())
		if err != nil {
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

		// 2. vote transition
		newValidators, err := v.voteTransition(parentNum.Uint64(), validators)
		if err != nil {
			logger.Error("Failed to handle vote data", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// 3. check VRank violation
		newValidators, err = v.GetVrankViolationTransition(newValidators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process vrank violation", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// 4. timeout transition
		newValidators = v.GetTimeoutTransition(newValidators)

		// 5. epoch transition
		newValidators, err = v.GetEpochTransition(newValidators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process epoch transition", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		if !prevCouncil.permlessVals.EqualState(newValidators) {
			// 6. write updated validators' state into checkpoint db
			v.writeValidators(header.Number.Uint64(), newValidators)
			// 7. write updated validators' state into contract
			msg, from, err := prepareValidatorWrite(backend, config, state, header, validatorStateAddr, newValidators)
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
	header *types.Header,
	validatorStateAddr common.Address,
	validators valset.ValidatorChartMap,
) (*types.Transaction, common.Address, error) {
	from, msg, err := system.EncodeWriteValidators(
		backend,
		config.Rules(header.Number),
		validatorStateAddr,
		validators,
	)
	if err != nil {
		logger.Error("Failed to encode WriteValidators", "number", header.Number.Uint64(), "err", err.Error(), "validators", validators.String())
		return nil, common.Address{}, err
	}
	return msg, from, err
}
