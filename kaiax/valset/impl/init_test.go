package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/kaiax/gov"
	hgmmock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/kaiax/staking"
	stakingmock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	chainmock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	tgn               = common.HexToAddress("0x8aD8F547fa00f58A8c4fb3B671Ee5f1A75bA028a") // governing node used in valSet test
	n2                = common.HexToAddress("0xB2AAda7943919e82143324296987f6091F3FDC9e")
	n3                = common.HexToAddress("0xD95c70710f07A3DaF7ae11cFBa10c789da3564D0")
	n4                = common.HexToAddress("0xC704765db1d21C2Ea6F7359dcB8FD5233DeD16b5")
	n5                = common.HexToAddress("0xcb7f556a77a9f7ec8dc16c7d6cf0aeff7c3ee80a")
	n6                = common.HexToAddress("0x0ce5cbdab931fa12821e2a71845fd284bee8914e")
	testGenesisValSet = []common.Address{tgn, n2, n3, n4, n5, n6} // if genesisValSet changed, the cExpectList of voteTestData must be changed

	testPUpdateInterval = uint64(36)
	testSUpdateInterval = uint64(72)
	testEpoch           = uint64(144)
	testSubGroupSize    = uint64(3)
	testProposerPolicy  = params.WeightedRandom
	testGovernanceMode  = "single"

	testIstanbulCompatibleNumber = big.NewInt(int64(testEpoch) + 10)
	testKoreCompatibleBlock      = big.NewInt(int64(testEpoch) + 20)
	testRandaoCompatibleBlock    = big.NewInt(int64(testEpoch) + 30)
	testKaiaCompatibleBlock      = big.NewInt(int64(testEpoch) + 40)

	testMixHash  = common.Hex2Bytes("0x365643b31592079285c66bfbf9215bccc5d977210ec0b28fb2f19b6b92773625")
	testPrevHash = common.HexToHash("0x2e43e99fe04c8664f93f21b255356cd2279b0e993cc367db8d3e629d8c745635")
)

func getValSetStakingInfoTestData() []*staking.StakingInfo {
	// test data for stakingInfoMock
	var (
		s1 = common.HexToAddress("0x4dd324F9821485caE941640B32c3Bcf1fA6E93E6")
		s2 = common.HexToAddress("0x0d5Df5086B5f86f748dFaed5779c3f862C075B1f")
		s3 = common.HexToAddress("0xD3Ff05f00491571E86A3cc8b0c320aA76D7413A5")
		s4 = common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f")
		s5 = common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f")
		s6 = common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f")

		r1 = common.HexToAddress("0x241c793A9AD555f52f6C3a83afe6178408796ab2")
		r2 = common.HexToAddress("0x79b427Fb77077A9716E08D049B0e8f36Abfc8E2E")
		r3 = common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7")
		r4 = common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09")
		r5 = common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7")
		r6 = common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09")

		kef = common.HexToAddress("0x136807B12327a8AfF9831F09617dA1B9D398cda2")
		kif = common.HexToAddress("0x46bA8F7538CD0749e572b2631F9FB4Ce3653AFB8")

		a0 uint64 = 0
		aL uint64 = 1000000  // less than minstaking
		aM uint64 = 2000000  // exactly minstaking (params.DefaultMinimumStake)
		a1 uint64 = 10000000 // above minstaking. Using 1,2,4,8 to uniquely spot errors
	)
	return []*staking.StakingInfo{
		{ // 1. StakingInfo has not been set
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn},
			StakingContracts: []common.Address{s1},
			RewardAddrs:      []common.Address{r1},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{aL},
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn, n2, n3, n4},
			StakingContracts: []common.Address{s1, s2, s3, s4},
			RewardAddrs:      []common.Address{r1, r2, r3, r4},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{a1, aM, aL, a0},
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn, n2, n4, n5, n6},
			StakingContracts: []common.Address{s1, s2, s4, s5, s6},
			RewardAddrs:      []common.Address{r1, r2, r4, r5, r6},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{a1, aM, aL, a0, a1},
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn, n2, n4, n5, n6},
			StakingContracts: []common.Address{s1, s2, s4, s5, s6},
			RewardAddrs:      []common.Address{r1, r2, r4, r5, r6},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{aL, aL, aL, aM, aL},
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn, n2, n4, n5, n6},
			StakingContracts: []common.Address{s1, s2, s4, s5, s6},
			RewardAddrs:      []common.Address{r1, r2, r4, r5, r6},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{aL, aL, aL, aL, aL},
		},
		{
			SourceBlockNum:   0,
			NodeIds:          []common.Address{tgn, n2, n4, n5, n6},
			StakingContracts: []common.Address{s1, s2, s4, s5, s6},
			RewardAddrs:      []common.Address{r1, r2, r4, r5, r6},
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   []uint64{aM, aL, aL, aL, aL},
		},
	}
}

