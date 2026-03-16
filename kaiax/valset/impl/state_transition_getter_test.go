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

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
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
	testEpochNum    = uint64(VRankEpoch)     // VRankEpoch % VRankEpoch == 0
	testNonEpochNum = uint64(VRankEpoch - 1) // (VRankEpoch-1) % VRankEpoch != 0
	testMaxValCount = 50

	highStake = uint64(30_000_000)
	midStake  = uint64(20_000_000)
	lowStake  = uint64(10_000_000)
)

var (
	testIdleTimeout    = 24 * time.Hour
	testPauseTimeout   = 8 * time.Hour
	defaultIdleTimeout = 30 * 24 * time.Hour
)

func newTestValsetModule(ctrl *gomock.Controller) *ValsetModule {
	mockGov := gov_mock.NewMockGovModule(ctrl)
	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: new(big.Int).SetUint64(testMinStake),
	}).AnyTimes()
	return &ValsetModule{
		InitOpts: InitOpts{
			GovModule: mockGov,
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
		{"T4b: CandReady + stake<min → CandInactive", valset.CandReady, belowMinStake, valset.CandInactive, false},
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
			result := v.getEpochTransition(testEpochNum, validators, testIdleTimeout, testMaxValCount)
			assert.Equal(t, tc.expectedState, result[addr1].State)
			if tc.expectTimeout {
				assert.False(t, result[addr1].IdleTimeout.IsZero())
			}
		})
	}
}

func TestGetEpochTransition_NonEpoch(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: aboveMinStake},
	}
	result := v.getEpochTransition(testNonEpochNum, validators, testIdleTimeout, testMaxValCount)
	assert.Equal(t, valset.ValActive, result[addr1].State)
}

func TestGetEpochTransition_MaxValidatorCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: highStake},
		addr2: {State: valset.ValActive, StakingAmount: midStake},
		addr3: {State: valset.ValActive, StakingAmount: lowStake},
	}
	result := v.getEpochTransition(testEpochNum, validators, testIdleTimeout, 2)
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
	result := v.getEpochTransition(testEpochNum, validators, testIdleTimeout, 2)

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
	result := v.getEpochTransition(testEpochNum, validators, testIdleTimeout, testMaxValCount)
	assert.Equal(t, valset.ValInactive, result[addr1].State)    // result is transitioned
	assert.Equal(t, valset.ValExiting, validators[addr1].State) // original is unchanged
}

// ============================================================
// TestGetTimeoutTransition
// ============================================================

func TestGetTimeoutTransition(t *testing.T) {
	expiredTimeout := time.Now().Add(-1 * time.Hour)
	futureTimeout := time.Now().Add(10 * 24 * time.Hour)

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
			"ValInactive: idle expired → CandInactive",
			&valset.ValidatorState{State: valset.ValInactive, IdleTimeout: expiredTimeout},
			valset.CandInactive, false, false, nil,
		},
		{
			"ValReady: idle expired → CandInactive",
			&valset.ValidatorState{State: valset.ValReady, IdleTimeout: expiredTimeout},
			valset.CandInactive, false, false, nil,
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
			result := v.getTimeoutTransition(validators, testPauseTimeout, defaultIdleTimeout)
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
		addr3: {State: valset.CandInactive, IdleTimeout: now, PausedTimeout: now},
		addr4: {State: valset.CandTesting, IdleTimeout: now, PausedTimeout: now},
		addr5: {State: valset.ValExiting, IdleTimeout: now, PausedTimeout: now},
	}
	result := v.getTimeoutTransition(validators, testPauseTimeout, defaultIdleTimeout)
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

	validators := valset.NodeStateMap{
		addr1: {State: valset.ValActive, StakingAmount: belowMinStake},
		addr2: {State: valset.ValActive, StakingAmount: aboveMinStake},
	}
	result := v.getViolationTransition(testEpochNum, validators)
	assert.Equal(t, valset.ValExiting, result[addr1].State)
	assert.Equal(t, valset.ValActive, result[addr2].State)
}

func TestGetViolationTransition_NonActiveNotAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := newTestValsetModule(ctrl)

	validators := valset.NodeStateMap{
		addr1: {State: valset.CandReady, StakingAmount: belowMinStake},
		addr2: {State: valset.ValInactive, StakingAmount: belowMinStake},
	}
	result := v.getViolationTransition(testEpochNum, validators)
	assert.Equal(t, valset.CandReady, result[addr1].State)
	assert.Equal(t, valset.ValInactive, result[addr2].State)
}

// ============================================================
// TestGetPermlessCouncil
// ============================================================

func TestGetPermlessCouncil_FiltersCorrectly(t *testing.T) {
	vl := newValidatorList(valset.NodeStateMap{
		addr1: {State: valset.ValActive},
		addr2: {State: valset.ValReady},
		addr3: {State: valset.ValPaused},
		addr4: {State: valset.CandReady},
		addr5: {State: valset.ValInactive},
	})
	council := vl.getPermlessCouncil()
	assert.Equal(t, 3, council.Len())
	assert.True(t, council.Contains(addr1))
	assert.True(t, council.Contains(addr2))
	assert.True(t, council.Contains(addr3))
	assert.False(t, council.Contains(addr4))
	assert.False(t, council.Contains(addr5))
}

