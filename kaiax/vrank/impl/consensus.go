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
//   - Before the permissionless fork: VRank must be empty.
//   - At epoch-start blocks: VRank must be RLPEncode(CandTesting(N)).
//   - Otherwise: VRank must be a valid encoded report whose addresses are sorted,
//     deduplicated, and all present in GetCandTesting(N-1), because header(N).VRank
//     reports candidate failures observed while building block N-1.
func (v *VRankModule) VerifyHeader(header *types.Header, _ *types.Header) error {
	number := header.Number.Uint64()
	permissionless := v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(number))

	if !permissionless {
		if len(header.VRank) > 0 {
			return vrank.ErrUnexpectedVRankBeforePermissionless
		}
		return nil
	}

	if number%v.vrankEpoch() != 0 && len(header.VRank) == 0 {
		return nil
	}

	report, err := vrank.DecodeReport(header.VRank)
	if err != nil {
		return vrank.ErrInvalidVRankFormat
	}

	if number%v.vrankEpoch() == 0 {
		candidates, err := v.Valset.GetCandTesting(number)
		if err != nil {
			return err
		}
		return validateEpochVRank(report, candidates)
	} else {
		candidates, err := v.Valset.GetCandTesting(number - 1)
		if err != nil {
			return err
		}
		return validateNonEpochVRank(report, candidates)
	}
}

func (v *VRankModule) PrepareHeader(*types.Header) error { return nil }

func validateEpochVRank(report, candidates []common.Address) error {
	if !slices.Equal(report, candidates) {
		return vrank.ErrMismatchedEpochStartVRank
	}
	return nil
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
