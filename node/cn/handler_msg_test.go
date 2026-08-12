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
	"crypto/sha256"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/auction"
	auction_mock "github.com/kaiachain/kaia/kaiax/auction/mock"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/networks/p2p"
	mocks2 "github.com/kaiachain/kaia/node/cn/mocks"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedErr = errors.New("some error")

// generateMsg creates a message struct for message handling tests.
func generateMsg(t *testing.T, msgCode uint64, data interface{}) p2p.Msg {
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		t.Fatal(err)
	}
	return p2p.Msg{
		Code:    msgCode,
		Size:    uint32(size),
		Payload: r,
	}
}

// prepareTestHandleNewBlockMsg creates structs for TestHandleNewBlockMsg_ tests.
func prepareTestHandleNewBlockMsg(t *testing.T, mockCtrl *gomock.Controller, blockNum int) (*types.Block, p2p.Msg, *MockPeer, *mocks2.MockProtocolManagerFetcher) {
	mockPeer := NewMockPeer(mockCtrl)

	newBlock := newBlock(blockNum)
	newBlock.ReceivedFrom = mockPeer
	msg := generateMsg(t, NewBlockMsg, newBlockData{Block: newBlock, TD: big.NewInt(int64(blockNum))})

	mockPeer.EXPECT().AddToKnownBlocks(newBlock.Hash()).Times(1)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

	mockFetcher := mocks2.NewMockProtocolManagerFetcher(mockCtrl)
	mockFetcher.EXPECT().Enqueue(nodeids[0].String(), newBlock).Times(1)

	return newBlock, msg, mockPeer, mockFetcher
}

func prepareDownloader(t *testing.T) (*gomock.Controller, *mocks2.MockProtocolManagerDownloader, *MockPeer, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockDownloader := mocks2.NewMockProtocolManagerDownloader(mockCtrl)

	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

	pm := &ProtocolManager{downloader: mockDownloader}

	return mockCtrl, mockDownloader, mockPeer, pm
}

