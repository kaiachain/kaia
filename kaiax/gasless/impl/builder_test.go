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
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestExtractTxBundles(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)

	g := NewGaslessModule()
	nodeKey, _ := crypto.GenerateKey()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     nodeKey,
		StateDB:     statedb,
	})

	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	key3, _ := crypto.GenerateKey()

	A1, err := makeTx(key1, 0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240"))
	require.NoError(t, err)
	S1, err := makeTx(key1, 1, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000001ed688"))
	require.NoError(t, err)
	A2, err := makeTx(key2, 0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240"))
	require.NoError(t, err)
	S2, err := makeTx(key2, 1, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000001ed688"))
	require.NoError(t, err)
	T1, err := makeTx(key3, 0, common.HexToAddress("0xAAAA"), big.NewInt(0), 1000000, big.NewInt(0), nil)
	require.NoError(t, err)
	T2, err := makeTx(key3, 0, common.HexToAddress("0xBBBB"), big.NewInt(0), 1000000, big.NewInt(0), nil)
	require.NoError(t, err)

	testcases := []struct {
		pending  []*types.Transaction
		expected []*builder.Bundle
	}{
		{
			[]*types.Transaction{A1, S1, T1, T2},
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
			},
		},
		{
			[]*types.Transaction{A1, T1, S1, T2},
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: T1.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, S1, A2, T1, S2},
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: T1.Hash(),
				},
			},
		},
		{
			[]*types.Transaction{A1, A2, S1, S2},
			[]*builder.Bundle{
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A1, S1), A1, S1},
					TargetTxHash: common.Hash{},
				},
				{
					BundleTxs:    []interface{}{g.GetLendTxGenerator(A2, S2), A2, S2},
					TargetTxHash: common.Hash{},
				},
			},
		},
		// GASLESS_TODO: conflict test
	}

	for _, tc := range testcases {
		bundles := g.ExtractTxBundles(tc.pending, nil)
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

// helper

type AccountKeyPickerForTest struct {
	AddrKeyMap map[common.Address]accountkey.AccountKey
}

func (a *AccountKeyPickerForTest) GetKey(addr common.Address) accountkey.AccountKey {
	return a.AddrKeyMap[addr]
}

func (a *AccountKeyPickerForTest) SetKey(addr common.Address, key accountkey.AccountKey) {
	a.AddrKeyMap[addr] = key
}

func (a *AccountKeyPickerForTest) Exist(addr common.Address) bool {
	return a.AddrKeyMap[addr] != nil
}

func makeTx(privKey *ecdsa.PrivateKey, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) (*types.Transaction, error) {
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	p := &AccountKeyPickerForTest{
		AddrKeyMap: make(map[common.Address]accountkey.AccountKey),
	}
	p.SetKey(addr, accountkey.NewAccountKeyLegacy())

	signer := types.LatestSignerForChainID(big.NewInt(1))
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	tx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, err
	}
	_, err = tx.ValidateSender(signer, p, 0)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func flattenBundleTxs(txs []interface{}) ([]common.Hash, error) {
	nodeNonce := uint64(0)
	hashes := []common.Hash{}
	for _, txi := range txs {
		var tx *types.Transaction
		var err error
		if genLendTx, ok := txi.(builder.TxGenerator); ok {
			tx, err = genLendTx(nodeNonce)
			if err != nil {
				return nil, err
			}
			nodeNonce += 1
		} else if tx, ok = txi.(*types.Transaction); ok {
		} else {
			err = errors.New("unsupported bundle tx")
		}
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, tx.Hash())
	}
	return hashes, nil
}
