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
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// #region Type aliases
type (
	Node      = valset.Node
	NodeMap   = valset.NodeMap
	NodeState = valset.NodeState
)

const (
	Registered  = valset.Registered
	CandReady   = valset.CandReady
	CandTesting = valset.CandTesting
	ValInactive = valset.ValInactive
	ValReady    = valset.ValReady
	ValActive   = valset.ValActive
	ValPaused   = valset.ValPaused
	ValExiting  = valset.ValExiting
)

// Test-wide constants and stake tiers, mirroring the permissionless impl tests.
const (
	testMinStake     = uint64(5_000_000)
	aboveMinStake    = uint64(10_000_000)
	belowMinStake    = uint64(4_999_999)
	testMaxValCount  = uint64(50)
	testCFSThreshold = uint64(10)

	highStake = uint64(30_000_000)
	midStake  = uint64(20_000_000)
	lowStake  = uint64(10_000_000)

	// Disable slot checks in tests that don't care about slot math.
	noSlotLimit = uint64(100)
	noMinActive = uint64(0)
)

var (
	testBlockTime      = time.Unix(1700000000, 0) // fixed timestamp for deterministic tests
	testIdleTimeout    = 24 * time.Hour
	testPauseTimeout   = 8 * time.Hour
	defaultIdleTimeout = 30 * 24 * time.Hour

	addr1 = common.HexToAddress("0x0001")
	addr2 = common.HexToAddress("0x0002")
	addr3 = common.HexToAddress("0x0003")
	addr4 = common.HexToAddress("0x0004")
	addr5 = common.HexToAddress("0x0005")
)

// ctxOpts collects every value the transition pipeline reads. Zero-valued
// fields are safe defaults (e.g. nil maps for vrank → fail-open).
type ctxOpts struct {
	MinStake                uint64
	PfsThreshold            uint64
	CfsThreshold            uint64
	IdleTimeout             time.Duration
	PauseTimeout            time.Duration
	MaxValActivePausedCount uint64
	MaxSlotAvailable        uint64
	MinActiveCount          uint64
	BlockTime               time.Time

	CFS      map[common.Address]uint64
	PFS      map[common.Address]uint64
	PfReport []common.Address
}

// buildCtx constructs a TransitionContext from opts, using testBlockTime if
// opts.BlockTime is zero.
func buildCtx(t *testing.T, opts ctxOpts) *TransitionContext {
	t.Helper()
	bt := opts.BlockTime
	if bt.IsZero() {
		bt = testBlockTime
	}
	header := &types.Header{
		Number: big.NewInt(999),
		Time:   big.NewInt(bt.Unix()),
	}
	ctx := NewTransitionContext()
	ctx.SetBlockCtx(header, false)
	ctx.SetGovCtx(gov.ParamSet{MinimumStake: new(big.Int).SetUint64(opts.MinStake)})
	ctx.SetABv2TransitionParam(system.ABv2TransitionParam{
		PfsThreshold:            opts.PfsThreshold,
		CfsThreshold:            opts.CfsThreshold,
		IdleTimeout:             opts.IdleTimeout,
		PauseTimeout:            opts.PauseTimeout,
		MaxValActivePausedCount: opts.MaxValActivePausedCount,
	})
	ctx.SetSlotsCtx(opts.MaxSlotAvailable, opts.MinActiveCount)
	ctx.SetVRankCtx(opts.CFS, opts.PFS, opts.PfReport)
	return ctx
}

// epochCtx returns a context with epoch-transition defaults.
func epochCtx(t *testing.T, cfs map[common.Address]uint64, maxValCount uint64) *TransitionContext {
	t.Helper()
	return buildCtx(t, ctxOpts{
		MinStake:                testMinStake,
		CfsThreshold:            testCFSThreshold,
		IdleTimeout:             testIdleTimeout,
		MaxValActivePausedCount: maxValCount,
		CFS:                     cfs,
	})
}

// timeoutCtx returns a context with timeout-transition defaults.
func timeoutCtx(t *testing.T) *TransitionContext {
	t.Helper()
	return buildCtx(t, ctxOpts{
		IdleTimeout:  defaultIdleTimeout,
		PauseTimeout: testPauseTimeout,
	})
}

