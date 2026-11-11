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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/require"
)

// TestAccountMapAdjustFeesHalvesRewardSincePermissionless pins the AccountMap half-burn mirror
// (creditedFee via AdjustFeesByReceipts) with Permissionless enabled. No Forks entry enables Permissionless,
// so this branch is otherwise never executed in ./tests. The rewardbase must end at creditedFee(actual),
// i.e. the adjustment is creditedFee(actual) - creditedFee(estimated): neither delta/2 nor creditedFee(delta)
// gives that for the odd/even combinations below (gasPrice 1, so fee == gas).
func TestAccountMapAdjustFeesHalvesRewardSincePermissionless(t *testing.T) {
	var (
		payer      = common.HexToAddress("0x1")
		rewardbase = common.HexToAddress("0xbeef")
		gasPrice   = big.NewInt(1)
	)
	for _, tc := range []struct {
		desc           string
		permissionless bool
		estimatedFee   int64
		gasUsed        uint64
		wantReward     int64
	}{
		{"before permissionless: rewardbase ends at the whole actual fee", false, 21000, 21001, 21001},
		{"since permissionless: odd estimate, even actual (credit(delta) would add 1 kei too much)", true, 21001, 21002, 10501},
		{"since permissionless: even estimate, odd actual (delta/2 would add 1 kei too little)", true, 21000, 21001, 10501},
	} {
		config := Forks["Osaka"].Copy()
		if tc.permissionless {
			config.PermissionlessCompatibleBlock = big.NewInt(0)
		}
		a := NewAccountMap()
		a.chainConfig = config
		a.rewardbase = rewardbase
		a.Set(payer, big.NewInt(1e9), 0)
		estimated := big.NewInt(tc.estimatedFee)
		a.Set(rewardbase, a.creditedFee(estimated, 1), 0) // what Update credited for the estimate

		tx := types.NewTransaction(0, common.HexToAddress("0x2"), big.NewInt(0), 30000, gasPrice, nil)
		a.txFeeCache[tx.Hash()] = &txFeeInfo{estimatedFee: estimated, from: payer, feePayer: payer}
		receipts := types.Receipts{{GasUsed: tc.gasUsed}}

		require.NoError(t, a.AdjustFeesByReceipts(types.Transactions{tx}, receipts, 1), tc.desc)
		require.Zero(t, big.NewInt(tc.wantReward).Cmp(a.m[rewardbase].balance), "%s: want %d got %s", tc.desc, tc.wantReward, a.m[rewardbase].balance)
	}
}
