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
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = common.HexToAddress("0x0001")
	addr2 = common.HexToAddress("0x0002")
	addr3 = common.HexToAddress("0x0003")
	addr4 = common.HexToAddress("0x0004")
	addr5 = common.HexToAddress("0x0005")
)

const (
	testMinStake    = uint64(5_000_000)
	aboveMinStake   = uint64(10_000_000)
	belowMinStake   = uint64(4_999_999)
	testVRankEpoch  = uint64(10)
	testEpochNum    = testVRankEpoch     // testVRankEpoch % testVRankEpoch == 0
	testNonEpochNum = testVRankEpoch - 1 // (testVRankEpoch-1) % testVRankEpoch != 0
	testMaxValCount = 50

	highStake = uint64(30_000_000)

	testCFSBlockNum  = uint64(100)
	testCFSThreshold = uint64(10)

	// noSlotLimit disables slot checks in tests that don't test slot math.
	noSlotLimit = uint64(100)
	noMinActive = uint64(0)
)

var (
	testBlockTime = time.Unix(1700000000, 0) // fixed block timestamp for deterministic tests
	midStake      = uint64(20_000_000)
	lowStake      = uint64(10_000_000)

	testIdleTimeout    = 24 * time.Hour
	testPauseTimeout   = 8 * time.Hour
	defaultIdleTimeout = 30 * 24 * time.Hour
)

func newTestValsetModule(ctrl *gomock.Controller) *ValsetModule {
	mockGov := gov_mock.NewMockGovModule(ctrl)
	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: new(big.Int).SetUint64(testMinStake),
	}).AnyTimes()
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(&params.ChainConfig{VRankEpoch: testVRankEpoch}).AnyTimes()
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetCFSWithSlotFactor(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	return &ValsetModule{
		InitOpts: InitOpts{
			Chain:       mockChain,
			GovModule:   mockGov,
			VRankModule: mockVRank,
		},
	}
}

// ============================================================
// TestGetEpochTransition
// ============================================================

func TestGetEpochTransition_StateTransitions(t *testing.T) {
	testcases := []struct {
		name          string
		inputState    valset.State
		stakingAmount uint64
		expectedState valset.State
		expectTimeout bool // IdleTimeout should be set
	}{
		{"T1: ValExiting → ValInactive", valset.ValExiting, aboveMinStake, valset.ValInactive, true},
		{"T4a: CandReady + stake>=min → CandTesting", valset.CandReady, aboveMinStake, valset.CandTesting, false},
		{"T4b: CandReady + stake<min → Registered", valset.CandReady, belowMinStake, valset.Registered, false},
		{"T3a: CandTesting + stake>=min → ValActive", valset.CandTesting, aboveMinStake, valset.ValActive, false},
		{"T3b: CandTesting + stake<min → ValInactive", valset.CandTesting, belowMinStake, valset.ValInactive, true},
		{"ValReady + stake>=min → ValActive", valset.ValReady, aboveMinStake, valset.ValActive, false},
		{"ValActive + stake>=min → ValActive", valset.ValActive, aboveMinStake, valset.ValActive, false},
		{"ValPaused + stake>=min → ValPaused (preserved)", valset.ValPaused, aboveMinStake, valset.ValPaused, false},
		{"T3b: ValActive + stake<min → ValInactive", valset.ValActive, belowMinStake, valset.ValInactive, true},
		{"T3b: ValReady + stake<min → ValInactive", valset.ValReady, belowMinStake, valset.ValInactive, true},
		{"T3b: ValPaused + stake<min → ValInactive", valset.ValPaused, belowMinStake, valset.ValInactive, true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)

			validators := valset.NodeStateMap{
				addr1: {State: tc.inputState, StakingAmount: tc.stakingAmount},
			}
			result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, testMaxValCount, testBlockTime, 0, 0, 0)
			assert.Equal(t, tc.expectedState, result[addr1].State)
			if tc.expectTimeout {
				assert.False(t, result[addr1].IdleTimeout.IsZero())
			}
		})
	}
}

