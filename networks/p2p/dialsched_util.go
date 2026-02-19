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
	"sync"

	"github.com/kaiachain/kaia/networks/p2p/discover"
)

func defaultConnTarget(selfType discover.NodeType) map[discover.NodeType]int {
	switch selfType {
	case discover.NodeTypeCN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 100,
		}
	case discover.NodeTypePN:
		return map[discover.NodeType]int{
			discover.NodeTypePN: 1,
		}
	case discover.NodeTypeEN:
		return map[discover.NodeType]int{
			discover.NodeTypePN: 2,
			discover.NodeTypeEN: 3,
		}
	default:
		return map[discover.NodeType]int{}
	}
}

type typedNodeSet struct {
	mu    sync.RWMutex
	nodes map[discover.NodeType]map[discover.NodeID]*discover.Node
}

func newTypedNodeSet() typedNodeSet {
	ns := make(map[discover.NodeType]map[discover.NodeID]*discover.Node)
	for _, nType := range []discover.NodeType{discover.NodeTypeCN, discover.NodeTypePN, discover.NodeTypeEN, discover.NodeTypeBN, discover.NodeTypeUnknown} {
		ns[nType] = make(map[discover.NodeID]*discover.Node)
	}
	return typedNodeSet{
		mu:    sync.RWMutex{},
		nodes: ns,
	}
}

func (t *typedNodeSet) add(n *discover.Node) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.nodes[n.NType][n.ID] = n
}

func (t *typedNodeSet) remove(id discover.NodeID) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for nType := range t.nodes {
		delete(t.nodes[nType], id)
	}
}

func (t *typedNodeSet) len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	count := 0
	for nType := range t.nodes {
		count += len(t.nodes[nType])
	}
	return count
}

func (t *typedNodeSet) all() []*discover.Node {
	t.mu.RLock()
	defer t.mu.RUnlock()
	nodes := make([]*discover.Node, 0, len(t.nodes))
	for nType := range t.nodes {
		for _, n := range t.nodes[nType] {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (t *typedNodeSet) contains(id discover.NodeID) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for nType := range t.nodes {
		if _, ok := t.nodes[nType][id]; ok {
			return true
		}
	}
	return false
}

func (t *typedNodeSet) count(nType discover.NodeType) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.nodes[nType])
}

// NodeDialer is used to connect to nodes in the network, typically by using
// an underlying net.Dialer but also using net.Pipe in tests.
type NodeDialer interface {
	Dial(*discover.Node) (net.Conn, error)
	DialMulti(*discover.Node) ([]net.Conn, error)
}

// TCPDialer implements the NodeDialer interface by using a net.Dialer to
// create TCP connections to nodes in the network.
type TCPDialer struct {
	*net.Dialer
}

// Dial creates a TCP connection to the node.
func (t TCPDialer) Dial(dest *discover.Node) (net.Conn, error) {
	addr := &net.TCPAddr{IP: dest.IP, Port: int(dest.TCP)}
	return t.Dialer.Dial("tcp", addr.String())
}

// DialMulti creates TCP connections to the node.
func (t TCPDialer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	var conns []net.Conn
	if dest.TCPs != nil || len(dest.TCPs) != 0 {
		conns = make([]net.Conn, 0, len(dest.TCPs))
		for _, tcp := range dest.TCPs {
			addr := &net.TCPAddr{IP: dest.IP, Port: int(tcp)}
			conn, err := t.Dialer.Dial("tcp", addr.String())
			conns = append(conns, conn)
			if err != nil {
				return nil, err
			}
		}
	}
	return conns, nil
}
