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
func (b *BuilderModule) IncorporateBundleTx(txs []*types.Transaction, bundles []*builder.Bundle) []interface{} {
	ret := make([]interface{}, len(txs))
	for i := range txs {
		ret[i] = txs[i]
	}

	for _, bundle := range bundles {
		ret = incorporate(ret, bundle)
	}
	return ret
}

func incorporate(txs []interface{}, bundle *builder.Bundle) []interface{} {
	ret := make([]interface{}, 0)

	// 1. collect txs that are not in bundle
	for _, txOrGen := range txs {
		switch tx := txOrGen.(type) {
		case *types.Transaction:
			if !bundle.Has(tx.Hash()) {
				ret = append(ret, tx)
			}
		case builder.TxGenerator:
			// append generator unconditionally
			ret = append(ret, tx)
		}
	}

	// 2. place bundle before TargetTxHash
	if bundle.TargetTxHash == (common.Hash{}) {
		ret = append(bundle.BundleTxs, ret...)
		return ret
	}

	for i, txOrGen := range ret {
		tx, ok := txOrGen.(*types.Transaction)
		if !ok {
			continue
		}
		if tx.Hash() == bundle.TargetTxHash {
			tmp := ret
			ret = make([]interface{}, len(ret)+len(bundle.BundleTxs))
			copy(ret[:i+1], tmp[:i+1])
			copy(ret[i+1:], bundle.BundleTxs)
			copy(ret[i+1+len(bundle.BundleTxs):], tmp[i+1:])
			return ret
		}
	}

	// must not reach here
	logger.Crit("failed to incorporate bundle")
	return ret
}

// Arrayify flattens transaction heaps into a single array
func (b *BuilderModule) Arrayify(heap *types.TransactionsByPriceAndNonce) []*types.Transaction {
	ret := make([]*types.Transaction, 0)
	copied := heap.Copy()
	for !copied.Empty() {
		ret = append(ret, copied.Peek())
		copied.Shift()
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
