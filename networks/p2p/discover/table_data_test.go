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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Table_addNode_typed(t *testing.T) {
	tab := newTestTable2(t, nil)

	en := newTestNode(t, NodeTypeEN)
	tab.addNode(en)
	assert.Equal(t, 1, tab.storages[NodeTypeEN].len(), "EN node should land in EN storage")
	assert.Equal(t, 0, tab.storages[NodeTypeCN].len())
	assert.Equal(t, 0, tab.storages[NodeTypePN].len())
	assert.Equal(t, 0, tab.storages[NodeTypeBN].len())

	cn := newTestNode(t, NodeTypeCN)
	tab.addNode(cn)
	assert.Equal(t, 1, tab.storages[NodeTypeCN].len(), "CN node should land in CN storage")
}

func Test_Table_addNode_BN(t *testing.T) {
	tab := newTestTable2(t, nil)

	bn := newTestBootnode(t, []byte{127, 0, 0, 1})
	tab.addNode(bn)

	// BN is added to every storage simultaneously.
	assert.Equal(t, 1, tab.storages[NodeTypeCN].len())
	assert.Equal(t, 1, tab.storages[NodeTypePN].len())
	assert.Equal(t, 1, tab.storages[NodeTypeEN].len())
	assert.Equal(t, 1, tab.storages[NodeTypeBN].len())
	assert.Equal(t, 4, tab.len())
}

func Test_Table_deleteNode_typed(t *testing.T) {
	tab := newTestTable2(t, nil)

	n := newTestNode(t, NodeTypePN)
	tab.addNode(n)
	require.Equal(t, 1, tab.storages[NodeTypePN].len())

	tab.deleteNode(n)
	assert.Equal(t, 0, tab.storages[NodeTypePN].len())
	// Other storages unaffected.
	assert.Equal(t, 0, tab.storages[NodeTypeEN].len())
}

func Test_Table_deleteNode_BN(t *testing.T) {
	tab := newTestTable2(t, nil)

	bn := newTestBootnode(t, []byte{127, 0, 0, 1})
	tab.addNode(bn)
	require.Equal(t, 4, tab.len())

	tab.deleteNode(bn)
	assert.Equal(t, 0, tab.len(), "BN should be removed from all storages")
}

func Test_Table_RandomNodes_unknownType(t *testing.T) {
	tab := newTestTable2(t, nil)
	buf := make([]*Node, 10)
	n := tab.RandomNodes(buf, NodeTypeUnknown)
	assert.Equal(t, 0, n, "nil storage should return 0")
}

func Test_Table_RandomNodes_happy(t *testing.T) {
	tab := newTestTable2(t, nil)
	for i := 0; i < 5; i++ {
		tab.addNode(newTestNode(t, NodeTypeEN))
	}

	buf := make([]*Node, 10)
	n := tab.RandomNodes(buf, NodeTypeEN)
	assert.Equal(t, 5, n)
	assert.Nil(t, buf[5], "buffer beyond returned count should be untouched")
}

func Test_Table_ClosestNodes_unknownType(t *testing.T) {
	tab := newTestTable2(t, nil)
	result := tab.ClosestNodes(NodeID{}, NodeTypeUnknown, 10)
	assert.Nil(t, result, "nil storage should return nil")
}

func Test_Table_ClosestNodes_happy(t *testing.T) {
	tab := newTestTable2(t, nil)
	for i := 0; i < 5; i++ {
		tab.addNode(newTestNode(t, NodeTypeEN))
	}

	result := tab.ClosestNodes(NodeID{}, NodeTypeEN, 3)
	assert.Len(t, result, 3, "should honour the max cap")

	result = tab.ClosestNodes(NodeID{}, NodeTypeEN, 100)
	assert.Len(t, result, 5, "should return all nodes when max > count")
}

func Test_Table_IsBonded_false(t *testing.T) {
	tab := newTestTable2(t, nil)
	n := newTestNode(t, NodeTypeEN)
	assert.False(t, tab.IsBonded(n.ID), "brand-new node should not be bonded")
}

func Test_Table_IsBonded_true(t *testing.T) {
	tab := newTestTable2(t, nil)
	n := newTestNode(t, NodeTypeEN)
	tab.recordBonded(n)
	assert.True(t, tab.IsBonded(n.ID), "node with recent bond time should be bonded")
}
