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

func (d *DialSched) Close() {
	for _, typedSched := range d.typedSchedulers {
		typedSched.Close()
	}
}

func (d *DialSched) AddStatic(n *discover.Node) {
	if tds := d.typedSchedulers[n.NType]; tds != nil {
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

func (td *TypedDialSched) Close() {
	if td.closed.Swap(true) {
		return
	}

	close(td.closeReq)
	td.wg.Wait()
}

func (td *TypedDialSched) AddStatic(n *discover.Node) {
	td.mu.Lock()
	defer td.mu.Unlock()

	td.static[n.ID] = n
}

func (td *TypedDialSched) RemoveStatic(id discover.NodeID) {
	td.mu.Lock()
	defer td.mu.Unlock()

	delete(td.static, id)
}

// Intake a dial success event.
func (td *TypedDialSched) OnSuccess(id discover.NodeID) {
	td.mu.Lock()
	defer td.mu.Unlock()

	td.connected[id] = true
	delete(td.dialing, id)
	delete(td.fails, id)
}

// Intake a dial failure or disconnection event.
func (td *TypedDialSched) OnFailure(id discover.NodeID) {
	td.mu.Lock()
	defer td.mu.Unlock()

	delete(td.connected, id)
	delete(td.dialing, id)
	td.fails[id]++
}

func (td *TypedDialSched) dialLoop() {
	defer td.wg.Done()
	resCh := make(chan error, td.targetConn)

	for {
		td.launchDialTasks(resCh)

		// Only back off when no dials are currently running.
		var idle <-chan time.Time
		td.mu.RLock()
		running := len(td.dialing)
		td.mu.RUnlock()
		if running == 0 {
			idle = time.After(10 * time.Second)
		}

		select {
		case <-resCh: // wait for a dial result
		case <-idle:
		case <-td.closeReq:
			return
		}
	}
}

// Launch dial tasks to reach the targetConn count.
func (td *TypedDialSched) launchDialTasks(resCh chan error) {
	td.mu.RLock()
	var (
		connected = len(td.connected)
		dialing   = len(td.dialing)
		want      = td.targetConn - connected - dialing
	)
	td.mu.RUnlock()
	if want <= 0 {
		return
	}

	launched := 0
	candidates := td.candidates(want)
	for _, n := range candidates {
		if launched >= want {
			break
		}
		td.mu.Lock()
		// Skip duplicates and already active/connected nodes.
		if td.connected[n.ID] || td.dialing[n.ID] {
			td.mu.Unlock()
			continue
		}
		td.dialing[n.ID] = true
		td.mu.Unlock()

		launched++
		nn := n
		go func() { resCh <- td.dialOnce(nn) }()
	}

	if launched > 0 {
		logger.Debug("DialSched launched", "target", td.targetConn, "connected", connected, "dialing", dialing, "want", want,
			"candidates", len(candidates), "launched", launched)
	}
}

// Returns the nodes to dial.
func (td *TypedDialSched) candidates(want int) []*discover.Node {
	if want <= 0 {
		return nil
	}

	// Overfetch because some candidates may be ineligible (already connected or dialing)
	td.mu.RLock()
	candidates := make([]*discover.Node, 0, len(td.static)+want*2)

	// Prioritize static nodes.
	for _, n := range td.static {
		candidates = append(candidates, n)
	}
	td.mu.RUnlock()

	// Add random nodes.
	random := make([]*discover.Node, want*2)
	if td.tab != nil {
		numRandom := td.tab.RandomNodes(random, td.targetType)
		candidates = append(candidates, random[:numRandom]...)
	}

	return candidates
}

func (td *TypedDialSched) dialOnce(n *discover.Node) error {
	// TODO: srv.SetupConn
	time.Sleep(time.Second)

	td.OnSuccess(n.ID) // TODO: it should be called inside srv.SetupConn. Calling here temporarily.
	logger.Debug("DialSched connected", "node", n.ID, "nType", n.NType)
	return nil
}
