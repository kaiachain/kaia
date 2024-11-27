package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func calcStakingAmounts(sInfo *staking.StakingInfo, cList []common.Address) map[common.Address]uint64 {
	stakingAmounts := make(map[common.Address]uint64, len(cList))
	for _, node := range cList {
		stakingAmounts[node] = uint64(0)
	}
	for _, consolidated := range sInfo.ConsolidatedNodes() {
		for _, nAddr := range consolidated.NodeIds {
			if _, ok := stakingAmounts[nAddr]; ok {
				stakingAmounts[nAddr] = consolidated.StakingAmount
			}
		}
	}
	return stakingAmounts
}
func TestCalcSlotsInProposers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test data
	qualified := valset.AddressList{tgn, n[1], n[2], n[3]}
	rules := params.Rules{ChainID: big.NewInt(1), IsIstanbul: true}
	pSet := prepareTestGovParam(testProposerPolicy, testSubGroupSize, tgn)
	sInfo := prepareTestStakingInfo(0, []uint64{0, 1, 2, 3}, []uint64{aL, aL, aL, aL})
	stakingAmounts := calcStakingAmounts(sInfo, qualified)

	// expect data
	expectProposersIndexes := []int{0, 1, 2, 3}
	assert.Equal(t, expectProposersIndexes, calsSlotsInProposers(qualified, rules, pSet, sInfo, stakingAmounts))
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
