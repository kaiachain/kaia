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

// addrN returns a deterministic address for test index n.
func addrN(n int) common.Address {
	return common.BigToAddress(big.NewInt(int64(n + 1)))
}

// newTestModuleWithHeaders wraps newTestModule for scoring/cache tests that need a minimal module
// with injectable chain state. It intentionally initializes the module with an empty chain so that
// catchUpScoreCaches exits early (keeping caches cold), then swaps in the real headers after Init.
// Tests that need pre-warmed caches should seed them manually via pfsCache.Add / cpMatrixCache.Add.
func newTestModuleWithHeaders(t *testing.T, valset *mock_valset.MockValsetModule, db database.Database, headers map[uint64]*types.Header) *VRankModule {
	t.Helper()

	_, module := newTestModule(t, valset, db, &testChain{headers: map[uint64]*types.Header{}})
	module.Chain = &testChain{headers: headers}
	return module
}

func computeCFS(v *VRankModule, start, end uint64) (map[common.Address]uint64, error) {
	cpMatrix, err := v.newCPMatrix(start)
	if err != nil {
		return nil, err
	}
	cpMatrix, err = v.computeCPMatrix(start, end, cpMatrix)
	if err != nil {
		return nil, err
	}
	return v.generateCFSFromCPMatrix(start, cpMatrix)
}

// makeHeaderWithVRank creates a header with a specific round and an encoded cfReport in VRank.
// cfAddrs may be nil/empty for an empty report.
func makeHeaderWithVRank(number uint64, round int64, cfAddrs []common.Address) *types.Header {
	h := makeHeaderWithRound(number, round)
	encoded, err := vrank.EncodeReport(vrank.Report(cfAddrs))
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

// TestGetPFS verifies incremental PFS accumulation over multiple blocks within one epoch.
// Each call seeds from the previous result via the nearby-cache path.
func TestGetPFS(t *testing.T) {
	var (
		P0, P1, P2    = addrN(0), addrN(1), addrN(2)
		proposerList  = []common.Address{P0, P1, P2}
		epochStart    = uint64(vrankEpoch)
		proposerCount = uint64(len(proposerList))
		round1Offset  = proposerCount     // keep (epochStart+round1Offset)%len(proposerList) == 0
		round2Offset  = proposerCount * 2 // keep (epochStart+round2Offset)%len(proposerList) == 0
		headers       = make(map[uint64]*types.Header, round2Offset+1)

		ctrl   = gomock.NewController(t)
		valset = mock_valset.NewMockValsetModule(ctrl)
		v      = createCN(t, valset).VRankModule
	)

	for i := uint64(0); i <= round2Offset; i++ {
		headers[epochStart+i] = makeHeaderWithRound(epochStart+i, 0)
	}
	// GetProposer(blockNum, round) = proposerList[blockNum%len(proposerList)+round].
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(
		func(blockNum, round uint64) (common.Address, error) {
			idx := int((blockNum % uint64(len(proposerList))) + round)
			if idx >= len(proposerList) {
				return common.Address{}, assert.AnError
			}
			return proposerList[idx], nil
		},
	).Times(4)

	// PfReport(epochStart+round1Offset)=[P0]
	headers[epochStart+round1Offset] = makeHeaderWithRound(epochStart+round1Offset, 1)
	// PfReport(epochStart+round2Offset)=[P0, P1]
	headers[epochStart+round2Offset] = makeHeaderWithRound(epochStart+round2Offset, 2)
	v.Chain = &testChain{headers: headers}

	pfs, err := v.GetPFS(epochStart)
	require.NoError(t, err)
	assert.Len(t, pfs, 0)

	pfs, err = v.GetPFS(epochStart + round1Offset)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs[P0])
	assert.Len(t, pfs, 1)

	pfs, err = v.GetPFS(epochStart + round2Offset)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), pfs[P0])
	assert.Equal(t, uint64(1), pfs[P1])
	assert.NotContains(t, pfs, P2)
	assert.Len(t, pfs, 2)
}

