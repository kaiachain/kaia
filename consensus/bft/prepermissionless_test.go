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

package bft_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
)

// TestSubjectRoundTrip verifies Subject encode→decode preserves all fields.
func TestSubjectRoundTrip(t *testing.T) {
	s := &bft.Subject{
		View:     &bft.View{Round: big.NewInt(1), Sequence: big.NewInt(2)},
		Digest:   common.HexToHash("0xdeadbeef"),
		PrevHash: common.HexToHash("0xfeedface"),
	}
	b, err := rlp.EncodeToBytes(s)
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.Subject
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.View.Cmp(s.View) != 0 || decoded.Digest != s.Digest || decoded.PrevHash != s.PrevHash {
		t.Fatalf("round-trip mismatch: got %v want %v", &decoded, s)
	}
}

// legacyMessage builds the exact pre-Permissionless wire encoding of a message
// whose payload is a Subject, signed by key.
func legacyMessage(t *testing.T, key *ecdsa.PrivateKey, code uint64, subject *bft.Subject, hash common.Hash, committedSeal []byte) []byte {
	t.Helper()
	legacy := &bft.PrePermissionlessMessage{
		Hash:          hash,
		Code:          code,
		Msg:           mustEncode(t, subject),
		Address:       crypto.PubkeyToAddress(key.PublicKey),
		CommittedSeal: committedSeal,
	}
	preimage := mustEncode(t, &bft.PrePermissionlessMessage{
		Hash: legacy.Hash, Code: legacy.Code, Msg: legacy.Msg, Address: legacy.Address, CommittedSeal: legacy.CommittedSeal,
	})
	var err error
	if legacy.Signature, err = crypto.Sign(crypto.Keccak256(preimage), key); err != nil {
		t.Fatal(err)
	}
	return mustEncode(t, legacy)
}

func neverPermissionless(uint64) bool  { return false }
func alwaysPermissionless(uint64) bool { return true }

// TestPrePermissionlessMessageNormalization verifies that a legacy message is
// authenticated over its legacy preimage, exposed with the post-Permissionless
// payload types, and re-encoded to the exact signed bytes so that it can be
// relayed to legacy peers.
func TestPrePermissionlessMessageNormalization(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	view := testView(7)
	prevHash := common.HexToHash("0xfeed")
	digest := common.HexToHash("0xd1")
	seal := bytes.Repeat([]byte{0x5a}, crypto.SignatureLength)

	t.Run("prepare", func(t *testing.T) {
		payload := legacyMessage(t, key, bft.MsgPrepare, &bft.Subject{View: view, Digest: digest, PrevHash: prevHash}, prevHash, []byte{})
		msg := decodeLegacy(t, payload)
		var prepare *bft.Prepare
		if err := msg.Decode(&prepare); err != nil {
			t.Fatal(err)
		}
		if prepare.View.Cmp(view) != 0 || prepare.Digest != digest || msg.PrevHash != prevHash {
			t.Fatalf("unexpected prepare: %+v prevHash %x", prepare, msg.PrevHash)
		}
		expectPayload(t, msg, neverPermissionless, payload)
	})
	t.Run("commit", func(t *testing.T) {
		payload := legacyMessage(t, key, bft.MsgCommit, &bft.Subject{View: view, Digest: digest, PrevHash: prevHash}, prevHash, seal)
		msg := decodeLegacy(t, payload)
		var commit *bft.Commit
		if err := msg.Decode(&commit); err != nil {
			t.Fatal(err)
		}
		if commit.View.Cmp(view) != 0 || commit.Digest != digest || !bytes.Equal(commit.CommittedSeal, seal) {
			t.Fatalf("unexpected commit: %+v", commit)
		}
		expectPayload(t, msg, neverPermissionless, payload)
	})
	t.Run("round change", func(t *testing.T) {
		payload := legacyMessage(t, key, bft.MsgRoundChange, &bft.Subject{View: view, PrevHash: prevHash}, prevHash, []byte{})
		msg := decodeLegacy(t, payload)
		var rc *bft.RoundChange
		if err := msg.Decode(&rc); err != nil {
			t.Fatal(err)
		}
		if rc.View.Cmp(view) != 0 || rc.Prepared != nil || msg.PrevHash != prevHash {
			t.Fatalf("unexpected round change: %+v", rc)
		}
		expectPayload(t, msg, neverPermissionless, payload)
	})
	t.Run("preprepare", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(7), ParentHash: prevHash})
		legacy := &bft.PrePermissionlessMessage{
			Hash:          prevHash,
			Code:          bft.MsgPreprepare,
			Msg:           mustEncode(t, &bft.Preprepare{View: view, Proposal: block}),
			Address:       crypto.PubkeyToAddress(key.PublicKey),
			CommittedSeal: []byte{},
		}
		preimage := mustEncode(t, &bft.PrePermissionlessMessage{Hash: legacy.Hash, Code: legacy.Code, Msg: legacy.Msg, Address: legacy.Address, CommittedSeal: []byte{}})
		if legacy.Signature, err = crypto.Sign(crypto.Keccak256(preimage), key); err != nil {
			t.Fatal(err)
		}
		payload := mustEncode(t, legacy)
		msg := decodeLegacy(t, payload)
		expectPayload(t, msg, neverPermissionless, payload)
	})
}

