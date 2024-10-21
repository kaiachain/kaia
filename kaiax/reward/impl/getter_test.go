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
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRewardSummary(t *testing.T) {
	// Test the wrapper
	// full deferred
	header, txs, receipts := makeTestKaiaBlock(61)
	r := makeTestRewardModule(t, false, true, header, txs, receipts)

	spec, err := r.GetRewardSummary(header.Number.Uint64())
	require.Nil(t, err)
	assert.Equal(t, &reward.RewardSummary{
		Minted:   big.NewInt(6.4e18),
		TotalFee: big.NewInt(0.0376e18),
		BurntFee: big.NewInt(0.0376e18),
	}, spec)

	// Test the core logic
	var (
		mintingAmount  = big.NewInt(6.4e18) // 6.40 KAIA
		rewardRatio, _ = reward.NewRewardRatio("50/20/30")
		kip82Ratio, _  = reward.NewRewardKip82Ratio("20/80")
		config         = &reward.RewardConfig{
			Rewardbase:    common.HexToAddress("0xfff"),
			MintingAmount: mintingAmount,
			MinimumStake:  big.NewInt(5_000_000),
			RewardRatio:   rewardRatio,
			Kip82Ratio:    kip82Ratio,
		}
		lowFee  = big.NewInt(7e16) // 0.07 KAIA
		highFee = big.NewInt(2e18) // 2.00 KAIA (F/2 > gpM)
	)
	testcases := []struct {
		desc     string
		simple   bool
		magma    bool
		kore     bool
		totalFee *big.Int
		expected *reward.RewardSummary
	}{
		{
			"simple pre-magma", true, false, false, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(0),
			},
		},
		{
			"simple magma", true, true, false, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(3.5e16),
			},
		},
		{
			"simple kore", true, true, true, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(3.5e16), // Kore logic not used in Simple policy.
			},
		},
		{
			"full pre-magma", false, false, false, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(0),
			},
		},
		{
			"full magma", false, true, false, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(3.5e16),
			},
		},
		{
			"full kore low traffic", false, true, true, lowFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(7e16),
			},
		},
		{
			"full kore high traffic", false, true, true, highFee, &reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(2e18),
				BurntFee: big.NewInt(1.64e18),
			},
		},
	}
	for _, tc := range testcases {
		config.IsSimple = tc.simple
		config.Rules.IsMagma = tc.magma
		config.Rules.IsKore = tc.kore
		spec := getRewardSummary(config, tc.totalFee)
		assert.Equal(t, tc.expected, spec, tc.desc)
	}
}

// TestGetBlockReward tests the "big branches" such as simple vs full, deferred vs non-deferred.
// Other branches such as hardfork and ratios are covered in smaller unit tests.
func TestGetBlockReward(t *testing.T) {
	testcases := []struct {
		desc     string
		simple   bool
		deferred bool
		expected *reward.RewardSpec
	}{
		{
			"simple non-deferred", true, false,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18), // includes non-deferred fee
					BurntFee: big.NewInt(0.0188e18),
				},
				Proposer: big.NewInt(6.4188e18),
				Stakers:  big.NewInt(0),
				KIF:      big.NewInt(0),
				KEF:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(6.4188e18),
				},
			},
		},
		{
			"simple deferred", true, true,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18),
					BurntFee: big.NewInt(0.0188e18),
				},
				Proposer: big.NewInt(6.4188e18),
				Stakers:  big.NewInt(0),
				KIF:      big.NewInt(0),
				KEF:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(6.4188e18),
				},
			},
		},
		{
			"full non-deferred", false, false,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18),
					BurntFee: big.NewInt(0.0188e18),
				},
				Proposer: big.NewInt(0.6588e18 + 1), // gpM + remainder + F/2
				Stakers:  big.NewInt(2.56e18 - 1),   // gsM - remainder
				KIF:      big.NewInt(1.28e18),
				KEF:      big.NewInt(1.92e18),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(0.6588e18 + 1),
					common.HexToAddress("0xd01"): big.NewInt(1.28e18),
					common.HexToAddress("0xd02"): big.NewInt(1.92e18),
					common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
					common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
					common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
				},
			},
		},
		{
			"full deferred", false, true,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18),
					BurntFee: big.NewInt(0.0376e18),
				},
				Proposer: big.NewInt(0.64e18 + 1), // gpM + remainder
				Stakers:  big.NewInt(2.56e18 - 1), // gsM - remainder
				KIF:      big.NewInt(1.28e18),
				KEF:      big.NewInt(1.92e18),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(0.64e18 + 1),
					common.HexToAddress("0xd01"): big.NewInt(1.28e18),
					common.HexToAddress("0xd02"): big.NewInt(1.92e18),
					common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
					common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
					common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
				},
			},
		},
	}
	for _, tc := range testcases {
		header, txs, receipts := makeTestKaiaBlock(61)
		r := makeTestRewardModule(t, tc.simple, tc.deferred, header, txs, receipts)

		spec, err := r.GetBlockReward(header.Number.Uint64())
		require.Nil(t, err)
		assert.Equal(t, tc.expected, spec, tc.desc)
	}
}