func TestGetEpochTransition_MaxValidatorCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: highStake},
		addr2: {State: valset.ValActive, StakingAmount: midStake},
		addr3: {State: valset.ValActive, StakingAmount: lowStake},
	}
	result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, 2, testBlockTime, 0, 0, 0)
	assert.Equal(t, valset.ValActive, result[addr1].State)
	assert.Equal(t, valset.ValActive, result[addr2].State)
	assert.Equal(t, valset.ValInactive, result[addr3].State)
	assert.False(t, result[addr3].IdleTimeout.IsZero())
}

func TestGetEpochTransition_TieBreakingByAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	// Same stake → lower address wins
	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: lowStake},
		addr2: {State: valset.ValActive, StakingAmount: lowStake},
		addr3: {State: valset.ValActive, StakingAmount: lowStake},
	}
	result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, 2, testBlockTime, 0, 0, 0)

	// addr1 (0x0001) < addr2 (0x0002) < addr3 (0x0003)
	assert.Equal(t, valset.ValActive, result[addr1].State)
	assert.Equal(t, valset.ValActive, result[addr2].State)
	assert.Equal(t, valset.ValInactive, result[addr3].State)
}

func TestGetEpochTransition_DoesNotMutateInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValExiting, StakingAmount: aboveMinStake},
	}
	result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, testMaxValCount, testBlockTime, 0, 0, 0)
	assert.Equal(t, valset.ValInactive, result[addr1].State)    // result is transitioned
	assert.Equal(t, valset.ValExiting, validators[addr1].State) // original is unchanged
}

// ============================================================
// TestGetFallbackTransition
// ============================================================

func TestGetFallbackTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl) // mockVRank returns nil scores → isPassVrankTest passes

	pausedTimeout := testBlockTime.Add(8 * time.Hour)
	// Staking amounts differ so the cap selects by stake rank.
	// Rank: addr3(highStake) > addr2(midStake) > addr1,addr4(lowStake, addr1<addr4 by address)
	validators := valset.NodeStateMap{
		addr1: {State: valset.CandTesting, StakingAmount: lowStake},
		addr2: {State: valset.ValReady, StakingAmount: midStake},
		addr3: {State: valset.ValActive, StakingAmount: highStake},
		addr4: {State: valset.ValPaused, StakingAmount: lowStake, PausedTimeout: pausedTimeout},
		addr5: {State: valset.Registered, StakingAmount: highStake},
	}

	// maxValidatorCount=3: top 3 are addr3, addr2, addr1 → ValActive; addr4 → ValInactive.
	result := v.getFallbackTransition(validators, testBlockTime, testIdleTimeout, 0, 0, 0, 3)

	assert.Equal(t, valset.ValActive, result[addr1].State) // CandTesting passes vrank (nil scores)
	assert.Equal(t, valset.ValActive, result[addr2].State)
	assert.Equal(t, valset.ValActive, result[addr3].State)
	assert.Equal(t, valset.ValInactive, result[addr4].State)                            // outside cap → ValInactive
	assert.True(t, result[addr4].IdleTimeout.Equal(testBlockTime.Add(testIdleTimeout))) // idleTimeout set
	assert.Equal(t, valset.Registered, result[addr5].State)                             // non-competition, unchanged
	// Input not mutated.
	assert.Equal(t, valset.ValPaused, validators[addr4].State)
}

// ============================================================
// TestCompareValidatorRank
// ============================================================

func TestCompareValidatorRank(t *testing.T) {
	sv := func(addr common.Address, stake uint64) sortableValidator {
		return sortableValidator{addr, &valset.ValidatorState{StakingAmount: stake}}
	}

	// Higher stake ranks first (negative result means a < b in sort order).
	assert.Negative(t, compareValidatorRank(sv(addr1, highStake), sv(addr2, lowStake)))
	assert.Positive(t, compareValidatorRank(sv(addr1, lowStake), sv(addr2, highStake)))

	// Equal stake: lower address ranks first.
	assert.Negative(t, compareValidatorRank(sv(addr1, highStake), sv(addr2, highStake)))
	assert.Positive(t, compareValidatorRank(sv(addr2, highStake), sv(addr1, highStake)))

	// Equal stake and address: tie.
	assert.Zero(t, compareValidatorRank(sv(addr1, highStake), sv(addr1, highStake)))
}

