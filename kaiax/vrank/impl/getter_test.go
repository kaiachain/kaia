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
	"github.com/kaiachain/kaia/consensus/istanbul"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testChain struct {
	headers map[uint64]*types.Header
}

func (c *testChain) CurrentHeader() *types.Header {
	var result *types.Header
	for num, h := range c.headers {
		if result == nil || num > result.Number.Uint64() {
			result = h
		}
	}
	return result
}

func (c *testChain) GetHeaderByNumber(number uint64) *types.Header {
	return c.headers[number]
}

func makeHeaderWithRound(number uint64, round int64) *types.Header {
	h := &types.Header{
		Number: big.NewInt(int64(number)),
		Extra:  make([]byte, types.IstanbulExtraVanity),
	}
	return types.SetRoundToHeader(h, round)
}

// ---------------------------------------------------------------------------
// TestGetCfReport
// ---------------------------------------------------------------------------

func TestGetCfReport(t *testing.T) {
	t.Run("returns decoded report from header.VRank", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		c1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
		c2 := common.HexToAddress("0x0000000000000000000000000000000000000002")
		encoded, err := vrank.EncodeReport(vrank.Report{c1, c2})
		require.NoError(t, err)
		h := makeHeaderWithRound(10, 0)
		h.VRank = encoded
		v.Chain = &testChain{headers: map[uint64]*types.Header{10: h}}

		report, err := v.GetCfReport(10)
		require.NoError(t, err)
		assert.Equal(t, vrank.Report{c1, c2}, report)

		report2, err := v.GetCfReport(10)
		require.NoError(t, err)
		assert.Equal(t, report, report2, "GetCfReport must be deterministic")
	})

	t.Run("nil header.VRank returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{headers: map[uint64]*types.Header{
			5: makeHeaderWithRound(5, 0), // header.VRank is nil
		}}

		report, err := v.GetCfReport(5)
		require.NoError(t, err)
		assert.Empty(t, report)
	})
}

func TestGetCfReport_Errors(t *testing.T) {
	t.Run("header not found returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{headers: map[uint64]*types.Header{}}

		report, err := v.GetCfReport(99)
		assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
		assert.Nil(t, report)
	})
}

// ---------------------------------------------------------------------------
// TestGetPfReport
// ---------------------------------------------------------------------------

func TestGetPfReport(t *testing.T) {
	t.Run("round zero returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				10: makeHeaderWithRound(10, 0),
			},
		}

		report, err := v.GetPfReport(10)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("round greater than zero returns proposers in round order", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				20: makeHeaderWithRound(20, 3),
			},
		}
		p0 := common.HexToAddress("0x0000000000000000000000000000000000000001")
		p1 := common.HexToAddress("0x0000000000000000000000000000000000000002")
		p2 := common.HexToAddress("0x0000000000000000000000000000000000000003")
		valset.EXPECT().GetProposer(uint64(20), uint64(0)).Return(p0, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(20), uint64(1)).Return(p1, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(20), uint64(2)).Return(p2, nil).Times(2)

		report, err := v.GetPfReport(20)
		require.NoError(t, err)
		assert.Equal(t, vrank.Report{p0, p1, p2}, report)

		report2, err := v.GetPfReport(20)
		require.NoError(t, err)
		assert.Equal(t, report, report2, "GetPfReport must be deterministic")
	})
}

func TestGetPfReport_Errors(t *testing.T) {
	t.Run("header not found returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{headers: map[uint64]*types.Header{}}

		report, err := v.GetPfReport(30)
		require.ErrorIs(t, err, vrank.ErrHeaderNotFound)
		assert.Nil(t, report)
	})

	t.Run("short extra returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				40: {Number: big.NewInt(40), Extra: make([]byte, types.IstanbulExtraVanity-1)},
			},
		}

		report, err := v.GetPfReport(40)
		require.ErrorIs(t, err, vrank.ErrHeaderExtraTooShort)
		assert.Nil(t, report)
	})

	t.Run("GetProposer error is propagated", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				50: makeHeaderWithRound(50, 2),
			},
		}
		p0 := common.HexToAddress("0x0000000000000000000000000000000000000011")
		valset.EXPECT().GetProposer(uint64(50), uint64(0)).Return(p0, nil).Times(1)
		valset.EXPECT().GetProposer(uint64(50), uint64(1)).Return(common.Address{}, assert.AnError).Times(1)

		report, err := v.GetPfReport(50)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, report)
	})
}

// ---------------------------------------------------------------------------
// TestTallyCfReport
// ---------------------------------------------------------------------------

