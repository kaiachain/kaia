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
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncorporateBundleTx(t *testing.T) {
	// Create test transactions
	txs := []*types.Transaction{
		types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(2, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(3, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
	}
	var g builder.TxGenerator = func(nonce uint64) (*types.Transaction, error) {
		return types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
	}

	b := NewBuilderModule()

	testCases := []struct {
		name     string
		bundles  []*builder.Bundle
		expected []interface{}
	}{
		{
			name: "incorporate multiple bundles",
			bundles: []*builder.Bundle{
				{BundleTxs: []interface{}{txs[0], txs[1]}, TargetTxHash: common.Hash{}},
				{BundleTxs: []interface{}{txs[2]}, TargetTxHash: txs[1].Hash()},
			},
			expected: []interface{}{txs[0], txs[1], txs[2], txs[3]},
		},
		{
			name:     "incorporate empty bundles",
			bundles:  []*builder.Bundle{},
			expected: []interface{}{txs[0], txs[1], txs[2], txs[3]},
		},
		{
			name: "incorporate bundle with generator",
			bundles: []*builder.Bundle{
				{BundleTxs: []interface{}{txs[0], g}, TargetTxHash: common.Hash{}},
			},
			expected: []interface{}{txs[0], g, txs[1], txs[2], txs[3]},
		},
		{
			name: "incorporate bundle with generator 2",
			bundles: []*builder.Bundle{
				{BundleTxs: []interface{}{g, txs[0]}, TargetTxHash: common.Hash{}},
				{BundleTxs: []interface{}{g, txs[1]}, TargetTxHash: txs[0].Hash()},
			},
			expected: []interface{}{g, txs[0], g, txs[1], txs[2], txs[3]},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := b.IncorporateBundleTx(txs, tc.bundles)
			require.Nil(t, err)
			require.Equal(t, len(tc.expected), len(ret))
			for i := range ret {
				require.Equal(t, reflect.TypeOf(tc.expected[i]), reflect.TypeOf(ret[i]), "type at ret[%d] is different", i)
				switch v := ret[i].(type) {
				case *types.Transaction:
					t.Logf("expected txs[%d]=hash %s", i, tc.expected[i].(*types.Transaction).Hash().Hex())
					t.Logf("actual txs[%d]=hash %s", i, v.Hash().Hex())
					expected, ok := tc.expected[i].(*types.Transaction)
					assert.True(t, ok, "tx %d", i)
					assert.Equal(t, expected.Hash(), v.Hash(), "tx %d", i, "nonce", v.Nonce)
				case builder.TxGenerator:
					_, ok := tc.expected[i].(builder.TxGenerator)
					assert.True(t, ok, "tx %d", i)
				}
			}
		})
	}
}

func TestIncorporate(t *testing.T) {
	// Create test transactions
	txs := []*types.Transaction{
		types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(2, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
	}
	txOrGenList := []interface{}{txs[0], txs[1], txs[2]}
	testCases := []struct {
		name     string
		bundle   *builder.Bundle
		expected []interface{}
	}{
		{
			name:     "incorporate first two transactions",
			bundle:   &builder.Bundle{BundleTxs: []interface{}{txs[0], txs[1]}, TargetTxHash: common.Hash{}},
			expected: []interface{}{txs[0], txs[1], txs[2]},
		},
		{
			name:     "incorporate last two transactions",
			bundle:   &builder.Bundle{BundleTxs: []interface{}{txs[1], txs[2]}, TargetTxHash: common.Hash{}},
			expected: []interface{}{txs[1], txs[2], txs[0]},
		},
		{
			name:     "incorporate with target hash",
			bundle:   &builder.Bundle{BundleTxs: []interface{}{txs[0]}, TargetTxHash: txs[2].Hash()},
			expected: []interface{}{txs[1], txs[2], txs[0]},
		},
		{
			name:     "incorporate single transaction",
			bundle:   &builder.Bundle{BundleTxs: []interface{}{txs[2]}, TargetTxHash: common.Hash{}},
			expected: []interface{}{txs[2], txs[0], txs[1]},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := incorporate(txOrGenList, tc.bundle)
			require.Nil(t, err)
			assert.Equal(t, tc.expected, ret)
		})
	}
}

func TestArrayify(t *testing.T) {
	// Generate a batch of accounts to start with
	keyLen := 10
	txLen := 30
	keys := make([]*ecdsa.PrivateKey, keyLen)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := types.LatestSignerForChainID(common.Big1)
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]types.Transactions{}
	hashes := map[common.Hash]bool{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < txLen; i++ {
			tx, _ := types.SignTx(types.NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil), signer, key)
			groups[addr] = append(groups[addr], tx)
			hashes[tx.Hash()] = true
		}
	}

	heap := types.NewTransactionsByPriceAndNonce(signer, groups, nil)
	b := NewBuilderModule()
	txs := b.Arrayify(heap)
	assert.Equal(t, keyLen*txLen, len(txs))
	for i := range txs {
		assert.Equal(t, true, hashes[txs[i].Hash()])
	}
	assert.False(t, heap.Empty()) // don't modify the original heap
}

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
			newBundles:  []*builder.Bundle{{BundleTxs: []interface{}{}, TargetTxHash: common.Hash{}}},
			expected:    true,
		},
		{
			name:        "TargetTxHash divides a bundle",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: []interface{}{}, TargetTxHash: txs[0].Hash()}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 1",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: []interface{}{txs[0], txs[2]}, TargetTxHash: defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 2",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: []interface{}{txs[1], txs[2], txs[3]}, TargetTxHash: defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Non-overlapping BundleTxs",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: []interface{}{txs[2], txs[3]}, TargetTxHash: defaultTargetHash}},
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
