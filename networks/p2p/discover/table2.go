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

package discover

import (
	"crypto/rand"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/crypto"
)

// Table is a self-maintained directory of nodes. Its key features are:
// - Refresh: Periodically or manually re-fills the table by looking up new nodes through FINDNODE requests.
// - Revalidate: Periodically checks the health of the nodes in the table by sending PING packets.
// - Bond: Establishes a UDP bond with a node by sending PING and PONG packets.
// - Resolve: Finds a node address by its ID by sending FINDNODE requests to other nodes.
// - RandomNodes & ClosestNodes: Serves as an API to retrieve nodes from the table.
type table struct {
	self *Node
	rand *sharedRand
	udp  transport

	// Node addresses
	nursery        []*Node
	storages       map[NodeType]tableStorage
	refreshTargets map[NodeType]int // target number of nodes to obtain for each node type. set 0 to disable discovery (i.e. only accept inbound connections).

	// UDP bonding
	bondmu    sync.Mutex
	bonding   map[NodeID]*bondtask // prevents concurrent bonding of the same node
	bondslots chan struct{}        // Semaphore to limit total number of active bonding processes

	// Coordination
	wg         sync.WaitGroup // to wait for maintenance loops to complete.
	closemu    sync.Mutex
	init       atomic.Bool
	closed     atomic.Bool
	refreshReq chan chan struct{} // to request a discretionary refresh, apart from the periodic refresh.
	closeReq   chan struct{}      // closed to stop loops.
}

type bondtask struct {
	done chan struct{}
	err  error
}

func newTable2(cfg *Config, udp transport) (*table, error) {
	var (
		self     = NewNode(cfg.Id, cfg.Addr.IP, uint16(cfg.Addr.Port), uint16(cfg.Addr.Port), nil, cfg.NodeType)
		rand     = newSharedRand()
		storages = map[NodeType]tableStorage{
			NodeTypeCN: newSimpleStorage2(rand),
			NodeTypePN: newSimpleStorage2(rand),
			NodeTypeEN: newKademliaStorage(self, rand),
			NodeTypeBN: newSimpleStorage2(rand),
		}
		refreshTargets = getRefreshTargets(cfg)
		bondslots      = make(chan struct{}, maxBondingPingPongs)

		err error
	)

	rand.Seed()
	nursery, err := getBootnodes(cfg)
	if err != nil {
		return nil, err
	}
	for i := 0; i < cap(bondslots); i++ {
		bondslots <- struct{}{}
	}

	tab := &table{
		self: self,
		rand: rand,
		udp:  udp,

		nursery:        nursery,
		storages:       storages,
		refreshTargets: refreshTargets,

		bonding:   make(map[NodeID]*bondtask),
		bondslots: bondslots,

		refreshReq: make(chan chan struct{}),
		closeReq:   make(chan struct{}),
	}

	// Insert the initial seeds into all storages.
	for _, n := range nursery {
		tab.addNode(n)
	}

	return tab, nil
}

//// Initialization

func getBootnodes(cfg *Config) ([]*Node, error) {
	// Filter bootnodes to add to nursery
	nursery := make([]*Node, 0, len(cfg.Bootnodes))
	for _, n := range cfg.Bootnodes {
		if err := n.validateComplete(); err != nil {
			return nil, fmt.Errorf("bad bootstrap node %q (%v)", n, err)
		}
		if cfg.NetRestrict != nil && !cfg.NetRestrict.Contains(n.IP) {
			logger.Warn("bootstrap node filtered by netrestrict", "node", n.String())
			continue
		}
		nursery = append(nursery, n)
	}
	return nursery, nil
}

// Returns whether the table's initial seeding procedure has completed.
func (tab *table) initialized() bool {
	return tab.init.Load()
}

func (tab *table) Start() {
	tab.wg.Add(2)
	go tab.refreshLoop()
	go tab.revalidateLoop()
}

