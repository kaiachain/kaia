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
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/rlp"
)

// wireMessage is the internal Subject representation, not a wire envelope.
type wireMessage struct {
	Hash          common.Hash
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
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
// uses rc.3's fork-aware codec. Prepared-certificate recovery remains unported.
var errUnsupportedRecovery = errors.New("kaiabft prepared-certificate recovery is not implemented")

func (m *wireMessage) normalized() (*bft.Message, error) {
	n := &bft.Message{PrevHash: m.Hash, Code: m.Code, Msg: m.Msg, Address: m.Address, Signature: m.Signature}
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
		value = &bft.RoundChange{View: s.View}
	default:
		return nil, bft.ErrInvalidMessage
	}
	var err error
	n.Msg, err = bft.Encode(value)
	return n, err
}
func (m *wireMessage) fromNormalized(n *bft.Message) error {
	*m = wireMessage{Hash: n.PrevHash, Code: n.Code, Msg: n.Msg, Address: n.Address, Signature: n.Signature}
	var s bft.Subject
	s.PrevHash = n.PrevHash
	switch n.Code {
	case bft.MsgPreprepare:
		var pp bft.Preprepare
		if err := n.Decode(&pp); err != nil {
			return err
		}
		if len(pp.RoundChangeCertificate) > 0 || len(pp.PreparedMessages) > 0 {
			return errUnsupportedRecovery
		}
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
		if r.Prepared != nil || n.Evidence != nil {
			return errUnsupportedRecovery
		}
		s.View = r.View
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
	// The state machine does not implement rc.3 certificate justification.
	// Reject every post-fork recovery round, including unjustified proposals.
	view, err := n.GetView()
	if err != nil {
		return err
	}
	if fork(view.Sequence.Uint64()) && view.Round.Sign() > 0 {
		return errUnsupportedRecovery
	}
	return m.fromNormalized(&n)
}

func (m *wireMessage) PayloadForFork(fork func(uint64) bool) ([]byte, error) {
	n, err := m.normalized()
	if err != nil {
		return nil, err
	}
	view, err := n.GetView()
	if err != nil {
		return nil, err
	}
	if fork(view.Sequence.Uint64()) && view.Round.Sign() > 0 {
		return nil, errUnsupportedRecovery
	}
	return n.PayloadForFork(fork)
}
func (m *wireMessage) PayloadNoSigForFork(fork func(uint64) bool) ([]byte, error) {
	n, err := m.normalized()
	if err != nil {
		return nil, err
	}
	view, err := n.GetView()
	if err != nil {
		return nil, err
	}
	if fork(view.Sequence.Uint64()) && view.Round.Sign() > 0 {
		return nil, errUnsupportedRecovery
	}
	return n.PayloadNoSigForFork(fork)
}
