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

	"github.com/kaiachain/kaia/kaiax/supply"
)

var (
	checkpointInterval    = uint64(128)
	accumulateLogInterval = uint64(102400) // Periodic log in accumulateCheckpoint().
)

// supplyCheckpoint represents the accumulated minted and burnt fee up to a specific block number.
type supplyCheckpoint struct {
	Minted   *big.Int
	BurntFee *big.Int
}

func (sc *supplyCheckpoint) Copy() *supplyCheckpoint {
	return &supplyCheckpoint{
		Minted:   new(big.Int).Set(sc.Minted),
		BurntFee: new(big.Int).Set(sc.BurntFee),
	}
}

func nearestCheckpointInterval(num uint64) uint64 {
	return num - (num % checkpointInterval)
}

// loadLastCheckpoint loads the last supply checkpoint from the database.
func (s *SupplyModule) loadLastCheckpoint() error {
	var (
		lastNum        = ReadLastSupplyCheckpointNumber(s.ChainKv)
		lastCheckpoint = ReadSupplyCheckpoint(s.ChainKv, lastNum)
	)

	// Something is wrong. Reset to genesis.
	if lastNum > 0 && lastCheckpoint == nil {
		logger.Error("Last supply checkpoint not found. Restarting supply catchup", "last", lastNum)
		WriteLastSupplyCheckpointNumber(s.ChainKv, 0) // soft reset to genesis
		lastNum = 0
	}

	// If we are at genesis (either data empty or soft reset), write the genesis supply.
	if lastNum == 0 {
		genesisTotalSupply, err := s.totalSupplyFromState(0)
		if err != nil {
			return err
		}
		lastCheckpoint = &supplyCheckpoint{
			Minted:   genesisTotalSupply,
			BurntFee: big.NewInt(0),
		}
		WriteLastSupplyCheckpointNumber(s.ChainKv, 0)
		WriteSupplyCheckpoint(s.ChainKv, 0, lastCheckpoint)
		logger.Info("Stored genesis total supply", "supply", genesisTotalSupply)
	}

	s.lastNum = lastNum
	s.lastCheckpoint = lastCheckpoint
	return nil
}

// accumulateCheckpoint accumulates the total supply increments from `fromNum` to `toNum`, inclusive.
// If `write` is true, the intermediate results at checkpointInterval will be written to the database.
func (s *SupplyModule) accumulateCheckpoint(fromNum, toNum uint64, fromCheckpoint *supplyCheckpoint, write bool) (*supplyCheckpoint, error) {
	checkpoint := fromCheckpoint.Copy() // make a copy because we're updating it in-place.

	for num := fromNum + 1; num <= toNum; num++ {
		if atomic.LoadUint32(&s.quit) == 1 { // Received quit signal
			return nil, supply.ErrSupplyModuleQuit
		}

		summary, err := s.RewardModule.GetRewardSummary(num)
		if err != nil {
			return nil, err
		}
		checkpoint.Minted.Add(checkpoint.Minted, summary.Minted)
		checkpoint.BurntFee.Add(checkpoint.BurntFee, summary.BurntFee)

		if write && (num%checkpointInterval) == 0 {
			WriteSupplyCheckpoint(s.ChainKv, num, checkpoint)
			WriteLastSupplyCheckpointNumber(s.ChainKv, num)
		}
		if (num % accumulateLogInterval) == 0 {
			logger.Info("Accumulated block rewards", "number", num, "minted", checkpoint.Minted.String(), "burntFee", checkpoint.BurntFee.String())
		}
	}
	return checkpoint, nil
}
