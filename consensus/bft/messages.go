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

// Package bft hosts wire-level BFT consensus types shared by Istanbul and
// KaiaBFT engines. Any change to these types is a wire-protocol change and must
// be considered against backward compatibility with deployed nodes.
package bft

import (
	"fmt"
	"io"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
)

// Proposal supports retrieving height and serialized block to be used during BFT consensus.
type Proposal interface {
	// Number retrieves the sequence number of this proposal.
	Number() *big.Int

	// Hash retrieves the hash of this proposal.
	Hash() common.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error

	String() string

	ParentHash() common.Hash

	Header() *types.Header

	WithSeal(header *types.Header) *types.Block
}

// Request wraps a proposal to be processed by a BFT engine.
type Request struct {
	Proposal Proposal
}

// View includes a round number and a sequence number.
//
// Sequence is the block number we'd like to commit.
// Each round has a number and is composed by 3 steps: preprepare, prepare and commit.
//
// If the given block is not accepted by validators, a round change will occur
// and the validators start a new round with round+1.
type View struct {
	Round    *big.Int
	Sequence *big.Int
}

// EncodeRLP serializes a View into the Kaia RLP format.
func (v *View) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{v.Round, v.Sequence})
}

// DecodeRLP deserializes a View from a Kaia RLP stream.
func (v *View) DecodeRLP(s *rlp.Stream) error {
	var view struct {
		Round    *big.Int
		Sequence *big.Int
	}
	if err := s.Decode(&view); err != nil {
		return err
	}
	v.Round, v.Sequence = view.Round, view.Sequence
	return nil
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Sequence: %d}", v.Round.Uint64(), v.Sequence.Uint64())
}

// Cmp compares v and y and returns:
//
//	-1 if v < y
//	 0 if v == y
//	+1 if v > y
func (v *View) Cmp(y *View) int {
	sdiff := v.Sequence.Cmp(y.Sequence)
	if sdiff != 0 {
		return sdiff
	}
	rdiff := v.Round.Cmp(y.Round)
	if rdiff != 0 {
		return rdiff
	}
	return 0
}

// ConsensusMsg is the envelope used by p2p for forwarded consensus messages.
type ConsensusMsg struct {
	PrevHash common.Hash
	Payload  []byte
}

// Preprepare is the message sent by the proposer to propose a new block.
//
// A PRE-PREPARE above round 0 is justified by
// RoundChangeCertificate, a quorum of signed ROUND CHANGE messages for its view
// stripped of their Evidence attachments. When any of them claims a
// prepared value, PreparedMessages are the PREPARE/COMMIT votes proving the
// highest claim for Proposal. The prepared block is therefore sent once, as the
// Proposal itself, instead of once per ROUND CHANGE.
type Preprepare struct {
	View                   *View
	Proposal               Proposal
	RoundChangeCertificate []*Message
	PreparedMessages       []*Message
}

// EncodeRLP serializes a Preprepare into the Kaia RLP format.
func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, struct {
		View                   *View
		Proposal               Proposal
		RoundChangeCertificate []*Message `rlp:"optional"`
		PreparedMessages       []*Message `rlp:"optional"`
	}{b.View, b.Proposal, b.RoundChangeCertificate, b.PreparedMessages})
}

// DecodeRLP deserializes a Preprepare from a Kaia RLP stream.
// Proposal is decoded as a *types.Block, which is the only concrete Proposal
// implementation exchanged on the wire.
func (b *Preprepare) DecodeRLP(s *rlp.Stream) error {
	var preprepare struct {
		View                   *View
		Proposal               *types.Block
		RoundChangeCertificate []*Message `rlp:"optional"`
		PreparedMessages       []*Message `rlp:"optional"`
	}
	if err := s.Decode(&preprepare); err != nil {
		return err
	}
	b.View, b.Proposal = preprepare.View, preprepare.Proposal
	b.RoundChangeCertificate, b.PreparedMessages = preprepare.RoundChangeCertificate, preprepare.PreparedMessages
	return nil
}

// PreparedCertificate proves that a proposal reached the prepared state in a
// prior round. Messages are the signed PREPARE/COMMIT envelopes that form the
// quorum; carrying the complete envelopes lets receivers authenticate every
// voter without persisting the certificate in the block header. It travels as
// ROUND CHANGE Evidence and is never signed by the ROUND CHANGE sender.
type PreparedCertificate struct {
	View     *View
	Proposal *types.Block
	Messages []*Message
}

// PreparedClaim is the signed summary of the PreparedCertificate a ROUND
// CHANGE sender holds: the round in which it prepared and the proposal hash.
type PreparedClaim struct {
	Round  *big.Int
	Digest common.Hash
}

// Prepare is the PREPARE payload.
type Prepare struct {
	View   *View
	Digest common.Hash
}

// Commit is the COMMIT payload. Its committed seal is part
// of Msg and is therefore authenticated by the outer Message signature.
type Commit struct {
	View          *View
	Digest        common.Hash
	CommittedSeal []byte
}

// RoundChange is the ROUND-CHANGE payload. The certificate proving Prepared is
// carried once as the outer Message Evidence.
type RoundChange struct {
	View     *View
	Prepared *PreparedClaim `rlp:"optional,nilList"`
}
