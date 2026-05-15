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

func newValsetModuleWithCouncil(t *testing.T, config *params.ChainConfig, council []common.Address) (*ValsetModule, *gov_mock.MockGovModule, *staking_mock.MockStakingModule) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	db := database.NewMemDB()
	writeLowestScannedVoteNum(db, 0)
	writeValidatorVoteBlockNums(db, []uint64{0})
	writeCouncil(db, 0, council)

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockGov := gov_mock.NewMockGovModule(ctrl)
	mockStaking := staking_mock.NewMockStakingModule(ctrl)
	mockChain.EXPECT().Config().Return(config).AnyTimes()

	v := &ValsetModule{InitOpts: InitOpts{
		ChainKv:       db,
		Chain:         mockChain,
		GovModule:     mockGov,
		StakingModule: mockStaking,
	}}
	return v, mockGov, mockStaking
}

func TestGetDemotedValidators(t *testing.T) {
	var (
		council = numsToAddrs(1, 2, 3, 4, 5)
		aL      = uint64(1000000) // Less than minstaking
		aM      = uint64(2000000) // Exactly minstaking
		pset    = gov.ParamSet{
			GovernanceMode: "none",
			MinimumStake:   big.NewInt(int64(aM)),
		}
		config = &params.ChainConfig{
			IstanbulCompatibleBlock: nil,
		}
		si = &staking.StakingInfo{
			NodeIds:          council,
			StakingContracts: council,
			RewardAddrs:      council,
			// "none mode, some demoted" case in TestFilterValidators
			StakingAmounts: []uint64{0, aL, aM, aM, aM},
		}
	)

	testcases := []struct {
		desc       string
		isIstanbul bool
		policy     istanbul.ProposerPolicy
		demoted    []common.Address
	}{
		{"RoundRobin", false, istanbul.RoundRobin, numsToAddrs()},
		{"Sticky", false, istanbul.Sticky, numsToAddrs()},
		{"WeightedRandom before istanbul", false, istanbul.WeightedRandom, numsToAddrs()},
		{"WeightedRandom after istanbul", true, istanbul.WeightedRandom, numsToAddrs(1, 2)},
	}

	var (
		ctrl        = gomock.NewController(t)
		db          = database.NewMemDB()
		mockChain   = chain_mock.NewMockBlockChain(ctrl)
		mockGov     = gov_mock.NewMockGovModule(ctrl)
		mockStaking = staking_mock.NewMockStakingModule(ctrl)
		v           = &ValsetModule{InitOpts: InitOpts{
			ChainKv:       db,
			Chain:         mockChain,
			GovModule:     mockGov,
			StakingModule: mockStaking,
		}}
	)
	defer ctrl.Finish()
	for _, tc := range testcases {
		if tc.isIstanbul {
			config.IstanbulCompatibleBlock = big.NewInt(1)
		}
		mockChain.EXPECT().Config().Return(config).Times(1)

		pset.ProposerPolicy = uint64(tc.policy)
		mockGov.EXPECT().GetParamSet(gomock.Any()).Return(pset).Times(1)
		mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(si, nil).AnyTimes()

		demoted, err := v.getDemotedValidatorsPermissioned(valset.NewAddressSet(council), 1)
		assert.NoError(t, err)
		assert.Equal(t, tc.demoted, demoted.List(), tc.desc)
	}
}

