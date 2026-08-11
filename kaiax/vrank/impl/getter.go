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
	"math/big"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// cfReport reads block blockNum's committed failed list from header.VRank (empty for empty /
// epoch-start blocks, which carry CandTesting not a cfReport).
// Returns ErrNotPermissionless before the fork.
func (v *VRankModule) cfReport(blockNum uint64) ([]common.Address, error) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, vrank.ErrNotPermissionless
	}
	if blockNum%v.vrankEpoch() == 0 {
		return []common.Address{}, nil
	}

	header := v.Chain.GetHeaderByNumber(blockNum)
	if header == nil {
		return nil, vrank.ErrHeaderNotFound
	}
	if len(header.VRank) == 0 {
		return []common.Address{}, nil
	}
	return vrank.DecodeReport(header.VRank)
}

// GetPfReport reads the committed pfReport for block blockNum from header.Extra.
// Returns an empty report if the block was finalized at round 0.
// Returns ErrNotPermissionless if blockNum is before the permissionless fork.
// Used by applyBlocksForPFS so that catchUp can process pre-fork blocks without error.
// GetPfReport returns the proposal failure report for the given block.
func (v *VRankModule) GetPfReport(blockNum uint64) ([]common.Address, error) {
	return v.pfReport(blockNum)
}

func (v *VRankModule) pfReport(blockNum uint64) ([]common.Address, error) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, vrank.ErrNotPermissionless
	}
	header := v.Chain.GetHeaderByNumber(blockNum)
	if header == nil {
		return nil, vrank.ErrHeaderNotFound
	}
	roundByte, err := v.Sealer.Round(header)
	if err != nil {
		return nil, err
	}
	round := uint64(roundByte)
	if round == 0 {
		return []common.Address{}, nil
	}

	pfReport := make([]common.Address, 0, round)
	for r := uint64(0); r < round; r++ {
		proposer, err := v.Valset.GetProposer(blockNum, r)
		if err != nil {
			return nil, err
		}
		pfReport = append(pfReport, proposer)
	}

	return pfReport, nil
}

// EvaluateCandidates computes the cfReport for block blockNum at the given round from in-memory collector state.
// The proposer fills its own header.VRank by evaluating its most recent prior proposal.
// Returns ErrNotPermissionless if header(blockNum+1) is before the permissionless fork.
func (v *VRankModule) EvaluateCandidates(blockNum, round uint64) ([]common.Address, error) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum + 1)) {
		return nil, vrank.ErrNotPermissionless
	}

	if round > maxRound {
		// Candidate messages above maxRound are dropped on receipt, though the
		// preprepare timestamp may still be recorded. Evaluating this view would
		// then see no candidate responses and mark every candidate as failed, so
		// return an empty report instead of that spurious result.
		return []common.Address{}, nil
	}

	vk := vrank.ViewKey{N: blockNum, R: uint8(round)}
	prepreparedAt, expectedBlockHash, viewMap := v.collector.GetViewData(vk)
	if prepreparedAt.IsZero() {
		// No preprepare data — either this node was not a validator for blockNum,
		// or it missed the PREPREPARE message. Either way, nothing to report.
		return []common.Address{}, nil
	}
	candidates, err := v.Valset.GetCandTesting(blockNum)
	if err != nil {
		logger.Error("GetCandTesting failed", "blockNum", blockNum, "err", err)
		return nil, vrank.ErrGetCandidateFailed
	}
	if len(candidates) == 0 {
		return []common.Address{}, nil
	}
	var cfReport []common.Address
	for _, addr := range candidates {
		msgWithTime, arrived := viewMap[addr]
		if !arrived ||
			msgWithTime.Msg == nil ||
			msgWithTime.Msg.BlockHash != expectedBlockHash ||
			msgWithTime.ReceivedAt.Sub(prepreparedAt).Milliseconds() > candidateMsgTimeoutMs {
			cfReport = append(cfReport, addr)
		}
	}

	slices.SortFunc(cfReport, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	return cfReport, nil
}
