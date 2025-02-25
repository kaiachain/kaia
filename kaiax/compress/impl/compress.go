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
	"sync/atomic"

	"github.com/kaiachain/kaia/common"
)

// ChunkItem is an element in a compression chunk.
type ChunkItem struct {
	Num          uint64
	Hash         common.Hash
	Uncompressed []byte // opaque uncompressed data
}

// compressContext is a context for continuously compressing items
type compressContext struct {
	// implementations
	dbm    blockHashReader
	schema ItemSchema
	codec  Codec

	// configuration
	chunkItemCap int
	chunkByteCap int

	// current chunk state
	fromNum    uint64 // the first number of the chunk
	currNum    uint64 // the number to be compressed at this step
	chunk      []ChunkItem
	chunkBytes int // sum of byte sizes of uncompressed items
	compacting atomic.Int32
}

func newChunkContext(dbm blockHashReader, schema ItemSchema, codec Codec, chunkItemCap, chunkByteCap int, fromNum uint64) *compressContext {
	return &compressContext{
		dbm:          dbm,
		schema:       schema,
		codec:        codec,
		chunkItemCap: chunkItemCap,
		chunkByteCap: chunkByteCap,
		fromNum:      fromNum,
		currNum:      fromNum,
	}
}

func (c *compressContext) chunkCapReached() bool {
	return len(c.chunk) >= c.chunkItemCap || c.chunkBytes >= c.chunkByteCap
}

func (c *compressContext) shouldCompact() bool {
	return c.currNum%compressCompactionPeriod == 0 && c.compacting.Load() == 0
}

func (c *compressContext) compactDetached() {
	go func() {
		c.compacting.Store(1)
		defer c.compacting.Store(0)
		compactUncompressed(c.schema)
	}()
}

// step appends one item to the context. If the capping condition is met, write the compressed chunk.
func (c *compressContext) step() error {
	// Read and append an uncompressed item
	hash := c.dbm.ReadCanonicalHash(c.currNum)
	uData := readUncompressed(c.schema, c.currNum, hash)
	c.chunk = append(c.chunk, ChunkItem{
		Num:          c.currNum,
		Hash:         hash,
		Uncompressed: uData,
	})
	c.chunkBytes += len(uData)

	if c.chunkCapReached() {
		cData, err := compressChunk(c.codec, c.chunk)
		if err != nil {
			return ErrCodecCompress(err)
		}
		// chunk = [fromNum, currNum]
		writeCompressed(c.schema, c.fromNum, c.currNum, cData)
		writeNextNum(c.schema, c.currNum+1)
		deleteUncompressedBatch(c.schema, c.chunk)

		// Start new chunk
		c.chunk = nil
		c.chunkBytes = 0
		c.currNum++
		c.fromNum = c.currNum
	} else {
		c.currNum++
	}

	if c.shouldCompact() {
		c.compactDetached()
	}
	return nil
}

// until loops until the current number reaches the given number or the quit signal is received.
func (c *compressContext) until(untilNum uint64, quit *atomic.Int32) error {
	for c.currNum <= untilNum {
		if quit.Load() != 0 {
			return nil
		}
		if err := c.step(); err != nil {
			return err
		}
		if (c.currNum % compressLogInterval) == 0 {
			logger.Info("Compressing chunk", "schema", c.schema.name(), "from", c.fromNum, "curr", c.currNum)
		}
	}
	return nil
}