// TestGetPFS_ErrFutureBlock verifies that GetPFS returns ErrFutureBlock for a block beyond the chain head.
func TestGetPFS_ErrFutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()

	_, err := v.GetPFS(5)
	assert.NoError(t, err)

	_, err = v.GetPFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)
}

// TestGetPFS_CacheHit verifies that a second GetPFS call for the same block returns the cached
// result, and that the returned map is a defensive copy (mutations do not affect the cache).
func TestGetPFS_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	headers := map[uint64]*types.Header{}
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[5] = makeHeaderWithRound(5, 1)

	v := newTestModuleWithHeaders(t, valset, db, headers)
	p0 := addrN(0)
	valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(p0, nil).Times(1)

	pfs1, err := v.GetPFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs1[p0])

	// Mutating the returned map must not affect the cached copy.
	pfs1[p0] = 999
	pfs2, err := v.GetPFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs2[p0])
}

// TestGetPFS_NearbyCacheHit verifies that GetPFS for block N seeds from the cached block N-1
// rather than recomputing from epoch start.
func TestGetPFS_NearbyCacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	headers := map[uint64]*types.Header{}
	for i := uint64(0); i <= 6; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[5] = makeHeaderWithRound(5, 1)
	headers[6] = makeHeaderWithRound(6, 1)

	v := newTestModuleWithHeaders(t, valset, db, headers)
	p0, p1 := addrN(0), addrN(1)
	valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(p0, nil).Times(1)
	valset.EXPECT().GetProposer(uint64(6), uint64(0)).Return(p1, nil).Times(1)

	pfs5, err := v.GetPFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs5[p0])
	_, ok := v.pfsCache.Get(uint64(5))
	assert.True(t, ok)

	pfs6, err := v.GetPFS(6)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs6[p0])
	assert.Equal(t, uint64(1), pfs6[p1])
}

// TestGetPFS_DBCheckpointHit verifies that when there is no in-memory cache but a DB checkpoint
// exists, GetPFS resumes from the checkpoint instead of epoch start.
// The checkpoint must be at a scoreCheckpointInterval multiple for the O(1) lookup to find it.
func TestGetPFS_DBCheckpointHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P0 := addrN(0)
	cp := scoreCheckpointInterval
	headers := map[uint64]*types.Header{
		cp:     makeHeaderWithRound(cp, 0),
		cp + 1: makeHeaderWithRound(cp+1, 1), // round 1 → pfReport = [P0]
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	WriteCheckpoint(db, cp, map[common.Address]uint64{}, map[common.Address]map[common.Address]uint64{})
	WriteLastCheckpoint(db, cp)

	// With the DB checkpoint, only block cp+1 needs to be computed.
	valset.EXPECT().GetProposer(cp+1, uint64(0)).Return(P0, nil).Times(1)

	pfs, err := v.GetPFS(cp + 1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs[P0])
}

// TestGetPFS_EpochScan verifies that with no in-memory cache and no DB checkpoint, GetPFS
// correctly scans the full epoch from epoch start.
func TestGetPFS_EpochScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P0 := addrN(0)
	headers := map[uint64]*types.Header{
		0: makeHeaderWithRound(0, 0),
		1: makeHeaderWithRound(1, 0),
		2: makeHeaderWithRound(2, 1), // round 1 → pfReport = [P0]
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	_, inCache := v.pfsCache.Get(uint64(2))
	require.False(t, inCache, "cache must be cold before first call")

	valset.EXPECT().GetProposer(uint64(2), uint64(0)).Return(P0, nil).Times(1)

	pfs, err := v.GetPFS(2)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs[P0])
}

