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
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		// Mirrors PrepareHeader's fallback: a node without a staking entry uses its
		// configured Rewardbase. We can't verify against staking, so accept the header.
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: unknown})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: evil}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})

	t.Run("staking lookup failure is treated as no entry", func(t *testing.T) {
		// VerifyHeader keeps the pre-existing behavior: an unreadable staking info leaves the
		// header unconstrained. Only the producer side (PrepareHeader) fails loudly.
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Proposer: node1, StakingErr: errors.New("staking unavailable")})
		header := &types.Header{Number: big.NewInt(100), Rewardbase: evil}
		assert.NoError(t, r.VerifyHeader(header, nil))
	})
}

// PrepareHeader must write the final reward address: with reward.deferredtxfee=false the fee is
// credited to header.Rewardbase while the transactions execute, so the value cannot change
// afterwards. The registered reward address wins whenever the staking info has one for this
// node under WeightedRandom; the configured Rewardbase is only the fallback.
func TestPrepareHeader(t *testing.T) {
	var (
		node1      = common.HexToAddress("0xa01") // stakingInfo.NodeIds[0]
		reward1    = common.HexToAddress("0xc01") // stakingInfo.RewardAddrs[0]
		unknown    = common.HexToAddress("0xdead")
		configured = common.HexToAddress("0xfff")
		number     = big.NewInt(100)
	)

	testcases := []struct {
		name string
		cfg  rewardModuleTestConfig
		want common.Address
	}{
		{"WeightedRandom with a staking entry uses the registered address", rewardModuleTestConfig{NodeAddress: node1, Rewardbase: configured}, reward1},
		{"WeightedRandom without a staking entry falls back to the configured address", rewardModuleTestConfig{NodeAddress: unknown, Rewardbase: configured}, configured},
		{"RoundRobin ignores the staking entry", rewardModuleTestConfig{Simple: true, NodeAddress: node1, Rewardbase: configured}, configured},
		{"configured address equal to the registered one", rewardModuleTestConfig{NodeAddress: node1, Rewardbase: reward1}, reward1},
		{"no configured address and no staking entry leaves the zero address", rewardModuleTestConfig{NodeAddress: unknown}, common.Address{}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// The proposer of the prepared block is this node: VerifyHeader must accept what PrepareHeader wrote.
			tc.cfg.Proposer = tc.cfg.NodeAddress
			r := makeTestRewardModuleWithConfig(t, tc.cfg)

			header := &types.Header{Number: number}
			require.NoError(t, r.PrepareHeader(header))
			assert.Equal(t, tc.want, header.Rewardbase)

			// Self-consistency: the rule PrepareHeader applies is the rule VerifyHeader enforces.
			assert.NoError(t, r.VerifyHeader(header, nil), "a header this node prepares must pass its own verification")
		})
	}

	t.Run("a staking entry with the zero reward address falls back to the configured address", func(t *testing.T) {
		// Both sides agree here: PrepareHeader uses the configured address and VerifyHeader leaves
		// a proposer whose registered address is zero unconstrained.
		si := makeTestStakingInfo([]uint64{5_000_001, 5_000_002, 5_000_003, 5_000_000}, false)
		si.RewardAddrs[0] = common.Address{}
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{NodeAddress: node1, Rewardbase: configured, Proposer: node1, StakingInfo: si})
		header := &types.Header{Number: number}
		require.NoError(t, r.PrepareHeader(header))
		assert.Equal(t, configured, header.Rewardbase)
		assert.NoError(t, r.VerifyHeader(header, nil))
	})

	t.Run("staking lookup failure aborts header preparation", func(t *testing.T) {
		// The producer must not fall back to the configured address when the staking info is
		// merely unreadable: validators that can read it would reject the block (ErrInvalidRewardbase).
		stakingErr := errors.New("staking unavailable")
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{NodeAddress: node1, Rewardbase: configured, StakingErr: stakingErr})

		header := &types.Header{Number: number, Rewardbase: common.HexToAddress("0x1234")}
		err := r.PrepareHeader(header)
		assert.ErrorIs(t, err, stakingErr)
		assert.Equal(t, common.HexToAddress("0x1234"), header.Rewardbase, "the header must be left untouched on failure")
	})

	t.Run("RoundRobin does not consult the staking info", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{Simple: true, NodeAddress: node1, Rewardbase: configured, StakingErr: errors.New("staking unavailable")})
		header := &types.Header{Number: number}
		require.NoError(t, r.PrepareHeader(header))
		assert.Equal(t, configured, header.Rewardbase)
	})
}
