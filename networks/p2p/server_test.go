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
// This file is derived from p2p/server_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package p2p

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/rand"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/rlpx"
)

func init() {
	// log.Root().SetHandler(logger.LvlFilterHandler(logger.LvlError, logger.StreamHandler(os.Stderr, logger.TerminalFormat(false))))
}

type testTransport struct {
	id discover.NodeID
	*rlpxTransport
	mutichannel bool

	closeErr error
}

func newTestTransport(id discover.NodeID, fd net.Conn, dialDest *ecdsa.PublicKey, mutichannel bool) transport {
	wrapped := newRLPX(fd, dialDest).(*rlpxTransport)
	wrapped.conn.InitWithSecrets(rlpx.Secrets{
		MAC:        make([]byte, 16),
		AES:        make([]byte, 16),
		IngressMAC: sha3.NewKeccak256(),
		EgressMAC:  sha3.NewKeccak256(),
	})
	return &testTransport{id: id, rlpxTransport: wrapped, mutichannel: mutichannel}
}

func (c *testTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	remoteKey, _ := c.id.Pubkey()
	return remoteKey, nil
}

func (c *testTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	return &protoHandshake{ID: c.id, Name: "test", Multichannel: c.mutichannel}, nil
}

func (c *testTransport) doConnTypeHandshake(myConnType common.ConnType) (common.ConnType, error) {
	return 1, nil
}

func (c *testTransport) close(err error) {
	c.conn.Close()
	c.closeErr = err
}

func startTestServer(t *testing.T, id discover.NodeID, pf func(*Peer), config *Config) Server {
	config.Name = "test"
	config.MaxPhysicalConnections = 10
	config.ListenAddr = "127.0.0.1:0"
	config.PrivateKey = newkey()
	server := &SingleChannelServer{
		&BaseServer{
			Config:      *config,
			newPeerHook: pf,
			newTransport: func(fd net.Conn, dialDest *ecdsa.PublicKey) transport {
				return newTestTransport(id, fd, dialDest, false)
			},
		},
	}
	if err := server.Start(); err != nil {
		t.Fatalf("Could not start server: %v", err)
	}
	return server
}

func startTestMultiChannelServer(t *testing.T, id discover.NodeID, pf func(*Peer), config *Config) Server {
	config.Name = "test"
	config.MaxPhysicalConnections = 10
	config.PrivateKey = newkey()

	listeners := make([]net.Listener, 0, len(config.SubListenAddr)+1)
	listenAddrs := make([]string, 0, len(config.SubListenAddr)+1)
	listenAddrs = append(listenAddrs, config.ListenAddr)
	listenAddrs = append(listenAddrs, config.SubListenAddr...)

	server := &MultiChannelServer{
		BaseServer: &BaseServer{
			Config:      *config,
			newPeerHook: pf,
			newTransport: func(fd net.Conn, dialDest *ecdsa.PublicKey) transport {
				return newTestTransport(id, fd, dialDest, true)
			},
		},
		listeners:      listeners,
		ListenAddrs:    listenAddrs,
		CandidateConns: make(map[discover.NodeID][]*conn),
	}
	if err := server.Start(); err != nil {
		t.Fatalf("Could not start server: %v", err)
	}
	return server
}

func makeconn(fd net.Conn, id discover.NodeID) *conn {
	dialDest, _ := id.Pubkey()
	tx := newTestTransport(id, fd, dialDest, false)
	return &conn{fd: fd, transport: tx, flags: staticDialedConn, conntype: common.ConnTypeUndefined, id: id}
}

func makeMultiChannelConn(fd net.Conn, id discover.NodeID) *conn {
	dialDest, _ := id.Pubkey()
	tx := newTestTransport(id, fd, dialDest, true)
	return &conn{fd: fd, transport: tx, flags: staticDialedConn, conntype: common.ConnTypeUndefined, id: id, multiChannel: true}
}

