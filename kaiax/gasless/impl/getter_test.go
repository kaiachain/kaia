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
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestIsApproveTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	testcases := []struct {
		tx *types.Transaction
		ok bool
	}{
		{ // Legacy TestToken.approve(SwapRouter, 1000000)
			types.NewTransaction(0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(0),
				hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240")),
			true,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		StateDB:     statedb,
	})
	for _, tc := range testcases {
		ok := g.IsApproveTx(tc.tx)
		require.Equal(t, tc.ok, ok)
	}
}

func TestIsSwapTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	testcases := []struct {
		tx *types.Transaction
		ok bool
	}{
		{ // Legacy TestRouter.swapForGas(Token, 10, 100, 2021000)
			types.NewTransaction(0, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1),
				hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000001ed688")),
			true,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		StateDB:     statedb,
	})
	for _, tc := range testcases {
		ok := g.IsSwapTx(tc.tx)
		require.Equal(t, tc.ok, ok)
	}
}

func TestIsExecutable(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	testcases := []struct {
		approve *types.Transaction
		swap    *types.Transaction
		ok      bool
	}{
		{
			types.NewTransaction(0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1),
				hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240")),
			types.NewTransaction(1, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1),
				hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000001ed688")),
			true,
		},
		{
			nil,
			types.NewTransaction(0, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1),
				hexutil.MustDecode("0x43bab9f7000000000000000000000000000000000000000000000000000000000000abcd000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000f9448")),
			true,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
		StateDB:     statedb,
	})
	for _, tc := range testcases {
		ok := g.IsExecutable(tc.approve, tc.swap)
		require.Equal(t, tc.ok, ok)
	}
}
