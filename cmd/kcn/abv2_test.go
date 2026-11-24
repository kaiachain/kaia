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

package main

import (
	"encoding/hex"
	"math/big"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestEndpointOf(t *testing.T) {
	endpoint := func(args ...string) string {
		var got string
		cmd := &cli.Command{
			Name:   "endpoint",
			Flags:  abv2Flags(),
			Action: func(ctx *cli.Context) error { got = endpointOf(ctx); return nil },
		}
		app := &cli.App{Name: "kcn", Commands: []*cli.Command{cmd}}
		require.NoError(t, app.Run(append([]string{"kcn", "endpoint"}, args...)))
		return got
	}

	got := endpoint("--datadir", "/data/kcnd")
	assert.Equal(t, "/data/kcnd/klay.ipc", got)

	got = endpoint("--datadir", "/data/kcnd", "--endpoint", "http://localhost:8551")
	assert.Equal(t, "http://localhost:8551", got)

	// Without flags, the IPC socket of an rpm install.
	got = endpoint()
	assert.Equal(t, "/var/kcnd/data/klay.ipc", got)
}

func TestDatadirDefault(t *testing.T) {
	if _, err := os.Stat(defaultDatadir); err == nil {
		t.Skip(defaultDatadir + " exists on this host")
	}
	app := &cli.App{Name: "kcn", Commands: []*cli.Command{ValOpsCommand}}
	run := func(args ...string) error { return app.Run(append([]string{"kcn", "valops"}, args...)) }

	// Without --datadir, every command that reads the node key looks under /var/kcnd/data.
	wantPath := "/var/kcnd/data/klay/nodekey"
	assert.ErrorContains(t, run("sign-create-node", "--chain-id", "1001", "--manager", "0x1111111111111111111111111111111111111111", "--staking-contract", "0x2222222222222222222222222222222222222222"), wantPath)
	assert.ErrorContains(t, run("ready-candidate"), wantPath)
	assert.ErrorContains(t, run("suspend-validator", "--node-id", "0x1111111111111111111111111111111111111111"), wantPath)
}

func TestDatadirUsedAndSource(t *testing.T) {
	check := func(args ...string) (used bool, source string) {
		cmd := &cli.Command{
			Name:  "check",
			Flags: abv2Flags(),
			Action: func(ctx *cli.Context) error {
				used, source = datadirUsed(ctx), datadirSource(ctx)
				return nil
			},
		}
		app := &cli.App{Name: "kcn", Commands: []*cli.Command{cmd}}
		require.NoError(t, app.Run(append([]string{"kcn", "check"}, args...)))
		return used, source
	}

	used, source := check()
	assert.True(t, used)
	assert.Equal(t, "default", source)

	used, source = check("--datadir", "/data/kcnd")
	assert.True(t, used)
	assert.Equal(t, "--datadir", source)

	// The key or the endpoint alone still comes from the data directory.
	used, _ = check("--private-key", "0x01")
	assert.True(t, used)
	used, _ = check("--endpoint", "http://localhost:8551")
	assert.True(t, used)

	// Neither the key nor the endpoint comes from the data directory.
	used, _ = check("--private-key", "0x01", "--endpoint", "http://localhost:8551")
	assert.False(t, used)
}

func TestParseAddress(t *testing.T) {
	addr, err := parseAddress("manager", "0x00000000000000000000000000000000000000Aa")
	require.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0xaa"), addr)

	_, err = parseAddress("manager", "00000000000000000000000000000000000000aa")
	assert.NoError(t, err, "0x prefix is optional")

	for _, bad := range []string{"", "0x1234", "0x00000000000000000000000000000000000000aa00", "0x000000000000000000000000000000000000zzzz"} {
		_, err := parseAddress("manager", bad)
		assert.ErrorContains(t, err, "invalid --manager", bad)
	}
	_, err = parseAddress("manager", "0x0000000000000000000000000000000000000000")
	assert.ErrorContains(t, err, "zero address")
}

func TestFormatKAIA(t *testing.T) {
	assert.Equal(t, "13.750000 KAIA", formatKAIA(new(big.Int).Mul(big.NewInt(1375), big.NewInt(1e16))))
	assert.Equal(t, "0.000000 KAIA", formatKAIA(big.NewInt(0)))
}

func TestSignCreateNode_ValidatesInputs(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	keyHex := "0x" + hex.EncodeToString(crypto.FromECDSA(key))
	nodeId := crypto.PubkeyToAddress(key.PublicKey).Hex()
	manager := "0x1111111111111111111111111111111111111111"
	staking := "0x2222222222222222222222222222222222222222"

	run := func(args ...string) error {
		app := &cli.App{Name: "kcn", Commands: []*cli.Command{SignCreateNodeCommand}}
		return app.Run(append([]string{"kcn", "sign-create-node", "--private-key", keyHex}, args...))
	}

	require.NoError(t, run("--chain-id", "1001", "--manager", manager, "--staking-contract", staking, "--node-id", nodeId))

	cases := []struct {
		name   string
		args   []string
		errStr string
	}{
		{"zero chain id", []string{"--chain-id", "0", "--manager", manager, "--staking-contract", staking}, "invalid --chain-id"},
		{"short manager", []string{"--chain-id", "1001", "--manager", "0x1234", "--staking-contract", staking}, "invalid --manager"},
		{"zero staking", []string{"--chain-id", "1001", "--manager", manager, "--staking-contract", "0x0000000000000000000000000000000000000000"}, "invalid --staking-contract"},
		{"malformed node id", []string{"--chain-id", "1001", "--manager", manager, "--staking-contract", staking, "--node-id", "0x12"}, "invalid --node-id"},
		{"other node id", []string{"--chain-id", "1001", "--manager", manager, "--staking-contract", staking, "--node-id", manager}, "does not match"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorContains(t, run(tc.args...), tc.errStr)
		})
	}
}

