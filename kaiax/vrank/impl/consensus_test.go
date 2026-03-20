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

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeVRankHeader builds a header at the given block number containing the encoded cfReport.
// cfAddrs may be nil for an empty VRank field.
func makeVRankHeader(t *testing.T, number uint64, cfAddrs []common.Address) *types.Header {
	t.Helper()
	h := &types.Header{Number: big.NewInt(int64(number))}
	if len(cfAddrs) > 0 {
		encoded, err := vrank.EncodeReport(vrank.Report(cfAddrs))
		require.NoError(t, err)
		h.VRank = encoded
	}
	return h
}

// TestVerifyHeader covers all branches of VRankModule.VerifyHeader.
func TestVerifyHeader(t *testing.T) {
	C1, C2, C3 := addrN(1), addrN(2), addrN(3)
	candidates := []common.Address{C1, C2, C3}

	// noCallValset returns a mock that must never be called.
	noCallValset := func(t *testing.T) *mock_valset.MockValsetModule {
		t.Helper()
		return mock_valset.NewMockValsetModule(gomock.NewController(t))
	}

	// withCandidates returns a mock that returns `candidates` for GetCandidates(number).
	withCandidates := func(t *testing.T, number uint64) *mock_valset.MockValsetModule {
		t.Helper()
		vs := mock_valset.NewMockValsetModule(gomock.NewController(t))
		vs.EXPECT().GetCandidates(number).Return(candidates, nil).Times(1)
		return vs
	}

	t.Run("pre-fork: empty VRank passes", func(t *testing.T) {
		v := newPreForkModule(t, noCallValset(t))
		h := &types.Header{Number: big.NewInt(100)}
		assert.NoError(t, v.VerifyHeader(h))
	})

	t.Run("pre-fork: non-empty VRank rejected", func(t *testing.T) {
		v := newPreForkModule(t, noCallValset(t))
		assert.ErrorIs(t, v.VerifyHeader(makeVRankHeader(t, 100, []common.Address{C1})),
			vrank.ErrUnexpectedVRankBeforePermissionless)
	})

	t.Run("epoch-start: empty VRank passes", func(t *testing.T) {
		vs := noCallValset(t)
		_, module := newTestModule(t, vs, database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		h := &types.Header{Number: big.NewInt(int64(vrank.Epoch))}
		assert.NoError(t, module.VerifyHeader(h))
	})

	t.Run("epoch-start: non-empty VRank rejected", func(t *testing.T) {
		vs := noCallValset(t)
		_, module := newTestModule(t, vs, database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		assert.ErrorIs(t, module.VerifyHeader(makeVRankHeader(t, vrank.Epoch, []common.Address{C1})),
			vrank.ErrUnexpectedVRankAtEpochStart)
	})

	t.Run("post-fork non-epoch: empty VRank passes", func(t *testing.T) {
		vs := noCallValset(t)
		_, module := newTestModule(t, vs, database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		h := &types.Header{Number: big.NewInt(100)}
		assert.NoError(t, module.VerifyHeader(h))
	})

	t.Run("valid sorted report with known candidates passes", func(t *testing.T) {
		const num = uint64(100)
		_, module := newTestModule(t, withCandidates(t, num), database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		assert.NoError(t, module.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, C2, C3})))
	})

	t.Run("unknown address rejected", func(t *testing.T) {
		const num = uint64(100)
		unknown := addrN(99)
		_, module := newTestModule(t, withCandidates(t, num), database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		assert.ErrorIs(t, module.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, unknown})),
			vrank.ErrInvalidVRankCandidate)
	})

	t.Run("duplicate address rejected", func(t *testing.T) {
		const num = uint64(100)
		_, module := newTestModule(t, withCandidates(t, num), database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		assert.ErrorIs(t, module.VerifyHeader(makeVRankHeader(t, num, []common.Address{C1, C1})),
			vrank.ErrDuplicateVRankCandidate)
	})

	t.Run("unsorted addresses rejected", func(t *testing.T) {
		const num = uint64(100)
		_, module := newTestModule(t, withCandidates(t, num), database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		// C3 > C2, so C3 before C2 is not ascending.
		assert.ErrorIs(t, module.VerifyHeader(makeVRankHeader(t, num, []common.Address{C3, C2})),
			vrank.ErrVRankNotSorted)
	})

	t.Run("invalid encoding rejected", func(t *testing.T) {
		vs := noCallValset(t)
		_, module := newTestModule(t, vs, database.NewMemDB(), &testChain{headers: map[uint64]*types.Header{}})
		h := &types.Header{Number: big.NewInt(100), VRank: []byte{0xff, 0xfe}} // garbage
		assert.ErrorIs(t, module.VerifyHeader(h), vrank.ErrInvalidVRankFormat)
	})
}
