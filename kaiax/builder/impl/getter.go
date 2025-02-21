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
	"github.com/kaiachain/kaia/kaiax/builder"
)

// IncorporateBundleTx incorporates bundle transactions into the transaction list
func (b *BuilderModule) IncorporateBundleTx(txs []*types.Transaction, bundles []*builder.Bundle) []*types.Transaction {
	// TODO: implement
	return nil
}

// Arrayify flattens transaction heaps into a single array
func (b *BuilderModule) Arrayify(heap *types.TransactionsByPriceAndNonce) []*types.Transaction {
	// TODO: deep copy heap
	ret := make([]*types.Transaction, 0)
	for !heap.Empty() {
		ret = append(ret, heap.Peek())
		heap.Pop()
	}
	return ret
}

// IsConflict checks if new bundles conflict with previous bundles
func (b *BuilderModule) IsConflict(prevBundles []*builder.Bundle, newBundles []*builder.Bundle) bool {
	for _, newBundle := range newBundles {
		for _, prevBundle := range prevBundles {
			if prevBundle.IsConflict(newBundle) {
				return true
			}
		}
	}

	return false
}
