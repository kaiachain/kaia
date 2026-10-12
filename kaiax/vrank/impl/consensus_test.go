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
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeSelfReportHeader builds a round-0 header at `number` whose VRank is a cfReport listing
// cfAddrs as failed candidates. Empty cfAddrs yields a nil VRank field.
func makeSelfReportHeader(t *testing.T, number uint64, cfAddrs []common.Address) *types.Header {
	t.Helper()
	h := makeHeaderWithRound(number, 0)
	if len(cfAddrs) > 0 {
		encoded, err := vrank.EncodeVRank(vrank.VRankPayload{Report: cfAddrs})
		require.NoError(t, err)
		h.VRank = encoded
	}
	return h
}

func makeEpochStartVRankHeader(t *testing.T, number uint64, candidates []common.Address) *types.Header {
	t.Helper()
	encoded, err := vrank.EncodeVRank(vrank.VRankPayload{Report: candidates})
	require.NoError(t, err)
	return &types.Header{Number: big.NewInt(int64(number)), VRank: encoded}
}

// TestVerifyHeader covers all branches of VRankModule.VerifyHeader.
func TestVerifyHeader(t *testing.T) {
	C1, C2, C3 := numToAddr(1), numToAddr(2), numToAddr(3)
	candidates := []common.Address{C1, C2, C3}
	const num = uint64(100)

	// newVerifier wires GetCandTesting(any)=candidates; VerifyHeader validates only the failed list.
	newVerifier := func(t *testing.T) *VRankModule {
		return newCN(t, withCandidates(candidates), withoutStart()).VRankModule
	}

	t.Run("pre-fork: VRank must be absent", func(t *testing.T) {
		v := newCN(t, withHardfork("osaka"), withoutStart()).VRankModule
		h := &types.Header{Number: big.NewInt(100)}
		assert.NoError(t, v.VerifyHeader(h, nil))
		h = makeSelfReportHeader(t, 100, []common.Address{C1})
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrUnexpectedVRankBeforePermissionless)
		// A non-nil zero-length VRank is present and must be rejected.
		h = &types.Header{Number: big.NewInt(100), VRank: []byte{}}
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

	t.Run("epoch-start: empty CandTesting leaves VRank absent", func(t *testing.T) {
		v := newCN(t, withCandidates(nil)).VRankModule
		h := makeEpochStartVRankHeader(t, params.DefaultVRankEpoch, nil)
		assert.Nil(t, h.VRank)
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
		assert.NoError(t, v.VerifyHeader(makeHeaderWithRound(100, 0), nil)) // nil VRank
	})

	t.Run("valid self-report passes", func(t *testing.T) {
		v := newVerifier(t)
		assert.NoError(t, v.VerifyHeader(makeSelfReportHeader(t, num, candidates), nil))
	})

	t.Run("non-candidate rejected", func(t *testing.T) {
		unknown := numToAddr(150)
		v := newVerifier(t)
		assert.ErrorIs(t, v.VerifyHeader(makeSelfReportHeader(t, num, []common.Address{C1, unknown}), nil),
			vrank.ErrInvalidVRankCandidate)
	})

	t.Run("duplicate address rejected", func(t *testing.T) {
		v := newVerifier(t)
		assert.ErrorIs(t, v.VerifyHeader(makeSelfReportHeader(t, num, []common.Address{C1, C1}), nil),
			vrank.ErrDuplicateVRankCandidate)
	})

	t.Run("unsorted addresses rejected", func(t *testing.T) {
		v := newVerifier(t)
		// C3 > C2, so C3 before C2 is not ascending.
		assert.ErrorIs(t, v.VerifyHeader(makeSelfReportHeader(t, num, []common.Address{C3, C2}), nil),
			vrank.ErrVRankNotSorted)
	})

	t.Run("invalid encoding rejected before validation", func(t *testing.T) {
		v := newCN(t, withoutStart()).VRankModule
		h := makeHeaderWithRound(100, 0)
		h.VRank = []byte{0xff, 0xfe} // garbage
		assert.ErrorIs(t, v.VerifyHeader(h, nil), vrank.ErrInvalidVRankFormat)
	})

	t.Run("candidate membership is checked against the reporting block's epoch", func(t *testing.T) {
		cn := newCN(t, withoutStart())
		// The target is an earlier block of the same epoch, and CandTesting is epoch-stable.
		cn.Valset.EXPECT().GetCandTesting(num).Return(candidates, nil).Times(1)

		assert.NoError(t, cn.VRankModule.VerifyHeader(makeSelfReportHeader(t, num, []common.Address{C1}), nil))
	})

	t.Run("candidate lookup failure is returned", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		valset.EXPECT().GetCandTesting(num).Return(nil, assert.AnError).AnyTimes()
		v := newCN(t, withValset(valset), withoutStart()).VRankModule
		assert.ErrorIs(t, v.VerifyHeader(makeSelfReportHeader(t, num, []common.Address{C1}), nil), assert.AnError)
	})
}

