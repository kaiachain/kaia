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

package vrank

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	blstypes "github.com/kaiachain/kaia/crypto/bls/types"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeVRank(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x15d34AAf54267DB7D7cC839724318F2730aC377B"),
		common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC"),
	}
	seals := [][]byte{make([]byte, 65), make([]byte, 65)}

	cases := []struct {
		name        string
		payload     VRankPayload
		expectEmpty bool
	}{
		// Decoding always yields non-nil slices, so the expectations spell out the empty ones.
		{"report only", VRankPayload{Report: addrs, ParentCommittedSeal: [][]byte{}}, false},
		{"parent round with seals", VRankPayload{Report: []common.Address{}, ParentRound: 3, ParentCommittedSeal: seals}, false},
		{"both", VRankPayload{Report: addrs, ParentRound: 1, ParentCommittedSeal: seals}, false},
		{"empty", VRankPayload{Report: []common.Address{}, ParentCommittedSeal: [][]byte{}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := EncodeVRank(tc.payload)
			require.NoError(t, err)
			if tc.expectEmpty {
				assert.Nil(t, enc, "an empty payload is encoded as absent bytes")
			} else {
				assert.NotEmpty(t, enc)
			}
			got, err := DecodeVRank(enc)
			require.NoError(t, err)
			assert.Equal(t, tc.payload.Report, got.Report)
			assert.Equal(t, tc.payload.ParentRound, got.ParentRound)
			assert.Equal(t, tc.payload.ParentCommittedSeal, got.ParentCommittedSeal)
		})
	}
}

// The pre-fixed-size shapes. A wrong signature length is only expressible from here.
type varSizeCandidate struct {
	BlockNumber uint64
	Round       uint8
	BlockHash   common.Hash
	Sig         []byte
	BlsSig      []byte
}

type varSizePreprepare struct {
	Block *types.Block
	View  *bft.View
	Sig   []byte
}

// The fixed-size fields are the only length guard: node/cn dropped its check for VRankCandidate
// and never had one for VRankPreprepare.
func TestVRankRLPSignatureLengths(t *testing.T) {
	const (
		sigLen = crypto.SignatureLength
		blsLen = blstypes.SignatureLength
	)
	var (
		block   = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view    = &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}
		bytesOf = func(n int) []byte { return bytes.Repeat([]byte{0x11}, n) }
		cand    = func(sig, bls int) any {
			return &varSizeCandidate{BlockNumber: 1, Sig: bytesOf(sig), BlsSig: bytesOf(bls)}
		}
		prep = func(sig int) any { return &varSizePreprepare{Block: block, View: view, Sig: bytesOf(sig)} }
	)

	for _, tc := range []struct {
		name    string
		msg     any
		target  any
		wantErr bool
	}{
		{"candidate exact", cand(sigLen, blsLen), new(VRankCandidate), false},
		{"candidate empty sig", cand(0, blsLen), new(VRankCandidate), true},
		{"candidate empty BLS sig", cand(sigLen, 0), new(VRankCandidate), true},
		{"candidate undersized sig", cand(sigLen-1, blsLen), new(VRankCandidate), true},
		{"candidate undersized BLS sig", cand(sigLen, blsLen-1), new(VRankCandidate), true},
		{"candidate oversized sig", cand(1<<20, blsLen), new(VRankCandidate), true},
		{"candidate oversized BLS sig", cand(sigLen, blsLen+1), new(VRankCandidate), true},
		{"preprepare exact", prep(sigLen), new(VRankPreprepare), false},
		{"preprepare empty sig", prep(0), new(VRankPreprepare), true},
		{"preprepare undersized sig", prep(sigLen - 1), new(VRankPreprepare), true},
		{"preprepare oversized sig", prep(1 << 20), new(VRankPreprepare), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := rlp.EncodeToBytes(tc.msg)
			require.NoError(t, err)

			err = rlp.DecodeBytes(enc, tc.target)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