// ============================================================
// TestIsPassVrankTest (CFS threshold check)
// ============================================================

func TestIsPassVrankTest(t *testing.T) {
	testcases := []struct {
		name         string
		cfsThreshold uint64
		cfsScores    map[common.Address]uint64
		cfsErr       error
		expected     bool
	}{
		{
			"CFS below threshold → pass",
			testCFSThreshold,
			map[common.Address]uint64{addr1: testCFSThreshold - 5},
			nil, true,
		},
		{
			"CFS equals threshold → fail",
			testCFSThreshold,
			map[common.Address]uint64{addr1: testCFSThreshold},
			nil, false,
		},
		{
			"CFS above threshold → fail",
			testCFSThreshold,
			map[common.Address]uint64{addr1: testCFSThreshold + 5},
			nil, false,
		},
		{
			"addr not in CFS scores → pass",
			testCFSThreshold,
			map[common.Address]uint64{addr2: testCFSThreshold + 5},
			nil, true,
		},
		{
			"GetCFS error → pass (fail-open)",
			testCFSThreshold,
			nil, errors.New("some error"), true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)
			mockVRank := vrank_mock.NewMockVRankModule(ctrl)
			mockVRank.EXPECT().GetCFSWithSlotFactor(testCFSBlockNum, uint64(0)).Return(tc.cfsScores, tc.cfsErr)
			v.VRankModule = mockVRank

			result := v.isPassVrankTest(addr1, testCFSBlockNum, tc.cfsThreshold, 0)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// ============================================================
// TestGetEpochTransition with CFS
// ============================================================

func TestGetEpochTransition_CandTestingCFS(t *testing.T) {
	testcases := []struct {
		name          string
		cfs           uint64
		expectedState valset.State
	}{
		{"CFS above threshold → Registered", testCFSThreshold + 5, valset.Registered},
		{"CFS below threshold → ValActive", testCFSThreshold - 5, valset.ValActive},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)
			mockVRank := vrank_mock.NewMockVRankModule(ctrl)
			mockVRank.EXPECT().GetCFSWithSlotFactor(testCFSBlockNum, uint64(0)).Return(map[common.Address]uint64{addr1: tc.cfs}, nil)
			v.VRankModule = mockVRank

			validators := valset.NodeStateMap{
				addr1: {State: valset.CandTesting, StakingAmount: aboveMinStake},
			}
			result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, testMaxValCount, testBlockTime, testCFSBlockNum, testCFSThreshold, 0)
			assert.Equal(t, tc.expectedState, result[addr1].State)
		})
	}
}

// ============================================================
// TestGetTimeoutTransition
// ============================================================

