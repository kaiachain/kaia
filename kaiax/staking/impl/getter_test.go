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

package staking

import (
	"testing"

	"gotest.tools/assert"
)

func TestSourceBlockNum(t *testing.T) {
	testcases := []struct {
		num      uint64
		isKaia   bool
		interval uint64
		expected uint64
	}{
		// Before Kaia
		{0, false, 1000, 0},
		{1, false, 1000, 0},
		{1000, false, 1000, 0},
		{1001, false, 1000, 0},
		{1999, false, 1000, 0},
		{2000, false, 1000, 0},
		{2001, false, 1000, 1000},
		{2999, false, 1000, 1000},
		{3000, false, 1000, 1000},
		{3001, false, 1000, 2000},

		// After Kaia
		{0, true, 1000, 0},
		{1, true, 1000, 0},
		{1000, true, 1000, 999},
		{1001, true, 1000, 1000},
		{1999, true, 1000, 1998},
	}

	for i, tc := range testcases {
		actual := sourceBlockNum(tc.num, tc.isKaia, tc.interval)
		assert.Equal(t, tc.expected, actual, i)
	}
}
