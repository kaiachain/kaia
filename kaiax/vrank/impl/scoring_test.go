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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	mock_randao "github.com/kaiachain/kaia/kaiax/randao/mock"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func computeCFS(v *VRankModule, start, end uint64) (map[common.Address]uint64, error) {
	cpMatrix, err := v.newCPMatrix(start)
	if err != nil {
		return nil, err
	}
	cpMatrix, err = v.applyBlocksForCPMatrix(start, end, cpMatrix)
	if err != nil {
		return nil, err
	}
	return generateCFSFromCPMatrix(cpMatrix), nil
}

// makeHeaderWithVRank creates a header with a specific round and an encoded cfReport in VRank.
// cfAddrs may be nil/empty for an empty report.
func makeHeaderWithVRank(number uint64, round int64, cfAddrs []common.Address) *types.Header {
	h := makeHeaderWithRound(number, round)
	encoded, err := vrank.EncodeReport(cfAddrs)
	if err != nil {
		panic(err)
	}
	h.VRank = encoded
	return h
}

// generateHeadersFromCPMatrix materializes per-block cfReports from a candidate/proposer matrix.
// cpMatrix is indexed as cpMatrix[candidateIndex][proposerIndex] and stores reporter totals.
func generateHeadersFromCPMatrix(
	t *testing.T,
	start uint64,
	candidates []common.Address,
	proposers []common.Address,
	cpMatrix [][]uint64,
) (map[uint64]*types.Header, map[uint64]common.Address, uint64, uint64) {
	t.Helper()
	require.Len(t, cpMatrix, len(candidates), "cpMatrix row count must equal number of candidates")
	for i := range cpMatrix {
		require.Len(t, cpMatrix[i], len(proposers), "cpMatrix col count mismatch at row %d", i)
	}

	columnMax := make([]uint64, len(proposers))
	for pi := range proposers {
		for ci := range candidates {
			if cpMatrix[ci][pi] > columnMax[pi] {
				columnMax[pi] = cpMatrix[ci][pi]
			}
		}
	}

	var totalBlocks uint64
	for _, n := range columnMax {
		totalBlocks += n
	}

	headers := make(map[uint64]*types.Header, int(totalBlocks))
	blockReporter := make(map[uint64]common.Address, int(totalBlocks))

	nextBlock := start
	for pi, p := range proposers {
		for j := uint64(0); j < columnMax[pi]; j++ {
			report := make([]common.Address, 0, len(candidates))
			for ci, c := range candidates {
				if j < cpMatrix[ci][pi] {
					report = append(report, c)
				}
			}
			headers[nextBlock] = makeHeaderWithVRank(nextBlock, 0, report)
			blockReporter[nextBlock] = p
			nextBlock++
		}
	}

	return headers, blockReporter, nextBlock - 1, totalBlocks
}

// # GetPFS
func TestGetPFS(t *testing.T) {
	proposer := numToAddr(1)
	t.Run("three proposers", func(t *testing.T) {
		p1, p2, p3 := numToAddr(1), numToAddr(2), numToAddr(3)
		headers := map[uint64]*types.Header{
			0: makeHeaderWithRound(0, 0),
			1: makeHeaderWithRound(1, 0),
			2: makeHeaderWithRound(2, 1),
			3: makeHeaderWithRound(3, 2),
		}

		cn := newCN(t, withHeaders(headers))
		cn.Valset.EXPECT().GetProposer(uint64(2), uint64(0)).Return(p1, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(uint64(3), uint64(0)).Return(p1, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(uint64(3), uint64(1)).Return(p2, nil).Times(1)

		v := cn.VRankModule
		pfs, err := v.GetPFS(1)
		require.NoError(t, err)
		assert.Len(t, pfs, 0)

		pfs, err = v.GetPFS(2)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), pfs[p1])
		assert.Len(t, pfs, 1)

		pfs, err = v.GetPFS(3)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), pfs[p1])
		assert.Equal(t, uint64(1), pfs[p2])
		assert.NotContains(t, pfs, p3)
		assert.Len(t, pfs, 2)
	})
	t.Run("new epoch - scan one header", func(t *testing.T) {
		epochStart := uint64(params.DefaultVRankEpoch)
		headers := map[uint64]*types.Header{epochStart: makeHeaderWithRound(epochStart, 1)}

		cn := newCN(t, withHeaders(headers))
		v := cn.VRankModule

		// Seed the previous epoch's last block with non-zero PFS.
		v.pfsCache.Add(epochStart-1, map[common.Address]uint64{proposer: 99})
		pfs, err := v.GetPFS(epochStart - 1)
		assert.Equal(t, uint64(99), pfs[proposer])

		cn.Valset.EXPECT().GetProposer(uint64(epochStart), uint64(0)).Return(proposer, nil).Times(1)
		pfs, err = v.GetPFS(epochStart)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), pfs[proposer], "PFS must reset at epoch")
	})
	t.Run("epoch boundary clamp", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		db := database.NewMemDB()

		epochStart := uint64(params.DefaultVRankEpoch)
		blockNum := epochStart + 5
		headers := map[uint64]*types.Header{}
		for i := epochStart; i <= blockNum; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		v := newCN(t, withValset(valset), withDB(db), withHeaders(headers)).VRankModule

		// Seed the last block of the previous epoch with a large PFS value.
		v.pfsCache.Add(epochStart-1, map[common.Address]uint64{numToAddr(0): 99})

		// The probe limit is min(64, blockNum-epochStart)=5, so epochStart-1 (distance=6) is
		// never reached. GetProposer must not be called because all blocks in [epochStart, blockNum]
		// have round=0.
		pfs, err := v.GetPFS(blockNum)
		require.NoError(t, err)
		assert.Empty(t, pfs, "previous epoch cache must not carry over into current epoch")
	})
	t.Run("defensive copy", func(t *testing.T) {
		headers := map[uint64]*types.Header{}
		for i := range uint64(5) {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[5] = makeHeaderWithRound(5, 1)

		proposer := numToAddr(1)
		v := newCN(t, withProposer(proposer), withHeaders(headers)).VRankModule

		pfs1, err := v.GetPFS(5)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), pfs1[proposer])

		// Mutating the returned map must not affect the cached copy.
		pfs1[proposer] = 999
		pfs2, err := v.GetPFS(5)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), pfs2[proposer])
	})
}

