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
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/mclock"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/nat"
)

// A peer dials all of its channels at once, so a set still incomplete after this is not coming.
const candidateAssemblyTimeout = 10 * time.Second

// MultiChannelServer is a server that uses a multi channel.
// It inherits BaseServer. This file only contains overriding methods.
// There must be a reason to override BaseServer methods.
type MultiChannelServer struct {
	*BaseServer
	listeners      []net.Listener                     // Extended TCP listener
	ListenAddrs    []string                           // Extended TCP listen addresses
	CandidateConns map[discover.NodeID][]*conn        // Subset of connections towards a premature peer
	candidateSince map[discover.NodeID]mclock.AbsTime // Open time of each CandidateConns entry
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
// Overrides BaseServer.Stop() to close multiple TCP listeners.
func (srv *MultiChannelServer) Stop() {
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
	for _, listener := range srv.listeners {
		listener.Close()
	}

	// Peers under assembly have no read loop to notice the shutdown.
	var candidates []*conn
	for id, connSet := range srv.CandidateConns {
		for _, c := range connSet {
			if c != nil {
				candidates = append(candidates, c)
			}
		}
		delete(srv.CandidateConns, id)
		delete(srv.candidateSince, id)
	}

	close(srv.quit) // Ask loops to terminate

	// Unlock here to allow Server loops and dialsched loop to finish their remaining jobs, dialsched to finish its dial goroutines.
	srv.lock.Unlock()

	for _, c := range candidates {
		c.close(DiscQuitting)
	}

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

// Start starts running the MultiChannelServer.
// MultiChannelServer can not be re-used after stopping.
func (srv *MultiChannelServer) Start() (err error) {
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
		if err := srv.startDiscovery(srv.ListenAddrs[ConnDefault]); err != nil {
			return err
		}
	}

	srv.ourHandshake = srv.makeOurHandshakeMsg()

	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}

	// Start listening
	if !srv.NoListen {
		if len(srv.ListenAddrs) != 0 && srv.ListenAddrs[ConnDefault] != "" {
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

	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", true)
	return nil
}

// Overrides BaseServer.initialize() to add multi channel specific fields.
func (srv *MultiChannelServer) initialize() error {
	if err := srv.BaseServer.initialize(); err != nil {
		return err
	}
	srv.listeners = make([]net.Listener, 0, len(srv.SubListenAddr)+1)
	srv.ListenAddrs = make([]string, 0, len(srv.SubListenAddr)+1)
	srv.ListenAddrs = append(srv.ListenAddrs, srv.ListenAddr)
	srv.ListenAddrs = append(srv.ListenAddrs, srv.SubListenAddr...)
	srv.CandidateConns = make(map[discover.NodeID][]*conn)
	srv.candidateSince = make(map[discover.NodeID]mclock.AbsTime)
	return nil
}

// Overrides BaseServer.makeOurHandshakeMsg() to advertise the multiple listen ports.
func (srv *MultiChannelServer) makeOurHandshakeMsg() *protoHandshake {
	hs := &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: true}
	for _, p := range srv.Protocols {
		hs.Caps = append(hs.Caps, p.cap())
	}
	for _, l := range srv.ListenAddrs {
		s := strings.Split(l, ":")
		if len(s) == 2 {
			if port, err := strconv.Atoi(s[1]); err == nil {
				hs.ListenPort = append(hs.ListenPort, uint64(port))
			}
		}
	}
	return hs
}

// startListening starts listening on the specified port on the server.
func (srv *MultiChannelServer) startListening() error {
	// Launch the TCP listener.
	for i, listenAddr := range srv.ListenAddrs {
		var lc net.ListenConfig
		listener, err := lc.Listen(context.Background(), "tcp", listenAddr)
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

// Overrides BaseServer.inboundAllowance(): a peer opens one connection per listener,
// plus the single-channel dial it made before the handshake told it that this node is
// multichannel.
func (srv *MultiChannelServer) inboundAllowance() int {
	return len(srv.ListenAddrs) + 1
}

// listenLoop waits for an external connection and connects it.
func (srv *MultiChannelServer) listenLoop(listener net.Listener) {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(listener))

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

		fd := srv.acceptInbound(listener, srv.inboundAllowance(), &lastLog)
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
func (srv *MultiChannelServer) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
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

// Overrides BaseServer.setupConn() to preserve multichannel port-order handling
// while moving the peer admission stages off the event loop.
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

// Override BaseServer.handleAddPeerConn because multichannel peer admission
// waits for the full candidate connection set and updates outbound metrics.
func (srv *MultiChannelServer) handleAddPeerConn(c *conn) error {
	// close writes a disconnect reason under a write deadline, so it must not run under srv.lock.
	// Registered before the unlock below to run after it.
	var (
		expired  []*conn
		replaced *conn
	)
	defer func() {
		for _, ec := range expired {
			ec.close(DiscUselessPeer)
		}
		if replaced != nil {
			replaced.close(DiscAlreadyConnected)
		}
	}()

	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return errServerStopped
	}
	expired = srv.expireCandidates(mclock.Now())
	err := srv.protoHandshakeChecks(srv.peers, c)
	if err != nil {
		return err
	}

	var p *Peer
	if c.multiChannel {
		// Out of range fills no slot; falling through would drop c unstored and unclosed.
		if int(c.portOrder) < 0 || int(c.portOrder) >= len(srv.ListenAddrs) {
			return DiscUnexpectedIdentity
		}

		connSet := srv.CandidateConns[c.id]
		if connSet == nil {
			connSet = make([]*conn, len(srv.ListenAddrs))
			srv.CandidateConns[c.id] = connSet
			srv.candidateSince[c.id] = mclock.Now()
		}

		replaced = connSet[c.portOrder]
		connSet[c.portOrder] = c

		pending := len(connSet)
		for _, conn := range connSet {
			if conn != nil {
				pending--
			}
		}
		if pending != 0 {
			return nil
		}

		p, err = newPeer(connSet, srv.Protocols, srv.Config.RWTimerConfig)
		delete(srv.CandidateConns, c.id)
		delete(srv.candidateSince, c.id)
	} else {
		p, err = newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
	}
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

// expireCandidates drops multichannel peers whose remaining channels never arrived and returns
// their connections for the caller to close. Nothing else does: CandidateConns is cleared only
// once every channel is present. Caller must hold srv.lock.
func (srv *MultiChannelServer) expireCandidates(now mclock.AbsTime) []*conn {
	var expired []*conn
	for id, connSet := range srv.CandidateConns {
		if now-srv.candidateSince[id] < mclock.AbsTime(candidateAssemblyTimeout) {
			continue
		}
		for _, c := range connSet {
			if c != nil {
				expired = append(expired, c)
			}
		}
		delete(srv.CandidateConns, id)
		delete(srv.candidateSince, id)
		srv.logger.Debug("Dropped incomplete multichannel peer", "id", id)
	}
	return expired
}

// Overrides BaseServer.runPeer() because multichannel peers run all socket
// read-writers and need multichannel-specific removal accounting.
func (srv *MultiChannelServer) runPeer(p *Peer) {
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
	remoteRequested, err := p.runWithRWs()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	srv.handleDelPeer(peerDrop{p, err, remoteRequested})
}

// Override BaseServer.handleDelPeer because multichannel teardown updates both
// inbound and outbound connection metrics.
func (srv *MultiChannelServer) handleDelPeer(pd peerDrop) {
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

func (srv *MultiChannelServer) reportPeerMetric() {
	peerCountGauge.Update(int64(len(srv.peers)))
	// In multi-channel, one peer may contain both inbound and outbound connections,
	// if two parties happened to dial each other at the same time. In this case,
	// it is unclear whether the peer should be counted as inbound or outbound.
	// Therefore, we do not report the peerIn and peerOut metrics for multi-channel.
	connectionCountGauge.Update(int64(srv.outboundCount + srv.inboundCount))
	connectionInCountGauge.Update(int64(srv.inboundCount))
	connectionOutCountGauge.Update(int64(srv.outboundCount))
}

// GetListenAddress returns the listen addresses of the server.
func (srv *MultiChannelServer) GetListenAddress() []string {
	return srv.ListenAddrs
}
