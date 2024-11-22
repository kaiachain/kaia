package impl

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
)

func TestCalcSlotsInProposers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, tm, _ := newTestVModule(ctrl)
	tm.prepareMockExpectGovParam(0, testProposerPolicy, testSubGroupSize, tgn)
	tm.prepareMockExpectStakingInfo(0, []uint64{0, 1, 2, 3}, []uint64{aL, aL, aL, aL})
	tm.prepareMockExpectHeader(0, nil, nil, tgn)

	valCtx, err := newValSetContext(vModule, 1)
	assert.NoError(t, err)

	qualified := valset.AddressList{tgn, n[1], n[2], n[3]}
	expectProposersIndexes := []int{0, 1, 2, 3}
	assert.Equal(t, expectProposersIndexes, calsSlotsInProposers(qualified, valCtx))
}

func TestShuffleProposers(t *testing.T) {
	qualified := []common.Address{tgn, n[1], n[2], n[3]}
	proposersIndexes := []int{0, 1, 2, 3}
	expectProposers := []common.Address{tgn, n[3], n[1], n[2]}
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
