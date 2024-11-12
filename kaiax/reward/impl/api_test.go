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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/reward"
	reward_mock "github.com/kaiachain/kaia/kaiax/reward/mock"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAPIGetRewards(t *testing.T) {
	mockCtrl, r, chain, kaiaAPI, _ := makeTestAPI(t)
	defer mockCtrl.Finish()

	var (
		lastNum  = rpc.LatestBlockNumber
		currNum  = rpc.BlockNumber(10)
		otherNum = rpc.BlockNumber(9)

		currBlock = types.NewBlock(&types.Header{Number: big.NewInt(int64(currNum))}, nil, nil)

		spec = &reward.RewardSpec{
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(1e18),
				TotalFee: big.NewInt(0.02e18),
				BurntFee: big.NewInt(0.01e18),
			},
			Proposer: big.NewInt(1.01e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(1.01e18),
			},
		}
	)
	chain.EXPECT().CurrentBlock().Return(currBlock).AnyTimes()
	r.EXPECT().GetBlockReward(uint64(currNum)).Return(spec, nil).AnyTimes()
	r.EXPECT().GetBlockReward(uint64(otherNum)).Return(nil, reward.ErrNoBlock).AnyTimes()

	{
		// Different ways to provide the block number.
		testcases := []*rpc.BlockNumber{nil, &lastNum, &currNum}
		for _, num := range testcases {
			result, err := kaiaAPI.GetRewards(num)
			assert.NoError(t, err)
			assert.Equal(t, spec, result)
		}
	}
	{
		// Failure cases.
		_, err := kaiaAPI.GetRewards(&otherNum)
		assert.ErrorIs(t, err, reward.ErrNoBlock)
	}
}

func TestAPIGetAccumulatedRewards(t *testing.T) {
	mockCtrl, r, chain, _, govAPI := makeTestAPI(t)
	defer mockCtrl.Finish()

	var (
		lastNum  = rpc.LatestBlockNumber
		lowerNum = rpc.BlockNumber(7)
		upperNum = rpc.BlockNumber(10)
		currNum  = rpc.BlockNumber(10)

		lowerTime = time.Unix(1700000007, 0)
		upperTime = time.Unix(1700000010, 0)
		currTime  = time.Unix(1700000010, 0)

		lowerHeader = &types.Header{Number: big.NewInt(int64(lowerNum)), Time: big.NewInt(lowerTime.Unix())}
		upperHeader = &types.Header{Number: big.NewInt(int64(upperNum)), Time: big.NewInt(upperTime.Unix())}
		currBlock   = types.NewBlock(&types.Header{Number: big.NewInt(int64(currNum)), Time: big.NewInt(currTime.Unix())}, nil, nil)

		spec1 = &reward.RewardSpec{ // blocks 10..10 (1 block)
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(1e18),
				TotalFee: big.NewInt(0.02e18),
				BurntFee: big.NewInt(0.01e18),
			},
			Proposer: big.NewInt(1.01e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(1.01e18),
			},
		}
		spec4 = &reward.RewardSpec{ // blocks 7..10 (4 blocks)
			RewardSummary: reward.RewardSummary{
				Minted:   big.NewInt(4e18),
				TotalFee: big.NewInt(0.08e18),
				BurntFee: big.NewInt(0.04e18),
			},
			Proposer: big.NewInt(4.04e18),
			Stakers:  big.NewInt(0),
			KIF:      big.NewInt(0),
			KEF:      big.NewInt(0),
			Rewards: map[common.Address]*big.Int{
				common.HexToAddress("0xfff"): big.NewInt(4.04e18),
			},
		}
		acc1 = spec1.ToAccumulatedResponse(upperHeader, upperHeader)
		acc4 = spec4.ToAccumulatedResponse(lowerHeader, upperHeader)
	)
	chain.EXPECT().GetHeaderByNumber(uint64(lowerNum)).Return(lowerHeader).AnyTimes()
	chain.EXPECT().GetHeaderByNumber(uint64(upperNum)).Return(upperHeader).AnyTimes()
	chain.EXPECT().CurrentBlock().Return(currBlock).AnyTimes()
	for num := lowerNum; num <= upperNum; num++ {
		r.EXPECT().GetBlockReward(uint64(num)).Return(spec1, nil).AnyTimes()
	}

	{
		// Different ways to provide the block range.
		testcases := []struct {
			lower    rpc.BlockNumber
			upper    rpc.BlockNumber
			expected *reward.AccumulatedRewardsResponse
		}{
			{lowerNum, upperNum, acc4},
			{lowerNum, lastNum, acc4},
			{lastNum, lastNum, acc1},
			{upperNum, upperNum, acc1},
		}
		for i, tc := range testcases {
			result, err := govAPI.GetRewardsAccumulated(tc.lower, tc.upper)
			assert.NoError(t, err, i)
			assert.Equal(t, tc.expected, result, i)
		}
	}
	{
		// Bad ways to provide the block range.
		_, err := govAPI.GetRewardsAccumulated(upperNum, lowerNum)
		assert.ErrorIs(t, err, reward.ErrInvalidBlockRange)
		_, err = govAPI.GetRewardsAccumulated(lowerNum, currNum+1)
		assert.ErrorIs(t, err, reward.ErrInvalidBlockRange)
	}
	{
		// Inject fault. Among the blocks 1..6, block 5 fails and the error must be propagated.
		for num := uint64(1); num <= 6; num++ {
			if num == 5 {
				r.EXPECT().GetBlockReward(num).Return(nil, reward.ErrNoReceipts).AnyTimes()
			} else {
				r.EXPECT().GetBlockReward(num).Return(spec1, nil).AnyTimes()
			}
		}
		chain.EXPECT().GetHeaderByNumber(uint64(1)).Return(&types.Header{Number: big.NewInt(1), Time: big.NewInt(1001)}).AnyTimes()
		chain.EXPECT().GetHeaderByNumber(uint64(6)).Return(&types.Header{Number: big.NewInt(6), Time: big.NewInt(1006)}).AnyTimes()
		_, err := govAPI.GetRewardsAccumulated(rpc.BlockNumber(1), rpc.BlockNumber(6))
		assert.ErrorIs(t, err, reward.ErrNoReceipts)
	}
}

func makeTestAPI(t *testing.T) (*gomock.Controller, *reward_mock.MockRewardModule, *mocks.MockBlockChain, *RewardKaiaAPI, *RewardGovAPI) {
	mockCtrl := gomock.NewController(t)

	r := reward_mock.NewMockRewardModule(mockCtrl)
	chain := mocks.NewMockBlockChain(mockCtrl)

	kaiaAPI := NewRewardKaiaAPI(r, chain)
	govAPI := NewRewardGovAPI(r, chain)
	return mockCtrl, r, chain, kaiaAPI, govAPI
}
