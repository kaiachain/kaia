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
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

var (
	maxConcurrentDials = 16 // Maximum number of concurrent dial attempts.
	dialMaxRetries     = 3  // Maximum number of retries to dial the same node.

	dialBackoff      = 30 * time.Second // Minimum interval between the dial for the same node ID.
	refreshBackoff   = 4 * time.Second  // Minimum interval between discretionary discovery table refreshes.
	idleDialInterval = 10 * time.Second // Spins down the dial loop when there are no ongoing dial attempts.

	candidateQueueSize = 100 // Number of nodes to fetch per dynamic candidate queue refill.
)

// DialBackend is the minimal server-side API dialsched needs for outbound dial execution.
type DialBackend interface {
	SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error
}

// discovery is the minimal table API dialsched needs.
type discovery interface {
	RandomNodes(buf []*discover.Node, nType discover.NodeType) int
	Refresh()
	Close()
}

type DialConfig struct {
	selfID   discover.NodeID
	selfType discover.NodeType

	staticNodes []*discover.Node // initial set of static nodes.
	netrestrict *netutil.Netlist

	// target number of outbound connections for each node type.
	// Set connTargets[T] = 0 to disable dynamic dial of node type T (i.e. only accept inbound connections).
	// Set connTargets = nil to use the KIP-311 defaults based on the selfType.
	// See also: discover.table.discoverTargets.
	connTarget map[discover.NodeType]int

	// maxENToENDialTarget is the dynamic-dial EN cap derived from
	// MaxPhysicalConnections / DialRatio. Pass Server.maxDialedConns() here.
	maxENToENDialTarget *int

	// maxPeers is the total peer capacity (inbound + outbound).
	// Dynamic dials are suppressed when connectedAll reaches this limit,
	// mirroring the server-side check in encHandshakeChecks.
	// Pass srv.Config.MaxPhysicalConnections here. 0 means no limit.
	maxPeers int

	// Hooks for testing.
	dialer NodeDialer
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
//
// All mutable fields are protected by the mutex `mu`.
// DialSched responds to the events from the Server, or make its own changes in the dialLoop().
// In both cases, any data modifications must go through the accessors with the `mu` held.
// The srv.lock - ds.mu order is always maintained, so there is no deadlock risk.
//
// [ Server layer ]                     [ dialLoop layer ]
//   - holds srv.lock                     - does not hold mu at the loop level (but subroutines do)
//   - AddPeer, RemovePeer                - getCandidates, launchDialTasks, go dialOnce, go refreshOnce
//     handleAddPeerConn, handleDelPeer
//
// [ Internal critical sections ]
//   - holds ds.mu
//   - addStatic, removeStatic, shouldRefresh, markRefreshStart, shouldDial, isDialing,
//     markDialStart, markDialEnd, markPeerConnected, markDialFailure, markPeerDisconnected
type DialSched struct {
	mu         sync.RWMutex
	selfID     discover.NodeID
	selfType   discover.NodeType
	connTarget map[discover.NodeType]int // Number of outbound connections to maintain for each node type
	maxPeers   int                       // Total peer capacity (inbound + outbound). 0 = no limit.
	dialer     NodeDialer
	backend    DialBackend

	// Node sources
	static typedNodeSet // The static nodes.
	tab    discovery    // The provider of dynamic nodes.
	// candidateQueue is a per-type queue of discovery candidates, drained sequentially across loop
	// iterations. This avoids re-offering the same unreachable nodes before their dialBackoff expires.
	candidateQueue map[discover.NodeType][]*discover.Node

	// Controls the discretionary discovery table refresh.
	refreshBackoff time.Time // Earliest time allowed to refresh the discovery table. Counted since the refresh started.

	// Dialing nodes
	dialing     typedNodeSet
	dialBackoff map[discover.NodeID]time.Time // Earliest times allowed to retry dialing. Counted since the dial ended.
	netrestrict *netutil.Netlist

	// Connected nodes
	connectedAll      typedNodeSet            // All connected peers, inbound + outbound.
	connectedOutbound typedNodeSet            // Outbound connected peers only.
	connFails         map[discover.NodeID]int // For static nodes, count the consecutive connection failures to limit retries.
	cnPeerAddrs       map[common.Address]struct{}

	// Coordination
	started  atomic.Bool
	closed   atomic.Bool
	closeReq chan struct{}
	wakeDial chan struct{}
	wg       sync.WaitGroup
}

// #region public interfaces

func NewDialSched(cfg DialConfig, tab discovery, backend DialBackend) *DialSched {
	connTarget := cfg.connTarget
	if connTarget == nil {
		connTarget = defaultConnTarget(cfg)
	}

	dialer := cfg.dialer
	if dialer == nil {
		dialer = TCPDialer{Dialer: &net.Dialer{Timeout: defaultDialTimeout}}
	}

	ds := &DialSched{
		selfID:            cfg.selfID,
		selfType:          discover.EffectiveNodeType(cfg.selfType),
		connTarget:        connTarget,
		maxPeers:          cfg.maxPeers,
		dialer:            dialer,
		backend:           backend,
		static:            newTypedNodeSet(),
		tab:               tab,
		candidateQueue:    make(map[discover.NodeType][]*discover.Node),
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

// SetCNPeers replaces the CN dynamic dial allowlist.
// nil disables filtering; an empty list rejects all dynamic CN candidates.
func (ds *DialSched) SetCNPeers(addrs []common.Address) {
	ds.mu.Lock()
	if addrs == nil {
		ds.cnPeerAddrs = nil
	} else {
		cnPeerAddrs := make(map[common.Address]struct{}, len(addrs))
		for _, addr := range addrs {
			cnPeerAddrs[addr] = struct{}{}
		}
		ds.cnPeerAddrs = cnPeerAddrs
	}
	ds.mu.Unlock()

	ds.signalDial()
}

// OnPeerConnected updates the `connected*` bookkeeping based on the actual peer type observed at handshake.
// If the peer is inbound (i.e. via Server.listenLoop()), it does not affect the `connTarget`.
// If the peer is outbound (i.e. via DialSched.dialOnce()), it counts towards the `connTarget` accounting.
func (ds *DialSched) OnPeerConnected(id discover.NodeID, nType discover.NodeType, inbound bool) {
	ds.markPeerConnected(id, nType, inbound)
	ds.signalDial()
}

// OnPeerDisconnected updates the `connected*` bookkeeping when it was disconnected.
func (ds *DialSched) OnPeerDisconnected(id discover.NodeID, nType discover.NodeType) {
	ds.markPeerDisconnected(id)
	ds.signalDial()
}

// #region dial loop

// Let the dialLoop know that the situation has changed and it should check again.
// This operation is non-blocking. Silently skips when ds.wakeDial is full or closed.
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

func (ds *DialSched) dialLoop() {
	defer ds.wg.Done()
	dialResCh := make(chan struct{}, maxConcurrentDials)
	refreshResCh := make(chan struct{}, 1)

	for {
		candidates, needRefresh := ds.getCandidates()
		ds.launchDialTasks(candidates, dialResCh)
		if needRefresh && ds.shouldRefresh() {
			ds.refreshOnce(refreshResCh)
		}

		var idle <-chan time.Time
		if !ds.isDialing() { // There is no ongoing dial attempt. Come back after idling.
			idle = time.After(idleDialInterval)
		}

		select { // Loop upon any of the following signals.
		case <-ds.wakeDial: // situation changed
		case <-dialResCh: // one dial task completed (either success or failure)
		case <-refreshResCh: // discovery refresh completed; next getCandidates will
		case <-idle: // done idling
		case <-ds.closeReq:
			return
		}
	}
}

// #region candidate collection

// Return the candidates for dialing from the mixture of static and dynamic nodes.
//
// Capacity math tells us that the dialing is unlikely to be saturated, so we don't separate the candidates by node type.
//
// The candidates are composed of [ static | dynamic CNs | dynamic PNs | dynamic ENs ], and they are consumed from the left.
// - The candidates are filtered with shouldDial() and then dialed up to maxConcurrentDials (16) at once.
// - Even if some candidates are unfortunately unreachable, they are skipped in later iterations because of the dialBackoff and connFails.
//
// If there are too many static nodes, the dynamic candidate dialing can be delayed. But it is reasonable and expected.
// - You don't need dynamic candidates when you have enough static nodes.
//
// If there are too many dynamic CN candidates, EN dialing can be delayed. But it is unlikely under the connTarget settings.
// - CN's connTarget = {CN: 100}. No starvation possible.
// - EN's connTarget = {CN: 2, PN: 2, EN: maxENToENDialTarget}.
// PN/EN dialing is delayed only if static nodes and CN candidates fill maxConcurrentDials first.
// But in reality, there are much less static nodes, and even so they will be filtered out by shouldDial() in later iterations quite soon.
func (ds *DialSched) getCandidates() (candidates []*discover.Node, needRefresh bool) {
	ds.mu.RLock()

	// 1. Determine the demand
	needs := make(map[discover.NodeType]int)
	for targetType := range ds.connTarget {
		// Note that the needs can be negative because we often over-dial the dynamic nodes in case of failures.
		needs[targetType] = ds.connTarget[targetType] - ds.connectedOutbound.count(targetType) - ds.dialing.count(targetType)
	}

	// 2. Add all static nodes, even if they are already connected or in the dialing list. Those will be filtered out later.
	// Since static nodes are exempt from the connTargets nor peer limits, all static nodes are going to be dialed even if needs[] is zero.
	candidates = ds.static.all()
	ds.mu.RUnlock()

	// 3. Dynamic nodes — skip when already at total peer capacity.
	// Even if we want more outbound peers (needs > 0), the total capacity might be reached (maxPeers) due to inbound peers.
	// Static dials (added above) bypass this check, mirroring the Server's capacity check.
	if ds.maxPeers == 0 || ds.connectedAll.len() < ds.maxPeers {
		for targetType, need := range needs {
			dynamicNodes, refresh := ds.getDynamicCandidates(targetType, need)
			candidates = append(candidates, dynamicNodes...)
			needRefresh = needRefresh || refresh
		}
	}

	if len(candidates) > 0 {
		logger.Debug("Dialing candidates", "needs", needs, "static", ds.static.len(),
			"connectedOutbound", ds.connectedOutbound.len(), "dialing", ds.dialing.len(), "candidates", len(candidates))
	}
	return candidates, needRefresh
}

func (ds *DialSched) getDynamicCandidates(targetType discover.NodeType, need int) (nodes []*discover.Node, needRefresh bool) {
	if ds.tab == nil || need <= 0 {
		return nil, false
	}
	// Try to refill from the table if the queue can't satisfy the demand.
	// The table may already have nodes (e.g. from local nodeDB) without needing a refresh
	if len(ds.candidateQueue[targetType]) < need {
		ds.refillCandidates()
	}

	// Disallowed CNs are discarded so they do not keep occupying the dynamic candidate queue.
	q := ds.candidateQueue[targetType]
	q = slices.DeleteFunc(q, func(n *discover.Node) bool {
		return !ds.dynamicCandidateAllowed(targetType, n)
	})

	// Pop up to `need` candidates from the filtered queue.
	end := min(need, len(q))
	nodes, ds.candidateQueue[targetType] = q[:end], q[end:]
	// If the queue is short even after a refill, signal the dial loop to run a discovery refresh.
	return nodes, len(nodes) < need
}

// #region dial task

func (ds *DialSched) launchDialTasks(candidates []*discover.Node, resCh chan struct{}) {
	for _, n := range candidates {
		if !ds.shouldDial(n) {
			continue
		}
		flags := ds.markDialStart(n)
		nn := cloneNode(n)
		go func() {
			ds.dialOnce(nn, flags)
			ds.signalResult(resCh)
		}()
	}
}

func (ds *DialSched) dialOnce(n *discover.Node, flags connFlag) {
	defer ds.markDialEnd(n.ID)

	if ds.dialer == nil || ds.backend == nil {
		ds.markDialFailure(n.ID)
		return
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
		ds.markDialFailure(n.ID)
	}
	// We don't call markPeerConnected() here. If dial()-SetupConn() succeeds, the Server will call OnPeerConnected()-markPeerConnected().
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

// #region refresh task

func (ds *DialSched) refreshOnce(refreshResCh chan struct{}) {
	ds.markRefreshStart()

	if ds.tab == nil {
		return
	}
	go func() {
		ds.tab.Refresh()
		ds.signalResult(refreshResCh)
	}()
}

// refillCandidates refills empty per-type candidate queues from the discovery table.
// Called by dialLoop after a refresh completes, so the table has fresh nodes.
// Accessed only by the dialLoop goroutine; no mutex required.
func (ds *DialSched) refillCandidates() {
	for targetType := range ds.connTarget {
		if len(ds.candidateQueue[targetType]) == 0 {
			buf := make([]*discover.Node, candidateQueueSize)
			n := ds.tab.RandomNodes(buf, targetType)
			ds.candidateQueue[targetType] = buf[:n]
		}
	}
}

func (ds *DialSched) dynamicCandidateAllowed(targetType discover.NodeType, n *discover.Node) bool {
	if ds.selfType != discover.NodeTypeCN || discover.EffectiveNodeType(targetType) != discover.NodeTypeCN {
		return true
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if ds.cnPeerAddrs == nil {
		return true
	}

	addr, err := addressFromNodeID(n.ID)
	if err != nil {
		return false
	}
	_, ok := ds.cnPeerAddrs[addr]
	return ok
}

// #region internal variable accessors

func (ds *DialSched) addStatic(n *discover.Node) {
	if n == nil {
		return
	}
	if n.Incomplete() {
		logger.Warn("Rejecting incomplete static node", "node", n.String())
		return
	}
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.static.add(cloneNode(n))

	// Additionally cleanup the dial speed bumps.
	// This allows admin_addPeer to force an immediate reconnect,
	// even if the node was previously removed due to repeated dial failures.
	delete(ds.dialBackoff, n.ID)
	delete(ds.connFails, n.ID)
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

func (ds *DialSched) markRefreshStart() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.refreshBackoff = time.Now().Add(refreshBackoff)
}

func (ds *DialSched) shouldDial(n *discover.Node) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if ds.closed.Load() {
		return false
	}
	if n == nil {
		return false
	}
	if n.Incomplete() {
		logger.Warn("Rejecting incomplete dial candidate", "node", n.String())
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

func (ds *DialSched) isDialing() bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.dialing.len() > 0
}

// Returns the connection flag to be used for dialing this node.
func (ds *DialSched) markDialStart(n *discover.Node) connFlag {
	dialTryCounter.Inc(1)

	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.dialing.add(n)
	if ds.static.contains(n.ID) {
		return staticDialedConn
	}
	return dynDialedConn
}

func (ds *DialSched) markDialEnd(id discover.NodeID) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.dialing.remove(id)
	ds.dialBackoff[id] = time.Now().Add(dialBackoff)
}

func (ds *DialSched) markPeerConnected(id discover.NodeID, nType discover.NodeType, inbound bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// A placeholder to populate the typedNodeSet. Other fields are irrelevant in DialSched.
	n := &discover.Node{ID: id, NType: nType}

	ds.connectedAll.add(n)
	if inbound {
		ds.connectedOutbound.remove(id)
	} else {
		ds.connectedOutbound.add(n)
	}

	// Reset consecutive failure counter.
	delete(ds.connFails, id)
}

func (ds *DialSched) markDialFailure(id discover.NodeID) {
	dialFailCounter.Inc(1)

	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.connectedAll.remove(id)
	ds.connectedOutbound.remove(id)

	if n := ds.static.get(id); n != nil {
		ds.connFails[id]++
		if ds.connFails[id] > dialMaxRetries {
			logger.Warn("Removing static node after too many connection failures", "node", n.String())
			ds.static.remove(id)
		}
	}
}

func (ds *DialSched) markPeerDisconnected(id discover.NodeID) {
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
