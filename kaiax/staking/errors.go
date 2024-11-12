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

package staking

import (
	"errors"
	"fmt"
)

var (
	ErrInitUnexpectedNil   = errors.New("unexpected nil during module init")
	ErrZeroStakingInterval = errors.New("staking interval cannot be zero")
	ErrAddressBookResult   = errors.New("invalid result from AddressBook")
)

func ErrAddressBookCall(err error) error {
	return fmt.Errorf("error calling AddressBook: %w", err)
}