// violationCtx returns a context with violation-transition defaults; opts
// fields override defaults (MinStake / IdleTimeout fall back if zero).
func violationCtx(t *testing.T, opts ctxOpts) *TransitionContext {
	t.Helper()
	if opts.MinStake == 0 {
		opts.MinStake = testMinStake
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = testIdleTimeout
	}
	return buildCtx(t, opts)
}

//#region applyEpochTransition

func TestApplyEpochTransition_StateTransitions(t *testing.T) {
	cases := []struct {
		name          string
		inputState    NodeState
		stakingAmount uint64
		expectedState NodeState
		expectTimeout bool
	}{
		{"T1: ValExiting → ValInactive", ValExiting, aboveMinStake, ValInactive, true},
		{"T4a: CandReady + stake above → CandTesting", CandReady, aboveMinStake, CandTesting, false},
		{"T4b: CandReady + stake below → Registered", CandReady, belowMinStake, Registered, false},
		{"T3a: CandTesting + stake above → ValActive", CandTesting, aboveMinStake, ValActive, false},
		{"T3b: CandTesting + stake below → ValInactive", CandTesting, belowMinStake, ValInactive, true},
		{"ValReady + stake above → ValActive", ValReady, aboveMinStake, ValActive, false},
		{"ValActive + stake above → ValActive (preserved)", ValActive, aboveMinStake, ValActive, false},
		{"ValPaused + stake above → ValPaused (preserved)", ValPaused, aboveMinStake, ValPaused, false},
		{"T3b: ValActive + stake below → ValInactive", ValActive, belowMinStake, ValInactive, true},
		{"T3b: ValReady + stake below → ValInactive", ValReady, belowMinStake, ValInactive, true},
		{"T3b: ValPaused + stake below → ValInactive", ValPaused, belowMinStake, ValInactive, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NodeMap{addr1: {State: tc.inputState, StakingAmount: tc.stakingAmount}}
			// CFS = 0 so CandTesting always passes vrank in this table; the
			// failing-vrank case lives in TestApplyEpochTransition_CandTestingCFS.
			ctx := epochCtx(t, map[common.Address]uint64{addr1: 0}, testMaxValCount)
			out := ctx.applyEpochTransition(m)
			assert.Equal(t, tc.expectedState, out[addr1].State, "not expected state")
			if tc.expectTimeout {
				assert.False(t, out[addr1].IdleTimeout.IsZero(), "IdleTimeout should be set")
				assert.Equal(t, ctx.BlockTime.Add(ctx.IdleTimeout), out[addr1].IdleTimeout, "IdleTimeout should be set")
			}
		})
	}
}

// TestApplyEpochTransition_BelowMinStakeDemoted verifies that VA/VR/VP+belowMin
// are demoted to VI at epoch (T3b), so newSF is computed only from validators
// with sufficient stake.
func TestApplyEpochTransition_BelowMinStakeDemoted(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: aboveMinStake}, // above-min: stays VA
		addr2: {State: ValActive, StakingAmount: belowMinStake}, // T3b → VI
		addr3: {State: ValReady, StakingAmount: belowMinStake},  // T3b → VI
		addr4: {State: ValPaused, StakingAmount: belowMinStake}, // T3b → VI
	}
	ctx := epochCtx(t, nil, testMaxValCount)
	out := ctx.applyEpochTransition(m)

	assert.Equal(t, ValActive, out[addr1].State)
	assert.Equal(t, ValInactive, out[addr2].State)
	assert.Equal(t, ValInactive, out[addr3].State)
	assert.Equal(t, ValInactive, out[addr4].State)
	assert.False(t, out[addr2].IdleTimeout.IsZero())
	assert.False(t, out[addr3].IdleTimeout.IsZero())
	assert.False(t, out[addr4].IdleTimeout.IsZero())
}

func TestApplyEpochTransition_MaxValActivePausedCount(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: highStake},
		addr2: {State: ValActive, StakingAmount: midStake},
		addr3: {State: ValActive, StakingAmount: lowStake},
	}
	ctx := epochCtx(t, nil, 2)
	out := ctx.applyEpochTransition(m)
	assert.Equal(t, ValActive, out[addr1].State)
	assert.Equal(t, ValActive, out[addr2].State)
	assert.Equal(t, ValInactive, out[addr3].State)
	assert.False(t, out[addr3].IdleTimeout.IsZero())
}

