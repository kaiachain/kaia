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

package auction

import "errors"

var (
	ErrInitUnexpectedNil     = errors.New("unexpected nil during module init")
	ErrBlockNotFound         = errors.New("block not found")
	ErrInvalidBlockNumber    = errors.New("invalid block number")
	ErrInvalidSignature      = errors.New("invalid signature")
	ErrInvalidSearcherSig    = errors.New("invalid searcher sig")
	ErrInvalidAuctioneerSig  = errors.New("invalid auctioneer sig")
	ErrNilChainId            = errors.New("chainId is nil")
	ErrNilVerifyingContract  = errors.New("verifyingContract is nil")
	ErrInvalidTargetTxHash   = errors.New("invalid target tx hash")
	ErrAuctionDisabled       = errors.New("auction is disabled")
	ErrExceedMaxCallGasLimit = errors.New("gas limit exceeds the maximum limit")
	ErrExceedMaxDataSize     = errors.New("data size exceeds the maximum limit")

	ErrBidAlreadyExists        = errors.New("bid already exists")
	ErrBidSenderExists         = errors.New("bid sender already exists")
	ErrBidInvalidSearcherSig   = errors.New("invalid searcher sig")
	ErrBidInvalidAuctioneerSig = errors.New("invalid auctioneer sig")
	ErrLowBid                  = errors.New("low bid")
	ErrZeroBid                 = errors.New("zero bid")
	ErrBidPoolFull             = errors.New("bid pool is full")

	ErrAuctionPaused = errors.New("auction is paused")
)
