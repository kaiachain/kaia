package impl

import (
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
)

func addressListToString(v []common.Address) string {
	res := ""
	for _, item := range v {
		res = res + item.String() + ","
	}
	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res
}

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

func TestValsetModule_HandleValidatorVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, tm, err := newTestVModule(ctrl)
	assert.NoError(t, err)

	tm.prepareMockExpectGovParam(1, testProposerPolicy, testSubGroupSize, tgn)

	for _, tc := range []struct {
		name                      string
		networkInfo               testNetworkInfo
		voteData                  headergov.VoteData
		expectHandleValidatorVote testHandleValidatorVoteExpectation
	}{
		{
			name:                      "remove a validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", addressListToString([]common.Address{n[1]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[2], n[3]}, nil},
		},
		{
			name:                      "add a validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", addressListToString([]common.Address{n[4]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3], n[4]}, nil},
		},
		{
			name:                      "remove multiple validators",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", addressListToString([]common.Address{n[1], n[2], n[3]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn}, nil},
		},
		{
			name:                      "add multiple validators",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", addressListToString([]common.Address{n[2], n[3]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, nil},
		},
		{
			name:                      "govnode cannot be removed",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", addressListToString([]common.Address{tgn})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, errInvalidVoteValue},
		},
		{
			name:                      "cannot add existing validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.addvalidator", addressListToString([]common.Address{n[1]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, errInvalidVoteValue},
		},
		{
			name:                      "cannot remove non-exist validator",
			networkInfo:               testNetworkInfo{[]common.Address{tgn, n[1], n[2], n[3]}, tgn},
			voteData:                  headergov.NewVoteData(tgn, "governance.removevalidator", addressListToString([]common.Address{n[4]})),
			expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n[1], n[2], n[3]}, errInvalidVoteValue},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err = WriteCouncilAddressListToDb(vModule.ChainKv, 0, tc.networkInfo.validators)
			assert.NoError(t, err)

			byteVote, err := tc.voteData.ToVoteBytes()
			assert.NoError(t, err)

			err = vModule.HandleValidatorVote(1, byteVote, tc.networkInfo.validators)
			assert.Equal(t, tc.expectHandleValidatorVote.expectError, err)

			cList, err := readCouncilAddressListFromValSetCouncilDB(vModule.ChainKv, uint64(1))
			sort.Sort(tc.expectHandleValidatorVote.expectCList)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectHandleValidatorVote.expectCList, valset.AddressList(cList))
		})
	}
}
