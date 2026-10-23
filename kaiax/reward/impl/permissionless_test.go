// Copyright 2026 The Kaia Authors
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

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testPermissionlessBlock = big.NewInt(200)

	// Committee member that is not reward-eligible at the source block, e.g. promoted to
	// ValActive by the epoch transition for the block. Its reward address is a PublicDelegation.
	promotedNode = common.HexToAddress("0xa05")
	promotedPD   = common.HexToAddress("0xc05")
)

// makeTestPermissionlessStakingInfo returns the default staking info (four reward-eligible
// nodes) plus the AddressBookV2 profile index, which also has promotedNode.
func makeTestPermissionlessStakingInfo() *staking.StakingInfo {
	si := makeTestStakingInfo([]uint64{5_000_001, 5_000_002, 5_000_003, 5_000_000}, false)
	si.ProfileRewardAddrs = map[common.Address]common.Address{promotedNode: promotedPD}
	for i, n := range si.NodeIds {
		si.ProfileRewardAddrs[n] = si.RewardAddrs[i]
	}
	return si
}

func TestVerifyHeader_Permissionless(t *testing.T) {
	var (
		node1   = common.HexToAddress("0xa01")
		reward1 = common.HexToAddress("0xc01")
		local   = common.HexToAddress("0xfff") // the proposer's --rewardbase
		unknown = common.HexToAddress("0xdead")
	)

	testcases := []struct {
		name       string
		num        int64
		proposer   common.Address
		rewardbase common.Address
		stakingErr error
		expected   error
	}{
		// Before the fork, a proposer without a staking entry is not checked (unchanged).
		{"pre-fork, proposer not in staking info", 199, promotedNode, local, nil, nil},
		{"pre-fork, staking info error", 199, node1, local, errors.New("no state"), nil},
		{"pre-fork, mismatch", 199, node1, local, nil, reward.ErrInvalidRewardbase},

		{"fork block, eligible proposer", 200, node1, reward1, nil, nil},
		{"fork block, eligible proposer mismatch", 200, node1, local, nil, reward.ErrInvalidRewardbase},
		{"promoted proposer, profile address", 300, promotedNode, promotedPD, nil, nil},
		{"promoted proposer, local rewardbase", 300, promotedNode, local, nil, reward.ErrInvalidRewardbase},
		{"proposer without profile", 300, unknown, local, nil, reward.ErrInvalidRewardbase},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{
				Proposer:            tc.proposer,
				PermissionlessBlock: testPermissionlessBlock,
				StakingInfo:         makeTestPermissionlessStakingInfo(),
				StakingErr:          tc.stakingErr,
			})
			header, _, _ := makeTestKaiaBlock(tc.num)
			header.Rewardbase = tc.rewardbase
			assert.ErrorIs(t, r.VerifyHeader(header, nil), tc.expected)
		})
	}

	t.Run("post-fork, staking info error", func(t *testing.T) {
		r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{
			Proposer:            node1,
			PermissionlessBlock: testPermissionlessBlock,
			StakingErr:          errors.New("no state"),
		})
		header, _, _ := makeTestKaiaBlock(300)
		header.Rewardbase = reward1
		assert.ErrorIs(t, r.VerifyHeader(header, nil), reward.ErrInvalidRewardbase)
	})
}

// A block producer that is in the committee but not reward-eligible at the source block
// must write its AddressBookV2 reward address, not its --rewardbase, after the fork.
func TestFinalizeState_Permissionless_Rewardbase(t *testing.T) {
	var (
		node1   = common.HexToAddress("0xa01")
		reward1 = common.HexToAddress("0xc01")
		local   = common.HexToAddress("0xfff")
	)
	committee := []common.Address{node1, common.HexToAddress("0xa02"), common.HexToAddress("0xa03"), common.HexToAddress("0xa04"), promotedNode}

	testcases := []struct {
		name     string
		num      int64
		node     common.Address
		expected common.Address
	}{
		{"pre-fork, node not in staking info keeps local", 199, promotedNode, local},
		{"pre-fork, eligible node", 199, node1, reward1},
		{"post-fork, eligible node", 200, node1, reward1},
		{"post-fork, promoted node", 200, promotedNode, promotedPD},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			header, txs, receipts := makeTestKaiaBlock(tc.num)
			r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{
				Deferred:            true,
				Proposer:            tc.node,
				PermissionlessBlock: testPermissionlessBlock,
				StakingInfo:         makeTestPermissionlessStakingInfo(),
				Qualified:           committee,
				Header:              header,
				Txs:                 txs,
				Receipts:            receipts,
			})
			r.NodeAddress = tc.node
			r.Rewardbase = local
			require.NoError(t, r.PrepareHeader(header))

			statedb, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
			require.NoError(t, err)
			require.NoError(t, r.FinalizeState(header, statedb, txs, receipts))
			assert.Equal(t, tc.expected, header.Rewardbase)

			// The header produced by FinalizeState passes VerifyHeader.
			assert.NoError(t, r.VerifyHeader(header, nil))

			// The proposer share is paid to header.Rewardbase, and the payouts add up.
			spec, err := r.GetDeferredReward(header, txs, receipts)
			require.NoError(t, err)
			sanityCheckRewardSpec(t, spec, tc.name)
			assert.Equal(t, spec.Rewards[tc.expected], statedb.GetBalance(tc.expected))
			assert.True(t, statedb.GetBalance(tc.expected).Cmp(spec.Proposer) >= 0)
			if tc.expected != local {
				assert.Zero(t, statedb.GetBalance(local).Sign())
			}
		})
	}
}
