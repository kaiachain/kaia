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
	"maps"
	"net"
	"slices"
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

var (
	errNotWhitelisted         = errors.New("not whitelisted in NetRestrict")
	errTooManyInboundAttempts = errors.New("too many inbound connection attempts")
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
	peerWG  sync.WaitGroup // peer goroutines

	// Relying components
	logger       log.Logger
	ntab         discover.Discovery2 // UDP discovery
	ourHandshake *protoHandshake     // our TCP handshake message
	listener     net.Listener        // TCP listener
	peerFeed     event.Feed          // Peer event feed as in peer.go:PeerEventType.

	// Peer lifecycle state
	selfID        discover.NodeID // precompulted Self().ID
	peers         map[discover.NodeID]*Peer
	inboundCount  int
	outboundCount int
	trusted       map[discover.NodeID]bool
	dialSched     *DialSched
	cnPeerAddrs   map[common.Address]struct{} // KIP-311 peerNodes allowlist for CN-CN connections.

	// Inbound connection throttling. inboundHistory tracks the last seen time
	// per remote IP; its mutex guards concurrent access from multiple
	// listenLoop goroutines (one per listener in MultiChannelServer).
	inboundHistory   expHeap
	inboundHistoryMu sync.Mutex
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

	close(srv.quit) // Ask loops to terminate

	// Unlock here to allow Server loops and dialsched loop to finish their remaining jobs, dialsched to finish its dial goroutines.
	srv.lock.Unlock()

	if srv.dialSched != nil {
		srv.dialSched.Close() // Wait for dial attempts to finish.
	}
	if srv.ntab != nil {
		srv.ntab.Close() // Terminate discovery now that dialsched does not need it.
	}
	srv.disconnectAllPeers() // Disconnect all peers. Now that dialsched stopped, no more connections will be established.
	srv.loopWG.Wait()        // Wait for loops to terminate
	srv.peerWG.Wait()        // Wait for peers to terminate
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

	if !srv.NoDial {
		maxENToENDialTarget := srv.maxDialedConns()
		srv.dialSched = NewDialSched(DialConfig{
			selfID:              srv.selfID,
			selfType:            ConvertNodeType(srv.ConnectionType),
			staticNodes:         srv.StaticNodes,
			netrestrict:         srv.NetRestrict,
			maxENToENDialTarget: &maxENToENDialTarget,
			maxPeers:            srv.MaxPeers(),
			dialer:              srv.Dialer,
		}, srv.ntab, srv)
		srv.dialSched.Start()
	}

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
	srv.quit = make(chan struct{})

	srv.selfID = discover.PubkeyID(&srv.PrivateKey.PublicKey)
	srv.peers = make(map[discover.NodeID]*Peer)
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
		Conn:          conn,
		Addr:          realaddr,
		Id:            discover.PubkeyID(&srv.PrivateKey.PublicKey),
		NodeType:      ConvertNodeType(srv.ConnectionType),
		NetworkID:     srv.NetworkID,
		DiscoverTypes: srv.DiscoverTypes,
	}

	ntab, err := discover.NewDiscovery2(&cfg)
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

func (srv *BaseServer) protoHandshakeChecks(peers map[discover.NodeID]*Peer, c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	return srv.encHandshakeChecks(peers, c)
}

func (srv *BaseServer) encHandshakeChecks(peers map[discover.NodeID]*Peer, c *conn) error {
	switch {
	case !c.isTrustedOrStatic() && len(peers) >= srv.Config.MaxPhysicalConnections:
		return DiscTooManyPeers
	case !c.is(trustedConn) && c.is(inboundConn) && countInboundPeers(peers) >= srv.maxInboundConns():
		return DiscTooManyPeers
	case !srv.admissiblePeerType(c):
		return DiscUselessPeer
	case srv.exceedsPeerTarget(peers, c):
		return DiscTooManyPeers
	case peers[c.id] != nil:
		return DiscAlreadyConnected
	case c.id == srv.selfID:
		return DiscSelf
	default:
		return srv.admitByCNPeers(c)
	}
}

