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

// TestGetEpochTransition_BelowMinStakeDemoted verifies that VA/VR/VP+belowMin are demoted to VI at epoch (T5).
// This ensures newSF is computed only from validators with sufficient stake,
// preventing SF inflation that would cause incorrect slot limits in violation transition.
func TestGetEpochTransition_BelowMinStakeDemoted(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: aboveMinStake}, // above-min: stays VA
		addr2: {State: valset.ValActive, StakingAmount: belowMinStake}, // T3b: → VI
		addr3: {State: valset.ValReady, StakingAmount: belowMinStake},  // T3b: → VI
		addr4: {State: valset.ValPaused, StakingAmount: belowMinStake}, // T3b: → VI
	}
	result := v.getEpochTransition(testMinStake, validators, testIdleTimeout, testMaxValCount, testBlockTime, 0, 0, 0)

	assert.Equal(t, valset.ValActive, result[addr1].State)
	assert.Equal(t, valset.ValInactive, result[addr2].State)
	assert.Equal(t, valset.ValInactive, result[addr3].State)
	assert.Equal(t, valset.ValInactive, result[addr4].State)
	assert.False(t, result[addr2].IdleTimeout.IsZero())
	assert.False(t, result[addr3].IdleTimeout.IsZero())
	assert.False(t, result[addr4].IdleTimeout.IsZero())
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
	result := v.getViolationTransition(testMinStake, validators, 0, 0, noSlotLimit, noMinActive, testIdleTimeout, testBlockTime)
	assert.Equal(t, valset.ValExiting, result[addr1].State)
	assert.Equal(t, valset.ValActive, result[addr2].State)
}

// TestGetViolationTransition_SuspendedNotExempt verifies that suspended ValActive validators
// are subject to the same violation transitions as non-suspended ones.
func TestGetViolationTransition_SuspendedNotExempt(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
	mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{addr2: 3}, nil)
	v.VRankModule = mockVRank

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: belowMinStake, Suspended: true},  // minStake violation
		addr2: {State: valset.ValActive, StakingAmount: aboveMinStake, Suspended: true},  // PFS severe
		addr3: {State: valset.ValActive, StakingAmount: aboveMinStake, Suspended: false}, // no violation
	}
	result := v.getViolationTransition(testMinStake, validators, 100, 2, noSlotLimit, noMinActive, testIdleTimeout, testBlockTime)
	assert.Equal(t, valset.ValExiting, result[addr1].State, "suspended + low staking → ValExiting")
	assert.Equal(t, valset.ValExiting, result[addr2].State, "suspended + PFS severe → ValExiting")
	assert.Equal(t, valset.ValActive, result[addr3].State, "non-suspended, no violation → unchanged")
}

// TestGetViolationTransition_Deterministic verifies that violation transitions produce
// consistent results regardless of Go map iteration order. Two ValActive validators both
// below minStake compete for a single ValExiting slot — the same one must always win.
func TestGetViolationTransition_Deterministic(t *testing.T) {
	var firstResult map[common.Address]valset.State
	for i := 0; i < 100; i++ {
		ctrl := gomock.NewController(t)
		v := newTestValsetModule(ctrl)
		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
		v.VRankModule = mockVRank

		validators := valset.NodeStateMap{
			addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
			addr2: {State: valset.ValActive, StakingAmount: belowMinStake},
			addr3: {State: valset.ValActive, StakingAmount: aboveMinStake}, // keeps ValActive count > minActiveCount
		}
		// maxSlotAvailable=1, minActiveCount=2 → only 1 of addr1/addr2 can transition
		result := v.getViolationTransition(testMinStake, validators, 0, 0, 1, 2, testIdleTimeout, testBlockTime)

		states := map[common.Address]valset.State{
			addr1: result[addr1].State,
			addr2: result[addr2].State,
			addr3: result[addr3].State,
		}
		if firstResult == nil {
			firstResult = states
		} else {
			require.Equal(t, firstResult, states, "nondeterministic result at iteration %d", i)
		}
		ctrl.Finish()
	}
}

