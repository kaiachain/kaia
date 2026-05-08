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

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// TransitionContext carries every input the pure transitions need, grouped by
// source of truth via embedded sub-structs.
// All fields are read-only during transition, except for EpochTransitionContext.
// EpochTransitionContext is overwritten after epoch transition.
//   - BlockContext           — header-derived
//   - GovContext             — governance parameters
//   - ABv2Context            — AddressBookV2 contract reads
//   - EpochTransitionContext — slot limits derived from epochVACount; recomputed
//     mid-pipeline at epoch boundaries by the orchestrator.
//   - VRankContext           — vrank scores pre-fetched at orchestrator entry
//
// Construct via NewTransitionContext + Set...Ctx setters.
type TransitionContext struct {
	BlockContext
	GovContext
	ABv2Context
	EpochTransitionContext
	VRankContext
}

// BlockContext is everything derived from the block header.
type BlockContext struct {
	BlockTime time.Time // header.Time. Used by: epoch, violation, timeout
	IsEpoch   bool      // Execute epoch transition if true. Used by: epoch (gating)
}

// GovContext holds required governance parameters.
type GovContext struct {
	MinStake uint64 // Used by: epoch, violation
}

// ABv2Context is a subset of ABv2StateSnapshot which holds values read from
// system.ReadNodeStates against the AddressBookV2 contract.
type ABv2Context struct {
	Nodes        valset.NodeMap // pipeline input. Used by: epoch, violation, timeout
	EpochVACount uint64         // pipeline initial value. Used by: epoch (carried through on non-epoch blocks)

	PfsThreshold            uint64        // Used by: violation
	CfsThreshold            uint64        // Used by: epoch
	IdleTimeout             time.Duration // Used by: epoch, violation, timeout
	PauseTimeout            time.Duration // Used by: timeout
	MaxValActivePausedCount uint64        // Used by: epoch
}

// EpochTransitionContext holds slot-limit values derived from the current
// epochVACount via slotLimitsFor (Go port of SlotMath.sol).
// During epoch transition, these values are *overwritten*.
type EpochTransitionContext struct {
	MaxSlotAvailable uint64 // Used by: violation
	MinActiveCount   uint64 // Used by: violation
}

// VRankContext holds vrank scores.
type VRankContext struct {
	CFS      map[common.Address]uint64 // Used by: epoch
	PFS      map[common.Address]uint64 // Used by: violation
	PfReport []common.Address          // Used by: violation
}

// TransitionResult is the output of the orchestrator's transition pipeline.
// epochVACountForWrite is set only for epoch blocks, where it is written back
// to AddressBookV2 via processSystemTransition.
type TransitionResult struct {
	Nodes                valset.NodeMap
	epochVACountForWrite uint64
}

func NewTransitionContext() *TransitionContext {
	return &TransitionContext{}
}

func (ctx *TransitionContext) SetBlockCtx(header *types.Header, isEpoch bool) {
	ctx.BlockContext = BlockContext{
		BlockTime: time.Unix(header.Time.Int64(), 0),
		IsEpoch:   isEpoch,
	}
}

func (ctx *TransitionContext) SetGovCtx(pset gov.ParamSet) {
	ctx.GovContext = GovContext{
		MinStake: pset.MinimumStake.Uint64(),
	}
}

func (ctx *TransitionContext) SetABv2Ctx(r *system.ABv2StateSnapshot) {
	ctx.ABv2Context = ABv2Context{
		Nodes:                   r.Nodes,
		EpochVACount:            r.EpochVACount,
		PfsThreshold:            r.PfsThreshold,
		CfsThreshold:            r.CfsThreshold,
		IdleTimeout:             r.IdleTimeout,
		PauseTimeout:            r.PauseTimeout,
		MaxValActivePausedCount: r.MaxValActivePausedCount,
	}
}

func (ctx *TransitionContext) SetSlotsCtx(maxSlot, minActive uint64) {
	ctx.EpochTransitionContext = EpochTransitionContext{
		MaxSlotAvailable: maxSlot,
		MinActiveCount:   minActive,
	}
}

func (ctx *TransitionContext) SetVRankCtx(cfs, pfs map[common.Address]uint64, pfReport []common.Address) {
	ctx.VRankContext = VRankContext{
		CFS:      cfs,
		PFS:      pfs,
		PfReport: pfReport,
	}
}