func TestServerListen(t *testing.T) {
	// start the test server
	connected := make(chan *Peer)
	remid := discover.PubkeyID(&newkey().PublicKey)
	srv := startTestServer(t, remid, func(p *Peer) {
		if p.ID() != remid {
			t.Error("peer func called with wrong node id")
		}
		if p == nil {
			t.Error("peer func called with nil conn")
		}
		connected <- p
	}, &Config{})
	defer close(connected)
	defer srv.Stop()

	// dial the test server
	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", srv.GetListenAddress()[ConnDefault])
	if err != nil {
		t.Fatalf("could not dial: %v", err)
	}
	c := makeconn(conn, randomID())
	c.doConnTypeHandshake(c.conntype)

	defer conn.Close()

	select {
	case peer := <-connected:
		if peer.LocalAddr().String() != conn.RemoteAddr().String() {
			t.Errorf("peer started with wrong conn: got %v, want %v",
				peer.LocalAddr(), conn.RemoteAddr())
		}

		peers := srv.Peers()
		if !reflect.DeepEqual(peers, []*Peer{peer}) {
			t.Errorf("Peers mismatch: got %v, want %v", peers, []*Peer{peer})
		}
	case <-time.After(5 * time.Second):
		t.Error("server did not accept within one second")
	}
}

func TestMultiChannelServerListen(t *testing.T) {
	// start the test server
	connected := make(chan *Peer)
	remid := discover.PubkeyID(&newkey().PublicKey)
	config := &Config{ListenAddr: "127.0.0.1:33331", SubListenAddr: []string{"127.0.0.1:33333"}}
	srv := startTestMultiChannelServer(t, remid, func(p *Peer) {
		if p.ID() != remid {
			t.Error("peer func called with wrong node id")
		}
		if p == nil {
			t.Error("peer func called with nil conn")
		}
		connected <- p
	}, config)
	defer close(connected)
	defer srv.Stop()

	// dial the test server
	var defaultConn net.Conn

	for i, address := range srv.GetListenAddress() {
		dialer := net.Dialer{Timeout: 5 * time.Second}
		conn, err := dialer.Dial("tcp", address)
		defer conn.Close()

		if i == ConnDefault {
			defaultConn = conn
		}

		if err != nil {
			t.Fatalf("could not dial: %v", err)
		}
	}

	select {
	case peer := <-connected:
		if peer.LocalAddr().String() != defaultConn.RemoteAddr().String() {
			t.Errorf("peer started with wrong conn: got %v, want %v",
				peer.LocalAddr(), defaultConn.RemoteAddr())
		}

		peers := srv.Peers()
		if !reflect.DeepEqual(peers, []*Peer{peer}) {
			t.Errorf("Peers mismatch: got %v, want %v", peers, []*Peer{peer})
		}
	case <-time.After(5 * time.Second):
		t.Error("server did not accept within five second")
	}
}

func TestServerNoListen(t *testing.T) {
	// start the test server
	connected := make(chan *Peer)
	remid := discover.PubkeyID(&newkey().PublicKey)
	srv := startTestServer(t, remid, func(p *Peer) {
		if p.ID() != remid {
			t.Error("peer func called with wrong node id")
		}
		if p == nil {
			t.Error("peer func called with nil conn")
		}
		connected <- p
	}, &Config{NoListen: true})
	defer close(connected)
	defer srv.Stop()

	// dial the test server that will be failed
	dialer := net.Dialer{Timeout: 10 * time.Millisecond}
	_, err := dialer.Dial("tcp", srv.GetListenAddress()[ConnDefault])
	if err == nil {
		t.Fatalf("server started with listening")
	}
}

func TestServerDial(t *testing.T) {
	// run a one-shot TCP server to handle the connection.
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not setup listener: %v", err)
	}
	defer listener.Close()
	accepted := make(chan net.Conn)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error("accept error:", err)
			return
		}

		c := makeconn(conn, discover.PubkeyID(&newkey().PublicKey))
		c.doConnTypeHandshake(c.conntype)
		accepted <- conn
	}()

	// start the server
	connected := make(chan *Peer)
	remid := discover.PubkeyID(&newkey().PublicKey)
	srv := startTestServer(t, remid, func(p *Peer) { connected <- p }, &Config{})
	defer close(connected)
	defer srv.Stop()

	// tell the server to connect
	tcpAddr := listener.Addr().(*net.TCPAddr)
	srv.AddPeer(&discover.Node{ID: remid, IP: tcpAddr.IP, TCP: uint16(tcpAddr.Port)})

	select {
	case conn := <-accepted:
		defer conn.Close()

		select {
		case peer := <-connected:
			if peer.ID() != remid {
				t.Errorf("peer has wrong id")
			}
			if peer.Name() != "test" {
				t.Errorf("peer has wrong name")
			}
			if peer.RemoteAddr().String() != conn.LocalAddr().String() {
				t.Errorf("peer started with wrong conn: got %v, want %v",
					peer.RemoteAddr(), conn.LocalAddr())
			}
			peers := srv.Peers()
			if !reflect.DeepEqual(peers, []*Peer{peer}) {
				t.Errorf("Peers mismatch: got %v, want %v", peers, []*Peer{peer})
			}
		case <-time.After(1 * time.Second):
			t.Error("server did not launch peer within one second")
		}

	case <-time.After(1 * time.Second):
		t.Error("server did not connect within one second")
	}
}

