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
	"github.com/kaiachain/kaia/kaiax"
)

var _ kaiax.TxPoolModule = (*BuilderModule)(nil)

func (b *BuilderModule) PreAddLocal(tx *types.Transaction) error {
	return nil
}

func (b *BuilderModule) PreAddRemote(tx *types.Transaction) error {
	return nil
}

func (b *BuilderModule) IsModuleTx(tx *types.Transaction) bool {
	return false
}

func (b *BuilderModule) GetCheckBalance() func(tx *types.Transaction) error {
	return nil
}

func (b *BuilderModule) IsReady(txs map[uint64]*types.Transaction, next uint64, ready types.Transactions) bool {
	return false
}

func (b *BuilderModule) Reset(pool kaiax.TxPoolForCaller, oldHead, newHead *types.Header) {
	pending, err := pool.Pending()
	if err != nil {
		logger.Error("Failed to get pending transactions", "err", err)
		return
	}

	signer := types.LatestSignerForChainID(b.Backend.ChainConfig().ChainID)
	baseFee := b.Backend.CurrentBlock().Header().BaseFee
	txs := types.NewTransactionsByPriceAndNonce(signer, pending, baseFee)
	arrayTxs := Arrayify(txs)
	_, bundles := ExtractBundlesAndIncorporate(arrayTxs, b.Modules)

	for _, bundle := range bundles {
		for _, txOrGen := range bundle.BundleTxs {
			if txOrGen.IsConcreteTx() {
				tx, err := txOrGen.GetTx(0)
				if err != nil {
					logger.Error("Failed to get tx from bundle", "err", err)
					continue
				}
				if _, ok := b.bundles[tx]; !ok {
					b.bundles[tx] = time.Now()
				}
			}
		}
	}

	for tx, tm := range b.bundles {
		if time.Since(tm) > BundleTxTimeout {
			tx.MarkUnexecutable(true)
			delete(b.bundles, tx)
		}
	}
}