// reservedSlots resolves this node's cross-type reservation, using the per-type default when unset.
func (srv *BaseServer) reservedSlots() int {
	return effectiveReservedSlots(srv.ConnectionType, srv.Config.ReservedCrossTypeSlots)
}

func (srv *BaseServer) exceedsPeerTarget(peers map[discover.NodeID]*Peer, c *conn) bool {
	if c.isTrustedOrStatic() {
		return false
	}

	selfType := EffectiveConnType(srv.ConnectionType)
	peerType := EffectiveConnType(c.conntype)
	target, ok := peerTargetFor(selfType, peerType, srv.Config.MaxPhysicalConnections, srv.reservedSlots())
	if !ok {
		return false
	}

	peersByType := slices.DeleteFunc(slices.Collect(maps.Values(peers)), func(p *Peer) bool {
		return EffectiveConnType(p.ConnType()) != peerType
	})
	if len(peersByType) >= target {
		return true
	}
	// Inbound peers don't count toward dialTargets, so reserve one slot for a dialed peer.
	if c.Inbound() && selfType == common.CONSENSUSNODE && selfType != peerType {
		inbound := 0
		for _, p := range peersByType {
			if p.Inbound() {
				inbound++
			}
		}
		return inbound >= target-crossTypeDialReserve(selfType, peerType, target)
	}
	return false
}

// admissiblePeerType reports whether the server accepts a peer of this type.
// A CN/EN server meshes only with CN/EN peers (keeping others out of the
// uncapped path in exceedsPeerTarget); trusted/static and non-CN/EN servers
// accept any type.
func (srv *BaseServer) admissiblePeerType(c *conn) bool {
	if c.isTrustedOrStatic() {
		return true
	}
	selfType := EffectiveConnType(srv.ConnectionType)
	if selfType != common.CONSENSUSNODE && selfType != common.ENDPOINTNODE {
		return true
	}
	peerType := EffectiveConnType(c.conntype)
	return peerType == common.CONSENSUSNODE || peerType == common.ENDPOINTNODE
}

// countInboundPeers returns how many peers are connected via an inbound
// connection. A multichannel peer holds several sockets but counts once, so
// inbound capacity is measured in peers (per KIP-311), not physical sockets.
// Direction is taken from the peer's default channel (Peer.Inbound); a peer
// whose channels are of mixed direction (a rare simultaneous-dial race) is
// classified by that channel.
func countInboundPeers(peers map[discover.NodeID]*Peer) int {
	n := 0
	for _, p := range peers {
		if p.Inbound() {
			n++
		}
	}
	return n
}

func (srv *BaseServer) admitByCNPeers(c *conn) error {
	if EffectiveConnType(c.conntype) != common.CONSENSUSNODE {
		return nil
	}
	if c.isTrustedOrStatic() {
		return nil
	}
	if srv.cnPeerAddrs == nil {
		return nil
	}

	addr, err := addressFromNodeID(c.id)
	if err != nil {
		return err
	}
	if _, ok := srv.cnPeerAddrs[addr]; !ok {
		return DiscUselessPeer
	}
	return nil
}

func (srv *BaseServer) maxInboundConns() int {
	return srv.Config.MaxPhysicalConnections - srv.maxDialedConns()
}

