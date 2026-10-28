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

package backend

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProposerRewardbaseIsFinalBeforeExecution builds a block the way a proposer does
// (PrepareHeader -> ApplyTransaction -> FinalizeState -> seal) on a chain that pays
// transaction fees immediately (reward.deferredtxfee=false) under WeightedRandom, with the
// node's configured --rewardbase different from the reward address registered for it in
// the staking info. The fee is credited to header.Rewardbase while the transaction runs, so
// the header value must already be the registered address at that point. Re-importing the
// sealed block re-executes it against its own header; if the proposer executed with a
// different address than it sealed, the state roots disagree and InsertChain rejects the
// block, which is the chain halt observed on homi-generated non-deferred networks at the
// Permissionless fork (every validator rejects every proposal with "invalid merkle root").
//
// The scenario runs without and with the Permissionless fork: before it the whole fee is
// credited, since it half of the fee is burnt (see StateTransition.TransitionDb).
//
// The intermediate header checks are non-fatal on purpose: on the code before this fix they
// report the configured address, and the test then goes on to show the consequence, the
// re-import being rejected with "invalid merkle root".
func TestProposerRewardbaseIsFinalBeforeExecution(t *testing.T) {
	for _, permissionless := range []bool{false, true} {
		t.Run(fmt.Sprintf("permissionless=%v", permissionless), func(t *testing.T) {
			testProposerRewardbaseIsFinalBeforeExecution(t, permissionless)
		})
	}
}

func testProposerRewardbaseIsFinalBeforeExecution(t *testing.T, permissionless bool) {
	const (
		validatorCount = 4
		minStake       = uint64(5_000_000)
	)

	setNodeKeys(validatorCount, nil)
	// Every validator is registered with a reward address that differs from its node
	// address and from any configured rewardbase (makeTestStakingInfo draws fresh keys).
	stakingInfo := makeTestStakingInfo([]uint64{minStake, minStake, minStake, minStake}, 0)
	localNode, registered := stakingInfo.NodeIds[0], stakingInfo.RewardAddrs[0] // node 0 is this backend (engine.address)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mStaking := staking_mock.NewMockStakingModule(ctrl)
	mStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(stakingInfo, nil).AnyTimes()

	// Magma (and Kore) are active from genesis, so the non-deferred fee goes to
	// header.Rewardbase (before Magma it went to Coinbase). Minting is disabled so the
	// registered address receives nothing but the fee credit.
	config := testRandaoConfig.Copy()
	config.PermissionlessCompatibleBlock = nil
	if permissionless {
		config.PermissionlessCompatibleBlock = big.NewInt(0)
	}
	config.Istanbul.ProposerPolicy = params.WeightedRandom
	config.Governance.Reward.DeferredTxFee = false
	config.Governance.Reward.MintingAmount = big.NewInt(0)
	config.Governance.Reward.MinimumStake = new(big.Int).SetUint64(minStake)
	config.Governance.Reward.StakingUpdateInterval = 1
	config.Governance.Reward.ProposerUpdateInterval = 1

	senderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	sender := crypto.PubkeyToAddress(senderKey.PublicKey)
	configured := common.HexToAddress("0x00000000000000000000000000000000c0ffee01")

	chain, engine := newBlockChain(t, validatorCount, config, mStaking, blockPeriod(0),
		rewardbase(configured),
		blockchain.GenesisAlloc{sender: {Balance: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000))}},
	)
	defer chain.Stop()
	defer engine.Stop()

	require.Equal(t, engine.address, localNode, "test setup: node 0 must be the local node")
	require.NotEqual(t, configured, registered, "test setup: the configured and registered addresses must differ")

	// 1. Prepare the header as the worker does (work/worker.go commitNewWork).
	parent := chain.CurrentBlock()
	header := makeHeader(parent, chain.Config())
	require.NoError(t, chain.PrepareHeader(header))
	// Keep the block one second after its parent regardless of the process-wide block period, so
	// the import-time timestamp check (parent.Time + BlockGenerationInterval <= header.Time) holds.
	header.Time = new(big.Int).Add(parent.Time(), common.Big1)
	assert.Equal(t, registered, header.Rewardbase,
		"PrepareHeader must decide the registered reward address before any transaction is executed")

	// 2. Execute one fee-paying transfer against the prepared header (work/worker.go commitTransaction).
	state, err := chain.StateAt(parent.Root())
	require.NoError(t, err)
	require.NoError(t, chain.Processor().InitializeState(header, state))

	recipient := common.HexToAddress("0x000000000000000000000000000000000000dead")
	gasPrice := big.NewInt(25e9)
	if header.BaseFee != nil {
		gasPrice = new(big.Int).Set(header.BaseFee)
	}
	tx, err := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
		AccountNonce: 0,
		Price:        gasPrice,
		GasLimit:     params.TxGas,
		Recipient:    &recipient,
		Amount:       big.NewInt(1),
	}), types.LatestSignerForChainID(chain.Config().ChainID), senderKey)
	require.NoError(t, err)

	author := engine.address
	receipt, _, err := chain.ApplyTransaction(chain.Config(), &author, state, header, tx, &header.GasUsed, &vm.Config{})
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	fee := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), gasPrice)
	credited := new(big.Int).Set(fee)
	if chain.Config().Rules(header.Number).IsPermissionless {
		credited.Sub(credited, new(big.Int).Div(fee, big.NewInt(2)))
	}
	require.Equal(t, permissionless, chain.Config().Rules(header.Number).IsPermissionless, "test setup: the burn follows the Permissionless fork")

	// 3. Finalize and seal. The header the proposer seals must be the one it executed with.
	block, err := chain.Processor().FinalizeState(header, state, []*types.Transaction{tx}, []*types.Receipt{receipt})
	require.NoError(t, err)
	assert.Equal(t, registered, block.Header().Rewardbase, "FinalizeState must not change the reward address after execution")
	sealed := sealBlock(engine, block)

	// 4. Import re-executes the block from the parent state with the sealed header, exactly
	// as every validator does. A header/execution mismatch fails here with "invalid merkle root".
	_, err = chain.InsertChain(types.Blocks{sealed})
	require.NoError(t, err, "a block produced by this node must be replayable from its own header")
	require.Equal(t, sealed.Hash(), chain.CurrentBlock().Hash())

	// 5. The fee credit went to the sealed reward address and nowhere else.
	imported, err := chain.StateAt(chain.CurrentBlock().Root())
	require.NoError(t, err)
	require.Zero(t, credited.Cmp(imported.GetBalance(registered)), "the registered reward address receives the non-deferred fee credit (whole fee before Permissionless, fee - fee/2 since)")
	assert.Zero(t, imported.GetBalance(configured).Sign(), "the configured --rewardbase receives nothing")
}
