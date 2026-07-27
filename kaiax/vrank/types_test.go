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
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeReport(t *testing.T) {
	cases := []struct {
		name   string
		report []common.Address
	}{
		{
			name:   "addresses",
			report: []common.Address{common.HexToAddress("0x15d34AAf54267DB7D7cC839724318F2730aC377B"), common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC")},
		},
		{
			name:   "empty",
			report: []common.Address{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := EncodeReport(tc.report)
			assert.NoError(t, err)
			if len(tc.report) == 0 {
				assert.Nil(t, enc)
				return
			}
			assert.NotEmpty(t, enc)
			dec, err := DecodeReport(enc)
			assert.NoError(t, err)
			assert.Equal(t, tc.report, dec)
		})
	}
}

func TestEncodeAddressList_EmptyList(t *testing.T) {
	enc, err := EncodeAddressList(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, enc)

	dec, err := DecodeReport(enc)
	assert.NoError(t, err)
	assert.Empty(t, dec)
}

func TestVRankCandidateRLPRejectsOversizedSignatures(t *testing.T) {
	// fixed-length signatures round-trip
	valid := VRankCandidate{
		BlockNumber: 1,
		Round:       0,
		BlockHash:   common.HexToHash("0x42"),
		Sig:         [65]byte{0x11},
		BlsSig:      [96]byte{0x22},
	}
	enc, err := rlp.EncodeToBytes(&valid)
	assert.NoError(t, err)
	var back VRankCandidate
	assert.NoError(t, rlp.DecodeBytes(enc, &back))
	assert.Equal(t, valid, back)

	// oversized signature bytes fail to decode into the fixed-size fields
	type wire struct {
		BlockNumber uint64
		Round       uint8
		BlockHash   common.Hash
		Sig         []byte
		BlsSig      []byte
	}
	oversized, err := rlp.EncodeToBytes(&wire{
		BlockNumber: 1,
		Sig:         bytes.Repeat([]byte{0x11}, 1<<20),
		BlsSig:      bytes.Repeat([]byte{0x22}, 96),
	})
	assert.NoError(t, err)
	assert.Error(t, rlp.DecodeBytes(oversized, new(VRankCandidate)))
}
