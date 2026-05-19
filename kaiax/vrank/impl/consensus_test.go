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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeVRankHeader builds a header at the given block number containing the encoded cfReport.
// cfAddrs may be nil for an empty VRank field.
func makeVRankHeader(t *testing.T, number uint64, cfAddrs []common.Address) *types.Header {
	t.Helper()
	h := &types.Header{Number: big.NewInt(int64(number))}
	if len(cfAddrs) > 0 {
		encoded, err := vrank.EncodeReport(cfAddrs)
		require.NoError(t, err)
		h.VRank = encoded
	}
	return h
}

func makeEpochStartVRankHeader(t *testing.T, number uint64, candidates []common.Address) *types.Header {
	t.Helper()
	encoded, err := vrank.EncodeAddressList(candidates)
	require.NoError(t, err)
	return &types.Header{Number: big.NewInt(int64(number)), VRank: encoded}
}

// TestVerifyHeader covers all branches of VRankModule.VerifyHeader.
func TestVerifyHeader(t *testing.T) {
	C1, C2, C3 := numToAddr(1), numToAddr(2), numToAddr(3)
	candidates := []common.Address{C1, C2, C3}
	const num = uint64(100)

	t.Run("pre-fork: VRank must be empty", func(t *testing.T) {
		v := newCN(t, withHardfork("osaka"), withoutStart()).VRankModule
		h := &types.Header{Number: big.NewInt(100)}
		assert.NoError(t, v.VerifyHeader(h, nil))
		h = makeVRankHeader(t, 100, []common.Address{C1})
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrUnexpectedVRankBeforePermissionless)
	})

	t.Run("epoch-start: VRank must match CandTesting", func(t *testing.T) {
		v := newCN(t, withCandidates([]common.Address{C1, C2})).VRankModule
		assert.NoError(t, v.VerifyHeader(makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C1, C2}), nil))
		assert.ErrorIs(t, v.VerifyHeader(makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C2, C1}), nil),
			vrank.ErrEpochStartVRankMismatch)
		assert.ErrorIs(t, v.VerifyHeader(makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C1, C1}), nil),
			vrank.ErrEpochStartVRankMismatch)
		assert.ErrorIs(t, v.VerifyHeader(makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C1, C3}), nil),
			vrank.ErrEpochStartVRankMismatch)
		assert.ErrorIs(t, v.VerifyHeader(makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C1}), nil),
			vrank.ErrEpochStartVRankMismatch)
		assert.ErrorIs(t, v.VerifyHeader(&types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}, nil),
			vrank.ErrEpochStartVRankMismatch)
	})

	t.Run("epoch-start: empty CandTesting is encoded as an RLP list", func(t *testing.T) {
		v := newCN(t, withCandidates(nil)).VRankModule
		h := makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, nil)
		assert.Equal(t, []byte{0xc0}, h.VRank)
		assert.NoError(t, v.VerifyHeader(h, nil))
	})

	t.Run("epoch-start: candidate membership is checked against block N", func(t *testing.T) {
		headerNum := uint64(params.DefaultVRankEpoch)
		cn := newCN(t, withoutStart())
		cn.Valset.EXPECT().GetCandTesting(headerNum).Return(candidates, nil).Times(1)

		assert.NoError(t, cn.VRankModule.VerifyHeader(makeEpochStartVRankHeader(t, headerNum, candidates), nil))
	})

	t.Run("epoch-start: candidate lookup failure is returned", func(t *testing.T) {
		headerNum := uint64(params.DefaultVRankEpoch)
		cn := newCN(t, withoutStart())
		cn.Valset.EXPECT().GetCandTesting(headerNum).Return(nil, assert.AnError).Times(1)

		assert.ErrorIs(t, cn.VRankModule.VerifyHeader(makeEpochStartVRankHeader(t, headerNum, candidates), nil), assert.AnError)
	})

	t.Run("post-fork non-epoch: empty VRank passes without reading candidates", func(t *testing.T) {
		v := newCN(t, withoutStart()).VRankModule
		h := &types.Header{Number: big.NewInt(100)}
		assert.NoError(t, v.VerifyHeader(h, nil))
	})

	t.Run("valid sorted report with known candidates passes", func(t *testing.T) {
		v := newCN(t, withCandidates(candidates)).VRankModule
		assert.NoError(t, v.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, C2, C3}), nil))
	})

	t.Run("non-candidate in VRank rejected", func(t *testing.T) {
		unknown := numToAddr(99)
		v := newCN(t, withCandidates(candidates)).VRankModule
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, unknown}), nil),
			vrank.ErrInvalidVRankCandidate)
	})

	t.Run("duplicate address rejected", func(t *testing.T) {
		v := newCN(t, withCandidates(candidates)).VRankModule
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, C1}), nil),
			vrank.ErrDuplicateVRankCandidate)
	})

	t.Run("unsorted addresses rejected", func(t *testing.T) {
		v := newCN(t, withCandidates(candidates)).VRankModule
		// C3 > C2, so C3 before C2 is not ascending.
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, num, []common.Address{C3, C2}), nil),
			vrank.ErrVRankNotSorted)
	})

	t.Run("invalid encoding rejected before reading candidates", func(t *testing.T) {
		v := newCN(t, withoutStart()).VRankModule
		h := &types.Header{Number: big.NewInt(100), VRank: []byte{0xff, 0xfe}} // garbage
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrInvalidVRankFormat)
	})

	t.Run("non-epoch candidate lookup failure is returned", func(t *testing.T) {
		headerNum := uint64(101)
		cn := newCN(t, withoutStart())
		cn.Valset.EXPECT().GetCandTesting(headerNum-1).Return(nil, assert.AnError).Times(1)

		assert.ErrorIs(t, cn.VRankModule.VerifyHeader(makeVRankHeader(t, headerNum, []common.Address{C1}), nil), assert.AnError)
	})

	t.Run("candidate membership is checked against the reported block N-1", func(t *testing.T) {
		headerNum := uint64(101)
		cand, futureCand := numToAddr(10), numToAddr(11)
		cn := newCN(t, withoutStart())
		v := cn.VRankModule

		cn.Valset.EXPECT().GetCandTesting(headerNum-1).Return([]common.Address{cand}, nil).Times(2)

		assert.NoError(t, v.VerifyHeader(makeVRankHeader(t, headerNum, []common.Address{cand}), nil))
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, headerNum, []common.Address{futureCand}), nil), vrank.ErrInvalidVRankCandidate)
	})
}

