package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validVoter = common.Address{1}
	extra      = hexutil.MustDecode("0xd883010703846b6c617988676f312e31352e37856c696e757800000000000000f90164f85494571e53df607be97431a5bbefca1dffe5aef56f4d945cb1a7dccbd0dc446e3640898ede8820368554c89499fb17d324fa0e07f23b49d09028ac0919414db694b74ff9dea397fe9e231df545eb53fe2adf776cb2b841acb7fcc5152506250d1ea49745e7d0d5930157724b410e6e62e0885e7978c81863647d90700dcf3e5d0727cb886f2cc2c63f8f6f3910b4341b302a0aa06eae4500f8c9b841d79c07fbee8861585a71af08a867546320ba804c49c7a3c8641b4d235fd50d5a29bf72d20f3ff1ddfb945ff193d7938967be694f3e602a1cffdea686acf2b0ea01b841dfcf5b5608ca86bc92e7fa3d88a8b25840a629234614ecb312621234ed665ae562ee64ea09fcc88080aaab1ee095acf705d7cc495732682ffee23023ed41feb200b841fefc3b618b2384ea5c7c519ddecc666c19e8a600a6e30c5d9831941c0d5af78d28250bab36ce29202e667c9c1681fd9930aab002988c7228b64caab003bd998100")
)

// genesis block must have the default governance params
func newHeaderGovModule(t *testing.T, config *params.ChainConfig) *headerGovModule {
	var (
		chain  = mocks.NewMockBlockChain(gomock.NewController(t))
		valSet = mock_valset.NewMockValsetModule(gomock.NewController(t))
		dbm    = database.NewMemoryDBManager()
		db     = dbm.GetMemDB()

		m      = gov.GetDefaultGovernanceParamSet().ToMap()
		gov, _ = headergov.NewGovData(m).ToGovBytes()
	)

	WriteVoteDataBlockNums(db, StoredUint64Array{})
	WriteGovDataBlockNums(db, StoredUint64Array{0})
	genesisHeader := &types.Header{
		Number:     big.NewInt(0),
		Governance: gov,
	}
	dbm.WriteHeader(genesisHeader)

	// mock accumulateVotesInEpoch
	chain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesisHeader).AnyTimes()
	for i := uint64(1); i < config.Istanbul.Epoch; i++ {
		chain.EXPECT().GetHeaderByNumber(i).Return(&types.Header{Number: big.NewInt(int64(i))})
	}

	cachingDb := state.NewDatabase(dbm)
	statedb, _ := state.New(common.Hash{}, cachingDb, nil, nil)
	chain.EXPECT().State().Return(statedb, nil).AnyTimes()
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(genesisHeader)).AnyTimes()

	valSet.EXPECT().GetCouncil(uint64(1)).Return(valset.AddressList{validVoter}, nil).AnyTimes()
	valSet.EXPECT().GetProposer(uint64(1), uint64(0)).Return(validVoter, nil).AnyTimes()

	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		Chain:       chain,
		ValSet:      valSet,
		ChainKv:     db,
		ChainConfig: config,
	})
	require.NoError(t, err)
	WriteLowestVoteScannedBlockNum(db, 0)
	h.Start()

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
	WriteVoteDataBlockNums(db, voteDataBlockNums)

	assert.Equal(t, votes, readVoteDataFromDB(chain, db))
}

func TestReadGovDataFromDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(mockCtrl)
	db := database.NewMemDB()

	ps1 := &gov.ParamSet{UnitPrice: uint64(100)}
	ps2 := &gov.ParamSet{UnitPrice: uint64(200)}

	WriteGovDataBlockNums(db, StoredUint64Array{1, 2})

	govs := map[uint64]headergov.GovData{
		1: headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: ps1.UnitPrice}),
		2: headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: ps2.UnitPrice}),
	}
	for num, govData := range govs {
		headerGovData, err := govData.ToGovBytes()
		require.NoError(t, err)
		chain.EXPECT().GetHeaderByNumber(num).Return(&types.Header{Governance: headerGovData})
	}

	assert.Equal(t, govs, readGovDataFromDB(chain, db))
}
