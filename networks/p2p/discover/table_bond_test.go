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
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Table_revalidateOnce_success(t *testing.T) {
	udp := newMockTransport()
	tab := newTestTable2(t, udp)

	n := newTestNode(t, NodeTypeEN)
	s := newFlatStorage(newSharedRand())
	s.add(n)
	require.Equal(t, 1, s.len())

	tab.revalidateOnce(s)

	// Node is still in storage (bumped to front on success).
	assert.Equal(t, 1, s.len())
	// Bond time is updated in the database.
	assert.True(t, tab.db.hasBond(n.ID), "bondTime should be recent after successful revalidation")
	assert.Equal(t, int32(1), udp.pingCnt.Load())
}

func Test_Table_revalidateOnce_failure(t *testing.T) {
	udp := newMockTransport()
	tab := newTestTable2(t, udp)

	n := newTestNode(t, NodeTypeEN)
	s := newFlatStorage(newSharedRand())
	s.add(n)
	require.Equal(t, 1, s.len())

	// Pre-store the node in the DB so we can verify it gets deleted.
	tab.db.updateNode(n)
	require.NotNil(t, tab.db.node(n.ID))

	udp.setPingErr(n.ID, errTimeout)
	tab.revalidateOnce(s)

	// Node is removed from storage on failure.
	assert.Equal(t, 0, s.len())
	// Node is deleted from the database.
	assert.Nil(t, tab.db.node(n.ID), "node should be deleted from DB after failed revalidation")
	assert.Equal(t, int32(1), udp.pingCnt.Load())
}

// Test_Table_bondall_nonBlocking verifies that bondall runs bonds concurrently:
// even if node A's ping is blocked, node B's ping can complete simultaneously.
// If bondall were serialising, B could never complete while A is stuck.
func Test_Table_bondall_nonBlocking(t *testing.T) {
	nodeA := newTestNode(t, NodeTypeEN)
	nodeB := newTestNode(t, NodeTypeEN)

	bDone := make(chan struct{})       // closed when B's ping finishes
	aCanProceed := make(chan struct{}) // closed to unblock A's ping

	udp := newMockTransport()
	udp.pingFn = func(id NodeID, addr *net.UDPAddr) error {
		if id == nodeA.ID {
			<-aCanProceed // blocked until B has already finished
			return errTimeout
		}
		// nodeB: succeed and signal completion
		close(bDone)
		return nil
	}
	tab := newTestTable2(t, udp)

	resultCh := make(chan []*Node, 1)
	go func() { resultCh <- tab.bondall([]*Node{nodeA, nodeB}) }()

	// B's ping must complete concurrently while A is still blocked.
	<-bDone
	close(aCanProceed) // release A

	result := <-resultCh
	// Only B bonded successfully.
	require.Len(t, result, 1)
	assert.Equal(t, nodeB.ID, result[0].ID)
}

func Test_Table_bondall_filtering(t *testing.T) {
	n1 := newTestNode(t, NodeTypeEN)
	n2 := newTestNode(t, NodeTypeEN)
	n3 := newTestNode(t, NodeTypeEN)

	udp := newMockTransport()
	udp.setPingErr(n2.ID, errTimeout) // n2 fails
	tab := newTestTable2(t, udp)

	result := tab.bondall([]*Node{n1, n2, n3})

	// Only n1 and n3 should be returned.
	assert.Len(t, result, 2)
	ids := make(map[NodeID]bool)
	for _, n := range result {
		ids[n.ID] = true
	}
	assert.True(t, ids[n1.ID], "n1 (success) should be in result")
	assert.False(t, ids[n2.ID], "n2 (failure) should not be in result")
	assert.True(t, ids[n3.ID], "n3 (success) should be in result")
}

func Test_Table_Bond_selfBonding(t *testing.T) {
	tab := newTestTable2(t, nil)
	// tab.selfID is NodeID{} (zero); bond to self must be rejected.
	self := NewNode(NodeID{}, net.IP{127, 0, 0, 1}, 30303, 30303, nil, NodeTypeEN)
	err := tab.Bond(false, self)
	assert.ErrorIs(t, err, errSelfBonding)
}

func Test_Table_Bond_notInitialized(t *testing.T) {
	tab := newTestTable2(t, nil)
	// init is false by default; pinged=true triggers the check.
	n := newTestNode(t, NodeTypeEN)
	err := tab.Bond(true, n)
	assert.ErrorIs(t, err, errTableNotInitialized)
}

func Test_Table_Bond_happy(t *testing.T) {
	udp := newMockTransport()
	tab := newTestTable2(t, udp)

	n := newTestNode(t, NodeTypeEN)
	err := tab.Bond(false, n)

	require.NoError(t, err)
	// Node should be added to the EN storage.
	assert.Equal(t, 1, tab.storages[NodeTypeEN].len(), "node should appear in EN storage")
	// DB should reflect successful bonding.
	assert.True(t, tab.db.hasBond(n.ID), "DB bond time should be set")
	assert.NotNil(t, tab.db.node(n.ID), "DB should store the node")
	assert.Equal(t, 0, tab.db.findFails(n.ID), "DB find failures should be 0")
	assert.Equal(t, int32(1), udp.pingCnt.Load())
}

