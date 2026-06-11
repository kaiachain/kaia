// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
//
// This file is derived from p2p/server.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package p2p

import (
	"errors"
	"net"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/p2p/discover"
)

// Default cross-type reservations (KIP-311). A node reserves this many slots for
// the *other* node type so its own-type mesh can't use every connection and leave
// no room for the CN<->EN link. Overridable per node via Config (0 = use default).
// e.g. a CN with maxconnections=110 meshes with up to 110-3=107 other CNs and
// always keeps 3 slots free for ENs.
const (
	defaultReservedENForCN = 3 // a CN keeps room for this many ENs by default
	defaultReservedCNForEN = 2 // an EN keeps room for this many CNs by default
)

// peerTargetFor returns how many peers of peerType a selfType node may keep, given
// maxconnections (maxPhys) and the node's cross-type reservations. Pass
// EffectiveConnType values. The reservation turns into accept caps:
//   - other type (CN->EN, EN->CN): the reserved number
//   - same type  (CN->CN, EN->EN): maxconnections minus the reservation
//
// ok=false means "no per-type cap" — only maxconnections / inbound limits apply.
func peerTargetFor(selfType, peerType common.ConnType, maxPhys, reservedENForCN, reservedCNForEN int) (target int, ok bool) {
	switch selfType {
	case common.CONSENSUSNODE:
		switch peerType {
		case common.CONSENSUSNODE:
			return max(0, maxPhys-reservedENForCN), true // CN mesh: room left after the EN reservation
		case common.ENDPOINTNODE:
			return reservedENForCN, true // CN->EN: the reserved slots
		}
	case common.ENDPOINTNODE:
		switch peerType {
		case common.ENDPOINTNODE:
			return max(0, maxPhys-reservedCNForEN), true // EN mesh: room left after the CN reservation
		case common.CONSENSUSNODE:
			return reservedCNForEN, true // EN->CN: the reserved slots
		}
	}
	return 0, false // unrelated pair (e.g. bootnode): no per-type cap
}

// minPhysicalConnections is the smallest maxconnections that still leaves at least
// 1 slot for the node's own mesh after the cross-type reservation. Below this the
// own-mesh cap would be 0 and the node couldn't mesh at all.
func minPhysicalConnections(selfType common.ConnType, reservedENForCN, reservedCNForEN int) int {
	switch EffectiveConnType(selfType) {
	case common.CONSENSUSNODE:
		return reservedENForCN + 1
	case common.ENDPOINTNODE:
		return reservedCNForEN + 1
	}
	return 1
}

// effectiveReservedENForCN and effectiveReservedCNForEN resolve a configured
// cross-type reservation, falling back to the package default when unset (<= 0).
func effectiveReservedENForCN(v int) int {
	if v > 0 {
		return v
	}
	return defaultReservedENForCN
}

func effectiveReservedCNForEN(v int) int {
	if v > 0 {
		return v
	}
	return defaultReservedCNForEN
}

// MinPhysicalConnections returns the smallest MaxPhysicalConnections a node of
// cfg's connection type may use while still honoring its cross-type reservation.
// Unset reservations fall back to the package defaults.
func MinPhysicalConnections(cfg Config) int {
	return minPhysicalConnections(cfg.ConnectionType, effectiveReservedENForCN(cfg.ReservedENForCN), effectiveReservedCNForEN(cfg.ReservedCNForEN))
}

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

type PortOrder int

const (
	PortOrderUndefined PortOrder = -1
)

type connFlag int

const (
	// inbound/dyndial/staticdial are mutually exclusive origins.
	// trusted is an independent overlay added after identity is known.
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
	trustedConn
)

