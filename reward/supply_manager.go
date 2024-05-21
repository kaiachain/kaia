// Copyright 2024 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package reward

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/contracts/system_contracts/rebalance"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/storage/database"
)

var (
	supplyCacheSize   = 86400          // A day; Some total supply consumers might want daily supply.
	supplyLogInterval = uint64(102400) // Periodic total supply log.
	zeroBurnAddress   = common.HexToAddress("0x0")
	deadBurnAddress   = common.HexToAddress("0xdead")

	errSupplyManagerQuit = errors.New("supply manager quit")
	errNoAccReward       = errors.New("accumulated reward not stored")
	errNoBlock           = errors.New("block not found")
	errNoRebalanceMemo   = errors.New("rebalance memo not yet stored")
)

func errNoCanonicalBurn(err error) error {
	return fmt.Errorf("cannot determine canonical (0x0, 0xdead) burn amount: %v", err)
}

func errNoRebalanceBurn(err error) error {
	return fmt.Errorf("cannot determine rebalance (kip103, kip160) burn amount: %v", err)
}

// SupplyManager tracks the total supply of native tokens.
// Note that SupplyManager only deals with the block rewards.
// Factors like KIP103, KIP160, 0xdead burning are not considered and must be accounted
// in the APIs that return the total supply.
type SupplyManager interface {
	// Start starts the supply manager goroutine that catches up or tracks the block rewards.
	Start()

	// Stop stops the supply manager goroutine.
	Stop()

	// GetTotalSupply returns the total supply amounts at the given block number,
	// broken down by minted amount and burnt amounts of each methods.
	GetTotalSupply(num uint64) (*TotalSupply, error)
}

type TotalSupply struct {
	TotalSupply *big.Int // The total supply of the native token. i.e. Minted - Burnt.
	TotalMinted *big.Int // Total minted amount.
	TotalBurnt  *big.Int // Total burnt amount. Sum of all burnt amounts below.
	BurntFee    *big.Int // from tx fee burn. ReadAccReward(num).BurntFee.
	ZeroBurn    *big.Int // balance of 0x0 (zero) address.
	DeadBurn    *big.Int // balance of 0xdead (dead) address.
	Kip103Burn  *big.Int // by KIP103 fork. Read from its memo.
	Kip160Burn  *big.Int // by KIP160 fork. Read from its memo.
}

type supplyManager struct {
	// Externally injected dependencies
	chain              blockChain
	chainHeadChan      chan blockchain.ChainHeadEvent
	chainHeadSub       event.Subscription
	gov                governanceHelper
	db                 database.DBManager
	checkpointInterval uint64

	// Internal data structures
	accRewardCache *lru.ARCCache  // Cache (number uint64) -> (accReward *database.AccReward)
	memoCache      *lru.ARCCache  // Cache (address Address) -> (memo.Burnt *big.Int)
	quit           uint32         // Stop the goroutine in initial catchup stage
	quitCh         chan struct{}  // Stop the goroutine in event subscription state
	wg             sync.WaitGroup // background goroutine wait group for shutting down
}

// NewSupplyManager creates a new supply manager.
// The TotalSupply data is stored every checkpointInterval blocks.
func NewSupplyManager(chain blockChain, gov governanceHelper, db database.DBManager, checkpointInterval uint) *supplyManager {
	accRewardCache, _ := lru.NewARC(supplyCacheSize)
	memoCache, _ := lru.NewARC(10)

	return &supplyManager{
		chain:              chain,
		chainHeadChan:      make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		gov:                gov,
		db:                 db,
		checkpointInterval: uint64(checkpointInterval),
		accRewardCache:     accRewardCache,
		memoCache:          memoCache,
		quitCh:             make(chan struct{}, 1), // make sure Stop() doesn't block if catchup() has exited before Stop()
	}
}

func (sm *supplyManager) Start() {
	sm.wg.Add(1)
	go sm.catchup()
}