func (tab *table) Close() {
	tab.closemu.Lock()
	defer tab.closemu.Unlock()
	if tab.closed.Load() {
		return
	}

	logger.Info("Discovery table closing")
	close(tab.closeReq) // Notify the loops to stop.
	tab.wg.Wait()       // Wait for the loops to finish.
	tab.closed.Store(true)
	logger.Info("Discovery table closed")
}

//// Periodic maintenance loops

func (tab *table) refreshLoop() {
	var (
		refreshTicker = time.NewTicker(refreshInterval)
	)
	defer tab.wg.Done()
	defer refreshTicker.Stop()

	// Initial refresh.
	tab.doRefresh()
	tab.init.Store(true)

	for {
		select {
		case <-refreshTicker.C:
			tab.rand.Seed() // reseed the randomness before consuming it during refresh.
			tab.doRefresh()
		case wait := <-tab.refreshReq:
			tab.doRefresh()
			close(wait) // Notify the Refresh() caller.
		case <-tab.closeReq:
			return
		}
	}
}

func (tab *table) revalidateLoop() {
	var (
		timer = time.NewTimer(revalidateInterval)
	)
	defer tab.wg.Done()
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			tab.doRevalidate()
			timer.Reset(revalidateInterval)
		case <-tab.closeReq:
			return
		}
	}
}

//// Refresh and lookup

func (tab *table) Refresh() {
	wait := make(chan struct{})
	select {
	case <-tab.closeReq:
		return
	case tab.refreshReq <- wait:
	}

	select {
	case <-tab.closeReq:
	case <-wait:
	}
}

func (tab *table) doRefresh() {
	logger.Info("Discovery table refreshing", "counts", tab.lenByNodeTypes())
	defer logger.Info("Discovery table refreshed", "counts", tab.lenByNodeTypes())

	tab.kademliaRefresh(NodeTypeEN)
	tab.simpleRefresh(NodeTypeCN)
	tab.simpleRefresh(NodeTypePN)
	tab.simpleRefresh(NodeTypeBN)
}

func (tab *table) kademliaRefresh(targetType NodeType) {
	if tab.refreshTargets[targetType] <= 0 {
		return
	}
	// Run self lookup to discover new neighbor nodes.
	tab.lookup(tab.self.ID, targetType, true, bucketSize)

	// The Kademlia paper specifies that the bucket refresh should
	// perform a lookup in the least recently used bucket. We cannot
	// adhere to this because the findnode target is a 512bit value
	// (not hash-sized) and it is not easily possible to generate a
	// sha3 preimage that falls into a chosen bucket.
	// We perform a few lookups with a random target instead.
	for i := 0; i < 3; i++ {
		var target NodeID
		rand.Read(target[:])
		tab.lookup(target, targetType, true, bucketSize)
	}
}

func (tab *table) simpleRefresh(targetType NodeType) {
	if tab.refreshTargets[targetType] <= 0 {
		return
	}
	tab.lookup(tab.self.ID, targetType, false, tab.refreshTargets[targetType])
}

// Query more nodes from the network by sending FINDNODE requests.
// Upon receiving a NEIGHBORS response, the nodes are first bonded, then added to the table.
func (tab *table) lookup(targetID NodeID, targetType NodeType, recurse bool, max int) []*Node {
	s := tab.storages[targetType]
	if s == nil {
		return nil
	}
	if max <= 0 {
		return nil
	}

	// Start from the closest nodes to the target Or bootnodes.
	targetHash := crypto.Keccak256Hash(targetID[:])
	seeds := s.closest(targetHash, max).entries
	if len(seeds) == 0 {
		for _, n := range tab.nursery {
			seeds = append(seeds, n)
		}
	}

	// Make sure to bond the seed nodes before sending FINDNODE requests.
	seeds = tab.bondall(seeds)

	// Find new nodes by querying to the seed nodes.
	return tab.findNodes(seeds, targetID, targetType, recurse, max)
}

