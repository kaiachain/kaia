package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostInsertBlock(t *testing.T) {
	h := newHeaderGovModule(t, getTestChainConfig())
	vote, _ := headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100)).ToVoteBytes()
	gov, _ := headergov.NewGovData(gov.PartialParamSet{
		gov.GovernanceUnitPrice: uint64(100),
	}).ToGovBytes()

	voteBlock := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(50),
		Vote:   vote,
	})
	govBlock := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(100),
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

	// PostInsertBlock does not validate the block, so non-epoch block with gov is also accepted
	govBlock.Header().Number = big.NewInt(123)
	err = h.PostInsertBlock(govBlock)
	assert.NoError(t, err)
}

func TestHandleVoteGov(t *testing.T) {
	h := newHeaderGovModule(t, getTestChainConfig())
	v := headergov.NewVoteData(common.Address{1}, string(gov.GovernanceUnitPrice), uint64(100))
	require.NotNil(t, v)
	require.Equal(t, 0, len(h.groupedVotes))
	require.Equal(t, 1, len(h.governances)) // gov at genesis
	voteBlock := uint64(50)
	govBlock := uint64(100)

	// test duplicate vote handling
	for range 2 {
		err := h.HandleVote(voteBlock, v)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(h.groupedVotes))
	}

	// test duplicate gov handling
	for range 2 {
		err := h.HandleGov(govBlock, headergov.NewGovData(gov.PartialParamSet{
			gov.GovernanceUnitPrice: uint64(100),
		}))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(h.governances))
	}
}

func TestAddVote(t *testing.T) {
	var (
		v      = headergov.NewVoteData(common.HexToAddress("0x1"), string(gov.GovernanceUnitPrice), uint64(100))
		n      = 10
		config = getTestChainConfig()
	)
	config.Istanbul.Epoch = 3
	h := newHeaderGovModule(t, config)

	for i := 0; i < n; i++ {
		h.AddVote(uint64(i), v)
	}

	assert.Equal(t, n, len(h.VoteBlockNums()))
	assert.Equal(t, 4, len(h.groupedVotes)) // ceil(n/epoch)
}

func TestAddGov(t *testing.T) {
	var (
		g      = headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(100)})
		n      = 10
		config = getTestChainConfig()
	)
	config.Istanbul.Epoch = 1
	h := newHeaderGovModule(t, config)

	for i := 0; i < n; i++ {
		h.AddGov(uint64(i), g)
	}

	assert.Equal(t, n, len(h.governances))
	assert.Equal(t, n, len(h.GovBlockNums()))
}
