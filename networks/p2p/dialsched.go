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
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

var (
	maxConcurrentDials = 16 // Maximum number of concurrent dial attempts.
	dialMaxRetries     = 3  // Maximum number of retries to dial the same node.

	dialBackoff      = 30 * time.Second // Minimum interval between the dial for the same node ID.
	refreshBackoff   = 4 * time.Second  // Minimum interval between discretionary discovery table refreshes.
	idleDialInterval = 10 * time.Second // Spins down the dial loop when there are no ongoing dial attempts.
)

// DialBackend is the minimal server-side API dialsched needs for outbound dial execution.
type DialBackend interface {
	SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error
}

// discovery is the minimal table API dialsched needs.
type discovery interface {
	GetNodes(targetType discover.NodeType, max int) []*discover.Node
	ReadRandomNodes(buf []*discover.Node, nType discover.NodeType) int
	Lookup(target discover.NodeID, targetType discover.NodeType) []*discover.Node
}

type DialConfig struct {
	selfID   discover.NodeID
	selfType discover.NodeType

	staticNodes []*discover.Node // initial set of static nodes.
	netrestrict *netutil.Netlist

	// target number of outbound connections for each node type.
	// Set connTargets[T] = 0 to disable dynamic dial of node type T (i.e. only accept inbound connections).
	// Set connTargets = nil to use the default value based on the selfType.
	// See also: discover.table.discoverTargets.
	connTarget map[discover.NodeType]int
}

// DialSched schedules outbound dial attempts to fulfill the desired amount of connections for each node type.
//
// The target nodes are either dynamic or static.
//   - Static nodes: Manually set by static-nodes.json or admin_addPeer.
//     Connections to the static nodes are always established, regardless of the connTargets or p2p.Server's peer quota.
//     Since we always attempt to dial static nodes, we have to limit infinite retries towards a dead static node.
//     If a static node fails many times, the node is removed from the `static` list.
//   - Dynamic nodes: Fetched from the discovery table. Learned from bootnodes and other peers over UDP discovery.
//     Discovery nodes are considered if static nodes are not enough to fulfill the connTargets.
//
// Note that NodeTypeUnknown is also a valid static node here. NodeTypeUnknown occurs when user-supplied KNI is missing the `?ntype=` parameter.
// In this case, the node's type is determined after the dialing and connection handshake.
type DialSched struct {
	mu         sync.RWMutex
	selfID     discover.NodeID
	connTarget map[discover.NodeType]int // Number of outbound connections to maintain for each node type
	dialer     NodeDialer
	backend    DialBackend

	// Node sources
	static typedNodeSet // The static nodes.
	tab    discovery    // The provider of dynamic nodes.

	// Controls the discretionary discovery table refresh.
	refreshBackoff time.Time // Earliest time allowed to refresh the discovery table.

	// Dialing nodes
	dialing     typedNodeSet
	dialBackoff map[discover.NodeID]time.Time // Earliest times allowed to retry dialing.
	netrestrict *netutil.Netlist

	// Connected nodes
	connectedAll      typedNodeSet            // All connected peers, inbound + outbound.
	connectedOutbound typedNodeSet            // Outbound connected peers only.
	connFails         map[discover.NodeID]int // For static nodes, count the consecutive connection failures to limit retries.

	// Coordination
	started  atomic.Bool
	closed   atomic.Bool
	closeReq chan struct{}
	wakeDial chan struct{}
	wg       sync.WaitGroup
}

func NewDialSched(cfg DialConfig, tab discovery, backend DialBackend) *DialSched {
	connTarget := cfg.connTarget
	if connTarget == nil {
		connTarget = defaultConnTarget(cfg.selfType)
	}

	ds := &DialSched{
		selfID:            cfg.selfID,
		connTarget:        connTarget,
		dialer:            TCPDialer{Dialer: &net.Dialer{Timeout: defaultDialTimeout}},
		backend:           backend,
		static:            newTypedNodeSet(),
		tab:               tab,
		dialing:           newTypedNodeSet(),
		dialBackoff:       make(map[discover.NodeID]time.Time),
		netrestrict:       cfg.netrestrict,
		connectedAll:      newTypedNodeSet(),
		connectedOutbound: newTypedNodeSet(),
		connFails:         make(map[discover.NodeID]int),
		closeReq:          make(chan struct{}),
		wakeDial:          make(chan struct{}, 1),
	}
	for _, n := range cfg.staticNodes {
		ds.addStatic(n)
	}

	return ds
}

func (ds *DialSched) Start() {
	if ds.closed.Load() || ds.started.Swap(true) {
		return
	}
	ds.wg.Add(1)
	go ds.dialLoop()
}

func (ds *DialSched) Close() {
	if ds.closed.Swap(true) {
		return
	}
	close(ds.closeReq)

	done := make(chan struct{})
	go func() {
		ds.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		logger.Warn("DialSched.Close: timed out waiting for dial tasks to finish")
	}
}