// ApplyAllTransitions runs the epoch → violation → timeout pipeline against
// the inputs already populated on ctx. Pure with respect to ctx's data
// fields (each ApplyXTransition copies before mutating); ctx itself gets its
// EpochTransitionContext refreshed mid-pipeline at epoch boundaries via
// slotLimitsFor(epochVACount).
//
// All ABv2 / gov / vrank / chain reads are expected to have happened before
// this call — the orchestrator pre-populates ctx and slot math is computed in
// Go (see slot_math.go), so this method stays free of I/O.
//
// Returns the final NodeMap and, on epoch blocks, the epoch VA count to write.
func (ctx *TransitionContext) ApplyAllTransitions() *TransitionResult {
	nodes := ctx.Nodes
	epochVACountForWrite := uint64(0)

	if ctx.IsEpoch {
		nodes = ctx.applyEpochTransition(nodes)
		// After the epoch transition the active-validator count changes, so
		// slot limits must be recomputed before violation runs.
		epochVACountForWrite = nodes.CountByState(valset.ValActive)
		ctx.SetSlotsCtx(slotLimitsFor(epochVACountForWrite))
	}
	nodes = ctx.applyViolationTransition(nodes)
	nodes = ctx.applyTimeoutTransition(nodes)

	return &TransitionResult{
		Nodes:                nodes,
		epochVACountForWrite: epochVACountForWrite,
	}
}

// applyEpochTransition returns a new NodeMap with the KIP-286 epoch transition
// applied to m. Pure: m is not mutated.
//
// Reads from ABv2Context:  IdleTimeout, MaxValActivePausedCount, CfsThreshold.
// Reads from GovContext:   MinStake.
// Reads from BlockContext: BlockTime.
// Reads from VRankContext: CFS (via isPassVrankTest).
//
// Transitions encoded (T-labels mirror KIP-286):
//
//	T1: ValExiting → ValInactive (drain)
//	T2: CandTesting → Registered (failed vrank test)
//	T3a: VR/VA/VP/CT → ValActive (won top-N stake competition)
//	T3b: VR/VA/VP/CT → ValInactive (lost top-N or below MinStake)
//	T4a: CandReady → CandTesting (above MinStake)
//	T4b: CandReady → Registered  (below MinStake)
//
// Internal use only — call via the orchestrator (applyAllTransitions). Calling
// applyEpochTransition out of order with applyViolationTransition/applyTimeoutTransition produces inconsistent state.
func (ctx *TransitionContext) applyEpochTransition(m valset.NodeMap) valset.NodeMap {
	// sortableValidator pairs an address with its mutable state so the epoch
	// transition can sort competitors by stake while still mutating their fields.
	type sortableValidator struct {
		addr common.Address
		*valset.Node
	}

	var (
		newValidators        = m.Copy()
		activeValCompetitors []sortableValidator
	)

	// T3a/T3b: compete for top-N or demote to ValInactive.
	// Shared by EpochCompetitors = {VA, VR, VP, CT}.
	competeOrDemote := func(addr common.Address, val *valset.Node) {
		if val.StakingAmount >= ctx.MinStake {
			activeValCompetitors = append(activeValCompetitors, sortableValidator{addr, val}) // T3a
		} else {
			val.State = valset.ValInactive // T3b
			val.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
		}
	}
	for addr, val := range newValidators {
		switch val.State {
		case valset.ValExiting:
			val.State = valset.ValInactive // T1
			val.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
		case valset.CandReady:
			if val.StakingAmount >= ctx.MinStake {
				val.State = valset.CandTesting // T4a
			} else {
				val.State = valset.Registered // T4b
			}
		case valset.CandTesting:
			if ctx.isPassVrankTest(addr) {
				competeOrDemote(addr, val)
			} else {
				val.State = valset.Registered // T2
			}
		case valset.ValReady, valset.ValActive, valset.ValPaused:
			competeOrDemote(addr, val)
		}
	}
	slices.SortFunc(activeValCompetitors, func(a, b sortableValidator) int {
		return cmp.Or(
			cmp.Compare(b.StakingAmount, a.StakingAmount),
			bytes.Compare(a.addr[:], b.addr[:]), // tie-breaking: address order
		)
	})
	for idx, potentialActiveVal := range activeValCompetitors {
		if uint64(idx) < ctx.MaxValActivePausedCount {
			if potentialActiveVal.State != valset.ValPaused {
				potentialActiveVal.State = valset.ValActive
				potentialActiveVal.IdleTimeout = time.Time{}
				potentialActiveVal.PausedTimeout = time.Time{}
			}
		} else {
			potentialActiveVal.State = valset.ValInactive // T3b
			potentialActiveVal.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
		}
	}
	return newValidators
}

