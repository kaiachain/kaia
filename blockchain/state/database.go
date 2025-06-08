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
// This file is derived from core/state/database.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package state

import (
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/lru"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
)

const (
	// Number of codehash->size associations to keep
	codeSizeCacheSize = 100000

	// Cache size granted for caching clean code.
	codeCacheSize = 64 * 1024 * 1024

	// Number of shards in cache
	shardsCodeSizeCache = 4096
)

// Database wraps access to tries and contract code.
type Database interface {
	// OpenTrie opens the main account trie.
	OpenTrie(root common.Hash, opts *statedb.TrieOpts) (Trie, error)

	// OpenStorageTrie opens the storage trie of an account.
	OpenStorageTrie(root common.ExtHash, opts *statedb.TrieOpts) (Trie, error)

	// CopyTrie returns an independent copy of the given trie.
	CopyTrie(Trie) Trie

	// ContractCode retrieves a particular contract's code.
	ContractCode(codeHash common.Hash) ([]byte, error)

	// DeleteCode deletes a particular contract's code.
	DeleteCode(codeHash common.Hash)

	// ContractCodeSize retrieves a particular contracts code's size.
	ContractCodeSize(codeHash common.Hash) (int, error)

	// TrieDB retrieves the low level trie database used for data storage.
	TrieDB() *statedb.Database

	// RLockGCCachedNode locks the GC lock of CachedNode.
	RLockGCCachedNode()

	// RUnlockGCCachedNode unlocks the GC lock of CachedNode.
	RUnlockGCCachedNode()
}

// Trie is a Kaia Merkle Patricia trie.
type Trie interface {
	// StartPruningSnapshot starts a pruning snapshot.
	StartPruningSnapshot()
	// EndPruningSnapshot ends a pruning snapshot.
	EndPruningSnapshot()
	// RevertPruningSnapshot reverts a pruning snapshot.
	RevertPruningSnapshot()

	// GetKey returns the sha3 preimage of a hashed key that was previously used
	// to store a value.
	//
	// TODO(fjl): remove this when SecureTrie is removed
	GetKey([]byte) []byte
	// TryGet returns the value for key stored in the trie. The value bytes must
	// not be modified by the caller. If a node was not found in the database, a
	// trie.MissingNodeError is returned.
	TryGet(key []byte) ([]byte, error)
	// TryUpdate associates key with value in the trie. If value has length zero, any
	// existing value is deleted from the trie. The value bytes must not be modified
	// by the caller while they are stored in the trie. If a node was not found in the
	// database, a trie.MissingNodeError is returned.
	TryUpdate(key, value []byte) error
	TryUpdateWithKeys(key, hashKey, hexKey, value []byte) error
	// TryDelete removes any existing value for key from the trie. If a node was not
	// found in the database, a trie.MissingNodeError is returned.
	TryDelete(key []byte) error
	// Hash returns the root hash of the trie. It does not write to the database and
	// can be used even if the trie doesn't have one.
	Hash() common.Hash
	HashExt() common.ExtHash
	// Commit writes all nodes to the trie's memory database, tracking the internal
	// and external (for account tries) references.
	Commit(onleaf statedb.LeafCallback) (common.Hash, error)
	CommitExt(onleaf statedb.LeafCallback) (common.ExtHash, error)
	// NodeIterator returns an iterator that returns nodes of the trie. Iteration
	// starts at the key after the given start key.
	NodeIterator(startKey []byte) statedb.NodeIterator
	// Prove constructs a Merkle proof for key. The result contains all encoded nodes
	// on the path to the value at key. The value itself is also included in the last
	// node and can be retrieved by verifying the proof.
	//
	// If the trie does not contain a value for key, the returned proof contains all
	// nodes of the longest existing prefix of the key (at least the root), ending
	// with the node that proves the absence of the key.
	Prove(key []byte, fromLevel uint, proofDb database.DBManager) error
}

// NewDatabase creates a backing store for state. The returned database is safe for
// concurrent use, but does not retain any recent trie nodes in memory. To keep some
// historical state in memory, use the NewDatabaseWithNewCache constructor.
func NewDatabase(db database.DBManager) Database {
	return NewDatabaseWithNewCache(db, statedb.GetEmptyTrieNodeCacheConfig())
}

func getCodeSizeCache() common.Cache {
	var cacheConfig common.CacheConfiger
	switch common.DefaultCacheType {
	case common.LRUShardCacheType:
		cacheConfig = common.LRUShardConfig{CacheSize: codeSizeCacheSize, NumShards: shardsCodeSizeCache}
	case common.LRUCacheType:
		cacheConfig = common.LRUConfig{CacheSize: codeSizeCacheSize}
	case common.FIFOCacheType:
		cacheConfig = common.FIFOCacheConfig{CacheSize: codeSizeCacheSize}
	default:
		cacheConfig = common.FIFOCacheConfig{CacheSize: codeSizeCacheSize}
	}

	return common.NewCache(cacheConfig)
}

