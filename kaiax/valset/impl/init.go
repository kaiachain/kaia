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

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ (valset.ValsetModule) = &ValsetModule{}

	istanbulCheckpointInterval = uint64(1024)

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

	validatorVoteBlockNumsCache []uint64
}

func NewValsetModule() *ValsetModule {
	pListCache, _ := lru.New(128)
	rVoteCache, _ := lru.New(128)
	return &ValsetModule{
		proposerListCache: pListCache,
		removeVotesCache:  rVoteCache,
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
	if pMinVoteNum := ReadLowestScannedVoteNum(v.ChainKv); pMinVoteNum == nil {
		// migration not started. Migrating the last interval and leave the rest to be migrated by background thread.
		currentNum := v.Chain.CurrentBlock().NumberU64()
		_, snapshotNum, err := v.getCouncilFromIstanbulSnapshot(currentNum, true)
		if err != nil {
			return err
		}
		if currentNum > 0 {
			// getCouncilFromIstanbulSnapshot() should have scanned until snapshotNum+1.
			writeLowestScannedVoteNum(v.ChainKv, snapshotNum+1)
		} else {
			writeLowestScannedVoteNum(v.ChainKv, 0)
		}
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

	pMinVoteNum := ReadLowestScannedVoteNum(v.ChainKv)
	if pMinVoteNum == nil {
		logger.Error("No lowest scanned snapshot num")
		return
	}

	targetNum := *pMinVoteNum
	for targetNum > 0 {
		if v.quit.Load() == 1 {
			break
		}
		// At each iteration, targetNum should decrease like ... -> 2048 -> 1024 -> 0.
		// get(2048,true) scans [1025, 2048] and returns snapshotNum=1024. So we write lowestScannedVoteNum=1025.
		// get(1024,true) scans [1, 1024] and returns snapshotNum=0. So we write lowestScannedVoteNum=1.
		_, snapshotNum, err := v.getCouncilFromIstanbulSnapshot(targetNum, true)
		if err != nil {
			logger.Error("Failed to migrate", "targetNum", targetNum, "err", err)
			break
		}
		// getCouncilFromIstanbulSnapshot() should have scanned until snapshotNum+1.
		writeLowestScannedVoteNum(v.ChainKv, snapshotNum+1)
		targetNum = snapshotNum
	}

	// Now the migration is complete.
	writeLowestScannedVoteNum(v.ChainKv, 0)
}