func TestApplyEpochTransition_TieBreakingByAddress(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: lowStake},
		addr2: {State: ValActive, StakingAmount: lowStake},
		addr3: {State: ValActive, StakingAmount: lowStake},
	}
	ctx := epochCtx(t, nil, 2)
	out := ctx.applyEpochTransition(m)

	// addr1 (0x0001) < addr2 (0x0002) < addr3 (0x0003)
	assert.Equal(t, ValActive, out[addr1].State)
	assert.Equal(t, ValActive, out[addr2].State)
	assert.Equal(t, ValInactive, out[addr3].State)
}

func TestApplyEpochTransition_DefensiveCopy(t *testing.T) {
	original := NodeMap{
		addr1: {State: ValExiting, StakingAmount: aboveMinStake},
	}
	ctx := epochCtx(t, nil, testMaxValCount)
	out := ctx.applyEpochTransition(original)

	assert.Equal(t, ValExiting, original[addr1].State, "input must not mutate")
	assert.Equal(t, ValInactive, out[addr1].State, "output reflects T1")
	assert.NotSame(t, original[addr1], out[addr1], "output points to a fresh Node")
}

func TestApplyEpochTransition_CandTestingCFS(t *testing.T) {
	cases := []struct {
		name          string
		cfs           uint64
		expectedState NodeState
	}{
		{"CFS above threshold → Registered (T2)", testCFSThreshold + 5, Registered},
		{"CFS below threshold → ValActive (T3a)", testCFSThreshold - 5, ValActive},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NodeMap{addr1: {State: CandTesting, StakingAmount: aboveMinStake}}
			ctx := epochCtx(t, map[common.Address]uint64{addr1: tc.cfs}, testMaxValCount)
			out := ctx.applyEpochTransition(m)
			assert.Equal(t, tc.expectedState, out[addr1].State)
		})
	}
}

//#endregion

//#region isPassVrankTest

func TestIsPassVrankTest(t *testing.T) {
	cases := []struct {
		name     string
		cfs      map[common.Address]uint64
		expected bool
	}{
		{"CFS below threshold → pass", map[common.Address]uint64{addr1: testCFSThreshold - 5}, true},
		{"CFS equals threshold → fail", map[common.Address]uint64{addr1: testCFSThreshold}, false},
		{"CFS above threshold → fail", map[common.Address]uint64{addr1: testCFSThreshold + 5}, false},
		{"addr not in CFS scores → pass (fail-open)", map[common.Address]uint64{addr2: testCFSThreshold + 5}, true},
		{"nil CFS map → pass (fail-open)", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := buildCtx(t, ctxOpts{CfsThreshold: testCFSThreshold, CFS: tc.cfs})
			assert.Equal(t, tc.expected, ctx.isPassVrankTest(addr1))
		})
	}
}

//#endregion

//#region applyViolationTransition

func TestApplyViolationTransition_ValActiveBelowMinStake(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: belowMinStake},
		addr2: {State: ValActive, StakingAmount: aboveMinStake},
	}
	ctx := violationCtx(t, ctxOpts{MaxSlotAvailable: noSlotLimit, MinActiveCount: noMinActive})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, ValExiting, out[addr1].State)
	assert.Equal(t, ValActive, out[addr2].State)
}

// TestApplyViolationTransition_SuspendedNotExempt verifies that suspended
// ValActive validators are subject to the same violation rules as non-suspended
// ones — Suspended is informational, not a violation exemption.
func TestApplyViolationTransition_SuspendedNotExempt(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: belowMinStake, Suspended: true},  // minStake violation
		addr2: {State: ValActive, StakingAmount: aboveMinStake, Suspended: true},  // PFS severe
		addr3: {State: ValActive, StakingAmount: aboveMinStake, Suspended: false}, // no violation
	}
	ctx := violationCtx(t, ctxOpts{
		PfsThreshold:     2,
		MaxSlotAvailable: noSlotLimit,
		MinActiveCount:   noMinActive,
		PFS:              map[common.Address]uint64{addr2: 3},
		PfReport:         []common.Address{addr1},
	})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, ValExiting, out[addr1].State, "suspended + low staking → ValExiting")
	assert.Equal(t, ValExiting, out[addr2].State, "suspended + PFS severe → ValExiting")
	assert.Equal(t, ValActive, out[addr3].State, "non-suspended, no violation → unchanged")
}

