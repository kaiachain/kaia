package impl

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/log"
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

func getTestChainConfig() *params.ChainConfig {
	cc := params.MainnetChainConfig.Copy()
	cc.Istanbul.Epoch = 100
	return cc
}

func getTestChainConfigKore() *params.ChainConfig {
	cc := params.MainnetChainConfig.Copy()
	cc.Istanbul.Epoch = 100
	cc.KoreCompatibleBlock = common.Big0
	return cc
}

// genesis block must have the default governance params
func newHeaderGovModule(t *testing.T, config *params.ChainConfig) *headerGovModule {
	var (
		chain  = mocks.NewMockBlockChain(gomock.NewController(t))
		valSet = mock_valset.NewMockValsetModule(gomock.NewController(t))
		dbm    = database.NewMemoryDBManager()
		db     = dbm.GetMemDB()
		b      = GetGenesisGovBytes(config)
	)

	WriteVoteDataBlockNums(db, StoredUint64Array{})
	WriteGovDataBlockNums(db, StoredUint64Array{0})
	genesisHeader := &types.Header{
		Number:     big.NewInt(0),
		Governance: b,
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

	valSet.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{validVoter}, nil).AnyTimes()
	valSet.EXPECT().GetProposer(uint64(1), uint64(0)).Return(validVoter, nil).AnyTimes()

	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		Chain:       chain,
		ValSet:      valSet,
		ChainKv:     db,
		ChainConfig: config,
		NodeAddress: config.Governance.GoverningNode,
	})
	require.NoError(t, err)
	WriteLowestVoteScannedEpochIdx(db, 0)
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

func TestInitialDB(t *testing.T) {
	config := getTestChainConfig()

	h := newHeaderGovModule(t, config)
	require.NotNil(t, h)

	assert.Nil(t, ReadLegacyIdxHistory(h.ChainKv))
	assert.Equal(t, StoredUint64Array{0}, ReadGovDataBlockNums(h.ChainKv))
	assert.Equal(t, StoredUint64Array{}, ReadVoteDataBlockNums(h.ChainKv))
	assert.Equal(t, uint64(0), *ReadLowestVoteScannedEpochIdx(h.ChainKv))
}

