package impl

import (
	"math/big"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/mocks"
	hgmmock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	stakingmock "github.com/kaiachain/kaia/kaiax/staking/mock"
	chainmock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestValsetModule_HandleValidatorVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		mockChain, mockEngine      = chainmock.NewMockBlockChain(ctrl), mocks.NewMockEngine(ctrl)
		mockStaking, mockHeaderGov = stakingmock.NewMockStakingModule(ctrl), hgmmock.NewMockHeaderGovModule(ctrl)
	)
	vModule, _, _, testParamSets, err := newTestVModule(mockChain, mockEngine, mockHeaderGov, mockStaking)
	mockHeaderGov.EXPECT().EffectiveParamSet(gomock.Any()).Return(testParamSets[0]).AnyTimes()

	assert.NoError(t, err)

	for blockNumber, tc := range testVotes {
		vData := &voteData{voter: tc.testVote.voter, name: tc.testVote.voteName, value: tc.testVote.validators}

		byteVote, err := vData.ToVoteBytes()
		assert.NoError(t, err)

		header := &types.Header{
			ParentHash: testPrevHash,
			Number:     big.NewInt(int64(blockNumber)),
			Vote:       byteVote,
		}
		mockEngine.EXPECT().Author(header).Return(tc.testVote.proposer, nil).AnyTimes()
		err = vModule.HandleValidatorVote(header, vData)
		assert.Equal(t, tc.expectHandleValidatorVote.expectError, err, "blockNumber: %d", blockNumber)

		cList, err := ReadCouncilAddressListFromDb(vModule.ChainKv, uint64(blockNumber))
		sort.Sort(tc.expectHandleValidatorVote.expectCList)
		assert.NoError(t, err, "blockNumber: %d", blockNumber)
		assert.Equal(t, tc.expectHandleValidatorVote.expectCList, subsetCouncilSlice(cList), "blockNumber: %d", blockNumber)
	}

	assert.Equal(t, []uint64{0, 1, 2, 3}, ReadValidatorVoteDataBlockNums(vModule.ChainKv))
}