func (ds *DialSched) AddStatic(n *discover.Node) {
	ds.addStatic(n)
	ds.signalDial()
}

func (ds *DialSched) RemoveStatic(id discover.NodeID) {
	ds.removeStatic(id)
	ds.signalDial()
}

// OnPeerConnected updates connected bookkeeping using the actual peer type observed at handshake.
// Inbound peers count for deduplication but not for outbound connTarget accounting.
func (ds *DialSched) OnPeerConnected(n *discover.Node, inbound bool) {
	ds.markConnSuccess(n, inbound)
	ds.signalDial()
}

func (ds *DialSched) OnPeerDisconnected(id discover.NodeID, nType discover.NodeType) {
	ds.markConnDisconnect(id)
	ds.signalDial()
}

// Let the dialLoop know that the situation has changed and it should check again.
func (ds *DialSched) signalDial() {
	select {
	case ds.wakeDial <- struct{}{}:
	default:
	}
}

// Best-effort completion signal for dialLoop.
// We intentionally drop redundant signals when full or closing.
func (ds *DialSched) signalResult(resCh chan struct{}) {
	select {
	case resCh <- struct{}{}:
	case <-ds.closeReq:
	default:
	}
}

//// Dial task loop

func (ds *DialSched) dialLoop() {
	defer ds.wg.Done()
	resCh := make(chan struct{}, maxConcurrentDials)

	for {
		candidates, needRefresh := ds.getCandidates()
		ds.launchDialTasks(candidates, resCh)
		if needRefresh && ds.shouldRefresh() {
			ds.refreshOnce(resCh)
		}

		var idle <-chan time.Time
		if ds.dialing.len() == 0 {
			idle = time.After(idleDialInterval)
		}

		select { // Loop upon any of the following signals.
		case <-ds.wakeDial: // situation changed
		case <-resCh: // one dial task completed (either success or failure)
		case <-idle: // done idling
		case <-ds.closeReq:
			return
		}
	}
}

// Return the candidates for dialing from the mixture of static and dynamic nodes.
func (ds *DialSched) getCandidates() (candidates []*discover.Node, needRefresh bool) {
	ds.mu.RLock()

	// 1. Determine the demand
	needs := make(map[discover.NodeType]int)
	for targetType := range ds.connTarget {
		needs[targetType] = ds.connTarget[targetType] - ds.connectedOutbound.count(targetType) - ds.dialing.count(targetType)
	}

	// 2. Add all static nodes, even if they are already connected or in the dialing list. Those will be filtered out later.
	// Since static nodes are exempt from the connTargets nor peer limits, all static nodes are going to be dialed.
	candidates = ds.static.all()
	ds.mu.RUnlock()

	// 3. Dynamic nodes
	for targetType, need := range needs {
		dynamicNodes, refresh := ds.getDynamicCandidates(targetType, need)
		candidates = append(candidates, dynamicNodes...)
		needRefresh = needRefresh || refresh
	}

	logger.Debug("Dialing candidates", "needs", needs, "candidates", len(candidates))
	return candidates, needRefresh
}

func (ds *DialSched) getDynamicCandidates(targetType discover.NodeType, need int) (nodes []*discover.Node, needRefresh bool) {
	if ds.tab == nil || need <= 0 {
		return nil, false
	}
	oversample := need * 4 // Oversample in case some nodes are not reachable.
	switch targetType {    // TODO: to be unified with tab.RandomNodes() and tab.Refresh()
	case discover.NodeTypeCN, discover.NodeTypePN:
		nodes = ds.tab.GetNodes(targetType, oversample)
		return nodes, len(nodes) < need
	case discover.NodeTypeEN:
		random := make([]*discover.Node, oversample)
		num := ds.tab.ReadRandomNodes(random, targetType)
		nodes := random[:min(num, len(random))]
		return nodes, len(nodes) < need
	default:
		return nil, false
	}
}

func (ds *DialSched) launchDialTasks(candidates []*discover.Node, resCh chan struct{}) {
	for _, n := range candidates {
		if !ds.shouldDial(n) {
			continue
		}
		ds.markDialingStart(n)
		nn := cloneNode(n)
		ds.wg.Add(1)
		go func() {
			defer ds.wg.Done()
			ds.dialOnce(nn)
			ds.signalResult(resCh)
		}()
	}
}