// findNodes finds new nodes by querying to the seed nodes.
func (tab *table) findNodes(seeds []*Node, targetID NodeID, targetType NodeType, recurse bool, max int) []*Node {
	var (
		excludeBn = targetType != NodeTypeBN && targetType != NodeTypeUnknown // Result can include bootnodes.

		pool    = make(map[NodeID]*Node)    // All nodes being processed in this function. Append only.
		isBn    = make(map[NodeID]bool)     // Bootnodes in the pool.
		unasked = make(map[NodeID]bool)     // Nodes in the pool that we did not ask yet.
		resCh   = make(chan []*Node, alpha) // Replies from findnode requests.
		running = 0
	)
	for _, seed := range seeds {
		pool[seed.ID] = seed
		isBn[seed.ID] = seed.NType == NodeTypeBN
		unasked[seed.ID] = true
	}

	// Launch the first round of up to alpha findnode requests.
	for id := range unasked { // pick random nodes to ask.
		if running >= alpha {
			break
		}
		delete(unasked, id)
		running++
		target := pool[id]
		go func() { resCh <- tab.findNodesOnce(target, targetID, targetType, max) }()
	}

	for running > 0 { // Wait for all running requests to finish.
		res := <-resCh // One request finished.
		running--

		// Add the nodes to the pool.
		for _, n := range res {
			if _, ok := pool[n.ID]; !ok {
				pool[n.ID] = n
				isBn[n.ID] = n.NType == NodeTypeBN
				unasked[n.ID] = true
			}
		}

		// If we need more and willing to recursively ask, launch a new findnode request.
		var enough bool
		if excludeBn {
			enough = len(pool)-len(isBn) >= max
		} else {
			enough = len(pool) >= max
		}
		if recurse && !enough {
			for id := range unasked { // pick a random node to ask.
				delete(unasked, id)
				running++
				target := pool[id]
				go func() { resCh <- tab.findNodesOnce(target, targetID, targetType, max) }()
				break // Launch only one new request.
			}
		}
	}

	// Return the non-bootnodes. The result may include the initial seeds.
	targetHash := crypto.Keccak256Hash(targetID[:])
	ordered := &nodesByDistance{target: targetHash}
	for id, n := range pool {
		if excludeBn && isBn[id] {
			continue
		}
		ordered.push(n, max)
	}
	return ordered.entries
}

func (tab *table) findNodesOnce(seed *Node, targetID NodeID, targetType NodeType, max int) []*Node {
	// TODO: update database counters
	r, err := tab.udp.findnode(seed.ID, seed.addr(), targetID, targetType, max)
	if err != nil {
		return nil
	} else {
		return tab.bondall(r)
	}
}

//// Bond methods

// bondall bonds with all given nodes concurrently and returns those
// nodes for which bonding has probably succeeded. Return the nodes
// for which bonding succeeded.
func (tab *table) bondall(nodes []*Node) []*Node {
	rc := make(chan *Node, len(nodes))
	wg := sync.WaitGroup{}
	for i := range nodes {
		wg.Add(1)
		n := nodes[i]
		go func(n *Node) {
			defer wg.Done()
			if err := tab.Bond(false, n); err == nil {
				rc <- n
			}
		}(n)
	}
	wg.Wait()
	close(rc)

	result := make([]*Node, 0, len(nodes))
	for n := range rc {
		result = append(result, n)
	}
	return result
}

