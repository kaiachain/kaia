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
	"encoding/binary"
	"errors"
	"net"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testNodeID(i uint32) discover.NodeID {
	var id discover.NodeID
	binary.BigEndian.PutUint32(id[:], i)
	return id
}

func testNode(i uint32, ip string, nType discover.NodeType) *discover.Node {
	return discover.NewNode(testNodeID(i), net.ParseIP(ip), 30303, 30303, nil, nType)
}

func testNodeWithAddress(ip string, nType discover.NodeType) (*discover.Node, common.Address) {
	key := newkey()
	return discover.NewNode(discover.PubkeyID(&key.PublicKey), net.ParseIP(ip), 30303, 30303, nil, nType),
		crypto.PubkeyToAddress(key.PublicKey)
}

func hasNodeID(nodes []*discover.Node, id discover.NodeID) bool {
	return slices.ContainsFunc(nodes, func(n *discover.Node) bool {
		return n.ID == id
	})
}

// Start and Close are idempotent.
func TestDialSched_Lifecycle_StartCloseTwice(t *testing.T) {
	var (
		ds = NewDialSched(DialConfig{
			selfID: testNodeID(100),
		}, nil, nil)
		done = make(chan struct{})
	)

	ds.Start()
	ds.Start() // second Start() is no-op
	assert.True(t, ds.started.Load(), "started should remain true after Start")

	go func() {
		ds.Close()
		ds.Close() // second Close() is no-op
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		require.Fail(t, "close should return quickly even when called twice")
	}

	assert.True(t, ds.closed.Load(), "closed should be true after Close")
}

// Rejects incomplete static nodes.
func TestDialSched_StaticNode_RejectIncomplete(t *testing.T) {
	var (
		incomplete = &discover.Node{ID: testNodeID(1), NType: discover.NodeTypePN}
		complete   = testNode(2, "10.0.0.2", discover.NodeTypePN)

		ds = NewDialSched(DialConfig{
			staticNodes: []*discover.Node{incomplete, complete},
		}, nil, nil)
	)

	assert.False(t, ds.static.contains(incomplete.ID), "incomplete static node should not be added")
	assert.True(t, ds.static.contains(complete.ID), "complete static node should be added")
}

// Correctly accounts for static node failures.
func TestDialSched_StaticNode_Retry(t *testing.T) {
	staticNode := testNode(1, "10.0.0.1", discover.NodeTypeEN)
	ds := NewDialSched(DialConfig{
		staticNodes: []*discover.Node{staticNode},
	}, nil, nil)

	// Disconnect does not remove the static node.
	for i := 0; i < dialMaxRetries+2; i++ {
		ds.OnPeerDisconnected(staticNode.ID, staticNode.NType)
	}
	assert.True(t, ds.static.contains(staticNode.ID), "static node must not be removed by disconnection events only")

	// But Failure does remove the static node.
	for i := 0; i < dialMaxRetries+2; i++ {
		ds.markDialFailure(staticNode.ID)
	}
	assert.False(t, ds.static.contains(staticNode.ID), "static node must be removed by failure events")
}

// Special feature: Removing and re-adding a static node re-dials immediately.
func TestDialSched_StaticNode_RedialImmediatelyAfterReAdd(t *testing.T) {
	node := testNode(51, "10.0.0.51", discover.NodeTypePN)
	ds := NewDialSched(DialConfig{
		staticNodes: []*discover.Node{node},
	}, nil, nil)

	ds.dialBackoff[node.ID] = time.Now().Add(time.Minute)
	ds.connFails[node.ID] = 2

	ds.RemoveStatic(node.ID)
	ds.AddStatic(node)

	_, hasBackoff := ds.dialBackoff[node.ID]
	_, hasConnFail := ds.connFails[node.ID]
	assert.False(t, hasBackoff, "removeStatic should clear dialBackoff")
	assert.False(t, hasConnFail, "removeStatic should clear connFails")
	assert.True(t, ds.shouldDial(node), "re-added static node should be dialable immediately")
}

