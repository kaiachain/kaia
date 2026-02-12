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
	"sort"
	"sync"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

const (
	VRankEpoch           = 50
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

	permlessVals := make(valset.ValidatorChartMap, len(vs.permlessVals))
	for addr, chart := range vs.permlessVals {
		chartCopy := &valset.ValidatorChart{
			State: chart.State,
			// PausedDuration: chart.PausedDuration,
			StakingAmount: chart.StakingAmount,
		}
		permlessVals[addr] = chartCopy
	}
	var permissionedVals *valset.AddressSet
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
	if vs.isLegacy {
		return vs.permissionedVals.Len()
	}
	// return the length of all state of validtaors
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	return len(vs.permlessVals)
}

func (vs *ValidatorList) Contains(targetAddr common.Address) bool {
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
	if vs.isLegacy {
		vs.permissionedVals.Add(addr)
		return
	}
	// In permissionless operation, `governance.vote('addvalidator', 0x...)` is not the expected usage. Thus, do nothing
}

func (vs *ValidatorList) Remove(targetAddr common.Address) bool {
	if vs.isLegacy {
		return vs.permissionedVals.Remove(targetAddr)
	}
	vs.permlessMu.Lock()
	defer vs.permlessMu.Unlock()

	_, exist := vs.permlessVals[targetAddr]
	delete(vs.permlessVals, targetAddr)
	return exist
}

