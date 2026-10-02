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

// TestOptionalCertificateWireCompatibility ensures ordinary round-0 and
// legacy ROUND-CHANGE messages keep their pre-certificate RLP bytes. The new
// fields appear only when the permissionless fork actually supplies evidence.
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

	subject := &bft.Subject{View: view, PrevHash: common.HexToHash("0x01")}
	legacyRoundChange, err := rlp.EncodeToBytes(subject)
	if err != nil {
		t.Fatal(err)
	}
	newRoundChange, err := rlp.EncodeToBytes(&bft.RoundChange{View: view, PrevHash: subject.PrevHash})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(legacyRoundChange, newRoundChange) {
		t.Fatalf("round-change wire drift: got %x want %x", newRoundChange, legacyRoundChange)
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
	if !s.Equal(&decoded) {
		t.Fatalf("round-trip mismatch")
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
		Hash:          common.HexToHash("0x01"),
		Code:          bft.MsgPrepare,
		Msg:           []byte{0xaa, 0xbb, 0xcc},
		Address:       common.HexToAddress("0xcafe"),
		Signature:     []byte{0x11, 0x22},
		CommittedSeal: []byte{0x33, 0x44},
	}
	b, err := rlp.EncodeToBytes(orig)
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.Message
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Hash != orig.Hash ||
		decoded.Code != orig.Code ||
		!equalBytes(decoded.Msg, orig.Msg) ||
		decoded.Address != orig.Address ||
		!equalBytes(decoded.Signature, orig.Signature) ||
		!equalBytes(decoded.CommittedSeal, orig.CommittedSeal) {
		t.Fatalf("round-trip mismatch: got %+v want %+v", &decoded, orig)
	}
}

func TestMessageFromPayloadRejectsUnexpectedCommittedSealLengthBeforeSignatureRecovery(t *testing.T) {
	tests := []struct {
		name string
		msg  *bft.Message
	}{
		{
			name: "prepare with committed seal",
			msg: &bft.Message{
				Code: bft.MsgPrepare, Signature: make([]byte, crypto.SignatureLength),
				CommittedSeal: []byte{1},
			},
		},
		{
			name: "round change with committed seal",
			msg: &bft.Message{
				Code: bft.MsgRoundChange, Signature: make([]byte, crypto.SignatureLength),
				CommittedSeal: []byte{1},
			},
		},
		{
			name: "commit without committed seal",
			msg:  &bft.Message{Code: bft.MsgCommit, Signature: make([]byte, crypto.SignatureLength)},
		},
		{
			name: "commit with oversized committed seal",
			msg: &bft.Message{
				Code: bft.MsgCommit, Signature: make([]byte, crypto.SignatureLength),
				CommittedSeal: make([]byte, crypto.SignatureLength+1),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.msg.Payload()
			if err != nil {
				t.Fatal(err)
			}
			called := false
			var decoded bft.Message
			err = decoded.FromPayload(payload, func([]byte, []byte) (common.Address, error) {
				called = true
				return tc.msg.Address, nil
			})
			if !errors.Is(err, bft.ErrInvalidMessage) {
				t.Fatalf("got %v, want ErrInvalidMessage", err)
			}
			if called {
				t.Fatal("signature recovery ran for an invalid committed seal")
			}
		})
	}
}

func TestMessageFromPayloadAcceptsValidEnvelopeShapes(t *testing.T) {
	tests := []*bft.Message{
		{Code: bft.MsgPreprepare, Signature: make([]byte, crypto.SignatureLength)},
		{Code: bft.MsgPrepare, Signature: make([]byte, crypto.SignatureLength)},
		{Code: bft.MsgCommit, Signature: make([]byte, crypto.SignatureLength), CommittedSeal: make([]byte, crypto.SignatureLength)},
		{Code: bft.MsgRoundChange, Signature: make([]byte, crypto.SignatureLength)},
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

// TestMessageEmptyJustificationKeepsLegacyEncoding pins the envelope bytes of
// a Message without a Justification to the pre-Justification six-field
// encoding, so that every ordinary message stays readable by deployed nodes.
func TestMessageEmptyJustificationKeepsLegacyEncoding(t *testing.T) {
	base := bft.Message{
		Hash:          common.HexToHash("0x01"),
		Code:          bft.MsgRoundChange,
		Msg:           []byte{0xaa, 0xbb, 0xcc},
		Address:       common.HexToAddress("0xcafe"),
		Signature:     make([]byte, crypto.SignatureLength),
		CommittedSeal: []byte{},
	}
	legacy, err := rlp.EncodeToBytes([]any{base.Hash, base.Code, base.Msg, base.Address, base.Signature, base.CommittedSeal})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name          string
		justification []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := base
			msg.Justification = tc.justification
			payload, err := msg.Payload()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, legacy) {
				t.Fatalf("envelope wire drift: got %x want %x", payload, legacy)
			}
		})
	}

	// Legacy bytes decode to an empty Justification and re-encode unchanged.
	var decoded bft.Message
	if err := rlp.DecodeBytes(legacy, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Justification) != 0 {
		t.Fatalf("legacy envelope decoded a justification: %x", decoded.Justification)
	}
	reencoded, err := decoded.Payload()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reencoded, legacy) {
		t.Fatalf("legacy envelope re-encode drift: got %x want %x", reencoded, legacy)
	}
}

