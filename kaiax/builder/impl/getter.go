// Copyright 2024 The Kaia Authors
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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
)

// IncorporateBundleTx incorporates bundle transactions into the transaction list
func (b *BuilderModule) IncorporateBundleTx(txs []*types.Transaction, bundles []*builder.Bundle) []*types.Transaction {
	// TODO: implement
	return nil
}

// Arrayify flattens transaction heaps into a single array
func (b *BuilderModule) Arrayify(heap []*types.TransactionsByPriceAndNonce) []*types.Transaction {
	// TODO: implement
	result := make([]*types.Transaction, 0)
	return result
}

// IsConflict checks if new bundles conflict with previous bundles
func (b *BuilderModule) IsConflict(prevBundles []*builder.Bundle, newBundles []*builder.Bundle) bool {
	// collect all tx hashes from previous bundles
	prevTxHashes := make(map[common.Hash]struct{})
	for _, bundle := range prevBundles {
		for _, txOrGen := range bundle.BundleTxs {
			if tx, ok := txOrGen.(*types.Transaction); ok {
				prevTxHashes[tx.Hash()] = struct{}{}
			}
		}
	}

	// check if any new bundle tx exists in previous bundles
	for _, bundle := range newBundles {
		for _, txOrGen := range bundle.BundleTxs {
			if tx, ok := txOrGen.(*types.Transaction); ok {
				if _, exists := prevTxHashes[tx.Hash()]; exists {
					return true
				}
			}
		}
	}

	return false
}
