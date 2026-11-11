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

package tests

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/stretchr/testify/require"
)

// TestRewardbaseAfterRestartWithoutHeadState restarts a WeightedRandom, non-deferred CN whose head
// state is not on disk and checks that it resumes producing blocks with the staking-registered reward
// address. Since the Kaia hardfork, GetStakingInfo(num) reads the state of num-1 and has no DB
// fallback, and RewardModule.PrepareHeader depends on that read. On startup NewBlockChain rewinds the
// head to the newest block whose state exists ("Head state missing, repairing chain"), so the read
// for the next block succeeds. The harness runs in archive mode (every block's state is committed and
// BlockChain.Stop does not commit again), so "stopped without committing the head state" is simulated
// by deleting the head state root before the stop.
func TestRewardbaseAfterRestartWithoutHeadState(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		forkNum    = big.NewInt(2) // Kaia fork: from here the staking info of block N comes from the state of N-1
		nodeId     = bind.NewKeyedTransactor(deriveTestAccount(0)).From
		owner      = bind.NewKeyedTransactor(deriveTestAccount(5))
		config     = testStateReexec_config(forkNum) // WeightedRandom, AddressBookMock
		alloc      = testStateReexec_alloc(t, owner, nodeId)
		registered = common.HexToAddress("0x2000") // testStateReexec_alloc registers (nodeId, 0x1000, 0x2000)
	)
	config.Governance.Reward.DeferredTxFee = false // the fee is credited to header.Rewardbase during execution

	ctx, err := newBlockchainTestContext(&blockchainTestOverrides{
		numNodes:    1,
		numAccounts: 8,
		config:      config,
		alloc:       alloc,
	})
	require.NoError(t, err)
	require.NoError(t, ctx.Start())
	defer ctx.Cleanup()

	ctx.WaitBlock(t, 5)
	bc := ctx.nodes[0].cn.BlockChain().(*blockchain.BlockChain)
	head := bc.CurrentBlock()
	require.NotEqual(t, registered, ctx.accountAddrs[0], "test setup: the configured rewardbase must differ from the registered one")
	require.Equal(t, registered, head.Header().Rewardbase, "PrepareHeader uses the registered reward address")

	// Lose the head state, then stop and restart.
	ctx.nodes[0].cn.ChainDB().DeleteTrieNode(head.Root().ExtendZero())
	require.NoError(t, ctx.Restart())

	bc = ctx.nodes[0].cn.BlockChain().(*blockchain.BlockChain)
	rewound := bc.CurrentBlock()
	require.Less(t, rewound.NumberU64(), head.NumberU64(), "NewBlockChain rewinds the head to the newest block with state")
	_, err = bc.StateAt(rewound.Root())
	require.NoError(t, err, "the rewound head has a readable state")

	// The node resumes on the rewound head; PrepareHeader resolves the registered address from that state.
	next := rewound.NumberU64() + 1
	ctx.WaitBlock(t, next+1)
	for n := next; n <= next+1; n++ {
		require.Equal(t, registered, bc.GetBlockByNumber(n).Header().Rewardbase, "block %d", n)
	}
}
