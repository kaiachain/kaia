// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
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

package cn

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/gxhash"
	consensusmocks "github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/node/cn/mocks"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	workmocks "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

const blockNum1 = 20190902

var td1 = big.NewInt(123)

const numVals = 6

var (
	addrs    []common.Address
	keys     []*ecdsa.PrivateKey
	nodeids  []discover.NodeID
	p2pPeers []*p2p.Peer
	blocks   []*types.Block
	hashes   []common.Hash
)

var (
	tx1 *types.Transaction
	txs types.Transactions
)

var hash1 common.Hash

func init() {
	addrs = make([]common.Address, numVals)
	keys = make([]*ecdsa.PrivateKey, numVals)
	nodeids = make([]discover.NodeID, numVals)
	p2pPeers = make([]*p2p.Peer, numVals)
	blocks = make([]*types.Block, numVals)
	hashes = make([]common.Hash, numVals)

	for i := range keys {
		keys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
		nodeids[i] = discover.PubkeyID(&keys[i].PublicKey)
		p2pPeers[i] = p2p.NewPeer(nodeids[i], nodeids[i].String(), []p2p.Cap{})
		blocks[i] = newBlock(i)
		hashes[i] = blocks[i].Hash()
	}

	signer := types.MakeSigner(params.BFTTestChainConfig, big.NewInt(2019))
	tx1 = types.NewTransaction(111, addrs[0], big.NewInt(111), 111, big.NewInt(111), addrs[0][:])

	tx1.Sign(signer, keys[0])
	tx1.Size()
	txs = types.Transactions{tx1}

	hash1 = tx1.Hash()
}

func newMocks(t *testing.T) (*gomock.Controller, *consensusmocks.MockEngine, *workmocks.MockBlockChain, *workmocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	mockEngine := consensusmocks.NewMockEngine(mockCtrl)
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	mockTxPool := workmocks.NewMockTxPool(mockCtrl)

	return mockCtrl, mockEngine, mockBlockChain, mockTxPool
}

func newBlock(blockNum int) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(blockNum)),
		BlockScore: big.NewInt(int64(1)),
		Extra:      addrs[0][:],
		Governance: addrs[0][:],
		Vote:       addrs[0][:],
		BaseFee:    big.NewInt(int64(params.DefaultLowerBoundBaseFee)),
	}
	header.Hash()
	block := types.NewBlockWithHeader(header)
	block = block.WithBody(types.Transactions{})
	block.Hash()
	block.Size()
	block.BlockScore()
	return block
}

func newReceipt(gasUsed int) *types.Receipt {
	rct := types.NewReceipt(uint(gasUsed), common.Hash{}, uint64(gasUsed))
	rct.Logs = []*types.Log{}
	rct.Bloom = types.Bloom{}
	return rct
}

func TestNewProtocolManager(t *testing.T) {
	// 1. If consensus.Engine returns an empty Protocol, NewProtocolManager throws an error.
	{
		mockCtrl, mockEngine, mockBlockChain, mockTxPool := newMocks(t)
		defer mockCtrl.Finish()

		block := newBlock(blockNum1)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)
		mockEngine.EXPECT().Protocol().Return(consensus.Protocol{}).Times(1)

		pm, err := NewProtocolManager(nil, downloader.FastSync, 0, nil, mockTxPool,
			mockEngine, mockBlockChain, nil, 1, -1, &Config{})

		assert.Nil(t, pm)
		assert.Equal(t, errIncompatibleConfig, err)
	}
}

