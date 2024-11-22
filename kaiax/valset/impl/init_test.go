package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
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
	// if genesisValSet changed, the cExpectList of voteTestData must be changed
	n = []common.Address{
		common.HexToAddress("0x8aD8F547fa00f58A8c4fb3B671Ee5f1A75bA028a"), // 138 216 245 71
		common.HexToAddress("0xB2AAda7943919e82143324296987f6091F3FDC9e"), // 178 170 218 121
		common.HexToAddress("0xD95c70710f07A3DaF7ae11cFBa10c789da3564D0"), // 217 92 112 113
		common.HexToAddress("0xC704765db1d21C2Ea6F7359dcB8FD5233DeD16b5"), // 199 4 118 93
		common.HexToAddress("0xcb7f556a77a9f7ec8dc16c7d6cf0aeff7c3ee80a"), // 203 127 85 106
		common.HexToAddress("0x0ce5cbdab931fa12821e2a71845fd284bee8914e"), // 12 229 203 218
	}
	tgn, myAddr = n[0], n[5]

	s = []common.Address{
		common.HexToAddress("0x4dd324F9821485caE941640B32c3Bcf1fA6E93E6"),
		common.HexToAddress("0x0d5Df5086B5f86f748dFaed5779c3f862C075B1f"),
		common.HexToAddress("0xD3Ff05f00491571E86A3cc8b0c320aA76D7413A5"),
		common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f"),
		common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f"),
		common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f"),
	}

	r = []common.Address{
		common.HexToAddress("0x241c793A9AD555f52f6C3a83afe6178408796ab2"),
		common.HexToAddress("0x79b427Fb77077A9716E08D049B0e8f36Abfc8E2E"),
		common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7"),
		common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09"),
		common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7"),
		common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09"),
	}

	kef = common.HexToAddress("0x136807B12327a8AfF9831F09617dA1B9D398cda2")
	kif = common.HexToAddress("0x46bA8F7538CD0749e572b2631F9FB4Ce3653AFB8")

	a0 uint64 = 0
	aL uint64 = 1000000  // less than minstaking
	aM uint64 = 2000000  // exactly minstaking (params.DefaultMinimumStake)
	a1 uint64 = 10000000 // above minstaking. Using 1,2,4,8 to uniquely spot errors

	testPUpdateInterval = uint64(36)
	testSUpdateInterval = uint64(72)
	testEpoch           = uint64(144)
	testSubGroupSize    = uint64(3)
	testProposerPolicy  = uint64(WeightedRandom)
	testGovernanceMode  = "single"

	testIstanbulCompatibleNumber = big.NewInt(int64(testEpoch) + 10)
	testKoreCompatibleBlock      = big.NewInt(int64(testEpoch) + 20)
	testRandaoCompatibleBlock    = big.NewInt(int64(testEpoch) + 30)
	testKaiaCompatibleBlock      = big.NewInt(int64(testEpoch) + 40)

	testMixHash  = common.Hex2Bytes("0x365643b31592079285c66bfbf9215bccc5d977210ec0b28fb2f19b6b92773625")
	testPrevHash = common.HexToHash("0x2e43e99fe04c8664f93f21b255356cd2279b0e993cc367db8d3e629d8c745635")
)

type testMocks struct {
	mockChain     *chainmock.MockBlockChain
	mockEngine    *mocks.MockEngine
	mockHeaderGov *hgmmock.MockHeaderGovModule
	mockStaking   *stakingmock.MockStakingModule
}

func (tm *testMocks) prepareMockExpectHeader(num uint64, mixHash, vote []byte, author common.Address) {
	header := &types.Header{
		ParentHash: testPrevHash,
		Number:     big.NewInt(int64(num)),
		MixHash:    mixHash,
		Vote:       vote,
	}
	if common.EmptyAddress(author) {
		header = nil // simulate the unmined header
	} else {
		tm.mockChain.EXPECT().GetHeaderByHash(header.Hash()).Return(header).AnyTimes()
		tm.mockEngine.EXPECT().Author(header).Return(author, nil).AnyTimes()
	}
	tm.mockChain.EXPECT().GetHeaderByNumber(num).Return(header).AnyTimes()
}

func (tm *testMocks) prepareMockExpectGovParam(num uint64, policy uint64, subSize uint64, gNode common.Address) {
	govParam := gov.GetDefaultGovernanceParamSet()
	// initialize the parameters which cannot votable later.
	_ = govParam.Set(gov.RewardProposerUpdateInterval, testPUpdateInterval)
	_ = govParam.Set(gov.RewardStakingUpdateInterval, testSUpdateInterval)
	_ = govParam.Set(gov.IstanbulEpoch, testEpoch)
	_ = govParam.Set(gov.GovernanceGovernanceMode, testGovernanceMode)
	_ = govParam.Set(gov.GovernanceGoverningNode, tgn)
	// set the parameters which can differs among test cases.
	_ = govParam.Set(gov.IstanbulCommitteeSize, subSize)
	_ = govParam.Set(gov.GovernanceGoverningNode, gNode)
	_ = govParam.Set(gov.IstanbulPolicy, policy)
	tm.mockHeaderGov.EXPECT().EffectiveParamSet(num).Return(*govParam).AnyTimes()
}

