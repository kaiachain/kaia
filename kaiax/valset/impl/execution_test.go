package impl

import (
	"math/big"
	"sort"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/mocks"
	hgmmock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	stakingmock "github.com/kaiachain/kaia/kaiax/staking/mock"
	chainmock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

// voteTestData is the testdata for vote tests.
// Not only votes, but also other functions such as VerifyHeader handles votes.
// To minimize the maintenance effort to manage the testdata, manage the expected result in voteTestData
type testNetworkInfo struct {
	validators []common.Address
	proposer   common.Address
}

type testHandleValidatorVoteExpectation struct {
	expectCList valset.AddressList
	expectError error
}

var testVotes = []struct {
	name                      string
	networkInfo               testNetworkInfo
	voteData                  headergov.VoteData
	expectHandleValidatorVote testHandleValidatorVoteExpectation
}{
	{
		name:                      "remove a validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", []common.Address{n2}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n3, n4}, nil},
	},
	{
		name:                      "add a validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", []common.Address{n5}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4, n5}, nil},
	},
	{
		name:                      "remove multiple validators",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", []common.Address{n2, n3, n4}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn}, nil},
	},
	{
		name:                      "add multiple validators",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", []common.Address{n3, n4}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, nil},
	},
	{
		name:                      "govnode cannot be removed",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", []common.Address{tgn}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
	{
		name:                      "proposer should be same with voter and it should be a gov node",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, n2},
		voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", []common.Address{n4}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoter},
	},
	{
		name:                      "voter and proposer should be a gov node",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, n2},
		voteData:                  headergov.NewVoteData(n2, "governance.removevalidator", []common.Address{n4}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoter},
	},
	{
		name:                      "cannot add existing validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", []common.Address{n2}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
	{
		name:                      "cannot remove non-exist validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", []common.Address{n5}),
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
}

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

			byteVote, err := tc.voteData.ToVoteBytes()
			assert.NoError(t, err)

			prevHeader := &types.Header{ParentHash: testPrevHash, Number: big.NewInt(0)}
			header := &types.Header{
				ParentHash: testPrevHash,
				Number:     big.NewInt(1),
				Vote:       byteVote,
			}
			mockEngine.EXPECT().Author(header).Return(tc.networkInfo.proposer, nil).AnyTimes()
			mockChain.EXPECT().GetHeaderByNumber(prevHeader.Number.Uint64()).Return(prevHeader).AnyTimes()
			err = vModule.HandleValidatorVote(header, tc.voteData)
			assert.Equal(t, tc.expectHandleValidatorVote.expectError, err)

			cList, err := ReadCouncilAddressListFromDb(vModule.ChainKv, uint64(1))
			sort.Sort(tc.expectHandleValidatorVote.expectCList)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectHandleValidatorVote.expectCList, valset.AddressList(cList))
		})
	}
}
