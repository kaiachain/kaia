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
// This file is derived from quorum/consensus/istanbul/core/roundstate.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/rlp"
)

// newRoundState creates a new roundState instance with the given view and validatorSet.
// A round change retains the accepted PRE-PREPARE, which is the proposal behind a
// local hash lock, and the prepared certificate behind that lock when known.
// The new round is recorded in round, independently of the PRE-PREPARE's
// original view.
func newRoundState(view *bft.View, qualified *valset.AddressSet, lockedHash common.Hash, preprepare *bft.Preprepare, preparedCertificate *bft.PreparedCertificate, pendingRequest *bft.Request, hasBadProposal func(hash common.Hash) bool) *roundState {
	return &roundState{
		round:               view.Round,
		sequence:            view.Sequence,
		Preprepare:          preprepare,
		Prepares:            newMessageSet(qualified),
		Commits:             newMessageSet(qualified),
		lockedHash:          lockedHash,
		preparedCertificate: preparedCertificate,
		mu:                  new(sync.RWMutex),
		pendingRequest:      pendingRequest,
		hasBadProposal:      hasBadProposal,
	}
}

// roundState stores the consensus state
type roundState struct {
	round               *big.Int
	sequence            *big.Int
	Preprepare          *bft.Preprepare
	Prepares            *messageSet
	Commits             *messageSet
	lockedHash          common.Hash
	preparedCertificate *bft.PreparedCertificate
	pendingRequest      *bft.Request

	mu             *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool

	// Ignore RLP ----------------------------------------------------------------------------
	qualified            *valset.AddressSet
	committee            *valset.AddressSet
	parentHash           common.Hash // parent of the height being decided; fixed across its rounds and equal to the ParentHash of every proposal Verify accepts
	proposer             common.Address
	committeeSize        uint64
	requiredMessageCount int
	f                    int
}

func (s *roundState) GetPrepareOrCommitSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := s.Prepares.Size() + s.Commits.Size()

	// find duplicate one
	for _, m := range s.Prepares.Values() {
		if s.Commits.Get(m.Address) != nil {
			result--
		}
	}
	return result
}

// proposalSubject identifies the proposal a PREPARE or COMMIT votes for. It is
// a local view of the round state, not a wire type.
type proposalSubject struct {
	View     *bft.View
	Digest   common.Hash
	PrevHash common.Hash
}

// matches reports whether a vote for view and digest is for s. The vote's
// PrevHash is checked against the height's parent by checkMessage.
func (s *proposalSubject) matches(view *bft.View, digest common.Hash) bool {
	return view != nil && view.Cmp(s.View) == 0 && digest == s.Digest
}

func (s *proposalSubject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v, ParentHash: %v}", s.View, s.Digest.String(), s.PrevHash.Hex())
}

func (s *roundState) Subject() *proposalSubject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare == nil {
		return nil
	}

	return &proposalSubject{
		View: &bft.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest:   s.Preprepare.Proposal.Hash(),
		PrevHash: s.Preprepare.Proposal.ParentHash(),
	}
}

func (s *roundState) SetPreprepare(preprepare *bft.Preprepare) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Preprepare = preprepare
}

func (s *roundState) Proposal() bft.Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare != nil {
		return s.Preprepare.Proposal
	}

	return nil
}

// PreparedCertificate returns the signed quorum that established the local
// lock. It is carried across rounds and advertised in ROUND-CHANGE messages.
func (s *roundState) PreparedCertificate() *bft.PreparedCertificate {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.preparedCertificate
}

func (s *roundState) LockedRound() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.preparedCertificate == nil || s.preparedCertificate.View == nil {
		return nil
	}
	return new(big.Int).Set(s.preparedCertificate.View.Round)
}

// AdoptPreparedCertificate updates the local lock from a strictly newer,
// independently verified certificate. This is used when the node did not
// observe the original PREPARE quorum itself but learns it from a justified
// PRE-PREPARE.
// Keeping the comparison with the mutation prevents any caller from
// accidentally downgrading a newer local lock.
func (s *roundState) AdoptPreparedCertificate(cert *bft.PreparedCertificate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cert == nil || cert.View == nil || cert.View.Round == nil || cert.Proposal == nil {
		return
	}
	if s.preparedCertificate != nil && s.preparedCertificate.View != nil &&
		s.preparedCertificate.View.Round != nil && cert.View.Round.Cmp(s.preparedCertificate.View.Round) <= 0 {
		return
	}
	s.lockedHash = cert.Proposal.Hash()
	s.preparedCertificate = cert
}

func (s *roundState) SetRound(r *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.round = new(big.Int).Set(r)
}

func (s *roundState) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState) SetSequence(seq *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence = seq
}

func (s *roundState) Sequence() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sequence
}

func (s *roundState) LockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Preprepare != nil {
		s.lockedHash = s.Preprepare.Proposal.Hash()
		block, ok := s.Preprepare.Proposal.(*types.Block)
		if !ok {
			s.preparedCertificate = nil
			return
		}
		messages := make([]*bft.Message, 0, s.Prepares.Size()+s.Commits.Size())
		seen := make(map[common.Address]struct{})
		for _, set := range []*messageSet{s.Prepares, s.Commits} {
			for _, msg := range set.Values() {
				if _, exists := seen[msg.Address]; exists {
					continue
				}
				seen[msg.Address] = struct{}{}
				messages = append(messages, msg)
			}
		}
		s.preparedCertificate = &bft.PreparedCertificate{
			View: &bft.View{
				Round:    new(big.Int).Set(s.round),
				Sequence: new(big.Int).Set(s.sequence),
			},
			Proposal: block,
			Messages: messages,
		}
	}
}

func (s *roundState) UnlockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockedHash = common.Hash{}
	s.preparedCertificate = nil
}

func (s *roundState) IsHashLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if common.EmptyHash(s.lockedHash) {
		return false
	}
	return !s.hasBadProposal(s.GetLockedHash())
}

func (s *roundState) GetLockedHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedHash
}

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var ss struct {
		Round               *big.Int
		Sequence            *big.Int
		Preprepare          *bft.Preprepare `rlp:"nil"`
		Prepares            *messageSet
		Commits             *messageSet
		LockedHash          common.Hash
		PendingRequest      *bft.Request             `rlp:"nil"`
		PreparedCertificate *bft.PreparedCertificate `rlp:"optional,nilList"`
	}

	if err := stream.Decode(&ss); err != nil {
		return err
	}
	s.round = ss.Round
	s.sequence = ss.Sequence
	s.Preprepare = ss.Preprepare
	s.Prepares = ss.Prepares
	s.Commits = ss.Commits
	s.lockedHash = ss.LockedHash
	s.preparedCertificate = ss.PreparedCertificate
	s.pendingRequest = ss.PendingRequest
	s.mu = new(sync.RWMutex)

	return nil
}

// EncodeRLP should write the RLP encoding of its receiver to w.
// If the implementation is a pointer method, it may also be
// called for nil pointers.
//
// Implementations should generate valid RLP. The data written is
// not verified at the moment, but a future version might. It is
// recommended to write only a single value but writing multiple
// values or no value at all is also permitted.
func (s *roundState) EncodeRLP(w io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return rlp.Encode(w, []interface{}{
		s.round,
		s.sequence,
		s.Preprepare,
		s.Prepares,
		s.Commits,
		s.lockedHash,
		s.pendingRequest,
		s.preparedCertificate,
	})
}