func (srv *BaseServer) maxDialedConns() int {
	switch srv.ConnectionType {
	case common.CONSENSUSNODE:
		return 0
	case common.PROXYNODE, common.ENDPOINTNODE:
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

// checkInboundConn reports whether an inbound connection from remoteIP may proceed to
// the handshake at time now.
func (srv *BaseServer) checkInboundConn(remoteIP net.IP, allowance int, now mclock.AbsTime) error {
	if remoteIP == nil {
		// No usable remote IP (e.g. in-process test pipes); nothing to throttle.
		return nil
	}
	if srv.NetRestrict != nil && !srv.NetRestrict.Contains(remoteIP) {
		return errNotWhitelisted
	}

	srv.inboundHistoryMu.Lock()
	defer srv.inboundHistoryMu.Unlock()
	srv.inboundHistory.expire(now, nil)
	// LAN addresses are exempt so co-located/internal nodes are never throttled.
	if !netutil.IsLAN(remoteIP) && srv.inboundHistory.count(remoteIP.String()) >= allowance {
		return errTooManyInboundAttempts
	}
	srv.inboundHistory.add(remoteIP.String(), now+mclock.AbsTime(inboundThrottleTime))
	return nil
}

// acceptWithBackoff accepts the next connection from listener. On temporary
// errors (e.g. fd exhaustion) it backs off to avoid busy-looping and log
// flooding while resources recover. It returns nil when the listener is
// permanently closed and the caller should stop the loop.
func (srv *BaseServer) acceptWithBackoff(listener net.Listener, lastLog *time.Time) net.Conn {
	for {
		fd, err := listener.Accept()
		if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
			if time.Since(*lastLog) > 1*time.Second {
				srv.logger.Debug("Temporary read error", "err", err)
				*lastLog = time.Now()
			}
			time.Sleep(200 * time.Millisecond)
			continue
		} else if err != nil {
			srv.logger.Debug("Read error", "err", err)
			return nil
		}
		return fd
	}
}

// inboundAllowance is how many inbound connections one remote IP may open within
// inboundThrottleTime. A single-channel node advertises one port, so one peer opens
// one connection.
func (srv *BaseServer) inboundAllowance() int {
	return 1
}

// acceptInbound returns the next inbound connection worth handshaking. It wraps
// acceptWithBackoff and applies the inbound checks (NetRestrict + per-IP
// throttle), silently dropping rejected connections. It returns nil when the
// listener is permanently closed and the caller should stop the loop.
//
// The handshake slot is held across rejects: listenLoop is the sole receiver of
// its slots channel (SetupConn goroutines only send), so releasing and
// re-acquiring a token on each reject would be a no-op.
func (srv *BaseServer) acceptInbound(listener net.Listener, allowance int, lastLog *time.Time) net.Conn {
	for {
		fd := srv.acceptWithBackoff(listener, lastLog)
		if fd == nil {
			return nil
		}
		if err := srv.checkInboundConn(tcpAddrIP(fd.RemoteAddr()), allowance, mclock.Now()); err != nil {
			srv.logger.Debug("Rejected inbound connection", "addr", fd.RemoteAddr(), "err", err)
			fd.Close()
			continue
		}
		return fd
	}
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *BaseServer) listenLoop() {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(srv.listener))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	var lastLog time.Time
	for {
		// Wait for a handshake slot before accepting.
		<-slots

		fd := srv.acceptInbound(srv.listener, srv.inboundAllowance(), &lastLog)
		if fd == nil {
			return
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

	err = srv.handleAddPeerConn(c)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by handleAddPeerConn.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

func (srv *BaseServer) handlePostHandshake(c *conn) error {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return errServerStopped
	}
	if srv.trusted[c.id] {
		c.flags |= trustedConn
	}
	return srv.encHandshakeChecks(srv.peers, c)
}

func (srv *BaseServer) handleAddPeerConn(c *conn) error {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return errServerStopped
	}
	err := srv.protoHandshakeChecks(srv.peers, c)
	if err != nil {
		return err
	}

	p, err := newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
	if err != nil {
		srv.logger.Error("Fail make a new peer", "err", err)
		return err
	}

	if srv.EnableMsgEvents {
		p.events = &srv.peerFeed
	}
	srv.peers[c.id] = p
	srv.incrementConnectionCounts(p)
	srv.reportPeerMetric()
	srv.logger.Debug("Adding p2p peer", "name", truncateName(c.name), "addr", c.fd.RemoteAddr(), "peers", len(srv.peers))
	if srv.dialSched != nil {
		srv.dialSched.OnPeerConnected(c.id, ConvertNodeType(c.conntype), c.Inbound())
	}

	srv.peerWG.Add(1)
	go srv.runPeer(p)
	return nil
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *BaseServer) runPeer(p *Peer) {
	defer srv.peerWG.Done()

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

	srv.handleDelPeer(peerDrop{p, err, remoteRequested})
}

func (srv *BaseServer) handleDelPeer(pd peerDrop) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	delete(srv.peers, pd.ID())

	d := common.PrettyDuration(mclock.Now() - pd.created)
	pd.logger.Debug("Removing p2p peer", "duration", d, "peers", len(srv.peers), "req", pd.requested, "err", pd.err)
	srv.decrementConnectionCounts(pd.Peer)
	srv.reportPeerMetric()
	if srv.dialSched != nil {
		srv.dialSched.OnPeerDisconnected(pd.ID(), ConvertNodeType(pd.ConnType()))
	}
}

