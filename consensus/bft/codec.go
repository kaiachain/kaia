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

// ErrInvalidMessage indicates the message code is not recognized.
var ErrInvalidMessage = errors.New("invalid message")

// Message is the envelope transmitted between BFT validators.
//
// Justification is an optional, unsigned attachment. Only a post-Permissionless
// ROUND CHANGE uses it, to carry the RLP-encoded PreparedCertificate behind the
// PreparedClaim it signs. Keeping the evidence outside the signature lets a
// proposer embed the small signed ROUND CHANGE in a certificate without the
// prepared block; receivers bind the attachment to the signed claim instead.
type Message struct {
	Hash          common.Hash
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
	Justification []byte
}

// EncodeRLP serializes m into the Kaia RLP format. An empty Justification is
// omitted, so ordinary messages keep their legacy encoding. rlp omits only a
// nil optional field, so a non-nil empty slice is normalized to nil first.
func (m *Message) EncodeRLP(w io.Writer) error {
	justification := m.Justification
	if len(justification) == 0 {
		justification = nil
	}
	return rlp.Encode(w, struct {
		Hash          common.Hash
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
		Justification []byte `rlp:"optional"`
	}{m.Hash, m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal, justification})
}

// DecodeRLP loads the consensus fields from a Kaia RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Hash          common.Hash
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
		Justification []byte `rlp:"optional"`
	}
	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Hash, m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Hash, msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	m.Justification = msg.Justification
	return nil
}

// FromPayload decodes b into m and, when validateFn is non-nil, verifies the
// signer matches m.Address.
func (m *Message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
	if err := rlp.DecodeBytes(b, &m); err != nil {
		return err
	}
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
		if len(m.CommittedSeal) != crypto.SignatureLength {
			return fmt.Errorf("%w: committed seal length %d", ErrInvalidMessage, len(m.CommittedSeal))
		}
	case MsgPreprepare, MsgPrepare, MsgRoundChange:
		if len(m.CommittedSeal) != 0 {
			return fmt.Errorf("%w: unexpected committed seal on message code %d", ErrInvalidMessage, m.Code)
		}
	}
	if len(m.Justification) != 0 && m.Code != MsgRoundChange {
		return fmt.Errorf("%w: unexpected justification on message code %d", ErrInvalidMessage, m.Code)
	}
	return nil
}

// Payload returns the RLP-encoded message including the signature.
func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

// PayloadNoSig returns the RLP-encoded message with Signature zeroed, used for
// recovering the signer address from the attached Signature. Justification is
// deliberately excluded: it is unsigned evidence checked against the payload.
func (m *Message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&Message{
		Hash:          m.Hash,
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// Decode unmarshals m.Msg into val.
func (m *Message) Decode(val any) error {
	return rlp.DecodeBytes(m.Msg, val)
}

// WithoutJustification returns a shallow copy of m without its unsigned
// attachment. The copy keeps the original signature.
func (m *Message) WithoutJustification() *Message {
	stripped := *m
	stripped.Justification = nil
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
	case MsgPrepare, MsgCommit:
		var subject *Subject
		if decodeErr := m.Decode(&subject); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = subject.View
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
