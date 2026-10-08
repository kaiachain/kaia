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

// TestOptionalCertificateWireCompatibility ensures an ordinary round-0
// PRE-PREPARE keeps its pre-certificate RLP bytes. The certificate fields
// appear only when the permissionless fork actually supplies evidence.
func TestOptionalCertificateWireCompatibility(t *testing.T) {
	view := &bft.View{Round: big.NewInt(0), Sequence: big.NewInt(1)}
	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})

	legacyPreprepare, err := rlp.EncodeToBytes([]any{view, block})
	if err != nil {
		t.Fatal(err)
	}
	newPreprepare, err := rlp.EncodeToBytes(&bft.Preprepare{View: view, Proposal: block})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(legacyPreprepare, newPreprepare) {
		t.Fatalf("preprepare wire drift: got %x want %x", newPreprepare, legacyPreprepare)
	}
}

// TestViewWireBytes pins the RLP byte output of bft.View{Round:7, Sequence:42}
// to a historical value computed from the pre-move istanbul.View encoding.
// Any change to View's field order, field types, or EncodeRLP implementation
// breaks this test — a loud signal that wire compatibility with already-
// deployed Istanbul nodes is at risk.
func TestViewWireBytes(t *testing.T) {
	v := &bft.View{Round: big.NewInt(7), Sequence: big.NewInt(42)}
	got, err := rlp.EncodeToBytes(v)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	// Expected bytes: a 2-item RLP list [0x07, 0x2a]:
	//   0xc2 — list prefix, 2 bytes follow
	//   0x07 — big.NewInt(7) encoded as a single byte
	//   0x2a — big.NewInt(42) encoded as a single byte
	const want = "c2072a"
	if hex.EncodeToString(got) != want {
		t.Fatalf("wire drift: got %x want %s", got, want)
	}
}

// TestViewRoundTrip verifies encode→decode preserves semantic equality.
func TestViewRoundTrip(t *testing.T) {
	orig := &bft.View{Round: big.NewInt(7), Sequence: big.NewInt(42)}
	b, err := rlp.EncodeToBytes(orig)
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.View
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Cmp(orig) != 0 {
		t.Fatalf("round-trip mismatch: %v vs %v", &decoded, orig)
	}
}

// TestMessageCodes pins msg code integer values to their historical values.
// The original unexported istanbul/core consts were: preprepare=0, prepare=1,
// commit=2, round-change=3, all=4. These MUST remain stable because existing
// deployed nodes exchange these integers on the wire.
func TestMessageCodes(t *testing.T) {
	cases := []struct {
		name string
		got  uint64
		want uint64
	}{
		{"preprepare", bft.MsgPreprepare, 0},
		{"prepare", bft.MsgPrepare, 1},
		{"commit", bft.MsgCommit, 2},
		{"roundchange", bft.MsgRoundChange, 3},
		{"all", bft.MsgAll, 4},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s: got %d want %d", tc.name, tc.got, tc.want)
		}
	}
}

// TestMessageRoundTrip verifies Message encode→decode preserves all fields.
func TestMessageRoundTrip(t *testing.T) {
	orig := &bft.Message{
		PrevHash:  common.HexToHash("0x01"),
		Code:      bft.MsgRoundChange,
		Msg:       []byte{0xaa, 0xbb, 0xcc},
		Address:   common.HexToAddress("0xcafe"),
		Signature: []byte{0x11, 0x22},
		Evidence:  []byte{0x33, 0x44},
	}
	b, err := rlp.EncodeToBytes(orig)
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.Message
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.PrevHash != orig.PrevHash ||
		decoded.Code != orig.Code ||
		!equalBytes(decoded.Msg, orig.Msg) ||
		decoded.Address != orig.Address ||
		!equalBytes(decoded.Signature, orig.Signature) ||
		!equalBytes(decoded.Evidence, orig.Evidence) {
		t.Fatalf("round-trip mismatch: got %+v want %+v", &decoded, orig)
	}
}

