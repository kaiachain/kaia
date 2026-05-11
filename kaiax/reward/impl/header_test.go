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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHeader(t *testing.T) {
	var (
		node1   = common.HexToAddress("0xa01") // matches stakingInfo.NodeIds[0]
		reward1 = common.HexToAddress("0xc01") // matches stakingInfo.RewardAddrs[0]
		unknown = common.HexToAddress("0xdead")
		evil    = common.HexToAddress("0xbeef")
	)

	t.Run("genesis is always accepted", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: node1})
		header := &types.Header{Number: big.NewInt(0), Rewardbase: evil}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})

	t.Run("matching rewardbase passes", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: node1})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: reward1}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})

	t.Run("mismatched rewardbase is rejected", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: node1})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: evil}
		assert.ErrorIs(t, r.VerifyHeader(header, nil), reward.ErrInvalidRewardbase)
	})

	t.Run("RoundRobin policy skips verification", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Simple: true, Proposer: node1})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: evil}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})

	t.Run("proposer not in staking info passes", func(t *testing.T) {
		// Mirrors FinalizeState's "else" branch where the node falls back to its
		// configured Rewardbase. We can't verify against staking, so accept the header.
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: unknown})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: evil}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})
}