// Bond ensures the local node has a bond with the given remote node.
// It also attempts to insert the node into the table if bonding succeeds.
//
// A bond is must be established before sending findnode requests.
// Both sides must have completed a ping/pong exchange for a bond to
// exist. The total number of active bonding processes is limited in
// order to restrain network use.
//
// bond is meant to operate idempotently in that bonding with a remote
// node which still remembers a previously established bond will work.
// The remote node will simply not send a ping back, causing waitping
// to time out.
//
// If pinged is true, the remote node has just pinged us and one half
// of the process can be skipped.
func (tab *table) Bond(pinged bool, n *Node) error {
	if n.ID == tab.self.ID {
		return errSelfBonding
	}
	// Do not accept inbound (i.e. remote pinged me) bonding requests before
	// populating from my seeds. It prevents remote nodes from filling up the table.
	if pinged && !tab.initialized() {
		return errTableNotInitialized
	}

	// Check for already ongoing bonding task.
	tab.bondmu.Lock()
	task := tab.bonding[n.ID]
	if task != nil {
		// Wait for an existing bonding task to complete.
		tab.bondmu.Unlock()
		<-task.done
		return task.err
	}

	// Register a new bonding task.
	task = &bondtask{done: make(chan struct{})}
	tab.bonding[n.ID] = task
	tab.bondmu.Unlock()

	// Execute the bonding task.
	tab.pingpong(task, pinged, n)

	// Unregister the task after it's done.
	tab.bondmu.Lock()
	delete(tab.bonding, n.ID)
	tab.bondmu.Unlock()

	if task.err != nil {
		return task.err
	}

	// If bond was successful, add the node to the table.
	tab.addNode(n)
	// TODO: update database
	// TODO: update metrics
	return nil
}

// Finish the 4-way handshake:
// 1. self --PING-> remote udp.ping():write(ping)
// 2. self <-PONG-- remote udp.ping():pending(pong)
// 3. self <-PING-- remote udp.waitping():pending(ping)  if !pinged
// 4. self --PONG-> remote ping.handle()                 if !pinged
func (tab *table) pingpong(task *bondtask, pinged bool, n *Node) {
	// Request a bonding slot to limit network usage.
	<-tab.bondslots
	defer func() { tab.bondslots <- struct{}{} }()
	defer close(task.done)

	// Ping the remote side and wait for a pong.
	if task.err = tab.udp.ping(n.ID, n.addr()); task.err != nil {
		return
	}
	if !pinged {
		// Give the remote node a chance to ping us before we start
		// sending findnode requests. If they still remember us,
		// waitping will simply time out.
		tab.udp.waitping(n.ID, n.IP)
	}
}

//// Revalidate

func (tab *table) doRevalidate() {
	for _, s := range tab.storages {
		tab.revalidateOnce(s)
	}
}

func (tab *table) revalidateOnce(s tableStorage) {
	n := s.oldest()
	if n == nil {
		return
	}

	// TODO: read last bond time

	err := tab.udp.ping(n.ID, n.addr())
	if err == nil {
		logger.Debug("Discovery revalidation passed", "node", n.addr())
		s.add(n) // bump the node
		// TODO: update last bond time
	} else {
		logger.Debug("Discovery revalidation failed", "node", n.addr())
		s.delete(n)
	}
}

//// Storage and database access wrappers

// Adds a node to one of the storages of the node's type.
// If the node is a bootstrap node, it is added to all storages.
func (tab *table) addNode(n *Node) {
	if n.NType == NodeTypeBN {
		for _, s := range tab.storages {
			s.add(n)
		}
		return
	} else {
		if s := tab.storages[n.NType]; s != nil {
			s.add(n)
		}
	}
}

//// Getters

func (tab *table) RandomNodes(buf []*Node, targetType NodeType) int {
	if s := tab.storages[targetType]; s == nil {
		return 0
	} else {
		return s.random(buf)
	}
}

func (tab *table) ClosestNodes(targetID NodeID, targetType NodeType, max int) []*Node {
	if s := tab.storages[targetType]; s == nil {
		return nil
	} else {
		return s.closest(crypto.Keccak256Hash(targetID[:]), max).entries
	}
}

// Count of all nodes in the table.
// Bootnodes are counted multiple times because they are added to all storages.
func (tab *table) len() int {
	n := 0
	for _, s := range tab.storages {
		n += s.len()
	}
	return n
}

// For logging
func (tab *table) lenByNodeTypes() map[string]int {
	r := make(map[string]int)
	for ty, s := range tab.storages {
		r[nodeTypeName(ty)] = s.len()
	}
	return r
}
