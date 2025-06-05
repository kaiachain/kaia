package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	var (
		vote               = headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100))
		voteBytes, _       = headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100)).ToVoteBytes()
		govBytes, _        = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(100)}).ToGovBytes()
		invalidGovBytes, _ = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(200)}).ToGovBytes()
		h                  = newHeaderGovModule(t, getTestChainConfig())
		invalidVoteRlp     = common.FromHex("0xea9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e74")
	)

	h.HandleVote(1, vote)

	tcs := []struct {
		desc          string
		header        *types.Header
		expectedError error
	}{
		{desc: "valid vote", header: &types.Header{Number: big.NewInt(1), Vote: voteBytes, Extra: extra}, expectedError: nil},
		{desc: "invalid vote rlp", header: &types.Header{Number: big.NewInt(1), Vote: []byte{1, 2, 3}, Extra: extra}, expectedError: headergov.ErrInvalidRlp},
		{desc: "invalid vote bytes", header: &types.Header{Number: big.NewInt(1), Vote: invalidVoteRlp, Extra: extra}, expectedError: headergov.ErrInvalidRlp},
		{desc: "valid gov", header: &types.Header{Number: big.NewInt(100), Governance: govBytes, Extra: extra}, expectedError: nil},
		{desc: "invalid gov rlp", header: &types.Header{Number: big.NewInt(100), Governance: []byte{1, 2, 3}, Extra: extra}, expectedError: headergov.ErrInvalidRlp},
		{desc: "gov must not be nil", header: &types.Header{Number: big.NewInt(100), Governance: nil}, expectedError: ErrGovVerification},
		{desc: "gov mismatch", header: &types.Header{Number: big.NewInt(100), Governance: invalidGovBytes}, expectedError: ErrGovVerification},
		{desc: "gov not on epoch", header: &types.Header{Number: big.NewInt(101), Governance: []byte{1, 2, 3}}, expectedError: ErrGovInNonEpochBlock},
		{desc: "gov must be nil", header: &types.Header{Number: big.NewInt(200), Governance: govBytes}, expectedError: ErrGovVerification},
		{desc: "valid gov", header: &types.Header{Number: big.NewInt(200), Governance: nil}, expectedError: nil},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			tc.header.Extra = extra
			err := h.VerifyHeader(tc.header)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestVerifyVote(t *testing.T) {
	config := getTestChainConfig()
	config.Governance.GoverningNode = validVoter // for "valid governingnode" test case

	var (
		h          = newHeaderGovModule(t, config)
		statedb, _ = h.Chain.State()

		eoa      = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
		contract = common.HexToAddress("0x0000000000000000000000000000000000000400")
	)

	statedb.SetNonce(eoa, 1)
	statedb.SetCode(contract, []byte{1})

	tcs := []struct {
		desc          string
		vote          headergov.VoteData
		expectedError error
	}{
		// governance.*
		{desc: "valid deriveshaimpl", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceDeriveShaImpl), uint64(1)), expectedError: nil},
		// GovernanceGovernanceMode vote is forbidden
		{desc: "valid governingnode", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceGoverningNode), validVoter.Hex()), expectedError: nil},
		{desc: "valid govparam", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceGovParamContract), contract), expectedError: nil},
		{desc: "valid unitprice", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceUnitPrice), uint64(25000000000)), expectedError: nil},
		{desc: "invalid govparam", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceGovParamContract), common.Address{}), expectedError: ErrGovParamNotAccount},
		{desc: "invalid govparam", vote: headergov.NewVoteData(validVoter, string(gov.GovernanceGovParamContract), eoa), expectedError: ErrGovParamNotContract},

		// istanbul.*
		{desc: "valid committeesize", vote: headergov.NewVoteData(validVoter, string(gov.IstanbulCommitteeSize), uint64(7)), expectedError: nil},
		// IstanbulEpoch vote is forbidden
		// IstanbulPolicy vote is forbidden

		// kip71.*
		{desc: "valid basefeedenominator", vote: headergov.NewVoteData(validVoter, string(gov.Kip71BaseFeeDenominator), uint64(8)), expectedError: nil},
		{desc: "valid gastarget", vote: headergov.NewVoteData(validVoter, string(gov.Kip71GasTarget), uint64(30000000)), expectedError: nil},
		{desc: "valid lowerboundbasefee", vote: headergov.NewVoteData(validVoter, string(gov.Kip71LowerBoundBaseFee), uint64(25000000000)), expectedError: nil},
		{desc: "valid maxblockgasusedforbasefee", vote: headergov.NewVoteData(validVoter, string(gov.Kip71MaxBlockGasUsedForBaseFee), uint64(60000000)), expectedError: nil},
		{desc: "valid upperboundbasefee", vote: headergov.NewVoteData(validVoter, string(gov.Kip71UpperBoundBaseFee), uint64(750000000000)), expectedError: nil},
		{desc: "invalid lower", vote: headergov.NewVoteData(validVoter, string(gov.Kip71LowerBoundBaseFee), uint64(1e18)), expectedError: ErrLowerBoundBaseFee},
		{desc: "invalid upper", vote: headergov.NewVoteData(validVoter, string(gov.Kip71UpperBoundBaseFee), uint64(1)), expectedError: ErrUpperBoundBaseFee},

		// reward.*
		// RewardDeferredTxFee vote is forbidden
		{desc: "valid kip82ratio", vote: headergov.NewVoteData(validVoter, string(gov.RewardKip82Ratio), "20/80"), expectedError: nil},
		{desc: "valid mintingamount", vote: headergov.NewVoteData(validVoter, string(gov.RewardMintingAmount), "6400000000000000000"), expectedError: nil},
		// RewardMinimumStake vote is forbidden
		// RewardProposerUpdateInterval vote is forbidden
		{desc: "valid ratio", vote: headergov.NewVoteData(validVoter, string(gov.RewardRatio), "50/30/20"), expectedError: nil},
		// RewardStakingUpdateInterval vote is forbidden
		// RewardUseGiniCoeff vote is forbidden
		{desc: "valid addvalidator", vote: headergov.NewVoteData(validVoter, string(gov.AddValidator), eoa), expectedError: nil},
		{desc: "valid removevalidator", vote: headergov.NewVoteData(validVoter, string(gov.RemoveValidator), eoa), expectedError: nil},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			vb, err := tc.vote.ToVoteBytes()
			assert.NoError(t, err)
			err = h.VerifyVote(&types.Header{Number: big.NewInt(1), Vote: vb, Extra: extra})
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetVotesInEpoch(t *testing.T) {
	h := newHeaderGovModule(t, getTestChainConfig())

	paramName := string(gov.GovernanceUnitPrice)
	v1 := headergov.NewVoteData(common.Address{1}, paramName, uint64(100))
	h.HandleVote(50, v1)
	v2 := headergov.NewVoteData(common.Address{2}, paramName, uint64(200))
	h.HandleVote(150, v2)

	assert.Equal(t, map[uint64]headergov.VoteData{50: v1}, h.getVotesInEpoch(0))
	assert.Equal(t, map[uint64]headergov.VoteData{150: v2}, h.getVotesInEpoch(1))
}

