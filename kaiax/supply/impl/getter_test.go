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
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSupply(t *testing.T) {
	suite.Run(t, new(SupplyTestSuite))
}

// Test individual private getters.
func (s *SupplyTestSuite) TestFromState() {
	t := s.T()
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

// Test reading canonical burn amounts from the state.
func (s *SupplyTestSuite) TestCanonicalBurn() {
	t := s.T()
	s.insertBlocks()

	// Delete state at 199
	root := s.db.ReadBlockByNumber(199).Root()
	s.db.DeleteTrieNode(root.ExtendZero())

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

// Test reading rebalance memo from the contract.
func (s *SupplyTestSuite) TestRebalanceMemo() {
	t := s.T()
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
