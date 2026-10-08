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

package bft

// This file holds everything that exists only for the pre-Permissionless wire
// format: the legacy envelope and Subject payload, their validation and signed
// preimage, and the conversion to and from the normalized Message. The
// consensus core never sees these types; it calls the *ForFork codec entry
// points, which select the format from the message sequence.
//
// Once Permissionless is active on every network, delete this file and switch
// the *ForFork call sites to FromPayload, Payload and PayloadNoSig.

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
)

// PrePermissionlessMessage is the legacy consensus-message envelope. It is
// decoded and verified only at the wire boundary and immediately normalized to
// Message, so the consensus core has a single representation.
type PrePermissionlessMessage struct {
	Hash          common.Hash
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

func (m *PrePermissionlessMessage) validateEnvelope() error {
	switch m.Code {
	case MsgCommit:
		if len(m.CommittedSeal) != crypto.SignatureLength {
			return fmt.Errorf("%w: committed seal length %d", ErrInvalidMessage, len(m.CommittedSeal))
		}
	case MsgPreprepare, MsgPrepare, MsgRoundChange:
		if len(m.CommittedSeal) != 0 {
			return fmt.Errorf("%w: unexpected committed seal on message code %d", ErrInvalidMessage, m.Code)
		}
	}
	return nil
}

func (m *PrePermissionlessMessage) payloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&PrePermissionlessMessage{
		Hash:          m.Hash,
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// Subject is the legacy payload of PREPARE, COMMIT and ROUND CHANGE messages.
type Subject struct {
	View     *View
	Digest   common.Hash
	PrevHash common.Hash
}

// EncodeRLP serializes a Subject into the Kaia RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{b.View, b.Digest, b.PrevHash})
}

// DecodeRLP deserializes a Subject from a Kaia RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View     *View
		Digest   common.Hash
		PrevHash common.Hash
	}
	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest, b.PrevHash = subject.View, subject.Digest, subject.PrevHash
	return nil
}

func (b *Subject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v, ParentHash: %v}", b.View, b.Digest.String(), b.PrevHash.Hex())
}

// FromPayloadForFork decodes the wire format selected by the message sequence.
// Legacy messages are normalized to the post-Permissionless in-memory types so
// the consensus core has a single representation.
func (m *Message) FromPayloadForFork(b []byte, isPermissionlessAt func(uint64) bool, validateFn func([]byte, []byte) (common.Address, error)) error {
	sequence, err := payloadSequence(b)
	if err != nil {
		return err
	}
	if isPermissionlessAt(sequence) {
		return m.FromPayload(b, validateFn)
	}
	var legacy PrePermissionlessMessage
	if err := rlp.DecodeBytes(b, &legacy); err != nil {
		return err
	}
	if err := legacy.validateEnvelope(); err != nil {
		return err
	}
	// Normalize before authenticating, as FromPayload does, so a caller that
	// gets ErrUnauthorizedAddress can still read the sender and view.
	if err := m.fromPrePermissionless(&legacy); err != nil {
		return err
	}
	if validateFn != nil {
		payload, err := legacy.payloadNoSig()
		if err != nil {
			return err
		}
		signerAddr, err := validateFn(payload, legacy.Signature)
		if err != nil {
			return err
		}
		if !bytes.Equal(signerAddr.Bytes(), legacy.Address.Bytes()) {
			return ErrInvalidSigner
		}
	}
	return nil
}

// PayloadForFork returns the wire encoding selected by the message sequence,
// mirroring FromPayloadForFork.
func (m *Message) PayloadForFork(isPermissionlessAt func(uint64) bool) ([]byte, error) {
	sequence, err := msgSequence(m.Msg)
	if err != nil {
		return nil, err
	}
	if isPermissionlessAt(sequence) {
		return m.Payload()
	}
	legacy, err := m.toPrePermissionless()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(legacy)
}

// PayloadNoSigForFork returns the signed preimage of the wire format selected
// by the message sequence.
func (m *Message) PayloadNoSigForFork(isPermissionlessAt func(uint64) bool) ([]byte, error) {
	sequence, err := msgSequence(m.Msg)
	if err != nil {
		return nil, err
	}
	if isPermissionlessAt(sequence) {
		return m.PayloadNoSig()
	}
	legacy, err := m.toPrePermissionless()
	if err != nil {
		return nil, err
	}
	return legacy.payloadNoSig()
}

func payloadSequence(payload []byte) (uint64, error) {
	var raw struct {
		Hash      common.Hash
		Code      uint64
		Msg       []byte
		Address   common.Address
		Signature []byte
		Tail      []byte `rlp:"optional"`
	}
	if err := rlp.DecodeBytes(payload, &raw); err != nil {
		return 0, err
	}
	return msgSequence(raw.Msg)
}