func TestGetTimeoutTransition(t *testing.T) {
	expiredTimeout := testBlockTime.Add(-1 * time.Hour)
	futureTimeout := testBlockTime.Add(10 * 24 * time.Hour)

	testcases := []struct {
		name                string
		input               *valset.ValidatorState
		expectedState       valset.State
		expectIdleSet       bool // IdleTimeout should be non-zero
		expectPausedSet     bool // PausedTimeout should be non-zero
		expectIdlePreserved *time.Time
	}{
		{
			"ValInactive: set idle timeout",
			&valset.ValidatorState{State: valset.ValInactive},
			valset.ValInactive, true, false, nil,
		},
		{
			"ValReady: set idle timeout",
			&valset.ValidatorState{State: valset.ValReady},
			valset.ValReady, true, false, nil,
		},
		{
			"ValInactive: preserve existing idle timeout",
			&valset.ValidatorState{State: valset.ValInactive, IdleTimeout: futureTimeout},
			valset.ValInactive, true, false, &futureTimeout,
		},
		{
			"ValReady: preserve existing idle timeout",
			&valset.ValidatorState{State: valset.ValReady, IdleTimeout: futureTimeout},
			valset.ValReady, true, false, &futureTimeout,
		},
		{
			"ValInactive: idle expired → Registered",
			&valset.ValidatorState{State: valset.ValInactive, IdleTimeout: expiredTimeout},
			valset.Registered, false, false, nil,
		},
		{
			"ValReady: idle expired → Registered",
			&valset.ValidatorState{State: valset.ValReady, IdleTimeout: expiredTimeout},
			valset.Registered, false, false, nil,
		},
		{
			"ValPaused: set paused timeout",
			&valset.ValidatorState{State: valset.ValPaused},
			valset.ValPaused, false, true, nil,
		},
		{
			"ValPaused: preserve existing paused timeout",
			&valset.ValidatorState{State: valset.ValPaused, PausedTimeout: futureTimeout},
			valset.ValPaused, false, true, nil,
		},
		{
			"ValPaused: paused expired → ValInactive + idle set",
			&valset.ValidatorState{State: valset.ValPaused, PausedTimeout: expiredTimeout},
			valset.ValInactive, true, false, nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)

			validators := valset.NodeStateMap{addr1: tc.input}
			result := v.getTimeoutTransition(validators, defaultIdleTimeout, testPauseTimeout, testBlockTime)
			r := result[addr1]

			assert.Equal(t, tc.expectedState, r.State)
			assert.Equal(t, tc.expectIdleSet, !r.IdleTimeout.IsZero(), "IdleTimeout")
			assert.Equal(t, tc.expectPausedSet, !r.PausedTimeout.IsZero(), "PausedTimeout")
			if tc.expectIdlePreserved != nil {
				assert.Equal(t, tc.expectIdlePreserved.Unix(), r.IdleTimeout.Unix())
			}
		})
	}
}

func TestGetTimeoutTransition_DefaultClearsAllTimeouts(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	now := time.Now()
	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, IdleTimeout: now, PausedTimeout: now},
		addr2: {State: valset.CandReady, IdleTimeout: now, PausedTimeout: now},
		addr3: {State: valset.Registered, IdleTimeout: now, PausedTimeout: now},
		addr4: {State: valset.CandTesting, IdleTimeout: now, PausedTimeout: now},
		addr5: {State: valset.ValExiting, IdleTimeout: now, PausedTimeout: now},
	}
	result := v.getTimeoutTransition(validators, defaultIdleTimeout, testPauseTimeout, testBlockTime)
	for addr, r := range result {
		assert.True(t, r.IdleTimeout.IsZero(), "IdleTimeout should be cleared for %s", addr.Hex())
		assert.True(t, r.PausedTimeout.IsZero(), "PausedTimeout should be cleared for %s", addr.Hex())
	}
}

// ============================================================
// TestGetViolationTransitionLessMinStakingAmount
// ============================================================

func TestGetViolationTransition_ValActiveBelowMinStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
	v.VRankModule = mockVRank

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
		addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
	}
	result := v.getViolationTransition(testMinStake, validators, 0, 0, noSlotLimit, noMinActive)
	assert.Equal(t, valset.ValExiting, result[addr1].State)
	assert.Equal(t, valset.ValActive, result[addr2].State)
}

func TestGetViolationTransition_NonActiveNotAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
	v.VRankModule = mockVRank

	validators := valset.NodeStateMap{
		addr1: {State: valset.CandReady, StakingAmount: belowMinStake},
		addr2: {State: valset.ValInactive, StakingAmount: belowMinStake},
	}
	result := v.getViolationTransition(testMinStake, validators, 0, 0, noSlotLimit, noMinActive)
	assert.Equal(t, valset.CandReady, result[addr1].State)
	assert.Equal(t, valset.ValInactive, result[addr2].State)
}