func TestEpochTimeLeft(t *testing.T) {
	interval := big.NewInt(86400)
	cases := []struct {
		head     int64
		next     int64
		timeLeft time.Duration
		late     bool
	}{
		{head: 86400, next: 172800, timeLeft: 24 * time.Hour}, // head is an epoch block
		{head: 86401, next: 172800, timeLeft: 86399 * time.Second},
		{head: 172800 - 300, next: 172800, timeLeft: 5 * time.Minute},
		{head: 172800 - 299, next: 172800, timeLeft: 299 * time.Second, late: true},
		{head: 172799, next: 172800, timeLeft: time.Second, late: true},
	}
	for _, tc := range cases {
		next, left := epochTimeLeft(big.NewInt(tc.head), interval)
		assert.Equal(t, tc.next, next.Int64(), "head %d", tc.head)
		assert.Equal(t, tc.timeLeft, left, "head %d", tc.head)
		assert.Equal(t, tc.late, left < minEpochTimeLeft, "head %d", tc.head)
	}
}

// fakeNode serves the RPC methods a valops command calls.
type fakeNode struct {
	peers, latest, pending uint64
	bumpPending            bool   // a pending-nonce read after the first returns one more
	lookup                 string // what the node says about the sent transaction: "pool", "block" or "" (unknown)

	pendingReads int
	sent         []*types.Transaction
}

func (n *fakeNode) PeerCount() hexutil.Uint { return hexutil.Uint(n.peers) }

func (n *fakeNode) ChainID() *hexutil.Big { return (*hexutil.Big)(big.NewInt(1001)) }

func (n *fakeNode) GetTransactionCount(_ common.Address, block string) hexutil.Uint64 {
	if block != "pending" {
		return hexutil.Uint64(n.latest)
	}
	n.pendingReads++
	if n.bumpPending && n.pendingReads > 1 {
		return hexutil.Uint64(n.pending + 1)
	}
	return hexutil.Uint64(n.pending)
}

