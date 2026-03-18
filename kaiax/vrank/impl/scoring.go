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
	"math/big"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

const scoreCacheProbeLookback = uint64(64)

func calcEpochStart(blockNum uint64) uint64 {
	return blockNum - (blockNum % vrankEpoch)
}

// GetPFS computes the running Proposal Failure Score up to blockNum.
// pfs(N) -> map[proposerAddr]score for blocks [epochBegin(N), N].
// Returns ErrNotPermissionless if blockNum is before the permissionless fork.
// Returns ErrFutureBlock if blockNum exceeds the current chain head.
func (v *VRankModule) GetPFS(blockNum uint64) (map[common.Address]uint64, error) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, vrank.ErrNotPermissionless
	}
	cur := v.Chain.CurrentHeader()
	if cur == nil || blockNum > cur.Number.Uint64() {
		return nil, vrank.ErrFutureBlock
	}

	var (
		epochStart = calcEpochStart(blockNum)
		start      uint64
		seed       map[common.Address]uint64
		err        error
	)

	if cachedNum, cachedScores, ok := v.loadNearbyPFSCache(blockNum); ok {
		// 1. exact cache hit
		if cachedNum == blockNum {
			return cloneMap(cachedScores), nil
		}

		// 2. nearby cache hit - use it as seed
		start = cachedNum + 1
		seed = cachedScores
	} else if cpNum, cpScores, _, ok := v.loadCheckpointInEpoch(blockNum); ok {
		// 3. DB checkpoint hit - use it as seed
		start = cpNum + 1
		seed = cpScores
	} else {
		// 4. no cache or DB checkpoint hit - new seed
		start = epochStart
		seed = make(map[common.Address]uint64)
	}

	seed, err = v.computePFS(start, blockNum, seed)
	if err != nil {
		return nil, err
	}

	v.pfsCache.Add(blockNum, cloneMap(seed))
	return cloneMap(seed), nil
}

// GetCFS computes the running Candidate Failure Score up to blockNum.
// cfs(N) -> map[candidateAddr]score for blocks [epochBegin(N), N].
// For the validator state transition at block N, use GetCFS(N-1).
// Returns ErrNotPermissionless if blockNum is before the permissionless fork.
// Returns ErrFutureBlock if blockNum exceeds the current chain head.
func (v *VRankModule) GetCFS(blockNum uint64) (map[common.Address]uint64, error) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, vrank.ErrNotPermissionless
	}
	cur := v.Chain.CurrentHeader()
	if cur == nil || blockNum > cur.Number.Uint64() {
		return nil, vrank.ErrFutureBlock
	}

	var (
		epochStart = calcEpochStart(blockNum)
		start      = epochStart
		end        = blockNum
		seed       map[common.Address]map[common.Address]uint64
		err        error
	)

	if cachedNum, cachedCPMatrix, ok := v.loadNearbyCPMatrix(blockNum); ok {
		// 1. exact cache hit
		if cachedNum == blockNum {
			return v.generateCFSFromCPMatrix(epochStart, cachedCPMatrix)
		}

		// 2. nearby cache hit - use it as seed
		start = cachedNum + 1
		seed = cachedCPMatrix
	} else if cpNum, _, cpMatrix, ok := v.loadCheckpointInEpoch(blockNum); ok {
		// 3. DB checkpoint hit - use it as seed
		start = cpNum + 1
		seed = cpMatrix
	} else {
		// 4. no cache or DB checkpoint hit - new seed
		start = epochStart
		seed, err = v.newCPMatrix(epochStart)
		if err != nil {
			return nil, err
		}
	}
	seed, err = v.computeCPMatrix(start, end, seed)
	if err != nil {
		return nil, err
	}

	cfs, err := v.generateCFSFromCPMatrix(epochStart, seed)
	if err != nil {
		return nil, err
	}

	v.cpMatrixCache.Add(blockNum, cloneCPMatrix(seed))
	return cloneMap(cfs), nil
}

func (v *VRankModule) computePFS(start, end uint64, seed map[common.Address]uint64) (map[common.Address]uint64, error) {
	pfs := cloneMap(seed)
	for blockNum := start; blockNum <= end; blockNum++ {
		report, err := v.pfReport(blockNum)
		if err != nil {
			return nil, err
		}
		for _, addr := range report {
			pfs[addr]++
		}
	}
	return pfs, nil
}