// This test checks that connections are disconnected
// just after the encryption handshake when the server is
// at capacity. Trusted connections should still be accepted.
func TestServerAtCap(t *testing.T) {
	trustedID := randomID()
	srv := &SingleChannelServer{
		BaseServer: &BaseServer{
			Config: Config{
				// BOOTNODE: no per-type reservation, isolates the global cap.
				ConnectionType:         common.BOOTNODE,
				PrivateKey:             newkey(),
				MaxPhysicalConnections: 10,
				NoDial:                 true,
				TrustedNodes:           []*discover.Node{{ID: trustedID}},
			},
		},
	}
	if err := srv.Start(); err != nil {
		t.Fatalf("could not start: %v", err)
	}
	defer srv.Stop()

	newconn := func(id discover.NodeID) *conn {
		fd, _ := net.Pipe()
		tx := newTestTransport(id, fd, nil, false)
		return &conn{fd: fd, transport: tx, flags: inboundConn, conntype: common.ConnTypeUndefined, id: id}
	}

	// Inject a few connections to fill up the peer set.
	for i := range 10 {
		c := newconn(randomID())
		if err := srv.handleAddPeerConn(c); err != nil {
			t.Fatalf("could not add conn %d: %v", i, err)
		}
	}
	// Try inserting a non-trusted connection.
	c := newconn(randomID())
	if err := srv.handlePostHandshake(c); err != DiscTooManyPeers {
		t.Error("wrong error for insert:", err)
	}
	// Try inserting a trusted connection.
	c = newconn(trustedID)
	if err := srv.handlePostHandshake(c); err != nil {
		t.Error("unexpected error for trusted conn @posthandshake:", err)
	}
	if !c.is(trustedConn) {
		t.Error("Server did not set trusted flag")
	}
}

func TestCountInboundPeers(t *testing.T) {
	// Inbound capacity is measured in peers, not sockets: a multichannel peer
	// counts once, and a mixed-direction peer is classified by its default
	// channel (ConnDefault), matching Peer.Inbound.
	peer := func(channelFlags ...connFlag) *Peer {
		rws := make([]*conn, len(channelFlags))
		for i, f := range channelFlags {
			rws[i] = &conn{flags: f}
		}
		return &Peer{rws: rws}
	}
	peers := map[discover.NodeID]*Peer{
		randomID(): peer(inboundConn, inboundConn),     // multichannel inbound -> counts
		randomID(): peer(inboundConn),                  // single-channel inbound -> counts
		randomID(): peer(inboundConn, dynDialedConn),   // default inbound, sub outbound -> counts
		randomID(): peer(dynDialedConn, dynDialedConn), // multichannel outbound -> not counted
		randomID(): peer(dynDialedConn, inboundConn),   // default outbound, sub inbound -> not counted
	}
	// Counted: the three peers whose default channel is inbound.
	if got := countInboundPeers(peers); got != 3 {
		t.Errorf("countInboundPeers = %d, want 3", got)
	}
}

