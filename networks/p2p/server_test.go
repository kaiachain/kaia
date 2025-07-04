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
	return &conn{fd: fd, transport: tx, flags: staticDialedConn, conntype: common.ConnTypeUndefined, id: id, cont: make(chan error)}
}

func makeMultiChannelConn(fd net.Conn, id discover.NodeID) *conn {
	dialDest, _ := id.Pubkey()
	tx := newTestTransport(id, fd, dialDest, true)
	return &conn{fd: fd, transport: tx, flags: staticDialedConn, conntype: common.ConnTypeUndefined, id: id, cont: make(chan error), multiChannel: true}
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
	conn, err := net.DialTimeout("tcp", srv.GetListenAddress()[ConnDefault], 5*time.Second)
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
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
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
	_, err := net.DialTimeout("tcp", srv.GetListenAddress()[ConnDefault], 10*time.Millisecond)
	if err == nil {
		t.Fatalf("server started with listening")
	}
}

func TestServerDial(t *testing.T) {
	// run a one-shot TCP server to handle the connection.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
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

// This test checks that tasks generated by dialstate are
// actually executed and taskdone is called for them.
func TestServerTaskScheduling(t *testing.T) {
	var (
		done           = make(chan *testTask)
		quit, returned = make(chan struct{}), make(chan struct{})
		tc             = 0
		tg             = taskgen{
			newFunc: func(running int, peers map[discover.NodeID]*Peer) []task {
				tc++
				return []task{&testTask{index: tc - 1}}
			},
			doneFunc: func(t task) {
				select {
				case done <- t.(*testTask):
				case <-quit:
				}
			},
		}
	)

	// The Server in this test isn't actually running
	// because we're only interested in what run does.
	srv := &SingleChannelServer{
		&BaseServer{
			Config:  Config{MaxPhysicalConnections: 10},
			quit:    make(chan struct{}),
			ntab:    fakeTable{},
			running: true,
			logger:  logger.NewWith(),
		},
	}
	srv.loopWG.Add(1)
	go func() {
		srv.run(tg)
		close(returned)
	}()

	var gotdone []*testTask
	for i := 0; i < 100; i++ {
		gotdone = append(gotdone, <-done)
	}
	for i, task := range gotdone {
		if task.index != i {
			t.Errorf("task %d has wrong index, got %d", i, task.index)
			break
		}
		if !task.called {
			t.Errorf("task %d was not called", i)
			break
		}
	}

	close(quit)
	srv.Stop()
	select {
	case <-returned:
	case <-time.After(500 * time.Millisecond):
		t.Error("Server.run did not return within 500ms")
	}
}

// This test checks that Server doesn't drop tasks,
// even if newTasks returns more than the maximum number of tasks.
func TestServerManyTasks(t *testing.T) {
	alltasks := make([]task, 300)
	for i := range alltasks {
		alltasks[i] = &testTask{index: i}
	}

	var (
		srv = &SingleChannelServer{
			&BaseServer{
				quit:    make(chan struct{}),
				ntab:    fakeTable{},
				running: true,
				logger:  logger.NewWith(),
			},
		}
		done       = make(chan *testTask)
		start, end = 0, 0
	)
	defer srv.Stop()
	srv.loopWG.Add(1)
	go srv.run(taskgen{
		newFunc: func(running int, peers map[discover.NodeID]*Peer) []task {
			start, end = end, end+maxActiveDialTasks+10
			if end > len(alltasks) {
				end = len(alltasks)
			}
			return alltasks[start:end]
		},
		doneFunc: func(tt task) {
			done <- tt.(*testTask)
		},
	})

	doneset := make(map[int]bool)
	timeout := time.After(2 * time.Second)
	for len(doneset) < len(alltasks) {
		select {
		case tt := <-done:
			if doneset[tt.index] {
				t.Errorf("task %d got done more than once", tt.index)
			} else {
				doneset[tt.index] = true
			}
		case <-timeout:
			t.Errorf("%d of %d tasks got done within 2s", len(doneset), len(alltasks))
			for i := 0; i < len(alltasks); i++ {
				if !doneset[i] {
					t.Logf("task %d not done", i)
				}
			}
			return
		}
	}
}

type taskgen struct {
	newFunc  func(running int, peers map[discover.NodeID]*Peer) []task
	doneFunc func(task)
}

func (tg taskgen) newTasks(running int, peers map[discover.NodeID]*Peer, now time.Time) []task {
	return tg.newFunc(running, peers)
}

func (tg taskgen) taskDone(t task, now time.Time) {
	tg.doneFunc(t)
}

func (tg taskgen) addStatic(*discover.Node) {
}

func (tg taskgen) removeStatic(*discover.Node) {
}

type testTask struct {
	index  int
	called bool
}

func (t *testTask) Do(srv Server) {
	t.called = true
}

// This test checks that connections are disconnected
// just after the encryption handshake when the server is
// at capacity. Trusted connections should still be accepted.
func TestServerAtCap(t *testing.T) {
	trustedID := randomID()
	srv := &SingleChannelServer{
		BaseServer: &BaseServer{
			Config: Config{
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
		return &conn{fd: fd, transport: tx, flags: inboundConn, conntype: common.ConnTypeUndefined, id: id, cont: make(chan error)}
	}

	// Inject a few connections to fill up the peer set.
	for i := 0; i < 10; i++ {
		c := newconn(randomID())
		if err := srv.checkpoint(c, srv.addpeer); err != nil {
			t.Fatalf("could not add conn %d: %v", i, err)
		}
	}
	// Try inserting a non-trusted connection.
	c := newconn(randomID())
	if err := srv.checkpoint(c, srv.posthandshake); err != DiscTooManyPeers {
		t.Error("wrong error for insert:", err)
	}
	// Try inserting a trusted connection.
	c = newconn(trustedID)
	if err := srv.checkpoint(c, srv.posthandshake); err != nil {
		t.Error("unexpected error for trusted conn @posthandshake:", err)
	}
	if !c.is(trustedConn) {
		t.Error("Server did not set trusted flag")
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
