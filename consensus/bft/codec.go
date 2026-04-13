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
type Message struct {
	Hash          common.Hash
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

// EncodeRLP serializes m into the Kaia RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{m.Hash, m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
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
	}
	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Hash, m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Hash, msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

// FromPayload decodes b into m and, when validateFn is non-nil, verifies the
// signer matches m.Address.
func (m *Message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
	if err := rlp.DecodeBytes(b, &m); err != nil {
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

// Payload returns the RLP-encoded message including the signature.
func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

// PayloadNoSig returns the RLP-encoded message with Signature zeroed, used for
// recovering the signer address from the attached Signature.
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
	case MsgPrepare, MsgCommit, MsgRoundChange:
		var subject *Subject
		if decodeErr := m.Decode(&subject); decodeErr != nil {
			return nil, decodeErr
		}
		msgView = subject.View
	default:
		return nil, ErrInvalidMessage
	}
	return msgView, nil
}

// Encode wraps rlp.EncodeToBytes for convenience at BFT call sites.
func Encode(val any) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