// TestGetPFS_EpochStart verifies that GetPFS at exactly the first block of an epoch
// always returns an empty map, even when the previous epoch's cache is populated.
func TestGetPFS_EpochStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	epochStart := uint64(vrankEpoch)
	v := newTestModuleWithHeaders(t, valset, db, map[uint64]*types.Header{
		epochStart: makeHeaderWithRound(epochStart, 0),
	})

	// Seed the previous epoch's last block with non-zero PFS.
	v.pfsCache.Add(epochStart-1, map[common.Address]uint64{addrN(0): 99})

	// Probe clamps to i<=0, so only epochStart itself is checked (no cache there).
	// DB checkpoint: lastCP=0 → ignored. Computes from epochStart: round=0 → no pfReport.
	pfs, err := v.GetPFS(epochStart)
	require.NoError(t, err)
	assert.Empty(t, pfs, "epoch start must reset PFS to empty regardless of prior epoch cache")
}

// TestGetPFS_EpochBoundaryClamp verifies that a nearby cache entry from the previous epoch
// is NOT used as a seed for the current epoch.
func TestGetPFS_EpochBoundaryClamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	epochStart := uint64(vrankEpoch)
	blockNum := epochStart + 5
	headers := map[uint64]*types.Header{}
	for i := epochStart; i <= blockNum; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed the last block of the previous epoch with a large PFS value.
	v.pfsCache.Add(epochStart-1, map[common.Address]uint64{addrN(0): 99})

	// The probe limit is min(64, blockNum-epochStart)=5, so epochStart-1 (distance=6) is
	// never reached. GetProposer must not be called because all blocks in [epochStart, blockNum]
	// have round=0.
	pfs, err := v.GetPFS(blockNum)
	require.NoError(t, err)
	assert.Empty(t, pfs, "previous epoch cache must not carry over into current epoch")
}

// TestGetPFS_NearbyProbe_ExactBoundary verifies that a cache entry exactly
// scoreCacheProbeLookback blocks in the past IS used as a seed.
func TestGetPFS_NearbyProbe_ExactBoundary(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	const blockNum = uint64(100)
	const boundary = blockNum - scoreCacheProbeLookback // = 36, distance exactly 64

	headers := map[uint64]*types.Header{}
	for i := boundary; i <= blockNum; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[blockNum] = makeHeaderWithRound(blockNum, 1) // one failure

	P0 := addrN(0)
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed cache at exactly the lookback boundary with one prior failure.
	v.pfsCache.Add(boundary, map[common.Address]uint64{P0: 1})

	// Only blockNum needs to be computed (blocks boundary+1..blockNum-1 are round=0).
	valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(P0, nil).Times(1)

	pfs, err := v.GetPFS(blockNum)
	require.NoError(t, err)
	// Seed contributes 1; blockNum adds 1 more.
	assert.Equal(t, uint64(2), pfs[P0], "cache at exact lookback boundary must be used as seed")
}

// TestGetPFS_NearbyProbe_BeyondBoundary verifies that a cache entry one block beyond
// scoreCacheProbeLookback is NOT found, so computation restarts from epoch start.
func TestGetPFS_NearbyProbe_BeyondBoundary(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	const blockNum = uint64(100)
	const beyond = blockNum - scoreCacheProbeLookback - 1 // = 35, distance 65

	headers := map[uint64]*types.Header{}
	for i := uint64(0); i <= blockNum; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[blockNum] = makeHeaderWithRound(blockNum, 1)

	P0 := addrN(0)
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed at distance-65 (outside the window) with a large value so that if the probe
	// incorrectly reaches it the result would be wrong.
	v.pfsCache.Add(beyond, map[common.Address]uint64{P0: 99})

	// Full recomputation from epoch start (0): only blockNum has round>0.
	valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(P0, nil).Times(1)

	pfs, err := v.GetPFS(blockNum)
	require.NoError(t, err)
	// Only blockNum contributed (the beyond-window seed must not be used).
	assert.Equal(t, uint64(1), pfs[P0], "cache beyond lookback boundary must not be used as seed")
}

