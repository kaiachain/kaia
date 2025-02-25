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
	Chain          chain
	DBM            dbm
	CompressConfig compress.CompressConfig
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

	if err := c.CompressConfig.Validate(); err != nil {
		return err
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
	logger.Info("Starting compress module", "compressing", c.CompressConfig.Enabled)

	// Reset the quit state.
	c.quit.Store(0)

	// Start the compression threads.
	if c.CompressConfig.Enabled {
		for _, schema := range c.schemas {
			c.wg.Add(1)
			go c.compress(schema)
		}
	}
	return nil
}

func (c *CompressModule) Stop() {
	c.quit.Store(1)
	c.wg.Wait()
	logger.Info("Compress module stopped")
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

	context := newChunkContext(c.DBM, schema, c.codec, c.CompressConfig.ChunkItemCap, c.CompressConfig.ChunkByteCap, *nextNum)

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
	if currNum > c.CompressConfig.Retention {
		endNum = currNum - c.CompressConfig.Retention
	}
	return endNum
}
