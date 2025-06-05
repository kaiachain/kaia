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

	"github.com/kaiachain/kaia/v2/accounts/abi/bind"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/system"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/common/hexutil"
	contracts "github.com/kaiachain/kaia/v2/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/stretchr/testify/require"
)

func TestPostInsertBlock(t *testing.T) {
	alloc := testAllocStorage()
	senderkey, _ := crypto.GenerateKey()
	senderAddr := crypto.PubkeyToAddress(senderkey.PublicKey)
	anotherGSR := common.HexToAddress("0x2345")

	alloc[senderAddr] = blockchain.GenesisAccount{
		Balance: new(big.Int).SetUint64(params.KAIA),
	}
	alloc[anotherGSR] = blockchain.GenesisAccount{
		Code:    hexutil.MustDecode(dummyGSRCode),
		Balance: big.NewInt(0),
		Storage: map[common.Hash]common.Hash{
			common.HexToHash("0x0"): common.HexToHash("0x1"),
			common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"): dummyTokenAddress1.Hash(),
		},
	}

	nodekey, _ := crypto.GenerateKey()

	tcs := map[string]struct {
		getCurrentBlock func() *GaslessModule
		swapRouter      common.Address
		tokens          []common.Address
	}{
		"add new token address": {
			func() *GaslessModule {
				g := NewGaslessModule()
				dbm := database.NewMemoryDBManager()
				backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
				err := g.Init(&InitOpts{
					ChainConfig:   testChainConfig,
					GaslessConfig: testGaslessConfig,
					NodeKey:       nodekey,
					Chain:         backend.BlockChain(),
					TxPool:        &testTxPool{},
					NodeType:      common.ENDPOINTNODE,
				})
				require.NoError(t, err)

				data := common.Hex2Bytes("d48bfca7000000000000000000000000000000000000000000000000000000000000ffff") // addToken(0x000000000000000000000000000000000000ffff)
				tx := types.NewTransaction(0, dummyGSRAddress, big.NewInt(0), 50000, big.NewInt(25000000000), data)
				signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(testChainConfig.ChainID), senderkey)
				require.NoError(t, err)

				err = backend.SendTransaction(context.Background(), signedTx)
				require.NoError(t, err)

				backend.Commit()
				return g
			},
			dummyGSRAddress,
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2, dummyTokenAddress3, common.HexToAddress("0xffff")},
		},
		"remove token address": {
			func() *GaslessModule {
				g := NewGaslessModule()
				dbm := database.NewMemoryDBManager()
				backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
				err := g.Init(&InitOpts{
					ChainConfig:   testChainConfig,
					GaslessConfig: testGaslessConfig,
					NodeKey:       nodekey,
					Chain:         backend.BlockChain(),
					TxPool:        &testTxPool{},
					NodeType:      common.ENDPOINTNODE,
				})
				require.NoError(t, err)
				data := common.Hex2Bytes("5fa7b584000000000000000000000000000000000000000000000000000000000000abcd") // removeToken(dummyTokenAddress1)
				tx := types.NewTransaction(0, dummyGSRAddress, big.NewInt(0), 50000, big.NewInt(25000000000), data)
				signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(testChainConfig.ChainID), senderkey)
				require.NoError(t, err)

				err = backend.SendTransaction(context.Background(), signedTx)
				require.NoError(t, err)

				backend.Commit()
				return g
			},
			dummyGSRAddress,
			[]common.Address{dummyTokenAddress2, dummyTokenAddress3},
		},
		"update gsr address": {
			func() *GaslessModule {
				g := NewGaslessModule()
				dbm := database.NewMemoryDBManager()
				backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)
				err := g.Init(&InitOpts{
					ChainConfig:   testChainConfig,
					GaslessConfig: testGaslessConfig,
					NodeKey:       nodekey,
					Chain:         backend.BlockChain(),
					TxPool:        &testTxPool{},
					NodeType:      common.ENDPOINTNODE,
				})
				require.NoError(t, err)

				sender := bind.NewKeyedTransactor(senderkey)
				contract, _ := contracts.NewRegistryMockTransactor(system.RegistryAddr, backend)
				contract.Register(sender, GaslessSwapRouterName, anotherGSR, big.NewInt(2))

				backend.Commit()
				backend.Commit()
				return g
			},
			anotherGSR,
			[]common.Address{dummyTokenAddress1},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			g := tc.getCurrentBlock()
			err := g.PostInsertBlock(g.Chain.CurrentBlock())
			require.NoError(t, err)

			require.Equal(t, tc.swapRouter, g.swapRouter)

			require.Equal(t, len(tc.tokens), len(g.allowedTokens))
			for _, addr := range tc.tokens {
				_, ok := g.allowedTokens[addr]
				require.True(t, ok)
			}
		})
	}
}
