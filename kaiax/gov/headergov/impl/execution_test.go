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
	h := newHeaderGovModule(t, &params.ChainConfig{Istanbul: &params.IstanbulConfig{Epoch: 10}})
	vote, _ := headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100)).ToVoteBytes()
	gov, _ := headergov.NewGovData(gov.PartialParamSet{
		gov.GovernanceUnitPrice: uint64(100),
	}).ToGovBytes()

	voteBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(5),
		Vote:   vote,
	})
	govBlock := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(10),
		Governance: gov,
	})

	require.Equal(t, 0, len(h.groupedVotes))
	require.Equal(t, 1, len(h.governances)) // gov at genesis

	err := h.PostInsertBlock(voteBlock)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(h.groupedVotes))

	err = h.PostInsertBlock(govBlock)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(h.governances))
}

func TestHandleVoteGov(t *testing.T) {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 10,
		},
	})

	v := headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100))
	require.NotNil(t, v)
	require.Equal(t, 0, len(h.groupedVotes))
	require.Equal(t, 1, len(h.governances)) // gov at genesis
	voteBlock := uint64(5)
	govBlock := uint64(10)

	// test duplicate vote handling
	for i := 0; i < 2; i++ {
		err := h.HandleVote(voteBlock, v)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(h.groupedVotes))
	}

	// test duplicate gov handling
	for i := 0; i < 2; i++ {
		err := h.HandleGov(govBlock, headergov.NewGovData(gov.PartialParamSet{
			gov.GovernanceUnitPrice: uint64(100),
		}))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(h.governances))
	}
}

func TestVote(t *testing.T) {
	var (
		v     = headergov.NewVoteData(common.HexToAddress("0x1"), string(gov.GovernanceUnitPrice), uint64(100))
		epoch = 3
		n     = 10
		h     = newHeaderGovModule(t, &params.ChainConfig{Istanbul: &params.IstanbulConfig{Epoch: uint64(epoch)}})
	)

	for i := 0; i < n; i++ {
		h.AddVote(uint64(i%epoch), uint64(i), v)
	}

	assert.Equal(t, n, len(h.VoteBlockNums()))
	assert.Equal(t, n/epoch, len(h.groupedVotes))
}

func TestGov(t *testing.T) {
	var (
		g = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(100)})
		n = 10
		h = newHeaderGovModule(t, &params.ChainConfig{Istanbul: &params.IstanbulConfig{Epoch: 1}})
	)

	for i := 0; i < n; i++ {
		h.AddGov(uint64(i), g)
	}

	assert.Equal(t, n, len(h.governances))
	assert.Equal(t, n, len(h.GovBlockNums()))
}