// TestPrePermissionlessGoldenPayloads pins messages produced by the deployed
// pre-Permissionless codec (core.message on kaiachain/kaia main at ac7c81f3a,
// signed with a fixed key). The other legacy tests build their bytes with
// PrePermissionlessMessage itself, so they cannot detect a change that both
// encodes and decodes differently from deployed nodes. Each payload must
// authenticate over the normalized signed preimage and re-encode to exactly
// the bytes a deployed node signed.
func TestPrePermissionlessGoldenPayloads(t *testing.T) {
	const (
		preprepare  = "f9022fa0000000000000000000000000000000000000000000000000000000000000feed80b901b1f901aec20107f901a8f901a4a0000000000000000000000000000000000000000000000000000000000000feed940000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000107800180808080c09471562b71999873db5b286df957af199ec94617f7b8414c04bda7733403a6fe6b3462921bc91edbc6ae1d45208c53b28121a88845abd112abe2fc9c7a078b1d160d127324c7af634e4aaf159ce1fdf8119354b04451300180"
		prepare     = "f8c4a0000000000000000000000000000000000000000000000000000000000000feed01b847f845c20107a000000000000000000000000000000000000000000000000000000000000000d1a0000000000000000000000000000000000000000000000000000000000000feed9471562b71999873db5b286df957af199ec94617f7b841596a4d3be092b9b9de9917540a30bebf5302ec9db92b7aa5a4b6cf94e98e433944fc265a31711bbb8928256fb53bbbf6b6bc2dc9d22679ba0081e40241182ffa0180"
		commit      = "f90106a0000000000000000000000000000000000000000000000000000000000000feed02b847f845c20107a000000000000000000000000000000000000000000000000000000000000000d1a0000000000000000000000000000000000000000000000000000000000000feed9471562b71999873db5b286df957af199ec94617f7b8412b36339718c8667bd9cf6dfd4b5ac550a66dcdf46a02243504c5b788d0daa0752e1033a7d414192ebf4adb62b6f6af8f7b61ac856f5ce983e8dde94da59e5cfc01b841b9b2212c70b25c13ce89937ef0ac58245c6bb40811360ff5bee94e305c05fb6a2ad624e709e4ab8c182a910cf1d529245254be9ffa0963b14bf4aa3ad57ad8ab00"
		roundChange = "f8c4a0000000000000000000000000000000000000000000000000000000000000feed03b847f845c20107a00000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000feed9471562b71999873db5b286df957af199ec94617f7b8414b4af5d1f14ac781ee0bceb81c576a6ca8ba9f89ce61c7941d63d14bd04a054c2c487b6ce542d2cbf4c555bd1f63aea8fb6439f5e5126b8fa0b194ffe2d16bb10180"
	)
	signer := common.HexToAddress("0x71562b71999873DB5b286dF957af199Ec94617F7")
	view := &bft.View{Round: big.NewInt(1), Sequence: big.NewInt(7)}
	prevHash := common.HexToHash("0xfeed")
	digest := common.HexToHash("0xd1")

	for _, tc := range []struct {
		name    string
		payload string
		check   func(t *testing.T, msg *bft.Message)
	}{
		{"preprepare", preprepare, func(t *testing.T, msg *bft.Message) {
			var pp *bft.Preprepare
			if err := msg.Decode(&pp); err != nil || pp.View.Cmp(view) != 0 ||
				pp.Proposal.ParentHash() != prevHash || pp.Proposal.Number().Cmp(view.Sequence) != 0 {
				t.Fatalf("unexpected preprepare %+v: %v", pp, err)
			}
		}},
		{"prepare", prepare, func(t *testing.T, msg *bft.Message) {
			var p *bft.Prepare
			if err := msg.Decode(&p); err != nil || p.View.Cmp(view) != 0 || p.Digest != digest {
				t.Fatalf("unexpected prepare %+v: %v", p, err)
			}
		}},
		{"commit", commit, func(t *testing.T, msg *bft.Message) {
			var c *bft.Commit
			if err := msg.Decode(&c); err != nil || c.View.Cmp(view) != 0 || c.Digest != digest || len(c.CommittedSeal) != crypto.SignatureLength {
				t.Fatalf("unexpected commit %+v: %v", c, err)
			}
		}},
		{"round change", roundChange, func(t *testing.T, msg *bft.Message) {
			var rc *bft.RoundChange
			if err := msg.Decode(&rc); err != nil || rc.View.Cmp(view) != 0 || rc.Prepared != nil {
				t.Fatalf("unexpected round change %+v: %v", rc, err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := hex.DecodeString(tc.payload)
			if err != nil {
				t.Fatal(err)
			}
			msg := decodeLegacy(t, payload)
			if msg.Address != signer || msg.PrevHash != prevHash || len(msg.Evidence) != 0 {
				t.Fatalf("unexpected envelope: address %s prevHash %x evidence %d", msg.Address, msg.PrevHash, len(msg.Evidence))
			}
			tc.check(t, msg)
			preimage, err := msg.PayloadNoSigForFork(neverPermissionless)
			if err != nil {
				t.Fatal(err)
			}
			if recovered, err := recoverSigner(preimage, msg.Signature); err != nil || recovered != signer {
				t.Fatalf("signed preimage drift: recovered %s, err %v", recovered, err)
			}
			expectPayload(t, msg, neverPermissionless, payload)
		})
	}
}

// TestPrePermissionlessMessageRejectsLossyShapes verifies that legacy
// messages which the normalized representation cannot reproduce exactly are
// rejected. Honest legacy senders never produce them.
func TestPrePermissionlessMessageRejectsLossyShapes(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	view := testView(7)
	prevHash := common.HexToHash("0xfeed")
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{"prepare with mismatched envelope hash", legacyMessage(t, key, bft.MsgPrepare, &bft.Subject{View: view, PrevHash: prevHash}, common.HexToHash("0x01"), []byte{})},
		{"round change with digest", legacyMessage(t, key, bft.MsgRoundChange, &bft.Subject{View: view, Digest: common.HexToHash("0x01"), PrevHash: prevHash}, prevHash, []byte{})},
		{"prepare with committed seal", legacyMessage(t, key, bft.MsgPrepare, &bft.Subject{View: view, PrevHash: prevHash}, prevHash, []byte{1})},
		{"commit with short committed seal", legacyMessage(t, key, bft.MsgCommit, &bft.Subject{View: view, PrevHash: prevHash}, prevHash, []byte{1})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var msg bft.Message
			if err := msg.FromPayloadForFork(tc.payload, neverPermissionless, recoverSigner); !errors.Is(err, bft.ErrInvalidMessage) {
				t.Fatalf("got %v, want ErrInvalidMessage", err)
			}
		})
	}
}

