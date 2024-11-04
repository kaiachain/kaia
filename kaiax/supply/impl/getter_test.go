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

package supply

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test individual getters
func (s *SupplyTestSuite) TestFromState() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	testcases := s.testcases()
	for _, tc := range testcases {
		fromState, err := s.s.totalSupplyFromState(tc.number)
		require.NoError(t, err)
		bigEqual(t, tc.expectFromState, fromState, tc.number)

		if s.T().Failed() {
			s.dumpState(tc.number)
			break
		}
	}
}

func (s *SupplyTestSuite) TestCheckpoint() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	for _, tc := range s.testcases() {
		checkpoint, err := s.s.getAccReward(tc.number)
		require.NoError(t, err)

		expected := tc.expectTotalSupply
		bigEqual(t, expected.TotalMinted, checkpoint.Minted, tc.number)
		bigEqual(t, expected.BurntFee, checkpoint.BurntFee, tc.number)
	}
}

func (s *SupplyTestSuite) TestCanonicalBurn() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	// Delete state at 199
	root := s.dbm.ReadBlockByNumber(199).Root()
	s.dbm.DeleteTrieNode(root.ExtendZero())

	// State unavailable at 199
	zero, dead, err := s.s.getCanonicalBurn(199)
	assert.Error(t, err)
	assert.Nil(t, zero)
	assert.Nil(t, dead)

	// State available at 200
	zero, dead, err = s.s.getCanonicalBurn(200)
	assert.NoError(t, err)
	assert.Equal(t, "1000000000000000000000000000", zero.String())
	assert.Equal(t, "2000000000000000000000000000", dead.String())
}

func (s *SupplyTestSuite) TestRebalanceMemo() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	// rebalance not configured
	amount, err := s.s.getRebalanceBurn(199, nil, common.Address{})
	require.NoError(t, err)
	assert.Equal(t, "0", amount.String())

	// num < forkNum
	amount, err = s.s.getRebalanceBurn(199, big.NewInt(200), addrKip103)
	require.NoError(t, err)
	assert.Equal(t, "0", amount.String())

	// num >= forkNum
	amount, err = s.s.getRebalanceBurn(200, big.NewInt(200), addrKip103)
	require.NoError(t, err)
	assert.Equal(t, "650511428500000000000", amount.String())

	amount, err = s.s.getRebalanceBurn(300, big.NewInt(300), addrKip160)
	require.NoError(t, err)
	assert.Equal(t, "-79200000000000000000", amount.String())

	// call failed: bad contract address
	amount, err = s.s.getRebalanceBurn(200, big.NewInt(200), addrFund1)
	require.Error(t, err)
	require.Nil(t, amount)
}

func (s *SupplyTestSuite) TestGetTotalSupply() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	for _, tc := range s.testcases() {
		ts, err := s.s.GetTotalSupply(tc.number)
		require.NoError(t, err)
		assert.Equal(t, tc.expectTotalSupply, ts)
	}
}

func (s *SupplyTestSuite) TestGetTotalSupply_PartialInfo() {
	t := s.T()
	require.Nil(t, s.s.loadLastCheckpoint())
	s.insertBlocks()

	// We will test on block 200.
	var num uint64 = 200
	var expected *supply.TotalSupply
	for _, tc := range s.testcases() {
		if tc.number == num {
			expected = tc.expectTotalSupply
			break
		}
	}

	// Missing state trie; returns partial data.
	root := s.dbm.ReadBlockByNumber(num).Root()
	s.dbm.DeleteTrieNode(root.ExtendZero())

	ts, err := s.s.GetTotalSupply(num)
	assert.ErrorContains(t, err, "missing trie node")
	partial := &supply.TotalSupply{
		TotalSupply: nil,
		TotalMinted: expected.TotalMinted,
		TotalBurnt:  nil,
		BurntFee:    expected.BurntFee,
		ZeroBurn:    nil,
		DeadBurn:    nil,
		Kip103Burn:  expected.Kip103Burn,
		Kip160Burn:  expected.Kip160Burn,
	}
	assert.Equal(t, partial, ts)

	// Misconfigured KIP-103; returns partial data.
	s.chain.Config().Kip103ContractAddress = addrFund1

	ts, err = s.s.GetTotalSupply(num)
	assert.ErrorContains(t, err, "missing trie node") // Errors are concatenated
	assert.ErrorContains(t, err, "no contract code")
	partial = &supply.TotalSupply{
		TotalSupply: nil,
		TotalMinted: expected.TotalMinted,
		TotalBurnt:  nil,
		BurntFee:    expected.BurntFee,
		ZeroBurn:    nil,
		DeadBurn:    nil,
		Kip103Burn:  nil,
		Kip160Burn:  expected.Kip160Burn,
	}
	assert.Equal(t, partial, ts)

	// No SupplyCheckpoint; returns nil.
	WriteLastSupplyCheckpointNumber(s.s.ChainKv, num-(num%128))
	DeleteSupplyCheckpoint(s.s.ChainKv, num-(num%128))
	s.s.accRewardCache.Purge()

	ts, err = s.s.GetTotalSupply(num)
	assert.ErrorIs(t, err, supply.ErrNoCheckpoint)
	assert.Nil(t, ts)
}
