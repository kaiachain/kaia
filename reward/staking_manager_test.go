// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package reward

import (
	"encoding/json"
	"math/big"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

type testGovernance struct {
	p *params.GovParamSet
}

func newTestGovernance(intMap map[int]interface{}) *testGovernance {
	p, _ := params.NewGovParamSetIntMap(intMap)
	return &testGovernance{p}
}

func newDefaultTestGovernance() *testGovernance {
	return newTestGovernance(map[int]interface{}{
		params.Epoch:               604800,
		params.Policy:              params.WeightedRandom,
		params.UnitPrice:           25000000000,
		params.MintingAmount:       "9600000000000000000",
		params.Ratio:               "34/54/12",
		params.UseGiniCoeff:        true,
		params.DeferredTxFee:       true,
		params.MinimumStake:        "5000000",
		params.StakeUpdateInterval: 86400,
	})
}

type stakingManagerTestCase struct {
	blockNum    uint64       // Requested num in GetStakingInfo(num)
	stakingNum  uint64       // Corresponding staking info block number
	stakingInfo *StakingInfo // Expected GetStakingInfo() output
}

// Note that Golang will correctly initialize these globals according to dependency.
// https://go.dev/ref/spec#Order_of_evaluation

// Note the testdata must not exceed maxStakingCache because otherwise cache test will fail.
var stakingManagerTestData = []*StakingInfo{
	stakingInfoTestCases[0].stakingInfo,
	stakingInfoTestCases[1].stakingInfo,
	stakingInfoTestCases[2].stakingInfo,
	stakingInfoTestCases[3].stakingInfo,
}
var stakingManagerTestCases = generateStakingManagerTestCases()

func generateStakingManagerTestCases() []stakingManagerTestCase {
	s := stakingManagerTestData

	return []stakingManagerTestCase{
		{1, 0, s[0]},
		{100, 0, s[0]},
		{86400, 0, s[0]},
		{86401, 0, s[0]},
		{172800, 0, s[0]},
		{172801, 86400, s[1]},
		{200000, 86400, s[1]},
		{259200, 86400, s[1]},
		{259201, 172800, s[2]},
		{300000, 172800, s[2]},
		{345600, 172800, s[2]},
		{345601, 259200, s[3]},
		{400000, 259200, s[3]},
	}
}

func newStakingManagerForTest(t *testing.T) {
	// test if nil
	assert.Nil(t, GetStakingManager())
	assert.Nil(t, GetStakingInfo(123))

	st, err := updateStakingInfo(456)
	assert.Nil(t, st)
	assert.EqualError(t, err, ErrStakingManagerNotSet.Error())

	assert.EqualError(t, checkStakingInfoStored(789), ErrStakingManagerNotSet.Error())

	// test if get same
	stNew := NewStakingManager(&blockchain.BlockChain{}, newDefaultTestGovernance(), nil)
	stGet := GetStakingManager()
	assert.NotNil(t, stNew)
	assert.Equal(t, stGet, stNew)
}

func resetStakingManagerForTest(t *testing.T) {
	sm := GetStakingManager()
	if sm == nil {
		newStakingManagerForTest(t)
		sm = GetStakingManager()
	}

	cache, _ := lru.NewARC(128)
	sm.stakingInfoCache = cache
	sm.stakingInfoDB = database.NewMemoryDBManager()
}

func TestStakingManager_NewStakingManager(t *testing.T) {
	newStakingManagerForTest(t)
}

// Check that appropriate StakingInfo is returned given various blockNum argument.
func checkGetStakingInfo(t *testing.T) {
	for _, testcase := range stakingManagerTestCases {
		expcectedInfo := testcase.stakingInfo
		actualInfo := GetStakingInfo(testcase.blockNum)

		assert.Equal(t, testcase.stakingNum, actualInfo.BlockNum)
		assert.Equal(t, expcectedInfo, actualInfo)
	}
}

// Check that StakingInfo are loaded from cache
func TestStakingManager_GetFromCache(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	resetStakingManagerForTest(t)

	for _, testdata := range stakingManagerTestData {
		GetStakingManager().stakingInfoCache.Add(testdata.BlockNum, testdata)
	}

	checkGetStakingInfo(t)
}

// Check that StakingInfo are loaded from database
func TestStakingManager_GetFromDB(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	resetStakingManagerForTest(t)

	for _, testdata := range stakingManagerTestData {
		AddStakingInfoToDB(testdata)
	}

	checkGetStakingInfo(t)
}

// Even if Gini was -1 in the cache, GetStakingInfo returns valid Gini
func TestStakingManager_FillGiniFromCache(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	resetStakingManagerForTest(t)

	for _, testdata := range stakingManagerTestData {
		// Insert a modified copy of testdata to cache
		copydata := &StakingInfo{}
		json.Unmarshal([]byte(testdata.String()), copydata)
		copydata.Gini = -1 // Suppose Gini was -1 in the cache
		GetStakingManager().stakingInfoCache.Add(copydata.BlockNum, copydata)
	}

	checkGetStakingInfo(t)
}

// Even if Gini was -1 in the DB, GetStakingInfo returns valid Gini
func TestStakingManager_FillGiniFromDB(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	resetStakingManagerForTest(t)

	for _, testdata := range stakingManagerTestData {
		// Insert a modified copy of testdata to cache
		copydata := &StakingInfo{}
		json.Unmarshal([]byte(testdata.String()), copydata)
		copydata.Gini = -1 // Suppose Gini was -1 in the cache
		AddStakingInfoToDB(copydata)
	}

	checkGetStakingInfo(t)
}

var expectedAddress = []common.Address{
	common.HexToAddress("0x0000000000000000000000000000000000000F00"), // CN 1's node id
	common.HexToAddress("0x0000000000000000000000000000000000000F01"), // CN 1's staking address
	common.HexToAddress("0x0000000000000000000000000000000000000F02"), // CN 1's reward address
	common.HexToAddress("0x0000000000000000000000000000000000000F03"), // CN 2's node id
	common.HexToAddress("0x0000000000000000000000000000000000000F04"), // CN 2's staking address
	common.HexToAddress("0x0000000000000000000000000000000000000F05"), // CN 2's reward address
	common.HexToAddress("0x0000000000000000000000000000000000000F06"), // KIF (POC, KFF)
	common.HexToAddress("0x0000000000000000000000000000000000000F07"), // KEF (KIR, KCF)
}

func TestStakingManager_GetFromAddressBook(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		alloc = blockchain.GenesisAlloc{
			system.AddressBookAddr: {
				Code:    system.AddressBookMockTwoCNCode,
				Balance: big.NewInt(0),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000F01"): {
				Balance: big.NewInt(0).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA)),
			},
			common.HexToAddress("0x0000000000000000000000000000000000000F04"): {
				Balance: big.NewInt(0).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA)),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)
	)
	defer func() {
		backend.Close()
	}()

	SetTestStakingManagerWithChain(backend.BlockChain(), newDefaultTestGovernance(), nil)

	stakingInfo := GetStakingInfo(0)

	actualAddress := []common.Address{
		stakingInfo.CouncilNodeAddrs[0],
		stakingInfo.CouncilStakingAddrs[0],
		stakingInfo.CouncilRewardAddrs[0],
		stakingInfo.CouncilNodeAddrs[1],
		stakingInfo.CouncilStakingAddrs[1],
		stakingInfo.CouncilRewardAddrs[1],
		stakingInfo.KIFAddr,
		stakingInfo.KEFAddr,
	}
	for i := 0; i < 8; i++ {
		assert.Equal(t, expectedAddress[i], actualAddress[i])
	}
	assert.Equal(t, uint64(5_000_000), stakingInfo.CouncilStakingAmounts[0])
	assert.Equal(t, uint64(5_000_000), stakingInfo.CouncilStakingAmounts[1])
}

