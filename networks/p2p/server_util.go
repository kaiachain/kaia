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
	"github.com/kaiachain/kaia/networks/p2p/discover"
)

// peerTargets from KIP-311.
var peerTargets = map[common.ConnType]map[common.ConnType]int{
	common.CONSENSUSNODE: {
		common.CONSENSUSNODE: 100,
		common.ENDPOINTNODE:  3,
	},
	common.ENDPOINTNODE: {
		common.CONSENSUSNODE: 2,
		common.ENDPOINTNODE:  10,
	},
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
