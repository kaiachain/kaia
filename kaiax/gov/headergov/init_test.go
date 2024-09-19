package headergov

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genesis block has the default governance params
func newHeaderGovModule(t *testing.T, config *params.ChainConfig) *headerGovModule {
	var (
		chain = mocks.NewMockBlockChain(gomock.NewController(t))
		db    = database.NewMemDB()

		m      = GetDefaultGovernanceParamSet().ToEnumMap()
		gov, _ = NewGovData(m).Serialize()
	)

	WriteVoteDataBlockNums(db, &StoredUint64Array{})
	WriteGovDataBlockNums(db, &StoredUint64Array{0})
	chain.EXPECT().GetHeaderByNumber(uint64(0)).Return(&types.Header{
		Number:     big.NewInt(0),
		Governance: gov,
	})

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
	paramName := Params[GovernanceUnitPrice].Name
	votes := map[uint64]VoteData{
		1:   NewVoteData(common.Address{1}, paramName, uint64(100)),
		50:  NewVoteData(common.Address{2}, paramName, uint64(200)),
		100: NewVoteData(common.Address{3}, paramName, uint64(300)),
	}

	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)

	db := database.NewMemDB()
	voteDataBlockNums := make(StoredUint64Array, 0, len(votes))
	for num, voteData := range votes {
		headerVoteData, err := voteData.Serialize()
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

	ps1 := &ParamSet{
		UnitPrice: uint64(100),
	}
	ps2 := &ParamSet{
		UnitPrice: uint64(200),
	}

	WriteGovDataBlockNums(db, &StoredUint64Array{1, 2})

	govs := map[uint64]GovData{
		1: NewGovData(map[ParamEnum]interface{}{GovernanceUnitPrice: ps1.UnitPrice}),
		2: NewGovData(map[ParamEnum]interface{}{GovernanceUnitPrice: ps2.UnitPrice}),
	}
	for num, govData := range govs {
		headerGovData, err := govData.Serialize()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Governance: headerGovData})
	}

	assert.Equal(t, govs, readGovDataFromDB(chain, db))
}
