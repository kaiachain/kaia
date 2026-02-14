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

import "github.com/kaiachain/kaia/networks/p2p/discover"

type DialConfig struct {
	selfID   discover.NodeID
	selfType discover.NodeType

	bootnodes  []*discover.Node
	static     []*discover.Node // initial set of static nodes.
	authorized []*discover.Node // initial set of authorized nodes.

	// target number of outbound connections for each node type.
	// Set connTargets[T] = 0 to disable dynamic dial of node type T (i.e. only accept inbound connections).
	// Set connTargets = nil to use the default value based on the selfType.
	// See also: discover.table.discoverTargets.
	connTargets map[discover.NodeType]int
}

type DialSched struct {
	DialConfig

	static    map[discover.NodeID]struct{}
	connected map[discover.NodeID]struct{}
	dialing   map[discover.NodeID]struct{}

	staticDialQueue []*discover.Node
}

func NewDialSched(cfg DialConfig) *DialSched {
	ds := &DialSched{
		DialConfig: cfg,
	}
	if ds.connTargets == nil {
		ds.connTargets = getConnTargets(cfg.selfType)
	}
	return ds
}

func getConnTargets(selfType discover.NodeType) map[discover.NodeType]int {
	switch selfType {
	case discover.NodeTypeCN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 100,
			discover.NodeTypePN: 0,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	case discover.NodeTypePN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 1,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	case discover.NodeTypeEN:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 2,
			discover.NodeTypeEN: 3,
			discover.NodeTypeBN: 0,
		}
	default:
		return map[discover.NodeType]int{
			discover.NodeTypeCN: 0,
			discover.NodeTypePN: 0,
			discover.NodeTypeEN: 0,
			discover.NodeTypeBN: 0,
		}
	}
}
