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
package governance

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/reward"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

type testBlockChain struct {
	num    uint64
	config *params.ChainConfig
}

func newTestBlockchain(config *params.ChainConfig) *testBlockChain {
	return &testBlockChain{
		config: config,
	}
}

func newTestGovernanceApi() *GovernanceAPI {
	config := params.MainnetChainConfig
	config.Governance.KIP71 = params.GetDefaultKIP71Config()
	govApi := NewGovernanceAPI(NewMixedEngine(config, database.NewMemoryDBManager()))
	govApi.governance.SetNodeAddress(common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"))
	bc := newTestBlockchain(config)
	govApi.governance.SetBlockchain(bc)
	return govApi
}

func TestUpperBoundBaseFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curLowerBoundBaseFee := govApi.governance.CurrentParams().LowerBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidUpperBoundBaseFee := curLowerBoundBaseFee - 100
	_, err := govApi.Vote("kip71.upperboundbasefee", invalidUpperBoundBaseFee)
	assert.Equal(t, err, errInvalidUpperBound)
}

func TestLowerBoundFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curUpperBoundBaseFee := govApi.governance.CurrentParams().UpperBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidLowerBoundBaseFee := curUpperBoundBaseFee + 100
	_, err := govApi.Vote("kip71.lowerboundbasefee", invalidLowerBoundBaseFee)
	assert.Equal(t, err, errInvalidLowerBound)
}

func TestGetRewards(t *testing.T) {
	type expected = map[int]uint64
	type strMap = map[string]interface{}
	type override struct {
		num    int
		strMap strMap
	}
	type testcase struct {
		length   int // total number of blocks to simulate
		override []override
		expected expected
	}

	var (
		mintAmount = uint64(1)
		koreBlock  = uint64(9)
		epoch      = 3
		latestNum  = rpc.BlockNumber(-1)
		proposer   = common.HexToAddress("0x0000000000000000000000000000000000000000")
		config     = getTestConfig()
	)

	testcases := []testcase{
		{
			12,
			[]override{
				{
					3,
					strMap{
						"reward.mintingamount": "2",
					},
				},
				{
					6,
					strMap{
						"reward.mintingamount": "3",
					},
				},
			},
			map[int]uint64{
				1:  1,
				2:  1,
				3:  1,
				4:  1,
				5:  1,
				6:  1,
				7:  2, // 2 is minted from now
				8:  2,
				9:  3, // 3 is minted from now
				10: 3,
				11: 3,
				12: 3,
				13: 3,
			},
		},
	}

	for _, tc := range testcases {
		config.Governance.Reward.MintingAmount = new(big.Int).SetUint64(mintAmount)
		config.Istanbul.Epoch = uint64(epoch)
		config.KoreCompatibleBlock = new(big.Int).SetUint64(koreBlock)

		bc := newTestBlockchain(config)

		dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})

		e := NewMixedEngine(config, dbm)
		e.SetBlockchain(bc)
		e.UpdateParams(bc.CurrentBlock().NumberU64())

		// write initial gov items and overrides to database
		pset, _ := params.NewGovParamSetChainConfig(config)
		gset := NewGovernanceSet()
		gset.Import(pset.StrMap())
		e.headerGov.WriteGovernance(0, NewGovernanceSet(), gset)
		for _, o := range tc.override {
			override := NewGovernanceSet()
			override.Import(o.strMap)
			e.headerGov.WriteGovernance(uint64(o.num), gset, override)
		}

		govKaiaApi := NewGovernanceKaiaAPI(e, bc)

		for num := 1; num <= tc.length; num++ {
			bc.SetBlockNum(uint64(num))

			rewardSpec, err := govKaiaApi.GetRewards(&latestNum)
			assert.Nil(t, err)

			minted := new(big.Int).SetUint64(tc.expected[num])
			expectedRewardSpec := &reward.RewardSpec{
				Minted:   minted,
				TotalFee: common.Big0,
				BurntFee: common.Big0,
				Proposer: minted,
				Stakers:  common.Big0,
				KIF:      common.Big0,
				KEF:      common.Big0,
				Rewards: map[common.Address]*big.Int{
					proposer: minted,
				},
			}
			assert.Equal(t, expectedRewardSpec, rewardSpec, "wrong at block %d", num)
		}
	}
}