// Updates accounting for inbound and outbound peers.
func TestDialSched_OnPeerConnected(t *testing.T) {
	var (
		candidate = testNode(1, "10.0.0.1", discover.NodeTypeEN)
		inbound   = testNode(2, "10.0.0.2", discover.NodeTypePN)
		outbound  = testNode(3, "10.0.0.3", discover.NodeTypePN)
		tab       = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{
				discover.NodeTypeEN: {candidate},
			},
		}
		ds = NewDialSched(DialConfig{
			connTarget: map[discover.NodeType]int{
				discover.NodeTypeEN: 1,
			},
		}, tab, nil)
	)

	ds.OnPeerConnected(inbound.ID, inbound.NType, true)
	assert.False(t, ds.connectedOutbound.contains(inbound.ID), "inbound peer should not be counted as connected outbound")
	assert.True(t, ds.connectedAll.contains(inbound.ID), "inbound peer should be counted as connected")
	cands, _ := ds.getCandidates()
	require.Len(t, cands, 1, "expected one outbound EN candidate when only inbound legacy PN exists")
	assert.Equal(t, candidate.ID, cands[0].ID)

	ds.OnPeerConnected(outbound.ID, outbound.NType, false)
	assert.True(t, ds.connectedOutbound.contains(outbound.ID), "outbound peer should be counted as connected outbound")
	assert.True(t, ds.connectedAll.contains(outbound.ID), "outbound peer should be counted as connected")
	cands, _ = ds.getCandidates()
	require.Len(t, cands, 1, "expected outbound PN not to fulfill exact EN target")
	assert.Equal(t, candidate.ID, cands[0].ID)
}

// Correctly enforces discovery table refresh throttling.
func TestDialSched_ShouldRefresh_RefreshBackoff(t *testing.T) {
	var (
		tab = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{},
		}
		ds = NewDialSched(DialConfig{
			selfID: testNodeID(1),
		}, tab, nil)
		resCh = make(chan struct{}, 1)
	)

	assert.True(t, ds.shouldRefresh(), "expected initial refresh to be allowed")
	ds.refreshOnce(resCh)
	select {
	case <-resCh:
	case <-time.After(2 * time.Second):
		require.Fail(t, "refreshOnce did not finish in time")
	}
	assert.Equal(t, 1, tab.refreshCalls, "expected one Refresh call")

	// Forcibly set the refresh backoff to be in the past.
	assert.False(t, ds.shouldRefresh(), "expected refresh to be throttled immediately after refreshOnce")
	ds.refreshBackoff = time.Now().Add(-refreshBackoff)

	assert.True(t, ds.shouldRefresh(), "expected refresh to be allowed after backoff expiry")
	ds.refreshOnce(resCh)
	select {
	case <-resCh:
	case <-time.After(2 * time.Second):
		require.Fail(t, "second refreshOnce did not finish in time")
	}
	assert.Equal(t, 2, tab.refreshCalls, "expected two Refresh calls total")
}

// Static nodes are always candidates, even if connTarget is met.
func TestDialSched_Candidates_StaticNodes(t *testing.T) {
	var (
		staticNode = testNode(61, "10.0.0.61", discover.NodeTypePN)
		ds         = NewDialSched(DialConfig{
			staticNodes: []*discover.Node{staticNode},
			connTarget: map[discover.NodeType]int{
				discover.NodeTypePN: 0, // Even if connTarget is zero, always attempt to connect to static nodes.
			},
		}, nil, nil)
	)

	cands, _ := ds.getCandidates()
	assert.Equal(t, staticNode.ID, cands[0].ID, "static node should be the only candidate")
}

