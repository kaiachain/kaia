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

	t.Run("epoch-start: VRank must be empty", func(t *testing.T) {
		v := newCN(t).VRankModule
		h := &types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}
		assert.NoError(t, v.VerifyHeader(h, nil))
		h = makeVRankHeader(t, params.DefaultVRankEpoch, []common.Address{C1})
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrUnexpectedVRankAtEpochStart)
	})

	t.Run("post-fork non-epoch: empty VRank passes", func(t *testing.T) {
		v := newCN(t).VRankModule
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

	t.Run("invalid encoding rejected", func(t *testing.T) {
		v := newCN(t).VRankModule
		h := &types.Header{Number: big.NewInt(100), VRank: []byte{0xff, 0xfe}} // garbage
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrInvalidVRankFormat)
	})

	t.Run("candidate membership is checked against the reported block N-1", func(t *testing.T) {
		headerNum := uint64(101)
		cand, futureCand := numToAddr(10), numToAddr(11)
		cn := newCN(t, withoutStart())
		v := cn.VRankModule

		// headerNum must be generated with CandidateSet at headerNum-1
		cn.Valset.EXPECT().GetCandTesting(headerNum-1).Return([]common.Address{cand}, nil).Times(2)

		assert.NoError(t, v.VerifyHeader(makeVRankHeader(t, headerNum, []common.Address{cand}), nil))
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, headerNum, []common.Address{futureCand}), nil), vrank.ErrInvalidVRankCandidate)
	})
}