func TestProtocolManager_RegisterValidator(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPeerSet := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeerSet

	val := &ByPassValidator{}

	mockPeerSet.EXPECT().RegisterValidator(common.CONSENSUSNODE, val).Times(1)
	mockPeerSet.EXPECT().RegisterValidator(common.ENDPOINTNODE, val).Times(1)
	mockPeerSet.EXPECT().RegisterValidator(common.PROXYNODE, val).Times(1)
	mockPeerSet.EXPECT().RegisterValidator(common.BOOTNODE, val).Times(1)
	mockPeerSet.EXPECT().RegisterValidator(common.UNKNOWNNODE, val).Times(1)

	pm.RegisterValidator(common.CONSENSUSNODE, val)
	pm.RegisterValidator(common.ENDPOINTNODE, val)
	pm.RegisterValidator(common.PROXYNODE, val)
	pm.RegisterValidator(common.BOOTNODE, val)
	pm.RegisterValidator(common.UNKNOWNNODE, val)
}

func TestProtocolManager_getWSEndPoint(t *testing.T) {
	pm := &ProtocolManager{}

	ws1 := "abc"
	ws2 := "123"

	pm.SetWsEndPoint(ws1)
	assert.Equal(t, ws1, pm.getWSEndPoint())

	pm.SetWsEndPoint(ws2)
	assert.Equal(t, ws2, pm.getWSEndPoint())
}

func TestProtocolManager_removePeer(t *testing.T) {
	peerID := nodeids[0].String()

	{
		pm := &ProtocolManager{}
		mockCtrl := gomock.NewController(t)

		mockPeerSet := NewMockPeerSet(mockCtrl)
		pm.peers = mockPeerSet

		mockPeerSet.EXPECT().Peer(peerID).Return(nil).Times(1)
		pm.removePeer(peerID)

		mockCtrl.Finish()
	}

	{
		pm := &ProtocolManager{}
		mockCtrl := gomock.NewController(t)

		mockPeerSet := NewMockPeerSet(mockCtrl)
		pm.peers = mockPeerSet

		mockPeer := NewMockPeer(mockCtrl)

		mockDownloader := mocks.NewMockProtocolManagerDownloader(mockCtrl)
		mockDownloader.EXPECT().UnregisterPeer(peerID).Times(1)
		pm.downloader = mockDownloader

		// Return
		mockPeer.EXPECT().ExistSnapExtension().Return(false).Times(1)

		mockPeerSet.EXPECT().Unregister(peerID).Return(expectedErr).Times(1)

		mockPeer.EXPECT().GetP2PPeer().Return(p2pPeers[0]).Times(1)

		mockPeerSet.EXPECT().Peer(peerID).Return(mockPeer).Times(1)
		pm.removePeer(peerID)

		mockCtrl.Finish()
	}

	{
		pm := &ProtocolManager{}
		mockCtrl := gomock.NewController(t)

		mockPeerSet := NewMockPeerSet(mockCtrl)
		pm.peers = mockPeerSet

		mockPeer := NewMockPeer(mockCtrl)

		mockDownloader := mocks.NewMockProtocolManagerDownloader(mockCtrl)
		mockDownloader.EXPECT().UnregisterPeer(peerID).Times(1)
		pm.downloader = mockDownloader

		// Return
		mockPeer.EXPECT().ExistSnapExtension().Return(false).Times(1)

		mockPeerSet.EXPECT().Unregister(peerID).Return(nil).Times(1)

		mockPeer.EXPECT().GetP2PPeer().Return(p2pPeers[0]).Times(1)

		mockPeerSet.EXPECT().Peer(peerID).Return(mockPeer).Times(1)
		pm.removePeer(peerID)

		mockCtrl.Finish()
	}
}

func TestProtocolManager_getChainID(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := &params.ChainConfig{ChainID: big.NewInt(12345)}

	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().Config().Return(cfg).AnyTimes()
	pm.blockchain = mockBlockChain

	assert.Equal(t, cfg.ChainID, pm.getChainID())
}

func TestProtocolManager_processMsg_panicRecover(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msgCh := make(chan p2p.Msg)
	errCh := make(chan error)
	addr := common.Address{}

	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetVersion().Do(
		func() { panic("panic test") },
	)

	// pm.processMsg will be panicked by the mockPeer
	go pm.processMsg(msgCh, mockPeer, addr, errCh)

	msgCh <- p2p.Msg{Code: NodeDataMsg}

	// panic will be recovered and errCh will receive an error
	err := <-errCh
	assert.Equal(t, errUnknownProcessingError, err)
}