func TestServerPeerTargets(t *testing.T) {
	newPeerWithType := func(connType common.ConnType) *Peer {
		return &Peer{rws: []*conn{{conntype: connType}}}
	}
	newPeers := func(connTypes ...common.ConnType) map[discover.NodeID]*Peer {
		peers := make(map[discover.NodeID]*Peer, len(connTypes))
		for _, connType := range connTypes {
			peers[randomID()] = newPeerWithType(connType)
		}
		return peers
	}
	newNPeers := func(connType common.ConnType, n int) map[discover.NodeID]*Peer {
		peers := make(map[discover.NodeID]*Peer, n)
		for range n {
			peers[randomID()] = newPeerWithType(connType)
		}
		return peers
	}

	cnSrv := &BaseServer{Config: Config{ConnectionType: common.CONSENSUSNODE, MaxPhysicalConnections: 10}}
	// CN-mesh accept cap = MaxPhysicalConnections - reserved cross-type slots (=7), reserving EN slots.
	if cnSrv.exceedsPeerTarget(newNPeers(common.CONSENSUSNODE, 6), &conn{conntype: common.CONSENSUSNODE}) {
		t.Fatal("CN mesh under the reserved cap should be accepted")
	}
	if !cnSrv.exceedsPeerTarget(newNPeers(common.CONSENSUSNODE, 7), &conn{conntype: common.CONSENSUSNODE}) {
		t.Fatal("CN mesh at the reserved cap (maxPhys - reserved cross-type slots) should be rejected")
	}
	cnPeers := newPeers(common.ENDPOINTNODE, common.PROXYNODE, common.ENDPOINTNODE)
	if !cnSrv.exceedsPeerTarget(cnPeers, &conn{conntype: common.ENDPOINTNODE}) {
		t.Fatal("CN should reject EN when EN-equivalent peer target is full")
	}
	if !cnSrv.exceedsPeerTarget(cnPeers, &conn{conntype: common.PROXYNODE}) {
		t.Fatal("CN should count PN against the EN-equivalent peer target")
	}
	if cnSrv.exceedsPeerTarget(cnPeers, &conn{conntype: common.ENDPOINTNODE, flags: trustedConn}) {
		t.Fatal("trusted peer should bypass peer target")
	}
	if cnSrv.exceedsPeerTarget(cnPeers, &conn{conntype: common.ENDPOINTNODE, flags: staticDialedConn}) {
		t.Fatal("static outbound peer should bypass peer target")
	}
	if !cnSrv.exceedsPeerTarget(cnPeers, &conn{conntype: common.ENDPOINTNODE, flags: inboundConn}) {
		t.Fatal("static/dynamic inbound peer should not bypass peer target unless trusted")
	}

	enSrv := &BaseServer{Config: Config{ConnectionType: common.ENDPOINTNODE, MaxPhysicalConnections: 10}}
	// EN-mesh accept cap = MaxPhysicalConnections - reserved cross-type slots (=8), reserving CN slots.
	if enSrv.exceedsPeerTarget(newNPeers(common.ENDPOINTNODE, 7), &conn{conntype: common.ENDPOINTNODE}) {
		t.Fatal("EN mesh under the reserved cap should be accepted")
	}
	if !enSrv.exceedsPeerTarget(newNPeers(common.ENDPOINTNODE, 8), &conn{conntype: common.ENDPOINTNODE}) {
		t.Fatal("EN mesh at the reserved cap (maxPhys - reserved cross-type slots) should be rejected")
	}
	enPeers := newPeers(common.CONSENSUSNODE, common.CONSENSUSNODE)
	if !enSrv.exceedsPeerTarget(enPeers, &conn{conntype: common.CONSENSUSNODE}) {
		t.Fatal("EN should reject CN when CN peer target is full")
	}
}

