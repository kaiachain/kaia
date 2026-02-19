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
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/mclock"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/nat"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

const (
	defaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	maxActiveDialTasks     = 16
	defaultMaxPendingPeers = 50
	defaultDialRatio       = 3

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second

	// Maximum number of times to retry typed static node discovery.
	typedStaticRetry = 3
)

var errServerStopped = errors.New("server stopped")

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

// NewServer returns a new Server interface.
func NewServer(config Config) Server {
	bServer := &BaseServer{
		Config: config,
	}

	if config.EnableMultiChannelServer {
		listeners := make([]net.Listener, 0, len(config.SubListenAddr)+1)
		listenAddrs := make([]string, 0, len(config.SubListenAddr)+1)
		listenAddrs = append(listenAddrs, config.ListenAddr)
		listenAddrs = append(listenAddrs, config.SubListenAddr...)
		return &MultiChannelServer{
			BaseServer:     bServer,
			listeners:      listeners,
			ListenAddrs:    listenAddrs,
			CandidateConns: make(map[discover.NodeID][]*conn),
		}
	} else {
		return &SingleChannelServer{
			BaseServer: bServer,
		}
	}
}

// Server manages all peer connections.
type Server interface {
	// GetProtocols returns a slice of protocols.
	GetProtocols() []Protocol

	// AddProtocols adds protocols to the server.
	AddProtocols(p []Protocol)

	// SetupConn runs the handshakes and attempts to add the connection
	// as a peer. It returns when the connection has been added as a peer
	// or the handshakes have failed.
	SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error

	// AddLastLookup adds lastLookup to duration.
	AddLastLookup() time.Time

	// SetLastLookupToNow sets LastLookup to the current time.
	SetLastLookupToNow()

	// CheckNilNetworkTable returns whether network table is nil.
	CheckNilNetworkTable() bool

	// GetNodes returns up to max alive nodes which a NodeType is nType
	GetNodes(nType discover.NodeType, max int) []*discover.Node

	// Lookup performs a network search for nodes close
	// to the given target. It approaches the target by querying
	// nodes that are closer to it on each iteration.
	// The given target does not need to be an actual node
	// identifier.
	Lookup(target discover.NodeID, nType discover.NodeType) []*discover.Node

	// Resolve searches for a specific node with the given ID and NodeType.
	// It returns nil if the node could not be found.
	Resolve(target discover.NodeID, nType discover.NodeType) *discover.Node

	// Start starts running the server.
	// Servers can not be re-used after stopping.
	Start() (err error)

	// Stop terminates the server and all active peer connections.
	// It blocks until all active connections are closed.
	Stop()

	// AddPeer connects to the given node and maintains the connection until the
	// server is shut down. If the connection fails for any reason, the server will
	// attempt to reconnect the peer.
	AddPeer(node *discover.Node)

	// RemovePeer disconnects from the given node.
	RemovePeer(node *discover.Node)

	// SubscribeEvents subscribes the given channel to peer events.
	SubscribeEvents(ch chan *PeerEvent) event.Subscription

	// PeersInfo returns an array of metadata objects describing connected peers.
	PeersInfo() []*PeerInfo

	// NodeInfo gathers and returns a collection of metadata known about the host.
	NodeInfo() *NodeInfo

	// Name returns name of server.
	Name() string

	// PeerCount returns the number of connected peers.
	PeerCount() int

	// PeerCountByType returns the number of connected specific tyeps of peers.
	PeerCountByType() map[string]uint

	// MaxPhysicalConnections returns maximum count of peers.
	MaxPeers() int

	// Disconnect tries to disconnect peer.
	Disconnect(destID discover.NodeID)

	// GetListenAddress returns the listen address list of the server.
	GetListenAddress() []string

	// Peers returns all connected peers.
	Peers() []*Peer

	// NodeDialer is used to connect to nodes in the network, typically by using
	// an underlying net.Dialer but also using net.Pipe in tests.
	NodeDialer
}

// MultiChannelServer is a server that uses a multi channel.
type MultiChannelServer struct {
	*BaseServer
	listeners      []net.Listener
	ListenAddrs    []string
	CandidateConns map[discover.NodeID][]*conn
}

