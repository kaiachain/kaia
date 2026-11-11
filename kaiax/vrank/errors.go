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

import "errors"

var (
	ErrInitUnexpectedNil      = errors.New("unexpected nil during module init")
	ErrVRankPreprepareNil     = errors.New("VRankPreprepare is nil")
	ErrVRankCandidateNil      = errors.New("VRankCandidate is nil")
	ErrGetCandidateFailed     = errors.New("valset.GetCandTesting failed")
	ErrViewMismatch           = errors.New("view mismatch")
	ErrBlockHashMismatch      = errors.New("block hash mismatch")
	ErrMsgFromNonProposer     = errors.New("message from non-proposer")
	ErrInvalidCandidateSig    = errors.New("invalid candidate ECDSA signature")
	ErrInvalidCandidateBlsSig = errors.New("invalid candidate BLS signature")
	ErrInvalidProposerSig     = errors.New("invalid proposer signature")
	ErrRoundOutOfRange        = errors.New("round out of range")
	ErrHeaderNotFound         = errors.New("header not found")
	ErrFutureBlock            = errors.New("block number is beyond the current chain head")
	ErrNotPermissionless      = errors.New("block is before the permissionless fork")

	ErrUnexpectedVRankBeforePermissionless = errors.New("unexpected vrank before permissionless fork")
	ErrEpochStartVRankMismatch             = errors.New("vrank at epoch start does not match CandTesting")
	ErrInvalidVRankFormat                  = errors.New("invalid vrank format")
	ErrInvalidVRankCandidate               = errors.New("invalid vrank: address is not a candidate")
	ErrDuplicateVRankCandidate             = errors.New("invalid vrank: duplicate candidate address")
	ErrVRankNotSorted                      = errors.New("invalid vrank: addresses not sorted")
	ErrUnexpectedParentRound               = errors.New("invalid vrank: unexpected parent round")
	ErrUnexpectedParentSeal                = errors.New("invalid vrank: unexpected parent committed seal")
	ErrTooManyParentSeals                  = errors.New("invalid vrank: more parent seals than committee members")
	ErrInvalidParentCertificate            = errors.New("invalid vrank: parent round is not backed by a committee quorum")
)