// TestMessageFromPayloadForForkSelectsCodec verifies the codec is chosen by the
// sequence of the carried view, and that a message in the other format fails.
func TestMessageFromPayloadForForkSelectsCodec(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	view := testView(7)
	prevHash := common.HexToHash("0xfeed")
	legacy := legacyMessage(t, key, bft.MsgRoundChange, &bft.Subject{View: view, PrevHash: prevHash}, prevHash, []byte{})

	current := &bft.Message{PrevHash: prevHash, Code: bft.MsgRoundChange, Msg: mustEncode(t, &bft.RoundChange{View: view}), Address: crypto.PubkeyToAddress(key.PublicKey)}
	sign(t, current, key)
	currentPayload, err := current.Payload()
	if err != nil {
		t.Fatal(err)
	}

	var sequence uint64
	forkAt := func(seq uint64) bool { sequence = seq; return true }
	var decoded bft.Message
	if err := decoded.FromPayloadForFork(currentPayload, forkAt, recoverSigner); err != nil {
		t.Fatal(err)
	}
	if sequence != view.Sequence.Uint64() {
		t.Fatalf("fork queried at sequence %d, want %d", sequence, view.Sequence.Uint64())
	}
	expectPayload(t, &decoded, alwaysPermissionless, currentPayload)

	if err := decoded.FromPayloadForFork(legacy, alwaysPermissionless, recoverSigner); err == nil {
		t.Fatal("legacy payload accepted after the fork")
	}
	if err := decoded.FromPayloadForFork(currentPayload, neverPermissionless, recoverSigner); err == nil {
		t.Fatal("post-Permissionless payload accepted before the fork")
	}
}

