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
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/rlp"
)

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