func TestTallyCfReport(t *testing.T) {
	var (
		valset                 = mock_valset.NewMockValsetModule(gomock.NewController(t))
		block1                 = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		block2                 = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		view1_0                = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
		view2_0                = &istanbul.View{Sequence: big.NewInt(2), Round: common.Big0}
		validators, candidates []*CN
		valAddrs               = make([]common.Address, 3)
		candAddrs              = make([]common.Address, 8)
		candMsgsBlock2         = make([]vrank.VRankCandidate, 8)
		earlybirdCands         = candAddrs[0:2] // sent VRankCandidate for block2 before validator preprepared block2
		ontimeCands            = candAddrs[2:4]
		liarCands              = candAddrs[4:6]
		lateCands              = candAddrs[6:8]
	)

	for i := 0; i < 3; i++ {
		validators = append(validators, createCN(t, valset))
		valAddrs[i] = validators[i].Addr
	}
	for i := 0; i < 8; i++ {
		candidates = append(candidates, createCN(t, valset))
		candAddrs[i] = candidates[i].Addr
	}

	valset.EXPECT().GetCouncil(gomock.Any()).Return(valAddrs, nil).AnyTimes()
	valset.EXPECT().GetCandidates(gomock.Any()).Return(candAddrs, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(validators[0].Addr, nil).AnyTimes()

	for i := 0; i < 8; i++ {
		sig := signVRankCandidate(t, candidates[i].VRankModule, candidates[i].Key, block2.NumberU64(), uint8(view2_0.Round.Uint64()), block2.Hash())
		candMsgsBlock2[i] = vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: block2.Hash(), Sig: sig}
	}

	// Initialize prepreparedView at least once to set `v.prepreparedView`.
	for _, v := range validators {
		v.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
	}

	// Earlybirds: candidates send VRankCandidate for block2 before validator has preprepared block2.
	for i := 0; i < 2; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	// Now validator preprepares for block2.
	for _, v := range validators {
		v.VRankModule.HandleIstanbulPreprepare(block2, view2_0)
	}

	// On-time: candidates send VRankCandidate for block2 after validator has preprepared block2 and before deadline.
	for i := 2; i < 4; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	// Liars: candidates send VRankCandidate for block2 with wrong BlockHash.
	for i := 4; i < 6; i++ {
		liarHash := common.Hash{byte(i)}
		sig := signVRankCandidate(t, candidates[i].VRankModule, candidates[i].Key, block2.NumberU64(), uint8(view2_0.Round.Uint64()), liarHash)
		liarMsg := vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: liarHash, Sig: sig}
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&liarMsg)
			assert.NoError(t, err)
		}
	}
	time.Sleep(candidateMsgTimeoutMs * time.Millisecond)
	// Late: candidates send VRankCandidate for block2 after deadline.
	for i := 6; i < 8; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	for _, v := range validators {
		report, err := v.VRankModule.TallyCfReport(2, 0)
		assert.NoError(t, err)
		assert.Len(t, report, 4, "cfReport: 2 liars + 2 late")
		for _, addr := range ontimeCands {
			assert.False(t, slices.Contains(report, addr))
		}
		for _, addr := range earlybirdCands {
			assert.False(t, slices.Contains(report, addr))
		}
		for _, addr := range liarCands {
			assert.True(t, slices.Contains(report, addr))
		}
		for _, addr := range lateCands {
			assert.True(t, slices.Contains(report, addr))
		}
		report2, err := v.VRankModule.TallyCfReport(2, 0)
		assert.NoError(t, err)
		assert.Equal(t, report, report2, "TallyCfReport must be deterministic")
	}
}

func TestTallyCfReport_Errors(t *testing.T) {
	block1 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	view1_0 := &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	candAddr := common.HexToAddress("0xc4nd1d473")

	t.Run("Report should contain candAddr", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.TallyCfReport(1, 0)
		require.NoError(t, err)
		assert.True(t, slices.Contains(report, candAddr))
	})

	t.Run("epoch header returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(vrankEpoch-1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(vrankEpoch-1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(vrankEpoch-1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(vrankEpoch - 1))})
		view := &istanbul.View{Sequence: big.NewInt(vrankEpoch - 1), Round: common.Big0}
		val.VRankModule.HandleIstanbulPreprepare(block, view)

		report, err := val.VRankModule.TallyCfReport(vrankEpoch-1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("round out of range returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.TallyCfReport(1, 11) // maxRound is 10
		require.ErrorIs(t, err, vrank.ErrRoundOutOfRange)
		assert.Nil(t, report)

		report, err = val.VRankModule.TallyCfReport(1, 10)
		assert.NotErrorIs(t, err, vrank.ErrRoundOutOfRange)
	})

	t.Run("non-validator returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val, otherVal := createCN(t, valset), createCN(t, valset)
		// This node is not in the council for block 1.
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{otherVal.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.TallyCfReport(1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("preprepared time not set returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		// skip HandleIstanbulPreprepare

		prepreparedTime, _, _ := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		report, err := val.VRankModule.TallyCfReport(1, 0)
		require.ErrorIs(t, err, vrank.ErrPrepreparedTimeNotSet)
		assert.Nil(t, report)
	})

	t.Run("GetCandidates failed returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return(nil, assert.AnError).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.TallyCfReport(1, 0)
		require.ErrorIs(t, err, vrank.ErrGetCandidateFailed)
		assert.Nil(t, report)
	})
}
