package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostInsertBlock(t *testing.T) {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 10,
		},
	})
	vote, _ := headergov.NewVoteData(common.Address{1}, gov.Params[gov.GovernanceUnitPrice].Name, uint64(100)).ToVoteBytes()
	gov, _ := headergov.NewGovData(map[gov.ParamEnum]any{
		gov.GovernanceUnitPrice: uint64(100),
	}).Serialize()

	voteBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(5),
		Vote:   vote,
	})
	govBlock := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(10),
		Governance: gov,
	})

	require.Equal(t, 0, len(h.cache.GroupedVotes()))
	require.Equal(t, 1, len(h.cache.Govs())) // gov at genesis

	err := h.PostInsertBlock(voteBlock)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(h.cache.GroupedVotes()))

	err = h.PostInsertBlock(govBlock)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(h.cache.Govs()))
}

func TestHandleVoteGov(t *testing.T) {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 10,
		},
	})

	v := headergov.NewVoteData(common.Address{1}, gov.Params[gov.GovernanceUnitPrice].Name, uint64(100))
	require.NotNil(t, v)
	require.Equal(t, 0, len(h.cache.GroupedVotes()))
	require.Equal(t, 1, len(h.cache.Govs())) // gov at genesis
	voteBlock := uint64(5)
	govBlock := uint64(10)

	// test duplicate vote handling
	for range 2 {
		err := h.HandleVote(voteBlock, v)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(h.cache.GroupedVotes()))
	}

	// test duplicate gov handling
	for range 2 {
		err := h.HandleGov(govBlock, headergov.NewGovData(map[gov.ParamEnum]any{
			gov.GovernanceUnitPrice: uint64(100),
		}))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(h.cache.Govs()))
	}
}