func (v *VRankModule) newCPMatrix(epochStart uint64) (map[common.Address]map[common.Address]uint64, error) {
	// cpMatrix[candidate][proposer] = total failures reported in epoch.
	cpMatrix := make(map[common.Address]map[common.Address]uint64)

	// Pre-seed all candidates so every candidate appears in the output, even with 0 CFS.
	candidates, err := v.Valset.GetCandidates(epochStart)
	if err != nil {
		return nil, err
	}
	for _, cand := range candidates {
		cpMatrix[cand] = make(map[common.Address]uint64)
	}
	return cpMatrix, nil
}

// computeCPMatrix accumulates the candidate-proposer matrix for N in [start, end].
func (v *VRankModule) computeCPMatrix(start, end uint64, seed map[common.Address]map[common.Address]uint64) (map[common.Address]map[common.Address]uint64, error) {
	cpMatrix := cloneCPMatrix(seed)
	for blockNum := start; blockNum <= end; blockNum++ {
		header := v.Chain.GetHeaderByNumber(blockNum)
		if header == nil {
			return nil, vrank.ErrHeaderNotFound
		}

		cfReport, err := v.cfReport(blockNum)
		if err != nil {
			return nil, err
		}
		if len(cfReport) == 0 {
			continue
		}

		round := uint64(header.Round())
		reporter, err := v.Valset.GetProposer(blockNum, round)
		if err != nil {
			return nil, err
		}

		for _, candidate := range cfReport {
			if _, ok := cpMatrix[candidate]; !ok {
				logger.Warn("cfReport contains address not in candidates list; skipping", "blockNum", blockNum, "candidate", candidate.Hex())
				continue
			}
			cpMatrix[candidate][reporter]++
		}
	}
	return cpMatrix, nil
}

func (v *VRankModule) generateCFSFromCPMatrix(epochStart uint64, cpMatrix map[common.Address]map[common.Address]uint64) (map[common.Address]uint64, error) {
	// Determine F from validator count at epoch start.
	committee, err := v.Valset.GetCommittee(epochStart, 0) // TODO: fetch value from AddressBookV2
	if err != nil {
		return nil, err
	}
	F := max(0, (len(committee)-1)/3) // make sure F>=0, just in case

	return byzantineFilter(cpMatrix, F), nil
}

// byzantineFilter computes CFS scores from pre-aggregated per-candidate failure data.
//
// failuresByCandidate[candidate][reporter] is the number of times reporter included
// candidate in cfReport over the epoch.
// F is the number of highest reporter totals to discard per candidate (MAX_BYZANTINE_NODES).
func byzantineFilter(
	cpMatrix map[common.Address]map[common.Address]uint64,
	F int,
) map[common.Address]uint64 {
	cfs := make(map[common.Address]uint64)
	for cand, reporterToScore := range cpMatrix {
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

func cloneMap(src map[common.Address]uint64) map[common.Address]uint64 {
	ret := make(map[common.Address]uint64, len(src))
	maps.Copy(ret, src)
	return ret
}

func cloneCPMatrix(src map[common.Address]map[common.Address]uint64) map[common.Address]map[common.Address]uint64 {
	ret := make(map[common.Address]map[common.Address]uint64, len(src))
	for candidate, reporterScores := range src {
		ret[candidate] = cloneMap(reporterScores)
	}
	return ret
}

func (v *VRankModule) loadNearbyPFSCache(blockNum uint64) (uint64, map[common.Address]uint64, bool) {
	epochStart := calcEpochStart(blockNum)
	for i := uint64(0); i <= scoreCacheProbeLookback && i <= blockNum-epochStart; i++ {
		candidateNum := blockNum - i
		if cached, ok := v.pfsCache.Get(candidateNum); ok {
			return candidateNum, cloneMap(cached.(map[common.Address]uint64)), true
		}
	}
	return 0, nil, false
}

func (v *VRankModule) loadNearbyCPMatrix(blockNum uint64) (uint64, map[common.Address]map[common.Address]uint64, bool) {
	epochStart := calcEpochStart(blockNum)
	for i := uint64(0); i <= scoreCacheProbeLookback && i <= blockNum-epochStart; i++ {
		candidateNum := blockNum - i
		if cached, ok := v.cpMatrixCache.Get(candidateNum); ok {
			cpMatrix := cloneCPMatrix(cached.(map[common.Address]map[common.Address]uint64))
			return candidateNum, cpMatrix, true
		}
	}
	return 0, nil, false
}
