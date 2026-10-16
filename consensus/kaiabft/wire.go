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
	"bytes"
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/rlp"
)

// wireMessage is the internal Subject representation, not a wire envelope.
// Prepared and Evidence carry what a post-Permissionless ROUND CHANGE adds to
// the Subject: the signed prepared claim and the unsigned certificate proving it.
type wireMessage struct {
	Hash          common.Hash
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
	Prepared      *bft.PreparedClaim
	Evidence      []byte
}

// Decode unmarshals m.Msg into val.
func (m *wireMessage) Decode(val any) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *wireMessage) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v}", m.Code, m.Address.String())
}

// GetView extracts the view from a wireMessage by decoding the inner payload.
func (m *wireMessage) GetView() (*bft.View, error) {
	var msgView *bft.View
	switch m.Code {
	case bft.MsgPreprepare:
		var preprepare *bft.Preprepare
		if decodeErr := m.Decode(&preprepare); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = preprepare.View
	case bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange:
		var subject *bft.Subject
		if decodeErr := m.Decode(&subject); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = subject.View
	default:
		return nil, bft.ErrInvalidMessage
	}
	return msgView, nil
}

// The engine keeps its Subject representation internally. Network encoding
// uses rc.3's fork-aware codec. The signed payload is re-derived from this
// representation whenever a message is relayed or embedded in a certificate,
// so ingress admits only messages that the conversion reproduces byte for byte.
var errNonCanonicalMessage = errors.New("kaiabft: message is not in canonical form")

func (m *wireMessage) normalized() (*bft.Message, error) {
	n := &bft.Message{PrevHash: m.Hash, Code: m.Code, Msg: m.Msg, Address: m.Address, Signature: m.Signature, Evidence: m.Evidence}
	if m.Code == bft.MsgPreprepare {
		return n, nil
	}
	var s bft.Subject
	if err := m.Decode(&s); err != nil {
		return nil, err
	}
	var value any
	switch m.Code {
	case bft.MsgPrepare:
		value = &bft.Prepare{View: s.View, Digest: s.Digest}
	case bft.MsgCommit:
		value = &bft.Commit{View: s.View, Digest: s.Digest, CommittedSeal: m.CommittedSeal}
	case bft.MsgRoundChange:
		value = &bft.RoundChange{View: s.View, Prepared: m.Prepared}
	default:
		return nil, bft.ErrInvalidMessage
	}
	var err error
	n.Msg, err = bft.Encode(value)
	return n, err
}

// withoutEvidence returns a shallow copy without the unsigned Evidence.
func (m *wireMessage) withoutEvidence() *wireMessage {
	stripped := *m
	stripped.Evidence = nil
	return &stripped
}

func (m *wireMessage) fromNormalized(n *bft.Message) error {
	*m = wireMessage{Hash: n.PrevHash, Code: n.Code, Msg: n.Msg, Address: n.Address, Signature: n.Signature}
	var s bft.Subject
	s.PrevHash = n.PrevHash
	switch n.Code {
	case bft.MsgPreprepare:
		// The payload, justification included, is kept as signed; the
		// handler decodes it.
		return nil
	case bft.MsgPrepare:
		var p bft.Prepare
		if err := n.Decode(&p); err != nil {
			return err
		}
		s.View, s.Digest = p.View, p.Digest
	case bft.MsgCommit:
		var c bft.Commit
		if err := n.Decode(&c); err != nil {
			return err
		}
		s.View, s.Digest, m.CommittedSeal = c.View, c.Digest, c.CommittedSeal
	case bft.MsgRoundChange:
		var r bft.RoundChange
		if err := n.Decode(&r); err != nil {
			return err
		}
		s.View = r.View
		m.Prepared, m.Evidence = r.Prepared, n.Evidence
	default:
		return bft.ErrInvalidMessage
	}
	var err error
	m.Msg, err = bft.Encode(&s)
	return err
}

func (m *wireMessage) FromPayloadForFork(data []byte, fork func(uint64) bool, validate func([]byte, []byte) (common.Address, error)) error {
	var n bft.Message
	if err := n.FromPayloadForFork(data, fork, validate); err != nil {
		return err
	}
	if err := m.fromNormalized(&n); err != nil {
		return err
	}
	if n.Code == bft.MsgPreprepare {
		return nil
	}
	// rlp accepts a few alternative encodings of one value, such as an empty
	// list for an absent prepared claim. A message in such a form would be
	// relayed, or embedded in a certificate, with bytes its sender never
	// signed, so it is refused here.
	back, err := m.normalized()
	if err != nil {
		return err
	}
	if !bytes.Equal(back.Msg, n.Msg) {
		return errNonCanonicalMessage
	}
	return nil
}

func (m *wireMessage) PayloadForFork(fork func(uint64) bool) ([]byte, error) {
	n, err := m.normalized()
	if err != nil {
		return nil, err
	}
	return n.PayloadForFork(fork)
}

func (m *wireMessage) PayloadNoSigForFork(fork func(uint64) bool) ([]byte, error) {
	n, err := m.normalized()
	if err != nil {
		return nil, err
	}
	return n.PayloadNoSigForFork(fork)
}
