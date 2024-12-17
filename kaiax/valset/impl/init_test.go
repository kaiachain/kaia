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
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartStop(t *testing.T) {
	var ( // Stuff large enough data for migration thread to run.
		council = numsToAddrs(1, 2, 3, 4)
		genesis = makeGenesisBlock(council).Header()
		current = makeEmptyBlock(12345)
		snap, _ = json.Marshal(istanbulSnapshotStorage{Validators: council})
	)

	var (
		ctrl        = gomock.NewController(t)
		db          = database.NewMemDB()
		mockChain   = chain_mock.NewMockBlockChain(ctrl)
		mockGov     = gov_mock.NewMockGovModule(ctrl)
		mockStaking = staking_mock.NewMockStakingModule(ctrl) // not used in GetCouncil() stage.
		v           = NewValsetModule()
	)
	defer ctrl.Finish()

	mockGov.EXPECT().EffectiveParamSet(gomock.Any()).Return(gov.ParamSet{}).AnyTimes()
	mockChain.EXPECT().CurrentBlock().Return(current).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesis).AnyTimes()
	for i := uint64(1); i <= current.NumberU64(); i++ {
		header := &types.Header{Number: big.NewInt(int64(i))}
		mockChain.EXPECT().GetHeaderByNumber(i).Return(header).AnyTimes()
		if i%1024 == 0 {
			db.Put(istanbulSnapshotKey(header.Hash()), snap)
		}
	}

	require.NoError(t, v.Init(&InitOpts{
		ChainKv:       db,
		Chain:         mockChain,
		GovModule:     mockGov,
		StakingModule: mockStaking,
	}))

	// Start() and Stop() must not block or panic.
	require.NoError(t, v.Start())
	time.Sleep(1 * time.Millisecond)
	v.Stop()
}

