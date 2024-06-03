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
package system

import (
	"math/big"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/contracts/testing/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestRebalanceTreasuryKIP103(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	rebalanceAddress := common.HexToAddress("0x1030")
	senderKey, _ := crypto.GenerateKey()
	sender := bind.NewKeyedTransactor(senderKey)
	rebalanceTreasury(t,
		sender,
		&params.ChainConfig{
			Kip103CompatibleBlock: big.NewInt(1),
			Kip103ContractAddress: rebalanceAddress,
		},
		rebalanceAddress,
		Kip103MockCode,
	)
}

func TestRebalanceTreasuryKIP160(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	rebalanceAddress := common.HexToAddress("0x1030")
	senderKey, _ := crypto.GenerateKey()
	sender := bind.NewKeyedTransactor(senderKey)
	rebalanceTreasury(t,
		sender,
		&params.ChainConfig{
			Kip160CompatibleBlock: big.NewInt(1),
			Kip160ContractAddress: rebalanceAddress,
		},
		rebalanceAddress,
		Kip160MockCode,
	)
}

func rebalanceTreasury(t *testing.T, sender *bind.TransactOpts, config *params.ChainConfig, rebalanceAddress common.Address, rebalanceCode []byte) {
	var (
		senderAddr = sender.From

		zeroeds = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xaa00"), big.NewInt(4_000_000)},
			{common.HexToAddress("0xaa11"), big.NewInt(2_000_000)},
			{common.HexToAddress("0xaa22"), big.NewInt(1_000_000)},
		}

		allocateds = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xbb00"), big.NewInt(8_012_345)},
			{common.HexToAddress("0xbb11"), big.NewInt(0)},
		}

		alloc = blockchain.GenesisAlloc{
			senderAddr:         {Balance: big.NewInt(params.KAIA)},
			rebalanceAddress:   {Code: rebalanceCode, Balance: common.Big0},
			zeroeds[0].addr:    {Balance: zeroeds[0].balance},
			zeroeds[1].addr:    {Balance: zeroeds[1].balance},
			zeroeds[2].addr:    {Balance: zeroeds[2].balance},
			allocateds[0].addr: {Balance: allocateds[0].balance},
			allocateds[1].addr: {Balance: allocateds[1].balance},
		}
	)

	testCases := []struct {
		rebalanceBlockNumber *big.Int
		status               uint8
		allocatedAmounts     []*big.Int

		expectedErr             error
		expectZeroedsAmounts    []*big.Int
		expectAllocatedsAmounts []*big.Int
		expectNonce             uint64
		expectBurnt             *big.Int
		expectSuccess           bool
		expectMemo              string
	}{
		{
			// normal case
			big.NewInt(1),
			EnumRebalanceStatus_Approved,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			nil,
			[]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},
			10,
			big.NewInt(8_012_345),
			true,
			"{\"retired\":{\"0x000000000000000000000000000000000000aa00\":4000000,\"0x000000000000000000000000000000000000aa11\":2000000,\"0x000000000000000000000000000000000000aa22\":1000000},\"newbie\":{\"0x000000000000000000000000000000000000bb00\":2000000,\"0x000000000000000000000000000000000000bb11\":5000000},\"burnt\":8012345,\"success\":true}",
		},
		{
			// failed case - rebalancing aborted due to wrong rebalanceBlockNumber
			big.NewInt(2),
			EnumRebalanceStatus_Approved,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			ErrRebalanceIncorrectBlock,
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
			8,
			big.NewInt(0),
			false,
			"{\"retired\":{\"0x000000000000000000000000000000000000aa00\":4000000,\"0x000000000000000000000000000000000000aa11\":2000000,\"0x000000000000000000000000000000000000aa22\":1000000},\"newbie\":{\"0x000000000000000000000000000000000000bb00\":2000000,\"0x000000000000000000000000000000000000bb11\":5000000},\"burnt\":0,\"success\":false}",
		},
		{
			// failed case - rebalancing aborted due to wrong state
			big.NewInt(1),
			EnumRebalanceStatus_Registered,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			ErrRebalanceBadStatus,
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
			9,
			big.NewInt(0),
			false,
			"{\"retired\":{\"0x000000000000000000000000000000000000aa00\":4000000,\"0x000000000000000000000000000000000000aa11\":2000000,\"0x000000000000000000000000000000000000aa22\":1000000},\"newbie\":{\"0x000000000000000000000000000000000000bb00\":2000000,\"0x000000000000000000000000000000000000bb11\":5000000},\"burnt\":0,\"success\":false}",
		},
		{
			// failed case - rebalancing aborted doe to bigger allocation than zeroeds
			big.NewInt(1),
			EnumRebalanceStatus_Registered,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000 + 1)},

			ErrRebalanceBadStatus,
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
			9,
			big.NewInt(0),
			false,
			"{\"retired\":{\"0x000000000000000000000000000000000000aa00\":4000000,\"0x000000000000000000000000000000000000aa11\":2000000,\"0x000000000000000000000000000000000000aa22\":1000000},\"newbie\":{\"0x000000000000000000000000000000000000bb00\":2000000,\"0x000000000000000000000000000000000000bb11\":5000001},\"burnt\":0,\"success\":false}",
		},
	}

	for _, tc := range testCases {
		var (
			db                          = database.NewMemoryDBManager()
			backend                     = backends.NewSimulatedBackendWithDatabase(db, alloc, config)
			chain                       = backend.BlockChain()
			zeroedAddrs, allocatedAddrs []common.Address
		)

		// Deploy TreasuryRebalanceMock contract at block 0 and transit to block 1
		for _, entry := range zeroeds {
			zeroedAddrs = append(zeroedAddrs, entry.addr)
		}
		for _, entry := range allocateds {
			allocatedAddrs = append(allocatedAddrs, entry.addr)
		}

		if chain.Config().Kip160CompatibleBlock != nil {
			contract, _ := system_contracts.NewTreasuryRebalanceMockV2Transactor(rebalanceAddress, backend)
			_, err := contract.TestSetAll(sender, zeroedAddrs, allocatedAddrs, tc.allocatedAmounts, tc.rebalanceBlockNumber, tc.status)
			assert.Nil(t, err)
		} else {
			contract, _ := system_contracts.NewTreasuryRebalanceMockTransactor(rebalanceAddress, backend)
			_, err := contract.TestSetAll(sender, zeroedAddrs, allocatedAddrs, tc.allocatedAmounts, tc.rebalanceBlockNumber, tc.status)
			assert.Nil(t, err)
		}

		backend.Commit()
		assert.Equal(t, uint64(1), chain.CurrentBlock().NumberU64())

		// Get state and check post-transition state
		state, err := chain.State()
		assert.Nil(t, err)

		res, err := RebalanceTreasury(state, chain, chain.CurrentHeader())
		assert.Equal(t, tc.expectedErr, err)

		for i, addr := range zeroedAddrs {
			assert.Equal(t, tc.expectZeroedsAmounts[i], state.GetBalance(addr))
		}
		for i, addr := range allocatedAddrs {
			assert.Equal(t, tc.expectAllocatedsAmounts[i], state.GetBalance(addr))
		}

		if chain.Config().Kip160CompatibleBlock != nil {
			assert.Equal(t, uint64(0x0), state.GetNonce(common.HexToAddress("0x0")))
		} else {
			assert.Equal(t, tc.expectNonce, state.GetNonce(common.HexToAddress("0x0")))
		}
		assert.Equal(t, tc.expectBurnt, res.Burnt)
		assert.Equal(t, tc.expectSuccess, res.Success)

		isKip160, isKip103 := chain.Config().Kip160CompatibleBlock != nil, chain.Config().Kip103CompatibleBlock != nil
		if isKip160 {
			tc.expectMemo = strings.Replace(tc.expectMemo, "retired", "zeroed", -1)
			tc.expectMemo = strings.Replace(tc.expectMemo, "newbie", "allocated", -1)
		}
		assert.Equal(t, tc.expectMemo, string(res.memo(isKip103)))
	}
}
