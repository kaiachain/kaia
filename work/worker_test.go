package work

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/builder"
)

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
			name:            "Pop transaction with dependencies",
			incorporatedTxs: []interface{}{tx2, tx3, tx4, tx5, tx6},
			numToPop:        2,
			bundles:         []*builder.Bundle{bundle1, bundle2},
			expectedTxs:     []interface{}{},
		},
		{
			name:            "Pop transaction with dependencies",
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
			popTxs(&tc.incorporatedTxs, tc.numToPop, &tc.bundles, signer)

			if len(tc.incorporatedTxs) != len(tc.expectedTxs) {
				t.Errorf("Expected %d transactions, got %d", len(tc.expectedTxs), len(tc.incorporatedTxs))
			}
		})
	}
}
