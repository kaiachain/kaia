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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
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
	g1 := NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := NewTxOrGenFromGen(gen, common.Hash{2})

	testCases := []struct {
		name     string
		bundles  []*Bundle
		expected []*TxOrGen
	}{
		{
			name: "incorporate multiple bundles",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(txs[0], txs[1]), common.Hash{}, false),
				NewBundle(NewTxOrGenList(txs[2]), txs[1].Hash(), false),
			},
			expected: NewTxOrGenList(txs[0], txs[1], txs[2], txs[3]),
		},
		{
			name:     "incorporate empty bundles",
			bundles:  []*Bundle{},
			expected: NewTxOrGenList(txs[0], txs[1], txs[2], txs[3]),
		},
		{
			name: "incorporate bundle with generator",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(txs[0], g1), common.Hash{}, false),
			},
			expected: NewTxOrGenList(txs[0], g1, txs[1], txs[2], txs[3]),
		},
		{
			name: "incorporate bundle with generator 2",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(g1, txs[0]), common.Hash{}, false),
				NewBundle(NewTxOrGenList(g2, txs[1]), txs[0].Hash(), false),
			},
			expected: NewTxOrGenList(g1, txs[0], g2, txs[1], txs[2], txs[3]),
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

func TestExtractBundles(t *testing.T) {
	// Create test transactions
	txs := []*types.Transaction{
		types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(2, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(3, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(4, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
		types.NewTransaction(5, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
	}

	gen := func(nonce uint64) (*types.Transaction, error) {
		return types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
	}
	g1 := NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := NewTxOrGenFromGen(gen, common.Hash{2})
	g3 := NewTxOrGenFromGen(gen, common.Hash{3})

	testCases := []struct {
		name                 string
		bundles              []*Bundle
		expectedList         []*TxOrGen
		expectedTargetTxHash []common.Hash
	}{
		{
			name: "correctly extract bundles with same target hash 1",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(txs[0], txs[1]), common.Hash{}, false),
				NewBundle(NewTxOrGenList(g1), common.Hash{}, true),
			},
			expectedList:         NewTxOrGenList(g1, txs[0], txs[1], txs[2], txs[3], txs[4], txs[5]),
			expectedTargetTxHash: []common.Hash{g1.Id, {}},
		},
		{
			name: "correctly extract bundles with same target hash 2",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(g1, txs[2], txs[3]), txs[1].Hash(), false),
				NewBundle(NewTxOrGenList(g2), txs[1].Hash(), true),
			},
			expectedList:         NewTxOrGenList(txs[0], txs[1], g2, g1, txs[2], txs[3], txs[4], txs[5]),
			expectedTargetTxHash: []common.Hash{g2.Id, txs[1].Hash()},
		},
		{
			name: "correctly extract bundles with same target hash 3",
			bundles: []*Bundle{
				NewBundle(NewTxOrGenList(g1), txs[0].Hash(), true),
				NewBundle(NewTxOrGenList(g2, txs[2], txs[3]), txs[1].Hash(), false),
				NewBundle(NewTxOrGenList(g3), txs[1].Hash(), true),
			},
			expectedList:         NewTxOrGenList(txs[0], g1, txs[1], g3, g2, txs[2], txs[3], txs[4], txs[5]),
			expectedTargetTxHash: []common.Hash{txs[0].Hash(), g3.Id, txs[1].Hash()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := IncorporateBundleTx(txs, tc.bundles)
			bundles := SetCorrectTargetTxHash(tc.bundles, ret)
			require.Nil(t, err)
			require.Equal(t, len(tc.expectedList), len(ret))
			for i := range ret {
				assert.Equal(t, tc.expectedList[i].Id.Hex(), ret[i].Id.Hex(), "mismatch at ret[%d]", i)
			}
			for i := range bundles {
				assert.Equal(t, tc.expectedTargetTxHash[i], bundles[i].TargetTxHash)
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
	txOrGenList := NewTxOrGenList(txs[0], txs[1], txs[2])
	testCases := []struct {
		name     string
		bundle   *Bundle
		expected []*TxOrGen
	}{
		{
			name:     "incorporate first two transactions",
			bundle:   NewBundle(NewTxOrGenList(txs[0], txs[1]), common.Hash{}, false),
			expected: NewTxOrGenList(txs[0], txs[1], txs[2]),
		},
		{
			name:     "incorporate last two transactions",
			bundle:   NewBundle(NewTxOrGenList(txs[1], txs[2]), common.Hash{}, false),
			expected: NewTxOrGenList(txs[1], txs[2], txs[0]),
		},
		{
			name:     "incorporate with target hash",
			bundle:   NewBundle(NewTxOrGenList(txs[0]), txs[2].Hash(), false),
			expected: NewTxOrGenList(txs[1], txs[2], txs[0]),
		},
		{
			name:     "incorporate single transaction",
			bundle:   NewBundle(NewTxOrGenList(txs[2]), common.Hash{}, false),
			expected: NewTxOrGenList(txs[2], txs[0], txs[1]),
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

	b0 := NewBundle(NewTxOrGenList(txs[0], txs[1]), common.Hash{}, true)
	defaultTargetHash := txs[1].Hash() // make TargetTxHash checks pass

	testcases := []struct {
		name        string
		prevBundles []*Bundle
		newBundles  []*Bundle
		expected    bool
	}{
		{
			name:        "Same TargetTxHash, TargetRequired",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(), common.Hash{}, true)},
			expected:    true,
		},
		{
			name:        "Same TargetTxHash, TargetRequired=false",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(), common.Hash{}, false)},
			expected:    false,
		},
		{
			name:        "TargetTxHash divides a bundle",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(), txs[0].Hash(), false)},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 1",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(txs[0], txs[2]), defaultTargetHash, false)},
			expected:    true,
		},
		{
			name:        "Overlapping BundleTxs 2",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(txs[1], txs[2], txs[3]), defaultTargetHash, false)},
			expected:    true,
		},
		{
			name:        "Non-overlapping BundleTxs",
			prevBundles: []*Bundle{b0},
			newBundles:  []*Bundle{NewBundle(NewTxOrGenList(txs[2], txs[3]), defaultTargetHash, false)},
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
	g1 := NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := NewTxOrGenFromGen(gen, common.Hash{2})
	g3 := NewTxOrGenFromGen(gen, common.Hash{3})

	for i := range keys {
		keys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
	}

	for i := range txs {
		addr, key := addrs[i/2], keys[i/2]
		txs[i], _ = types.SignTx(types.NewTransaction(uint64(i), addr, big.NewInt(1), 21000, big.NewInt(1), nil), signer, key)
	}

	// Create test bundles
	bundles := []*Bundle{
		NewBundle(NewTxOrGenList(txs[1], txs[2]), common.Hash{}, false),
		NewBundle(NewTxOrGenList(txs[3], txs[4]), txs[2].Hash(), false),
		NewBundle(NewTxOrGenList(g1, txs[5]), txs[4].Hash(), false),
		NewBundle(NewTxOrGenList(g1, txs[1], txs[2]), common.Hash{}, false),
		NewBundle(NewTxOrGenList(txs[2]), txs[1].Hash(), true),  // If target is popped, the bundle should be popped
		NewBundle(NewTxOrGenList(txs[2]), txs[1].Hash(), false), // Bundle is not popped if target is popped
		NewBundle(NewTxOrGenList(g3), txs[2].Hash(), true),      // If target in bundle is popped, the bundle should be popped
	}

	testCases := []struct {
		name            string
		incorporatedTxs []*TxOrGen
		numToPop        int
		bundles         []*Bundle
		expectedTxs     []*TxOrGen
	}{
		{
			name:            "Without any dependencies",
			incorporatedTxs: NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        1,
			bundles:         []*Bundle{},
			expectedTxs:     NewTxOrGenList(txs[2], txs[3], txs[4], txs[5]),
		},
		{
			name:            "No bundles, tx0 and tx1 dependency (same sender)",
			incorporatedTxs: NewTxOrGenList(txs[0], txs[1], txs[2], txs[3], txs[4]),
			numToPop:        1,
			bundles:         []*Bundle{},
			expectedTxs:     NewTxOrGenList(txs[2], txs[3], txs[4]),
		},
		{
			name:            "One bundle - first tx is generator",
			incorporatedTxs: NewTxOrGenList(g1, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        1,
			bundles:         []*Bundle{bundles[3]},
			expectedTxs:     NewTxOrGenList(txs[4], txs[5], txs[6]),
		},
		{
			name:            "Two bundles - chaining dependency",
			incorporatedTxs: NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        2,
			bundles:         []*Bundle{bundles[0], bundles[1]},
			expectedTxs:     NewTxOrGenList(),
		},
		{
			name:            "Two bundles - one independent tx (tx6)",
			incorporatedTxs: NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*Bundle{bundles[0], bundles[1]},
			expectedTxs:     NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - change order",
			incorporatedTxs: NewTxOrGenList(txs[2], txs[3], txs[4], txs[6], txs[5]), // 6 is before 5
			numToPop:        2,
			bundles:         []*Bundle{bundles[0], bundles[1]},
			expectedTxs:     NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - one independent tx (tx6) with one generator",
			incorporatedTxs: NewTxOrGenList(g1, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*Bundle{bundles[0], bundles[1]},
			expectedTxs:     NewTxOrGenList(txs[6]),
		},
		{
			name:            "Two bundles - two generators",
			incorporatedTxs: NewTxOrGenList(g1, g2, txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
			numToPop:        2,
			bundles:         []*Bundle{bundles[0], bundles[1]},
			expectedTxs:     NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5], txs[6]),
		},
		{
			name:            "Two bundles - target is popped",
			incorporatedTxs: NewTxOrGenList(g1, txs[1], txs[2], g3, txs[3], txs[4], txs[5], txs[6]),
			numToPop:        3,
			bundles:         []*Bundle{bundles[3], bundles[6]},
			expectedTxs:     NewTxOrGenList(txs[4], txs[5], txs[6]),
		},
		{
			name:            "One bundle, target is popped with TargetRequired=true",
			incorporatedTxs: NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        1,
			bundles:         []*Bundle{bundles[4]},
			expectedTxs:     NewTxOrGenList(txs[4], txs[5]),
		},
		{
			name:            "One bundle, target is popped with TargetRequired=false",
			incorporatedTxs: NewTxOrGenList(txs[1], txs[2], txs[3], txs[4], txs[5]),
			numToPop:        1,
			bundles:         []*Bundle{bundles[5]},
			expectedTxs:     NewTxOrGenList(txs[2], txs[3], txs[4], txs[5]),
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

type mockTxbundlingModule struct{}

func (m *mockTxbundlingModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*Bundle) []*Bundle {
	return nil
}

func (m *mockTxbundlingModule) FilterTxs(txs map[common.Address]types.Transactions) {
	for addr, list := range txs {
		for i, tx := range list {
			if tx.Nonce() >= 3 {
				if i == 0 {
					delete(txs, addr)
				} else {
					txs[addr] = list[:i]
				}
				break
			}
		}
	}
}

func TestFilterTxs(t *testing.T) {
	txs := make([]*types.Transaction, 5)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{}, common.Big0, 0, common.Big0, nil)
	}

	txsMap := map[common.Address]types.Transactions{
		{1}: {txs[0], txs[1]},
		{2}: {txs[2], txs[3]},
		{3}: {txs[4]},
	}

	mockTxBundlingModule := &mockTxbundlingModule{}
	FilterTxs(txsMap, []TxBundlingModule{mockTxBundlingModule})

	assert.Equal(t, txsMap, map[common.Address]types.Transactions{
		{1}: {txs[0], txs[1]},
		{2}: {txs[2]},
	})
}

func TestCoordinateTargetTxHash(t *testing.T) {
	txs := make([]*types.Transaction, 6)
	for i := range txs {
		txs[i] = types.NewTransaction(uint64(i), common.Address{}, common.Big0, 0, common.Big0, nil)
	}

	gen := func(nonce uint64) (*types.Transaction, error) {
		return types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
	}
	g1 := NewTxOrGenFromGen(gen, common.Hash{1})
	g2 := NewTxOrGenFromGen(gen, common.Hash{2})
	g3 := NewTxOrGenFromGen(gen, common.Hash{3})
	g4 := NewTxOrGenFromGen(gen, common.Hash{4})

	b0 := NewBundle(NewTxOrGenList(g1, txs[3], txs[4]), txs[2].Hash(), false)
	b1 := NewBundle(NewTxOrGenList(g2), txs[2].Hash(), true)
	b2 := NewBundle(NewTxOrGenList(g3), txs[2].Hash(), false)
	b3 := NewBundle(NewTxOrGenList(g4, txs[5]), txs[4].Hash(), false)

	// Copy bundles to avoid modifying the original bundles
	copyBundles := func(bundles []*Bundle) []*Bundle {
		copiedBundles := make([]*Bundle, len(bundles))
		for i, bundle := range bundles {
			copiedBundles[i] = NewBundle(bundle.BundleTxs, bundle.TargetTxHash, bundle.TargetRequired)
		}
		return copiedBundles
	}

	compareBundle := func(a, b *Bundle) bool {
		for _, tx := range a.BundleTxs {
			if !b.Has(tx) {
				return false
			}
		}
		return true
	}

	testCases := []struct {
		name                   string
		bundles                []*Bundle
		expectedTargetTxHashes []common.Hash // Same order as bundles
		expectedTxs            []*TxOrGen
	}{
		{
			name:                   "One bundles: True",
			bundles:                copyBundles([]*Bundle{b1}),
			expectedTargetTxHashes: []common.Hash{txs[2].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g2, txs[3], txs[4], txs[5]),
		},
		{
			name:                   "One bundles: False",
			bundles:                copyBundles([]*Bundle{b0}),
			expectedTargetTxHashes: []common.Hash{txs[2].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g1, txs[3], txs[4], txs[5]),
		},
		{
			name:                   "Two bundles: True, False",
			bundles:                copyBundles([]*Bundle{b1, b0}),
			expectedTargetTxHashes: []common.Hash{txs[2].Hash(), g2.Id},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g2, g1, txs[3], txs[4], txs[5]),
		},
		{
			name:                   "Two bundles: False, True",
			bundles:                copyBundles([]*Bundle{b0, b1}),
			expectedTargetTxHashes: []common.Hash{g2.Id, txs[2].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g2, g1, txs[3], txs[4], txs[5]),
		},
		{
			name:                   "Two bundles: False, False",
			bundles:                copyBundles([]*Bundle{b0, b2}),
			expectedTargetTxHashes: []common.Hash{txs[2].Hash(), txs[4].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g1, txs[3], txs[4], g3, txs[5]),
		},
		{
			name:                   "Three bundles with same target tx hash",
			bundles:                copyBundles([]*Bundle{b0, b1, b2}),
			expectedTargetTxHashes: []common.Hash{g2.Id, txs[2].Hash(), txs[4].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g2, g1, txs[3], txs[4], g3, txs[5]),
		},
		{
			name:                   "Three bundles with different target tx hash",
			bundles:                copyBundles([]*Bundle{b0, b1, b3}),
			expectedTargetTxHashes: []common.Hash{g2.Id, txs[2].Hash(), txs[4].Hash()},
			expectedTxs:            NewTxOrGenList(txs[0], txs[1], txs[2], g2, g1, txs[3], txs[4], g4, txs[5]),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coordinatedBundles := coordinateTargetTxHash(tc.bundles)
			for _, bundle := range coordinatedBundles {
				idx := 0
				for _, b := range tc.bundles {
					if compareBundle(b, bundle) {
						break
					}
					idx++
				}

				assert.Equal(t, tc.expectedTargetTxHashes[idx], bundle.TargetTxHash)
			}

			incorporatedTxs, err := IncorporateBundleTx(txs, coordinatedBundles)
			require.Nil(t, err)
			assert.Equal(t, tc.expectedTxs, incorporatedTxs)
		})
	}
}