// Start starts running the MultiChannelServer.
// MultiChannelServer can not be re-used after stopping.
func (srv *MultiChannelServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.logger = srv.Config.Logger
	if srv.logger == nil {
		srv.logger = logger.NewWith()
	}
	srv.logger.Info("Starting P2P networking")

	// static fields
	if srv.PrivateKey == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if !srv.ConnectionType.Valid() {
		return fmt.Errorf("Invalid connection type speficied")
	}

	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.schedWake = make(chan struct{}, 1)
	srv.initPeerState()

	var (
		conn      *net.UDPConn
		realaddr  *net.UDPAddr
		unhandled chan discover.ReadPacket
	)

	if !srv.NoDiscovery {
		addr, err := net.ResolveUDPAddr("udp", srv.ListenAddrs[ConnDefault])
		if err != nil {
			return err
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		realaddr = conn.LocalAddr().(*net.UDPAddr)
		if srv.NAT != nil {
			if !realaddr.IP.IsLoopback() {
				go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "klaytn discovery")
			}
			// TODO: react to external IP changes over time.
			if ext, err := srv.NAT.ExternalIP(); err == nil {
				realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
			}
		}
	}

	// node table
	if !srv.NoDiscovery {
		cfg := discover.Config{
			PrivateKey:    srv.PrivateKey,
			AnnounceAddr:  realaddr,
			NodeDBPath:    srv.NodeDatabase,
			NetRestrict:   srv.NetRestrict,
			Bootnodes:     srv.BootstrapNodes,
			Unhandled:     unhandled,
			Conn:          conn,
			Addr:          realaddr,
			Id:            discover.PubkeyID(&srv.PrivateKey.PublicKey),
			NodeType:      ConvertNodeType(srv.ConnectionType),
			NetworkID:     srv.NetworkID,
			DiscoverTypes: srv.DiscoverTypes,
		}

		ntab, err := discover.ListenUDP(&cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}

	srv.dialstate = newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())

	// handshake
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: true}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	for _, l := range srv.ListenAddrs {
		s := strings.Split(l, ":")
		if len(s) == 2 {
			if port, err := strconv.Atoi(s[1]); err == nil {
				srv.ourHandshake.ListenPort = append(srv.ourHandshake.ListenPort, uint64(port))
			}
		}
	}

	// listen/dial
	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}
	if !srv.NoListen {
		if srv.ListenAddrs != nil && len(srv.ListenAddrs) != 0 && srv.ListenAddrs[ConnDefault] != "" {
			if err := srv.startListening(); err != nil {
				return err
			}
		} else {
			srv.logger.Error("P2P server might be useless, listening address is missing")
		}
	}

	srv.loopWG.Add(1)
	go srv.run()
	srv.running = true
	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", true)
	return nil
}

// startListening starts listening on the specified port on the server.
func (srv *MultiChannelServer) startListening() error {
	// Launch the TCP listener.
	for i, listenAddr := range srv.ListenAddrs {
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}
		laddr := listener.Addr().(*net.TCPAddr)
		srv.ListenAddrs[i] = laddr.String()
		srv.listeners = append(srv.listeners, listener)
		srv.loopWG.Add(1)
		go srv.listenLoop(listener)
		// Map the TCP listening port if NAT is configured.
		if !laddr.IP.IsLoopback() && srv.NAT != nil {
			srv.loopWG.Go(func() {
				nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "klaytn p2p")
			})
		}
	}
	return nil
}

// listenLoop waits for an external connection and connects it.
func (srv *MultiChannelServer) listenLoop(listener net.Listener) {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(listener, srv.ntab))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.logger.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.logger.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.logger.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		fd = newMeteredConn(fd, true)
		srv.logger.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *MultiChannelServer) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}

	c := &conn{fd: fd, flags: flags, conntype: common.ConnTypeUndefined, portOrder: PortOrderUndefined}

	if dialDest != nil {
		// retrieve pubkey. if err occurs, dialPubkey is automatically set as nil
		dialPubkey, _ := dialDest.ID.Pubkey()
		c.transport = srv.newTransport(fd, dialPubkey)
		c.portOrder = PortOrder(dialDest.PortOrder)
	} else {
		c.transport = srv.newTransport(fd, nil)
		for i, addr := range srv.ListenAddrs {
			s1 := strings.Split(addr, ":")                    // string format example, [::]:30303 or 123.123.123.123:30303
			s2 := strings.Split(fd.LocalAddr().String(), ":") // string format example, 123.123.123.123:30303
			if s1[len(s1)-1] == s2[len(s2)-1] {
				c.portOrder = PortOrder(i)
				break
			}
		}
	}

	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.logger.Trace("close connection", "id", c.id, "err", err)
	}
	return err
}

// setupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *MultiChannelServer) setupConn(c *conn, flags connFlag, dialDest *discover.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	// Run the connection type handshake
	var err error
	if c.conntype, err = c.doConnTypeHandshake(srv.ConnectionType); err != nil {
		srv.logger.Warn("Failed doConnTypeHandshake", "addr", c.fd.RemoteAddr(), "conn", c.flags,
			"conntype", c.conntype, "err", err)
		return err
	}
	srv.logger.Trace("Connection Type Trace", "addr", c.fd.RemoteAddr(), "conn", c.flags, "ConnType", c.conntype.String())

	// Run the RLPx handshake.
	remotePubkey, err := c.doEncHandshake(srv.PrivateKey)
	if err != nil {
		srv.logger.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}

	c.id = discover.PubkeyID(remotePubkey)
	clog := srv.logger.NewWith("id", c.id, "addr", c.fd.RemoteAddr(), "conn", c.flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.id != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return DiscUnexpectedIdentity
	}
	select {
	case <-srv.quit:
		return errServerStopped
	default:
	}
	err = srv.handlePostHandshake(c)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed protobuf handshake", "err", err)
		return err
	}
	if phs.ID != c.id {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return DiscUnexpectedIdentity
	}
	c.caps, c.name, c.multiChannel = phs.Caps, phs.Name, phs.Multichannel

	if c.multiChannel && dialDest != nil && (dialDest.TCPs == nil || len(dialDest.TCPs) < 2) && len(dialDest.TCPs) < len(phs.ListenPort) {
		logger.Debug("[Dial] update and retry the dial candidate as a multichannel",
			"id", dialDest.ID, "addr", dialDest.IP, "previous", dialDest.TCPs, "new", phs.ListenPort)

		dialDest.TCPs = make([]uint16, 0, len(phs.ListenPort))
		for _, listenPort := range phs.ListenPort {
			dialDest.TCPs = append(dialDest.TCPs, uint16(listenPort))
		}
		return errUpdateDial
	}

	select {
	case <-srv.quit:
		return errServerStopped
	default:
	}
	err = srv.handleAddPeer(c)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

// run is the main loop that the server runs.
func (srv *MultiChannelServer) run() {
	logger.Debug("[p2p.Server] start MultiChannel p2p server")
	defer srv.loopWG.Done()
	var (
		taskdone     = make(chan task, maxActiveDialTasks)
		runningTasks []task
		queuedTasks  []task // tasks that can't run yet
	)

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.logger.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := srv.dialstate.newTasks(len(runningTasks)+len(queuedTasks), srv.allPeers(), time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.logger.Trace("Dial task done", "task", t)
			srv.dialstate.taskDone(t, time.Now())
			delTask(t)
		case <-srv.schedWake:
		}
	}

	srv.cleanupAfterRun()
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
func (srv *MultiChannelServer) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	for _, listener := range srv.listeners {
		listener.Close()
	}
	close(srv.quit)
	srv.loopWG.Wait()
}

// GetListenAddress returns the listen addresses of the server.
func (srv *MultiChannelServer) GetListenAddress() []string {
	return srv.ListenAddrs
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *MultiChannelServer) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.runWithRWs()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	srv.handlePeerDrop(peerDrop{p, err, remoteRequested})
}

// SingleChannelServer is a server that uses a single channel.
type SingleChannelServer struct {
	*BaseServer
}

// AddLastLookup adds lastLookup to duration.
func (srv *BaseServer) AddLastLookup() time.Time {
	srv.lastLookupMu.Lock()
	defer srv.lastLookupMu.Unlock()
	return srv.lastLookup.Add(lookupInterval)
}

// SetLastLookupToNow sets LastLookup to the current time.
func (srv *BaseServer) SetLastLookupToNow() {
	srv.lastLookupMu.Lock()
	defer srv.lastLookupMu.Unlock()
	srv.lastLookup = time.Now()
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) Dial(dest *discover.Node) (net.Conn, error) {
	return srv.Dialer.Dial(dest)
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	return srv.Dialer.DialMulti(dest)
}

// BaseServer is a common data structure used by implementation of Server.
type BaseServer struct {
	// Config fields may not be modified while the server is running.
	Config

	// Hooks for testing. These are useful because we can inhibit
	// the whole protocol stack.
	newTransport func(net.Conn, *ecdsa.PublicKey) transport
	newPeerHook  func(*Peer)

	lock    sync.Mutex // protects running
	running bool

	ntab         discover.Discovery
	listener     net.Listener
	ourHandshake *protoHandshake
	lastLookup   time.Time
	lastLookupMu sync.Mutex
	// DiscV5       *discv5.Network

	peerMu        sync.RWMutex
	peers         map[discover.NodeID]*Peer
	inboundCount  int
	outboundCount int
	trusted       map[discover.NodeID]bool
	// peerCond is broadcast whenever a peer is removed from peers.
	// cleanupAfterRun waits on it until the peer map is empty.
	// peerCond.L is &peerMu (write lock).
	peerCond *sync.Cond
	// stopping is set to true during shutdown to prevent new peers from being added.
	// Protected by peerMu.
	stopping bool

	quit      chan struct{}
	schedWake chan struct{}
	dialstate dialer
	loopWG    sync.WaitGroup // loop, listenLoop
	peerFeed  event.Feed
	logger    log.Logger
}

