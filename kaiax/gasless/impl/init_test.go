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
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestInitSuccess(t *testing.T) {
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	key, _ := crypto.GenerateKey()

	tcs := map[string]struct {
		getOpts  func() *InitOpts
		disabled bool
		err      error
	}{
		"correct chain conifg": {
			func() *InitOpts {
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
				return &InitOpts{
					testChainConfig,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			false,
			nil,
		},
		"chain conifg is nil": {
			func() *InitOpts {
				return nil
			},
			true,
			ErrInitUnexpectedNil,
		},
		"node key is nil": {
			func() *InitOpts {
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
				return &InitOpts{
					testChainConfig,
					nil,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			true,
			ErrInitUnexpectedNil,
		},
		"chain is nil": {
			func() *InitOpts {
				return &InitOpts{
					testChainConfig,
					key,
					nil,
					&testTxPool{},
				}
			},
			true,
			ErrInitUnexpectedNil,
		},
		"tx pool is nil": {
			func() *InitOpts {
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
				return &InitOpts{
					testChainConfig,
					key,
					backend.BlockChain(),
					nil,
				}
			},
			true,
			ErrInitUnexpectedNil,
		},
		"gasless config is nil": {
			func() *InitOpts {
				cfg := testChainConfig.Copy()
				cfg.Gasless = nil
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, cfg)
				return &InitOpts{
					cfg,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			true,
			ErrInitUnexpectedNil,
		},
		"gasless is disabled": {
			func() *InitOpts {
				cfg := testChainConfig.Copy()
				cfg.Gasless.IsDisabled = true
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, cfg)
				return &InitOpts{
					cfg,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			true,
			nil,
		},
	}

	for _, tc := range tcs {
		g := NewGaslessModule()
		disabled, err := g.Init(tc.getOpts())
		require.Equal(t, tc.disabled, disabled)
		require.ErrorIs(t, tc.err, err)
	}
}

func TestInitAllowedTokens(t *testing.T) {
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	key, _ := crypto.GenerateKey()

	tcs := map[string]struct {
		getOpts       func() *InitOpts
		allowedTokens []common.Address
	}{
		"allowed tokens are nil": {
			func() *InitOpts {
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)
				return &InitOpts{
					testChainConfig,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2, dummyTokenAddress3},
		},
		"allowed tokens are empty slice": {
			func() *InitOpts {
				cfg := testChainConfig.Copy()
				cfg.Gasless.AllowedTokens = []common.Address{}
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, cfg)
				return &InitOpts{
					cfg,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			[]common.Address{},
		},
		"allowed tokens have existing addresses": {
			func() *InitOpts {
				cfg := testChainConfig.Copy()
				cfg.Gasless.AllowedTokens = []common.Address{dummyTokenAddress1, dummyTokenAddress2}
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, cfg)
				return &InitOpts{
					cfg,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2},
		},
		"allowed tokens have a non-existing address": {
			func() *InitOpts {
				cfg := testChainConfig.Copy()
				cfg.Gasless.AllowedTokens = []common.Address{common.HexToAddress("0xffff")}
				backend := backends.NewSimulatedBackendWithDatabase(db, alloc, cfg)
				return &InitOpts{
					cfg,
					key,
					backend.BlockChain(),
					&testTxPool{},
				}
			},
			[]common.Address{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			g := NewGaslessModule()
			disabled, err := g.Init(tc.getOpts())
			require.NoError(t, err)
			require.False(t, disabled)

			require.Equal(t, dummyGSRAddress.Bytes(), g.swapRouter.Bytes())
			require.Equal(t, len(tc.allowedTokens), len(g.allowedTokens))
			for _, addr := range tc.allowedTokens {
				_, ok := g.allowedTokens[addr]
				require.True(t, ok)
			}
		})
	}
}
