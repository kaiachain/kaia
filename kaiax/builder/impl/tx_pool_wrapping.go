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
	mu               sync.RWMutex
}

func NewBuilderWrappingModule(txBundlingModule builder.TxBundlingModule) *BuilderWrappingModule {
	txPoolModule, _ := txBundlingModule.(kaiax.TxPoolModule)
	return &BuilderWrappingModule{
		txBundlingModule: txBundlingModule,
		txPoolModule:     txPoolModule,
		knownTxs:         make(map[common.Hash]txAndTime),
		mu:               sync.RWMutex{},
	}
}

// PreAddTx is a wrapper function that calls the PreAddTx method of the underlying module.
func (b *BuilderWrappingModule) PreAddTx(tx *types.Transaction, local bool) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	txTime, ok := b.knownTxs[tx.Hash()]
	if ok && time.Since(txTime.time) < KnownTxTimeout {
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

	tx, isReady := txs[next]
	if b.txPoolModule != nil {
		isReady = b.txPoolModule.IsReady(txs, next, ready)
	}
	if isReady && b.txBundlingModule.IsBundleTx(tx) {
		if _, ok := b.knownTxs[tx.Hash()]; !ok {
			// if maxBundleSize is max uint64, it means no limit
			if maxBundleNum := b.txBundlingModule.GetMaxBundleNum(); maxBundleNum != math.MaxUint64 {
				numExecutable := uint(0)
				for _, tx := range b.knownTxs {
					if !tx.tx.IsMarkedUnexecutable() {
						numExecutable++
					}
				}
				// it's too much cost to check the size of the bundle here, so we just check the number of txs
				// if the number of txs is greater than the max bundle size, we reject the tx
				// but if there are ready txs, we should add the tx to the pool to complete the bundle
				if numExecutable >= maxBundleNum && len(ready) == 0 {
					logger.Info("Pending tx pool is full of bundle txs, rejecting tx", "maxBundleNum", maxBundleNum)
					return false
				}
			}
			b.knownTxs[tx.Hash()] = txAndTime{
				tx:   tx,
				time: time.Now(),
			}
		}
	}
	return isReady
}

// PreReset is a wrapper function that removes timed out tx from the tx pool and knownTxs.
func (b *BuilderWrappingModule) PreReset(oldHead, newHead *types.Header) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for hash, txAndTime := range b.knownTxs {
		// remove pending timed out tx from tx pool
		if time.Since(txAndTime.time) >= PendingTimeout {
			b.knownTxs[hash].tx.MarkUnexecutable(true)
		}
		// remove known timed out tx from knownTxs
		if time.Since(txAndTime.time) >= KnownTxTimeout {
			delete(b.knownTxs, hash)
		}
	}
	if b.txPoolModule != nil {
		b.txPoolModule.PreReset(oldHead, newHead)
	}
}

// PostReset is a wrapper function that calls the PostReset method of the underlying module.
func (b *BuilderWrappingModule) PostReset(oldHead, newHead *types.Header) {}
