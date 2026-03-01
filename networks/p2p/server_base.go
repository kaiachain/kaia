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
	"context"
	"crypto/ecdsa"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/mclock"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/nat"
)

// BaseServer is a common data structure used by implementations of Server.
type BaseServer struct {
	// Config fields may not be modified while the server is running.
	Config

	// Hooks for testing. These are useful because we can inhibit
	// the whole protocol stack.
	newTransport func(net.Conn, *ecdsa.PublicKey) transport
	newPeerHook  func(*Peer)

	// Lifecycle
	lock    sync.Mutex
	running bool
	quit    chan struct{}
	loopWG  sync.WaitGroup // loop, listenLoop

	// Relying components
	logger       log.Logger
	ntab         discover.Discovery // UDP discovery
	ourHandshake *protoHandshake    // our TCP handshake message
	listener     net.Listener       // TCP listener
	peerFeed     event.Feed         // Peer event feed as in peer.go:PeerEventType.

	// LastLookup memory TODO: dialsched should manage it inside.
	lastLookup   time.Time
	lastLookupMu sync.Mutex

	// Peer lifecycle state
	selfID       discover.NodeID // precompulted Self().ID
	peers        map[discover.NodeID]*Peer
	inboundCount int
	trusted      map[discover.NodeID]bool
	dialstate    dialer

	// Request channels
	addstatic     chan *discover.Node
	removestatic  chan *discover.Node
	peerOp        chan peerOpFunc
	peerOpDone    chan struct{}
	posthandshake chan *conn
	addpeer       chan *conn
	delpeer       chan peerDrop
	discpeer      chan discover.NodeID
}

// SingleChannelServer is a server that uses a single channel.
type SingleChannelServer struct {
	*BaseServer
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
func (srv *BaseServer) Stop() {
	srv.lock.Lock()
	if !srv.running {
		srv.lock.Unlock()
		return
	}
	srv.running = false

	// Stop TCP listeners
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}

	close(srv.quit)   // Ask loops to terminate
	srv.lock.Unlock() // Unlock to allow loops to finish their remaining jobs.
	srv.loopWG.Wait() // Wait for loops to terminate
	srv.logger.Info("Stopped P2P server")
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

	if err := srv.initialize(); err != nil {
		return err
	}

	if !srv.NoDiscovery {
		if err := srv.startDiscovery(srv.ListenAddr); err != nil {
			return err
		}
	}

	srv.ourHandshake = srv.makeOurHandshakeMsg()

	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}

	// Start listening
	if !srv.NoListen {
		if srv.ListenAddr != "" {
			if err := srv.startListening(); err != nil {
				return err
			}
		} else {
			srv.logger.Error("P2P server might be useless, listening address is missing")
		}
	}

	// Start dialing. TODO: Only if !NoDial, start DialSched.
	srv.dialstate = newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())
	srv.loopWG.Add(1)
	go srv.run(srv.dialstate)

	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", false)
	return nil
}

func (srv *BaseServer) initialize() error {
	srv.logger = srv.Config.Logger
	if srv.logger == nil {
		srv.logger = logger.NewWith()
	}
	srv.logger.Info("Starting P2P networking")

	if srv.PrivateKey == nil {
		return errors.New("Server.PrivateKey must be set to a non-nil key")
	}
	if !srv.ConnectionType.Valid() {
		return errors.New("Invalid connection type speficied")
	}
	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *conn)
	srv.addstatic = make(chan *discover.Node)
	srv.removestatic = make(chan *discover.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})
	srv.discpeer = make(chan discover.NodeID)

	srv.selfID = discover.PubkeyID(&srv.PrivateKey.PublicKey)
	srv.peers = make(map[discover.NodeID]*Peer)
	srv.inboundCount = 0
	srv.trusted = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
	for _, n := range srv.TrustedNodes {
		srv.trusted[n.ID] = true
	}

	return nil
}

func (srv *BaseServer) startDiscovery(listenAddr string) error {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if srv.NAT != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "klaytn discovery")
		}
		// TODO: react to external IP changes over time.
		if ext, err := srv.NAT.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}
	cfg := discover.Config{
		PrivateKey:    srv.PrivateKey,
		AnnounceAddr:  realaddr,
		NodeDBPath:    srv.NodeDatabase,
		NetRestrict:   srv.NetRestrict,
		Bootnodes:     srv.BootstrapNodes,
		Unhandled:     nil,
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
	return nil
}

func (srv *BaseServer) makeOurHandshakeMsg() *protoHandshake {
	hs := &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: false}
	for _, p := range srv.Protocols {
		hs.Caps = append(hs.Caps, p.cap())
	}
	return hs
}

