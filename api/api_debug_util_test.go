// Copyright 2025 The Kaia Authors
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

package api

import (
	"testing"

	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestDebugUtilAPI_ChaindbCompact_IterateRange(t *testing.T) {
	testcases := []struct {
		desc        string // given
		startHex    string
		stepHex     string
		endHex      string
		segStartHex []string // expected
		segEndHex   []string
		err         error
	}{
		{
			"one-byte range starting from zero", "0x", "0x01", "0x10",
			[]string{"0x00", "0x01", "0x02", "0x03", "0x04", "0x05", "0x06", "0x07", "0x08", "0x09", "0x0a", "0x0b", "0x0c", "0x0d", "0x0e", "0x0f"},
			[]string{"0x01", "0x02", "0x03", "0x04", "0x05", "0x06", "0x07", "0x08", "0x09", "0x0a", "0x0b", "0x0c", "0x0d", "0x0e", "0x0f", "0x10"},
			nil,
		},
		{
			"one-byte range ending with nil", "0xf0", "0x01", "0x",
			[]string{"0xf0", "0xf1", "0xf2", "0xf3", "0xf4", "0xf5", "0xf6", "0xf7", "0xf8", "0xf9", "0xfa", "0xfb", "0xfc", "0xfd", "0xfe", "0xff"},
			[]string{"0xf1", "0xf2", "0xf3", "0xf4", "0xf5", "0xf6", "0xf7", "0xf8", "0xf9", "0xfa", "0xfb", "0xfc", "0xfd", "0xfe", "0xff", "0x"},
			nil,
		},
		{
			"step overruns end", "0x", "0x30", "0x77",
			[]string{"0x00", "0x30", "0x60"},
			[]string{"0x30", "0x60", "0x77"}, // next segment end (0x90) is clipped to range end (0x77)
			nil,
		},
		{
			"step overruns nil", "0x", "0x30", "0x",
			[]string{"0x00", "0x30", "0x60", "0x90", "0xc0", "0xf0"},
			[]string{"0x30", "0x60", "0x90", "0xc0", "0xf0", "0x"}, // next segment end overflows, clipped to nil.
			nil,
		},
		{
			"two-byte range", "0x42", "0x0010", "0x43",
			[]string{"0x4200", "0x4210", "0x4220", "0x4230", "0x4240", "0x4250", "0x4260", "0x4270", "0x4280", "0x4290", "0x42a0", "0x42b0", "0x42c0", "0x42d0", "0x42e0", "0x42f0"},
			[]string{"0x4210", "0x4220", "0x4230", "0x4240", "0x4250", "0x4260", "0x4270", "0x4280", "0x4290", "0x42a0", "0x42b0", "0x42c0", "0x42d0", "0x42e0", "0x42f0", "0x4300"},
			nil,
		},
		{"step is zero", "0x", "0x00", "0x", nil, nil, errCompactRangeZeroStep},
		{"start is greater than end", "0x42", "0x01", "0x41", nil, nil, errCompactRangeStartEnd},
	}

	for _, tc := range testcases {
		segStartHex := make([]string, 0)
		segEndHex := make([]string, 0)

		err := iterateRange(hexutil.MustDecode(tc.startHex), hexutil.MustDecode(tc.stepHex), hexutil.MustDecode(tc.endHex), func(segStart, segEnd []byte) error {
			segStartHex = append(segStartHex, hexutil.Encode(segStart))
			segEndHex = append(segEndHex, hexutil.Encode(segEnd))
			return nil
		})
		if tc.err != nil {
			assert.ErrorIs(t, err, tc.err, tc.desc)
		} else {
			assert.NoError(t, err, tc.desc)
			assert.Equal(t, tc.segStartHex, segStartHex, tc.desc)
			assert.Equal(t, tc.segEndHex, segEndHex, tc.desc)
		}
	}
}
