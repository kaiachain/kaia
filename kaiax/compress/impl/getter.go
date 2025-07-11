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
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
)

func (c *CompressModule) FindFromCompressedHeader(num uint64, hash common.Hash) ([]byte, bool) {
	return c.findFromCompressed(c.schemas[schemaNameHeader], num, hash)
}

func (c *CompressModule) FindFromCompressedBody(num uint64, hash common.Hash) ([]byte, bool) {
	return c.findFromCompressed(c.schemas[schemaNameBody], num, hash)
}

func (c *CompressModule) FindFromCompressedReceipts(num uint64, hash common.Hash) ([]byte, bool) {
	return c.findFromCompressed(c.schemas[schemaNameReceipts], num, hash)
}

func (c *CompressModule) findFromCompressed(schema ItemSchema, num uint64, hash common.Hash) ([]byte, bool) {
	// TODO: cache the decompressed chunk
	cData, _, _, ok := readCompressed(schema, num)
	if !ok {
		return nil, false
	}
	chunk, err := decompressChunk(c.codec, cData)
	if err != nil {
		return nil, false
	}
	for _, item := range chunk {
		if item.Hash == hash {
			return item.Uncompressed, true
		}
	}
	return nil, false
}

// chunkCache remembers the recent mappings from (num, hash) to uncompressed data.
type chunkCache struct {
	items  *lru.Cache // itemCacheKey -> data []byte
	chunks *lru.Cache // chunkCacheKey -> chunk []ChunkItem
}

type itemCacheKey struct {
	num  uint64
	hash common.Hash
}

type chunkCacheKey struct {
	from uint64
	to   uint64
}

func newChunkCache(chunkCacheSize, itemCacheSize int) *chunkCache {
	chunks, _ := lru.New(chunkCacheSize)
	items, _ := lru.New(itemCacheSize)
	return &chunkCache{chunks: chunks, items: items}
}

func (c *chunkCache) get(num uint64, hash common.Hash) ([]byte, bool) {
	// 1. Try exact item match.
	itemKey := itemCacheKey{num: num, hash: hash}
	item, ok := c.items.Get(itemKey)
	if ok {
		return item.([]byte), true
	}

	// 2. Try containing chunk match.
	for _, key := range c.chunks.Keys() { // search for cached chunk
		chunkKey := key.(chunkCacheKey)
		if num < chunkKey.from || num > chunkKey.to {
			continue
		}
		chunk, ok := c.chunks.Get(key)
		if !ok {
			continue
		}
		// Found the chunk. Pick the item and cache it.
		for _, item := range chunk.([]ChunkItem) {
			if item.Num == num && item.Hash == hash {
				itemKey := itemCacheKey{num: num, hash: hash}
				c.items.Add(itemKey, item.Uncompressed)
				return item.Uncompressed, true
			}
		}
	}

	return nil, false
}

func (c *chunkCache) add(from, to uint64, chunk []ChunkItem, num uint64, hash common.Hash, data []byte) {
	chunkKey := chunkCacheKey{from: from, to: to}
	c.chunks.Add(chunkKey, chunk)
	itemKey := itemCacheKey{num: num, hash: hash}
	c.items.Add(itemKey, data)
}
