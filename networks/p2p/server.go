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
	"crypto/ecdsa"
	"errors"
	"net"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/nat"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

const (
	defaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	defaultMaxPendingPeers = 50
	defaultDialRatio       = 3

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second
)

var (
	errServerStopped = errors.New("server stopped")
	errUpdateDial    = errors.New("updated to be multichannel peer")
)

// Config holds Server options.
type Config struct {
	// This field must be set to a valid secp256k1 private key.
	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	// MaxPhysicalConnections is the maximum number of physical connections.
	// A peer uses one connection if single channel peer and uses two connections if
	// multi channel peer. It must be greater than zero.
	MaxPhysicalConnections int

	// ConnectionType is a type of connection like Consensus or Normal
	// described at common.ConnType
	// When the connection is established, each peer exchange each connection type
	ConnectionType common.ConnType

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DialRatio controls the ratio of inbound to dialed connections.
	// Example: a DialRatio of 2 allows 1/2 of connections to be dialed.
	// Setting DialRatio to zero defaults it to 3.
	DialRatio int `toml:",omitempty"`

	// NoDiscovery can be used to disable the peer discovery mechanism.
	// Disabling is useful for protocol debugging (manual topology).
	NoDiscovery bool

	// DiscoverTypes is list of node type to enable discovery.
	DiscoverTypes discover.DiscoverTypesConfig

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// BootstrapNodes are used to establish connectivity
	// with the rest of the network.
	BootstrapNodes []*discover.Node

	//// BootstrapNodesV5 are used to establish connectivity
	//// with the rest of the network using the V5 discovery
	//// protocol.
	//BootstrapNodesV5 []*discv5.Node `toml:",omitempty"`

	// Static nodes are used as pre-configured connections which are always
	// maintained and re-connected on disconnects.
	StaticNodes []*discover.Node

	// Trusted nodes are used as pre-configured connections which are always
	// allowed to connect, even above the peer limit.
	TrustedNodes []*discover.Node

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	Protocols []Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// NoListen can be used to disable the listening for incoming connections.
	NoListen bool

	// SubListenAddr is the list of the secondary listen address used for peer-to-peer connections.
	SubListenAddr []string

	// If EnableMultiChannelServer is true, multichannel can communicate with other nodes
	EnableMultiChannelServer bool

	// If set to a non-nil value, the given NAT port mapper
	// is used to make the listening port available to the
	// Internet.
	NAT nat.Interface `toml:",omitempty"`

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	Dialer NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`

	// RWTimerConfig is a configuration for interval based timer for rw.
	// It checks if a rw successfully writes its task in given time.
	RWTimerConfig RWTimerConfig

	// NetworkID to use for selecting peers to connect to
	NetworkID uint64
}

// Server manages all peer connections.
type Server interface {
	// Start starts running the server.
	// Servers can not be re-used after stopping.
	Start() (err error)

	// Stop terminates the server and all active peer connections.
	// It blocks until all active connections are closed.
	Stop()

	// SetupConn runs the handshakes and attempts to add the connection
	// as a peer. It returns when the connection has been added as a peer
	// or the handshakes have failed.
	SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error

	// Disconnect tries to disconnect peer.
	Disconnect(destID discover.NodeID)

	// Underlying protocol handlers
	GetProtocols() []Protocol
	AddProtocols(p []Protocol)

	// Inject static nodes to connect to
	// AddPeer connects to the given node and maintains the connection until the
	// server is shut down. If the connection fails for any reason, the server will
	// attempt to reconnect the peer.
	AddPeer(node *discover.Node)
	// RemovePeer disconnects from the given node.
	RemovePeer(node *discover.Node)

	// Query connected peers
	SubscribeEvents(ch chan *PeerEvent) event.Subscription
	PeersInfo() []*PeerInfo
	PeerCount() int
	PeerCountByType() map[string]uint
	Peers() []*Peer

	// Describe the host (self)
	NodeInfo() *NodeInfo
	Name() string
	GetListenAddress() []string
	MaxPeers() int // returns MaxPhysicalConnections
}

// NewServer returns a new Server interface.
func NewServer(config Config) Server {
	bServer := &BaseServer{
		Config: config,
	}

	if config.EnableMultiChannelServer {
		return &MultiChannelServer{
			BaseServer: bServer,
		}
	} else {
		return &SingleChannelServer{
			BaseServer: bServer,
		}
	}
}
