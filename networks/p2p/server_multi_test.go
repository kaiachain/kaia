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

package p2p

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/mclock"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCandidateTestServer builds a two-channel server holding no peers yet.
func newCandidateTestServer(t *testing.T, maxPeers int) *MultiChannelServer {
	t.Helper()
	return &MultiChannelServer{
		BaseServer: &BaseServer{
			Config:  Config{ConnectionType: common.CONSENSUSNODE, MaxPhysicalConnections: maxPeers},
			running: true,
			peers:   make(map[discover.NodeID]*Peer),
			logger:  log.NewModuleLogger(log.NetworksP2P),
		},
		ListenAddrs:    []string{":32323", ":32324"},
		CandidateConns: make(map[discover.NodeID][]*conn),
		candidateSince: make(map[discover.NodeID]mclock.AbsTime),
		pendingDials:   make(map[*conn]struct{}),
	}
}

// newCandidateConn builds an inbound multichannel connection for node `n` on one channel.
func newCandidateConn(t *testing.T, n int, portOrder PortOrder) (*conn, *testTransport) {
	t.Helper()
	var id discover.NodeID
	id[0] = byte(n)
	return newCandidateConnFor(t, id, portOrder)
}

func newCandidateConnFor(t *testing.T, id discover.NodeID, portOrder PortOrder) (*conn, *testTransport) {
	t.Helper()
	serverFD, clientFD := net.Pipe()
	t.Cleanup(func() { serverFD.Close(); clientFD.Close() })

	tt := newTestTransport(id, serverFD, nil, true).(*testTransport)
	return &conn{
		fd:           serverFD,
		transport:    tt,
		flags:        inboundConn,
		conntype:     common.ENDPOINTNODE,
		id:           id,
		portOrder:    portOrder,
		multiChannel: true,
	}, tt
}

// lockProbeTransport records whether srv.lock was free while close ran. close writes a disconnect
// reason under a 1s deadline, so a sweep holding the lock across it stalls the whole server.
type lockProbeTransport struct {
	transport
	srv      *MultiChannelServer
	lockFree bool
}

func (t *lockProbeTransport) close(err error) {
	if t.srv.lock.TryLock() {
		t.lockFree = true
		t.srv.lock.Unlock()
	}
	t.transport.close(err)
}

func TestMultiChannelCandidatesExpire(t *testing.T) {
	srv := newCandidateTestServer(t, 4)

	stale, staleTransport := newCandidateConn(t, 1, ConnDefault)
	probe := &lockProbeTransport{transport: stale.transport, srv: srv}
	stale.transport = probe
	require.NoError(t, srv.handleAddPeerConn(stale))
	require.Len(t, srv.CandidateConns, 1)
	require.Empty(t, srv.peers, "half-assembled peers must not be registered")

	srv.candidateSince[stale.id] = mclock.Now() - mclock.AbsTime(candidateAssemblyTimeout) - 1

	fresh, _ := newCandidateConn(t, 2, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(fresh))

	assert.Equal(t, DiscUselessPeer, staleTransport.closeErr, "the expired candidate must be closed")
	assert.True(t, probe.lockFree, "the sweep must not close connections under srv.lock")
	require.Len(t, srv.CandidateConns, 1)
	assert.NotContains(t, srv.CandidateConns, stale.id)
	assert.Contains(t, srv.CandidateConns, fresh.id)
}

func TestMultiChannelDuplicateChannelClosesPrevious(t *testing.T) {
	srv := newCandidateTestServer(t, 4)

	first, firstTransport := newCandidateConn(t, 1, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(first))

	second, _ := newCandidateConn(t, 1, ConnDefault) // same node, same channel
	require.NoError(t, srv.handleAddPeerConn(second))

	assert.Equal(t, DiscAlreadyConnected, firstTransport.closeErr)
	require.Len(t, srv.CandidateConns, 1)
	assert.Same(t, second, srv.CandidateConns[second.id][ConnDefault])
}

