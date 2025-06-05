// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package downloader

import (
	"math/big"

	kaia "github.com/kaiachain/kaia/v2"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/staking"
	"github.com/kaiachain/kaia/v2/node/cn/snap"
	"github.com/kaiachain/kaia/v2/params"
)

// fakeDownloader do nothing
type FakeDownloader struct{}

func NewFakeDownloader() *FakeDownloader {
	logger.Warn("downloader is disabled; no data will be downloaded from peers")
	return &FakeDownloader{}
}

func (*FakeDownloader) RegisterPeer(id string, version int, peer Peer) error { return nil }
func (*FakeDownloader) UnregisterPeer(id string) error                       { return nil }

func (*FakeDownloader) DeliverBodies(id string, transactions [][]*types.Transaction) error {
	return nil
}
func (*FakeDownloader) DeliverHeaders(id string, headers []*types.Header) error      { return nil }
func (*FakeDownloader) DeliverNodeData(id string, data [][]byte) error               { return nil }
func (*FakeDownloader) DeliverReceipts(id string, receipts [][]*types.Receipt) error { return nil }
func (*FakeDownloader) DeliverStakingInfos(id string, stakingInfos []*staking.P2PStakingInfo) error {
	return nil
}

func (*FakeDownloader) DeliverSnapPacket(peer *snap.Peer, packet snap.Packet) error {
	return nil
}

func (*FakeDownloader) Terminate()          {}
func (*FakeDownloader) Synchronising() bool { return false }
func (*FakeDownloader) Synchronise(id string, head common.Hash, td *big.Int, mode SyncMode) error {
	return nil
}
func (*FakeDownloader) Progress() kaia.SyncProgress { return kaia.SyncProgress{} }
func (*FakeDownloader) Cancel()                     {}

func (*FakeDownloader) GetSnapSyncer() *snap.Syncer                      { return nil }
func (*FakeDownloader) SyncStakingInfo(id string, from, to uint64) error { return nil }
func (*FakeDownloader) SyncStakingInfoStatus() *SyncingStatus            { return nil }

func (*FakeDownloader) Config() *params.ChainConfig { return params.TestChainConfig }
