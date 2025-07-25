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
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/work/builder"
)

var _ kaiax.TxBundlingModule = (*AuctionModule)(nil)

func (a *AuctionModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle {
	bundles := []*builder.Bundle{}
	curBlock := a.Chain.CurrentBlock()
	if curBlock == nil || atomic.LoadUint32(&a.bidPool.running) == 0 {
		return bundles
	}

	miningBlock := curBlock.NumberU64() + 1
	bidTargetMap := a.bidPool.GetTargetTxMap(miningBlock)
	if bidTargetMap == nil {
		return bundles
	}

	for _, tx := range txs {
		txHash := tx.Hash()
		bid, ok := bidTargetMap[txHash]
		if !ok {
			continue
		}
		b := builder.NewBundle(
			builder.NewTxOrGenList(a.GetBidTxGenerator(tx, bid)),
			txHash,
			true,
		)

		isConflict := false
		for _, prev := range append(prevBundles, bundles...) {
			if prev.IsConflict(b) {
				isConflict = true
				break
			}
		}
		if isConflict {
			continue
		}
		bundles = append(bundles, b)
	}

	return bundles
}

func (g *AuctionModule) IsBundleTx(tx *types.Transaction) bool {
	return false
}

func (g *AuctionModule) GetMaxBundleTxsInPending() uint {
	return 0
}

func (g *AuctionModule) GetMaxBundleTxsInQueue() uint {
	return 0
}

func (a *AuctionModule) FilterTxs(txs map[common.Address]types.Transactions) {
	if atomic.LoadUint32(&a.bidPool.running) == 0 {
		return
	}

	curBlock := a.Chain.CurrentBlock()
	if curBlock == nil {
		return
	}
	targetTxHashMap := a.bidPool.getTargetTxHashMap(curBlock.NumberU64() + 1)

	now := time.Now()
	deadline := now.Add(-AuctionEarlyDeadlineOffset)
	// filter txs that are after the auction early deadline
	for addr, list := range txs {
		for i, tx := range list {
			if tx.Time().After(deadline) && !a.isGaslessTx(tx) {
				// if the tx is a target tx, skip it
				if _, ok := targetTxHashMap[tx.Hash()]; ok {
					continue
				}

				if i == 0 {
					// if first transaction exceeds deadline, remove the address
					delete(txs, addr)
				} else {
					// keep only transactions before the deadline
					txs[addr] = list[:i]
				}
				break
			}
		}
	}
}

func (a *AuctionModule) isGaslessTx(tx *types.Transaction) bool {
	if a.gaslessModule == nil {
		return false
	}

	return a.gaslessModule.IsBundleTx(tx)
}
