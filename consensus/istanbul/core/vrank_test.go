// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package core

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
)

func TestVrank(t *testing.T) {
	var (
		N            = 6
		quorum       = 4
		committee, _ = genValidators(N)
		view         = istanbul.View{Sequence: big.NewInt(0), Round: big.NewInt(0)}
		msg          = &istanbul.Subject{View: &view}
		vrank        = NewVrank(view, committee)

		expectedAssessList  []uint8
		expectedLateCommits []time.Duration
	)

	committee = valset.NewAddressSet(committee).List() // sort it
	for i := 0; i < quorum; i++ {
		vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedEarly)
	}
	vrank.HandleCommitted(view.Sequence)
	for i := quorum; i < N; i++ {
		vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedLate)
		expectedLateCommits = append(expectedLateCommits, vrank.commitArrivalTimeMap[committee[i]])
	}

	bitmap, late := vrank.Bitmap(), vrank.LateCommits()
	assert.Equal(t, hex.EncodeToString(compress(expectedAssessList)), bitmap)
	assert.Equal(t, expectedLateCommits, late)
}

func TestVrankAssessBatch(t *testing.T) {
	arr := []time.Duration{0, 4, 1, vrankNotArrivedPlaceholder, 2}
	threshold := time.Duration(2)
	expected := []uint8{
		vrankArrivedEarly, vrankArrivedLate, vrankArrivedEarly, vrankNotArrived, vrankArrivedEarly,
	}
	assert.Equal(t, expected, assessBatch(arr, threshold))
}

func TestVrankSerialize(t *testing.T) {
	var (
		N            = 6
		committee, _ = genValidators(N)
		timeMap      = make(map[common.Address]time.Duration)
		expected     = make([]time.Duration, len(committee))
	)

	committee = valset.NewAddressSet(committee).List() // sort it
	for i, val := range committee {
		t := time.Duration((i * 100) * int(time.Millisecond))
		timeMap[val] = t
		expected[i] = t
	}

	assert.Equal(t, expected, serialize(committee, timeMap))
}

func TestVrankCompress(t *testing.T) {
	tcs := []struct {
		input    []uint8
		expected []byte
	}{
		{
			input:    []uint8{2},
			expected: []byte{0b10_000000},
		},
		{
			input:    []uint8{2, 1},
			expected: []byte{0b10_01_0000},
		},
		{
			input:    []uint8{0, 2, 1},
			expected: []byte{0b00_10_01_00},
		},
		{
			input:    []uint8{0, 2, 1, 1},
			expected: []byte{0b00_10_01_01},
		},
		{
			input:    []uint8{1, 2, 1, 2, 1},
			expected: []byte{0b01_10_01_10, 0b01_000000},
		},
		{
			input:    []uint8{1, 2, 1, 2, 1, 2},
			expected: []byte{0b01_10_01_10, 0b01_10_0000},
		},
		{
			input:    []uint8{1, 2, 1, 2, 1, 2, 1},
			expected: []byte{0b01_10_01_10, 0b01_10_01_00},
		},
		{
			input:    []uint8{1, 2, 1, 2, 1, 2, 0, 2},
			expected: []byte{0b01_10_01_10, 0b01_10_00_10},
		},
		{
			input:    []uint8{1, 1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 1, 1},
			expected: []byte{0b01011010, 0b01011010, 0b01011010, 0b01011010, 0b01010000},
		},
	}
	for i, tc := range tcs {
		assert.Equal(t, tc.expected, compress(tc.input), "tc %d failed", i)
	}
}
