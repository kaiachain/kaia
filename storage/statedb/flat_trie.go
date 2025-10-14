// Copyright 2025 The Kaia Authors
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

package statedb

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/erigontech/erigon-lib/kaiatrie"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

type FlatAccountTrie struct {
	dm *kaiatrie.DomainsManager
	dt *kaiatrie.DeferredAccountTrie

	baseNum uint64
}

func NewFlatAccountTrie(dm *kaiatrie.DomainsManager, root common.Hash, opts *TrieOpts) (*FlatAccountTrie, error) {
	if opts == nil {
		opts = &TrieOpts{}
	}
	dt := kaiatrie.NewDeferredAccountTrie(dm, root.Bytes(), opts.BaseBlockNumber, opts.CommitGenesis)
	return &FlatAccountTrie{dm: dm, dt: dt, baseNum: opts.BaseBlockNumber}, nil
}

func (t *FlatAccountTrie) GetKey(key []byte) []byte {
	// Because FlatNodeIterator returns the unhashed key, no need to lookup preimage.
	return key
}

func (t *FlatAccountTrie) TryGet(key []byte) ([]byte, error) {
	return t.dt.Get(key)
}

func (t *FlatAccountTrie) TryUpdate(key, value []byte) error {
	return t.dt.Put(key, value)
}

func (t *FlatAccountTrie) TryUpdateWithKeys(key, hashKey, hexKey, value []byte) error {
	return t.TryUpdate(key, value)
}

func (t *FlatAccountTrie) TryDelete(key []byte) error {
	// Wipe contract storage
	if err := t.dt.DeleteAccountStorage(key); err != nil {
		return err
	}
	return t.dt.Put(key, nil)
}

func (t *FlatAccountTrie) Hash() common.Hash {
	h, err := t.dt.Hash()
	if err != nil {
		logger.Error("Failed to hash account trie", "err", err)
		return common.Hash{}
	}
	return common.BytesToHash(h[:])
}

func (t *FlatAccountTrie) HashExt() common.ExtHash {
	return t.Hash().ExtendZero()
}

func (t *FlatAccountTrie) Commit(onleaf LeafCallback) (common.Hash, error) {
	h, err := t.dt.Commit()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(h[:]), nil
}

func (t *FlatAccountTrie) CommitExt(onleaf LeafCallback) (common.ExtHash, error) {
	h, err := t.Commit(onleaf)
	if err != nil {
		return common.ExtHash{}, err
	}
	return h.ExtendZero(), nil
}

func (t *FlatAccountTrie) NodeIterator(start []byte) NodeIterator {
	dit, err := kaiatrie.NewAccountIterator(t.dm, t.baseNum)
	if err != nil {
		logger.Error("Failed to create FlatAccountTrie.NodeIterator", "err", err)
		return &EmptyNodeIterator{}
	}
	return &FlatNodeIterator{dit: dit}
}

func (t *FlatAccountTrie) Prove(key []byte, fromLevel uint, proofDb database.DBManager) error {
	logger.Error("FlatAccountTrie.Prove is not implemented")
	return errors.New("not implemented")
}

type FlatStorageTrie struct {
	dm *kaiatrie.DomainsManager
	dt *kaiatrie.DeferredStorageTrie

	addr    common.Address
	baseNum uint64
}

func NewFlatStorageTrie(dm *kaiatrie.DomainsManager, addr common.Address, storageRoot common.Hash, opts *TrieOpts) (*FlatStorageTrie, error) {
	if opts == nil {
		opts = &TrieOpts{}
	}
	if opts.AccountTrie == nil {
		return nil, errors.New("account trie is not set")
	}
	dt := kaiatrie.NewDeferredStorageTrie(opts.AccountTrie.dt, addr.Bytes(), storageRoot.Bytes())
	return &FlatStorageTrie{dm: dm, dt: dt, addr: addr, baseNum: opts.BaseBlockNumber}, nil
}

func (t *FlatStorageTrie) GetKey(key []byte) []byte {
	// Because FlatNodeIterator returns the unhashed key, no need to lookup preimage.
	return key
}

func (t *FlatStorageTrie) TryGet(key []byte) ([]byte, error) {
	value, err := t.dt.Get(key)
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00"))
}

func (t *FlatStorageTrie) TryUpdate(key, value []byte) error {
	_, slot, _, err := rlp.Split(value)
	if err != nil {
		return fmt.Errorf("failed to rlp decode: %w", err)
	}
	return t.dt.Put(key, slot)
}

