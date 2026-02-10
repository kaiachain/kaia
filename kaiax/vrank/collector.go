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

// ViewKey identifies a (sequence, round) for collection.
type ViewKey struct {
	N uint64
	R uint8
}

// Cmp returns -1 if a < b, 0 if a == b, 1 if a > b (order: Seq then Round).
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
	viewMap        map[ViewKey]map[common.Address]CandidateMsg
}

// NewCollector creates a collector. maxWindow limits how far ahead of currentView messages are accepted (0 = unbounded).
func NewCollector() *Collector {
	return &Collector{
		prepreparedMap: make(map[ViewKey]time.Time),
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
	for vk := range c.viewMap {
		if vk.Cmp(threshold) < 0 {
			delete(c.viewMap, vk)
		}
	}
}

func (c *Collector) AddPrepreparedTime(vk ViewKey, prepreparedAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prepreparedMap[vk] = prepreparedAt
}

// AddCandMsg stores a VRankCandidate message for the given view. No verification is done here.
// Returns false if msg is nil or vk is too far (behind currentView or beyond currentView+maxWindow).
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

// GetViewData returns the raw data for the view: start time and all stored messages.
func (c *Collector) GetViewData(vk ViewKey) (time.Time, map[common.Address]CandidateMsg) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	prepreparedAt := c.prepreparedMap[vk]
	candMap := c.viewMap[vk]
	return prepreparedAt, candMap
}