// Static nodes dialed in the same batch consume the outbound demand of their type.
func TestDialSched_GetCandidates_StaticNodesConsumeDemand(t *testing.T) {
	var (
		staticCNs = []*discover.Node{
			testNode(61, "10.0.0.61", discover.NodeTypeCN),
			testNode(62, "10.0.0.62", discover.NodeTypeCN),
		}
		dynamicCNs = []*discover.Node{
			testNode(201, "10.0.0.201", discover.NodeTypeCN),
			testNode(202, "10.0.0.202", discover.NodeTypeCN),
		}
		tab = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{
				discover.NodeTypeCN: dynamicCNs,
			},
		}
	)

	filled := NewDialSched(DialConfig{
		staticNodes: staticCNs,
		connTarget:  map[discover.NodeType]int{discover.NodeTypeCN: 2},
	}, tab, nil)
	cands, _ := filled.getCandidates()
	assert.Len(t, cands, len(staticCNs), "pending static CNs alone should fill the CN target")
	for _, dynamic := range dynamicCNs {
		assert.False(t, hasNodeID(cands, dynamic.ID), "no dynamic CN is needed while statics fill the target")
	}

	// A connected static is already counted in connectedOutbound, so it must not reserve a second slot.
	connected := NewDialSched(DialConfig{
		staticNodes: staticCNs[:1],
		connTarget:  map[discover.NodeType]int{discover.NodeTypeCN: 2},
	}, tab, nil)
	connected.OnPeerConnected(staticCNs[0].ID, staticCNs[0].NType, false)
	cands, _ = connected.getCandidates()
	assert.True(t, hasNodeID(cands, dynamicCNs[0].ID), "the demand left by a connected static should pull a dynamic CN")
}

// Dynamic nodes are candidates if connTarget is not met.
func TestDialSched_GetCandidates_FromDiscovery(t *testing.T) {
	var (
		cn  = testNode(201, "10.0.0.201", discover.NodeTypeCN)
		tab = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{
				discover.NodeTypeCN: {cn},
			},
		}
		ds = NewDialSched(DialConfig{
			connTarget: map[discover.NodeType]int{
				discover.NodeTypeCN: 1,
			},
		}, tab, nil)
	)

	cands, needRefresh := ds.getCandidates()

	assert.True(t, hasNodeID(cands, cn.ID), "CN candidates should include discovery result")
	assert.False(t, needRefresh, "No need to refresh when candidates are enough")
}

func TestDialSched_GetCandidates_ENTargetDoesNotQueryLegacyPN(t *testing.T) {
	var (
		pn  = testNode(201, "10.0.0.201", discover.NodeTypePN)
		en  = testNode(202, "10.0.0.202", discover.NodeTypeEN)
		tab = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{
				discover.NodeTypePN: {pn},
				discover.NodeTypeEN: {en},
			},
		}
		ds = NewDialSched(DialConfig{
			connTarget: map[discover.NodeType]int{
				discover.NodeTypeEN: 1,
			},
		}, tab, nil)
	)

	cands, needRefresh := ds.getCandidates()

	assert.True(t, hasNodeID(cands, en.ID), "EN candidates should fill the EN target")
	assert.False(t, hasNodeID(cands, pn.ID), "upgraded nodes should not dial PN candidates dynamically")
	assert.False(t, needRefresh)
}

func TestDialSched_GetCandidates_CNPeerAllowlist(t *testing.T) {
	allowed, allowedAddr := testNodeWithAddress("10.0.0.201", discover.NodeTypeCN)
	blocked, _ := testNodeWithAddress("10.0.0.202", discover.NodeTypeCN)
	tab := &mockTable{
		nodesByType: map[discover.NodeType][]*discover.Node{
			discover.NodeTypeCN: {allowed, blocked},
		},
	}
	ds := NewDialSched(DialConfig{
		selfType: discover.NodeTypeCN,
		connTarget: map[discover.NodeType]int{
			discover.NodeTypeCN: 2,
		},
	}, tab, nil)
	ds.SetCNPeers([]common.Address{allowedAddr})

	cands, needRefresh := ds.getCandidates()

	assert.True(t, hasNodeID(cands, allowed.ID), "CN peer in allowlist should remain dialable")
	assert.False(t, hasNodeID(cands, blocked.ID), "CN peer outside allowlist should be filtered")
	assert.True(t, needRefresh, "filtering below target should request discovery refresh")
}

