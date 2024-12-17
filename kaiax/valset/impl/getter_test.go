package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func numToAddr(n int) common.Address {
	return common.BigToAddress(big.NewInt(int64(n)))
}

func numsToAddrs(n ...int) []common.Address {
	addrs := make([]common.Address, len(n))
	for i, num := range n {
		addrs[i] = common.BigToAddress(big.NewInt(int64(num)))
	}
	return addrs
}

func addrsToNums(addrs []common.Address) []int {
	nums := make([]int, len(addrs))
	for i, addr := range addrs {
		b := new(big.Int).SetBytes(addr.Bytes())
		nums[i] = int(b.Int64())
	}
	return nums
}

func hexToAddrs(s ...string) []common.Address {
	addrs := make([]common.Address, len(s))
	for i, addr := range s {
		addrs[i] = common.HexToAddress(addr)
	}
	return addrs
}

func TestGetCouncilGenesis(t *testing.T) {
	var (
		ctrl      = gomock.NewController(t)
		mockChain = chain_mock.NewMockBlockChain(ctrl)
		v         = &ValsetModule{InitOpts: InitOpts{Chain: mockChain}}
	)
	defer ctrl.Finish()
	// Kairos block 0
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(&types.Header{
		Extra: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f89af8549499fb17d324fa0e07f23b49d09028ac0919414db694b74ff9dea397fe9e231df545eb53fe2adf776cb294571e53df607be97431a5bbefca1dffe5aef56f4d945cb1a7dccbd0dc446e3640898ede8820368554c8b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
	}).AnyTimes()

	council, err := v.getCouncilGenesis()
	assert.NoError(t, err)
	assert.Equal(t, hexToAddrs(
		"0x571e53df607be97431a5bbefca1dffe5aef56f4d",
		"0x5cb1a7dccbd0dc446e3640898ede8820368554c8",
		"0x99fb17d324fa0e07f23b49d09028ac0919414db6",
		"0xb74ff9dea397fe9e231df545eb53fe2adf776cb2",
	), council.List())
}

func TestGetCouncilDB(t *testing.T) {
	var (
		db       = database.NewMemDB()
		voteNums = []uint64{0, 2, 5, 6}
		council0 = numsToAddrs(0) // affects blocks 0,1,2
		council2 = numsToAddrs(2) // affects blocks 3,4,5
		council5 = numsToAddrs(5) // affects block 6
		council6 = numsToAddrs(6) // affects blocks 7,8,9,...
	)

	writeValidatorVoteBlockNums(db, voteNums)
	writeCouncil(db, 0, council0)
	writeCouncil(db, 2, council2)
	writeCouncil(db, 5, council5)
	writeCouncil(db, 6, council6)

	testcases := []struct {
		num     uint64
		voteNum uint64
		council []common.Address
	}{
		{0, 0, council0},
		{1, 0, council0},
		{2, 0, council0},
		{3, 2, council2},
		{4, 2, council2},
		{5, 2, council2},
		{6, 5, council5},
		{7, 6, council6},
		{8, 6, council6},
	}
	for _, tc := range testcases {
		voteNum := lastNumLessThan(voteNums, tc.num)
		assert.Equal(t, tc.voteNum, voteNum, tc.num)

		v := &ValsetModule{InitOpts: InitOpts{ChainKv: db}}
		council, err := v.getCouncilDB(tc.num)
		assert.NoError(t, err)
		assert.Equal(t, tc.council, council.List(), tc.num)
	}
}

