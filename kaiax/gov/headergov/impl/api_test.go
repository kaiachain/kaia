package impl

import (
	"errors"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHeaderGovAPI(t *testing.T) *headerGovAPI {
	chainConfig := getTestChainConfig()
	h := newHeaderGovModule(t, chainConfig)
	return NewHeaderGovAPI(h)
}

func TestVoteDeprecated(t *testing.T) {
	api := newHeaderGovAPI(t)
	for name := range gov.AlwaysDeprecated {
		_, err := api.Vote(string(name), gov.Params[name].DefaultValue)
		assert.Equal(t, ErrDeprecatedVote, err, "expected ErrDeprecatedVote for %s", name)
	}
}

func TestVoteUpperBoundBaseFee(t *testing.T) {
	api := newHeaderGovAPI(t)
	s, err := api.Vote("kip71.upperboundbasefee", uint64(1))
	assert.Equal(t, ErrUpperBoundBaseFee, err)
	assert.Equal(t, "", s)
}

func TestVoteLowerBoundBaseFee(t *testing.T) {
	api := newHeaderGovAPI(t)
	s, err := api.Vote("kip71.lowerboundbasefee", uint64(1e18))
	assert.Equal(t, ErrLowerBoundBaseFee, err)
	assert.Equal(t, "", s)
}

func TestValidatorVotes(t *testing.T) {
	testCases := []struct {
		key         string
		value       string
		expectedErr error
	}{
		{key: "governance.addvalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{key: "governance.addvalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{key: "governance.addvalidator", value: ",0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", expectedErr: ErrInvalidKeyValue},
		{key: "governance.removevalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{key: "governance.removevalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{key: "governance.removevalidator", value: ",0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", expectedErr: ErrInvalidKeyValue},
	}

	api := newHeaderGovAPI(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_, err := api.Vote(tc.key, tc.value)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestVoteGovParamContract(t *testing.T) {
	var (
		api        = newHeaderGovAPI(t)
		statedb, _ = api.h.Chain.State()

		eoa           = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
		contract      = common.HexToAddress("0x0000000000000000000000000000000000000400")
		emptyContract = common.HexToAddress("0x0000000000000000000000000000000000000401")
		nonexistent   = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	)

	statedb.SetNonce(eoa, 1)
	statedb.SetCode(contract, []byte{1})
	// A smart contract account with zero nonce, zero balance and no code is Empty().
	statedb.SetCode(emptyContract, nil)
	statedb.SetNonce(emptyContract, 0)

	testCases := []struct {
		desc        string
		value       any
		expectedErr error
	}{
		{desc: "contract", value: contract.Hex(), expectedErr: nil},
		{desc: "eoa", value: eoa.Hex(), expectedErr: ErrGovParamNotContract},
		{desc: "empty contract", value: emptyContract.Hex(), expectedErr: ErrGovParamNotContract},
		{desc: "zero address", value: common.Address{}.Hex(), expectedErr: ErrGovParamNotAccount},
		{desc: "nonexistent", value: nonexistent.Hex(), expectedErr: ErrGovParamNotAccount},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			before := len(api.h.myVotesSnapshot())
			s, err := api.Vote(string(gov.GovernanceGovParamContract), tc.value)
			assert.Equal(t, tc.expectedErr, err)

			after := api.h.myVotesSnapshot()
			if tc.expectedErr == nil {
				assert.NotEqual(t, "", s)
				require.Len(t, after, before+1)
				assert.Equal(t, gov.GovernanceGovParamContract, after[len(after)-1].Name())
				assert.Equal(t, common.HexToAddress(tc.value.(string)), after[len(after)-1].Value())
			} else {
				assert.Equal(t, "", s)
				assert.Len(t, after, before)
			}
		})
	}
}

func TestVoteGovParamContractStateError(t *testing.T) {
	var (
		api      = newHeaderGovAPI(t)
		chain    = mocks.NewMockBlockChain(gomock.NewController(t))
		stateErr = errors.New("state unavailable")
		contract = common.HexToAddress("0x0000000000000000000000000000000000000400")
	)

	chain.EXPECT().CurrentBlock().Return(api.h.Chain.CurrentBlock()).AnyTimes()
	chain.EXPECT().State().Return(nil, stateErr).AnyTimes()
	api.h.Chain = chain

	before := len(api.h.myVotesSnapshot())
	s, err := api.Vote(string(gov.GovernanceGovParamContract), contract.Hex())
	assert.Equal(t, stateErr, err)
	assert.Equal(t, "", s)
	assert.Len(t, api.h.myVotesSnapshot(), before)
}