func TestHandleBlockHeadersMsg(t *testing.T) {
	headers := []*types.Header{blocks[0].Header(), blocks[1].Header()}
	{
		mockCtrl, _, mockPeer, pm := prepareDownloader(t)
		msg := generateMsg(t, BlockHeadersMsg, blocks[0].Header())

		assert.Error(t, handleBlockHeadersMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockDownloader, mockPeer, pm := prepareDownloader(t)
		msg := generateMsg(t, BlockHeadersMsg, headers)
		mockDownloader.EXPECT().DeliverHeaders(nodeids[0].String(), gomock.Eq(headers)).Return(expectedErr).Times(1)

		assert.NoError(t, handleBlockHeadersMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockDownloader, mockPeer, pm := prepareDownloader(t)
		msg := generateMsg(t, BlockHeadersMsg, headers)
		mockDownloader.EXPECT().DeliverHeaders(nodeids[0].String(), gomock.Eq(headers)).Return(nil).Times(1)

		assert.NoError(t, handleBlockHeadersMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}

func prepareBlockChain(t *testing.T) (*gomock.Controller, *mocks.MockBlockChain, *MockPeer, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockAuctionModule := auction_mock.NewMockAuctionModule(mockCtrl)

	mockAuctionModule.EXPECT().HandleBid(gomock.Any(), gomock.Any()).AnyTimes()

	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

	pm := &ProtocolManager{blockchain: mockBlockChain, auctionModule: mockAuctionModule}

	return mockCtrl, mockBlockChain, mockPeer, pm
}

func TestHandleBlockBodiesRequestMsg(t *testing.T) {
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, BlockBodiesRequestMsg, uint64(123)) // Non-list value to invoke an error

		bodies, err := handleBlockBodiesRequest(pm, mockPeer, msg)
		assert.Nil(t, bodies)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		requestedHashes := []common.Hash{hashes[0], hashes[1]}
		returnedData := []rlp.RawValue{hashes[1][:], hashes[0][:]}

		mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, BlockBodiesRequestMsg, requestedHashes)

		mockBlockChain.EXPECT().GetBodyRLP(gomock.Eq(hashes[0])).Return(returnedData[0]).Times(1)
		mockBlockChain.EXPECT().GetBodyRLP(gomock.Eq(hashes[1])).Return(returnedData[1]).Times(1)

		bodies, err := handleBlockBodiesRequest(pm, mockPeer, msg)
		assert.Equal(t, returnedData, bodies)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestHandleBlockBodiesMsg(t *testing.T) {
	{
		mockCtrl, _, mockPeer, pm := prepareDownloader(t)
		msg := generateMsg(t, BlockBodiesMsg, blocks[0].Header())

		assert.Error(t, handleBlockBodiesMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}

func TestNodeDataRequestMsg(t *testing.T) {
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, NodeDataRequestMsg, uint64(123)) // Non-list value to invoke an error

		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	{
		requestedHashes := []common.Hash{hashes[0], hashes[1]}
		returnedData := [][]byte{hashes[1][:], hashes[0][:]}

		mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, NodeDataRequestMsg, requestedHashes)

		mockBlockChain.EXPECT().TrieNode(gomock.Eq(hashes[0])).Return(returnedData[0], nil).Times(1)
		mockBlockChain.EXPECT().TrieNode(gomock.Eq(hashes[1])).Return(returnedData[1], nil).Times(1)

		mockPeer.EXPECT().SendNodeData(returnedData).Return(nil).Times(1)

		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func TestHandleReceiptsRequestMsg(t *testing.T) {
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, ReceiptsRequestMsg, uint64(123)) // Non-list value to invoke an error

		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	{
		requestedHashes := []common.Hash{hashes[0], hashes[1]}

		rct1 := newReceipt(123)

		mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, ReceiptsRequestMsg, requestedHashes)

		mockBlockChain.EXPECT().GetReceiptsByBlockHash(gomock.Eq(hashes[0])).Return(types.Receipts{rct1}).Times(1)
		mockBlockChain.EXPECT().GetReceiptsByBlockHash(gomock.Eq(hashes[1])).Return(nil).Times(1)
		mockBlockChain.EXPECT().GetHeaderByHash(gomock.Eq(hashes[1])).Return(nil).Times(1)

		mockPeer.EXPECT().SendReceiptsRLP(gomock.Any()).Return(nil).Times(1)

		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func TestHandleNewBlockMsg_LargeLocalPeerBlockScore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	_, msg, mockPeer, mockFetcher := prepareTestHandleNewBlockMsg(t, mockCtrl, blockNum1)

	pm := &ProtocolManager{}
	pm.fetcher = mockFetcher

	mockPeer.EXPECT().Head().Return(hash1, big.NewInt(blockNum1+1)).AnyTimes()

	assert.NoError(t, handleNewBlockMsg(pm, mockPeer, msg))
}

func TestHandleNewBlockMsg_SmallLocalPeerBlockScore_NoSynchronise(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	block, msg, mockPeer, mockFetcher := prepareTestHandleNewBlockMsg(t, mockCtrl, blockNum1)

	pm := &ProtocolManager{}
	pm.fetcher = mockFetcher

	mockPeer.EXPECT().Head().Return(hash1, big.NewInt(blockNum1-2)).AnyTimes()
	mockPeer.EXPECT().SetHead(block.ParentHash(), big.NewInt(blockNum1-1)).Times(1)

	currBlock := newBlock(blockNum1 - 1)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().CurrentBlock().Return(currBlock).Times(1)
	mockBlockChain.EXPECT().GetTd(currBlock.Hash(), currBlock.NumberU64()).Return(big.NewInt(blockNum1)).Times(1)

	pm.blockchain = mockBlockChain

	assert.NoError(t, handleNewBlockMsg(pm, mockPeer, msg))
}

func TestHandleTxMsg(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()

	txs := types.Transactions{tx1}
	msg := generateMsg(t, TxMsg, txs)

	// If pm.acceptTxs == 0, nothing happens.
	{
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
	}
	// If pm.acceptTxs == 1, TxPool.HandleTxMsg is called.
	{
		pm.acceptTxs.Store(1)
		mockTxPool := mocks.NewMockTxPool(mockCtrl)

		// The time field in received transaction through pm.handleMsg() has different value from generated transaction(`tx1`).
		// It can check whether the transaction created `HandleTxMsg()` is the same as `tx1` through `AddToKnownTxs(txs[0].Hash())`.
		mockTxPool.EXPECT().HandleTxMsg(gomock.Any()).AnyTimes()
		pm.txpool = mockTxPool

		mockPeer.EXPECT().AddToKnownTxs(txs[0].Hash()).Times(1)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
	}
}

// prepareBlobTxMsg returns a protocol manager, a peer and a signed blob transaction
// with a valid v1 sidecar. Callers corrupt the sidecar via blobTx.BlobTxSidecar().
func prepareBlobTxMsg(t *testing.T, mockCtrl *gomock.Controller) (*ProtocolManager, *MockPeer, *types.Transaction) {
	pm := &ProtocolManager{verifiedBlobTxs: newKnownHashSet(maxVerifiedBlobTxs)}
	pm.acceptTxs.Store(1)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().HandleTxMsg(gomock.Any()).AnyTimes()
	pm.txpool = mockTxPool

	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
	mockPeer.EXPECT().GetID().Return("test-peer").AnyTimes()
	mockPeer.EXPECT().AddToKnownTxs(gomock.Any()).AnyTimes()

	sidecar, hashes := newBlobSidecar(t)
	return pm, mockPeer, newBlobTx(t, 0, hashes, sidecar)
}

// newBlobSidecar returns a valid v1 sidecar and the blob hashes it commits to.
func newBlobSidecar(t *testing.T) (*types.BlobTxSidecar, []common.Hash) {
	blob := kzg4844.Blob{}
	commitment, err := kzg4844.BlobToCommitment(&blob)
	require.NoError(t, err)
	proofs, err := kzg4844.ComputeCellProofs(&blob)
	require.NoError(t, err)
	return &types.BlobTxSidecar{
		Version:     types.BlobSidecarVersion1,
		Blobs:       []kzg4844.Blob{blob},
		Commitments: []kzg4844.Commitment{commitment},
		Proofs:      proofs,
	}, []common.Hash{common.Hash(kzg4844.CalcBlobHashV1(sha256.New(), &commitment))}
}

// newBlobTx signs a blob transaction carrying the given hashes and sidecar. The nonce is
// a parameter so a caller can replay one sidecar under a different transaction hash.
func newBlobTx(t *testing.T, nonce uint64, hashes []common.Hash, sidecar *types.BlobTxSidecar) *types.Transaction {
	tx, err := types.NewTransactionWithMap(types.TxTypeEthereumBlob, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      nonce,
		types.TxValueKeyTo:         crypto.PubkeyToAddress(keys[0].PublicKey),
		types.TxValueKeyAmount:     big.NewInt(0),
		types.TxValueKeyGasLimit:   uint64(10000000),
		types.TxValueKeyGasFeeCap:  big.NewInt(25),
		types.TxValueKeyGasTipCap:  big.NewInt(25),
		types.TxValueKeyData:       []byte{},
		types.TxValueKeyAccessList: types.AccessList{},
		types.TxValueKeyBlobFeeCap: big.NewInt(25),
		types.TxValueKeyBlobHashes: hashes,
		types.TxValueKeySidecar:    sidecar,
		types.TxValueKeyChainID:    params.TestChainConfig.ChainID,
	})
	require.NoError(t, err)
	require.NoError(t, tx.Sign(types.MakeSigner(params.TestChainConfig, common.Big0), keys[0]))
	return tx
}

func TestHandleTxMsg_KZGVerificationError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pm, mockPeer, blobTx := prepareBlobTxMsg(t, mockCtrl)
	blobTx.BlobTxSidecar().Commitments[0][0] ^= 0xFF

	// Should return error and disconnect peer
	err := handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{blobTx}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), errKZGVerificationError.Error())
}

