package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	gov_mock "github.com/kaiachain/kaia/v2/kaiax/gov/mock"
	"github.com/kaiachain/kaia/v2/kaiax/staking"
	staking_mock "github.com/kaiachain/kaia/v2/kaiax/staking/mock"
	"github.com/kaiachain/kaia/v2/kaiax/valset"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
	chain_mock "github.com/kaiachain/kaia/v2/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetProposers_GetRemoveVotesInInterval(t *testing.T) {
	var (
		ctrl                = gomock.NewController(t)
		mockChain           = chain_mock.NewMockBlockChain(ctrl)
		mockStaking         = staking_mock.NewMockStakingModule(ctrl)
		mockGovModule       = gov_mock.NewMockGovModule(ctrl)
		cache, _            = lru.New(128)
		removeVotesCache, _ = lru.New(128)
		v                   = &ValsetModule{
			InitOpts: InitOpts{
				ChainKv:       database.NewMemoryDBManager().GetMiscDB(),
				Chain:         mockChain,
				StakingModule: mockStaking,
				GovModule:     mockGovModule,
			},
			proposerListCache: cache,
			removeVotesCache:  removeVotesCache,
		}
	)

	// Test Data
	// BlockNum     | Council | qualified | demoted | proposers
	// -------------+---------+-----------+---------+-----------
	// 3600         |  1,2,3,4 |  1,2,3,4 |    -    |  1,2,3,4
	// 3601         |  1,2,3,4 |  1,2,3,4 |    -    |  1,2,3,4
	// 3602         |  1,2,3,4 |  1,2,3,4 |    -    |  1,2,3,4
	// 3603         |  1,2,3,4 |  1,2,3,4 |    -    |  1,2,3,4
	// 3604         |  1,2,3,4 |  1,2,4   |    3    |  1,2,3,4 <----- 3 demoted
	// 3605 (R 3,4) |  1,2,3,4 |  1,2,4   |    3    |  1,2,3,4
	// 3606         |  1,2     |  1,2     |    -    |  1,2,3
	// 3607 (A 3,6) |  1,2     |  1,2     |    -    |  1,2,3
	// 3608         |  1,2,3,6 |  1,2,6   |    3    |  1,2,3
	// 3609 (R 3,7) |  1,2,3,6 |  1,2,3,6 |    -    |  1,2,3    <---- 3 requalified
	// 3610         |  1,2,6   |  1,2,6   |    -    |  1,2
	var (
		genesisCouncil          = numsToAddrs(1, 2, 3, 4)
		mockQualifiedValidators = map[uint64][]common.Address{
			3600: numsToAddrs(1, 2, 3, 4),
			3604: numsToAddrs(1, 2, 4), // 3 demoted
			3605: numsToAddrs(1, 2, 4),
			3606: numsToAddrs(1, 2), // 3,4 removeval
			3607: numsToAddrs(1, 2),
			3608: numsToAddrs(1, 2, 6),    // 3,6 addval
			3609: numsToAddrs(1, 2, 3, 6), // 3 requalified
			3610: numsToAddrs(1, 2, 6),    // 3,7 removeval
		}
		mockBaseProposers = map[uint64][]common.Address{
			0:    numsToAddrs(1, 2, 3, 4),
			3600: numsToAddrs(1, 2, 3, 4),
		}
		voteData = []struct {
			blkNum    uint64
			key       gov.ParamName
			addresses []common.Address
		}{
			{3605, gov.RemoveValidator, numsToAddrs(3, 4)},
			{3607, gov.AddValidator, numsToAddrs(3, 6)},
			{3609, gov.RemoveValidator, numsToAddrs(3, 7)},
		}
		expectedProposers = map[uint64][]common.Address{
			3600: numsToAddrs(1, 2, 3, 4),
			3601: numsToAddrs(1, 2, 3, 4),
			3602: numsToAddrs(1, 2, 3, 4),
			3603: numsToAddrs(1, 2, 3, 4),
			3604: numsToAddrs(1, 2, 3, 4),
			3605: numsToAddrs(1, 2, 3, 4),
			3606: numsToAddrs(1, 2, 3),
			3607: numsToAddrs(1, 2, 3),
			3608: numsToAddrs(1, 2, 3),
			3609: numsToAddrs(1, 2, 3),
			3610: numsToAddrs(1, 2),
			3611: numsToAddrs(1, 2),
		}
	)
	// Mock v.InitSchema
	writeLowestScannedVoteNum(v.ChainKv, 0)
	writeValidatorVoteBlockNums(v.ChainKv, []uint64{0})
	writeCouncil(v.ChainKv, 0, genesisCouncil)

	// Mock gov module
	mockGovModule.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		StakingUpdateInterval:  1,
		ProposerUpdateInterval: params.DefaultProposerRefreshInterval,
		ProposerPolicy:         uint64(params.WeightedRandom),
		MinimumStake:           big.NewInt(5000000),
		GovernanceMode:         "single",
		GoverningNode:          numToAddr(0),
	}).AnyTimes()

	// Mock chain headers
	for _, vote := range voteData {
		voteBytes, err := headergov.NewVoteData(numToAddr(0), string(vote.key), vote.addresses).ToVoteBytes()
		assert.NoError(t, err)

		header := &types.Header{Vote: voteBytes, Number: big.NewInt(int64(vote.blkNum))}
		mockChain.EXPECT().GetHeaderByNumber(vote.blkNum).Return(header).AnyTimes()
		assert.NoError(t, v.PostInsertBlock(types.NewBlockWithHeader(header)))
	}
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(&types.Header{}).AnyTimes() // For unused blocks

	// Mock chainConfig
	mockChain.EXPECT().Config().Return(&params.ChainConfig{IstanbulCompatibleBlock: big.NewInt(0)}).AnyTimes()

	// Mock qualified validators
	for blkNum, data := range mockQualifiedValidators {
		stakingAmounts := make([]uint64, len(data))
		for i := 0; i < len(data); i++ {
			stakingAmounts[i] = uint64(5000000)
		}
		mockStaking.EXPECT().GetStakingInfo(blkNum).Return(&staking.StakingInfo{
			NodeIds:          data,
			StakingContracts: data,
			RewardAddrs:      data,
			StakingAmounts:   stakingAmounts,
		}, nil).AnyTimes()
	}
	// Proposer list cache setup
	for updateNum, proposers := range mockBaseProposers {
		v.proposerListCache.Add(updateNum, proposers)
	}

	for currentNum, mockQualified := range mockQualifiedValidators {
		qualified, err := v.getQualifiedValidators(currentNum)
		assert.NoError(t, err)
		assert.Equal(t, mockQualified, qualified.List())
	}

	// Execute test for each block
	for currentNum, expectedProposer := range expectedProposers {
		updateNum := roundDown(currentNum-1, params.DefaultProposerRefreshInterval)
		removeVotes := v.getRemoveVotesInInterval(updateNum, params.DefaultProposerRefreshInterval)
		assert.Equal(t, expectedProposer, removeVotes.filteredProposerList(currentNum, mockBaseProposers[updateNum]))
	}
}

