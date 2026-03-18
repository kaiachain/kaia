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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/rlp"
)

type Report []common.Address

const (
	VRankPreprepareMsg = 0x17
	VRankCandidateMsg  = 0x18
)

const (
	// Epoch is the number of blocks in one VRank scoring epoch.
	Epoch = uint64(86400)
	// MaxRound is the maximum allowed consensus round per block (range [0, MaxRound]).
	MaxRound = 10
)

type VRankBroadcastEvent struct {
	Targets []common.Address
	Code    int // VRankPreprepareMsg or VRankCandidateMsg
	Msg     any // VRankPreprepare or VRankCandidate
}

type VRankPreprepare struct {
	Block *types.Block
	View  *istanbul.View
	Sig   []byte // proposer's signature over vrankPreprepareSigHash(chainID, blockNum, round, blockHash)
}

type VRankCandidate struct {
	BlockNumber uint64
	Round       uint8
	BlockHash   common.Hash
	Sig         []byte
}

func EncodeReport(report Report) ([]byte, error) {
	if len(report) == 0 {
		return nil, nil
	}

	return rlp.EncodeToBytes(report)
}

func DecodeReport(data []byte) (Report, error) {
	var report []common.Address
	if err := rlp.DecodeBytes(data, &report); err != nil {
		return nil, err
	}
	return report, nil
}