func TestFilterValidators(t *testing.T) {
	var (
		governingNode = numToAddr(3)
		aL            = uint64(1000000) // Less than minstaking
		aM            = uint64(2000000) // Exactly minstaking
		pset          = gov.ParamSet{
			GoverningNode: governingNode,
			MinimumStake:  big.NewInt(int64(aM)),
		}
	)

	testcases := []struct {
		desc    string
		council []common.Address
		amounts []uint64
		single  bool
		demoted []common.Address
	}{
		{
			desc:    "none mode, all qualified",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{aM, aM, aM, aM, aM},
			single:  false,
			demoted: numsToAddrs(),
		},
		{
			desc:    "none mode, some demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, aL, aM, aM, aM},
			single:  false,
			demoted: numsToAddrs(1, 2),
		},
		{
			desc:    "none mode, all demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, aL, aL},
			single:  false,
			demoted: numsToAddrs(), // If all are demoted, none are demoted.
		},
		{
			desc:    "single mode, all qualified",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{aM, aM, aM, aM, aM},
			single:  true,
			demoted: numsToAddrs(),
		},
		{
			desc:    "single mode, some understaked",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, aL, aM, aM, aM},
			single:  true,
			demoted: numsToAddrs(1, 2),
		},
		{
			desc:    "single mode, governingNode and others are understaked",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, aM, aM},
			single:  true,
			demoted: numsToAddrs(1, 2), // despite governingNode(3) understaked, it is not demoted.
		},
		{
			desc:    "single mode, only governingNode is staked enough",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, aM, aL, aL},
			single:  true,
			demoted: numsToAddrs(1, 2, 4, 5),
		},
		{
			desc:    "single mode, all demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, 0, 0},
			single:  true,
			demoted: numsToAddrs(), // If all are demoted, none are demoted.
		},
	}
	for _, tc := range testcases {
		council := valset.NewAddressSet(tc.council)
		si := &staking.StakingInfo{
			NodeIds:          tc.council,
			StakingContracts: tc.council,
			RewardAddrs:      tc.council,
			StakingAmounts:   tc.amounts,
		}
		if tc.single {
			pset.GovernanceMode = "single"
		} else {
			pset.GovernanceMode = "none"
		}

		demoted := getDemotedValidatorsIstanbul(council, si, pset)
		assert.Equal(t, tc.demoted, demoted.List(), tc.desc)
	}
}