func TestGetViolationTransition_PFSAboveThreshold(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
	mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr1: 3, addr2: 1}, nil)
	v.VRankModule = mockVRank

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: aboveMinStake},
		addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
	}
	result := v.getViolationTransition(testMinStake, validators, 100, 2, noSlotLimit, noMinActive)
	assert.Equal(t, valset.ValExiting, result[addr1].State, "PFS(3) >= threshold(2) → ValExiting (severe)")
	assert.Equal(t, valset.ValPaused, result[addr2].State, "PFS(1) > 0, < threshold(2) → ValPaused (minor)")
}

func TestGetViolationTransition_PFSNonActiveNotAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
	mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr1: 10, addr2: 10}, nil)
	v.VRankModule = mockVRank

	validators := valset.NodeStateMap{
		addr1: {State: valset.CandReady, StakingAmount: aboveMinStake},
		addr2: {State: valset.ValPaused, StakingAmount: aboveMinStake},
	}
	result := v.getViolationTransition(testMinStake, validators, 100, 2, noSlotLimit, noMinActive)
	assert.Equal(t, valset.CandReady, result[addr1].State, "CandReady not affected by PFS")
	assert.Equal(t, valset.ValPaused, result[addr2].State, "ValPaused not affected by PFS")
}

// ============================================================
// TestGetViolationTransition with SlotLimits
//
// Uses slotFactor=4 (4 validators at epoch start).
// SlotMath (see contracts/libraries/SlotMath.sol):
//   f(n)              = (n-1)/3            → f(4)=1
//   maxSlotAvailable  = max(1, f(n)/2)     → max(1, 0)=1  (up to 1 ValPaused, up to 1 ValExiting independently)
//   minActiveCount    = max(2f+1, n-2*max) → max(3, 2)=3  (at least 3 ValActive required)
//
// With n=4, maxSlot=1 for each but minActive=3 means only 1 total can leave ValActive
// (either 1 paused or 1 exiting, not both), since 4-2=2 < minActive=3.
// ============================================================

func countState(validators valset.NodeStateMap, state valset.State) int {
	count := 0
	for _, val := range validators {
		if val.State == state {
			count++
		}
	}
	return count
}

func TestGetViolationTransition_SlotLimits(t *testing.T) {
	const (
		slotMax      = uint64(1) // maxSlotAvailable(4)
		minActive    = uint64(3) // minActiveCount(4)
		pfsThreshold = uint64(2)
	)

	t.Run("PFS minor: only 1 paused when 2 violate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		v := newTestValsetModule(ctrl)
		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
		mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr1: pfsThreshold - 1, addr2: pfsThreshold - 1}, nil)
		v.VRankModule = mockVRank

		validators := valset.NodeStateMap{
			addr1: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr3: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr4: {State: valset.ValActive, StakingAmount: aboveMinStake},
		}
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive)

		assert.Equal(t, 1, countState(result, valset.ValPaused), "only 1 slot available for ValPaused")
		assert.Equal(t, 3, countState(result, valset.ValActive), "remaining 3 should stay ValActive")
	})

	t.Run("PFS severe: only 1 exited when 2 violate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		v := newTestValsetModule(ctrl)
		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
		mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr1: pfsThreshold + 1, addr2: pfsThreshold + 1}, nil)
		v.VRankModule = mockVRank

		validators := valset.NodeStateMap{
			addr1: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr3: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr4: {State: valset.ValActive, StakingAmount: aboveMinStake},
		}
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive)

		assert.Equal(t, 1, countState(result, valset.ValExiting), "only 1 slot available for ValExiting")
		assert.Equal(t, 3, countState(result, valset.ValActive), "remaining 3 should stay ValActive")
	})

	t.Run("minStake violation: skip when slot full", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		v := newTestValsetModule(ctrl)
		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
		v.VRankModule = mockVRank

		validators := valset.NodeStateMap{
			addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
			addr2: {State: valset.ValActive, StakingAmount: belowMinStake},
			addr3: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr4: {State: valset.ValActive, StakingAmount: aboveMinStake},
		}
		result := v.getViolationTransition(testMinStake, validators, 0, 0, slotMax, minActive)

		assert.Equal(t, 1, countState(result, valset.ValExiting), "only 1 can exit due to slot limit")
		assert.Equal(t, 3, countState(result, valset.ValActive), "remaining 3 should stay ValActive")
	})

	t.Run("slot already occupied: no more transitions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		v := newTestValsetModule(ctrl)
		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
		mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr1: pfsThreshold - 1}, nil)
		v.VRankModule = mockVRank

		validators := valset.NodeStateMap{
			addr1: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr2: {State: valset.ValPaused, StakingAmount: aboveMinStake}, // already paused
			addr3: {State: valset.ValActive, StakingAmount: aboveMinStake},
			addr4: {State: valset.ValActive, StakingAmount: aboveMinStake},
		}
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive)
		assert.Equal(t, valset.ValActive, result[addr1].State, "paused slot full, minor violation skipped")
	})
}

