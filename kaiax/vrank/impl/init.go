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
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
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

	broadcastChSize = 2048
	maxRound        = vrank.MaxRound
	maxWindow       = uint64(10) // max collection window [N-10, N+10]

	scoreCacheSize = 1024
)

var (
	_ vrank.VRankModule = &VRankModule{}

	logger = log.NewModuleLogger(log.KaiaxVrank)
)

type InitOpts struct {
	Valset      valset.ValsetModule
	Randao      vrank.BlsPubkeyGetter
	RoundReader vrank.RoundReader
	NodeKey     *ecdsa.PrivateKey
	BlsKey      bls.SecretKey
	ChainConfig *params.ChainConfig
	Chain       vrank.Chain
	ChainKv     database.Database
}

type VRankModule struct {
	InitOpts

	lifecycleMu sync.Mutex
	running     bool

	broadcastCh   chan *vrank.VRankBroadcastEvent
	broadcastFeed event.Feed
	stopCh        chan struct{}

	nodeID common.Address

	// only for the proposer, which collects VRankCandidate replies to its own VRankPreprepare
	collector *vrank.Collector

	// ownProposals: blocks this node proposed this epoch, pending report. Their collector views are
	// kept past the prune window until reported. In-memory — a restart drops them (fail-safe).
	ownProposals   map[uint64]struct{}
	ownProposalsMu sync.Mutex

	// only for candidates
	seenPreprepare   map[vrank.ViewKey]common.Hash
	seenPreprepareMu sync.Mutex

	pfsCache      *lru.ARCCache // map[proposer]score
	cpMatrixCache *lru.ARCCache // map[candidate][proposer]score

	// TODO-Permissionless: Testing only. Remove before production release.
	skipCandidate atomic.Bool
}

func NewVRankModule() *VRankModule {
	pfsCache, _ := lru.NewARC(scoreCacheSize)
	cpMatrixCache, _ := lru.NewARC(scoreCacheSize)
	return &VRankModule{
		broadcastCh:    make(chan *vrank.VRankBroadcastEvent, broadcastChSize),
		stopCh:         make(chan struct{}),
		collector:      vrank.NewCollector(),
		ownProposals:   make(map[uint64]struct{}),
		seenPreprepare: make(map[vrank.ViewKey]common.Hash),
		pfsCache:       pfsCache,
		cpMatrixCache:  cpMatrixCache,
	}
}

// TODO-Permissionless: Testing only. Remove before production release.
func (v *VRankModule) SetSkipCandidate(skip bool) {
	v.skipCandidate.Store(skip)
}

func (v *VRankModule) vrankEpoch() uint64 {
	return v.ChainConfig.VRankEpoch
}

// scoreCheckpointInterval returns the block interval at which PFS/CFS scores are persisted to disk.
// Divides VRankEpoch into 8 segments so that cold-start replay covers at most 1/8 of an epoch.
func (v *VRankModule) scoreCheckpointInterval() uint64 {
	return v.ChainConfig.VRankEpoch / 8
}

func (v *VRankModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Valset == nil || opts.Randao == nil || opts.RoundReader == nil || opts.NodeKey == nil || opts.BlsKey == nil || opts.ChainConfig == nil || opts.ChainConfig.ChainID == nil || opts.Chain == nil || opts.ChainKv == nil {
		return vrank.ErrInitUnexpectedNil
	}
	v.InitOpts = *opts
	v.nodeID = crypto.PubkeyToAddress(opts.NodeKey.PublicKey)
	if err := v.catchUpScoreCaches(); err != nil {
		logger.Warn("Failed to catch up score caches, starting cold", "err", err)
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

	if _, err := v.GetPFS(headNum); err != nil {
		return err
	}
	if _, err := v.getCPMatrix(headNum); err != nil {
		return err
	}
	return nil
}

func (v *VRankModule) loadPFSCheckpointInEpoch(blockNum uint64) (uint64, map[common.Address]uint64, bool) {
	cpNum := calcCheckpointBlock(blockNum, v.scoreCheckpointInterval())
	pfs := ReadCheckpointPFS(v.ChainKv, cpNum)
	if pfs == nil {
		return 0, nil, false
	}
	return cpNum, pfs, true
}

func (v *VRankModule) loadCPMatrixCheckpointInEpoch(blockNum uint64) (uint64, vrank.CPMatrix, bool) {
	cpNum := calcCheckpointBlock(blockNum, v.scoreCheckpointInterval())
	cpMatrix := ReadCheckpointCPMatrix(v.ChainKv, cpNum)
	if cpMatrix == nil {
		return 0, nil, false
	}
	return cpNum, cpMatrix, true
}