// Check that StakingInfo are loaded from multicall contract correctly
func TestStakingManager_GetFromMultiCall(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		alloc = blockchain.GenesisAlloc{
			system.AddressBookAddr: {
				Code:    system.AddressBookMockTwoCNCode,
				Balance: big.NewInt(0),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)
	)
	originCode := system.MultiCallCode
	// Temporary code injection
	system.MultiCallCode = system.MultiCallMockCode
	defer func() {
		system.MultiCallCode = originCode
		backend.Close()
	}()

	backend.BlockChain().Config().KaiaCompatibleBlock = big.NewInt(0)
	SetTestStakingManagerWithChain(backend.BlockChain(), newDefaultTestGovernance(), nil)

	stakingInfo := GetStakingInfo(0)

	actualAddress := []common.Address{
		stakingInfo.CouncilNodeAddrs[0],
		stakingInfo.CouncilStakingAddrs[0],
		stakingInfo.CouncilRewardAddrs[0],
		stakingInfo.CouncilNodeAddrs[1],
		stakingInfo.CouncilStakingAddrs[1],
		stakingInfo.CouncilRewardAddrs[1],
		stakingInfo.KIFAddr,
		stakingInfo.KEFAddr,
	}
	for i := 0; i < 8; i++ {
		assert.Equal(t, expectedAddress[i], actualAddress[i])
	}
	assert.Equal(t, uint64(5_000_000), stakingInfo.CouncilStakingAmounts[0])
	assert.Equal(t, uint64(20_000_000), stakingInfo.CouncilStakingAmounts[1])
}