func TestDialSched_GetCandidates_CNPeerAllowlistNilAndEmpty(t *testing.T) {
	cn, _ := testNodeWithAddress("10.0.0.201", discover.NodeTypeCN)
	tab := &mockTable{
		nodesByType: map[discover.NodeType][]*discover.Node{
			discover.NodeTypeCN: {cn},
		},
	}

	nilFilter := NewDialSched(DialConfig{
		selfType: discover.NodeTypeCN,
		connTarget: map[discover.NodeType]int{
			discover.NodeTypeCN: 1,
		},
	}, tab, nil)
	nilFilter.SetCNPeers(nil)
	cands, needRefresh := nilFilter.getCandidates()
	assert.True(t, hasNodeID(cands, cn.ID), "nil allowlist should disable CN filtering")
	assert.False(t, needRefresh)

	emptyFilter := NewDialSched(DialConfig{
		selfType: discover.NodeTypeCN,
		connTarget: map[discover.NodeType]int{
			discover.NodeTypeCN: 1,
		},
	}, tab, nil)
	emptyFilter.SetCNPeers([]common.Address{})
	cands, needRefresh = emptyFilter.getCandidates()
	assert.False(t, hasNodeID(cands, cn.ID), "empty allowlist should reject dynamic CN candidates")
	assert.True(t, needRefresh)
}

func TestDialSched_GetCandidates_CNPeerAllowlistStaticBypass(t *testing.T) {
	staticCN, _ := testNodeWithAddress("10.0.0.201", discover.NodeTypeCN)
	ds := NewDialSched(DialConfig{
		selfType:    discover.NodeTypeCN,
		staticNodes: []*discover.Node{staticCN},
		connTarget: map[discover.NodeType]int{
			discover.NodeTypeCN: 1,
		},
	}, &mockTable{}, nil)
	ds.SetCNPeers([]common.Address{})

	cands, _ := ds.getCandidates()

	assert.True(t, hasNodeID(cands, staticCN.ID), "static outbound CN should bypass dynamic CN filtering")
}

func TestDialSched_GetCandidates_CNPeerAllowlistOnlyAppliesToCN(t *testing.T) {
	cn, _ := testNodeWithAddress("10.0.0.201", discover.NodeTypeCN)
	tab := &mockTable{
		nodesByType: map[discover.NodeType][]*discover.Node{
			discover.NodeTypeCN: {cn},
		},
	}
	ds := NewDialSched(DialConfig{
		selfType: discover.NodeTypeEN,
		connTarget: map[discover.NodeType]int{
			discover.NodeTypeCN: 1,
		},
	}, tab, nil)
	ds.SetCNPeers([]common.Address{})

	cands, needRefresh := ds.getCandidates()

	assert.True(t, hasNodeID(cands, cn.ID), "EN dynamic CN dials should not use CN peer allowlist")
	assert.False(t, needRefresh)
}

func TestDialSched_GetCandidates_PNTargetQueriesPN(t *testing.T) {
	var (
		pn  = testNode(201, "10.0.0.201", discover.NodeTypePN)
		en  = testNode(202, "10.0.0.202", discover.NodeTypeEN)
		tab = &mockTable{
			nodesByType: map[discover.NodeType][]*discover.Node{
				discover.NodeTypePN: {pn},
				discover.NodeTypeEN: {en},
			},
		}
		ds = NewDialSched(DialConfig{
			connTarget: map[discover.NodeType]int{
				discover.NodeTypePN: 1,
			},
		}, tab, nil)
	)

	cands, needRefresh := ds.getCandidates()

	assert.True(t, hasNodeID(cands, pn.ID), "PN target should include PN candidates")
	assert.False(t, hasNodeID(cands, en.ID), "PN target should not query EN candidates")
	assert.False(t, needRefresh)
}

