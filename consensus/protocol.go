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
// This file is derived from quorum/consensus/protocol.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package consensus

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p"
)

// ConsensusMsgCode is the p2p message code for BFT consensus messages.
const ConsensusMsgCode = 0x11

// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *types.Block)
	// FindPeers retrives peers by addresses
	FindPeers(map[common.Address]bool) map[common.Address]Peer

	FindCNPeers(map[common.Address]bool) map[common.Address]Peer

	GetCNPeers() map[common.Address]Peer

	GetENPeers() map[common.Address]Peer

	RegisterValidator(conType common.ConnType, validator p2p.PeerTypeValidator)
}

// Peer defines the interface to communicate with peer
type Peer interface {
	// Send sends the message to this peer
	Send(msgcode uint64, data interface{}) error

	// RegisterConsensusMsgCode registers the channel of consensus msg.
	RegisterConsensusMsgCode(msgCode uint64) error
}
