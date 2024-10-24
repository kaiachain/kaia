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

type refCountedState struct {
	info *staking.StakingInfo
	refs map[uint64]struct{}
}

// sideStates remembers temporary StateDBs for StakingModule to refer to.
// The temporary states are created during state reexec(regen) when the states are not available in the database.
type sideStates struct {
	states    map[common.Hash]*refCountedState // keyed by state root
	nextRefId uint64                           // uint64 should be enough during the node process runtime
	mu        sync.RWMutex
}

func NewSideStates() *sideStates {
	return &sideStates{
		states:    make(map[common.Hash]*refCountedState),
		nextRefId: 1,
	}
}

func (ss *sideStates) AllocRefId() uint64 {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	id := ss.nextRefId
	ss.nextRefId++
	return id
}

func (ss *sideStates) FreeRefId(refId uint64) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for root, state := range ss.states {
		delete(state.refs, refId) // delete reference
		if len(state.refs) == 0 { // no more references
			delete(ss.states, root)
		}
	}
}

func (ss *sideStates) GetInfo(root common.Hash) *staking.StakingInfo {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if ss.states[root] == nil {
		return nil
	}
	return ss.states[root].info
}

func (s *StakingModule) AllocSideStateRef() uint64 {
	return s.sideStates.AllocRefId()
}

func (s *StakingModule) FreeSideStateRef(refId uint64) {
	s.sideStates.FreeRefId(refId)
}

func (s *StakingModule) AddSideState(refId uint64, header *types.Header, statedb *state.StateDB) error {
	ss := s.sideStates
	root := statedb.IntermediateRoot(false)
	// Sanity check
	if header.Root != root {
		return fmt.Errorf("header root mismatch: %s != %s", header.Root, root)
	}

	// Quickly check if the info is already stored.
	ss.mu.RLock()
	_, ok := ss.states[root]
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

	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Check again, if info is still not stored, store it.
	if ss.states[root] == nil {
		ss.states[root] = &refCountedState{
			info: info,
			refs: make(map[uint64]struct{}),
		}
	}
	ss.states[root].refs[refId] = struct{}{}
	return nil
}
