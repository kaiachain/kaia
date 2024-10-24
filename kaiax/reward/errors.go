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

package reward

import (
	"errors"
	"fmt"
)

var (
	ErrInitUnexpectedNil     = errors.New("unexpected nil during module init")
	ErrTxReceiptsLenMismatch = errors.New("txs and receipts length mismatch")
	ErrNoBlock               = errors.New("block not found")
	ErrNoReceipts            = errors.New("receipts not found")
	ErrInvalidBlockRange     = errors.New("invalid block number range")
	ErrBlockRangeLimit       = errors.New("exceeds block number range limit")
)

func errMalformedRewardRatio(ratio string) error {
	return fmt.Errorf("malformed reward.ratio: %s", ratio)
}

func errMalformedRewardKip82Ratio(ratio string) error {
	return fmt.Errorf("malformed reward.kip82ratio: %s", ratio)
}
