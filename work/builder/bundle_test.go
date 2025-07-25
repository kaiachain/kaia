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

	b0 := NewBundle(
		NewTxOrGenList(txs[0], txs[1]),
		common.Hash{},
		true,
	)
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
			newBundle: NewBundle(
				NewTxOrGenList(),
				common.Hash{},
				true,
			),
			expected: true,
		},
		{
			name:   "Same TargetTxHash (empty TargetHash) but newBundle is not required",
			bundle: b0,
			newBundle: NewBundle(
				NewTxOrGenList(),
				common.Hash{},
				false,
			),
			expected: false,
		},
		{
			name:   "TargetTxHash divides a bundle",
			bundle: b0,
			newBundle: NewBundle(
				NewTxOrGenList(),
				txs[0].Hash(),
				true,
			),
			expected: true,
		},
		{
			name:   "Overlapping BundleTxs",
			bundle: b0,
			newBundle: NewBundle(
				NewTxOrGenList(txs[0]),
				defaultTargetHash,
				true,
			),
			expected: true,
		},
		{
			name:   "Non-overlapping BundleTxs",
			bundle: b0,
			newBundle: NewBundle(
				NewTxOrGenList(txs[2], txs[3]),
				defaultTargetHash,
				true,
			),
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
