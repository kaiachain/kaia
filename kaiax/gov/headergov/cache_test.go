package headergov

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVote(t *testing.T) {
	var (
		cache = NewHeaderGovCache()
		v     = NewVoteData(common.HexToAddress("0x1"), string(gov.GovernanceUnitPrice), uint64(100))
		epoch = 3
		n     = 10
	)

	require.NotNil(t, v)
	for i := 0; i < n; i++ {
		cache.AddVote(uint64(i%epoch), uint64(i), v)
	}

	assert.Equal(t, n, len(cache.VoteBlockNums()))
	assert.Equal(t, n/epoch, len(cache.GroupedVotes()))
}

func TestGov(t *testing.T) {
	var (
		cache = NewHeaderGovCache()
		g     = NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(100)})
		n     = 10
	)

	require.NotNil(t, g)
	for i := 0; i < n; i++ {
		cache.AddGov(uint64(i), g)
	}

	assert.Equal(t, n, len(cache.Govs()))
	assert.Equal(t, n, len(cache.GovBlockNums()))
}
