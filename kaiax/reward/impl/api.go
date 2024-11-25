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
	"runtime"
	"sync"

	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/networks/rpc"
)

var accumulatedRewardsRangeLimit = uint64(604800) // 7 days. naive resource protection

func (r *RewardModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewRewardKaiaAPI(r, r.Chain),
			Public:    true,
		},
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   NewRewardGovAPI(r, r.Chain),
			Public:    true,
		},
	}
}

// RewardKaiaAPI defines the kaia namespace APIs.
type RewardKaiaAPI struct {
	r     reward.RewardModule
	chain blockChain
}

func NewRewardKaiaAPI(r reward.RewardModule, chain blockChain) *RewardKaiaAPI {
	return &RewardKaiaAPI{r: r, chain: chain}
}

func (api *RewardKaiaAPI) GetRewards(num *rpc.BlockNumber) (*reward.RewardResponse, error) {
	var blockNum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNum = api.chain.CurrentBlock().NumberU64()
	} else {
		blockNum = num.Uint64()
	}
	return api.r.GetBlockReward(blockNum)
}

// RewardGovAPI defines the governance namespace APIs.
type RewardGovAPI struct {
	r     reward.RewardModule
	chain blockChain
}

func NewRewardGovAPI(r reward.RewardModule, chain blockChain) *RewardGovAPI {
	return &RewardGovAPI{r: r, chain: chain}
}

func (api *RewardGovAPI) GetRewardsAccumulated(lower, upper rpc.BlockNumber) (*reward.AccumulatedRewardsResponse, error) {
	// normalize block numbers
	currentNum := api.chain.CurrentBlock().NumberU64()
	var lowerNum uint64
	var upperNum uint64
	if lower == rpc.LatestBlockNumber || lower == rpc.PendingBlockNumber {
		lowerNum = currentNum
	} else {
		lowerNum = lower.Uint64()
	}
	if upper == rpc.LatestBlockNumber || upper == rpc.PendingBlockNumber {
		upperNum = currentNum
	} else {
		upperNum = upper.Uint64()
	}
	if lowerNum > upperNum || upperNum > currentNum {
		return nil, reward.ErrInvalidBlockRange
	}
	count := upperNum - lowerNum + 1
	if count > accumulatedRewardsRangeLimit {
		return nil, reward.ErrBlockRangeLimit
	}

	// Fetch block timestamps
	// Fail fast before accumulating rewards
	lowerHeader := api.chain.GetHeaderByNumber(lowerNum)
	upperHeader := api.chain.GetHeaderByNumber(upperNum)
	if lowerHeader == nil || upperHeader == nil {
		return nil, reward.ErrNoBlock
	}

	// Accumulate the block rewards
	spec, err := api.accumulateRewards(lowerNum, upperNum)
	if err != nil {
		return nil, err
	}
	return spec.ToAccumulatedResponse(lowerHeader, upperHeader), nil
}

func (api *RewardGovAPI) accumulateRewards(lower, upper uint64) (*reward.RewardSpec, error) {
	var (
		accSpec = reward.NewRewardSpec()
		mu      sync.Mutex

		numWorkers = runtime.NumCPU()
		reqCh      = make(chan uint64, numWorkers)
		errCh      = make(chan error, numWorkers)
		wg         sync.WaitGroup
	)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range reqCh {
				spec, err := api.r.GetBlockReward(num)
				if err != nil {
					errCh <- err
					return
				}
				mu.Lock()
				accSpec.Add(spec)
				mu.Unlock()
			}
		}()
	}

	for num := lower; num <= upper; num++ {
		reqCh <- num
	}
	close(reqCh)

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}
	return accSpec, nil
}
