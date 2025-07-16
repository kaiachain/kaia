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
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/statedb"
)

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// StartWarmUp retrieves all state/storage tries of the latest committed state root and caches the tries.
func (api *DebugCNStorageAPI) StartWarmUp(minLoad uint) error {
	return api.cn.blockchain.StartWarmUp(minLoad)
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// StartContractWarmUp retrieves a storage trie of the latest state root and caches the trie
// corresponding to the given contract address.
func (api *DebugCNStorageAPI) StartContractWarmUp(contractAddr common.Address, minLoad uint) error {
	return api.cn.blockchain.StartContractWarmUp(contractAddr, minLoad)
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// StopWarmUp stops the warming up process.
func (api *DebugCNStorageAPI) StopWarmUp() error {
	return api.cn.blockchain.StopWarmUp()
}

// TODO-Kaia: Rearrange DebugCNAPI and DebugCNStorageAPI receivers
// StartCollectingTrieStats  collects state/storage trie statistics and print in the log.
func (api *DebugCNStorageAPI) StartCollectingTrieStats(contractAddr common.Address) error {
	return api.cn.blockchain.StartCollectingTrieStats(contractAddr)
}

// DebugCNStorageAPI is the collection of CN full node APIs exposed over
// the private debugging endpoint.
type DebugCNStorageAPI struct {
	config *params.ChainConfig
	cn     *CN
}

// NewDebugCNStorageAPI creates a new API definition for the full node-related
// private debug methods of the CN service.
func NewDebugCNStorageAPI(config *params.ChainConfig, cn *CN) *DebugCNStorageAPI {
	return &DebugCNStorageAPI{config: config, cn: cn}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *DebugCNStorageAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := api.cn.ChainDB().ReadPreimage(hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the statedb.
}

type storageMap map[common.Hash]storageEntry

type storageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

// StorageRangeAt returns the storage at the given block height and transaction index.
func (api *DebugCNStorageAPI) StorageRangeAt(ctx context.Context, blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	// Retrieve the block
	block := api.cn.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return StorageRangeResult{}, fmt.Errorf("block %#x not found", blockHash)
	}
	_, _, _, statedb, release, err := api.cn.stateAtTransaction(block, txIndex, 0, nil, true, false)
	if err != nil {
		return StorageRangeResult{}, err
	}
	defer release()

	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}

func storageRangeAt(st state.Trie, start []byte, maxResult int) (StorageRangeResult, error) {
	it := statedb.NewIterator(st.NodeIterator(start))
	result := StorageRangeResult{Storage: storageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return StorageRangeResult{}, err
		}
		e := storageEntry{Value: common.BytesToHash(content)}
		if preimage := st.GetKey(it.Key); preimage != nil {
			preimage := common.BytesToHash(preimage)
			e.Key = &preimage
		}
		result.Storage[common.BytesToHash(it.Key)] = e
	}
	// Add the 'next key' so clients can continue downloading.
	if it.Next() {
		next := common.BytesToHash(it.Key)
		result.NextKey = &next
	}
	return result, nil
}
