package work

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/builder"
)

func TestCheckBundlesComplete(t *testing.T) {
	// Create test transactions and bundles
	tx1 := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	tx2 := types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	tx3 := types.NewTransaction(2, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)

	bundle1 := &builder.Bundle{
		BundleTxs:    []interface{}{tx1, tx2},
		TargetTxHash: tx1.Hash(),
	}

	bundle2 := &builder.Bundle{
		BundleTxs:    []interface{}{tx3},
		TargetTxHash: tx2.Hash(),
	}

	bundles := []*builder.Bundle{bundle1, bundle2}

	testCases := []struct {
		name            string
		incorporatedTxs []interface{}
		expectedBundles int
	}{
		{
			name:            "All bundles incomplete",
			incorporatedTxs: []interface{}{tx1},
			expectedBundles: 2,
		},
		{
			name:            "Bundle1 complete, Bundle2 incomplete",
			incorporatedTxs: []interface{}{tx1, tx2},
			expectedBundles: 1,
		},
		{
			name:            "All bundles complete",
			incorporatedTxs: []interface{}{tx1, tx2, tx3},
			expectedBundles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incompleteBundles := getIncompleteBundles(tc.incorporatedTxs, bundles)
			if len(incompleteBundles) != tc.expectedBundles {
				t.Errorf("Expected %d incomplete bundles, got %d", tc.expectedBundles, len(incompleteBundles))
			}
		})
	}

}

