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
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/work"
)

// AdminCNChainAPI is the collection of CN full node-related APIs
// exposed over the private admin endpoint.
type AdminCNChainAPI struct {
	cn *CN
}

// NewAdminCNChainAPI creates a new API definition for the full node private
// admin methods of the CN service.
func NewAdminCNChainAPI(cn *CN) *AdminCNChainAPI {
	return &AdminCNChainAPI{cn: cn}
}

// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil.
func (api *AdminCNChainAPI) ExportChain(file string, first, last *rpc.BlockNumber) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		// File already exists. Allowing overwrite could be a DoS vecotor,
		// since the 'file' may point to arbitrary paths on the drive
		return false, errors.New("location would overwrite an existing file")
	}
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first == nil {
		zero := rpc.EarliestBlockNumber
		first = &zero
	}
	if last == nil || *last == rpc.LatestBlockNumber {
		head := rpc.BlockNumber(api.cn.BlockChain().CurrentBlock().NumberU64())
		last = &head
	}

	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if err := api.cn.BlockChain().ExportN(writer, first.Uint64(), last.Uint64()); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain work.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *AdminCNChainAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}
	stream := rlp.NewStream(reader, 0)

	return api.importChain(stream)
}

func (api *AdminCNChainAPI) ImportChainFromString(blockRlp string) (bool, error) {
	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(bytes.NewReader(common.FromHex(blockRlp)), 0)

	return api.importChain(stream)
}

func (api *AdminCNChainAPI) importChain(stream *rlp.Stream) (bool, error) {
	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.cn.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.cn.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// StartStateMigration starts state migration.
func (api *AdminCNChainAPI) StartStateMigration() error {
	return api.cn.blockchain.PrepareStateMigration()
}

// StopStateMigration stops state migration and removes stateMigrationDB.
func (api *AdminCNChainAPI) StopStateMigration() error {
	return api.cn.BlockChain().StopStateMigration()
}

// StateMigrationStatus returns the status information of state trie migration.
func (api *AdminCNChainAPI) StateMigrationStatus() map[string]interface{} {
	isMigration, blkNum, read, committed, pending, progress, err := api.cn.BlockChain().StateMigrationStatus()

	errStr := "null"
	if err != nil {
		errStr = err.Error()
	}

	return map[string]interface{}{
		"isMigration":          isMigration,
		"migrationBlockNumber": blkNum,
		"read":                 read,
		"committed":            committed,
		"pending":              pending,
		"progress":             progress,
		"err":                  errStr,
	}
}

func (api *AdminCNChainAPI) SaveTrieNodeCacheToDisk() error {
	return api.cn.BlockChain().SaveTrieNodeCacheToDisk()
}

func (api *AdminCNChainAPI) SpamThrottlerConfig(ctx context.Context) (*blockchain.ThrottlerConfig, error) {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return nil, errors.New("spam throttler is not running")
	}
	return throttler.GetConfig(), nil
}

func (api *AdminCNChainAPI) StopSpamThrottler(ctx context.Context) error {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return errors.New("spam throttler was already stopped")
	}
	api.cn.txPool.StopSpamThrottler()
	return nil
}

func (api *AdminCNChainAPI) StartSpamThrottler(ctx context.Context, config *blockchain.ThrottlerConfig) error {
	throttler := blockchain.GetSpamThrottler()
	if throttler != nil {
		return errors.New("spam throttler is already running")
	}
	return api.cn.txPool.StartSpamThrottler(config)
}

func (api *AdminCNChainAPI) SetSpamThrottlerWhiteList(ctx context.Context, addrs []common.Address) error {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return errors.New("spam throttler is not running")
	}
	throttler.SetAllowed(addrs)
	return nil
}

func (api *AdminCNChainAPI) GetSpamThrottlerWhiteList(ctx context.Context) ([]common.Address, error) {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return nil, errors.New("spam throttler is not running")
	}
	return throttler.GetAllowed(), nil
}

func (api *AdminCNChainAPI) GetSpamThrottlerThrottleList(ctx context.Context) ([]common.Address, error) {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return nil, errors.New("spam throttler is not running")
	}
	return throttler.GetThrottled(), nil
}

func (api *AdminCNChainAPI) GetSpamThrottlerCandidateList(ctx context.Context) (map[common.Address]int, error) {
	throttler := blockchain.GetSpamThrottler()
	if throttler == nil {
		return nil, errors.New("spam throttler is not running")
	}
	return throttler.GetCandidates(), nil
}

func (s *AdminCNChainAPI) NodeConfig(ctx context.Context) interface{} {
	return *s.cn.config
}