// msgSequence reads the sequence of the View that every consensus payload,
// legacy or not, carries as its first field, without decoding the rest.
func msgSequence(msg []byte) (uint64, error) {
	stream := rlp.NewStream(bytes.NewReader(msg), uint64(len(msg)))
	if _, err := stream.List(); err != nil {
		return 0, ErrInvalidMessage
	}
	var view View
	if err := stream.Decode(&view); err != nil || view.Sequence == nil || !view.Sequence.IsUint64() {
		return 0, ErrInvalidMessage
	}
	return view.Sequence.Uint64(), nil
}

// fromPrePermissionless normalizes a verified legacy message. The conversion
// must be lossless: a retained legacy message is re-encoded by PayloadForFork before
// it is gossiped, and that encoding has to reproduce the signed bytes. Honest
// legacy senders always repeat the envelope Hash as Subject.PrevHash and leave
// the ROUND CHANGE digest empty, so any other shape is rejected.
func (m *Message) fromPrePermissionless(legacy *PrePermissionlessMessage) error {
	m.PrevHash, m.Code, m.Address, m.Signature = legacy.Hash, legacy.Code, legacy.Address, legacy.Signature
	m.Evidence = nil
	if legacy.Code == MsgPreprepare {
		// The legacy codec decodes exactly {View, Proposal}. Preprepare also
		// accepts the post-Permissionless justification fields, so reject them
		// here or a legacy PRE-PREPARE could carry fields deployed nodes refuse.
		var preprepare struct {
			View     *View
			Proposal *types.Block
		}
		if err := rlp.DecodeBytes(legacy.Msg, &preprepare); err != nil {
			return fmt.Errorf("%w: legacy preprepare: %v", ErrInvalidMessage, err)
		}
		m.Msg = legacy.Msg
		return nil
	}
	var subject *Subject
	if err := rlp.DecodeBytes(legacy.Msg, &subject); err != nil {
		return err
	}
	if subject == nil || subject.View == nil || subject.PrevHash != legacy.Hash {
		return fmt.Errorf("%w: inconsistent legacy subject", ErrInvalidMessage)
	}
	var (
		payload any
		err     error
	)
	switch legacy.Code {
	case MsgPrepare:
		payload = &Prepare{View: subject.View, Digest: subject.Digest}
	case MsgCommit:
		payload = &Commit{View: subject.View, Digest: subject.Digest, CommittedSeal: legacy.CommittedSeal}
	case MsgRoundChange:
		if !common.EmptyHash(subject.Digest) {
			return fmt.Errorf("%w: legacy round change carries a digest", ErrInvalidMessage)
		}
		payload = &RoundChange{View: subject.View}
	default:
		return ErrInvalidMessage
	}
	m.Msg, err = Encode(payload)
	return err
}

// toPrePermissionless is the inverse of fromPrePermissionless. Fields that the
// legacy format cannot express, such as Evidence or a prepared claim, are
// rejected instead of being silently dropped from the signed payload.
func (m *Message) toPrePermissionless() (*PrePermissionlessMessage, error) {
	if len(m.Evidence) != 0 {
		return nil, fmt.Errorf("%w: evidence in legacy message", ErrInvalidMessage)
	}
	legacy := &PrePermissionlessMessage{Hash: m.PrevHash, Code: m.Code, Address: m.Address, Signature: m.Signature}
	subject := &Subject{PrevHash: m.PrevHash}
	switch m.Code {
	case MsgPreprepare:
		legacy.Msg = m.Msg
		return legacy, nil
	case MsgPrepare:
		var prepare *Prepare
		if err := m.Decode(&prepare); err != nil {
			return nil, err
		}
		subject.View, subject.Digest = prepare.View, prepare.Digest
	case MsgCommit:
		var commit *Commit
		if err := m.Decode(&commit); err != nil {
			return nil, err
		}
		subject.View, subject.Digest = commit.View, commit.Digest
		legacy.CommittedSeal = commit.CommittedSeal
	case MsgRoundChange:
		var roundChange *RoundChange
		if err := m.Decode(&roundChange); err != nil {
			return nil, err
		}
		if roundChange.Prepared != nil {
			return nil, fmt.Errorf("%w: prepared claim in legacy round change", ErrInvalidMessage)
		}
		subject.View = roundChange.View
	default:
		return nil, ErrInvalidMessage
	}
	var err error
	legacy.Msg, err = Encode(subject)
	return legacy, err
}
