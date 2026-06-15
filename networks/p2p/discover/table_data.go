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

	"github.com/kaiachain/kaia/crypto"
)

func (tab *Table2) RandomNodes(buf []*Node, targetType NodeType) int {
	if s := tab.storages[targetType]; s == nil {
		return 0
	} else {
		return s.random(buf)
	}
}

func (tab *Table2) ClosestNodes(targetID NodeID, targetType NodeType, max int) []*Node {
	if s := tab.storages[targetType]; s == nil {
		return nil
	} else {
		return s.closest(crypto.Keccak256Hash(targetID[:]), max).entries
	}
}

func (tab *Table2) IsBonded(id NodeID, ip net.IP) bool {
	return tab.bondAge(id, ip) < nodeDBNodeExpiration
}

//// Wrappers for underlying storages.

// Note that Bootnodes are counted multiple times because they are added to all storages.
func (tab *Table2) len() int {
	n := 0
	for _, s := range tab.storages {
		n += s.len()
	}
	return n
}

func (tab *Table2) numReplacements() int {
	n := 0
	for _, s := range tab.storages {
		n += s.numReplacements()
	}
	return n
}

// For logging
func (tab *Table2) lenByNodeTypes() map[string]int {
	r := make(map[string]int)
	for ty, s := range tab.storages {
		r[nodeTypeName(ty)] = s.len()
	}
	return r
}

// Adds a node to one of the storages of the node's type.
// If the node is a bootstrap node, it is added to all storages.
func (tab *Table2) addNode(n *Node) {
	if n.NType == NodeTypeBN {
		for _, s := range tab.storages {
			s.add(n)
		}
		return
	} else {
		if s := tab.storages[n.NType]; s != nil {
			s.add(n)
		}
	}
}

func (tab *Table2) deleteNode(n *Node) {
	if n.NType == NodeTypeBN {
		for _, s := range tab.storages {
			s.delete(n)
		}
		return
	} else {
		if s := tab.storages[n.NType]; s != nil {
			s.delete(n)
		}
	}
}
