// Modifications Copyright 2022 The Klaytn Authors
// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/gasprice/feehistory.go (2021/11/09).
// Modified and improved for the klaytn development

package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestFeeHistory(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	cases := []struct {
		maxHeader, maxBlock int
		count               int
		last                rpc.BlockNumber
		percent             []float64
		expFirst            uint64
		expCount            int
		expErr              error
	}{
		{1000, 1000, 10, 30, nil, 21, 10, nil},
		{1000, 1000, 10, 30, []float64{0, 10}, 21, 10, nil},
		{1000, 1000, 10, 30, []float64{20, 10}, 0, 0, fmt.Errorf("%w: #%d:%f > #%d:%f", errInvalidPercentile, 0, 20.0000, 1, 10.0000)},
		{1000, 1000, 1000000000, 30, nil, 0, 31, nil},
		{1000, 1000, 1000000000, rpc.LatestBlockNumber, nil, 0, 33, nil},
		{1000, 1000, 10, 40, nil, 0, 0, fmt.Errorf("%w: requested %d, head %d", errRequestBeyondHead, 40, 32)},
		{20, 2, 100, rpc.LatestBlockNumber, nil, 13, 20, nil},
		{20, 2, 100, rpc.LatestBlockNumber, []float64{0, 10}, 31, 2, nil},
		{20, 2, 100, 32, []float64{0, 10}, 31, 2, nil},
		{1000, 1000, 1, rpc.PendingBlockNumber, nil, 0, 0, nil},
		{1000, 1000, 2, rpc.PendingBlockNumber, nil, 32, 1, nil},
	}
	magmaBlock, kaiaBlock := int64(16), int64(20)
	backend, govModule := newTestBackend(t, big.NewInt(magmaBlock), big.NewInt(kaiaBlock))

	defer backend.teardown()
	for i, c := range cases {
		config := Config{
			MaxHeaderHistory: c.maxHeader,
			MaxBlockHistory:  c.maxBlock,
			MaxPrice:         big.NewInt(500000000000),
		}
		oracle := NewOracle(backend, config, nil, govModule)

		first, reward, baseFee, ratio, err := oracle.FeeHistory(context.Background(), c.count, c.last, c.percent)

		expReward := c.expCount
		if len(c.percent) == 0 {
			expReward = 0
		}
		expBaseFee := c.expCount
		if expBaseFee != 0 {
			expBaseFee++
		}

		assert.Equal(t, c.expFirst, first.Uint64(), "Test case %d: first block mismatch, want %d, got %d", i, c.expFirst, first.Uint64())
		assert.Equal(t, expReward, len(reward), "Test case %d: reward array length mismatch, want %d, got %d", i, expReward, len(reward))
		assert.Equal(t, expBaseFee, len(baseFee), "Test case %d: baseFee array length mismatch, want %d, got %d", i, expBaseFee, len(baseFee))
		assert.Equal(t, c.expCount, len(ratio), "Test case %d: baseFee array length mismatch, want %d, got %d", i, c.expCount, len(ratio))
		assert.Equal(t, c.expErr, err, "Test case %d: error mismatch, want %v, got %v", i, c.expErr, err)
	}

	// Last check. Check the value of Reward, BaseFee and GasUsedRatio of every block.
	var (
		config = Config{
			MaxHeaderHistory: 1000,
			MaxBlockHistory:  1000,
			MaxPrice:         big.NewInt(500000000000),
		}
		oracle = NewOracle(backend, config, nil, govModule)

		atMagmaPset    = oracle.govModule.GetParamSet(uint64(magmaBlock))
		afterMagmaPset = oracle.govModule.GetParamSet(uint64(magmaBlock + 1))

		beforeMagmaExpectedBaseFee = big.NewInt(0)
		atMagmaExpectedBaseFee     = big.NewInt(int64(atMagmaPset.LowerBoundBaseFee))
		afterMagmaExpectedBaseFee  = big.NewInt(int64(afterMagmaPset.LowerBoundBaseFee))

		beforeMagmaExpectedGasUsedRatio = float64(21000) / float64(params.UpperGasLimit)
		atMagmaExpectedGasUsedRatio     = float64(21000) / float64(atMagmaPset.MaxBlockGasUsedForBaseFee)
		afterMagmaExpectedGasUsedRatio  = float64(21000) / float64(afterMagmaPset.MaxBlockGasUsedForBaseFee)
	)

	first, _, baseFee, ratio, err := oracle.FeeHistory(context.Background(), 32, rpc.LatestBlockNumber, nil)
	assert.Equal(t, first, big.NewInt(1))
	assert.Nil(t, err)

	// magma hardfork
	assert.Equal(t, []*big.Int{beforeMagmaExpectedBaseFee, atMagmaExpectedBaseFee, afterMagmaExpectedBaseFee}, baseFee[14:17])
	assert.Equal(t, []float64{beforeMagmaExpectedGasUsedRatio, atMagmaExpectedGasUsedRatio, afterMagmaExpectedGasUsedRatio}, ratio[14:17])

	// kaia hardfork
	// check the value of reward
	// assert.Equal(t, <impl>, gasTip[18:21])
}
