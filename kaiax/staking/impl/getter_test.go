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

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func testPragueForkChainConfig(forkNum *big.Int) *params.ChainConfig {
	var config *params.ChainConfig

	config = &params.ChainConfig{
		ChainID: common.Big1,
		Governance: &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				UseGiniCoeff:          true,
				StakingUpdateInterval: 86400,
			},
			KIP71: params.GetDefaultKIP71Config(),
		},
	}
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = nil
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.PragueCompatibleBlock = forkNum

	return config
}

func copyAlloc(alloc blockchain.GenesisAlloc) blockchain.GenesisAlloc {
	copiedAlloc := blockchain.GenesisAlloc{}
	for addr, acc := range alloc {
		copiedAlloc[addr] = acc
	}
	return copiedAlloc
}

func TestGetStakingInfo_Uncached(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db    = database.NewMemoryDBManager()
		alloc = blockchain.GenesisAlloc{
			system.AddressBookAddr: {
				Code:    system.AddressBookMockTwoCNCode,
				Balance: big.NewInt(0),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000F01"): { // staking1
				Balance: new(big.Int).Mul(big.NewInt(42_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000f04"): { // staking2
				Balance: new(big.Int).Mul(big.NewInt(99_000_000), big.NewInt(params.KAIA)),
			},
		}
		config = testPragueForkChainConfig(nil)

		// Addresses are already stored in AddressBookMock.sol:AddressBookMockTwoCN
		// The balances are given at the GenesisAlloc above
		expected = &staking.StakingInfo{
			SourceBlockNum: 0,
			NodeIds: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000F00"),
				common.HexToAddress("0x0000000000000000000000000000000000000F03"),
			},
			StakingContracts: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000F01"),
				common.HexToAddress("0x0000000000000000000000000000000000000f04"),
			},
			RewardAddrs: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000f02"),
				common.HexToAddress("0x0000000000000000000000000000000000000f05"),
			},
			KIFAddr:        common.HexToAddress("0x0000000000000000000000000000000000000F06"),
			KEFAddr:        common.HexToAddress("0x0000000000000000000000000000000000000f07"),
			StakingAmounts: []uint64{42_000_000, 99_000_000},
			CLStakingInfos: nil,
		}
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	// Test GetStakingInfo()
	mStaking := NewStakingModule()
	mStaking.Init(&InitOpts{
		ChainKv:     db.GetMiscDB(),
		ChainConfig: config,
		Chain:       backend.BlockChain(),
	})
	si, err := mStaking.GetStakingInfo(0)
	assert.NoError(t, err)
	assert.Equal(t, expected, si)
}