func TestGetPFS_Errors(t *testing.T) {
	t.Run("pre-fork block returns ErrNotPermissionless", func(t *testing.T) {
		module := newCN(t, withHardfork("osaka"), withHeaders(map[uint64]*types.Header{10: makeHeaderWithRound(10, 0)})).VRankModule
		_, err := module.GetPFS(10)
		assert.ErrorIs(t, err, vrank.ErrNotPermissionless)
	})
	t.Run("future block (beyond head) returns ErrFutureBlock", func(t *testing.T) {
		headers := make(map[uint64]*types.Header, 6)
		for i := uint64(0); i <= 5; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		v := newCN(t, withHeaders(headers), withProposer(numToAddr(0))).VRankModule

		_, err := v.GetPFS(5)
		assert.NoError(t, err)

		_, err = v.GetPFS(6)
		assert.ErrorIs(t, err, vrank.ErrFutureBlock)
	})
	t.Run("GetProposer error", func(t *testing.T) {
		headers := map[uint64]*types.Header{
			0: makeHeaderWithRound(0, 0),
			1: makeHeaderWithRound(1, 1), // round=1 → needs GetProposer(1, 0)
		}
		cn := newCN(t, withHeaders(headers))
		v := cn.VRankModule

		cn.Valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(common.Address{}, assert.AnError)

		_, err := v.GetPFS(1)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("missing header error", func(t *testing.T) {
		// Header exists for block 1 (so CurrentHeader()=1 and ErrFutureBlock is not triggered)
		// but NOT for block 0, which applyBlocksForPFS will try to fetch first.
		headers := map[uint64]*types.Header{1: makeHeaderWithRound(1, 0)}
		cn := newCN(t, withHeaders(headers))
		v := cn.VRankModule

		_, err := v.GetPFS(1)
		assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
	})
}

func TestGetPFS_CacheFallback(t *testing.T) {
	proposer := numToAddr(1)

	t.Run("nearby cache hit", func(t *testing.T) {
		headers := map[uint64]*types.Header{}
		for i := range uint64(5) {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[5] = makeHeaderWithRound(5, 1)
		headers[6] = makeHeaderWithRound(6, 1)

		v := newCN(t, withHeaders(headers), withProposer(proposer)).VRankModule

		_, err := v.GetPFS(5)
		require.NoError(t, err)
		pfs5, ok := v.pfsCache.Get(uint64(5))
		assert.True(t, ok)
		assert.Equal(t, uint64(1), pfs5.(map[common.Address]uint64)[proposer])
		_, ok = v.pfsCache.Get(uint64(6))
		assert.False(t, ok)

		_, err = v.GetPFS(6)
		require.NoError(t, err)
		pfs6, ok := v.pfsCache.Get(uint64(6))
		assert.Equal(t, uint64(2), pfs6.(map[common.Address]uint64)[proposer])
		assert.True(t, ok)
	})
	t.Run("no cache hit - scan all headers", func(t *testing.T) {
		headers := map[uint64]*types.Header{
			0: makeHeaderWithRound(0, 0),
			1: makeHeaderWithRound(1, 1),
			2: makeHeaderWithRound(2, 2),
		}

		// ensure no DB, no cache. Scans all headers in this case.
		cn := newCN(t, withHeaders(headers))
		for i := range 3 {
			_, inCache := cn.VRankModule.pfsCache.Get(uint64(i))
			require.False(t, inCache, "cache must be cold before first call")
		}

		cn.Valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(uint64(2), uint64(0)).Return(proposer, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(uint64(2), uint64(1)).Return(proposer, nil).Times(1)

		pfs, err := cn.VRankModule.GetPFS(2)
		require.NoError(t, err)
		assert.Equal(t, uint64(3), pfs[proposer])
	})
	t.Run("nearby cache hit - cache at boundary must be used", func(t *testing.T) {
		const blockNum = uint64(100)
		const boundary = blockNum - scoreCacheProbeLookback // distance exactly 64

		headers := map[uint64]*types.Header{}
		for i := range blockNum {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[blockNum] = makeHeaderWithRound(blockNum, 1) // one failure

		cn := newCN(t, withHeaders(headers))

		// forge cache - it must be used
		cn.VRankModule.pfsCache.Add(boundary, map[common.Address]uint64{proposer: 12345})

		cn.Valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(proposer, nil).Times(1)

		pfs, err := cn.VRankModule.GetPFS(blockNum)
		require.NoError(t, err)
		assert.Equal(t, uint64(12346), pfs[proposer], "cache at exact boundary must be used as seed")
	})
	t.Run("no cache hit - cache beyond boundary must not be used", func(t *testing.T) {
		const blockNum = uint64(100)
		const boundary = blockNum - scoreCacheProbeLookback // distance exactly 64

		headers := map[uint64]*types.Header{}
		for i := range blockNum {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[blockNum] = makeHeaderWithRound(blockNum, 1) // one failure

		cn := newCN(t, withHeaders(headers))

		// forge cache - it must NOT be used
		cn.VRankModule.pfsCache.Add(boundary-1, map[common.Address]uint64{proposer: 12345})

		cn.Valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(proposer, nil).Times(1)

		pfs, err := cn.VRankModule.GetPFS(blockNum)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), pfs[proposer], "cache beyond boundary must NOT be used as seed")
	})
}

func TestGetPFS_DBFallback(t *testing.T) {
	proposer := numToAddr(1)
	ckpt := testCheckpointInterval

	t.Run("checkpoint hit - returns seed verbatim", func(t *testing.T) {
		headers := map[uint64]*types.Header{ckpt: makeHeaderWithRound(ckpt, 0)}

		db := database.NewMemDB()
		WriteCheckpoint(db, ckpt, map[common.Address]uint64{proposer: 12345}, vrank.CPMatrix{})
		WriteLastCheckpoint(db, ckpt)

		cn := newCN(t, withDB(db), withHeaders(headers))

		pfs, err := cn.VRankModule.GetPFS(ckpt)
		require.NoError(t, err)
		assert.Equal(t, uint64(12345), pfs[proposer], "GetPFS(ckpt) must return the DB-stored checkpoint seed verbatim")
	})
	t.Run("checkpoint hit - advance one block accumulates", func(t *testing.T) {
		headers := map[uint64]*types.Header{ckpt + 1: makeHeaderWithRound(ckpt+1, 1)}

		db := database.NewMemDB()
		WriteCheckpoint(db, ckpt, map[common.Address]uint64{proposer: 12345}, vrank.CPMatrix{})
		WriteLastCheckpoint(db, ckpt)

		cn := newCN(t, withDB(db), withHeaders(headers))
		// ckpt+1 has round=1 → adds one pfReport on top of the checkpoint seed.
		// No epoch reset because ckpt is not an epoch boundary (ckpt = vrankEpoch/8).
		cn.Valset.EXPECT().GetProposer(uint64(ckpt+1), uint64(0)).Return(proposer, nil).Times(1)

		pfs, err := cn.VRankModule.GetPFS(ckpt + 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(12346), pfs[proposer], "12345 from checkpoint + 1 from ckpt+1")
	})
}

// # GetCFS
func TestGetCFS(t *testing.T) {
	t.Run("epoch boundary", func(t *testing.T) {
		P1, C1 := numToAddr(1), numToAddr(10)
		epochStart := uint64(params.DefaultVRankEpoch)
		epochEnd := uint64(2*params.DefaultVRankEpoch - 1)

		headers := make(map[uint64]*types.Header, params.DefaultVRankEpoch)
		for i := uint64(0); i < params.DefaultVRankEpoch; i++ {
			headers[epochStart+i] = makeHeaderWithRound(epochStart+i, 0)
		}
		headers[epochStart+1] = makeHeaderWithVRank(epochStart+1, 0, []common.Address{C1})
		v := newCN(t, withGenesis(), withCandidates([]common.Address{C1}), withProposer(P1)).VRankModule
		v.Chain = &testChain{headers: headers}

		cfs, err := v.GetCFS(epochStart)
		require.NoError(t, err)
		assert.Equal(t, uint64(0), cfs[C1], "C1 should have 0 candidate failures at epoch start")

		cfs, err = v.GetCFS(epochEnd)
		require.NoError(t, err)
		// ProposerCount = 1 (P1 only), F = 1 - ceil(2/3) = 0, so no filtering; C1 appears once.
		assert.Equal(t, uint64(1), cfs[C1], "C1 should have 1 candidate failure")
	})

	t.Run("multi-reporter epoch", func(t *testing.T) {
		epochStart := uint64(params.DefaultVRankEpoch)
		P1, P2, P3, P4 := numToAddr(30), numToAddr(31), numToAddr(32), numToAddr(33)
		C1, C2, C3 := numToAddr(40), numToAddr(41), numToAddr(42)
		candidates := []common.Address{C1, C2, C3}
		b := newCN(t, withGenesis(), withCandidates(candidates))
		v := b.VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				epochStart:     makeHeaderWithVRank(epochStart, 0, nil),
				epochStart + 1: makeHeaderWithVRank(epochStart+1, 0, []common.Address{C1, C2}),
				epochStart + 2: makeHeaderWithVRank(epochStart+2, 1, []common.Address{C1}),
				epochStart + 3: makeHeaderWithVRank(epochStart+3, 0, []common.Address{C1, C2, C3}),
				epochStart + 4: makeHeaderWithVRank(epochStart+4, 0, []common.Address{C1, C2}),
				epochStart + 5: makeHeaderWithVRank(epochStart+5, 0, []common.Address{C1, C2}),
			},
		}

		// Per-block proposer assignments; can't use withProposer because each block has a different one.
		b.Valset.EXPECT().GetProposer(epochStart, uint64(0)).Return(P1, nil)
		b.Valset.EXPECT().GetProposer(epochStart+1, uint64(0)).Return(P1, nil)
		b.Valset.EXPECT().GetProposer(epochStart+2, uint64(1)).Return(P2, nil)
		b.Valset.EXPECT().GetProposer(epochStart+3, uint64(0)).Return(P3, nil)
		b.Valset.EXPECT().GetProposer(epochStart+4, uint64(0)).Return(P4, nil)
		b.Valset.EXPECT().GetProposer(epochStart+5, uint64(0)).Return(P4, nil)
		cfs, err := v.GetCFS(epochStart + 5)
		require.NoError(t, err)
		// ProposerCount = 4, F = 4 - ceil(8/3) = 1. Sorted scores per candidate:
		// C1: [1,1,1,2]  -> drop 1 -> sum 3
		// C2: [0,1,1,2]  -> drop 1 -> sum 2
		// C3: [0,0,0,1]  -> drop 1 -> sum 0
		assert.Equal(t, uint64(3), cfs[C1])
		assert.Equal(t, uint64(2), cfs[C2])
		assert.Equal(t, uint64(0), cfs[C3])
		assert.Len(t, cfs, 3)
	})

	t.Run("defensive copy", func(t *testing.T) {
		P1, C1 := numToAddr(1), numToAddr(10)
		headers := map[uint64]*types.Header{}
		for i := uint64(0); i <= 5; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[5] = makeHeaderWithVRank(5, 0, []common.Address{C1})

		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}), withProposer(P1))

		cfs1, err := cn.VRankModule.GetCFS(5)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs1[C1])

		// Mutating the returned map must not affect the cached copy.
		cfs1[C1] = 999
		cfs2, err := cn.VRankModule.GetCFS(5)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs2[C1])
	})

	t.Run("epoch start clears previous-epoch state", func(t *testing.T) {
		epochStart := uint64(params.DefaultVRankEpoch)
		C1, C2, P1 := numToAddr(10), numToAddr(11), numToAddr(1)
		headers := map[uint64]*types.Header{
			epochStart: makeHeaderWithVRank(epochStart, 0, nil), // epoch start: VRank must be nil
		}
		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1, C2}), withProposer(P1))

		// Seed the previous epoch's last block with large CFS values.
		cn.VRankModule.cpMatrixCache.Add(epochStart-1, vrank.CPMatrix{
			C1: {P1: 50}, C2: {P1: 50},
		})

		cfs, err := cn.VRankModule.GetCFS(epochStart)
		require.NoError(t, err)
		assert.Equal(t, uint64(0), cfs[C1], "C1 must start at 0 in new epoch")
		assert.Equal(t, uint64(0), cfs[C2], "C2 must start at 0 in new epoch")
		assert.Len(t, cfs, 2, "all candidates must appear even with zero score")
	})
}

