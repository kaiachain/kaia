// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/types.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package istanbul

import (
	"fmt"
	"io"
	"math/big"

	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/rlp"
)

// Proposal supports retrieving height and serialized block to be used during Istanbul consensus.
type Proposal interface {
	// Number retrieves the sequence number of this proposal
	Number() *big.Int

	// Hash retrieves the hash of this proposal
	Hash() common.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error

	String() string

	ParentHash() common.Hash

	Header() *types.Header

	WithSeal(header *types.Header) *types.Block
}

type Request struct {
	Proposal Proposal
}

// View includes a round number and a sequence number.
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
	return rlp.Encode(w, []interface{}{v.Round, v.Sequence})
}

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
// -1 if v < y
//
//	0 if v == y
//
// +1 if v > y
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

type Preprepare struct {
	View     *View
	Proposal Proposal
}

func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Proposal})
}

func (b *Preprepare) DecodeRLP(s *rlp.Stream) error {
	var preprepare struct {
		View     *View
		Proposal *types.Block
	}

	if err := s.Decode(&preprepare); err != nil {
		return err
	}
	b.View, b.Proposal = preprepare.View, preprepare.Proposal

	return nil
}

type Subject struct {
	View     *View
	Digest   common.Hash
	PrevHash common.Hash
}

func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest, b.PrevHash})
}

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

func (a *Subject) Equal(b *Subject) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Digest == b.Digest &&
		a.PrevHash == b.PrevHash &&
		a.View.Cmp(b.View) == 0
}

func (b *Subject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v, ParentHash: %v}", b.View, b.Digest.String(), b.PrevHash.Hex())
}

type ConsensusMsg struct {
	PrevHash common.Hash
	Payload  []byte
}
