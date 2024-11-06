package impl

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

// voteTestData is the testdata for vote tests.
// Not only votes, but also other functions such as VerifyHeader handles votes.
// To minimize the maintenance effort to manage the testdata, manage the expected result in voteTestData
type testVoteData struct {
	voter      common.Address
	proposer   common.Address
	voteName   string
	validators []common.Address
}

type testHandleValidatorVoteExpectation struct {
	expectCList subsetCouncilSlice
	expectError error
}

type testVerifyHeaderExpectation struct {
	expectError error
}

var testVotes = []struct {
	name                      string
	testVote                  testVoteData
	expectHandleValidatorVote testHandleValidatorVoteExpectation
	expectVerifyHeader        testVerifyHeaderExpectation
}{
	{
		name:                      "remove a validator",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.removevalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n3, n4}, nil},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "add a validator",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.addvalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n4}, nil},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "remove multiple validators",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.removevalidator", []common.Address{n2, n3, n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode}, nil},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "add multiple validators",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.addvalidator", []common.Address{n2, n3, n5}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, nil},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "govnode cannot be removed",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.removevalidator", []common.Address{testGoverningNode}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoteValue},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "proposer should be a gov node",
		testVote:                  testVoteData{testGoverningNode, n2, "governance.addvalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoter},
		expectVerifyHeader:        testVerifyHeaderExpectation{errInvalidVoter},
	},
	{
		name:                      "voter should be a gov node",
		testVote:                  testVoteData{n2, testGoverningNode, "governance.removevalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoter},
		expectVerifyHeader:        testVerifyHeaderExpectation{errInvalidVoter},
	},
	{
		name:                      "proposer should be a gov node",
		testVote:                  testVoteData{n2, n2, "governance.removevalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoter},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "cannot add existing validator",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.addvalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoteValue},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "cannot remove non-exist validator",
		testVote:                  testVoteData{testGoverningNode, testGoverningNode, "governance.removevalidator", []common.Address{n4}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoteValue},
		expectVerifyHeader:        testVerifyHeaderExpectation{},
	},
	{
		name:                      "voter is not in council list",
		testVote:                  testVoteData{n6, n6, "governance.removevalidator", []common.Address{n2}},
		expectHandleValidatorVote: testHandleValidatorVoteExpectation{[]common.Address{testGoverningNode, n2, n3, n5}, errInvalidVoter},
		expectVerifyHeader:        testVerifyHeaderExpectation{errInvalidVoter},
	},
}

func TestVoteData_ToVoteBytes(t *testing.T) {
	for idx, tc := range []struct {
		voteData
		voteBytes   string
		expectError error
	}{
		{voteData{testGoverningNode, "governance.removevalidator", []common.Address{n2}}, "f846948ad8f547fa00f58a8c4fb3b671ee5f1a75ba028a9a676f7665726e616e63652e72656d6f766576616c696461746f72d594b2aada7943919e82143324296987f6091f3fdc9e", nil},
	} {
		voteBytes, err := tc.voteData.ToVoteBytes()
		assert.Equal(t, tc.expectError, err, "test case: %d", idx)
		assert.Equal(t, common.Hex2Bytes(tc.voteBytes), voteBytes, "test case: %d", idx)
	}
}
