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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/builder"
	mock_builder "github.com/kaiachain/kaia/kaiax/builder/mock"
	mock_kaiax "github.com/kaiachain/kaia/kaiax/mock"
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

	gen := func(nonce uint64) (*types.Transaction, error) {
		return types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
	}
	g1 := builder.NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := builder.NewTxOrGenFromGen(gen, common.Hash{2})

	testCases := []struct {
		name     string
		bundles  []*builder.Bundle
		expected []*builder.TxOrGen
	}{
		{
			name: "incorporate multiple bundles",
			bundles: []*builder.Bundle{
				{BundleTxs: builder.NewTxOrGenList(txs[0], txs[1]), TargetTxHash: common.Hash{}},
				{BundleTxs: builder.NewTxOrGenList(txs[2]), TargetTxHash: txs[1].Hash()},
			},
			expected: builder.NewTxOrGenList(txs[0], txs[1], txs[2], txs[3]),
		},
		{
			name:     "incorporate empty bundles",
			bundles:  []*builder.Bundle{},
			expected: builder.NewTxOrGenList(txs[0], txs[1], txs[2], txs[3]),
		},
		{
			name: "incorporate bundle with generator",
			bundles: []*builder.Bundle{
				{BundleTxs: builder.NewTxOrGenList(txs[0], g1), TargetTxHash: common.Hash{}},
			},
			expected: builder.NewTxOrGenList(txs[0], g1, txs[1], txs[2], txs[3]),
		},
		{
			name: "incorporate bundle with generator 2",
			bundles: []*builder.Bundle{
				{BundleTxs: builder.NewTxOrGenList(g1, txs[0]), TargetTxHash: common.Hash{}},
				{BundleTxs: builder.NewTxOrGenList(g2, txs[1]), TargetTxHash: txs[0].Hash()},
			},
			expected: builder.NewTxOrGenList(g1, txs[0], g2, txs[1], txs[2], txs[3]),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := IncorporateBundleTx(txs, tc.bundles)
			require.Nil(t, err)
			require.Equal(t, len(tc.expected), len(ret))
			for i := range ret {
				assert.Equal(t, tc.expected[i].Id.Hex(), ret[i].Id.Hex(), "mismatch at ret[%d]", i)
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
	txOrGenList := builder.NewTxOrGenList(txs[0], txs[1], txs[2])
	testCases := []struct {
		name     string
		bundle   *builder.Bundle
		expected []*builder.TxOrGen
	}{
		{
			name:     "incorporate first two transactions",
			bundle:   &builder.Bundle{BundleTxs: builder.NewTxOrGenList(txs[0], txs[1]), TargetTxHash: common.Hash{}},
			expected: builder.NewTxOrGenList(txs[0], txs[1], txs[2]),
		},
		{
			name:     "incorporate last two transactions",
			bundle:   &builder.Bundle{BundleTxs: builder.NewTxOrGenList(txs[1], txs[2]), TargetTxHash: common.Hash{}},
			expected: builder.NewTxOrGenList(txs[1], txs[2], txs[0]),
		},
		{
			name:     "incorporate with target hash",
			bundle:   &builder.Bundle{BundleTxs: builder.NewTxOrGenList(txs[0]), TargetTxHash: txs[2].Hash()},
			expected: builder.NewTxOrGenList(txs[1], txs[2], txs[0]),
		},
		{
			name:     "incorporate single transaction",
			bundle:   &builder.Bundle{BundleTxs: builder.NewTxOrGenList(txs[2]), TargetTxHash: common.Hash{}},
			expected: builder.NewTxOrGenList(txs[2], txs[0], txs[1]),
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
	txs := Arrayify(heap)
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
		BundleTxs:    builder.NewTxOrGenList(txs[0], txs[1]),
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
			newBundles:  []*builder.Bundle{{BundleTxs: []*builder.TxOrGen{}, TargetTxHash: common.Hash{}}},
			expected:    true,
		},
		{
			name:        "TargetTxHash divides a bundle",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: []*builder.TxOrGen{}, TargetTxHash: txs[0].Hash()}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 1",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: builder.NewTxOrGenList(txs[0], txs[2]), TargetTxHash: defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 2",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: builder.NewTxOrGenList(txs[1], txs[2], txs[3]), TargetTxHash: defaultTargetHash}},
			expected:    true,
		},
		{
			name:        "Non-overlapping BundleTxs",
			prevBundles: []*builder.Bundle{b0},
			newBundles:  []*builder.Bundle{{BundleTxs: builder.NewTxOrGenList(txs[2], txs[3]), TargetTxHash: defaultTargetHash}},
			expected:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotConflict := IsConflict(tc.prevBundles, tc.newBundles)
			assert.Equal(t, tc.expected, gotConflict)
		})
	}
}

