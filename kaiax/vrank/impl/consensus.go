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

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// VerifyHeader checks the VRank field in the header:
//   - Before the permissionless fork: VRank must be empty.
//   - At epoch-start blocks: VRank must be empty.
//   - Otherwise: VRank must be a valid encoded report whose addresses are sorted,
//     deduplicated, and all present in GetCandidates(calcEpochStart(N)).
func (v *VRankModule) VerifyHeader(header *types.Header) error {
	number := header.Number.Uint64()
	permissionless := v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(number))

	if !permissionless {
		if len(header.VRank) > 0 {
			return vrank.ErrUnexpectedVRankBeforePermissionless
		}
		return nil
	}

	if number%vrank.Epoch == 0 && len(header.VRank) > 0 {
		return vrank.ErrUnexpectedVRankAtEpochStart
	}

	if len(header.VRank) == 0 {
		return nil
	}

	report, err := vrank.DecodeReport(header.VRank)
	if err != nil {
		return vrank.ErrInvalidVRankFormat
	}

	// Validate every address in the cfReport is an actual candidate.
	// header(N).VRank is produced by TallyCfReport(N-1, ...) which uses
	// GetCandidates(calcEpochStart(N)), so mirror that here.
	// This prevents a malicious proposer from injecting arbitrary addresses (e.g. validator
	// addresses) into header.VRank to manipulate CFS scores.
	candidates, err := v.Valset.GetCandidates(calcEpochStart(number))
	if err != nil {
		return err
	}
	candSet := make(map[common.Address]struct{}, len(candidates))
	for _, c := range candidates {
		candSet[c] = struct{}{}
	}
	for i, addr := range report {
		if _, ok := candSet[addr]; !ok {
			return vrank.ErrInvalidVRankCandidate
		}
		if i > 0 {
			cmp := bytes.Compare(report[i-1].Bytes(), addr.Bytes())
			if cmp == 0 {
				return vrank.ErrDuplicateVRankCandidate
			}
			if cmp > 0 {
				return vrank.ErrVRankNotSorted
			}
		}
	}
	return nil
}

func (v *VRankModule) PrepareHeader(header *types.Header) error { return nil }

func (v *VRankModule) FinalizeHeader(header *types.Header, _ *state.StateDB, _ []*types.Transaction, _ []*types.Receipt) error {
	return nil
}