func TestGetGenesisParamNames(t *testing.T) {
	magmaGenesisConfig := params.MainnetChainConfig.Copy()
	magmaGenesisConfig.MagmaCompatibleBlock = new(big.Int).SetUint64(0)
	magmaGenesisConfig.Governance.KIP71 = params.GetDefaultKIP71Config()

	koreGenesisConfig := magmaGenesisConfig.Copy()
	koreGenesisConfig.KoreCompatibleBlock = new(big.Int).SetUint64(0)
	koreGenesisConfig.Governance.GovParamContract = common.HexToAddress("0x123")
	koreGenesisConfig.Governance.Reward.Kip82Ratio = "20/80"

	testcases := []struct {
		desc     string
		config   *params.ChainConfig
		expected []gov.ParamName
	}{
		{
			desc:   "Mainnet config",
			config: params.MainnetChainConfig.Copy(),
			expected: []gov.ParamName{
				gov.GovernanceGovernanceMode, gov.GovernanceGoverningNode, gov.GovernanceUnitPrice,
				gov.RewardMintingAmount, gov.RewardRatio, gov.RewardUseGiniCoeff,
				gov.RewardDeferredTxFee, gov.RewardMinimumStake,
				gov.RewardStakingUpdateInterval, gov.RewardProposerUpdateInterval,
				gov.IstanbulEpoch, gov.IstanbulPolicy, gov.IstanbulCommitteeSize,
			},
		},
		{
			desc:   "Kairos config",
			config: params.KairosChainConfig.Copy(),
			expected: []gov.ParamName{
				gov.GovernanceGovernanceMode, gov.GovernanceGoverningNode, gov.GovernanceUnitPrice,
				gov.RewardMintingAmount, gov.RewardRatio, gov.RewardUseGiniCoeff,
				gov.RewardDeferredTxFee, gov.RewardMinimumStake,
				gov.RewardStakingUpdateInterval, gov.RewardProposerUpdateInterval,
				gov.IstanbulEpoch, gov.IstanbulPolicy, gov.IstanbulCommitteeSize,
			},
		},
		{
			desc:   "Private net config - genesis is Magma",
			config: magmaGenesisConfig,
			expected: []gov.ParamName{
				gov.GovernanceGovernanceMode, gov.GovernanceGoverningNode, gov.GovernanceUnitPrice,
				gov.RewardMintingAmount, gov.RewardRatio, gov.RewardUseGiniCoeff,
				gov.RewardDeferredTxFee, gov.RewardMinimumStake,
				gov.RewardStakingUpdateInterval, gov.RewardProposerUpdateInterval,
				gov.IstanbulEpoch, gov.IstanbulPolicy, gov.IstanbulCommitteeSize,
				gov.Kip71LowerBoundBaseFee, gov.Kip71UpperBoundBaseFee, gov.Kip71GasTarget,
				gov.Kip71BaseFeeDenominator, gov.Kip71MaxBlockGasUsedForBaseFee,
			},
		},
		{
			desc:   "Private net config - genesis is Kore",
			config: koreGenesisConfig,
			expected: []gov.ParamName{
				gov.GovernanceGovernanceMode, gov.GovernanceGoverningNode, gov.GovernanceUnitPrice,
				gov.RewardMintingAmount, gov.RewardRatio, gov.RewardUseGiniCoeff,
				gov.RewardDeferredTxFee, gov.RewardMinimumStake,
				gov.RewardStakingUpdateInterval, gov.RewardProposerUpdateInterval,
				gov.IstanbulEpoch, gov.IstanbulPolicy, gov.IstanbulCommitteeSize,
				gov.Kip71LowerBoundBaseFee, gov.Kip71UpperBoundBaseFee, gov.Kip71GasTarget,
				gov.Kip71BaseFeeDenominator, gov.Kip71MaxBlockGasUsedForBaseFee,
				gov.GovernanceDeriveShaImpl, gov.GovernanceGovParamContract, gov.RewardKip82Ratio,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, getGenesisParamNames(tc.config))
		})
	}

	// this prevents forgetting to update getGenesisParamNames after adding a new governance parameter
	t.Run("getGenesisParamNames must include all governance parameters when all hardforks are enabled", func(t *testing.T) {
		latestGenesisConfig := koreGenesisConfig.Copy()

		// Set all *CompatibleBlock fields to zero.
		configValue := reflect.ValueOf(latestGenesisConfig).Elem()
		configType := configValue.Type()

		for i := 0; i < configType.NumField(); i++ {
			field := configType.Field(i)
			if strings.HasSuffix(field.Name, "CompatibleBlock") {
				fieldValue := configValue.Field(i)
				if fieldValue.Type() == reflect.TypeOf((*big.Int)(nil)) {
					fieldValue.Set(reflect.ValueOf(big.NewInt(0)))
				}
			}
		}

		assert.Equal(t, len(gov.Params), len(getGenesisParamNames(latestGenesisConfig)))
	})
}

func TestKairosGenesisHash(t *testing.T) {
	kairosHash := params.KairosGenesisHash
	genesis := blockchain.DefaultKairosGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(common.Hash{}, db)
	if block.Hash() != kairosHash {
		t.Errorf("Generated hash is not equal to Kairos's hash. Want %v, Have %v", kairosHash.String(), block.Hash().String())
	}
}

func TestMainnetGenesisHash(t *testing.T) {
	mainnetHash := params.MainnetGenesisHash
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(common.Hash{}, db)
	if block.Hash() != mainnetHash {
		t.Errorf("Generated hash is not equal to Mainnet's hash. Want %v, Have %v", mainnetHash.String(), block.Hash().String())
	}
}

func makeEmptyBlock(num uint64) *types.Block {
	return types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(num))})
}

func TestAccumulateVotesInEpoch(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	var (
		config    = getTestChainConfigKore()
		mockChain = mocks.NewMockBlockChain(gomock.NewController(t))
		dbm       = database.NewMemoryDBManager()
		db        = dbm.GetMemDB()
	)

	for i := uint64(0); i <= 222; i++ {
		header := makeEmptyBlock(i).Header()
		mockChain.EXPECT().GetHeaderByNumber(i).Return(header).AnyTimes()
	}
	mockChain.EXPECT().CurrentBlock().Return(makeEmptyBlock(222)).AnyTimes()
	assert.Nil(t, ReadLowestVoteScannedEpochIdx(db))

	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		ChainConfig: config,
		ChainKv:     db,
		Chain:       mockChain,
		NodeAddress: config.Governance.GoverningNode,
	})
	require.Nil(t, err)
	// Init calls h.accumulateVotesInEpoch(2)
	assert.Equal(t, uint64(2), *ReadLowestVoteScannedEpochIdx(db))

	h.accumulateVotesInEpoch(1)
	assert.Equal(t, uint64(1), *ReadLowestVoteScannedEpochIdx(db))

	h.accumulateVotesInEpoch(0)
	assert.Equal(t, uint64(0), *ReadLowestVoteScannedEpochIdx(db))
}

