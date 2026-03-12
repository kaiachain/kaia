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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

const (
	candidateMsgTimeoutMs   = 500
	vrankCandidateSigDomain = "VRANK_CANDIDATE_V1"

	broadcastChSize    = 2048
	vrankEpoch         = 86400
	maxRound           = 10 // round range [0, 10]
	maxCollectorWindow = 10 // max collection window [N-10, N+10]
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
}

type chain interface {
	CurrentHeader() *types.Header
	GetHeaderByNumber(number uint64) *types.Header
}

type VRankModule struct {
	InitOpts

	broadcastCh   chan *vrank.VRankBroadcastEvent
	broadcastFeed event.Feed
	stopCh        chan struct{}

	nodeID common.Address

	// only for validators
	prepreparedView   istanbul.View // for collection window management
	prepreparedViewMu sync.RWMutex
	collector         *vrank.Collector
}

func NewVRankModule() *VRankModule {
	return &VRankModule{
		broadcastCh: make(chan *vrank.VRankBroadcastEvent, broadcastChSize),
		stopCh:      make(chan struct{}),
		collector:   vrank.NewCollector(),
	}
}

func (v *VRankModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Valset == nil || opts.NodeKey == nil || opts.ChainConfig == nil || opts.ChainConfig.ChainID == nil || opts.Chain == nil {
		return vrank.ErrInitUnexpectedNil
	}
	v.InitOpts = *opts
	v.nodeID = crypto.PubkeyToAddress(opts.NodeKey.PublicKey)
	return nil
}

func (v *VRankModule) Start() error {
	go v.handleBroadcastLoop()
	logger.Info("VRankModule started")

	return nil
}

func (v *VRankModule) Stop() {
	logger.Info("VRankModule stopped")
	close(v.stopCh)
	close(v.broadcastCh)
}

func (v *VRankModule) SubscribeVRank(sink chan<- *vrank.VRankBroadcastEvent) event.Subscription {
	return v.broadcastFeed.Subscribe(sink)
}
