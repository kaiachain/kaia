package staking

import "github.com/kaiachain/kaia/networks/rpc"

func (s *StakingModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   newStakingAPI(s),
			Public:    true,
		},
	}
}

type stakingAPI struct {
	s *StakingModule
}

func newStakingAPI(s *StakingModule) *stakingAPI {
	return &stakingAPI{s}
}

func (api *stakingAPI) GetStakingInfo(num *rpc.BlockNumber) (*StakingInfo, error) {
	var blockNum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNum = api.s.Chain.CurrentBlock().NumberU64()
	} else {
		blockNum = uint64(*num)
	}

	return api.s.GetStakingInfo(blockNum)
}
