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
	ErrInitUnexpectedNil    = errors.New("unexpected nil during module init")
	ErrInvalidBlockNumber   = errors.New("invalid block number")
	ErrInvalidSearcherSig   = errors.New("invalid searcher sig")
	ErrInvalidAuctioneerSig = errors.New("invalid auctioneer sig")

	ErrBidAlreadyExists        = errors.New("bid already exists")
	ErrBidSenderExists         = errors.New("bid sender already exists")
	ErrBidInvalidSearcherSig   = errors.New("invalid searcher sig")
	ErrBidInvalidAuctioneerSig = errors.New("invalid auctioneer sig")
	ErrLowBid                  = errors.New("low bid")
	ErrZeroBid                 = errors.New("zero bid")

	ErrAuctionPaused = errors.New("auction is paused")
)