func (srv *BaseServer) initPeerState() {
	srv.peerMu.Lock()
	defer srv.peerMu.Unlock()

	srv.peers = make(map[discover.NodeID]*Peer)
	srv.inboundCount = 0
	srv.outboundCount = 0
	srv.stopping = false
	srv.trusted = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
	for _, n := range srv.TrustedNodes {
		srv.trusted[n.ID] = true
	}
	srv.peerCond = sync.NewCond(&srv.peerMu)
}

func (srv *BaseServer) allPeers() map[discover.NodeID]*Peer {
	srv.peerMu.RLock()
	defer srv.peerMu.RUnlock()

	peers := make(map[discover.NodeID]*Peer, len(srv.peers))
	for id, p := range srv.peers {
		peers[id] = p
	}
	return peers
}

func (srv *BaseServer) signalSchedule() {
	select {
	case srv.schedWake <- struct{}{}:
	default:
	}
}

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

type connFlag int

const (
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
	trustedConn
)

type PortOrder int

const (
	PortOrderUndefined PortOrder = -1
)

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

type transport interface {
	doConnTypeHandshake(myConnType common.ConnType) (common.ConnType, error)
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error)
	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)
	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because NsgPipe doesn't use it.
	close(err error)
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

func (c *conn) is(f connFlag) bool {
	return c.flags&f != 0
}

// GetProtocols returns a slice of protocols.
func (srv *BaseServer) GetProtocols() []Protocol {
	return srv.Protocols
}

// AddProtocols adds protocols to the server.
func (srv *BaseServer) AddProtocols(p []Protocol) {
	srv.Protocols = append(srv.Protocols, p...)
}

// Peers returns all connected peers.
func (srv *BaseServer) Peers() []*Peer {
	peers := srv.allPeers()
	ps := make([]*Peer, 0, len(peers))
	for _, p := range peers {
		ps = append(ps, p)
	}
	return ps
}

// PeerCount returns the number of connected peers.
func (srv *BaseServer) PeerCount() int {
	return len(srv.allPeers())
}

func (srv *BaseServer) PeerCountByType() map[string]uint {
	peers := srv.allPeers()
	pc := map[string]uint{
		"total": 0,
	}
	for _, peer := range peers {
		key := ConvertConnTypeToString(peer.ConnType())
		pc[key]++
		pc["total"]++
	}
	return pc
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *BaseServer) AddPeer(node *discover.Node) {
	if node == nil {
		return
	}
	select {
	case <-srv.quit:
		return
	default:
	}

	srv.dialstate.addStatic(node)
	srv.logger.Debug("Adding static node", "node", node)
	srv.signalSchedule()
}

// RemovePeer disconnects from the given node.
func (srv *BaseServer) RemovePeer(node *discover.Node) {
	if node == nil {
		return
	}
	select {
	case <-srv.quit:
		return
	default:
	}

	srv.dialstate.removeStatic(node)
	srv.logger.Debug("Removing static node", "node", node)
	srv.signalSchedule()

	srv.peerMu.RLock()
	p := srv.peers[node.ID]
	srv.peerMu.RUnlock()
	if p != nil {
		p.Disconnect(DiscRequested)
	}
}

// SubscribeEvents subscribes the given channel to peer events.
func (srv *BaseServer) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
}

// Self returns the local node's endpoint information.
func (srv *BaseServer) Self() *discover.Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return &discover.Node{IP: net.ParseIP("0.0.0.0")}
	}
	return srv.makeSelf(srv.listener, srv.ntab)
}

func (srv *BaseServer) makeSelf(listener net.Listener, discovery discover.Discovery) *discover.Node {
	// If the server's not running, return an empty node.
	// If the node is running but discovery is off, manually assemble the node infos.
	if discovery == nil {
		// Inbound connections disabled, use zero address.
		if listener == nil {
			return &discover.Node{IP: net.ParseIP("0.0.0.0"), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey)}
		}
		// Otherwise inject the listener address too
		addr := listener.Addr().(*net.TCPAddr)
		return &discover.Node{
			ID:  discover.PubkeyID(&srv.PrivateKey.PublicKey),
			IP:  addr.IP,
			TCP: uint16(addr.Port),
		}
	}
	// Otherwise return the discovery node.
	return discovery.Self()
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
func (srv *BaseServer) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	close(srv.quit)
	srv.loopWG.Wait()
}

// GetListenAddress returns the listen address of the server.
func (srv *BaseServer) GetListenAddress() []string {
	return []string{srv.ListenAddr}
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
		return 0, nil, fmt.Errorf("Connection was closed")
	}
	l := min(len(packet.Data), len(b))
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}

