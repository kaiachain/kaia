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
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_findNodesPool_push_dedup(t *testing.T) {
	n := newTestNode(t, NodeTypeEN)
	p := newFindNodesPool([]*Node{n})
	added := p.push([]*Node{n}) // same node again
	assert.Equal(t, 0, added, "duplicate node should not be added")
	assert.Len(t, p.all, 1)
}

func Test_findNodesPool_count(t *testing.T) {
	en := newTestNode(t, NodeTypeEN)
	bn := newTestBootnode(t, net.IP{127, 0, 0, 1})
	p := newFindNodesPool([]*Node{en, bn})
	assert.Equal(t, 2, p.count(false), "total count including bootnodes")
	assert.Equal(t, 1, p.count(true), "count excluding bootnodes")
}

func Test_findNodesPool_popUnasked(t *testing.T) {
	n := newTestNode(t, NodeTypeEN)
	p := newFindNodesPool([]*Node{n})
	got := p.popUnasked()
	assert.Equal(t, n.ID, got.ID)
	assert.Nil(t, p.popUnasked(), "pool is exhausted after one pop")
}

func Test_findNodesPool_harvest(t *testing.T) {
	en1 := newTestNode(t, NodeTypeEN)
	en2 := newTestNode(t, NodeTypeEN)
	bn := newTestBootnode(t, net.IP{127, 0, 0, 1})
	p := newFindNodesPool([]*Node{en1, en2, bn})

	// excludeBn=true: the bootnode should be absent from results.
	result := p.harvest(NodeID{}, true, 10)
	assert.Len(t, result, 2)
	for _, n := range result {
		assert.NotEqual(t, NodeTypeBN, n.NType)
	}

	// max=1 caps the output regardless of pool size.
	result = p.harvest(NodeID{}, true, 1)
	assert.Len(t, result, 1)
}

func Test_Table_recordFindFailure_increments(t *testing.T) {
	tab := newTestTable2(t, nil)
	n := newTestNode(t, NodeTypeEN)
	tab.addNode(n)

	require.Equal(t, 0, tab.db.findFails(n.ID))
	tab.recordFindFailure(n)
	assert.Equal(t, 1, tab.db.findFails(n.ID))
	// Node still in the table (1 < maxFindnodeFailures).
	assert.Equal(t, 1, tab.storages[NodeTypeEN].len())
}

func Test_Table_recordFindFailure_drops(t *testing.T) {
	tab := newTestTable2(t, nil)
	n := newTestNode(t, NodeTypeEN)
	tab.addNode(n)
	tab.db.updateNode(n)

	for i := 0; i < maxFindnodeFailures; i++ {
		tab.recordFindFailure(n)
	}
	// Node must be evicted from storage and database.
	assert.Equal(t, 0, tab.storages[NodeTypeEN].len())
	assert.Nil(t, tab.db.node(n.ID))
}

func Test_Table_findNodesOnce_error(t *testing.T) {
	udp := newMockTransport()
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		return nil, errTimeout
	}
	tab := newTestTable2(t, udp)
	seed := newTestNode(t, NodeTypeEN)
	tab.addNode(seed)
	tab.db.updateNode(seed)

	result := tab.findNodesOnce(seed, NodeID{}, NodeTypeEN, 10)
	assert.Nil(t, result)
	assert.Equal(t, 1, tab.db.findFails(seed.ID), "failure should be recorded in DB")
}

func Test_Table_findNodesOnce_success(t *testing.T) {
	peer := newTestNode(t, NodeTypeEN)
	udp := newMockTransport()
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		return []*Node{peer}, nil
	}
	tab := newTestTable2(t, udp)
	seed := newTestNode(t, NodeTypeEN)

	result := tab.findNodesOnce(seed, NodeID{}, NodeTypeEN, 10)
	// Peer is bonded (ping returns nil by default) so it appears in results.
	require.Len(t, result, 1)
	assert.Equal(t, peer.ID, result[0].ID)
}

func Test_Table_lookup_emptyTable(t *testing.T) {
	var findCalled int32
	udp := newMockTransport()
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		atomic.AddInt32(&findCalled, 1)
		return nil, nil
	}
	tab := newTestTable2(t, udp)

	// Neither seeds nor nursery → lookup terminates immediately.
	result := tab.lookup(NodeID{}, NodeTypeEN, false, 10)
	assert.Nil(t, result)
	assert.Equal(t, int32(0), atomic.LoadInt32(&findCalled), "no findnode calls expected on an empty table")
}

func Test_Table_lookup_withNursery(t *testing.T) {
	var findCalled int32
	udp := newMockTransport()
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		atomic.AddInt32(&findCalled, 1)
		return nil, nil
	}
	tab := newTestTable2(t, udp)
	tab.nursery = []*Node{newTestBootnode(t, net.IP{127, 0, 0, 1})}

	tab.lookup(NodeID{}, NodeTypeEN, false, 10)
	assert.Greater(t, atomic.LoadInt32(&findCalled), int32(0), "nursery BN should trigger at least one findnode call")
}