// A CN keeps one of its cross-type slots for an EN it dials itself, so inbound
// EN/PN peers alone cannot make dialTargets[CN][EN] unsatisfiable.
func TestServerCrossTypeDialReserve(t *testing.T) {
	newPeer := func(connType common.ConnType, flags connFlag) *Peer {
		return &Peer{rws: []*conn{{conntype: connType, flags: flags}}}
	}
	peersOf := func(ps ...*Peer) map[discover.NodeID]*Peer {
		peers := make(map[discover.NodeID]*Peer, len(ps))
		for _, p := range ps {
			peers[randomID()] = p
		}
		return peers
	}

	// defaultCNReservedSlots = 3 and dialTargets[CN][EN] = 1, so a CN takes at
	// most two inbound EN-equivalent peers.
	cnSrv := &BaseServer{Config: Config{ConnectionType: common.CONSENSUSNODE, MaxPhysicalConnections: 10}}
	inboundEN := &conn{conntype: common.ENDPOINTNODE, flags: inboundConn}
	dialedEN := &conn{conntype: common.ENDPOINTNODE, flags: dynDialedConn}

	oneInbound := peersOf(newPeer(common.ENDPOINTNODE, inboundConn))
	if cnSrv.exceedsPeerTarget(oneInbound, inboundEN) {
		t.Fatal("a second inbound EN should be accepted")
	}

	twoInbound := peersOf(
		newPeer(common.ENDPOINTNODE, inboundConn),
		newPeer(common.PROXYNODE, inboundConn),
	)
	if !cnSrv.exceedsPeerTarget(twoInbound, inboundEN) {
		t.Fatal("a third inbound EN-equivalent should not take the slot reserved for a dialed peer")
	}
	if cnSrv.exceedsPeerTarget(twoInbound, dialedEN) {
		t.Fatal("the dialed EN should still be admitted into the reserved slot")
	}
	if cnSrv.exceedsPeerTarget(twoInbound, &conn{conntype: common.ENDPOINTNODE, flags: inboundConn | trustedConn}) {
		t.Fatal("trusted peers should stay exempt from the inbound sub-cap")
	}

	// The sub-cap counts inbound peers only, so the bucket can still fill to the cap.
	withDialed := peersOf(
		newPeer(common.ENDPOINTNODE, dynDialedConn),
		newPeer(common.ENDPOINTNODE, inboundConn),
	)
	if cnSrv.exceedsPeerTarget(withDialed, inboundEN) {
		t.Fatal("an inbound EN should fill the last slot once a dialed EN holds the reservation")
	}
}

// CN/EN servers reject untrusted non-CN/EN peers (which would escape the
// per-type cap); trusted/static and non-CN/EN servers accept any type.
func TestAdmissiblePeerType(t *testing.T) {
	cnSrv := &BaseServer{Config: Config{ConnectionType: common.CONSENSUSNODE}}
	enSrv := &BaseServer{Config: Config{ConnectionType: common.ENDPOINTNODE}}
	bnSrv := &BaseServer{Config: Config{ConnectionType: common.BOOTNODE}}

	for _, peerType := range []common.ConnType{common.CONSENSUSNODE, common.ENDPOINTNODE, common.PROXYNODE} {
		if !cnSrv.admissiblePeerType(&conn{conntype: peerType, flags: inboundConn}) {
			t.Fatalf("CN server should admit peer type %v", peerType)
		}
		if !enSrv.admissiblePeerType(&conn{conntype: peerType, flags: inboundConn}) {
			t.Fatalf("EN server should admit peer type %v", peerType)
		}
	}

	if cnSrv.admissiblePeerType(&conn{conntype: common.BOOTNODE, flags: inboundConn}) {
		t.Fatal("CN server should reject an inbound bootnode peer")
	}
	if enSrv.admissiblePeerType(&conn{conntype: common.UNKNOWNNODE, flags: inboundConn}) {
		t.Fatal("EN server should reject an inbound unknown-type peer")
	}

	if !cnSrv.admissiblePeerType(&conn{conntype: common.BOOTNODE, flags: trustedConn}) {
		t.Fatal("trusted peer should bypass the peer-type check")
	}
	if !cnSrv.admissiblePeerType(&conn{conntype: common.BOOTNODE, flags: staticDialedConn}) {
		t.Fatal("static outbound peer should bypass the peer-type check")
	}

	if !bnSrv.admissiblePeerType(&conn{conntype: common.BOOTNODE, flags: inboundConn}) {
		t.Fatal("bootnode server should keep accepting any peer type")
	}
}