// Close implements discv5.conn
func (s *sharedUDPConn) Close() error {
	return nil
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *BaseServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.logger = srv.Config.Logger
	if srv.logger == nil {
		srv.logger = logger.NewWith()
	}
	srv.logger.Info("Starting P2P networking")

	// static fields
	if srv.PrivateKey == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if !srv.ConnectionType.Valid() {
		return fmt.Errorf("Invalid connection type speficied")
	}

	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.schedWake = make(chan struct{}, 1)
	srv.initPeerState()

	var (
		conn      *net.UDPConn
		realaddr  *net.UDPAddr
		unhandled chan discover.ReadPacket
	)

	if !srv.NoDiscovery {
		addr, err := net.ResolveUDPAddr("udp", srv.ListenAddr)
		if err != nil {
			return err
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		realaddr = conn.LocalAddr().(*net.UDPAddr)
		if srv.NAT != nil {
			if !realaddr.IP.IsLoopback() {
				go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "klaytn discovery")
			}
			// TODO: react to external IP changes over time.
			if ext, err := srv.NAT.ExternalIP(); err == nil {
				realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
			}
		}
	}

	// node table
	if !srv.NoDiscovery {
		cfg := discover.Config{
			PrivateKey:    srv.PrivateKey,
			AnnounceAddr:  realaddr,
			NodeDBPath:    srv.NodeDatabase,
			NetRestrict:   srv.NetRestrict,
			Bootnodes:     srv.BootstrapNodes,
			Unhandled:     unhandled,
			Conn:          conn,
			Addr:          realaddr,
			Id:            discover.PubkeyID(&srv.PrivateKey.PublicKey),
			NodeType:      ConvertNodeType(srv.ConnectionType),
			NetworkID:     srv.NetworkID,
			DiscoverTypes: srv.DiscoverTypes,
		}

		cfgForLog := cfg
		cfgForLog.PrivateKey = nil

		logger.Info("Create udp", "config", cfgForLog)

		ntab, err := discover.ListenUDP(&cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}

	srv.dialstate = newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())

	// handshake
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: false}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	// listen/dial
	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}
	if !srv.NoListen {
		if srv.ListenAddr != "" {
			if err := srv.startListening(); err != nil {
				return err
			}
		} else {
			srv.logger.Error("P2P server might be useless, listening address is missing")
		}
	}

	srv.loopWG.Add(1)
	go srv.run()
	srv.running = true
	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", false)
	return nil
}

func (srv *BaseServer) startListening() error {
	// Launch the TCP listener.
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.listener = listener
	srv.loopWG.Add(1)
	go srv.listenLoop()
	// Map the TCP listening port if NAT is configured.
	if !laddr.IP.IsLoopback() && srv.NAT != nil {
		srv.loopWG.Go(func() {
			nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "klaytn p2p")
		})
	}
	return nil
}

type dialer interface {
	newTasks(running int, peers map[discover.NodeID]*Peer, now time.Time) []task
	taskDone(task, time.Time)
	addStatic(*discover.Node)
	removeStatic(*discover.Node)
}

func (srv *BaseServer) handlePostHandshake(c *conn) error {
	srv.peerMu.RLock()
	defer srv.peerMu.RUnlock()

	if srv.trusted[c.id] {
		// Ensure that the trusted flag is set before checking against MaxPhysicalConnections.
		c.flags |= trustedConn
	}

	return srv.encHandshakeChecksUnlocked(c)
}

func (srv *BaseServer) handleAddPeer(c *conn) error {
	// At this point the connection is past the protocol handshake.
	// Its capabilities are known and the remote identity is verified.
	err := srv.protoHandshakeChecks(c)
	if err != nil {
		return err
	}

	p, err := newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
	if err != nil {
		srv.logger.Error("Fail make a new peer", "err", err)
		return nil
	}
	// If message events are enabled, pass the peerFeed to the peer.
	if srv.EnableMsgEvents {
		p.events = &srv.peerFeed
	}

	srv.peerMu.Lock()
	if srv.stopping {
		srv.peerMu.Unlock()
		return errServerStopped
	}
	// Repeat the capacity check because the peer set might have changed between the handshakes.
	if err := srv.encHandshakeChecksUnlocked(c); err != nil {
		srv.peerMu.Unlock()
		return err
	}
	srv.peers[c.id] = p
	srv.incrementConnectionCounts(p)
	peerCount := len(srv.peers)
	srv.peerMu.Unlock()

	go srv.runPeer(p)

	srv.logger.Debug("Adding p2p peer", "name", truncateName(c.name), "addr", c.fd.RemoteAddr(), "peers", peerCount)
	srv.reportPeerMetric()
	return nil
}