func TestGetCFS_Errors(t *testing.T) {
	t.Run("pre-fork block returns ErrNotPermissionless", func(t *testing.T) {
		module := newCN(t, withHardfork("osaka"), withHeaders(map[uint64]*types.Header{10: makeHeaderWithRound(10, 0)})).VRankModule
		_, err := module.GetCFS(10)
		assert.ErrorIs(t, err, vrank.ErrNotPermissionless)
	})
	t.Run("future block (beyond head) returns ErrFutureBlock", func(t *testing.T) {
		headers := make(map[uint64]*types.Header, 6)
		for i := uint64(0); i <= 5; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		v := newCN(t, withHeaders(headers), withCandidates(nil), withProposer(numToAddr(0))).VRankModule

		_, err := v.GetCFS(5)
		assert.NoError(t, err)

		_, err = v.GetCFS(6)
		assert.ErrorIs(t, err, vrank.ErrFutureBlock)
	})
	t.Run("GetProposer error", func(t *testing.T) {
		C1 := numToAddr(10)
		headers := map[uint64]*types.Header{
			0: makeHeaderWithRound(0, 0),
			1: makeHeaderWithVRank(1, 0, []common.Address{C1}),
		}
		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}))
		// Block 0 is the first block applied for the [0, 1] range, so its GetProposer is the first call — surface the error there and verify it propagates.
		cn.Valset.EXPECT().GetProposer(uint64(0), uint64(0)).Return(common.Address{}, assert.AnError).Times(1)

		_, err := cn.VRankModule.GetCFS(1)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("missing header error", func(t *testing.T) {
		C1 := numToAddr(10)
		// Header exists for block 1 (so CurrentHeader()=1 and ErrFutureBlock is not triggered)
		// but NOT for block 0, which applyCPMatrix will try to fetch first.
		headers := map[uint64]*types.Header{1: makeHeaderWithRound(1, 0)}
		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}))

		_, err := cn.VRankModule.GetCFS(1)
		assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
	})
}

