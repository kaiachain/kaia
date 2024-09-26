package impl

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genesis block must have the default governance params
func newHeaderGovModule(t *testing.T, config *params.ChainConfig) *headerGovModule {
	var (
		chain = mocks.NewMockBlockChain(gomock.NewController(t))
		dbm   = database.NewMemoryDBManager()
		db    = dbm.GetMemDB()

		m      = gov.GetDefaultGovernanceParamSet().ToEnumMap()
		gov, _ = headergov.NewGovData(m).ToGovBytes()
	)

	WriteVoteDataBlockNums(db, &StoredUint64Array{})
	WriteGovDataBlockNums(db, &StoredUint64Array{0})
	genesisHeader := &types.Header{
		Number:     big.NewInt(0),
		Governance: gov,
	}
	dbm.WriteHeader(genesisHeader)
	chain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesisHeader)

	cachingDb := state.NewDatabase(dbm)
	statedb, _ := state.New(common.Hash{}, cachingDb, nil, nil)
	chain.EXPECT().State().Return(statedb, nil).AnyTimes()
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(genesisHeader)).AnyTimes()

	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		Chain:       chain,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)

	return h
}

func TestReadGovVoteBlockNumsFromDB(t *testing.T) {
	paramName := string(gov.GovernanceUnitPrice)
	votes := map[uint64]headergov.VoteData{
		1:   headergov.NewVoteData(common.Address{1}, paramName, uint64(100)),
		50:  headergov.NewVoteData(common.Address{2}, paramName, uint64(200)),
		100: headergov.NewVoteData(common.Address{3}, paramName, uint64(300)),
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)

	db := database.NewMemDB()
	voteDataBlockNums := make(StoredUint64Array, 0, len(votes))
	for num, voteData := range votes {
		headerVoteData, err := voteData.ToVoteBytes()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Vote: headerVoteData})
		voteDataBlockNums = append(voteDataBlockNums, num)
	}
	WriteVoteDataBlockNums(db, &voteDataBlockNums)

	assert.Equal(t, votes, readVoteDataFromDB(chain, db))
}

func TestReadGovDataFromDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()

	ps1 := &gov.ParamSet{UnitPrice: uint64(100)}
	ps2 := &gov.ParamSet{UnitPrice: uint64(200)}

	WriteGovDataBlockNums(db, &StoredUint64Array{1, 2})

	govs := map[uint64]headergov.GovData{
		1: headergov.NewGovData(map[gov.ParamName]any{gov.GovernanceUnitPrice: ps1.UnitPrice}),
		2: headergov.NewGovData(map[gov.ParamName]any{gov.GovernanceUnitPrice: ps2.UnitPrice}),
	}
	for num, govData := range govs {
		headerGovData, err := govData.ToGovBytes()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Governance: headerGovData})
	}

	assert.Equal(t, govs, readGovDataFromDB(chain, db))
}