// Skips discovery when connTarget is met.
func TestDialSched_GetCandidates_SkipDiscovery(t *testing.T) {
	var (
		pn  = testNode(211, "10.0.0.211", discover.NodeTypePN)
		en  = testNode(212, "10.0.0.212", discover.NodeTypeEN)
		tab = &mockTable{}
		ds  = NewDialSched(DialConfig{
			connTarget: map[discover.NodeType]int{
				discover.NodeTypePN: 1,
				discover.NodeTypeEN: 1,
			},
		}, tab, nil)
	)

	// Already fulfilled the connTarget demand.
	ds.connectedOutbound.add(pn)
	ds.connectedOutbound.add(en)

	cands, needRefresh := ds.getCandidates()

	assert.Len(t, cands, 0) // Don't fetch dynamic nodes.
	assert.False(t, needRefresh, "needRefresh should remain false when demand is already satisfied")
}

// Various conditions that determine whether a node should be dialed.
func TestDialSched_ShouldDial(t *testing.T) {
	var (
		ds = NewDialSched(DialConfig{
			selfID:      testNodeID(900),
			netrestrict: mustParseNetlist(t, "10.0.0.0/8"),
		}, nil, nil)

		nodeSelf       = testNode(900, "10.0.0.9", discover.NodeTypeEN)
		nodeRestricted = testNode(901, "11.0.0.1", discover.NodeTypeEN)
		nodeConnected  = testNode(902, "10.0.0.2", discover.NodeTypeEN)
		nodeDialing    = testNode(903, "10.0.0.3", discover.NodeTypeEN)
		nodeBackoff    = testNode(904, "10.0.0.4", discover.NodeTypeEN)
		nodeAllowed    = testNode(905, "10.0.0.5", discover.NodeTypeEN)

		cases = []struct {
			name string
			node *discover.Node
			want bool
		}{
			{name: "nil", node: nil, want: false},
			{name: "incomplete", node: &discover.Node{ID: testNodeID(906), NType: discover.NodeTypeEN}, want: false},
			{name: "self", node: nodeSelf, want: false},
			{name: "netrestrict", node: nodeRestricted, want: false},
			{name: "connected", node: nodeConnected, want: false},
			{name: "dialing", node: nodeDialing, want: false},
			{name: "backoff", node: nodeBackoff, want: false},
			{name: "allowed", node: nodeAllowed, want: true},
		}
	)

	ds.connectedAll.add(nodeConnected)
	ds.dialing.add(nodeDialing)
	ds.dialBackoff[nodeBackoff.ID] = time.Now().Add(time.Minute)

	for _, tc := range cases {
		assert.Equal(t, tc.want, ds.shouldDial(tc.node), tc.name)
	}
}

// Launches dial tasks and waits for them to complete.
func TestDialSched_LaunchDialTasks(t *testing.T) {
	var (
		staticNode = testNode(11, "10.0.0.11", discover.NodeTypeEN)
		dynNode    = testNode(21, "10.0.0.21", discover.NodeTypeEN)
		md         = &mockDialer{dialErr: make(map[discover.NodeID]error)}
		mb         = &mockSetupBackend{setupErr: make(map[discover.NodeID]error)}
		ds         = NewDialSched(DialConfig{
			staticNodes: []*discover.Node{staticNode},
		}, nil, mb)
		resCh = make(chan struct{}, 2)
	)
	ds.dialer = md
	mb.dialsched = ds

	// Launch and wait for both dial tasks to complete.
	ds.launchDialTasks([]*discover.Node{staticNode, dynNode}, resCh)
	for i := 0; i < 2; i++ {
		select {
		case <-resCh:
		case <-time.After(2 * time.Second):
			require.Fail(t, "launchDialTasks did not finish in time")
		}
	}

	// Inspect the SetupConn() arguments (could be in different order due to concurrent execution)
	require.Len(t, mb.calls, 2)
	if mb.calls[0].id == staticNode.ID {
		assert.Equal(t, setupConnArg{id: staticNode.ID, flags: staticDialedConn, port: 0}, mb.calls[0])
		assert.Equal(t, setupConnArg{id: dynNode.ID, flags: dynDialedConn, port: 0}, mb.calls[1])
	} else {
		assert.Equal(t, setupConnArg{id: dynNode.ID, flags: dynDialedConn, port: 0}, mb.calls[0])
		assert.Equal(t, setupConnArg{id: staticNode.ID, flags: staticDialedConn, port: 0}, mb.calls[1])
	}

	// Inspect the internal states
	assert.Zero(t, ds.dialing.len(), "expected no nodes to be in the dialing state")
	assert.True(t, ds.connectedAll.contains(staticNode.ID))
	assert.True(t, ds.connectedAll.contains(dynNode.ID))
	assert.True(t, ds.connectedOutbound.contains(staticNode.ID))
	assert.True(t, ds.connectedOutbound.contains(dynNode.ID))
}

