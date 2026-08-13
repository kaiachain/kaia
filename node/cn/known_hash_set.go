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
	"sync"

	"github.com/kaiachain/kaia/common"
)

// knownHashSet is a bounded, FIFO set of common.Hash. It is used per peer to
// track which transaction/block/bid hashes a peer is already known to have, so
// the same item is not re-sent (gossip de-duplication).
//
// Bounding is safe because the set is only an optimization, never a correctness
// mechanism: it can only suppress a send. Evicting an entry can therefore at
// worst cause one redundant send (the peer already has the item, dedups it
// locally, and does not re-announce it), never a missed delivery. Entries are
// only useful during an item's brief propagation window, so a fixed bound that
// covers that window gives the full dedup benefit; FIFO eviction drops the
// oldest-inserted (already aged-out) entries first, which is exactly what we
// want. Access-recency (LRU) is irrelevant here, and the previous cache was in
// fact configured FIFO (Get == Peek). Callers use only membership (KnowsX) and
// insertion-with-eviction (AddToKnownX).
//
// Why a custom type instead of golang-lru: golang-lru stores entries in a
// container/list (doubly-linked list) keyed by a map, and the generic
// common.Cache wrapper boxes the key into an interface. That is ~3-4 heap
// pointers per entry, so the Go GC mark phase must traverse the whole structure
// every cycle (runtime.findObject). With one set per peer scaled by RAM, this
// reached tens of millions of pointer-rich live objects and made GC the
// dominant CPU consumer on long-running, many-peer nodes.
//
// knownHashSet keeps hashes in a preallocated ring (for O(1) FIFO eviction) and
// a map[common.Hash]struct{} for O(1) membership. common.Hash is a [32]byte
// value and struct{} is empty, so BOTH the ring and the map are pointer-free
// ("noscan"): the GC never traverses pointers inside this structure regardless
// of how many entries it holds. It also does no per-entry allocation at steady
// state and uses less memory than the equivalent LRU (no list nodes, no
// interface boxing of keys).
type knownHashSet struct {
	mu   sync.RWMutex // reads (Contains) dominate via KnowsX during broadcast filtering
	max  int
	set  map[common.Hash]struct{}
	ring []common.Hash // preallocated; len == max
	next int           // index of the oldest entry / next write position
	full bool
}

// newKnownHashSet returns an empty set that holds at most max hashes.
func newKnownHashSet(max int) *knownHashSet {
	if max < 1 {
		max = 1
	}
	return &knownHashSet{
		max:  max,
		set:  make(map[common.Hash]struct{}, max),
		ring: make([]common.Hash, max),
	}
}

// Add records hash as known. When the set is full, the oldest inserted hash is
// evicted. Adding a hash that is already present is a no-op: insertion order
// (and therefore eviction order) is preserved. This matches the de-duplication
// intent and avoids unnecessary churn.
func (s *knownHashSet) Add(hash common.Hash) {
	s.mu.Lock()
	if _, ok := s.set[hash]; !ok {
		if s.full {
			delete(s.set, s.ring[s.next])
		}
		s.ring[s.next] = hash
		s.set[hash] = struct{}{}
		s.next++
		if s.next == s.max {
			s.next = 0
			s.full = true
		}
	}
	s.mu.Unlock()
}

// Contains reports whether hash is currently in the set.
func (s *knownHashSet) Contains(hash common.Hash) bool {
	s.mu.RLock()
	_, ok := s.set[hash]
	s.mu.RUnlock()
	return ok
}

// Len returns the current number of hashes held.
func (s *knownHashSet) Len() int {
	s.mu.RLock()
	n := len(s.set)
	s.mu.RUnlock()
	return n
}