func TestServerCNPeersAdmission(t *testing.T) {
	key := newkey()
	id := discover.PubkeyID(&key.PublicKey)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	srv := &BaseServer{}

	c := &conn{conntype: common.CONSENSUSNODE, id: id}
	srv.SetCNPeers(nil)
	if err := srv.admitByCNPeers(c); err != nil {
		t.Fatalf("nil CN peers should disable CN filtering: %v", err)
	}

	srv.SetCNPeers([]common.Address{})
	if err := srv.admitByCNPeers(c); err != DiscUselessPeer {
		t.Fatalf("empty CN peers should reject CN claims, got %v", err)
	}

	if err := srv.admitByCNPeers(&conn{conntype: common.ENDPOINTNODE, id: id}); err != nil {
		t.Fatalf("EN should bypass CN peer admission: %v", err)
	}
	if err := srv.admitByCNPeers(&conn{conntype: common.PROXYNODE, id: id}); err != nil {
		t.Fatalf("legacy PN should bypass CN peer admission as EN-equivalent: %v", err)
	}
	if err := srv.admitByCNPeers(&conn{conntype: common.CONSENSUSNODE, id: id, flags: trustedConn}); err != nil {
		t.Fatalf("trusted CN should bypass CN peer admission: %v", err)
	}
	if err := srv.admitByCNPeers(&conn{conntype: common.CONSENSUSNODE, id: id, flags: staticDialedConn}); err != nil {
		t.Fatalf("static outbound CN should bypass CN peer admission: %v", err)
	}
	if err := srv.admitByCNPeers(&conn{conntype: common.CONSENSUSNODE, id: id, flags: inboundConn}); err != DiscUselessPeer {
		t.Fatalf("inbound CN should not bypass CN peer admission, got %v", err)
	}

	srv.SetCNPeers([]common.Address{addr})
	if err := srv.admitByCNPeers(c); err != nil {
		t.Fatalf("CN in CN peers should pass admission: %v", err)
	}
}

// Only CNs enforce the CN allowlist. EN/PN must keep serving any peer, so their
// SetCNPeers is a no-op and admitByCNPeers admits CNs outside any allowlist.
func TestServerENBypassesCNPeerFilter(t *testing.T) {
	allowed := crypto.PubkeyToAddress(newkey().PublicKey)

	for _, nodeType := range []common.ConnType{common.ENDPOINTNODE, common.PROXYNODE} {
		srv := &BaseServer{Config: Config{ConnectionType: nodeType}}
		srv.SetCNPeers([]common.Address{allowed}) // no-op on EN/PN
		if srv.cnPeerAddrs != nil {
			t.Fatalf("%v must not populate cnPeerAddrs", nodeType)
		}
		outsider := &conn{conntype: common.CONSENSUSNODE, id: randomID(), flags: inboundConn}
		if err := srv.admitByCNPeers(outsider); err != nil {
			t.Fatalf("%v must admit a CN outside any allowlist, got %v", nodeType, err)
		}
	}
}

func TestPeerTargetFor(t *testing.T) {
	const (
		maxPhys    = 100
		cnReserved = defaultCNReservedSlots // slots a CN keeps for ENs
		enReserved = defaultENReservedSlots // slots an EN keeps for CNs
	)
	cases := []struct {
		self, peer common.ConnType
		reserved   int
		want       int
		wantOK     bool
	}{
		{common.CONSENSUSNODE, common.CONSENSUSNODE, cnReserved, maxPhys - cnReserved, true}, // CN mesh
		{common.CONSENSUSNODE, common.ENDPOINTNODE, cnReserved, cnReserved, true},            // CN->EN reservation
		{common.ENDPOINTNODE, common.ENDPOINTNODE, enReserved, maxPhys - enReserved, true},   // EN mesh
		{common.ENDPOINTNODE, common.CONSENSUSNODE, enReserved, enReserved, true},            // EN->CN reservation
		{common.CONSENSUSNODE, common.BOOTNODE, cnReserved, 0, false},                        // unrelated peer type: no cap
		{common.BOOTNODE, common.CONSENSUSNODE, cnReserved, 0, false},                        // node with no reservation: no cap
	}
	for _, c := range cases {
		if got, ok := peerTargetFor(c.self, c.peer, maxPhys, c.reserved); got != c.want || ok != c.wantOK {
			t.Errorf("peerTargetFor(%v,%v)=(%d,%v), want (%d,%v)", c.self, c.peer, got, ok, c.want, c.wantOK)
		}
	}
	// own-mesh cap clamps to 0 when maxPhys is below the reservation.
	if got, _ := peerTargetFor(common.CONSENSUSNODE, common.CONSENSUSNODE, cnReserved-1, cnReserved); got != 0 {
		t.Errorf("own-mesh cap should clamp to 0, got %d", got)
	}
}

