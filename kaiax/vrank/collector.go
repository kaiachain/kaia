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

package vrank

import (
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
)

// ViewKey identifies a view (sequence N, round R) for collection.
type ViewKey struct {
	N uint64
	R uint8
}

// Cmp returns -1 if a < b, 0 if a == b, 1 if a > b (order: N then R).
func (a ViewKey) Cmp(b ViewKey) int {
	if a.N != b.N {
		if a.N < b.N {
			return -1
		}
		return 1
	}
	if a.R != b.R {
		if a.R < b.R {
			return -1
		}
		return 1
	}
	return 0
}

// CandidateMsg is a VRankCandidate plus when it was received. Caller uses this for verification and elapsed.
type CandidateMsg struct {
	ReceivedAt time.Time
	Msg        *VRankCandidate
}

// Collector stores VRankCandidate messages per view.
type Collector struct {
	mu             sync.RWMutex
	prepreparedMap map[ViewKey]time.Time
	blockHashMap   map[ViewKey]common.Hash // expected block hash per view (for validation at report time)
	viewMap        map[ViewKey]map[common.Address]CandidateMsg
}

// NewCollector creates a collector.
func NewCollector() *Collector {
	return &Collector{
		prepreparedMap: make(map[ViewKey]time.Time),
		blockHashMap:   make(map[ViewKey]common.Hash),
		viewMap:        make(map[ViewKey]map[common.Address]CandidateMsg),
	}
}

// RemoveOldViews deletes views that are strictly behind threshold.
func (c *Collector) RemoveOldViews(threshold ViewKey) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for vk := range c.prepreparedMap {
		if vk.Cmp(threshold) < 0 {
			delete(c.prepreparedMap, vk)
		}
	}
	for vk := range c.blockHashMap {
		if vk.Cmp(threshold) < 0 {
			delete(c.blockHashMap, vk)
		}
	}
	for vk := range c.viewMap {
		if vk.Cmp(threshold) < 0 {
			delete(c.viewMap, vk)
		}
	}
}

// AddPrepreparedTime records the start time and expected block hash for the view.
// expectedBlockHash is used at GetViewData/report time to validate VRankCandidate.BlockHash (reject liars).
func (c *Collector) AddPrepreparedTime(vk ViewKey, prepreparedAt time.Time, expectedBlockHash common.Hash) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prepreparedMap[vk] = prepreparedAt
	c.blockHashMap[vk] = expectedBlockHash
}

// AddCandMsg stores a VRankCandidate message for the given view. No verification is done here.
// Returns false if the sender already has a message stored for this view (duplicate).
func (c *Collector) AddCandMsg(vk ViewKey, sender common.Address, receivedAt time.Time, msg *VRankCandidate) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	candMap, ok := c.viewMap[vk]
	if !ok {
		c.viewMap[vk] = make(map[common.Address]CandidateMsg)
		candMap = c.viewMap[vk]
	}
	// no duplicate messages
	if _, ok := candMap[sender]; ok {
		return false
	}
	candMap[sender] = CandidateMsg{ReceivedAt: receivedAt, Msg: msg}
	return true
}

// GetViewData returns the raw data for the view: start time, expected block hash, and all stored messages.
// Caller should only count a message as valid/on-time if msg.Msg.BlockHash == expectedBlockHash.
func (c *Collector) GetViewData(vk ViewKey) (prepreparedAt time.Time, expectedBlockHash common.Hash, candMap map[common.Address]CandidateMsg) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	prepreparedAt = c.prepreparedMap[vk]
	expectedBlockHash = c.blockHashMap[vk]
	candMap = c.viewMap[vk]
	return prepreparedAt, expectedBlockHash, candMap
}
