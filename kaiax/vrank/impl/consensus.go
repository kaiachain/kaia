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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// VerifyHeader checks the VRank field in the header:
//   - Before the permissionless fork: VRank must be absent.
//   - At epoch-start blocks: VRank must be RLPEncode(CandTesting(N)) or RLPEncode([]).
//   - Otherwise: a cfReport whose failed addresses are sorted, deduped, and ⊆ CandTesting.
func (v *VRankModule) VerifyHeader(header *types.Header, _ *types.Header) error {
	number := header.Number.Uint64()
	permissionless := v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(number))

	if !permissionless {
		// Reject presence, not content: an explicit empty RLP string decodes to a non-nil zero-length slice.
		if header.VRank != nil {
			return vrank.ErrUnexpectedVRankBeforePermissionless
		}
		return nil
	}

	// Epoch-start block
	if number%v.vrankEpoch() == 0 {
		expected, err := v.encodeEpochStartVRank(number)
		if err != nil {
			return err
		}
		if !bytes.Equal(header.VRank, expected) {
			return vrank.ErrEpochStartVRankMismatch
		}
		return nil
	}

	// Non-epoch-start block
	report, err := vrank.DecodeReport(header.VRank)
	if err != nil {
		return vrank.ErrInvalidVRankFormat
	}
	if len(report) == 0 {
		return nil
	}
	// Failures score against the reporter's own byzantine-filterable column regardless of content,
	// so only the failed list is checked (CandTesting is epoch-stable within the epoch).
	candidates, err := v.Valset.GetCandTesting(number - 1)
	if err != nil {
		return err
	}
	return validateNonEpochVRank(report, candidates)
}

// PrepareHeader fills header.VRank per KIP-227.
//
//   - At epoch-start blocks: VRank is CandTesting(N), encoded even when empty.
//   - Otherwise: VRank is a cfReport about this proposer's own most recent prior proposal in
//     the current epoch, or nil when there is no such block or no failures.
func (v *VRankModule) PrepareHeader(header *types.Header) error {
	number := header.Number.Uint64()
	if !v.ChainConfig.IsPermissionlessForkEnabled(header.Number) {
		header.VRank = nil
		return nil
	}

	if number%v.vrankEpoch() == 0 {
		encoded, err := v.encodeEpochStartVRank(number)
		if err != nil {
			return err
		}
		header.VRank = encoded
	} else {
		encoded, err := v.encodeCandidateFailureVRank(number)
		if err != nil {
			return err
		}
		header.VRank = encoded
	}

	return nil
}

// encodeEpochStartVRank returns RLPEncode(CandTesting(N)) or RLPEncode([]) = 0xc0 when empty.
func (v *VRankModule) encodeEpochStartVRank(number uint64) ([]byte, error) {
	candidates, err := v.Valset.GetCandTesting(number)
	if err != nil {
		logger.Error("epoch-start VRank fill: GetCandTesting failed", "num", number, "err", err)
		return nil, err
	}
	encoded, err := vrank.EncodeAddressList(candidates)
	if err != nil {
		logger.Error("epoch-start VRank fill: EncodeAddressList failed", "num", number, "err", err)
		return nil, err
	}
	return encoded, nil
}

func (v *VRankModule) encodeCandidateFailureVRank(number uint64) ([]byte, error) {
	targetNum, round, ok := v.selectReportTarget(number)
	if !ok {
		// No own prior proposal this epoch (first proposal, or restart). Empty report — fail-safe.
		return nil, nil
	}
	report, err := v.EvaluateCandidates(targetNum, round)
	if err != nil {
		logger.Error("Failed to evaluate VRank candidates", "err", err, "targetNum", targetNum, "round", round)
		return nil, err
	}
	v.collector.PruneReported(targetNum) // drop views older than targetNum; targetNum stays for re-report
	if len(report) == 0 {
		return nil, nil
	}
	encoded, err := vrank.EncodeReport(report)
	if err != nil {
		logger.Error("Failed to encode VRank report", "err", err, "report", report)
		return nil, err
	}
	return encoded, nil
}

func validateNonEpochVRank(report, candidates []common.Address) error {
	if isNonCandContained(report, candidates) {
		return vrank.ErrInvalidVRankCandidate
	}
	if !isSorted(report) {
		return vrank.ErrVRankNotSorted
	}
	if hasDuplicate(report) {
		return vrank.ErrDuplicateVRankCandidate
	}
	return nil
}

func isNonCandContained(report, candidates []common.Address) bool {
	candSet := make(map[common.Address]struct{}, len(candidates))
	for _, c := range candidates {
		candSet[c] = struct{}{}
	}
	for _, addr := range report {
		if _, ok := candSet[addr]; !ok {
			return true
		}
	}
	return false
}

func isSorted(addrs []common.Address) bool {
	sorted := append([]common.Address(nil), addrs...)
	slices.SortFunc(sorted, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	return slices.Equal(sorted, addrs)
}

func hasDuplicate(addrs []common.Address) bool {
	seen := make(map[common.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		if _, ok := seen[addr]; ok {
			return true
		}
		seen[addr] = struct{}{}
	}
	return false
}
