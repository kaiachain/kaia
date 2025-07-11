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
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestBundle_IsConflict(t *testing.T) {
	txs := make([]*types.Transaction, 4)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{}, common.Big0, 0, common.Big0, nil)
	}

	b0 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[0], txs[1]),
		TargetTxHash:   common.Hash{},
		TargetRequired: true,
	}
	defaultTargetHash := txs[1].Hash()

	testcases := []struct {
		name      string
		bundle    *Bundle
		newBundle *Bundle
		expected  bool
	}{
		{
			name:   "Same TargetTxHash (empty TargetHash)",
			bundle: b0,
			newBundle: &Bundle{
				BundleTxs:      NewTxOrGenList(),
				TargetTxHash:   common.Hash{},
				TargetRequired: true,
			},
			expected: true,
		},
		{
			name:   "Same TargetTxHash (empty TargetHash) but newBundle is not required",
			bundle: b0,
			newBundle: &Bundle{
				BundleTxs:      NewTxOrGenList(),
				TargetTxHash:   common.Hash{},
				TargetRequired: false,
			},
			expected: false,
		},
		{
			name:   "TargetTxHash divides a bundle",
			bundle: b0,
			newBundle: &Bundle{
				BundleTxs:      NewTxOrGenList(),
				TargetTxHash:   txs[0].Hash(),
				TargetRequired: true,
			},
			expected: true,
		},
		{
			name:   "Overlapping BundleTxs",
			bundle: b0,
			newBundle: &Bundle{
				BundleTxs:      NewTxOrGenList(txs[0]),
				TargetTxHash:   defaultTargetHash,
				TargetRequired: true,
			},
			expected: true,
		},
		{
			name:   "Non-overlapping BundleTxs",
			bundle: b0,
			newBundle: &Bundle{
				BundleTxs:      NewTxOrGenList(txs[2], txs[3]),
				TargetTxHash:   defaultTargetHash,
				TargetRequired: true,
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bundle.IsConflict(tc.newBundle)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestBundle_CoordinateTargetTxHash(t *testing.T) {
	txs := make([]*types.Transaction, 7)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{}, common.Big0, 0, common.Big0, nil)
	}

	// All bundles have same target tx hash
	b0 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[0]),
		TargetTxHash:   common.Hash{},
		TargetRequired: false,
	}
	b1 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[1]),
		TargetTxHash:   common.Hash{},
		TargetRequired: true,
	}
	b2 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[2]),
		TargetTxHash:   common.Hash{},
		TargetRequired: false,
	}
	b3 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[3]),
		TargetTxHash:   txs[2].Hash(),
		TargetRequired: false,
	}
	b4 := &Bundle{
		BundleTxs:      NewTxOrGenList(txs[4]),
		TargetTxHash:   txs[2].Hash(),
		TargetRequired: false,
	}

	bundles := []*Bundle{b0, b1, b2, b3, b4}

	coordinatedBundles := CoordinateTargetTxHash(bundles)

	expectedTargetTxHashes := []common.Hash{
		{},
		txs[1].Hash(),
		txs[0].Hash(),
		txs[2].Hash(),
		txs[3].Hash(),
	}

	for i, bundle := range coordinatedBundles {
		assert.Equal(t, expectedTargetTxHashes[i], bundle.TargetTxHash)
	}
}
