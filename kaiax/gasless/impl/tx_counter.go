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
	"github.com/kaiachain/kaia/common"
	"github.com/rcrowley/go-metrics"
)

var (
	// Metrics for knownTxs
	numQueueGauge   = metrics.NewRegisteredGauge("txpool/knowntxs/num/queue", nil)
	numPendingGauge = metrics.NewRegisteredGauge("txpool/knowntxs/num/pending", nil)
	// numExecutable = numPending - MarkedUnexecutable (by the local miner)
	numExecutableGauge          = metrics.NewRegisteredGauge("txpool/knowntxs/num/executable", nil)
	oldestTxTimeInKnownTxsGauge = metrics.NewRegisteredGauge("txpool/knowntxs/oldesttime/seconds", nil)
)

func updateMetrics(knownTxs *knownTxs) {
	numQueueGauge.Update(int64(knownTxs.numQueue()))
	numPendingGauge.Update(int64(knownTxs.numPending()))
	numExecutableGauge.Update(int64(knownTxs.numExecutable()))
	oldestTxTimeInKnownTxsGauge.Update(knownTxs.getTimeOfOldestKnownTx())
}

type knownTxs map[common.Hash]*knownTx

func (k knownTxs) add(tx *types.Transaction, status int) {
	if tx == nil {
		return
	}

	if ktx, ok := k.get(tx.Hash()); ok {
		ktx.status = status
	} else {
		k[tx.Hash()] = &knownTx{
			tx:           tx,
			addedTime:    time.Time{},
			promotedTime: time.Time{},
			status:       status,
		}
	}

	if status == TxStatusQueue {
		k[tx.Hash()].startAddedTimeIfZero()
	} else if status == TxStatusPending {
		k[tx.Hash()].startPromotedTimeIfZero()
	}

	updateMetrics(&k)
}

func (k knownTxs) addKnownTx(knownTx *knownTx) {
	if ktx, ok := k.get(knownTx.tx.Hash()); ok {
		ktx.status = knownTx.status
	} else {
		k[knownTx.tx.Hash()] = knownTx
	}
	updateMetrics(&k)
}

func (k knownTxs) get(hash common.Hash) (*knownTx, bool) {
	tx, ok := k[hash]
	return tx, ok
}

func (k knownTxs) has(hash common.Hash) bool {
	_, ok := k[hash]
	return ok
}

func (k knownTxs) delete(hash common.Hash) {
	delete(k, hash)
	updateMetrics(&k)
}

func (k knownTxs) numPending() int {
	num := 0
	for _, knownTx := range k {
		if knownTx.status == TxStatusPending {
			num++
		}
	}
	return num
}

func (k knownTxs) numExecutable() int {
	num := 0
	for _, knownTx := range k {
		if knownTx.status == TxStatusPending && !knownTx.tx.IsMarkedUnexecutable() {
			num++
		}
	}
	return num
}

func (k knownTxs) numQueue() int {
	num := 0
	for _, knownTx := range k {
		if knownTx.status == TxStatusQueue {
			num++
		}
	}
	return num
}

func (k knownTxs) getTimeOfOldestKnownTx() int64 {
	var oldestTime float64 = 0
	for _, knownTx := range k {
		if oldestTime < knownTx.elapsedAddedTime().Seconds() {
			oldestTime = knownTx.elapsedAddedTime().Seconds()
		}
	}
	return int64(oldestTime)
}

func (k knownTxs) Copy() *knownTxs {
	newMap := &knownTxs{}
	for _, knownTx := range k {
		newMap.addKnownTx(knownTx)
	}
	return newMap
}

const (
	TxStatusQueue   = iota // exists in txpool.queue
	TxStatusPending        // exists in txpool.pending
	TxStatusDemoted        // not exist in txpool.pending and txpool.queue
)

// A metadata of a known bundle tx that has been submitted to txpool during the last window (KnownTxTimeout).
type knownTx struct {
	tx           *types.Transaction
	addedTime    time.Time
	promotedTime time.Time

	// The location in txpool, of this knownTx.
	// Refreshed at pool.reset()-PostReset(), pool.addTx()-PreAddTx(), pool.promoteExecutables()-IsReady()
	status int
}

func (t *knownTx) elapsedAddedTime() time.Duration {
	return time.Since(t.addedTime)
}

func (t *knownTx) elapsedPromotedTime() time.Duration {
	return time.Since(t.promotedTime)
}

func (t *knownTx) elapsedPromotedOrAddedTime() time.Duration {
	if t.promotedTime.IsZero() {
		return t.elapsedAddedTime()
	}
	return t.elapsedPromotedTime()
}

func (t *knownTx) startAddedTimeIfZero() time.Time {
	if t.addedTime.IsZero() {
		t.addedTime = time.Now()
	}
	return t.addedTime
}

func (t *knownTx) startPromotedTimeIfZero() time.Time {
	t.startAddedTimeIfZero()
	if t.promotedTime.IsZero() {
		t.promotedTime = time.Now()
	}
	return t.promotedTime
}