// TestMessageJustificationRoundTrip verifies a non-empty Justification survives
// encode→decode and is appended as a seventh envelope field.
func TestMessageJustificationRoundTrip(t *testing.T) {
	orig := &bft.Message{
		Code:          bft.MsgRoundChange,
		Msg:           []byte{0xaa},
		Signature:     make([]byte, crypto.SignatureLength),
		Justification: []byte{0xde, 0xad, 0xbe, 0xef},
	}
	b, err := orig.Payload()
	if err != nil {
		t.Fatal(err)
	}
	var decoded bft.Message
	if err := rlp.DecodeBytes(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded.Justification, orig.Justification) {
		t.Fatalf("justification mismatch: got %x want %x", decoded.Justification, orig.Justification)
	}
	stripped, err := orig.WithoutJustification().Payload()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) <= len(stripped) {
		t.Fatalf("justification not encoded: %d <= %d bytes", len(b), len(stripped))
	}
}

// TestMessageJustificationIsNotSigned verifies Justification is outside the
// signed preimage: replacing or stripping the attachment of a signed ROUND
// CHANGE keeps the original signature valid, while changing a signed field
// does not.
func TestMessageJustificationIsNotSigned(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signed := &bft.Message{
		Code:    bft.MsgRoundChange,
		Msg:     []byte{0xaa, 0xbb},
		Address: crypto.PubkeyToAddress(key.PublicKey),
	}
	preimage, err := signed.PayloadNoSig()
	if err != nil {
		t.Fatal(err)
	}
	signed.Signature, err = crypto.Sign(crypto.Keccak256(preimage), key)
	if err != nil {
		t.Fatal(err)
	}

	justified := *signed
	justified.Justification = []byte{0x01, 0x02, 0x03}
	justifiedPreimage, err := justified.PayloadNoSig()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(justifiedPreimage, preimage) {
		t.Fatalf("justification leaked into the signed preimage: got %x want %x", justifiedPreimage, preimage)
	}

	replaced := justified
	replaced.Justification = []byte{0xff}
	signedField := justified
	signedField.Msg = []byte{0xaa, 0xbc}

	for _, tc := range []struct {
		name    string
		msg     *bft.Message
		wantErr error
	}{
		{"signed without justification", signed, nil},
		{"justification attached after signing", &justified, nil},
		{"justification replaced", &replaced, nil},
		{"justification stripped", justified.WithoutJustification(), nil},
		{"signed field changed", &signedField, bft.ErrInvalidSigner},
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
			if err == nil && !bytes.Equal(decoded.Justification, tc.msg.Justification) {
				t.Fatalf("justification mismatch: got %x want %x", decoded.Justification, tc.msg.Justification)
			}
		})
	}
}

// TestMessageFromPayloadRejectsJustificationOutsideRoundChange verifies only a
// ROUND CHANGE may carry a Justification, and that the envelope is rejected
// before signature recovery.
func TestMessageFromPayloadRejectsJustificationOutsideRoundChange(t *testing.T) {
	sig := make([]byte, crypto.SignatureLength)
	justification := []byte{0x01}
	tests := []struct {
		name string
		msg  *bft.Message
	}{
		{"preprepare", &bft.Message{Code: bft.MsgPreprepare, Signature: sig, Justification: justification}},
		{"prepare", &bft.Message{Code: bft.MsgPrepare, Signature: sig, Justification: justification}},
		{"commit", &bft.Message{
			Code: bft.MsgCommit, Signature: sig,
			CommittedSeal: make([]byte, crypto.SignatureLength), Justification: justification,
		}},
		{"unknown code", &bft.Message{Code: bft.MsgAll, Signature: sig, Justification: justification}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.msg.Payload()
			if err != nil {
				t.Fatal(err)
			}
			called := false
			var decoded bft.Message
			err = decoded.FromPayload(payload, func([]byte, []byte) (common.Address, error) {
				called = true
				return tc.msg.Address, nil
			})
			if !errors.Is(err, bft.ErrInvalidMessage) {
				t.Fatalf("got %v, want ErrInvalidMessage", err)
			}
			if called {
				t.Fatal("signature recovery ran for a misplaced justification")
			}
		})
	}
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
