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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPostInsertBlock(t *testing.T) {
	var (
		governingNode  = numToAddr(3)
		pset           = gov.ParamSet{GoverningNode: governingNode}
		genesisCouncil = numsToAddrs(1, 2, 3)

		voteAdd1, _ = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(1)).ToVoteBytes()
		block1      = types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
			Vote:   voteAdd1,
		})
		voteAdd6, _ = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(6)).ToVoteBytes()
		block2      = types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(2),
			Vote:   voteAdd6,
		})
	)

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
	writeCouncil(db, 0, genesisCouncil)
	writeValidatorVoteBlockNums(db, []uint64{0})
	writeLowestScannedVoteNum(db, 0)
	mockChain.EXPECT().Config().Return(&params.ChainConfig{}).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(makeGenesisBlock(genesisCouncil).Header()).AnyTimes()
	mockGov.EXPECT().GetParamSet(uint64(1)).Return(pset).AnyTimes()
	mockGov.EXPECT().GetParamSet(uint64(2)).Return(pset).AnyTimes()

	// Ineffective vote (adding already existing address)
	assert.NoError(t, v.PostInsertBlock(block1))

	// Effective vote (adding new address)
	assert.NoError(t, v.PostInsertBlock(block2))

	// Check the DB
	assert.Equal(t, []uint64{0, 2}, ReadValidatorVoteBlockNums(db))
	assert.Equal(t, numsToAddrs(1, 2, 3, 6), ReadCouncil(db, 2))
}

func Test_AddRemove(t *testing.T) {
	type vote struct {
		key   gov.ParamName
		value interface{}
	}
	type expected struct {
		validators []int
	}
	type testcase struct {
		length   int
		votes    map[int]vote
		expected map[int]expected
	}

	testcases := []testcase{
		{5, map[int]vote{1: {gov.RemoveValidator, 3}, 3: {gov.AddValidator, 3}}, map[int]expected{0: {[]int{0, 1, 2, 3}}, 1: {[]int{0, 1, 2, 3}}, 2: {[]int{0, 1, 2}}, 3: {[]int{0, 1, 2}}, 4: {[]int{0, 1, 2, 3}}}},
		{5, map[int]vote{1: {gov.RemoveValidator, []int{1, 2, 3}}, 3: {gov.AddValidator, []int{1, 2}}}, map[int]expected{0: {[]int{0, 1, 2, 3}}, 1: {[]int{0, 1, 2, 3}}, 2: {[]int{0}}, 3: {[]int{0}}, 4: {[]int{0, 1, 2}}}},
		{params.CheckpointInterval + 10, map[int]vote{params.CheckpointInterval - 5: {gov.RemoveValidator, 3}, params.CheckpointInterval - 1: {gov.RemoveValidator, 2}, params.CheckpointInterval + 0: {gov.RemoveValidator, 1}, params.CheckpointInterval + 1: {gov.AddValidator, 1}, params.CheckpointInterval + 2: {gov.AddValidator, 2}, params.CheckpointInterval + 3: {gov.AddValidator, 3}}, map[int]expected{0: {[]int{0, 1, 2, 3}}, 1: {[]int{0, 1, 2, 3}}, params.CheckpointInterval - 4: {[]int{0, 1, 2}}, params.CheckpointInterval + 0: {[]int{0, 1}}, params.CheckpointInterval + 1: {[]int{0}}, params.CheckpointInterval + 2: {[]int{0, 1}}, params.CheckpointInterval + 3: {[]int{0, 1, 2}}, params.CheckpointInterval + 4: {[]int{0, 1, 2, 3}}, params.CheckpointInterval + 9: {[]int{0, 1, 2, 3}}}},
		{10, map[int]vote{0: {gov.RemoveValidator, 3}, 2: {gov.AddValidator, 3}, 4: {gov.AddValidator, 3}, 6: {gov.RemoveValidator, 3}, 8: {gov.RemoveValidator, 3}}, map[int]expected{1: {[]int{0, 1, 2}}, 3: {[]int{0, 1, 2, 3}}, 5: {[]int{0, 1, 2, 3}}, 7: {[]int{0, 1, 2}}, 9: {[]int{0, 1, 2}}}},
		{10, map[int]vote{0: {gov.RemoveValidator, 3}, 2: {gov.RemoveValidator, 3}, 4: {gov.AddValidator, 3}, 6: {gov.AddValidator, 3}}, map[int]expected{1: {[]int{0, 1, 2}}, 3: {[]int{0, 1, 2}}, 5: {[]int{0, 1, 2, 3}}, 7: {[]int{0, 1, 2, 3}}}},
		{10, map[int]vote{0: {gov.RemoveValidator, []int{2, 3}}, 2: {gov.AddValidator, []int{2, 3}}, 4: {gov.AddValidator, []int{2, 3}}, 6: {gov.RemoveValidator, []int{2, 3}}, 8: {gov.RemoveValidator, []int{2, 3}}}, map[int]expected{1: {[]int{0, 1}}, 3: {[]int{0, 1, 2, 3}}, 5: {[]int{0, 1, 2, 3}}, 7: {[]int{0, 1}}, 9: {[]int{0, 1}}}},
		{10, map[int]vote{0: {gov.RemoveValidator, []int{2, 3}}, 2: {gov.RemoveValidator, []int{2, 3}}, 4: {gov.AddValidator, []int{2, 3}}, 6: {gov.AddValidator, []int{2, 3}}}, map[int]expected{1: {[]int{0, 1}}, 3: {[]int{0, 1}}, 5: {[]int{0, 1, 2, 3}}, 7: {[]int{0, 1, 2, 3}}}},
	}

	genesisCouncil := numsToAddrs(0, 1, 2, 3)
	governingNode := numToAddr(0)

	for _, tc := range testcases {
		ctrl := gomock.NewController(t)
		db := database.NewMemDB()
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockGov := gov_mock.NewMockGovModule(ctrl)
		mockStaking := staking_mock.NewMockStakingModule(ctrl)
		v := &ValsetModule{InitOpts: InitOpts{
			ChainKv:       db,
			Chain:         mockChain,
			GovModule:     mockGov,
			StakingModule: mockStaking,
		}}
		writeCouncil(db, 0, genesisCouncil)
		writeValidatorVoteBlockNums(db, []uint64{0})
		writeLowestScannedVoteNum(db, 0)
		mockChain.EXPECT().Config().Return(&params.ChainConfig{}).AnyTimes()
		mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{GoverningNode: governingNode}).AnyTimes()

		for i := 0; i < tc.length; i++ {
			header := &types.Header{Number: big.NewInt(int64(i + 1))}
			if vtc, ok := tc.votes[i]; ok {
				voteValue := vtc.value
				if idx, ok := voteValue.(int); ok {
					voteValue = numToAddr(idx)
				}
				if indices, ok := voteValue.([]int); ok {
					voteValue = numsToAddrs(indices...)
				}
				voteBytes, err := headergov.NewVoteData(governingNode, string(vtc.key), voteValue).ToVoteBytes()
				assert.NoError(t, err)
				header.Vote = voteBytes
			}
			assert.NoError(t, v.PostInsertBlock(types.NewBlockWithHeader(header)))
		}

		for i, e := range tc.expected {
			council, err := v.GetCouncil(uint64(i) + 1)
			assert.NoError(t, err)
			assert.Equal(t, numsToAddrs(e.validators...), council)
		}
		ctrl.Finish()
	}
}