func (sm *supplyManager) Stop() {
	atomic.StoreUint32(&sm.quit, 1)
	sm.quitCh <- struct{}{}
	sm.wg.Wait()
	if sm.chainHeadSub != nil {
		sm.chainHeadSub.Unsubscribe()
	}
}

func (sm *supplyManager) GetAccReward(num uint64) (*database.AccReward, error) {
	if accReward, ok := sm.accRewardCache.Get(num); ok {
		return accReward.(*database.AccReward), nil
	}

	lastNum := sm.db.ReadLastAccRewardBlockNumber()
	if lastNum < num { // soft deleted
		return nil, errNoAccReward
	}

	accReward := sm.db.ReadAccReward(num)
	if accReward == nil {
		return nil, errNoAccReward
	}

	sm.accRewardCache.Add(num, accReward.Copy())
	return accReward, nil
}

func (sm *supplyManager) GetCanonicalBurn(num uint64) (zero *big.Int, dead *big.Int, err error) {
	header := sm.chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, nil, errNoCanonicalBurn(errNoBlock)
	}
	state, err := sm.chain.StateAt(header.Root)
	if err != nil {
		return nil, nil, errNoCanonicalBurn(err)
	}
	return state.GetBalance(zeroBurnAddress), state.GetBalance(deadBurnAddress), nil
}

// GetRebalanceBurn attempts to read the burnt amount from the rebalance memo.
// 1. Rebalance is not configured or the fork block is not reached: return 0.
// 2. Found the memo: return the burnt amount.
// 3. Rebalance is configured but the memo is not found: return nil.
// - 3a. the rebalance is misconfigured so that system.RebalanceTreasury did not change the state.
// - 3b. the memo is not yet submitted to the contract.
// 4. Something else went wrong: return nil.
//
// The case 3a returning 0 would be more accurate representation of the rebalance burn amount,
// but 3a is indistinguishable from 3b or 4. Therefore it returns nil and an error for safety.
// If we actually create case 3a by accident (i.e. rebalance skipped at the fork block), fix this function to return 0.
func (sm *supplyManager) GetRebalanceBurn(num uint64, forkNum *big.Int, addr common.Address) (*big.Int, error) {
	bigNum := new(big.Int).SetUint64(num)

	if forkNum == nil || forkNum.Cmp(bigNum) > 0 {
		// 1. rebalance is not configured or the rebalance forkNum has not passed (at the given block number).
		return big.NewInt(0), nil
	}

	if burnt, ok := sm.memoCache.Get(addr); ok {
		// 2. memo is in cache.
		return burnt.(*big.Int), nil
	}

	// Load the state at latest block, not the rebalance fork block.
	// The memo is manually stored in the contract after-the-fact by calling the finalizeContract function.
	// Therefore it's safest to read from the latest state.
	backend := backends.NewBlockchainContractBackend(sm.chain, nil, nil)
	caller, err := rebalance.NewTreasuryRebalanceV2Caller(addr, backend)
	if err != nil {
		// 4. contract call failed.
		return nil, errNoRebalanceBurn(err)
	}
	memo, err := caller.Memo(&bind.CallOpts{BlockNumber: nil}) // call at the latest block
	if err != nil {
		// 3a. the contract reverted or the contract is not there.
		// 4. contract call failed for other unknown reasons.
		return nil, errNoRebalanceBurn(err)
	}
	if memo == "" {
		// 3a. the contract is intact but the rebalance did not happen due to e.g. insufficient funds.
		// 3b. the memo is not yet submitted to the contract.
		// Return nil to prevent totalSupply calculation, therefore prevents misinformation.
		return nil, errNoRebalanceBurn(errNoRebalanceMemo)
	}

	result := struct { // See system.rebalanceResult struct
		Burnt *big.Int `json:"burnt"`
	}{}
	if err := json.Unmarshal([]byte(memo), &result); err != nil {
		// 4. memo is malformed
		return nil, errNoRebalanceBurn(err)
	}

	// 2. found the memo
	sm.memoCache.Add(addr, result.Burnt)
	return result.Burnt, nil
}

