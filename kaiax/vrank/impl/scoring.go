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

	start, seed := v.lookupPFSSeed(blockNum)
	if start <= blockNum {
		var err error
		seed, err = v.applyBlocksForPFS(start, blockNum, seed)
		if err != nil {
			return nil, err
		}
	}

	v.pfsCache.Add(blockNum, seed)
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

	start, seed, err := v.lookupCFSSeed(blockNum)
	if err != nil {
		return nil, err
	}
	if start <= blockNum {
		seed, err = v.applyBlocksForCPMatrix(start, blockNum, seed)
		if err != nil {
			return nil, err
		}
	}

	cfs, err := v.generateCFSFromCPMatrix(blockNum, seed)
	if err != nil {
		return nil, err
	}

	v.cpMatrixCache.Add(blockNum, seed)
	return cfs, nil
}

// lookupPFSSeed returns (start, seed) where seed holds PFS accumulated up to start-1.
// Callers should apply [start, blockNum] on top of seed to reach blockNum.
func (v *VRankModule) lookupPFSSeed(blockNum uint64) (start uint64, seed map[common.Address]uint64) {
	if cachedNum, pfs, ok := v.lookupPFSCache(blockNum); ok {
		return cachedNum + 1, pfs
	}
	if cpNum, pfs, ok := v.loadPFSCheckpointInEpoch(blockNum); ok {
		return cpNum + 1, pfs
	}
	return calcEpochStart(blockNum), make(map[common.Address]uint64)
}

// lookupCFSSeed returns (start, seed) where seed holds the CP matrix accumulated up to start-1.
// Callers should apply [start, blockNum] on top of seed to reach blockNum.
func (v *VRankModule) lookupCFSSeed(blockNum uint64) (start uint64, seed vrank.CPMatrix, err error) {
	if cachedNum, cpMatrix, ok := v.lookupCPMatrixCache(blockNum); ok {
		return cachedNum + 1, cpMatrix, nil
	}
	if cpNum, cpMatrix, ok := v.loadCPMatrixCheckpointInEpoch(blockNum); ok {
		return cpNum + 1, cpMatrix, nil
	}
	cpMatrix, err := v.newCPMatrix(blockNum)
	if err != nil {
		return 0, nil, err
	}
	return calcEpochStart(blockNum), cpMatrix, nil
}

func (v *VRankModule) newCPMatrix(blockNum uint64) (vrank.CPMatrix, error) {
	if blockNum == 0 {
		// Block 0 is genesis; there is no parent block, so GetCandidates(0) would
		// underflow to GetHeaderByNumber(MaxUint64). Start with an empty matrix —
		// the first epoch begins scoring from block 1 onward.
		return vrank.NewCPMatrix(nil), nil
	}
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil {
		return nil, err
	}

	return vrank.NewCPMatrix(candidates), nil
}

// applyBlocksForPFS accumulates the pfReport for blocks in [start, end].
func (v *VRankModule) applyBlocksForPFS(start, end uint64, seed map[common.Address]uint64) (map[common.Address]uint64, error) {
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

// applyBlocksForCPMatrix accumulates the candidate-proposer matrix for blocks in [start, end].
func (v *VRankModule) applyBlocksForCPMatrix(start, end uint64, seed vrank.CPMatrix) (vrank.CPMatrix, error) {
	cpMatrix := seed.Clone()
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
			cpMatrix.Increment(candidate, reporter)
		}
	}
	return cpMatrix, nil
}

func (v *VRankModule) generateCFSFromCPMatrix(blockNum uint64, cpMatrix vrank.CPMatrix) (map[common.Address]uint64, error) {
	// Determine F from validator count at epoch start.
	committee, err := v.Valset.GetCommittee(blockNum, 0) // TODO: fetch value from AddressBookV2
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
	cpMatrix vrank.CPMatrix,
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

func (v *VRankModule) lookupPFSCache(blockNum uint64) (uint64, map[common.Address]uint64, bool) {
	epochStart := calcEpochStart(blockNum)
	for i := uint64(0); i <= min(scoreCacheProbeLookback, blockNum-epochStart); i++ {
		candidateNum := blockNum - i
		if cached, ok := v.pfsCache.Get(candidateNum); ok {
			return candidateNum, cached.(map[common.Address]uint64), true
		}
	}
	return 0, nil, false
}

func (v *VRankModule) lookupCPMatrixCache(blockNum uint64) (uint64, vrank.CPMatrix, bool) {
	epochStart := calcEpochStart(blockNum)
	for i := uint64(0); i <= min(scoreCacheProbeLookback, blockNum-epochStart); i++ {
		candidateNum := blockNum - i
		if cached, ok := v.cpMatrixCache.Get(candidateNum); ok {
			return candidateNum, cached.(vrank.CPMatrix), true
		}
	}
	return 0, nil, false
}
