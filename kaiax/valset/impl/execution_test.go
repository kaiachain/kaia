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
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPostInsertBlock(t *testing.T) {
	var (
		governingNode  = numToAddr(3)
		pset           = gov.ParamSet{GoverningNode: governingNode}
		genesisCouncil = valset.NewCommonAddressSet(numsToAddrs(1, 2, 3))

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
	mockChain.EXPECT().Config().Return(&params.ChainConfig{}).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(makeGenesisBlock(genesisCouncil.Council()).Header()).AnyTimes()

	writeCouncilPermissioned(db, 0, genesisCouncil)
	writeValidatorVoteBlockNums(db, []uint64{0})
	writeLowestScannedVoteNum(db, 0)

	mockGov.EXPECT().GetParamSet(uint64(1)).Return(pset).AnyTimes()
	mockGov.EXPECT().GetParamSet(uint64(2)).Return(pset).AnyTimes()

	// Ineffective vote (adding already existing address)
	assert.NoError(t, v.PostInsertBlock(block1))

	// Effective vote (adding new address)
	assert.NoError(t, v.PostInsertBlock(block2))

	// Check the DB
	assert.Equal(t, []uint64{0, 2}, ReadValidatorVoteBlockNums(db))
	assert.Equal(t, numsToAddrs(1, 2, 3, 6), readCouncilPermissioned(db, 2))
}

// TestPostInsertBlock_PermissionlessIgnoresVote verifies that header votes
// (add/remove validator) are ignored when permissionless fork is enabled.
func TestPostInsertBlock_PermissionlessIgnoresVote(t *testing.T) {
	var (
		governingNode = numToAddr(3)

		voteAdd6, _ = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(6)).ToVoteBytes()
		block1      = types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
			Vote:   voteAdd6,
		})
	)

	var (
		ctrl               = gomock.NewController(t)
		db                 = database.NewMemDB()
		mockChain          = chain_mock.NewMockBlockChain(ctrl)
		nodeStatesCache, _ = lru.New(128)
		v                  = &ValsetModule{
			InitOpts: InitOpts{
				ChainKv: db,
				Chain:   mockChain,
			},
			nodeStatesCache: nodeStatesCache,
		}
	)

	// PermissionlessCompatibleBlock=0 → all blocks are permissionless
	mockChain.EXPECT().Config().Return(&params.ChainConfig{
		PermissionlessCompatibleBlock: big.NewInt(0),
	}).AnyTimes()
	// getOrComputeNodeStates will fail (no header/state), but that's fine —
	// the point is that applyVote is never reached.
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(nil).AnyTimes()

	// Seed initial vote state so we can verify it's unchanged
	writeValidatorVoteBlockNums(db, []uint64{0})

	// PostInsertBlock takes the permissionless path; vote is ignored
	_ = v.PostInsertBlock(block1)

	// Vote was NOT applied — DB unchanged
	assert.Equal(t, []uint64{0}, ReadValidatorVoteBlockNums(db))
	assert.Nil(t, readCouncilPermissioned(db, 1))
}