func (sm *supplyManager) GetTotalSupply(num uint64) (*TotalSupply, error) {
	errs := make([]error, 0)
	ts := new(TotalSupply)

	// Read accumulated rewards (minted, burntFee)
	// This is an essential component, so failure to read it immediately aborts the function.
	accReward, err := sm.GetAccReward(num)
	if err != nil {
		return nil, err
	}
	ts.TotalMinted = accReward.Minted
	ts.BurntFee = accReward.BurntFee

	// Read canonical burn address balances
	// Leave them nil if the historic state isn't available.
	ts.ZeroBurn, ts.DeadBurn, err = sm.GetCanonicalBurn(num)
	if err != nil {
		errs = append(errs, err)
	}

	// Read KIP103 and KIP160 burns
	// 1. Leave them zero if the rebalance is not configured or the fork block is not reached.
	// 2. Leave them nil if the rebalance should have been executed but rebalance memo isn't available yet.
	config := sm.chain.Config()
	ts.Kip103Burn, err = sm.GetRebalanceBurn(num, config.Kip103CompatibleBlock, config.Kip103ContractAddress)
	if err != nil {
		errs = append(errs, err)
	}
	ts.Kip160Burn, err = sm.GetRebalanceBurn(num, config.Kip160CompatibleBlock, config.Kip160ContractAddress)
	if err != nil {
		errs = append(errs, err)
	}

	// TotalBurnt and TotalSupply is only calculated if all components are available.
	if ts.BurntFee != nil && ts.ZeroBurn != nil && ts.DeadBurn != nil && ts.Kip103Burn != nil && ts.Kip160Burn != nil {
		ts.TotalBurnt = new(big.Int)
		ts.TotalBurnt.Add(ts.TotalBurnt, ts.BurntFee)
		ts.TotalBurnt.Add(ts.TotalBurnt, ts.ZeroBurn)
		ts.TotalBurnt.Add(ts.TotalBurnt, ts.DeadBurn)
		ts.TotalBurnt.Add(ts.TotalBurnt, ts.Kip103Burn)
		ts.TotalBurnt.Add(ts.TotalBurnt, ts.Kip160Burn)
		ts.TotalSupply = new(big.Int).Sub(ts.TotalMinted, ts.TotalBurnt)
	}

	return ts, errors.Join(errs...)
}

