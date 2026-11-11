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
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/stretchr/testify/require"
)

// A static KNI must assign the same channels to the dialed sockets as the
// listener does. Exercise the real dial scheduler, TCP sockets and RLPx transport.
func TestMultiChannelStaticPeerPortOrder(t *testing.T) {
	started := make(chan *Peer, 2)
	handshakes := make(chan error, 2)
	protocol := Protocol{
		Name: "portorder", Version: 1, Length: 1,
		RunWithRWs: func(peer *Peer, rws []MsgReadWriter) error {
			started <- peer
			for channel, rw := range rws {
				// Like the consensus Status exchange, start by sending and
				// receiving on channel 0. Also check the transaction channel.
				if err := Send(rw, 0, uint64(channel)); err != nil {
					handshakes <- err
					return err
				}
				if err := ExpectMsg(rw, 0, uint64(channel)); err != nil {
					handshakes <- err
					return err
				}
			}
			handshakes <- nil
			<-peer.closed
			return io.EOF
		},
	}
	newServer := func() *MultiChannelServer {
		return &MultiChannelServer{BaseServer: &BaseServer{Config: Config{
			PrivateKey:             newkey(),
			ConnectionType:         common.CONSENSUSNODE,
			MaxPhysicalConnections: 10,
			NoDiscovery:            true,
			ListenAddr:             "127.0.0.1:0",
			SubListenAddr:          []string{"127.0.0.1:0"},
			Protocols:              []Protocol{protocol},
			RWTimerConfig:          defaultRWTimerConfig,
		}}}
	}
	listener := newServer()
	listener.NoDial = true
	t.Cleanup(listener.Stop)
	require.NoError(t, listener.Start())

	_, subport, err := net.SplitHostPort(listener.ListenAddrs[1])
	require.NoError(t, err)
	url := fmt.Sprintf("kni://%s@%s?subport=%s&discport=0&ntype=cn",
		listener.selfID, listener.ListenAddrs[ConnDefault], subport)
	staticNode, err := discover.ParseNode(url)
	require.NoError(t, err)

	dialer := newServer()
	dialer.StaticNodes = []*discover.Node{staticNode}
	t.Cleanup(dialer.Stop)
	require.NoError(t, dialer.Start())

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	for range 2 {
		select {
		case peer := <-started:
			require.Len(t, peer.rws, 2)
			for channel, c := range peer.rws {
				require.Equal(t, PortOrder(channel), c.portOrder)
				addr := c.fd.RemoteAddr()
				if peer.Inbound() {
					addr = c.fd.LocalAddr()
				}
				require.Equal(t, listener.ListenAddrs[channel], addr.String(),
					"channel %d must use the same listener port on both ends", channel)
			}
		case <-timeout.C:
			t.Fatal("static multichannel peer did not start on both ends")
		}
	}
	for range 2 {
		select {
		case err := <-handshakes:
			require.NoError(t, err)
		case <-timeout.C:
			t.Fatal("protocol handshake did not complete on both channels")
		}
	}
}
