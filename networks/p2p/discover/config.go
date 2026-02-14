// Copyright 2024 The Kaia Authors
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
	"errors"
	"math"
)

var (
	errSelfBonding         = errors.New("cannot bond with self")
	errTableNotInitialized = errors.New("table not yet initialized")
)

// TODO: bring udp.go:Config
// TODO: bring contants from table.go

// Returns the number of nodes to actively discover depending of self node type.
// Note that this number is not the number of connections to maintain.
func getDiscoverTargets(cfg *Config) map[NodeType]int {
	switch cfg.NodeType {
	case NodeTypeCN:
		// Try to discover all CNs and BNs
		return map[NodeType]int{
			NodeTypeCN: 100,
			NodeTypePN: 0,
			NodeTypeEN: 0,
			NodeTypeBN: 3,
		}
	case NodeTypePN:
		// No active discovery. Only accepts inbound connections.
		return map[NodeType]int{
			NodeTypeCN: 0,
			NodeTypePN: 0,
			NodeTypeEN: 0,
			NodeTypeBN: 0,
		}
	case NodeTypeEN:
		// Try to discover all ENs and BNs.
		return map[NodeType]int{
			NodeTypeCN: 0,
			NodeTypePN: 0,
			NodeTypeEN: math.MaxInt32,
			NodeTypeBN: 3,
		}
	case NodeTypeBN:
		// Try to discover all other BNs.
		return map[NodeType]int{
			NodeTypeCN: 0,
			NodeTypePN: 0,
			NodeTypeEN: 0,
			NodeTypeBN: 3,
		}
	default:
		logger.Error("Unsupported node type", "NodeType", cfg.NodeType)
		return map[NodeType]int{}
	}
}
