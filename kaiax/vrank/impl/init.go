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
	"time"

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
	candidatePrepareDeadlineMs = 200 * time.Millisecond

	VRankPreprepareMsg = 0x17
	VRankCandidateMsg  = 0x18

	broadcastChSize = 2048
)

var (
	_ vrank.VRankModule = &VRankModule{}

	logger = log.NewModuleLogger(log.KaiaxVrank)
)

type InitOpts struct {
	Valset      valset.ValsetModule
	NodeKey     *ecdsa.PrivateKey
	ChainConfig *params.ChainConfig
}

type VRankModule struct {
	InitOpts

	broadcastCh   chan *vrank.BroadcastRequest
	broadcastFeed event.Feed

	nodeId common.Address

	// only for validators
	prepreparedTime      time.Time
	prepreparedView      istanbul.View
	prepreparedBlockHash common.Hash
	candResponses        sync.Map // map[common.Address]time.Duration
}

func NewVRankModule() *VRankModule {
	return &VRankModule{
		broadcastCh: make(chan *vrank.BroadcastRequest, broadcastChSize),
	}
}

func (v *VRankModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Valset == nil || opts.NodeKey == nil || opts.ChainConfig == nil {
		return vrank.ErrInitUnexpectedNil
	}
	v.InitOpts = *opts
	v.nodeId = crypto.PubkeyToAddress(opts.NodeKey.PublicKey)
	return nil
}

func (v *VRankModule) Start() error {
	go v.handleBroadcastLoop()
	logger.Info("VRankModule started")

	return nil
}

func (v *VRankModule) Stop() {
	logger.Info("VRankModule stopped")
}
