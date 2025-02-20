package impl

import (
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/stretchr/testify/require"
)

func makeBundle(nonces []uint64) *builder.Bundle {
	bundleTxs := make([]interface{}, len(nonces))
	to := common.Address{}
	for i, nonce := range nonces {
		bundleTxs[i] = types.NewTransaction(nonce, to, common.Big0, 0, common.Big0, nil)
	}
	return &builder.Bundle{
		BundleTxs:    bundleTxs,
		TargetTxHash: common.Hash{},
	}
}

func TestIsConflict(t *testing.T) {
	testcases := []struct {
		name        string
		prevBundles []*builder.Bundle
		newBundles  []*builder.Bundle
		expected    bool
	}{
		{
			name:        "Non-overlapping txs",
			prevBundles: []*builder.Bundle{makeBundle([]uint64{1, 2})},
			newBundles:  []*builder.Bundle{makeBundle([]uint64{3, 4})},
			expected:    false,
		},
		{
			name:        "Overlapping tx 1",
			prevBundles: []*builder.Bundle{makeBundle([]uint64{1, 2})},
			newBundles:  []*builder.Bundle{makeBundle([]uint64{1, 3})},
			expected:    true,
		},
		{
			name:        "Overlapping tx 2",
			prevBundles: []*builder.Bundle{makeBundle([]uint64{1, 2})},
			newBundles:  []*builder.Bundle{makeBundle([]uint64{2, 3, 4})},
			expected:    true,
		},
		{
			name:        "Overlapping tx 3",
			prevBundles: []*builder.Bundle{makeBundle([]uint64{1, 2}), makeBundle([]uint64{3, 4})},
			newBundles:  []*builder.Bundle{makeBundle([]uint64{3})},
			expected:    true,
		},
	}

	b := NewBuilderModule()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotConflict := b.IsConflict(tc.prevBundles, tc.newBundles)
			require.Equal(t, tc.expected, gotConflict)
		})
	}
}
