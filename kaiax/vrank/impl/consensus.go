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
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// VerifyHeader checks the VRank field in the header:
//   - Before the permissionless fork: VRank must be absent.
//   - At epoch-start blocks: Report must equal CandTesting(N).
//   - Otherwise: Report is a cfReport whose failed addresses are sorted, deduped, and ⊆ CandTesting.
//   - Always: ParentRound must be backed by a committee quorum that committed the parent there.
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

	payload, err := vrank.DecodeVRank(header.VRank)
	if err != nil {
		return vrank.ErrInvalidVRankFormat
	}
	if err := v.verifyParentRound(header, payload); err != nil {
		return err
	}

	// Epoch-start block
	if number%v.vrankEpoch() == 0 {
		candidates, err := v.Valset.GetCandTesting(number)
		if err != nil {
			return err
		}
		if !slices.Equal(payload.Report, candidates) {
			return vrank.ErrEpochStartVRankMismatch
		}
		return nil
	}

	// Non-epoch-start block
	if len(payload.Report) == 0 {
		return nil
	}
	// Failures score against the reporter's own byzantine-filterable column regardless of content,
	// so only the failed list is checked (CandTesting is epoch-stable within the epoch).
	candidates, err := v.Valset.GetCandTesting(number)
	if err != nil {
		return err
	}
	return validateNonEpochVRank(payload.Report, candidates)
}

// verifyParentRound checks that the seals prove a committee quorum committed the parent at the
// claimed round. It never compares against this node's own stored round: two honest nodes can hold
// different rounds for one block, so comparing would make them reject each other's headers.
// The seals prove that the round happened, not that it was the last one, so a proposer can claim a
// lower round that really happened but cannot invent a higher one.
func (v *VRankModule) verifyParentRound(header *types.Header, payload vrank.VRankPayload) error {
	if !v.hasParentCertificate(header.Number.Uint64()) || payload.ParentRound == 0 {
		if payload.ParentRound != 0 {
			return vrank.ErrUnexpectedParentRound
		}
		if len(payload.ParentCommittedSeal) > 0 {
			return vrank.ErrUnexpectedParentSeal
		}
		return nil
	}

	parentNum := header.Number.Uint64() - 1
	committee, err := v.Valset.GetCommittee(parentNum, uint64(payload.ParentRound))
	if err != nil {
		return err
	}
	// Bound the recovery work: seals beyond the committee size can never help.
	if len(payload.ParentCommittedSeal) > len(committee) {
		return vrank.ErrTooManyParentSeals
	}
	committers, err := v.Sealer.RecoverCommitters(parentNum, header.ParentHash, payload.ParentRound, payload.ParentCommittedSeal)
	if err != nil {
		return vrank.ErrInvalidParentCertificate
	}
	committeeSet := valset.NewAddressSet(committee)
	valid := 0
	for _, addr := range committers {
		if !committeeSet.Remove(addr) {
			return vrank.ErrInvalidParentCertificate
		}
		valid++
	}
	// The committee equals the qualified set post-permissionless, so pass its size for both.
	if valid < v.Sealer.Quorum(parentNum, len(committee), len(committee)) {
		return vrank.ErrInvalidParentCertificate
	}
	return nil
}

// hasParentCertificate reports whether block num records its parent's round. Two blocks do not:
// an epoch start, whose parent's failures belong to the previous epoch, and the first
// permissionless block, whose parent's seals predate the round binding.
func (v *VRankModule) hasParentCertificate(num uint64) bool {
	if num == 0 || num%v.vrankEpoch() == 0 {
		return false
	}
	return v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(num - 1))
}

// PrepareHeader fills header.VRank per KIP-227.
//
//   - At epoch-start blocks: Report is CandTesting(N).
//   - Otherwise: Report is a cfReport about this proposer's own most recent prior proposal in
//     the current epoch, or empty when there is no such block or no failures.
//   - ParentRound and its seals are copied from this node's own stored parent header.
func (v *VRankModule) PrepareHeader(header *types.Header) error {
	number := header.Number.Uint64()
	if !v.ChainConfig.IsPermissionlessForkEnabled(header.Number) {
		header.VRank = nil
		return nil
	}

	var payload vrank.VRankPayload
	if number%v.vrankEpoch() == 0 {
		candidates, err := v.Valset.GetCandTesting(number)
		if err != nil {
			logger.Error("epoch-start VRank fill: GetCandTesting failed", "num", number, "err", err)
			return err
		}
		payload.Report = candidates
	} else {
		report, err := v.evaluateCandidateFailures(number)
		if err != nil {
			return err
		}
		payload.Report = report
	}

	round, seals, err := v.parentRoundCertificate(number)
	if err != nil {
		return err
	}
	payload.ParentRound = round
	payload.ParentCommittedSeal = seals

	encoded, err := vrank.EncodeVRank(payload)
	if err != nil {
		logger.Error("Failed to encode VRank payload", "err", err, "num", number)
		return err
	}
	header.VRank = encoded
	return nil
}

// parentRoundCertificate reads the round and seals from this node's own stored parent, never from
// seals gathered off the network: the proof of that round is the parent this node accepted.
func (v *VRankModule) parentRoundCertificate(number uint64) (uint8, [][]byte, error) {
	if !v.hasParentCertificate(number) {
		return 0, nil, nil
	}
	parent := v.Chain.GetHeaderByNumber(number - 1)
	if parent == nil {
		return 0, nil, vrank.ErrHeaderNotFound
	}
	round, err := v.Sealer.Round(parent)
	if err != nil {
		return 0, nil, err
	}
	if round == 0 {
		return 0, nil, nil
	}
	_, seals, err := v.Sealer.RawSeals(parent)
	if err != nil {
		return 0, nil, err
	}
	return round, seals, nil
}

// evaluateCandidateFailures reports on this proposer's own most recent prior proposal this epoch.
func (v *VRankModule) evaluateCandidateFailures(number uint64) ([]common.Address, error) {
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
	return report, nil
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