// ============================================================
// TestApplyAllTransitions
// ============================================================

// newTestApplyAllTransitions creates a ValsetModule with default timeouts and max counts.
func newTestApplyAllTransitions(ctrl *gomock.Controller) *ValsetModule {
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockGov := gov_mock.NewMockGovModule(ctrl)

	chainConfig := &params.ChainConfig{VRankEpoch: testVRankEpoch}
	mockChain.EXPECT().Config().Return(chainConfig).AnyTimes()

	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: new(big.Int).SetUint64(testMinStake),
	}).AnyTimes()

	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
	mockVRank.EXPECT().GetPFS(gomock.Any()).Return(nil, nil).AnyTimes()
	mockVRank.EXPECT().GetCFSWithSlotFactor(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	v := &ValsetModule{
		InitOpts: InitOpts{
			Chain:       mockChain,
			GovModule:   mockGov,
			VRankModule: mockVRank,
		},
	}
	return v
}

func TestApplyAllTransitions(t *testing.T) {
	testcases := []struct {
		name     string
		num      uint64
		input    valset.NodeStateMap
		expected map[common.Address]valset.State
	}{
		{
			// addr1: ValActive+belowMin → violation(ValExiting) → timeout(noop) → epoch(ValInactive)
			// addr2: ValActive+aboveMin → violation(noop) → timeout(noop) → epoch(ValActive)
			"violation→epoch pipeline",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive, addr2: valset.ValActive},
		},
		{
			// addr1: ValActive+belowMin → violation(ValExiting) → timeout(noop) → epoch(noop, non-epoch)
			// Proves: violation applies regardless of epoch, but epoch transition only fires at VRankEpoch
			"non-epoch: violation fires but epoch noop",
			testNonEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValExiting},
		},
		{
			// addr1: CandReady+aboveMin → violation(noop) → timeout(noop) → epoch(CandTesting)
			// addr2: ValActive+aboveMin → epoch(ValActive) — prevents fallback from firing
			"candidate promotion at epoch",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.CandReady, StakingAmount: aboveMinStake},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.CandTesting, addr2: valset.ValActive},
		},
		{
			// addr1: ValActive+belowMin  → violation(ValExiting) → timeout(noop) → epoch(ValInactive)
			// addr2: ValActive+aboveMin  → violation(noop) → timeout(noop) → epoch(ValActive)
			// addr3: CandReady+aboveMin  → violation(noop) → timeout(noop) → epoch(CandTesting)
			// addr4: ValPaused+aboveMin  → violation(noop) → timeout(set PausedTimeout) → epoch(ValPaused preserved)
			// addr5: CandReady+belowMin  → violation(noop) → timeout(noop) → epoch(Registered)
			"mixed states at epoch",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
				addr3: {State: valset.CandReady, StakingAmount: aboveMinStake},
				addr4: {State: valset.ValPaused, StakingAmount: aboveMinStake},
				addr5: {State: valset.CandReady, StakingAmount: belowMinStake},
			},
			map[common.Address]valset.State{
				addr1: valset.ValInactive,
				addr2: valset.ValActive,
				addr3: valset.CandTesting,
				addr4: valset.ValPaused,
				addr5: valset.Registered,
			},
		},
		{
			// addr1: ValInactive+expiredIdle → violation(noop) → timeout(Registered) → epoch(noop)
			// addr2: ValActive+aboveMin → epoch(ValActive) — prevents fallback from firing
			"timeout fires: expired idle → Registered",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValInactive, StakingAmount: aboveMinStake, IdleTimeout: testBlockTime.Add(-1 * time.Hour)},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.Registered, addr2: valset.ValActive},
		},
		{
			// addr1: ValPaused+expiredPause → violation(noop) → timeout(ValInactive+IdleTimeout set) → epoch(noop, non-epoch)
			"timeout chain (non-epoch): expired pause → ValInactive",
			testNonEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: aboveMinStake, PausedTimeout: testBlockTime.Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive},
		},
		{
			// addr1: ValPaused+expiredPause → violation(noop) → timeout(ValInactive) → epoch(ValInactive not in competition → stays)
			// addr2: ValActive+aboveMin → epoch(ValActive) — prevents fallback from firing
			"timeout→epoch chain (epoch): expired pause → ValInactive",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: aboveMinStake, PausedTimeout: testBlockTime.Add(-1 * time.Hour)},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive, addr2: valset.ValActive},
		},
		{
			// All validators below minStake → epoch transition produces no ValActive → fallback fires.
			// getFallbackTransition(original) promotes epoch-1 competition group to ValActive.
			// addr1: ValActive+below   → violation(ValExiting) → epoch(ValInactive) → fallback(ValActive)
			// addr2: ValReady+below    → epoch(T3b → ValInactive)                   → fallback(ValActive)
			// addr3: ValPaused+below   → epoch(T3b → ValInactive)                   → fallback(ValActive)
			// addr4: CandTesting+below → epoch(pass vrank, T3b → ValInactive)       → fallback(pass vrank → ValActive)
			"fallback: all below minStake → epoch-1 committee restored",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
				addr2: {State: valset.ValReady, StakingAmount: belowMinStake},
				addr3: {State: valset.ValPaused, StakingAmount: belowMinStake},
				addr4: {State: valset.CandTesting, StakingAmount: belowMinStake},
			},
			map[common.Address]valset.State{
				addr1: valset.ValActive,
				addr2: valset.ValActive,
				addr3: valset.ValActive,
				addr4: valset.ValActive,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestApplyAllTransitions(ctrl)
			parentHeader := &types.Header{Number: big.NewInt(int64(tc.num - 1)), Time: big.NewInt(testBlockTime.Unix())}

			result, err := v.applyAllTransitions(tc.input, parentHeader, DefaultValPausedTimeout, DefaultValIdleTimeout, DefaultActiveValidatorCount, 0, 0, noSlotLimit, noMinActive, 0)
			assert.NoError(t, err)
			for addr, expectedState := range tc.expected {
				assert.Equal(t, expectedState, result[addr].State, "addr=%s", addr.Hex())
			}
		})
	}
}

