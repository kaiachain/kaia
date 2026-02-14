// Modifications Copyright 2026 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/networks/p2p/discover"
)

type DialConfig struct {
	selfID   discover.NodeID
	selfType discover.NodeType

	staticNodes []*discover.Node // initial set of static nodes.

	// target number of outbound connections for each node type.
	// Set connTargets[T] = 0 to disable dynamic dial of node type T (i.e. only accept inbound connections).
	// Set connTargets = nil to use the default value based on the selfType.
	// See also: discover.table.discoverTargets.
	connTargets map[discover.NodeType]int
}

type DialSched struct {
	typedSchedulers map[discover.NodeType]*TypedDialSched
}

func NewDialSched(cfg DialConfig, tab discover.Discovery2) *DialSched {
	typedSchedulers := make(map[discover.NodeType]*TypedDialSched)
	connTargets := cfg.connTargets
	if connTargets == nil {
		connTargets = getConnTargets(cfg.selfType)
	}
	for targetType, targetConn := range connTargets {
		if targetConn > 0 {
			typedSchedulers[targetType] = newTypedDialSched(cfg.selfID, targetType, targetConn, tab)
		}
	}
	ds := &DialSched{
		typedSchedulers: typedSchedulers,
	}

	for _, n := range cfg.staticNodes {
		ds.AddStatic(n)
	}

	return ds
}

func (ds *DialSched) Close() {
	for _, typedSched := range ds.typedSchedulers {
		typedSched.Close()
	}
}

func (ds *DialSched) AddStatic(n *discover.Node) {
	if tds := ds.typedSchedulers[n.NType]; tds != nil {
		tds.AddStatic(n)
	}
}

func getConnTargets(selfType discover.NodeType) map[discover.NodeType]int {
	switch selfType {
	case discover.NodeTypeCN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 100,
			discover.NodeTypePN: 0,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	case discover.NodeTypePN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 1,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	case discover.NodeTypeEN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 2,
			discover.NodeTypeEN: 3,
			discover.NodeTypeBN: 0,
		}
	default:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 0,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	}
}

type TypedDialSched struct {
	selfID discover.NodeID
	tab    discover.Discovery2

	targetType discover.NodeType
	targetConn int

	mu        sync.RWMutex
	static    map[discover.NodeID]*discover.Node
	fails     map[discover.NodeID]int // dial or connection failures of each static node.
	dialing   map[discover.NodeID]bool
	connected map[discover.NodeID]bool

	wg       sync.WaitGroup
	closed   atomic.Bool
	closeReq chan struct{}
}

func newTypedDialSched(selfID discover.NodeID, targetType discover.NodeType, targetConn int, tab discover.Discovery2) *TypedDialSched {
	return &TypedDialSched{
		selfID:     selfID,
		targetType: targetType,
		targetConn: targetConn,
		tab:        tab,

		static:    make(map[discover.NodeID]*discover.Node),
		dialing:   make(map[discover.NodeID]bool),
		connected: make(map[discover.NodeID]bool),
		fails:     make(map[discover.NodeID]int),

		closeReq: make(chan struct{}),
	}
}

func (ds *TypedDialSched) Close() {
	if ds.closed.Swap(true) {
		return
	}

	close(ds.closeReq)
	ds.wg.Wait()
}

func (ds *TypedDialSched) AddStatic(n *discover.Node) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.static[n.ID] = n
}

func (ds *TypedDialSched) RemoveStatic(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.static, id)
}

// Intake a dial success event.
func (ds *TypedDialSched) OnSuccess(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.connected[id] = true
	delete(ds.dialing, id)
	delete(ds.fails, id)
}

// Intake a dial failure or disconnection event.
func (ds *TypedDialSched) OnFailure(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.connected, id)
	delete(ds.dialing, id)
	ds.fails[id]++
}

func (ds *TypedDialSched) dialLoop() {
	defer ds.wg.Done()
	resCh := make(chan error, ds.targetConn)

	for {
		ds.launchDialTasks(resCh)

		// Only back off when no dials are currently running.
		var idle <-chan time.Time
		ds.mu.RLock()
		running := len(ds.dialing)
		ds.mu.RUnlock()
		if running == 0 {
			idle = time.After(10 * time.Second)
		}

		select {
		case <-resCh: // wait for a dial result
		case <-idle:
		case <-ds.closeReq:
			return
		}
	}
}

// Launch dial tasks to reach the targetConn count.
func (ds *TypedDialSched) launchDialTasks(resCh chan error) {
	ds.mu.RLock()
	var (
		connected = len(ds.connected)
		dialing   = len(ds.dialing)
		want      = ds.targetConn - connected - dialing
	)
	ds.mu.RUnlock()
	if want <= 0 {
		return
	}

	launched := 0
	candidates := ds.candidates(want)
	for _, n := range candidates {
		if launched >= want {
			break
		}
		ds.mu.Lock()
		// Skip duplicates and already active/connected nodes.
		if ds.connected[n.ID] || ds.dialing[n.ID] {
			ds.mu.Unlock()
			continue
		}
		ds.dialing[n.ID] = true
		ds.mu.Unlock()

		launched++
		nn := n
		go func() { resCh <- ds.dialOnce(nn) }()
	}

	if launched > 0 {
		logger.Debug("DialSched launched", "target", ds.targetConn, "connected", connected, "dialing", dialing, "want", want,
			"candidates", len(candidates), "launched", launched)
	}
}

// Returns the nodes to dial.
func (ds *TypedDialSched) candidates(want int) []*discover.Node {
	if want <= 0 {
		return nil
	}

	// Overfetch because some candidates may be ineligible (already connected or dialing)
	ds.mu.RLock()
	candidates := make([]*discover.Node, 0, len(ds.static)+want*2)

	// Prioritize static nodes.
	for _, n := range ds.static {
		candidates = append(candidates, n)
	}
	ds.mu.RUnlock()

	// Add random nodes.
	random := make([]*discover.Node, want*2)
	if ds.tab != nil {
		numRandom := ds.tab.RandomNodes(random, ds.targetType)
		candidates = append(candidates, random[:numRandom]...)
	}

	return candidates
}

func (ds *TypedDialSched) dialOnce(n *discover.Node) error {
	// TODO: srv.SetupConn
	time.Sleep(time.Second)

	ds.OnSuccess(n.ID) // TODO: it should be called inside srv.SetupConn. Calling here temporarily.
	logger.Debug("DialSched connected", "node", n.ID, "nType", n.NType)
	return nil
}
