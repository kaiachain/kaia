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

package supply

import (
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ supply.SupplyModule = &SupplyModule{}

	logger = log.NewModuleLogger(log.KaiaxSupply)

	// A day; Some total supply consumers might want daily supply.
	// It may be evicted by the historic query though.
	checkpointCacheSize = 86400
)

type InitOpts struct {
	ChainKv      database.Database
	ChainConfig  *params.ChainConfig
	Chain        backends.BlockChainForCaller
	RewardModule reward.RewardModule
}

type SupplyModule struct {
	InitOpts

	// Accumulated supply checkpoint so far.
	// This in-memory variables advance every block, but the database is updated every checkpointInterval.
	mu             sync.Mutex
	lastNum        uint64            // Last in-memory supply checkpoint number
	lastCheckpoint *supplyCheckpoint // Last in-memory supply checkpoint

	// Stops long-running tasks.
	quit   uint32         // stops the synchronous loop in accumulateCheckpoint
	quitCh chan struct{}  // stops the goroutine in select loop
	wg     sync.WaitGroup // wait for the goroutine to finish

	checkpointCache *lru.ARCCache // (number uint64) -> (checkpoint *supplyCheckpoint)
	memoCache       *lru.ARCCache // (contract Address) -> (memo.Burnt *big.Int)
}

func NewSupplyModule() *SupplyModule {
	checkpointCache, _ := lru.NewARC(checkpointCacheSize)
	memoCache, _ := lru.NewARC(10)
	return &SupplyModule{
		checkpointCache: checkpointCache,
		memoCache:       memoCache,
	}
}

func (s *SupplyModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainKv == nil || opts.ChainConfig == nil || opts.Chain == nil || opts.RewardModule == nil {
		return supply.ErrInitUnexpectedNil
	}
	s.InitOpts = *opts
	return nil
}

func (s *SupplyModule) Start() error {
	s.checkpointCache.Purge()
	s.memoCache.Purge()
	return s.loadLastCheckpoint()
}

func (s *SupplyModule) Stop() {
}
