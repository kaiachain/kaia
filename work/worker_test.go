package work

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
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
			incompleteBundles := checkBundlesComplete(tc.incorporatedTxs, bundles)
			if len(incompleteBundles) != tc.expectedBundles {
				t.Errorf("Expected %d incomplete bundles, got %d", tc.expectedBundles, len(incompleteBundles))
			}
		})
	}

}