func TestGetViolationTransition_MinStakeMigrated(t *testing.T) {
	testcases := []struct {
		name           string
		validators     valset.NodeStateMap
		maxSlot        uint64
		expectedStates map[common.Address]valset.State
		checkTimeout   *common.Address // if set, assert IdleTimeout is non-zero
	}{
		{
			"ValReady + low staking → ValInactive",
			valset.NodeStateMap{
				addr1: {State: valset.ValReady, StakingAmount: belowMinStake},
				addr2: {State: valset.ValReady, StakingAmount: aboveMinStake},
			},
			noSlotLimit,
			map[common.Address]valset.State{addr1: valset.ValInactive, addr2: valset.ValReady},
			&addr1,
		},
		{
			"ValPaused + low staking → ValExiting",
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: belowMinStake},
				addr2: {State: valset.ValPaused, StakingAmount: aboveMinStake},
			},
			noSlotLimit,
			map[common.Address]valset.State{addr1: valset.ValExiting, addr2: valset.ValPaused},
			nil,
		},
		{
			"ValPaused + low staking + slot full → stays ValPaused",
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: belowMinStake},
				addr2: {State: valset.ValExiting, StakingAmount: aboveMinStake}, // slot occupied
				addr3: {State: valset.ValActive, StakingAmount: aboveMinStake},
				addr4: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			1, // maxSlotAvailable=1, already occupied
			map[common.Address]valset.State{addr1: valset.ValPaused, addr2: valset.ValExiting, addr3: valset.ValActive, addr4: valset.ValActive},
			nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)
			mockVRank := vrank_mock.NewMockVRankModule(ctrl)
			mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
			v.VRankModule = mockVRank

			result := v.getViolationTransition(testMinStake, tc.validators, 0, 0, tc.maxSlot, noMinActive, testIdleTimeout, testBlockTime)
			for addr, expected := range tc.expectedStates {
				assert.Equal(t, expected, result[addr].State, "addr=%s", addr.Hex())
			}
			if tc.checkTimeout != nil {
				assert.False(t, result[*tc.checkTimeout].IdleTimeout.IsZero(), "IdleTimeout should be set")
			}
		})
	}
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
	result := v.getViolationTransition(testMinStake, validators, 0, 0, noSlotLimit, noMinActive, testIdleTimeout, testBlockTime)
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
	result := v.getViolationTransition(testMinStake, validators, 100, 2, noSlotLimit, noMinActive, testIdleTimeout, testBlockTime)
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
	result := v.getViolationTransition(testMinStake, validators, 100, 2, noSlotLimit, noMinActive, testIdleTimeout, testBlockTime)
	assert.Equal(t, valset.CandReady, result[addr1].State, "CandReady not affected by PFS")
	assert.Equal(t, valset.ValPaused, result[addr2].State, "ValPaused not affected by PFS")
}