// TestApplyViolationTransition_Deterministic verifies that violation transitions
// produce the same result regardless of map iteration order. Two ValActive
// below minStake compete for one ValExiting slot — the same one must always win.
func TestApplyViolationTransition_Deterministic(t *testing.T) {
	var firstResult map[common.Address]NodeState
	for i := range 100 {
		m := NodeMap{
			addr1: {State: ValActive, StakingAmount: belowMinStake},
			addr2: {State: ValActive, StakingAmount: belowMinStake},
			addr3: {State: ValActive, StakingAmount: aboveMinStake}, // keeps VA count > minActiveCount
		}
		// maxSlot=1, minActive=2 → only 1 of addr1/addr2 can transition
		ctx := violationCtx(t, ctxOpts{MaxSlotAvailable: 1, MinActiveCount: 2})
		out := ctx.applyViolationTransition(m)
		states := map[common.Address]NodeState{
			addr1: out[addr1].State,
			addr2: out[addr2].State,
			addr3: out[addr3].State,
		}
		if firstResult == nil {
			firstResult = states
		} else {
			require.Equal(t, firstResult, states, "nondeterministic result at iteration %d", i)
		}
	}
}

func TestApplyViolationTransition_MinStakeMigrated(t *testing.T) {
	cases := []struct {
		name           string
		validators     NodeMap
		maxSlot        uint64
		expectedStates map[common.Address]NodeState
		checkTimeout   *common.Address
	}{
		{
			"ValReady + low staking → ValInactive",
			NodeMap{
				addr1: {State: ValReady, StakingAmount: belowMinStake},
				addr2: {State: ValReady, StakingAmount: aboveMinStake},
			},
			noSlotLimit,
			map[common.Address]NodeState{addr1: ValInactive, addr2: ValReady},
			&addr1,
		},
		{
			"ValPaused + low staking → ValExiting",
			NodeMap{
				addr1: {State: ValPaused, StakingAmount: belowMinStake},
				addr2: {State: ValPaused, StakingAmount: aboveMinStake},
			},
			noSlotLimit,
			map[common.Address]NodeState{addr1: ValExiting, addr2: ValPaused},
			nil,
		},
		{
			"ValPaused + low staking + slot full → stays ValPaused",
			NodeMap{
				addr1: {State: ValPaused, StakingAmount: belowMinStake},
				addr2: {State: ValExiting, StakingAmount: aboveMinStake}, // slot occupied
				addr3: {State: ValActive, StakingAmount: aboveMinStake},
				addr4: {State: ValActive, StakingAmount: aboveMinStake},
			},
			1, // maxSlotAvailable=1, already occupied
			map[common.Address]NodeState{addr1: ValPaused, addr2: ValExiting, addr3: ValActive, addr4: ValActive},
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := violationCtx(t, ctxOpts{MaxSlotAvailable: tc.maxSlot, MinActiveCount: noMinActive})
			out := ctx.applyViolationTransition(tc.validators)
			for addr, expected := range tc.expectedStates {
				assert.Equal(t, expected, out[addr].State, "addr=%s", addr.Hex())
			}
			if tc.checkTimeout != nil {
				assert.False(t, out[*tc.checkTimeout].IdleTimeout.IsZero(), "IdleTimeout should be set")
			}
		})
	}
}

func TestApplyViolationTransition_NonActiveNotAffected(t *testing.T) {
	m := NodeMap{
		addr1: {State: CandReady, StakingAmount: belowMinStake},
		addr2: {State: ValInactive, StakingAmount: belowMinStake},
	}
	ctx := violationCtx(t, ctxOpts{MaxSlotAvailable: noSlotLimit, MinActiveCount: noMinActive})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, CandReady, out[addr1].State)
	assert.Equal(t, ValInactive, out[addr2].State)
}

func TestApplyViolationTransition_PFSAboveThreshold(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: aboveMinStake},
		addr2: {State: ValActive, StakingAmount: aboveMinStake},
	}
	ctx := violationCtx(t, ctxOpts{
		PfsThreshold:     2,
		MaxSlotAvailable: noSlotLimit,
		MinActiveCount:   noMinActive,
		PFS:              map[common.Address]uint64{addr1: 3, addr2: 1},
		PfReport:         []common.Address{addr1},
	})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, ValExiting, out[addr1].State, "PFS(3) >= threshold(2) → ValExiting (severe)")
	assert.Equal(t, ValPaused, out[addr2].State, "PFS(1) > 0, < threshold(2) → ValPaused (minor)")
}

