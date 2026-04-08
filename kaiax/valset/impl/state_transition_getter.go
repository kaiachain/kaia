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
	minStake uint64,
	validators valset.NodeStateMap,
	idleTimeout time.Duration,
	maxValidatorCount int,
	now time.Time,
	num, cfsThreshold, slotFactor uint64,
) valset.NodeStateMap {
	var (
		newValidators        = validators.Copy()
		activeValCompetitors []sortableValidator
	)
	for addr, val := range newValidators {
		switch val.State {
		case valset.ValExiting:
			val.State = valset.ValInactive // T1
			val.IdleTimeout = now.Add(idleTimeout)
		case valset.CandReady:
			if val.StakingAmount >= minStake {
				val.State = valset.CandTesting // T4a
			} else {
				val.State = valset.Registered // T4b
			}
		case valset.CandTesting:
			if v.isPassVrankTest(addr, num, cfsThreshold, slotFactor) {
				if val.StakingAmount >= minStake {
					activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val}) // T3a
				} else {
					val.State = valset.ValInactive // T3b
					val.IdleTimeout = now.Add(idleTimeout)
				}
			} else {
				val.State = valset.Registered // T2
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
				potentialActiveVal.IdleTimeout = time.Time{}
				potentialActiveVal.PausedTimeout = time.Time{}
			}
		} else {
			potentialActiveVal.State = valset.ValInactive
			potentialActiveVal.IdleTimeout = now.Add(idleTimeout)
		}
	}
	return newValidators
}

// getFallbackTransition is the last-resort committee fallback when getEpochTransition produces no ValActive.
// It promotes the top maxValidatorCount epoch-1 competition group members (sorted by stake desc, address asc)
// to ValActive regardless of stake. CandTesting validators are still subject to the vrank test.
func (v *ValsetModule) getFallbackTransition(validators valset.NodeStateMap, num, cfsThreshold, slotFactor uint64, maxValidatorCount int) valset.NodeStateMap {
	newValidators := validators.Copy()

	var activeValCompetitors []sortableValidator
	for addr, val := range newValidators {
		switch val.State {
		case valset.CandTesting:
			if v.isPassVrankTest(addr, num, cfsThreshold, slotFactor) {
				activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val})
			}
		case valset.ValReady, valset.ValActive, valset.ValPaused:
			activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val})
		}
	}
	slices.SortFunc(activeValCompetitors, func(a, b sortableValidator) int {
		return cmp.Or(
			cmp.Compare(b.StakingAmount, a.StakingAmount),
			bytes.Compare(a.addr[:], b.addr[:]), // tie-breaking: address order
		)
	})
	for idx, potentialActiveVal := range activeValCompetitors {
		if idx >= maxValidatorCount {
			break
		}
		potentialActiveVal.State = valset.ValActive
		potentialActiveVal.IdleTimeout = time.Time{}
		potentialActiveVal.PausedTimeout = time.Time{}
	}
	return newValidators
}

// getTimeoutTransition returns new validators after applying timeout transition.
// For nodes newly entering ValReady/ValInactive or ValPaused (timeout not yet set),
// it sets the timeout. If from state was already ValReady/ValInactive/ValPaused,
// the existing timeout is preserved.
func (v *ValsetModule) getTimeoutTransition(validators valset.NodeStateMap, idleTimeout, pauseTimeout time.Duration, now time.Time) valset.NodeStateMap {
	newValidators := validators.Copy()
	for _, val := range newValidators {
		switch val.State {
		case valset.ValReady, valset.ValInactive:
			if val.IdleTimeout.IsZero() {
				val.IdleTimeout = now.Add(idleTimeout)
			}
			if now.After(val.IdleTimeout) {
				val.State = valset.Registered
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

// getViolationTransition transitions ValActive validators to ValExiting when they violate rules:
// rule1: staking amount dropped below MinimumStake
// rule2: PFS >= pfsThreshold (vrank violation, anytime)
func (v *ValsetModule) getViolationTransition(minStake uint64, validators valset.NodeStateMap, num, pfsThreshold, maxSlotAvailable, minActiveCount uint64) valset.NodeStateMap {
	var (
		newValidators = validators.Copy()
		countByState  = func(state valset.State) uint64 {
			var count uint64
			for _, val := range newValidators {
				if val.State == state {
					count++
				}
			}
			return count
		}
		canTransition = func(targetState valset.State) bool {
			return countByState(targetState) < maxSlotAvailable && countByState(valset.ValActive) > minActiveCount
		}
	)

	// rule1: staking amount dropped below MinimumStake → ValExiting
	for addr, val := range newValidators {
		if val.State != valset.ValActive || val.StakingAmount >= minStake {
			continue
		}
		if canTransition(valset.ValExiting) {
			logger.Info("MinStake violation: transitioning to ValExiting", "addr", addr, "staking", val.StakingAmount, "minStake", minStake, "num", num)
			val.State = valset.ValExiting
		} else {
			logger.Warn("MinStake violation: slot full, skipping transition", "addr", addr, "staking", val.StakingAmount, "num", num)
		}
	}

	// rule2: PFS violation — only triggered when a proposal failure occurred at this block.
	// - PFS >= pfsThreshold (severe) → ValActive → ValExiting
	// - PFS > 0 (minor) → ValActive → ValPaused
	pfReport, err := v.VRankModule.GetPfReport(num)
	if err != nil {
		logger.Warn("getViolationTransition: GetPfReport failed", "num", num, "err", err)
		return newValidators
	}
	if len(pfReport) == 0 {
		return newValidators
	}
	pfsScores, err := v.VRankModule.GetPFS(num)
	if err != nil {
		logger.Warn("getViolationTransition: GetPFS failed", "num", num, "err", err)
		return newValidators
	}
	for addr, val := range newValidators {
		if val.State != valset.ValActive {
			continue
		}
		pfs, ok := pfsScores[addr]
		if !ok || pfs == 0 {
			continue
		}
		if pfs >= pfsThreshold {
			if canTransition(valset.ValExiting) {
				logger.Info("PFS severe violation: transitioning to ValExiting", "addr", addr, "pfs", pfs, "pfsThreshold", pfsThreshold, "num", num)
				val.State = valset.ValExiting
			} else {
				logger.Warn("PFS severe violation: slot full, skipping transition", "addr", addr, "pfs", pfs, "num", num)
			}
		} else {
			if canTransition(valset.ValPaused) {
				logger.Info("PFS minor violation: transitioning to ValPaused", "addr", addr, "pfs", pfs, "pfsThreshold", pfsThreshold, "num", num)
				val.State = valset.ValPaused
			} else {
				logger.Warn("PFS minor violation: slot full, skipping transition", "addr", addr, "pfs", pfs, "num", num)
			}
		}
	}
	return newValidators
}
