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

package discover

import (
	"net"
	"os"
	"testing"

	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiscovery_MainnetBootnodes is an end-to-end integration test that
// connects to the real Kaia mainnet bootnodes and discovers peers from an EN's perspective.
func TestDiscovery_MainnetBootnodes(t *testing.T) {
	// Enable only if P2P_DISCOVER_E2E_TEST is set
	if os.Getenv("P2P_DISCOVER_E2E_TEST") == "" {
		t.Skip("skipping p2p/discover end-to-end test")
	}
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)

	// Load mainnet bootnodes. KIP-311 BNs are shared by all node types.
	urls := params.MainnetBootnodes
	oldLookupIPFunc := lookupIPFunc
	lookupIPFunc = net.LookupIP
	defer func() {
		lookupIPFunc = oldLookupIPFunc
	}()
	bootnodes := make([]*Node, 0, len(urls))
	for _, url := range urls {
		node, err := ParseNode(url)
		require.NoError(t, err, "failed to parse bootnode URL: %s", url)
		bootnodes = append(bootnodes, node)
	}

	// Generate an ephemeral key for this test node
	key, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Listen on a random UDP port
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	assert.NoError(t, err)
	defer udpConn.Close()

	addr := udpConn.LocalAddr().(*net.UDPAddr)
	t.Logf("Listening on %v", addr)

	// Start the discovery stack.
	cfg := &Config{
		NetworkID:  params.MainnetNetworkId,
		PrivateKey: key,
		Bootnodes:  bootnodes,
		Id:         PubkeyID(&key.PublicKey),
		Addr:       addr,
		Conn:       udpConn,
		NodeType:   NodeTypeEN,
		DiscoverTypes: DiscoverTypesConfig{
			EN: true,
		},
	}
	disc, err := NewDiscovery2(cfg)
	assert.NoError(t, err)
	defer disc.Close()

	// Refresh contacts the bootnodes and populates the table
	disc.Refresh()
	t.Logf("Refreshed the discovery table")

	// Print nodes in the table
	for _, ty := range []NodeType{NodeTypeCN, NodeTypePN, NodeTypeEN, NodeTypeBN} {
		buf := make([]*Node, 100)
		found := disc.RandomNodes(buf, ty)
		t.Logf("Found %d %s nodes (including initial bootnodes)", found, StringNodeType(ty))
		for i := 0; i < found && i < 10; i++ {
			n := buf[i]
			t.Logf("  [%d] %x@%s:%d (type=%s)", i, n.ID[:8], n.IP, n.UDP, StringNodeType(n.NType))
		}
		if found > 10 {
			t.Logf("  ... and more")
		}
	}
}
