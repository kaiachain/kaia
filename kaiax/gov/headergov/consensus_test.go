package headergov

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyHeader(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	paramName := Params[GovernanceUnitPrice].Name
	config := &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	}

	h := newHeaderGovModule(t, config)
	vote := NewVoteData(common.Address{1}, paramName, uint64(100))
	require.NotNil(t, vote)
	h.HandleVote(500, vote)

	gov := NewGovData(map[ParamEnum]interface{}{
		GovernanceUnitPrice: uint64(100),
	})
	require.NotNil(t, gov)

	govBytes, err := gov.Serialize()
	require.NoError(t, err)

	tcs := []struct {
		blockNum uint64
		gov      []byte
		isError  bool
	}{
		// no errors
		{999, nil, false},
		{1000, govBytes, false},
		{1001, nil, false},

		// errors
		{999, govBytes, true},
		{1000, nil, true},
		{1001, govBytes, true},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("BlockNum=%d,HasGov=%v", tc.blockNum, tc.gov != nil), func(t *testing.T) {
			err := h.VerifyHeader(&types.Header{
				Number:     big.NewInt(int64(tc.blockNum)),
				Governance: tc.gov,
			})
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
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
