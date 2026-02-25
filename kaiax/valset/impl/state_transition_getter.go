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
	"bytes"
	"cmp"
	"slices"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
)

type sortableValidator struct {
	addr common.Address
	*valset.ValidatorState
}

// getEpochTransition returns new validators after applying epoch transition
func (v *ValsetModule) getEpochTransition(si *staking.StakingInfo, num uint64, validators valset.ValidatorStateMap) valset.ValidatorStateMap {
	devLog := func(vals valset.ValidatorStateMap) {
		for addr, val := range vals {
			logger.Info("TODO-Permissionless: Remove this log", "num", num, "addr", addr.String(), "state", val.State.String(), "stakingamount", val.StakingAmount, "idleTimeout", val.IdleTimeout.String(), "pausedtimeout", val.PausedTimeout.String())
		}
	}
	if !isVrankEpoch(num) {
		// return early if the `num` is not vrank epoch number
		devLog(validators)
		return validators.Copy()
	}
	if len(si.NodeIds) == 0 && len(si.StakingContracts) == 0 && len(si.RewardAddrs) == 0 {
		devLog(validators)
		// Not ABook activated yet
		return nil
	}

	var (
		newValidators        = validators.Copy()
		activeValCompetitors []sortableValidator
		pset                 = v.GovModule.GetParamSet(num - 1) // read gov param from parent number
		minStake             = pset.MinimumStake.Uint64()       // in KAIA
	)
	defer devLog(newValidators)
	for addr, val := range newValidators {
		switch val.State {
		case valset.ValExiting:
			val.State = valset.ValInactive // T1
		case valset.CandReady:
			if val.StakingAmount >= minStake {
				val.State = valset.CandTesting // T4a
			} else {
				val.State = valset.CandInactive // T4b
			}
		case valset.CandTesting:
			if v.isPassVrankTest() {
				if val.StakingAmount >= minStake {
					activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val}) // T3a
				} else {
					val.State = valset.ValInactive // T3b
				}
			} else {
				val.State = valset.CandInactive // T2
			}
		case valset.ValReady, valset.ValActive, valset.ValPaused:
			if val.StakingAmount >= minStake {
				activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val}) // T3a
			} else {
				val.State = valset.ValInactive // T3b
			}
		}
	}
	slices.SortFunc(activeValCompetitors, func(a, b sortableValidator) int {
		return cmp.Or(
			cmp.Compare(b.StakingAmount, a.StakingAmount),
			bytes.Compare(a.addr[:], b.addr[:]), // tie-breaking: address order
		)
	})
	for idx, potentialActiveVal := range activeValCompetitors {
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

// getTimeoutTransition returns new validators after applying timeout transition
func (v *ValsetModule) getTimeoutTransition(validators valset.ValidatorStateMap) valset.ValidatorStateMap {
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

// getPermlessCouncil returns a new ValidatorList containing only council members
// (validators with ValActive, ValReady, or ValPaused state).
func (vs *ValidatorList) getPermlessCouncil() *ValidatorList {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return newValidatorList(valset.ValidatorStateMap{})
	}
	vs.permlessMu.RLock()
	defer vs.permlessMu.RUnlock()
	councilVals := make(valset.ValidatorStateMap)
	for addr, val := range vs.permlessVals {
		switch val.State {
		case valset.ValActive, valset.ValReady, valset.ValPaused:
			councilVals[addr] = val
		}
	}
	return newValidatorList(councilVals.Copy())
}

func (v *ValsetModule) deactiveStakersLessMinStakingAmount(si *staking.StakingInfo, num uint64, validators valset.ValidatorStateMap) valset.ValidatorStateMap {
	if len(si.NodeIds) == 0 && len(si.StakingContracts) == 0 && len(si.RewardAddrs) == 0 {
		// Not ABook activated yet
		return nil
	}

	var (
		pset          = v.GovModule.GetParamSet(num - 1) // read gov param from parent number
		minStake      = pset.MinimumStake.Uint64()       // in KAIA
		newValidators = validators.Copy()
	)

	// rule1: check if staking amount become less than minimum staking amount
	for _, val := range newValidators {
		if val.State == valset.ValActive {
			if val.StakingAmount < minStake {
				val.State = valset.ValExiting
			}
		}
	}
	// rule2: RC
	// TODO-Permissionless: implement me
	return newValidators
}
