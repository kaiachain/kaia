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
	P0 := addrN(0)
	P1 := addrN(1)

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

// -------------------------------------------------------------------
// TestComputeCFS – end-to-end over a chain range
// -------------------------------------------------------------------
//
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

func TestComputeCFS_ChainRange(t *testing.T) {
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

// -------------------------------------------------------------------
// TestComputeCFS_EpochStartEmpty
// -------------------------------------------------------------------
//
// cfReport at epoch-start block must be nil/empty.  The scoring engine should
// accept a nil VRank field without error and count nothing for that block.

func TestComputeCFS_EpochStartEmpty(t *testing.T) {
	P1 := addrN(0)
	C1 := addrN(10)
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
// TestGetPFS_EpochBoundary
// -------------------------------------------------------------------
//
// GetPFS(N) covers [epochBegin(N), N]. Calling with the last block of epoch 1
// (2*vrankEpoch-1) should cover [vrankEpoch, 2*vrankEpoch-1].
// Block epochStart has round=1, so pfReport = [P0].

func TestGetPFS_EpochBoundary(t *testing.T) {
	P0 := addrN(0)
	epochStart := uint64(vrankEpoch)
	epochEnd := uint64(2*vrankEpoch - 1)

	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule

	headers := make(map[uint64]*types.Header, vrankEpoch)
	for i := uint64(0); i < vrankEpoch; i++ {
		headers[epochStart+i] = makeHeaderWithRound(epochStart+i, 0)
	}
	headers[epochStart] = makeHeaderWithRound(epochStart, 1) // round 1 → pfReport = [P0]
	v.Chain = &testChain{headers: headers}

	valset.EXPECT().GetProposer(epochStart, uint64(0)).Return(P0, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(1), nil).AnyTimes()

	pfs, err := v.GetPFS(epochEnd)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), pfs[P0], "P0 should have 1 proposal failure")
}

// -------------------------------------------------------------------
// TestGetCFS_EpochBoundary
// -------------------------------------------------------------------
//
// GetCFS(N) covers [epochBegin(N), N]. Calling with the last block of epoch 1
// (2*vrankEpoch-1) should cover [vrankEpoch, 2*vrankEpoch-1].
// Block epochStart+1 is proposed by P1 and reports C1; F=0, so CFS[C1]=1.

func TestGetCFS_EpochBoundary(t *testing.T) {
	P1 := addrN(1)
	C1 := addrN(10)
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

	valset.EXPECT().GetCandidates(epochStart).Return([]common.Address{C1}, nil)
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(P1, nil).AnyTimes()
	valset.EXPECT().GetCommittee(epochStart, uint64(0)).Return([]common.Address{addrN(0), P1}, nil)

	cfs, err := v.GetCFS(epochEnd)
	require.NoError(t, err)
	// F = (2-1)/3 = 0, so no filtering; C1 appears once.
	assert.Equal(t, uint64(1), cfs[C1], "C1 should have 1 candidate failure")
}

// -------------------------------------------------------------------
// TestGetPFS_FutureBlock / TestGetCFS_FutureBlock
// -------------------------------------------------------------------

func TestGetPFS_FutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0) // current head = 5
	}
	v.Chain = &testChain{headers: headers}

	_, err := v.GetPFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)

	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()
	_, err = v.GetPFS(5)
	assert.NoError(t, err)
}

func TestGetCFS_FutureBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	v := createCN(t, valset).VRankModule

	headers := make(map[uint64]*types.Header, 6)
	for i := uint64(0); i <= 5; i++ {
		headers[i] = makeHeaderWithRound(i, 0) // current head = 5
	}
	v.Chain = &testChain{headers: headers}

	_, err := v.GetCFS(6)
	assert.ErrorIs(t, err, vrank.ErrFutureBlock)

	valset.EXPECT().GetCandidates(uint64(0)).Return(nil, nil)
	valset.EXPECT().GetProposer(gomock.Any(), uint64(0)).Return(addrN(0), nil).AnyTimes()
	valset.EXPECT().GetCommittee(uint64(0), uint64(0)).Return([]common.Address{addrN(0)}, nil)
	_, err = v.GetCFS(5)
	assert.NoError(t, err)
}
