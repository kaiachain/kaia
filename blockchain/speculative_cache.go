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
	"sync"
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
)

// SpeculativeResult holds the outputs of StateProcessor.Process, cached so
// that InsertChain can skip re-execution when the cache hits.
//
// The producer MUST populate this from Process() — the same function that
// InsertChain calls on the normal path. This guarantees identical results
// including FinalizeState effects (block rewards, state root). Do NOT
// mirror Process() manually; call it directly.
type SpeculativeResult struct {
	State            *state.StateDB
	Receipts         types.Receipts
	Logs             []*types.Log
	UsedGas          uint64
	InternalTxTraces []*vm.InternalTxTrace
	ProcessStats     ProcessStats // timing from speculative Process() call
}

// speculativeEntry is the internal cache slot. It acts as a future:
// the producer fills result/err and closes done to signal readiness.
// Thread safety: result and err are written before done is closed.
// The Go memory model guarantees that a channel close happens-before
// a receive on that channel completes, so readers see the writes.
type speculativeEntry struct {
	blockHash common.Hash
	once      sync.Once
	done      chan struct{}
	result    *SpeculativeResult // set before closing done
	err       error              // set before closing done
}

// Complete fills the entry's result and signals waiters.
// Safe to call multiple times — only the first call takes effect.
func (e *speculativeEntry) Complete(result *SpeculativeResult, err error) {
	e.once.Do(func() {
		e.result = result
		e.err = err
		close(e.done)
	})
}

// SpeculativeResultCache is a single-entry cache for speculative execution
// results. Within a consensus sequence, only the latest round's proposed
// block matters (round changes discard earlier proposals), so a single
// slot keyed by block hash is sufficient. See the analysis:
//
//   - Within a sequence: each round change overwrites the stale entry.
//   - Cross-sequence: speculative execution of block N+1 cannot start until
//     InsertChain(N) commits N's state, so there is no concurrent overlap.
//
// Thread safety: Reserve is called by the consensus goroutine;
// TryGet is called by the chain insertion goroutine. Both are
// synchronized by mu. The done channel provides the handoff between
// the producer (speculative executor) and the consumer (InsertChain).
type SpeculativeResultCache struct {
	mu    sync.Mutex
	entry *speculativeEntry
	hits  uint64
}

// NewSpeculativeResultCache creates a new single-entry speculative cache.
func NewSpeculativeResultCache() *SpeculativeResultCache {
	return &SpeculativeResultCache{}
}

// Reserve creates a new cache entry for the given block hash and returns it.
// The caller (speculative executor) must populate entry via Complete when
// execution finishes.
//
// If an existing entry is present (e.g., from a previous round), it is
// replaced. In-flight waiters on the old entry will time out gracefully
// via the timeout parameter of TryGet.
func (c *SpeculativeResultCache) Reserve(hash common.Hash) *speculativeEntry {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := &speculativeEntry{
		blockHash: hash,
		done:      make(chan struct{}),
	}
	c.entry = e
	return e
}

// TryGet attempts to retrieve a cached result for the given block hash.
//
// Returns non-nil on cache hit. Returns nil on:
//   - no entry exists
//   - hash mismatch (entry is for a different block)
//   - entry is in-flight and does not complete within timeout
//   - entry completed with an error
//
// Pass timeout <= 0 for a non-blocking check: returns the result only if
// execution already completed. This is the mode used by InsertChain to
// avoid blocking the chain mutex.
//
// On a successful hit, the entry is automatically cleared to free memory.
func (c *SpeculativeResultCache) TryGet(hash common.Hash, timeout time.Duration) *SpeculativeResult {
	c.mu.Lock()
	e := c.entry
	c.mu.Unlock()

	if e == nil || e.blockHash != hash {
		return nil
	}

	// Wait for the entry to complete or timeout.
	if timeout <= 0 {
		// Non-blocking: only succeed if already done.
		select {
		case <-e.done:
		default:
			return nil
		}
	} else {
		select {
		case <-e.done:
		case <-time.After(timeout):
			return nil
		}
	}

	if e.err != nil {
		return nil
	}

	c.mu.Lock()
	// Only clear if the entry hasn't been replaced by a new Reserve.
	if c.entry == e {
		c.entry = nil
	}
	c.hits++
	c.mu.Unlock()

	return e.result
}

// Clear removes the current entry. Safe to call when no entry exists.
func (c *SpeculativeResultCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = nil
}

// Hits returns the total number of successful cache hits.
func (c *SpeculativeResultCache) Hits() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits
}