func TestMultiChannelRejectsOutOfRangePortOrder(t *testing.T) {
	srv := newCandidateTestServer(t, 4)

	c, _ := newCandidateConn(t, 1, PortOrder(len(srv.ListenAddrs)))
	assert.ErrorIs(t, srv.handleAddPeerConn(c), DiscUnexpectedIdentity)
	assert.Empty(t, srv.CandidateConns)
}

// One multichannel peer opens a single-channel dial that ends in errUpdateDial, then
// one connection per listener, all from the same IP within microseconds. The per-IP
// inbound throttle has to fit them all, since handleAddPeerConn only builds a Peer
// once every channel has arrived.
func TestMultiChannelInboundThrottleFitsOnePeer(t *testing.T) {
	srv := newCandidateTestServer(t, 10)
	allowance := srv.inboundAllowance()
	require.Equal(t, 3, allowance, "a two-channel node must fit one dial per listener plus the single-channel dial")

	// TEST-NET-3 is public, so netutil.IsLAN does not exempt it the way loopback does.
	ip := net.ParseIP("203.0.113.7")
	now := mclock.AbsTime(0)

	for i := 0; i < allowance; i++ {
		require.NoError(t, srv.checkInboundConn(ip, allowance, now), "connection %d of one peer was rejected", i)
	}
	require.Error(t, srv.checkInboundConn(ip, allowance, now), "the throttle must still bound connections beyond one peer")
}

// failingProtoTransport fails the protocol handshake.
type failingProtoTransport struct {
	*testTransport
	err error
}

func (t *failingProtoTransport) doProtoHandshake(*protoHandshake) (*protoHandshake, error) {
	return nil, t.err
}

// newDialFixture prepares one dial of a fresh node. A non-nil failErr fails channel failAt.
func newDialFixture(t *testing.T, srv *MultiChannelServer, failAt PortOrder, failErr error) ([]net.Conn, *discover.Node, []*testTransport) {
	t.Helper()
	srv.PrivateKey = newkey()
	id := discover.PubkeyID(&newkey().PublicKey)

	fds := make([]net.Conn, len(srv.ListenAddrs))
	transports := make([]*testTransport, len(fds))
	pipes := make([]net.Conn, 0, 2*len(fds))
	for i := range fds {
		serverFD, clientFD := net.Pipe()
		fds[i] = serverFD
		tt := newTestTransport(id, serverFD, nil, true).(*testTransport)
		tt.listenPorts = []uint64{30303, 30304}
		transports[i] = tt
		pipes = append(pipes, serverFD, clientFD)
	}
	// Closing the pipes ends any peer running on them.
	t.Cleanup(func() {
		closeConns(pipes)
		srv.peerWG.Wait()
	})

	srv.newTransport = func(fd net.Conn, _ *ecdsa.PublicKey) transport {
		for i := range fds {
			if fds[i] == fd {
				if failErr != nil && PortOrder(i) == failAt {
					return &failingProtoTransport{testTransport: transports[i], err: failErr}
				}
				return transports[i]
			}
		}
		t.Errorf("unexpected connection")
		return nil
	}
	dest := &discover.Node{ID: id, IP: net.ParseIP("10.0.0.1"), TCPs: []uint16{30303, 30304}}
	return fds, dest, transports
}

func peerDirections(srv *MultiChannelServer, id discover.NodeID) (inbound, outbound int, ok bool) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	p := srv.peers[id]
	if p == nil {
		return 0, 0, false
	}
	inbound, outbound = p.GetNumberInboundAndOutbound()
	return inbound, outbound, true
}

func TestMultiChannelDialAddsOutboundPeer(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	fds, dest, _ := newDialFixture(t, srv, 0, nil)

	require.NoError(t, srv.SetupMultiConn(fds, dynDialedConn, dest))

	inbound, outbound, ok := peerDirections(srv, dest.ID)
	require.True(t, ok, "a dial whose channels all succeeded must add the peer")
	assert.Equal(t, 0, inbound)
	assert.Equal(t, len(fds), outbound)
	assert.Empty(t, srv.CandidateConns, "a dial must not use the candidate set")
}

