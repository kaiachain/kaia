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
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/builder"
)

var _ kaiax.TxPoolModule = (*BuilderWrappingModule)(nil)

type BuilderWrappingModule struct {
	txBundlingModule builder.TxBundlingModule
	txPoolModule     kaiax.TxPoolModule // either nil or same object as txBundlingModule
	knownTxs         map[common.Hash]txAndTime
}

func NewBuilderWrappingModule(txBundlingModule builder.TxBundlingModule) *BuilderWrappingModule {
	txPoolModule, _ := txBundlingModule.(kaiax.TxPoolModule)
	return &BuilderWrappingModule{
		txBundlingModule: txBundlingModule,
		txPoolModule:     txPoolModule,
		knownTxs:         make(map[common.Hash]txAndTime),
	}
}

// PreAddTx is a wrapper function that calls the PreAddTx method of the underlying module.
func (b *BuilderWrappingModule) PreAddTx(tx *types.Transaction, local bool) error {
	txTime, ok := b.knownTxs[tx.Hash()]
	if ok && time.Since(txTime.time) < KnownTxTimeout {
		return errors.New("Unable to add known bundle tx into tx pool during lock period")
	}
	if b.txPoolModule != nil {
		return b.txPoolModule.PreAddTx(tx, local)
	}
	return nil
}

// IsModuleTx is a wrapper function that calls the IsModuleTx method of the underlying module.
func (b *BuilderWrappingModule) IsModuleTx(tx *types.Transaction) bool {
	if b.txPoolModule != nil {
		return b.txPoolModule.IsModuleTx(tx)
	}
	return b.txBundlingModule.IsBundleTx(tx)
}

// GetCheckBalance is a wrapper function that calls the GetCheckBalance method of the underlying module.
func (b *BuilderWrappingModule) GetCheckBalance() func(tx *types.Transaction) error {
	if b.txPoolModule != nil {
		return b.txPoolModule.GetCheckBalance()
	}
	return nil
}

// IsReady is a wrapper function that checks if the transaction is ready to be added to the tx pool.
func (b *BuilderWrappingModule) IsReady(txs map[uint64]*types.Transaction, next uint64, ready types.Transactions) bool {
	tx, isReady := txs[next]
	if b.txPoolModule != nil {
		isReady = b.txPoolModule.IsReady(txs, next, ready)
	}
	if isReady && b.txBundlingModule.IsBundleTx(tx) {
		newTxTime := txAndTime{
			tx:   tx,
			time: time.Now(),
		}
		txTime, ok := b.knownTxs[tx.Hash()]
		if ok {
			newTxTime.time = txTime.time
		}
		b.knownTxs[tx.Hash()] = newTxTime
	}
	return isReady
}

// PreReset is a wrapper function that removes timed out tx from the tx pool and knownTxs.
func (b *BuilderWrappingModule) PreReset(oldHead, newHead *types.Header) {
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
	if b.txPoolModule != nil {
		b.txPoolModule.PreReset(oldHead, newHead)
	}
}