func TestApplyViolationTransition_PFSNonActiveNotAffected(t *testing.T) {
	m := NodeMap{
		addr1: {State: CandReady, StakingAmount: aboveMinStake},
		addr2: {State: ValPaused, StakingAmount: aboveMinStake},
	}
	ctx := violationCtx(t, ctxOpts{
		PfsThreshold:     2,
		MaxSlotAvailable: noSlotLimit,
		MinActiveCount:   noMinActive,
		PFS:              map[common.Address]uint64{addr1: 10, addr2: 10},
		PfReport:         []common.Address{addr1},
	})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, CandReady, out[addr1].State, "CandReady not affected by PFS")
	assert.Equal(t, ValPaused, out[addr2].State, "ValPaused not affected by PFS")
}

func TestApplyViolationTransition_EmptyPfReport_NoOp(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, StakingAmount: aboveMinStake},
	}
	// Even if PFS map has an entry, no PfReport means rule2 doesn't fire.
	ctx := violationCtx(t, ctxOpts{
		PfsThreshold:     2,
		MaxSlotAvailable: noSlotLimit,
		MinActiveCount:   noMinActive,
		PFS:              map[common.Address]uint64{addr1: 100},
	})
	out := ctx.applyViolationTransition(m)
	assert.Equal(t, ValActive, out[addr1].State, "no PfReport → no PFS transition")
}

