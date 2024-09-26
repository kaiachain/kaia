package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	var (
		vote               = headergov.NewVoteData(common.Address{1}, gov.Params[gov.GovernanceUnitPrice].Name, uint64(100))
		voteBytes, _       = headergov.NewVoteData(common.Address{1}, gov.Params[gov.GovernanceUnitPrice].Name, uint64(100)).Serialize()
		govBytes, _        = headergov.NewGovData(map[gov.ParamEnum]any{gov.GovernanceUnitPrice: uint64(100)}).Serialize()
		invalidGovBytes, _ = headergov.NewGovData(map[gov.ParamEnum]any{gov.GovernanceUnitPrice: uint64(200)}).Serialize()
		h                  = newHeaderGovModule(t, &params.ChainConfig{Istanbul: &params.IstanbulConfig{Epoch: 1000}})
		invalidVoteRlp     = common.FromHex("0xea9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e74")
	)

	h.HandleVote(1, vote)

	tcs := []struct {
		desc          string
		header        *types.Header
		expectedError error
	}{
		{desc: "valid vote", header: &types.Header{Number: big.NewInt(1), Vote: voteBytes}, expectedError: nil},
		{desc: "invalid vote rlp", header: &types.Header{Number: big.NewInt(1), Vote: []byte{1, 2, 3}}, expectedError: headergov.ErrInvalidRlp},
		{desc: "invalid vote bytes", header: &types.Header{Number: big.NewInt(1), Vote: invalidVoteRlp}, expectedError: headergov.ErrInvalidRlp},
		{desc: "valid gov", header: &types.Header{Number: big.NewInt(1000), Governance: govBytes}, expectedError: nil},
		{desc: "invalid gov rlp", header: &types.Header{Number: big.NewInt(1000), Governance: []byte{1, 2, 3}}, expectedError: headergov.ErrInvalidRlp},
		{desc: "gov must not be nil", header: &types.Header{Number: big.NewInt(1000), Governance: nil}, expectedError: ErrGovVerification},
		{desc: "gov mismatch", header: &types.Header{Number: big.NewInt(1000), Governance: invalidGovBytes}, expectedError: ErrGovVerification},
		{desc: "gov not on epoch", header: &types.Header{Number: big.NewInt(1001), Governance: []byte{1, 2, 3}}, expectedError: ErrGovInNonEpochBlock},
		{desc: "gov must be nil", header: &types.Header{Number: big.NewInt(2000), Governance: govBytes}, expectedError: ErrGovVerification},
		{desc: "valid gov", header: &types.Header{Number: big.NewInt(2000), Governance: nil}, expectedError: nil},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := h.VerifyHeader(tc.header)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestVerifyVote(t *testing.T) {
	var (
		h = newHeaderGovModule(t, &params.ChainConfig{
			Istanbul: &params.IstanbulConfig{
				Epoch: 1000,
			},
		})
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
		{desc: "valid govparam", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.GovernanceGovParamContract].Name, contract), expectedError: nil},
		{desc: "valid lower", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.Kip71LowerBoundBaseFee].Name, uint64(1)), expectedError: nil},
		{desc: "valid upper", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.Kip71UpperBoundBaseFee].Name, uint64(1e18)), expectedError: nil},
		{desc: "invalid govparam", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.GovernanceGovParamContract].Name, common.Address{}), expectedError: ErrGovParamNotAccount},
		{desc: "invalid govparam", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.GovernanceGovParamContract].Name, eoa), expectedError: ErrGovParamNotContract},
		{desc: "invalid lower", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.Kip71LowerBoundBaseFee].Name, uint64(1e18)), expectedError: ErrLowerBoundBaseFee},
		{desc: "invalid upper", vote: headergov.NewVoteData(common.Address{}, gov.Params[gov.Kip71UpperBoundBaseFee].Name, uint64(1)), expectedError: ErrUpperBoundBaseFee},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := h.VerifyVote(1, tc.vote)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetVotesInEpoch(t *testing.T) {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	})

	paramName := gov.Params[gov.GovernanceUnitPrice].Name
	v1 := headergov.NewVoteData(common.Address{1}, paramName, uint64(100))
	h.HandleVote(500, v1)
	v2 := headergov.NewVoteData(common.Address{2}, paramName, uint64(200))
	h.HandleVote(1500, v2)

	assert.Equal(t, map[uint64]headergov.VoteData{500: v1}, h.getVotesInEpoch(0))
	assert.Equal(t, map[uint64]headergov.VoteData{1500: v2}, h.getVotesInEpoch(1))
}

func TestGetExpectedGovernance(t *testing.T) {
	var (
		paramName = gov.Params[gov.GovernanceUnitPrice].Name
		config    = &params.ChainConfig{
			Istanbul: &params.IstanbulConfig{
				Epoch: 1000,
			},
		}
		h  = newHeaderGovModule(t, config)
		v1 = headergov.NewVoteData(common.Address{1}, paramName, uint64(100))
		v2 = headergov.NewVoteData(common.Address{2}, paramName, uint64(200))
		g1 = headergov.NewGovData(map[gov.ParamEnum]any{
			gov.GovernanceUnitPrice: uint64(100),
		})
		g2 = headergov.NewGovData(map[gov.ParamEnum]any{
			gov.GovernanceUnitPrice: uint64(200),
		})
	)

	h.HandleVote(500, v1)
	h.HandleVote(1500, v2)

	h.HandleGov(1000, g1)
	h.HandleGov(2000, g2)

	assert.Equal(t, g1, h.getExpectedGovernance(1000))
	assert.Equal(t, g2, h.getExpectedGovernance(2000))
}

func TestPrepareHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	var (
		h      = newHeaderGovModule(t, &params.ChainConfig{Istanbul: &params.IstanbulConfig{Epoch: 1000}})
		vote   = headergov.NewVoteData(h.nodeAddress, gov.Params[gov.GovernanceUnitPrice].Name, uint64(100))
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
	ps, _ := h.EffectiveParamSet(2001)
	assert.Equal(t, ps.UnitPrice, uint64(100))
}