func TestHandleTxMsg_BlobSidecarVerifiedOnce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pm, mockPeer, blobTx := prepareBlobTxMsg(t, mockCtrl)
	require.NoError(t, handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{blobTx})))
	require.Equal(t, 1, pm.verifiedBlobTxs.Len())

	// The same sidecar under a different transaction hash must reuse the entry, or a
	// sender replays one sidecar for free by bumping the nonce.
	replay := newBlobTx(t, 1, blobTx.BlobHashes(), blobTx.BlobTxSidecar())
	require.NotEqual(t, blobTx.Hash(), replay.Hash())
	require.NoError(t, handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{replay})))
	assert.Equal(t, 1, pm.verifiedBlobTxs.Len(), "the replay was verified again")

	// A sidecar swapped under an already verified transaction hash must not reuse the
	// entry, whichever input of the verification was altered.
	for name, tamper := range map[string]func(*types.BlobTxSidecar){
		"blob":       func(sc *types.BlobTxSidecar) { sc.Blobs[0][0] ^= 0xFF },
		"proof":      func(sc *types.BlobTxSidecar) { sc.Proofs[0][0] ^= 0xFF },
		"version":    func(sc *types.BlobTxSidecar) { sc.Version = 0 },
		"commitment": func(sc *types.BlobTxSidecar) { sc.Commitments[0][0] ^= 0xFF },
	} {
		t.Run(name, func(t *testing.T) {
			pm, mockPeer, blobTx := prepareBlobTxMsg(t, mockCtrl)
			require.NoError(t, handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{blobTx})))

			tamper(blobTx.BlobTxSidecar())
			err := handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{blobTx}))
			require.Error(t, err)
			assert.Contains(t, err.Error(), errKZGVerificationError.Error())
		})
	}
}

func TestHandleTxMsg_BloblessBlobTx(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pm, mockPeer, _ := prepareBlobTxMsg(t, mockCtrl)

	// Every check in ValidateWithBlobHashes compares against the declared hashes, so
	// it passes vacuously when there are none. The pool rejects such a transaction
	// unconditionally, so the handler must not forward it.
	sidecar, _ := newBlobSidecar(t)
	sidecar.Blobs, sidecar.Commitments, sidecar.Proofs = nil, nil, nil
	blobless := newBlobTx(t, 0, nil, sidecar)
	require.NoError(t, sidecar.ValidateWithBlobHashes(nil), "expected the vacuous pass this guards")

	err := handleTxMsg(pm, mockPeer, generateMsg(t, TxMsg, types.Transactions{blobless}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), errBloblessBlobTx.Error())
}

func prepareTestHandleBlockHeaderFetchRequestMsg(t *testing.T) (*gomock.Controller, *MockPeer, *mocks.MockBlockChain, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()

	return mockCtrl, mockPeer, mockBlockChain, &ProtocolManager{blockchain: mockBlockChain}
}

func TestHandleBlockHeaderFetchRequestMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl, mockPeer, _, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, newBlock(blockNum1)) // use message data as a block, not a hash

		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// GetHeaderByHash returns nil, an error is returned.
	{
		mockCtrl, mockPeer, mockBlockChain, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)
		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(nil).AnyTimes()
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, hash1)

		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// GetHeaderByHash returns a header, p.SendFetchedBlockHeader(header) should be called.
	{
		mockCtrl, mockPeer, mockBlockChain, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)

		header := newBlock(blockNum1).Header()

		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(header).AnyTimes()
		mockPeer.EXPECT().SendFetchedBlockHeader(header).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, hash1)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func prepareTestHandleBlockHeaderFetchResponseMsg(t *testing.T) (*gomock.Controller, *MockPeer, *mocks2.MockProtocolManagerFetcher, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()

	mockFetcher := mocks2.NewMockProtocolManagerFetcher(mockCtrl)
	pm := &ProtocolManager{fetcher: mockFetcher}

	return mockCtrl, mockPeer, mockFetcher, pm
}

