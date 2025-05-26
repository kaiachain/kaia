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

		demoted, err := v.getDemotedValidators(valset.NewAddressSet(council), 1)
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
