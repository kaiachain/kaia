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

package supply

import (
	"errors"
	"fmt"
)

var (
	ErrInitUnexpectedNil = errors.New("unexpected nil during module init")
	ErrNoBlock           = errors.New("block not found")
	ErrNoRebalanceMemo   = errors.New("rebalance memo empty")
	ErrSupplyModuleQuit  = errors.New("supply module quit")
	ErrNoCheckpoint      = errors.New("supply checkpoint not found")
)

func ErrNoCanonicalBurn(err error) error {
	return fmt.Errorf("cannot determine canonical (0x0, 0xdead) burn amount: %w", err)
}

func ErrNoRebalanceBurn(err error) error {
	return fmt.Errorf("cannot determine rebalance (kip103, kip160) burn amount: %w", err)
}