func TestPrepareHeader(t *testing.T) {
	C1, C2 := numToAddr(1), numToAddr(2)
	candidates := []common.Address{C1, C2}

	t.Run("pre-fork clears VRank", func(t *testing.T) {
		v := newCN(t, withHardfork("osaka"), withoutStart()).VRankModule
		header := makeSelfReportHeader(t, 100, []common.Address{C1})

		require.NoError(t, v.PrepareHeader(header))
		assert.Nil(t, header.VRank)
	})

	t.Run("epoch-start fills CandTesting", func(t *testing.T) {
		v := newCN(t, withCandidates(candidates), withoutStart(), withParentAt(params.DefaultVRankEpoch)).VRankModule
		header := &types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}

		require.NoError(t, v.PrepareHeader(header))
		expected, err := vrank.EncodeVRank(vrank.VRankPayload{Report: candidates})
		require.NoError(t, err)
		assert.Equal(t, expected, header.VRank)
		assert.NoError(t, v.VerifyHeader(header, nil))
	})

	t.Run("epoch-start with no candidates leaves VRank nil", func(t *testing.T) {
		v := newCN(t, withCandidates(nil), withoutStart(), withParentAt(params.DefaultVRankEpoch)).VRankModule
		header := &types.Header{Number: big.NewInt(int64(params.DefaultVRankEpoch))}

		require.NoError(t, v.PrepareHeader(header))
		assert.Nil(t, header.VRank)
		assert.NoError(t, v.VerifyHeader(header, nil))
	})

	t.Run("non-epoch with no own proposal leaves VRank nil", func(t *testing.T) {
		v := newCN(t, withoutStart(), withParentAt(11)).VRankModule
		header := &types.Header{Number: big.NewInt(11)}

		require.NoError(t, v.PrepareHeader(header))
		assert.Nil(t, header.VRank)
	})

	// setupOwnProposal builds a CN that produced block `target` at `round` and collected its
	// own-proposal preprepare, so PrepareHeader can report it.
	setupOwnProposal := func(t *testing.T, target uint64, round int64) *VRankModule {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		valset.EXPECT().GetCandTesting(gomock.Any()).Return(candidates, nil).AnyTimes()
		targetHeader := makeHeaderWithRound(target, round)
		cn := newCN(t, withValset(valset), withoutStart(), withHeaders(map[uint64]*types.Header{
			target: targetHeader,
			// PrepareHeader(target+1) reads its parent for the certificate.
			target + 1: makeHeaderWithRound(target+1, 0),
		}))
		valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(cn.Addr, nil).AnyTimes()
		cn.VRankModule.collector.AddPrepreparedTime(
			vrank.ViewKey{N: target, R: uint8(round)}, time.Now(), types.NewBlockWithHeader(targetHeader).Hash())
		return cn.VRankModule
	}

	t.Run("non-epoch reports candidate failures of own prior proposal", func(t *testing.T) {
		v := setupOwnProposal(t, 10, 0)
		header := makeHeaderWithRound(11, 0)

		require.NoError(t, v.PrepareHeader(header))
		payload, err := vrank.DecodeVRank(header.VRank)
		require.NoError(t, err)
		assert.Equal(t, candidates, payload.Report) // no responses collected ⇒ all candidates failed
		assert.NoError(t, v.VerifyHeader(header, nil))
	})

	t.Run("high-round own proposal yields an empty report", func(t *testing.T) {
		// An own proposal committed above MaxRound is evaluated as empty (candidate msgs were dropped),
		// so it must not fail PrepareHeader.
		v := setupOwnProposal(t, 10, int64(vrank.MaxRound+1))
		header := &types.Header{Number: big.NewInt(11)}

		require.NoError(t, v.PrepareHeader(header))
		payload, err := vrank.DecodeVRank(header.VRank)
		require.NoError(t, err)
		assert.Empty(t, payload.Report)
	})
}