// Test the special case of non-deferred mode and before Magma,
// where non-deferred fees are assigned to evm.Coinbase (= block author).
func TestSpecWithNonDeferredFeeAuthor(t *testing.T) {
	header, txs, receipts := makeTestPreMagmaBlock(1)
	r := makeTestRewardModule(t, true, false, header, txs, receipts)
	author, _ := r.Chain.Engine().Author(header)

	expected := &reward.RewardSpec{
		RewardSummary: reward.RewardSummary{
			Minted:   big.NewInt(6.4e18),
			TotalFee: big.NewInt(0.05e18), // Simple & Pre-Magma: F (non-deferred) + F (deferred)
			BurntFee: big.NewInt(0),
		},
		Proposer: big.NewInt(6.45e18), // Simple & Pre-Magma: F (non-deferred) + M + F (deferred)
		Stakers:  big.NewInt(0),
		KIF:      big.NewInt(0),
		KEF:      big.NewInt(0),
		Rewards: map[common.Address]*big.Int{
			author:                       big.NewInt(0.025e18), // F (non-deferred) to Coinbase (author)
			common.HexToAddress("0xfff"): big.NewInt(6.425e18), // M + F (deferred) to Rewardbase
		},
	}

	spec, err := r.GetBlockReward(header.Number.Uint64())
	require.Nil(t, err)
	assert.Equal(t, expected, spec)
}

// TestGetDeferredReward tests the "big branches" such as simple vs full, deferred vs non-deferred.
// Other branches such as hardfork and ratios are covered in smaller unit tests.
func TestGetDeferredReward(t *testing.T) {
	testcases := []struct {
		desc     string
		simple   bool
		deferred bool
		expected *reward.RewardSpec
	}{
		{
			"simple non-deferred", true, false,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0), // no deferred fee
					BurntFee: big.NewInt(0),
				},
				Proposer: big.NewInt(6.4e18),
				Stakers:  big.NewInt(0),
				KIF:      big.NewInt(0),
				KEF:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(6.4e18),
				},
			},
		},
		{
			"simple deferred", true, true,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18),
					BurntFee: big.NewInt(0.0188e18),
				},
				Proposer: big.NewInt(6.4188e18),
				Stakers:  big.NewInt(0),
				KIF:      big.NewInt(0),
				KEF:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(6.4188e18),
				},
			},
		},
		{
			"full non-deferred", false, false,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0),
					BurntFee: big.NewInt(0),
				},
				Proposer: big.NewInt(0.64e18 + 1), // gpM + remainder
				Stakers:  big.NewInt(2.56e18 - 1), // gsM - remainder
				KIF:      big.NewInt(1.28e18),
				KEF:      big.NewInt(1.92e18),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(0.64e18 + 1),
					common.HexToAddress("0xd01"): big.NewInt(1.28e18),
					common.HexToAddress("0xd02"): big.NewInt(1.92e18),
					common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
					common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
					common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
				},
			},
		},
		{
			"full deferred", false, true,
			&reward.RewardSpec{
				RewardSummary: reward.RewardSummary{
					Minted:   big.NewInt(6.4e18),
					TotalFee: big.NewInt(0.0376e18),
					BurntFee: big.NewInt(0.0376e18),
				},
				Proposer: big.NewInt(0.64e18 + 1), // gpM + remainder
				Stakers:  big.NewInt(2.56e18 - 1), // gsM - remainder
				KIF:      big.NewInt(1.28e18),
				KEF:      big.NewInt(1.92e18),
				Rewards: map[common.Address]*big.Int{
					common.HexToAddress("0xfff"): big.NewInt(0.64e18 + 1),
					common.HexToAddress("0xd01"): big.NewInt(1.28e18),
					common.HexToAddress("0xd02"): big.NewInt(1.92e18),
					common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
					common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
					common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
				},
			},
		},
	}
	for _, tc := range testcases {
		header, txs, receipts := makeTestKaiaBlock(61)
		r := makeTestRewardModule(t, tc.simple, tc.deferred, header, txs, receipts)

		spec, err := r.GetDeferredReward(header, txs, receipts)
		require.Nil(t, err)
		assert.Equal(t, tc.expected, spec, tc.desc)
	}
}

