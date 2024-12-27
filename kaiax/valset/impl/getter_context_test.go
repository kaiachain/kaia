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

	"github.com/golang/mock/gomock"
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

// TestGetCommittee tests various branches of getCommittee().
// Individual selectors are tested as separate tests.
func TestGetCommittee(t *testing.T) {
	var (
		qualified = numsToAddrs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
		aM        = uint64(2000000)
		si        = &staking.StakingInfo{
			NodeIds:          qualified,
			StakingContracts: qualified,
			RewardAddrs:      qualified,
			StakingAmounts:   []uint64{aM, aM, aM, aM, aM, aM, aM, aM, aM, aM},
		}
		c = &blockContext{
			num:       101,
			qualified: valset.NewAddressSet(qualified),
			prevHeader: &types.Header{
				Number:  big.NewInt(100),
				MixHash: common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84").Bytes(),
			},
			prevProposer: numToAddr(1),
			rules: params.Rules{
				IsIstanbul:  true,
				IsLondon:    true,
				IsEthTxType: true,
				IsMagma:     true,
			},
			pset: gov.ParamSet{
				CommitteeSize:          6,
				UseGiniCoeff:           true,
				ProposerUpdateInterval: 100,
			},
		}
	)

	testcases := []struct {
		desc     string
		isKore   bool
		isRandao bool
		policy   istanbul.ProposerPolicy
		expected []common.Address
	}{
		// Refer to other Test*Proposer tests for the expected results.
		{"RoundRobin", false, false, istanbul.RoundRobin, numsToAddrs(2, 3, 4, 1, 7, 5)},                      // prev=1, curr=2, next=3
		{"Sticky", false, false, istanbul.Sticky, numsToAddrs(1, 2, 4, 3, 7, 5)},                              // prev=1, curr=1, next=2
		{"WeightedRandom, before Kore", false, false, istanbul.WeightedRandom, numsToAddrs(4, 6, 2, 1, 7, 3)}, // curr=4, next=6, list=[4,6,1,7,...]
		{"WeightedRandom, after Kore", true, false, istanbul.WeightedRandom, numsToAddrs(7, 1, 3, 2, 6, 4)},   // curr=7, next=1, list=[7,1,0,4,...]
		{"WeightedRandom, after Randao", true, true, istanbul.WeightedRandom, numsToAddrs(8, 3, 5, 1, 0, 9)},  // committee=[8,3,5,1,0,9]
	}

	var (
		ctrl          = gomock.NewController(t)
		mockChain     = chain_mock.NewMockBlockChain(ctrl)
		mockGov       = gov_mock.NewMockGovModule(ctrl)
		mockStaking   = staking_mock.NewMockStakingModule(ctrl)
		pListCache, _ = lru.New(128)
		rVoteCache, _ = lru.New(128)
		v             = &ValsetModule{
			InitOpts: InitOpts{
				ChainKv:       database.NewMemoryDBManager().GetMiscDB(),
				Chain:         mockChain,
				GovModule:     mockGov,
				StakingModule: mockStaking,
			},
			proposerListCache: pListCache,
			removeVotesCache:  rVoteCache,
		}
	)
	v.validatorVoteBlockNumsCache = []uint64{}
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(c.prevHeader).AnyTimes()
	mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(si, nil).AnyTimes()

	for _, tc := range testcases {
		c.pset.ProposerPolicy = uint64(tc.policy)
		c.rules.IsKore = tc.isKore
		c.rules.IsShanghai = tc.isRandao
		c.rules.IsCancun = tc.isRandao
		c.rules.IsRandao = tc.isRandao

		v.proposerListCache.Purge()
		v.removeVotesCache.Purge()
		committee, err := v.getCommittee(c, 0)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, committee, tc.desc)
	}
}

