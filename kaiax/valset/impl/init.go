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
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
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

	// cache for validatorVoteBlockNums
	validatorVoteBlockNumsCache []uint64
}

func NewValsetModule() *ValsetModule {
	cache, _ := lru.New(128)
	return &ValsetModule{
		proposerListCache: cache,
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
		// getCouncilFromIstanbulSnapshot() should have scanned until snapshotNum+1.
		writeLowestScannedVoteNum(v.ChainKv, snapshotNum+1)
	}

	return nil
}

func (v *ValsetModule) Start() error {
	logger.Info("ValsetModule Started")

	// Reset all caches
	v.proposerListCache.Purge()
	v.validatorVoteBlockNumsCache = nil

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

// getCouncilFromIstanbulSnapshot re-generates the council at the given targetNum.
// Returns the council at targetNum, the nearest snapshot number, and error if any.
//
// The council is calculated from the nearest istanbul snapshot (at snapshotNum)
// plus the validator votes in the range [snapshotNum+1, targetNum-1]. Note that
// snapshot at snapshotNum already reflects the validator vote at snapshotNum,
// so we apply the votes starting from snapshotNum+1.
//
// If write is true, ValidatorVoteBlockNums and Council in the extended range
// [snapshotNum+1, targetNum] are written to the database. Note that this time
// the targetNum is included in the range for completeness. This property is
// useful for snapshot interval-wise migration.
func (v *ValsetModule) getCouncilFromIstanbulSnapshot(targetNum uint64, write bool) (*valset.AddressSet, uint64, error) {
	if targetNum == 0 {
		council, err := v.getCouncilGenesis()
		return council, 0, err
	}
	// Load council at the nearest istanbul snapshot. This is the result
	// applying the votes up to the snapshotNum.
	snapshotNum := roundDown(targetNum-1, istanbulCheckpointInterval)
	header := v.Chain.GetHeaderByNumber(snapshotNum)
	if header == nil {
		return nil, 0, errNoHeader
	}
	council := valset.NewAddressSet(ReadIstanbulSnapshot(v.ChainKv, header.Hash()))
	if council.Len() == 0 {
		return nil, 0, ErrNoIstanbulSnapshot(snapshotNum)
	}

	// Apply the votes in the interval [snapshotNum+1, targetNum-1].
	for n := snapshotNum + 1; n < targetNum; n++ {
		if err := v.replayBlock(council, n, write); err != nil {
			return nil, 0, err
		}
	}
	// Apply the vote at targetNum to write to database, but do not return the modified council.
	if write {
		// Apply the vote at targetNum and write to database, but do not affect the returning council.
		if err := v.replayBlock(council.Copy(), targetNum, true); err != nil {
			return nil, 0, err
		}
	}
	return council, snapshotNum, nil
}

func (v *ValsetModule) replayBlock(council *valset.AddressSet, num uint64, write bool) error {
	header := v.Chain.GetHeaderByNumber(num)
	if header == nil {
		return errNoHeader
	}
	governingNode := v.GovModule.EffectiveParamSet(num).GoverningNode
	if applyVote(header, council, governingNode) && write {
		insertValidatorVoteBlockNums(v.ChainKv, num)
		writeCouncil(v.ChainKv, num, council.List())
	}
	return nil
}

// applyVote modifies the given council *in-place* by the validator vote in the given header.
// governingNode, if specified, is not affected by the vote.
// Returns true if the council is modified, false otherwise.
func applyVote(header *types.Header, council *valset.AddressSet, governingNode common.Address) bool {
	voteKey, addresses, ok := parseValidatorVote(header)
	if !ok {
		return false
	}

	originalSize := council.Len()
	for _, address := range addresses {
		if address == governingNode {
			continue
		}
		switch voteKey {
		case gov.AddValidator:
			if !council.Contains(address) {
				council.Add(address)
			}
		case gov.RemoveValidator:
			if council.Contains(address) {
				council.Remove(address)
			}
		}
	}
	return originalSize != council.Len()
}

func parseValidatorVote(header *types.Header) (gov.ValidatorParamName, []common.Address, bool) {
	// Check that a vote exists and is a validator vote.
	voteBytes := headergov.VoteBytes(header.Vote)
	if len(voteBytes) == 0 {
		return "", nil, false
	}
	vote, err := voteBytes.ToVoteData()
	if err != nil {
		return "", nil, false
	}
	voteKey := gov.ValidatorParamName(vote.Name())
	_, isValidatorVote := gov.ValidatorParams[voteKey]
	if !isValidatorVote {
		return "", nil, false
	}

	// Type cast the vote value. It can be a single address or a list of addresses.
	var addresses []common.Address
	switch voteValue := vote.Value().(type) {
	case common.Address:
		addresses = []common.Address{voteValue}
	case []common.Address:
		addresses = voteValue
	default:
		return "", nil, false
	}

	return voteKey, addresses, true
}

func roundDown(n, p uint64) uint64 {
	return n - (n % p)
}

// TODO-kaiax: move the feature into staking_info.go
func collectStakingAmounts(nodes []common.Address, si *staking.StakingInfo) map[common.Address]float64 {
	cns := si.ConsolidatedNodes()
	stakingAmounts := make(map[common.Address]float64, len(nodes))
	for _, node := range nodes {
		stakingAmounts[node] = 0
	}
	for _, cn := range cns {
		for _, node := range cn.NodeIds {
			if _, ok := stakingAmounts[node]; ok {
				stakingAmounts[node] = float64(cn.StakingAmount)
			}
		}
	}
	return stakingAmounts
}
