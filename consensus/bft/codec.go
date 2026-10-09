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

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
)

// Message codes carried by a BFT Message. Values MUST match the pre-existing
// istanbul/core constants so that pre-move wire bytes remain valid.
const (
	MsgPreprepare uint64 = iota
	MsgPrepare
	MsgCommit
	MsgRoundChange
	MsgAll
)

// ErrInvalidSigner is returned when the signer address of a Message does not
// match the address recovered from its signature.
var ErrInvalidSigner = errors.New("message not signed by the sender")

// ErrInvalidMessage indicates a malformed envelope or payload, including an
// unrecognized message code.
var ErrInvalidMessage = errors.New("invalid message")

// Message is the envelope transmitted between BFT validators. PrevHash, Code,
// Msg and Address, in that order, are authenticated by Signature. Evidence is
// deliberately excluded from that signature; it must hash to the EvidenceHash
// of the signed claim carried by Msg.
type Message struct {
	PrevHash  common.Hash
	Code      uint64
	Msg       []byte
	Address   common.Address
	Signature []byte
	Evidence  []byte
}

// EncodeRLP serializes m into the Kaia RLP format. An empty Evidence is
// omitted because rlp omits only a nil optional field.
func (m *Message) EncodeRLP(w io.Writer) error {
	evidence := m.Evidence
	if len(evidence) == 0 {
		evidence = nil
	}
	return rlp.Encode(w, struct {
		PrevHash  common.Hash
		Code      uint64
		Msg       []byte
		Address   common.Address
		Signature []byte
		Evidence  []byte `rlp:"optional"`
	}{m.PrevHash, m.Code, m.Msg, m.Address, m.Signature, evidence})
}

// DecodeRLP loads the consensus fields from a Kaia RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		PrevHash  common.Hash
		Code      uint64
		Msg       []byte
		Address   common.Address
		Signature []byte
		Evidence  []byte `rlp:"optional"`
	}
	if err := s.Decode(&msg); err != nil {
		return err
	}
	// A missing optional field decodes as nil and an explicit empty one as
	// []byte{}. EncodeRLP never writes the latter, so reject it: otherwise a
	// message without Evidence would have a second valid encoding, and a relay
	// could re-encode it past gossip deduplication.
	if msg.Evidence != nil && len(msg.Evidence) == 0 {
		return fmt.Errorf("%w: empty evidence field", ErrInvalidMessage)
	}
	m.PrevHash, m.Code, m.Msg, m.Address, m.Signature, m.Evidence = msg.PrevHash, msg.Code, msg.Msg, msg.Address, msg.Signature, msg.Evidence
	return nil
}

// FromPayload decodes b into m and, when validateFn is non-nil, verifies its
// signer.
func (m *Message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
	if err := rlp.DecodeBytes(b, &m); err != nil {
		return err
	}
	return m.validateAndRecover(validateFn)
}

func (m *Message) validateAndRecover(validateFn func([]byte, []byte) (common.Address, error)) error {
	if err := m.validateEnvelope(); err != nil {
		return err
	}
	if validateFn != nil {
		payload, err := m.PayloadNoSig()
		if err != nil {
			return err
		}
		signerAddr, err := validateFn(payload, m.Signature)
		if err != nil {
			return err
		}
		if !bytes.Equal(signerAddr.Bytes(), m.Address.Bytes()) {
			return ErrInvalidSigner
		}
	}
	return nil
}

// validateEnvelope checks the envelope shape before signature recovery.
func (m *Message) validateEnvelope() error {
	switch m.Code {
	case MsgCommit:
		var commit *Commit
		if err := m.Decode(&commit); err != nil || commit == nil {
			return fmt.Errorf("%w: undecodable commit payload: %v", ErrInvalidMessage, err)
		}
		if len(commit.CommittedSeal) != crypto.SignatureLength {
			return fmt.Errorf("%w: committed seal length %d", ErrInvalidMessage, len(commit.CommittedSeal))
		}
	}
	if len(m.Evidence) != 0 && m.Code != MsgRoundChange {
		return fmt.Errorf("%w: unexpected evidence on message code %d", ErrInvalidMessage, m.Code)
	}
	return nil
}

// Payload returns the RLP-encoded message including the signature.
func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

// PayloadNoSig returns the signed preimage, used for recovering the signer
// address from the attached Signature. Evidence is deliberately excluded: it
// is unsigned and checked against the signed claim.
func (m *Message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(struct {
		PrevHash common.Hash
		Code     uint64
		Msg      []byte
		Address  common.Address
	}{
		PrevHash: m.PrevHash,
		Code:     m.Code,
		Msg:      m.Msg,
		Address:  m.Address,
	})
}

// Decode unmarshals m.Msg into val.
func (m *Message) Decode(val any) error {
	return rlp.DecodeBytes(m.Msg, val)
}

// WithoutEvidence returns a shallow copy without its unsigned Evidence.
func (m *Message) WithoutEvidence() *Message {
	stripped := *m
	stripped.Evidence = nil
	return &stripped
}

func (m *Message) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v}", m.Code, m.Address.String())
}

// GetView extracts the view from a Message by decoding the inner payload.
func (m *Message) GetView() (*View, error) {
	var msgView *View
	switch m.Code {
	case MsgPreprepare:
		var preprepare *Preprepare
		if decodeErr := m.Decode(&preprepare); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = preprepare.View
	case MsgPrepare:
		var prepare *Prepare
		if decodeErr := m.Decode(&prepare); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = prepare.View
	case MsgCommit:
		var commit *Commit
		if decodeErr := m.Decode(&commit); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = commit.View
	case MsgRoundChange:
		var roundChange *RoundChange
		if decodeErr := m.Decode(&roundChange); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = roundChange.View
	default:
		return nil, ErrInvalidMessage
	}
	return msgView, nil
}

// Encode wraps rlp.EncodeToBytes for convenience at BFT call sites.
func Encode(val any) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