func TestHandleBlockHeaderFetchResponseMsg(t *testing.T) {
	header := newBlock(blockNum1).Header()
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		pm := &ProtocolManager{}
		msg := generateMsg(t, BlockHeaderFetchResponseMsg, newBlock(blockNum1)) // use message data as a block, not a header
		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// FilterHeaders returns nil, error is not returned.
	{
		mockCtrl, mockPeer, mockFetcher, pm := prepareTestHandleBlockHeaderFetchResponseMsg(t)
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()
		mockFetcher.EXPECT().FilterHeaders(nodeids[0].String(), gomock.Eq([]*types.Header{header}), gomock.Any()).Return(nil).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchResponseMsg, header)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// FilterHeaders returns not-nil, peer.GetID() is called twice to leave a log.
	{
		mockCtrl, mockPeer, mockFetcher, pm := prepareTestHandleBlockHeaderFetchResponseMsg(t)
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()
		mockFetcher.EXPECT().FilterHeaders(nodeids[0].String(), gomock.Eq([]*types.Header{header}), gomock.Any()).Return([]*types.Header{header}).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchResponseMsg, header)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func preparePeerAndDownloader(t *testing.T) (*gomock.Controller, *MockPeer, *mocks2.MockProtocolManagerDownloader, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()
	mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()

	mockDownloader := mocks2.NewMockProtocolManagerDownloader(mockCtrl)
	pm := &ProtocolManager{downloader: mockDownloader}

	return mockCtrl, mockPeer, mockDownloader, pm
}

func TestHandleReceiptMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()

		pm := &ProtocolManager{}
		msg := generateMsg(t, ReceiptsMsg, newBlock(blockNum1)) // use message data as a block, not a header
		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// DeliverReceipts returns nil, error is not returned.
	{
		receipts := make([][]*types.Receipt, 1)
		receipts[0] = []*types.Receipt{newReceipt(123)}

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverReceipts(nodeids[0].String(), gomock.Eq(receipts)).Times(1).Return(nil)

		msg := generateMsg(t, ReceiptsMsg, receipts)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// DeliverReceipts returns an error, but the error is not returned.
	{
		receipts := make([][]*types.Receipt, 1)
		receipts[0] = []*types.Receipt{newReceipt(123)}

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverReceipts(nodeids[0].String(), gomock.Eq(receipts)).Times(1).Return(expectedErr)

		msg := generateMsg(t, ReceiptsMsg, receipts)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func TestHandleNodeDataMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().GetVersion().Return(kaia63).AnyTimes()
		pm := &ProtocolManager{}
		msg := generateMsg(t, NodeDataMsg, newBlock(blockNum1)) // use message data as a block, not a node data
		assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// DeliverNodeData returns nil, error is not returned.
	{
		nodeData := make([][]byte, 1)
		nodeData[0] = hash1[:]

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverNodeData(nodeids[0].String(), gomock.Eq(nodeData)).Times(1).Return(nil)

		msg := generateMsg(t, NodeDataMsg, nodeData)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
	// DeliverNodeData returns an error, but the error is not returned.
	{
		nodeData := make([][]byte, 1)
		nodeData[0] = hash1[:]

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverNodeData(nodeids[0].String(), gomock.Eq(nodeData)).Times(1).Return(expectedErr)

		msg := generateMsg(t, NodeDataMsg, nodeData)
		assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg))
		mockCtrl.Finish()
	}
}

func TestHandleStakingInfoRequestMsg(t *testing.T) {
	testChainConfig := params.TestChainConfig.Copy()

	{
		// test if chain config istanbul is nil
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = nil
		pm.chainconfig = testChainConfig

		err := handleStakingInfoRequestMsg(pm, mockPeer, p2p.Msg{})
		assert.Error(t, err)
		assert.Equal(t, err, errResp(ErrUnsupportedEnginePolicy, "the engine is not istanbul or the policy is not weighted random"))
		mockCtrl.Finish()
	}
	{
		// test if chain config istanbul is not nil, but proposer policy is not weighted random
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = params.GetDefaultIstanbulConfig()
		testChainConfig.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
		pm.chainconfig = testChainConfig

		err := handleStakingInfoRequestMsg(pm, mockPeer, p2p.Msg{})
		assert.Error(t, err)
		assert.Equal(t, err, errResp(ErrUnsupportedEnginePolicy, "the engine is not istanbul or the policy is not weighted random"))
		mockCtrl.Finish()
	}
	{
		// test if message does not contain expected data
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = params.GetDefaultIstanbulConfig()
		testChainConfig.Istanbul.ProposerPolicy = uint64(istanbul.WeightedRandom)
		pm.chainconfig = testChainConfig
		msg := generateMsg(t, StakingInfoRequestMsg, uint64(123)) // Non-list value to invoke an error

		err := handleStakingInfoRequestMsg(pm, mockPeer, msg)
		assert.Error(t, err)
		assert.Equal(t, err, rlp.ErrExpectedList)
		mockCtrl.Finish()
	}

	// Setup governance items for testing
	{
		requestedHashes := []common.Hash{hashes[0], hashes[1]}

		mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)

		mStaking := staking_mock.NewMockStakingModule(mockCtrl)
		si := &staking.StakingInfo{
			SourceBlockNum:   4,
			NodeIds:          []common.Address{{0x1}, {0x1}},
			StakingContracts: []common.Address{{0x2}, {0x2}},
			RewardAddrs:      []common.Address{{0x3}, {0x3}},
			StakingAmounts:   []uint64{2, 5, 6},
		}
		mStaking.EXPECT().GetStakingInfoFromDB(gomock.Eq(uint64(4))).Return(si).Times(1)
		mStaking.EXPECT().GetStakingInfoFromDB(gomock.Eq(uint64(5))).Return(nil).Times(1)
		pm.stakingModule = mStaking

		testChainConfig.KaiaCompatibleBlock = nil
		testChainConfig.Istanbul = &params.IstanbulConfig{ProposerPolicy: uint64(istanbul.WeightedRandom)}
		testChainConfig.Governance = params.GetDefaultGovernanceConfig()
		testChainConfig.Governance.Reward.StakingUpdateInterval = 4
		pm.chainconfig = testChainConfig

		msg := generateMsg(t, StakingInfoRequestMsg, requestedHashes)

		mockBlockChain.EXPECT().GetHeaderByHash(gomock.Eq(hashes[0])).Return(&types.Header{Number: big.NewInt(4)}).Times(1) // on staking interval
		mockBlockChain.EXPECT().GetHeaderByHash(gomock.Eq(hashes[1])).Return(&types.Header{Number: big.NewInt(5)}).Times(1) // not on staking interval

		useGini, minStake := testChainConfig.Governance.Reward.UseGiniCoeff, testChainConfig.Governance.Reward.MinimumStake.Uint64()
		expectedResult := staking.FromStakingInfoWithGini(si, useGini, minStake)
		data, _ := rlp.EncodeToBytes(expectedResult)
		expectedRlpList := []rlp.RawValue{data}
		mockPeer.EXPECT().SendStakingInfoRLP(gomock.Eq(expectedRlpList)).Return(nil).Times(1)

		err := handleStakingInfoRequestMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestHandleStakingInfoRequestMsgAfterKaia(t *testing.T) {
	testChainConfig := params.TestChainConfig.Copy()

	{
		requestedHashes := []common.Hash{hashes[0], hashes[1]}

		mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)

		mStaking := staking_mock.NewMockStakingModule(mockCtrl)
		siBeforeKaia := &staking.StakingInfo{
			SourceBlockNum:   4,
			NodeIds:          []common.Address{{0x1}, {0x1}},
			StakingContracts: []common.Address{{0x2}, {0x2}},
			RewardAddrs:      []common.Address{{0x3}, {0x3}},
			StakingAmounts:   []uint64{2, 5, 6},
		}
		siAfterKaia := &staking.StakingInfo{
			SourceBlockNum:   5,
			NodeIds:          []common.Address{{0x1}, {0x1}},
			StakingContracts: []common.Address{{0x2}, {0x2}},
			RewardAddrs:      []common.Address{{0x3}, {0x3}},
			StakingAmounts:   []uint64{2, 5, 6},
		}
		mStaking.EXPECT().GetStakingInfoFromDB(gomock.Eq(uint64(4))).Return(siBeforeKaia).Times(1)
		mStaking.EXPECT().GetStakingInfo(gomock.Eq(uint64(6))).Return(siAfterKaia, nil).Times(1)
		pm.stakingModule = mStaking

		testChainConfig.KaiaCompatibleBlock = big.NewInt(5)
		testChainConfig.Istanbul = &params.IstanbulConfig{ProposerPolicy: uint64(istanbul.WeightedRandom)}
		testChainConfig.Governance = params.GetDefaultGovernanceConfig()
		testChainConfig.Governance.Reward.StakingUpdateInterval = 4
		pm.chainconfig = testChainConfig

		msg := generateMsg(t, StakingInfoRequestMsg, requestedHashes)

		mockBlockChain.EXPECT().Config().Return(pm.chainconfig).AnyTimes()
		mockBlockChain.EXPECT().GetHeaderByHash(gomock.Eq(hashes[0])).Return(&types.Header{Number: big.NewInt(4)}).Times(1) // should return StakingInfo(4)
		mockBlockChain.EXPECT().GetHeaderByHash(gomock.Eq(hashes[1])).Return(&types.Header{Number: big.NewInt(6)}).Times(1) // should return StakingInfo(5)

		useGini, minStake := testChainConfig.Governance.Reward.UseGiniCoeff, testChainConfig.Governance.Reward.MinimumStake.Uint64()
		dataBeforeKaia, _ := rlp.EncodeToBytes(staking.FromStakingInfoWithGini(siBeforeKaia, useGini, minStake))
		dataAfterKaia, _ := rlp.EncodeToBytes(staking.FromStakingInfoWithGini(siAfterKaia, useGini, minStake))
		expectedRlpList := []rlp.RawValue{dataBeforeKaia, dataAfterKaia}
		mockPeer.EXPECT().SendStakingInfoRLP(gomock.Eq(expectedRlpList)).Return(nil).Times(1)

		err := handleStakingInfoRequestMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestHandleStakingInfoMsg(t *testing.T) {
	testChainConfig := params.TestChainConfig.Copy()
	{
		// test if chain config istanbul is nil
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = nil
		pm.chainconfig = testChainConfig

		err := handleStakingInfoMsg(pm, mockPeer, p2p.Msg{})
		assert.Error(t, err)
		assert.Equal(t, err, errResp(ErrUnsupportedEnginePolicy, "the engine is not istanbul or the policy is not weighted random"))
		mockCtrl.Finish()
	}
	{
		// test if chain config istanbul is not nil, but proposer policy is not weighted random
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = params.GetDefaultIstanbulConfig()
		testChainConfig.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
		pm.chainconfig = testChainConfig

		err := handleStakingInfoMsg(pm, mockPeer, p2p.Msg{})
		assert.Error(t, err)
		assert.Equal(t, err, errResp(ErrUnsupportedEnginePolicy, "the engine is not istanbul or the policy is not weighted random"))
		mockCtrl.Finish()
	}
	{
		// test if message does not contain expected data
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		testChainConfig.Istanbul = params.GetDefaultIstanbulConfig()
		testChainConfig.Istanbul.ProposerPolicy = uint64(istanbul.WeightedRandom)
		pm.chainconfig = testChainConfig
		msg := generateMsg(t, StakingInfoRequestMsg, uint64(123)) // Non-list value to invoke an error

		err := handleStakingInfoMsg(pm, mockPeer, msg)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), errCode(ErrDecode).String()))
		mockCtrl.Finish()
	}

	{
		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)

		testChainConfig.Istanbul = params.GetDefaultIstanbulConfig()
		testChainConfig.Istanbul.ProposerPolicy = uint64(istanbul.WeightedRandom)
		pm.chainconfig = testChainConfig

		si := &staking.StakingInfo{
			SourceBlockNum:   4,
			NodeIds:          []common.Address{{0x1}, {0x1}},
			StakingContracts: []common.Address{{0x2}, {0x2}},
			RewardAddrs:      []common.Address{{0x3}, {0x3}},
			StakingAmounts:   []uint64{2, 5, 6},
		}
		stakingInfos := []*staking.P2PStakingInfo{
			staking.FromStakingInfoWithGini(si, false, 5000000),
		}
		mockDownloader.EXPECT().DeliverStakingInfos(gomock.Eq(nodeids[0].String()), gomock.Eq(stakingInfos)).Times(1).Return(expectedErr)

		msg := generateMsg(t, StakingInfoMsg, stakingInfos)
		err := handleStakingInfoMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestHandleBidMsg(t *testing.T) {
	mockCtrl, _, mockPeer, pm := prepareBlockChain(t)

	bidData := auction.BidData{
		TargetTxHash:  common.HexToHash("0xf3c03c891206b24f5d2ff65b460df9b58c652279a3e0faed865dde4c46fe9dab"),
		BlockNumber:   11,
		Sender:        common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
		To:            common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"),
		Nonce:         0,
		Bid:           new(big.Int).SetBytes(common.Hex2Bytes("8ac7230489e80000")),
		CallGasLimit:  10000000,
		Data:          common.Hex2Bytes("d09de08a"),
		SearcherSig:   common.Hex2Bytes("2162312ceb6a69efdb73c98ee96e56d0aea1ea019184c372022ab378151112c0747066e9a9d224a822dbf31d59de492502d69d7cfc789464fa84aaac0d53f6a11b"),
		AuctioneerSig: common.Hex2Bytes("63ca36c4f6a3522b59070539453ff92011463940f98930b34a80b06a5b6b45fa136f8e79957e56e41de19cb340f2f1f7db31f964e5d5f26b1d8df13aeb2b390c1b"),
	}

	testBid := &auction.Bid{
		BidData: bidData,
	}

	msg := generateMsg(t, BidMsg, testBid)

	mockPeer.EXPECT().GetVersion().Return(kaia63).Times(11)
	assert.Error(t, pm.handleMsg(mockPeer, addrs[0], msg), "should return error when protocol version is not kaia66")

	mockPeer.EXPECT().GetVersion().Return(kaia66).AnyTimes()
	assert.NoError(t, pm.handleMsg(mockPeer, addrs[0], msg), "should not return error when protocol version is kaia66")

	mockCtrl.Finish()
}

func TestHandleBlobSidecarsRequestMsg(t *testing.T) {
	// test if the blob sidecars are retrieved from the block chain
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		requestData := []blobSidecarsRequestData{
			{BlockNum: 10, TxIndex: 0, Hash: tx1.Hash()},
			{BlockNum: 20, TxIndex: 1, Hash: tx2.Hash()},
		}
		msg := generateMsg(t, BlobSidecarsRequestMsg, requestData)

		sidecar0 := generateTestSidecar(tx1.Hash())
		sidecar1 := generateTestSidecar(tx2.Hash())
		rlp0, _ := rlp.EncodeToBytes(blobSidecarsData{
			BlockNum: requestData[0].BlockNum,
			TxIndex:  requestData[0].TxIndex,
			TxHash:   requestData[0].Hash,
			Sidecar:  sidecar0,
		})
		rlp1, _ := rlp.EncodeToBytes(blobSidecarsData{
			BlockNum: requestData[1].BlockNum,
			TxIndex:  requestData[1].TxIndex,
			TxHash:   requestData[1].Hash,
			Sidecar:  sidecar1,
		})
		expectedRlpList := []rlp.RawValue{rlp0, rlp1}

		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Eq(big.NewInt(int64(requestData[0].BlockNum))), gomock.Eq(int(requestData[0].TxIndex))).Return(sidecar0, nil).Times(1)
		mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Eq(big.NewInt(int64(requestData[1].BlockNum))), gomock.Eq(int(requestData[1].TxIndex))).Return(sidecar1, nil).Times(1)
		pm.txpool = mockTxPool
		mockPeer.EXPECT().SendBlobSidecarsRLP(gomock.Eq(expectedRlpList)).Return(nil).Times(1)
		err := handleBlobSidecarsRequestMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
	// test if the blob sidecars are retrieved from the tx pool
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		requestData := []blobSidecarsRequestData{
			{BlockNum: 10, TxIndex: 0, Hash: tx1.Hash()},
			{BlockNum: 20, TxIndex: 1, Hash: tx2.Hash()},
		}
		msg := generateMsg(t, BlobSidecarsRequestMsg, requestData)

		sidecar0 := generateTestSidecar(tx1.Hash())
		sidecar1 := generateTestSidecar(tx2.Hash())
		rlp0, _ := rlp.EncodeToBytes(blobSidecarsData{
			BlockNum: requestData[0].BlockNum,
			TxIndex:  requestData[0].TxIndex,
			TxHash:   requestData[0].Hash,
			Sidecar:  sidecar0,
		})
		rlp1, _ := rlp.EncodeToBytes(blobSidecarsData{
			BlockNum: requestData[1].BlockNum,
			TxIndex:  requestData[1].TxIndex,
			TxHash:   requestData[1].Hash,
			Sidecar:  sidecar1,
		})
		expectedRlpList := []rlp.RawValue{rlp0, rlp1}

		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
		mockTxPool.EXPECT().GetBlobSidecarFromPool(gomock.Eq(requestData[0].Hash)).Return(sidecar0, nil).Times(1)
		mockTxPool.EXPECT().GetBlobSidecarFromPool(gomock.Eq(requestData[1].Hash)).Return(sidecar1, nil).Times(1)
		pm.txpool = mockTxPool
		mockPeer.EXPECT().SendBlobSidecarsRLP(gomock.Eq(expectedRlpList)).Return(nil).Times(1)
		err := handleBlobSidecarsRequestMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
	// test if the blob sidecars are not retrieved from the block chain or tx pool
	{
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		requestData := []blobSidecarsRequestData{
			{BlockNum: 10, TxIndex: 0, Hash: tx1.Hash()},
			{BlockNum: 20, TxIndex: 1, Hash: tx2.Hash()},
		}
		msg := generateMsg(t, BlobSidecarsRequestMsg, requestData)

		sidecar1 := generateTestSidecar(requestData[1].Hash)
		rlp1, _ := rlp.EncodeToBytes(blobSidecarsData{
			BlockNum: requestData[1].BlockNum,
			TxIndex:  requestData[1].TxIndex,
			TxHash:   requestData[1].Hash,
			Sidecar:  sidecar1,
		})
		expectedRlpList := []rlp.RawValue{rlp1}

		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Eq(big.NewInt(int64(requestData[0].BlockNum))), gomock.Eq(int(requestData[0].TxIndex))).Return(nil, nil).Times(1)
		mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Eq(big.NewInt(int64(requestData[1].BlockNum))), gomock.Eq(int(requestData[1].TxIndex))).Return(sidecar1, nil).Times(1)
		mockTxPool.EXPECT().GetBlobSidecarFromPool(gomock.Eq(requestData[0].Hash)).Return(nil, nil).Times(1)
		pm.txpool = mockTxPool
		mockPeer.EXPECT().SendBlobSidecarsRLP(gomock.Eq(expectedRlpList)).Return(nil).Times(1)
		err := handleBlobSidecarsRequestMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestHandleBlobSidecarsMsg(t *testing.T) {
	{
		// test if message does not contain expected data (decode error)
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, BlobSidecarsMsg, uint64(123)) // Non-list value to invoke an error

		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "blob sidecar request manager is not initialized"))
		mockCtrl.Finish()
	}

	{
		// test if message does not contain expected data (decode error)
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		msg := generateMsg(t, BlobSidecarsMsg, uint64(123)) // Non-list value to invoke an error
		pm.blobSidecarReqManager = &sidecarReqManager{
			list:     map[common.Hash]*sidecarReq{},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}

		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), errCode(ErrDecode).String()))
		mockCtrl.Finish()
	}

	{
		// test if message does not contain sidecar.Sidecar (decode error)
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		d := &blobSidecarsData{
			BlockNum: 1,
			TxIndex:  1,
			TxHash:   tx1.Hash(),
			Sidecar:  nil,
		}
		pm.blobSidecarReqManager = &sidecarReqManager{
			list:     map[common.Hash]*sidecarReq{},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{d})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.Error(t, err)
		mockCtrl.Finish()
	}

	{
		// empty sidecars list, nothing happens
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		pm.blobSidecarReqManager = &sidecarReqManager{
			list:     map[common.Hash]*sidecarReq{},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}

	{
		// blobSidecarReqManager does not contain hash
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		d := &blobSidecarsData{
			BlockNum: 1,
			TxIndex:  1,
			TxHash:   tx1.Hash(),
			Sidecar:  generateTestSidecar(tx1.Hash()),
		}
		pm.blobSidecarReqManager = &sidecarReqManager{
			list:     map[common.Hash]*sidecarReq{},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{d})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}

	{
		// blobSidecarReqManager entry exists but from wrong peer
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		d := &blobSidecarsData{
			BlockNum: 1,
			TxIndex:  1,
			TxHash:   tx1.Hash(),
			Sidecar:  generateTestSidecar(tx1.Hash()),
		}
		pm.blobSidecarReqManager = &sidecarReqManager{
			list: map[common.Hash]*sidecarReq{
				d.TxHash: {
					peer: "different-peer-id",
					try:  1,
					time: time.Now(),
				},
			},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{d})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}

	{
		// test: sidecar is present, peerId matches, save succeeds
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		d := &blobSidecarsData{
			BlockNum: 1,
			TxIndex:  1,
			TxHash:   tx1.Hash(),
			Sidecar:  generateTestSidecar(tx1.Hash()),
		}
		pm.blobSidecarReqManager = &sidecarReqManager{
			list: map[common.Hash]*sidecarReq{
				d.TxHash: {
					peer: mockPeer.GetID(),
					try:  1,
					time: time.Now(),
				},
			},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		// Setup expected save call
		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().SaveBlobSidecar(
			d.TxHash,
			d.Sidecar,
		).Return(nil).Times(1)
		pm.txpool = mockTxPool
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{d})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		req := pm.blobSidecarReqManager.get(d.TxHash)
		assert.Nilf(t, req, "should be deleted from blobSidecarReqManager")
		mockCtrl.Finish()
	}

	{
		// test: save fails (error is logged, not returned), req entry should not be deleted
		mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
		d := &blobSidecarsData{
			BlockNum: 1,
			TxIndex:  1,
			TxHash:   tx1.Hash(),
			Sidecar:  generateTestSidecar(tx1.Hash()),
		}
		pm.blobSidecarReqManager = &sidecarReqManager{
			list: map[common.Hash]*sidecarReq{
				d.TxHash: {
					peer: mockPeer.GetID(),
					try:  1,
					time: time.Now(),
				},
			},
			cooldown: 10 * time.Second,
			maxTry:   5,
		}
		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().SaveBlobSidecar(
			d.TxHash,
			d.Sidecar,
		).Return(errors.New("expected error")).Times(1)
		pm.txpool = mockTxPool
		msg := generateMsg(t, BlobSidecarsMsg, []*blobSidecarsData{d})
		err := handleBlobSidecarsMsg(pm, mockPeer, msg)
		assert.NoError(t, err)
		req := pm.blobSidecarReqManager.get(d.TxHash)
		assert.NotNil(t, req, "should not be deleted from blobSidecarReqManager even on error")
		mockCtrl.Finish()
	}
}