func (srv *MultiChannelServer) handleAddPeer(c *conn) error {
	// At this point the connection is past the protocol handshake.
	// Its capabilities are known and the remote identity is verified.
	if err := srv.protoHandshakeChecks(c); err != nil {
		return err
	}

	var (
		p *Peer
		e error
	)
	if c.multiChannel {
		readyConns, err := srv.collectMultiChannelConns(c)
		if err != nil {
			return err
		}
		if len(readyConns) == 0 {
			// The connections for the peer is not yet ready.
			// Revisit after another connection is established (and handleAddPeer() is called again).
			return nil
		}
		p, e = newPeer(readyConns, srv.Protocols, srv.Config.RWTimerConfig)
	} else {
		// The handshakes are done and it passed all checks.
		p, e = newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
	}

	if e != nil {
		srv.logger.Error("Fail make a new peer", "err", e)
		return nil
	}
	if p == nil {
		return nil
	}

	// If message events are enabled, pass the peerFeed to the peer.
	if srv.EnableMsgEvents {
		p.events = &srv.peerFeed
	}

	srv.peerMu.Lock()
	if srv.stopping {
		srv.peerMu.Unlock()
		return errServerStopped
	}
	// Repeat the capacity check because the peer set might have changed between the handshakes.
	if err := srv.encHandshakeChecksUnlocked(c); err != nil {
		srv.peerMu.Unlock()
		return err
	}
	srv.peers[c.id] = p
	srv.incrementConnectionCounts(p)
	peerCount := len(srv.peers)
	srv.peerMu.Unlock()

	go srv.runPeer(p)

	srv.logger.Debug("Adding p2p peer", "name", truncateName(c.name), "addr", c.fd.RemoteAddr(), "peers", peerCount)
	srv.reportPeerMetric()
	return nil
}

// collectMultiChannelConns collects multi-channel connections that is required to launch a peer.
// If the required number of connections are collected, return the list of connections for the peer.
// If the required number of connections are not collected, return nil.
// If the given new connection has a problem, return an error.
func (srv *MultiChannelServer) collectMultiChannelConns(c *conn) ([]*conn, error) {
	srv.peerMu.Lock()
	defer srv.peerMu.Unlock()

	connSet := srv.CandidateConns[c.id]
	if connSet == nil {
		connSet = make([]*conn, len(srv.ListenAddrs))
		srv.CandidateConns[c.id] = connSet
	}
	if c.portOrder < 0 || int(c.portOrder) >= len(connSet) {
		return nil, fmt.Errorf("invalid multi-channel port order: %d", c.portOrder)
	}
	connSet[c.portOrder] = c

	for _, conn := range connSet {
		if conn == nil {
			return nil, nil
		}
	}

	readyConns := make([]*conn, len(connSet))
	copy(readyConns, connSet)
	delete(srv.CandidateConns, c.id)
	return readyConns, nil
}

func (srv *BaseServer) handlePeerDrop(pd peerDrop) {
	// A peer disconnected.
	d := common.PrettyDuration(mclock.Now() - pd.created)

	srv.peerMu.Lock()
	delete(srv.peers, pd.ID())
	srv.decrementConnectionCounts(pd.Peer)
	peerCount := len(srv.peers)
	srv.peerMu.Unlock()

	// Wake up cleanupAfterRun if it is waiting for the peer map to drain.
	srv.peerCond.Broadcast()

	pd.logger.Debug("Removing p2p peer", "duration", d, "peers", peerCount, "req", pd.requested, "err", pd.err)
	srv.reportPeerMetric()
	srv.signalSchedule()
}

func (srv *BaseServer) handleDisconnectRequest(nid discover.NodeID) {
	srv.peerMu.RLock()
	p := srv.peers[nid]
	srv.peerMu.RUnlock()
	if p != nil {
		p.Disconnect(DiscRequested)
		p.logger.Debug("disconnect peer")
	}
}

func (srv *BaseServer) cleanupAfterRun() {
	srv.logger.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}

	// Set stopping=true under the lock so that handleAddPeer sees it atomically,
	// then disconnect all current peers and wait for them to finish.
	// peerCond.Wait() releases the write lock while sleeping, allowing
	// handlePeerDrop to acquire it and delete peers from the map.
	srv.peerMu.Lock()
	srv.stopping = true
	for _, p := range srv.peers {
		p.Disconnect(DiscQuitting)
	}
	for len(srv.peers) > 0 {
		logger.Warn("Waiting for peers to finish", "peers", len(srv.peers))
		srv.peerCond.Wait() // sleep until someone calls peerCond.Broadcast()
	}
	srv.peerMu.Unlock()
}

func (srv *BaseServer) run() {
	defer srv.loopWG.Done()
	var (
		taskdone     = make(chan task, maxActiveDialTasks)
		runningTasks []task
		queuedTasks  []task // tasks that can't run yet
	)

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.logger.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := srv.dialstate.newTasks(len(runningTasks)+len(queuedTasks), srv.allPeers(), time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.logger.Trace("Dial task done", "task", t)
			srv.dialstate.taskDone(t, time.Now())
			delTask(t)
		case <-srv.schedWake:
		}
	}

	srv.cleanupAfterRun()
}

func (srv *BaseServer) protoHandshakeChecks(c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	return nil
}

