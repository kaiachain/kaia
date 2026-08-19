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
	"net"
	"testing"

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
	}
}

// newCandidateConn builds an inbound multichannel connection for node `n` on one channel.
func newCandidateConn(t *testing.T, n int, portOrder PortOrder) (*conn, *testTransport) {
	t.Helper()
	serverFD, clientFD := net.Pipe()
	t.Cleanup(func() { serverFD.Close(); clientFD.Close() })

	var id discover.NodeID
	id[0] = byte(n)
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

func TestMultiChannelStopClosesCandidates(t *testing.T) {
	srv := newCandidateTestServer(t, 4)
	srv.quit = make(chan struct{})

	c, tt := newCandidateConn(t, 1, ConnDefault)
	require.NoError(t, srv.handleAddPeerConn(c))

	srv.Stop()

	assert.Equal(t, DiscQuitting, tt.closeErr)
	assert.Empty(t, srv.CandidateConns)
}
