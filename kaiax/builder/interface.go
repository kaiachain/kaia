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

package builder

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

type TxGenerator func(nonce uint64) (*types.Transaction, error)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/builder BuilderModule
type BuilderModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
}

//go:generate mockgen -destination=./mock/tx_bundling_module.go -package=mock github.com/kaiachain/kaia/kaiax/builder TxBundlingModule
type TxBundlingModule interface {
	// The function finds transactions to be bundled.
	// New transactions can be injected.
	// returned bundles must not have conflict with `prevBundles`.
	// `txs` and `prevBundles` is read-only; it is only to check if there's conflict between new bundles.
	ExtractTxBundles(txs []*types.Transaction, prevBundles []*Bundle) []*Bundle

	// IsBundleTx returns true if the tx is a potential bundle tx.
	IsBundleTx(tx *types.Transaction) bool

	// GetMaxBundleTxsInPending returns the maximum number of transactions that can be bundled in pending.
	// if zero or minus value is returned, it means no limit.
	// this limitation works properly only when a module bundles only sequential txs by the same sender.
	GetMaxBundleTxsInPending() uint

	// GetMaxBundleTxsInPool returns the maximum number of transactions that can be bundled in txpool.
	// if zero or minus value is returned, it means no limit.
	GetMaxBundleTxsInPool() uint
}

// Any component or module that accomodate tx bundling modules.
type TxBundlingModuleHost interface {
	RegisterTxBundlingModule(modules ...TxBundlingModule)
}
