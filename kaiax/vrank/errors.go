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

const (
	// Epoch is the number of blocks in one VRank scoring epoch.
	Epoch = uint64(86400)
	// MaxRound is the maximum allowed consensus round per block (range [0, MaxRound]).
	MaxRound = 10
)

var (
	ErrInitUnexpectedNil     = errors.New("unexpected nil during module init")
	ErrVRankPreprepareNil    = errors.New("VRankPreprepare is nil")
	ErrVRankCandidateNil     = errors.New("VRankCandidate is nil")
	ErrGetCandidateFailed    = errors.New("valset.GetCandidates failed")
ErrPrepreparedViewNotSet = errors.New("preprepared view is not set")
	ErrViewMismatch          = errors.New("view mismatch")
	ErrBlockHashMismatch     = errors.New("block hash mismatch")
	ErrMsgFromNonCandidate   = errors.New("message from non-candidate")
	ErrMsgFromNonProposer    = errors.New("message from non-proposer")
	ErrInvalidCandidateSig   = errors.New("invalid candidate signature")
	ErrInvalidProposerSig    = errors.New("invalid proposer signature")
	ErrTooFar                = errors.New("too far in the future")
	ErrRoundOutOfRange       = errors.New("round out of range")
	ErrHeaderNotFound        = errors.New("header not found")
	ErrHeaderExtraTooShort   = errors.New("header extra is too short")
	ErrFutureBlock           = errors.New("block number is beyond the current chain head")
	ErrNotPermissionless     = errors.New("block is before the permissionless fork")
)