func TestFilterValidators(t *testing.T) {
	var (
		governingNode = numToAddr(3)
		aL            = uint64(1000000) // Less than minstaking
		aM            = uint64(2000000) // Exactly minstaking
		pset          = gov.ParamSet{
			GoverningNode: governingNode,
			MinimumStake:  big.NewInt(int64(aM)),
		}
	)

	testcases := []struct {
		desc    string
		council []common.Address
		amounts []uint64
		single  bool
		demoted []common.Address
	}{
		{
			desc:    "none mode, all qualified",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{aM, aM, aM, aM, aM},
			single:  false,
			demoted: numsToAddrs(),
		},
		{
			desc:    "none mode, some demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, aL, aM, aM, aM},
			single:  false,
			demoted: numsToAddrs(1, 2),
		},
		{
			desc:    "none mode, all demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, aL, aL},
			single:  false,
			demoted: numsToAddrs(), // If all are demoted, none are demoted.
		},
		{
			desc:    "single mode, all qualified",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{aM, aM, aM, aM, aM},
			single:  true,
			demoted: numsToAddrs(),
		},
		{
			desc:    "single mode, some understaked",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, aL, aM, aM, aM},
			single:  true,
			demoted: numsToAddrs(1, 2),
		},
		{
			desc:    "single mode, governingNode and others are understaked",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, aM, aM},
			single:  true,
			demoted: numsToAddrs(1, 2), // despite governingNode(3) understaked, it is not demoted.
		},
		{
			desc:    "single mode, only governingNode is staked enough",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, aM, aL, aL},
			single:  true,
			demoted: numsToAddrs(1, 2, 4, 5),
		},
		{
			desc:    "single mode, all demoted",
			council: numsToAddrs(1, 2, 3, 4, 5),
			amounts: []uint64{0, 0, 0, 0, 0},
			single:  true,
			demoted: numsToAddrs(), // If all are demoted, none are demoted.
		},
	}
	for _, tc := range testcases {
		council := valset.NewAddressSet(tc.council)
		si := &staking.StakingInfo{
			NodeIds:          tc.council,
			StakingContracts: tc.council,
			RewardAddrs:      tc.council,
			StakingAmounts:   tc.amounts,
		}
		if tc.single {
			pset.GovernanceMode = "single"
		} else {
			pset.GovernanceMode = "none"
		}

		demoted := filterValidatorsIstanbul(council, si, pset)
		assert.Equal(t, tc.demoted, demoted.List(), tc.desc)
	}
}

func TestRandomCommittee(t *testing.T) {
}

