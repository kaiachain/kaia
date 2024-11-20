package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestCalcSlotsInProposers(t *testing.T) {
	testHeaders, testStakingInfos, testParamSets := getValSetChainHeadersTestData(), getValSetStakingInfoTestData(), getValSetParamSetTestData()

	valSet := &valSetContext{
		blockNumber: uint64(1),
		rules: params.Rules{
			ChainID:    big.NewInt(0),
			IsIstanbul: true, IsLondon: true, IsEthTxType: true, IsMagma: true, IsKore: true, IsShanghai: true, IsCancun: true, IsKaia: true, IsRandao: true,
		},
		prevBlockResult: &blockResult{
			proposerPolicy: WeightedRandom,
			staking:        testStakingInfos[0],
			header:         testHeaders[0],
			author:         tgn,
			pSet:           testParamSets[0],
		},
	}
	qualified := valset.AddressList{tgn, n2, n3, n4}
	expectProposersIndexes := []int{0, 1, 2, 3}
	assert.Equal(t, expectProposersIndexes, calsSlotsInProposers(qualified, valSet))
}

func TestShuffleProposers(t *testing.T) {
	qualified := []common.Address{tgn, n2, n3, n4}
	proposersIndexes := []int{0, 1, 2, 3}
	expectProposers := []common.Address{tgn, n4, n2, n3}
	assert.Equal(t, expectProposers, shuffleProposers(qualified, proposersIndexes, testPrevHash))
}

func TestProposersBlockNum(t *testing.T) {
	for _, tc := range []struct {
		name                string
		blockNumber         uint64
		expectProposerBlock uint64
	}{
		{"block number is zero", 0, 0},
		{"block number is one", 1, 0},
		{"block number is a multiple of proposersUpdateInterval", testPUpdateInterval * 72, testPUpdateInterval * 71},
		{"block number is not a multiple of proposersUpdateInterval", testPUpdateInterval*72 - 3, testPUpdateInterval * 71},
	} {
		assert.Equal(t, tc.expectProposerBlock, calcProposerBlockNumber(tc.blockNumber, testPUpdateInterval), "tc name:%s", tc.name)
	}
}