// TestMessageEmptyEvidenceIsOmitted pins the envelope of a Message without
// Evidence to five fields, so that ordinary messages carry no attachment slot.
func TestMessageEmptyEvidenceIsOmitted(t *testing.T) {
	base := bft.Message{
		PrevHash:  common.HexToHash("0x01"),
		Code:      bft.MsgRoundChange,
		Msg:       []byte{0xaa, 0xbb, 0xcc},
		Address:   common.HexToAddress("0xcafe"),
		Signature: make([]byte, crypto.SignatureLength),
	}
	want, err := rlp.EncodeToBytes([]any{base.PrevHash, base.Code, base.Msg, base.Address, base.Signature})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name     string
		evidence []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := base
			msg.Evidence = tc.evidence
			payload, err := msg.Payload()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, want) {
				t.Fatalf("envelope wire drift: got %x want %x", payload, want)
			}
		})
	}
}

// TestMessageEvidenceRoundTrip verifies a non-empty Evidence survives
// encode→decode and is appended as a sixth envelope field.
func TestMessageEvidenceRoundTrip(t *testing.T) {
	orig := &bft.Message{
		Code:      bft.MsgRoundChange,
		Msg:       []byte{0xaa},
		Signature: make([]byte, crypto.SignatureLength),
		Evidence:  []byte{0xde, 0xad, 0xbe, 0xef},
	}
	b, err := orig.Payload()
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.Message
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded.Evidence, orig.Evidence) {
		t.Fatalf("evidence mismatch: got %x want %x", decoded.Evidence, orig.Evidence)
	}
	stripped, err := orig.WithoutEvidence().Payload()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) <= len(stripped) {
		t.Fatalf("evidence not encoded: %d <= %d bytes", len(b), len(stripped))
	}
}

// TestMessageEvidenceIsNotSigned verifies Evidence is outside the signed
// preimage: replacing or stripping the attachment of a signed ROUND CHANGE
// keeps the original signature valid, while changing a signed field does not.
func TestMessageEvidenceIsNotSigned(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signed := &bft.Message{
		Code:    bft.MsgRoundChange,
		Msg:     []byte{0xaa, 0xbb},
		Address: crypto.PubkeyToAddress(key.PublicKey),
	}
	sign(t, signed, key)
	preimage, err := signed.PayloadNoSig()
	if err != nil {
		t.Fatal(err)
	}

	attached := *signed
	attached.Evidence = []byte{0x01, 0x02, 0x03}
	attachedPreimage, err := attached.PayloadNoSig()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(attachedPreimage, preimage) {
		t.Fatalf("evidence leaked into the signed preimage: got %x want %x", attachedPreimage, preimage)
	}

	replaced := attached
	replaced.Evidence = []byte{0xff}
	signedField := attached
	signedField.Msg = []byte{0xaa, 0xbc}
	prevHashChanged := attached
	prevHashChanged.PrevHash = common.HexToHash("0x02")

	for _, tc := range []struct {
		name    string
		msg     *bft.Message
		wantErr error
	}{
		{"signed without evidence", signed, nil},
		{"evidence attached after signing", &attached, nil},
		{"evidence replaced", &replaced, nil},
		{"evidence stripped", attached.WithoutEvidence(), nil},
		{"signed field changed", &signedField, bft.ErrInvalidSigner},
		{"prev hash changed", &prevHashChanged, bft.ErrInvalidSigner},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.msg.Payload()
			if err != nil {
				t.Fatal(err)
			}
			var decoded bft.Message
			err = decoded.FromPayload(payload, recoverSigner)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got %v, want %v", err, tc.wantErr)
			}
			if err == nil && !bytes.Equal(decoded.Evidence, tc.msg.Evidence) {
				t.Fatalf("evidence mismatch: got %x want %x", decoded.Evidence, tc.msg.Evidence)
			}
		})
	}
}

// TestMessageFromPayloadRejectsEvidenceOutsideRoundChange verifies only a
// ROUND CHANGE may carry Evidence, and that the envelope is rejected before
// signature recovery.
func TestMessageFromPayloadRejectsEvidenceOutsideRoundChange(t *testing.T) {
	sig := make([]byte, crypto.SignatureLength)
	evidence := []byte{0x01}
	commit := mustEncode(t, &bft.Commit{View: testView(1), CommittedSeal: make([]byte, crypto.SignatureLength)})
	tests := []struct {
		name string
		msg  *bft.Message
	}{
		{"preprepare", &bft.Message{Code: bft.MsgPreprepare, Signature: sig, Evidence: evidence}},
		{"prepare", &bft.Message{Code: bft.MsgPrepare, Signature: sig, Evidence: evidence}},
		{"commit", &bft.Message{Code: bft.MsgCommit, Msg: commit, Signature: sig, Evidence: evidence}},
		{"unknown code", &bft.Message{Code: bft.MsgAll, Signature: sig, Evidence: evidence}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expectRejectedBeforeRecovery(t, tc.msg)
		})
	}
}

