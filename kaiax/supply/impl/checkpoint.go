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
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/kaiax/supply"
)

var (
	checkpointInterval    = uint64(128)    // AccReward checkpoint interval
	accumulateLogInterval = uint64(102400) // Periodic log in accumulateCheckpoint().
)

func nearestCheckpointInterval(num uint64) uint64 {
	return num - (num % checkpointInterval)
}

// loadLastCheckpoint loads the last supply checkpoint from the database.
func (s *SupplyModule) loadLastCheckpoint() error {
	var (
		lastAccNum    = ReadLastSupplyCheckpointNumber(s.ChainKv)
		lastAccReward = ReadSupplyCheckpoint(s.ChainKv, lastAccNum)
	)

	// Something is wrong. Reset to genesis.
	if lastAccNum > 0 && lastAccReward == nil {
		logger.Error("Last supply checkpoint not found. Restarting supply catchup", "last", lastAccNum)
		WriteLastSupplyCheckpointNumber(s.ChainKv, 0) // soft reset to genesis
		lastAccNum = 0
	}

	// If we are at genesis (either data empty or soft reset), write the genesis supply.
	if lastAccNum == 0 {
		genesisTotalSupply, err := s.totalSupplyFromState(0)
		if err != nil {
			return err
		}
		lastAccReward = &supply.AccReward{
			Minted:   genesisTotalSupply,
			BurntFee: big.NewInt(0),
		}
		WriteLastSupplyCheckpointNumber(s.ChainKv, 0)
		WriteSupplyCheckpoint(s.ChainKv, 0, lastAccReward)
		logger.Info("Stored genesis total supply", "supply", genesisTotalSupply)
	}

	s.lastAccNum = lastAccNum
	s.lastAccReward = lastAccReward
	return nil
}

// accumulateRewards accumulates the reward increments from `fromNum` to `toNum`, inclusive.
// If `write` is true, the intermediate results at checkpointInterval will be written to the database.
func (s *SupplyModule) accumulateRewards(fromNum, toNum uint64, fromAccReward *supply.AccReward, write bool) (*supply.AccReward, error) {
	accReward := fromAccReward.Copy() // make a copy because we're updating it in-place.

	for num := fromNum + 1; num <= toNum; num++ {
		if atomic.LoadUint32(&s.quit) == 1 { // Received quit signal
			return nil, supply.ErrSupplyModuleQuit
		}

		summary, err := s.RewardModule.GetRewardSummary(num)
		if err != nil {
			return nil, err
		}
		accReward.Minted.Add(accReward.Minted, summary.Minted)
		accReward.BurntFee.Add(accReward.BurntFee, summary.BurntFee)

		if write && (num%checkpointInterval) == 0 {
			WriteSupplyCheckpoint(s.ChainKv, num, accReward)
			WriteLastSupplyCheckpointNumber(s.ChainKv, num)
		}
		if (num % accumulateLogInterval) == 0 {
			logger.Info("Accumulated block rewards", "number", num, "minted", accReward.Minted.String(), "burntFee", accReward.BurntFee.String())
		}
	}
	return accReward, nil
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
