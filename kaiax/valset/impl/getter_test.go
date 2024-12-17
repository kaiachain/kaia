// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
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

func TestLastNumLessThan(t *testing.T) {
	var (
		nums    = []uint64{0, 2, 4}
		inputs  = []uint64{0, 1, 2, 3, 4, 5}
		outputs = []uint64{0, 0, 0, 2, 2, 4}
	)
	for i, input := range inputs {
		assert.Equal(t, outputs[i], lastNumLessThan(nums, input))
	}
}

func TestGetCouncilDB(t *testing.T) {
	var (
		// An hypothetical chain where:
		// voteNums         0     2     4
		// num              0  1  2  3  4  5  6
		// council          A  A  A  B  B  C  C
		// lastNumLessThan  0  0  0  2  2  4  4
		voteNums = []uint64{0, 2, 4}
		setA     = numsToAddrs(1)
		setB     = numsToAddrs(1, 2)
		setC     = numsToAddrs(1, 2, 3)
		db       = database.NewMemDB()
	)
	writeValidatorVoteBlockNums(db, voteNums)
	writeCouncil(db, 0, setA)
	writeCouncil(db, 2, setB)
	writeCouncil(db, 4, setC)

	// Test getCouncilDB under various LowestScannedVoteNum.
	type expected struct { // expected result of getCouncilDB()
		council []common.Address
		ok      bool
	}
	testcases := []struct {
		desc                 string
		lowestScannedVoteNum uint64
		expected             [7]expected // expected output for blocks 0..6
	}{
		{
			"migration complete",
			0,
			[7]expected{{setA, true}, {setA, true}, {setA, true}, {setB, true}, {setB, true}, {setC, true}, {setC, true}},
		},
		{
			"migration incomplete",
			4,
			[7]expected{{nil, false}, {nil, false}, {nil, false}, {setB, false}, {setB, false}, {setC, true}, {setC, true}},
		},
	}
	for _, tc := range testcases {
		v := &ValsetModule{InitOpts: InitOpts{ChainKv: db}}
		writeLowestScannedVoteNum(db, tc.lowestScannedVoteNum)

		for i := uint64(0); i < 7; i++ {
			council, ok, err := v.getCouncilDB(i)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected[i].ok, ok)
			if ok {
				assert.Equal(t, tc.expected[i].council, council.List())
			} else {
				assert.Nil(t, council)
			}
		}
	}
}

func TestParseValidatorVote(t *testing.T) {
	testcases := []struct {
		voteHex string
		voteKey gov.ValidatorParamName
		voteVal []common.Address
		ok      bool
	}{
		{ // Empty
			"0x",
			"", nil, false,
		},
		{ // Malformed
			"0xabcd",
			"", nil, false,
		},
		{ // Kairos block 83863326 (not a validator vote)
			"0xf09499fb17d324fa0e07f23b49d09028ac0919414db694676f7665726e616e63652e756e6974707269636585ae9f7bcc00",
			"", nil, false,
		},
		{ // Kairos block 4202779 (add one address)
			"0xf8429499fb17d324fa0e07f23b49d09028ac0919414db697676f7665726e616e63652e61646476616c696461746f72948a88a093c05376886754a9b70b0d0a826a5e64be",
			"governance.addvalidator", hexToAddrs("0x8a88a093c05376886754a9b70b0d0a826a5e64be"), true,
		},
		{ // Kairos block 4740968 (remove one address)
			"0xf8459499fb17d324fa0e07f23b49d09028ac0919414db69a676f7665726e616e63652e72656d6f766576616c696461746f72949419fa2e3b9eb1158de31be66c586a52f49c5de7",
			"governance.removevalidator", hexToAddrs("0x9419fa2e3b9eb1158de31be66c586a52f49c5de7"), true,
		},
		{ // Mainnet block 90897408 (remove one address in hex string)
			"0xf85b9452d41ca72af615a1ac3301b0a93efa222ecc75419a676f7665726e616e63652e72656d6f766576616c696461746f72aa307831366331393235383561306162323462353532373833623462663764386463396636383535633335",
			"governance.removevalidator", hexToAddrs("0x16c192585a0ab24b552783b4bf7d8dc9f6855c35"), true,
		},
	}
	for _, tc := range testcases {
		header := &types.Header{
			Vote: hexutil.MustDecode(tc.voteHex),
		}
		voteKey, voteVal, ok := parseValidatorVote(header)
		assert.Equal(t, tc.voteKey, voteKey)
		assert.Equal(t, tc.voteVal, voteVal)
		assert.Equal(t, tc.ok, ok)
	}
}

func TestApplyVote(t *testing.T) {
	var (
		governingNode  = numToAddr(3)
		initialCouncil = numsToAddrs(1, 2, 3)

		voteAdd1, _    = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(1)).ToVoteBytes()
		voteAdd6, _    = headergov.NewVoteData(governingNode, string(gov.AddValidator), numToAddr(6)).ToVoteBytes()
		voteRemove2, _ = headergov.NewVoteData(governingNode, string(gov.RemoveValidator), numToAddr(2)).ToVoteBytes()
		voteRemove3, _ = headergov.NewVoteData(governingNode, string(gov.RemoveValidator), numToAddr(3)).ToVoteBytes()
		voteRemove7, _ = headergov.NewVoteData(governingNode, string(gov.RemoveValidator), numToAddr(7)).ToVoteBytes()
	)
	testcases := []struct {
		voteData []byte
		council  []common.Address
		modified bool
	}{
		{nil, numsToAddrs(1, 2, 3), false},
		{voteAdd1, numsToAddrs(1, 2, 3), false},    // Cannot add already existing one
		{voteAdd6, numsToAddrs(1, 2, 3, 6), true},  // Add one
		{voteRemove2, numsToAddrs(1, 3), true},     // Remove one
		{voteRemove3, numsToAddrs(1, 2, 3), false}, // Cannot remove governingNode
		{voteRemove7, numsToAddrs(1, 2, 3), false}, // Cannot remove non-existing one
	}
	for i, tc := range testcases {
		header := &types.Header{
			Vote: tc.voteData,
		}
		council := valset.NewAddressSet(initialCouncil)
		modified := applyVote(header, council, governingNode)
		assert.Equal(t, tc.modified, modified)
		assert.Equal(t, tc.council, council.List(), i) // council is modified in-place
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