func TestSampleSize(t *testing.T) {
	peers := make([]Peer, minNumPeersToSendBlock-1)
	assert.Equal(t, len(peers), sampleSize(peers))

	peers = make([]Peer, 4)
	assert.Equal(t, minNumPeersToSendBlock, sampleSize(peers))

	peers = make([]Peer, 16)
	assert.Equal(t, 4, sampleSize(peers))
}

func TestSamplingPeers(t *testing.T) {
	peers := make([]Peer, 10)
	assert.Equal(t, peers, samplingPeers(peers, 20))
	assert.Equal(t, peers[:5], samplingPeers(peers, 5))
}

func TestBroadcastBlock_NoParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := int64(100)
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(nil).Times(1)
	mockBlockChain.EXPECT().GetTd(block.ParentHash(), block.NumberU64()-1).Return(big.NewInt(td)).Times(0)
	pm.blockchain = mockBlockChain

	mockPeers := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeers

	mockPeer := NewMockPeer(mockCtrl)
	mockPeers.EXPECT().SamplePeersToSendBlock(block, pm.nodetype).Return([]Peer{mockPeer}).Times(0)
	mockPeer.EXPECT().AsyncSendNewBlock(block, new(big.Int).Add(block.BlockScore(), big.NewInt(td))).Times(0)

	pm.BroadcastBlock(block)
}

func TestBroadcastBlock_ParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := int64(100)
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(block).Times(1)
	mockBlockChain.EXPECT().GetTd(block.ParentHash(), block.NumberU64()-1).Return(big.NewInt(td)).Times(1)
	pm.blockchain = mockBlockChain

	mockPeers := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeers

	mockPeer := NewMockPeer(mockCtrl)
	mockPeers.EXPECT().SamplePeersToSendBlock(block, pm.nodetype).Return([]Peer{mockPeer}).Times(1)
	mockPeer.EXPECT().AsyncSendNewBlock(block, new(big.Int).Add(block.BlockScore(), big.NewInt(td))).Times(1)

	pm.BroadcastBlock(block)
}

func TestBroadcastBlockHash(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// When the given block doesn't exist.
	{
		mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(false).Times(1)
		pm.blockchain = mockBlockChain
		pm.BroadcastBlockHash(block)
	}

	// When the given block exists.
	{
		mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(true).Times(1)
		pm.blockchain = mockBlockChain

		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().AsyncSendNewBlockHash(block).Times(1)

		mockPeers := NewMockPeerSet(mockCtrl)
		mockPeers.EXPECT().PeersWithoutBlock(block.Hash()).Return([]Peer{mockPeer}).Times(1)
		pm.peers = mockPeers

		pm.BroadcastBlockHash(block)
	}
}

func TestProtocolManager_txBroadcastLoop_FromCN_CN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	txsCh := make(chan blockchain.NewTxsEvent, txChanSize)
	pm.txsCh = txsCh

	feed := &event.Feed{}
	pm.txsSub = feed.Subscribe(txsCh)

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	// Using gomock.Eq(txs) for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)

	go pm.txBroadcastLoop()

	txsCh <- blockchain.NewTxsEvent{Txs: txs}

	time.Sleep(500 * time.Millisecond)

	pm.txsSub.Unsubscribe()
}