// This function is used to check the capacity of the peer set.
// Caller must hold the peerMu Lock or RLock.
func (srv *BaseServer) encHandshakeChecksUnlocked(c *conn) error {
	switch {
	case !c.is(trustedConn|staticDialedConn) && len(srv.peers) >= srv.Config.MaxPhysicalConnections:
		return DiscTooManyPeers
	case !c.is(trustedConn) && c.is(inboundConn) && srv.inboundCount >= srv.maxInboundConns():
		return DiscTooManyPeers
	case srv.peers[c.id] != nil:
		return DiscAlreadyConnected
	case c.id == srv.Self().ID:
		return DiscSelf
	default:
		return nil
	}
}

func (srv *BaseServer) maxInboundConns() int {
	return srv.Config.MaxPhysicalConnections - srv.maxDialedConns()
}

func (srv *BaseServer) maxDialedConns() int {
	switch srv.ConnectionType {
	case common.CONSENSUSNODE:
		return 0
	case common.PROXYNODE:
		return 0
	case common.ENDPOINTNODE:
		if srv.NoDiscovery || srv.NoDial {
			return 0
		}
		r := srv.DialRatio
		if r == 0 {
			r = defaultDialRatio
		}
		return srv.Config.MaxPhysicalConnections / r
	case common.BOOTNODE:
		return 0 // TODO check the bn for en
	default:
		logger.Crit("[p2p.Server] UnSupported Connection Type:", "ConnectionType", srv.ConnectionType)
		return 0
	}
}

func (srv *BaseServer) getTypeStatics() map[dialType]typedStatic {
	switch srv.ConnectionType {
	case common.CONSENSUSNODE:
		tsMap := make(map[dialType]typedStatic)
		if srv.DiscoverTypes.CN {
			tsMap[DT_CN] = typedStatic{discover.MaxCNCNCount, typedStaticRetry}
		}
		return tsMap
	case common.PROXYNODE:
		tsMap := make(map[dialType]typedStatic)
		if srv.DiscoverTypes.PN {
			tsMap[DT_PN] = typedStatic{discover.MaxPNPNCount, typedStaticRetry}
		}
		return tsMap
	case common.ENDPOINTNODE:
		tsMap := make(map[dialType]typedStatic)
		if srv.DiscoverTypes.PN {
			tsMap[DT_PN] = typedStatic{discover.MaxENPNCount, typedStaticRetry}
		}
		return tsMap
	case common.BOOTNODE:
		return nil
	default:
		logger.Crit("[p2p.Server] UnSupported Connection Type:", "ConnectionType", srv.ConnectionType)
		return nil
	}
}

// Caller must hold the peerMu Lock.
func (srv *BaseServer) incrementConnectionCounts(p *Peer) {
	peerInbound, peerOutbound := p.GetNumberInboundAndOutbound()
	srv.inboundCount += peerInbound
	srv.outboundCount += peerOutbound
}

// Caller must hold the peerMu Lock.
func (srv *BaseServer) decrementConnectionCounts(p *Peer) {
	peerInbound, peerOutbound := p.GetNumberInboundAndOutbound()
	srv.inboundCount -= peerInbound
	srv.outboundCount -= peerOutbound
}

func (srv *BaseServer) reportPeerMetric() {
	srv.peerMu.RLock()
	defer srv.peerMu.RUnlock()

	peerCountGauge.Update(int64(len(srv.peers)))
	peerInCountGauge.Update(int64(srv.inboundCount))
	peerOutCountGauge.Update(int64(srv.outboundCount))

	connectionCountGauge.Update(int64(srv.outboundCount + srv.inboundCount))
	connectionInCountGauge.Update(int64(srv.inboundCount))
	connectionOutCountGauge.Update(int64(srv.outboundCount))
}

func (srv *MultiChannelServer) reportPeerMetric() {
	srv.peerMu.RLock()
	defer srv.peerMu.RUnlock()

	peerCountGauge.Update(int64(len(srv.peers)))
	// In multi-channel, one peer may contain both inbound and outbound connections,
	// if two parties happened to dial each other at the same time. In this case,
	// it is unclear whether the peer should be counted as inbound or outbound.
	// Therefore, we do not report the peerIn and peerOut metrics for multi-channel.

	connectionCountGauge.Update(int64(srv.outboundCount + srv.inboundCount))
	connectionInCountGauge.Update(int64(srv.inboundCount))
	connectionOutCountGauge.Update(int64(srv.outboundCount))
}

type tempError interface {
	Temporary() bool
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *BaseServer) listenLoop() {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(srv.listener, srv.ntab))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.logger.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.logger.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.logger.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		fd = newMeteredConn(fd, true)
		srv.logger.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *BaseServer) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}

	c := &conn{fd: fd, flags: flags, conntype: common.ConnTypeUndefined, portOrder: ConnDefault}

	// if err occurs, dialPubkey is automatically set as nil
	if dialDest != nil {
		dialPubkey, _ := dialDest.ID.Pubkey()
		c.transport = srv.newTransport(fd, dialPubkey)
	} else {
		c.transport = srv.newTransport(fd, nil)
	}

	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.logger.Trace("Setting up connection failed", "id", c.id, "err", err)
	}
	return err
}