func TestMultiChannelDialIsAllOrNothing(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	failErr := errors.New("handshake failed")
	fds, dest, transports := newDialFixture(t, srv, 1, failErr)

	require.ErrorIs(t, srv.SetupMultiConn(fds, dynDialedConn, dest), failErr)

	_, _, ok := peerDirections(srv, dest.ID)
	assert.False(t, ok)
	assert.Empty(t, srv.CandidateConns)
	for i, tt := range transports {
		assert.Equal(t, failErr, tt.closeErr, "channel %d must be closed", i)
	}
}

func TestMultiChannelDialNeverJoinsAcceptedChannels(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	fds, dest, _ := newDialFixture(t, srv, 0, nil)

	// A shared candidate set would pair this channel with dialed channel 0.
	accepted, acceptedTransport := newCandidateConnFor(t, dest.ID, PortOrder(1))
	probe := &lockProbeTransport{transport: accepted.transport, srv: srv}
	accepted.transport = probe
	require.NoError(t, srv.handleAddPeerConn(accepted))
	require.Contains(t, srv.CandidateConns, dest.ID)

	require.NoError(t, srv.SetupMultiConn(fds, dynDialedConn, dest))

	inbound, outbound, ok := peerDirections(srv, dest.ID)
	require.True(t, ok)
	assert.Equal(t, 0, inbound, "an accepted channel must not complete a dialed peer")
	assert.Equal(t, len(fds), outbound)
	assert.Equal(t, DiscAlreadyConnected, acceptedTransport.closeErr)
	assert.True(t, probe.lockFree, "the accepted channel must not be closed under srv.lock")
	assert.Empty(t, srv.CandidateConns)
	assert.Empty(t, srv.candidateSince)
	assert.Empty(t, srv.pendingDials)
}

func TestMultiChannelDialExpiresStaleCandidates(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	fds, dest, _ := newDialFixture(t, srv, 0, nil)

	stale, staleTransport := newCandidateConn(t, 1, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(stale))
	srv.candidateSince[stale.id] = mclock.Now() - mclock.AbsTime(candidateAssemblyTimeout) - 1

	require.NoError(t, srv.SetupMultiConn(fds, dynDialedConn, dest))

	assert.Equal(t, DiscUselessPeer, staleTransport.closeErr)
	assert.Empty(t, srv.CandidateConns)
}

func TestMultiChannelDialRejectsChannelCountMismatch(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	fds, dest, transports := newDialFixture(t, srv, 0, nil)
	srv.ListenAddrs = append(srv.ListenAddrs, ":32325")

	require.ErrorIs(t, srv.SetupMultiConn(fds, dynDialedConn, dest), errMultiChannelMismatch)

	assert.Empty(t, srv.peers)
	for i, tt := range transports {
		assert.Equal(t, errMultiChannelMismatch, tt.closeErr, "channel %d must be closed", i)
	}
}

// blockingProtoTransport blocks the protocol handshake until its connection closes.
type blockingProtoTransport struct {
	*testTransport
	entered chan struct{}
}

func (t *blockingProtoTransport) doProtoHandshake(*protoHandshake) (*protoHandshake, error) {
	close(t.entered)
	_, _, err := t.conn.Read()
	return nil, err
}

func TestMultiChannelStopClosesUnfinishedDial(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	srv.quit = make(chan struct{})
	fds, dest, transports := newDialFixture(t, srv, 0, nil)

	blocked := &blockingProtoTransport{testTransport: transports[0], entered: make(chan struct{})}
	newTransport := srv.newTransport
	srv.newTransport = func(fd net.Conn, key *ecdsa.PublicKey) transport {
		if fd == fds[0] {
			return blocked
		}
		return newTransport(fd, key)
	}

	done := make(chan error, 1)
	go func() { done <- srv.SetupMultiConn(fds, dynDialedConn, dest) }()
	<-blocked.entered

	srv.Stop()

	for i, fd := range fds {
		_, err := fd.Write([]byte{0})
		assert.ErrorIs(t, err, io.ErrClosedPipe, "channel %d must be closed", i)
	}
	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Stop must end a dial that is still in a handshake")
	}
	assert.Empty(t, srv.peers)
	assert.Empty(t, srv.pendingDials)
}

