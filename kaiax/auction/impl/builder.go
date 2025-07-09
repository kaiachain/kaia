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
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
)

var _ builder.TxBundlingModule = (*AuctionModule)(nil)

func (a *AuctionModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle {
	// TODO: implement me
	return nil
}

func (a *AuctionModule) FilterTxs(txs map[common.Address]types.Transactions) {
	now := time.Now()
	// filter txs that are after the auction early deadline
	for addr, list := range txs {
		temp := list
		for i, tx := range list {
			if tx.Time().Add(AuctionEarlyDeadline).After(now) {
				temp = list[:i]
				break
			}
		}
		if len(temp) > 0 {
			txs[addr] = temp
		}
	}
}
