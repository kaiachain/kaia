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
	"fmt"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FinalizeState distributes the reward to the addresses named by the header it is given and
// never rewrites the header. Until v3.0.0 it replaced header.Rewardbase with the node's
// staking-registered reward address when header.Root was still empty (the proposer's path),
// after the transactions had already been executed with the previous value; that decision now
// belongs to PrepareHeader. The fixture reproduces the old trigger: WeightedRandom, empty Root,
// a node with a staking entry (0xa01 -> 0xc01) and a different configured Rewardbase (0xfff).
func TestFinalizeStateLeavesHeaderUntouched(t *testing.T) {
	var (
		node1      = common.HexToAddress("0xa01") // stakingInfo.NodeIds[0]
		configured = common.HexToAddress("0xfff") // makeTestMagmaBlock's header.Rewardbase
	)

	for _, deferred := range []bool{false, true} {
		t.Run(fmt.Sprintf("deferred=%v", deferred), func(t *testing.T) {
			header, txs, receipts := makeTestMagmaBlock(100)
			require.True(t, common.EmptyHash(header.Root), "fixture: the proposer's header has no state root yet")
			require.Equal(t, configured, header.Rewardbase)

			r := makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{
				Deferred: deferred, NodeAddress: node1, Rewardbase: configured,
				Header: header, Txs: txs, Receipts: receipts,
			})
			statedb, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
			require.NoError(t, err)

			before := types.CopyHeader(header)
			expected, err := r.GetDeferredReward(before, txs, receipts)
			require.NoError(t, err)

			require.NoError(t, r.FinalizeState(header, statedb, txs, receipts))

			assert.Equal(t, before, types.CopyHeader(header), "FinalizeState must not modify the header")
			assert.Equal(t, before.Hash(), header.Hash())
			require.NotEmpty(t, expected.Rewards)
			for addr, amount := range expected.Rewards {
				assert.Equal(t, amount, statedb.GetBalance(addr), "reward of %s", addr.Hex())
			}
			assert.Contains(t, expected.Rewards, configured, "the proposer's share goes to the sealed header's Rewardbase")
		})
	}
}