// catchup accumulates the block rewards until the current block.
// The result will be written to the database.
func (sm *supplyManager) catchup() {
	defer sm.wg.Done()

	var (
		headNum = sm.chain.CurrentBlock().NumberU64()
		lastNum = sm.db.ReadLastAccRewardBlockNumber()
	)

	if lastNum > 0 && sm.db.ReadAccReward(lastNum) == nil {
		logger.Error("Last accumulated reward not found. Restarting supply catchup", "last", lastNum, "head", headNum)
		sm.db.WriteLastAccRewardBlockNumber(0) // soft reset to genesis
	}

	// Store genesis supply if not exists
	// ReadLastAccRewardBlockNumber is 0 when either (a) the database is empty or (b) was soft reset to 0.
	if sm.db.ReadLastAccRewardBlockNumber() == 0 {
		genesisTotalSupply, err := sm.totalSupplyFromState(0)
		if err != nil {
			logger.Error("totalSupplyFromState failed", "number", 0, "err", err)
			return
		}
		sm.db.WriteAccReward(0, &database.AccReward{
			Minted:   genesisTotalSupply,
			BurntFee: big.NewInt(0),
		})
		sm.db.WriteLastAccRewardBlockNumber(0)
		lastNum = 0
		logger.Info("Stored genesis total supply", "supply", genesisTotalSupply)
	}

	lastNum = sm.db.ReadLastAccRewardBlockNumber()
	lastAccReward := sm.db.ReadAccReward(lastNum)

	// Big-step catchup; accumulate until the head block as of now.
	// The head block can be obsolete by the time catchup finished, so the big-step can end up being a bit short.
	// Repeat until the difference is close enough so that the headNum stays the same after one iteration.
	for lastNum < headNum {
		logger.Info("Total supply big step catchup", "last", lastNum, "head", headNum, "minted", lastAccReward.Minted.String(), "burntFee", lastAccReward.BurntFee.String())

		accReward, err := sm.accumulateReward(lastNum, headNum, lastAccReward)
		if err != nil {
			if err != errSupplyManagerQuit {
				logger.Error("Total supply accumulate failed", "from", lastNum, "to", headNum, "err", err)
			}
			return
		}

		lastNum = headNum
		lastAccReward = accReward
		headNum = sm.chain.CurrentBlock().NumberU64()
	}
	logger.Info("Total supply big step catchup done", "last", lastNum, "minted", lastAccReward.Minted.String(), "burntFee", lastAccReward.BurntFee.String())

	// Subscribe to chain head events and accumulate on demand.
	sm.chainHeadSub = sm.chain.SubscribeChainHeadEvent(sm.chainHeadChan)
	for {
		select {
		case <-sm.quitCh:
			return
		case head := <-sm.chainHeadChan:
			headNum = head.Block.NumberU64()

			supply, err := sm.accumulateReward(lastNum, headNum, lastAccReward)
			if err != nil {
				if err != errSupplyManagerQuit {
					logger.Error("Total supply accumulate failed", "from", lastNum, "to", headNum, "err", err)
				}
				return
			}

			lastNum = headNum
			lastAccReward = supply
		}
	}
}

// totalSupplyFromState calculates the ground truth total supply by iterating over all accounts.
// This is extremely inefficient and should only be used for the genesis block and testing.
func (sm *supplyManager) totalSupplyFromState(num uint64) (*big.Int, error) {
	header := sm.chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errors.New("header not found")
	}
	stateDB, err := sm.chain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	dump := stateDB.RawDump() // Extremely inefficient but okay for genesis block.

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

// accumulateReward calculates the total supply from the last block to the current block.
// Given supply at `from` is `fromSupply`, calculate the supply until `to`, inclusive.
func (sm *supplyManager) accumulateReward(from, to uint64, fromAcc *database.AccReward) (*database.AccReward, error) {
	accReward := fromAcc.Copy() // make a copy because we're updating it in-place.

	for num := from + 1; num <= to; num++ {
		// Abort upon quit signal
		if atomic.LoadUint32(&sm.quit) != 0 {
			return nil, errSupplyManagerQuit
		}

		// Accumulate one block
		var (
			header    = sm.chain.GetHeaderByNumber(num)
			rules     = sm.chain.Config().Rules(new(big.Int).SetUint64(num))
			pset, err = sm.gov.EffectiveParams(num)
		)
		if err != nil {
			return nil, err
		}
		blockTotal, err := GetTotalReward(header, rules, pset)
		if err != nil {
			return nil, err
		}
		accReward.Minted.Add(accReward.Minted, blockTotal.Minted)
		accReward.BurntFee.Add(accReward.BurntFee, blockTotal.BurntFee)

		// Store to database, print progress log.
		sm.accRewardCache.Add(num, accReward.Copy())
		if (num % sm.checkpointInterval) == 0 {
			sm.db.WriteAccReward(num, accReward)
			sm.db.WriteLastAccRewardBlockNumber(num)
		}
		if (num % supplyLogInterval) == 0 {
			logger.Info("Accumulated block rewards", "number", num, "minted", accReward.Minted.String(), "burntFee", accReward.BurntFee.String())
		}
	}
	return accReward, nil
}