// Updates accounting for dial failures.
func TestDialSched_Dial_DialFailure(t *testing.T) {
	var (
		staticNode = testNode(301, "10.0.1.31", discover.NodeTypePN)
		md         = &mockDialer{dialErr: make(map[discover.NodeID]error)}
		mb         = &mockSetupBackend{
			setupErr: map[discover.NodeID]error{
				staticNode.ID: errors.New("mock setup failure"),
			},
		}
		ds = NewDialSched(DialConfig{
			staticNodes: []*discover.Node{staticNode},
		}, nil, mb)
	)
	ds.dialer = md
	mb.dialsched = ds

	ds.dialOnce(staticNode, staticDialedConn)

	require.Len(t, mb.calls, 1)
	assert.Equal(t, staticNode.ID, mb.calls[0].id)
	assert.Equal(t, staticDialedConn, mb.calls[0].flags)
	assert.Equal(t, 1, ds.connFails[staticNode.ID], "static dial failure should increment connFails")
	assert.False(t, ds.connectedAll.contains(staticNode.ID), "static node should not be counted as connected")
	assert.False(t, ds.connectedOutbound.contains(staticNode.ID), "static node should not be counted as connected")
}

// Special feature: Retries dial after upgrading from single-channel to multi-channel.
func TestDialSched_Dial_RetryOnErrUpdateDial(t *testing.T) {
	var (
		node = testNode(321, "10.0.1.32", discover.NodeTypePN)
		md   = &mockDialer{dialErr: make(map[discover.NodeID]error)}
		mb   = &updateToMultiBackend{
			ports: []uint16{30304, 30305},
		}
		ds = NewDialSched(DialConfig{}, nil, mb)
	)
	ds.dialer = md
	mb.dialsched = ds

	ds.dialOnce(node, dynDialedConn)

	require.Len(t, mb.calls, 3, "should retry with DialMulti after errUpdateDial")

	// First SetupConn resulted in errUpdateDial.
	assert.Equal(t, node.ID, mb.calls[0].id)
	assert.Equal(t, dynDialedConn, mb.calls[0].flags)

	// Second and third SetupConns with the updated ports.
	assert.Equal(t, node.ID, mb.calls[1].id)
	assert.Equal(t, node.ID, mb.calls[2].id)
	assert.Equal(t, dynDialedConn, mb.calls[1].flags)
	assert.Equal(t, dynDialedConn, mb.calls[2].flags)
	assert.Equal(t, uint16(0), mb.calls[1].port)
	assert.Equal(t, uint16(1), mb.calls[2].port)

	// Eventually connected.
	assert.True(t, ds.connectedAll.contains(node.ID))
	assert.True(t, ds.connectedOutbound.contains(node.ID))
}

// Correctly dials multi-channel peers.
func TestDialSched_DialMulti(t *testing.T) {
	var (
		multiNode = discover.NewNode(testNodeID(401), net.ParseIP("10.0.0.41"), 30303, 30303, []uint16{30304, 30305}, discover.NodeTypeEN)
		md        = &mockDialer{dialErr: make(map[discover.NodeID]error)}
		mb        = &mockSetupBackend{setupErr: make(map[discover.NodeID]error)}
		ds        = NewDialSched(DialConfig{}, nil, mb)
	)
	ds.dialer = md
	mb.dialsched = ds

	ds.dialOnce(multiNode, dynDialedConn)

	require.Len(t, mb.calls, len(multiNode.TCPs), "SetupConn should be called for each dialed connection")
	for i := range mb.calls {
		assert.Equal(t, multiNode.ID, mb.calls[i].id)
		assert.Equal(t, dynDialedConn, mb.calls[i].flags)
		assert.Equal(t, uint16(i), mb.calls[i].port, "PortOrder should increment per connection")
	}
}

