// Copyright 2024 The Kaia Authors
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
package validator

import (
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func TestCalcSeed(t *testing.T) {
	type testCase struct {
		hash        common.Hash
		expected    int64
		expectedErr error
	}
	testCases := []testCase{
		{
			hash:        common.HexToHash("0x1111"),
			expected:    0,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234000000000000000000000000000000000000000000000000000000000000"),
			expected:    81979586966978560,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341230000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234123412341234123412341234123412341234123412341234"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x0034123412341234123412341234123412341234123412341234123412341234"),
			expected:    916044602622243,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0xabcdef3412341234123412341234123412341234123412341234123412341234"),
			expected:    773738372352131363,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		actual, err := ConvertHashToSeed(tc.hash)
		assert.Equal(t, err, tc.expectedErr)
		assert.Equal(t, actual, tc.expected)
	}
}
