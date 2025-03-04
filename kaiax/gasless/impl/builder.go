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

package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
)

var _ builder.TxBundlingModule = (*GaslessModule)(nil)

func (g *GaslessModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle {
	// there are only at most two gasless transactions in pending for a sender
	bundles := []*builder.Bundle{}
	approveTxs := map[common.Address]*types.Transaction{}
	targetTxHash := common.Hash{}
	for _, tx := range txs {
		addr := tx.ValidatedSender()
		if g.IsApproveTx(tx) {
			approveTxs[addr] = tx
		} else if g.IsSwapTx(tx) && g.IsExecutable(approveTxs[addr], tx) {
			b := &builder.Bundle{
				BundleTxs: []interface{}{g.GetLendTxGenerator(approveTxs[addr], tx)},
			}
			if approveTxs[addr] != nil {
				b.BundleTxs = append(b.BundleTxs, approveTxs[addr])
			}
			b.TargetTxHash = targetTxHash
			b.BundleTxs = append(b.BundleTxs, tx)
			conflict := false
			for _, prev := range prevBundles {
				conflict = prev.IsConflict(b)
				if conflict {
					break
				}
			}
			if !conflict {
				bundles = append(bundles, b)
			}
		} else {
			targetTxHash = tx.Hash()
		}
	}
	return bundles
}