func TestGetCFS_CacheFallback(t *testing.T) {
	P1, C1 := numToAddr(1), numToAddr(10)

	t.Run("nearby cache hit", func(t *testing.T) {
		P2 := numToAddr(2)
		headers := map[uint64]*types.Header{}
		for i := uint64(0); i <= 6; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[5] = makeHeaderWithVRank(5, 0, []common.Address{C1})
		headers[6] = makeHeaderWithVRank(6, 0, []common.Address{C1})

		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}))
		// Specific overrides for blocks 5 and 6 (different reporters); P1 wildcard catches the rest.
		cn.Valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(P1, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(uint64(6), uint64(0)).Return(P2, nil).Times(1)
		cn.Valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(P1, nil).AnyTimes()

		cfs5, err := cn.VRankModule.GetCFS(5)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs5[C1])
		_, ok := cn.VRankModule.cpMatrixCache.Get(uint64(5))
		assert.True(t, ok)

		cfs6, err := cn.VRankModule.GetCFS(6)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), cfs6[C1])
	})

	t.Run("no cache hit - scan all headers", func(t *testing.T) {
		headers := map[uint64]*types.Header{
			0: makeHeaderWithRound(0, 0),
			1: makeHeaderWithVRank(1, 0, []common.Address{C1}),
			2: makeHeaderWithVRank(2, 0, []common.Address{C1}),
		}
		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}), withProposer(P1))
		_, inCache := cn.VRankModule.cpMatrixCache.Get(uint64(2))
		require.False(t, inCache, "cache must be cold before first call")

		cfs, err := cn.VRankModule.GetCFS(2)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), cfs[C1])
	})

	t.Run("nearby cache hit - cache at boundary must be used", func(t *testing.T) {
		const blockNum = uint64(100)
		const boundary = blockNum - scoreCacheProbeLookback // distance exactly 64

		headers := map[uint64]*types.Header{}
		for i := range blockNum {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[blockNum] = makeHeaderWithVRank(blockNum, 0, []common.Address{C1})

		cn := newCN(t, withHeaders(headers), withProposer(P1))

		// forge cache - it must be used
		cn.VRankModule.cpMatrixCache.Add(boundary, vrank.CPMatrix{C1: {P1: 12345}})

		cfs, err := cn.VRankModule.GetCFS(blockNum)
		require.NoError(t, err)
		assert.Equal(t, uint64(12346), cfs[C1], "cache at exact boundary must be used as seed")
	})

	t.Run("no cache hit - cache beyond boundary must not be used", func(t *testing.T) {
		const blockNum = uint64(100)
		const boundary = blockNum - scoreCacheProbeLookback

		headers := map[uint64]*types.Header{}
		for i := range blockNum {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		headers[blockNum] = makeHeaderWithVRank(blockNum, 0, []common.Address{C1})

		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}), withProposer(P1))

		// forge cache - it must NOT be used (beyond boundary by 1)
		cn.VRankModule.cpMatrixCache.Add(boundary-1, vrank.CPMatrix{C1: {P1: 12345}})

		cfs, err := cn.VRankModule.GetCFS(blockNum)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs[C1], "cache beyond boundary must NOT be used as seed")
	})

	t.Run("nearby cache hit - same epoch only", func(t *testing.T) {
		epochStart := uint64(params.DefaultVRankEpoch)
		headers := map[uint64]*types.Header{
			epochStart:     makeHeaderWithRound(epochStart, 0),
			epochStart + 1: makeHeaderWithVRank(epochStart+1, 0, []common.Address{C1}),
			epochStart + 2: makeHeaderWithRound(epochStart+2, 0),
		}
		cn := newCN(t, withHeaders(headers), withProposer(P1))

		// Seed cache at epochStart+1 (same epoch, distance=1) — must be used.
		cn.VRankModule.cpMatrixCache.Add(epochStart+1, vrank.CPMatrix{C1: {P1: 1}})
		// Seed previous-epoch entry — must NOT be used.
		cn.VRankModule.cpMatrixCache.Add(epochStart-1, vrank.CPMatrix{C1: {P1: 99}})

		cfs, err := cn.VRankModule.GetCFS(epochStart + 2)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs[C1], "score must come only from same-epoch nearby cache")
	})

	t.Run("no cache hit - cache from previous epoch must not be used", func(t *testing.T) {
		epochStart := uint64(params.DefaultVRankEpoch)
		blockNum := epochStart + 3
		headers := map[uint64]*types.Header{}
		for i := epochStart; i <= blockNum; i++ {
			headers[i] = makeHeaderWithRound(i, 0)
		}
		cn := newCN(t, withHeaders(headers), withCandidates([]common.Address{C1}), withProposer(P1))

		// Probe limit = min(64, blockNum-epochStart) = 3; epochStart-1 is at distance 4 — not reached.
		// Even within scoreCacheProbeLookback range, the epoch clamp blocks the probe from reading
		// across the epoch boundary.
		cn.VRankModule.cpMatrixCache.Add(epochStart-1, vrank.CPMatrix{C1: {P1: 99}})

		cfs, err := cn.VRankModule.GetCFS(blockNum)
		require.NoError(t, err)
		assert.Equal(t, uint64(0), cfs[C1], "previous-epoch cache must not carry over into current epoch")
	})
}

