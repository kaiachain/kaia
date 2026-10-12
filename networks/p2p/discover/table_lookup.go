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
	"time"

	"github.com/kaiachain/kaia/crypto"
)

func (tab *Table2) refreshLoop() {
	refreshTicker := time.NewTicker(refreshInterval)
	defer tab.wg.Done()
	defer refreshTicker.Stop()

	// Initial refresh.
	tab.Refresh()
	tab.init.Store(true)

	for {
		select {
		case <-refreshTicker.C:
			tab.rand.Seed() // reseed the randomness before consuming it during refresh.
			tab.doRefresh()
		case <-tab.closeReq:
			return
		}
	}
}

func (tab *Table2) Refresh() {
	if tab.closed.Load() {
		return
	}
	tab.doRefresh()
}

// doRefresh performs a full table refresh. Concurrent calls are deduplicated and wait for the single in-flight refresh.
func (tab *Table2) doRefresh() {
	tab.refreshGroup.Do("refresh", func() (interface{}, error) {
		logger.Debug("Discovery table refreshing", "counts", tab.lenByNodeTypes())
		tab.kademliaRefresh(NodeTypeEN)
		tab.simpleRefresh(NodeTypeCN)
		tab.simpleRefresh(NodeTypePN)
		tab.simpleRefresh(NodeTypeBN)
		logger.Debug("Discovery table refreshed", "counts", tab.lenByNodeTypes())
		return nil, nil
	})
}

func (tab *Table2) kademliaRefresh(targetType NodeType) {
	if tab.discoverTargets[targetType] <= 0 {
		return
	}

	// Run self lookup to discover new neighbor nodes.
	tab.lookup(tab.selfID, targetType, true, bucketSize)

	// The Kademlia paper specifies that the bucket refresh should
	// perform a lookup in the least recently used bucket. We cannot
	// adhere to this because the findnode target is a 512bit value
	// (not hash-sized) and it is not easily possible to generate a
	// sha3 preimage that falls into a chosen bucket.
	// We perform a few lookups with a random target instead.
	for i := 0; i < 3; i++ {
		var target NodeID
		tab.rand.Read(target[:])
		tab.lookup(target, targetType, true, bucketSize)
	}
}

func (tab *Table2) simpleRefresh(targetType NodeType) {
	if tab.discoverTargets[targetType] <= 0 {
		return
	}
	tab.lookup(tab.selfID, targetType, false, tab.discoverTargets[targetType])
}

// Query more nodes from the network by sending FINDNODE requests.
// Upon receiving a NEIGHBORS response, the nodes are first bonded, then added to the table.
func (tab *Table2) lookup(targetID NodeID, targetType NodeType, recurse bool, max int) []*Node {
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
func (tab *Table2) findNodes(seeds []*Node, targetID NodeID, targetType NodeType, recurse bool, max int) []*Node {
	var (
		// Bootnodes could be in the seeds but they are not usually the result we want.
		excludeBn = targetType != NodeTypeBN
		pool      = newFindNodesPool(seeds)
		resCh     = make(chan []*Node, alpha)
		running   = 0
	)

	// Launch the first round of findnode requests.
	for running < alpha {
		seed := pool.popUnasked()
		if seed == nil {
			break
		}
		running++
		go func() { resCh <- tab.findNodesOnce(seed, targetID, targetType, max) }()
		logger.Trace("Discovery request sent", "targetType", targetType, "have", pool.count(excludeBn), "max", max)
	}

	// Wait for all running requests to finish.
	for running > 0 {
		// Wait for one request to finish.
		res := <-resCh
		running--
		fresh := pool.push(res)
		logger.Trace("Discovery response received", "targetType", targetType, "response", len(res), "fresh", fresh,
			"pool", len(pool.all), "bootnodes", pool.numBootnodes, "unasked", len(pool.unasked))

		// If we don't have enough nodes and willing to recursively ask, launch a new findnode request.
		if recurse && pool.count(excludeBn) < max {
			if seed := pool.popUnasked(); seed != nil {
				running++
				go func() { resCh <- tab.findNodesOnce(seed, targetID, targetType, max) }()
				logger.Trace("Discovery recursive request sent", "targetType", targetType, "have", pool.count(excludeBn), "max", max)
			}
		}
	}

	result := pool.harvest(targetID, excludeBn, max)
	logger.Debug("Discovery result", "self", targetID == tab.selfID, "targetType", targetType, "recurse", recurse, "count", len(result))
	return result
}

func (tab *Table2) findNodesOnce(seed *Node, targetID NodeID, targetType NodeType, max int) []*Node {
	r, err := tab.udp.findnode(seed.ID, seed.addr(), targetID, targetType, max)
	// A timeout with no NEIGHBORS at all means the seed is unresponsive
	if errors.Is(err, errTimeout) && len(r) == 0 {
		tab.recordFindFailure(seed)
		return nil
	}
	// Drop bootnodes from the response unless the caller is explicitly looking for them
	if targetType != NodeTypeBN {
		r = removeBn(r)
	}
	// Out of the find results, return only reachable (bonded) nodes.
	return tab.bondall(r)
}

// Records a findnode failure.
func (tab *Table2) recordFindFailure(seed *Node) {
	fails := tab.db.findFails(seed.ID, seed.IP) + 1
	tab.db.updateFindFails(seed.ID, seed.IP, fails)
	logger.Trace("Findnode failure", "id", seed.ID, "failcount", fails)

	if fails >= maxFindnodeFailures {
		tab.deleteNode(seed)
		tab.db.deleteNode(seed.ID)
		logger.Trace("Too many findnode failures, dropping", "id", seed.ID, "failcount", fails)
	}
}

// findNodesPool remembers the seed nodes and fresh nodes being processed in findNodes().
type findNodesPool struct {
	all          map[NodeID]*Node    // Cumulative list of nodes being processed in findNodes(). Append only and never removed.
	unasked      map[NodeID]struct{} // Nodes in the pool that we did not ask findnode yet.
	numBootnodes int                 // Number of bootnodes in the pool.
}

func newFindNodesPool(seeds []*Node) *findNodesPool {
	p := &findNodesPool{
		all:     make(map[NodeID]*Node),
		unasked: make(map[NodeID]struct{}),
	}
	p.push(seeds)
	return p
}

func (p *findNodesPool) count(excludeBn bool) int {
	if excludeBn {
		return len(p.all) - p.numBootnodes
	} else {
		return len(p.all)
	}
}

// Pushes new nodes to the pool. Returns the number of new nodes added.
func (p *findNodesPool) push(nodes []*Node) int {
	added := 0
	for _, n := range nodes {
		if _, ok := p.all[n.ID]; ok {
			continue
		}
		p.all[n.ID] = n
		p.unasked[n.ID] = struct{}{}
		if n.NType == NodeTypeBN {
			p.numBootnodes++
		}
		added++
	}
	return added
}

// Pops a random unasked node.
func (p *findNodesPool) popUnasked() *Node {
	for id := range p.unasked {
		delete(p.unasked, id)
		return p.all[id]
	}
	return nil
}

// Returns up to max nodes ordered by distance to the target.
func (p *findNodesPool) harvest(targetID NodeID, excludeBn bool, max int) []*Node {
	ordered := &nodesByDistance{target: crypto.Keccak256Hash(targetID[:])}
	for _, n := range p.all {
		if excludeBn && n.NType == NodeTypeBN {
			continue
		}
		ordered.push(n, max)
	}
	return ordered.entries
}
