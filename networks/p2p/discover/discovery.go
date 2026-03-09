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
	"crypto/ecdsa"
	"net"

	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

type Discovery2 interface {
	// Lifecycle. Used by p2p.Server.
	Close()

	// Trigger a refresh of the discovery table and wait for completion.
	Refresh()

	// Serve discovered nodes. Used by p2p.DialSched.
	RandomNodes(buf []*Node, nType NodeType) int

	// Update authorized nodes. To be used by permless.
	PutAuthorizedNodes(nodes []*Node)
}

// Config holds Table-related settings.
type Config struct {
	NetworkID uint64
	// These settings are required and configure the UDP listener:
	PrivateKey *ecdsa.PrivateKey

	// These settings are optional:
	AnnounceAddr *net.UDPAddr     // local address announced in the DHT
	NodeDBPath   string           // if set, the node database is stored at this filesystem location
	NetRestrict  *netutil.Netlist // network whitelist
	Bootnodes    []*Node          // list of bootstrap nodes

	// These settings are required for create Table and UDP
	Id       NodeID
	Addr     *net.UDPAddr
	udp      transport
	Conn     conn
	NodeType NodeType

	// These settings are required for discovery packet control
	MaxNeighborsNode uint
	AuthorizedNodes  []*Node

	// DiscoverNodetype is list of node type to enable discovery.
	DiscoverTypes DiscoverTypesConfig
}

func NewDiscovery2(cfg *Config) (Discovery2, error) {
	return nil, nil
}
