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

package discover

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable2_Bond(t *testing.T) {
	config := &Config{
		Id:         NodeID{},
		Addr:       &net.UDPAddr{},
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	transport := newPingRecorder()
	tab, err := newTable2(config, transport)
	require.NoError(t, err)
	defer tab.Close()

	aliveNode := NewNode(MustHexID("a502af0f59b2aab7746995408c79e9ca312d2793cc997e44fc55eda62f0150bbb8c59a6f9269ba3a081518b62699ee807c7c19c20125ddfccca872608af9e370"),
		net.IP{}, 99, 99, nil, NodeTypeUnknown)
	transport.dead[aliveNode.ID] = false

	err = tab.Bond(false, aliveNode)
	require.NoError(t, err)
	assert.True(t, transport.pinged[aliveNode.ID])
	assert.Nil(t, tab.bonding[aliveNode.ID])

	deadNode := NewNode(MustHexID("b502af0f59b2aab7746995408c79e9ca312d2793cc997e44fc55eda62f0150bbb8c59a6f9269ba3a081518b62699ee807c7c19c20125ddfccca872608af9e370"),
		net.IP{}, 99, 99, nil, NodeTypeUnknown)
	transport.dead[deadNode.ID] = true

	err = tab.Bond(false, deadNode)
	require.ErrorIs(t, err, errTimeout)
	assert.True(t, transport.pinged[deadNode.ID]) // still attempt to ping.
}

func TestTable2_Lookup(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)
	self := nodeAtDistance(common.Hash{}, 0, NodeTypeEN)
	config := &Config{
		Id:         self.ID,
		Addr:       &net.UDPAddr{},
		NodeType:   NodeTypeEN,
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	transport := lookupTestnet
	tab, err := newTable2(config, transport)
	require.NoError(t, err)
	defer tab.Close()

	// Lookup without any seed nodes returns no nodes.
	results := tab.lookup(lookupTestnet.target, NodeTypeEN, true, bucketSize)
	assert.Empty(t, results)

	// Push a seed node and retry.
	seed := NewNode(lookupTestnet.dists[256][0], net.IP{}, 256, 0, nil, NodeTypeEN)
	tab.addNode(seed)
	require.Equal(t, 1, tab.len())

	results = tab.lookup(lookupTestnet.target, NodeTypeEN, true, bucketSize)
	assert.Equal(t, bucketSize, len(results))

	// max=0 returns no nodes.
	results = tab.lookup(lookupTestnet.target, NodeTypeEN, true, 0)
	assert.Empty(t, results)
}

// manually test bonding against a local EN.
func TestTable_KairosExample(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	// t.Skip("manual test; comment me to run")

	oldLookupIPFunc := lookupIPFunc
	lookupIPFunc = net.LookupIP
	defer func() { lookupIPFunc = oldLookupIPFunc }()
	bn := MustParseNode(params.KairosBootnodes[common.ENDPOINTNODE].Addrs[0])

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:32326")
	require.NoError(t, err)
	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)

	config := &Config{
		NetworkID:  1001,
		PrivateKey: newkey(),

		NodeDBPath: "",
		Id:         NodeID{},
		Addr:       conn.LocalAddr().(*net.UDPAddr), // Local address
		Conn:       conn,
		NodeType:   NodeTypeEN,
		Bootnodes:  []*Node{bn},
	}
	d, err := newDiscovery(config)
	require.NoError(t, err)
	defer d.Close()

	d.tab.Refresh()
	fmt.Println("Refresh done")

	time.Sleep(30 * time.Second)

}
