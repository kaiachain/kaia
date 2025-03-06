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
	"github.com/kaiachain/kaia/common"
)

type Bundle struct {
	// each element can be either *types.Transaction, or TxGenerator
	BundleTxs []interface{}

	// BundleTxs is placed AFTER the target tx. If empty hash, it is placed at the very front.
	TargetTxHash common.Hash
}

// Has checks if the bundle contains a tx with the given hash.
func (b *Bundle) Has(hash common.Hash) bool {
	return b.FindIdx(hash) != -1
}

// FindIdx returns if the bundle contains a tx with the given hash and its index in bundle.
func (b *Bundle) FindIdx(hash common.Hash) int {
	for i, txOrGen := range b.BundleTxs {
		switch v := txOrGen.(type) {
		case *types.Transaction:
			if v.Hash() == hash {
				return i
			}
		}
	}
	return -1
}

// IsConflict checks if newBundle conflicts with current bundle.
func (b *Bundle) IsConflict(newBundle *Bundle) bool {
	// 1. Check for same target tx hash
	if b.TargetTxHash == newBundle.TargetTxHash {
		return true
	}

	// 2-1. Empty bundleTxs does not conflict with other transactions
	if len(b.BundleTxs) == 0 {
		return false
	}

	// 2-2. Build a map of TxHash -> IndexInBundle
	hashes := make(map[common.Hash]int)
	for i, txOrGen := range b.BundleTxs {
		tx, ok := txOrGen.(*types.Transaction)
		if !ok {
			continue
		}
		hashes[tx.Hash()] = i
	}

	// 2-3. Check for TargetTxHash breaking current bundle.
	// If newBundle.TargetTxHash is equal to the last tx of current bundle, it is NOT a conflict.
	// e.g.) b.txs = [0x1, 0x2] and newBundle's TargetTxHash is 0x2.
	if idx, ok := hashes[newBundle.TargetTxHash]; ok && idx != len(b.BundleTxs)-1 {
		return true
	}

	// 2-4. Check for overlapping txs
	for _, txOrGen := range newBundle.BundleTxs {
		if tx, ok := txOrGen.(*types.Transaction); ok {
			if _, has := hashes[tx.Hash()]; has {
				return true
			}
		}
	}

	return false
}
