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
	"github.com/kaiachain/kaia/kaiax/gasless"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitSuccessAndFailure(t *testing.T) {
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	key, _ := crypto.GenerateKey()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)

	tcs := map[string]struct {
		opts     *InitOpts
		disabled bool
		err      error
	}{
		"correct chain conifg": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			false,
			nil,
		},
		"init option is nil": {
			nil,
			true,
			ErrInitUnexpectedNil,
		},
		"cn config is nil": {
			&InitOpts{
				testChainConfig,
				nil,
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			true,
			ErrInitUnexpectedNil,
		},
		"node key is nil": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				nil,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			true,
			ErrInitUnexpectedNil,
		},
		"chain is nil": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				key,
				nil,
				common.ENDPOINTNODE,
			},
			true,
			ErrInitUnexpectedNil,
		},
		"gasless is disabled": {
			&InitOpts{
				testChainConfig,
				&gasless.GaslessConfig{
					AllowedTokens: nil,
					Disable:       true,
				},
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			true,
			nil,
		},
		"cn has insufficient balance": {
			&InitOpts{
				testChainConfig,
				&gasless.GaslessConfig{
					AllowedTokens: nil,
					Disable:       true,
				},
				key,
				backend.BlockChain(),
				common.CONSENSUSNODE,
			},
			true,
			nil,
		},
		"system contract is not prepared": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				key,
				backends.NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), nil, testChainConfig).BlockChain(),
				common.ENDPOINTNODE,
			},
			false,
			nil,
		},
	}

	for tcName, tc := range tcs {
		g := NewGaslessModule()
		err := g.Init(tc.opts)
		require.ErrorIs(t, tc.err, err, tcName)
		if err == nil {
			disabled := g.IsDisabled()
			assert.Equal(t, tc.disabled, disabled, tcName)
		}
	}
}

func TestInitGSRAndAllowedTokens(t *testing.T) {
	db := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	key, _ := crypto.GenerateKey()
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, testChainConfig)

	tcs := map[string]struct {
		opts          *InitOpts
		swapRouter    common.Address
		allowedTokens []common.Address
	}{
		"allowed tokens are nil": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			dummyGSRAddress,
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2, dummyTokenAddress3},
		},
		"allowed tokens are empty slice": {
			&InitOpts{
				testChainConfig,
				&gasless.GaslessConfig{
					AllowedTokens: []common.Address{},
					Disable:       false,
				},
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			dummyGSRAddress,
			[]common.Address{},
		},
		"allowed tokens have existing addresses": {
			&InitOpts{
				testChainConfig,
				&gasless.GaslessConfig{
					AllowedTokens: []common.Address{dummyTokenAddress1, dummyTokenAddress2},
					Disable:       false,
				},
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			dummyGSRAddress,
			[]common.Address{dummyTokenAddress1, dummyTokenAddress2},
		},
		"allowed tokens have a non-existing address": {
			&InitOpts{
				testChainConfig,
				&gasless.GaslessConfig{
					AllowedTokens: []common.Address{common.HexToAddress("0xffff")},
					Disable:       false,
				},
				key,
				backend.BlockChain(),
				common.ENDPOINTNODE,
			},
			dummyGSRAddress,
			[]common.Address{},
		},
		"system contract is not prepared": {
			&InitOpts{
				testChainConfig,
				testGaslessConfig,
				key,
				backends.NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), nil, testChainConfig).BlockChain(),
				common.ENDPOINTNODE,
			},
			common.Address{},
			nil,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			g := NewGaslessModule()
			err := g.Init(tc.opts)
			require.NoError(t, err)
			require.Equal(t, tc.swapRouter, g.swapRouter)
			require.Equal(t, len(tc.allowedTokens), len(g.allowedTokens))
			for _, addr := range tc.allowedTokens {
				_, ok := g.allowedTokens[addr]
				require.True(t, ok)
			}
		})
	}
}
