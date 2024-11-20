package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/mocks"
	hgmmock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	stakingmock "github.com/kaiachain/kaia/kaiax/staking/mock"
	chainmock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCouncilAddressList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		mockChain, mockEngine      = chainmock.NewMockBlockChain(ctrl), mocks.NewMockEngine(ctrl)
		mockStaking, mockHeaderGov = stakingmock.NewMockStakingModule(ctrl), hgmmock.NewMockHeaderGovModule(ctrl)
	)
	vModule, _, _, _, err := newTestVModule(mockChain, mockEngine, mockHeaderGov, mockStaking)
	assert.NoError(t, err)

	genesisValSet := make([]common.Address, len(testGenesisValSet))
	copy(genesisValSet, testGenesisValSet)

	// prepare vote db
	assert.NoError(t, WriteCouncilAddressListToDb(vModule.ChainKv, 0, genesisValSet[:4]))
	assert.NoError(t, WriteCouncilAddressListToDb(vModule.ChainKv, 2, genesisValSet[:5]))
	assert.NoError(t, WriteCouncilAddressListToDb(vModule.ChainKv, 4, genesisValSet[:6]))

	// check council
	for blockNumber, expectCList := range [][]common.Address{
		{tgn, n2, n4, n3},
		{tgn, n2, n4, n3},
		{tgn, n2, n4, n3, n5},
		{tgn, n2, n4, n3, n5},
		{n6, tgn, n2, n4, n3, n5},
		{n6, tgn, n2, n4, n3, n5},
	} {
		cList, err := readCouncilAddressListFromValSetCouncilDB(vModule.ChainKv, uint64(blockNumber))
		assert.NoError(t, err, "tc(blockNumber):%d", blockNumber)
		assert.Equal(t, expectCList, cList)
	}
}

func TestGetCommitteeAddressList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defaultBlockNumber, defaultRound := testRandaoCompatibleBlock.Uint64()+1, uint64(0)

	for idx, tc := range []struct {
		name            string
		blockNumber     uint64
		round           uint64
		testStakingInfo int
		testParams      int
		testAuthor      common.Address

		expectCommitteeList []common.Address
		expectError         error
	}{
		// per committeeSize
		{"committeesize is zero", defaultBlockNumber, defaultRound, 0, 0, tgn, nil, errInvalidCommitteeSize},
		{"committeesize is one", defaultBlockNumber, defaultRound, 0, 1, tgn, []common.Address{n3}, nil},
		{"committeesize is three", defaultBlockNumber, defaultRound, 0, 2, tgn, []common.Address{n3, tgn, n6}, nil},
		{"committeesize is six", defaultBlockNumber, defaultRound, 0, 3, tgn, []common.Address{n6, tgn, n2, n4, n3, n5}, nil},
		{"committeesize is seven", defaultBlockNumber, defaultRound, 0, 4, tgn, []common.Address{n6, tgn, n2, n4, n3, n5}, nil},
		// per proposerPolicy
		{"proposerPolicy roundrobin", defaultBlockNumber, defaultRound, 0, 5, tgn, []common.Address{n3, tgn, n6}, nil},
		{"proposerPolicy sticky", defaultBlockNumber, defaultRound, 0, 6, tgn, []common.Address{n3, tgn, n6}, nil},
		// per HF
		{"genesis block", 0, defaultRound, 0, 2, tgn, []common.Address{n6, tgn, n2}, nil},
		{"block 1", 1, defaultRound, 0, 2, tgn, []common.Address{n6, n5, tgn}, nil},
		{"istanbul hf activated", testIstanbulCompatibleNumber.Uint64() + 1, defaultRound, 0, 2, tgn, []common.Address{n6, tgn, n4}, nil},
		{"kore hf activated", testKoreCompatibleBlock.Uint64() + 1, defaultRound, 0, 2, tgn, []common.Address{n2, n5, n6}, nil},
		{"randao hf activated", testRandaoCompatibleBlock.Uint64() + 1, defaultRound, 0, 2, tgn, []common.Address{n3, tgn, n6}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				mockChain, mockEngine      = chainmock.NewMockBlockChain(ctrl), mocks.NewMockEngine(ctrl)
				mockStaking, mockHeaderGov = stakingmock.NewMockStakingModule(ctrl), hgmmock.NewMockHeaderGovModule(ctrl)
			)
			vModule, _, testStaking, testParamSets, err := newTestVModule(mockChain, mockEngine, mockHeaderGov, mockStaking)
			assert.NoError(t, err)

			// custom setting for mocks
			mockEngine.EXPECT().Author(gomock.Any()).Return(tgn, nil).AnyTimes()
			// proposerBlock set

			// testBlockSet
			mockHeaderGov.EXPECT().EffectiveParamSet(gomock.Any()).Return(testParamSets[tc.testParams]).AnyTimes()
			mockStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(testStaking[tc.testStakingInfo], nil).AnyTimes()
			mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(&types.Header{Number: big.NewInt(int64(tc.blockNumber))}).AnyTimes()

			committee, err := vModule.GetCommitteeAddressList(tc.blockNumber, tc.round)
			assert.Equal(t, tc.expectError, err, "testcase: %d", idx)
			assert.Equal(t, tc.expectCommitteeList, committee, "testcase: %d", idx)
		})
	}
}
