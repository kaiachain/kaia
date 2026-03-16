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
	"github.com/kaiachain/kaia/kaiax/valset"
)

type sortableValidator struct {
	addr common.Address
	*valset.ValidatorState
}

// getEpochTransition returns new validators after applying epoch transition
func (v *ValsetModule) getEpochTransition(
	num uint64,
	validators valset.NodeStateMap,
	idleTimeout time.Duration,
	maxValidatorCount int,
) valset.NodeStateMap {
	devLog := func(vals valset.NodeStateMap) {
		for addr, val := range vals {
			logger.Info("TODO-Permissionless: Remove this log", "num", num, "addr", addr.String(), "state", val.State.String(), "stakingamount", val.StakingAmount, "idleTimeout", val.IdleTimeout.String(), "pausedtimeout", val.PausedTimeout.String())
		}
	}
	if !isVrankEpoch(num) {
		// return early if the `num` is not vrank epoch number
		devLog(validators)
		return validators.Copy()
	}
	var (
		now                  = time.Now()
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
			val.IdleTimeout = now.Add(idleTimeout)
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
					val.IdleTimeout = now.Add(idleTimeout)
				}
			} else {
				val.State = valset.CandInactive // T2
			}
		case valset.ValReady, valset.ValActive, valset.ValPaused:
			if val.StakingAmount >= minStake {
				activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val}) // T3a
			} else {
				val.State = valset.ValInactive // T3b
				val.IdleTimeout = now.Add(idleTimeout)
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
		if idx < maxValidatorCount {
			if potentialActiveVal.State != valset.ValPaused {
				potentialActiveVal.State = valset.ValActive
			}
		} else {
			potentialActiveVal.State = valset.ValInactive
			potentialActiveVal.IdleTimeout = now.Add(idleTimeout)
		}
	}
	return newValidators
}

// getTimeoutTransition returns new validators after applying timeout transition.
// For nodes newly entering ValReady/ValInactive or ValPaused (timeout not yet set),
// it sets the timeout. If from state was already ValReady/ValInactive/ValPaused,
// the existing timeout is preserved.
func (v *ValsetModule) getTimeoutTransition(validators valset.NodeStateMap, pauseTimeout, idleTimeout time.Duration) valset.NodeStateMap {
	newValidators := validators.Copy()
	now := time.Now()
	for _, val := range newValidators {
		switch val.State {
		case valset.ValReady, valset.ValInactive:
			if val.IdleTimeout.IsZero() {
				val.IdleTimeout = now.Add(idleTimeout)
			}
			if now.After(val.IdleTimeout) {
				val.State = valset.CandInactive
				val.IdleTimeout = time.Time{}
			}
		case valset.ValPaused:
			if val.PausedTimeout.IsZero() {
				val.PausedTimeout = now.Add(pauseTimeout)
			}
			if now.After(val.PausedTimeout) {
				val.State = valset.ValInactive
				val.PausedTimeout = time.Time{}
				val.IdleTimeout = now.Add(idleTimeout)
			}
		default:
			val.IdleTimeout = time.Time{}
			val.PausedTimeout = time.Time{}
		}
	}
	return newValidators
}

// getPermlessCouncil returns a new ValidatorList containing only council members
// (validators with ValActive, ValReady, or ValPaused state).
func (vs *ValidatorList) getPermlessCouncil() *ValidatorList {
	if vs == nil {
		logger.Error("ValidatorList is nil")
		return newValidatorList(valset.NodeStateMap{})
	}
	councilVals := make(valset.NodeStateMap)
	for addr, val := range vs.permlessVals {
		switch val.State {
		case valset.ValActive, valset.ValReady, valset.ValPaused:
			councilVals[addr] = val
		}
	}
	return newValidatorList(councilVals.Copy())
}

// getViolationTransition transitions ValActive validators to ValExiting when they violate rules:
// rule1: staking amount dropped below MinimumStake, rule2: vrank violation (TODO).
func (v *ValsetModule) getViolationTransition(num uint64, validators valset.NodeStateMap) valset.NodeStateMap {
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