func TestMultiChannelAcceptedPeerRejectsDial(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	fds, dest, transports := newDialFixture(t, srv, 0, nil)

	for i := range srv.ListenAddrs {
		accepted, _ := newCandidateConnFor(t, dest.ID, PortOrder(i))
		require.NoError(t, srv.handleAddPeerConn(accepted))
	}

	require.ErrorIs(t, srv.SetupMultiConn(fds, dynDialedConn, dest), DiscAlreadyConnected)

	inbound, outbound, ok := peerDirections(srv, dest.ID)
	require.True(t, ok)
	assert.Equal(t, len(srv.ListenAddrs), inbound)
	assert.Equal(t, 0, outbound, "a dialed channel must not replace one of an accepted peer")
	assert.Equal(t, DiscAlreadyConnected, transports[0].closeErr)
	_, err := fds[1].Write([]byte{0})
	assert.ErrorIs(t, err, io.ErrClosedPipe, "the channel that never started its handshake must be closed too")
}

func TestMultiChannelDialedChannelStaysOutOfCandidates(t *testing.T) {
	srv := newCandidateTestServer(t, 4)

	accepted, _ := newCandidateConn(t, 1, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(accepted))

	dialed, _ := newCandidateConn(t, 1, 1)
	dialed.flags = dynDialedConn
	require.ErrorIs(t, srv.handleAddPeerConn(dialed), errMultiChannelMismatch)

	assert.Empty(t, srv.peers)
	require.Contains(t, srv.CandidateConns, accepted.id)
	assert.Same(t, accepted, srv.CandidateConns[accepted.id][ConnDefault])
	assert.Nil(t, srv.CandidateConns[accepted.id][1])
}

// capsTransport reports the given capabilities in the protocol handshake.
type capsTransport struct {
	*testTransport
	caps []Cap
}

func (t *capsTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	phs, err := t.testTransport.doProtoHandshake(our)
	phs.Caps = t.caps
	return phs, err
}

func TestMultiChannelDialChecksProtocolsPerChannel(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	srv.Protocols = []Protocol{{Name: "a", Version: 1}}
	fds, dest, transports := newDialFixture(t, srv, 0, nil)
	srv.newTransport = func(fd net.Conn, _ *ecdsa.PublicKey) transport {
		if fd == fds[0] {
			return &capsTransport{testTransport: transports[0], caps: []Cap{{Name: "a", Version: 1}}}
		}
		return transports[1]
	}

	require.ErrorIs(t, srv.SetupMultiConn(fds, dynDialedConn, dest), DiscUselessPeer)

	assert.Empty(t, srv.peers)
	for i, tt := range transports {
		assert.Equal(t, DiscUselessPeer, tt.closeErr, "channel %d must be closed", i)
	}
}

// The scheduler gets the peer's type and direction, whichever channel completed it.
func TestMultiChannelSchedulerGetsPeerClassification(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	t.Cleanup(srv.peerWG.Wait)
	ds := NewDialSched(DialConfig{maxPeers: 4}, &mockTable{}, srv)
	srv.dialSched = ds

	main, _ := newCandidateConn(t, 1, ConnDefault)
	sub, _ := newCandidateConn(t, 1, 1)
	sub.conntype = common.CONSENSUSNODE
	require.NoError(t, srv.handleAddPeerConn(main))
	require.NoError(t, srv.handleAddPeerConn(sub))

	require.Len(t, srv.peers, 1)
	assert.Equal(t, 1, ds.connectedAll.count(discover.NodeTypeEN))
	assert.Zero(t, ds.connectedAll.count(discover.NodeTypeCN))
}