func (ctx *TransitionContext) isPassVrankTest(addr common.Address) bool {
	cfs, ok := ctx.CFS[addr]
	if !ok { // no score recorded → pass
		return true
	}
	return cfs < ctx.CfsThreshold
}

// applyViolationTransition returns a new NodeMap with min-stake and PFS violation
// transitions applied to m. Pure: m is not mutated. Runs every block.
//
// Reads from ABv2Context:            PfsThreshold, IdleTimeout.
// Reads from EpochTransitionContext: MaxSlotAvailable, MinActiveCount.
// Reads from GovContext:             MinStake.
// Reads from BlockContext:           BlockTime.
// Reads from VRankContext:           PFS, PfReport.
//
// Rule 1 — staking < MinimumStake (only fires on non-epoch blocks since
// applyEpochTransition already demotes below-min validators to VI via T3b):
//   - ValActive → ValExiting (slot + minActiveCount floor)
//   - ValPaused → ValExiting (slot only)
//   - ValReady  → ValInactive (unconditional)
//
// Rule 2 — PFS violation (only when a proposal failure occurred *this block*):
//   - ValActive → ValExiting (PFS >= pfsThreshold, severe violation. slot + minActiveCount floor)
//   - ValActive → ValPaused (0 < PFS < pfsThreshold, minor violation. slot + minActiveCount floor)
//
// Internal use only — call via the orchestrator after applyEpochTransition on epoch blocks.
func (ctx *TransitionContext) applyViolationTransition(m valset.NodeMap) valset.NodeMap {
	var (
		newValidators = m.Copy()
		// Slot/count helpers check the in-progress state of newValidators.
		// Counts change as validators transition within the loop, so these cannot be replaced with a contract call.
		hasSlot = func(state valset.NodeState) bool {
			return newValidators.CountByState(state) < ctx.MaxSlotAvailable
		}
		// canDemoteActive additionally ensures enough ValActive remain for consensus.
		// Used only when transitioning FROM ValActive (reducing active count).
		canDemoteActive = func(targetState valset.NodeState) bool {
			return hasSlot(targetState) && newValidators.CountByState(valset.ValActive) > ctx.MinActiveCount
		}
	)

	// Iterate in deterministic address order. Slot-limited transitions depend on
	// which validator is processed first, so random map iteration would be nondeterministic.
	sortedAddrs := newValidators.Addresses()

	// rule1: staking amount dropped below MinimumStake
	for _, addr := range sortedAddrs {
		val := newValidators[addr]
		if val.StakingAmount >= ctx.MinStake {
			continue
		}
		switch val.State {
		case valset.ValActive:
			// ValActive → ValExiting (slot + minActiveCount: removing an active validator reduces consensus participants)
			if canDemoteActive(valset.ValExiting) {
				logger.Trace("MinStake violation: ValActive → ValExiting", "addr", addr, "staking", val.StakingAmount, "minStake", ctx.MinStake)
				val.State = valset.ValExiting
			} else {
				logger.Trace("MinStake violation: slot full, skipping ValActive transition", "addr", addr, "staking", val.StakingAmount)
			}
		case valset.ValPaused:
			// ValPaused → ValExiting (slot only: ValPaused is already not in active set)
			if hasSlot(valset.ValExiting) {
				logger.Trace("MinStake violation: ValPaused → ValExiting", "addr", addr, "staking", val.StakingAmount, "minStake", ctx.MinStake)
				val.State = valset.ValExiting
			} else {
				logger.Trace("MinStake violation: slot full, skipping ValPaused transition", "addr", addr, "staking", val.StakingAmount)
			}
		case valset.ValReady:
			// ValReady → ValInactive (no slot check, not in active set)
			logger.Trace("MinStake violation: ValReady → ValInactive", "addr", addr, "staking", val.StakingAmount, "minStake", ctx.MinStake)
			val.State = valset.ValInactive
			val.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
		}
	}

	// rule2: PFS violation — only triggered when a proposal failure occurred at this block.
	// - PFS >= pfsThreshold (severe): ValActive → ValExiting
	// - PFS > 0 (minor): ValActive → ValPaused
	if len(ctx.PfReport) == 0 {
		return newValidators
	}
	for _, addr := range sortedAddrs {
		val := newValidators[addr]
		if val.State != valset.ValActive {
			continue
		}
		pfs, ok := ctx.PFS[addr]
		if !ok || pfs == 0 {
			continue
		}
		if pfs >= ctx.PfsThreshold {
			if canDemoteActive(valset.ValExiting) {
				logger.Trace("PFS severe violation: transitioning to ValExiting", "addr", addr, "pfs", pfs, "pfsThreshold", ctx.PfsThreshold)
				val.State = valset.ValExiting
			} else {
				logger.Trace("PFS severe violation: slot full, skipping transition", "addr", addr, "pfs", pfs)
			}
		} else {
			if canDemoteActive(valset.ValPaused) {
				logger.Trace("PFS minor violation: transitioning to ValPaused", "addr", addr, "pfs", pfs, "pfsThreshold", ctx.PfsThreshold)
				val.State = valset.ValPaused
			} else {
				logger.Trace("PFS minor violation: slot full, skipping transition", "addr", addr, "pfs", pfs)
			}
		}
	}
	return newValidators
}