func (vs *ValidatorList) Subtract(other *valset.AddressSet) *valset.AddressSet {
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

func updateStakingAmount(si *staking.StakingInfo, validators valset.ValidatorChartMap) {
	if si == nil {
		return
	}
	stakingAmountM := make(map[common.Address]float64)
	for i, nodeId := range si.NodeIds {
		stakingAmountM[nodeId] = float64(si.StakingAmounts[i])
	}
	for addr, val := range validators {
		val.StakingAmount = stakingAmountM[addr]
	}
}

// at vrank epoch:
//   - [O] ValExiting -> ValInactive
//   - [O] CandReady -> CandTesting
//   - [O] CandTesting -> CandInactive

// - [O] CandTesting -> CandInactive
// - [O] ValReady -> ValInactive
// - [O] ValActive -> ValInactive
// - [O] ValPaused -> ValInactive
// - [O] CandTesting -> ValActive
// - [O] ValReady -> ValActive
// - [O] ValActive -> ValActive
func (v *ValsetModule) epochTransition(si *staking.StakingInfo, num uint64, validators valset.ValidatorChartMap) valset.ValidatorChartMap {
	defer func() {
		for addr, val := range validators {
			logger.Info("TODO-Permissionless: Remove this log", "num", num, "addr", addr.String(), "state", val.State.String(), "stakingamount", val.StakingAmount, "idleTimeout", val.IdleTimeout.String(), "pausedtimeout", val.PausedTimeout.String())
		}
	}()
	if !isVrankEpoch(num) {
		// return early if the `num` is not vrank epoch number
		return validators
	}
	if len(si.NodeIds) == 0 && len(si.StakingContracts) == 0 && len(si.RewardAddrs) == 0 {
		// Not ABook activated yet
		return nil
	}

	newValidators := validators.Copy()
	// update staking amount of council
	updateStakingAmount(si, newValidators)

	var (
		potentialActiveValSet []*valset.ValidatorChart
		pset                  = v.GovModule.GetParamSet(num - 1)    // read gov param from parent number
		minStake              = float64(pset.MinimumStake.Uint64()) // in KAIA
	)
	for _, val := range newValidators {
		switch val.State {
		case valset.ValExiting:
			val.State = valset.ValInactive
		case valset.CandReady:
			if val.StakingAmount >= minStake {
				val.State = valset.CandTesting
			} else {
				val.State = valset.CandInactive
			}
		case valset.CandTesting:
			if v.isPassVrankTest() {
				if val.StakingAmount >= minStake {
					potentialActiveValSet = append(potentialActiveValSet, val)
				} else {
					val.State = valset.ValInactive
				}
			} else {
				val.State = valset.CandInactive
			}
		case valset.ValReady, valset.ValActive, valset.ValPaused:
			if val.StakingAmount >= minStake {
				potentialActiveValSet = append(potentialActiveValSet, val)
			} else {
				val.State = valset.ValInactive
			}
		}
	}
	sort.Slice(potentialActiveValSet, func(i, j int) bool {
		return potentialActiveValSet[i].StakingAmount > potentialActiveValSet[j].StakingAmount
	})
	for idx, potentialActiveVal := range potentialActiveValSet {
		if idx < ActiveValidatorCount {
			if potentialActiveVal.State != valset.ValPaused {
				potentialActiveVal.State = valset.ValActive
			}
		} else {
			potentialActiveVal.State = valset.ValInactive
		}
	}
	return newValidators
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

func (v *ValsetModule) timeoutTransition(validators valset.ValidatorChartMap) valset.ValidatorChartMap {
	newValidators := validators.Copy()
	for _, val := range newValidators {
		switch val.State {
		case valset.ValReady, valset.ValInactive:
			if time.Now().After(val.IdleTimeout) {
				val.State = valset.CandInactive
			}
		case valset.ValPaused:
			if time.Now().After(val.PausedTimeout) {
				val.State = valset.ValInactive
			}
		}
	}
	return newValidators
}

func (v *ValsetModule) deactiveStakersLessMinStakingAmount(si *staking.StakingInfo, num uint64, validators valset.ValidatorChartMap) valset.ValidatorChartMap {
	if len(si.NodeIds) == 0 && len(si.StakingContracts) == 0 && len(si.RewardAddrs) == 0 {
		// Not ABook activated yet
		return nil
	}

	var (
		pset          = v.GovModule.GetParamSet(num - 1)    // read gov param from parent number
		minStake      = float64(pset.MinimumStake.Uint64()) // in KAIA
		newValidators = validators.Copy()
	)

	// update staking amount of council
	updateStakingAmount(si, newValidators)
	for _, val := range newValidators {
		if val.State == valset.ValActive {
			if val.StakingAmount < minStake {
				val.State = valset.ValExiting
			}
		}
	}
	return newValidators
}

func (v *ValsetModule) ProcessTransition(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	config := v.Chain.Config()
	if config.IsPermissionlessForkEnabled(header.Number) {
		// 0 self-state transition(user tx) might have been executed header.Number - 1
		// 1. TODO-Permissionless: read all validators from contrcat on every block
		var (
			parentNum = new(big.Int).Sub(header.Number, common.Big1)
			backend   = backends.NewStateBlockchainContractBackend(v.Chain, nil, nil, state)
		)

		validatorStateAddr, err := readValidatorStateAddr(backend, parentNum)
		if err != nil {
			// TODO-Permissionless: Change the log level to Debug
			logger.Error("Failed to fetch ValidatorState contract adress", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}
		validators, err := system.ReadGetAllValidators(backend, validatorStateAddr, parentNum)
		if err != nil {
			logger.Error("Failed to fetch all validators' state", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// 2. check VRank violation
		newValidators, err := v.VrankViolationTransition(validators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process vrank violation", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		// 3. timeout transition
		newValidators = v.TimeoutTransition(newValidators)

		// 4. epoch transition
		newValidators, err = v.EpochTransition(newValidators, header.Number.Uint64(), state)
		if err != nil {
			logger.Error("Failed to process epoch transition", "number", header.Number.Uint64(), "err", err.Error())
			return err
		}

		if !validators.Equal(newValidators) {
			// 5. write updated validators' state into checkpoint db
			v.writeValidators(header.Number.Uint64(), newValidators)
			// 6. write updated validators' state into contract
			msg, from, err := prepareValidatorWrite(backend, config, state, header, parentNum, validatorStateAddr, newValidators)
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
	parentNum *big.Int,
	validatorStateAddr common.Address,
	validators valset.ValidatorChartMap,
) (*types.Transaction, common.Address, error) {
	from, msg, err := system.EncodeWriteValidators(
		backend,
		config.Rules(parentNum),
		parentNum,
		validatorStateAddr,
		validators,
	)
	if err != nil {
		logger.Error("Failed to encode WriteValidators", "number", header.Number.Uint64(), "err", err.Error(), "validators", validators.String())
		return nil, common.Address{}, err
	}
	return msg, from, err
}
