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

func TestPickRoundRobinProposer(t *testing.T) {
	cList := valset.AddressList{
		tgn, n[1], n[2], n[3],
	}

	tests := []struct {
		name           string
		policy         ProposerPolicy
		block          uint64
		round          uint64
		prevAuthor     common.Address
		expectProposer common.Address
		expectIdx      int
	}{
		// RoundRobin Policy Cases
		{"RoundRobin: block 1, round 0", RoundRobin, 1, 0, tgn, n[1], 1},
		{"RoundRobin: block 2, round 0", RoundRobin, 2, 0, n[1], n[2], 2},
		{"RoundRobin: block 3, round 0", RoundRobin, 3, 0, n[2], n[3], 3},
		{"RoundRobin: block 4, round 0", RoundRobin, 4, 0, n[3], tgn, 0},
		{"RoundRobin: block 11, round 0", RoundRobin, 11, 0, n[1], n[2], 2},
		{"RoundRobin: block 11, round 1", RoundRobin, 11, 1, n[1], n[3], 3},
		{"RoundRobin: block 11, round 2", RoundRobin, 11, 2, n[1], tgn, 0},

		// Sticky Policy Cases
		{"Sticky: block 1, round 0", Sticky, 1, 0, tgn, tgn, 0},
		{"Sticky: block 2, round 0", Sticky, 2, 0, tgn, tgn, 0},
		{"Sticky: block 11, round 0", Sticky, 11, 0, tgn, tgn, 0},
		{"Sticky: block 11, round 1", Sticky, 11, 1, tgn, n[1], 1},
		{"Sticky: block 11, round 2", Sticky, 11, 2, tgn, n[2], 2},
		{"Sticky: block 12, round 0", Sticky, 12, 0, n[2], n[2], 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proposer, idx := pickRoundRobinProposer(cList, tc.policy, tc.prevAuthor, tc.round)

			assert.Equal(t, tc.expectProposer, proposer, "unexpected proposer for %s", tc.name)
			assert.Equal(t, tc.expectIdx, idx, "unexpected index for %s", tc.name)
		})
	}
}
