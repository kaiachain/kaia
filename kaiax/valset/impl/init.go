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
	"sync"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/consensus"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/staking"
	"github.com/kaiachain/kaia/v2/kaiax/valset"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
)

var (
	_ (valset.ValsetModule) = &ValsetModule{}

	// Istanbul snapshots are generated every 1024 blocks.
	// Note: snapshots are not persisted after the migration process.
	istanbulCheckpointInterval = uint64(1024)

	// A background migration thread transfers Istanbul snapshots to the valset council.
	// To prevent high CPU usage, the migration loop is throttled with a 50ms delay per iteration.
	// For example, if migration starts at block 180,000,000, the entire process will take at least 3.2 hours.
	migrationThrottlingDelay = 50 * time.Millisecond
	migrateLogInterval       = uint64(102400)

	logger = log.NewModuleLogger(log.KaiaxValset)
)

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	CurrentBlock() *types.Block
	Config() *params.ChainConfig
	Engine() consensus.Engine
}

type InitOpts struct {
	ChainKv       database.Database
	Chain         chain
	GovModule     gov.GovModule
	StakingModule staking.StakingModule
}

type ValsetModule struct {
	InitOpts

	quit atomic.Int32 // stop the migration thread
	wg   sync.WaitGroup

	// cache for weightedRandom and uniformRandom proposerLists.
	proposerListCache *lru.Cache // uint64 -> []common.Address
	removeVotesCache  *lru.Cache // uint64 -> removeVoteList
	councilCache      *lru.Cache // uint64 -> *valset.AddressSet

	validatorVoteBlockNumsCache []uint64
	lowestScannedVoteNumCache   *uint64
}

func NewValsetModule() *ValsetModule {
	pListCache, _ := lru.New(128)
	rVoteCache, _ := lru.New(128)
	councilCache, _ := lru.New(128)
	return &ValsetModule{
		proposerListCache: pListCache,
		removeVotesCache:  rVoteCache,
		councilCache:      councilCache,
	}
}

func (v *ValsetModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainKv == nil || opts.Chain == nil || opts.GovModule == nil || opts.StakingModule == nil {
		return errInitUnexpectedNil
	}
	v.InitOpts = *opts

	return v.initSchema()
}

func (v *ValsetModule) initSchema() error {
	// Ensure mandatory schema at block 0
	if voteBlockNums := ReadValidatorVoteBlockNums(v.ChainKv); voteBlockNums == nil {
		writeValidatorVoteBlockNums(v.ChainKv, []uint64{0})
		v.validatorVoteBlockNumsCache = nil
	}
	if council := ReadCouncil(v.ChainKv, 0); council == nil {
		genesisCouncil, err := v.getCouncilGenesis()
		if err != nil {
			return err
		}
		writeCouncil(v.ChainKv, 0, genesisCouncil.List())
	}

	// Ensure mandatory schema lowestScannedCheckpointInterval
	if pMinVoteNum := v.readLowestScannedVoteNumCached(); pMinVoteNum == nil {
		// migration not started. Migrating the last interval and leave the rest to be migrated by background thread.
		currentNum := v.Chain.CurrentBlock().NumberU64()
		_, snapshotNum, err := v.getCouncilFromIstanbulSnapshot(currentNum, true)
		if err != nil {
			return err
		}
		logger.Info("ValsetModule migrate latest interval", "currentNum", currentNum, "snapshotNum", snapshotNum)
		if currentNum > 0 {
			// getCouncilFromIstanbulSnapshot() should have scanned until snapshotNum+1.
			writeLowestScannedVoteNum(v.ChainKv, snapshotNum+1)
		} else {
			writeLowestScannedVoteNum(v.ChainKv, 0)
		}
		v.lowestScannedVoteNumCache = nil
	}

	return nil
}

func (v *ValsetModule) Start() error {
	logger.Info("ValsetModule Started")

	// Reset all caches
	v.proposerListCache.Purge()
	v.removeVotesCache.Purge()

	// Reset the quit state.
	v.quit.Store(0)
	v.wg.Add(1)
	go v.migrate()
	return nil
}

func (v *ValsetModule) Stop() {
	logger.Info("ValsetModule Stopped")
	v.quit.Store(1)
	v.wg.Wait()
}

func (v *ValsetModule) migrate() {
	defer v.wg.Done()

	pMinVoteNum := v.readLowestScannedVoteNumCached()
	if pMinVoteNum == nil {
		logger.Error("No lowest scanned snapshot num")
		return
	}

	targetNum := *pMinVoteNum
	logger.Info("ValsetModule migrate start", "targetNum", targetNum)

	for targetNum > 0 {
		if v.quit.Load() == 1 {
			break
		}

		time.Sleep(migrationThrottlingDelay)

		// At each iteration, targetNum should decrease like ... -> 2048 -> 1024 -> 0.
		// get(2048,true) scans [1025, 2048] and returns snapshotNum=1024. So we write lowestScannedVoteNum=1025.
		// get(1024,true) scans [1, 1024] and returns snapshotNum=0. So we write lowestScannedVoteNum=1.
		_, snapshotNum, err := v.getCouncilFromIstanbulSnapshot(targetNum, true)
		if err != nil {
			logger.Error("Failed to migrate", "targetNum", targetNum, "err", err)
			break
		}
		if targetNum%migrateLogInterval == 0 {
			logger.Info("ValsetModule migrate", "targetNum", targetNum, "snapshotNum", snapshotNum)
		}
		// getCouncilFromIstanbulSnapshot() should have scanned until snapshotNum+1.
		writeLowestScannedVoteNum(v.ChainKv, snapshotNum+1)
		v.lowestScannedVoteNumCache = nil
		targetNum = snapshotNum
	}

	if targetNum == 0 {
		logger.Info("ValsetModule migrate complete")
		// Now the migration is complete.
		writeLowestScannedVoteNum(v.ChainKv, 0)
		v.lowestScannedVoteNumCache = nil
	}
}