func TestMaxPhysicalConnectionsLowerBound(t *testing.T) {
	const (
		cnReserved = defaultCNReservedSlots
		enReserved = defaultENReservedSlots
	)
	// An unset reservation falls back to the per-node-type default.
	if got := MaxPhysicalConnectionsLowerBound(Config{ConnectionType: common.CONSENSUSNODE}); got != cnReserved+1 {
		t.Errorf("CN lower bound with default reservation = %d, want %d", got, cnReserved+1)
	}
	if got := MaxPhysicalConnectionsLowerBound(Config{ConnectionType: common.ENDPOINTNODE}); got != enReserved+1 {
		t.Errorf("EN lower bound with default reservation = %d, want %d", got, enReserved+1)
	}
	// PROXYNODE resolves as EN.
	if got := MaxPhysicalConnectionsLowerBound(Config{ConnectionType: common.PROXYNODE}); got != enReserved+1 {
		t.Errorf("PN lower bound (as EN) = %d, want %d", got, enReserved+1)
	}
	// An explicit reservation overrides the default.
	if got := MaxPhysicalConnectionsLowerBound(Config{ConnectionType: common.CONSENSUSNODE, ReservedCrossTypeSlots: 5}); got != 6 {
		t.Errorf("CN lower bound with reserved 5 = %d, want 6", got)
	}
}

func TestServerSetCNPeersReconcilesConnectedCNPeers(t *testing.T) {
	newPeerWithKey := func(key *ecdsa.PrivateKey, connType common.ConnType, flags connFlag) *Peer {
		return &Peer{
			rws:  []*conn{{id: discover.PubkeyID(&key.PublicKey), conntype: connType, flags: flags}},
			disc: make(chan DiscReason, 1),
		}
	}
	assertNotDisconnected := func(t *testing.T, p *Peer) {
		t.Helper()
		select {
		case reason := <-p.disc:
			t.Fatalf("peer should remain connected, got disconnect %v", reason)
		default:
		}
	}

	allowedKey := newkey()
	dropKey := newkey()
	trustedKey := newkey()
	staticKey := newkey()
	enKey := newkey()

	allowed := newPeerWithKey(allowedKey, common.CONSENSUSNODE, inboundConn)
	drop := newPeerWithKey(dropKey, common.CONSENSUSNODE, inboundConn)
	trusted := newPeerWithKey(trustedKey, common.CONSENSUSNODE, inboundConn|trustedConn)
	staticOutbound := newPeerWithKey(staticKey, common.CONSENSUSNODE, staticDialedConn)
	en := newPeerWithKey(enKey, common.ENDPOINTNODE, inboundConn)

	srv := &BaseServer{
		logger: logger.NewWith(),
		peers: map[discover.NodeID]*Peer{
			allowed.ID():        allowed,
			drop.ID():           drop,
			trusted.ID():        trusted,
			staticOutbound.ID(): staticOutbound,
			en.ID():             en,
		},
	}

	srv.SetCNPeers([]common.Address{crypto.PubkeyToAddress(allowedKey.PublicKey)})

	select {
	case reason := <-drop.disc:
		if reason != DiscRequested {
			t.Fatalf("wrong disconnect reason: %v", reason)
		}
	default:
		t.Fatal("CN peer outside allowlist should be disconnected")
	}
	assertNotDisconnected(t, allowed)
	assertNotDisconnected(t, trusted)
	assertNotDisconnected(t, staticOutbound)
	assertNotDisconnected(t, en)
}

func TestServerSetCNPeersUpdatesDialSched(t *testing.T) {
	allowedKey := newkey()
	blockedKey := newkey()
	allowed := discover.NewNode(discover.PubkeyID(&allowedKey.PublicKey), net.ParseIP("10.0.0.1"), 30303, 30303, nil, discover.NodeTypeCN)
	blocked := discover.NewNode(discover.PubkeyID(&blockedKey.PublicKey), net.ParseIP("10.0.0.2"), 30303, 30303, nil, discover.NodeTypeCN)
	ds := NewDialSched(DialConfig{selfType: discover.NodeTypeCN}, nil, nil)
	srv := &BaseServer{
		dialSched: ds,
		peers:     make(map[discover.NodeID]*Peer),
	}

	srv.SetCNPeers([]common.Address{crypto.PubkeyToAddress(allowedKey.PublicKey)})
	if !ds.dynamicCandidateAllowed(discover.NodeTypeCN, allowed) {
		t.Fatal("allowed CN should remain dialable")
	}
	if ds.dynamicCandidateAllowed(discover.NodeTypeCN, blocked) {
		t.Fatal("blocked CN should not remain dialable")
	}

	srv.SetCNPeers(nil)
	if !ds.dynamicCandidateAllowed(discover.NodeTypeCN, blocked) {
		t.Fatal("nil CN peer allowlist should disable dynamic dial filtering")
	}
}

