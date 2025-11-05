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
	"github.com/kaiachain/kaia/common"
)

type Bundle struct {
	// each element can be either *types.Transaction, or TxGenerator
	BundleTxs []*TxOrGen

	// BundleTxs is placed AFTER the target tx. If empty hash, it is placed at the very front.
	TargetTxHash common.Hash

	// TargetRequired is true if the bundle must be executed after the target tx
	// and only if the target tx is successfully executed.
	TargetRequired bool

	// lookup map for O(1) membership checks (lazy initialized)
	txLookup map[common.Hash]int
}

func NewBundle(txs []*TxOrGen, targetTxHash common.Hash, targetRequired bool) *Bundle {
	b := &Bundle{
		BundleTxs:      txs,
		TargetTxHash:   targetTxHash,
		TargetRequired: targetRequired,
	}
	b.buildLookup()
	return b
}

func (b *Bundle) buildLookup() {
	if b.txLookup == nil {
		b.txLookup = make(map[common.Hash]int, len(b.BundleTxs))
		for i, txOrGen := range b.BundleTxs {
			b.txLookup[txOrGen.Id] = i
		}
	}
}

// Has checks if the bundle contains a tx with the given hash.
func (b *Bundle) Has(txOrGen *TxOrGen) bool {
	_, exists := b.txLookup[txOrGen.Id]
	return exists
}

// FindIdx returns if the bundle contains a tx with the given hash and its index in bundle.
func (b *Bundle) FindIdx(id common.Hash) int {
	if idx, exists := b.txLookup[id]; exists {
		return idx
	}
	return -1
}

// IsConflict checks if newBundle conflicts with current bundle.
func (b *Bundle) IsConflict(newBundle *Bundle) bool {
	// 1. Check for same target tx hash and both are required
	// If both are required, it discards the new bundle.
	if b.TargetTxHash == newBundle.TargetTxHash && b.TargetRequired && newBundle.TargetRequired {
		return true
	}

	// 2-1. Empty bundleTxs does not conflict with other transactions
	if len(b.BundleTxs) == 0 {
		return false
	}

	// 2-2. Check for overlapping txs
	for _, txOrGen := range newBundle.BundleTxs {
		if b.Has(txOrGen) {
			return true
		}
	}

	// 2-3. Check for TargetTxHash breaking current bundle.
	// If newBundle.TargetTxHash is equal to the last tx of current bundle, it is NOT a conflict.
	// Check both direction to guarantee symmetry.
	// e.g.) b.txs = [0x1, 0x2] and newBundle's TargetTxHash is 0x2.
	if idx := b.FindIdx(newBundle.TargetTxHash); idx != -1 && idx != len(b.BundleTxs)-1 {
		return true
	}
	if idx := newBundle.FindIdx(b.TargetTxHash); idx != -1 && idx != len(newBundle.BundleTxs)-1 {
		return true
	}

	return false
}