func TestGetCFS_DBFallback(t *testing.T) {
	P1, C1 := numToAddr(1), numToAddr(10)
	ckpt := testCheckpointInterval

	t.Run("checkpoint hit - returns seed verbatim", func(t *testing.T) {
		C2 := numToAddr(11)
		headers := map[uint64]*types.Header{ckpt: makeHeaderWithRound(ckpt, 0)}

		db := database.NewMemDB()
		WriteCheckpoint(db, ckpt, map[common.Address]uint64{}, vrank.CPMatrix{
			C1: {P1: 1},
			C2: {},
		})
		WriteLastCheckpoint(db, ckpt)

		cn := newCN(t, withDB(db), withHeaders(headers))

		cfs, err := cn.VRankModule.GetCFS(ckpt)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), cfs[C1])
		assert.Equal(t, uint64(0), cfs[C2], "zero-failure candidates must remain in CFS after checkpoint load")

		cpCached, ok := cn.VRankModule.cpMatrixCache.Get(ckpt)
		require.True(t, ok, "checkpoint-backed cpMatrix must be cached")
		cpMatrix := cpCached.(vrank.CPMatrix)
		assert.Contains(t, cpMatrix, C2)
		assert.Equal(t, uint64(0), cpMatrix[C2][P1])
	})

	t.Run("checkpoint hit - advance one block accumulates", func(t *testing.T) {
		headers := map[uint64]*types.Header{
			ckpt + 1: makeHeaderWithVRank(ckpt+1, 0, []common.Address{C1}),
		}

		db := database.NewMemDB()
		WriteCheckpoint(db, ckpt, map[common.Address]uint64{}, vrank.CPMatrix{
			C1: {P1: 1},
		})
		WriteLastCheckpoint(db, ckpt)

		cn := newCN(t, withDB(db), withHeaders(headers), withProposer(P1))

		cfs, err := cn.VRankModule.GetCFS(ckpt + 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), cfs[C1], "1 from checkpoint + 1 from ckpt+1")
	})
}