// ============================================================
// TestGetViolationTransition with SlotLimits
//
// Uses epochVACount=4 (4 validators at epoch start).
// SlotMath (see contracts/libraries/SlotMath.sol):
//
//	minActiveCount    = ceil(2*4/3)        → 3
//	totalBudget       = 4 - 3              → 1
//	maxSlotAvailable  = ceil(1/2)          → 1  (up to 1 ValPaused, up to 1 ValExiting independently)
//
// With n=4, maxSlot=1 for each but minActive=3 means only 1 total can leave ValActive
// (either 1 paused or 1 exiting, not both), since 4-2=2 < minActive=3.
// ============================================================
func TestApplyViolationTransition_SlotLimits(t *testing.T) {
	const (
		slotMax      = uint64(1) // maxSlotAvailable(4)
		minActive    = uint64(3) // minActiveCount(4)
		pfsThreshold = uint64(2)
	)

	t.Run("PFS minor: only 1 paused when 2 violate", func(t *testing.T) {
		m := NodeMap{
			addr1: {State: ValActive, StakingAmount: aboveMinStake},
			addr2: {State: ValActive, StakingAmount: aboveMinStake},
			addr3: {State: ValActive, StakingAmount: aboveMinStake},
			addr4: {State: ValActive, StakingAmount: aboveMinStake},
		}
		ctx := violationCtx(t, ctxOpts{
			PfsThreshold:     pfsThreshold,
			MaxSlotAvailable: slotMax,
			MinActiveCount:   minActive,
			PFS:              map[common.Address]uint64{addr1: pfsThreshold - 1, addr2: pfsThreshold - 1},
			PfReport:         []common.Address{addr1},
		})
		out := ctx.applyViolationTransition(m)
		assert.Equal(t, 1, int(out.CountByState(ValPaused)), "only 1 slot for ValPaused")
		assert.Equal(t, 3, int(out.CountByState(ValActive)), "remaining 3 stay ValActive")
	})

	t.Run("PFS severe: only 1 exited when 2 violate", func(t *testing.T) {
		m := NodeMap{
			addr1: {State: ValActive, StakingAmount: aboveMinStake},
			addr2: {State: ValActive, StakingAmount: aboveMinStake},
			addr3: {State: ValActive, StakingAmount: aboveMinStake},
			addr4: {State: ValActive, StakingAmount: aboveMinStake},
		}
		ctx := violationCtx(t, ctxOpts{
			PfsThreshold:     pfsThreshold,
			MaxSlotAvailable: slotMax,
			MinActiveCount:   minActive,
			PFS:              map[common.Address]uint64{addr1: pfsThreshold + 1, addr2: pfsThreshold + 1},
			PfReport:         []common.Address{addr1},
		})
		out := ctx.applyViolationTransition(m)
		assert.Equal(t, 1, int(out.CountByState(ValExiting)), "only 1 slot for ValExiting")
		assert.Equal(t, 3, int(out.CountByState(ValActive)), "remaining 3 stay ValActive")
	})

	t.Run("minStake violation: skip when slot full", func(t *testing.T) {
		m := NodeMap{
			addr1: {State: ValActive, StakingAmount: belowMinStake},
			addr2: {State: ValActive, StakingAmount: belowMinStake},
			addr3: {State: ValActive, StakingAmount: aboveMinStake},
			addr4: {State: ValActive, StakingAmount: aboveMinStake},
		}
		ctx := violationCtx(t, ctxOpts{
			MaxSlotAvailable: slotMax,
			MinActiveCount:   minActive,
		})
		out := ctx.applyViolationTransition(m)
		assert.Equal(t, 1, int(out.CountByState(ValExiting)), "only 1 can exit")
		assert.Equal(t, 3, int(out.CountByState(ValActive)), "3 remain ValActive")
	})

	t.Run("slot already occupied: no more transitions", func(t *testing.T) {
		m := NodeMap{
			addr1: {State: ValActive, StakingAmount: aboveMinStake},
			addr2: {State: ValPaused, StakingAmount: aboveMinStake}, // already paused
			addr3: {State: ValActive, StakingAmount: aboveMinStake},
			addr4: {State: ValActive, StakingAmount: aboveMinStake},
		}
		ctx := violationCtx(t, ctxOpts{
			PfsThreshold:     pfsThreshold,
			MaxSlotAvailable: slotMax,
			MinActiveCount:   minActive,
			PFS:              map[common.Address]uint64{addr1: pfsThreshold - 1},
			PfReport:         []common.Address{addr1},
		})
		out := ctx.applyViolationTransition(m)
		assert.Equal(t, ValActive, out[addr1].State, "paused slot full, minor violation skipped")
	})

	// n=10: SlotMath:
	//   minActiveCount   = ceil(2*10/3)       → 7
	//   totalBudget      = 10 - 7             → 3
	//   maxSlotAvailable = ceil(3/2)          → 2
	// Up to 3 can leave ValActive (2 paused + 1 exiting or vice versa), but minActive=7 caps total.
	t.Run("n=10", func(t *testing.T) {
		const (
			slotMax10      = uint64(2) // maxSlotAvailable(10)
			minActive10    = uint64(7) // minActiveCount(10)
			pfsThreshold10 = uint64(2)
		)
		addrs := [10]common.Address{
			addr1, addr2, addr3, addr4, addr5,
			common.HexToAddress("0x0006"),
			common.HexToAddress("0x0007"),
			common.HexToAddress("0x0008"),
			common.HexToAddress("0x0009"),
			common.HexToAddress("0x000a"),
		}

		// 4 violators: 2 severe + 2 minor. budget=3, maxSlot=2.
		// addrs[0]: severe → VE (exitSlot 1/2, VA 10→9)
		// addrs[1]: severe → VE (exitSlot 2/2, VA 9→8)
		// addrs[2]: minor  → VP (pauseSlot 1/2, VA 8→7)
		// addrs[3]: minor  → blocked (VA=7 = minActive, canDemoteActive fails)
		t.Run("PFS 4 violators: 3 transition, 4th blocked by minActive", func(t *testing.T) {
			m := NodeMap{}
			for _, addr := range addrs {
				m[addr] = &Node{State: ValActive, StakingAmount: aboveMinStake}
			}
			ctx := violationCtx(t, ctxOpts{
				PfsThreshold:     pfsThreshold10,
				MaxSlotAvailable: slotMax10,
				MinActiveCount:   minActive10,
				PFS: map[common.Address]uint64{
					addrs[0]: pfsThreshold10,     // severe → ValExiting
					addrs[1]: pfsThreshold10,     // severe → ValExiting
					addrs[2]: pfsThreshold10 - 1, // minor  → ValPaused
					addrs[3]: pfsThreshold10 - 1, // minor  → blocked by minActive
				},
				PfReport: []common.Address{addr1},
			})
			out := ctx.applyViolationTransition(m)

			assert.Equal(t, 2, int(out.CountByState(ValExiting)), "2 severe → 2 exiting (maxSlot=2)")
			assert.Equal(t, 1, int(out.CountByState(ValPaused)), "1 minor paused, 2nd blocked by minActive")
			assert.Equal(t, 7, int(out.CountByState(ValActive)), "7 remain active (minActive=7)")
		})

		t.Run("3 belowMinStake: 2 exit (exitSlot cap=2), 3rd blocked", func(t *testing.T) {
			m := NodeMap{}
			for i, addr := range addrs {
				stake := aboveMinStake
				if i < 3 {
					stake = belowMinStake // 3 validators below min stake
				}
				m[addr] = &Node{State: ValActive, StakingAmount: stake}
			}
			ctx := violationCtx(t, ctxOpts{
				MaxSlotAvailable: slotMax10,
				MinActiveCount:   minActive10,
			})
			out := ctx.applyViolationTransition(m)

			assert.Equal(t, 2, int(out.CountByState(ValExiting)), "maxSlot=2 → 2 can exit")
			assert.Equal(t, 8, int(out.CountByState(ValActive)), "8 remain active")
		})
	})
}

