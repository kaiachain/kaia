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
	"slices"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	mock_randao "github.com/kaiachain/kaia/kaiax/randao/mock"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// TestGetCfReport
// ---------------------------------------------------------------------------
func TestGetCfReport(t *testing.T) {
	t.Run("returns decoded report from header.VRank", func(t *testing.T) {
		c1, c2 := numToAddr(1), numToAddr(2)
		encoded, err := vrank.EncodeVRank(vrank.VRankPayload{Report: []common.Address{c1, c2}})
		require.NoError(t, err)
		h := makeHeaderWithRound(10, 0)
		h.VRank = encoded
		v := newCN(t, withHeaders(map[uint64]*types.Header{10: h})).VRankModule

		report, err := v.cfReport(10)
		require.NoError(t, err)
		assert.Equal(t, []common.Address{c1, c2}, report)

		// must be deterministic
		report2, err := v.cfReport(10)
		require.NoError(t, err)
		assert.Equal(t, report, report2, "cfReport must be deterministic")
	})

	t.Run("nil header.VRank returns empty report", func(t *testing.T) {
		v := newCN(t, withHeaders(map[uint64]*types.Header{
			5: makeHeaderWithRound(5, 0), // header.VRank is nil
		})).VRankModule

		report, err := v.cfReport(5)
		require.NoError(t, err)
		assert.Empty(t, report)
	})
}

func TestGetCfReport_Errors(t *testing.T) {
	t.Run("header not found returns error", func(t *testing.T) {
		v := newCN(t).VRankModule

		report, err := v.cfReport(99)
		assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
		assert.Nil(t, report)
	})

	t.Run("pre-fork block returns ErrNotPermissionless", func(t *testing.T) {
		v := newCN(t,
			withHardfork("osaka"),
			withHeaders(map[uint64]*types.Header{10: makeHeaderWithRound(10, 0)}),
			withoutStart(),
		).VRankModule

		report, err := v.cfReport(10)
		assert.ErrorIs(t, err, vrank.ErrNotPermissionless)
		assert.Nil(t, report)
	})
}

// ---------------------------------------------------------------------------
// TestGetPfReport
// ---------------------------------------------------------------------------
func TestGetPfReport(t *testing.T) {
	t.Run("parent round zero returns empty report", func(t *testing.T) {
		v := newCN(t, withHeaders(map[uint64]*types.Header{
			10: makeHeaderWithRound(10, 0),
		})).VRankModule

		report, err := v.pfReport(10)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("parent round greater than zero returns proposers in round order", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		p0, p1, p2 := numToAddr(1), numToAddr(2), numToAddr(3)
		valset.EXPECT().GetProposer(uint64(19), uint64(0)).Return(p0, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(19), uint64(1)).Return(p1, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(19), uint64(2)).Return(p2, nil).Times(2)

		v := newCN(t, withValset(valset), withHeaders(map[uint64]*types.Header{
			20: makeHeaderWithParentRound(20, 3),
		})).VRankModule

		report, err := v.pfReport(20)
		require.NoError(t, err)
		assert.Equal(t, []common.Address{p0, p1, p2}, report)

		// must be deterministic
		report2, err := v.pfReport(20)
		require.NoError(t, err)
		assert.Equal(t, report, report2, "pfReport must be deterministic")
	})
}

func TestGetPfReport_Errors(t *testing.T) {
	t.Run("pre-fork block returns ErrNotPermissionless", func(t *testing.T) {
		v := newCN(t,
			withHardfork("osaka"),
			withHeaders(map[uint64]*types.Header{10: makeHeaderWithRound(10, 0)}),
			withoutStart(),
		).VRankModule

		report, err := v.pfReport(10)
		assert.ErrorIs(t, err, vrank.ErrNotPermissionless)
		assert.Nil(t, report)
	})

	t.Run("header not found returns error", func(t *testing.T) {
		v := newCN(t).VRankModule // default: empty chain

		report, err := v.pfReport(30)
		require.ErrorIs(t, err, vrank.ErrHeaderNotFound)
		assert.Nil(t, report)
	})

	t.Run("undecodable vrank returns error", func(t *testing.T) {
		header := makeHeaderWithRound(40, 0)
		header.VRank = []byte{0x01, 0x02, 0x03}
		v := newCN(t, withHeaders(map[uint64]*types.Header{40: header})).VRankModule

		report, err := v.pfReport(40)
		require.Error(t, err)
		assert.Nil(t, report)
	})

	t.Run("GetProposer error is propagated", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		p0 := numToAddr(11)
		valset.EXPECT().GetProposer(uint64(49), uint64(0)).Return(p0, nil).Times(1)
		valset.EXPECT().GetProposer(uint64(49), uint64(1)).Return(common.Address{}, assert.AnError).Times(1)

		v := newCN(t, withValset(valset), withHeaders(map[uint64]*types.Header{
			50: makeHeaderWithParentRound(50, 2),
		})).VRankModule

		report, err := v.pfReport(50)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, report)
	})
}