// TestApplyBlocksForPFS verifies the raw proposer-failure accumulation over a small block range.
func TestApplyBlocksForPFS(t *testing.T) {
	// Set up three consecutive blocks:
	//   block 10, round 0: no pfReport (first proposer succeeded)
	//   block 11, round 2: pfReport = [p1, p2]  (rounds 0 and 1 failed)
	//   block 12, round 1: pfReport = [p1]       (round 0 failed)
	//
	// Expected PFS: {p1: 2, p2: 1}
	p1, p2 := numToAddr(1), numToAddr(2)

	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	randao := mock_randao.NewMockRandaoModule(ctrl)
	v := newCN(t, withValset(valset), withRandao(randao), withGenesis()).VRankModule
	v.Chain = &testChain{
		headers: map[uint64]*types.Header{
			10: makeHeaderWithRound(10, 0),
			11: makeHeaderWithRound(11, 2),
			12: makeHeaderWithRound(12, 1),
		},
	}

	// block 10, round 0: no GetProposer calls (loop [0,0) is empty)
	// block 11, round 2: proposers for rounds 0 and 1
	valset.EXPECT().GetProposer(uint64(11), uint64(0)).Return(p1, nil)
	valset.EXPECT().GetProposer(uint64(11), uint64(1)).Return(p2, nil)
	// block 12, round 1: proposer for round 0
	valset.EXPECT().GetProposer(uint64(12), uint64(0)).Return(p1, nil)

	pfs, err := v.applyBlocksForPFS(10, 12, make(map[common.Address]uint64))
	require.NoError(t, err)
	assert.Equal(t, uint64(2), pfs[p1])
	assert.Equal(t, uint64(1), pfs[p2])
	assert.Len(t, pfs, 2)
}