func TestGetTotalFee(t *testing.T) {
	testcases := []struct {
		desc     string
		config   *reward.RewardConfig
		header   *types.Header
		txs      []*types.Transaction
		receipts []*types.Receipt
		expected *big.Int
	}{
		{
			"pre-magma",
			&reward.RewardConfig{UnitPrice: big.NewInt(25e9)},
			&types.Header{BaseFee: nil, GasUsed: 1_000_000},
			[]*types.Transaction{
				makeTestTx_type0(25e9, 200_000),
				makeTestTx_type0(25e9, 400_000),
				makeTestTx_type0(25e9, 600_000),
				makeTestTx_type0(25e9, 800_000),
			},
			[]*types.Receipt{
				{GasUsed: 100_000},
				{GasUsed: 200_000},
				{GasUsed: 300_000},
				{GasUsed: 400_000},
			},
			big.NewInt(0.025e18), // unitPrice * gasUsed
		},
		{
			"magma",
			&reward.RewardConfig{Rules: params.Rules{IsMagma: true}, UnitPrice: big.NewInt(25e9)},
			&types.Header{BaseFee: big.NewInt(27e9), GasUsed: 1_000_000},
			[]*types.Transaction{
				makeTestTx_type0(27e9, 200_000),
				makeTestTx_type0(28e9, 400_000),
				makeTestTx_type0(29e9, 600_000),
				makeTestTx_type0(30e9, 800_000),
			},
			[]*types.Receipt{
				{GasUsed: 100_000},
				{GasUsed: 200_000},
				{GasUsed: 300_000},
				{GasUsed: 400_000},
			},
			big.NewInt(0.027e18), // baseFee * gasUsed
		},
		{
			"kaia",
			&reward.RewardConfig{Rules: params.Rules{IsMagma: true, IsKaia: true}, UnitPrice: big.NewInt(25e9)},
			&types.Header{BaseFee: big.NewInt(27e9), GasUsed: 1_000_000},
			[]*types.Transaction{
				makeTestTx_type2(50e9, 1e9, 200_000), // effectivePrice = 28e9, effectiveTip = 1e9
				makeTestTx_type2(29e9, 7e9, 400_000), // effectivePrice = 29e9, effectiveTip = 7e9
				makeTestTx_type0(30e9, 600_000),      // effectivePrice = 30e9, effectiveTip = 3e9
				makeTestTx_type0(50e9, 800_000),      // effectivePrice = 50e9, effectiveTip = 33e9
			},
			[]*types.Receipt{
				{GasUsed: 100_000},
				{GasUsed: 200_000},
				{GasUsed: 300_000},
				{GasUsed: 400_000},
			},
			big.NewInt(0.0376e18), // sum{ (effectiveGasTip[i] + baseFee) * gasUsed[i] }
		},
	}
	for _, tc := range testcases {
		totalFee, err := getTotalFee(tc.config, tc.header, tc.txs, tc.receipts)
		require.Nil(t, err)
		assert.Equal(t, tc.expected, totalFee, tc.desc)
	}
}

