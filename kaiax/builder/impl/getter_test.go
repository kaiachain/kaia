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
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	"gotest.tools/assert"
)

func TestIsConflict(t *testing.T) {
	txs := make([]*types.Transaction, 4)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{}, common.Big0, 0, common.Big0, nil)
	}

	b0 := &builder.Bundle{
		BundleTxs:    []interface{}{txs[0], txs[1]},
		TargetTxHash: common.Hash{},
	}
	defaultTargetHash := txs[1].Hash() // make TargetTxHash checks pass

	testcases := []struct {
		name        string
		prevBundles []*builder.Bundle
		newBundles  []*builder.Bundle
		expected    bool
	}{
		{
			name:        "Same TargetTxHash",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{[]interface{}{}, common.Hash{}}},
			expected:    true,
		},
		{
			name:        "TargetTxHash divides a bundle",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{[]interface{}{}, txs[0].Hash()}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 1",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{[]interface{}{txs[0], txs[2]}, defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 2",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{[]interface{}{txs[1], txs[2], txs[3]}, defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Non-overlapping BundleTxs",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{[]interface{}{txs[2], txs[3]}, defaultTargetHash}},
			expected:    false,
		},
	}

	b := NewBuilderModule()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotConflict := b.IsConflict(tc.prevBundles, tc.newBundles)
			assert.Equal(t, tc.expected, gotConflict)
		})
	}
}