func TestPrepareHeader(t *testing.T) {
	C1, C2 := numToAddr(1), numToAddr(2)
	candidates := []common.Address{C1, C2}

	t.Run("pre-fork clears VRank", func(t *testing.T) {
		v := newCN(t, withHardfork("osaka"), withoutStart()).VRankModule
		header := makeVRankHeader(t, 100, []common.Address{C1})

		require.NoError(t, v.PrepareHeader(header))
		assert.Nil(t, header.VRank)
	})

	t.Run("epoch-start fills CandTesting", func(t *testing.T) {
		v := newCN(t, withCandidates(candidates), withoutStart()).VRankModule
		header := &types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}

		require.NoError(t, v.PrepareHeader(header))
		report, err := vrank.DecodeReport(header.VRank)
		require.NoError(t, err)
		assert.Equal(t, candidates, report)
		assert.NoError(t, v.VerifyHeader(header, nil))
	})

	t.Run("epoch-start with no candidates writes encoded empty list", func(t *testing.T) {
		v := newCN(t, withCandidates(nil), withoutStart()).VRankModule
		header := &types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}

		require.NoError(t, v.PrepareHeader(header))
		assert.Equal(t, []byte{0xc0}, header.VRank)
		assert.NoError(t, v.VerifyHeader(header, nil))
	})

	t.Run("non-epoch empty evaluation leaves VRank nil", func(t *testing.T) {
		parent := makeHeaderWithRound(10, 1)
		v := newCN(t, withCandidates(candidates), withHeaders(map[uint64]*types.Header{10: parent}), withoutStart()).VRankModule
		header := &types.Header{Number: big.NewInt(11)}

		require.NoError(t, v.PrepareHeader(header))
		assert.Nil(t, header.VRank)
	})

	t.Run("non-epoch fills evaluated candidate failures", func(t *testing.T) {
		parent := makeHeaderWithRound(10, 1)
		parentBlock := types.NewBlockWithHeader(parent)
		v := newCN(t, withCandidates(candidates), withHeaders(map[uint64]*types.Header{10: parent}), withoutStart()).VRankModule
		v.collector.AddPrepreparedTime(vrank.ViewKey{N: 10, R: 1}, time.Now(), parentBlock.Hash())
		header := &types.Header{Number: big.NewInt(11)}

		require.NoError(t, v.PrepareHeader(header))
		report, err := vrank.DecodeReport(header.VRank)
		require.NoError(t, err)
		assert.Equal(t, candidates, report)
	})

	t.Run("non-epoch parent lookup failure is returned", func(t *testing.T) {
		v := newCN(t, withoutStart()).VRankModule
		header := &types.Header{Number: big.NewInt(11)}

		assert.ErrorIs(t, v.PrepareHeader(header), vrank.ErrHeaderNotFound)
	})
}