func TestServerSetupConn(t *testing.T) {
	var (
		id     = discover.PubkeyID(&newkey().PublicKey)
		srvkey = newkey()
		srvid  = discover.PubkeyID(&srvkey.PublicKey)
	)

	tests := []struct {
		dontstart bool
		tt        *setupTransport
		flags     connFlag
		dialDest  *discover.Node

		wantCloseErr error
		wantCalls    string
	}{
		{
			dontstart:    true,
			tt:           &setupTransport{id: id},
			wantCalls:    "close,",
			wantCloseErr: errServerStopped,
		},
		{
			tt:           &setupTransport{id: id, encHandshakeErr: errors.New("read error")},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,close,",
			wantCloseErr: errors.New("read error"),
		},
		{
			tt:           &setupTransport{id: id, phs: &protoHandshake{ID: randomID()}},
			dialDest:     &discover.Node{ID: id, NType: discover.NodeType(common.ENDPOINTNODE)},
			flags:        dynDialedConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: DiscUnexpectedIdentity,
		},
		{
			tt:           &setupTransport{id: id, protoHandshakeErr: errors.New("foo")},
			dialDest:     &discover.Node{ID: id, NType: discover.NodeType(common.ENDPOINTNODE)},
			flags:        dynDialedConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: errors.New("foo"),
		},
		{
			tt:           &setupTransport{id: srvid, phs: &protoHandshake{ID: srvid}},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,close,",
			wantCloseErr: DiscSelf,
		},
		{
			tt:           &setupTransport{id: id, phs: &protoHandshake{ID: id}},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: DiscUselessPeer,
		},
	}

	for i, test := range tests {
		srv := &SingleChannelServer{
			&BaseServer{
				Config: Config{
					PrivateKey:             srvkey,
					MaxPhysicalConnections: 10,
					NoDial:                 true,
					Protocols:              []Protocol{discard},
					ConnectionType:         1, // ENDPOINTNODE
				},
				newTransport: func(fd net.Conn, dialDest *ecdsa.PublicKey) transport { return test.tt },
				logger:       logger.NewWith(),
			},
		}
		if !test.dontstart {
			if err := srv.Start(); err != nil {
				t.Fatalf("couldn't start server: %v", err)
			}
		}
		p1, _ := net.Pipe()
		srv.SetupConn(p1, test.flags, test.dialDest)
		if !reflect.DeepEqual(test.tt.closeErr, test.wantCloseErr) {
			t.Errorf("test %d: close error mismatch: got %q, want %q", i, test.tt.closeErr, test.wantCloseErr)
		}
		if test.tt.calls != test.wantCalls {
			t.Errorf("test %d: calls mismatch: got %q, want %q", i, test.tt.calls, test.wantCalls)
		}
	}
}

type setupTransport struct {
	id              discover.NodeID
	encHandshakeErr error

	phs               *protoHandshake
	protoHandshakeErr error

	calls    string
	closeErr error
}

func (c *setupTransport) doConnTypeHandshake(myConnType common.ConnType) (common.ConnType, error) {
	return 1, nil
}

func (c *setupTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	c.calls += "doEncHandshake,"
	pubkey, _ := c.id.Pubkey()
	return pubkey, c.encHandshakeErr
}

func (c *setupTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	c.calls += "doProtoHandshake,"
	if c.protoHandshakeErr != nil {
		return nil, c.protoHandshakeErr
	}
	return c.phs, nil
}

func (c *setupTransport) close(err error) {
	c.calls += "close,"
	c.closeErr = err
}

// WriteMsg shouldn't write to/read from the connection.
func (c *setupTransport) WriteMsg(Msg) error {
	panic("WriteMsg called on setupTransport")
}

func (c *setupTransport) ReadMsg() (Msg, error) {
	panic("ReadMsg called on setupTransport")
}

func newkey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic("couldn't generate key: " + err.Error())
	}
	return key
}

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}
