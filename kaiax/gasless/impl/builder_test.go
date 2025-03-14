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
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestExtractTxBundles(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)

	g := NewGaslessModule()
	dbm := database.NewMemoryDBManager()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(dbm), nil, nil)
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
	nodeKey, _ := crypto.GenerateKey()
	disabled, err := g.Init(&InitOpts{
		ChainConfig: testChainConfig,
		NodeKey:     nodeKey,
		Chain:       backend.BlockChain(),
		TxPool:      &testTxPool{sdb},
	})
	require.NoError(t, err)
	require.False(t, disabled)

	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()

	A1 := makeApproveTx(t, key1, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)})
	S1 := makeSwapTx(t, key1, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)})

	A2 := makeApproveTx(t, key2, 0, ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)})
	S2 := makeSwapTx(t, key2, 1, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)})

	S3 := makeSwapTx(t, nil, 0, SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(1021000)})

	T4 := makeTx(t, nil, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)
	T5 := makeTx(t, nil, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(1), nil)

	testcases := []struct {
		pending  []*types.Transaction
		pre      []*builder.Bundle
		expected []*builder.Bundle
	}{
		{
			[]*types.Transaction{A1, S1, T4, T5},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
			},
		},
		{
			[]*types.Transaction{A1, T4, S1, T5},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: T4.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{T4, A1, S1, T5},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: T4.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, S1, A2, T4, S2},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: T4.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, A2, S1, S2},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: S1.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, S1, A2, S2},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: S1.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, A2, S2, S1},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: S2.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, S1, S3},
			nil,
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(nil, S3), S3},
					TargetTxHash: S1.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, S1, T4, T5},
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{},
					TargetTxHash: common.Hash{},
				},
			},
			[]*builder.Bundle{},
		},
	}

	for _, tc := range testcases {
		bundles := g.ExtractTxBundles(tc.pending, tc.pre)
		require.Equal(t, len(tc.expected), len(bundles))

		for i, e := range tc.expected {
			// check TargetTxHash
			require.Equal(t, e.TargetTxHash.String(), bundles[i].TargetTxHash.String())

			// check BundleTxs
			require.Equal(t, len(e.BundleTxs), len(bundles[i].BundleTxs))
			ehashes, err := flattenBundleTxs(e.BundleTxs)
			require.NoError(t, err)
			hashes, err := flattenBundleTxs(bundles[i].BundleTxs)
			require.NoError(t, err)
			for j, ehash := range ehashes {
				require.Equal(t, ehash, hashes[j])
			}
		}
	}
}
