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
	"errors"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
)

var (
	ErrFailedToIncorporateBundle = errors.New("failed to incorporate bundle")

	logger = log.NewModuleLogger(log.KaiaxBuilder)
)

// buildDependencyIndices builds a dependency indices of txs.
// Two txs with the same sender has an edge.
// Two txs in the same bundle has an edge.
func buildDependencyIndices(txs []*TxOrGen, bundles []*Bundle, signer types.Signer) (map[common.Address][]int, map[int][]int, error) {
	senderToIndices := make(map[common.Address][]int)
	bundleToIndices := make(map[int][]int)

	for i, txOrGen := range txs {
		if txOrGen.IsConcreteTx() {
			tx, _ := txOrGen.GetTx(0)
			from, err := types.Sender(signer, tx)
			if err != nil {
				return nil, nil, err
			}
			senderToIndices[from] = append(senderToIndices[from], i)
		}
		if bundleIdx := FindBundleIdx(bundles, txOrGen); bundleIdx != -1 {
			bundleToIndices[bundleIdx] = append(bundleToIndices[bundleIdx], i)
		}
	}

	return senderToIndices, bundleToIndices, nil
}

// IncorporateBundleTx incorporates bundle transactions into the transaction list.
// Caller must ensure that there is no conflict between bundles.
func IncorporateBundleTx(txs []*types.Transaction, bundles []*Bundle) ([]*TxOrGen, error) {
	ret := make([]*TxOrGen, len(txs))
	for i, tx := range txs {
		ret[i] = NewTxOrGenFromTx(tx)
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
func incorporate(txs []*TxOrGen, bundle *Bundle) ([]*TxOrGen, error) {
	ret := make([]*TxOrGen, 0, len(txs)+len(bundle.BundleTxs))
	targetFound := false

	// 1. place bundle at the beginning
	if bundle.TargetTxHash == (common.Hash{}) {
		ret = append(ret, bundle.BundleTxs...)
		targetFound = true
	}

	// 2. place bundle after TargetTxHash
	for _, txOrGen := range txs {
		// if tx-in-bundle, the tx will be appended when target is found.
		if bundle.Has(txOrGen) {
			continue
		}
		ret = append(ret, txOrGen)
		if txOrGen.Id == bundle.TargetTxHash {
			targetFound = true
			ret = append(ret, bundle.BundleTxs...)
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
func IsConflict(prevBundles []*Bundle, newBundles []*Bundle) bool {
	for _, newBundle := range newBundles {
		for _, prevBundle := range prevBundles {
			if prevBundle.IsConflict(newBundle) {
				return true
			}
		}
	}

	return false
}

// Filter removes elements at indices specified in `toRemove`.
func Filter[T any](slice *[]T, toRemove map[int]bool) []T {
	ret := make([]T, 0)
	for i := 0; i < len(*slice); i++ {
		if !toRemove[i] {
			ret = append(ret, (*slice)[i])
		}
	}
	return ret
}

func FindBundleIdx(bundles []*Bundle, txOrGen *TxOrGen) int {
	for i, bundle := range bundles {
		if bundle.Has(txOrGen) {
			return i
		}
	}
	return -1
}

func SetCorrectTargetTxHash(bundles []*Bundle, txs []*TxOrGen) []*Bundle {
	ret := make([]*Bundle, 0)
	for _, bundle := range bundles {
		bundle.TargetTxHash = FindTargetTxHash(bundle, txs)
		ret = append(ret, bundle)
	}
	return ret
}

func FindTargetTxHash(bundle *Bundle, txOrGens []*TxOrGen) common.Hash {
	for i := range txOrGens {
		if bundle.BundleTxs[0].Equals(txOrGens[i]) {
			if i == 0 {
				return common.Hash{}
			} else {
				return txOrGens[i-1].Id
			}
		}
	}
	return common.Hash{}
}

func ShiftTxs(txs *[]*TxOrGen, num int) {
	if len(*txs) <= num {
		*txs = (*txs)[:0]
		return
	}
	*txs = (*txs)[num:]
}

func PopTxs(txOrGens *[]*TxOrGen, num int, bundles *[]*Bundle, signer types.Signer) {
	if len(*txOrGens) == 0 || num == 0 {
		return
	}

	senderToIndices, bundleToIndices, err := buildDependencyIndices(*txOrGens, *bundles, signer)
	if err != nil {
		logger.Error("Failed to build dependency indices", "err", err)
		ShiftTxs(txOrGens, num)
		return
	}

	toRemove := make(map[int]bool)
	queue := make([]int, 0, num)

	for i := 0; i < min(num, len(*txOrGens)); i++ {
		toRemove[i] = true
		queue = append(queue, i)
	}

	for len(queue) > 0 {
		curIdx := queue[0]
		queue = queue[1:]
		txOrGen := (*txOrGens)[curIdx]

		if txOrGen.IsConcreteTx() {
			tx, _ := txOrGen.GetTx(0)
			from, _ := types.Sender(signer, tx)
			for _, idx := range senderToIndices[from] {
				if idx > curIdx && !toRemove[idx] {
					toRemove[idx] = true
					queue = append(queue, idx)
				}
			}
		}
		if bundleIdx := FindBundleIdx(*bundles, txOrGen); bundleIdx != -1 {
			for _, idx := range bundleToIndices[bundleIdx] {
				if !toRemove[idx] {
					toRemove[idx] = true
					queue = append(queue, idx)
				}
			}
		}
	}

	newTxs := Filter(txOrGens, toRemove)

	bundleIdxToRemove := map[int]bool{}
	for bundleIdx, txIndices := range bundleToIndices {
		for _, txIdx := range txIndices {
			if toRemove[txIdx] {
				bundleIdxToRemove[bundleIdx] = true
				break
			}
		}
	}

	newBundles := SetCorrectTargetTxHash(Filter(bundles, bundleIdxToRemove), newTxs)

	*txOrGens = newTxs
	*bundles = newBundles
}

func ExtractBundlesAndIncorporate(arrayTxs []*types.Transaction, txBundlingModules []TxBundlingModule) ([]*TxOrGen, []*Bundle) {
	// Detect bundles and add them to bundles
	bundles := []*Bundle{}
	flattenedTxs := []*TxOrGen{}
	if txBundlingModules == nil {
		for _, tx := range arrayTxs {
			flattenedTxs = append(flattenedTxs, NewTxOrGenFromTx(tx))
		}
		return flattenedTxs, nil
	}

	for _, txBundlingModule := range txBundlingModules {
		newBundles := txBundlingModule.ExtractTxBundles(arrayTxs, bundles)
		for _, newBundle := range newBundles {
			isConflict := false
			// Check for conflicts with all previous bundles
			for _, prevBundle := range bundles {
				isConflict = prevBundle.IsConflict(newBundle)
				if isConflict {
					break
				}
			}
			if !isConflict {
				bundles = append(bundles, newBundle)
			}
		}
	}

	incorporatedTxs, err := IncorporateBundleTx(arrayTxs, bundles)
	if err != nil {
		return flattenedTxs, nil
	}

	return incorporatedTxs, bundles
}
