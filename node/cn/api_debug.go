// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package cn

import (
	"context"
	"fmt"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/storage/statedb"
)

// DebugCNAPI is the collection of Kaia full node APIs exposed
// over the public debugging endpoint.
type DebugCNAPI struct {
	cn *CN
}

// NewDebugCNAPI creates a new API definition for the full node-
// related public debug methods of the Kaia service.
func NewDebugCNAPI(cn *CN) *DebugCNAPI {
	return &DebugCNAPI{cn: cn}
}

// DumpBlock retrieves the entire state of the database at a given block.
func (api *DebugCNAPI) DumpBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (state.Dump, error) {
	if *blockNrOrHash.BlockNumber == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, _, stateDb := api.cn.miner.Pending()
		if stateDb == nil {
			return state.Dump{}, fmt.Errorf("pending block is not prepared yet")
		}
		return stateDb.RawDump(), nil
	}

	var block *types.Block
	var err error
	if *blockNrOrHash.BlockNumber == rpc.LatestBlockNumber {
		block = api.cn.APIBackend.CurrentBlock()
	} else {
		block, err = api.cn.APIBackend.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			blockNrOrHashString, _ := blockNrOrHash.NumberOrHashString()
			return state.Dump{}, fmt.Errorf("block %v not found", blockNrOrHashString)
		}
	}
	stateDb, err := api.cn.BlockChain().StateAtWithPersistent(block.Root())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(), nil
}

type Trie struct {
	Type   string `json:"type"`
	Hash   string `json:"hash"`
	Parent string `json:"parent"`
	Path   string `json:"path"`
}

type DumpStateTrieResult struct {
	Root  string `json:"root"`
	Tries []Trie `json:"tries"`
}

// DumpStateTrie retrieves all state/storage tries of the given state root.
func (api *DebugCNAPI) DumpStateTrie(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (DumpStateTrieResult, error) {
	block, err := api.cn.APIBackend.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		blockNrOrHashString, _ := blockNrOrHash.NumberOrHashString()
		return DumpStateTrieResult{}, fmt.Errorf("block #%v not found", blockNrOrHashString)
	}

	result := DumpStateTrieResult{
		Root:  block.Root().String(),
		Tries: make([]Trie, 0),
	}

	db := state.NewDatabaseWithExistingCache(api.cn.chainDB, api.cn.blockchain.StateCache().TrieDB().TrieNodeCache())
	stateDB, err := state.New(block.Root(), db, nil, nil)
	if err != nil {
		return DumpStateTrieResult{}, err
	}
	it := state.NewNodeIterator(stateDB)
	for it.Next() {
		t := Trie{
			it.Type,
			it.Hash.String(),
			it.Parent.String(),
			statedb.HexPathToString(it.Path),
		}

		result.Tries = append(result.Tries, t)
	}
	return result, nil
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// GetBadBLocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *DebugCNAPI) GetBadBlocks(ctx context.Context) ([]blockchain.BadBlockArgs, error) {
	return api.cn.BlockChain().BadBlocks()
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *DebugCNAPI) GetModifiedAccountsByNumber(ctx context.Context, startNum rpc.BlockNumber, endNum *rpc.BlockNumber) ([]common.Address, error) {
	startBlock, endBlock, err := api.getStartAndEndBlock(ctx, startNum, endNum)
	if err != nil {
		return nil, err
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *DebugCNAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.cn.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.cn.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", startHash)
		}
	} else {
		endBlock = api.cn.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
func (api *DebugCNAPI) getModifiedAccounts(startBlock, endBlock *types.Block) ([]common.Address, error) {
	trieDB := api.cn.blockchain.StateCache().TrieDB()

	oldTrie, err := statedb.NewSecureTrie(startBlock.Root(), trieDB, nil)
	if err != nil {
		return nil, err
	}
	newTrie, err := statedb.NewSecureTrie(endBlock.Root(), trieDB, nil)
	if err != nil {
		return nil, err
	}

	diff, _ := statedb.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := statedb.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// getStartAndEndBlock returns start and end block based on the given startNum and endNum.
func (api *DebugCNAPI) getStartAndEndBlock(ctx context.Context, startNum rpc.BlockNumber, endNum *rpc.BlockNumber) (*types.Block, *types.Block, error) {
	var startBlock, endBlock *types.Block

	startBlock, err := api.cn.APIBackend.BlockByNumber(ctx, startNum)
	if err != nil {
		return nil, nil, fmt.Errorf("start block number #%d not found", startNum.Uint64())
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock, err = api.cn.APIBackend.BlockByHash(ctx, startBlock.ParentHash())
		if err != nil {
			return nil, nil, fmt.Errorf("block number #%d has no parent", startNum.Uint64())
		}
	} else {
		endBlock, err = api.cn.APIBackend.BlockByNumber(ctx, *endNum)
		if err != nil {
			return nil, nil, fmt.Errorf("end block number #%d not found", (*endNum).Uint64())
		}
	}

	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}

	return startBlock, endBlock, nil
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// GetModifiedStorageNodesByNumber returns the number of storage nodes of a contract account
// that have been changed between the two blocks specified.
//
// With the first two parameters, it returns the number of storage trie nodes modified in the specified block.
func (api *DebugCNAPI) GetModifiedStorageNodesByNumber(ctx context.Context, contractAddr common.Address, startNum rpc.BlockNumber, endNum *rpc.BlockNumber, printDetail *bool) (int, error) {
	startBlock, endBlock, err := api.getStartAndEndBlock(ctx, startNum, endNum)
	if err != nil {
		return 0, err
	}
	return api.getModifiedStorageNodes(contractAddr, startBlock, endBlock, printDetail)
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
func (api *DebugCNAPI) getModifiedStorageNodes(contractAddr common.Address, startBlock, endBlock *types.Block, printDetail *bool) (int, error) {
	startBlockRoot, err := api.cn.blockchain.GetContractStorageRoot(startBlock, api.cn.blockchain.StateCache(), contractAddr)
	if err != nil {
		return 0, err
	}
	endBlockRoot, err := api.cn.blockchain.GetContractStorageRoot(endBlock, api.cn.blockchain.StateCache(), contractAddr)
	if err != nil {
		return 0, err
	}

	trieDB := api.cn.blockchain.StateCache().TrieDB()
	oldTrie, err := statedb.NewSecureStorageTrie(startBlockRoot, trieDB, nil)
	if err != nil {
		return 0, err
	}
	newTrie, err := statedb.NewSecureStorageTrie(endBlockRoot, trieDB, nil)
	if err != nil {
		return 0, err
	}

	diff, _ := statedb.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := statedb.NewIterator(diff)

	logger.Info("Start collecting the modified storage nodes", "contractAddr", contractAddr.String(),
		"startBlock", startBlock.NumberU64(), "endBlock", endBlock.NumberU64())
	start := time.Now()
	numModifiedNodes := 0
	for iter.Next() {
		numModifiedNodes++
		if printDetail != nil && *printDetail {
			logger.Info("modified storage trie nodes", "contractAddr", contractAddr.String(),
				"nodeHash", common.BytesToHash(iter.Key).String())
		}
	}
	logger.Info("Finished collecting the modified storage nodes", "contractAddr", contractAddr.String(),
		"startBlock", startBlock.NumberU64(), "endBlock", endBlock.NumberU64(), "numModifiedNodes", numModifiedNodes, "elapsed", time.Since(start))
	return numModifiedNodes, nil
}