// applyTimeoutTransition returns a new NodeMap with idle/paused timeout
// transitions applied to m. Pure: m is not mutated.
//
// For nodes newly entering ValReady/ValInactive or ValPaused (timeout not yet
// set), the timeout is initialized. Existing timeouts are preserved.
//
// Reads from ABv2Context:  IdleTimeout, PauseTimeout.
// Reads from BlockContext: BlockTime.
//
// Internal use only — call via the orchestrator after applyViolationTransition.
func (ctx *TransitionContext) applyTimeoutTransition(m valset.NodeMap) valset.NodeMap {
	newValidators := m.Copy()
	for _, val := range newValidators {
		switch val.State {
		case valset.ValReady, valset.ValInactive:
			if val.IdleTimeout.IsZero() {
				val.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
			}
			if ctx.BlockTime.After(val.IdleTimeout) {
				val.State = valset.Registered
				val.IdleTimeout = time.Time{}
			}
		case valset.ValPaused:
			if val.PausedTimeout.IsZero() {
				val.PausedTimeout = ctx.BlockTime.Add(ctx.PauseTimeout)
			}
			if ctx.BlockTime.After(val.PausedTimeout) {
				val.State = valset.ValInactive
				val.PausedTimeout = time.Time{}
				val.IdleTimeout = ctx.BlockTime.Add(ctx.IdleTimeout)
			}
		default:
			val.IdleTimeout = time.Time{}
			val.PausedTimeout = time.Time{}
		}
	}
	return newValidators
}

// #region SlotMath
//
// Slot limit calculations for node state transitions. Go port of
// contracts_permissionless/contracts/libraries/SlotMath.sol — kept in lockstep
// with that contract since the orchestrator and ABv2 must agree on slot math.
//
// For n >= 4 (standard BFT with f = floor(n/3) > 0):
//
//	minActiveCount   = ceil(2n/3)             = (2n + 2) / 3
//	totalBudget      = n - minActiveCount     = floor(n/3) = n / 3
//	maxSlotAvailable = ceil(totalBudget / 2)  = (n/3 + 1) / 2
//
// For n < 4 (BFT fault tolerance f = 0; all nodes must remain active):
//
//	minActiveCount   = n
//	maxSlotAvailable = 0

// minActiveCount returns the minimum number of active nodes required for
// BFT liveness given epochVACount n.
func minActiveCount(n uint64) uint64 {
	if n < 4 {
		return n
	}
	return (2*n + 2) / 3
}

// maxSlotAvailable returns the maximum number of nodes allowed in
// ValPaused or ValExiting state each given epochVACount n.
func maxSlotAvailable(n uint64) uint64 {
	if n < 4 {
		return 0
	}
	return (n/3 + 1) / 2
}

// slotLimitsFor returns (maxSlotAvailable, minActiveCount) for epochVACount n.
func slotLimitsFor(n uint64) (uint64, uint64) {
	return maxSlotAvailable(n), minActiveCount(n)
}
