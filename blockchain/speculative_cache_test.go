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

package blockchain

import (
	"errors"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpeculativeCache_ReserveThenGet(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")
	result := &SpeculativeResult{UsedGas: 42}

	entry := cache.Reserve(hash)
	entry.Complete(result, nil)

	got := cache.TryGet(hash, time.Second)
	require.NotNil(t, got)
	assert.Equal(t, uint64(42), got.UsedGas)
	assert.Equal(t, uint64(1), cache.Hits())
}

func TestSpeculativeCache_HashMismatch(t *testing.T) {
	cache := NewSpeculativeResultCache()
	entry := cache.Reserve(common.HexToHash("0xaaa"))
	entry.Complete(&SpeculativeResult{}, nil)

	got := cache.TryGet(common.HexToHash("0xbbb"), time.Second)
	assert.Nil(t, got, "hash mismatch should return nil")
	assert.Equal(t, uint64(0), cache.Hits())
}

func TestSpeculativeCache_NoEntry(t *testing.T) {
	cache := NewSpeculativeResultCache()
	got := cache.TryGet(common.HexToHash("0xaaa"), time.Second)
	assert.Nil(t, got)
}

func TestSpeculativeCache_OverwriteOnReserve(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hashA := common.HexToHash("0xaaa")
	hashB := common.HexToHash("0xbbb")

	entryA := cache.Reserve(hashA)
	entryA.Complete(&SpeculativeResult{UsedGas: 1}, nil)

	// Reserve for new hash overwrites old entry.
	entryB := cache.Reserve(hashB)
	entryB.Complete(&SpeculativeResult{UsedGas: 2}, nil)

	// Old hash misses, new hash hits.
	assert.Nil(t, cache.TryGet(hashA, time.Second))
	got := cache.TryGet(hashB, time.Second)
	require.NotNil(t, got)
	assert.Equal(t, uint64(2), got.UsedGas)
}

func TestSpeculativeCache_TimeoutOnInFlight(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	// Reserve but never complete.
	cache.Reserve(hash)

	got := cache.TryGet(hash, 50*time.Millisecond)
	assert.Nil(t, got, "should timeout on in-flight entry")
	assert.Equal(t, uint64(0), cache.Hits())
}

func TestSpeculativeCache_WaitForInFlight(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")
	result := &SpeculativeResult{UsedGas: 99}

	entry := cache.Reserve(hash)

	// Complete after a short delay.
	go func() {
		time.Sleep(20 * time.Millisecond)
		entry.Complete(result, nil)
	}()

	got := cache.TryGet(hash, time.Second)
	require.NotNil(t, got)
	assert.Equal(t, uint64(99), got.UsedGas)
}

func TestSpeculativeCache_ErrorEntry(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	entry := cache.Reserve(hash)
	entry.Complete(nil, errors.New("exec failed"))

	got := cache.TryGet(hash, time.Second)
	assert.Nil(t, got, "error entry should return nil")
}

func TestSpeculativeCache_Clear(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	entry := cache.Reserve(hash)
	entry.Complete(&SpeculativeResult{}, nil)

	cache.Clear()
	assert.Nil(t, cache.TryGet(hash, time.Second))
}

func TestSpeculativeCache_ClearNoEntry(t *testing.T) {
	cache := NewSpeculativeResultCache()
	cache.Clear() // should not panic
}

func TestSpeculativeCache_NonBlockingMiss(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	// Reserve but don't complete — non-blocking TryGet must return nil immediately.
	cache.Reserve(hash)

	got := cache.TryGet(hash, 0)
	assert.Nil(t, got, "non-blocking TryGet on in-flight entry should return nil")
}

func TestSpeculativeCache_NonBlockingHit(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	entry := cache.Reserve(hash)
	entry.Complete(&SpeculativeResult{UsedGas: 77}, nil)

	// Non-blocking TryGet should succeed because entry is already done.
	got := cache.TryGet(hash, 0)
	require.NotNil(t, got)
	assert.Equal(t, uint64(77), got.UsedGas)
}

func TestSpeculativeCache_DoubleComplete(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	entry := cache.Reserve(hash)
	entry.Complete(&SpeculativeResult{UsedGas: 1}, nil)
	// Second Complete must not panic (sync.Once guards the close).
	entry.Complete(&SpeculativeResult{UsedGas: 2}, nil)

	got := cache.TryGet(hash, time.Second)
	require.NotNil(t, got)
	assert.Equal(t, uint64(1), got.UsedGas, "first Complete's value should win")
}

func TestSpeculativeCache_HitClearsEntry(t *testing.T) {
	cache := NewSpeculativeResultCache()
	hash := common.HexToHash("0xaaa")

	entry := cache.Reserve(hash)
	entry.Complete(&SpeculativeResult{UsedGas: 1}, nil)

	// First get succeeds and auto-clears.
	got := cache.TryGet(hash, time.Second)
	require.NotNil(t, got)

	// Second get misses because TryGet cleared the consumed entry.
	got2 := cache.TryGet(hash, time.Second)
	assert.Nil(t, got2)
	assert.Equal(t, uint64(1), cache.Hits())
}
