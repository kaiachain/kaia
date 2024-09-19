package headergov

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	headergov_types "github.com/kaiachain/kaia/kaiax/gov/headergov/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

var (
	ErrInvalidRlp      = headergov_types.ErrInvalidRlp
	ErrInvalidVoteData = headergov_types.ErrInvalidVoteData
	ErrInvalidGovData  = headergov_types.ErrInvalidGovData
)

func TestVerifyHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	var (
		vote               = NewVoteData(common.Address{1}, Params[GovernanceUnitPrice].Name, uint64(100))
		voteBytes, _       = NewVoteData(common.Address{1}, Params[GovernanceUnitPrice].Name, uint64(100)).Serialize()
		govBytes, _        = NewGovData(map[ParamEnum]interface{}{GovernanceUnitPrice: uint64(100)}).Serialize()
		invalidGovBytes, _ = NewGovData(map[ParamEnum]interface{}{GovernanceUnitPrice: uint64(200)}).Serialize()
		h                  = newHeaderGovModule(t, &params.ChainConfig{
			Istanbul: &params.IstanbulConfig{
				Epoch: 1000,
			},
		})
		twoTuple = common.FromHex("0xea9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e74")
	)

	h.HandleVote(1, vote)

	tcs := []struct {
		desc          string
		header        *types.Header
		expectedError error
	}{
		{desc: "valid vote", header: &types.Header{Number: big.NewInt(1), Vote: voteBytes}, expectedError: nil},
		{desc: "invalid vote rlp", header: &types.Header{Number: big.NewInt(1), Vote: []byte{1, 2, 3}}, expectedError: ErrInvalidRlp},
		{desc: "invalid vote bytes", header: &types.Header{Number: big.NewInt(1), Vote: twoTuple}, expectedError: ErrInvalidRlp},
		{desc: "valid gov", header: &types.Header{Number: big.NewInt(1000), Governance: govBytes}, expectedError: nil},
		{desc: "invalid gov rlp", header: &types.Header{Number: big.NewInt(1000), Governance: []byte{1, 2, 3}}, expectedError: ErrInvalidRlp},
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

func TestGetVotesInEpoch(t *testing.T) {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	})

	paramName := Params[GovernanceUnitPrice].Name
	v1 := NewVoteData(common.Address{1}, paramName, uint64(100))
	h.HandleVote(500, v1)
	v2 := NewVoteData(common.Address{2}, paramName, uint64(200))
	h.HandleVote(1500, v2)

	assert.Equal(t, map[uint64]VoteData{500: v1}, h.getVotesInEpoch(0))
	assert.Equal(t, map[uint64]VoteData{1500: v2}, h.getVotesInEpoch(1))
}

func TestGetExpectedGovernance(t *testing.T) {
	var (
		paramName = Params[GovernanceUnitPrice].Name
		config    = &params.ChainConfig{
			Istanbul: &params.IstanbulConfig{
				Epoch: 1000,
			},
		}
		h  = newHeaderGovModule(t, config)
		v1 = NewVoteData(common.Address{1}, paramName, uint64(100))
		v2 = NewVoteData(common.Address{2}, paramName, uint64(200))
		g1 = NewGovData(map[ParamEnum]interface{}{
			GovernanceUnitPrice: uint64(100),
		})
		g2 = NewGovData(map[ParamEnum]interface{}{
			GovernanceUnitPrice: uint64(200),
		})
	)

	h.HandleVote(500, v1)
	h.HandleVote(1500, v2)

	h.HandleGov(1000, g1)
	h.HandleGov(2000, g2)

	assert.Equal(t, g1, h.getExpectedGovernance(1000))
	assert.Equal(t, g2, h.getExpectedGovernance(2000))
}
