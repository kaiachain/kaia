// Modifications Copyright 2024 The Kaia Authors
// Copyright 2022 The klaytn Authors
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

package snap

import (
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/snapshot"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
)

type testMsgRW struct {
	reader func() (p2p.Msg, error)
	writer func(msg p2p.Msg) error
}

func (rw *testMsgRW) ReadMsg() (p2p.Msg, error)  { return rw.reader() }
func (rw *testMsgRW) WriteMsg(msg p2p.Msg) error { return nil }

type testDownloader struct{}

func (d *testDownloader) DeliverSnapPacket(peer *Peer, packet Packet) error { return nil }

type testSnapshotReader struct {
	db   state.Database
	snap *snapshot.Tree
}

type testKV struct {
	k []byte
	v []byte
}

func NewTestSnapshotReader(items []*testKV) (*testSnapshotReader, common.Hash) {
	memdb := database.NewMemoryDBManager()
	db := state.NewDatabase(memdb)
	trie, _ := statedb.NewTrie(common.Hash{}, db.TrieDB(), nil)
	for _, kv := range items {
		trie.Update(kv.k, kv.v)
	}
	root, _ := trie.Commit(nil)
	db.TrieDB().Commit(root, false, 0)

	snap, _ := snapshot.New(memdb, db.TrieDB(), 256, root, false, true, false)

	return &testSnapshotReader{
		db,
		snap,
	}, root
}

func (r *testSnapshotReader) StateCache() state.Database {
	return r.db
}

func (r *testSnapshotReader) Snapshots() *snapshot.Tree {
	return r.snap
}

func (r *testSnapshotReader) ContractCode(hash common.Hash) ([]byte, error) {
	return r.db.ContractCode(hash)
}

func (r *testSnapshotReader) ContractCodeWithPrefix(hash common.Hash) ([]byte, error) {
	return nil, nil
}

func createMsg(msgcode uint64, data interface{}) (p2p.Msg, error) {
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		return p2p.Msg{}, err
	}
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}, nil
}

func createAccountRangeReqMsg(root common.Hash) p2p.Msg {
	msg, _ := createMsg(GetAccountRangeMsg, &GetAccountRangePacket{
		ID:     1,
		Root:   root,
		Origin: common.Hash{},
		Limit:  common.Hash{},
		Bytes:  softResponseLimit,
	})
	return msg
}

func mockPeer(msg p2p.Msg) *Peer {
	mockPeer := NewFakePeer(1, common.BytesToHash([]byte{0x1}).String(), &testMsgRW{reader: func() (p2p.Msg, error) {
		return msg, nil
	}})
	return mockPeer
}

func TestMessageDecoding(t *testing.T) {
	var (
		msg p2p.Msg
		err error
	)
	msg, err = createMsg(GetAccountRangeMsg, &GetAccountRangePacket{
		ID:     0,
		Root:   common.Hash{},
		Origin: common.Hash{},
		Limit:  common.Hash{},
		Bytes:  0,
	})
	assert.NoError(t, err)
	var req1 GetAccountRangePacket
	assert.NoError(t, msg.Decode(&req1))

	msg, err = createMsg(AccountRangeMsg, &AccountRangePacket{
		ID:       0,
		Accounts: nil,
		Proof:    nil,
	})
	assert.NoError(t, err)
	var req2 AccountRangePacket
	assert.NoError(t, msg.Decode(&req2))

	msg, err = createMsg(GetStorageRangesMsg, &GetStorageRangesPacket{
		ID:       0,
		Root:     common.Hash{},
		Accounts: nil,
		Origin:   nil,
		Limit:    nil,
		Bytes:    0,
	})
	assert.NoError(t, err)
	var req3 GetStorageRangesPacket
	assert.NoError(t, msg.Decode(&req3))

	msg, err = createMsg(StorageRangesMsg, &StorageRangesPacket{
		ID:    0,
		Slots: nil,
		Proof: nil,
	})
	assert.NoError(t, err)
	var req4 StorageRangesPacket
	assert.NoError(t, msg.Decode(&req4))

	msg, err = createMsg(GetByteCodesMsg, &GetByteCodesPacket{
		ID:     0,
		Hashes: nil,
		Bytes:  0,
	})
	assert.NoError(t, err)
	var req5 GetByteCodesPacket
	assert.NoError(t, msg.Decode(&req5))

	msg, err = createMsg(ByteCodesMsg, &ByteCodesPacket{
		ID:    0,
		Codes: nil,
	})
	assert.NoError(t, err)
	var req6 ByteCodesPacket
	assert.NoError(t, msg.Decode(&req6))

	msg, err = createMsg(GetTrieNodesMsg, &GetTrieNodesPacket{
		ID:    0,
		Root:  common.Hash{},
		Paths: nil,
		Bytes: 0,
	})
	assert.NoError(t, err)
	var req7 GetTrieNodesPacket
	assert.NoError(t, msg.Decode(&req7))

	msg, err = createMsg(TrieNodesMsg, &TrieNodesPacket{
		ID:    0,
		Nodes: nil,
	})
	assert.NoError(t, err)
	var req8 TrieNodesPacket
	assert.NoError(t, msg.Decode(&req8))
}