func mustParseNetlist(t *testing.T, s string) *netutil.Netlist {
	t.Helper()
	nl, err := netutil.ParseNetlist(s)
	require.NoError(t, err)
	return nl
}

// TestDialSched_Close_BlockedDial verifies that Close() returns promptly even when
// the dialOnce() is blocked at dialer.Dial(). Close() should not wait for the in-flight dialOnce.
func TestDialSched_Close_BlockedDial(t *testing.T) {
	invokedCh := make(chan struct{}, 1)
	unblockCh := make(chan struct{})

	node := testNode(1, "10.0.0.1", discover.NodeTypeEN)
	ds := NewDialSched(DialConfig{
		staticNodes: []*discover.Node{node},
		dialer:      &blockingDialer{invokedCh: invokedCh, unblockCh: unblockCh}, // Dial() blocks until unblockCh.
	}, nil, &rejectBackend{}) // SetupConn() always immediatelyreturns an error.
	ds.Start()

	select {
	case <-invokedCh:
	case <-time.After(2 * time.Second):
		require.Fail(t, "Dial() was not invoked")
	}

	closeDone := make(chan struct{})
	go func() {
		ds.Close()
		close(closeDone)
	}()
	select {
	case <-closeDone:
	case <-time.After(2 * time.Second):
		require.Fail(t, "Close() must return quickly without waiting for the in-flight Dial()")
	}

	// Unblock so the blocked goroutine can finish.
	// Note that Close() didn't wait for the dialOnce goroutine to finish.
	close(unblockCh)
}

// TestDialSched_Close_BlockedSetupConn verifies that Close() returns promptly even when
// the dialOnce() is blocked at backend.SetupConn(). Close() should not wait for the in-flight dialOnce.
func TestDialSched_Close_BlockedSetupConn(t *testing.T) {
	invokedCh := make(chan struct{}, 1)
	unblockCh := make(chan struct{})

	node := testNode(2, "10.0.0.2", discover.NodeTypeEN)
	ds := NewDialSched(DialConfig{
		staticNodes: []*discover.Node{node},
		dialer:      &mockDialer{dialErr: make(map[discover.NodeID]error)}, // Dial() succeeds immediately
	}, nil, &blockingBackend{invokedCh: invokedCh, unblockCh: unblockCh}) // SetupConn() blocks until unblockCh.
	ds.Start()

	select {
	case <-invokedCh:
	case <-time.After(2 * time.Second):
		require.Fail(t, "SetupConn was not invoked")
	}

	closeDone := make(chan struct{})
	go func() {
		ds.Close()
		close(closeDone)
	}()
	select {
	case <-closeDone:
	case <-time.After(2 * time.Second):
		require.Fail(t, "Close() must return without waiting for the in-flight SetupConn()")
	}

	// Unblock so the blocked goroutine can finish.
	// Note that Close() didn't wait for the dialOnce goroutine to finish.
	close(unblockCh)
}

//// Discovery table mocks.

// Generic mock discovery table.
type mockTable struct {
	nodesByType  map[discover.NodeType][]*discover.Node
	refreshCalls int
}

func (m *mockTable) RandomNodes(buf []*discover.Node, nType discover.NodeType) int {
	return copy(buf, m.nodesByType[nType])
}

func (m *mockTable) Refresh() {
	m.refreshCalls++
}

func (m *mockTable) Close() {
}

//// Dialer (Dial, DialMulti) mocks.

type mockDialer struct {
	mu      sync.Mutex
	dialErr map[discover.NodeID]error // injected errors, discriminated by NodeID.
}