func (srv *BaseServer) incrementConnectionCounts(p *Peer) {
	peerInbound, peerOutbound := p.GetNumberInboundAndOutbound()
	srv.inboundCount += peerInbound
	srv.outboundCount += peerOutbound
}

func (srv *BaseServer) decrementConnectionCounts(p *Peer) {
	peerInbound, peerOutbound := p.GetNumberInboundAndOutbound()
	srv.inboundCount -= peerInbound
	srv.outboundCount -= peerOutbound
}

func (srv *BaseServer) reportPeerMetric() {
	peerCountGauge.Update(int64(len(srv.peers)))
	peerInCountGauge.Update(int64(srv.inboundCount))
	peerOutCountGauge.Update(int64(srv.outboundCount))
	connectionCountGauge.Update(int64(srv.outboundCount + srv.inboundCount))
	connectionInCountGauge.Update(int64(srv.inboundCount))
	connectionOutCountGauge.Update(int64(srv.outboundCount))
}

// SetCNPeers replaces the address-keyed CN admission allowlist.
// nil disables filtering; an empty list rejects all CN claims.
func (srv *BaseServer) SetCNPeers(addrs []common.Address) {
	// Every observer keeps only CN discovery entries it can authenticate.
	if srv.ntab != nil {
		srv.ntab.SetCNPeers(addrs)
	}
	// Only CNs enforce the validator allowlist. ENs/PNs serve everyone, so they
	// skip it: leaving cnPeerAddrs nil makes the admit and drop checks no-op.
	if srv.ConnectionType != common.CONSENSUSNODE {
		return
	}
	// Warn if the full CN mesh (len(addrs)-1 peers) no longer fits under the CN-mesh cap
	// (maxconnections - reservation). Only a CN in CNPeers forms the mesh; a nonCNPeers
	// CN (e.g. ValInactive) syncs via ENs, so it's exempt.
	if selfAddr, err := addressFromNodeID(srv.selfID); err == nil && slices.Contains(addrs, selfAddr) {
		reserved := srv.reservedSlots()
		if mesh := len(addrs) - 1; mesh > srv.Config.MaxPhysicalConnections-reserved {
			logger.Warn("CN count exceeds capacity for full mesh + EN reservation; raise maxconnections",
				"cnPeers", len(addrs), "maxConn", srv.Config.MaxPhysicalConnections, "reservedSlots", reserved)
		}
	}
	srv.lock.Lock()
	if addrs == nil {
		srv.cnPeerAddrs = nil
		if srv.dialSched != nil {
			srv.dialSched.SetCNPeers(nil)
		}
		srv.lock.Unlock()
		return
	}

	cnPeerAddrs := make(map[common.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		cnPeerAddrs[addr] = struct{}{}
	}
	srv.cnPeerAddrs = cnPeerAddrs
	if srv.dialSched != nil {
		srv.dialSched.SetCNPeers(addrs)
	}
	drops := srv.peersOutsideCNPeerAddrs()
	srv.lock.Unlock()

	for _, p := range drops {
		srv.logger.Info("Disconnecting CN peer outside allowlist", "id", p.ID())
		p.Disconnect(DiscRequested)
	}
}

func (srv *BaseServer) peersOutsideCNPeerAddrs() []*Peer {
	drops := make([]*Peer, 0)
	for id, p := range srv.peers {
		if p == nil || EffectiveConnType(p.ConnType()) != common.CONSENSUSNODE {
			continue
		}
		if p.rws[ConnDefault].isTrustedOrStatic() {
			continue
		}
		addr, err := addressFromNodeID(id)
		if err != nil {
			continue
		}
		if _, ok := srv.cnPeerAddrs[addr]; !ok {
			drops = append(drops, p)
		}
	}
	return drops
}

// Disconnect tries to disconnect peer.
func (srv *BaseServer) Disconnect(destID discover.NodeID) {
	srv.lock.Lock()
	if !srv.running {
		srv.lock.Unlock()
		return
	}
	p := srv.peers[destID]
	srv.lock.Unlock()

	if p != nil {
		p.Disconnect(DiscRequested)
		p.logger.Debug("disconnect peer")
	}
}

// GetProtocols returns a slice of protocols.
func (srv *BaseServer) GetProtocols() []Protocol {
	return srv.Protocols
}

// AddProtocols adds protocols to the server.
func (srv *BaseServer) AddProtocols(p []Protocol) {
	srv.Protocols = append(srv.Protocols, p...)
}

//// Inject static nodes to connect to

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *BaseServer) AddPeer(node *discover.Node) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}

	srv.logger.Debug("Adding static node", "node", node)
	if srv.dialSched != nil {
		srv.dialSched.AddStatic(node)
	}
}

