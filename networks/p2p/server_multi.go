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

// MultiChannelServer is a server that uses a multi channel.
// It inherits BaseServer. This file only contains overriding methods.
// There must be a reason to override BaseServer methods.
type MultiChannelServer struct {
	*BaseServer
	listeners      []net.Listener              // Extended TCP listener
	ListenAddrs    []string                    // Extended TCP listen addresses
	CandidateConns map[discover.NodeID][]*conn // Subset of connections towards a premature peer
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

	close(srv.quit)   // Ask loops to terminate
	srv.lock.Unlock() // Unlock to allow loops to finish their remaining jobs.
	srv.loopWG.Wait() // Wait for loops to terminate
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

	// Start dialing. TODO: Only if !NoDial, start DialSched.
	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())
	srv.loopWG.Add(1)
	go srv.run(dialer)

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

	c := &conn{fd: fd, flags: flags, conntype: common.ConnTypeUndefined, cont: make(chan error), portOrder: PortOrderUndefined}

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

	if c.multiChannel && dialDest != nil && (dialDest.TCPs == nil || len(dialDest.TCPs) < 2) && len(dialDest.TCPs) < len(phs.ListenPort) {
		logger.Debug("[Dial] update and retry the dial candidate as a multichannel",
			"id", dialDest.ID, "addr", dialDest.IP, "previous", dialDest.TCPs, "new", phs.ListenPort)

		dialDest.TCPs = make([]uint16, 0, len(phs.ListenPort))
		for _, listenPort := range phs.ListenPort {
			dialDest.TCPs = append(dialDest.TCPs, uint16(listenPort))
		}
		return errUpdateDial
	}

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

// run is the main loop that the server runs.
func (srv *MultiChannelServer) run(dialstate dialer) {
	logger.Debug("[p2p.Server] start MultiChannel p2p server")
	defer srv.loopWG.Done()
	var (
		peers         = make(map[discover.NodeID]*Peer)
		inboundCount  = 0
		outboundCount = 0
		trusted       = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
		taskdone      = make(chan task, maxActiveDialTasks)
		runningTasks  []task
		queuedTasks   []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup and cannot be
	// modified while the server is running.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID] = true
	}

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
			case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			var p *Peer
			var e error
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			err := srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				if c.multiChannel {
					connSet := srv.CandidateConns[c.id]
					if connSet == nil {
						connSet = make([]*conn, len(srv.ListenAddrs))
						srv.CandidateConns[c.id] = connSet
					}

					if int(c.portOrder) < len(connSet) {
						connSet[c.portOrder] = c
					}

					count := len(connSet)
					for _, conn := range connSet {
						if conn != nil {
							count--
						}
					}

					if count == 0 {
						p, e = newPeer(connSet, srv.Protocols, srv.Config.RWTimerConfig)
						srv.CandidateConns[c.id] = nil
					}
				} else {
					// The handshakes are done and it passed all checks.
					p, e = newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
				}

				if e != nil {
					srv.logger.Error("Fail make a new peer", "err", e)
				} else if p != nil {
					// If message events are enabled, pass the peerFeed
					// to the peer
					if srv.EnableMsgEvents {
						p.events = &srv.peerFeed
					}
					name := truncateName(c.name)
					srv.logger.Debug("Adding p2p peer", "name", name, "addr", c.fd.RemoteAddr(), "peers", len(peers)+1)
					go srv.runPeer(p)
					peers[c.id] = p

					peerCountGauge.Update(int64(len(peers)))
					inboundCount, outboundCount = increasesConnectionMetric(inboundCount, outboundCount, p)
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

			peerCountGauge.Update(int64(len(peers)))
			inboundCount, outboundCount = decreasesConnectionMetric(inboundCount, outboundCount, pd.Peer)
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

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *MultiChannelServer) AddPeer(node *discover.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node.
func (srv *MultiChannelServer) RemovePeer(node *discover.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

// GetListenAddress returns the listen addresses of the server.
func (srv *MultiChannelServer) GetListenAddress() []string {
	return srv.ListenAddrs
}