// makePhantomHashes returns n distinct hashes whose i-th byte is i+1, suitable
// for phantom-flood tests where every lookup is expected to miss.
func makePhantomHashes(n int) []common.Hash {
	out := make([]common.Hash, n)
	for i := range out {
		out[i][0] = byte(i & 0xff)
		out[i][1] = byte((i >> 8) & 0xff)
		out[i][2] = byte((i >> 16) & 0xff)
	}
	return out
}

// The next set of tests verifies the per-iteration lookup counter that bounds
// DB work when a peer floods the handler with hashes that all miss on disk.
// Each test sends 4*cap phantom hashes and asserts the handler runs at most
// 2*cap iterations before returning, matching the pattern geth adopted in
// eth/66.

func TestHandleBlockBodiesRequestMsg_PhantomFlood(t *testing.T) {
	cap := downloader.MaxBlockFetch
	hashes := makePhantomHashes(4 * cap)

	mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
	defer mockCtrl.Finish()

	msg := generateMsg(t, BlockBodiesRequestMsg, hashes)

	// With the fix, lookups (incremented once per iteration) stops at 2*cap,
	// so GetBodyRLP is invoked exactly 2*cap times for all-miss input.
	mockBlockChain.EXPECT().GetBodyRLP(gomock.Any()).Return(nil).Times(2 * cap)

	bodies, err := handleBlockBodiesRequest(pm, mockPeer, msg)
	assert.NoError(t, err)
	assert.Empty(t, bodies)
}

