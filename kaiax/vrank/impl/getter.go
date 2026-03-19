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

package impl

import (
	"bytes"
	"slices"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// GetCfReport reads the committed cfReport for block blockNum from header.VRank.
// Returns an empty report if header.VRank is nil (e.g. epoch-start block).
func (v *VRankModule) GetCfReport(blockNum uint64) (vrank.Report, error) {
	header := v.Chain.GetHeaderByNumber(blockNum)
	if header == nil {
		return nil, vrank.ErrHeaderNotFound
	}
	if len(header.VRank) == 0 {
		return vrank.Report{}, nil
	}
	return vrank.DecodeReport(header.VRank)
}

// GetPfReport reads the committed pfReport for block blockNum from header.Extra.
// Returns an empty report if the block was finalized at round 0.
func (v *VRankModule) GetPfReport(blockNum uint64) (vrank.Report, error) {
	header := v.Chain.GetHeaderByNumber(blockNum)
	if header == nil {
		return nil, vrank.ErrHeaderNotFound
	}
	if len(header.Extra) < types.IstanbulExtraVanity {
		return nil, vrank.ErrHeaderExtraTooShort
	}

	round := uint64(header.Round())
	if round == 0 {
		return vrank.Report{}, nil
	}

	pfReport := make(vrank.Report, 0, round)
	for r := uint64(0); r < round; r++ {
		proposer, err := v.Valset.GetProposer(blockNum, r)
		if err != nil {
			return nil, err
		}
		pfReport = append(pfReport, proposer)
	}

	return pfReport, nil
}

// TallyCfReport computes the cfReport for block blockNum at the given round from in-memory collector state.
// To fill `header(N).VRank`, the proposer of block N must use `TallyCfReport(N-1, last round of N-1)`.
func (v *VRankModule) TallyCfReport(blockNum, round uint64) (vrank.Report, error) {
	// epoch header's VRank should be nil
	if (blockNum+1)%vrankEpoch == 0 {
		return vrank.Report{}, nil
	}
	if round > maxRound {
		return nil, vrank.ErrRoundOutOfRange
	}

	// if I was not a validator for blockNum, I couldn't have collected cfReport for blockNum.
	if !v.isValidator(blockNum) {
		return vrank.Report{}, nil
	}

	vk := vrank.ViewKey{N: blockNum, R: uint8(round)}
	prepreparedAt, expectedBlockHash, viewMap := v.collector.GetViewData(vk)
	if prepreparedAt.IsZero() {
		return nil, vrank.ErrPrepreparedTimeNotSet
	}
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return nil, vrank.ErrGetCandidateFailed
	}
	// safeCands = candidates who responded with expectedBlockHash before deadline.
	safeCands := make(map[common.Address]struct{})
	for sender, msgWithTime := range viewMap {
		if !slices.Contains(candidates, sender) {
			continue
		}
		if msgWithTime.Msg == nil || msgWithTime.Msg.BlockHash != expectedBlockHash {
			continue
		}
		elapsed := msgWithTime.ReceivedAt.Sub(prepreparedAt).Milliseconds()
		// early birds who sent before validator preprepared are safe
		if elapsed <= candidateMsgTimeoutMs {
			safeCands[sender] = struct{}{}
		}
	}

	// cfReport = candidates - safeCands = candidates who did not respond, responded after deadline, or lied (wrong BlockHash).
	var cfReport vrank.Report
	for _, addr := range candidates {
		if _, ok := safeCands[addr]; !ok {
			cfReport = append(cfReport, addr)
		}
	}
	slices.SortFunc(cfReport, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	return cfReport, nil
}
