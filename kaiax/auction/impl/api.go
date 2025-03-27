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

package impl

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

func (a *AuctionModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "auction",
			Version:   "1.0",
			Service:   newAuctionAPI(a),
			Public:    false,
		},
	}
}

type AuctionAPI struct {
	a *AuctionModule
}

func newAuctionAPI(a *AuctionModule) *AuctionAPI {
	return &AuctionAPI{a: a}
}

// TODO-kaiax: replace this with a correct implementation
func (a *AuctionAPI) SubmitBid(input hexutil.Bytes) (common.Hash, error) {
	bid := new(auction.Bid)
	err := rlp.DecodeBytes(input, bid)
	if err != nil {
		return common.Hash{}, err
	}

	return a.a.bidPool.AddBid(bid)
}