func TestWeightedRandomProposer_ListWeighted(t *testing.T) {
	var (
		blockHash = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84")
		aM        = uint64(2000000) // Exactly minstaking
	)

	testcases := []struct {
		desc         string
		qualified    []common.Address
		amounts      []uint64
		useGini      bool
		expectedFreq []int            // expected appearance frequency
		expectedList []common.Address // expected proposer list. Leave nil to skip check.
	}{
		{
			desc:      "positive stakes, gini",
			qualified: numsToAddrs(0, 1, 2, 3),
			amounts:   []uint64{1 * aM, 2 * aM, 3 * aM, 4 * aM},
			useGini:   true,
			// gini=0.25, exponent=0.8 -> Adjusted staking amounts = [109856,191270,264558,333021] -> Percentile weights = [12,21,29,37]
			expectedFreq: []int{12, 21, 29, 37},
		},
		{
			desc:         "positive stakes, no gini",
			qualified:    numsToAddrs(0, 1, 2, 3),
			amounts:      []uint64{1 * aM, 2 * aM, 3 * aM, 4 * aM},
			useGini:      false,
			expectedFreq: []int{10, 20, 30, 40},
			expectedList: numsToAddrs(1, 1, 3, 2, 0, 3, 2, 3, 1, 1, 3, 1, 3, 2, 3, 1, 2, 2, 0, 3, 3, 2, 3, 3, 2, 1, 1, 1, 3, 3, 2, 3, 1, 2, 3, 1, 2, 3, 2, 2, 0, 3, 2, 2, 1, 3, 1, 0, 2, 2, 2, 1, 1, 3, 2, 3, 1, 2, 3, 3, 0, 3, 3, 3, 2, 3, 3, 3, 2, 3, 3, 3, 0, 3, 2, 0, 1, 3, 2, 2, 1, 2, 3, 1, 3, 3, 0, 0, 3, 1, 2, 3, 3, 2, 2, 3, 2, 0, 3, 2),
		},
		{
			desc:         "zero stakes",
			qualified:    numsToAddrs(0, 1, 2, 3), // Note: validators can be qualified with zero stakes, if all are understaked.
			amounts:      []uint64{0, 0, 0, 0},
			useGini:      false,
			expectedFreq: []int{1, 1, 1, 1},
			expectedList: numsToAddrs(1, 3, 0, 2),
		},
		{
			desc:         "severe inequality",
			qualified:    numsToAddrs(0, 1, 2, 3),
			amounts:      []uint64{aM, aM, aM, 200 * aM},
			useGini:      false,
			expectedFreq: []int{1, 1, 1, 99}, // at least 1 slot is guaranteed for each validator.
		},
		{
			desc:         "identical stakingAmounts",
			qualified:    numsToAddrs(0, 1, 2),
			amounts:      []uint64{aM, aM, aM},
			useGini:      true,
			expectedFreq: []int{33, 33, 33},
		},
		{
			desc:         "multiple of minimum staking",
			qualified:    numsToAddrs(0, 1, 2, 3),
			amounts:      []uint64{aM, 2 * aM, aM, aM},
			useGini:      true,
			expectedFreq: []int{21, 38, 21, 21},
		},
		{
			desc:         "non-multiple of minimum staking",
			qualified:    numsToAddrs(0, 1, 2, 3, 4),
			amounts:      []uint64{324946, 560845, 771786, 967997, 1153934},
			useGini:      true,
			expectedFreq: []int{10, 16, 21, 25, 29},
		},
	}
	for _, tc := range testcases {
		qualified := valset.NewAddressSet(tc.qualified)
		si := &staking.StakingInfo{
			NodeIds:          tc.qualified,
			StakingContracts: tc.qualified,
			RewardAddrs:      tc.qualified,
			StakingAmounts:   tc.amounts,
		}
		proposerList := generateProposerListWeighted(qualified, si, tc.useGini, blockHash)

		freq := make(map[common.Address]int)
		for _, addr := range proposerList {
			freq[addr]++
		}
		for idx, addr := range tc.qualified {
			assert.Equal(t, tc.expectedFreq[idx], freq[addr], tc.desc)
		}

		if tc.expectedList != nil && !assert.Equal(t, tc.expectedList, proposerList, tc.desc) {
			t.Logf("expected: %v", addrsToNums(tc.expectedList))
			t.Logf("actual: %v", addrsToNums(proposerList))
		}
	}
}

func TestWeightedRandomProposer_ListUniform(t *testing.T) {
	var (
		blockHash    = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84")
		qualified    = valset.NewAddressSet(numsToAddrs(0, 1, 2, 3))
		proposerList = generateProposerListUniform(qualified, blockHash)
		expectedList = numsToAddrs(1, 3, 0, 2)
	)
	assert.Equal(t, expectedList, proposerList)
}
