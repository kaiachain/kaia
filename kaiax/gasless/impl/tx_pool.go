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
	"math"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

var _ kaiax.TxPoolModule = (*GaslessModule)(nil)

func (g *GaslessModule) PreAddLocal(tx *types.Transaction) error {
	return nil
}

func (g *GaslessModule) PreAddRemote(tx *types.Transaction) error {
	return nil
}

func (g *GaslessModule) IsModuleTx(tx *types.Transaction) bool {
	return g.IsApproveTx(tx) || g.IsSwapTx(tx)
}

func (g *GaslessModule) GetCheckBalance() func(tx *types.Transaction) error {
	return func(*types.Transaction) error { return nil }
}

func (g *GaslessModule) IsReady(txs map[uint64]*types.Transaction, i uint64, ready types.Transactions) bool {
	tx := txs[i]
	addr := tx.ValidatedSender()
	nonce := g.StateDB.GetNonce(addr)

	if i > nonce+1 {
		return false
	}

	isApproveTx := g.IsApproveTx(tx)
	isSwapTx := g.IsSwapTx(tx)

	if isApproveTx && i == nonce && i+1 < uint64(math.MaxUint64) {
		if next := txs[i+1]; next != nil && g.IsSwapTx(next) {
			return g.IsExecutable(tx, next)
		}
	}

	if isSwapTx {
		if i == nonce {
			return g.IsExecutable(nil, tx)
		}

		if i == nonce+1 && len(ready) > 0 {
			if prev := ready[len(ready)-1]; prev != nil && g.IsApproveTx(prev) {
				return g.IsExecutable(prev, tx)
			}
		}
	}

	return false
}