// A dialed peer uses an outbound slot of the scheduler; an accepted one does not.
func TestMultiChannelDialDirectionReachesScheduler(t *testing.T) {
	for _, tc := range []struct {
		name    string
		inbound bool
	}{{"dial wins", false}, {"accepted set wins", true}} {
		t.Run(tc.name, func(t *testing.T) {
			srv := newCandidateTestServer(t, 4)
			fds, dest, _ := newDialFixture(t, srv, 0, nil)
			candidate := testNode(42, "10.0.0.42", discover.NodeTypeEN)
			table := &mockTable{nodesByType: map[discover.NodeType][]*discover.Node{discover.NodeTypeEN: {candidate}}}
			ds := NewDialSched(DialConfig{connTarget: map[discover.NodeType]int{discover.NodeTypeEN: 1}, maxPeers: 4}, table, srv)
			srv.dialSched = ds

			accepted, _ := newCandidateConnFor(t, dest.ID, PortOrder(1))
			require.NoError(t, srv.handleAddPeerConn(accepted))
			if tc.inbound {
				main, _ := newCandidateConnFor(t, dest.ID, ConnDefault)
				require.NoError(t, srv.handleAddPeerConn(main))
				require.ErrorIs(t, srv.SetupMultiConn(fds, dynDialedConn, dest), DiscAlreadyConnected)
			} else {
				require.NoError(t, srv.SetupMultiConn(fds, dynDialedConn, dest))
			}

			require.Len(t, srv.peers, 1)
			require.Equal(t, tc.inbound, srv.peers[dest.ID].Inbound())
			candidates, _ := ds.getCandidates()
			if tc.inbound {
				assert.Zero(t, ds.connectedOutbound.count(discover.NodeTypeEN))
				assert.Equal(t, []*discover.Node{candidate}, candidates)
			} else {
				assert.Equal(t, 1, ds.connectedOutbound.count(discover.NodeTypeEN))
				assert.Empty(t, candidates)
			}
		})
	}
}

func TestMultiChannelStopClosesCandidates(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	srv.quit = make(chan struct{})

	c, tt := newCandidateConn(t, 1, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(c))

	srv.Stop()

	assert.Equal(t, DiscQuitting, tt.closeErr)
	assert.Empty(t, srv.CandidateConns)
}

