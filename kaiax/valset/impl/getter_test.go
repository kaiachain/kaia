package impl

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestGetDemotedAndQualified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, tm, err := newTestVModule(ctrl)
	assert.NoError(t, err)

	blockNum := testIstanbulCompatibleNumber.Uint64()
	tm.prepareMockExpectGovParam(blockNum, uint64(WeightedRandom), testSubGroupSize, tgn)
	tm.prepareMockExpectStakingInfo(blockNum, []uint64{0, 1, 2, 3, 4, 5}, []uint64{aM, aM, 0, 0, aL, aL})

	demoted, err := vModule.getDemotedValidators(blockNum)
	assert.NoError(t, err)
	assert.Equal(t, valset.AddressList{n[3], n[2]}, demoted)

	qualified, err := vModule.getQualifiedValidators(blockNum)
	assert.NoError(t, err)
	assert.Equal(t, valset.AddressList{n[0], n[1]}, qualified)
}

func TestGetCouncil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, _, err := newTestVModule(ctrl)
	assert.NoError(t, err)

	genesisValSet := make([]common.Address, len(n))
	copy(genesisValSet, n)

	// prepare vote db
	assert.NoError(t, writeCouncil(vModule.ChainKv, 0, genesisValSet[:4]))
	assert.NoError(t, writeCouncil(vModule.ChainKv, 2, genesisValSet[:5]))
	assert.NoError(t, writeCouncil(vModule.ChainKv, 4, genesisValSet[:6]))

	// check council
	for blockNumber, expectCList := range [][]common.Address{
		{tgn, n[1], n[2], n[3]}, // genesis council or committee shouldn't be modified
		{tgn, n[1], n[3], n[2]},
		{tgn, n[1], n[3], n[2], n[4]},
		{tgn, n[1], n[3], n[2], n[4]},
		{n[5], tgn, n[1], n[3], n[2], n[4]},
		{n[5], tgn, n[1], n[3], n[2], n[4]},
	} {
		cList, err := readCouncil(vModule.ChainKv, uint64(blockNumber))
		assert.NoError(t, err, "tc(blockNumber):%d", blockNumber)
		assert.Equal(t, expectCList, cList, "tc(blockNumber):%d", blockNumber)
	}
}

func TestGetCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defaultBN, defaultRound := testRandaoCompatibleBlock.Uint64()+1, uint64(0)

	for idx, tc := range []struct {
		name        string
		blockNumber uint64
		round       uint64

		govParamPolicy  uint64
		govParamSubSize uint64
		govParamGovNode common.Address

		expectCommitteeList valset.AddressList
		expectError         error
	}{
		// per committeeSize
		{"committeesize is zero", defaultBN, defaultRound, testProposerPolicy, 0, tgn, nil, errInvalidCommitteeSize},
		{"committeesize is one", defaultBN, defaultRound, testProposerPolicy, 1, tgn, []common.Address{tgn}, nil},
		{"committeesize is three", defaultBN, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{n[0], n[1], n[2]}, nil},
		{"committeesize is six", defaultBN, defaultRound, testProposerPolicy, testSubGroupSize + 3, tgn, []common.Address{tgn, n[1], n[3], n[2]}, nil},
		{"committeesize is seven", defaultBN, defaultRound, testProposerPolicy, testSubGroupSize + 4, tgn, []common.Address{tgn, n[1], n[3], n[2]}, nil},
		// per proposerPolicy
		{"proposerPolicy roundrobin", defaultBN, defaultRound, params.RoundRobin, testSubGroupSize, tgn, []common.Address{n[1], n[3], tgn}, nil},
		{"proposerPolicy sticky", defaultBN, defaultRound, params.Sticky, testSubGroupSize, tgn, []common.Address{tgn, n[1], n[3]}, nil},
		// per HF
		{"genesis block", 0, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{tgn, n[1], n[2], n[3]}, nil},
		{"block 1", 1, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{n[3], n[1], tgn}, nil},
		{"istanbul hf activated", testIstanbulCompatibleNumber.Uint64() + 1, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{n[1], n[2], n[3]}, nil},
		{"kore hf activated", testKoreCompatibleBlock.Uint64() + 1, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{n[3], tgn, n[1]}, nil},
		{"randao hf activated", testRandaoCompatibleBlock.Uint64() + 1, defaultRound, testProposerPolicy, testSubGroupSize, tgn, []common.Address{tgn, n[1], n[2]}, nil},
		// TODO-kaia-valset: add mainnet,testnet testcases
	} {
		t.Run(tc.name, func(t *testing.T) {
			vModule, tm, err := newTestVModule(ctrl)
			assert.NoError(t, err)

			mixHash := testMixHash
			if tc.blockNumber < testRandaoCompatibleBlock.Uint64() {
				mixHash = nil
			}

			if tc.blockNumber > 0 {
				tm.mockChain.EXPECT().GetHeaderByNumber(tc.blockNumber).Return(nil).AnyTimes()
				tm.prepareMockExpectHeader(tc.blockNumber-1, mixHash, nil, tgn)
				tm.prepareMockExpectStakingInfo(tc.blockNumber, nil, nil)
				tm.prepareMockExpectGovParam(tc.blockNumber, tc.govParamPolicy, tc.govParamSubSize, tc.govParamGovNode)

				proposersPSet := vModule.governance.EffectiveParamSet(tc.blockNumber)
				proposersBlock := calcProposerBlockNumber(tc.blockNumber, proposersPSet.ProposerUpdateInterval)
				tm.prepareMockExpectHeader(proposersBlock-1, mixHash, nil, tgn)
				tm.prepareMockExpectStakingInfo(proposersBlock, nil, nil)
				tm.prepareMockExpectGovParam(proposersBlock, tc.govParamPolicy, tc.govParamSubSize, tc.govParamGovNode)
			}

			committee, err := vModule.GetCommittee(tc.blockNumber, tc.round)
			assert.Equal(t, tc.expectError, err, "testcase: %d", idx)
			assert.Equal(t, tc.expectCommitteeList, committee, "testcase: %d", idx)
		})
	}
}

func TestGetProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range []struct {
		name        string
		blockNumber uint64
		round       uint64

		govParamPolicy  uint64
		prevBlockAuthor common.Address

		expectProposer common.Address
		expectError    error
	}{
		{"sticky round 0", 3, 0, uint64(Sticky), n[0], n[0], nil},
		{"sticky round 1", 3, 1, uint64(Sticky), n[0], n[1], nil},
	} {
		vModule, tm, err := newTestVModule(ctrl)
		assert.NoError(t, err)

		tm.prepareMockExpectHeader(tc.blockNumber-1, nil, nil, n[0])
		tm.prepareMockExpectStakingInfo(tc.blockNumber, []uint64{0, 1, 2, 3, 4, 5}, []uint64{aM, aM, aM, aM, aM, aM})
		tm.prepareMockExpectGovParam(tc.blockNumber, tc.govParamPolicy, testSubGroupSize, tgn)

		tm.prepareMockExpectHeader(tc.blockNumber, nil, nil, common.Address{})

		proposer, err := vModule.GetProposer(tc.blockNumber, tc.round)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectProposer, proposer)
	}
}
