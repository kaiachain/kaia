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

package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
)

type Pair struct {
	u, v int
}

// two txs with same sender has an edge
// two txs in the same bundle has an edge
func BuildGraph(txs []interface{}, bundles []*builder.Bundle, signer types.Signer) ([][]int, error) {
	edges := make([][]int, len(txs))
	for i := range edges {
		edges[i] = make([]int, len(txs))
	}

	// Group txs by sender
	senderToIndices := make(map[common.Address][]int)
	for i, txOrGen := range txs {
		if tx, ok := txOrGen.(*types.Transaction); ok {
			from, err := types.Sender(signer, tx)
			if err != nil {
				return nil, err
			}
			senderToIndices[from] = append(senderToIndices[from], i)
		}
	}

	// Add edges between txs with same sender
	for _, indices := range senderToIndices {
		for _, u := range indices {
			for _, v := range indices {
				edges[u][v] = 1
			}
		}
	}

	for _, bundle := range bundles {
		indices := make([]int, 0)
		for i, txOrGen := range txs {
			if tx, ok := txOrGen.(*types.Transaction); ok {
				if bundle.Has(tx.Hash(), 0) {
					indices = append(indices, i)
				}
			}
		}

		// Add edges between all txs in bundle
		for _, u := range indices {
			for _, v := range indices {
				edges[u][v] = 1
			}
		}
	}

	return edges, nil
}

func Bfs(edges [][]int, indices []int) map[int]bool {
	// Initialize visited array
	visited := make([]bool, len(edges))
	for _, i := range indices {
		visited[i] = true
	}

	// Initialize queue with starting indices
	queue := make([]int, len(indices))
	copy(queue, indices)

	// Result slice to store nodes in BFS order
	result := make(map[int]bool)
	for _, i := range indices {
		result[i] = true
	}

	// BFS
	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]

		// Check all neighbors
		for neighbor := 0; neighbor < len(edges[current]); neighbor++ {
			if edges[current][neighbor] == 1 && !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				result[neighbor] = true
			}
		}
	}

	return result
}

// IncorporateBundleTx incorporates bundle transactions into the transaction list.
// Caller must ensure that there is no conflict between bundles.
func IncorporateBundleTx(txs []*types.Transaction, bundles []*builder.Bundle) ([]interface{}, error) {
	ret := make([]interface{}, len(txs))
	for i := range txs {
		ret[i] = txs[i]
	}

	for _, bundle := range bundles {
		var err error
		ret, err = incorporate(ret, bundle)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// incorporate assumes that `txs` does not contain any bundle transactions.
func incorporate(txs []interface{}, bundle *builder.Bundle) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(txs)+len(bundle.BundleTxs))
	targetFound := false

	// 1. place bundle at the beginning
	if bundle.TargetTxHash == (common.Hash{}) {
		ret = append(ret, bundle.BundleTxs...)
		targetFound = true
	}

	// 2. place bundle after TargetTxHash
	for _, txOrGen := range txs {
		switch tx := txOrGen.(type) {
		case *types.Transaction:
			// if tx-in-bundle, the tx will be appended when target is found.
			if bundle.Has(tx.Hash(), 0) {
				continue
			}
			// Because bundle.TargetTxHash cannot be in the bundle, we only check tx-not-in-bundle case.
			ret = append(ret, tx)
			if tx.Hash() == bundle.TargetTxHash {
				targetFound = true
				ret = append(ret, bundle.BundleTxs...)
			}
		default: // if tx is TxGenerator, unconditionally append
			ret = append(ret, txOrGen)
		}
	}

	if !targetFound {
		return nil, ErrFailedToIncorporateBundle
	}

	return ret, nil
}

// Arrayify flattens transaction heaps into a single array
func Arrayify(heap *types.TransactionsByPriceAndNonce) []*types.Transaction {
	ret := make([]*types.Transaction, 0)
	copied := heap.Copy()
	for !copied.Empty() {
		ret = append(ret, copied.Peek())
		copied.Shift()
	}
	return ret
}

// IsConflict checks if new bundles conflict with previous bundles
func IsConflict(prevBundles []*builder.Bundle, newBundles []*builder.Bundle) bool {
	for _, newBundle := range newBundles {
		for _, prevBundle := range prevBundles {
			if prevBundle.IsConflict(newBundle) {
				return true
			}
		}
	}

	return false
}
