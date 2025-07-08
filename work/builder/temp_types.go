// TODO-Kaia: This file is temporarily used during refactoring.
package builder

import (
	"github.com/kaiachain/kaia/blockchain/types"
)

// TxBundlingModule can intervene how miner/proposer orders transactions in a block.
// TODO-Kaia: Move this to kaiax/interface.go
//
//go:generate mockgen -destination=./mock/tx_bundling_module.go -package=mock github.com/kaiachain/kaia/work/builder TxBundlingModule
type TxBundlingModule interface {
	// The function finds transactions to be bundled.
	// New transactions can be injected.
	// returned bundles must not have conflict with `prevBundles`.
	// `txs` and `prevBundles` is read-only; it is only to check if there's conflict between new bundles.
	ExtractTxBundles(txs []*types.Transaction, prevBundles []*Bundle) []*Bundle

	// IsBundleTx returns true if the tx is a potential bundle tx.
	IsBundleTx(tx *types.Transaction) bool

	// GetMaxBundleTxsInPending returns the maximum number of transactions that can be bundled in pending.
	// This limitation works properly only when a module bundles only sequential txs by the same sender.
	GetMaxBundleTxsInPending() uint

	// GetMaxBundleTxsInQueue returns the maximum number of transactions that can be bundled in queue.
	// This limitation works properly only when a module bundles only sequential txs by the same sender.
	GetMaxBundleTxsInQueue() uint
}

// Any component or module that accomodate tx bundling modules.
type TxBundlingModuleHost interface {
	RegisterTxBundlingModule(modules ...TxBundlingModule)
}
