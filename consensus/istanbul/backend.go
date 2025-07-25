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
// This file is derived from quorum/consensus/istanbul/backend.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package istanbul

import (
	"math/big"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
)

// Backend provides application specific functions for Istanbul core
//
//go:generate mockgen -destination=./mocks/backend_mock.go -package=mock_istanbul github.com/kaiachain/kaia/consensus/istanbul Backend
type Backend interface {
	// Address returns the owner's address
	Address() common.Address

	// EventMux returns the event mux in backend
	EventMux() *event.TypeMux

	// Broadcast sends a message to all validators (include self)
	Broadcast(prevHash common.Hash, payload []byte) error

	// Gossip sends a message to all validators (exclude self)
	Gossip(payload []byte) error

	GossipSubPeer(prevHash common.Hash, payload []byte)

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposal Proposal, seals [][]byte) error

	// Verify verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	Verify(Proposal) (time.Duration, error)

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	// CheckSignature verifies the signature by checking if it's signed by
	// the given validator
	CheckSignature(data []byte, addr common.Address, sig []byte) error

	// LastProposal retrieves latest committed proposal and the address of proposer
	LastProposal() (Proposal, common.Address)

	// HasPropsal checks if the combination of the given hash and height matches any existing blocks
	HasPropsal(hash common.Hash, number *big.Int) bool

	// GetProposer returns the proposer of the given block height
	GetProposer(number uint64) common.Address

	// HasBadProposal returns whether the proposal with the hash is a bad proposal
	HasBadProposal(hash common.Hash) bool

	GetRewardBase() common.Address

	SetCurrentView(view *View)

	NodeType() common.ConnType

	GetValidatorSet(num uint64) (*BlockValSet, error)

	GetCommitteeState(num uint64) (*RoundCommitteeState, error)

	GetCommitteeStateByRound(num uint64, round uint64) (*RoundCommitteeState, error)

	GetProposerByRound(num uint64, round uint64) (common.Address, error)
}
