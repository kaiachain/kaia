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
	"maps"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// GetPFS computes the running Proposal Failure Score up to blockNum.
// pfs(N) -> map[proposerAddr]score for blocks [epochBegin(N), N].
// Returns ErrFutureBlock if blockNum exceeds the current chain head.
func (v *VRankModule) GetPFS(blockNum uint64) (map[common.Address]uint64, error) {
	cur := v.Chain.CurrentHeader()
	if cur == nil || blockNum > cur.Number.Uint64() {
		return nil, vrank.ErrFutureBlock
	}
	epochStart := blockNum - (blockNum % vrankEpoch)
	return v.computePFS(epochStart, blockNum)
}

// GetCFS computes the running Candidate Failure Score up to blockNum.
// cfs(N) -> map[candidateAddr]score for blocks [epochBegin(N), N].
// For the validator state transition at block N, use GetCFS(N-1).
// Returns ErrFutureBlock if blockNum exceeds the current chain head.
func (v *VRankModule) GetCFS(blockNum uint64) (map[common.Address]uint64, error) {
	cur := v.Chain.CurrentHeader()
	if cur == nil || blockNum > cur.Number.Uint64() {
		return nil, vrank.ErrFutureBlock
	}
	epochStart := blockNum - (blockNum % vrankEpoch)
	return v.computeCFS(epochStart, blockNum)
}

// computePFS accumulates proposal failure counts from pfReport(N) for N in [start, end].
func (v *VRankModule) computePFS(start, end uint64) (map[common.Address]uint64, error) {
	pfs := make(map[common.Address]uint64)
	for blockNum := start; blockNum <= end; blockNum++ {
		report, err := v.GetPfReport(blockNum)
		if err != nil {
			return nil, err
		}
		for _, addr := range report {
			pfs[addr]++
		}
	}
	return pfs, nil
}

// computeCFS accumulates candidate failure counts with Byzantine filtering for N in [start, end].
func (v *VRankModule) computeCFS(start, end uint64) (map[common.Address]uint64, error) {
	// failuresByCandidate[candidate][reporter] = total failures reported in epoch.
	failuresByCandidate := make(map[common.Address]map[common.Address]uint64)

	// Pre-seed all candidates so every candidate appears in the output, even with 0 CFS.
	candidates, err := v.Valset.GetCandidates(start)
	if err != nil {
		return nil, err
	}
	for _, cand := range candidates {
		failuresByCandidate[cand] = make(map[common.Address]uint64)
	}

	for blockNum := start; blockNum <= end; blockNum++ {
		header := v.Chain.GetHeaderByNumber(blockNum)
		if header == nil {
			return nil, vrank.ErrHeaderNotFound
		}

		round := uint64(header.Round())
		reporter, err := v.Valset.GetProposer(blockNum, round)
		if err != nil {
			return nil, err
		}

		cfReport, err := v.GetCfReport(blockNum)
		if err != nil {
			return nil, err
		}
		if len(cfReport) == 0 {
			continue
		}
		for _, candidate := range cfReport {
			if _, ok := failuresByCandidate[candidate]; !ok {
				logger.Warn("cfReport contains candidate not in GetCandidates list", "blockNum", blockNum, "candidate", candidate.Hex())
				failuresByCandidate[candidate] = make(map[common.Address]uint64)
			}
			failuresByCandidate[candidate][reporter]++
		}
	}

	// Determine F from validator count at epoch start.
	committee, err := v.Valset.GetCommittee(start, 0) // TODO: fetch value from AddressBookV2
	if err != nil {
		return nil, err
	}
	F := max(0, (len(committee)-1)/3) // make sure F>=0, just in case

	return byzantineFilter(failuresByCandidate, F), nil
}

// byzantineFilter computes CFS scores from pre-aggregated per-candidate failure data.
//
// failuresByCandidate[candidate][reporter] is the number of times reporter included
// candidate in cfReport over the epoch.
// F is the number of highest reporter totals to discard per candidate (MAX_BYZANTINE_NODES).
func byzantineFilter(
	failuresByCandidate map[common.Address]map[common.Address]uint64,
	F int,
) map[common.Address]uint64 {
	cfs := make(map[common.Address]uint64)
	for cand, reporterToScore := range failuresByCandidate {
		scores := slices.Collect(maps.Values(reporterToScore))
		slices.Sort(scores)
		if F >= len(scores) {
			// since `scores` contain non-zero scores only, F >= len(scores) can happen, in which case all scores are discarded.
			scores = nil
		} else {
			scores = scores[:len(scores)-F]
		}
		var sum uint64
		for _, t := range scores {
			sum += t
		}
		cfs[cand] = sum
	}
	return cfs
}