// TestGetCouncilPermissionless tests GetCouncil filters by council states via GetNodeByState.
func TestGetCouncilPermissionless(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config, _ := system.MakeTestPermissionlessConfig(7)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	simBackend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	defer simBackend.Close()
	chain := simBackend.BlockChain()
	chain.Config().PermissionlessCompatibleBlock = big.NewInt(0)

	v := NewValsetModule()
	v.Chain = chain

	// Seed cache: 7 validators covering all states
	nodes := valset.NodeStateMap{
		config.NodeIds[0]: {State: valset.ValActive},   // included
		config.NodeIds[1]: {State: valset.ValPaused},   // included
		config.NodeIds[2]: {State: valset.ValReady},    // included
		config.NodeIds[3]: {State: valset.ValInactive}, // excluded
		config.NodeIds[4]: {State: valset.ValExiting},  // excluded
		config.NodeIds[5]: {State: valset.CandReady},   // excluded
		config.NodeIds[6]: {State: valset.Registered},  // excluded
	}
	v.nodeStatesCache.Add(uint64(1), nodes)

	council, err := v.GetCouncil(1)
	require.NoError(t, err)

	// Only ValActive, ValPaused, ValReady should be in council (3 of 7)
	councilAddrs := council
	assert.Len(t, councilAddrs, 3)
	assert.Contains(t, councilAddrs, config.NodeIds[0])
	assert.Contains(t, councilAddrs, config.NodeIds[1])
	assert.Contains(t, councilAddrs, config.NodeIds[2])
	assert.NotContains(t, councilAddrs, config.NodeIds[3])
	assert.NotContains(t, councilAddrs, config.NodeIds[4])
	assert.NotContains(t, councilAddrs, config.NodeIds[5])
	assert.NotContains(t, councilAddrs, config.NodeIds[6])
}

