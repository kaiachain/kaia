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

package impl

import (
	"fmt"
	"sync"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
)

type refCountedInfo struct {
	info *staking.StakingInfo
	refs map[uint64]struct{}
}

// PreloadBuffer remembers temporary StakingInfo for StakingModule to refer to.
// The temporary info are created during state reexec(regen) when the states are not available in the database.
type PreloadBuffer struct {
	preloaded map[common.Hash]*refCountedInfo // keyed by state root
	nextRefId uint64                          // uint64 should be enough during the node process runtime
	mu        sync.RWMutex
}

func NewPreloadBuffer() *PreloadBuffer {
	return &PreloadBuffer{
		preloaded: make(map[common.Hash]*refCountedInfo),
		nextRefId: 1,
	}
}

func (ss *PreloadBuffer) AllocRefId() uint64 {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	id := ss.nextRefId
	ss.nextRefId++
	return id
}

func (ss *PreloadBuffer) FreeRefId(refId uint64) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for root, state := range ss.preloaded {
		delete(state.refs, refId) // delete reference
		if len(state.refs) == 0 { // no more references
			delete(ss.preloaded, root)
		}
	}
}

func (ss *PreloadBuffer) GetInfo(root common.Hash) *staking.StakingInfo {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if ss.preloaded[root] == nil {
		return nil
	}
	return ss.preloaded[root].info
}

func (s *StakingModule) AllocPreloadRef() uint64 {
	return s.preloadBuffer.AllocRefId()
}

func (s *StakingModule) FreePreloadRef(refId uint64) {
	s.preloadBuffer.FreeRefId(refId)
}

func (s *StakingModule) PreloadFromState(refId uint64, header *types.Header, statedb *state.StateDB) error {
	ss := s.preloadBuffer
	root := statedb.IntermediateRoot(false)
	// Sanity check
	if header.Root != root {
		return fmt.Errorf("header root mismatch: %s != %s", header.Root, root)
	}

	// Quickly check if the info is already stored.
	ss.mu.RLock()
	_, ok := ss.preloaded[root]
	ss.mu.RUnlock()
	if ok {
		return nil
	}

	// Calculate staking info from the state.
	// Do not lock here because it may take time.
	info, err := s.getFromState(header, statedb)
	if err != nil {
		return err
	}

	// Check again, if info is still not stored, store it.
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.preloaded[root] == nil {
		ss.preloaded[root] = &refCountedInfo{
			info: info,
			refs: make(map[uint64]struct{}),
		}
	}
	ss.preloaded[root].refs[refId] = struct{}{}
	return nil
}