// TestMessagePayloadForForkSelectsCodec verifies that the encoder picks the
// wire format from the message sequence, and that each encoding decodes back
// with the same fork rule.
func TestMessagePayloadForForkSelectsCodec(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	view := testView(7)
	msg := &bft.Message{
		PrevHash: common.HexToHash("0xfeed"),
		Code:     bft.MsgCommit,
		Msg:      mustEncode(t, &bft.Commit{View: view, Digest: common.HexToHash("0xd1"), CommittedSeal: make([]byte, crypto.SignatureLength)}),
		Address:  crypto.PubkeyToAddress(key.PublicKey),
	}
	for _, tc := range []struct {
		name               string
		isPermissionlessAt func(uint64) bool
		legacy             bool
	}{
		{"fork at the message sequence", func(n uint64) bool { return n >= 7 }, false},
		{"fork after the message sequence", func(n uint64) bool { return n >= 8 }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			signed := *msg
			preimage, err := signed.PayloadNoSigForFork(tc.isPermissionlessAt)
			if err != nil {
				t.Fatal(err)
			}
			if signed.Signature, err = crypto.Sign(crypto.Keccak256(preimage), key); err != nil {
				t.Fatal(err)
			}
			payload, err := signed.PayloadForFork(tc.isPermissionlessAt)
			if err != nil {
				t.Fatal(err)
			}
			var legacy bft.PrePermissionlessMessage
			if isLegacy := rlp.DecodeBytes(payload, &legacy) == nil; isLegacy != tc.legacy {
				t.Fatalf("legacy encoding = %v, want %v", isLegacy, tc.legacy)
			}
			var decoded bft.Message
			if err := decoded.FromPayloadForFork(payload, tc.isPermissionlessAt, recoverSigner); err != nil {
				t.Fatal(err)
			}
			expectPayload(t, &decoded, tc.isPermissionlessAt, payload)
		})
	}
}

// TestPrePermissionlessPayloadRejectsUnrepresentableFields verifies that the
// legacy codec refuses to silently drop fields it cannot carry.
func TestPrePermissionlessPayloadRejectsUnrepresentableFields(t *testing.T) {
	view := testView(1)
	for _, tc := range []struct {
		name string
		msg  *bft.Message
	}{
		{"evidence", &bft.Message{Code: bft.MsgRoundChange, Msg: mustEncode(t, &bft.RoundChange{View: view}), Evidence: []byte{1}}},
		{"prepared claim", &bft.Message{Code: bft.MsgRoundChange, Msg: mustEncode(t, &bft.RoundChange{
			View: view, Prepared: &bft.PreparedClaim{Round: big.NewInt(0)},
		})}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.msg.PayloadForFork(neverPermissionless); !errors.Is(err, bft.ErrInvalidMessage) {
				t.Fatalf("Payload: got %v, want ErrInvalidMessage", err)
			}
			if _, err := tc.msg.PayloadNoSigForFork(neverPermissionless); !errors.Is(err, bft.ErrInvalidMessage) {
				t.Fatalf("PayloadNoSig: got %v, want ErrInvalidMessage", err)
			}
		})
	}
}

func decodeLegacy(t *testing.T, payload []byte) *bft.Message {
	t.Helper()
	var msg bft.Message
	if err := msg.FromPayloadForFork(payload, neverPermissionless, recoverSigner); err != nil {
		t.Fatal(err)
	}
	return &msg
}
