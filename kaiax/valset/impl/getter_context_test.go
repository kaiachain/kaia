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
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRandomCommittee(t *testing.T) {
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
		{"RoundRobin", false, false, istanbul.RoundRobin, numToAddr(2)},
		{"Sticky", false, false, istanbul.Sticky, numToAddr(1)},
		{"WeightedRandom, before Kore", false, false, istanbul.WeightedRandom, numToAddr(3)}, // weighted random (list is mostly 3)
		{"WeightedRandom, after Kore", true, false, istanbul.WeightedRandom, numToAddr(2)},   // uniform random (list is [2,1,3,0])
		{"WeightedRandom, after Randao", true, true, istanbul.WeightedRandom, numToAddr(0)},  // randao (committee is [0,1,2,3])
	}

	var (
		ctrl        = gomock.NewController(t)
		mockChain   = chain_mock.NewMockBlockChain(ctrl)
		mockGov     = gov_mock.NewMockGovModule(ctrl)
		mockStaking = staking_mock.NewMockStakingModule(ctrl)
		cache, _    = lru.New(128)
		v           = &ValsetModule{
			InitOpts: InitOpts{
				Chain:         mockChain,
				GovModule:     mockGov,
				StakingModule: mockStaking,
			},
			proposerListCache: cache,
		}
	)
	mockChain.EXPECT().GetHeaderByNumber(uint64(100)).Return(c.prevHeader).AnyTimes()
	mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(si, nil).AnyTimes()

	for _, tc := range testcases {
		c.pset.ProposerPolicy = uint64(tc.policy)
		c.rules.IsKore = tc.isKore
		c.rules.IsShanghai = tc.isRandao
		c.rules.IsCancun = tc.isRandao
		c.rules.IsRandao = tc.isRandao

		v.proposerListCache.Purge()
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