// TestMessageFromPayloadValidatesCommitSealBeforeSignatureRecovery verifies
// the committed seal inside a post-Permissionless COMMIT payload is
// shape-checked before signature recovery.
func TestMessageFromPayloadValidatesCommitSealBeforeSignatureRecovery(t *testing.T) {
	sig := make([]byte, crypto.SignatureLength)
	for _, tc := range []struct {
		name string
		seal []byte
	}{
		{"missing", nil},
		{"short", []byte{1}},
		{"oversized", make([]byte, crypto.SignatureLength+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := &bft.Message{
				Code:      bft.MsgCommit,
				Msg:       mustEncode(t, &bft.Commit{View: testView(1), CommittedSeal: tc.seal}),
				Signature: sig,
			}
			expectRejectedBeforeRecovery(t, msg)
		})
	}
	t.Run("undecodable payload", func(t *testing.T) {
		expectRejectedBeforeRecovery(t, &bft.Message{Code: bft.MsgCommit, Signature: sig})
	})
}

func TestMessageFromPayloadAcceptsValidEnvelopeShapes(t *testing.T) {
	sig := make([]byte, crypto.SignatureLength)
	tests := []*bft.Message{
		{Code: bft.MsgPreprepare, Signature: sig},
		{Code: bft.MsgPrepare, Signature: sig},
		{Code: bft.MsgCommit, Signature: sig, Msg: mustEncode(t, &bft.Commit{View: testView(1), CommittedSeal: make([]byte, crypto.SignatureLength)})},
		{Code: bft.MsgRoundChange, Signature: sig},
		{Code: bft.MsgRoundChange, Signature: sig, Evidence: []byte{0x01}},
	}
	for _, msg := range tests {
		payload, err := msg.Payload()
		if err != nil {
			t.Fatal(err)
		}
		var decoded bft.Message
		if err := decoded.FromPayload(payload, func([]byte, []byte) (common.Address, error) {
			return msg.Address, nil
		}); err != nil {
			t.Fatalf("message code %d: %v", msg.Code, err)
		}
	}
}

func expectPayload(t *testing.T, msg *bft.Message, isPermissionlessAt func(uint64) bool, want []byte) {
	t.Helper()
	got, err := msg.PayloadForFork(isPermissionlessAt)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("payload drift: got %x want %x", got, want)
	}
}

func expectRejectedBeforeRecovery(t *testing.T, msg *bft.Message) {
	t.Helper()
	payload, err := msg.Payload()
	if err != nil {
		t.Fatal(err)
	}
	called := false
	var decoded bft.Message
	err = decoded.FromPayload(payload, func([]byte, []byte) (common.Address, error) {
		called = true
		return msg.Address, nil
	})
	if !errors.Is(err, bft.ErrInvalidMessage) {
		t.Fatalf("got %v, want ErrInvalidMessage", err)
	}
	if called {
		t.Fatal("signature recovery ran for an invalid envelope")
	}
}

func sign(t *testing.T, msg *bft.Message, key *ecdsa.PrivateKey) {
	t.Helper()
	preimage, err := msg.PayloadNoSig()
	if err != nil {
		t.Fatal(err)
	}
	if msg.Signature, err = crypto.Sign(crypto.Keccak256(preimage), key); err != nil {
		t.Fatal(err)
	}
}

func mustEncode(t *testing.T, val any) []byte {
	t.Helper()
	b, err := rlp.EncodeToBytes(val)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func testView(sequence int64) *bft.View {
	return &bft.View{Round: big.NewInt(0), Sequence: big.NewInt(sequence)}
}

// recoverSigner mirrors istanbul.GetSignatureAddress without importing the
// engine package into the wire-format tests.
func recoverSigner(data []byte, sig []byte) (common.Address, error) {
	pubkey, err := crypto.SigToPub(crypto.Keccak256(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