func TestGetDeferredRewardSimple(t *testing.T) {
	var (
		mintingAmount = big.NewInt(6.4e18) // 6.40 KAIA
		config        = &reward.RewardConfig{
			Rewardbase:    common.HexToAddress("0xfff"),
			MintingAmount: mintingAmount,
		}
		totalFee = big.NewInt(7e16) // 0.07 KAIA
	)
	testcases := []struct {
		desc     string
		deferred bool
		magma    bool
		expected *reward.RewardSpec
	}{
		{"non-deferred pre-magma", false, false, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(0),
			},
			Proposer: big.NewInt(6.47e18), // Buggy behavior; M+F despite being non-deferred.
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(6.47e18),
			},
		}},
		{"non-deferred magma", false, true, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(0),
				BurntFee: big.NewInt(0),
			},
			Proposer: big.NewInt(6.4e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(6.4e18),
			},
		}},
		{"deferred pre-magma", true, false, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(0),
			},
			Proposer: big.NewInt(6.47e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(6.47e18),
			},
		}},
		{"deferred magma", true, true, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(3.5e16),
			},
			Proposer: big.NewInt(6.435e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(6.435e18),
			},
		}},
	}
	for _, tc := range testcases {
		config.DeferredTxFee = tc.deferred
		config.Rules.IsMagma = tc.magma

		spec, err := getDeferredRewardSimple(config, totalFee)
		require.Nil(t, err, tc.desc)
		assert.Equal(t, tc.expected, spec, tc.desc)
		sanityCheckRewardSpec(t, spec, tc.desc)
	}
}

func TestGetDeferredRewardFull(t *testing.T) {
	var (
		mintingAmount  = big.NewInt(6.4e18) // 6.40 KAIA
		rewardRatio, _ = reward.NewRewardRatio("50/20/30")
		kip82Ratio, _  = reward.NewRewardKip82Ratio("20/80")
		config         = &reward.RewardConfig{
			Rewardbase:    common.HexToAddress("0xfff"),
			MintingAmount: mintingAmount,
			MinimumStake:  big.NewInt(5_000_000),
			RewardRatio:   rewardRatio,
			Kip82Ratio:    kip82Ratio,
		}
		lowFee  = big.NewInt(7e16) // 0.07 KAIA
		highFee = big.NewInt(2e18) // 2.00 KAIA (F/2 > gpM)
	)
	testcases := []struct {
		desc     string
		deferred bool
		magma    bool
		kore     bool
		totalFee *big.Int
		expected *reward.RewardSpec
	}{
		{"non-deferred", false, false, false, lowFee, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(0),
				BurntFee: big.NewInt(0),
			},
			Proposer: big.NewInt(3.2e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(1.28e18),
			KEF:      big.NewInt(1.92e18),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(3.2e18),
				common.HexToAddress("0xd01"): big.NewInt(1.28e18),
				common.HexToAddress("0xd02"): big.NewInt(1.92e18),
			},
		}},
		{"pre-magma", true, false, false, lowFee, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(0),
			},
			Proposer: big.NewInt(3.235e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(1.294e18),
			KEF:      big.NewInt(1.941e18),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(3.235e18),
				common.HexToAddress("0xd01"): big.NewInt(1.294e18),
				common.HexToAddress("0xd02"): big.NewInt(1.941e18),
			},
		}},
		{"magma", true, true, false, lowFee, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(3.5e16),
			},
			Proposer: big.NewInt(3.2175e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(1.287e18),
			KEF:      big.NewInt(1.9305e18),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(3.2175e18),
				common.HexToAddress("0xd01"): big.NewInt(1.287e18),
				common.HexToAddress("0xd02"): big.NewInt(1.9305e18),
			},
		}},
		{"kore low traffic", true, true, true, lowFee, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(7e16),
				BurntFee: big.NewInt(7e16), // F
			},
			Proposer: big.NewInt(0.64e18 + 1), // gpM + remainder
			Stakers:  big.NewInt(2.56e18 - 1), // gsM - remainder
			KIF:      big.NewInt(1.28e18),
			KEF:      big.NewInt(1.92e18),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(0.64e18 + 1),
				common.HexToAddress("0xd01"): big.NewInt(1.28e18),
				common.HexToAddress("0xd02"): big.NewInt(1.92e18),
				common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
				common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
				common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
			},
		}},
		{"kore high traffic", true, true, true, highFee, &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(6.4e18),
				TotalFee: big.NewInt(2.00e18),
				BurntFee: big.NewInt(1.64e18), // F/2 + gpM
			},
			Proposer: big.NewInt(1.00e18 + 1), // 0.64 (gpM) + 0.36 (F/2 - gpM) + remainder
			Stakers:  big.NewInt(2.56e18 - 1), // gsM - remainder
			KIF:      big.NewInt(1.28e18),
			KEF:      big.NewInt(1.92e18),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(1.00e18 + 1),
				common.HexToAddress("0xd01"): big.NewInt(1.28e18),
				common.HexToAddress("0xd02"): big.NewInt(1.92e18),
				common.HexToAddress("0xc01"): big.NewInt(426666666666666666),
				common.HexToAddress("0xc02"): big.NewInt(853333333333333333),
				common.HexToAddress("0xc03"): big.NewInt(1280000000000000000),
			},
		}},
	}
	for _, tc := range testcases {
		config.DeferredTxFee = tc.deferred
		config.Rules.IsMagma = tc.magma
		config.Rules.IsKore = tc.kore
		si := makeTestStakingInfo([]uint64{5_000_001, 5_000_002, 5_000_003, 5_000_000}) // 1:2:3:0

		spec, err := getDeferredRewardFull(config, tc.totalFee, si)
		require.Nil(t, err, tc.desc)
		assert.Equal(t, tc.expected, spec, tc.desc)
		sanityCheckRewardSpec(t, spec, tc.desc)
	}
}

