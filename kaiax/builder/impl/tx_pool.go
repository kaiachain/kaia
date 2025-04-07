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
	"errors"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

var _ kaiax.TxPoolModule = (*BuilderModule)(nil)

func (b *BuilderModule) PreAddTx(tx *types.Transaction, local bool) error {
	txTime, ok := b.knownTxs[tx.Hash()]
	if ok && time.Since(txTime.time) < KnownTxTimeout {
		return errors.New("Unable to add known bundle tx into tx pool during lock period")
	}
	for _, module := range b.Modules {
		if module.IsBundleTx(tx) {
			newTxTime := txAndTime{
				tx:   tx,
				time: time.Now(),
			}
			if ok {
				newTxTime.time = txTime.time
			}
			b.knownTxs[tx.Hash()] = newTxTime
			break
		}
	}
	return nil
}

func (b *BuilderModule) IsModuleTx(tx *types.Transaction) bool {
	if _, ok := b.knownTxs[tx.Hash()]; ok {
		return true
	}
	for _, module := range b.Modules {
		if module.IsBundleTx(tx) {
			return true
		}
	}
	return false
}

func (b *BuilderModule) GetCheckBalance() func(tx *types.Transaction) error {
	return nil
}

func (b *BuilderModule) IsReady(txs map[uint64]*types.Transaction, next uint64, ready types.Transactions) bool {
	return true
}

func (b *BuilderModule) PreReset(oldHead, newHead *types.Header) {
	for hash, txAndTime := range b.knownTxs {
		// remove pending timed out tx from tx pool
		if time.Since(txAndTime.time) > PendingTimeout {
			b.knownTxs[hash].tx.MarkUnexecutable(true)
		}
		// remove known timed out tx from knownTxs
		if time.Since(txAndTime.time) > KnownTxTimeout {
			delete(b.knownTxs, hash)
		}
	}
}