// NewDatabaseWithNewCache creates a backing store for state. The returned database
// is safe for concurrent use and retains a lot of collapsed RLP trie nodes in a
// large memory cache.
func NewDatabaseWithNewCache(db database.DBManager, cacheConfig *statedb.TrieNodeCacheConfig) Database {
	return &cachingDB{
		db:            statedb.NewDatabaseWithNewCache(db, cacheConfig),
		codeSizeCache: getCodeSizeCache(),
		codeCache:     lru.NewSizeConstrainedCache[common.Hash, []byte](codeCacheSize),
	}
}

// NewDatabaseWithExistingCache creates a backing store for state with given cache. The returned database
// is safe for concurrent use and retains a lot of collapsed RLP trie nodes in a
// large memory cache.
func NewDatabaseWithExistingCache(db database.DBManager, cache statedb.TrieNodeCache) Database {
	return &cachingDB{
		db:            statedb.NewDatabaseWithExistingCache(db, cache),
		codeSizeCache: getCodeSizeCache(),
		codeCache:     lru.NewSizeConstrainedCache[common.Hash, []byte](codeCacheSize),
	}
}

type cachingDB struct {
	db            *statedb.Database
	codeSizeCache common.Cache
	codeCache     *lru.SizeConstrainedCache[common.Hash, []byte]
}

// OpenTrie opens the main account trie at a specific root hash.
func (db *cachingDB) OpenTrie(root common.Hash, opts *statedb.TrieOpts) (Trie, error) {
	return statedb.NewSecureTrie(root, db.db, opts)
}

// OpenStorageTrie opens the storage trie of an account.
func (db *cachingDB) OpenStorageTrie(root common.ExtHash, opts *statedb.TrieOpts) (Trie, error) {
	return statedb.NewSecureStorageTrie(root, db.db, opts)
}

// CopyTrie returns an independent copy of the given trie.
func (db *cachingDB) CopyTrie(t Trie) Trie {
	switch t := t.(type) {
	case *statedb.SecureTrie:
		return t.Copy()
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

// ContractCode retrieves a particular contract's code.
func (db *cachingDB) ContractCode(codeHash common.Hash) ([]byte, error) {
	if code, _ := db.codeCache.Get(codeHash); len(code) > 0 {
		return code, nil
	}
	code := db.db.DiskDB().ReadCode(codeHash)
	if len(code) > 0 {
		db.codeCache.Add(codeHash, code)
		db.codeSizeCache.Add(codeHash, len(code))
		return code, nil
	}
	return nil, errors.New("not found")
}

// DeleteCode deletes a particular contract's code.
func (db *cachingDB) DeleteCode(codeHash common.Hash) {
	db.codeCache.DeleteCode(codeHash)
	db.db.DiskDB().DeleteCode(codeHash)
}

// ContractCodeWithPrefix retrieves a particular contract's code. If the
// code can't be found in the cache, then check the existence with **new**
// db scheme.
func (db *cachingDB) ContractCodeWithPrefix(codeHash common.Hash) ([]byte, error) {
	if code, _ := db.codeCache.Get(codeHash); len(code) > 0 {
		return code, nil
	}
	code := db.db.DiskDB().ReadCodeWithPrefix(codeHash)
	if len(code) > 0 {
		db.codeCache.Add(codeHash, code)
		db.codeSizeCache.Add(codeHash, len(code))
		return code, nil
	}
	return nil, errors.New("not found")
}

// ContractCodeSize retrieves a particular contracts code's size.
func (db *cachingDB) ContractCodeSize(codeHash common.Hash) (int, error) {
	if cached, ok := db.codeSizeCache.Get(codeHash); ok {
		return cached.(int), nil
	}
	code, err := db.ContractCode(codeHash)
	return len(code), err
}

// TrieDB retrieves the low level trie database used for data storage.
func (db *cachingDB) TrieDB() *statedb.Database {
	return db.db
}

// RLockGCCachedNode locks the GC lock of CachedNode.
func (db *cachingDB) RLockGCCachedNode() {
	db.db.RLockGCCachedNode()
}

// RUnlockGCCachedNode unlocks the GC lock of CachedNode.
func (db *cachingDB) RUnlockGCCachedNode() {
	db.db.RUnlockGCCachedNode()
}
