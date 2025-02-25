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

package impl

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ compress.CompressModule = (*CompressModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxCompress)

	chunkCacheSize = 1000
	itemCacheSize  = 1000

	// Periodic log in compressContext.until(). Once in few minutes.
	compressLogInterval = uint64(1024000)
	// Periodic compaction in compressContext.until(). After 1 month worth blocks.
	compressCompactionPeriod = uint64(86400 * 30)
)

type chain interface {
	CurrentBlock() *types.Block
}

type blockHashReader interface {
	ReadCanonicalHash(num uint64) common.Hash
}

type dbm interface {
	blockHashReader
	GetHeaderDB() database.Database
	GetBodyDB() database.Database
	GetReceiptsDB() database.Database
	GetCompressedHeaderDB() database.Database
	GetCompressedBodyDB() database.Database
	GetCompressedReceiptsDB() database.Database
}

type InitOpts struct {
	Chain        chain
	DBM          dbm
	Retention    uint64 // number of blocks to keep in the uncompressed database
	ChunkItemCap int    // maximum number of items in a chunk
	ChunkByteCap int    // maximum size of uncompressed data in a chunk
}

type CompressModule struct {
	InitOpts

	codec   Codec
	schemas map[string]ItemSchema
	caches  map[string]*chunkCache

	// Stops long-running tasks.
	quit atomic.Int32
	wg   sync.WaitGroup
}

func NewCompressModule() *CompressModule {
	return &CompressModule{
		schemas: make(map[string]ItemSchema),
		caches:  make(map[string]*chunkCache),
	}
}

func (c *CompressModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Chain == nil || opts.DBM == nil {
		return ErrInitUnexpectedNil
	}
	c.InitOpts = *opts

	if c.ChunkItemCap < compress.MinChunkItemCap || c.ChunkItemCap > compress.MaxChunkItemCap {
		return ErrInvalidConfig(c.Retention, c.ChunkItemCap, c.ChunkByteCap)
	}
	if c.ChunkByteCap < compress.MinChunkByteCap || c.ChunkByteCap > compress.MaxChunkByteCap {
		return ErrInvalidConfig(c.Retention, c.ChunkItemCap, c.ChunkByteCap)
	}
	if c.Retention < compress.MinRetention {
		return ErrInvalidConfig(c.Retention, c.ChunkItemCap, c.ChunkByteCap)
	}

	c.codec = NewZstdCodec()
	c.setSchemas(
		NewHeaderSchema(c.DBM.GetHeaderDB(), c.DBM.GetCompressedHeaderDB()),
		NewBodySchema(c.DBM.GetBodyDB(), c.DBM.GetCompressedBodyDB()),
		NewReceiptSchema(c.DBM.GetReceiptsDB(), c.DBM.GetCompressedReceiptsDB()),
	)
	for _, schema := range c.schemas {
		c.initSchema(schema)
	}
	return nil
}

func (c *CompressModule) setSchemas(s ...ItemSchema) {
	c.schemas = make(map[string]ItemSchema)
	for _, schema := range s {
		c.schemas[schema.name()] = schema
		c.caches[schema.name()] = newChunkCache(chunkCacheSize, itemCacheSize)
	}
}

func (c *CompressModule) Start() error {
	// Reset the quit state.
	c.quit.Store(0)

	// Start the compression threads.
	for _, schema := range c.schemas {
		c.wg.Add(1)
		go c.compress(schema)
	}
	return nil
}

func (c *CompressModule) Stop() {
	c.quit.Store(1)
	c.wg.Wait()
}

func (c *CompressModule) initSchema(schema ItemSchema) {
	nextNum := readNextNum(schema)
	if nextNum == nil {
		writeNextNum(schema, 1) // do not delete block 0.
	}
}

func (c *CompressModule) compress(schema ItemSchema) {
	defer c.wg.Done()

	nextNum := readNextNum(schema)
	if nextNum == nil {
		logger.Error("No next compression number for schema %s", schema.name())
		return
	}
	logger.Info("Start compressing", "schema", schema.name(), "nextNum", *nextNum)

	context := newChunkContext(c.DBM, schema, c.codec, c.ChunkItemCap, c.ChunkByteCap, *nextNum)

	for {
		endNum := c.compressEndNum()
		if err := context.until(endNum, &c.quit); err != nil {
			logger.Error("Database compression failed", "schema", schema.name(), "err", err)
			return
		}

		// No gap detected. Sleep a while and check again.
		time.Sleep(time.Second)
		if c.quit.Load() != 0 {
			return
		}
	}
}

// compressEndNum returns the number of the last block to be compressed.
// The current block number minus retention or 0.
func (c *CompressModule) compressEndNum() uint64 {
	currNum := c.Chain.CurrentBlock().NumberU64()
	endNum := uint64(0)
	if currNum > c.Retention {
		endNum = currNum - c.Retention
	}
	return endNum
}