func TestRoundRobinProposer(t *testing.T) {
	testcases := []struct {
		qualified    []common.Address
		prevProposer common.Address
		round        uint64
		expected     common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 0, numToAddr(3)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 1, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 2, numToAddr(1)}, // wrap-around by round
		{numsToAddrs(1, 2, 3, 4), numToAddr(4), 0, numToAddr(1)}, // wrap-around by prevProposer
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 0, numToAddr(2)}, // fallback to prevIdx = 0
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 1, numToAddr(3)}, // fallback to prevIdx = 0
	}
	for _, tc := range testcases {
		currProposer := selectRoundRobinProposer(valset.NewAddressSet(tc.qualified), tc.prevProposer, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

func TestStickyProposer(t *testing.T) {
	testcases := []struct {
		qualified    []common.Address
		prevProposer common.Address
		round        uint64
		expected     common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 0, numToAddr(2)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 1, numToAddr(3)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 2, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(2), 3, numToAddr(1)}, // wrap-around by round
		{numsToAddrs(1, 2, 3, 4), numToAddr(4), 0, numToAddr(4)},
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 0, numToAddr(1)}, // fallback to prevIdx = 0
		{numsToAddrs(1, 2, 3, 4), numToAddr(7), 1, numToAddr(2)}, // fallback to prevIdx = 0
	}
	for _, tc := range testcases {
		currProposer := selectStickyProposer(valset.NewAddressSet(tc.qualified), tc.prevProposer, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

func TestWeightedRandomProposer_Select(t *testing.T) {
	testcases := []struct {
		proposerList  []common.Address
		listSourceNum uint64
		num           uint64
		round         uint64
		expected      common.Address
	}{
		{numsToAddrs(1, 2, 3, 4), 100, 101, 0, numToAddr(1)},
		{numsToAddrs(1, 2, 3, 4), 100, 106, 0, numToAddr(2)}, // wrap-around by num
		{numsToAddrs(1, 2, 3, 4), 100, 103, 3, numToAddr(2)}, // wrap-around by round
	}
	for _, tc := range testcases {
		currProposer := selectWeightedRandomProposer(tc.proposerList, tc.listSourceNum, tc.num, tc.round)
		assert.Equal(t, tc.expected, currProposer)
	}
}

func TestWeightedRandomProposer_ListWeighted(t *testing.T) {
	var (
		blockHash = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84")
		aM        = uint64(2000000) // Exactly minstaking
	)

	testcases := []struct {
		desc         string
		qualified    []common.Address
		amounts      []uint64
		useGini      bool
		expectedFreq []int            // expected appearance frequency
		expectedList []common.Address // expected proposer list. Leave nil to skip check.
	}{
		{
			desc:      "positive stakes, gini",
			qualified: numsToAddrs(0, 1, 2, 3),
			amounts:   []uint64{1 * aM, 2 * aM, 3 * aM, 4 * aM},
			useGini:   true,
			// gini=0.25, exponent=0.8 -> Adjusted staking amounts = [109856,191270,264558,333021] -> Percentile weights = [12,21,29,37]
			expectedFreq: []int{12, 21, 29, 37},
		},
		{
			desc:         "positive stakes, no gini",
			qualified:    numsToAddrs(0, 1, 2, 3),
			amounts:      []uint64{1 * aM, 2 * aM, 3 * aM, 4 * aM},
			useGini:      false,
			expectedFreq: []int{10, 20, 30, 40},
			expectedList: numsToAddrs(1, 1, 3, 2, 0, 3, 2, 3, 1, 1, 3, 1, 3, 2, 3, 1, 2, 2, 0, 3, 3, 2, 3, 3, 2, 1, 1, 1, 3, 3, 2, 3, 1, 2, 3, 1, 2, 3, 2, 2, 0, 3, 2, 2, 1, 3, 1, 0, 2, 2, 2, 1, 1, 3, 2, 3, 1, 2, 3, 3, 0, 3, 3, 3, 2, 3, 3, 3, 2, 3, 3, 3, 0, 3, 2, 0, 1, 3, 2, 2, 1, 2, 3, 1, 3, 3, 0, 0, 3, 1, 2, 3, 3, 2, 2, 3, 2, 0, 3, 2),
		},
		{
			desc:         "zero stakes",
			qualified:    numsToAddrs(0, 1, 2, 3), // Note: validators can be qualified with zero stakes, if all are understaked.
			amounts:      []uint64{0, 0, 0, 0},
			useGini:      false,
			expectedFreq: []int{1, 1, 1, 1},
			expectedList: numsToAddrs(1, 3, 0, 2),
		},
		{
			desc:         "severe inequality",
			qualified:    numsToAddrs(0, 1, 2, 3),
			amounts:      []uint64{aM, aM, aM, 200 * aM},
			useGini:      false,
			expectedFreq: []int{1, 1, 1, 99}, // at least 1 slot is guaranteed for each validator.
		},
	}
	for _, tc := range testcases {
		qualified := valset.NewAddressSet(tc.qualified)
		si := &staking.StakingInfo{
			NodeIds:          tc.qualified,
			StakingContracts: tc.qualified,
			RewardAddrs:      tc.qualified,
			StakingAmounts:   tc.amounts,
		}
		proposerList := generateProposerListWeighted(qualified, si, tc.useGini, blockHash)

		freq := make(map[common.Address]int)
		for _, addr := range proposerList {
			freq[addr]++
		}
		for idx, addr := range tc.qualified {
			assert.Equal(t, tc.expectedFreq[idx], freq[addr], tc.desc)
		}

		if tc.expectedList != nil && !assert.Equal(t, tc.expectedList, proposerList, tc.desc) {
			t.Logf("expected: %v", addrsToNums(tc.expectedList))
			t.Logf("actual: %v", addrsToNums(proposerList))
		}
	}
}

func TestWeightedRandomProposer_ListUniform(t *testing.T) {
	var (
		blockHash    = common.HexToHash("0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84")
		qualified    = valset.NewAddressSet(numsToAddrs(0, 1, 2, 3))
		proposerList = generateProposerListUniform(qualified, blockHash)
		expectedList = numsToAddrs(1, 3, 0, 2)
	)
	assert.Equal(t, expectedList, proposerList)
}