// 1. Setup DB and Headers
// 2. Run Init(), check initial schema
// 3. Before migration, check GetCouncil() results
// 4. Run migrate(), check migrated schema
// 5. After migration, check GetCouncil() results
func TestMigration(t *testing.T) {
	// Prepare sample data and testcases.
	// Blocks 1024-1029 has votes. All other blocks have no votes.
	// This blockchain has advanced to block 2050 and 3 istanbul snapshots were produced by old implementation.
	var (
		governingNode = numToAddr(3)
		pset          = gov.ParamSet{GoverningNode: governingNode}

		voteAdd1, _    = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(1)).ToVoteBytes()
		voteAdd2, _    = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(2)).ToVoteBytes()
		voteRemove2, _ = headergov.NewVoteData(governingNode, string(gov.RemoveValidator), numToAddr(2)).ToVoteBytes()
		voteRemove3, _ = headergov.NewVoteData(governingNode, string(gov.RemoveValidator), numToAddr(3)).ToVoteBytes()
		voteAdd4, _    = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(4)).ToVoteBytes()
		votes          = map[uint64][]byte{
			1024: voteAdd1,    // +1 since 1025
			1025: voteRemove2, // -2 since 1026
			1026: voteRemove3, // ignored because 3 is governingNode
			1027: voteAdd4,    // +4 since 1028
			1028: voteAdd4,    // ignored because 4 is already in council
			1029: voteRemove2, // ignored because 2 is already not in council
			2049: voteAdd2,    // +2 since 2050
		}

		snap0, _ = json.Marshal(istanbulSnapshotStorage{ // start with [2,3,5] at block 0
			Validators: numsToAddrs(2, 3, 5),
		})
		snap1024, _ = json.Marshal(istanbulSnapshotStorage{ // snapshot of [1,2,3,5] at block 1024 (applied votes up to 1024; effective since 1025)
			Validators:        numsToAddrs(1, 2, 3),
			DemotedValidators: numsToAddrs(5),
		})
		snap2048, _ = json.Marshal(istanbulSnapshotStorage{ // snapshot of [1,3,4,5] at block 2048 (applied votes up to 2048; effective since 2049)
			Validators:        numsToAddrs(1, 3, 4),
			DemotedValidators: numsToAddrs(5),
		})
		genesis   = makeGenesisBlock(numsToAddrs(2, 3, 5)).Header()
		block2050 = makeEmptyBlock(2050)

		expectedCouncils = []struct {
			num     uint64
			council []common.Address
		}{
			{0, numsToAddrs(2, 3, 5)},          // snap0
			{1024, numsToAddrs(2, 3, 5)},       // snap0 (snap1024 includes vote at 1024 and effective since 1025. Therefore block 1024 refers to snap0, not snap1024)
			{1025, numsToAddrs(1, 2, 3, 5)},    // snap1024 + add1
			{1026, numsToAddrs(1, 3, 5)},       // snap1024 + add1 + remove2
			{1027, numsToAddrs(1, 3, 5)},       // snap1024 + add1 + remove2
			{1028, numsToAddrs(1, 3, 4, 5)},    // snap1024 + add1 + remove2 + add4
			{1029, numsToAddrs(1, 3, 4, 5)},    // snap1024 + add1 + remove2 + add4
			{1030, numsToAddrs(1, 3, 4, 5)},    // snap1024 + add1 + remove2 + add4
			{2048, numsToAddrs(1, 3, 4, 5)},    // snap1024 + add1 + remove2 + add4
			{2049, numsToAddrs(1, 3, 4, 5)},    // snap2048
			{2050, numsToAddrs(1, 2, 3, 4, 5)}, // snap2048 + add2
		}
	)

	// 1. Setup DB and Headers
	var (
		ctrl        = gomock.NewController(t)
		db          = database.NewMemDB()
		mockChain   = chain_mock.NewMockBlockChain(ctrl)
		mockGov     = gov_mock.NewMockGovModule(ctrl)
		mockStaking = staking_mock.NewMockStakingModule(ctrl) // not used in GetCouncil() stage.
		v           = NewValsetModule()
	)
	defer ctrl.Finish()

	mockGov.EXPECT().EffectiveParamSet(gomock.Any()).Return(pset).AnyTimes()
	mockChain.EXPECT().CurrentBlock().Return(block2050).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesis).AnyTimes()
	for i := uint64(1); i <= 2050; i++ {
		header := &types.Header{Number: big.NewInt(int64(i)), Vote: votes[i]}
		mockChain.EXPECT().GetHeaderByNumber(i).Return(header).AnyTimes()
	}

	db.Put(istanbulSnapshotKey(mockChain.GetHeaderByNumber(0).Hash()), snap0)
	db.Put(istanbulSnapshotKey(mockChain.GetHeaderByNumber(1024).Hash()), snap1024)
	db.Put(istanbulSnapshotKey(mockChain.GetHeaderByNumber(2048).Hash()), snap2048)

	// 2. Run Init() and check resulting initial schema
	require.NoError(t, v.Init(&InitOpts{
		ChainKv:       db,
		Chain:         mockChain,
		GovModule:     mockGov,
		StakingModule: mockStaking,
	}))
	// After initSchema: DB has mandatory schema + last istanbul snapshot interval (2048..2050)
	assert.Equal(t, []uint64{0, 2049}, ReadValidatorVoteBlockNums(db))
	assert.Equal(t, uint64(2048), *ReadLowestScannedVoteNum(db))
	assert.Equal(t, numsToAddrs(2, 3, 5), ReadCouncil(db, 0))
	assert.Equal(t, numsToAddrs(1, 2, 3, 4, 5), ReadCouncil(db, 2049))

	// 3. Before migration, check GetCouncil() results
	for _, tc := range expectedCouncils {
		// council, _, err := v.replayFromIstanbulSnapshot(tc.num, false)
		council, err := v.getCouncil(tc.num)
		require.NoError(t, err)
		assert.Equal(t, tc.council, council.List(), tc.num)
	}

	// 4. Run migrate() and check resulting migrated schema
	v.wg.Add(1)
	v.migrate()
	assert.Equal(t, []uint64{0, 1024, 1025, 1027, 2049}, ReadValidatorVoteBlockNums(db)) // valid votes
	assert.Equal(t, uint64(0), *ReadLowestScannedVoteNum(db))
	assert.Equal(t, numsToAddrs(2, 3, 5), ReadCouncil(db, 0))          // genesis council
	assert.Equal(t, numsToAddrs(1, 2, 3, 5), ReadCouncil(db, 1024))    // after vote at 1024 (+1)
	assert.Equal(t, numsToAddrs(1, 3, 5), ReadCouncil(db, 1025))       // after vote at 1025 (-2)
	assert.Equal(t, numsToAddrs(1, 3, 4, 5), ReadCouncil(db, 1027))    // after vote at 1027 (+4)
	assert.Equal(t, numsToAddrs(1, 2, 3, 4, 5), ReadCouncil(db, 2049)) // after vote at 2049 (+2)

	// 5. After migration, check GetCouncil() results
	for _, tc := range expectedCouncils {
		council, _, err := v.getCouncilFromIstanbulSnapshot(tc.num, false)
		require.NoError(t, err)
		assert.Equal(t, tc.council, council.List(), tc.num)
	}
}

func makeGenesisBlock(council []common.Address) *types.Block {
	genesisExtra, _ := rlp.EncodeToBytes(&types.IstanbulExtra{
		Validators:    council,
		Seal:          make([]byte, types.IstanbulExtraSeal),
		CommittedSeal: [][]byte{},
	})
	genesisExtraBytes := append(make([]byte, types.IstanbulExtraVanity), genesisExtra...)
	return types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0), Extra: genesisExtraBytes})
}

func makeEmptyBlock(num uint64) *types.Block {
	return types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(num))})
}
