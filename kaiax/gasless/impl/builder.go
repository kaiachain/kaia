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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/work/builder"
)

var _ kaiax.TxBundlingModule = (*GaslessModule)(nil)

func (g *GaslessModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle {
	// there are only at most two gasless transactions in pending for a sender
	bundles := []*builder.Bundle{}
	approveTxs := map[common.Address]*types.Transaction{}
	targetTxHash := common.Hash{}
	for _, tx := range txs {
		addr, err := types.Sender(g.signer, tx)
		if err != nil {
			continue
		}
		if g.IsApproveTx(tx) {
			approveTxs[addr] = tx
		} else if g.IsSwapTx(tx) && g.IsExecutable(approveTxs[addr], tx) {
			b := builder.NewBundle(
				builder.NewTxOrGenList(g.GetLendTxGenerator(approveTxs[addr], tx)),
				targetTxHash,
				false,
			)

			if approveTxs[addr] != nil {
				b.BundleTxs = append(b.BundleTxs, builder.NewTxOrGenFromTx(approveTxs[addr]))
			}
			b.BundleTxs = append(b.BundleTxs, builder.NewTxOrGenFromTx(tx))

			targetTxHash = tx.Hash()

			isConflict := false
			for _, prev := range append(prevBundles, bundles...) {
				if prev.IsConflict(b) {
					isConflict = true
					break
				}
			}
			if isConflict {
				// Gasless transactions will just fail even if they aren't bundled.
				continue
			}
			bundles = append(bundles, b)
		} else {
			targetTxHash = tx.Hash()
		}
	}
	return bundles
}

func (g *GaslessModule) IsBundleTx(tx *types.Transaction) bool {
	return g.IsModuleTx(tx)
}

func (g *GaslessModule) GetMaxBundleTxsInPending() uint {
	return g.GaslessConfig.MaxBundleTxsInPending
}

func (g *GaslessModule) GetMaxBundleTxsInQueue() uint {
	return g.GaslessConfig.MaxBundleTxsInQueue
}

func (g *GaslessModule) FilterTxs(txs map[common.Address]types.Transactions) {
	// do nothing
}