func (n *fakeNode) Call(map[string]interface{}, string) hexutil.Bytes { return nil }

func (n *fakeNode) EstimateGas(map[string]interface{}) hexutil.Uint64 { return 100000 }

func (n *fakeNode) GasPrice() *hexutil.Big { return (*hexutil.Big)(big.NewInt(25e9)) }

func (n *fakeNode) GetBalance(common.Address, string) *hexutil.Big {
	return (*hexutil.Big)(new(big.Int).SetUint64(params.KAIA))
}

func (n *fakeNode) SendRawTransaction(raw hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(raw, tx); err != nil {
		return common.Hash{}, err
	}
	n.sent = append(n.sent, tx)
	return tx.Hash(), nil
}

func (n *fakeNode) GetTransactionReceipt(common.Hash) map[string]interface{} { return nil }

func (n *fakeNode) GetTransactionByHash(hash common.Hash) map[string]interface{} {
	if len(n.sent) == 0 || hash != n.sent[0].Hash() {
		return nil
	}
	switch n.lookup {
	case "pool":
		return api.NewRPCPendingTransaction(n.sent[0], nil)
	case "block":
		return api.NewRPCTransaction(nil, n.sent[0], common.HexToHash("0x17"), 23, 0, nil)
	}
	return nil
}

func TestCheckEndpoint(t *testing.T) {
	from := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cases := []struct {
		name  string
		node  fakeNode
		noNet bool
		want  []string // substrings of the refusal; none for no refusal
	}{
		{name: "connected, nothing pending", node: fakeNode{peers: 3, latest: 5, pending: 5}},
		{name: "no peers", node: fakeNode{peers: 0, latest: 5, pending: 5}, want: []string{"has no peers"}},
		{name: "no peers, net API not served", node: fakeNode{peers: 0, latest: 5, pending: 5}, noNet: true},
		{name: "pending transaction", node: fakeNode{peers: 3, latest: 5, pending: 7}, want: []string{"already has 2 pending transaction(s)"}},
		{name: "no peers and pending transaction", node: fakeNode{peers: 0, latest: 5, pending: 6}, want: []string{"has no peers", "already has 1 pending transaction(s)"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := rpc.NewServer()
			defer server.Stop()
			require.NoError(t, server.RegisterName("kaia", &tc.node))
			if !tc.noNet {
				require.NoError(t, server.RegisterName("net", &tc.node))
			}
			ec := client.NewClient(rpc.DialInProc(server))
			defer ec.Close()

			latest, pending, refusal, err := abv2CheckEndpoint(ec, "node.ipc", from)
			require.NoError(t, err)
			assert.Equal(t, tc.node.latest, latest)
			assert.Equal(t, tc.node.pending, pending)
			if len(tc.want) == 0 {
				assert.Empty(t, refusal)
			}
			for _, want := range tc.want {
				assert.Contains(t, refusal, want)
			}
		})
	}
}

