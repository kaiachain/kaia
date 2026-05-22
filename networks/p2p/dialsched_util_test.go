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
	"net"
	"testing"

	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/stretchr/testify/assert"
)

func typedSetNode(id uint32, ip string, nType discover.NodeType) *discover.Node {
	var nodeID discover.NodeID
	nodeID[0] = byte(id >> 24)
	nodeID[1] = byte(id >> 16)
	nodeID[2] = byte(id >> 8)
	nodeID[3] = byte(id)
	return discover.NewNode(nodeID, net.ParseIP(ip), 30303, 30303, nil, nType)
}

func hasNode(nodes []*discover.Node, id discover.NodeID) bool {
	for _, n := range nodes {
		if n != nil && n.ID == id {
			return true
		}
	}
	return false
}

func TestTypedNodeSet_BasicOps(t *testing.T) {
	set := newTypedNodeSet()
	n1 := typedSetNode(1, "10.0.0.1", discover.NodeTypePN)
	n2 := typedSetNode(2, "10.0.0.2", discover.NodeTypeEN)

	set.add(n1)
	set.add(n2)

	assert.Equal(t, 2, set.len())
	assert.Equal(t, 2, set.count(discover.NodeTypePN))
	assert.Equal(t, 2, set.count(discover.NodeTypeEN))
	assert.Equal(t, 0, set.count(discover.NodeTypeCN))
	assert.True(t, set.contains(n1.ID))
	assert.True(t, set.contains(n2.ID))
	assert.Equal(t, n1, set.get(n1.ID))
	assert.Equal(t, n2, set.get(n2.ID))

	all := set.all()
	assert.Len(t, all, 2)
	assert.True(t, hasNode(all, n1.ID))
	assert.True(t, hasNode(all, n2.ID))
}

func TestTypedNodeSet_CountTreatsPNAsEN(t *testing.T) {
	set := newTypedNodeSet()
	pn := typedSetNode(1, "10.0.0.1", discover.NodeTypePN)
	en := typedSetNode(2, "10.0.0.2", discover.NodeTypeEN)
	cn := typedSetNode(3, "10.0.0.3", discover.NodeTypeCN)

	set.add(pn)
	set.add(en)
	set.add(cn)

	assert.Equal(t, 2, set.count(discover.NodeTypePN))
	assert.Equal(t, 2, set.count(discover.NodeTypeEN))
	assert.Equal(t, 1, set.count(discover.NodeTypeCN))
}

func TestTypedNodeSet_RetypeSameID(t *testing.T) {
	set := newTypedNodeSet()
	n1 := typedSetNode(1, "10.0.0.1", discover.NodeTypePN)
	n2 := typedSetNode(1, "10.0.0.1", discover.NodeTypeEN)

	set.add(n1)
	set.add(n2)

	assert.Equal(t, 1, set.len())
	assert.Equal(t, 0, set.counts[discover.NodeTypePN])
	assert.Equal(t, 1, set.count(discover.NodeTypeEN))
	assert.Equal(t, n2, set.get(n1.ID))
}

func TestTypedNodeSet_RemoveNoopAndNilAdd(t *testing.T) {
	set := newTypedNodeSet()
	n := typedSetNode(3, "10.0.0.3", discover.NodeTypeCN)
	var unknownID discover.NodeID
	unknownID[0] = 99

	set.add(nil) // no-op
	set.remove(unknownID)

	set.add(n)
	assert.Equal(t, 1, set.len())
	assert.Equal(t, 1, set.count(discover.NodeTypeCN))

	set.remove(n.ID)
	assert.Equal(t, 0, set.len())
	assert.Equal(t, 0, set.count(discover.NodeTypeCN))
	assert.False(t, set.contains(n.ID))
	assert.Nil(t, set.get(n.ID))
}