// TestSelectReportTarget follows one node through a realistic sequence: it reports the most
// recent block it committed, skips a round it proposed but another validator committed, re-reports
// idempotently across a round change, and moves on once a newer own block supersedes the target.
func TestSelectReportTarget(t *testing.T) {
	t.Run("reports most recent own block, skips lost round, superseded by newer", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		other := numToAddr(777)
		cn := newCN(t, withValset(valset), withoutStart(), withHeaders(map[uint64]*types.Header{
			7:  makeHeaderWithRound(7, 0),
			9:  makeHeaderWithRound(9, 0),
			11: makeHeaderWithRound(11, 0),
		}))
		v := cn.VRankModule
		// propose records the view exactly as HandleIstanbulPreprepare does when this node proposes.
		propose := func(n uint64) {
			v.collector.AddPrepreparedTime(vrank.ViewKey{N: n, R: 0}, time.Now(), common.Hash{})
		}
		// 7 and 11 are committed by this node; 9 is proposed by this node but committed by `other`.
		valset.EXPECT().GetProposer(uint64(7), uint64(0)).Return(cn.Addr, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(9), uint64(0)).Return(other, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(11), uint64(0)).Return(cn.Addr, nil).AnyTimes()

		// Block 7 is this node's first proposal in the epoch — nothing prior to report.
		_, _, ok := v.selectReportTarget(7)
		assert.False(t, ok, "first proposal in the epoch has no prior target")
		propose(7)

		// This node then builds block 9: that build reports block 7, and prune keeps 7 (strict <).
		// Keeping it matters because this node loses block 9 to `other`, so 7's report is discarded.
		target, round, ok := v.selectReportTarget(9)
		require.True(t, ok)
		assert.Equal(t, uint64(7), target)
		assert.Equal(t, uint64(0), round)
		v.collector.PruneReported(target) // keeps 7
		propose(9)

		// Building block 11: block 9 was lost, so 7 is re-reported (idempotent) and 9 is skipped.
		target, _, ok = v.selectReportTarget(11)
		require.True(t, ok)
		assert.Equal(t, uint64(7), target)
		v.collector.PruneReported(target) // still keeps 7

		// Once block 11 commits it becomes the most recent own block: 7 is superseded, 9 still skipped.
		propose(11)
		target, _, ok = v.selectReportTarget(13)
		require.True(t, ok)
		assert.Equal(t, uint64(11), target)
	})

	// Under a hash lock the committed block keeps the round-0 proposer's seal, but the candidates
	// that matter answered the round that committed it. The target belongs to that round's
	// proposer, so the round-0 proposer must not claim it and report a view nobody asked for.
	t.Run("hash-locked block belongs to the committing round", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		reproposer := numToAddr(777)
		cn := newCN(t, withValset(valset), withoutStart(), withHeaders(map[uint64]*types.Header{
			7: makeHeaderWithRound(7, 2), // proposed here at round 0, committed at round 2
		}))
		v := cn.VRankModule
		v.collector.AddPrepreparedTime(vrank.ViewKey{N: 7, R: 0}, time.Now(), common.Hash{})
		valset.EXPECT().GetProposer(uint64(7), uint64(2)).Return(reproposer, nil).AnyTimes()

		_, _, ok := v.selectReportTarget(9)
		assert.False(t, ok, "the round-0 proposer holds no view for the round that committed")
	})

	t.Run("prior-epoch proposal is not selected", func(t *testing.T) {
		epoch := uint64(params.DefaultVRankEpoch)
		v := newCN(t, withoutStart()).VRankModule
		v.collector.AddPrepreparedTime(vrank.ViewKey{N: epoch - 1, R: 0}, time.Now(), common.Hash{}) // previous epoch

		_, _, ok := v.selectReportTarget(epoch + 1)
		assert.False(t, ok)
	})
}

// certifiedParent returns a parent committed at `round`, the hash the child must carry as
// ParentHash, and the signers' seals over that (hash, round).
func certifiedParent(t *testing.T, number uint64, round byte, signers []*ecdsa.PrivateKey) (*types.Header, common.Hash, [][]byte) {
	t.Helper()
	parent := makeHeaderWithRound(number, int64(round))
	hash := newTestSealer().HeaderHash(parent)
	seals := make([][]byte, 0, len(signers))
	for _, key := range signers {
		seal, err := istanbul.NewSealerImpl(key).MakeCommittedSealFromHashWithRound(hash, round)
		require.NoError(t, err)
		seals = append(seals, seal)
	}
	return parent, hash, seals
}

