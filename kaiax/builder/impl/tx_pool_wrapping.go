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
	"github.com/rcrowley/go-metrics"
)

var _ kaiax.TxPoolModule = (*BuilderWrappingModule)(nil)

var (
	// Metrics for knownTxs
	numQueueGauge   = metrics.NewRegisteredGauge("txpool/knowntxs/num/queue", nil)
	numPendingGauge = metrics.NewRegisteredGauge("txpool/knowntxs/num/pending", nil)
	// numExecutable = numPending - MarkedUnexecutable (by the local miner)
	numExecutableGauge          = metrics.NewRegisteredGauge("txpool/knowntxs/num/executable", nil)
	oldestTxTimeInKnownTxsGauge = metrics.NewRegisteredGauge("txpool/knowntxs/oldesttime/seconds", nil)
)

type BuilderWrappingModule struct {
	txBundlingModule builder.TxBundlingModule
	txPoolModule     kaiax.TxPoolModule // either nil or same object as txBundlingModule
	knownTxs         *knownTxs
	mu               sync.RWMutex
}

func NewBuilderWrappingModule(txBundlingModule builder.TxBundlingModule) *BuilderWrappingModule {
	txPoolModule, _ := txBundlingModule.(kaiax.TxPoolModule)
	return &BuilderWrappingModule{
		txBundlingModule: txBundlingModule,
		txPoolModule:     txPoolModule,
		knownTxs:         &knownTxs{},
		mu:               sync.RWMutex{},
	}
}

// PreAddTx is a wrapper function that calls the PreAddTx method of the underlying module.
func (b *BuilderWrappingModule) PreAddTx(tx *types.Transaction, local bool) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if knownTx, ok := b.knownTxs.get(tx.Hash()); ok && knownTx.elapsedAddedOrPromotedTime() < KnownTxTimeout {
		return ErrUnableToAddKnownBundleTx
	}

	if b.txPoolModule != nil {
		err := b.txPoolModule.PreAddTx(tx, local)
		if err != nil {
			return err
		}
	}

	if b.txBundlingModule.IsBundleTx(tx) {
		if b.knownTxs.numQueue() >= int(b.txBundlingModule.GetMaxBundleTxsInQueue()) {
			return ErrBundleTxQueueFull
		}
		b.knownTxs.add(tx, TxStatusQueue)
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
	if b.txBundlingModule.IsBundleTx(tx) {
		// If prev tx is bundle tx, there's no need to check the knownTxs limit because it has been checked in the previous `IsReady()` execution.
		isPrevTxBundleTx := len(ready) != 0 && b.txBundlingModule.IsBundleTx(ready[len(ready)-1])
		if isPrevTxBundleTx {
			b.knownTxs.add(tx, TxStatusPending)
			return true
		}

		maxBundleTxsInPending := b.txBundlingModule.GetMaxBundleTxsInPending()
		if maxBundleTxsInPending != math.MaxUint64 {
			numExecutable := uint(b.knownTxs.numExecutable())

			numSeqTxs := uint(1)
			for i := next + 1; i < next+uint64(len(txs)); i++ {
				if tx, ok := txs[i]; ok && b.txBundlingModule.IsBundleTx(tx) {
					numSeqTxs++
				} else {
					break
				}
			}

			// false if there is possibility of exceeding max bundle tx num
			if numExecutable+numSeqTxs > maxBundleTxsInPending {
				logger.Trace("Not promoting a tx because of exceeding max bundle tx num", "tx", tx.Hash().String(), "numExecutable", numExecutable, "maxBundleTxsInPending", maxBundleTxsInPending)
				return false
			}
		}

		b.knownTxs.add(tx, TxStatusPending)
	}

	return true
}

// PreReset is a wrapper function that removes timed out tx from the tx pool and knownTxs.
func (b *BuilderWrappingModule) PreReset(oldHead, newHead *types.Header) []common.Hash {
	b.mu.Lock()
	defer b.mu.Unlock()

	drops := make([]common.Hash, 0)

	for hash, knownTx := range *b.knownTxs {
		// remove pending timed out tx from tx pool
		if knownTx.status == TxStatusPending && knownTx.elapsedPromotedTime() >= PendingTimeout {
			drops = append(drops, hash)
		}
		// remove queue timed out tx from tx pool
		if knownTx.status == TxStatusQueue && knownTx.elapsedAddedTime() >= QueueTimeout {
			drops = append(drops, hash)
		}
		// remove known timed out tx from knownTxs
		if knownTx.elapsedAddedOrPromotedTime() >= KnownTxTimeout {
			b.knownTxs.delete(hash)
		}
	}
	if b.txPoolModule != nil {
		moduleDrops := b.txPoolModule.PreReset(oldHead, newHead)
		drops = append(drops, moduleDrops...)
	}

	return drops
}

// PostReset is a wrapper function that calls the PostReset method of the underlying module.
func (b *BuilderWrappingModule) PostReset(oldHead, newHead *types.Header, queue, pending map[common.Address]types.Transactions) {
	b.mu.Lock()
	defer b.mu.Unlock()

	flattenedQueue := make(map[common.Hash]*types.Transaction)
	flattenedPending := make(map[common.Hash]*types.Transaction)
	for _, txs := range queue {
		for _, tx := range txs {
			flattenedQueue[tx.Hash()] = tx
		}
	}
	for _, txs := range pending {
		for _, tx := range txs {
			flattenedPending[tx.Hash()] = tx
		}
	}

	for _, knownTx := range *b.knownTxs {
		if _, ok := flattenedQueue[knownTx.tx.Hash()]; ok {
			b.knownTxs.add(knownTx.tx, TxStatusQueue)
		} else if _, ok := flattenedPending[knownTx.tx.Hash()]; ok {
			b.knownTxs.add(knownTx.tx, TxStatusPending)
		} else {
			b.knownTxs.add(knownTx.tx, TxStatusDemoted)
		}
	}

	if b.txPoolModule != nil {
		b.txPoolModule.PostReset(oldHead, newHead, queue, pending)
	}
}
