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
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"sync/atomic"

	"github.com/kaiachain/kaia/v2/accounts/abi/bind"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/contracts/contracts/system_contracts/rebalance"
	"github.com/kaiachain/kaia/v2/kaiax/supply"
)

var (
	zeroBurnAddress = common.HexToAddress("0x0")
	deadBurnAddress = common.HexToAddress("0xdead")
)

func (s *SupplyModule) GetTotalSupply(num uint64) (*supply.TotalSupply, error) {
	if cached, ok := s.supplyCache.Get(num); ok {
		return cached.(*supply.TotalSupply), nil
	}

	// Read accumulated supply checkpoint (minted, burntFee)
	// This is an essential component, so failure to read it immediately aborts the function.
	accReward, err := s.getAccReward(num)
	if err != nil {
		return nil, err
	}

	// Below components may fail to read, but try to return what it can.
	var (
		zeroBurn   *big.Int
		deadBurn   *big.Int
		kip103Burn *big.Int
		kip160Burn *big.Int
		errs       = make([]error, 0)
	)
	zeroBurn, deadBurn, err = s.getCanonicalBurn(num)
	if err != nil {
		errs = append(errs, err)
	}
	config := s.ChainConfig
	kip103Burn, err = s.getRebalanceBurn(num, config.Kip103CompatibleBlock, config.Kip103ContractAddress)
	if err != nil {
		errs = append(errs, err)
	}
	kip160Burn, err = s.getRebalanceBurn(num, config.Kip160CompatibleBlock, config.Kip160ContractAddress)
	if err != nil {
		errs = append(errs, err)
	}

	ts := accReward.ToTotalSupply(
		zeroBurn,
		deadBurn,
		kip103Burn,
		kip160Burn,
	)
	s.supplyCache.Add(num, ts)
	return ts, errors.Join(errs...)
}

// totalSupplyFromState exhausitively traverses all accounts in the state at the given block number.
func (s *SupplyModule) totalSupplyFromState(num uint64) (*big.Int, error) {
	header := s.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, supply.ErrNoBlock
	}
	stateDB, err := s.Chain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	dump := stateDB.RawDump()

	totalSupply := new(big.Int)
	for _, account := range dump.Accounts {
		balance, ok := new(big.Int).SetString(account.Balance, 10)
		if !ok {
			return nil, errors.New("malformed state dump")
		}
		totalSupply.Add(totalSupply, balance)
	}
	return totalSupply, nil
}

// getAccReward reads the supply checkpoint at the given block number.
// If the checkpoint is not found, it will re-accumulate from the nearest checkpoint.
func (s *SupplyModule) getAccReward(num uint64) (*supply.AccReward, error) {
	// Find from the database.
	accReward := ReadAccReward(s.ChainKv, num)
	if accReward != nil {
		return accReward, nil
	}

	// If not found, re-accumulate from the nearest checkpoint.
	fromNum := nearestCheckpointNumber(num)
	fromAccReward := ReadAccReward(s.ChainKv, fromNum)
	if fromAccReward == nil {
		return nil, supply.ErrNoSupplyCheckpoint
	}
	logger.Trace("on-demand reaccumulating supply checkpoint", "from", fromNum, "to", num)
	return s.accumulateRewards(fromNum, num, fromAccReward, false) // do not commit when re-accumulating
}

// getCanonicalBurn reads the balances of the canonical burn addresses (0x0, 0xdead) from the state.
func (s *SupplyModule) getCanonicalBurn(num uint64) (*big.Int, *big.Int, error) {
	header := s.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, nil, supply.ErrNoCanonicalBurn(supply.ErrNoBlock)
	}
	state, err := s.Chain.StateAt(header.Root)
	if err != nil {
		return nil, nil, supply.ErrNoCanonicalBurn(err)
	}
	return state.GetBalance(zeroBurnAddress), state.GetBalance(deadBurnAddress), nil
}

// getRebalanceBurn attempts to read the burnt amount from the rebalance memo.
// 1. Rebalance is not configured or the fork block is not reached: return 0.
// 2. Found the memo: return the burnt amount.
// 3. Rebalance is configured but the memo is not found: return nil.
// - 3a. the rebalance is misconfigured so that system.RebalanceTreasury did not change the state.
// - 3b. the memo is not yet submitted to the contract.
// 4. Something else went wrong: return nil.
//
// The case 3a returning 0 would be more accurate representation of the rebalance burn amount (no burn happened),
// but 3a is indistinguishable from 3b or 4. Therefore it returns nil and an error for safety.
// If we actually create case 3a by accident (i.e. rebalance actually not happened), fix this function to return 0 for case 3a.
func (s *SupplyModule) getRebalanceBurn(num uint64, forkNum *big.Int, addr common.Address) (*big.Int, error) {
	bigNum := new(big.Int).SetUint64(num)

	if forkNum == nil || forkNum.Sign() == 0 || (addr == common.Address{}) || forkNum.Cmp(bigNum) > 0 {
		// 1. rebalance is not configured or the rebalance forkNum has not passed (at the given block number).
		return big.NewInt(0), nil
	}

	if burnt, ok := s.memoCache.Get(addr); ok {
		// 2. found the memo in cache.
		return burnt.(*big.Int), nil
	}

	// Load the state at latest block, not the rebalance fork block.
	// The memo is manually stored in the contract after-the-fact by calling the finalizeContract function.
	// Therefore it's safest to read from the latest state.
	backend := backends.NewBlockchainContractBackend(s.Chain, nil, nil)
	caller, err := rebalance.NewTreasuryRebalanceV2Caller(addr, backend)
	if err != nil {
		// 4. contract call failed.
		return nil, supply.ErrNoRebalanceBurn(err)
	}
	memo, err := caller.Memo(&bind.CallOpts{BlockNumber: nil}) // call at the latest block
	if err != nil {
		// 3a. the contract reverted or the contract is not there.
		// 4. contract call failed for other unknown reasons.
		return nil, supply.ErrNoRebalanceBurn(err)
	}
	if memo == "" {
		// 3a. the contract is intact but the rebalance did not happen due to e.g. insufficient funds.
		// 3b. the memo is not yet submitted to the contract.
		// Return nil to prevent totalSupply calculation, therefore prevents misinformation.
		return nil, supply.ErrNoRebalanceBurn(supply.ErrNoRebalanceMemo)
	}

	result := struct { // See system.rebalanceResult struct
		Burnt *big.Int `json:"burnt"`
	}{}

	if s.ChainConfig.ChainID.Uint64() == 1001 && strings.HasPrefix(memo, "before") {
		// 2. override for Kairos testnet
		result.Burnt = new(big.Int)
		result.Burnt.SetString("-3704329462904320084000000000", 10)
	} else {
		if err := json.Unmarshal([]byte(memo), &result); err != nil {
			// 4. memo is malformed
			return nil, supply.ErrNoRebalanceBurn(err)
		}
	}
	// 2. found the memo in state.
	s.memoCache.Add(addr, result.Burnt)
	return result.Burnt, nil
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
		accReward.TotalMinted.Add(accReward.TotalMinted, summary.Minted)
		accReward.BurntFee.Add(accReward.BurntFee, summary.BurntFee)

		if write && (num%checkpointInterval) == 0 {
			WriteAccReward(s.ChainKv, num, accReward)
			WriteLastAccRewardNumber(s.ChainKv, num)
		}
		if (num % accumulateLogInterval) == 0 {
			logger.Info("Accumulated block rewards", "number", num, "minted", accReward.TotalMinted.String(), "burntFee", accReward.BurntFee.String())
		}
	}
	return accReward, nil
}
