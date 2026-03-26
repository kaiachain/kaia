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
	"github.com/kaiachain/kaia/kaiax/vrank"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostInsertBlock verifies the checkpoint-writing behaviour of PostInsertBlock.
func TestPostInsertBlock(t *testing.T) {
	cp := scoreCheckpointInterval
	P1, C1 := addrN(1), addrN(10)

	// at_interval: block lands exactly on a scoreCheckpointInterval boundary →
	// combined PFS+CFS checkpoint and lastCheckpoint pointer must be written.
	t.Run("at_interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		db := database.NewMemDB()

		// Build just enough headers around the checkpoint block.
		// Seed the cache at cp-1 so GetPFS/GetCFS use the nearby-hit branch
		// instead of recomputing the full epoch from block 0.
		headers := map[uint64]*types.Header{
			cp - 1: makeHeaderWithRound(cp-1, 0),
			cp:     makeHeaderWithVRank(cp, 0, []common.Address{C1}),
		}
		v := newTestModuleWithHeaders(t, valset, db, headers)
		v.pfsCache.Add(cp-1, map[common.Address]uint64{})
		v.cpMatrixCache.Add(cp-1, vrank.CPMatrix{C1: {}})

		// Block cp has a non-empty cfReport, so GetProposer is called once.
		valset.EXPECT().GetProposer(cp, uint64(0)).Return(P1, nil).Times(1)
		valset.EXPECT().GetCommittee(cp, uint64(0)).Return([]common.Address{P1}, nil).Times(1)

		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(cp))})
		require.NoError(t, v.PostInsertBlock(block))

		_, ok := v.pfsCache.Get(cp)
		assert.True(t, ok, "pfsCache should be populated at checkpoint")
		_, ok = v.cpMatrixCache.Get(cp)
		assert.True(t, ok, "cpMatrixCache should be populated at checkpoint")

		pfsCP, cpMatrixCP := ReadCheckpoint(db, cp)
		assert.NotNil(t, pfsCP, "checkpoint should be written at checkpoint block")
		assert.NotNil(t, cpMatrixCP, "checkpoint should be written at checkpoint block")
		lastCP, hasCP := ReadLastCheckpoint(db)
		assert.True(t, hasCP, "lastCheckpoint should be set")
		assert.Equal(t, cp, lastCP, "lastCheckpoint should point to the checkpoint block")
	})

	// not_at_interval: block does NOT land on a boundary →
	// no new checkpoint must be written and lastCheckpoint must stay at cp.
	t.Run("not_at_interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		db := database.NewMemDB()

		// Provide headers for cp and cp+1 (cp+1 has no VRank → empty cfReport → GetProposer not called).
		headers := map[uint64]*types.Header{
			cp:     makeHeaderWithVRank(cp, 0, []common.Address{C1}),
			cp + 1: makeHeaderWithRound(cp+1, 0),
		}
		v := newTestModuleWithHeaders(t, valset, db, headers)

		// Pre-seed the cache and DB to simulate PostInsertBlock(cp) having already run.
		v.pfsCache.Add(cp, map[common.Address]uint64{P1: 1})
		v.cpMatrixCache.Add(cp, vrank.CPMatrix{C1: {P1: 1}})
		WriteCheckpoint(db, cp,
			map[common.Address]uint64{P1: 1},
			vrank.CPMatrix{C1: {P1: 1}},
		)
		WriteLastCheckpoint(db, cp)

		// Block cp+1 has no VRank → empty cfReport → GetProposer not called.
		// generateCFSFromCPMatrix still calls GetCommittee to determine F.
		valset.EXPECT().GetCommittee(cp+1, uint64(0)).Return([]common.Address{P1}, nil).Times(1)

		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(cp + 1))})
		require.NoError(t, v.PostInsertBlock(block))

		pfsNone, _ := ReadCheckpoint(db, cp+1)
		assert.Nil(t, pfsNone, "checkpoint should NOT be written for a non-interval block")
		lastCP, hasCP := ReadLastCheckpoint(db)
		assert.True(t, hasCP, "lastCheckpoint pointer should still exist")
		assert.Equal(t, cp, lastCP, "lastCheckpoint should remain at cp, not advance to cp+1")
	})
}