func (srv *BaseServer) setupConn(c *conn, flags connFlag, dialDest *discover.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	// Run the connection type handshake
	var err error
	if c.conntype, err = c.doConnTypeHandshake(srv.ConnectionType); err != nil {
		srv.logger.Warn("Failed doConnTypeHandshake", "addr", c.fd.RemoteAddr(), "conn", c.flags,
			"conntype", c.conntype, "err", err)
		return err
	}
	srv.logger.Trace("Connection Type Trace", "addr", c.fd.RemoteAddr(), "conn", c.flags, "ConnType", c.conntype.String())

	// Run the encryption handshake.
	remotePubkey, err := c.doEncHandshake(srv.PrivateKey)
	if err != nil {
		srv.logger.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}

	c.id = discover.PubkeyID(remotePubkey)
	clog := srv.logger.NewWith("id", c.id, "addr", c.fd.RemoteAddr(), "conn", c.flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.id != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return DiscUnexpectedIdentity
	}
	select {
	case <-srv.quit:
		return errServerStopped
	default:
	}
	err = srv.handlePostHandshake(c)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed protobuf handshake", "err", err)
		return err
	}
	if phs.ID != c.id {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return DiscUnexpectedIdentity
	}
	c.caps, c.name, c.multiChannel = phs.Caps, phs.Name, phs.Multichannel

	select {
	case <-srv.quit:
		return errServerStopped
	default:
	}
	err = srv.handleAddPeer(c)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *BaseServer) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.run()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	srv.handlePeerDrop(peerDrop{p, err, remoteRequested})
}

// NodeInfo represents a short summary of the information known about the host.
type NodeInfo struct {
	ID    string `json:"id"`   // Unique node identifier (also the encryption key)
	Name  string `json:"name"` // Name of the node, including client type, version, OS, custom data
	Enode string `json:"kni"`  // Enode URL for adding this peer from remote peers
	IP    string `json:"ip"`   // IP address of the node
	Ports struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (srv *BaseServer) NodeInfo() *NodeInfo {
	node := srv.Self()

	// Gather and assemble the generic node infos
	info := &NodeInfo{
		Name:       srv.Name(),
		Enode:      node.String(),
		ID:         node.ID.String(),
		IP:         node.IP.String(),
		ListenAddr: srv.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}
	info.Ports.Discovery = int(node.UDP)
	info.Ports.Listener = int(node.TCP)

	// Gather all the running protocol infos (only once per protocol type)
	for _, proto := range srv.Protocols {
		if _, ok := info.Protocols[proto.Name]; !ok {
			nodeInfo := interface{}("unknown")
			if query := proto.NodeInfo; query != nil {
				nodeInfo = proto.NodeInfo()
			}
			info.Protocols[proto.Name] = nodeInfo
		}
	}
	return info
}

// PeersInfo returns an array of metadata objects describing connected peers.
func (srv *BaseServer) PeersInfo() []*PeerInfo {
	// Gather all the generic and sub-protocol specific infos
	infos := make([]*PeerInfo, 0, srv.PeerCount())
	for _, peer := range srv.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	// Sort the result array alphabetically by node identifier
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}

// Disconnect tries to disconnect peer.
func (srv *BaseServer) Disconnect(destID discover.NodeID) {
	srv.handleDisconnectRequest(destID)
}

// CheckNilNetworkTable returns whether network table is nil.
func (srv *BaseServer) CheckNilNetworkTable() bool {
	return srv.ntab == nil
}

// Lookup performs a network search for nodes close
// to the given target. It approaches the target by querying
// nodes that are closer to it on each iteration.
// The given target does not need to be an actual node
// identifier.
func (srv *BaseServer) Lookup(target discover.NodeID, nType discover.NodeType) []*discover.Node {
	return srv.ntab.Lookup(target, nType)
}

// Resolve searches for a specific node with the given ID and NodeType.
// It returns nil if the node could not be found.
func (srv *BaseServer) Resolve(target discover.NodeID, nType discover.NodeType) *discover.Node {
	return srv.ntab.Resolve(target, nType)
}

func (srv *BaseServer) GetNodes(nType discover.NodeType, max int) []*discover.Node {
	return srv.ntab.GetNodes(nType, max)
}

// Name returns name of server.
func (srv *BaseServer) Name() string {
	return srv.Config.Name
}

// MaxPhysicalConnections returns maximum count of peers.
func (srv *BaseServer) MaxPeers() int {
	return srv.Config.MaxPhysicalConnections
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
