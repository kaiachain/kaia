// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/backend/handler.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/p2p"
)

const (
	// chainInitTimeout is the maximum time ValidatePeerType waits for the
	// consensus engine to be ready before rejecting the peer connection.
	chainInitTimeout = 30 * time.Second
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed       = errors.New("fail to decode istanbul message")
	errNoChainReader      = errors.New("sb.chain is nil! --mine option might be missing")
	errInvalidPeerAddress = errors.New("invalid address")
)

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code != consensus.ConsensusMsgCode {
		return false, nil
	}

	event, shouldPost, err := sb.prepareConsensusMessageEvent(addr, msg)
	if err != nil || !shouldPost {
		return true, err
	}

	// Post outside coreMu so a stalled event consumer applies backpressure to
	// this peer without blocking Start or Stop from changing the core lifecycle.
	return true, sb.istanbulEventMux.Post(event)
}

// prepareConsensusMessageEvent validates an incoming consensus envelope and
// records its duplicate-cache state. The event is posted by HandleMsg after
// coreMu is released, so a slow core cannot turn inbound messages into an
// unbounded number of blocked goroutines.
func (sb *backend) prepareConsensusMessageEvent(addr common.Address, msg p2p.Msg) (istanbul.MessageEvent, bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if !sb.coreStarted.Load() {
		return istanbul.MessageEvent{}, false, istanbul.ErrStoppedEngine
	}

	var cmsg bft.ConsensusMsg

	if err := msg.Decode(&cmsg); err != nil {
		return istanbul.MessageEvent{}, false, errDecodeFailed
	}
	data := cmsg.Payload
	hash := istanbul.RLPHash(data)

	// Mark peer's message
	var m *lru.ARCCache
	ms, ok := sb.recentMessages.Get(addr)
	if ok {
		m, _ = ms.(*lru.ARCCache)
	} else {
		m, _ = lru.NewARC(inmemoryMessages)
		sb.recentMessages.Add(addr, m)
	}
	m.Add(hash, true)

	// Mark self known message
	if _, ok := sb.knownMessages.Get(hash); ok {
		return istanbul.MessageEvent{}, false, nil
	}
	sb.knownMessages.Add(hash, true)

	return istanbul.MessageEvent{
		Payload: data,
		Hash:    cmsg.PrevHash,
	}, true, nil
}

func (sb *backend) ValidatePeerType(addr common.Address) error {
	// Wait for engine initialization instead of immediately rejecting peers.
	// P2P server starts before engine.Start() sets sb.chain, so early peer
	// connections would be rejected and must wait for the P2P retry interval
	// (~30s) before reconnecting, delaying the first block consensus.
	// Timeout is a safety net against goroutine leaks if signal is never sent.
	select {
	case <-sb.chainInitCh:
	case <-time.After(chainInitTimeout):
		return errNoChainReader
	}
	if sb.chain == nil {
		return errNoChainReader
	}
	if sb.valsetModule == nil {
		return errInvalidPeerAddress
	}
	num := sb.chain.CurrentHeader().Number.Uint64() + 1
	cnPeers, err := sb.valsetModule.GetCNPeers(num)
	if err != nil {
		sb.logger.Trace("Failed to read CN peers for peer type validation", "addr", addr, "num", num, "err", err)
		return errInvalidPeerAddress
	}
	if cnPeers == nil {
		sb.logger.Trace("CN peer validation disabled", "addr", addr, "num", num)
		return nil
	}
	if valset.NewAddressSet(cnPeers).Contains(addr) {
		sb.logger.Trace("CN peer type validation accepted", "addr", addr, "num", num)
		return nil
	}
	sb.logger.Trace("CN peer type validation rejected", "addr", addr, "num", num)
	return errInvalidPeerAddress
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.broadcaster = broadcaster
	if sb.nodetype == common.CONSENSUSNODE {
		sb.broadcaster.RegisterValidator(common.CONSENSUSNODE, sb)
	}
}

func (sb *backend) NewChainHead() error {
	// Do not take coreMu here. NewChainHead runs on the worker loop, and a
	// coreMu holder can block waiting on that loop, so locking risks a deadlock.
	if !sb.coreStarted.Load() {
		return istanbul.ErrStoppedEngine
	}

	go sb.istanbulEventMux.Post(istanbul.ChainHeadEvent{})
	return nil
}
