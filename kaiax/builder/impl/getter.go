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
	"bytes"
	"slices"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/builder"
)

// buildDependencyIndices builds a dependency indices of txs.
// Two txs with the same sender has an edge.
// Two txs in the same bundle has an edge.
func buildDependencyIndices(txs []interface{}, bundles []*builder.Bundle, signer types.Signer) (map[common.Address][]int, map[int][]int, error) {
	senderToIndices := make(map[common.Address][]int)
	bundleToIndices := make(map[int][]int)

	for i, txOrGen := range txs {
		if tx, ok := txOrGen.(*types.Transaction); ok {
			from, err := types.Sender(signer, tx)
			if err != nil {
				return nil, nil, err
			}
			senderToIndices[from] = append(senderToIndices[from], i)
		}
		if bundleIdx := FindBundleIdxAsTxOrGen(bundles, txOrGen); bundleIdx != -1 {
			bundleToIndices[bundleIdx] = append(bundleToIndices[bundleIdx], i)
		}
	}

	return senderToIndices, bundleToIndices, nil
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
			if bundle.Has(tx.Hash()) {
				continue
			}
			// Because bundle.TargetTxHash cannot be in the bundle, we only check tx-not-in-bundle case.
			ret = append(ret, tx)
			if tx.Hash() == bundle.TargetTxHash {
				targetFound = true
				for i, txInBundleI := range bundle.BundleTxs {
					switch txInBundle := txInBundleI.(type) {
					case *types.Transaction:
						ret = append(ret, txInBundleI)
					case builder.TxGenerator:
						txInBundle.Hash = crypto.Keccak256Hash(bundle.TargetTxHash[:], common.Int64ToByteLittleEndian(uint64(i)))
						txInBundleI = txInBundle
						ret = append(ret, txInBundleI)
					}
				}
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

func FindBundleIdx(bundles []*builder.Bundle, tx *types.Transaction) int {
	for i, bundle := range bundles {
		if bundle.Has(tx.Hash()) {
			return i
		}
	}
	return -1
}

func FindBundleIdxAsTxOrGen(bundles []*builder.Bundle, txOrGen interface{}) int {
	for i, bundle := range bundles {
		for _, txOrGenInBundle := range bundle.BundleTxs {
			if EqualTxOrGen(txOrGenInBundle, txOrGen) {
				return i
			}
		}
	}
	return -1
}

func EqualTxOrGen(txOrGenIX, txOrGenIY interface{}) bool {
	var (
		txOrGenXHash common.Hash
		txOrGenYHash common.Hash
	)

	switch txOrGenX := txOrGenIX.(type) {
	case *types.Transaction:
		txOrGenXHash = txOrGenX.Hash()
	case builder.TxGenerator:
		txOrGenXHash = txOrGenX.Hash
	}

	switch txOrGenY := txOrGenIY.(type) {
	case *types.Transaction:
		txOrGenYHash = txOrGenY.Hash()
	case builder.TxGenerator:
		txOrGenYHash = txOrGenY.Hash
	}

	return bytes.Equal(txOrGenXHash.Bytes(), txOrGenYHash.Bytes())
}

func SetCorrectTargetTxHash(bundles []*builder.Bundle, txs []interface{}) []*builder.Bundle {
	ret := make([]*builder.Bundle, 0)
	for _, bundle := range bundles {
		bundle.TargetTxHash = FindTargetTxHash(bundle, txs)
		ret = append(ret, bundle)
	}
	return ret
}

func FindTargetTxHash(bundle *builder.Bundle, txs []interface{}) common.Hash {
	// If this is never updated it is expected to return an empty hash.
	targetTxHash := common.Hash{}
	for i, tx := range txs {
		// For index greater than 0, if it can be cast to *types.Transaction then record this as the TargetTxHash.
		if i > 0 {
			if txTarget, ok := txs[i-1].(*types.Transaction); ok {
				targetTxHash = txTarget.Hash()
			}
		}
		// If tx is the first tx in the bundle then there is no need to look further.
		if EqualTxOrGen(bundle.BundleTxs[0], tx) {
			break
		}
	}
	return targetTxHash
}

func ShiftTxs(txs *[]interface{}, num int) {
	if len(*txs) <= num {
		*txs = (*txs)[:0]
		return
	}
	*txs = (*txs)[num:]
}

func PopTxs(txs *[]interface{}, num int, bundles *[]*builder.Bundle, signer types.Signer) {
	if len(*txs) == 0 || num == 0 {
		return
	}

	senderToIndices, bundleToIndices, err := buildDependencyIndices(*txs, *bundles, signer)
	if err != nil {
		logger.Error("Failed to build dependency indices", "err", err)
		ShiftTxs(txs, num)
		return
	}

	toRemove := make(map[int]bool)
	queue := make([]int, 0, num)

	for i := 0; i < min(num, len(*txs)); i++ {
		toRemove[i] = true
		queue = append(queue, i)
	}

	for len(queue) > 0 {
		curIdx := queue[0]
		queue = queue[1:]

		if tx, ok := (*txs)[curIdx].(*types.Transaction); ok {
			from, _ := types.Sender(signer, tx)
			for _, idx := range senderToIndices[from] {
				if idx > curIdx && !toRemove[idx] {
					toRemove[idx] = true
					queue = append(queue, idx)
				}
			}
		}
		if bundleIdx := FindBundleIdxAsTxOrGen(*bundles, (*txs)[curIdx]); bundleIdx != -1 {
			for _, idx := range bundleToIndices[bundleIdx] {
				if !toRemove[idx] {
					toRemove[idx] = true
					queue = append(queue, idx)
				}
			}
		}
	}

	newTxs := Filter(txs, toRemove)

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

	*txs = newTxs
	*bundles = newBundles
}

func ExtractBundlesAndIncorporate(arrayTxs []*types.Transaction, txBundlingModules []builder.TxBundlingModule) ([]interface{}, []*builder.Bundle) {
	// Detect bundles and add them to bundles
	bundles := []*builder.Bundle{}
	flattenedTxs := []interface{}{}
	if txBundlingModules == nil {
		for _, tx := range arrayTxs {
			var itx interface{} = tx
			flattenedTxs = append(flattenedTxs, itx)
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

// WrapAndConcatenateBundlingModules wraps bundling modules and concatenates them.
// given: mTxPool = [ A, B, C ], mTxBundling = [ B, D ]
// want : mTxPool = [ A, WB, C, WD ] (W: wrapped)
func WrapAndConcatenateBundlingModules(mTxBundling []builder.TxBundlingModule, mTxPool []kaiax.TxPoolModule, txPool kaiax.TxPoolForCaller) []kaiax.TxPoolModule {
	ret := make([]kaiax.TxPoolModule, 0, len(mTxBundling)+len(mTxPool))

	for _, module := range mTxPool {
		if txb, ok := module.(builder.TxBundlingModule); ok {
			ret = append(ret, NewBuilderWrappingModule(txb, txPool))
		} else {
			ret = append(ret, module)
		}
	}

	for _, module := range mTxBundling {
		if txp, ok := module.(kaiax.TxPoolModule); !(ok && slices.Contains(mTxPool, txp)) {
			ret = append(ret, NewBuilderWrappingModule(module, txPool))
		}
	}

	return ret
}
