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
	"sync"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/builder"
)

var _ kaiax.TxPoolModule = (*BuilderWrappingModule)(nil)

type BuilderWrappingModule struct {
	txPool           kaiax.TxPoolForCaller
	txBundlingModule builder.TxBundlingModule
	txPoolModule     kaiax.TxPoolModule // either nil or same object as txBundlingModule
	knownTxs         knownTxs
	mu               sync.RWMutex
}

func NewBuilderWrappingModule(txBundlingModule builder.TxBundlingModule, txPool kaiax.TxPoolForCaller) *BuilderWrappingModule {
	txPoolModule, _ := txBundlingModule.(kaiax.TxPoolModule)
	return &BuilderWrappingModule{
		txPool:           txPool,
		txBundlingModule: txBundlingModule,
		txPoolModule:     txPoolModule,
		knownTxs:         knownTxs{},
		mu:               sync.RWMutex{},
	}
}

// PreAddTx is a wrapper function that calls the PreAddTx method of the underlying module.
func (b *BuilderWrappingModule) PreAddTx(tx *types.Transaction, local bool) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if knownTx, ok := b.knownTxs.get(tx.Hash()); ok && knownTx.elapsedTime() < KnownTxTimeout {
		return ErrUnableToAddKnownBundleTx
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
	b.mu.Lock()
	defer b.mu.Unlock()

	// false if tx is not in txs
	tx, ok := txs[next]
	if !ok {
		return false
	}

	// false if tx is not ready in child module
	if b.txPoolModule != nil && !b.txPoolModule.IsReady(txs, next, ready) {
		return false
	}

	// add tx to knownTxs if it is a bundle tx and not in knownTxs
	if b.txBundlingModule.IsBundleTx(tx) && !b.knownTxs.has(tx.Hash()) {
		// If prev tx is bundle tx, there's no need to check the knownTxs limit because it has been checked in the previous `IsReady()` execution.
		isPrevTxBundleTx := len(ready) != 0 && b.txBundlingModule.IsBundleTx(ready[len(ready)-1])
		if isPrevTxBundleTx {
			b.knownTxs.add(tx)
			return true
		}

		maxBundleTxsInPending := b.txBundlingModule.GetMaxBundleTxsInPending()
		if maxBundleTxsInPending != math.MaxUint64 {
			numExecutable := uint(b.knownTxs.numExecutable())

			numSeqTxs := uint(1)
			for i := next + 1; i < next+uint64(len(txs)); i++ {
				if _, ok := txs[i]; !ok {
					break
				}
				if !b.txBundlingModule.IsBundleTx(txs[i]) {
					break
				}
				numSeqTxs++
			}

			// false if there is possibility of exceeding max bundle tx num
			if numExecutable+numSeqTxs > maxBundleTxsInPending {
				logger.Info("Exceed max bundle tx num", "numExecutable", "maxBundleTxsInPending", maxBundleTxsInPending)
				return false
			}
		}

		b.knownTxs.add(tx)
	}

	return true
}

// PreReset is a wrapper function that removes timed out tx from the tx pool and knownTxs.
func (b *BuilderWrappingModule) PreReset(oldHead, newHead *types.Header) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for hash, knownTx := range b.knownTxs {
		// remove pending timed out tx from tx pool
		if knownTx.elapsedTime() >= PendingTimeout {
			b.knownTxs.markUnexecutable(hash)
		}
		// remove known timed out tx from knownTxs
		if knownTx.elapsedTime() >= KnownTxTimeout {
			b.knownTxs.delete(hash)
		}
	}
	if b.txPoolModule != nil {
		b.txPoolModule.PreReset(oldHead, newHead)
	}
}

// PostReset is a wrapper function that calls the PostReset method of the underlying module.
func (b *BuilderWrappingModule) PostReset(oldHead, newHead *types.Header) {
	b.mu.Lock()
	defer b.mu.Unlock()

	pending, err := b.txPool.PendingUnlocked()
	if err != nil {
		logger.Error("Failed to get pending txs", "error", err)
		return
	}
	flattened := make(map[common.Hash]*types.Transaction)
	for _, txs := range pending {
		for _, tx := range txs {
			flattened[tx.Hash()] = tx
		}
	}
	for hash := range b.knownTxs {
		if _, ok := flattened[hash]; !ok {
			b.knownTxs.markDemoted(hash)
		}
	}
}