func (tm *testMocks) prepareMockExpectStakingInfo(num uint64, stakingNodeIds, stakingAmount []uint64) {
	sInfo := &staking.StakingInfo{}
	if len(stakingNodeIds) != 0 {
		var (
			nodeIds          = make([]common.Address, len(stakingNodeIds))
			stakingContracts = make([]common.Address, len(stakingNodeIds))
			rewardContracts  = make([]common.Address, len(stakingNodeIds))
		)
		for _, idx := range stakingNodeIds {
			nodeIds[idx] = n[idx]
			stakingContracts[idx] = s[idx]
			rewardContracts[idx] = r[idx]
		}
		sInfo = &staking.StakingInfo{
			SourceBlockNum:   num,
			NodeIds:          nodeIds,
			StakingContracts: stakingContracts,
			RewardAddrs:      rewardContracts,
			KEFAddr:          kef,
			KIFAddr:          kif,
			StakingAmounts:   stakingAmount,
		}
	}
	tm.mockStaking.EXPECT().GetStakingInfo(num).Return(sInfo, nil).AnyTimes()
}

func newTestVModule(ctrl *gomock.Controller) (*ValsetModule, *testMocks, error) {
	// create mocks for valSet test
	tm := &testMocks{
		mockChain:     chainmock.NewMockBlockChain(ctrl),
		mockEngine:    mocks.NewMockEngine(ctrl),
		mockStaking:   stakingmock.NewMockStakingModule(ctrl),
		mockHeaderGov: hgmmock.NewMockHeaderGovModule(ctrl),
	}
	tm.mockChain.EXPECT().Engine().Return(tm.mockEngine).AnyTimes()

	// set chainConfig
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
	tm.mockChain.EXPECT().Config().Return(config).AnyTimes()

	// init valSetModule
	vModule := NewValsetModule()
	if err := vModule.Init(&InitOpts{
		database.NewMemoryDBManager().GetMiscDB(),
		tm.mockChain,
		tm.mockHeaderGov,
		tm.mockStaking,
		myAddr,
	}); err != nil {
		return nil, nil, err
	}

	genesisHeader := &types.Header{
		Number: big.NewInt(0),
		Extra:  hexutil.MustDecode("0xd883010703846b6c617988676f312e31352e37856c696e757800000000000000f90164f85494571e53df607be97431a5bbefca1dffe5aef56f4d945cb1a7dccbd0dc446e3640898ede8820368554c89499fb17d324fa0e07f23b49d09028ac0919414db694b74ff9dea397fe9e231df545eb53fe2adf776cb2b841acb7fcc5152506250d1ea49745e7d0d5930157724b410e6e62e0885e7978c81863647d90700dcf3e5d0727cb886f2cc2c63f8f6f3910b4341b302a0aa06eae4500f8c9b841d79c07fbee8861585a71af08a867546320ba804c49c7a3c8641b4d235fd50d5a29bf72d20f3ff1ddfb945ff193d7938967be694f3e602a1cffdea686acf2b0ea01b841dfcf5b5608ca86bc92e7fa3d88a8b25840a629234614ecb312621234ed665ae562ee64ea09fcc88080aaab1ee095acf705d7cc495732682ffee23023ed41feb200b841fefc3b618b2384ea5c7c519ddecc666c19e8a600a6e30c5d9831941c0d5af78d28250bab36ce29202e667c9c1681fd9930aab002988c7228b64caab003bd998100"),
	}
	tm.mockChain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(genesisHeader))
	tm.mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesisHeader)

	// start vModule
	if err := vModule.Start(); err != nil {
		return nil, nil, err
	}
	// set initial db
	if err := WriteCouncilAddressListToDb(vModule.ChainKv, 0, n[:4]); err != nil {
		return nil, nil, err
	}

	return vModule, tm, nil
}

func TestValsetModule_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, tm, err := newTestVModule(ctrl)
	assert.NoError(t, err)

	num := uint64(0)
	tm.prepareMockExpectHeader(num, nil, nil, tgn)
	tm.prepareMockExpectGovParam(num, testProposerPolicy, testSubGroupSize, tgn)
	tm.prepareMockExpectStakingInfo(num, []uint64{0, 1, 2, 3}, []uint64{aM, aM, aM, a1})

	// check the result
	header := vModule.chain.GetHeaderByNumber(0)
	assert.Equal(t, big.NewInt(0), header.Number)
	assert.Equal(t, header, vModule.chain.GetHeaderByHash(header.Hash()))

	author, _ := vModule.chain.Engine().Author(header)
	assert.Equal(t, tgn, author)

	sInfo, _ := vModule.stakingInfo.GetStakingInfo(0)
	assert.Equal(t, []common.Address{n[0], n[1], n[2], n[3]}, sInfo.NodeIds)
	assert.Equal(t, []uint64{aM, aM, aM, a1}, sInfo.StakingAmounts)

	pSet := vModule.headerGov.EffectiveParamSet(0)
	assert.Equal(t, testProposerPolicy, pSet.ProposerPolicy)
	assert.Equal(t, testSubGroupSize, pSet.CommitteeSize)
	assert.Equal(t, tgn, pSet.GoverningNode)
}