func TestHandleMessage_ReadMsgErr(t *testing.T) {
	reader := &testSnapshotReader{}

	// create test message
	msg, _ := createMsg(GetAccountRangeMsg, []byte{0x1})
	msg.Size = maxMessageSize + 1
	peer := mockPeer(msg)
	testErr := errors.New("test error")
	peer.rw = &testMsgRW{reader: func() (p2p.Msg, error) { return p2p.Msg{}, testErr }}

	// failed to handle message due to read msg error
	err := HandleMessage(reader, &testDownloader{}, peer)
	assert.Equal(t, err, testErr)
}

func TestHandleMessage_LargeMessageErr(t *testing.T) {
	reader := &testSnapshotReader{}

	// create test message
	msg, _ := createMsg(GetAccountRangeMsg, []byte{0x1})
	msg.Size = maxMessageSize + 1
	peer := mockPeer(msg)

	// failed to handle message due to too large message size
	err := HandleMessage(reader, &testDownloader{}, peer)
	assert.True(t, strings.Contains(err.Error(), errMsgTooLarge.Error()))
}

func TestHandleMessage_LargeMessageInvalidMsgErr(t *testing.T) {
	reader := &testSnapshotReader{}

	// create test message
	msg, _ := createMsg(0x08, []byte{0x1})
	peer := mockPeer(msg)

	// failed to handle message due to too large message size
	err := HandleMessage(reader, &testDownloader{}, peer)
	assert.True(t, strings.Contains(err.Error(), errInvalidMsgCode.Error()))
}

func TestHandleMessage_GetAccountRange_EmptyItem(t *testing.T) {
	items := []*testKV{}
	reader, root := NewTestSnapshotReader(items)

	err := HandleMessage(reader, &testDownloader{}, mockPeer(createAccountRangeReqMsg(root)))
	assert.NoError(t, err)
}

func TestHandleMessage_Success(t *testing.T) {
	items := []*testKV{}
	for i := uint64(1); i <= 100; i++ {
		acc, _ := genExternallyOwnedAccount(1, big.NewInt(1))
		serializer := account.NewAccountSerializerWithAccount(acc)
		bytes, _ := rlp.EncodeToBytes(serializer)
		items = append(items, &testKV{key32(i), bytes})
	}

	reader, root := NewTestSnapshotReader(items)
	var (
		msgs []p2p.Msg
		msg  p2p.Msg
		err  error
	)

	msg, err = createMsg(GetAccountRangeMsg, &GetAccountRangePacket{ID: 1, Root: root})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(AccountRangeMsg, &AccountRangePacket{})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(GetStorageRangesMsg, &GetStorageRangesPacket{Root: root})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(StorageRangesMsg, &StorageRangesPacket{})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(GetByteCodesMsg, &GetByteCodesPacket{})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(ByteCodesMsg, &ByteCodesPacket{})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(GetTrieNodesMsg, &GetTrieNodesPacket{Root: root})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	msg, err = createMsg(TrieNodesMsg, &TrieNodesPacket{})
	assert.NoError(t, err)
	msgs = append(msgs, msg)

	for _, msg := range msgs {
		err := HandleMessage(reader, &testDownloader{}, mockPeer(msg))
		assert.NoError(t, err)
	}
}

// TestServiceGetTrieNodesQuery_LookupCapCountsCachedLookups checks that the
// per-request lookup cap (maxTrieNodeLookups) also counts lookups that are
// answered from the already-expanded in-memory trie and cost no database read.
// Before ethereum/go-ethereum#35856 such lookups were free, so a single request
// could carry an unbounded number of them.
func TestServiceGetTrieNodesQuery_LookupCapCountsCachedLookups(t *testing.T) {
	items := []*testKV{}
	for i := uint64(1); i <= 100; i++ {
		acc, _ := genExternallyOwnedAccount(1, big.NewInt(1))
		serializer := account.NewAccountSerializerWithAccount(acc)
		blob, _ := rlp.EncodeToBytes(serializer)
		items = append(items, &testKV{key32(i), blob})
	}
	reader, root := NewTestSnapshotReader(items)

	// Sanity check: the root node (empty path) is served.
	nodes, err := ServiceGetTrieNodesQuery(reader, &GetTrieNodesPacket{
		Root:  root,
		Paths: []TrieNodePathSet{{[]byte{}}},
		Bytes: softResponseLimit,
	}, time.Now())
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.NotEmpty(t, nodes[0])

	// The keys start with bytes 0x01..0x64, so the root is a full node whose
	// child at nibble 0xf is empty. The compact path 0x1f (the single nibble
	// 0xf) is absent from the trie, and its lookup is answered from the
	// in-memory root node without touching the database. Request it more
	// often than the lookup cap allows.
	missing := TrieNodePathSet{[]byte{0x1f}}
	paths := make([]TrieNodePathSet, 2*maxTrieNodeLookups)
	for i := range paths {
		paths[i] = missing
	}
	nodes, err = ServiceGetTrieNodesQuery(reader, &GetTrieNodesPacket{
		Root:  root,
		Paths: paths,
		Bytes: softResponseLimit,
	}, time.Now())
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(nodes), maxTrieNodeLookups+1)
}
