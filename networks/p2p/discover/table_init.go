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
	"sync"
	"sync/atomic"

	"golang.org/x/sync/singleflight"
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
	bondGroup singleflight.Group // deduplicates concurrent Bond() calls for the same node.
	bondslots chan struct{}      // semaphore to limit total number of active bonding processes

	// UDP refresh
	refreshGroup singleflight.Group // deduplicates concurrent Refresh() calls.
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
	nursery, err := getBootnodes(cfg)
	if err != nil {
		return nil, err
	}

	db, err := newNodeDB(cfg.NodeDBPath, Version, cfg.Id)
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
		selfID:          cfg.Id,
		nursery:         nursery,
		storages:        storages,
		discoverTargets: discoverTargets,
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
	tab.db.close()
	logger.Info("Discovery table closed")
}

// Returns whether the table's initial seeding procedure has completed.
func (tab *Table2) initialized() bool {
	return tab.init.Load()
}

// getDiscoverTargets returns this node's discovery goals by remote node type.
// These targets are not connection counts; p2p.Server and p2p.DialSched regulate peers.
func getDiscoverTargets(cfg *Config) map[NodeType]int {
	targets := discoverTargets[EffectiveNodeType(cfg.NodeType)]
	if len(targets) == 0 {
		logger.Error("Unsupported node type", "NodeType", cfg.NodeType)
	}
	return cloneNodeTypeTargets(targets)
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
