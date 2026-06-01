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
	"math"
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

func Test_Table_getDiscoverTargets(t *testing.T) {
	assert.Equal(t, map[NodeType]int{NodeTypeCN: 100, NodeTypeEN: 1, NodeTypeBN: 3},
		getDiscoverTargets(&Config{NodeType: NodeTypeCN}))
	assert.Equal(t, map[NodeType]int{NodeTypeCN: 100, NodeTypeEN: math.MaxInt32, NodeTypeBN: 3},
		getDiscoverTargets(&Config{NodeType: NodeTypeEN}))
	assert.Equal(t, getDiscoverTargets(&Config{NodeType: NodeTypeEN}),
		getDiscoverTargets(&Config{NodeType: NodeTypePN}))
	assert.Equal(t, map[NodeType]int{NodeTypeCN: 100, NodeTypeBN: 3},
		getDiscoverTargets(&Config{NodeType: NodeTypeBN}))
}

func Test_Table_StartClose(t *testing.T) {
	tab := newTestTable2(t, nil)
	tab.Start()
	tab.Close()
	tab.Close() // idempotent
}