func TestHandleNodeDataRequestMsg_PhantomFlood(t *testing.T) {
	cap := downloader.MaxStateFetch
	hashes := makePhantomHashes(4 * cap)

	mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
	defer mockCtrl.Finish()

	msg := generateMsg(t, NodeDataRequestMsg, hashes)

	// lookups increments once per iteration, capped at 2*cap. Every miss falls
	// through TrieNode to ContractCodeWithPrefix, so both fire 2*cap times.
	mockBlockChain.EXPECT().TrieNode(gomock.Any()).Return(nil, nil).Times(2 * cap)
	mockBlockChain.EXPECT().ContractCodeWithPrefix(gomock.Any()).Return(nil, nil).Times(2 * cap)
	mockPeer.EXPECT().SendNodeData([][]byte(nil)).Return(nil).Times(1)

	assert.NoError(t, handleNodeDataRequestMsg(pm, mockPeer, msg))
}

func TestHandleReceiptsRequestMsg_PhantomFlood(t *testing.T) {
	cap := downloader.MaxReceiptFetch
	hashes := makePhantomHashes(4 * cap)

	mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
	defer mockCtrl.Finish()

	msg := generateMsg(t, ReceiptsRequestMsg, hashes)

	// lookups increments once per iteration, capped at 2*cap. Every miss falls
	// through GetReceiptsByBlockHash to GetHeaderByHash, so both fire 2*cap times.
	mockBlockChain.EXPECT().GetReceiptsByBlockHash(gomock.Any()).Return(nil).Times(2 * cap)
	mockBlockChain.EXPECT().GetHeaderByHash(gomock.Any()).Return(nil).Times(2 * cap)
	mockPeer.EXPECT().SendReceiptsRLP(gomock.Any()).Return(nil).Times(1)

	assert.NoError(t, handleReceiptsRequestMsg(pm, mockPeer, msg))
}