// TestGetNodeByState tests GetNodeByState filtering by state.
func TestGetNodeByState(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config, _ := system.MakeTestPermissionlessConfig(5)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	chainConfig := params.TestChainConfig.Copy()
	chainConfig.PermissionlessCompatibleBlock = big.NewInt(0) // permissionless from genesis

	simBackend := backends.NewSimulatedBackendWithChainConfig(blockchain.GenesisAlloc(alloc), chainConfig)
	defer simBackend.Close()
	chain := simBackend.BlockChain()

	v := NewValsetModule()
	v.Chain = chain

	// Seed cache with mixed states
	nodes := valset.NodeStateMap{
		config.NodeIds[0]: {State: valset.ValActive, StakingAmount: 5_000_000},
		config.NodeIds[1]: {State: valset.ValPaused, StakingAmount: 5_000_000},
		config.NodeIds[2]: {State: valset.ValInactive, StakingAmount: 5_000_000},
		config.NodeIds[3]: {State: valset.CandReady, StakingAmount: 5_000_000},
		config.NodeIds[4]: {State: valset.ValExiting, StakingAmount: 5_000_000},
	}
	v.nodeStatesCache.Add(uint64(1), nodes)

	t.Run("filter single state", func(t *testing.T) {
		result, err := v.GetNodeByState(1, []valset.State{valset.ValActive})
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Contains(t, result, config.NodeIds[0])
	})

	t.Run("filter multiple states", func(t *testing.T) {
		result, err := v.GetNodeByState(1, []valset.State{valset.ValActive, valset.ValPaused})
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Contains(t, result, config.NodeIds[0])
		assert.Contains(t, result, config.NodeIds[1])
	})

	t.Run("empty states returns all", func(t *testing.T) {
		result, err := v.GetNodeByState(1, []valset.State{})
		require.NoError(t, err)
		assert.Len(t, result, 5)
	})

	t.Run("nil states returns all", func(t *testing.T) {
		result, err := v.GetNodeByState(1, nil)
		require.NoError(t, err)
		assert.Len(t, result, 5)
	})

	t.Run("no match returns empty", func(t *testing.T) {
		result, err := v.GetNodeByState(1, []valset.State{valset.ValReady})
		require.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("permissionless fork not enabled", func(t *testing.T) {
		noPermConfig := params.TestChainConfig.Copy()
		noPermConfig.PermissionlessCompatibleBlock = nil
		noPermBackend := backends.NewSimulatedBackendWithChainConfig(blockchain.GenesisAlloc(alloc), noPermConfig)
		defer noPermBackend.Close()

		v2 := NewValsetModule()
		v2.Chain = noPermBackend.BlockChain()

		_, err := v2.GetNodeByState(1, nil)
		assert.Error(t, err)
	})
}
