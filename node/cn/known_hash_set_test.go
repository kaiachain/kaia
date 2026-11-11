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

package cn

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/kaiachain/kaia/common"
)

func bhash(i int) common.Hash {
	var h common.Hash
	binary.BigEndian.PutUint64(h[0:8], uint64(i))
	binary.BigEndian.PutUint64(h[24:32], uint64(i)) // spread low/high bytes
	return h
}

func TestKnownHashSet(t *testing.T) {
	s := newKnownHashSet(4)
	for i := 0; i < 4; i++ {
		s.Add(bhash(i))
	}
	if s.Len() != 4 {
		t.Fatalf("len=%d, want 4", s.Len())
	}
	for i := 0; i < 4; i++ {
		if !s.Contains(bhash(i)) {
			t.Fatalf("hash %d missing", i)
		}
	}
	// Re-adding an existing hash must be a no-op (no eviction, order preserved).
	s.Add(bhash(1))
	if s.Len() != 4 || !s.Contains(bhash(0)) {
		t.Fatalf("re-add changed the set: len=%d", s.Len())
	}
	// Adding a new hash to a full set evicts the oldest (0), FIFO.
	s.Add(bhash(4))
	if s.Contains(bhash(0)) {
		t.Fatal("oldest entry (0) should have been evicted")
	}
	if !s.Contains(bhash(4)) || !s.Contains(bhash(1)) {
		t.Fatal("expected 1 and 4 to be present")
	}
	if s.Len() != 4 {
		t.Fatalf("len=%d, want 4 after eviction", s.Len())
	}
	// Next eviction must be the next-oldest (1, then 2), preserving FIFO order.
	s.Add(bhash(5))
	if s.Contains(bhash(1)) {
		t.Fatal("entry 1 should have been evicted next (FIFO)")
	}
}

func TestKnownHashSetMinSize(t *testing.T) {
	s := newKnownHashSet(0) // must clamp to >=1, never divide-by-zero / panic
	s.Add(bhash(1))
	s.Add(bhash(2))
	if s.Len() != 1 || !s.Contains(bhash(2)) || s.Contains(bhash(1)) {
		t.Fatalf("min-size set misbehaved: len=%d", s.Len())
	}
}

// Run with -race to validate the locking.
func TestKnownHashSetConcurrent(t *testing.T) {
	s := newKnownHashSet(2048)
	var wg sync.WaitGroup
	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 20000; i++ {
				h := bhash(g*1_000_000 + i)
				s.Add(h)
				_ = s.Contains(h)
			}
		}(g)
	}
	wg.Wait()
	if s.Len() > 2048 {
		t.Fatalf("exceeded capacity: len=%d", s.Len())
	}
}

// --- Benchmarks: current golang-lru-backed FIFO cache vs knownHashSet ---
//
// maxKnownTxs is 32768 today; on a 128GB host the previous IsScaled:true made
// it 262144 per peer. Both sizes are benchmarked.

func newOldCache(capacity int) common.Cache {
	return common.NewCache(common.FIFOCacheConfig{CacheSize: capacity, IsScaled: false})
}

// BenchmarkKnownAdd measures steady-state Add (cache pre-filled to capacity, so
// every Add both evicts the oldest and inserts a new hash) — the hot path on a
// node relaying transactions. Watch allocs/op and B/op.
func BenchmarkKnownAdd(b *testing.B) {
	const capacity = 32768
	b.Run("golang-lru", func(b *testing.B) {
		c := newOldCache(capacity)
		for i := 0; i < capacity; i++ {
			c.Add(bhash(i), struct{}{})
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Add(bhash(capacity+i), struct{}{})
		}
	})
	b.Run("knownHashSet", func(b *testing.B) {
		s := newKnownHashSet(capacity)
		for i := 0; i < capacity; i++ {
			s.Add(bhash(i))
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Add(bhash(capacity + i))
		}
	})
}