// TestGetPFS_GetProposerError verifies that an error from GetProposer propagates out of GetPFS.
func TestGetPFS_GetProposerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	headers := map[uint64]*types.Header{
		0: makeHeaderWithRound(0, 0),
		1: makeHeaderWithRound(1, 1), // round=1 → needs GetProposer(1, 0)
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(common.Address{}, assert.AnError)

	_, err := v.GetPFS(1)
	assert.ErrorIs(t, err, assert.AnError)
}

// TestGetPFS_MissingHeader verifies that a missing block header propagates ErrHeaderNotFound.
func TestGetPFS_MissingHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	// Header exists for block 1 (so CurrentHeader()=1 and ErrFutureBlock is not triggered)
	// but NOT for block 0, which computePFS will try to fetch first.
	headers := map[uint64]*types.Header{
		1: makeHeaderWithRound(1, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	_, err := v.GetPFS(1)
	assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
}

// TestGetCFS verifies CFS computation across epoch boundaries and with multiple reporters.
func TestGetCFS(t *testing.T) {
	t.Run("epoch boundary", func(t *testing.T) {
		P1, C1 := addrN(1), addrN(10)
		epochStart := uint64(vrankEpoch)
		epochEnd := uint64(2*vrankEpoch - 1)

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		v := createCN(t, valset).VRankModule

		headers := make(map[uint64]*types.Header, vrankEpoch)
		for i := uint64(0); i < vrankEpoch; i++ {
			headers[epochStart+i] = makeHeaderWithRound(epochStart+i, 0)
		}
		headers[epochStart+1] = makeHeaderWithVRank(epochStart+1, 0, []common.Address{C1})
		v.Chain = &testChain{headers: headers}

		valset.EXPECT().GetCandidates(gomock.Any()).Return([]common.Address{C1}, nil).AnyTimes()
		valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(P1, nil).AnyTimes()
		valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return([]common.Address{addrN(0), P1}, nil).Times(2)

		cfs, err := v.GetCFS(epochStart)
		require.NoError(t, err)
		assert.Equal(t, uint64(0), cfs[C1], "C1 should have 0 candidate failures at epoch start")

		cfs, err = v.GetCFS(epochEnd)
		require.NoError(t, err)
		// F = (2-1)/3 = 0, so no filtering; C1 appears once.
		assert.Equal(t, uint64(1), cfs[C1], "C1 should have 1 candidate failure")
	})

	t.Run("multi-reporter epoch", func(t *testing.T) {
		epochStart := uint64(vrankEpoch)
		P1, P2, P3, P4 := addrN(30), addrN(31), addrN(32), addrN(33)
		C1, C2, C3 := addrN(40), addrN(41), addrN(42)
		candidates := []common.Address{C1, C2, C3}
		committee := []common.Address{P1, P2, P3, P4}

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		v := createCN(t, valset).VRankModule
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

		valset.EXPECT().GetCandidates(gomock.Any()).Return(candidates, nil).AnyTimes()
		valset.EXPECT().GetProposer(epochStart, uint64(0)).Return(P1, nil)
		valset.EXPECT().GetProposer(epochStart+1, uint64(0)).Return(P1, nil)
		valset.EXPECT().GetProposer(epochStart+2, uint64(1)).Return(P2, nil)
		valset.EXPECT().GetProposer(epochStart+3, uint64(0)).Return(P3, nil)
		valset.EXPECT().GetProposer(epochStart+4, uint64(0)).Return(P4, nil)
		valset.EXPECT().GetProposer(epochStart+5, uint64(0)).Return(P4, nil)
		valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return(committee, nil)

		cfs, err := v.GetCFS(epochStart + 5)
		require.NoError(t, err)
		// F = 1. Raw totals:
		// C1: {P1:1, P2:1, P3:1, P4:2} -> drop 2 -> score 3
		// C2: {P1:1, P3:1, P4:2}       -> drop 2 -> score 2
		// C3: {P3:1}                   -> drop 1 -> score 0
		assert.Equal(t, uint64(3), cfs[C1])
		assert.Equal(t, uint64(2), cfs[C2])
		assert.Equal(t, uint64(0), cfs[C3])
		assert.Len(t, cfs, 3)
	})
}

// TestGetCFS_ErrFutureBlock verifies that GetCFS returns ErrFutureBlock for a block beyond the chain head.
func TestGetCFS_ErrFutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	valset.EXPECT().GetCandidates(gomock.Any()).Return(nil, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{addrN(0)}, nil)

	_, err := v.GetCFS(5)
	assert.NoError(t, err)

	_, err = v.GetCFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)
}