func TestHandleStakingInfoRequestMsg_PhantomFlood(t *testing.T) {
	cap := downloader.MaxStakingInfoFetch
	hashes := makePhantomHashes(4 * cap)

	mockCtrl, mockBlockChain, mockPeer, pm := prepareBlockChain(t)
	defer mockCtrl.Finish()

	cfg := params.TestChainConfig.Copy()
	cfg.Istanbul = &params.IstanbulConfig{ProposerPolicy: uint64(istanbul.WeightedRandom)}
	cfg.Governance = params.GetDefaultGovernanceConfig()
	pm.chainconfig = cfg

	msg := generateMsg(t, StakingInfoRequestMsg, hashes)

	// lookups increments once per iteration, capped at 2*cap. Phantom header
	// triggers `continue`, so only GetHeaderByHash fires (2*cap times).
	mockBlockChain.EXPECT().GetHeaderByHash(gomock.Any()).Return(nil).Times(2 * cap)
	mockPeer.EXPECT().SendStakingInfoRLP(gomock.Any()).Return(nil).Times(1)

	assert.NoError(t, handleStakingInfoRequestMsg(pm, mockPeer, msg))
}

func TestHandleBlobSidecarsRequestMsg_PhantomFlood(t *testing.T) {
	cap := downloader.MaxBlobSidecarsFetch
	hashes := makePhantomHashes(4 * cap)
	requestData := make([]blobSidecarsRequestData, 4*cap)
	for i := range requestData {
		requestData[i] = blobSidecarsRequestData{
			BlockNum: hexutil.Uint64(i + 1),
			TxIndex:  hexutil.Uint(i),
			Hash:     hashes[i],
		}
	}

	mockCtrl, _, mockPeer, pm := prepareBlockChain(t)
	defer mockCtrl.Finish()

	msg := generateMsg(t, BlobSidecarsRequestMsg, requestData)

	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	// lookups increments once per iteration, capped at 2*cap. Every miss falls
	// through storage to pool, so both fire 2*cap times.
	mockTxPool.EXPECT().GetBlobSidecarFromStorage(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2 * cap)
	mockTxPool.EXPECT().GetBlobSidecarFromPool(gomock.Any()).Return(nil, nil).Times(2 * cap)
	pm.txpool = mockTxPool
	mockPeer.EXPECT().SendBlobSidecarsRLP(gomock.Any()).Return(nil).Times(1)

	assert.NoError(t, handleBlobSidecarsRequestMsg(pm, mockPeer, msg))
}

// generateTestSidecar generates a test BlobTxSidecar for a given transaction hash.
func generateTestSidecar(hash common.Hash) *types.BlobTxSidecar {
	var blob kzg4844.Blob
	var commitment kzg4844.Commitment
	var proof kzg4844.Proof

	copy(blob[:], hash.Bytes())
	copy(commitment[:], hash.Bytes())
	copy(proof[:], hash.Bytes())
	return &types.BlobTxSidecar{
		Version:     types.BlobSidecarVersion1,
		Blobs:       []kzg4844.Blob{blob},
		Commitments: []kzg4844.Commitment{commitment},
		Proofs:      []kzg4844.Proof{proof},
	}
}
