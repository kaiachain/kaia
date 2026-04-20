// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package istanbul

import (
	"errors"
)

var (
	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("started engine")
	// errInvalidProposal is returned when a prposal is malformed.
	ErrInvalidProposal = errors.New("invalid proposal")
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	ErrInvalidSignature = errors.New("invalid signature")
	// errNoValidator is returned when the validator is not set.
	ErrNoValidator = errors.New("no validator")
	// errNoEssentialModule is returned when essential module is not registered.
	ErrNoEssentialModule = errors.New("no essential module")
	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	ErrInvalidExtraDataFormat = errors.New("invalid extra data format")
	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	ErrInvalidVotingChain = errors.New("invalid voting chain")
	// ErrUnknownBlock is returned when a required block/header context is missing.
	ErrUnknownBlock = errors.New("unknown block")
	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidCommittedSeals = errors.New("invalid committed seals")
	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	ErrMismatchTxhashes = errors.New("mismatch transactions hashes")
	// errInternalError is returned when an internal error occurs.
	ErrInternalError = errors.New("internal error")
	// errPendingNotAllowed is returned when pending block is not allowed.
	ErrPendingNotAllowed = errors.New("pending is not allowed")
	// errNoBlobSidecarForBlobTx is returned if the blob sidecar is not found for a blob transaction.
	ErrNoBlobSidecarForBlobTx = errors.New("no blob sidecar for blob transaction")
	// errInvalidBlobTxWithSidecar is returned if the blob transaction has an invalid sidecar.
	ErrInvalidBlobTxWithSidecar = errors.New("invalid blob transaction with sidecar")
)
