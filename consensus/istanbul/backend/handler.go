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
	"math/big"
	"slices"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/networks/p2p"
)

const (
	IstanbulMsg = 0x11

	// chainInitTimeout is the maximum time ValidatePeerType waits for the
	// consensus engine to be ready before rejecting the peer connection.
	chainInitTimeout = 30 * time.Second
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed       = errors.New("fail to decode istanbul message")
	errNoChainReader      = errors.New("sb.chain is nil! --mine option might be missing")
	errInvalidPeerAddress = errors.New("invalid address")

	// TODO-Kaia-Istanbul: define Versions and Lengths with correct values.
	IstanbulProtocol = consensus.Protocol{
		Name:     "istanbul",
		Versions: []uint{68, 67, 66, 65, 64},
		Lengths:  []uint64{28, 26, 24, 23, 21},
	}
)

// Protocol implements consensus.Engine.Protocol
func (sb *backend) Protocol() consensus.Protocol {
	return IstanbulProtocol
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == IstanbulMsg {
		if !sb.coreStarted {
			return true, istanbul.ErrStoppedEngine
		}

		var cmsg istanbul.ConsensusMsg

		// var data []byte
		if err := msg.Decode(&cmsg); err != nil {
			return true, errDecodeFailed
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
			return true, nil
		}
		sb.knownMessages.Add(hash, true)

		go sb.istanbulEventMux.Post(istanbul.MessageEvent{
			Payload: data,
			Hash:    cmsg.PrevHash,
		})

		return true, nil
	}
	return false, nil
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
	num := sb.chain.CurrentHeader().Number.Uint64() + 1

	// Permissionless: use CNPeers (VA+VR+VP+CR+CT) for P2P connection validation
	if sb.chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		cnPeers, err := sb.valsetModule.GetCNPeers(num)
		if err != nil {
			return err
		}
		if slices.Contains(cnPeers, addr) {
			return nil
		}
		return errInvalidPeerAddress
	}

	valSet, err := sb.GetValidatorSet(num)
	if err != nil {
		return errInvalidPeerAddress
	}
	if valSet.Council().Contains(addr) {
		return nil
	}
	return errInvalidPeerAddress
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *backend) SetBroadcaster(broadcaster consensus.Broadcaster, nodetype common.ConnType) {
	sb.broadcaster = broadcaster
	if nodetype == common.CONSENSUSNODE {
		sb.broadcaster.RegisterValidator(common.CONSENSUSNODE, sb)
	}
}

// RegisterConsensusMsgCode registers the channel of consensus msg.
func (sb *backend) RegisterConsensusMsgCode(peer consensus.Peer) {
	if err := peer.RegisterConsensusMsgCode(IstanbulMsg); err != nil {
		logger.Error("RegisterConsensusMsgCode failed", "err", err)
	}
}

func (sb *backend) NewChainHead() error {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()
	if !sb.coreStarted {
		return istanbul.ErrStoppedEngine
	}

	go sb.istanbulEventMux.Post(istanbul.FinalCommittedEvent{})
	// TODO: Discuss in team whether to enable active CN peer disconnect.
	// Concern: VI/VE nodes lose block sync since CN discovery doesn't cover EN/PN.
	// sb.disconnectNonCNPeers()
	return nil
}

// disconnectNonCNPeers disconnects CN peers that are no longer in the CNPeers set.
func (sb *backend) disconnectNonCNPeers() {
	if sb.broadcaster == nil || sb.valsetModule == nil || sb.chain == nil {
		return
	}
	header := sb.chain.CurrentHeader()
	if header == nil {
		return
	}
	num := header.Number.Uint64()
	if !sb.chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return
	}
	cnPeers, err := sb.valsetModule.GetCNPeers(num)
	if err != nil {
		return
	}
	cnPeerSet := make(map[common.Address]bool, len(cnPeers))
	for _, addr := range cnPeers {
		cnPeerSet[addr] = true
	}
	var toDisconnect []common.Address
	for addr := range sb.broadcaster.GetCNPeers() {
		if !cnPeerSet[addr] {
			toDisconnect = append(toDisconnect, addr)
		}
	}
	if len(toDisconnect) > 0 {
		sb.broadcaster.DisconnectCNPeers(toDisconnect)
	}
}
