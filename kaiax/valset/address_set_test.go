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

package valset

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func numsToAddrs(n ...int) []common.Address {
	addrs := make([]common.Address, len(n))
	for i, num := range n {
		addrs[i] = common.BigToAddress(big.NewInt(int64(num)))
	}
	return addrs
}

func TestAddressSet_List(t *testing.T) {
	var (
		unsorted = numsToAddrs(0x120c, 0x120b, 0x120a)
		sorted   = numsToAddrs(0x120a, 0x120c, 0x120b)
	)

	as := NewAddressSet(unsorted)
	assert.Equal(t, sorted, []common.Address(as.list))
	assert.Equal(t, sorted, []common.Address(as.Copy().list))
	assert.Equal(t, sorted, as.List())
	assert.Equal(t, 3, as.Len())

	assert.Equal(t, sorted[0], as.At(0))
	assert.Equal(t, sorted[1], as.At(1))
	assert.Equal(t, sorted[2], as.At(2))

	assert.Equal(t, 0, as.IndexOf(sorted[0]))
	assert.Equal(t, 1, as.IndexOf(sorted[1]))
	assert.Equal(t, 2, as.IndexOf(sorted[2]))
	assert.Equal(t, -1, as.IndexOf(common.HexToAddress("0x120d")))

	assert.True(t, as.Contains(sorted[0]))
	assert.True(t, as.Contains(sorted[1]))
	assert.True(t, as.Contains(sorted[2]))
	assert.False(t, as.Contains(common.HexToAddress("0x120d")))
}

func TestAddressSet_Set(t *testing.T) {
	// because of mixed-case checksum
	// the sorted order is 120A-120C-120D-120b.
	var (
		unsorted = numsToAddrs(0x120c, 0x120b, 0x120a)
		sorted1  = numsToAddrs(0x120a, 0x120c, 0x120d, 0x120b)
		sorted2  = numsToAddrs(0x120a, 0x120c, 0x120b)
	)

	as := NewAddressSet(unsorted)

	as.Add(common.HexToAddress("0x120d"))
	assert.Equal(t, sorted1, as.List())

	as.Remove(common.HexToAddress("0x120d"))
	assert.Equal(t, sorted2, as.List())
}

func TestAddressSet_Subtract(t *testing.T) {
	var (
		from   = numsToAddrs(1, 2, 3, 4, 5)
		to     = numsToAddrs(3, 4, 5, 6, 7)
		result = numsToAddrs(1, 2)
	)

	a := NewAddressSet(from)
	b := NewAddressSet(to)
	c := a.Subtract(b)
	assert.Equal(t, result, c.List())
}

func TestAddressSet_Shuffle(t *testing.T) {
	var (
		// Example from the README:SelectRandomCommittee
		unsorted       = numsToAddrs(0, 1, 2, 4, 5, 6, 8, 9)
		seed           = int64(0x112233445566778)
		shuffledLegacy = numsToAddrs(6, 8, 1, 2, 4, 0, 9, 5)
		shuffled       = numsToAddrs(5, 6, 4, 8, 9, 0, 2, 1)
	)

	as := NewAddressSet(unsorted)
	assert.Equal(t, shuffledLegacy, as.ShuffledListLegacy(seed))
	assert.Equal(t, shuffled, as.ShuffledList(seed))
}

func TestHashToSeed(t *testing.T) {
	testcases := []struct {
		hash       common.Hash
		seedLegacy int64
		seed       int64
	}{
		{common.HexToHash("0x1111"), 0x0, 0x0},
		{common.HexToHash("0x1234000000000000000000000000000000000000000000000000000000000000"), 0x123400000000000, 0x1234000000000000},
		{common.HexToHash("0x1234123412341230000000000000000000000000000000000000000000000000"), 0x123412341234123, 0x1234123412341230},
		{common.HexToHash("0x1234123412341234000000000000000000000000000000000000000000000000"), 0x123412341234123, 0x1234123412341234},
		{common.HexToHash("0x1234123412341234123412341234123412341234123412341234123412341234"), 0x123412341234123, 0x1234123412341234},
		{common.HexToHash("0x0034123412341234123412341234123412341234123412341234123412341234"), 0x003412341234123, 0x0034123412341234},
		{common.HexToHash("0xabcdef3412341234123412341234123412341234123412341234123412341234"), 0xabcdef341234123, -6066930116075449804}, // int64(0xABCDEF3412341234)
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.seedLegacy, HashToSeedLegacy(tc.hash))
		assert.Equal(t, tc.seed, HashToSeed(tc.hash.Bytes()))
	}
}
