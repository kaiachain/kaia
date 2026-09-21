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

	t.Run("score warming failure is ignored", func(t *testing.T) {
		const (
			forkBlock = uint64(120)
			epoch     = uint64(50)
		)
		cn := newCN(t, withHeaders(map[uint64]*types.Header{
			forkBlock: makeHeaderWithRound(forkBlock, 0),
		}))
		cn.VRankModule.ChainConfig.PermissionlessCompatibleBlock = new(big.Int).SetUint64(forkBlock)
		cn.VRankModule.ChainConfig.VRankEpoch = epoch

		// This invalid configuration puts the fork in the middle of the epoch
		// [100, 149], so score warming reaches into pre-fork blocks.
		_, err := cn.VRankModule.GetPFS(forkBlock)
		require.ErrorIs(t, err, vrank.ErrNotPermissionless)

		block := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(forkBlock)})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		_, ok := cn.VRankModule.pfsCache.Get(forkBlock)
		assert.False(t, ok, "failed PFS warming should not populate the cache")
		_, ok = cn.VRankModule.cpMatrixCache.Get(forkBlock)
		assert.False(t, ok, "CP matrix warming should be skipped after PFS failure")
		assert.Nil(t, ReadCheckpointPFS(cn.DB, forkBlock), "partial checkpoint should not be written")
		_, hasCP := ReadLastCheckpoint(cn.DB)
		assert.False(t, hasCP, "lastCheckpoint should not advance after a warming failure")
	})

	t.Run("CP matrix warming failure is ignored", func(t *testing.T) {
		cn := newCN(t, withHeaders(map[uint64]*types.Header{
			cp: makeHeaderWithRound(cp, 0),
		}))
		cn.VRankModule.pfsCache.Add(cp, map[common.Address]uint64{})
		// Block 0 is intentionally absent, so epochCandidates cannot fall back to
		// the epoch-start header when GetCandTesting fails.
		cn.Valset.EXPECT().GetCandTesting(cp).Return(nil, assert.AnError)

		block := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(cp)})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		_, ok := cn.VRankModule.cpMatrixCache.Get(cp)
		assert.False(t, ok, "failed CP matrix warming should not populate the cache")
		assert.Nil(t, ReadCheckpointPFS(cn.DB, cp), "partial checkpoint should not be written")
		_, hasCP := ReadLastCheckpoint(cn.DB)
		assert.False(t, hasCP, "lastCheckpoint should not advance after a warming failure")
	})

	t.Run("later checkpoint succeeds after warming failure", func(t *testing.T) {
		const (
			epoch          = uint64(16)
			checkpoint     = epoch / 8
			nextCheckpoint = 2 * checkpoint
		)
		headers := map[uint64]*types.Header{
			1:          makeHeaderWithRound(1, 0),
			checkpoint: makeHeaderWithRound(checkpoint, 0),
		}
		cn := newCN(t, withHeaders(headers), withProposer(P1))
		cn.VRankModule.ChainConfig.VRankEpoch = epoch
		cn.VRankModule.pfsCache.Add(checkpoint, map[common.Address]uint64{})

		// The first checkpoint cannot initialize its CP matrix because neither
		// CandTesting state nor the epoch-start header is available.
		cn.Valset.EXPECT().GetCandTesting(checkpoint).Return(nil, assert.AnError)
		block := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(checkpoint)})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))
		assert.Nil(t, ReadCheckpointPFS(cn.DB, checkpoint))

		// Once the missing inputs become available, the next checkpoint succeeds
		// and advances lastCheckpoint across the gap.
		headers[0] = makeHeaderWithRound(0, 0)
		headers[3] = makeHeaderWithRound(3, 0)
		headers[nextCheckpoint] = makeHeaderWithRound(nextCheckpoint, 0)
		cn.Valset.EXPECT().GetCandTesting(nextCheckpoint).Return([]common.Address{C1}, nil)
		block = types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(nextCheckpoint)})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		assert.NotNil(t, ReadCheckpointPFS(cn.DB, nextCheckpoint))
		assert.NotNil(t, ReadCheckpointCPMatrix(cn.DB, nextCheckpoint))
		lastCP, hasCP := ReadLastCheckpoint(cn.DB)
		assert.True(t, hasCP)
		assert.Equal(t, nextCheckpoint, lastCP)
	})

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
		cn := newCN(t, withHeaders(headers), withProposer(P1))

		// Pre-seed the cache and DB to simulate PostInsertBlock(cp) having already run.
		cn.VRankModule.pfsCache.Add(cp, map[common.Address]uint64{P1: 1})
		cn.VRankModule.cpMatrixCache.Add(cp, vrank.CPMatrix{C1: {P1: 1}})
		WriteCheckpoint(cn.DB, cp,
			map[common.Address]uint64{P1: 1},
			vrank.CPMatrix{C1: {P1: 1}},
		)
		WriteLastCheckpoint(cn.DB, cp)

		// Block cp+1 has no VRank → empty cfReport.
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(cp + 1))})
		require.NoError(t, cn.VRankModule.PostInsertBlock(block))

		assert.Nil(t, ReadCheckpointPFS(cn.DB, cp+1), "checkpoint should NOT be written for a non-interval block")
		lastCP, hasCP := ReadLastCheckpoint(cn.DB)
		assert.True(t, hasCP, "lastCheckpoint pointer should still exist")
		assert.Equal(t, cp, lastCP, "lastCheckpoint should remain at cp, not advance to cp+1")
	})
}