// TestGetCFS_CacheHit verifies that a second GetCFS call for the same block returns the cached
// result, and that the returned map is a defensive copy (mutations do not affect the cache).
func TestGetCFS_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P1, C1 := addrN(1), addrN(10)
	headers := map[uint64]*types.Header{}
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[5] = makeHeaderWithVRank(5, 0, []common.Address{C1})

	v := newTestModuleWithHeaders(t, valset, db, headers)
	// GetProposer is only called for blocks with non-empty cfReports (block 5 has VRank=[C1]).
	valset.EXPECT().GetCandidates(uint64(0)).Return([]common.Address{C1}, nil).Times(1)
	valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(P1, nil).Times(1)
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{P1}, nil).Times(2)

	cfs1, err := v.GetCFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), cfs1[C1])

	// Mutating the returned map must not affect the cached copy.
	cfs1[C1] = 999
	cfs2, err := v.GetCFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), cfs2[C1])
}

// TestGetCFS_NearbyCacheHit verifies that GetCFS for block N seeds from the cached block N-1
// rather than recomputing from epoch start.
func TestGetCFS_NearbyCacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P1, P2, C1 := addrN(1), addrN(2), addrN(10)
	headers := map[uint64]*types.Header{}
	for i := uint64(0); i <= 6; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	headers[5] = makeHeaderWithVRank(5, 0, []common.Address{C1})
	headers[6] = makeHeaderWithVRank(6, 0, []common.Address{C1})

	v := newTestModuleWithHeaders(t, valset, db, headers)
	valset.EXPECT().GetCandidates(gomock.Any()).Return([]common.Address{C1}, nil).AnyTimes()
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{P1, P2}, nil).Times(2)
	// GetProposer only called for blocks with non-empty cfReports (blocks 5 and 6).
	valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(P1, nil).Times(1)
	valset.EXPECT().GetProposer(uint64(6), uint64(0)).Return(P2, nil).Times(1)

	cfs5, err := v.GetCFS(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), cfs5[C1])
	_, ok := v.cpMatrixCache.Get(uint64(5))
	assert.True(t, ok)

	cfs6, err := v.GetCFS(6)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), cfs6[C1])
}

// TestGetCFS_DBCheckpointHit verifies that when there is no in-memory cache but a DB checkpoint
// exists, GetCFS resumes from the checkpoint instead of epoch start.
// The checkpoint must be at a scoreCheckpointInterval multiple for the O(1) lookup to find it.
func TestGetCFS_DBCheckpointHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P1, C1 := addrN(1), addrN(10)
	cp := scoreCheckpointInterval
	headers := map[uint64]*types.Header{
		cp:     makeHeaderWithRound(cp, 0),
		cp + 1: makeHeaderWithVRank(cp+1, 0, []common.Address{C1}),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Pre-write a DB checkpoint at the first checkpoint block with C1 already having 1 failure.
	WriteCheckpoint(db, cp, map[common.Address]uint64{}, map[common.Address]map[common.Address]uint64{
		C1: {P1: 1},
	})
	WriteLastCheckpoint(db, cp)

	// With the DB checkpoint, only block cp+1 needs to be computed.
	// GetCandidates (called in newCPMatrix for epoch start) must NOT be called.
	valset.EXPECT().GetProposer(cp+1, uint64(0)).Return(P1, nil).Times(1)
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{P1}, nil).Times(1)

	cfs, err := v.GetCFS(cp + 1)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), cfs[C1])
}

