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
	"github.com/kaiachain/kaia/networks/p2p/discover"
)

// dialTargets from KIP-311.
var dialTargets = map[discover.NodeType]map[discover.NodeType]int{
	discover.NodeTypeCN: {discover.NodeTypeCN: 100, discover.NodeTypeEN: 1},
	discover.NodeTypeEN: {
		discover.NodeTypeCN: 2,
		discover.NodeTypePN: 2,
		discover.NodeTypeEN: 3,
	},
}

// defaultConnTarget applies KIP-311 dialTargets plus PN compatibility.
func defaultConnTarget(cfg DialConfig) map[discover.NodeType]int {
	targets := dialTargets[discover.EffectiveNodeType(cfg.selfType)]
	targets = cloneTargets(targets)
	if discover.EffectiveNodeType(cfg.selfType) == discover.NodeTypeEN && cfg.maxENToENDialTarget != nil {
		targets[discover.NodeTypeEN] = *cfg.maxENToENDialTarget
	}
	return targets
}

func cloneTargets(targets map[discover.NodeType]int) map[discover.NodeType]int {
	cloned := make(map[discover.NodeType]int, len(targets))
	for nType, target := range targets {
		cloned[nType] = target
	}
	return cloned
}

// typedNodeSet is a set of discover.Nodes, which are counted by NodeType.
type typedNodeSet struct {
	nodes  map[discover.NodeID]*discover.Node
	counts map[discover.NodeType]int
}

func newTypedNodeSet() typedNodeSet {
	return typedNodeSet{
		nodes:  make(map[discover.NodeID]*discover.Node),
		counts: make(map[discover.NodeType]int),
	}
}

// Adding the same node ID with different NodeType is allowed, in which case the counts are adjusted.
// This can happen when the static node KNI given without/wrong node type, later connected with different node type.
func (t *typedNodeSet) add(n *discover.Node) {
	if n == nil {
		return
	}
	if old := t.nodes[n.ID]; old != nil {
		t.counts[old.NType]--
	}
	t.nodes[n.ID] = n
	t.counts[n.NType]++
}

func (t *typedNodeSet) get(id discover.NodeID) *discover.Node {
	return t.nodes[id]
}

func (t *typedNodeSet) remove(id discover.NodeID) {
	if n := t.nodes[id]; n != nil {
		delete(t.nodes, id)
		t.counts[n.NType]--
	}
}

func (t *typedNodeSet) len() int {
	return len(t.nodes)
}

func (t *typedNodeSet) all() []*discover.Node {
	nodes := make([]*discover.Node, 0, len(t.nodes))
	for _, n := range t.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (t *typedNodeSet) contains(id discover.NodeID) bool {
	_, ok := t.nodes[id]
	return ok
}

func (t *typedNodeSet) count(nType discover.NodeType) int {
	return t.counts[nType]
}
