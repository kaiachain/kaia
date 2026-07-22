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

// CfReport is the non-epoch header.VRank payload. TargetBlock is the proposer's own most recent
// prior block this epoch — reporting only its own block keeps any failure in its own
// byzantine-filterable column. Failed lists CandTesting addresses that missed TargetBlock.
type CfReport struct {
	TargetBlock uint64
	Failed      []common.Address
}

// EncodeReport returns nil for an empty report so an absent VRank stays the "nothing" encoding.
func EncodeReport(targetBlock uint64, failed []common.Address) ([]byte, error) {
	if len(failed) == 0 {
		return nil, nil
	}
	return rlp.EncodeToBytes(&CfReport{TargetBlock: targetBlock, Failed: failed})
}

func EncodeAddressList(addrs []common.Address) ([]byte, error) {
	return rlp.EncodeToBytes(addrs)
}

func DecodeReport(data []byte) (uint64, []common.Address, error) {
	var report CfReport
	if err := rlp.DecodeBytes(data, &report); err != nil {
		return 0, nil, err
	}
	return report.TargetBlock, report.Failed, nil
}
