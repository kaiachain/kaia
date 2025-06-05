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

	"github.com/kaiachain/kaia/v2/api"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/builder"
	"github.com/kaiachain/kaia/v2/log"
)

var (
	_      builder.BuilderModule = (*BuilderModule)(nil)
	logger                       = log.NewModuleLogger(log.KaiaxBuilder)
)

type InitOpts struct {
	Backend api.Backend
}

type knownTxs map[common.Hash]knownTx

func (k knownTxs) add(tx *types.Transaction) {
	k[tx.Hash()] = knownTx{
		tx:        tx,
		time:      time.Now(),
		isDemoted: false,
	}
}

func (k knownTxs) get(hash common.Hash) (knownTx, bool) {
	tx, ok := k[hash]
	return tx, ok
}

func (k knownTxs) has(hash common.Hash) bool {
	_, ok := k[hash]
	return ok
}

func (k knownTxs) delete(hash common.Hash) {
	delete(k, hash)
}

func (k knownTxs) numExecutable() int {
	num := 0
	for _, knownTx := range k {
		if !knownTx.tx.IsMarkedUnexecutable() && !knownTx.isDemoted {
			num++
		}
	}
	return num
}

func (k knownTxs) markUnexecutable(hash common.Hash) {
	if tx, ok := k[hash]; ok {
		tx.tx.MarkUnexecutable(true)
	}
}

func (k knownTxs) markDemoted(hash common.Hash) {
	if tx, ok := k[hash]; ok {
		k[hash] = knownTx{
			tx:        tx.tx,
			time:      tx.time,
			isDemoted: true,
		}
	}
}

type knownTx struct {
	tx        *types.Transaction
	time      time.Time
	isDemoted bool
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