func TestBroadcastTxsFromCN_CN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	// Using gomock.Eq(txs) for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromCN_CN_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	// Using gomock.Eq(txs) for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromPN_PN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.PROXYNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(common.CONSENSUSNODE).Times(1)
	pnPeer.EXPECT().ConnType().Return(common.PROXYNODE).Times(1)
	enPeer.EXPECT().ConnType().Return(common.ENDPOINTNODE).Times(1)

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromPN_PN_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.PROXYNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(common.CONSENSUSNODE).Times(1)
	pnPeer.EXPECT().ConnType().Return(common.PROXYNODE).Times(1)
	enPeer.EXPECT().ConnType().Return(common.ENDPOINTNODE).Times(1)

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromEN_ALL_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().ConnType().Return(common.CONSENSUSNODE).AnyTimes()
	pnPeer.EXPECT().ConnType().Return(common.PROXYNODE).AnyTimes()
	enPeer.EXPECT().ConnType().Return(common.ENDPOINTNODE).AnyTimes()

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	enPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromEN_ALL_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().ConnType().Return(common.CONSENSUSNODE).AnyTimes()
	pnPeer.EXPECT().ConnType().Return(common.PROXYNODE).AnyTimes()
	enPeer.EXPECT().ConnType().Return(common.ENDPOINTNODE).AnyTimes()

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	enPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFrom_DefaultCase(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.BOOTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	// There are no expected calls for the mocks.
	pm.nodetype = common.BOOTNODE
	pm.BroadcastTxs(txs)

	pm.nodetype = common.UNKNOWNNODE
	pm.BroadcastTxs(txs)
}

func TestProtocolManager_txResendLoop(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	pm.quitResendCh = make(chan struct{})

	maxTxCount := 100
	mockTxPool := workmocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().CachedPendingTxsByCount(maxTxCount).Return(txs).Times(1)

	pm.txpool = mockTxPool

	go pm.txResendLoop(1, maxTxCount)

	time.Sleep(1500 * time.Millisecond)

	pm.quitResendCh <- struct{}{}
}

func TestProtocolManager_txResend(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	pm.txResend(txs)
}

func TestReBroadcastTxs_CN(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	pm.ReBroadcastTxs(txs)
}

func TestReBroadcastTxs_PN(t *testing.T) {
	// CN Peer=0, PN Peer=1
	{
		pm := &ProtocolManager{}
		pm.nodetype = common.PROXYNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(common.PROXYNODE).Times(2)
		enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.enpeers[addrs[2]] = enPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
	// CN Peer=1, PN Peer=0
	{
		pm := &ProtocolManager{}
		pm.nodetype = common.PROXYNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(common.CONSENSUSNODE).Times(1)
		pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.pnpeers[addrs[2]] = pnPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = pnPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
}

func TestReBroadcastTxs_EN(t *testing.T) {
	// PN Peer=0, EN Peer=1
	{
		pm := &ProtocolManager{}
		pm.nodetype = common.ENDPOINTNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(common.ENDPOINTNODE).Times(3)
		enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.enpeers[addrs[2]] = enPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
	// PN Peer=1, EN Peer=0
	{
		pm := &ProtocolManager{}
		pm.nodetype = common.ENDPOINTNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(common.PROXYNODE).Times(3)
		pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.pnpeers[addrs[2]] = pnPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = pnPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
}

func TestUseTxResend(t *testing.T) {
	testSet := [...]struct {
		pm     *ProtocolManager
		result bool
	}{
		{&ProtocolManager{nodetype: common.CONSENSUSNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: common.ENDPOINTNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: common.PROXYNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: common.BOOTNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: common.UNKNOWNNODE, txResendUseLegacy: true}, false},

		{&ProtocolManager{nodetype: common.CONSENSUSNODE, txResendUseLegacy: false}, false},
		{&ProtocolManager{nodetype: common.ENDPOINTNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: common.PROXYNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: common.BOOTNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: common.UNKNOWNNODE, txResendUseLegacy: false}, true},
	}

	for _, tc := range testSet {
		assert.Equal(t, tc.result, tc.pm.useTxResend())
	}
}

func TestNodeInfo(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	pm.blockchain = mockBlockChain

	genesis := newBlock(0)
	block := newBlock(blockNum1)
	config := &params.ChainConfig{ChainID: td1}

	pm.networkId = 1234
	mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)
	mockBlockChain.EXPECT().GetTd(block.Hash(), block.NumberU64()).Return(td1).Times(1)
	mockBlockChain.EXPECT().Genesis().Return(genesis).Times(1)
	mockBlockChain.EXPECT().Config().Return(config).Times(1)

	expected := &NodeInfo{
		Network:    pm.networkId,
		BlockScore: td1,
		Genesis:    genesis.Hash(),
		Config:     config,
		Head:       block.Hash(),
	}

	assert.Equal(t, *expected, *pm.NodeInfo())
}

func TestGetCNPeersAndGetENPeers(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peers.cnpeers[addrs[0]] = cnPeer
	peers.pnpeers[addrs[1]] = pnPeer
	peers.enpeers[addrs[2]] = enPeer

	cnPeers := pm.GetCNPeers()
	enPeers := pm.GetENPeers()

	assert.Equal(t, 1, len(cnPeers))
	assert.Equal(t, 1, len(enPeers))

	assert.Equal(t, cnPeer, cnPeers[addrs[0]])
	assert.Equal(t, enPeer, enPeers[addrs[2]])
}

func TestFindPeers_AddrExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addrs[0]).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addrs[1]).Times(1)
	enPeer.EXPECT().GetAddr().Return(addrs[2]).Times(1)

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addrs[0]])
	assert.EqualValues(t, pnPeer, foundPeers[addrs[1]])
	assert.Nil(t, foundPeers[addrs[2]])
}