// BenchmarkKnownGCMark measures the cost of a full GC while the structure is
// live and full. ns/op is the time per runtime.GC(); the gap between the two
// implementations at the same n is the per-structure mark (pointer-scan) cost.
func BenchmarkKnownGCMark(b *testing.B) {
	for _, n := range []int{32768, 262144, 1_000_000} {
		b.Run(fmt.Sprintf("golang-lru/n=%d", n), func(b *testing.B) {
			c := newOldCache(n)
			for i := 0; i < n; i++ {
				c.Add(bhash(i), struct{}{})
			}
			runtime.GC()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				runtime.GC()
			}
			b.StopTimer()
			runtime.KeepAlive(c)
		})
		b.Run(fmt.Sprintf("knownHashSet/n=%d", n), func(b *testing.B) {
			s := newKnownHashSet(n)
			for i := 0; i < n; i++ {
				s.Add(bhash(i))
			}
			runtime.GC()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				runtime.GC()
			}
			b.StopTimer()
			runtime.KeepAlive(s)
		})
	}
}

// BenchmarkKnownGCMarkProd reproduces the production configuration that exposed
// the issue: 54 peers, each with a knownTxs set scaled to 32768*8 = 262144
// entries (IsScaled:true on a 128 GB host), modeled as 54 separate instances.
// ns/op is the time for one full runtime.GC() with all of them live and full.
func BenchmarkKnownGCMarkProd(b *testing.B) {
	const peers = 54
	const perPeer = 32768 * 8 // 262144: the IsScaled:true size on a 128 GB host

	b.Run("golang-lru", func(b *testing.B) {
		caches := make([]common.Cache, peers)
		for p := 0; p < peers; p++ {
			c := newOldCache(perPeer)
			for i := 0; i < perPeer; i++ {
				c.Add(bhash(p*perPeer+i), struct{}{})
			}
			caches[p] = c
		}
		runtime.GC()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runtime.GC()
		}
		b.StopTimer()
		runtime.KeepAlive(caches)
	})
	b.Run("knownHashSet", func(b *testing.B) {
		sets := make([]*knownHashSet, peers)
		for p := 0; p < peers; p++ {
			s := newKnownHashSet(perPeer)
			for i := 0; i < perPeer; i++ {
				s.Add(bhash(p*perPeer + i))
			}
			sets[p] = s
		}
		runtime.GC()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runtime.GC()
		}
		b.StopTimer()
		runtime.KeepAlive(sets)
	})
}

// TestFootprint reports the resident heap cost and live-object count per entry
// for both structures (informational; run with: go test -run TestFootprint -v).
// The per-entry live-object count is why GC mark cost differs: golang-lru keeps
// ~3 pointer-rich objects per entry, knownHashSet keeps ~0.
func TestFootprint(t *testing.T) {
	const n = 262144
	report := func(name string, before, after *runtime.MemStats) {
		dBytes := int64(after.HeapInuse) - int64(before.HeapInuse)
		dObjs := int64(after.HeapObjects) - int64(before.HeapObjects)
		t.Logf("%-14s n=%d  heapInuse %+.1f MB (%d B/entry)  liveObjects %+d (%.2f objs/entry)",
			name, n, float64(dBytes)/1e6, dBytes/int64(n), dObjs, float64(dObjs)/float64(n))
	}
	{
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		c := newOldCache(n)
		for i := 0; i < n; i++ {
			c.Add(bhash(i), struct{}{})
		}
		runtime.GC()
		runtime.ReadMemStats(&after)
		report("golang-lru", &before, &after)
		runtime.KeepAlive(c)
	}
	runtime.GC() // release the first structure before measuring the second
	{
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		s := newKnownHashSet(n)
		for i := 0; i < n; i++ {
			s.Add(bhash(i))
		}
		runtime.GC()
		runtime.ReadMemStats(&after)
		report("knownHashSet", &before, &after)
		runtime.KeepAlive(s)
	}
}