func getValSetParamSetTestData() []gov.ParamSet {
	baseGovParam := func() gov.ParamSet {
		govParam := gov.GetDefaultGovernanceParamSet()
		// initialize the parameters which cannot votable later.
		_ = govParam.Set(gov.RewardProposerUpdateInterval, testPUpdateInterval)
		_ = govParam.Set(gov.RewardStakingUpdateInterval, testSUpdateInterval)
		_ = govParam.Set(gov.IstanbulEpoch, testEpoch)
		_ = govParam.Set(gov.GovernanceGovernanceMode, testGovernanceMode)
		_ = govParam.Set(gov.GovernanceGoverningNode, tgn)
		return *govParam
	}

	var paramSets []gov.ParamSet
	for _, params := range []struct {
		proposerPolicy int
		subGroupSize   uint64
		governingNode  common.Address
	}{
		{testProposerPolicy, 0, tgn},
		{testProposerPolicy, 1, tgn},
		{testProposerPolicy, testSubGroupSize, tgn},
		{testProposerPolicy, testSubGroupSize + 3, tgn},
		{testProposerPolicy, testSubGroupSize + 4, tgn},
		{params.RoundRobin, testSubGroupSize, tgn},
		{params.Sticky, testSubGroupSize, tgn},
	} {
		// initialize the parameters which can be votable later.
		govParam := baseGovParam()
		_ = govParam.Set(gov.IstanbulCommitteeSize, params.subGroupSize)
		_ = govParam.Set(gov.GovernanceGoverningNode, params.governingNode)
		_ = govParam.Set(gov.IstanbulPolicy, params.proposerPolicy)
		paramSets = append(paramSets, govParam)
	}

	return paramSets
}

func getValSetChainHeadersTestData() []*types.Header {
	return []*types.Header{
		{
			ParentHash: testPrevHash,
			MixHash:    nil,
			Number:     big.NewInt(1),
			Vote:       nil,
		},
		{
			ParentHash: testPrevHash,
			MixHash:    nil,
			Number:     testIstanbulCompatibleNumber,
			Vote:       nil,
		},
		{
			ParentHash: testPrevHash,
			MixHash:    nil,
			Number:     testKoreCompatibleBlock,
			Vote:       nil,
		},
		{
			ParentHash: testPrevHash,
			MixHash:    testMixHash,
			Number:     testRandaoCompatibleBlock,
			Vote:       nil,
		},
		{},
	}
}

func newTestVModule(mockChain *chainmock.MockBlockChain, mockEngine *mocks.MockEngine,
	mockHeaderGov *hgmmock.MockHeaderGovModule, mockStaking *stakingmock.MockStakingModule) (
	*ValsetModule, []*types.Header, []*staking.StakingInfo, []gov.ParamSet, error,
) {
	// init vModule
	vModule := NewValsetModule()
	if err := vModule.Init(&InitOpts{
		database.NewMemoryDBManager().GetMiscDB(),
		mockChain,
		mockHeaderGov,
		mockStaking,
		n5,
	}); err != nil {
		return nil, nil, nil, nil, err
	}

	// set initial db
	if err := WriteCouncilAddressListToDb(vModule.ChainKv, 0, testGenesisValSet); err != nil {
		return nil, nil, nil, nil, err
	}

	// set mockEngines
	config := params.KairosChainConfig.Copy()
	config.Governance.Reward.ProposerUpdateInterval = testPUpdateInterval
	config.Governance.Reward.StakingUpdateInterval = testSUpdateInterval
	config.Governance.GoverningNode = tgn
	config.Istanbul.Epoch = testEpoch
	config.Istanbul.SubGroupSize = testSubGroupSize
	config.Istanbul.ProposerPolicy = params.WeightedRandom
	config.IstanbulCompatibleBlock = testIstanbulCompatibleNumber
	config.KoreCompatibleBlock = testKoreCompatibleBlock
	config.RandaoCompatibleBlock = testRandaoCompatibleBlock
	config.KaiaCompatibleBlock = testKaiaCompatibleBlock

	// customize the rest of test data
	testHeaders, testStakingInfos, testParamSets := getValSetChainHeadersTestData(), getValSetStakingInfoTestData(), getValSetParamSetTestData()
	mockChain.EXPECT().Config().Return(config).AnyTimes()
	mockChain.EXPECT().Engine().Return(mockEngine).AnyTimes()
	return vModule, testHeaders, testStakingInfos, testParamSets, nil
}

func TestValsetModule_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		mockChain, mockEngine      = chainmock.NewMockBlockChain(ctrl), mocks.NewMockEngine(ctrl)
		mockStaking, mockHeaderGov = stakingmock.NewMockStakingModule(ctrl), hgmmock.NewMockHeaderGovModule(ctrl)
	)

	vModule, testHeaders, testStakingInfos, testParamSets, err := newTestVModule(mockChain, mockEngine, mockHeaderGov, mockStaking)
	assert.NoError(t, err)

	mockEngine.EXPECT().Author(gomock.Any()).Return(tgn, nil).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(testHeaders[0]).AnyTimes()
	mockChain.EXPECT().GetHeaderByHash(testHeaders[0].Hash()).Return(testHeaders[0]).AnyTimes()
	mockChain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(testHeaders[0])).AnyTimes()
	mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(testStakingInfos[0], nil).AnyTimes()
	mockHeaderGov.EXPECT().EffectiveParamSet(gomock.Any()).Return(testParamSets[0]).AnyTimes()

	// check the result
	sInfo, err := vModule.stakingInfo.GetStakingInfo(0)
	assert.NoError(t, err)
	assert.Equal(t, testStakingInfos[0], sInfo)

	pSet := vModule.headerGov.EffectiveParamSet(0)
	assert.Equal(t, testParamSets[0], pSet)

	author, err := vModule.chain.Engine().Author(testHeaders[0])
	assert.NoError(t, err)
	assert.Equal(t, tgn, author)

	header := vModule.chain.GetHeaderByNumber(0)
	assert.Equal(t, testHeaders[0], header)
	assert.Equal(t, testHeaders[0], vModule.chain.GetHeaderByHash(header.Hash()))
	assert.Equal(t, testRandaoCompatibleBlock, vModule.chain.Config().RandaoCompatibleBlock)
	assert.Equal(t, types.NewBlockWithHeader(testHeaders[0]), vModule.chain.CurrentBlock())
}
