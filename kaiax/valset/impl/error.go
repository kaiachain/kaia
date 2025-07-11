// Copyright 2024 The Kaia Authors
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
	"errors"
	"fmt"
)

var (
	errInitUnexpectedNil      = errors.New("unexpected nil during module init")
	errInvalidProposerPolicy  = errors.New("invalid proposer policy")
	errNoNextDistinctProposer = errors.New("cannot determine next distinct proposer")
	errNoHeader               = errors.New("no header found")
	errNoBlock                = errors.New("no block found")
	errNoLowestScannedNum     = errors.New("no lowest scanned validator vote num")
	errNoVoteBlockNums        = errors.New("no validator vote block nums")

	// rpc related errors
	errPendingNotAllowed       = errors.New("pending is not allowed")
	errInternalError           = errors.New("internal error")
	errStartNotPositive        = errors.New("start block number should be positive")
	errEndLargetThanLatest     = errors.New("end block number should be smaller than the latest block number")
	errStartLargerThanEnd      = errors.New("start should be smaller than end")
	errRequestedBlocksTooLarge = errors.New("number of requested blocks should be smaller than 50")
	errRangeNil                = errors.New("range values should not be nil")
	errNoBlockNumber           = errors.New("block number is not assigned")
	errUnknownBlock            = errors.New("unknown block")
	errUnknownProposer         = errors.New("unknown proposer")
	errNoChainReader           = errors.New("sb.chain is nil! --mine option might be missing")
)

func ErrNoIstanbulSnapshot(num uint64) error {
	return fmt.Errorf("no istanbul snapshot at block %d", num)
}