func TestMigration(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	var (
		config        = getTestChainConfigKore()
		mockChain     = mocks.NewMockBlockChain(gomock.NewController(t))
		dbm           = database.NewMemoryDBManager()
		db            = dbm.GetMemDB()
		governingNode = config.Governance.GoverningNode
		vote50, _     = headergov.NewVoteData(governingNode, string(gov.GovernanceUnitPrice), uint64(100)).ToVoteBytes()
		vote150, _    = headergov.NewVoteData(governingNode, string(gov.GovernanceUnitPrice), uint64(200)).ToVoteBytes()
		vote200, _    = headergov.NewVoteData(governingNode, string(gov.GovernanceUnitPrice), uint64(300)).ToVoteBytes()
		vote222, _    = headergov.NewVoteData(governingNode, string(gov.GovernanceUnitPrice), uint64(400)).ToVoteBytes()
		gov0          = GetGenesisGovBytes(config)
		gov100, _     = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(100)}).ToGovBytes()
		gov200, _     = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(200)}).ToGovBytes()
		votes         = map[uint64][]byte{50: vote50, 150: vote150, 200: vote200, 222: vote222}
		govs          = map[uint64][]byte{0: gov0, 100: gov100, 200: gov200}
		expected      = map[uint64]uint64{0: config.UnitPrice, 100: config.UnitPrice, 200: 100}
	)

	// 1. Setup DB and Headers
	legacyIdxHistory, _ := json.Marshal([]uint64{0, 100, 200})
	db.Put(legacyIdxHistoryKey, legacyIdxHistory)

	for i := uint64(0); i <= 222; i++ {
		header := &types.Header{Number: big.NewInt(int64(i)), Vote: votes[i], Governance: govs[i]}
		mockChain.EXPECT().GetHeaderByNumber(i).Return(header).AnyTimes()
	}
	mockChain.EXPECT().CurrentBlock().Return(makeEmptyBlock(222)).AnyTimes()

	assert.Nil(t, ReadLowestVoteScannedEpochIdx(db))
	assert.Nil(t, ReadVoteDataBlockNums(db))
	assert.Nil(t, ReadGovDataBlockNums(db))
	assert.Equal(t, StoredUint64Array{0, 100, 200}, ReadLegacyIdxHistory(db))

	// 2. Run Init() and check resulting initial schema
	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		ChainConfig: config,
		ChainKv:     db,
		Chain:       mockChain,
		NodeAddress: config.Governance.GoverningNode,
	})
	require.Nil(t, err)

	assert.Equal(t, uint64(2), *ReadLowestVoteScannedEpochIdx(db))
	assert.Equal(t, StoredUint64Array{200, 222}, ReadVoteDataBlockNums(db))
	assert.Equal(t, []uint64{200, 222}, h.VoteBlockNums())
	assert.Equal(t, StoredUint64Array{0, 100, 200}, ReadGovDataBlockNums(db))
	assert.Equal(t, []uint64{0, 100, 200}, h.GovBlockNums())

	// 3. Before migration, check GetParamSet() results
	for num, expectedValue := range expected {
		paramSet := h.GetParamSet(num)
		assert.Equal(t, expectedValue, paramSet.UnitPrice, "wrong at block %d", num)
	}

	// 4. Run migrate() and check resulting migrated schema
	h.wg.Add(1)
	h.migrate()

	assert.Equal(t, uint64(0), *ReadLowestVoteScannedEpochIdx(db))
	assert.Equal(t, StoredUint64Array{50, 150, 200, 222}, ReadVoteDataBlockNums(db))
	assert.Equal(t, []uint64{50, 150, 200, 222}, h.VoteBlockNums())
	assert.Equal(t, StoredUint64Array{0, 100, 200}, ReadGovDataBlockNums(db))
	assert.Equal(t, []uint64{0, 100, 200}, h.GovBlockNums())

	for num, expectedValue := range expected {
		paramSet := h.GetParamSet(num)
		assert.Equal(t, expectedValue, paramSet.UnitPrice, "wrong at block %d", num)
	}
}
