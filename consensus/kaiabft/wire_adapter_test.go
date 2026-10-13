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

func TestWireAdapterRejectsPostForkRecovery(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	view := &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(1)}
	block := newEmptyBlock(10, common.Hash{1})
	seal, err := crypto.Sign(crypto.Keccak256(b.committedSealPreimage(view, block.Hash())), key)
	require.NoError(t, err)
	for _, tc := range []struct {
		code    uint64
		payload any
	}{
		{bft.MsgPreprepare, &bft.Preprepare{View: view, Proposal: block}},
		{bft.MsgPrepare, &bft.Prepare{View: view, Digest: block.Hash()}},
		{bft.MsgCommit, &bft.Commit{View: view, Digest: block.Hash(), CommittedSeal: seal}},
		{bft.MsgRoundChange, &bft.RoundChange{View: view}},
	} {
		shared := &bft.Message{PrevHash: block.ParentHash(), Code: tc.code, Msg: mustEncode(t, tc.payload), Address: b.address}
		preimage, err := shared.PayloadNoSigForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		shared.Signature, err = b.sign(preimage)
		require.NoError(t, err)
		raw, err := shared.PayloadForFork(b.isPermissionlessAt)
		require.NoError(t, err)
		var local wireMessage
		require.ErrorIs(t, local.FromPayloadForFork(raw, b.isPermissionlessAt, bft.GetSignatureAddress), errUnsupportedRecovery)
		// Even an unjustified bare round-1 PRE-PREPARE cannot enter the old core.
		require.NoError(t, local.fromNormalized(shared))
		_, err = local.PayloadForFork(b.isPermissionlessAt)
		require.ErrorIs(t, err, errUnsupportedRecovery)
		_, err = local.PayloadNoSigForFork(b.isPermissionlessAt)
		require.ErrorIs(t, err, errUnsupportedRecovery)
	}
}