func TestGetRewardsAccumulated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBlockchain := mocks.NewMockBlockChain(mockCtrl)
	mockGovEngine := NewMockEngine(mockCtrl)
	db := database.NewMemoryDBManager()

	// prepare configurations and data for the test environment
	chainConfig := params.MainnetChainConfig.Copy()
	chainConfig.KoreCompatibleBlock = big.NewInt(0)
	chainConfig.Governance.Reward.Ratio = "50/20/30"
	chainConfig.Governance.Reward.Kip82Ratio = params.DefaultKip82Ratio

	govParamSet, err := params.NewGovParamSetChainConfig(chainConfig)
	if err != nil {
		t.Fatal(err)
	}

	oldSm := reward.GetStakingManager()
	defer reward.SetTestStakingManager(oldSm)
	reward.SetTestStakingManagerWithChain(mockBlockchain, mockGovEngine, db)

	testAddrList := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
	}

	testStakingAmountList := []uint64{
		uint64(5000000),
		uint64(10000000),
		uint64(15000000),
		uint64(20000000),
	}

	stInfo := reward.StakingInfo{
		BlockNum:              0,
		CouncilNodeAddrs:      testAddrList,
		CouncilStakingAddrs:   testAddrList,
		CouncilRewardAddrs:    testAddrList,
		KEFAddr:               common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"),
		KIFAddr:               common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
		CouncilStakingAmounts: testStakingAmountList,
	}

	siBytes, _ := json.Marshal(stInfo)
	if err := db.WriteStakingInfo(stInfo.BlockNum, siBytes); err != nil {
		t.Fatal(err)
	}

	startBlockNum := 0
	endBlockNum := 10
	blocks := make([]*types.Block, endBlockNum-startBlockNum+1)
	receipts := make([]types.Receipts, endBlockNum-startBlockNum+1)
	// set testing data for mock instances
	for i := startBlockNum; i <= endBlockNum; i++ {
		blocks[i] = types.NewBlockWithHeader(&types.Header{
			Number:     big.NewInt(int64(i)),
			Rewardbase: stInfo.CouncilRewardAddrs[i%4], // round-robin way
			GasUsed:    uint64(1000),
			BaseFee:    big.NewInt(25 * params.Gkei),
			Time:       big.NewInt(int64(1000 + i)),
		})
		receipts[i] = types.Receipts{
			&types.Receipt{GasUsed: uint64(1000)},
		}
		mockBlockchain.EXPECT().GetBlock(blocks[i].Hash(), uint64(i)).Return(blocks[i]).AnyTimes()
		mockBlockchain.EXPECT().GetHeaderByNumber(uint64(i)).Return(blocks[i].Header()).AnyTimes()
		mockBlockchain.EXPECT().GetReceiptsByBlockHash(blocks[i].Hash()).Return(receipts[i]).AnyTimes()
	}

	mockBlockchain.EXPECT().Config().Return(chainConfig).AnyTimes()
	mockBlockchain.EXPECT().CurrentBlock().Return(blocks[endBlockNum]).AnyTimes()
	mockGovEngine.EXPECT().EffectiveParams(gomock.Any()).Return(govParamSet, nil).AnyTimes()
	mockGovEngine.EXPECT().BlockChain().Return(mockBlockchain).AnyTimes()

	// execute a target function
	govAPI := NewGovernanceAPI(mockGovEngine)
	ret, err := govAPI.GetRewardsAccumulated(rpc.BlockNumber(startBlockNum), rpc.BlockNumber(endBlockNum))
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ret)

	// pre-calculated estimated rewards per a block
	blockMinted, _ := new(big.Int).SetString("9600000000000000000", 10)  // 9.6 KAIA
	blockProposer, _ := new(big.Int).SetString("960000000000000000", 10) // 0.96 KAIA = 9.6 KAIA * 0.5 * 0.2
	blockStaking, _ := new(big.Int).SetString("3840000000000000000", 10) // 3.84 KAIA = 9.6 KAIA * 0.5 * 0.8
	blockTxFee, _ := new(big.Int).SetString("25000000000000", 10)        // 25000 gkei = 1000 * 25 gkei
	blockTxBurnt := blockTxFee
	blockKIF, _ := new(big.Int).SetString("1920000000000000000", 10) //  1.92 KAIA = 9.6 KAIA * 0.2
	blockKEF, _ := new(big.Int).SetString("2880000000000000000", 10) //  2.88 KAIA = 9.6 KAIA * 0.3

	// check the execution result
	assert.Equal(t, time.Unix(blocks[startBlockNum].Time().Int64(), 0).String(), ret.FirstBlockTime)
	assert.Equal(t, time.Unix(blocks[endBlockNum].Time().Int64(), 0).String(), ret.LastBlockTime)
	assert.Equal(t, uint64(startBlockNum), ret.FirstBlock.Uint64())
	assert.Equal(t, uint64(endBlockNum), ret.LastBlock.Uint64())

	blockCount := big.NewInt(int64(endBlockNum - startBlockNum + 1))
	assert.Equal(t, new(big.Int).Mul(blockMinted, blockCount), ret.TotalMinted)
	assert.Equal(t, new(big.Int).Mul(blockTxFee, blockCount), ret.TotalTxFee)
	assert.Equal(t, new(big.Int).Mul(blockTxBurnt, blockCount), ret.TotalBurntTxFee)
	assert.Equal(t, new(big.Int).Mul(blockProposer, blockCount), ret.TotalProposerRewards)
	assert.Equal(t, new(big.Int).Mul(blockStaking, blockCount), ret.TotalStakingRewards)
	assert.Equal(t, new(big.Int).Mul(blockKIF, blockCount), ret.TotalKIFRewards)
	assert.Equal(t, new(big.Int).Mul(blockKEF, blockCount), ret.TotalKEFRewards)

	gcReward := big.NewInt(0)
	for acc, bal := range ret.Rewards {
		if acc != stInfo.KIFAddr && acc != stInfo.KEFAddr {
			gcReward.Add(gcReward, bal)
		}
	}
	assert.Equal(t, gcReward, new(big.Int).Add(ret.TotalStakingRewards, ret.TotalProposerRewards))
}