func Test_Table_Bond_duplicate(t *testing.T) {
	n := newTestNode(t, NodeTypeEN)
	pingEntered := make(chan struct{}, 2)
	pingBlocker := make(chan struct{})
	udp := newMockTransport()
	udp.pingFn = func(id NodeID, addr *net.UDPAddr) error {
		pingEntered <- struct{}{}
		<-pingBlocker
		return nil
	}
	tab := newTestTable2(t, udp)

	var err1, err2 error
	var wg sync.WaitGroup
	secondStarted := make(chan struct{})

	// The first Bond call should register a task and execute ping.
	wg.Add(1)
	go func() {
		defer wg.Done()
		err1 = tab.Bond(false, n)
	}()
	<-pingEntered // Wait until the first Bond has registered its task and is executing ping.

	// The second Bond call should find the in-flight task and wait on task.done.
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(secondStarted)
		err2 = tab.Bond(false, n)
	}()
	<-secondStarted
	for i := 0; i < 100; i++ {
		// Give the second goroutine a scheduling chance to reach Bond().
		runtime.Gosched()
	}

	// Verify the bonding map still holds the task.
	// Thread 1 is stuck at pingpong() waiting for transport.ping(), Thread 2 is stuck at Bond() waiting for task.done.
	tab.bondmu.Lock()
	assert.NotNil(t, tab.bonding[n.ID], "bonding task should be registered")
	tab.bondmu.Unlock()

	// If duplicate suppression is broken, a second ping enters while the first is blocked.
	select {
	case <-pingEntered:
		t.Fatal("duplicate Bond should not issue a second ping while first ping is blocked")
	default:
	}

	close(pingBlocker) // release the in-flight ping
	wg.Wait()

	// The bonding map must be cleared after completion.
	tab.bondmu.Lock()
	assert.Nil(t, tab.bonding[n.ID], "bonding task should be removed after completion")
	tab.bondmu.Unlock()

	// Both callers should receive the same (nil) error.
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Only one actual ping must have been issued.
	assert.Equal(t, int32(1), udp.pingCnt.Load())
}

func Test_Table_Bond_recentlyBonded(t *testing.T) {
	udp := newMockTransport()
	tab := newTestTable2(t, udp)

	n := newTestNode(t, NodeTypeEN)
	// Simulate a previous successful bond: write bond time + node + findFails=0 into DB.
	tab.recordBonded(n)

	// canAssumeBonded requires node to be stored in DB too.
	require.NotNil(t, tab.db.node(n.ID))
	require.True(t, tab.db.hasBond(n.ID))
	require.Equal(t, 0, tab.db.findFails(n.ID))

	err := tab.Bond(false, n)

	require.NoError(t, err)
	// ping must not be called because the node is assumed to be bonded.
	assert.Equal(t, int32(0), udp.pingCnt.Load(), "ping should be skipped for a recently bonded node")
}

// mockTransport is a transport whose ping and findnode behaviour is configurable per NodeID.
type mockTransport struct {
	// pingFn, if non-nil, is called instead of the per-node error map.
	// Useful for tests that need blocking, concurrency signals, etc.
	pingFn func(NodeID, *net.UDPAddr) error
	// findnodeFn, if non-nil, overrides the default (return nil, nil).
	findnodeFn func(NodeID, *net.UDPAddr, NodeID, NodeType, int) ([]*Node, error)

	mu      sync.Mutex
	pingErr map[NodeID]error // nil by default, indicating success

	pingCnt atomic.Int32 // total ping() invocations
}

func newMockTransport() *mockTransport {
	return &mockTransport{pingErr: make(map[NodeID]error)}
}

// setPingErr configures the error returned by ping() for the given node.
func (m *mockTransport) setPingErr(id NodeID, err error) {
	m.mu.Lock()
	m.pingErr[id] = err
	m.mu.Unlock()
}

func (m *mockTransport) ping(toid NodeID, toaddr *net.UDPAddr) error {
	m.pingCnt.Add(1)
	if m.pingFn != nil {
		return m.pingFn(toid, toaddr)
	}
	m.mu.Lock()
	err := m.pingErr[toid]
	m.mu.Unlock()
	return err
}

func (m *mockTransport) waitping(from NodeID, fromIP net.IP) error { return nil }
func (m *mockTransport) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID, nType NodeType, max int) ([]*Node, error) {
	if m.findnodeFn != nil {
		return m.findnodeFn(toid, toaddr, target, nType, max)
	}
	return nil, nil
}
func (m *mockTransport) close() {}
