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
	"github.com/kaiachain/kaia/rlp"
)

type CfReport []common.Address

type BroadcastRequest struct {
	Targets []common.Address
	Code    int
	Msg     any
}

type VRankPreprepare struct {
	Block *types.Block
}

type VRankCandidate struct {
	BlockNumber uint64
	Round       uint8
	BlockHash   common.Hash
	Sig         []byte
}

func EncodeCfReport(cfReport CfReport) ([]byte, error) {
	if len(cfReport) == 0 {
		return nil, nil
	}

	return rlp.EncodeToBytes(cfReport)
}

func DecodeCfReport(data []byte) (CfReport, error) {
	var cfReport []common.Address
	if err := rlp.DecodeBytes(data, &cfReport); err != nil {
		return nil, err
	}
	return cfReport, nil
}