func (ds *DialSched) dialOnce(n *discover.Node) {
	defer ds.markDialingEnd(n.ID)

	if ds.dialer == nil || ds.backend == nil {
		ds.markConnFailure(n.ID)
		return
	}

	flags := dynDialedConn
	if ds.static.contains(n.ID) {
		flags = staticDialedConn
	}

	var err error
	if len(n.TCPs) > 1 {
		err = ds.dialMulti(n, flags)
	} else {
		err = ds.dial(n, flags)
	}

	// Multi-channel peers can request redial after updating listen ports.
	if err == errUpdateDial {
		err = ds.dialMulti(n, flags)
	}
	logger.Debug("Dialed node", "id", n.ID, "type", n.NType, "addresses", n.TCPs, "err", err)

	if err != nil {
		ds.markConnFailure(n.ID)
	}
	// We don't call markConnSuccess() here. If dial()-SetupConn() succeeds, the Server will call OnPeerConnected()-markConnSuccess().
}

func (ds *DialSched) dial(dest *discover.Node, flags connFlag) error {
	fd, err := ds.dialer.Dial(dest)
	if err != nil {
		return err
	}
	mfd := newMeteredConn(fd, false)
	return ds.backend.SetupConn(mfd, flags, dest)
}

func (ds *DialSched) dialMulti(dest *discover.Node, flags connFlag) error {
	fds, err := ds.dialer.DialMulti(dest)
	if err != nil {
		return err
	}

	var errBackup error
	for portOrder, fd := range fds {
		mfd := newMeteredConn(fd, false)
		destByPort := cloneNode(dest)
		destByPort.PortOrder = uint16(portOrder)
		if err := ds.backend.SetupConn(mfd, flags, destByPort); err != nil {
			errBackup = err
		}
	}
	if errBackup != nil {
		for _, fd := range fds {
			fd.Close()
		}
	}
	return errBackup
}

func (ds *DialSched) refreshOnce(resCh chan struct{}) {
	ds.mu.Lock()
	ds.refreshBackoff = time.Now().Add(refreshBackoff)
	ds.mu.Unlock()

	if ds.tab == nil {
		return
	}
	ds.wg.Add(1)
	go func() {
		defer ds.wg.Done()
		ds.tab.Lookup(ds.selfID, discover.NodeTypeCN)
		ds.tab.Lookup(ds.selfID, discover.NodeTypePN)
		ds.tab.Lookup(ds.selfID, discover.NodeTypeEN)
		ds.signalResult(resCh)
	}()
}

//// Internal variable accessors

func (ds *DialSched) addStatic(n *discover.Node) {
	if n == nil {
		return
	}
	if n.Incomplete() {
		logger.Warn("Rejecting incomplete static node", "id", n.ID, "type", n.NType)
		return
	}
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.static.add(cloneNode(n))
}

func (ds *DialSched) removeStatic(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.static.remove(id)

	// Additionally cleanup the dial speed bumps.
	// This allows admin_removePeer -> admin_addPeer combo to force an immediate reconnect.
	delete(ds.dialBackoff, id)
	delete(ds.connFails, id)
}

func (ds *DialSched) shouldRefresh() bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.refreshBackoff.Before(time.Now())
}

func (ds *DialSched) shouldDial(n *discover.Node) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if n == nil {
		return false
	}
	if n.Incomplete() {
		logger.Warn("Rejecting incomplete dial candidate", "id", n.ID, "type", n.NType)
		return false
	}
	if ds.dialing.len() >= maxConcurrentDials {
		return false
	}
	if n.ID == ds.selfID {
		return false
	}
	if ds.netrestrict != nil && !ds.netrestrict.Contains(n.IP) {
		return false
	}
	if ds.connectedAll.contains(n.ID) {
		return false
	}
	if ds.dialing.contains(n.ID) {
		return false
	}
	if until, ok := ds.dialBackoff[n.ID]; ok && time.Now().Before(until) {
		return false
	}
	return true
}

func (ds *DialSched) markDialingStart(n *discover.Node) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.dialing.add(n)
	ds.dialBackoff[n.ID] = time.Now().Add(dialBackoff)
}

func (ds *DialSched) markDialingEnd(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.dialing.remove(id)
}

func (ds *DialSched) markConnSuccess(n *discover.Node, inbound bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if n == nil {
		return
	}
	// Retype existing entry (if any) to fix the source type vs. connected type discrepancy.
	ds.connectedAll.remove(n.ID)
	ds.connectedAll.add(n)

	ds.connectedOutbound.remove(n.ID)
	if !inbound {
		ds.connectedOutbound.add(n)
	}

	delete(ds.connFails, n.ID)
	delete(ds.dialBackoff, n.ID)
}

func (ds *DialSched) markConnFailure(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.connectedAll.remove(id)
	ds.connectedOutbound.remove(id)
	if ds.static.contains(id) {
		ds.connFails[id]++
		if ds.connFails[id] > dialMaxRetries {
			ds.static.remove(id)
		}
	}
}

func (ds *DialSched) markConnDisconnect(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.connectedAll.remove(id)
	ds.connectedOutbound.remove(id)
}

func cloneNode(n *discover.Node) *discover.Node {
	if n == nil {
		return nil
	}
	nn := *n
	nn.TCPs = append([]uint16(nil), n.TCPs...)
	return &nn
}