func TestGetBurnAmountMagma(t *testing.T) {
	assert.Equal(t, big.NewInt(1234), getBurnAmountMagma(big.NewInt(2468)))
}

func TestGetBurnAmountKore(t *testing.T) {
	var (
		mintingAmount  = big.NewInt(1000)
		rewardRatio, _ = reward.NewRewardRatio("50/20/30")
		kip82Ratio, _  = reward.NewRewardKip82Ratio("20/80")
		config         = &reward.RewardConfig{MintingAmount: mintingAmount, RewardRatio: rewardRatio, Kip82Ratio: kip82Ratio}
	)

	testcases := []struct {
		totalFee *big.Int
		expected *big.Int
	}{
		// burnt = F/2 + min(F/2, gpM) where gpM (Proposer's minting reward) = 0.5*0.2*1000 = 100

		// F/2 <= gpM, or F <= 200. burnt = F.
		{big.NewInt(0), big.NewInt(0)},
		{big.NewInt(100), big.NewInt(100)},
		{big.NewInt(200), big.NewInt(200)},

		// F/2 > gpM, or F > 200. burnt = F/2 + 100.
		{big.NewInt(201), big.NewInt(200)}, // odd number case
		{big.NewInt(300), big.NewInt(250)},
	}
	for _, tc := range testcases {
		assert.Equal(t, tc.expected, getBurnAmountKore(config, tc.totalFee), tc.totalFee.String())
	}
}

func TestAssignStakingRewards(t *testing.T) {
	var (
		min    = uint64(5_000_000)
		config = &reward.RewardConfig{MinimumStake: big.NewInt(int64(min))}
		reward = big.NewInt(1e18)
	)
	testcases := []struct {
		desc              string
		stakingAmounts    []uint64                    // in KAIA
		expectedAlloc     map[common.Address]*big.Int // in kei
		expectedRemainder *big.Int                    // in kei
	}{
		{
			"no one eligible",
			[]uint64{0, 1, min - 1, min},
			map[common.Address]*big.Int{},
			big.NewInt(1e18),
		},
		{
			"only one eligible",
			[]uint64{min + 1, min, min, min},
			map[common.Address]*big.Int{
				common.HexToAddress("0xc01"): big.NewInt(1e18),
			},
			big.NewInt(0),
		},
		{
			"no remainder",
			[]uint64{min + 1, min + 2, min + 3, min + 4},
			map[common.Address]*big.Int{
				common.HexToAddress("0xc01"): big.NewInt(1e17),
				common.HexToAddress("0xc02"): big.NewInt(2e17),
				common.HexToAddress("0xc03"): big.NewInt(3e17),
				common.HexToAddress("0xc04"): big.NewInt(4e17),
			},
			big.NewInt(0),
		},
		{
			"remainder",
			[]uint64{min + 1000, min + 2000, min + 4000, 0},
			map[common.Address]*big.Int{
				common.HexToAddress("0xc01"): big.NewInt(142857142857142857),
				common.HexToAddress("0xc02"): big.NewInt(285714285714285714),
				common.HexToAddress("0xc03"): big.NewInt(571428571428571428),
			},
			big.NewInt(1),
		},
	}
	for _, tc := range testcases {
		si := makeTestStakingInfo(tc.stakingAmounts)
		alloc, remainder := assignStakingRewards(config, reward, si)

		assert.Equal(t, tc.expectedAlloc, alloc, tc.desc)
		assert.Equal(t, tc.expectedRemainder.String(), remainder.String(), tc.desc)
	}
}

