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
	"math/big"
	"sync"
	"sync/atomic"
	"time"

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
	supplyCacheSize       = 86400
	checkpointInterval    = uint64(128)    // AccReward checkpoint interval
	accumulateLogInterval = uint64(102400) // Periodic log in accumulateRewards().
)

func nearestCheckpointNumber(num uint64) uint64 {
	return num - (num % checkpointInterval)
}

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
	mu            sync.RWMutex
	lastAccNum    uint64            // Last AccReward number
	lastAccReward *supply.AccReward // Last AccReward

	// Stops long-running tasks.
	quit   uint32         // stops the synchronous loop in accumulateCheckpoint
	quitCh chan struct{}  // stops the goroutine in select loop
	wg     sync.WaitGroup // wait for the goroutine to finish

	supplyCache *lru.ARCCache // (number uint64) -> (totalSupply *TotalSupply)
	memoCache   *lru.ARCCache // (contract Address) -> (memo.Burnt *big.Int)
}

func NewSupplyModule() *SupplyModule {
	supplyCache, _ := lru.NewARC(supplyCacheSize)
	memoCache, _ := lru.NewARC(10)
	return &SupplyModule{
		supplyCache: supplyCache,
		memoCache:   memoCache,
		quitCh:      make(chan struct{}, 1),
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
	// Reload the last checkpoint from database.
	if err := s.loadLastCheckpoint(); err != nil {
		return err
	}

	// Reset the caches.
	s.supplyCache.Purge()
	s.memoCache.Purge()

	// Reset the quit state.
	atomic.StoreUint32(&s.quit, 0)
	s.quitCh = make(chan struct{}, 1)
	s.wg.Add(1)
	go s.catchup()
	return nil
}

func (s *SupplyModule) Stop() {
	atomic.StoreUint32(&s.quit, 1)
	s.quitCh <- struct{}{}
	s.wg.Wait()
}

// loadLastCheckpoint loads the last supply checkpoint from the database.
func (s *SupplyModule) loadLastCheckpoint() error {
	var (
		lastAccNum    = ReadLastAccRewardNumber(s.ChainKv)
		lastAccReward = ReadAccReward(s.ChainKv, lastAccNum)
	)

	// Something is wrong. Reset to genesis.
	if lastAccNum > 0 && lastAccReward == nil {
		logger.Error("Last supply checkpoint not found. Restarting supply catchup", "last", lastAccNum)
		WriteLastAccRewardNumber(s.ChainKv, 0) // soft reset to genesis
		lastAccNum = 0
	}

	// If we are at genesis (either data empty or soft reset), write the genesis supply.
	if lastAccNum == 0 {
		genesisTotalSupply, err := s.totalSupplyFromState(0)
		if err != nil {
			return err
		}
		lastAccReward = &supply.AccReward{
			TotalMinted: genesisTotalSupply,
			BurntFee:    big.NewInt(0),
		}
		WriteLastAccRewardNumber(s.ChainKv, 0)
		WriteAccReward(s.ChainKv, 0, lastAccReward)
		logger.Info("Stored genesis total supply", "supply", genesisTotalSupply)
	}

	s.lastAccNum = lastAccNum
	s.lastAccReward = lastAccReward
	return nil
}

// catchup is a long-running goroutine that accumulates the supply checkpoint until the current head block.
func (s *SupplyModule) catchup() {
	defer s.wg.Done()

	for {
		newAccNum := s.Chain.CurrentBlock().NumberU64()
		s.mu.RLock()
		lastAccNum := s.lastAccNum
		lastAccReward := s.lastAccReward.Copy()
		s.mu.RUnlock()

		// A gap detected. Accumulate to the current block.
		if lastAccNum < newAccNum {
			newAccReward, err := s.accumulateRewards(lastAccNum, newAccNum, lastAccReward, true)
			if err != nil {
				if err != supply.ErrSupplyModuleQuit {
					logger.Error("Total supply accumulate failed", "from", lastAccNum, "to", newAccNum, "err", err)
				}
				return
			}
			s.mu.Lock()
			s.lastAccNum = newAccNum
			s.lastAccReward = newAccReward
			s.mu.Unlock()
			// Because current head may have increased while we accumulate, we need to check again.
			continue
		}

		// No gap detected. Sleep a while and check again just in case.
		// If PostInsertBlock() is filling in the gap, this loop would do nothing but waiting.
		timer := time.NewTimer(time.Second)
		select {
		case <-s.quitCh:
			return
		case <-timer.C:
		}
	}
}