func TestFindPeers_AddrNotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	pnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	enPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeids[0]).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeids[1]).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeids[2]).Times(1)

	cnPeer.EXPECT().SetAddr(addrs[0]).Times(1)
	pnPeer.EXPECT().SetAddr(addrs[1]).Times(1)
	enPeer.EXPECT().SetAddr(addrs[2]).Times(1)

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addrs[0]])
	assert.EqualValues(t, pnPeer, foundPeers[addrs[1]])
	assert.Nil(t, foundPeers[addrs[2]])
}

func TestFindCNPeers(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers

	cnPeer1 := NewMockPeer(mockCtrl)
	cnPeer2 := NewMockPeer(mockCtrl)
	cnPeer3 := NewMockPeer(mockCtrl)

	peers.cnpeers[addrs[0]] = cnPeer1
	peers.cnpeers[addrs[1]] = cnPeer2
	peers.cnpeers[addrs[2]] = cnPeer3

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundCNPeers := pm.FindCNPeers(targets)

	assert.Equal(t, 2, len(foundCNPeers))
	assert.EqualValues(t, cnPeer1, foundCNPeers[addrs[0]])
	assert.EqualValues(t, cnPeer2, foundCNPeers[addrs[1]])
	assert.Nil(t, foundCNPeers[addrs[2]])
}

func TestGetPeers_AddrExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addrs[0]).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addrs[1]).Times(1)
	enPeer.EXPECT().GetAddr().Return(addrs[2]).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addrs[0]))
	assert.True(t, contains(foundAddrs, addrs[1]))
	assert.True(t, contains(foundAddrs, addrs[2]))
}

func TestGetPeers_AddrNotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	pnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	enPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeids[0]).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeids[1]).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeids[2]).Times(1)

	cnPeer.EXPECT().SetAddr(addrs[0]).Times(1)
	pnPeer.EXPECT().SetAddr(addrs[1]).Times(1)
	enPeer.EXPECT().SetAddr(addrs[2]).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addrs[0]))
	assert.True(t, contains(foundAddrs, addrs[1]))
	assert.True(t, contains(foundAddrs, addrs[2]))
}

func TestEnqueue(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fetcherMock := mocks.NewMockProtocolManagerFetcher(mockCtrl)
	pm.fetcher = fetcherMock

	block := newBlock(blockNum1)
	id := nodeids[0].String()

	fetcherMock.EXPECT().Enqueue(id, block).Times(1)
	pm.Enqueue(id, block)
}