func TestPopTxs(t *testing.T) {
	var (
		signer = types.LatestSignerForChainID(big.NewInt(1))
		keys   = make([]*ecdsa.PrivateKey, 4)
		addrs  = make([]common.Address, 4)
		txs    = make([]*types.Transaction, 7)
	)
	gen := func(nonce uint64) (*types.Transaction, error) {
		return types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
	}
	g1 := builder.NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := builder.NewTxOrGenFromGen(gen, common.Hash{2})

	for i := range keys {
		keys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
	}

	for i := range txs {
		addr, key := addrs[i/2], keys[i/2]
		txs[i], _ = types.SignTx(types.NewTransaction(uint64(i), addr, big.NewInt(1), 21000, big.NewInt(1), nil), signer, key)
	}

	// Create test bundles
	bundles := []*builder.Bundle{
		{
			BundleTxs: builder.NewTxOrGenList(txs[1], txs[2]),
		},
		{
			BundleTxs:    builder.NewTxOrGenList(txs[3], txs[4]),
			TargetTxHash: txs[2].Hash(),
		},
		{
			BundleTxs:    builder.NewTxOrGenList(g1, txs[5]),
			TargetTxHash: txs[4].Hash(),
		},
		{
			BundleTxs: builder.NewTxOrGenList(g1, txs[1], txs[2]),
		},
	}

	testCases := []struct {
		name            string
		incorporatedTxs []*builder.TxOrGen
		numToPop        int
		bundles         []*builder.Bundle
		expectedTxs     []*builder.TxOrGen
	}{
		{
			name:            "Without any dependencies",
			incorporatedTxs: builder.NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        1,
			bundles:         []*builder.Bundle{},
			expectedTxs:     builder.NewTxOrGenList(txs[2], txs[3], txs[4], txs[5]),
		},
		{
			name:            "No bundles, tx0 and tx1 dependency (same sender)",
			incorporatedTxs: builder.NewTxOrGenList(txs[0], txs[1], txs[2], txs[3], txs[4]),
			numToPop:        1,
			bundles:         []*builder.Bundle{},
			expectedTxs:     builder.NewTxOrGenList(txs[2], txs[3], txs[4]),
		},
		{
			name:            "One bundle - first tx is generator",
			incorporatedTxs: builder.NewTxOrGenList(g1, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        1,
			bundles:         []*builder.Bundle{bundles[3]},
			expectedTxs:     builder.NewTxOrGenList(txs[4], txs[5], txs[6]),
		},
		{
			name:            "Two bundles - chaining dependency",
			incorporatedTxs: builder.NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        2,
			bundles:         []*builder.Bundle{bundles[0], bundles[1]},
			expectedTxs:     builder.NewTxOrGenList(),
		},
		{
			name:            "Two bundles - one independent tx (tx6)",
			incorporatedTxs: builder.NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*builder.Bundle{bundles[0], bundles[1]},
			expectedTxs:     builder.NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - change order",
			incorporatedTxs: builder.NewTxOrGenList(txs[2], txs[3], txs[4], txs[6], txs[5]), // 6 is before 5
			numToPop:        2,
			bundles:         []*builder.Bundle{bundles[0], bundles[1]},
			expectedTxs:     builder.NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - one independent tx (tx6) with one generator",
			incorporatedTxs: builder.NewTxOrGenList(g1, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*builder.Bundle{bundles[0], bundles[1]},
			expectedTxs:     builder.NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - two generators",
			incorporatedTxs: builder.NewTxOrGenList(g1, g2, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*builder.Bundle{bundles[0], bundles[1]},
			expectedTxs:     builder.NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			PopTxs(&tc.incorporatedTxs, tc.numToPop, &tc.bundles, signer)

			if len(tc.incorporatedTxs) != len(tc.expectedTxs) {
				t.Errorf("Expected %d transactions, got %d", len(tc.expectedTxs), len(tc.incorporatedTxs))
			}
			for i := range tc.expectedTxs {
				assert.Equal(t, tc.expectedTxs[i].Id.Hex(), tc.incorporatedTxs[i].Id.Hex())
			}
		})
	}
}

func TestWrapAndConcatenateBundlingModules(t *testing.T) {
	// Create mock modules
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create test modules
	mockTxPool1 := mock_kaiax.NewMockTxPoolModule(ctrl)
	mockTxPool2 := mock_kaiax.NewMockTxPoolModule(ctrl)
	mockTxBundling1 := mock_builder.NewMockTxBundlingModule(ctrl)
	mockTxBundling2 := mock_builder.NewMockTxBundlingModule(ctrl)

	// Create a module that implements both interfaces
	mockBoth1 := struct {
		*mock_kaiax.MockTxPoolModule
		*mock_builder.MockTxBundlingModule
	}{
		MockTxPoolModule:     mock_kaiax.NewMockTxPoolModule(ctrl),
		MockTxBundlingModule: mock_builder.NewMockTxBundlingModule(ctrl),
	}
	mockBoth2 := struct {
		*mock_kaiax.MockTxPoolModule
		*mock_builder.MockTxBundlingModule
	}{
		MockTxPoolModule:     mock_kaiax.NewMockTxPoolModule(ctrl),
		MockTxBundlingModule: mock_builder.NewMockTxBundlingModule(ctrl),
	}

	testCases := []struct {
		name        string
		mTxBundling []builder.TxBundlingModule
		mTxPool     []kaiax.TxPoolModule
		expected    []interface{}
	}{
		{
			name:        "No modules",
			mTxBundling: []builder.TxBundlingModule{},
			mTxPool:     []kaiax.TxPoolModule{},
			expected:    []interface{}{},
		},
		{
			name:        "Only TxPool modules",
			mTxBundling: []builder.TxBundlingModule{},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockTxPool2},
			expected:    []interface{}{mockTxPool1, mockTxPool2},
		},
		{
			name:        "Only TxBundling modules",
			mTxBundling: []builder.TxBundlingModule{mockTxBundling1, mockTxBundling2},
			mTxPool:     []kaiax.TxPoolModule{},
			expected:    []interface{}{mockTxBundling1, mockTxBundling2},
		},
		{
			name:        "Mixed modules",
			mTxBundling: []builder.TxBundlingModule{mockTxBundling1, mockTxBundling2},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockTxPool2},
			expected:    []interface{}{mockTxPool1, mockTxPool2, mockTxBundling1, mockTxBundling2},
		},
		{
			name:        "Overlapping modules",
			mTxBundling: []builder.TxBundlingModule{mockBoth2, mockTxBundling1, mockBoth1},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockBoth1, mockBoth2},
			expected:    []interface{}{mockTxPool1, mockBoth1, mockBoth2, mockTxBundling1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := WrapAndConcatenateBundlingModules(tc.mTxBundling, tc.mTxPool, nil)

			// Check the length of the result
			assert.Equal(t, len(tc.expected), len(result))

			// Verify that modules are properly wrapped
			for i, module := range result {
				if wrapped, ok := module.(*BuilderWrappingModule); ok {
					assert.Equal(t, wrapped.txBundlingModule, tc.expected[i])
				} else {
					assert.Equal(t, module, tc.expected[i])
				}
			}
		})
	}
}