// RemovePeer disconnects from the given node.
func (srv *BaseServer) RemovePeer(node *discover.Node) {
	srv.lock.Lock()
	if !srv.running {
		srv.lock.Unlock()
		return
	}

	srv.logger.Debug("Removing static node", "node", node)
	if srv.dialSched != nil {
		srv.dialSched.RemoveStatic(node.ID)
	}
	p := srv.peers[node.ID]
	srv.lock.Unlock()

	if p != nil {
		p.Disconnect(DiscRequested)
	}
}

func (srv *BaseServer) disconnectAllPeers() {
	srv.lock.Lock()
	ps := make([]*Peer, 0, len(srv.peers))
	for _, p := range srv.peers {
		ps = append(ps, p)
	}
	srv.lock.Unlock()

	for _, p := range ps {
		p.Disconnect(DiscQuitting)
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
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return 0
	}
	return len(srv.peers)
}

func (srv *BaseServer) PeerCountByType() map[string]uint {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	pc := map[string]uint{"total": 0}
	if !srv.running {
		return pc
	}
	for _, peer := range srv.peers {
		key := ConvertConnTypeToString(peer.ConnType())
		pc[key]++
		pc["total"]++
	}
	return pc
}

// Peers returns all connected peers.
func (srv *BaseServer) Peers() []*Peer {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	ps := make([]*Peer, 0, len(srv.peers))
	if !srv.running {
		return ps
	}
	for _, p := range srv.peers {
		ps = append(ps, p)
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
	return srv.makeSelf(srv.listener)
}

// makeSelf returns the local node's TCP endpoint information.
// The `listener` shall be the single listener (srv.listener) for BaseServer,
// and one of the listeners (srv.listeners[*]) for MultiChannelServer.
func (srv *BaseServer) makeSelf(listener net.Listener) *discover.Node {
	id := discover.PubkeyID(&srv.PrivateKey.PublicKey)
	if listener == nil {
		return &discover.Node{IP: net.ParseIP("0.0.0.0"), ID: id}
	} else {
		addr := listener.Addr().(*net.TCPAddr)
		return discover.NewNode(id, addr.IP, uint16(addr.Port), uint16(addr.Port), nil, ConvertNodeType(srv.ConnectionType))
	}
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
