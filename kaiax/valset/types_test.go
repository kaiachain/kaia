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

package valset

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

var (
	addr1 = common.HexToAddress("0x0001")
	addr2 = common.HexToAddress("0x0002")
	addr3 = common.HexToAddress("0x0003")
)

func TestMarkSuspended(t *testing.T) {
	m := NodeStateMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive},
		addr3: {State: ValPaused},
	}
	m.MarkSuspended([]common.Address{addr1, addr3})

	assert.True(t, m[addr1].Suspended, "addr1 should be suspended")
	assert.False(t, m[addr2].Suspended, "addr2 should not be suspended")
	assert.True(t, m[addr3].Suspended, "addr3 should be suspended")
}

func TestMarkSuspended_Empty(t *testing.T) {
	m := NodeStateMap{
		addr1: {State: ValActive},
	}
	m.MarkSuspended(nil)

	assert.False(t, m[addr1].Suspended, "no suspended list → all false")
}

func TestExcludeSuspended(t *testing.T) {
	m := NodeStateMap{
		addr1: {State: ValActive, Suspended: true},
		addr2: {State: ValActive, Suspended: false},
		addr3: {State: ValActive, Suspended: true},
	}
	filtered := m.ExcludeSuspended()

	assert.Len(t, filtered, 1)
	assert.NotNil(t, filtered[addr2])
	assert.Nil(t, filtered[addr1])
	assert.Nil(t, filtered[addr3])
}

func TestCountByState(t *testing.T) {
	m := NodeStateMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive},
		addr3: {State: ValPaused},
	}
	assert.Equal(t, uint64(2), m.CountByState(ValActive))
	assert.Equal(t, uint64(1), m.CountByState(ValPaused))
	assert.Equal(t, uint64(0), m.CountByState(ValExiting))
}
