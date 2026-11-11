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

package blockchain

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

// TestNonDeferredFeeAcrossPermissionless generates a two-block chain in non-deferred mode
// (reward.deferredtxfee=false) with the Permissionless hardfork at block 2 and checks the fee
// credited to Rewardbase at the exact boundary: the whole fee at block 1 (the v2.2.2 behavior that
// must stay replayable, #1097) and the remainder after the half burn at block 2 (#816).
// GenerateChain does not run the reward module, so the Rewardbase balance moves by the fee only.
func TestNonDeferredFeeAcrossPermissionless(t *testing.T) {
	config := params.TestKaiaConfig("osaka")
	config.PermissionlessCompatibleBlock = big.NewInt(2)
	require.NoError(t, fork.SetHardForkBlockNumberConfig(config))
	t.Cleanup(func() { require.NoError(t, fork.SetHardForkBlockNumberConfig(params.TestChainConfig)) })
	require.False(t, config.Governance.DeferredTxFee(), "non-deferred mode expected")
	require.False(t, config.Rules(big.NewInt(1)).IsPermissionless)
	require.True(t, config.Rules(big.NewInt(2)).IsPermissionless)

	var (
		key, _     = crypto.GenerateKey()
		sender     = crypto.PubkeyToAddress(key.PublicKey)
		to         = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		rewardbase = common.HexToAddress("0x000000000000000000000000000000000000beef")
		db         = database.NewMemoryDBManager()
		gspec      = &Genesis{
			Config: config,
			Alloc:  GenesisAlloc{sender: {Balance: new(big.Int).SetUint64(params.KAIA)}},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(config.ChainID)
	)

	blocks, receipts := GenerateChain(config, genesis, faker.NewFaker(), db, 2, func(i int, b *BlockGen) {
		b.SetRewardbase(rewardbase)
		tx, err := types.SignTx(types.NewTransaction(uint64(i), to, big.NewInt(1000), 21000, big.NewInt(1), nil), signer, key)
		require.NoError(t, err)
		b.AddTx(tx)
	})
	require.Len(t, blocks, 2)
	importGenerated(t, config, gspec, blocks)

	balanceAt := func(block *types.Block) *big.Int {
		statedb, err := state.New(block.Root(), state.NewDatabase(db), nil, nil)
		require.NoError(t, err)
		return statedb.GetBalance(rewardbase)
	}
	feeOf := func(i int) *big.Int {
		require.Len(t, receipts[i], 1)
		tx := blocks[i].Transactions()[0]
		return new(big.Int).Mul(new(big.Int).SetUint64(receipts[i][0].GasUsed), tx.EffectiveGasPrice(blocks[i].Header(), config))
	}

	fee1, fee2 := feeOf(0), feeOf(1)
	require.Positive(t, fee1.Sign())
	credited1 := balanceAt(blocks[0])
	credited2 := new(big.Int).Sub(balanceAt(blocks[1]), credited1)

	// Block 1 (before Permissionless): the whole fee is credited.
	require.Zero(t, fee1.Cmp(credited1), "block 1: want %s got %s", fee1, credited1)
	// Block 2 (Permissionless): half the fee is burnt, the rest is credited (fee - fee/2).
	want2 := new(big.Int).Sub(fee2, new(big.Int).Div(fee2, big.NewInt(2)))
	require.Zero(t, want2.Cmp(credited2), "block 2: want %s got %s", want2, credited2)
}

// TestNonDeferredFeeOddRounding pins the rounding direction with an odd fee: the call target holds
// [JUMPDEST STOP], so gasUsed = 21001 and the fee is odd (gasPrice = 1). Before Permissionless the whole fee is credited; since Permissionless the burnt half is
// rounded down (fee/2) and the proposer receives fee - fee/2 (the larger half), matching TransitionDb.
func TestNonDeferredFeeOddRounding(t *testing.T) {
	config := params.TestKaiaConfig("osaka")
	config.PermissionlessCompatibleBlock = big.NewInt(2)
	require.NoError(t, fork.SetHardForkBlockNumberConfig(config))
	t.Cleanup(func() { require.NoError(t, fork.SetHardForkBlockNumberConfig(params.TestChainConfig)) })

	var (
		key, _     = crypto.GenerateKey()
		sender     = crypto.PubkeyToAddress(key.PublicKey)
		to         = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		rewardbase = common.HexToAddress("0x000000000000000000000000000000000000beef")
		db         = database.NewMemoryDBManager()
		gspec      = &Genesis{
			Config: config,
			Alloc: GenesisAlloc{
				sender: {Balance: new(big.Int).SetUint64(params.KAIA)},
				to:     {Balance: big.NewInt(0), Code: []byte{0x5b, 0x00}}, // JUMPDEST(1 gas), STOP
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(config.ChainID)
	)

	nonce := uint64(0)
	addCall := func(b *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(nonce, to, big.NewInt(0), 30000, big.NewInt(1), nil), signer, key)
		require.NoError(t, err)
		nonce++
		b.AddTx(tx)
	}
	blocks, receipts := GenerateChain(config, genesis, faker.NewFaker(), db, 2, func(i int, b *BlockGen) {
		b.SetRewardbase(rewardbase)
		addCall(b)
		if i == 1 {
			addCall(b) // two odd fees in the Permissionless block: the burn is rounded per transaction
		}
	})
	require.Len(t, blocks, 2)
	importGenerated(t, config, gspec, blocks)

	balanceAt := func(block *types.Block) *big.Int {
		statedb, err := state.New(block.Root(), state.NewDatabase(db), nil, nil)
		require.NoError(t, err)
		return statedb.GetBalance(rewardbase)
	}
	feesOf := func(i int) []*big.Int {
		fees := make([]*big.Int, 0, len(blocks[i].Transactions()))
		for j, tx := range blocks[i].Transactions() {
			require.Equal(t, uint64(21001), receipts[i][j].GasUsed)
			fee := new(big.Int).Mul(new(big.Int).SetUint64(receipts[i][j].GasUsed), tx.EffectiveGasPrice(blocks[i].Header(), config))
			require.Equal(t, int64(1), new(big.Int).Mod(fee, big.NewInt(2)).Int64(), "fee must be odd for this test")
			fees = append(fees, fee)
		}
		return fees
	}
	credit := func(fee *big.Int) *big.Int { return new(big.Int).Sub(fee, new(big.Int).Div(fee, big.NewInt(2))) }

	fees1, fees2 := feesOf(0), feesOf(1)
	require.Len(t, fees1, 1)
	require.Len(t, fees2, 2)

	credited1 := balanceAt(blocks[0])
	credited2 := new(big.Int).Sub(balanceAt(blocks[1]), credited1)
	require.Zero(t, fees1[0].Cmp(credited1), "block 1: whole fee, want %s got %s", fees1[0], credited1)

	// Block 2: each transaction credits fee - fee/2 (the larger half), NOT fee/2, and the rounding is per
	// transaction: with two odd fees the sum of the two credits exceeds the credit of the block total by one kei.
	perTx := new(big.Int).Add(credit(fees2[0]), credit(fees2[1]))
	total := credit(new(big.Int).Add(fees2[0], fees2[1]))
	require.Zero(t, perTx.Cmp(credited2), "block 2: want %s got %s", perTx, credited2)
	require.NotZero(t, total.Cmp(credited2), "block 2: the credit must be rounded per transaction, not on the block total")
}

// TestNonDeferredFeeAcrossPermissionlessFeeDelegated repeats the boundary check with a fee-delegated value
// transfer with a fee ratio. The payer split only changes who pays the fee; the rewardbase still receives the
// whole fee before Permissionless and fee - fee/2 since, because TransitionDb credits gasUsed * gasPrice
// regardless of how buyGas/returnGas split it between the sender and the fee payer.
func TestNonDeferredFeeAcrossPermissionlessFeeDelegated(t *testing.T) {
	config := params.TestKaiaConfig("osaka")
	config.PermissionlessCompatibleBlock = big.NewInt(2)
	require.NoError(t, fork.SetHardForkBlockNumberConfig(config))
	t.Cleanup(func() { require.NoError(t, fork.SetHardForkBlockNumberConfig(params.TestChainConfig)) })

	var (
		senderKey, _ = crypto.GenerateKey()
		payerKey, _  = crypto.GenerateKey()
		sender       = crypto.PubkeyToAddress(senderKey.PublicKey)
		payer        = crypto.PubkeyToAddress(payerKey.PublicKey)
		to           = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		rewardbase   = common.HexToAddress("0x000000000000000000000000000000000000beef")
		ratio        = types.FeeRatio(30)
		db           = database.NewMemoryDBManager()
		gspec        = &Genesis{
			Config: config,
			Alloc: GenesisAlloc{
				sender: {Balance: new(big.Int).SetUint64(params.KAIA)},
				payer:  {Balance: new(big.Int).SetUint64(params.KAIA)},
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(config.ChainID)
	)

	blocks, receipts := GenerateChain(config, genesis, faker.NewFaker(), db, 2, func(i int, b *BlockGen) {
		b.SetRewardbase(rewardbase)
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(i),
			types.TxValueKeyFrom:               sender,
			types.TxValueKeyTo:                 to,
			types.TxValueKeyAmount:             big.NewInt(1000),
			types.TxValueKeyGasLimit:           uint64(100000),
			types.TxValueKeyGasPrice:           big.NewInt(1),
			types.TxValueKeyFeePayer:           payer,
			types.TxValueKeyFeeRatioOfFeePayer: ratio,
		})
		require.NoError(t, err)
		require.NoError(t, tx.SignWithKeys(signer, []*ecdsa.PrivateKey{senderKey}))
		require.NoError(t, tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{payerKey}))
		b.AddTx(tx)
	})
	require.Len(t, blocks, 2)
	importGenerated(t, config, gspec, blocks)

	stateAt := func(block *types.Block) *state.StateDB {
		statedb, err := state.New(block.Root(), state.NewDatabase(db), nil, nil)
		require.NoError(t, err)
		return statedb
	}
	feeOf := func(i int) *big.Int {
		require.Len(t, receipts[i], 1)
		tx := blocks[i].Transactions()[0]
		return new(big.Int).Mul(new(big.Int).SetUint64(receipts[i][0].GasUsed), tx.EffectiveGasPrice(blocks[i].Header(), config))
	}

	genesisState, state1, state2 := stateAt(genesis), stateAt(blocks[0]), stateAt(blocks[1])
	fee1, fee2 := feeOf(0), feeOf(1)
	require.Positive(t, fee1.Sign())

	// The fee payer pays its ratio of the whole fee in both regimes; the burn does not change who pays.
	payerFee1, _ := types.CalcFeeWithRatio(ratio, fee1)
	require.Zero(t, payerFee1.Cmp(new(big.Int).Sub(genesisState.GetBalance(payer), state1.GetBalance(payer))), "block 1 payer")
	payerFee2, _ := types.CalcFeeWithRatio(ratio, fee2)
	require.Zero(t, payerFee2.Cmp(new(big.Int).Sub(state1.GetBalance(payer), state2.GetBalance(payer))), "block 2 payer")

	// Block 1 (before Permissionless): the whole fee is credited; block 2 (Permissionless): fee - fee/2.
	credited1 := state1.GetBalance(rewardbase)
	credited2 := new(big.Int).Sub(state2.GetBalance(rewardbase), credited1)
	require.Zero(t, fee1.Cmp(credited1), "block 1: want %s got %s", fee1, credited1)
	want2 := new(big.Int).Sub(fee2, new(big.Int).Div(fee2, big.NewInt(2)))
	require.Zero(t, want2.Cmp(credited2), "block 2: want %s got %s", want2, credited2)
}

// importGenerated re-executes the generated blocks through InsertChain on a fresh database: #1097 surfaced
// on this path as an import-time "invalid merkle root", not in isolated transaction execution.
func importGenerated(t *testing.T, config *params.ChainConfig, gspec *Genesis, blocks []*types.Block) {
	chaindb := database.NewMemoryDBManager()
	gspec.MustCommit(chaindb)
	chain, err := NewBlockChain(chaindb, nil, config, faker.NewFaker(), vm.Config{})
	require.NoError(t, err)
	defer chain.Stop()
	n, err := chain.InsertChain(blocks)
	require.NoError(t, err, "import failed at block %d", n)
	require.Equal(t, blocks[len(blocks)-1].Hash(), chain.CurrentBlock().Hash())
}