// TestGetCFS_EpochScan verifies that with no in-memory cache and no DB checkpoint, GetCFS
// correctly scans the full epoch from epoch start.
func TestGetCFS_EpochScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	P1, C1 := addrN(1), addrN(10)
	headers := map[uint64]*types.Header{
		0: makeHeaderWithRound(0, 0),
		1: makeHeaderWithVRank(1, 0, []common.Address{C1}),
		2: makeHeaderWithVRank(2, 0, []common.Address{C1}),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	_, inCache := v.cpMatrixCache.Get(uint64(2))
	require.False(t, inCache, "cache must be cold before first call")

	// GetProposer only called for blocks with non-empty cfReports (blocks 1 and 2).
	valset.EXPECT().GetCandidates(uint64(0)).Return([]common.Address{C1}, nil).Times(1)
	valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(P1, nil).Times(1)
	valset.EXPECT().GetProposer(uint64(2), uint64(0)).Return(P1, nil).Times(1)
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{P1}, nil).Times(1)

	cfs, err := v.GetCFS(2)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), cfs[C1])
}

// TestGetCFS_EpochStart verifies that GetCFS at the first block of an epoch returns a zero score
// for every candidate (pre-seeded from GetCandidates) and does NOT carry over state from the
// previous epoch.
func TestGetCFS_EpochStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	epochStart := uint64(vrankEpoch)
	C1, C2, P1 := addrN(10), addrN(11), addrN(1)
	headers := map[uint64]*types.Header{
		epochStart: makeHeaderWithVRank(epochStart, 0, nil), // epoch start: VRank must be nil
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed the previous epoch's last block with large CFS values.
	v.cpMatrixCache.Add(epochStart-1, map[common.Address]map[common.Address]uint64{
		C1: {P1: 50}, C2: {P1: 50},
	})

	// GetProposer is NOT called: the epoch-start block has a nil/empty cfReport (VRank=nil encoded).
	valset.EXPECT().GetCandidates(epochStart).Return([]common.Address{C1, C2}, nil).Times(1)
	valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return([]common.Address{P1}, nil).Times(1)

	cfs, err := v.GetCFS(epochStart)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), cfs[C1], "C1 must start at 0 in new epoch")
	assert.Equal(t, uint64(0), cfs[C2], "C2 must start at 0 in new epoch")
	assert.Len(t, cfs, 2, "all candidates must appear even with zero score")
}

// TestGetCFS_EpochBoundaryClamp verifies that a nearby cpMatrix cache entry from the
// previous epoch is NOT used as a seed for the current epoch.
func TestGetCFS_EpochBoundaryClamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	epochStart := uint64(vrankEpoch)
	blockNum := epochStart + 3
	C1, P1 := addrN(10), addrN(1)
	headers := map[uint64]*types.Header{}
	for i := epochStart; i <= blockNum; i++ {
		headers[i] = makeHeaderWithRound(i, 0)
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed the last block of the previous epoch with a large accumulated value.
	v.cpMatrixCache.Add(epochStart-1, map[common.Address]map[common.Address]uint64{
		C1: {P1: 99},
	})

	// Probe limit = min(64, blockNum-epochStart)=3; epochStart-1 is at distance 4, not reached.
	// Must call GetCandidates to start a fresh epoch (Times(1)).
	valset.EXPECT().GetCandidates(epochStart).Return([]common.Address{C1}, nil).Times(1)
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(P1, nil).AnyTimes()
	valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return([]common.Address{P1}, nil).Times(1)

	cfs, err := v.GetCFS(blockNum)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), cfs[C1], "previous epoch cpMatrix must not carry over into current epoch")
}

