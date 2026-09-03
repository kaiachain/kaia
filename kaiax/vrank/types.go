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

package vrank

import (
	"maps"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	blstypes "github.com/kaiachain/kaia/crypto/bls/types"
	"github.com/kaiachain/kaia/rlp"
)

// CPMatrix is a candidate-proposer matrix tracking per-block report counts.
// cpMatrix[candidate][reporter] is the number of times reporter reported candidate.
type CPMatrix map[common.Address]map[common.Address]uint64

func NewCPMatrix(candidates []common.Address) CPMatrix {
	cpMatrix := make(CPMatrix, len(candidates))
	for _, candidate := range candidates {
		cpMatrix[candidate] = make(map[common.Address]uint64)
	}
	return cpMatrix
}

func (m CPMatrix) Clone() CPMatrix {
	clone := make(CPMatrix, len(m))
	for candidate, reporters := range m {
		inner := make(map[common.Address]uint64, len(reporters))
		maps.Copy(inner, reporters)
		clone[candidate] = inner
	}
	return clone
}

func (m CPMatrix) Increment(candidate, reporter common.Address) {
	m[candidate][reporter]++
}

// AddProposer records that `proposer` produced a block in the current epoch.
func (m CPMatrix) AddProposer(proposer common.Address) {
	for candidate := range m {
		if _, ok := m[candidate][proposer]; !ok {
			m[candidate][proposer] = 0
		}
	}
}

// ProposerCount returns the number of distinct proposers seen in the current epoch.
func (m CPMatrix) ProposerCount() int {
	for _, sub := range m {
		return len(sub) // any key should be the same
	}
	return 0
}

const (
	// MaxRound is the maximum allowed consensus round per block (range [0, MaxRound]).
	MaxRound = 10
)

type VRankBroadcastEvent struct {
	Targets []common.Address
	Msg     any // *VRankPreprepare or *VRankCandidate; node/cn routes on the type
}

type VRankPreprepare struct {
	Block *types.Block
	View  *bft.View
	Sig   [crypto.SignatureLength]byte // Sign(vrankPreprepareSigHash(), nodeKey)
}

type VRankCandidate struct {
	BlockNumber uint64
	Round       uint8
	BlockHash   common.Hash
	Sig         [crypto.SignatureLength]byte   // Sign(vrankCandidateSigHash(), nodeKey)
	BlsSig      [blstypes.SignatureLength]byte // Sign(vrankCandidateSigHash(), blsKey)
}

// VRankPayload is the decoded header.VRank field.
type VRankPayload struct {
	// Report is CandTesting(N) at epoch-start blocks and the cfReport elsewhere.
	Report []common.Address
	// ParentRound is the canonical round of block N-1, proven by ParentCommittedSeal.
	// Round 0 claims no proposal failure, so it carries no seals.
	ParentRound         uint8
	ParentCommittedSeal [][]byte
}

// EncodeVRank returns nil for an empty payload so an absent VRank stays the "nothing" encoding.
func EncodeVRank(payload VRankPayload) ([]byte, error) {
	if len(payload.Report) == 0 && payload.ParentRound == 0 && len(payload.ParentCommittedSeal) == 0 {
		return nil, nil
	}
	return rlp.EncodeToBytes(&payload)
}

func DecodeVRank(data []byte) (VRankPayload, error) {
	// Decoded payloads always carry non-nil slices.
	if len(data) == 0 {
		return VRankPayload{Report: []common.Address{}, ParentCommittedSeal: [][]byte{}}, nil
	}

	var payload VRankPayload
	if err := rlp.DecodeBytes(data, &payload); err != nil {
		return VRankPayload{}, err
	}
	if payload.Report == nil {
		payload.Report = []common.Address{}
	}
	if payload.ParentCommittedSeal == nil {
		payload.ParentCommittedSeal = [][]byte{}
	}
	return payload, nil
}