// The endpoint checks as a valops command applies them: what is sent, with which nonce,
// and what a receipt timeout reports.
func TestValopsEndpointChecks(t *testing.T) {
	cases := []struct {
		name      string
		node      fakeNode
		force     bool
		want      string // substring of the error
		wantNonce []uint64
	}{
		{name: "no peers", node: fakeNode{peers: 0, latest: 5, pending: 5}, want: "has no peers"},
		{name: "no peers, forced", node: fakeNode{peers: 0, latest: 5, pending: 5, lookup: "pool"}, force: true, want: "still pending", wantNonce: []uint64{5}},
		{name: "pending transaction", node: fakeNode{peers: 3, latest: 5, pending: 6}, want: "already has 1 pending transaction(s)"},
		{name: "pending transaction, forced", node: fakeNode{peers: 3, latest: 5, pending: 6, lookup: "pool"}, force: true, want: "behind 1 earlier transaction(s)", wantNonce: []uint64{6}},
		// Another run sent a transaction after the check: the nonce the check saw is kept.
		{name: "nonce is pinned", node: fakeNode{peers: 3, latest: 5, pending: 5, bumpPending: true, lookup: "block"}, want: "included in a block", wantNonce: []uint64{5}},
		{name: "outcome unknown", node: fakeNode{peers: 3, latest: 5, pending: 5}, want: "outcome is unknown", wantNonce: []uint64{5}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := rpc.NewServer()
			defer server.Stop()
			require.NoError(t, server.RegisterName("kaia", &tc.node))
			require.NoError(t, server.RegisterName("net", &tc.node))
			httpServer := httptest.NewServer(server)
			defer httpServer.Close()

			key := "0x" + strings.Repeat("11", 32)
			args := []string{"kcn", "valops", "pause", "--private-key", key, "--endpoint", httpServer.URL, "--timeout", "10ms"}
			if tc.force {
				args = append(args, "--force")
			}
			app := &cli.App{Name: "kcn", Commands: []*cli.Command{ValOpsCommand}}
			err := app.Run(args)

			assert.ErrorContains(t, err, tc.want)
			var nonces []uint64
			for _, tx := range tc.node.sent {
				nonces = append(nonces, tx.Nonce())
			}
			assert.Equal(t, tc.wantNonce, nonces)
			if len(tc.wantNonce) == 0 {
				assert.ErrorContains(t, err, "(no transaction sent)")
			}
		})
	}
}

func TestReadyWarnings(t *testing.T) {
	next := big.NewInt(172800)
	day := 24 * time.Hour
	dur := func(d time.Duration) *time.Duration { return &d }

	cases := []struct {
		name      string
		method    string
		timeLeft  time.Duration
		timeoutIn *time.Duration
		want      []string // substrings, one per warning, in order
	}{
		{name: "candidate, enough time", method: methodReadyCandidate, timeLeft: time.Hour},
		{name: "candidate, late", method: methodReadyCandidate, timeLeft: 4 * time.Minute, want: []string{"one epoch later"}},
		{name: "candidate ignores timeout", method: methodReadyCandidate, timeLeft: time.Hour, timeoutIn: dur(time.Minute)},
		{name: "validator, no timeout", method: methodReadyValidator, timeLeft: time.Hour},
		{name: "validator, timeout after next epoch", method: methodReadyValidator, timeLeft: time.Hour, timeoutIn: dur(2 * time.Hour)},
		{name: "validator, timeout before next epoch", method: methodReadyValidator, timeLeft: time.Hour, timeoutIn: dur(30 * time.Minute), want: []string{"before epoch block 172800"}},
		{name: "validator, timeout on the next epoch block", method: methodReadyValidator, timeLeft: time.Hour, timeoutIn: dur(time.Hour)}, // epoch transition runs first
		{name: "validator, timeout one block before next epoch", method: methodReadyValidator, timeLeft: time.Hour, timeoutIn: dur(time.Hour - time.Second), want: []string{"become Registered"}},
		{name: "validator, late, timeout within the following epoch", method: methodReadyValidator, timeLeft: 4 * time.Minute, timeoutIn: dur(12 * time.Hour), want: []string{"one epoch later", "if the transaction is included after"}},
		{name: "validator, late, timeout after the following epoch", method: methodReadyValidator, timeLeft: 4 * time.Minute, timeoutIn: dur(2 * day), want: []string{"one epoch later"}},
		{name: "validator, late, timeout before next epoch", method: methodReadyValidator, timeLeft: 4 * time.Minute, timeoutIn: dur(time.Minute), want: []string{"one epoch later", "before epoch block 172800"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := readyWarnings(tc.method, next, tc.timeLeft, day, tc.timeoutIn)
			require.Len(t, got, len(tc.want), "%q", got)
			for i, w := range tc.want {
				assert.Contains(t, got[i], w)
			}
		})
	}
}