// ============================================================
// TestGetViolationTransition with SlotLimits
//
// Uses slotFactor=4 (4 validators at epoch start).
// SlotMath (see contracts/libraries/SlotMath.sol):
//   minActiveCount    = ceil(2*4/3)        → 3
//   totalBudget       = 4 - 3             → 1
//   maxSlotAvailable  = ceil(1/2)          → 1  (up to 1 ValPaused, up to 1 ValExiting independently)
//
// With n=4, maxSlot=1 for each but minActive=3 means only 1 total can leave ValActive
// (either 1 paused or 1 exiting, not both), since 4-2=2 < minActive=3.
// ============================================================

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
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive, testIdleTimeout, testBlockTime)

		assert.Equal(t, 1, int(result.CountByState(valset.ValPaused)), "only 1 slot available for ValPaused")
		assert.Equal(t, 3, int(result.CountByState(valset.ValActive)), "remaining 3 should stay ValActive")
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
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive, testIdleTimeout, testBlockTime)

		assert.Equal(t, 1, int(result.CountByState(valset.ValExiting)), "only 1 slot available for ValExiting")
		assert.Equal(t, 3, int(result.CountByState(valset.ValActive)), "remaining 3 should stay ValActive")
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
		result := v.getViolationTransition(testMinStake, validators, 0, 0, slotMax, minActive, testIdleTimeout, testBlockTime)

		assert.Equal(t, 1, int(result.CountByState(valset.ValExiting)), "only 1 can exit due to slot limit")
		assert.Equal(t, 3, int(result.CountByState(valset.ValActive)), "remaining 3 should stay ValActive")
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
		result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold, slotMax, minActive, testIdleTimeout, testBlockTime)
		assert.Equal(t, valset.ValActive, result[addr1].State, "paused slot full, minor violation skipped")
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
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)
			mockVRank := vrank_mock.NewMockVRankModule(ctrl)
			mockVRank.EXPECT().GetPfReport(uint64(100)).Return([]common.Address{addr1}, nil)
			mockVRank.EXPECT().GetPFS(uint64(100)).Return(map[common.Address]uint64{
				addrs[0]: pfsThreshold10,     // severe → ValExiting
				addrs[1]: pfsThreshold10,     // severe → ValExiting
				addrs[2]: pfsThreshold10 - 1, // minor → ValPaused
				addrs[3]: pfsThreshold10 - 1, // minor → blocked by minActive
			}, nil)
			v.VRankModule = mockVRank

			validators := valset.NodeStateMap{}
			for _, addr := range addrs {
				validators[addr] = &valset.ValidatorState{State: valset.ValActive, StakingAmount: aboveMinStake}
			}

			result := v.getViolationTransition(testMinStake, validators, 100, pfsThreshold10, slotMax10, minActive10, testIdleTimeout, testBlockTime)

			assert.Equal(t, 2, int(result.CountByState(valset.ValExiting)), "2 severe → 2 exiting (maxSlot=2)")
			assert.Equal(t, 1, int(result.CountByState(valset.ValPaused)), "1 minor paused, 2nd blocked by minActive")
			assert.Equal(t, 7, int(result.CountByState(valset.ValActive)), "7 remain active (minActive=7)")
		})

		t.Run("3 belowMinStake: 2 exit (exitSlot cap=2), 3rd blocked", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestValsetModule(ctrl)
			mockVRank := vrank_mock.NewMockVRankModule(ctrl)
			mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
			v.VRankModule = mockVRank

			validators := valset.NodeStateMap{}
			for i, addr := range addrs {
				stake := aboveMinStake
				if i < 3 {
					stake = belowMinStake // 3 validators below min stake
				}
				validators[addr] = &valset.ValidatorState{State: valset.ValActive, StakingAmount: stake}
			}

			result := v.getViolationTransition(testMinStake, validators, 0, 0, slotMax10, minActive10, testIdleTimeout, testBlockTime)

			assert.Equal(t, 2, int(result.CountByState(valset.ValExiting)), "maxSlot=2 → 2 can exit")
			assert.Equal(t, 8, int(result.CountByState(valset.ValActive)), "8 remain active")
		})
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
			// addr1: ValActive+belowMin → epoch(T3b: ValInactive) → violation(noop) → timeout(noop)
			// addr2: ValActive+aboveMin → epoch(ValActive) → violation(noop) → timeout(noop)
			"epoch→T5 pipeline",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive, addr2: valset.ValActive},
		},
		{
			// addr1: ValActive+belowMin → violation(ValExiting) → timeout(noop) → epoch(noop, non-epoch)
			// addr2: ValActive+aboveMin → violation(noop) → timeout(noop)
			// Proves: violation applies regardless of epoch, but epoch transition only fires at VRankEpoch
			"non-epoch: violation fires but epoch noop",
			testNonEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
				addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValExiting, addr2: valset.ValActive},
		},
		{
			// addr1: ValActive+belowMin → violation would fire but minActiveCount=1 blocks demotion
			// Proves: last VA is protected from demotion when minActiveCount=1
			"non-epoch: last VA protected by minActiveCount",
			testNonEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
			},
			map[common.Address]valset.State{addr1: valset.ValActive},
		},
		{
			// addr1: CandReady+aboveMin → violation(noop) → timeout(noop) → epoch(CandTesting)
			"candidate promotion at epoch",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.CandReady, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.CandTesting},
		},
		{
			// addr1: ValActive+belowMin  → epoch(T3b: ValInactive) → violation(noop) → timeout(noop)
			// addr2: ValActive+aboveMin  → epoch(ValActive) → violation(noop) → timeout(noop)
			// addr3: CandReady+aboveMin  → epoch(CandTesting) → violation(noop) → timeout(noop)
			// addr4: ValPaused+aboveMin  → epoch(ValPaused preserved) → violation(noop) → timeout(set PausedTimeout)
			// addr5: CandReady+belowMin  → epoch(Registered) → violation(noop) → timeout(noop)
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
			"timeout fires: expired idle → Registered",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValInactive, StakingAmount: aboveMinStake, IdleTimeout: testBlockTime.Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.Registered},
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
			"timeout→epoch chain (epoch): expired pause → ValInactive",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: aboveMinStake, PausedTimeout: testBlockTime.Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := newTestApplyAllTransitions(ctrl)
			parentHeader := &types.Header{Number: big.NewInt(int64(tc.num - 1)), Time: big.NewInt(testBlockTime.Unix())}

			noopSlotLimitsFn := func(sf uint64) (uint64, uint64, error) { return noSlotLimit, 1, nil }
			res := &system.NodeStatesResult{
				Validators:           tc.input,
				PauseTimeout:         DefaultValPausedTimeout,
				IdleTimeout:          DefaultValIdleTimeout,
				ActiveValidatorCount: DefaultActiveValidatorCount,
				MaxSlotAvailable:     noSlotLimit,
				MinActiveCount:       1,
			}
			result, _, err := v.applyAllTransitions(res, parentHeader, noopSlotLimitsFn)
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
	v.nodeStatesCache.Add(uint64(1), nodeStatesCacheEntry{validators: nodes})

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
	v.nodeStatesCache.Add(uint64(1), nodeStatesCacheEntry{validators: nodes})

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