func TestSpecWithProposerAndFunds(t *testing.T) {
	var (
		zeroAddr   = common.Address{}
		kifAddr    = common.HexToAddress("0xd01")
		kefAddr    = common.HexToAddress("0xd02")
		rewardbase = common.HexToAddress("0xfff")

		proposer = int64(500)
		kif      = int64(200)
		kef      = int64(300)

		config = &reward.RewardConfig{Rewardbase: rewardbase}
	)

	testcases := []struct {
		kifAddr          common.Address
		kefAddr          common.Address
		expectedProposer *big.Int
		expectedKIF      *big.Int
		expectedKEF      *big.Int
		expectedRewards  map[common.Address]*big.Int
	}{
		{zeroAddr, zeroAddr, big.NewInt(proposer + kif + kef), big.NewInt(0), big.NewInt(0), map[common.Address]*big.Int{
			rewardbase: big.NewInt(proposer + kif + kef),
		}},
		{kifAddr, zeroAddr, big.NewInt(kef + proposer), big.NewInt(kif), big.NewInt(0), map[common.Address]*big.Int{
			rewardbase: big.NewInt(proposer + kef),
			kifAddr:    big.NewInt(kif),
		}},
		{zeroAddr, kefAddr, big.NewInt(proposer + kif), big.NewInt(0), big.NewInt(kef), map[common.Address]*big.Int{
			rewardbase: big.NewInt(proposer + kif),
			kefAddr:    big.NewInt(kef),
		}},
		{kifAddr, kefAddr, big.NewInt(proposer), big.NewInt(kif), big.NewInt(kef), map[common.Address]*big.Int{
			rewardbase: big.NewInt(proposer),
			kifAddr:    big.NewInt(kif),
			kefAddr:    big.NewInt(kef),
		}},
		// KIF and KEF addresses are the same.
		{kifAddr, kifAddr, big.NewInt(proposer), big.NewInt(kif), big.NewInt(kef), map[common.Address]*big.Int{
			rewardbase: big.NewInt(proposer),
			kifAddr:    big.NewInt(kif + kef),
		}},
	}
	for i, tc := range testcases {
		spec := reward.NewRewardSpec()
		si := &staking.StakingInfo{KIFAddr: tc.kifAddr, KEFAddr: tc.kefAddr}
		spec = specWithProposerAndFunds(spec, config, big.NewInt(proposer), big.NewInt(kif), big.NewInt(kef), si)

		assert.Equal(t, tc.expectedKIF, spec.KIF, i)
		assert.Equal(t, tc.expectedKEF, spec.KEF, i)
		assert.Equal(t, tc.expectedProposer, spec.Proposer, i)
		assert.Equal(t, tc.expectedRewards, spec.Rewards, i)
	}
}

func sanityCheckRewardSpec(t *testing.T, spec *reward.RewardSpec, msg interface{}) {
	sumSummary := new(big.Int).Add(spec.Minted, spec.TotalFee)
	sumSummary.Sub(sumSummary, spec.BurntFee)

	sumParts := new(big.Int).Add(spec.Proposer, spec.Stakers)
	sumParts.Add(sumParts, spec.KIF)
	sumParts.Add(sumParts, spec.KEF)

	sumRewards := new(big.Int)
	for _, amount := range spec.Rewards {
		sumRewards.Add(sumRewards, amount)
	}

	assert.Equal(t, sumSummary, sumParts, msg)
	assert.Equal(t, sumSummary, sumRewards, msg)
}
