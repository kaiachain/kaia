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

	for _, tc := range testVotes {
		t.Run(tc.name, func(t *testing.T) {
			err = WriteCouncilAddressListToDb(vModule.ChainKv, 0, tc.networkInfo.validators)
			assert.NoError(t, err)

			vData := &tc.voteData

			byteVote, err := vData.ToVoteBytes()
			assert.NoError(t, err)

			prevHeader := &types.Header{ParentHash: testPrevHash, Number: big.NewInt(0)}
			header := &types.Header{
				ParentHash: testPrevHash,
				Number:     big.NewInt(1),
				Vote:       byteVote,
			}
			mockEngine.EXPECT().Author(header).Return(tc.networkInfo.proposer, nil).AnyTimes()
			mockChain.EXPECT().GetHeaderByNumber(prevHeader.Number.Uint64()).Return(prevHeader).AnyTimes()
			err = vModule.HandleValidatorVote(header, vData)
			assert.Equal(t, tc.expectHandleValidatorVote.expectError, err)

			cList, err := ReadCouncilAddressListFromDb(vModule.ChainKv, uint64(1))
			sort.Sort(tc.expectHandleValidatorVote.expectCList)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectHandleValidatorVote.expectCList, subsetCouncilSlice(cList))
		})
	}
}