func (m *mockDialer) Dial(dest *discover.Node) (net.Conn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.dialErr[dest.ID]; err != nil {
		return nil, err
	}
	c1, c2 := net.Pipe()
	_ = c2.Close()
	return c1, nil
}

func (m *mockDialer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.dialErr[dest.ID]; err != nil {
		return nil, err
	}
	conns := make([]net.Conn, 0, len(dest.TCPs))
	for range dest.TCPs {
		c1, c2 := net.Pipe()
		_ = c2.Close()
		conns = append(conns, c1)
	}
	return conns, nil
}

// blockingDialer blocks inside Dial() until unblockCh is closed.
// It signals the first call on invokedCh, so the test can inspect the invocation.
type blockingDialer struct {
	invokedCh chan struct{} // receives once when the first Dial() is entered
	unblockCh chan struct{} // Dial() returns when closed
}

func (m *blockingDialer) Dial(dest *discover.Node) (net.Conn, error) {
	select {
	case m.invokedCh <- struct{}{}:
	default:
	}
	<-m.unblockCh
	return nil, errors.New("dial unblocked")
}

func (m *blockingDialer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	select {
	case m.invokedCh <- struct{}{}:
	default:
	}
	<-m.unblockCh
	return nil, errors.New("dial unblocked")
}

//// Backend (SetupConn) mocks.

type setupConnArg struct {
	id    discover.NodeID
	flags connFlag
	port  uint16
}

type mockSetupBackend struct {
	mu        sync.Mutex
	dialsched *DialSched
	setupErr  map[discover.NodeID]error // injected errors, discriminated by NodeID.
	calls     []setupConnArg            // recorded SetupConn calls.
}

func (m *mockSetupBackend) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	if fd != nil {
		_ = fd.Close()
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	call := setupConnArg{flags: flags}
	if dialDest != nil {
		call.id = dialDest.ID
		call.port = dialDest.PortOrder
		if err := m.setupErr[dialDest.ID]; err != nil {
			m.calls = append(m.calls, call)
			return err
		}
	}
	m.calls = append(m.calls, call)
	m.dialsched.markPeerConnected(dialDest.ID, dialDest.NType, false)
	return nil
}

// rejectBackend always returns an error from SetupConn. Used when the test only
// cares about the dialing side, not the connection setup side.
type rejectBackend struct{}

func (m *rejectBackend) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	if fd != nil {
		fd.Close()
	}
	return errors.New("rejectBackend: connection rejected")
}

// blockingBackend blocks inside SetupConn() until unblockCh is closed.
// It signals the first call on invokedCh, so the test can inspect the invocation.
type blockingBackend struct {
	invokedCh chan struct{} // receives once when the first SetupConn() is entered
	unblockCh chan struct{} // SetupConn() returns when closed
}

func (m *blockingBackend) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	if fd != nil {
		fd.Close()
	}
	select {
	case m.invokedCh <- struct{}{}:
	default:
	}
	<-m.unblockCh
	return errors.New("SetupConn unblocked")
}

// A SetupConn backend that returns errUpdateDial, telling me (the dialsched) to use different ports.
// Returns errUpdateDial at first call, then succeeds afterwards.
type updateToMultiBackend struct {
	mu         sync.Mutex
	dialsched  *DialSched
	calledOnce bool
	ports      []uint16
	calls      []setupConnArg
}

func (m *updateToMultiBackend) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	if fd != nil {
		_ = fd.Close()
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	call := setupConnArg{flags: flags}
	if dialDest != nil {
		call.id = dialDest.ID
		call.port = dialDest.PortOrder
	}
	m.calls = append(m.calls, call)

	if !m.calledOnce {
		m.calledOnce = true
		if dialDest != nil {
			dialDest.TCPs = append([]uint16(nil), m.ports...)
		}
		return errUpdateDial
	}
	m.dialsched.markPeerConnected(dialDest.ID, dialDest.NType, false)
	return nil
}