// TestApplyBlocksForCFS verifies CFS computation via applyCPMatrix + generateCFSFromCPMatrix.
func TestApplyBlocksForCFS(t *testing.T) {
	t.Run("KIP227_Example1", func(t *testing.T) {
		// Reproduces KIP-227 Example 1 through the chain-reading path.
		// Blocks 5–9 are the epoch range (matching the example's block numbers).
		//
		//   proposer(5)=P1, cfReport(5)=[]        (epoch start: must be empty per spec)
		//   proposer(6)=P2, cfReport(6)=[]
		//   proposer(7)=P3, cfReport(7)=[C1,C2,C3]
		//   proposer(8)=P4, cfReport(8)=[C1,C2]
		//   proposer(9)=P4, cfReport(9)=[C1,C2]
		//
		// 4 distinct proposers → F = N - ceil(2N/3) = 4 - 3 = 1.
		// Expected CFS: C1=1, C2=1, C3=0 (present in map with zero score)
		P1, P2, P3, P4 := numToAddr(10), numToAddr(11), numToAddr(12), numToAddr(13)
		C1, C2, C3 := numToAddr(20), numToAddr(21), numToAddr(22)
		candidates := []common.Address{C1, C2, C3}

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		randao := mock_randao.NewMockRandaoModule(ctrl)
		v := newCN(t, withValset(valset), withRandao(randao), withGenesis()).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				5: makeHeaderWithVRank(5, 0, nil),                          // epoch start, empty
				6: makeHeaderWithVRank(6, 0, nil),                          // empty report
				7: makeHeaderWithVRank(7, 0, []common.Address{C1, C2, C3}), // P3 reports all
				8: makeHeaderWithVRank(8, 0, []common.Address{C1, C2}),     // P4 reports C1,C2
				9: makeHeaderWithVRank(9, 0, []common.Address{C1, C2}),     // P4 reports C1,C2 again
			},
		}
		valset.EXPECT().GetCandTesting(uint64(5)).Return(candidates, nil)
		valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(P1, nil)
		valset.EXPECT().GetProposer(uint64(6), uint64(0)).Return(P2, nil)
		valset.EXPECT().GetProposer(uint64(7), uint64(0)).Return(P3, nil)
		valset.EXPECT().GetProposer(uint64(8), uint64(0)).Return(P4, nil)
		valset.EXPECT().GetProposer(uint64(9), uint64(0)).Return(P4, nil)
		cfs, err := computeCFS(v, 5, 9)
		require.NoError(t, err)

		assert.Equal(t, uint64(1), cfs[C1], "C1 should be 1")
		assert.Equal(t, uint64(1), cfs[C2], "C2 should be 1")
		assert.Equal(t, uint64(0), cfs[C3], "C3 should be 0")
		_, hasC3 := cfs[C3]
		assert.True(t, hasC3, "C3 should appear in CFS map")
	})

	t.Run("KIP227_Example2", func(t *testing.T) {
		// Reproduces KIP-227 Example 2 through the chain-reading path.
		// Raw per-candidate/per-reporter totals are materialized as per-block cfReports.
		// For each reporter P_i, we generate max_c(raw[c][i]) blocks and include candidate C_j
		// in the first raw[j][i] blocks for that reporter.
		// With committee size 10, F = (10-1)/3 = 3.
		// Expected filtered CFS: C1=139, C2=289, C3=283, C4=221, C5=116.
		proposers := make([]common.Address, 10)
		for i := range proposers {
			proposers[i] = numToAddr(i)
		}
		candidates := make([]common.Address, 5)
		for i := range candidates {
			candidates[i] = numToAddr(10 + i)
		}
		// KIP-227 Example 2 matrix (rows=candidates C1..C5, cols=reporters P1..P10).
		cpMatrix := [][]uint64{
			{14, 12, 15, 34, 12, 32, 20, 8640, 8637, 8634}, // C1
			{48, 10, 59, 33, 49, 49, 41, 8640, 8637, 8634}, // C2
			{48, 22, 40, 41, 44, 27, 61, 8640, 8637, 8634}, // C3
			{50, 29, 45, 30, 23, 2, 42, 56, 56, 64},        // C4
			{71, 34, 62, 5, 11, 20, 18, 30, 19, 13},        // C5
		}

		start := uint64(1000)
		headers, blockReporter, end, totalBlocks := generateHeadersFromCPMatrix(t, start, candidates, proposers, cpMatrix)

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		randao := mock_randao.NewMockRandaoModule(ctrl)
		v := newCN(t, withValset(valset), withRandao(randao), withGenesis()).VRankModule
		v.Chain = &testChain{headers: headers}

		valset.EXPECT().GetCandTesting(start).Return(candidates, nil).Times(1)
		valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(
			func(blockNum, round uint64) (common.Address, error) {
				if round != 0 {
					return common.Address{}, assert.AnError
				}
				reporter, ok := blockReporter[blockNum]
				if !ok {
					return common.Address{}, assert.AnError
				}
				return reporter, nil
			},
		).Times(int(totalBlocks))

		cfs, err := computeCFS(v, start, end)
		require.NoError(t, err)

		assert.Equal(t, uint64(139), cfs[candidates[0]], "C1")
		assert.Equal(t, uint64(289), cfs[candidates[1]], "C2")
		assert.Equal(t, uint64(283), cfs[candidates[2]], "C3")
		assert.Equal(t, uint64(221), cfs[candidates[3]], "C4")
		assert.Equal(t, uint64(116), cfs[candidates[4]], "C5")
	})

	t.Run("epoch_start", func(t *testing.T) {
		// Verifies that a range containing only the epoch-start block produces a zero-score
		// map pre-seeded with all candidates (no failures yet).
		C1 := numToAddr(10)
		candidates := []common.Address{C1}

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		randao := mock_randao.NewMockRandaoModule(ctrl)
		v := newCN(t, withValset(valset), withRandao(randao), withGenesis()).VRankModule

		// Only one block in the range: the epoch-start block with nil VRank.
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				0: makeHeaderWithRound(0, 0), // VRank is nil
			},
		}
		valset.EXPECT().GetCandTesting(uint64(0)).Return(candidates, nil)
		valset.EXPECT().GetProposer(uint64(0), uint64(0)).Return(numToAddr(1), nil)

		cfs, err := computeCFS(v, 0, 0)
		require.NoError(t, err)
		// C1 is pre-seeded from GetCandTesting but no block reported any failure, so CFS=0.
		score, hasC1 := cfs[C1]
		assert.True(t, hasC1, "C1 should appear in CFS map (pre-seeded from candidates)")
		assert.Equal(t, uint64(0), score, "C1 CFS should be 0 (no failures reported)")
	})
}

