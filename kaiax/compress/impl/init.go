// Copyright 2024 The Kaia Authors
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

package compress

import (
	"sync"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
)

const MB_1 = uint64(1000000)

var (
	_      compress.CompressionModule = &CompressModule{}
	logger                            = log.NewModuleLogger(log.KaiaxCompress)
)

//go:generate mockgen -destination=mock/blockchain_mock.go github.com/kaiachain/kaia/kaiax/compress/impl BlockChain
type BlockChain interface {
	CurrentBlock() *types.Block
}

type InitOpts struct {
	ChunkBlockSize uint64
	ChunkCap       uint64
	Retention      uint64
	Enable         bool
	Chain          BlockChain
	Dbm            database.DBManager
}

type IdleState struct {
	isIdle   bool
	idleTime time.Duration
}

type CompressModule struct {
	InitOpts

	loopIdleTime time.Duration

	terminateCompress chan any
	wg                sync.WaitGroup

	compressChunkMu   sync.RWMutex
	compressMaxSizeMu sync.RWMutex
	compressRetention sync.RWMutex

	headerIdleState     *IdleState
	bodyIdleState       *IdleState
	receiptsIdleState   *IdleState
	headerIdleStateMu   sync.RWMutex
	bodyIdleStateMu     sync.RWMutex
	receiptsIdleStateMu sync.RWMutex
}

func NewCompression() *CompressModule {
	return &CompressModule{}
}

func (c *CompressModule) setMaxSize(v uint64) {
	c.compressMaxSizeMu.Lock()
	defer c.compressMaxSizeMu.Unlock()
	c.InitOpts.ChunkCap = v
}

func (c *CompressModule) getChunkCap() uint64 {
	c.compressMaxSizeMu.RLock()
	defer c.compressMaxSizeMu.RUnlock()
	if c.InitOpts.ChunkCap == 0 {
		return blockchain.DefaultCompressChunkCap
	}
	return c.InitOpts.ChunkCap
}

func (c *CompressModule) setCompressChunk(v uint64) {
	c.compressChunkMu.Lock()
	defer c.compressChunkMu.Unlock()
	c.InitOpts.ChunkBlockSize = v
}

func (c *CompressModule) getCompressChunk() uint64 {
	c.compressChunkMu.RLock()
	defer c.compressChunkMu.RUnlock()
	if c.InitOpts.ChunkBlockSize == 0 {
		return blockchain.DefaultChunkBlockSize
	}
	return c.InitOpts.ChunkBlockSize
}

func (c *CompressModule) setCompressRetention(v uint64) {
	c.compressRetention.Lock()
	defer c.compressRetention.Unlock()
	c.InitOpts.Retention = v
}

func (c *CompressModule) getCompressRetention() uint64 {
	c.compressRetention.RLock()
	defer c.compressRetention.RUnlock()
	return c.InitOpts.Retention
}

func (c *CompressModule) setIdleState(compressTyp CompressionType, is *IdleState) {
	switch compressTyp {
	case HeaderCompressType:
		c.headerIdleStateMu.Lock()
		defer c.headerIdleStateMu.Unlock()
		c.headerIdleState = is
	case BodyCompressType:
		c.bodyIdleStateMu.Lock()
		defer c.bodyIdleStateMu.Unlock()
		c.bodyIdleState = is
	case ReceiptCompressType:
		c.receiptsIdleStateMu.Lock()
		defer c.receiptsIdleStateMu.Unlock()
		c.receiptsIdleState = is
	}
}

func (c *CompressModule) getIdleState(compressTyp CompressionType) IdleState {
	switch compressTyp {
	case HeaderCompressType:
		c.headerIdleStateMu.RLock()
		defer c.headerIdleStateMu.RUnlock()
		if c.headerIdleState == nil {
			return IdleState{}
		}
		return *c.headerIdleState
	case BodyCompressType:
		c.bodyIdleStateMu.RLock()
		defer c.bodyIdleStateMu.RUnlock()
		if c.bodyIdleState == nil {
			return IdleState{}
		}
		return *c.bodyIdleState
	case ReceiptCompressType:
		c.receiptsIdleStateMu.RLock()
		defer c.receiptsIdleStateMu.RUnlock()
		if c.receiptsIdleState == nil {
			return IdleState{}
		}
		return *c.receiptsIdleState
	default:
		panic("unreachable")
	}
}

func (c *CompressModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Chain == nil || opts.Dbm == nil {
		return compress.ErrInitNil
	}
	c.InitOpts = *opts
	c.loopIdleTime = time.Second
	c.terminateCompress = make(chan any, TotalCompressTypeSize)
	initCache(c.InitOpts.ChunkCap)
	return nil
}

func (c *CompressModule) Start() error {
	if !c.Enable {
		logger.Info("[Compression] Compression disabled")
		return nil
	}
	logger.Info("[Compression] Compression started")
	c.Compress()
	return nil
}

func (c *CompressModule) Stop() {
	if !c.Enable {
		return
	}
	c.stopCompress()
	clearCache()
	logger.Info("[Compression] Compression Stopped")
}
