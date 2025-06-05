// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/state/sync.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package state

import (
	"bytes"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/v2/blockchain/types/account"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/rlp"
	"github.com/kaiachain/kaia/v2/storage/statedb"
)

// NewStateSync create a new state trie download scheduler.
// LRU cache is mendatory when state syncing and block processing are executed simultaneously
func NewStateSync(root common.Hash, database statedb.StateTrieReadDB, bloom *statedb.SyncBloom, lruCache *lru.Cache, onLeaf func(paths [][]byte, leaf []byte) error) *statedb.TrieSync {
	// Register the storage slot callback if the external callback is specified.
	var onSlot statedb.LeafCallback
	if onLeaf != nil {
		onSlot = func(paths [][]byte, _ []byte, leaf []byte, _ common.ExtHash, _ int) error {
			return onLeaf(paths, leaf)
		}
	}
	// Register the account callback to connect the state trie and the storage
	// trie belongs to the contract.
	var syncer *statedb.TrieSync
	onAccount := func(paths [][]byte, hexpath []byte, leaf []byte, parent common.ExtHash, parentDepth int) error {
		if onLeaf != nil {
			if err := onLeaf(paths, leaf); err != nil {
				return err
			}
		}
		serializer := account.NewAccountSerializer()
		if err := rlp.Decode(bytes.NewReader(leaf), serializer); err != nil {
			return err
		}
		obj := serializer.GetAccount()
		if pa := account.GetProgramAccount(obj); pa != nil {
			syncer.AddSubTrie(pa.GetStorageRoot().Unextend(), hexpath, parentDepth+1, parent.Unextend(), onSlot)
			syncer.AddCodeEntry(common.BytesToHash(pa.GetCodeHash()), hexpath, parentDepth+1, parent.Unextend())
		}
		return nil
	}
	syncer = statedb.NewTrieSync(root, database, onAccount, bloom, lruCache)
	return syncer
}
