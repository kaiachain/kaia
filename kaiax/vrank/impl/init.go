// Copyright 2026 The Kaia Authors
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
	"crypto/ecdsa"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

const (
	candidateMsgTimeoutMs    = 500
	vrankPreprepareSigDomain = "VRANK_PREPREPARE_V1"
	vrankCandidateSigDomain  = "VRANK_CANDIDATE_V1"

	broadcastChSize    = 2048
	vrankEpoch         = vrank.Epoch
	maxRound           = vrank.MaxRound
	maxCollectorWindow = uint64(10) // max collection window [N-10, N+10]

	scoreCacheSize          = 1024
	scoreCheckpointInterval = uint64(vrankEpoch / 8) // 10,800 blocks
)

var (
	_ vrank.VRankModule = &VRankModule{}

	logger = log.NewModuleLogger(log.KaiaxVrank)
)

type InitOpts struct {
	Valset      valset.ValsetModule
	NodeKey     *ecdsa.PrivateKey
	ChainConfig *params.ChainConfig
	Chain       chain
	ChainKv     database.Database
}

type chain interface {
	CurrentHeader() *types.Header
	GetHeaderByNumber(number uint64) *types.Header
}

type VRankModule struct {
	InitOpts

	lifecycleMu sync.Mutex
	running     bool

	broadcastCh   chan *vrank.VRankBroadcastEvent
	broadcastFeed event.Feed
	stopCh        chan struct{}

	nodeID common.Address

	// only for validators
	prepreparedView   istanbul.View // for collection window management
	prepreparedViewMu sync.RWMutex
	collector         *vrank.Collector

	pfsCache      *lru.ARCCache // map[proposer]score
	cpMatrixCache *lru.ARCCache // map[candidate][proposer]score
}

func NewVRankModule() *VRankModule {
	pfsCache, _ := lru.NewARC(scoreCacheSize)
	cpMatrixCache, _ := lru.NewARC(scoreCacheSize)
	return &VRankModule{
		broadcastCh:   make(chan *vrank.VRankBroadcastEvent, broadcastChSize),
		stopCh:        make(chan struct{}),
		collector:     vrank.NewCollector(),
		pfsCache:      pfsCache,
		cpMatrixCache: cpMatrixCache,
	}
}

func (v *VRankModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Valset == nil || opts.NodeKey == nil || opts.ChainConfig == nil || opts.ChainConfig.ChainID == nil || opts.Chain == nil || opts.ChainKv == nil {
		return vrank.ErrInitUnexpectedNil
	}
	v.InitOpts = *opts
	v.nodeID = crypto.PubkeyToAddress(opts.NodeKey.PublicKey)
	if err := v.catchUpScoreCaches(); err != nil {
		return err
	}
	return nil
}

func (v *VRankModule) Start() error {
	v.lifecycleMu.Lock()
	defer v.lifecycleMu.Unlock()

	if v.running {
		return nil
	}

	v.stopCh = make(chan struct{})
	v.running = true
	go v.handleBroadcastLoop(v.stopCh)
	logger.Info("VRankModule started")

	return nil
}

func (v *VRankModule) Stop() {
	v.lifecycleMu.Lock()
	if !v.running {
		v.lifecycleMu.Unlock()
		return
	}
	stopCh := v.stopCh
	v.running = false
	v.lifecycleMu.Unlock()

	logger.Info("VRankModule stopped")
	close(stopCh)
}

func (v *VRankModule) SubscribeVRank(sink chan<- *vrank.VRankBroadcastEvent) event.Subscription {
	return v.broadcastFeed.Subscribe(sink)
}

func (v *VRankModule) catchUpScoreCaches() error {
	v.pfsCache.Purge()
	v.cpMatrixCache.Purge()

	head := v.Chain.CurrentHeader()
	if head == nil || head.Number == nil {
		return nil
	}
	if !v.ChainConfig.IsPermissionlessForkEnabled(head.Number) {
		return nil
	}
	headNum := head.Number.Uint64()

	if err := v.catchUp(headNum); err != nil {
		return err
	}
	return nil
}

func (v *VRankModule) catchUp(headNum uint64) error {
	epochStart := calcEpochStart(headNum)

	var (
		start    uint64
		pfs      map[common.Address]uint64
		cpMatrix map[common.Address]map[common.Address]uint64
	)

	if cpNum, storedPFS, storedCpMatrix, ok := v.loadCheckpointInEpoch(headNum); ok {
		start = cpNum + 1
		pfs = storedPFS
		cpMatrix = storedCpMatrix
		v.pfsCache.Add(cpNum, cloneMap(storedPFS))
		v.cpMatrixCache.Add(cpNum, cloneCPMatrix(storedCpMatrix))
	} else {
		start = epochStart
		pfs = make(map[common.Address]uint64)
		var err error
		cpMatrix, err = v.newCPMatrix(epochStart)
		if err != nil {
			return err
		}
	}

	for blockNum := start; blockNum <= headNum; blockNum++ {
		var err error
		pfs, err = v.applyBlocksForPFS(blockNum, blockNum, pfs)
		if err != nil {
			logger.Error("Failed to compute PFS", "blockNum", blockNum, "err", err)
			return err
		}
		v.pfsCache.Add(blockNum, cloneMap(pfs))

		cpMatrix, err = v.applyBlocksForCPMatrix(blockNum, blockNum, cpMatrix)
		if err != nil {
			logger.Error("Failed to compute CFS", "blockNum", blockNum, "err", err)
			return err
		}
		v.cpMatrixCache.Add(blockNum, cloneCPMatrix(cpMatrix))

		if blockNum%scoreCheckpointInterval == 0 {
			WriteCheckpoint(v.ChainKv, blockNum, pfs, cpMatrix)
			WriteLastCheckpoint(v.ChainKv, blockNum)
		}

		if blockNum%10000 == 0 && start != headNum {
			logger.Info("VRank module catchup progress", "blockNum", blockNum, "progress", float64(blockNum-start)/float64(headNum-start))
		}
	}
	return nil
}

func (v *VRankModule) loadCheckpointInEpoch(blockNum uint64) (uint64, map[common.Address]uint64, map[common.Address]map[common.Address]uint64, bool) {
	cpNum := calcCheckpointBlock(blockNum)
	pfs, cpMatrix := ReadCheckpoint(v.ChainKv, cpNum)
	if pfs == nil {
		return 0, nil, nil, false
	}
	return cpNum, pfs, cpMatrix, true
}
