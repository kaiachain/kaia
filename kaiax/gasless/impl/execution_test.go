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
	"context"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestPostInsertBlock(t *testing.T) {
	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	senderkey, _ := crypto.GenerateKey()
	senderAddr := crypto.PubkeyToAddress(senderkey.PublicKey)
	alloc[senderAddr] = blockchain.GenesisAccount{
		Balance: new(big.Int).SetUint64(params.KAIA),
	}
	nodekey, _ := crypto.GenerateKey()

	tcs := map[string]struct {
		getCurrentBlock func() *GaslessModule
		expected        []common.Address
	}{
		"add new token address": {
			func() *GaslessModule {
				g := NewGaslessModule()
				backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
				disabled, err := g.Init(&InitOpts{
					ChainConfig: testChainConfig,
					CNConfig:    testCNConfig,
					NodeKey:     nodekey,
					Chain:       backend.BlockChain(),
					TxPool:      &testTxPool{},
				})
				require.NoError(t, err)
				require.False(t, disabled)

				data := common.Hex2Bytes("d48bfca7000000000000000000000000000000000000000000000000000000000000ffff") // addToken(0x000000000000000000000000000000000000ffff)
				tx := types.NewTransaction(0, dummyGSRAddress, big.NewInt(0), 50000, big.NewInt(25000000000), data)
				signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(testChainConfig.ChainID), senderkey)
				require.NoError(t, err)

				err = backend.SendTransaction(context.Background(), signedTx)
				require.NoError(t, err)

				backend.Commit()
				return g
			},
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2, dummyTokenAddress3, common.HexToAddress("0xffff")},
		},
		"remove token address": {
			func() *GaslessModule {
				g := NewGaslessModule()
				backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
				disabled, err := g.Init(&InitOpts{
					ChainConfig: testChainConfig,
					CNConfig:    testCNConfig,
					NodeKey:     nodekey,
					Chain:       backend.BlockChain(),
					TxPool:      &testTxPool{},
				})
				require.NoError(t, err)
				require.False(t, disabled)

				data := common.Hex2Bytes("5fa7b584000000000000000000000000000000000000000000000000000000000000abcd") // removeToken(dummyTokenAddress1)
				tx := types.NewTransaction(0, dummyGSRAddress, big.NewInt(0), 50000, big.NewInt(25000000000), data)
				signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(testChainConfig.ChainID), senderkey)
				require.NoError(t, err)

				err = backend.SendTransaction(context.Background(), signedTx)
				require.NoError(t, err)

				backend.Commit()
				return g
			},
			[]common.Address{dummyTokenAddress2, dummyTokenAddress3},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			g := tc.getCurrentBlock()
			err := g.PostInsertBlock(g.Chain.CurrentBlock())
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(g.allowedTokens))
			for _, addr := range tc.expected {
				_, ok := g.allowedTokens[addr]
				require.True(t, ok)
			}
		})
	}
}
