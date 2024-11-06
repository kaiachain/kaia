package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/mocks"
	hgmmock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	stakingmock "github.com/kaiachain/kaia/kaiax/staking/mock"
	chainmock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestValsetModule_HandleVerifyVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		mockChain, mockEngine      = chainmock.NewMockBlockChain(ctrl), mocks.NewMockEngine(ctrl)
		mockStaking, mockHeaderGov = stakingmock.NewMockStakingModule(ctrl), hgmmock.NewMockHeaderGovModule(ctrl)
	)
	vModule, _, _, _, err := newTestVModule(mockChain, mockEngine, mockHeaderGov, mockStaking)
	assert.NoError(t, err)

	for idx, tc := range testVotes {
		vData := &voteData{tc.testVote.voter, tc.testVote.voteName, tc.testVote.validators}

		byteVote, err := vData.ToVoteBytes()
		assert.NoError(t, err)

		header := &types.Header{
			ParentHash: testPrevHash,
			Number:     big.NewInt(int64(idx + 1)),
			Vote:       byteVote,
		}
		mockEngine.EXPECT().Author(header).Return(tc.testVote.proposer, nil).AnyTimes()
		err = vModule.VerifyHeader(header)
		assert.Equal(t, tc.expectVerifyHeader.expectError, err, "tcnum: %d", idx)
	}
}