func (srv *BaseServer) startListening() error {
	// Launch the TCP listener.
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", srv.ListenAddr)
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

func (srv *BaseServer) run(dialstate dialer) {
	defer srv.loopWG.Done()
	var (
		peers        = srv.peers
		trusted      = srv.trusted
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
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
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
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.logger.Debug("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected
			srv.logger.Debug("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID]; ok {
				p.Disconnect(DiscRequested)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.logger.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.id] {
				// Ensure that the trusted flag is set before checking against MaxPhysicalConnections.
				c.flags |= trustedConn
			}
			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			case c.cont <- srv.encHandshakeChecks(peers, srv.inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			var err error
			err = srv.protoHandshakeChecks(peers, srv.inboundCount, c)
			if err == nil {
				// The handshakes are done and it passed all checks.
				p, err := newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
				if err != nil {
					srv.logger.Error("Fail make a new peer", "err", err)
				} else {
					// If message events are enabled, pass the peerFeed
					// to the peer
					if srv.EnableMsgEvents {
						p.events = &srv.peerFeed
					}
					name := truncateName(c.name)
					srv.logger.Debug("Adding p2p peer", "name", name, "addr", c.fd.RemoteAddr(), "peers", len(peers)+1)
					go srv.runPeer(p)
					peers[c.id] = p

					if p.Inbound() {
						srv.inboundCount++
					}
					peerCountGauge.Update(int64(len(peers)))
					peerInCountGauge.Update(int64(srv.inboundCount))
					peerOutCountGauge.Update(int64(len(peers) - srv.inboundCount))
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := common.PrettyDuration(mclock.Now() - pd.created)
			pd.logger.Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())

			if pd.Inbound() {
				srv.inboundCount--
			}

			peerCountGauge.Update(int64(len(peers)))
			peerInCountGauge.Update(int64(srv.inboundCount))
			peerOutCountGauge.Update(int64(len(peers) - srv.inboundCount))
		case nid := <-srv.discpeer:
			if p, ok := peers[nid]; ok {
				p.Disconnect(DiscRequested)
				p.logger.Debug("disconnect peer")
			}
		}
	}

	srv.logger.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}
	//if srv.DiscV5 != nil {
	//	srv.DiscV5.Close()
	//}
	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.logger.Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

func (srv *BaseServer) protoHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	return srv.encHandshakeChecks(peers, inboundCount, c)
}

func (srv *BaseServer) encHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	switch {
	case !c.is(trustedConn|staticDialedConn) && len(peers) >= srv.Config.MaxPhysicalConnections:
		return DiscTooManyPeers
	case !c.is(trustedConn) && c.is(inboundConn) && inboundCount >= srv.maxInboundConns():
		return DiscTooManyPeers
	case peers[c.id] != nil:
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

	c := &conn{fd: fd, flags: flags, conntype: common.ConnTypeUndefined, cont: make(chan error), portOrder: ConnDefault}

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
	err = srv.checkpoint(c, srv.posthandshake)
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

	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

// checkpoint sends the conn to run, which performs the
// post-handshake checks for the stage (posthandshake, addpeer).
func (srv *BaseServer) checkpoint(c *conn, stage chan<- *conn) error {
	select {
	case stage <- c:
	case <-srv.quit:
		return errServerStopped
	}
	select {
	case err := <-c.cont:
		return err
	case <-srv.quit:
		return errServerStopped
	}
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

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) Dial(dest *discover.Node) (net.Conn, error) {
	return srv.Dialer.Dial(dest)
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	return srv.Dialer.DialMulti(dest)
}

// Disconnect tries to disconnect peer.
func (srv *BaseServer) Disconnect(destID discover.NodeID) {
	srv.discpeer <- destID
}

// GetProtocols returns a slice of protocols.
func (srv *BaseServer) GetProtocols() []Protocol {
	return srv.Protocols
}

// AddProtocols adds protocols to the server.
func (srv *BaseServer) AddProtocols(p []Protocol) {
	srv.Protocols = append(srv.Protocols, p...)
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

//// Wrapper to discovery methods

// CheckNilNetworkTable returns whether network table is nil.
func (srv *BaseServer) CheckNilNetworkTable() bool {
	return srv.ntab == nil
}

func (srv *BaseServer) GetNodes(nType discover.NodeType, max int) []*discover.Node {
	return srv.ntab.GetNodes(nType, max)
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

//// Inject static nodes to connect to

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *BaseServer) AddPeer(node *discover.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node.
func (srv *BaseServer) RemovePeer(node *discover.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

//// Query connected peers

// SubscribeEvents subscribes the given channel to peer events.
func (srv *BaseServer) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
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

// PeerCount returns the number of connected peers.
func (srv *BaseServer) PeerCount() int {
	var count int
	select {
	case srv.peerOp <- func(ps map[discover.NodeID]*Peer) { count = len(ps) }:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return count
}

func (srv *BaseServer) PeerCountByType() map[string]uint {
	pc := make(map[string]uint)
	pc["total"] = 0
	select {
	case srv.peerOp <- func(ps map[discover.NodeID]*Peer) {
		for _, peer := range ps {
			key := ConvertConnTypeToString(peer.ConnType())
			pc[key]++
			pc["total"]++
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return pc
}

// Peers returns all connected peers.
func (srv *BaseServer) Peers() []*Peer {
	var ps []*Peer
	select {
	// Note: We'd love to put this function into a variable but
	// that seems to cause a weird compiler error in some
	// environments.
	case srv.peerOp <- func(peers map[discover.NodeID]*Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return ps
}

//// Describe the host (self)

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

// GetListenAddress returns the listen address of the server.
func (srv *BaseServer) GetListenAddress() []string {
	return []string{srv.ListenAddr}
}

// Name returns name of server.
func (srv *BaseServer) Name() string {
	return srv.Config.Name
}

// MaxPhysicalConnections returns maximum count of peers.
func (srv *BaseServer) MaxPeers() int {
	return srv.Config.MaxPhysicalConnections
}
