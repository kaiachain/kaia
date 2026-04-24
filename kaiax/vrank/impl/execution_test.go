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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostInsertBlock verifies the checkpoint-writing behaviour of PostInsertBlock.
func TestPostInsertBlock(t *testing.T) {
	cp := testCheckpointInterval
	P1, C1 := numToAddr(1), numToAddr(10)

	t.Run("DB/cache write at checkpoint", func(t *testing.T) {
		// Build just enough headers around the checkpoint block.
		// Seed the cache at cp-1 so GetPFS/GetCFS use the nearby-hit branch
		// instead of recomputing the full epoch from block 0.
		headers := map[uint64]*types.Header{
			cp - 1: makeHeaderWithRound(cp-1, 0),
			cp:     makeHeaderWithVRank(cp, 0, []common.Address{C1}),
		}
		cn := newCN(t, withHeaders(headers))
		cn.VRankModule.pfsCache.Add(cp-1, map[common.Address]uint64{})
		cn.VRankModule.cpMatrixCache.Add(cp-1, vrank.CPMatrix{C1: {}})

		// Block cp has a non-empty cfReport, so GetProposer is called once.
		cn.Valset.EXPECT().GetProposer(cp, uint64(0)).Return(P1, nil).Times(1)

		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(cp))})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		_, ok := cn.VRankModule.pfsCache.Get(cp)
		assert.True(t, ok, "pfsCache should be populated at checkpoint")
		_, ok = cn.VRankModule.cpMatrixCache.Get(cp)
		assert.True(t, ok, "cpMatrixCache should be populated at checkpoint")

		assert.NotNil(t, ReadCheckpointPFS(cn.DB, cp), "checkpoint should be written at checkpoint block")
		assert.NotNil(t, ReadCheckpointCPMatrix(cn.DB, cp), "checkpoint should be written at checkpoint block")
		lastCP, hasCP := ReadLastCheckpoint(cn.DB)
		assert.True(t, hasCP, "lastCheckpoint should be set")
		assert.Equal(t, cp, lastCP, "lastCheckpoint should point to the checkpoint block")
	})

	t.Run("no DB write at checkpoint, only cache", func(t *testing.T) {
		headers := map[uint64]*types.Header{
			cp:     makeHeaderWithVRank(cp, 0, []common.Address{C1}),
			cp + 1: makeHeaderWithRound(cp+1, 0),
		}
		cn := newCN(t, withHeaders(headers))

		// Pre-seed the cache and DB to simulate PostInsertBlock(cp) having already run.
		cn.VRankModule.pfsCache.Add(cp, map[common.Address]uint64{P1: 1})
		cn.VRankModule.cpMatrixCache.Add(cp, vrank.CPMatrix{C1: {P1: 1}})
		WriteCheckpoint(cn.DB, cp,
			map[common.Address]uint64{P1: 1},
			vrank.CPMatrix{C1: {P1: 1}},
		)
		WriteLastCheckpoint(cn.DB, cp)

		// Block cp+1 has no VRank → empty cfReport → GetProposer not called.
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(cp + 1))})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		assert.Nil(t, ReadCheckpointPFS(cn.DB, cp+1), "checkpoint should NOT be written for a non-interval block")
		lastCP, hasCP := ReadLastCheckpoint(cn.DB)
		assert.True(t, hasCP, "lastCheckpoint pointer should still exist")
		assert.Equal(t, cp, lastCP, "lastCheckpoint should remain at cp, not advance to cp+1")
	})
}
