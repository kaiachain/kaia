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

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ staking.StakingModule = &StakingModule{}

	logger = log.NewModuleLogger(log.KaiaxStaking)
)

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       backends.BlockChainForCaller
}

type StakingModule struct {
	InitOpts

	// Genesis configurations
	stakingInterval uint64
	useGiniCoeff    bool
	minimumStake    *big.Int

	stakingInfoCache *lru.ARCCache // cached by sourceNum
}

func NewStakingModule() *StakingModule {
	cache, _ := lru.NewARC(128)
	return &StakingModule{
		stakingInfoCache: cache,
	}
}

func (s *StakingModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.Chain == nil || opts.ChainKv == nil {
		return staking.ErrInitUnexpectedNil
	}
	s.InitOpts = *opts

	if s.ChainConfig.Governance == nil || s.ChainConfig.Governance.Reward == nil {
		return staking.ErrInitUnexpectedNil
	}
	if s.ChainConfig.Governance.Reward.StakingUpdateInterval == 0 {
		return staking.ErrZeroStakingInterval
	}
	s.stakingInterval = s.ChainConfig.Governance.Reward.StakingUpdateInterval
	s.useGiniCoeff = s.ChainConfig.Governance.Reward.UseGiniCoeff
	s.minimumStake = s.ChainConfig.Governance.Reward.MinimumStake

	return nil
}

func (s *StakingModule) Start() error {
	// This module may have restarted after a rewind. Purge the cache.
	s.stakingInfoCache.Purge()
	return nil
}

func (s *StakingModule) Stop() {
}
