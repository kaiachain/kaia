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

	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/log"
)

var (
	_      builder.BuilderModule = (*BuilderModule)(nil)
	logger                       = log.NewModuleLogger(log.KaiaxBuilder)
)

type InitOpts struct {
	Backend api.Backend
}

type knownTxs map[common.Hash]*knownTx

func (k knownTxs) add(tx *types.Transaction, status int) {
	if tx == nil {
		return
	}

	if ktx, ok := k[tx.Hash()]; ok {
		ktx.status = status
	} else {
		k[tx.Hash()] = &knownTx{
			tx:     tx,
			time:   time.Now(),
			status: status,
		}
	}
	updateMetrics(&k)
}

func (k knownTxs) addKnownTx(knownTx *knownTx) {
	if ktx, ok := k[knownTx.tx.Hash()]; ok {
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
		if oldestTime < knownTx.elapsedTime().Seconds() {
			oldestTime = knownTx.elapsedTime().Seconds()
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
	tx   *types.Transaction
	time time.Time

	// The location in txpool, of this knownTx.
	// Refreshed at pool.reset()-PostReset(), pool.addTx()-PreAddTx(), pool.promoteExecutables()-IsReady()
	status int
}

func (t *knownTx) elapsedTime() time.Duration {
	return time.Since(t.time)
}

type BuilderModule struct {
	InitOpts
}

func NewBuilderModule() *BuilderModule {
	return &BuilderModule{}
}

func (b *BuilderModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Backend == nil {
		return ErrInitUnexpectedNil
	}
	b.InitOpts = *opts
	return nil
}

func (b *BuilderModule) Start() error {
	return nil
}

func (b *BuilderModule) Stop() {
}

func updateMetrics(knownTxs *knownTxs) {
	numQueueGauge.Update(int64(knownTxs.numQueue()))
	numPendingGauge.Update(int64(knownTxs.numPending()))
	numExecutableGauge.Update(int64(knownTxs.numExecutable()))
	oldestTxTimeInKnownTxsGauge.Update(knownTxs.getTimeOfOldestKnownTx())
}
