// Modifications Copyright 2024 The Kaia Authors

package blockchain

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

type chainTxLookup struct {
	fn func(common.Hash) *types.Transaction
}

// SetTxLookup registers the hash→tx resolver used by sender warming. Nil is ignored.
func (bc *BlockChain) SetTxLookup(fn func(common.Hash) *types.Transaction) {
	if bc == nil || fn == nil {
		return
	}
	bc.txLookup.Store(&chainTxLookup{fn: fn})
}

// TxLookup returns the registered hash→tx resolver, or nil if not set.
func (bc *BlockChain) TxLookup() func(common.Hash) *types.Transaction {
	if bc == nil {
		return nil
	}
	if lookup := bc.txLookup.Load(); lookup != nil {
		return lookup.fn
	}
	return nil
}