// ---------------------------------------------------------------------------
// TestEvaluateCandidates
// ---------------------------------------------------------------------------
func TestEvaluateCandidates(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		valset = mock_valset.NewMockValsetModule(ctrl)
		randao = mock_randao.NewMockRandaoModule(ctrl)

		block2  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		view2_0 = &bft.View{Sequence: big.NewInt(2), Round: common.Big0}

		candidates     []*CN
		candAddrs      = make([]common.Address, 8)
		candMsgsBlock2 = make([]vrank.VRankCandidate, 8)
		ontimeCands    = candAddrs[0:4] // sent VRankCandidate for block2 after the proposer preprepared it, before deadline
		liarCands      = candAddrs[4:6]
		lateCands      = candAddrs[6:8]
	)

	// Only the proposer collects candidate replies; it records the view before any reply can arrive.
	proposer := newCN(t, withValset(valset), withRandao(randao), withGenesis())
	for i := 0; i < 8; i++ {
		candidates = append(candidates, newCN(t, withValset(valset), withRandao(randao), withGenesis()))
		candAddrs[i] = candidates[i].Addr
	}

	valset.EXPECT().GetCandTesting(gomock.Any()).Return(candAddrs, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(proposer.Addr, nil).AnyTimes()

	for i := 0; i < 8; i++ {
		sig, blsSig := signVRankCandidate(t, candidates[i].VRankModule, candidates[i].Key, candidates[i].BlsKey, block2.NumberU64(), uint8(view2_0.Round.Uint64()), block2.Hash())
		candMsgsBlock2[i] = vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: block2.Hash(), Sig: sig, BlsSig: blsSig}
	}

	proposer.VRankModule.HandleIstanbulPreprepare(block2, view2_0)

	for i, cand := range candidates {
		candMsg := &candMsgsBlock2[i]
		if i == 6 {
			time.Sleep(candidateMsgTimeoutMs * time.Millisecond)
		}
		switch {
		case i < 4:
			// On-time: candidates send VRankCandidate for block2 after the proposer preprepared it and before deadline.
			err := proposer.VRankModule.HandleVRankCandidate(candMsg)
			assert.NoError(t, err)
		case i < 6:
			// Liars: candidates send VRankCandidate for block2 with wrong BlockHash.
			liarHash := common.Hash{byte(i)}
			sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block2.NumberU64(), uint8(view2_0.Round.Uint64()), liarHash)
			liarMsg := vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: liarHash, Sig: sig, BlsSig: blsSig}
			err := proposer.VRankModule.HandleVRankCandidate(&liarMsg)
			assert.NoError(t, err)
		default:
			// Late: candidates send VRankCandidate for block2 after deadline.
			err := proposer.VRankModule.HandleVRankCandidate(candMsg)
			assert.NoError(t, err)
		}
	}

	report, err := proposer.VRankModule.EvaluateCandidates(2, 0)
	assert.NoError(t, err)
	assert.Len(t, report, 4, "cfReport: 2 liars + 2 late")
	for _, addr := range ontimeCands {
		assert.False(t, slices.Contains(report, addr))
	}
	for _, addr := range liarCands {
		assert.True(t, slices.Contains(report, addr))
	}
	for _, addr := range lateCands {
		assert.True(t, slices.Contains(report, addr))
	}
	report2, err := proposer.VRankModule.EvaluateCandidates(2, 0)
	assert.NoError(t, err)
	assert.Equal(t, report, report2, "EvaluateCandidates must be deterministic")
}

func TestEvaluateCandidates_Errors(t *testing.T) {
	block1 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	view1_0 := &bft.View{Sequence: big.NewInt(1), Round: common.Big0}
	candAddr := common.HexToAddress("0xc4nd1d473")

	newCNWithDefaults := func() *CN {
		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		cn := newCN(t, withValset(valset), withGenesis())
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{cn.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(cn.Addr, nil).AnyTimes()
		return cn
	}

	t.Run("pre-fork block returns ErrNotPermissionless", func(t *testing.T) {
		// EvaluateCandidates(blockNum=0, ...) targets header(1).VRank; fork check is on blockNum+1=1.
		// With osaka config, fork is never enabled, so block 1 is pre-fork.
		v := newCN(t, withHardfork("osaka"), withoutStart()).VRankModule

		report, err := v.EvaluateCandidates(0, 0)
		assert.ErrorIs(t, err, vrank.ErrNotPermissionless)
		assert.Nil(t, report)
	})

	t.Run("Report should contain candAddr", func(t *testing.T) {
		cn := newCNWithDefaults()
		cn.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := cn.VRankModule.EvaluateCandidates(1, 0)
		require.NoError(t, err)
		assert.True(t, slices.Contains(report, candAddr))
	})

	t.Run("round above maxRound returns empty report", func(t *testing.T) {
		cn := newCNWithDefaults()
		cn.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		// round > maxRound must degrade gracefully.
		report, err := cn.VRankModule.EvaluateCandidates(1, 11) // maxRound is 10
		require.NoError(t, err)
		assert.NotErrorIs(t, err, vrank.ErrRoundOutOfRange)
		assert.Empty(t, report)

		report, err = cn.VRankModule.EvaluateCandidates(1, 10)
		require.NoError(t, err)
		assert.NotErrorIs(t, err, vrank.ErrRoundOutOfRange)
		assert.Empty(t, report)
	})

	t.Run("non-validator returns empty report", func(t *testing.T) {
		// Two CNs share the same valset mocks.
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		cn1, cn2 := newCN(t, withValset(valset), withGenesis()), newCN(t, withValset(valset), withGenesis())

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{cn2.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(cn1.Addr, nil).AnyTimes()
		cn1.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := cn2.VRankModule.EvaluateCandidates(1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("no preprepare data returns empty report", func(t *testing.T) {
		val := newCN(t, withGenesis())
		// skip HandleIstanbulPreprepare — no preprepare data in collector

		prepreparedTime, _, _ := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		report, err := val.VRankModule.EvaluateCandidates(1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("GetCandTesting failed returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		cn := newCN(t, withValset(valset), withGenesis())
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{cn.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(cn.Addr, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return(nil, assert.AnError).AnyTimes()
		cn.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := cn.VRankModule.EvaluateCandidates(1, 0)
		require.ErrorIs(t, err, vrank.ErrGetCandidateFailed)
		assert.Nil(t, report)
	})
}