func TestGetStakingInfo_Prague_Uncached(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		alloc = blockchain.GenesisAlloc{
			system.AddressBookAddr: {
				Code:    system.AddressBookMockTwoCNCode,
				Balance: big.NewInt(0),
			},
			system.RegistryAddr: {
				Code:    system.RegistryMockForCLCode,
				Balance: big.NewInt(0),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000F01"): { // staking1
				Balance: new(big.Int).Mul(big.NewInt(42_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000f04"): { // staking2
				Balance: new(big.Int).Mul(big.NewInt(99_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000e00"): { // CLPool1
				Balance: new(big.Int).Mul(big.NewInt(20_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000e01"): { // CLPool2
				Balance: new(big.Int).Mul(big.NewInt(23_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000e02"): { // CLPool3
				Balance: new(big.Int).Mul(big.NewInt(30_000_000), big.NewInt(params.KAIA)),
			},
		}
		expected = &staking.StakingInfo{
			SourceBlockNum: 0,
			NodeIds: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000F00"),
				common.HexToAddress("0x0000000000000000000000000000000000000F03"),
			},
			StakingContracts: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000F01"),
				common.HexToAddress("0x0000000000000000000000000000000000000f04"),
			},
			RewardAddrs: []common.Address{
				common.HexToAddress("0x0000000000000000000000000000000000000f02"),
				common.HexToAddress("0x0000000000000000000000000000000000000f05"),
			},
			KIFAddr:        common.HexToAddress("0x0000000000000000000000000000000000000F06"),
			KEFAddr:        common.HexToAddress("0x0000000000000000000000000000000000000f07"),
			StakingAmounts: []uint64{42_000_000, 99_000_000},
			CLStakingInfos: staking.CLStakingInfos{
				{
					CLNodeId:        common.HexToAddress("0x0000000000000000000000000000000000000F00"),
					CLPoolAddr:      common.HexToAddress("0x0000000000000000000000000000000000000e00"),
					CLStakingAmount: 20_000_000,
				},
				{
					CLNodeId:        common.HexToAddress("0x0000000000000000000000000000000000000F03"),
					CLPoolAddr:      common.HexToAddress("0x0000000000000000000000000000000000000e01"),
					CLStakingAmount: 23_000_000,
				},
				{
					CLNodeId:        common.HexToAddress("0x0000000000000000000000000000000000000F06"),
					CLPoolAddr:      common.HexToAddress("0x0000000000000000000000000000000000000e02"),
					CLStakingAmount: 30_000_000,
				},
			},
		}
	)

	testGetStakingInfo_CL_NoCLRegistry(t, alloc)
	testGetStakingInfo_CL(t, alloc, expected)
}

func testGetStakingInfo_CL_NoCLRegistry(t *testing.T, alloc blockchain.GenesisAlloc) {
	db := database.NewMemoryDBManager()

	config := testPragueForkChainConfig(big.NewInt(0))
	alloc = copyAlloc(alloc)
	alloc[system.RegistryAddr] = blockchain.GenesisAccount{
		Code:    system.RegistryMockZero,
		Balance: big.NewInt(0),
	}
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	// Test GetStakingInfo()
	mStaking := NewStakingModule()
	mStaking.Init(&InitOpts{
		ChainKv:     db.GetMiscDB(),
		ChainConfig: config,
		Chain:       backend.BlockChain(),
	})
	si, err := mStaking.GetStakingInfo(0)
	assert.NoError(t, err)
	assert.EqualValues(t, staking.CLStakingInfos(nil), si.CLStakingInfos)
}

func testGetStakingInfo_CL(t *testing.T, alloc blockchain.GenesisAlloc, expected *staking.StakingInfo) {
	db := database.NewMemoryDBManager()

	config := testPragueForkChainConfig(big.NewInt(0))
	alloc = copyAlloc(alloc)
	alloc[system.CLRegistryMockThreeCLAddr] = blockchain.GenesisAccount{
		Code:    system.CLRegistryMockThreeCLCode,
		Balance: big.NewInt(0),
	}
	alloc[system.WrappedKaiaMockAddr] = blockchain.GenesisAccount{
		Code:    system.WrappedKaiaMockCode,
		Balance: big.NewInt(0),
	}
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	// Test GetStakingInfo()
	mStaking := NewStakingModule()
	mStaking.Init(&InitOpts{
		ChainKv:     db.GetMiscDB(),
		ChainConfig: config,
		Chain:       backend.BlockChain(),
	})
	si, err := mStaking.GetStakingInfo(0)
	assert.NoError(t, err)
	assert.Equal(t, expected, si)
}

func TestSourceBlockNum(t *testing.T) {
	testcases := []struct {
		num      uint64
		isKaia   bool
		interval uint64
		expected uint64
	}{
		// Before Kaia
		{0, false, 1000, 0},
		{1, false, 1000, 0},
		{1000, false, 1000, 0},
		{1001, false, 1000, 0},
		{1999, false, 1000, 0},
		{2000, false, 1000, 0},
		{2001, false, 1000, 1000},
		{2999, false, 1000, 1000},
		{3000, false, 1000, 1000},
		{3001, false, 1000, 2000},

		// After Kaia
		{0, true, 1000, 0},
		{1, true, 1000, 0},
		{1000, true, 1000, 999},
		{1001, true, 1000, 1000},
		{1999, true, 1000, 1998},
	}

	for i, tc := range testcases {
		actual := sourceBlockNum(tc.num, tc.isKaia, tc.interval)
		assert.Equal(t, tc.expected, actual, i)
	}
}