//#endregion

//#region applyTimeoutTransition

func TestApplyTimeoutTransition(t *testing.T) {
	expiredTimeout := testBlockTime.Add(-1 * time.Hour)
	futureTimeout := testBlockTime.Add(10 * 24 * time.Hour)

	cases := []struct {
		name                string
		input               *Node
		expectedState       NodeState
		expectIdleSet       bool
		expectPausedSet     bool
		expectIdlePreserved *time.Time
	}{
		{"ValInactive: set idle timeout", &Node{State: ValInactive}, ValInactive, true, false, nil},
		{"ValReady: set idle timeout", &Node{State: ValReady}, ValReady, true, false, nil},
		{"ValInactive: preserve existing idle timeout", &Node{State: ValInactive, IdleTimeout: futureTimeout}, ValInactive, true, false, &futureTimeout},
		{"ValReady: preserve existing idle timeout", &Node{State: ValReady, IdleTimeout: futureTimeout}, ValReady, true, false, &futureTimeout},
		{"ValInactive: idle expired → Registered", &Node{State: ValInactive, IdleTimeout: expiredTimeout}, Registered, false, false, nil},
		{"ValReady: idle expired → Registered", &Node{State: ValReady, IdleTimeout: expiredTimeout}, Registered, false, false, nil},
		{"ValPaused: set paused timeout", &Node{State: ValPaused}, ValPaused, false, true, nil},
		{"ValPaused: preserve existing paused timeout", &Node{State: ValPaused, PausedTimeout: futureTimeout}, ValPaused, false, true, nil},
		{"ValPaused: paused expired → ValInactive + idle set", &Node{State: ValPaused, PausedTimeout: expiredTimeout}, ValInactive, true, false, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NodeMap{addr1: tc.input}
			ctx := timeoutCtx(t)
			out := ctx.applyTimeoutTransition(m)
			r := out[addr1]

			assert.Equal(t, tc.expectedState, r.State)
			assert.Equal(t, tc.expectIdleSet, !r.IdleTimeout.IsZero(), "IdleTimeout")
			assert.Equal(t, tc.expectPausedSet, !r.PausedTimeout.IsZero(), "PausedTimeout")
			if tc.expectIdlePreserved != nil {
				assert.Equal(t, tc.expectIdlePreserved.Unix(), r.IdleTimeout.Unix())
			}
		})
	}
}

// TestApplyTimeoutTransition_DefaultClearsAllTimeouts verifies that the default
// branch (states that don't set timeouts) zeroes any pre-existing timeouts.
func TestApplyTimeoutTransition_DefaultClearsAllTimeouts(t *testing.T) {
	now := time.Now()
	m := NodeMap{
		addr1: {State: ValActive, IdleTimeout: now, PausedTimeout: now},
		addr2: {State: CandReady, IdleTimeout: now, PausedTimeout: now},
		addr3: {State: Registered, IdleTimeout: now, PausedTimeout: now},
		addr4: {State: CandTesting, IdleTimeout: now, PausedTimeout: now},
		addr5: {State: ValExiting, IdleTimeout: now, PausedTimeout: now},
	}
	ctx := timeoutCtx(t)
	out := ctx.applyTimeoutTransition(m)
	for addr, r := range out {
		assert.True(t, r.IdleTimeout.IsZero(), "IdleTimeout should be cleared for %s", addr.Hex())
		assert.True(t, r.PausedTimeout.IsZero(), "PausedTimeout should be cleared for %s", addr.Hex())
	}
}