// TestGetCFS_NearbyProbe_SameEpoch verifies that the nearby-probe path uses a same-epoch cache
// entry and ignores a previous-epoch entry at the same distance.
func TestGetCFS_NearbyProbe_SameEpoch(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	epochStart := uint64(vrankEpoch)
	C1, P1 := addrN(10), addrN(1)

	headers := map[uint64]*types.Header{
		epochStart:     makeHeaderWithRound(epochStart, 0),
		epochStart + 1: makeHeaderWithVRank(epochStart+1, 0, []common.Address{C1}),
		epochStart + 2: makeHeaderWithRound(epochStart+2, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Seed cache at epochStart+1: C1 has 1 failure from that block.
	v.cpMatrixCache.Add(epochStart+1, map[common.Address]map[common.Address]uint64{
		C1: {P1: 1},
	})
	// Also seed the previous epoch to make sure it is NOT used.
	v.cpMatrixCache.Add(epochStart-1, map[common.Address]map[common.Address]uint64{
		C1: {P1: 99},
	})

	// Nearby hit at epochStart+1 (distance=1). GetCandidates must NOT be called.
	// epochStart+2 has no cfReport, so GetProposer is NOT called (cfReport checked first).
	valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return([]common.Address{P1}, nil).Times(1)

	cfs, err := v.GetCFS(epochStart + 2)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), cfs[C1], "score must come only from same-epoch nearby cache")
}

// TestGetCFS_GetProposerError verifies that an error from GetProposer propagates out of GetCFS.
func TestGetCFS_GetProposerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	C1 := addrN(10)
	headers := map[uint64]*types.Header{
		0: makeHeaderWithRound(0, 0),
		// Block 1 has a non-empty cfReport so computeCPMatrix will call GetProposer for it.
		1: makeHeaderWithVRank(1, 0, []common.Address{C1}),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	valset.EXPECT().GetCandidates(uint64(0)).Return([]common.Address{C1}, nil).Times(1)
	// GetProposer is NOT called for block 0 (empty cfReport); only for block 1.
	valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(common.Address{}, assert.AnError).Times(1)

	_, err := v.GetCFS(1)
	assert.ErrorIs(t, err, assert.AnError)
}

