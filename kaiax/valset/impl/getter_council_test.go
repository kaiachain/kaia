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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

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
			[7]expected{{nil, false}, {nil, false}, {nil, false}, {nil, false}, {nil, false}, {setC, true}, {setC, true}},
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
		voteKey gov.ParamName
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