func TestProtocolManager_Downloader(t *testing.T) {
	pm := &ProtocolManager{}
	assert.Nil(t, pm.Downloader())

	downloader := &downloader.Downloader{}
	pm.downloader = downloader

	assert.Equal(t, downloader, pm.Downloader())
}

func TestProtocolManager_SetWsEndPoint(t *testing.T) {
	pm := &ProtocolManager{}
	assert.Equal(t, "", pm.wsendpoint)

	wsep := "wsep"
	pm.SetWsEndPoint(wsep)
	assert.Equal(t, wsep, pm.wsendpoint)
}

func TestBroadcastTxsSortedByPriceAndTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 10)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	// txs[0:4] - legacyTxType with gasPrice 25ston. txs[5:10] - dynamicFeeTxType with gasPrice 25ston + tip.
	txs := types.Transactions{}
	for _, key := range keys[0:5] {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(25*params.Gkei), nil), signer, key)
		txs = append(txs, tx)
	}
	for i, key := range keys[5:10] {
		tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			AccountNonce: uint64(0),
			Recipient:    &common.Address{},
			Amount:       big.NewInt(100),
			GasLimit:     100,
			GasFeeCap:    big.NewInt(int64(25*params.Gkei + i + 1)),
			GasTipCap:    big.NewInt(int64(i + 1)),
			Payload:      nil,
		}), signer, key)
		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sortedTxs = types.SortTxsByPriceAndTime(sortedTxs, big.NewInt(25*params.Gkei))

	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE

	peers := newPeerSet()
	basePeer, _, oppositePipe := newBasePeer()

	pm.peers = peers
	pm.peers.Register(basePeer, nil)

	go func(t *testing.T) {
		pm.BroadcastTxs(txs)
	}(t)

	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.True(t, basePeer.KnowsTx(tx.Hash()))
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}

func TestReBroadcastTxsSortedByTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 10)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	// txs[0:4] - legacyTxType with gasPrice 25ston. txs[5:10] - dynamicFeeTxType with gasPrice 25ston + tip.
	txs := types.Transactions{}
	for _, key := range keys[0:4] {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(25*params.Gkei), nil), signer, key)
		txs = append(txs, tx)
	}
	for i, key := range keys[5:10] {
		tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			AccountNonce: uint64(0),
			Recipient:    &common.Address{},
			Amount:       big.NewInt(100),
			GasLimit:     100,
			GasFeeCap:    big.NewInt(int64(25*params.Gkei + i + 1)),
			GasTipCap:    big.NewInt(int64(i + 1)),
			Payload:      nil,
		}), signer, key)
		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sortedTxs = types.SortTxsByPriceAndTime(sortedTxs, big.NewInt(25*params.Gkei))

	pm := &ProtocolManager{}
	pm.nodetype = common.ENDPOINTNODE

	peers := newPeerSet()
	basePeer, _, oppositePipe := newBasePeer()

	pm.peers = peers
	pm.peers.Register(basePeer, nil)

	go func(t *testing.T) {
		pm.ReBroadcastTxs(txs)
	}(t)

	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}

func contains(addrs []common.Address, item common.Address) bool {
	for _, a := range addrs {
		if a == item {
			return true
		}
	}
	return false
}

func createAndRegisterPeers(mockCtrl *gomock.Controller, peers *peerSet) (*MockPeer, *MockPeer, *MockPeer) {
	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peers.cnpeers[addrs[0]] = cnPeer
	peers.pnpeers[addrs[1]] = pnPeer
	peers.enpeers[addrs[2]] = enPeer

	peers.peers[fmt.Sprintf("%x", nodeids[0][:8])] = cnPeer
	peers.peers[fmt.Sprintf("%x", nodeids[1][:8])] = pnPeer
	peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

	return cnPeer, pnPeer, enPeer
}

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress(testKey.PublicKey)
)