func TestGetPermlessCouncil_NilValidatorList(t *testing.T) {
	var vl *ValidatorList
	council := vl.getPermlessCouncil()
	assert.Equal(t, 0, council.Len())
}

func TestGetPermlessCouncil_EmptyValidators(t *testing.T) {
	vl := newValidatorList(valset.NodeStateMap{})
	council := vl.getPermlessCouncil()
	assert.Equal(t, 0, council.Len())
}

// ============================================================
// TestApplyAllTransitions
// ============================================================

// newTestApplyAllTransitions creates a ValsetModule with empty statedb (no ABv2 contract → default fallback).
func newTestApplyAllTransitions(ctrl *gomock.Controller) (*ValsetModule, *state.StateDB) {
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockGov := gov_mock.NewMockGovModule(ctrl)

	chainConfig := &params.ChainConfig{}
	mockChain.EXPECT().Config().Return(chainConfig).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(nil).AnyTimes()

	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: new(big.Int).SetUint64(testMinStake),
	}).AnyTimes()

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)

	v := &ValsetModule{
		InitOpts: InitOpts{
			Chain:     mockChain,
			GovModule: mockGov,
		},
	}
	return v, statedb
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
			"candidate promotion at epoch",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.CandReady, StakingAmount: aboveMinStake},
			},
			map[common.Address]valset.State{addr1: valset.CandTesting},
		},
		{
			// addr1: ValActive+belowMin  → violation(ValExiting) → timeout(noop) → epoch(ValInactive)
			// addr2: ValActive+aboveMin  → violation(noop) → timeout(noop) → epoch(ValActive)
			// addr3: CandReady+aboveMin  → violation(noop) → timeout(noop) → epoch(CandTesting)
			// addr4: ValPaused+aboveMin  → violation(noop) → timeout(set PausedTimeout) → epoch(ValPaused preserved)
			// addr5: CandReady+belowMin  → violation(noop) → timeout(noop) → epoch(CandInactive)
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
				addr5: valset.CandInactive,
			},
		},
		{
			// addr1: ValInactive+expiredIdle → violation(noop) → timeout(CandInactive) → epoch(noop)
			"timeout fires: expired idle → CandInactive",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValInactive, StakingAmount: aboveMinStake, IdleTimeout: time.Now().Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.CandInactive},
		},
		{
			// addr1: ValPaused+expiredPause → violation(noop) → timeout(ValInactive+IdleTimeout set) → epoch(noop, non-epoch)
			"timeout chain (non-epoch): expired pause → ValInactive",
			testNonEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: aboveMinStake, PausedTimeout: time.Now().Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive},
		},
		{
			// addr1: ValPaused+expiredPause → violation(noop) → timeout(ValInactive) → epoch(ValInactive not in competition → stays)
			"timeout→epoch chain (epoch): expired pause → ValInactive",
			testEpochNum,
			valset.NodeStateMap{
				addr1: {State: valset.ValPaused, StakingAmount: aboveMinStake, PausedTimeout: time.Now().Add(-1 * time.Hour)},
			},
			map[common.Address]valset.State{addr1: valset.ValInactive},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v, statedb := newTestApplyAllTransitions(ctrl)
			header := &types.Header{Number: big.NewInt(int64(tc.num))}

			result, err := v.applyAllTransitions(tc.input, tc.num, statedb, header)
			assert.NoError(t, err)
			for addr, expectedState := range tc.expected {
				assert.Equal(t, expectedState, result[addr].State, "addr=%s", addr.Hex())
			}
		})
	}
}

// TestGetCouncilPermissionless tests getCouncilPermissionless filters by council states.
func TestGetCouncilPermissionless(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config := system.MakeTestPermissionlessConfig(t, 7)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	simBackend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	defer simBackend.Close()
	chain := simBackend.BlockChain()

	v := NewValsetModule()
	v.Chain = chain

	// Seed cache: 7 validators covering all states
	nodes := valset.NodeStateMap{
		config.NodeIds[0]: {State: valset.ValActive},    // included
		config.NodeIds[1]: {State: valset.ValPaused},    // included
		config.NodeIds[2]: {State: valset.ValReady},     // included
		config.NodeIds[3]: {State: valset.ValInactive},  // excluded
		config.NodeIds[4]: {State: valset.ValExiting},   // excluded
		config.NodeIds[5]: {State: valset.CandReady},    // excluded
		config.NodeIds[6]: {State: valset.CandInactive}, // excluded
	}
	v.nodeStatesCache.Add(uint64(1), nodes)

	council, err := v.getCouncilPermissionless(1)
	require.NoError(t, err)

	// Only ValActive, ValPaused, ValReady should be in council (3 of 7)
	councilAddrs := council.Council()
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
	config := system.MakeTestPermissionlessConfig(t, 5)

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
