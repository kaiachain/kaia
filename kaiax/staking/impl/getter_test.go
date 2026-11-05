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
	"errors"
	"fmt"
	"maps"
	"math/big"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
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
	config.RandaoCompatibleBlock = big.NewInt(0)
	config.RandaoRegistry = &params.RegistryConfig{
		Records: map[string]common.Address{
			"KIP113": system.Kip113LogicAddrMock,
		},
		Owner: common.HexToAddress("0x0000000000000000000000000000000000000000"),
	}
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.PragueCompatibleBlock = forkNum

	return config
}

func copyAlloc(alloc blockchain.GenesisAlloc) blockchain.GenesisAlloc {
	copiedAlloc := blockchain.GenesisAlloc{}
	maps.Copy(copiedAlloc, alloc)
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
	mStaking.fillComputedFields(expected) // the cached StakingInfo has them filled
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
	mStaking.fillComputedFields(expected) // the cached StakingInfo has them filled
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

// newPermissionlessTestStakingModule returns a StakingModule post-permissionless-fork from genesis,
// whose MultiCall returns the 6 profiles of MultiCallContractMockPermissionless:
//
//	[0] node=0xF00  staking=0xF01  reward=0xF02  ValActive    5M
//	[1] node=0xF03  staking=0xF04  reward=0xF05  ValPaused   10M
//	[2] node=0x1000 staking=0x1001 reward=0x1002 ValReady     8M
//	[3] node=0x2000 staking=0x2001 reward=0x2002 ValExiting   7M
//	[4] node=0x3000 staking=0x3001 reward=0x3002 CandReady    6M
//	[5] node=0x4000 staking=0x4001 reward=0x4002 ValActive    9M
func newPermissionlessTestStakingModule(t *testing.T, mValset valset.ValsetModule) *StakingModule {
	originCode := system.MultiCallCode
	system.MultiCallCode = system.MultiCallPermlessMockCode
	t.Cleanup(func() { system.MultiCallCode = originCode })

	db := database.NewMemoryDBManager()
	config := testPragueForkChainConfig(big.NewInt(0))
	config.PermissionlessCompatibleBlock = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		system.AddressBookAddr: {Code: []byte{0x01}, Balance: big.NewInt(0)}, // non-nil to pass bail-out check
	}
	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	mStaking := NewStakingModule()
	assert.NoError(t, mStaking.Init(&InitOpts{
		ChainKv:      db.GetMiscDB(),
		ChainConfig:  config,
		Chain:        backend.BlockChain(),
		ValsetModule: mValset,
	}))
	return mStaking
}

// Post-permissionless-fork, GetStakingInfo keeps only the nodes reward-eligible at the block (KIP-286).
func TestGetStakingInfo_Permissionless_RewardEligible(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		n0, n2, n5 = common.HexToAddress("0xF00"), common.HexToAddress("0x1000"), common.HexToAddress("0x4000")
		ctrl       = gomock.NewController(t)
		mValset    = valset_mock.NewMockValsetModule(ctrl)
	)
	mValset.EXPECT().GetRewardEligible(uint64(1)).Return([]common.Address{n5, n2, n0}, nil).Times(1)
	mStaking := newPermissionlessTestStakingModule(t, mValset)

	si, err := mStaking.GetStakingInfo(1)
	assert.NoError(t, err)

	again, err := mStaking.GetStakingInfo(1)
	assert.NoError(t, err)
	assert.Same(t, si, again)
	assert.Equal(t, []common.Address{n0, n2, n5}, si.NodeIds) // in the AddressBookV2 order
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xF01"), common.HexToAddress("0x1001"), common.HexToAddress("0x4001"),
	}, si.StakingContracts)
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xF02"), common.HexToAddress("0x1002"), common.HexToAddress("0x4002"),
	}, si.RewardAddrs)
	assert.Equal(t, []uint64{5_000_000, 8_000_000, 9_000_000}, si.StakingAmounts)
	assert.Equal(t, common.HexToAddress("0x0a01"), si.KEFAddr)

	entry, err := mStaking.getStakingInfoEntry(0)
	assert.NoError(t, err)
	assert.Len(t, entry.si.NodeIds, 6)
}

func TestGetStakingInfo_Permissionless_Concurrent(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	mValset := valset_mock.NewMockValsetModule(gomock.NewController(t))
	mValset.EXPECT().GetRewardEligible(uint64(1)).Return([]common.Address{common.HexToAddress("0xF00")}, nil).AnyTimes()
	mStaking := newPermissionlessTestStakingModule(t, mValset)
	mStaking.minimumStake = big.NewInt(5_000_000)
	_, err := mStaking.getStakingInfoEntry(0)
	assert.NoError(t, err)

	var (
		results = make([]*staking.StakingInfo, 8)
		wg      sync.WaitGroup
	)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i], _ = mStaking.GetStakingInfo(1)
			results[i].ConsolidatedNodes()
			results[i].ToResponse(false, mStaking.minimumStake.Uint64()) // as kaia_getStakingInfo does
		}()
	}
	wg.Wait()
	for _, r := range results {
		assert.Same(t, results[0], r)
	}
}

