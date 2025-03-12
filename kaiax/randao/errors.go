// Copyright 2025 The Kaia Authors
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

package randao

import "errors"

var (
	ErrInitUnexpectedNil   = errors.New("unexpected nil during module init")
	ErrZeroBlockNumber     = errors.New("block number cannot be zero")
	ErrMissingKIP113       = errors.New("kip113 address not set in ChainConfig")
	ErrBeforeRandaoFork    = errors.New("cannot read kip113 address from registry before randao fork")
	ErrNoBlsKey            = errors.New("bls key not configured")
	ErrNoBlsPub            = errors.New("bls pubkey not found for the proposer")
	ErrInvalidRandaoFields = errors.New("invalid randao fields")
	ErrUnexpectedRandao    = errors.New("unexpected randao fields")
)