func TestRandomCommittee(t *testing.T) {
	var (
		qualified            = valset.NewAddressSet(numsToAddrs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
		prevHash             = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84")
		currProposer         = numToAddr(3)
		nextDistinctProposer = numToAddr(7)
	)
	assert.Equal(t, numsToAddrs(3, 7), selectRandomCommittee(qualified, 2, prevHash, currProposer, nextDistinctProposer))
	assert.Equal(t, numsToAddrs(3, 7, 6, 8, 1, 2), selectRandomCommittee(qualified, 6, prevHash, currProposer, nextDistinctProposer))
}

func TestRandaoCommittee(t *testing.T) {
	var (
		qualified = valset.NewAddressSet(numsToAddrs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
		mixHash   = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84").Bytes()
	)
	assert.Equal(t, numsToAddrs(8), selectRandaoCommittee(qualified, 1, mixHash))
	assert.Equal(t, numsToAddrs(8, 3, 5, 1, 0, 9), selectRandaoCommittee(qualified, 6, mixHash))
}

// TestGetProposer tests various branches of getProposer().
// Individual selectors are tested as separate tests.
func TestGetProposer(t *testing.T) {
	var (
		qualified = numsToAddrs(0, 1, 2, 3)
		aM        = uint64(2000000)
		si        = &staking.StakingInfo{
			NodeIds:          qualified,
			StakingContracts: qualified,
			RewardAddrs:      qualified,
			StakingAmounts:   []uint64{aM, aM, aM, 200 * aM},
		}
		c = &blockContext{
			num:       101,
			qualified: valset.NewAddressSet(qualified),
			prevHeader: &types.Header{
				Number:  big.NewInt(100),
				MixHash: common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84").Bytes(),
			},
			prevProposer: numToAddr(1),
			rules: params.Rules{
				IsIstanbul:  true,
				IsLondon:    true,
				IsEthTxType: true,
				IsMagma:     true,
			},
			pset: gov.ParamSet{
				CommitteeSize:          4,
				UseGiniCoeff:           true,
				ProposerUpdateInterval: 100,
			},
		}
	)

	testcases := []struct {
		desc     string
		isKore   bool
		isRandao bool
		policy   istanbul.ProposerPolicy
		expected common.Address
	}{
		// Refer to other Test*Proposer tests for the expected results.
		{"RoundRobin", false, false, istanbul.RoundRobin, numToAddr(2)},                      // prev=1, curr=2
		{"Sticky", false, false, istanbul.Sticky, numToAddr(1)},                              // prev=1, curr=1
		{"WeightedRandom, before Kore", false, false, istanbul.WeightedRandom, numToAddr(3)}, // list=[(mostly 9)]
		{"WeightedRandom, after Kore", true, false, istanbul.WeightedRandom, numToAddr(2)},   // list=[2,1,3,0]
		{"WeightedRandom, after Randao", true, true, istanbul.WeightedRandom, numToAddr(0)},  // committee=[0,1,2,3]
	}

	var (
		ctrl          = gomock.NewController(t)
		mockChain     = chain_mock.NewMockBlockChain(ctrl)
		mockGov       = gov_mock.NewMockGovModule(ctrl)
		mockStaking   = staking_mock.NewMockStakingModule(ctrl)
		pListCache, _ = lru.New(128)
		rVoteCache, _ = lru.New(128)
		v             = &ValsetModule{
			InitOpts: InitOpts{
				ChainKv:       database.NewMemoryDBManager().GetMiscDB(),
				Chain:         mockChain,
				GovModule:     mockGov,
				StakingModule: mockStaking,
			},
			proposerListCache: pListCache,
			removeVotesCache:  rVoteCache,
		}
	)
	v.validatorVoteBlockNumsCache = []uint64{}
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(c.prevHeader).AnyTimes()
	mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(si, nil).AnyTimes()

	for _, tc := range testcases {
		c.pset.ProposerPolicy = uint64(tc.policy)
		c.rules.IsKore = tc.isKore
		c.rules.IsShanghai = tc.isRandao
		c.rules.IsCancun = tc.isRandao
		c.rules.IsRandao = tc.isRandao

		v.proposerListCache.Purge()
		v.removeVotesCache.Purge()
		proposer, err := v.getProposer(c, 0)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, proposer, tc.desc)
	}
}

func TestRoundRobinProposer(t *testing.T) {
	testcases := []struct {
		qualified    []common.Address
		prevProposer common.Address
		round        uint64
		expected     common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 0, numToAddr(3)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 1, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 2, numToAddr(1)}, // wrap-around by round
		{numsToAddrs(1, 2, 3, 4), numToAddr(4), 0, numToAddr(1)}, // wrap-around by prevProposer
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 0, numToAddr(2)}, // fallback to prevIdx = 0
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 1, numToAddr(3)}, // fallback to prevIdx = 0
	}
	for _, tc := range testcases {
		currProposer := selectRoundRobinProposer(valset.NewAddressSet(tc.qualified), tc.prevProposer, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

func TestStickyProposer(t *testing.T) {
	testcases := []struct {
		qualified    []common.Address
		prevProposer common.Address
		round        uint64
		expected     common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 0, numToAddr(2)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 1, numToAddr(3)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 2, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 3, numToAddr(1)}, // wrap-around by round
		{numsToAddrs(1, 2, 3, 4), numToAddr(4), 0, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 0, numToAddr(1)}, // fallback to prevIdx = 0
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 1, numToAddr(2)}, // fallback to prevIdx = 0
	}
	for _, tc := range testcases {
		currProposer := selectStickyProposer(valset.NewAddressSet(tc.qualified), tc.prevProposer, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

func TestWeightedRandomProposer_Select(t *testing.T) {
	testcases := []struct {
		proposerList  []common.Address
		listSourceNum uint64
		num           uint64
		round         uint64
		expected      common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), 100, 101, 0, numToAddr(1)},
		{numsToAddrs(1, 2, 3, 4), 100, 106, 0, numToAddr(2)}, // wrap-around by num
		{numsToAddrs(1, 2, 3, 4), 100, 103, 3, numToAddr(2)}, // wrap-around by round
	}
	for _, tc := range testcases {
		currProposer := selectWeightedRandomProposer(tc.proposerList, tc.listSourceNum, tc.num, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

// TestCollectStakingAmounts checks if validators and stakingAmounts from a stakingInfo are matched well.
// stakingAmounts of multiple staking contracts will be added to stakingAmounts of validators which have the same reward address.
// input
//   - validator and stakingInfo is matched by a nodeAddress.
//
// output
//   - weightedValidators are sorted by nodeAddress
//   - stakingAmounts should be same as expectedStakingAmounts
func TestCollectStakingAmounts(t *testing.T) {
	uintMS, floatMS := uint64(5000000), float64(5000000)
	testCases := []struct {
		validators             []common.Address
		stakingInfo            *staking.StakingInfo
		expectedStakingAmounts []float64
	}{
		{
			numsToAddrs(1, 2, 3),
			&staking.StakingInfo{
				NodeIds:          numsToAddrs(1, 2, 3),
				StakingContracts: numsToAddrs(1, 2, 3),
				RewardAddrs:      numsToAddrs(4, 5, 6),
				StakingAmounts:   []uint64{2 * uintMS, uintMS, uintMS},
			},
			[]float64{2 * floatMS, floatMS, floatMS},
		},
		{
			numsToAddrs(1, 2, 3, 4),
			&staking.StakingInfo{
				NodeIds:          numsToAddrs(1, 2, 3, 4, 5),
				StakingContracts: numsToAddrs(1, 2, 3, 4, 5),
				RewardAddrs:      numsToAddrs(6, 7, 8, 9, 6),
				StakingAmounts:   []uint64{uintMS, uintMS, uintMS, uintMS, uintMS},
			},
			[]float64{2 * floatMS, floatMS, floatMS, floatMS},
		},
		{
			numsToAddrs(1, 2, 3, 4),
			&staking.StakingInfo{
				NodeIds:          numsToAddrs(1, 2, 3, 4, 5, 6),
				StakingContracts: numsToAddrs(1, 2, 3, 4, 5, 6),
				RewardAddrs:      numsToAddrs(7, 8, 9, 10, 7, 8),
				StakingAmounts:   []uint64{uintMS, uintMS, uintMS, uintMS, uintMS, uintMS},
			},
			[]float64{2 * floatMS, 2 * floatMS, floatMS, floatMS},
		},
	}
	for _, tc := range testCases {
		stakingAmounts := collectStakingAmounts(tc.validators, tc.stakingInfo)
		for idx, nodeId := range tc.validators {
			assert.Equal(t, stakingAmounts[nodeId], tc.expectedStakingAmounts[idx])
		}
	}
}
