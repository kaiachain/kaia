// Copyright 2026 The Kaia Authors
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