func newTestBackendWithGenerator(blocks int, generator func(int, *blockchain.BlockGen)) (database.DBManager, *blockchain.BlockChain) {
	var (
		// Create a database pre-initialize with a genesis block
		db     = database.NewMemoryDBManager()
		config = params.TestChainConfig
		engine = gxhash.NewFaker()
	)

	gspec := &blockchain.Genesis{
		Config: config,
		Alloc:  blockchain.GenesisAlloc{testAddr: {Balance: big.NewInt(100_000_000_000_000_000)}},
	}
	genesis := gspec.MustCommit(db)
	cacheConfig := &blockchain.CacheConfig{
		CacheSize:           512,
		BlockInterval:       blockchain.DefaultBlockInterval,
		TriesInMemory:       blockchain.DefaultTriesInMemory,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:   512,
		ArchiveMode:         true, // Archive mode
	}
	chain, _ := blockchain.NewBlockChain(db, cacheConfig, config, engine, vm.Config{})

	bs, _ := blockchain.GenerateChain(config, genesis, engine, db, blocks, generator)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	for i, block := range bs {
		chain.StateCache().TrieDB().Commit(block.Root(), false, uint64(i))
	}

	return db, chain
}

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders(t *testing.T) {
	t.Parallel()

	db, backend := newTestBackendWithGenerator(downloader.MaxHeaderFetch+15, nil)
	peer, _, net := newBasePeer()

	// Create a "random" unknown hash for testing
	var unknown common.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	getHashes := func(from, limit uint64) (hashes []common.Hash) {
		for i := uint64(0); i < limit; i++ {
			hashes = append(hashes, backend.GetBlockByNumber(from-1-i).Hash())
		}
		return hashes
	}

	limit := uint64(downloader.MaxHeaderFetch)
	tests := []struct {
		query  *getBlockHeadersData // The query to execute for header retrieval
		expect []common.Hash        // The hashes of the block whose headers are expected
	}{
		// A single random block should be retrievable by hash
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: backend.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{backend.GetBlockByNumber(limit / 2).Hash()},
		},
		// A single random block should be retrievable by number
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 1},
			[]common.Hash{backend.GetBlockByNumber(limit / 2).Hash()},
		},
		// Multiple headers should be retrievable in both directions
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3},
			[]common.Hash{
				backend.GetBlockByNumber(limit / 2).Hash(),
				backend.GetBlockByNumber(limit/2 + 1).Hash(),
				backend.GetBlockByNumber(limit/2 + 2).Hash(),
			},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.GetBlockByNumber(limit / 2).Hash(),
				backend.GetBlockByNumber(limit/2 - 1).Hash(),
				backend.GetBlockByNumber(limit/2 - 2).Hash(),
			},
		},
		// Multiple headers with skip lists should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.GetBlockByNumber(limit / 2).Hash(),
				backend.GetBlockByNumber(limit/2 + 4).Hash(),
				backend.GetBlockByNumber(limit/2 + 8).Hash(),
			},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.GetBlockByNumber(limit / 2).Hash(),
				backend.GetBlockByNumber(limit/2 - 4).Hash(),
				backend.GetBlockByNumber(limit/2 - 8).Hash(),
			},
		},
		// The chain endpoints should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 0}, Amount: 1},
			[]common.Hash{backend.GetBlockByNumber(0).Hash()},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64()}, Amount: 1},
			[]common.Hash{backend.CurrentBlock().Hash()},
		},
		// If the peer requests a bit into the future, we deliver what we have
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64()}, Amount: 10},
			[]common.Hash{backend.CurrentBlock().Hash()},
		},
		// Ensure protocol limits are honored
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64() - 1}, Amount: limit + 10, Reverse: true},
			getHashes(backend.CurrentBlock().Number().Uint64(), limit),
		},
		// Check that requesting more than available is handled gracefully
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64() - 4}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.GetBlockByNumber(backend.CurrentBlock().Number().Uint64() - 4).Hash(),
				backend.GetBlockByNumber(backend.CurrentBlock().Number().Uint64()).Hash(),
			},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.GetBlockByNumber(4).Hash(),
				backend.GetBlockByNumber(0).Hash(),
			},
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64() - 4}, Skip: 2, Amount: 3},
			[]common.Hash{
				backend.GetBlockByNumber(backend.CurrentBlock().Number().Uint64() - 4).Hash(),
				backend.GetBlockByNumber(backend.CurrentBlock().Number().Uint64() - 1).Hash(),
			},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.GetBlockByNumber(4).Hash(),
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]common.Hash{
				backend.GetBlockByNumber(2).Hash(),
				backend.GetBlockByNumber(1).Hash(),
				backend.GetBlockByNumber(0).Hash(),
			},
		},
		// Check a corner case where skipping causes overflow with reverse=false
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 1}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where skipping causes overflow with reverse=true
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 1}, Amount: 2, Reverse: true, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check another corner case where skipping causes overflow with reverse=false
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 1}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check another corner case where skipping causes overflow with reverse=true
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 1}, Amount: 2, Reverse: true, Skip: math.MaxUint64},
			[]common.Hash{
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: backend.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				backend.GetBlockByNumber(3).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: backend.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				backend.GetBlockByNumber(1).Hash(),
			},
		},
		// Check that non existing headers aren't returned
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: unknown}, Amount: 1},
			[]common.Hash{},
		},
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: backend.CurrentBlock().Number().Uint64() + 1}, Amount: 1},
			[]common.Hash{},
		},
	}

	pm, err := NewProtocolManager(params.TestChainConfig, downloader.FullSync, 1, nil, nil, backend.Engine(), backend, db, 1, common.ENDPOINTNODE, &Config{TxResendUseLegacy: false, TxResendInterval: 1, TxResendCount: 0})
	if err != nil {
		t.Fatal(err)
	}

	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		var headers []*types.Header
		for _, hash := range tt.expect {
			headers = append(headers, backend.GetBlockByHash(hash).Header())
		}
		// Send the hash request and verify the response
		_, r, _ := rlp.EncodeToReader(tt.query)
		go func() {
			handleBlockHeadersRequestMsg(pm, peer, p2p.Msg{Code: BlockHeadersRequestMsg, Payload: r})
		}()
		msg, err := net.ReadMsg()
		if err != nil {
			t.Fatal(err)
		}
		var receivedHeaders []*types.Header
		if err := msg.Decode(&receivedHeaders); err != nil {
			t.Fatal(err)
		}
		encodedHeaders, _ := rlp.EncodeToBytes(headers)
		encodedReceivedHeaders, _ := rlp.EncodeToBytes(receivedHeaders)
		if !bytes.Equal(encodedHeaders, encodedReceivedHeaders) {
			t.Fatalf("test %d: headers mismatch, expected: %v, received: %v", i, headers, receivedHeaders)
		}

		// If the test used number origins, repeat with hashes as the origin too
		if tt.query.Origin.Hash == (common.Hash{}) {
			if origin := backend.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
				tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0
				_, r, _ := rlp.EncodeToReader(tt.query)
				go func() {
					handleBlockHeadersRequestMsg(pm, peer, p2p.Msg{Code: BlockHeadersRequestMsg, Payload: r})
				}()
				msg, err := net.ReadMsg()
				if err != nil {
					t.Error(err)
				}
				if err := msg.Decode(&receivedHeaders); err != nil {
					t.Fatal(err)
				}
				encodedHeaders, _ := rlp.EncodeToBytes(headers)
				encodedReceivedHeaders, _ := rlp.EncodeToBytes(receivedHeaders)
				if !bytes.Equal(encodedHeaders, encodedReceivedHeaders) {
					t.Errorf("test %d: headers mismatch, expected: %v, received: %v", i, headers, receivedHeaders)
				}
			}
		}
	}
}
