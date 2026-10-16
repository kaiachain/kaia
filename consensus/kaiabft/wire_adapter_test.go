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
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

// Verify both boundaries against the shared rc.3 codec, including signature
// recovery; a private encoder/decoder round trip alone cannot establish this.
func TestWireAdapterRoundZeroSharedCodec(t *testing.T) {
	for _, cfg := range []*params.ChainConfig{legacyConfig(), permissionlessConfig()} {
		b, key := newTestBackend(t, &fakeChain{cfg: cfg})
		block := newEmptyBlock(10, common.Hash{1})
		m := newTestMachine(b, 10, 0, block, []common.Address{b.address})
		for _, code := range []uint64{bft.MsgPreprepare, bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange} {
			view := &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(0)}
			subject := &bft.Subject{View: view, Digest: block.Hash(), PrevHash: block.ParentHash()}
			var inner any = subject
			if code == bft.MsgPreprepare {
				inner = &bft.Preprepare{View: view, Proposal: block}
			}
			if code == bft.MsgRoundChange {
				subject.Digest = common.Hash{}
			}
			msg := &wireMessage{Hash: block.ParentHash(), Code: code, Msg: mustEncode(t, inner)}
			raw := m.signPayload(msg)
			require.NotNil(t, raw, "code %d", code)
			var shared bft.Message
			require.NoError(t, shared.FromPayloadForFork(raw, b.isPermissionlessAt, bft.GetSignatureAddress))
			require.Equal(t, crypto.PubkeyToAddress(key.PublicKey), shared.Address)
			require.Equal(t, block.ParentHash(), shared.PrevHash)
			encoded, err := shared.PayloadForFork(b.isPermissionlessAt)
			require.NoError(t, err)
			require.Equal(t, raw, encoded)
			var decoded wireMessage
			require.NoError(t, decoded.FromPayloadForFork(encoded, b.isPermissionlessAt, bft.GetSignatureAddress))
			if len(msg.CommittedSeal) == 0 {
				msg.CommittedSeal = nil
			}
			require.Equal(t, msg, &decoded)
		}
	}
}

// Post-fork recovery messages, a prepared claim with its Evidence and a
// justified PRE-PREPARE included, pass the adapter losslessly: the payload and
// signed preimage the machine re-derives are the ones the sender signed, so a
// relayed or embedded message stays verifiable. Whether such a message is
// justified is decided by the machine (recovery_test.go), not the adapter.
func TestWireAdapterPostForkRecoveryRoundTrip(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	view := &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(1)}
	block := newEmptyBlock(10, common.Hash{1})
	seal, err := crypto.Sign(crypto.Keccak256(b.committedSealPreimage(view, block.Hash())), key)
	require.NoError(t, err)
	evidence := []byte{0xc1, 0x80}
	claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: block.Hash(), EvidenceHash: crypto.Keccak256Hash(evidence)}
	embedded := &bft.Message{PrevHash: block.ParentHash(), Code: bft.MsgRoundChange, Msg: mustEncode(t, &bft.RoundChange{View: view, Prepared: claim}), Address: b.address, Signature: make([]byte, 65)}
	for _, tc := range []struct {
		code     uint64
		payload  any
		evidence []byte
	}{
		{bft.MsgPreprepare, &bft.Preprepare{View: view, Proposal: block}, nil},
		{bft.MsgPreprepare, &bft.Preprepare{View: view, Proposal: block, RoundChangeCertificate: []*bft.Message{embedded}, PreparedMessages: []*bft.Message{embedded}}, nil},
		{bft.MsgPrepare, &bft.Prepare{View: view, Digest: block.Hash()}, nil},
		{bft.MsgCommit, &bft.Commit{View: view, Digest: block.Hash(), CommittedSeal: seal}, nil},
		{bft.MsgRoundChange, &bft.RoundChange{View: view}, nil},
		{bft.MsgRoundChange, &bft.RoundChange{View: view, Prepared: claim}, evidence},
	} {
		shared := &bft.Message{PrevHash: block.ParentHash(), Code: tc.code, Msg: mustEncode(t, tc.payload), Address: b.address, Evidence: tc.evidence}
		preimage, err := shared.PayloadNoSigForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		shared.Signature, err = b.sign(preimage)
		require.NoError(t, err)
		raw, err := shared.PayloadForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		var local wireMessage
		require.NoError(t, local.FromPayloadForFork(raw, b.isPermissionlessAt, bft.GetSignatureAddress))
		encoded, err := local.PayloadForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		require.Equal(t, raw, encoded, "code %d", tc.code)
		unsigned, err := local.PayloadNoSigForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		require.Equal(t, preimage, unsigned, "code %d", tc.code)
		stripped, err := local.withoutEvidence().normalized()
		require.NoError(t, err)
		require.Equal(t, shared.WithoutEvidence(), stripped, "code %d", tc.code)
	}
}