// TestByzantineFilter verifies the byzantine-reporter filtering logic against KIP-227 spec examples.
func TestByzantineFilter(t *testing.T) {
	t.Run("KIP227_Example1", func(t *testing.T) {
		// From the spec: epoch=5, len(validator)=4, len(candidates)=3, F=1
		//
		//   Candidate \ Reporter | P1 | P2 | P3 | P4
		//   C1                   |  0 |  0 |  1 |  2  → Filtered = 1  (drop P4)
		//   C2                   |  0 |  0 |  1 |  2  → Filtered = 1  (drop P4)
		//   C3                   |  0 |  0 |  1 |  0  → Filtered = 0  (drop P3)
		P3, P4 := numToAddr(12), numToAddr(13)
		C1, C2, C3 := numToAddr(20), numToAddr(21), numToAddr(22)

		failuresByCandidate := vrank.CPMatrix{
			C1: {P3: 1, P4: 2},
			C2: {P3: 1, P4: 2},
			C3: {P3: 1},
		}

		cfs := byzantineFilter(failuresByCandidate, 1 /* F = (4-1)/3 = 1 */)

		assert.Equal(t, uint64(1), cfs[C1], "C1 CFS should be 1")
		assert.Equal(t, uint64(1), cfs[C2], "C2 CFS should be 1")
		score, hasC3 := cfs[C3]
		assert.True(t, hasC3, "C3 should appear in CFS map")
		assert.Equal(t, uint64(0), score, "C3 CFS should be 0 (filtered out)")
	})

	t.Run("KIP227_Example2", func(t *testing.T) {
		// From the spec: n=10 validators, 5 candidates (C1–C5), F=(10-1)/3=3.
		// P8, P9, P10 are Byzantine (report abnormally high counts for C1–C3).
		// Expected filtered CFS: C1=139, C2=289, C3=283, C4=221, C5=116.

		// Proposers P1–P10 indexed as numToAddr(0)–numToAddr(9)
		proposers := make([]common.Address, 10)
		for i := range proposers {
			proposers[i] = numToAddr(i)
		}
		// Candidates C1–C5 indexed as numToAddr(10)–numToAddr(14)
		candidates := make([]common.Address, 5)
		for i := range candidates {
			candidates[i] = numToAddr(10 + i)
		}

		// Raw data from KIP-227 Table (rows = candidates, columns = proposers):
		// Columns: P1  P2   P3   P4   P5   P6   P7   P8    P9    P10
		raw := [5][10]uint64{
			{14, 12, 15, 34, 12, 32, 20, 8640, 8637, 8634}, // C1
			{48, 10, 59, 33, 49, 49, 41, 8640, 8637, 8634}, // C2
			{48, 22, 40, 41, 44, 27, 61, 8640, 8637, 8634}, // C3
			{50, 29, 45, 30, 23, 2, 42, 56, 56, 64},        // C4
			{8640, 8637, 8634, 5, 11, 20, 18, 30, 19, 13},  // C5
		}

		failuresByCandidate := make(vrank.CPMatrix)
		for ci, c := range candidates {
			for pi, p := range proposers {
				if raw[ci][pi] > 0 {
					if failuresByCandidate[c] == nil {
						failuresByCandidate[c] = make(map[common.Address]uint64)
					}
					failuresByCandidate[c][p] = raw[ci][pi]
				}
			}
		}

		F := (10 - 1) / 3 // = 3
		cfs := byzantineFilter(failuresByCandidate, F)

		assert.Equal(t, uint64(139), cfs[candidates[0]], "C1")
		assert.Equal(t, uint64(289), cfs[candidates[1]], "C2")
		assert.Equal(t, uint64(283), cfs[candidates[2]], "C3")
		assert.Equal(t, uint64(221), cfs[candidates[3]], "C4")
		assert.Equal(t, uint64(116), cfs[candidates[4]], "C5")
	})
}
