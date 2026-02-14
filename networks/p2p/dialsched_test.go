// Modifications Copyright 2026 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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

package p2p

import (
	"net"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

func TestDialSched_Candidates(t *testing.T) {

}

// TODO: make gomock
type mockDiscovery2 struct {
	randomCount int
}

func (m *mockDiscovery2) RandomNodes(buf []*discover.Node, targetType discover.NodeType) int {
	for i := 0; i < m.randomCount && i < len(buf); i++ {
		buf[i] = &discover.Node{ID: discover.NodeID{byte(i + 1)}, NType: targetType}
	}
	return m.randomCount
}

func (m *mockDiscovery2) Close() {
}

func (m *mockDiscovery2) Refresh() {
	time.Sleep(time.Second * 5)
}

func TestDialSched_KairosExample(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	selfkey := newkey()
	selfid := discover.PubkeyID(&selfkey.PublicKey)
	bn := discover.MustParseNode(params.KairosBootnodes[common.ENDPOINTNODE].Addrs[0])

	// Start Discovery stack
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:32326")
	require.NoError(t, err)
	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)

	config := &discover.Config{
		NetworkID:  1001,
		PrivateKey: selfkey,

		NodeDBPath: "",
		Id:         selfid,
		Addr:       conn.LocalAddr().(*net.UDPAddr), // Local address
		Conn:       conn,
		NodeType:   discover.NodeTypeEN,
		Bootnodes:  []*discover.Node{bn},
	}
	d, err := discover.NewDiscovery2(config)
	require.NoError(t, err)
	defer d.Close()

	cfg := DialConfig{
		selfID:   selfid,
		selfType: discover.NodeTypeEN,
		connTargets: map[discover.NodeType]int{
			discover.NodeTypeEN: 10,
		},
	}
	ds := NewDialSched(cfg, d)
	ENds := ds.typedSchedulers[discover.NodeTypeEN]
	ENds.wg.Add(1)
	go ENds.dialLoop()
	defer ds.Close()

	time.Sleep(time.Second * 10)
	var id discover.NodeID
	for id = range ENds.connected {
		break
	}
	require.NotEmpty(t, id)

	ENds.OnFailure(id)

	time.Sleep(time.Second * 10)
}
