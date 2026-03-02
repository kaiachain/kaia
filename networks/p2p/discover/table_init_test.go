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
	"testing"

	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Table_getBootnodes_valid(t *testing.T) {
	n1 := newTestBootnode(t, net.IP{10, 0, 0, 1})
	n2 := newTestBootnode(t, net.IP{10, 0, 0, 2})
	cfg := &Config{Bootnodes: []*Node{n1, n2}}

	nursery, err := getBootnodes(cfg)
	require.NoError(t, err)
	require.Len(t, nursery, 2)
	assert.Equal(t, n1.ID, nursery[0].ID)
	assert.Equal(t, n2.ID, nursery[1].ID)
}

func Test_Table_getBootnodes_invalid(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	// UDP port 0 fails validateComplete ("missing UDP port").
	bad := NewNode(PubkeyID(&key.PublicKey), net.IP{127, 0, 0, 1}, 0, 0, nil, NodeTypeBN)

	_, err = getBootnodes(&Config{Bootnodes: []*Node{bad}})
	assert.ErrorIs(t, err, errBadBootnode)
}

func Test_Table_getBootnodes_netRestrict(t *testing.T) {
	restrict, err := netutil.ParseNetlist("10.0.0.0/8")
	require.NoError(t, err)

	inside := newTestBootnode(t, net.IP{10, 0, 0, 1})
	outside := newTestBootnode(t, net.IP{192, 168, 0, 1})
	cfg := &Config{
		Bootnodes:   []*Node{inside, outside},
		NetRestrict: restrict,
	}

	nursery, err := getBootnodes(cfg)
	require.NoError(t, err)
	require.Len(t, nursery, 1, "node outside netrestrict subnet should be filtered")
	assert.Equal(t, inside.ID, nursery[0].ID)
}

func Test_Table_StartClose(t *testing.T) {
	tab := newTestTable2(t, nil)
	tab.Start()
	tab.Close()
	tab.Close() // idempotent
}

//// Test helpers

// newTestNode generates a cryptographically valid Node with a loopback IP.
// The loopback IP is treated as LAN and bypasses the kademliaStorage IP limit.
func newTestNode(t *testing.T, nType NodeType) *Node {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return NewNode(PubkeyID(&key.PublicKey), net.IP{127, 0, 0, 1}, 30303, 30303, nil, nType)
}

// newTestTable2 constructs a Table2 backed by an in-memory node database and
// a no-op UDP transport. tab.Close is registered as a test cleanup, so callers
// do not have to call it explicitly (though doing so a second time is safe).
func newTestTable2(t *testing.T, transport transport) *Table2 {
	t.Helper()
	cfg := &Config{
		Id:         NodeID{},
		Addr:       &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 30303},
		NodeType:   NodeTypeEN,
		NodeDBPath: "", // In-memory temporary database
	}
	if transport == nil {
		transport = &noopTransport{}
	}
	tab, err := newTable2(cfg, transport)
	require.NoError(t, err)
	t.Cleanup(tab.Close)
	return tab
}

// newTestBootnode generates a cryptographically valid Node at the given IP,
// suitable for use as a bootstrap node.
func newTestBootnode(t *testing.T, ip net.IP) *Node {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return NewNode(PubkeyID(&key.PublicKey), ip, 30303, 30303, nil, NodeTypeBN)
}

//// Mock UDP transports

// noopTransport is a minimal no-op transport for Table2 unit tests.
type noopTransport struct{}

func (m *noopTransport) ping(toid NodeID, toaddr *net.UDPAddr) error { return nil }
func (m *noopTransport) waitping(from NodeID, fromIP net.IP) error   { return nil }
func (m *noopTransport) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID, nType NodeType, max int) ([]*Node, error) {
	return nil, nil
}
func (m *noopTransport) close() {}