func TestRemoveIncompleteBundle(t *testing.T) {
	// Create test transactions and bundles
	tx1 := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	tx2 := types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	tx3 := types.NewTransaction(2, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	tx4 := types.NewTransaction(3, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	signer := types.NewEIP155Signer(big.NewInt(1))

	bundle1 := &builder.Bundle{
		BundleTxs:    []interface{}{tx1, tx2},
		TargetTxHash: tx1.Hash(),
	}

	bundle2 := &builder.Bundle{
		BundleTxs:    []interface{}{tx3},
		TargetTxHash: tx3.Hash(),
	}

	testCases := []struct {
		name              string
		incorporatedTxs   []interface{}
		incompleteBundles []*builder.Bundle
		expectedTxs       []interface{}
		expectedTxsLength int
	}{
		{
			name:              "Remove transactions from incomplete bundle1",
			incorporatedTxs:   []interface{}{tx1, tx2, tx3, tx4},
			incompleteBundles: []*builder.Bundle{bundle1},
			expectedTxs:       []interface{}{tx3, tx4},
			expectedTxsLength: 2,
		},
		{
			name:              "Remove transactions from incomplete bundle2",
			incorporatedTxs:   []interface{}{tx1, tx2, tx3, tx4},
			incompleteBundles: []*builder.Bundle{bundle2},
			expectedTxs:       []interface{}{tx1, tx2, tx4},
			expectedTxsLength: 3,
		},

		{
			name:              "Remove transactions from both incomplete bundles",
			incorporatedTxs:   []interface{}{tx1, tx2, tx3, tx4},
			incompleteBundles: []*builder.Bundle{bundle1, bundle2},
			expectedTxs:       []interface{}{tx4},
			expectedTxsLength: 1,
		},

		{
			name:              "No incomplete bundles",
			incorporatedTxs:   []interface{}{tx1, tx2, tx3, tx4},
			incompleteBundles: []*builder.Bundle{},
			expectedTxs:       []interface{}{tx1, tx2, tx3, tx4},
			expectedTxsLength: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incorporatedTxs := make([]interface{}, len(tc.incorporatedTxs))
			copy(incorporatedTxs, tc.incorporatedTxs)

			removeTxsFromIncompleteBundles(&incorporatedTxs, tc.incompleteBundles, signer)

			if len(incorporatedTxs) != tc.expectedTxsLength {
				t.Errorf("Expected %d transactions, got %d", tc.expectedTxsLength, len(incorporatedTxs))
			}

			for _, expectedTx := range tc.expectedTxs {
				found := false
				for _, remainingTx := range incorporatedTxs {
					if tx1, ok1 := expectedTx.(*types.Transaction); ok1 {
						if tx2, ok2 := remainingTx.(*types.Transaction); ok2 {
							if tx1.Hash() == tx2.Hash() {
								found = true
								break
							}
						}
					}
				}
				if !found {
					t.Errorf("Expected transaction not found in remaining transactions")
				}
			}
		})
	}
}

func TestPopXXXX(t *testing.T) {
	// Setup test transactions
	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	signer := types.NewEIP155Signer(big.NewInt(1))

	tx1 := types.NewTransaction(0, addr1, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx1, _ = types.SignTx(tx1, signer, key1)
	tx2 := types.NewTransaction(1, addr1, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx2, _ = types.SignTx(tx2, signer, key1)
	tx3 := types.NewTransaction(2, addr2, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx3, _ = types.SignTx(tx3, signer, key2)
	tx4 := types.NewTransaction(3, addr2, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx4, _ = types.SignTx(tx4, signer, key2)

	t.Log(addr1.Hex(), addr2.Hex())
	a1, _ := types.Sender(signer, tx1)
	a2, _ := types.Sender(signer, tx2)
	a3, _ := types.Sender(signer, tx3)
	a4, _ := types.Sender(signer, tx4)
	t.Log(a1.Hex(), a2.Hex(), a3.Hex(), a4.Hex())

	testCases := []struct {
		name            string
		incorporatedTxs []interface{}
		numToPop        int
		expectedTxs     []interface{}
	}{
		{
			name:            "Pop single transaction",
			incorporatedTxs: []interface{}{tx1, tx2, tx3, tx4},
			numToPop:        1,
			expectedTxs:     []interface{}{tx3, tx4},
		},
		{
			name:            "Pop multiple transactions",
			incorporatedTxs: []interface{}{tx1, tx2, tx3, tx4},
			numToPop:        3,
			expectedTxs:     []interface{}{},
		},
		{
			name:            "Pop zero transactions",
			incorporatedTxs: []interface{}{tx1, tx2, tx3},
			numToPop:        0,
			expectedTxs:     []interface{}{tx1, tx2, tx3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			popXXXX(&tc.incorporatedTxs, tc.numToPop, signer)

			if len(tc.incorporatedTxs) != len(tc.expectedTxs) {
				t.Errorf("Expected %d transactions, got %d", len(tc.expectedTxs), len(tc.incorporatedTxs))
			}
		})
	}
}

func TestPopTxs(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	key3, _ := crypto.GenerateKey()
	key4, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	addr3 := crypto.PubkeyToAddress(key3.PublicKey)
	addr4 := crypto.PubkeyToAddress(key4.PublicKey)
	signer := types.NewEIP155Signer(big.NewInt(1))

	// Create test transactions
	tx1 := types.NewTransaction(0, addr1, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx1, _ = types.SignTx(tx1, signer, key1)
	tx2 := types.NewTransaction(1, addr1, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx2, _ = types.SignTx(tx2, signer, key1)
	tx3 := types.NewTransaction(2, addr2, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx3, _ = types.SignTx(tx3, signer, key2)
	tx4 := types.NewTransaction(3, addr2, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx4, _ = types.SignTx(tx4, signer, key2)
	tx5 := types.NewTransaction(4, addr3, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx5, _ = types.SignTx(tx5, signer, key3)
	tx6 := types.NewTransaction(5, addr3, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx6, _ = types.SignTx(tx6, signer, key3)
	tx7 := types.NewTransaction(6, addr4, big.NewInt(1), 21000, big.NewInt(1), nil)
	tx7, _ = types.SignTx(tx7, signer, key4)

	a3, _ := types.Sender(signer, tx3)
	t.Log(a3.Hex())

	// Create test bundles
	bundle1 := &builder.Bundle{
		BundleTxs: []interface{}{tx2, tx3},
	}
	bundle2 := &builder.Bundle{
		BundleTxs:    []interface{}{tx4, tx5},
		TargetTxHash: tx3.Hash(),
	}

	testCases := []struct {
		name            string
		incorporatedTxs []interface{}
		numToPop        int
		bundles         []*builder.Bundle
		expectedTxs     []interface{}
	}{
		{
			name:            "Pop single transaction with no bundles",
			incorporatedTxs: []interface{}{tx2, tx3, tx4, tx5, tx6},
			numToPop:        2,
			bundles:         []*builder.Bundle{bundle1, bundle2},
			expectedTxs:     []interface{}{},
		},
		{
			name:            "Pop single transaction with no bundles",
			incorporatedTxs: []interface{}{tx2, tx3, tx4, tx5, tx7, tx6},
			numToPop:        2,
			bundles:         []*builder.Bundle{bundle1, bundle2},
			expectedTxs:     []interface{}{tx7},
		},
		/*
			{
				name:            "Pop transaction from bundle",
				incorporatedTxs: []interface{}{tx1, tx2, tx3, tx4},
				numToPop:        1,
				bundles:         []*builder.Bundle{bundle1},
				expectedTxs:     []interface{}{tx3, tx4},
			},
			{
				name:            "Pop multiple transactions with multiple bundles",
				incorporatedTxs: []interface{}{tx1, tx2, tx3, tx4},
				numToPop:        2,
				bundles:         []*builder.Bundle{bundle1, bundle2},
				expectedTxs:     []interface{}{},
			},
			{
				name:            "Pop zero transactions",
				incorporatedTxs: []interface{}{tx1, tx2, tx3},
				numToPop:        0,
				bundles:         []*builder.Bundle{bundle1},
				expectedTxs:     []interface{}{tx1, tx2, tx3},
			},
		*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			popTxs(&tc.incorporatedTxs, tc.numToPop, tc.bundles, signer)

			if len(tc.incorporatedTxs) != len(tc.expectedTxs) {
				t.Errorf("Expected %d transactions, got %d", len(tc.expectedTxs), len(tc.incorporatedTxs))
			}
		})
	}
}
