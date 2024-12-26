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