// Two nodes dial each other at the same moment over real TCP and RLPx. A link must never
// have channels of both directions, and a link that holds must work on every channel.
func TestMultiChannelMutualDial(t *testing.T) {
	const rounds, attemptsPerRound = 5, 50

	var exchanged sync.Map // *Peer -> struct{}, set once every channel passed a message both ways
	protocol := Protocol{
		Name: "mutualdial", Version: 1, Length: 1,
		RunWithRWs: func(peer *Peer, rws []MsgReadWriter) error {
			for channel, rw := range rws {
				if err := Send(rw, 0, uint64(channel)); err != nil {
					return err
				}
				if err := ExpectMsg(rw, 0, uint64(channel)); err != nil {
					return err
				}
			}
			exchanged.Store(peer, struct{}{})
			<-peer.closed
			return io.EOF
		},
	}
	newServer := func() (*MultiChannelServer, *discover.Node) {
		srv := &MultiChannelServer{BaseServer: &BaseServer{Config: Config{
			PrivateKey:             newkey(),
			ConnectionType:         common.CONSENSUSNODE,
			MaxPhysicalConnections: 10,
			NoDiscovery:            true,
			NoDial:                 true,
			ListenAddr:             "127.0.0.1:0",
			SubListenAddr:          []string{"127.0.0.1:0"},
			Protocols:              []Protocol{protocol},
			RWTimerConfig:          defaultRWTimerConfig,
		}}}
		t.Cleanup(srv.Stop)
		require.NoError(t, srv.Start())
		_, subport, err := net.SplitHostPort(srv.ListenAddrs[1])
		require.NoError(t, err)
		node, err := discover.ParseNode(fmt.Sprintf("kni://%s@%s?subport=%s&discport=0&ntype=cn",
			srv.selfID, srv.ListenAddrs[ConnDefault], subport))
		require.NoError(t, err)
		return srv, node
	}
	a, nodeA := newServer()
	b, nodeB := newServer()

	dial := func(from, to *MultiChannelServer, dest *discover.Node) {
		fds := make([]net.Conn, 0, len(to.ListenAddrs))
		for _, addr := range to.ListenAddrs {
			fd, err := (&net.Dialer{Timeout: time.Second}).DialContext(t.Context(), "tcp", addr)
			if err != nil {
				closeConns(fds)
				return
			}
			fds = append(fds, fd)
		}
		from.SetupMultiConn(fds, staticDialedConn, dest)
	}
	peerOf := func(srv *MultiChannelServer, id discover.NodeID) *Peer {
		srv.lock.Lock()
		defer srv.lock.Unlock()
		return srv.peers[id]
	}
	// settle reports the peers of a link that works on both ends, or nils once neither end holds a peer.
	settle := func() (*Peer, *Peer) {
		deadline, quiet := time.Now().Add(5*time.Second), 0
		for time.Now().Before(deadline) {
			pa, pb := peerOf(a, b.selfID), peerOf(b, a.selfID)
			for _, p := range []*Peer{pa, pb} {
				if p != nil {
					in, out := p.GetNumberInboundAndOutbound()
					require.True(t, in == 0 || out == 0, "a peer has %d inbound and %d outbound channels", in, out)
				}
			}
			if pa != nil && pb != nil {
				_, okA := exchanged.Load(pa)
				_, okB := exchanged.Load(pb)
				if okA && okB {
					return pa, pb
				}
			}
			if quiet = quiet + 1; pa != nil || pb != nil {
				quiet = 0
			}
			if quiet >= 5 {
				return nil, nil
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatal("the link neither held nor dropped")
		return nil, nil
	}

	firstTry := 0
	for round := range rounds {
		var pa, pb *Peer
		attempt := 0
		for ; pa == nil && attempt < attemptsPerRound; attempt++ {
			var wg sync.WaitGroup
			wg.Go(func() { dial(a, b, nodeB) })
			wg.Go(func() { dial(b, a, nodeA) })
			wg.Wait()
			pa, pb = settle()
		}
		require.NotNil(t, pa, "round %d: no link after %d simultaneous dials", round, attemptsPerRound)
		require.NotEqual(t, pa.Inbound(), pb.Inbound(), "round %d: the two ends disagree on who dialed", round)
		if attempt == 1 {
			firstTry++
		}

		pa.Disconnect(DiscRequested)
		require.Eventually(t, func() bool { return peerOf(a, b.selfID) == nil && peerOf(b, a.selfID) == nil },
			5*time.Second, 5*time.Millisecond)
	}
	t.Logf("%d of %d rounds linked on the first simultaneous dial", firstTry, rounds)
}

// A peer's advertised port list drives one dial each, so it is bounded before it is used.
func TestListenPortsToTCPs(t *testing.T) {
	const maxChannels = 2
	testcases := []struct {
		name  string
		ports []uint64
		want  []uint16
	}{
		{"single channel", []uint64{30303}, []uint16{30303}},
		{"all channels", []uint64{30303, 30304}, []uint16{30303, 30304}},
		{"more ports than listeners", []uint64{30303, 30304, 30305}, nil},
		{"port zero", []uint64{30303, 0}, nil},
		{"port above uint16", []uint64{30303, 65536 + 22}, nil},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tcps, err := listenPortsToTCPs(tc.ports, maxChannels)
			if tc.want == nil {
				assert.Equal(t, DiscProtocolError, err)
				assert.Nil(t, tcps)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, tcps)
		})
	}
}