func (t *FlatStorageTrie) TryUpdateWithKeys(key, hashKey, hexKey, value []byte) error {
	return t.TryUpdate(key, value)
}

func (t *FlatStorageTrie) TryDelete(key []byte) error {
	return t.dt.Put(key, nil)
}

func (t *FlatStorageTrie) Hash() common.Hash {
	h, err := t.dt.Hash()
	if err != nil {
		logger.Error("Failed to hash storage trie", "err", err)
		return common.Hash{}
	}
	return common.BytesToHash(h[:])
}

func (t *FlatStorageTrie) HashExt() common.ExtHash {
	return t.Hash().ExtendZero()
}

func (t *FlatStorageTrie) Commit(onleaf LeafCallback) (common.Hash, error) {
	// Do not commit from StorageTrie. AccountTrie.Commit will eventually commit the storage updates.
	return t.Hash(), nil
}

func (t *FlatStorageTrie) CommitExt(onleaf LeafCallback) (common.ExtHash, error) {
	// Do not commit from StorageTrie. AccountTrie.Commit will eventually commit the storage updates.
	return t.HashExt(), nil
}

func (t *FlatStorageTrie) NodeIterator(start []byte) NodeIterator {
	dit, err := kaiatrie.NewStorageIterator(t.dm, t.addr.Bytes(), t.baseNum)
	if err != nil {
		logger.Error("Failed to create FlatStorageTrie.NodeIterator", "err", err)
		return &EmptyNodeIterator{}
	}
	return &FlatNodeIterator{dit: dit, isStorage: true}
}

func (t *FlatStorageTrie) Prove(key []byte, fromLevel uint, proofDb database.DBManager) error {
	logger.Error("FlatStorageTrie.Prove is not implemented")
	return errors.New("not implemented")
}

type EmptyNodeIterator struct{}

func (nit *EmptyNodeIterator) Next(bool) bool {
	return false
}

func (nit *EmptyNodeIterator) Error() error {
	return nil
}

func (nit *EmptyNodeIterator) Hash() common.Hash {
	return common.Hash{}
}

func (nit *EmptyNodeIterator) Parent() common.Hash {
	return common.Hash{}
}

func (nit *EmptyNodeIterator) Path() []byte {
	return nil
}

func (nit *EmptyNodeIterator) Leaf() bool {
	return true
}

func (nit *EmptyNodeIterator) LeafKey() []byte {
	return nil
}

func (nit *EmptyNodeIterator) LeafBlob() []byte {
	return nil
}

func (nit *EmptyNodeIterator) LeafProof() [][]byte {
	return nil
}

func (nit *EmptyNodeIterator) AddResolver(database.DBManager) {
	// do nothing
}

// TODO: Add close() method
type FlatNodeIterator struct {
	dit       kaiatrie.DomainsIterator
	isStorage bool

	lastKey   []byte
	lastValue []byte
	lastErr   error
}

func NewFlatAccountIterator(t *FlatAccountTrie) (*FlatNodeIterator, error) {
	dit, err := kaiatrie.NewAccountIterator(t.dm, t.baseNum)
	if err != nil {
		return nil, err
	}
	return &FlatNodeIterator{dit: dit}, nil
}

func (nit *FlatNodeIterator) Next(bool) bool {
	var ok bool
	nit.lastKey, nit.lastValue, ok, nit.lastErr = nit.dit.Next()
	if nit.isStorage {
		nit.lastValue, _ = rlp.EncodeToBytes(nit.lastValue)
	}
	if ok && nit.lastErr == nil {
		return true
	} else {
		nit.dit.Close()
		return false
	}
}

func (nit *FlatNodeIterator) Error() error {
	return nit.lastErr
}

func (nit *FlatNodeIterator) Hash() common.Hash {
	return common.BytesToHash(crypto.Keccak256(nit.lastValue))
}

func (nit *FlatNodeIterator) Parent() common.Hash {
	return common.Hash{} // not supported
}

func (nit *FlatNodeIterator) Path() []byte {
	return nil // not supported
}

func (nit *FlatNodeIterator) Leaf() bool {
	return true // always leaf
}

func (nit *FlatNodeIterator) LeafKey() []byte {
	return nit.lastKey
}

func (nit *FlatNodeIterator) LeafBlob() []byte {
	return nit.lastValue
}

func (nit *FlatNodeIterator) LeafProof() [][]byte {
	return nil // not supported
}

func (nit *FlatNodeIterator) AddResolver(database.DBManager) {
	// do nothing
}
