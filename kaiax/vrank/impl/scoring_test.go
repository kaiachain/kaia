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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// addrN returns a deterministic address for test index n.
func addrN(n int) common.Address {
	return common.BigToAddress(big.NewInt(int64(n + 1)))
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

// -------------------------------------------------------------------
// TestGetPFS
// -------------------------------------------------------------------

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

func TestGetPFS_ErrFutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0) // current head = 5
	}
	v.Chain = &testChain{headers: headers}
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()

	_, err := v.GetPFS(5)
	assert.NoError(t, err)

	_, err = v.GetPFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)
}

// -------------------------------------------------------------------
// TestGetCFS
// -------------------------------------------------------------------

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

		valset.EXPECT().GetCandidates(epochStart).Return([]common.Address{C1}, nil).Times(2)
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

		valset.EXPECT().GetCandidates(epochStart).Return(candidates, nil)
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

func TestGetCFS_ErrFutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0) // current head = 5
	}
	v.Chain = &testChain{headers: headers}

	valset.EXPECT().GetCandidates(uint64(0)).Return(nil, nil)
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{addrN(0)}, nil)

	_, err := v.GetCFS(5)
	assert.NoError(t, err)

	_, err = v.GetCFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)
}

// -------------------------------------------------------------------
// TestComputePFS
// -------------------------------------------------------------------

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

	pfs, err := v.computePFS(10, 12)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), pfs[P0])
	assert.Equal(t, uint64(1), pfs[P1])
	assert.Len(t, pfs, 2)
}

// -------------------------------------------------------------------
// TestComputeCFS
// -------------------------------------------------------------------

func TestComputeCFS_KIP227Example1(t *testing.T) {
	// This reproduces KIP-227 Example 1 through the chain-reading path.
	// Blocks 5–9 are the epoch range (matching the example's block numbers).
	//
	//   proposer(5)=P1, cfReport(5)=[]        (epoch start: must be empty per spec)
	//   proposer(6)=P2, cfReport(6)=[]
	//   proposer(7)=P3, cfReport(7)=[C1,C2,C3]
	//   proposer(8)=P4, cfReport(8)=[C1,C2]
	//   proposer(9)=P4, cfReport(9)=[C1,C2]
	//
	// council size = 4 → F = 1
	// Expected CFS: C1=1, C2=1, C3=0 (in map, zero-scored)
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
	valset.EXPECT().GetCandidates(uint64(5)).Return(candidates, nil)
	valset.EXPECT().GetProposer(uint64(5), uint64(0)).Return(P1, nil)
	valset.EXPECT().GetProposer(uint64(6), uint64(0)).Return(P2, nil)
	valset.EXPECT().GetProposer(uint64(7), uint64(0)).Return(P3, nil)
	valset.EXPECT().GetProposer(uint64(8), uint64(0)).Return(P4, nil)
	valset.EXPECT().GetProposer(uint64(9), uint64(0)).Return(P4, nil)
	valset.EXPECT().GetCommittee(uint64(5), uint64(0)).Return(council, nil)

	cfs, err := v.computeCFS(5, 9)
	require.NoError(t, err)

	assert.Equal(t, uint64(1), cfs[C1], "C1 should be 1")
	assert.Equal(t, uint64(1), cfs[C2], "C2 should be 1")
	assert.Equal(t, uint64(0), cfs[C3], "C3 should be 0")
	_, hasC3 := cfs[C3]
	assert.True(t, hasC3, "C3 should appear in CFS map")
}

func TestComputeCFS_KIP227Example2(t *testing.T) {
	// This reproduces KIP-227 Example 2 through the chain-reading path.
	// Raw per-candidate/per-reporter totals are materialized as per-block cfReports.
	// For each reporter P_i, we generate max_c(raw[c][i]) blocks and include candidate C_j
	// in the first raw[j][i] blocks for that reporter.
	// With committee size 10, F = (10-1)/3 = 3.
	// Expected filtered CFS:
	//   C1=139, C2=289, C3=283, C4=221, C5=116.
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

	cfs, err := v.computeCFS(start, end)
	require.NoError(t, err)

	assert.Equal(t, uint64(139), cfs[candidates[0]], "C1")
	assert.Equal(t, uint64(289), cfs[candidates[1]], "C2")
	assert.Equal(t, uint64(283), cfs[candidates[2]], "C3")
	assert.Equal(t, uint64(221), cfs[candidates[3]], "C4")
	assert.Equal(t, uint64(116), cfs[candidates[4]], "C5")
}

func TestComputeCFS_EpochStartEmpty(t *testing.T) {
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
	valset.EXPECT().GetCandidates(uint64(0)).Return(candidates, nil)
	valset.EXPECT().GetProposer(uint64(0), uint64(0)).Return(P1, nil)
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return(council, nil)

	cfs, err := v.computeCFS(0, 0)
	require.NoError(t, err)
	// C1 is pre-seeded from GetCandidates but no block reported any failure, so CFS=0.
	score, hasC1 := cfs[C1]
	assert.True(t, hasC1, "C1 should appear in CFS map (pre-seeded from candidates)")
	assert.Equal(t, uint64(0), score, "C1 CFS should be 0 (no failures reported)")
}

// -------------------------------------------------------------------
// TestByzantineFilter – KIP-227 Example 1
// -------------------------------------------------------------------
//
// From the spec:
//   epoch = 5, len(validator) = 4, len(candidates) = 3, F = 1
//
//   Candidate \ Reporter | P1 | P2 | P3 | P4
//   C1                   |  0 |  0 |  1 |  2  → Filtered = 1  (drop P4)
//   C2                   |  0 |  0 |  1 |  2  → Filtered = 1  (drop P4)
//   C3                   |  0 |  0 |  1 |  0  → Filtered = 0  (drop P3)

func TestByzantineFilter_KIP227Example1(t *testing.T) {
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
}

// -------------------------------------------------------------------
// TestByzantineFilter – KIP-227 Example 2
// -------------------------------------------------------------------
//
// From the spec:
//   n = 10 validators, 5 candidates (C1–C5), F = (10-1)/3 = 3
//   P8, P9, P10 are Byzantine (report abnormally high counts for C1–C3).
//
//   Expected filtered CFS:
//     C1 = 139, C2 = 289, C3 = 283, C4 = 221, C5 = 116

func TestByzantineFilter_KIP227Example2(t *testing.T) {
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
}