// TestVerifyHeaderParentRound covers the parent-round certificate: a committee quorum must have
// committed the parent at the claimed round.
func TestVerifyHeaderParentRound(t *testing.T) {
	const (
		parentNum   = uint64(20)
		parentRound = byte(2)
	)

	keys := make([]*ecdsa.PrivateKey, 4)
	committee := make([]common.Address, 4)
	for i := range keys {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i] = key
		committee[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	strangerKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// quorum over a 4-member committee is 3
	_, parentHash, quorumSeals := certifiedParent(t, parentNum, parentRound, keys[:3])
	_, _, wrongRoundSeals := certifiedParent(t, parentNum, parentRound+1, keys[:3])
	_, _, strangerSeals := certifiedParent(t, parentNum, parentRound, []*ecdsa.PrivateKey{keys[0], keys[1], strangerKey})
	_, _, dupSeals := certifiedParent(t, parentNum, parentRound, []*ecdsa.PrivateKey{keys[0], keys[0], keys[1]})
	_, _, overlongSeals := certifiedParent(t, parentNum, parentRound, append(keys, keys[0]))

	testcases := []struct {
		name        string
		payload     vrank.VRankPayload
		expectedErr error
	}{
		{"quorum of the round's committee is accepted", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: quorumSeals}, nil},
		{"seals bound to another round are rejected", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: wrongRoundSeals}, vrank.ErrInvalidParentCertificate},
		{"a signer outside the committee is rejected", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: strangerSeals}, vrank.ErrInvalidParentCertificate},
		{"a repeated signer is rejected", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: dupSeals}, vrank.ErrInvalidParentCertificate},
		{"below quorum is rejected", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: quorumSeals[:2]}, vrank.ErrInvalidParentCertificate},
		{"a round without seals is rejected", vrank.VRankPayload{ParentRound: parentRound}, vrank.ErrInvalidParentCertificate},
		{"more seals than the committee is rejected", vrank.VRankPayload{ParentRound: parentRound, ParentCommittedSeal: overlongSeals}, vrank.ErrTooManyParentSeals},
		{"round zero with seals is rejected", vrank.VRankPayload{ParentCommittedSeal: quorumSeals}, vrank.ErrUnexpectedParentSeal},
		{"round zero without seals is accepted", vrank.VRankPayload{}, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
			valset.EXPECT().GetCommittee(parentNum, uint64(parentRound)).Return(committee, nil).AnyTimes()
			valset.EXPECT().GetCandTesting(gomock.Any()).Return(nil, nil).AnyTimes()
			v := newCN(t, withValset(valset), withoutStart()).VRankModule

			encoded, err := vrank.EncodeVRank(tc.payload)
			require.NoError(t, err)
			header := &types.Header{
				Number:     big.NewInt(int64(parentNum + 1)),
				ParentHash: parentHash,
				VRank:      encoded,
			}

			err = v.VerifyHeader(header, nil)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}

// TestPrepareHeaderCopiesParentCertificate checks that the certificate is copied from this node's
// own stored parent and that the result verifies.
func TestPrepareHeaderCopiesParentCertificate(t *testing.T) {
	const (
		parentNum   = uint64(30)
		parentRound = byte(2)
	)

	keys := make([]*ecdsa.PrivateKey, 4)
	committee := make([]common.Address, 4)
	for i := range keys {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i] = key
		committee[i] = crypto.PubkeyToAddress(key.PublicKey)
	}

	parent, parentHash, seals := certifiedParent(t, parentNum, parentRound, keys[:3])
	require.NoError(t, newTestSealer().WriteCommittedSeals(parent, seals))

	valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
	valset.EXPECT().GetCommittee(parentNum, uint64(parentRound)).Return(committee, nil).AnyTimes()
	valset.EXPECT().GetCandTesting(gomock.Any()).Return(nil, nil).AnyTimes()
	v := newCN(t, withValset(valset), withoutStart(),
		withHeaders(map[uint64]*types.Header{parentNum: parent})).VRankModule

	header := &types.Header{Number: big.NewInt(int64(parentNum + 1)), ParentHash: parentHash}
	require.NoError(t, v.PrepareHeader(header))

	payload, err := vrank.DecodeVRank(header.VRank)
	require.NoError(t, err)
	assert.Equal(t, parentRound, payload.ParentRound)
	assert.Equal(t, seals, payload.ParentCommittedSeal)
	assert.NoError(t, v.VerifyHeader(header, nil))
}

// TestParentRoundAbsentAtEpochStart checks that an epoch start records no parent round: its parent
// belongs to the previous epoch, whose failures may not score into this one.
func TestParentRoundAbsentAtEpochStart(t *testing.T) {
	const epochStart = uint64(params.DefaultVRankEpoch)

	t.Run("PrepareHeader drops a non-zero parent round", func(t *testing.T) {
		parent := makeHeaderWithRound(epochStart-1, 2)
		v := newCN(t, withCandidates(nil), withoutStart(),
			withHeaders(map[uint64]*types.Header{epochStart - 1: parent})).VRankModule

		header := &types.Header{Number: new(big.Int).SetUint64(epochStart)}
		require.NoError(t, v.PrepareHeader(header))

		payload, err := vrank.DecodeVRank(header.VRank)
		require.NoError(t, err)
		assert.Zero(t, payload.ParentRound)
		assert.Empty(t, payload.ParentCommittedSeal)
	})

	t.Run("VerifyHeader rejects a recorded parent round", func(t *testing.T) {
		v := newCN(t, withCandidates(nil), withoutStart()).VRankModule

		encoded, err := vrank.EncodeVRank(vrank.VRankPayload{ParentRound: 2})
		require.NoError(t, err)
		header := &types.Header{Number: new(big.Int).SetUint64(epochStart), VRank: encoded}

		assert.ErrorIs(t, v.VerifyHeader(header, nil), vrank.ErrUnexpectedParentRound)
	})
}
