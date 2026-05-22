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
	"github.com/kaiachain/kaia/networks/rpc"
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

	// APIs returns the RPC APIs for the discovery service.
	APIs() []rpc.API
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
	Conn     conn
	NodeType NodeType

	// These settings are required for discovery packet control
	MaxNeighborsNode uint
	AuthorizedNodes  []*Node

	// DiscoverNodetype is list of node type to enable discovery.
	DiscoverTypes DiscoverTypesConfig
}

// discovery2 combines a UDP transport and a Table2 into a Discovery2 implementation.
//
// Lifecycle: UDP is started before Table (udp.Start -> tab.Start) and stopped after
// (tab.Close -> udp.close), treating UDP as the transport layer and Table as the application.
//
// Table -> UDP calls (originating from the maintenance loops):
//   - revalidateLoop tick -> revalidateOnce()     -> udp.ping()
//   - refreshLoop         -> Bond() -> pingpong() -> udp.ping(), udp.waitping()
//   - refreshLoop         -> findNodesOnce()      -> udp.findnode()
//
// UDP -> Table calls (originating from udp.readLoop dispatching packet handlers):
//   - incoming PING     -> ping.preverify()     -> tab.IsAuthorized()  works even before tab.Start or after tab.Close
//   - incoming FINDNODE -> findnode.preverify() -> tab.IsBonded()      same as above
//   - incoming FINDNODE -> findnode.handle()    -> tab.ClosestNodes()  same as above
//   - incoming PING     -> go t.bond(...)       -> go tab.Bond()       no-op before tab.Start or after tab.Close
type discovery2 struct {
	tab *Table2
	udp *udp
}

func NewDiscovery2(cfg *Config) (Discovery2, error) {
	udp := newUDPv4(cfg)
	tab, err := newTable2(cfg, udp)
	if err != nil {
		return nil, err
	}
	udp.Start(tab)
	tab.Start()
	return &discovery2{tab: tab, udp: udp}, nil
}

func (d *discovery2) Close() {
	d.tab.Close()
	d.udp.close()
}

func (d *discovery2) Refresh() {
	d.tab.Refresh()
}

func (d *discovery2) RandomNodes(buf []*Node, nType NodeType) int {
	return d.tab.RandomNodes(buf, nType)
}

func (d *discovery2) PutAuthorizedNodes(nodes []*Node) {
	d.tab.PutAuthorizedNodes(nodes)
}
