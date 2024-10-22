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
	"sync"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
)

type refCountedState struct {
	state *state.StateDB
	refs  map[uint64]struct{}
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

func (ss *sideStates) AddState(refId uint64, statedb *state.StateDB) {
	root := statedb.IntermediateRoot(false)

	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.states[root] == nil {
		ss.states[root] = &refCountedState{
			state: statedb.Copy(),
			refs:  make(map[uint64]struct{}),
		}
	}
	ss.states[root].refs[refId] = struct{}{}
}

func (ss *sideStates) GetState(root common.Hash) *state.StateDB {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if ss.states[root] == nil {
		return nil
	}
	return ss.states[root].state
}

func (s *StakingModule) AllocSideStateRef() uint64 {
	return s.sideStates.AllocRefId()
}

func (s *StakingModule) FreeSideStateRef(refId uint64) {
	s.sideStates.FreeRefId(refId)
}

func (s *StakingModule) AddSideState(refId uint64, statedb *state.StateDB) {
	s.sideStates.AddState(refId, statedb)
}