func Test_AfterMinimumStakingVotes(t *testing.T) {
	type vote struct {
		key   gov.ParamName
		value interface{}
	}
	type expected struct {
		blocks     []uint64
		validators []int
		demoted    []int
	}
	type testcase struct {
		stakingAmounts []uint64
		votes          []vote
		expected       []expected
	}

	testcases := []testcase{
		{
			stakingAmounts: []uint64{8000000, 7000000, 6000000, 5000000},
			votes: []vote{
				{gov.GovernanceGovernanceMode, "none"},
				{gov.RewardMinimumStake, uint64(5500000)},
				{gov.RewardMinimumStake, uint64(6500000)},
				{gov.RewardMinimumStake, uint64(7500000)},
				{gov.RewardMinimumStake, uint64(8500000)},
				{gov.RewardMinimumStake, uint64(7500000)},
				{gov.RewardMinimumStake, uint64(6500000)},
				{gov.RewardMinimumStake, uint64(5500000)},
				{gov.RewardMinimumStake, uint64(4500000)},
			},
			expected: []expected{
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
				{[]uint64{12, 13, 14}, []int{0, 1}, []int{2, 3}},
				{[]uint64{15, 16, 17}, []int{0}, []int{1, 2, 3}},
				{[]uint64{18, 19, 20}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{21, 22, 23}, []int{0}, []int{1, 2, 3}},
				{[]uint64{24, 25, 26}, []int{0, 1}, []int{2, 3}},
				{[]uint64{27, 28, 29}, []int{0, 1, 2}, []int{3}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			stakingAmounts: []uint64{5000000, 6000000, 7000000, 8000000},
			votes: []vote{
				{gov.RewardMinimumStake, uint64(8500000)},
				{gov.RewardMinimumStake, uint64(7500000)},
				{gov.RewardMinimumStake, uint64(6500000)},
				{gov.RewardMinimumStake, uint64(5500000)},
				{gov.RewardMinimumStake, uint64(4500000)},
				{gov.RewardMinimumStake, uint64(5500000)},
				{gov.RewardMinimumStake, uint64(6500000)},
				{gov.RewardMinimumStake, uint64(7500000)},
				{gov.RewardMinimumStake, uint64(8500000)},
			},
			expected: []expected{
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 3}, []int{1, 2}},
				{[]uint64{12, 13, 14}, []int{0, 2, 3}, []int{1}},
				{[]uint64{15, 16, 17, 18, 19, 20, 21, 22, 23}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{24, 25, 26}, []int{0, 2, 3}, []int{1}},
				{[]uint64{27, 28, 29}, []int{0, 3}, []int{1, 2}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			stakingAmounts: []uint64{6000000, 6000000, 5000000, 5000000},
			votes: []vote{
				{gov.RewardMinimumStake, uint64(5500000)},
				{gov.GovernanceGoverningNode, 2},
			},
			expected: []expected{
				{[]uint64{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{6, 7, 8}, []int{0, 1}, []int{2, 3}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
			},
		},
	}

	const testEpoch = uint64(3)
	addrs := numsToAddrs(0, 1, 2, 3)
	config := &params.ChainConfig{IstanbulCompatibleBlock: big.NewInt(0)}

	for _, tc := range testcases {
		v, mockGov, mockStaking := newValsetModuleWithCouncil(t, config, addrs)
		mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(&staking.StakingInfo{
			SourceBlockNum:   0,
			NodeIds:          addrs,
			StakingContracts: addrs,
			RewardAddrs:      addrs,
			StakingAmounts:   tc.stakingAmounts,
		}, nil).AnyTimes()

		basePset := gov.ParamSet{
			ProposerPolicy: uint64(istanbul.WeightedRandom),
			GovernanceMode: "single",
			GoverningNode:  addrs[0],
			MinimumStake:   big.NewInt(4_000_000),
		}
		mockGov.EXPECT().GetParamSet(gomock.Any()).DoAndReturn(func(blockNum uint64) gov.ParamSet {
			pset := basePset
			for i, vote := range tc.votes {
				effectiveFrom := uint64(1) + uint64(i+2)*testEpoch
				if blockNum < effectiveFrom {
					break
				}
				switch vote.key {
				case gov.GovernanceGovernanceMode:
					pset.GovernanceMode = vote.value.(string)
				case gov.RewardMinimumStake:
					pset.MinimumStake = new(big.Int).SetUint64(vote.value.(uint64))
				case gov.GovernanceGoverningNode:
					pset.GoverningNode = addrs[vote.value.(int)]
				default:
					t.Fatalf("unexpected vote key: %s", vote.key)
				}
			}
			return pset
		}).AnyTimes()

		for _, e := range tc.expected {
			for _, num := range e.blocks {
				qualified, err := v.GetQualifiedValidators(num + 1)
				assert.NoError(t, err)
				demoted, err := v.GetDemotedValidators(num + 1)
				assert.NoError(t, err)
				assert.Equal(t, numsToAddrs(e.validators...), qualified, "blockNum:%d", num+1)
				assert.Equal(t, numsToAddrs(e.demoted...), demoted, "blockNum:%d", num+1)
			}
		}
	}
}

func Test_AfterKaia_BasedOnStaking(t *testing.T) {
	type testcase struct {
		stakingAmounts     []uint64
		isKaiaCompatible   bool
		expectedValidators []int
		expectedDemoted    []int
	}

	genesisStakingAmounts := []uint64{5000000, 5000000, 5000000, 5000000}
	testcases := []testcase{
		{[]uint64{5000000, 5000000, 5000000, 6000000}, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 6000000, 6000000}, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 5000000, 6000000}, true, []int{3}, []int{0, 1, 2}},
		{[]uint64{5000000, 5000000, 6000000, 6000000}, true, []int{2, 3}, []int{0, 1}},
	}

	addrs := numsToAddrs(0, 1, 2, 3)
	for _, tc := range testcases {
		config := &params.ChainConfig{
			IstanbulCompatibleBlock: big.NewInt(0),
		}
		if tc.isKaiaCompatible {
			config.KaiaCompatibleBlock = big.NewInt(0)
		}

		v, mockGov, mockStaking := newValsetModuleWithCouncil(t, config, addrs)
		mockGov.EXPECT().GetParamSet(uint64(2)).Return(gov.ParamSet{
			ProposerPolicy: uint64(istanbul.WeightedRandom),
			GovernanceMode: "none",
			MinimumStake:   big.NewInt(5_500_000),
		}).AnyTimes()

		amounts := genesisStakingAmounts
		if tc.isKaiaCompatible {
			amounts = tc.stakingAmounts
		}
		mockStaking.EXPECT().GetStakingInfo(uint64(2)).Return(&staking.StakingInfo{
			SourceBlockNum:   1,
			NodeIds:          addrs,
			StakingContracts: addrs,
			RewardAddrs:      addrs,
			StakingAmounts:   amounts,
		}, nil).AnyTimes()

		qualified, err := v.GetQualifiedValidators(2)
		assert.NoError(t, err)
		demoted, err := v.GetDemotedValidators(2)
		assert.NoError(t, err)

		assert.Equal(t, numsToAddrs(tc.expectedValidators...), qualified)
		assert.Equal(t, numsToAddrs(tc.expectedDemoted...), demoted)
	}
}

func Test_BasedOnStaking(t *testing.T) {
	type testcase struct {
		stakingAmounts       []uint64
		isIstanbulCompatible bool
		isSingleMode         bool
		expectedValidators   []int
		expectedDemoted      []int
	}

	testcases := []testcase{
		{[]uint64{5000000, 5000000, 5000000, 5000000}, false, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 5000000, 6000000}, false, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 6000000, 6000000}, false, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 6000000, 6000000, 6000000}, false, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{6000000, 6000000, 6000000, 6000000}, false, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 5000000, 5000000}, true, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{6000000, 5000000, 5000000, 5000000}, true, false, []int{0}, []int{1, 2, 3}},
		{[]uint64{6000000, 5000000, 5000000, 6000000}, true, false, []int{0, 3}, []int{1, 2}},
		{[]uint64{6000000, 5000000, 6000000, 6000000}, true, false, []int{0, 2, 3}, []int{1}},
		{[]uint64{6000000, 6000000, 6000000, 6000000}, true, false, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5500001, 5500000, 5499999, 0}, true, false, []int{0, 1}, []int{2, 3}},
		{[]uint64{6000000, 6000000, 6000000, 6000000}, true, true, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 6000000, 6000000, 6000000}, true, true, []int{0, 1, 2, 3}, []int{}},
		{[]uint64{5000000, 5000000, 6000000, 6000000}, true, true, []int{0, 2, 3}, []int{1}},
		{[]uint64{5000000, 5000000, 5000000, 6000000}, true, true, []int{0, 3}, []int{1, 2}},
		{[]uint64{5000000, 5000000, 5000000, 5000000}, true, true, []int{0, 1, 2, 3}, []int{}},
	}

	addrs := numsToAddrs(0, 1, 2, 3)
	for _, tc := range testcases {
		config := &params.ChainConfig{}
		if tc.isIstanbulCompatible {
			config.IstanbulCompatibleBlock = big.NewInt(0)
		}

		v, mockGov, mockStaking := newValsetModuleWithCouncil(t, config, addrs)
		pset := gov.ParamSet{
			ProposerPolicy: uint64(istanbul.WeightedRandom),
			GovernanceMode: "none",
			GoverningNode:  addrs[0],
			MinimumStake:   big.NewInt(5_500_000),
		}
		if tc.isSingleMode {
			pset.GovernanceMode = "single"
		}
		mockGov.EXPECT().GetParamSet(uint64(1)).Return(pset).AnyTimes()
		mockStaking.EXPECT().GetStakingInfo(uint64(1)).Return(&staking.StakingInfo{
			SourceBlockNum:   0,
			NodeIds:          addrs,
			StakingContracts: addrs,
			RewardAddrs:      addrs,
			StakingAmounts:   tc.stakingAmounts,
		}, nil).AnyTimes()

		qualified, err := v.GetQualifiedValidators(1)
		assert.NoError(t, err)
		demoted, err := v.GetDemotedValidators(1)
		assert.NoError(t, err)

		assert.Equal(t, numsToAddrs(tc.expectedValidators...), qualified)
		assert.Equal(t, numsToAddrs(tc.expectedDemoted...), demoted)
	}
}
