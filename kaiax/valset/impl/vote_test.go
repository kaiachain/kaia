package impl

import (
	"testing"

	"github.com/kaiachain/kaia/common"
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
	expectCList subsetCouncilSlice
	expectError error
}

var testVotes = []struct {
	name                      string
	networkInfo               testNetworkInfo
	voteData                  voteData
	expectHandleValidatorVote testHandleValidatorVoteExpectation
}{
	{
		name:                      "remove a validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.removevalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n3, n4}, nil},
	},
	{
		name:                      "add a validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.addvalidator", []common.Address{n5}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4, n5}, nil},
	},
	{
		name:                      "remove multiple validators",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.removevalidator", []common.Address{n2, n3, n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn}, nil},
	},
	{
		name:                      "add multiple validators",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2}, tgn},
		voteData:                  voteData{tgn, "governance.addvalidator", []common.Address{n3, n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, nil},
	},
	{
		name:                      "govnode cannot be removed",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.removevalidator", []common.Address{tgn}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
	{
		name:                      "proposer should be same with voter and it should be a gov node",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, n2},
		voteData:                  voteData{tgn, "governance.addvalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoter},
	},
	{
		name:                      "voter and proposer should be a gov node",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, n2},
		voteData:                  voteData{n2, "governance.removevalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoter},
	},
	{
		name:                      "cannot add existing validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.addvalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
	{
		name:                      "cannot remove non-exist validator",
		networkInfo:               testNetworkInfo{[]common.Address{tgn, n2, n3, n4}, tgn},
		voteData:                  voteData{tgn, "governance.removevalidator", []common.Address{n5}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{tgn, n2, n3, n4}, errInvalidVoteValue},
	},
}

func TestVoteData_ToVoteBytes(t *testing.T) {
	for idx, tc := range []struct {
		voteData
		voteBytes   string
		expectError error
	}{
		{voteData{tgn, "governance.removevalidator", []common.Address{n2}}, "f846948ad8f547fa00f58a8c4fb3b671ee5f1a75ba028a9a676f7665726e616e63652e72656d6f766576616c696461746f72d594b2aada7943919e82143324296987f6091f3fdc9e", nil},
	} {
		voteBytes, err := tc.voteData.ToVoteBytes()
		assert.Equal(t, tc.expectError, err, "test case: %d", idx)
		assert.Equal(t, common.Hex2Bytes(tc.voteBytes), voteBytes, "test case: %d", idx)
	}
}