// TestSuspendedFallback tests the safety fallback in getQualifiedValidators:
// when all ValActive are suspended, the suspended set is ignored to prevent consensus halt.
func TestSuspendedFallback(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config, _ := system.MakeTestPermissionlessConfig(5)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	chainConfig := params.TestChainConfig.Copy()
	chainConfig.PermissionlessCompatibleBlock = big.NewInt(0)

	simBackend := backends.NewSimulatedBackendWithChainConfig(blockchain.GenesisAlloc(alloc), chainConfig)
	defer simBackend.Close()
	chain := simBackend.BlockChain()

	v := NewValsetModule()
	v.Chain = chain

	t.Run("all ValActive suspended → fallback returns all ValActive", func(t *testing.T) {
		nodes := valset.NodeStateMap{
			config.NodeIds[0]: {State: valset.ValActive, Suspended: true},
			config.NodeIds[1]: {State: valset.ValActive, Suspended: true},
			config.NodeIds[2]: {State: valset.ValReady},
		}
		v.nodeStatesCache.Add(uint64(1), nodeStatesCacheEntry{validators: nodes})

		qualified, err := v.getQualifiedValidators(1)
		require.NoError(t, err)
		assert.Equal(t, 2, qualified.Len(), "fallback: all suspended ValActive returned")
		assert.True(t, qualified.Contains(config.NodeIds[0]))
		assert.True(t, qualified.Contains(config.NodeIds[1]))
	})

	t.Run("some ValActive suspended → only non-suspended returned", func(t *testing.T) {
		nodes := valset.NodeStateMap{
			config.NodeIds[0]: {State: valset.ValActive, Suspended: true},
			config.NodeIds[1]: {State: valset.ValActive, Suspended: false},
			config.NodeIds[2]: {State: valset.ValReady},
		}
		v.nodeStatesCache.Add(uint64(1), nodeStatesCacheEntry{validators: nodes})

		qualified, err := v.getQualifiedValidators(1)
		require.NoError(t, err)
		assert.Equal(t, 1, qualified.Len())
		assert.True(t, qualified.Contains(config.NodeIds[1]))
	})

	t.Run("demoted = council - qualified with suspended", func(t *testing.T) {
		nodes := valset.NodeStateMap{
			config.NodeIds[0]: {State: valset.ValActive, Suspended: true},
			config.NodeIds[1]: {State: valset.ValActive, Suspended: false},
			config.NodeIds[2]: {State: valset.ValReady},
			config.NodeIds[3]: {State: valset.ValPaused},
		}
		v.nodeStatesCache.Add(uint64(1), nodeStatesCacheEntry{validators: nodes})

		demoted, err := v.GetDemotedValidators(1)
		require.NoError(t, err)
		// demoted = council(4) - qualified(1 non-suspended ValActive) = 3
		assert.Len(t, demoted, 3)
		assert.Contains(t, demoted, config.NodeIds[0], "suspended ValActive is demoted")
		assert.Contains(t, demoted, config.NodeIds[2], "ValReady is demoted")
		assert.Contains(t, demoted, config.NodeIds[3], "ValPaused is demoted")
		assert.NotContains(t, demoted, config.NodeIds[1], "non-suspended ValActive is not demoted")
	})
}
