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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/kaiachain/kaia/params"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPermissionlessConfig(permissionlessBlock int64, vrankEpoch uint64) *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	config.PermissionlessCompatibleBlock = big.NewInt(permissionlessBlock)
	config.VRankEpoch = vrankEpoch
	return config
}

func TestGetNodesCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		v := NewValsetModule()
		cached := NodeMap{
			addr1: {State: ValActive, StakingAmount: aboveMinStake},
		}
		v.transitionResultCache.Add(uint64(10), &TransitionResult{Nodes: cached})

		nodes, err := v.getNodes(10)
		require.NoError(t, err)
		require.Same(t, cached[addr1], nodes[addr1])
		assert.Equal(t, cached, nodes)
	})

	t.Run("no cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		root := common.HexToHash("0x1234")
		parentHeader := &types.Header{Number: big.NewInt(9), Root: root}
		stateErr := errors.New("state unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(parentHeader)
		mockChain.EXPECT().StateAt(root).Return(nil, stateErr)

		v := NewValsetModule()
		v.Chain = mockChain

		_, inCache := v.transitionResultCache.Get(uint64(10))
		require.False(t, inCache, "cache must be cold before first call")

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, stateErr)
	})
}

func TestGetTransitionResultCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		v := NewValsetModule()
		cached := &TransitionResult{
			Nodes: NodeMap{
				addr1: {State: ValActive, StakingAmount: aboveMinStake},
			},
		}
		v.transitionResultCache.Add(uint64(10), cached)

		result, err := v.getTransitionResult(10, nil)
		require.NoError(t, err)
		require.Same(t, cached, result)
	})

	t.Run("no cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(nil)

		v := NewValsetModule()
		v.Chain = mockChain

		_, inCache := v.transitionResultCache.Get(uint64(10))
		require.False(t, inCache, "cache must be cold before first call")

		result, err := v.getTransitionResult(10, nil)
		require.Nil(t, result)
		require.ErrorContains(t, err, "parent header not found for block 10")
	})
}

func TestGetNodesParentReadErrors(t *testing.T) {
	t.Run("missing parent header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(nil)

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorContains(t, err, "parent header not found for block 10")
	})

	t.Run("StateAt error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		root := common.HexToHash("0x1234")
		parentHeader := &types.Header{Number: big.NewInt(9), Root: root}
		stateErr := errors.New("state unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(parentHeader)
		mockChain.EXPECT().StateAt(root).Return(nil, stateErr)

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, stateErr)
	})
}

func TestWriteTransitionToABv2PreForkNoOp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

	v := NewValsetModule()
	v.Chain = mockChain

	err := v.WriteTransitionToABv2(nil, &types.Header{Number: big.NewInt(99)}, nil)
	require.NoError(t, err)
}

func TestInstallABv2OutsideForkParentNoOp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

	v := NewValsetModule()
	v.Chain = mockChain

	err := v.InstallABv2(nil, &types.Header{Number: big.NewInt(98)}, nil)
	require.NoError(t, err)
}

func TestFetchVRankCtx(t *testing.T) {
	t.Run("epoch with pf report reads all vrank inputs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cfs := map[common.Address]uint64{addr1: 1}
		pfs := map[common.Address]uint64{addr1: 2}
		pfReport := []common.Address{addr1}

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetCFS(uint64(9)).Return(cfs, nil)
		mockVRank.EXPECT().GetPfReport(uint64(9)).Return(pfReport, nil)
		mockVRank.EXPECT().GetPFS(uint64(9)).Return(pfs, nil)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(9)
		assert.Equal(t, cfs, gotCFS)
		assert.Equal(t, pfs, gotPFS)
		assert.Equal(t, pfReport, gotPfReport)
	})

	t.Run("non-epoch without pf report skips CFS and PFS", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(uint64(8)).Return(nil, nil)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(8)
		assert.Nil(t, gotCFS)
		assert.Nil(t, gotPFS)
		assert.Nil(t, gotPfReport)
	})

	t.Run("errors fail open", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pfReport := []common.Address{addr1}
		vrankErr := errors.New("vrank unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetCFS(uint64(9)).Return(nil, vrankErr)
		mockVRank.EXPECT().GetPfReport(uint64(9)).Return(pfReport, nil)
		mockVRank.EXPECT().GetPFS(uint64(9)).Return(nil, vrankErr)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(9)
		assert.Nil(t, gotCFS)
		assert.Nil(t, gotPFS)
		assert.Equal(t, pfReport, gotPfReport)
	})
}

func TestDiffNodeStates(t *testing.T) {
	oldIdleTimeout := time.Unix(100, 0)
	newIdleTimeout := time.Unix(200, 0)

	parent := NodeMap{
		addr1: {State: ValActive, StakingAmount: highStake},
		addr2: {State: ValPaused, PausedTimeout: oldIdleTimeout},
		addr3: {State: ValReady, IdleTimeout: oldIdleTimeout},
		addr4: {State: ValActive},
	}
	current := NodeMap{
		addr1: {State: ValActive, StakingAmount: lowStake}, // staking amount is not written by this path
		addr2: {State: ValActive},
		addr3: {State: ValReady, IdleTimeout: newIdleTimeout},
		addr4: {State: ValActive, Suspended: true}, // suspended is derived, not written
		addr5: {State: CandReady},
	}

	diff := diffNodeStates(parent, current)
	assert.Equal(t, []common.Address{addr2, addr3, addr5}, diff.Addresses())
	assert.Equal(t, ValActive, diff[addr2].State)
	assert.Equal(t, newIdleTimeout, diff[addr3].IdleTimeout)
	assert.Equal(t, CandReady, diff[addr5].State)

	diff[addr2].State = Registered
	assert.Equal(t, ValActive, current[addr2].State, "diff must not alias current nodes")
}
