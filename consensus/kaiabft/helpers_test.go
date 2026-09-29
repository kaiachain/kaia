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

package kaiabft

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// types.NewBlock derives the tx root through the process-global DeriveSha.
	blockchain.InitDeriveSha(params.TestChainConfig)
	os.Exit(m.Run())
}

// fakeChain serves the consensus.ChainReader calls these tests reach. Any other
// call panics through the nil embedded interface, which makes an unexpected
// dependency visible.
type fakeChain struct {
	consensus.ChainReader
	cfg     *params.ChainConfig
	headers map[common.Hash]*types.Header
	current *types.Block
}

func (c *fakeChain) Config() *params.ChainConfig { return c.cfg }

func (c *fakeChain) HasBadBlock(common.Hash) bool { return false }

func (c *fakeChain) ValidateHeader(*types.Header) error { return nil }

func (c *fakeChain) ValidateProposalHeader(*types.Header) error { return nil }

func (c *fakeChain) GetHeader(hash common.Hash, _ uint64) *types.Header { return c.headers[hash] }

func (c *fakeChain) CurrentBlock() *types.Block { return c.current }

func (c *fakeChain) CurrentHeader() *types.Header { return c.current.Header() }

func permissionlessConfig() *params.ChainConfig {
	return &params.ChainConfig{PermissionlessCompatibleBlock: big.NewInt(0)}
}

func legacyConfig() *params.ChainConfig {
	return &params.ChainConfig{}
}

func newTestBackend(t *testing.T, chain consensus.ChainReader) (*backend, *ecdsa.PrivateKey) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	b := New(&Opts{
		Timeout:    1000,
		PrivateKey: key,
		NodeType:   common.CONSENSUSNODE,
		Sealer:     istanbul.NewSealerImpl(key),
	}).(*backend)
	b.chain = chain
	return b, key
}

func newEmptyBlock(number int64, parent common.Hash) *types.Block {
	return types.NewBlock(&types.Header{Number: big.NewInt(number), ParentHash: parent, Time: big.NewInt(0)}, nil, nil)
}

// newTestMachine returns a machine in the Preprepared state at (seq, round)
// with block accepted. Thresholds are out of reach so handlers never move the
// state during a test.
func newTestMachine(b *backend, seq, round int64, block *types.Block, committee []common.Address) *machine {
	m := newMachine(b)
	m.sequence = big.NewInt(seq)
	m.round = big.NewInt(round)
	m.state = statePreprepared
	if block != nil {
		m.preprepare = &bft.Preprepare{View: m.currentView(), Proposal: block}
	}
	m.qualified = valset.NewAddressSet(committee)
	m.committee = valset.NewAddressSet(committee)
	m.prepares = newMessageSet(m.qualified)
	m.commits = newMessageSet(m.qualified)
	m.requiredMessageCount = len(committee) + 1
	return m
}

func mustEncode(t *testing.T, v any) []byte {
	enc, err := bft.Encode(v)
	require.NoError(t, err)
	return enc
}

func sealOver(t *testing.T, key *ecdsa.PrivateKey, preimage []byte) []byte {
	sig, err := crypto.Sign(crypto.Keccak256(preimage), key)
	require.NoError(t, err)
	return sig
}

func recovers(preimage, seal []byte, want common.Address) bool {
	got, err := bft.GetSignatureAddress(preimage, seal)
	return err == nil && got == want
}

// captureMessages drains the backend's self-delivered messages so a test can
// see what the machine broadcast.
func captureMessages(b *backend) (<-chan messageEvent, func()) {
	sub := b.eventMux.Subscribe(messageEvent{})
	out := make(chan messageEvent, 16)
	go func() {
		for ev := range sub.Chan() {
			if me, ok := ev.Data.(messageEvent); ok {
				out <- me
			}
		}
	}()
	return out, sub.Unsubscribe
}

func expectNoMessage(t *testing.T, msgs <-chan messageEvent) {
	select {
	case ev := <-msgs:
		t.Fatalf("unexpected broadcast: %x", ev.Payload)
	case <-time.After(200 * time.Millisecond):
	}
}

func recvMessage(t *testing.T, msgs <-chan messageEvent) *bft.Message {
	select {
	case ev := <-msgs:
		var msg bft.Message
		require.NoError(t, msg.FromPayload(ev.Payload, nil))
		return &msg
	case <-time.After(2 * time.Second):
		t.Fatal("no broadcast within 2s")
		return nil
	}
}

func subjectMsg(t *testing.T, code uint64, src common.Address, seq, round int64) *bft.Message {
	sub := &bft.Subject{View: &bft.View{Sequence: big.NewInt(seq), Round: big.NewInt(round)}}
	return &bft.Message{Code: code, Msg: mustEncode(t, sub), Address: src}
}
