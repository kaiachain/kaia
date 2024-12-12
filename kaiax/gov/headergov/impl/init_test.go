package impl

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
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

	h := NewHeaderGovModule()
	err := h.Init(&InitOpts{
		Chain:       chain,
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

func TestInitialDB(t *testing.T) {
	config := params.TestChainConfig.Copy()
	config.Istanbul = &params.IstanbulConfig{Epoch: 10}

	h := newHeaderGovModule(t, config)
	require.NotNil(t, h)

	assert.Nil(t, ReadLegacyIdxHistory(h.ChainKv))
	assert.Equal(t, StoredUint64Array{0}, ReadGovDataBlockNums(h.ChainKv))
	assert.Nil(t, StoredUint64Array(nil), ReadVoteDataBlockNums(h.ChainKv))
	assert.Equal(t, uint64(0), *ReadLowestVoteScannedBlockNum(h.ChainKv))
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