func TestGetExpectedGovernance(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	var (
		paramName = string(gov.GovernanceUnitPrice)
		h         = newHeaderGovModule(t, getTestChainConfig())
		v1        = headergov.NewVoteData(common.Address{1}, paramName, uint64(100))
		v2        = headergov.NewVoteData(common.Address{2}, paramName, uint64(200))
		g         = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(200)})
	)

	// v2 overrides v1
	h.HandleVote(50, v1)
	h.HandleVote(60, v2)

	// test many times for deterministic result
	for range 100 {
		assert.Equal(t, g, h.getExpectedGovernance(100))
	}
}

func TestPrepareHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	var (
		h      = newHeaderGovModule(t, getTestChainConfig())
		vote   = headergov.NewVoteData(h.nodeAddress, string(gov.GovernanceUnitPrice), uint64(100))
		header = &types.Header{}
	)

	h.PushMyVotes(vote)
	header.Number = big.NewInt(999)
	err := h.PrepareHeader(header)
	assert.Nil(t, err)
	assert.NotNil(t, header.Vote)
	assert.Nil(t, header.Governance)

	h.PostInsertBlock(types.NewBlockWithHeader(header))
	header = &types.Header{}
	header.Number = big.NewInt(1000)
	err = h.PrepareHeader(header)
	assert.Nil(t, err)
	assert.Nil(t, header.Vote)
	assert.NotNil(t, header.Governance)

	h.PostInsertBlock(types.NewBlockWithHeader(header))
	ps := h.GetParamSet(2001)
	assert.Equal(t, ps.UnitPrice, uint64(100))
}