func (bc *testBlockChain) Engine() consensus.Engine                    { return nil }
func (bc *testBlockChain) GetHeader(common.Hash, uint64) *types.Header { return nil }
func (bc *testBlockChain) GetHeaderByNumber(val uint64) *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(val),
	}
}

func (bc *testBlockChain) GetReceiptsByBlockHash(hash common.Hash) types.Receipts {
	return types.Receipts{
		&types.Receipt{GasUsed: 10},
		&types.Receipt{GasUsed: 10},
	}
}

func (bc *testBlockChain) GetBlockByNumber(num uint64) *types.Block {
	return types.NewBlockWithHeader(bc.GetHeaderByNumber(num))
}
func (bc *testBlockChain) StateAt(root common.Hash) (*state.StateDB, error) { return nil, nil }
func (bc *testBlockChain) State() (*state.StateDB, error)                   { return nil, nil }
func (bc *testBlockChain) Config() *params.ChainConfig {
	return bc.config
}

func (bc *testBlockChain) Processor() blockchain.Processor {
	return nil
}

func (bc *testBlockChain) StateCache() state.Database {
	return nil
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(bc.CurrentHeader())
}

func (bc *testBlockChain) CurrentHeader() *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(bc.num),
	}
}

func (bc *testBlockChain) SetBlockNum(num uint64) {
	bc.num = num
}

func (bc *testBlockChain) GetBlock(hash common.Hash, num uint64) *types.Block {
	return bc.GetBlockByNumber(num)
}