func Test_Table_lookup_returnsSeedResults(t *testing.T) {
	discovered := newTestNode(t, NodeTypeEN)
	udp := newMockTransport()
	// First findnode call returns a new node; further calls return nothing.
	var once sync.Once
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		var res []*Node
		once.Do(func() { res = []*Node{discovered} })
		return res, nil
	}
	tab := newTestTable2(t, udp)
	tab.nursery = []*Node{newTestBootnode(t, net.IP{127, 0, 0, 1})}

	result := tab.lookup(NodeID{}, NodeTypeEN, false, 10)
	ids := make(map[NodeID]bool)
	for _, n := range result {
		ids[n.ID] = true
	}
	assert.True(t, ids[discovered.ID], "newly discovered node should appear in lookup results")
}

func Test_Table_lookup_singleRound(t *testing.T) {
	tab := newTestTable2(t, lookupTestnet)

	// Seed with dists[256][0].  Its UDP port (256) tells the adapter which
	// distance bucket to serve when Table2 issues the first findnode request.
	seed := NewNode(lookupTestnet.dists[256][0], net.ParseIP("127.0.0.1"), 256, 0, nil, NodeTypeEN)
	tab.storages[NodeTypeEN].add(seed)

	results := tab.lookup(lookupTestnet.target, NodeTypeEN, false, bucketSize)

	t.Logf("results (%d):", len(results))
	for _, n := range results {
		t.Logf("  ld=%d  %x", logdist(lookupTestnet.targetSha, n.sha), n.sha[:8])
	}
	assert.Equal(t, bucketSize, len(results), "should return exactly bucketSize results")
	assert.False(t, hasDuplicates(results), "result set contains duplicate entries")
	assert.True(t, sortedByDistanceTo(lookupTestnet.targetSha, results), "result set not sorted by distance to target")
}

func Test_Table_lookup_recursive(t *testing.T) {
	tab := newTestTable2(t, lookupTestnet)

	seed := NewNode(lookupTestnet.dists[256][0], net.ParseIP("127.0.0.1"), 256, 0, nil, NodeTypeEN)
	tab.storages[NodeTypeEN].add(seed)

	// lookupTestnet.findnode would return at most bucketSize results.
	// lookup()-findNodes() would call 3 findOnce() calls, yielding only at most 3*bucketSize results.
	// The recursive lookup would call findOnce() more times, fulfilling the max results.
	results := tab.lookup(lookupTestnet.target, NodeTypeEN, true, bucketSize*4)

	t.Logf("results (%d):", len(results))
	for _, n := range results {
		t.Logf("  ld=%d  %x", logdist(lookupTestnet.targetSha, n.sha), n.sha[:8])
	}
	assert.Equal(t, bucketSize*4, len(results), "should return exactly 2×bucketSize results")
	assert.False(t, hasDuplicates(results), "result set contains duplicate entries")
	assert.True(t, sortedByDistanceTo(lookupTestnet.targetSha, results), "result set not sorted by distance to target")
}

// Test_Table_doRefreshOnce_dedup verifies that a second concurrent call to
// doRefreshOnce() waits for the in-flight refresh to finish rather than
// starting its own.
func Test_Table_doRefreshOnce_dedup(t *testing.T) {
	var once sync.Once
	firstFindnode := make(chan struct{})
	blockFindnode := make(chan struct{})

	udp := newMockTransport()
	udp.findnodeFn = func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error) {
		// Block only the very first call; all subsequent ones pass through.
		once.Do(func() {
			close(firstFindnode)
			<-blockFindnode
		})
		return nil, nil
	}
	tab := newTestTable2(t, udp)
	tab.nursery = []*Node{newTestBootnode(t, net.IP{127, 0, 0, 1})}

	// g1 starts a refresh and gets blocked inside the first findnode call.
	var g1 sync.WaitGroup
	g1.Add(1)
	go func() {
		defer g1.Done()
		tab.doRefreshOnce()
	}()
	<-firstFindnode // g1 is now blocking; refreshDone is registered

	// The in-progress refresh channel must exist.
	tab.refreshmu.Lock()
	assert.NotNil(t, tab.refreshDone, "refresh should be in progress")
	tab.refreshmu.Unlock()

	// g2 should join g1's refresh and not start a new one.
	g2done := make(chan struct{})
	go func() {
		tab.doRefreshOnce()
		close(g2done)
	}()

	// Give g2 time to call doRefreshOnce and block on refreshDone.
	time.Sleep(20 * time.Millisecond)
	select {
	case <-g2done:
		t.Error("g2 completed before g1 was released — dedup did not work")
	default:
		// expected: g2 is waiting on the refresh channel
	}

	// Release g1; both goroutines should now complete.
	close(blockFindnode)
	g1.Wait()
	select {
	case <-g2done:
		// expected
	case <-time.After(time.Second):
		t.Error("g2 did not complete within 1 s after g1 finished")
	}

	// refreshDone must be cleared after completion.
	tab.refreshmu.Lock()
	assert.Nil(t, tab.refreshDone)
	tab.refreshmu.Unlock()
}