//#endregion

//#region ApplyAllTransitions

func TestApplyAllTransitions(t *testing.T) {
	expiredTimeout := testBlockTime.Add(-1 * time.Hour)

	cases := []struct {
		name                 string
		isEpoch              bool
		input                NodeMap
		expected             map[common.Address]NodeState
		expectedEpochVACount uint64
	}{
		{
			name:    "epoch pipeline demotes below-min before violation",
			isEpoch: true,
			input: NodeMap{
				addr1: {State: ValActive, StakingAmount: belowMinStake},
				addr2: {State: ValActive, StakingAmount: aboveMinStake},
			},
			expected: map[common.Address]NodeState{
				addr1: ValInactive,
				addr2: ValActive,
			},
			expectedEpochVACount: 1,
		},
		{
			name:    "non-epoch violation fires while epoch transition is skipped",
			isEpoch: false,
			input: NodeMap{
				addr1: {State: ValActive, StakingAmount: belowMinStake},
				addr2: {State: ValActive, StakingAmount: aboveMinStake},
			},
			expected: map[common.Address]NodeState{
				addr1: ValExiting,
				addr2: ValActive,
			},
		},
		{
			name:    "non-epoch last ValActive protected by minActiveCount",
			isEpoch: false,
			input: NodeMap{
				addr1: {State: ValActive, StakingAmount: belowMinStake},
			},
			expected: map[common.Address]NodeState{
				addr1: ValActive,
			},
		},
		{
			name:    "candidate promotion only happens at epoch",
			isEpoch: true,
			input: NodeMap{
				addr1: {State: CandReady, StakingAmount: aboveMinStake},
			},
			expected: map[common.Address]NodeState{
				addr1: CandTesting,
			},
		},
		{
			name:    "mixed states at epoch run through epoch and timeout",
			isEpoch: true,
			input: NodeMap{
				addr1: {State: ValActive, StakingAmount: belowMinStake},
				addr2: {State: ValActive, StakingAmount: aboveMinStake},
				addr3: {State: CandReady, StakingAmount: aboveMinStake},
				addr4: {State: ValPaused, StakingAmount: aboveMinStake},
				addr5: {State: CandReady, StakingAmount: belowMinStake},
			},
			expected: map[common.Address]NodeState{
				addr1: ValInactive,
				addr2: ValActive,
				addr3: CandTesting,
				addr4: ValPaused,
				addr5: Registered,
			},
			expectedEpochVACount: 1,
		},
		{
			name:    "timeout fires after transition pipeline",
			isEpoch: true,
			input: NodeMap{
				addr1: {State: ValInactive, StakingAmount: aboveMinStake, IdleTimeout: expiredTimeout},
			},
			expected: map[common.Address]NodeState{
				addr1: Registered,
			},
		},
		{
			name:    "expired pause transitions to ValInactive",
			isEpoch: false,
			input: NodeMap{
				addr1: {State: ValPaused, StakingAmount: aboveMinStake, PausedTimeout: expiredTimeout},
			},
			expected: map[common.Address]NodeState{
				addr1: ValInactive,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := buildCtx(t, ctxOpts{
				MinStake:                testMinStake,
				IdleTimeout:             testIdleTimeout,
				PauseTimeout:            testPauseTimeout,
				MaxValActivePausedCount: testMaxValCount,
				MaxSlotAvailable:        noSlotLimit,
				MinActiveCount:          1,
			})
			ctx.IsEpoch = tc.isEpoch

			result := ctx.ApplyAllTransitions(tc.input)
			for addr, expectedState := range tc.expected {
				assert.Equal(t, expectedState, result.Nodes[addr].State, "addr=%s", addr.Hex())
			}
			assert.Equal(t, tc.expectedEpochVACount, result.epochVACountForWrite)
		})
	}
}

//#endregion
