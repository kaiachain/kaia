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

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/vrank VRankModule

package vrank

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax"
)

// Chain is the narrow view VRank needs of the blockchain. State access and
// contract reads are not required: scoring paths that needed them are gated
// behind IsPermissionlessForkEnabled and will be wired to a real backend once
// the permissionless system contracts land.
type Chain interface {
	CurrentHeader() *types.Header
	GetHeaderByNumber(number uint64) *types.Header
}

// BlsPubkeyGetter is a narrow interface for looking up a node's registered BLS public key.
// Typically satisfied by the Randao module, which reads KIP-113 state.
type BlsPubkeyGetter interface {
	GetBlsPubkey(nodeId common.Address, num *big.Int) (bls.PublicKey, error)
}

// Sealer is a narrow view of the seal data VRank reads and verifies.
// Typically satisfied by consensus.Sealer via the istanbul backend.
type Sealer interface {
	Round(header *types.Header) (byte, error)
	Author(header *types.Header) (common.Address, error)
	RawSeals(header *types.Header) ([]byte, [][]byte, error)
	RecoverCommitters(hash common.Hash, round byte, seals [][]byte) ([]common.Address, error)
	Quorum(blockNum uint64, qualifiedlen, committeeSize int) int
}

type VRankModule interface {
	kaiax.BaseModule
	kaiax.HeaderModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	HandleIstanbulPreprepare(block *types.Block, view *bft.View)
	HandleVRankPreprepare(msg *VRankPreprepare) error
	HandleVRankCandidate(msg *VRankCandidate) error
	EvaluateCandidates(blockNum, round uint64) ([]common.Address, error)
	GetPfReport(blockNum uint64) ([]common.Address, error)
	GetPFS(blockNum uint64) (map[common.Address]uint64, error)
	GetCFS(blockNum uint64) (map[common.Address]uint64, error)

	SubscribeVRank(sink chan<- *VRankBroadcastEvent) event.Subscription
}

// VRankModuleHost is implemented by the protocol manager to accept the VRank
// module for p2p send/receive of VRank messages.
type VRankModuleHost interface {
	RegisterVRankModule(module VRankModule)
}