func (f connFlag) String() string {
	s := ""
	if f&trustedConn != 0 {
		s += "-trusted"
	}
	if f&dynDialedConn != 0 {
		s += "-dyndial"
	}
	if f&staticDialedConn != 0 {
		s += "-staticdial"
	}
	if f&inboundConn != 0 {
		s += "-inbound"
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

// conn wraps a network connection with information gathered
// during the two handshakes.
type conn struct {
	fd net.Conn
	transport
	flags        connFlag
	conntype     common.ConnType // valid after the encryption handshake at the inbound connection case
	id           discover.NodeID // valid after the encryption handshake
	caps         []Cap           // valid after the protocol handshake
	name         string          // valid after the protocol handshake
	portOrder    PortOrder       // portOrder is the order of the ports that should be connected in multi-channel.
	multiChannel bool            // multiChannel is whether the peer is using multi-channel.
}

func (c *conn) String() string {
	s := c.flags.String()
	s += " " + c.conntype.String()
	if (c.id != discover.NodeID{}) {
		s += " " + c.id.String()
	}
	s += " " + c.fd.RemoteAddr().String()
	return s
}

func (c *conn) Inbound() bool {
	return c.flags&inboundConn != 0
}

func (c *conn) is(f connFlag) bool {
	return c.flags&f != 0
}

// sharedUDPConn implements a shared connection. Write sends messages to the underlying connection while read returns
// messages that were found unprocessable and sent to the unhandled channel by the primary listener.
type sharedUDPConn struct {
	*net.UDPConn
	unhandled chan discover.ReadPacket
}

// ReadFromUDP implements discv5.conn
func (s *sharedUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	packet, ok := <-s.unhandled
	if !ok {
		return 0, nil, errors.New("Connection was closed")
	}
	l := min(len(packet.Data), len(b))
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}

// Close implements discv5.conn
func (s *sharedUDPConn) Close() error {
	return nil
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
	if len(dest.TCPs) == 0 {
		return nil, nil
	}
	conns := make([]net.Conn, 0, len(dest.TCPs))
	for _, tcp := range dest.TCPs {
		addr := &net.TCPAddr{IP: dest.IP, Port: int(tcp)}
		conn, err := t.Dialer.Dial("tcp", addr.String())
		if err != nil {
			// Close any connections opened so far to avoid leaking sockets
			// when only some of the ports could be reached.
			for _, c := range conns {
				c.Close()
			}
			return nil, err
		}
		conns = append(conns, conn)
	}
	return conns, nil
}

type tempError interface {
	Temporary() bool
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

func ConvertNodeType(ct common.ConnType) discover.NodeType {
	switch ct {
	case common.CONSENSUSNODE:
		return discover.NodeTypeCN
	case common.PROXYNODE:
		return discover.NodeTypePN
	case common.ENDPOINTNODE:
		return discover.NodeTypeEN
	case common.BOOTNODE:
		return discover.NodeTypeBN
	default:
		return discover.NodeTypeUnknown // TODO-Kaia-Node Maybe, call panic() func or Crit()
	}
}

// EffectiveConnType returns the KIP-311 role bucket for a wire-level connection type.
func EffectiveConnType(ct common.ConnType) common.ConnType {
	if ct == common.PROXYNODE {
		return common.ENDPOINTNODE
	}
	return ct
}

func addressFromNodeID(id discover.NodeID) (common.Address, error) {
	pub, err := id.Pubkey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pub), nil
}

func ConvertConnType(nt discover.NodeType) common.ConnType {
	switch nt {
	case discover.NodeTypeCN:
		return common.CONSENSUSNODE
	case discover.NodeTypePN:
		return common.PROXYNODE
	case discover.NodeTypeEN:
		return common.ENDPOINTNODE
	case discover.NodeTypeBN:
		return common.BOOTNODE
	default:
		return common.UNKNOWNNODE
	}
}

func ConvertConnTypeToString(ct common.ConnType) string {
	switch ct {
	case common.CONSENSUSNODE:
		return "cn"
	case common.PROXYNODE:
		return "pn"
	case common.ENDPOINTNODE:
		return "en"
	case common.BOOTNODE:
		return "bn"
	default:
		return "unknown"
	}
}

func ConvertStringToConnType(s string) common.ConnType {
	st := strings.ToLower(s)
	switch st {
	case "cn":
		return common.CONSENSUSNODE
	case "pn":
		return common.PROXYNODE
	case "en":
		return common.ENDPOINTNODE
	case "bn":
		return common.BOOTNODE
	default:
		return common.UNKNOWNNODE
	}
}
