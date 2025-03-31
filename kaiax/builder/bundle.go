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
}

func NewBundle(txs []*TxOrGen, targetTxHash common.Hash, targetRequired bool) *Bundle {
	return &Bundle{
		BundleTxs:      txs,
		TargetTxHash:   targetTxHash,
		TargetRequired: targetRequired,
	}
}

// Has checks if the bundle contains a tx with the given hash.
func (b *Bundle) Has(txOrGen *TxOrGen) bool {
	return b.FindIdx(txOrGen.Id) != -1
}

// FindIdx returns if the bundle contains a tx with the given hash and its index in bundle.
func (b *Bundle) FindIdx(id common.Hash) int {
	for i, txOrGen := range b.BundleTxs {
		if txOrGen.Id == id {
			return i
		}
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
	// e.g.) b.txs = [0x1, 0x2] and newBundle's TargetTxHash is 0x2.
	if idx := b.FindIdx(newBundle.TargetTxHash); idx == -1 || idx == len(b.BundleTxs)-1 {
		return false
	}

	return true
}

// CoordinateTargetTxHash coordinates the target tx hash of bundles.
// It assumes there's only one bundle with TargetRequired = true among the bundles with the same TargetTxHash
// and no zero-length bundle.
// e.g.) bundles = [
//
//	{TargetTxHash: 0x1, TargetRequired: false, BundleTxs: []*TxOrGen{tx3, tx4}},
//	{TargetTxHash: 0x1, TargetRequired: true, BundleTxs: []*TxOrGen{tx1, tx2}},
//
// ]
// -> returns [
//
//	{TargetTxHash: 0x1, TargetRequired: true, BundleTxs: []*TxOrGen{tx1, tx2}},
//	{TargetTxHash: 0x2, TargetRequired: false, BundleTxs: []*TxOrGen{tx3, tx4}},
//
// ]
func CoordinateTargetTxHash(bundles []*Bundle) []*Bundle {
	if len(bundles) <= 1 {
		return bundles
	}

	newBundles := make([]*Bundle, 0, len(bundles))
	sameTargetTxHashBundles := make(map[common.Hash][]*Bundle)

	for _, bundle := range bundles {
		sameTargetTxHashBundles[bundle.TargetTxHash] = append(sameTargetTxHashBundles[bundle.TargetTxHash], bundle)
	}

	for _, list := range sameTargetTxHashBundles {
		if len(list) == 1 {
			newBundles = append(newBundles, list[0])
			continue
		}

		// Find the bundle with TargetRequired = true and move it to the front.
		// This is needed because #incorporate assumes that targetTxHash is already in the txs.
		for i, bundle := range list {
			if bundle.TargetRequired {
				list[0], list[i] = list[i], list[0]
				break
			}
		}

		for i, bundle := range list {
			if i == 0 {
				continue
			}
			bundle.TargetTxHash = lastBundleTx(list[i-1]).Id
		}
		newBundles = append(newBundles, list...)
	}

	return newBundles
}

func lastBundleTx(bundle *Bundle) *TxOrGen {
	return bundle.BundleTxs[len(bundle.BundleTxs)-1]
}