// Blocks 0 and 1 share the staking info of block 0, but each keeps its own reward-eligible nodes
// regardless of which one is queried first.
func TestGetStakingInfo_Permissionless_Genesis(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		n0, n2, n5 = common.HexToAddress("0xF00"), common.HexToAddress("0x1000"), common.HexToAddress("0x4000")
		eligible0  = []common.Address{n0, n2, n5}
		eligible1  = []common.Address{n2, n5}
	)
	for _, order := range [][]uint64{{0, 1}, {1, 0}} {
		t.Run(fmt.Sprintf("%v", order), func(t *testing.T) {
			mValset := valset_mock.NewMockValsetModule(gomock.NewController(t))
			mValset.EXPECT().GetRewardEligible(uint64(0)).Return(eligible0, nil).AnyTimes()
			mValset.EXPECT().GetRewardEligible(uint64(1)).Return(eligible1, nil).Times(1)
			mStaking := newPermissionlessTestStakingModule(t, mValset)

			expected := map[uint64][]common.Address{0: eligible0, 1: eligible1}
			for _, num := range append(order, order...) {
				si, err := mStaking.GetStakingInfo(num)
				assert.NoError(t, err)
				assert.Equal(t, expected[num], si.NodeIds, "block %d", num)
			}
		})
	}
}

func TestGetStakingInfo_Permissionless_RewardEligibleError(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	t.Run("valset module not set", func(t *testing.T) {
		db := database.NewMemoryDBManager()
		config := testPragueForkChainConfig(big.NewInt(0))
		config.PermissionlessCompatibleBlock = big.NewInt(100)
		backend := backends.NewSimulatedBackendWithDatabase(db, blockchain.GenesisAlloc{}, config)
		err := NewStakingModule().Init(&InitOpts{
			ChainKv:     db.GetMiscDB(),
			ChainConfig: config,
			Chain:       backend.BlockChain(),
		})
		assert.ErrorIs(t, err, staking.ErrValsetModuleNotSet)
	})

	t.Run("reward-eligible nodes unavailable", func(t *testing.T) {
		errNoState := errors.New("no state")
		mValset := valset_mock.NewMockValsetModule(gomock.NewController(t))
		mValset.EXPECT().GetRewardEligible(uint64(1)).Return(nil, errNoState).Times(2)
		mStaking := newPermissionlessTestStakingModule(t, mValset)
		_, err := mStaking.GetStakingInfo(1)
		assert.ErrorIs(t, err, errNoState)
		_, err = mStaking.GetStakingInfo(1) // an error is not kept
		assert.ErrorIs(t, err, errNoState)

		_, err = mStaking.getStakingInfoEntry(0)
		assert.NoError(t, err)
	})
}

// Post-permissionless-fork, the unfiltered staking info must source effective stake from
// AddressBookV2 (MultiCallStakingInfoPermissionless) and keep every profile regardless of
// its state, in the AddressBookV2 order.
func TestGetStakingInfo_Permissionless_AllProfiles(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	mStaking := newPermissionlessTestStakingModule(t, valset_mock.NewMockValsetModule(gomock.NewController(t)))
	entry, err := mStaking.getStakingInfoEntry(0)
	assert.NoError(t, err)
	si := entry.si

	assert.Equal(t, []common.Address{
		common.HexToAddress("0xF00"), common.HexToAddress("0xF03"), common.HexToAddress("0x1000"),
		common.HexToAddress("0x2000"), common.HexToAddress("0x3000"), common.HexToAddress("0x4000"),
	}, si.NodeIds)
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xF01"), common.HexToAddress("0xF04"), common.HexToAddress("0x1001"),
		common.HexToAddress("0x2001"), common.HexToAddress("0x3001"), common.HexToAddress("0x4001"),
	}, si.StakingContracts)
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xF02"), common.HexToAddress("0xF05"), common.HexToAddress("0x1002"),
		common.HexToAddress("0x2002"), common.HexToAddress("0x3002"), common.HexToAddress("0x4002"),
	}, si.RewardAddrs)
	assert.Equal(t, []uint64{5_000_000, 10_000_000, 8_000_000, 7_000_000, 6_000_000, 9_000_000}, si.StakingAmounts)
	assert.Equal(t, common.HexToAddress("0x0a01"), si.KEFAddr)
	assert.Equal(t, common.HexToAddress("0x0a02"), si.KIFAddr)
	assert.Equal(t, common.HexToAddress("0x0a03"), si.KPFAddr)
}
