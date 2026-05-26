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
	"sync"
	"time"
)

var (
	revalidateBackoff = 10 * time.Second // Do not revalidate (ping) a node within this duration.

	errSelfBonding   = errors.New("cannot bond with self")
	errTableNotReady = errors.New("table not ready")
)

//// Periodic revalidation

func (tab *Table2) revalidateLoop() {
	timer := time.NewTimer(revalidateInterval)
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

func (tab *Table2) doRevalidate() {
	for _, s := range tab.storages {
		tab.revalidateOnce(s)
	}
}

func (tab *Table2) revalidateOnce(s tableStorage) {
	n := s.oldest() // pick a node to revalidate.
	if n == nil {
		return
	}

	if tab.bondAge(n.ID) < revalidateBackoff {
		logger.Trace("skip revalidate", "node", n.addr())
		return
	}

	err := tab.udp.ping(n.ID, n.addr())
	if err == nil {
		s.add(n) // bump the node
		tab.db.updateBondTime(n.ID, time.Now())
		logger.Debug("Discovery revalidation passed", "node", n.addr())
	} else {
		s.delete(n)
		tab.db.deleteNode(n.ID)
		logger.Debug("Discovery revalidation failed", "node", n.addr())
	}
}

// Returns the time since the last bond with the given node.
// If the node is not bonded, return a very large duration.
func (tab *Table2) bondAge(id NodeID) time.Duration {
	return time.Since(tab.db.bondTime(id))
}

//// Bonding primitives

// bondall bonds with all given nodes concurrently and returns those
// nodes for which bonding has probably succeeded. Return the nodes
// for which bonding succeeded.
func (tab *Table2) bondall(nodes []*Node) []*Node {
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
func (tab *Table2) Bond(pinged bool, n *Node) error {
	if n.ID == tab.selfID {
		return errSelfBonding
	}
	// Do not accept inbound (i.e. remote pinged me) bonding requests before
	// populating from my seeds. It prevents remote nodes from filling up the table by spamming.
	if pinged && !tab.initialized() {
		return errTableNotReady
	}
	// Do not attempt bonding if the table is closed. Could happen after tab.Close() and before udp.close()
	// You can't write to the closed DB.
	if tab.closed.Load() {
		return errTableNotReady
	}

	// Skip if we're really sure that this node is reachable and completed the bonding process.
	if tab.canAssumeBonded(n.ID) {
		return nil
	}

	// Deduplicate concurrent Bond() calls for the same node.
	_, err, _ := tab.bondGroup.Do(n.ID.String(), func() (interface{}, error) {
		err := tab.pingpong(pinged, n)
		if err == nil {
			tab.addNode(n)
			tab.recordBonded(n)
		}
		return nil, err
	})
	return err
}

// Finish the 4-way handshake:
// 1. self --PING-> remote udp.ping():write(ping)
// 2. self <-PONG-- remote udp.ping():pending(pong)
// 3. self <-PING-- remote udp.waitping():pending(ping)  if !pinged
// 4. self --PONG-> remote ping.handle()                 if !pinged
func (tab *Table2) pingpong(pinged bool, n *Node) error {
	<-tab.bondslots                                // Acquire a bonding slot to limit network usage.
	defer func() { tab.bondslots <- struct{}{} }() // Release the bonding slot.

	// Ping the remote side and wait for a pong.
	if err := tab.udp.ping(n.ID, n.addr()); err != nil {
		return err
	}
	if !pinged {
		// Give the remote node a chance to ping us before we start
		// sending findnode requests. If they still remember us,
		// waitping will simply time out.
		tab.udp.waitping(n.ID, n.IP)
	}
	return nil
}

// A node is assumed to be bonded if all values are recently written by recordBonded().
func (tab *Table2) canAssumeBonded(id NodeID) bool {
	recentlyBonded := tab.db.hasBond(id)
	neverFailed := tab.db.findFails(id) == 0
	storedInDB := tab.db.node(id) != nil
	return recentlyBonded && neverFailed && storedInDB
}

func (tab *Table2) recordBonded(n *Node) {
	tab.db.updateBondTime(n.ID, time.Now())
	tab.db.updateFindFails(n.ID, 0)
	tab.db.updateNode(n)

	bucketEntriesGauge.Update(int64(tab.len()))
	bucketReplacementsGauge.Update(int64(tab.numReplacements()))
}
