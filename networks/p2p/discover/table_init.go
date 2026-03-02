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

package discover

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

var errBadBootnode = errors.New("bad bootstrap node")

// Table is a self-maintained directory of nodes. Its key features are:
// - Bond: Establishes a UDP bond with a node by sending PING and PONG packets.
// - Refresh: Periodically or manually re-fills the table by looking up new nodes through FINDNODE requests.
// - Revalidate: Periodically checks the health of the nodes in the table by sending PING packets.
// - RandomNodes & ClosestNodes: Serves as an API to retrieve nodes from the table.
type Table2 struct {
	// Relying components
	rand *sharedRand
	db   *nodeDB
	udp  transport

	// Lifecycle
	wg       sync.WaitGroup
	init     atomic.Bool
	closed   atomic.Bool
	closeReq chan struct{}

	// Node addresses
	selfID   NodeID
	nursery  []*Node
	storages map[NodeType]tableStorage

	// target number of nodes to obtain for each node type.
	// set discoverTargets[T] = 0 to disable discovery of node type T (i.e. only accept inbound connections).
	// see also: p2p.DialSched.connTargets.
	discoverTargets map[NodeType]int

	// UDP bonding
	bondmu    sync.Mutex
	bonding   map[NodeID]*bondtask // prevents concurrent bonding of the same node
	bondslots chan struct{}        // semaphore to limit total number of active bonding processes

	// UDP refresh
	refreshmu   sync.Mutex
	refreshDone chan struct{} // prevents concurrent refresh. refreshDone is nil when no refresh is running.
}

func newTable2(cfg *Config, udp transport) (*Table2, error) {
	var (
		self     = NewNode(cfg.Id, cfg.Addr.IP, uint16(cfg.Addr.Port), uint16(cfg.Addr.Port), nil, cfg.NodeType)
		rand     = newSharedRand()
		storages = map[NodeType]tableStorage{
			NodeTypeCN: newFlatStorage(rand),
			NodeTypePN: newFlatStorage(rand),
			NodeTypeEN: newKademliaStorage(self, rand),
			NodeTypeBN: newFlatStorage(rand),
		}
		discoverTargets = getDiscoverTargets(cfg)
		bondslots       = make(chan struct{}, maxBondingPingPongs)

		err error
	)

	rand.Seed()
	db, err := newNodeDB(cfg.NodeDBPath, Version, cfg.Id)
	if err != nil {
		return nil, err
	}

	nursery, err := getBootnodes(cfg)
	if err != nil {
		return nil, err
	}
	for i := 0; i < cap(bondslots); i++ {
		bondslots <- struct{}{}
	}

	tab := &Table2{
		rand:            rand,
		db:              db,
		udp:             udp,
		closeReq:        make(chan struct{}),
		nursery:         nursery,
		storages:        storages,
		discoverTargets: discoverTargets,
		bonding:         make(map[NodeID]*bondtask),
		bondslots:       bondslots,
	}

	// Insert the initial seeds into storages.
	for _, n := range tab.nursery {
		tab.addNode(n)
	}
	for _, n := range removeBn(tab.db.querySeeds(seedCount, seedMaxAge)) {
		tab.addNode(n)
	}

	return tab, nil
}

func (tab *Table2) Start() {
	tab.wg.Add(2)
	go tab.refreshLoop()
	go tab.revalidateLoop()
	tab.db.ensureExpirer()
}

func (tab *Table2) Close() {
	if tab.closed.Swap(true) {
		return
	}

	close(tab.closeReq) // Ask the loops to terminate.
	tab.wg.Wait()       // Wait for the loops to terminate.
	logger.Info("Discovery table closed")
}

// Returns whether the table's initial seeding procedure has completed.
func (tab *Table2) initialized() bool {
	return tab.init.Load()
}

// Returns the number of nodes to actively discover depending of self node type.
// Note that this number is NOT the number of connections to maintain.
// The connection counts (or peer counts) are regulated by the p2p.Server and p2p.DialSched.
func getDiscoverTargets(cfg *Config) map[NodeType]int {
	switch cfg.NodeType {
	case NodeTypeCN:
		// CN discovers each other in a mesh structure. Assuming up to 100 validators.
		return map[NodeType]int{
			NodeTypeCN: 100,
			NodeTypeBN: 3,
		}
	case NodeTypePN:
		// PN discovers at least one neighboring PN. However, in production, PN connections are usually specified via static nodes
		// to make up a certain network structure, in which case this number is irrelevant.
		return map[NodeType]int{
			NodeTypePN: 1,
			NodeTypeBN: 3,
		}
	case NodeTypeEN:
		// EN discovers to up to 2 PNs for high availability.
		// EN try to learn as many ENs as possible via Kademlia algorithm.
		return map[NodeType]int{
			NodeTypeEN: math.MaxInt32,
			NodeTypePN: 2,
			NodeTypeBN: 3,
		}
	case NodeTypeBN:
		// BN has to learn every CN and PN.
		return map[NodeType]int{
			NodeTypeCN: 100,
			NodeTypePN: 100,
			NodeTypeBN: 3,
		}
	default:
		logger.Error("Unsupported node type", "NodeType", cfg.NodeType)
		return map[NodeType]int{}
	}
}

// Filter bootnodes to be used as initial seeds.
func getBootnodes(cfg *Config) ([]*Node, error) {
	nursery := make([]*Node, 0, len(cfg.Bootnodes))
	for _, n := range cfg.Bootnodes {
		if err := n.validateComplete(); err != nil {
			return nil, fmt.Errorf("%w: %v", errBadBootnode, err)
		}
		if cfg.NetRestrict != nil && !cfg.NetRestrict.Contains(n.IP) {
			logger.Warn("bootstrap node filtered by netrestrict", "node", n.String())
			continue
		}
		nursery = append(nursery, n)
	}
	return nursery, nil
}
