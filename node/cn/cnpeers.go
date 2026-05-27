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

package cn

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const cnPeerHeadChanSize = 16

type cnPeerReader interface {
	GetCNPeers(num uint64) ([]common.Address, error)
}

type cnPeerSink interface {
	SetCNPeers(addrs []common.Address)
}

type cnPeerUpdater struct {
	reader      cnPeerReader
	sink        cnPeerSink
	chainConfig *params.ChainConfig

	lastHash *common.Hash // Last hashed SetCNPeers input; nil means never pushed.
}

type cnPeerInput struct {
	AllCNPeersAllowed bool
	Addrs             []common.Address
}

func newCNPeerUpdater(reader cnPeerReader, sink cnPeerSink, chainConfig *params.ChainConfig) *cnPeerUpdater {
	return &cnPeerUpdater{
		reader:      reader,
		sink:        sink,
		chainConfig: chainConfig,
	}
}

func (u *cnPeerUpdater) syncHead(head *types.Block) {
	if head == nil {
		return
	}
	num := head.NumberU64() + 1
	if !u.shouldSync(num) {
		u.setCNPeersIfChanged(nil)
		return
	}
	u.sync(num)
}

func (u *cnPeerUpdater) shouldSync(num uint64) bool {
	if u == nil || u.chainConfig == nil || u.chainConfig.PermissionlessCompatibleBlock == nil {
		return false
	}
	// Avoid next-block ABv2 reads until the HF block has been committed.
	syncAfter := new(big.Int).Add(u.chainConfig.PermissionlessCompatibleBlock, common.Big1)
	return new(big.Int).SetUint64(num).Cmp(syncAfter) >= 0
}

func (u *cnPeerUpdater) sync(num uint64) {
	if u == nil || u.reader == nil || u.sink == nil {
		return
	}
	addrs, err := u.reader.GetCNPeers(num)
	if err != nil {
		logger.Warn("syncCNPeers: disabling CN peer filter after initial GetCNPeers failed", "num", num, "err", err)
		u.setCNPeersIfChanged(nil)
		return
	}
	u.setCNPeersIfChanged(addrs)
}

func (u *cnPeerUpdater) setCNPeersIfChanged(addrs []common.Address) {
	hash := hashCNPeerInput(addrs)
	if u.lastHash != nil && hash == *u.lastHash {
		return
	}
	u.lastHash = &hash
	u.sink.SetCNPeers(addrs)
}

// hashCNPeerInput hashes `AllCNPeersAllowed || addrs`. The order of `addrs` matters.
func hashCNPeerInput(addrs []common.Address) common.Hash {
	input := cnPeerInput{AllCNPeersAllowed: addrs == nil}
	if addrs != nil {
		input.Addrs = addrs
	}
	encoded, err := rlp.EncodeToBytes(input)
	if err != nil {
		panic(err)
	}
	return crypto.Keccak256Hash(encoded)
}

func (s *CN) startCNPeerSync(srvr p2p.Server) {
	if srvr == nil || s.valsetModule == nil || s.blockchain == nil {
		return
	}
	s.cnPeerUpdater = newCNPeerUpdater(s.valsetModule, srvr, s.chainConfig)

	ch := make(chan blockchain.ChainHeadEvent, cnPeerHeadChanSize)
	sub := s.blockchain.SubscribeChainHeadEvent(ch)
	s.cnPeerHeadSub = sub

	s.cnPeerUpdater.syncHead(s.blockchain.CurrentBlock())

	s.cnPeerSyncWg.Add(1)
	go s.cnPeerSyncLoop(ch, sub)
}

func (s *CN) cnPeerSyncLoop(ch <-chan blockchain.ChainHeadEvent, sub event.Subscription) {
	defer s.cnPeerSyncWg.Done()

	for {
		select {
		case ev := <-ch:
			s.cnPeerUpdater.syncHead(ev.Block)
		case err, ok := <-sub.Err():
			if ok && err != nil {
				logger.Warn("CN peer sync stopped", "err", err)
			}
			return
		}
	}
}

func (s *CN) stopCNPeerSync() {
	if s.cnPeerHeadSub != nil {
		s.cnPeerHeadSub.Unsubscribe()
		s.cnPeerHeadSub = nil
	}
	s.cnPeerSyncWg.Wait()
	s.cnPeerUpdater = nil
}
