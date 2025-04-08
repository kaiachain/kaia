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
	"github.com/kaiachain/kaia/kaiax/builder"
)

var _ kaiax.TxPoolModule = (*BuilderWrappingModule)(nil)

type BuilderWrappingModule struct {
	builder *BuilderModule
	module  builder.TxBundlingModule
}

func NewBuilderWrappingModule(builder *BuilderModule, module builder.TxBundlingModule) *BuilderWrappingModule {
	return &BuilderWrappingModule{
		builder: builder,
		module:  module,
	}
}

func (b *BuilderWrappingModule) PreAddTx(tx *types.Transaction, local bool) error {
	txTime, ok := b.builder.knownTxs[tx.Hash()]
	if ok && time.Since(txTime.time) < KnownTxTimeout {
		return errors.New("Unable to add known bundle tx into tx pool during lock period")
	}
	if mTxpool, ok := b.module.(kaiax.TxPoolModule); ok {
		return mTxpool.PreAddTx(tx, local)
	}
	return nil
}

func (b *BuilderWrappingModule) IsModuleTx(tx *types.Transaction) bool {
	if mTxpool, ok := b.module.(kaiax.TxPoolModule); ok {
		return mTxpool.IsModuleTx(tx)
	}
	return b.module.IsBundleTx(tx)
}

func (b *BuilderWrappingModule) GetCheckBalance() func(tx *types.Transaction) error {
	if mTxpool, ok := b.module.(kaiax.TxPoolModule); ok {
		return mTxpool.GetCheckBalance()
	}
	return nil
}

func (b *BuilderWrappingModule) IsReady(txs map[uint64]*types.Transaction, next uint64, ready types.Transactions) bool {
	isReady := true
	if mTxpool, ok := b.module.(kaiax.TxPoolModule); ok {
		isReady = mTxpool.IsReady(txs, next, ready)
	}
	if isReady {
		tx := ready[next]
		if b.module.IsBundleTx(tx) {
			newTxTime := txAndTime{
				tx:   tx,
				time: time.Now(),
			}
			txTime, ok := b.builder.knownTxs[tx.Hash()]
			if ok {
				newTxTime.time = txTime.time
			}
			b.builder.knownTxs[tx.Hash()] = newTxTime
		}
	}
	return isReady
}

func (b *BuilderWrappingModule) PreReset(oldHead, newHead *types.Header) {
	if mTxpool, ok := b.module.(kaiax.TxPoolModule); ok {
		mTxpool.PreReset(oldHead, newHead)
	}
}