// TestGetCFS_MissingHeader verifies that a missing block header propagates ErrHeaderNotFound.
func TestGetCFS_MissingHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	db := database.NewMemDB()

	C1 := addrN(10)
	// Block 1 exists (so CurrentHeader()=1), but block 0 is missing.
	// computeCPMatrix will try GetHeaderByNumber(0) first and get nil.
	headers := map[uint64]*types.Header{
		1: makeHeaderWithRound(1, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	valset.EXPECT().GetCandidates(uint64(0)).Return([]common.Address{C1}, nil).Times(1)

	_, err := v.GetCFS(1)
	assert.ErrorIs(t, err, vrank.ErrHeaderNotFound)
}

// TestComputePFS verifies the raw proposer-failure accumulation over a small block range.
func TestComputePFS(t *testing.T) {
	// Set up three consecutive blocks:
	//   block 10, round 0: no pfReport (first proposer succeeded)
	//   block 11, round 2: pfReport = [P0, P1]  (rounds 0 and 1 failed)
	//   block 12, round 1: pfReport = [P0]       (round 0 failed)
	//
	// Expected PFS: {P0: 2, P1: 1}
	P0, P1 := addrN(0), addrN(1)

	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule
	v.Chain = &testChain{
		headers: map[uint64]*types.Header{
			10: makeHeaderWithRound(10, 0),
			11: makeHeaderWithRound(11, 2),
			12: makeHeaderWithRound(12, 1),
		},
	}

	// block 10, round 0: no GetProposer calls (loop [0,0) is empty)
	// block 11, round 2: proposers for rounds 0 and 1
	valset.EXPECT().GetProposer(uint64(11), uint64(0)).Return(P0, nil)
	valset.EXPECT().GetProposer(uint64(11), uint64(1)).Return(P1, nil)
	// block 12, round 1: proposer for round 0
	valset.EXPECT().GetProposer(uint64(12), uint64(0)).Return(P0, nil)

	pfs, err := v.computePFS(10, 12, make(map[common.Address]uint64))
	require.NoError(t, err)
	assert.Equal(t, uint64(2), pfs[P0])
	assert.Equal(t, uint64(1), pfs[P1])
	assert.Len(t, pfs, 2)
}

// TestComputeCFS verifies CFS computation via computeCPMatrix + generateCFSFromCPMatrix.
func TestComputeCFS(t *testing.T) {
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
		// council size = 4 → F = 1
		// Expected CFS: C1=1, C2=1, C3=0 (present in map with zero score)
		P1, P2, P3, P4 := addrN(10), addrN(11), addrN(12), addrN(13)
		C1, C2, C3 := addrN(20), addrN(21), addrN(22)
		council := []common.Address{P1, P2, P3, P4}
		candidates := []common.Address{C1, C2, C3}

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				5: makeHeaderWithVRank(5, 0, nil),                          // epoch start, empty
				6: makeHeaderWithVRank(6, 0, nil),                          // empty report
				7: makeHeaderWithVRank(7, 0, []common.Address{C1, C2, C3}), // P3 reports all
				8: makeHeaderWithVRank(8, 0, []common.Address{C1, C2}),     // P4 reports C1,C2
				9: makeHeaderWithVRank(9, 0, []common.Address{C1, C2}),     // P4 reports C1,C2 again
			},
		}
		// GetProposer only called for blocks with non-empty cfReports (blocks 7–9).
		// Blocks 5 and 6 have nil cfReports so GetProposer is skipped for them.
		valset.EXPECT().GetCandidates(uint64(5)).Return(candidates, nil)
		valset.EXPECT().GetProposer(uint64(7), uint64(0)).Return(P3, nil)
		valset.EXPECT().GetProposer(uint64(8), uint64(0)).Return(P4, nil)
		valset.EXPECT().GetProposer(uint64(9), uint64(0)).Return(P4, nil)
		valset.EXPECT().GetCommittee(uint64(5), uint64(0)).Return(council, nil)

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
			proposers[i] = addrN(i)
		}
		candidates := make([]common.Address, 5)
		for i := range candidates {
			candidates[i] = addrN(10 + i)
		}
		committee := append([]common.Address(nil), proposers...)

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
		v := createCN(t, valset).VRankModule
		v.Chain = &testChain{headers: headers}

		valset.EXPECT().GetCandidates(start).Return(candidates, nil).Times(1)
		valset.EXPECT().GetCommittee(start, uint64(0)).Return(committee, nil).Times(1)
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
		P1, C1 := addrN(0), addrN(10)
		council := []common.Address{P1}
		candidates := []common.Address{C1}

		ctrl := gomock.NewController(t)
		valset := mock_valset.NewMockValsetModule(ctrl)
		v := createCN(t, valset).VRankModule

		// Only one block in the range: the epoch-start block with nil VRank.
		v.Chain = &testChain{
			headers: map[uint64]*types.Header{
				0: makeHeaderWithRound(0, 0), // VRank is nil
			},
		}
		// GetProposer is NOT called: block 0 has no VRank set (nil → empty cfReport).
		valset.EXPECT().GetCandidates(uint64(0)).Return(candidates, nil)
		valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return(council, nil)

		cfs, err := computeCFS(v, 0, 0)
		require.NoError(t, err)
		// C1 is pre-seeded from GetCandidates but no block reported any failure, so CFS=0.
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
		P3, P4 := addrN(12), addrN(13)
		C1, C2, C3 := addrN(20), addrN(21), addrN(22)

		failuresByCandidate := map[common.Address]map[common.Address]uint64{
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

		// Proposers P1–P10 indexed as addrN(0)–addrN(9)
		proposers := make([]common.Address, 10)
		for i := range proposers {
			proposers[i] = addrN(i)
		}
		// Candidates C1–C5 indexed as addrN(10)–addrN(14)
		candidates := make([]common.Address, 5)
		for i := range candidates {
			candidates[i] = addrN(10 + i)
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

		failuresByCandidate := make(map[common.Address]map[common.Address]uint64)
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
