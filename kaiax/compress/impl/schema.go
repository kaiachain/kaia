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
	"encoding/binary"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ ItemSchema = (*HeaderSchema)(nil)

	schemaNameHeader             = "header"
	schemaNameBody               = "body"
	schemaNameReceipts           = "receipts"
	compressedHeaderKeyPrefix    = []byte("Compressed-h")
	compressedBodyKeyPrefix      = []byte("Compressed-b")
	compressedReceiptsKeyPrefix  = []byte("Compressed-r")
	compressedHeaderNextNumKey   = []byte("NextCompressingNum-h")
	compressedBodyNextNumKey     = []byte("NextCompressingNum-b")
	compressedReceiptsNextNumKey = []byte("NextCompressingNum-r")
)

// chunkKey = prefix || to || from
// Be careful with the order of `to` and `from`. This is for quick search using iterator.
func chunkKey(prefix []byte, from, to uint64) []byte {
	bFrom := make([]byte, 8)
	binary.BigEndian.PutUint64(bFrom, from)
	bTo := make([]byte, 8)
	binary.BigEndian.PutUint64(bTo, to)
	return append(append(prefix, bTo...), bFrom...)
}

func parseChunkKey(prefix, key []byte) (from, to uint64, ok bool) {
	if len(key) != len(prefix)+16 {
		return 0, 0, false
	}
	to = binary.BigEndian.Uint64(key[len(prefix):])
	from = binary.BigEndian.Uint64(key[len(prefix)+8:])
	return from, to, true
}

// ItemSchema is a generic interface for a pair of uncompressed and compressed databases.
type ItemSchema interface {
	name() string                                        // name that appears in logs
	uncompressedDb() database.Database                   // uncompressed database handle (e.g. HeaderDB)
	compressedDb() database.Database                     // compressed database handle (e.g. CompressedHeaderDB)
	uncompressedKey(num uint64, hash common.Hash) []byte // key for uncompressed item
	compressedKeyPrefix() []byte                         // key prefix for compressed chunk
	nextNumKey() []byte                                  // key for next number to be compressed
}

func readUncompressed(schema ItemSchema, num uint64, hash common.Hash) []byte {
	uDb := schema.uncompressedDb()
	uKey := schema.uncompressedKey(num, hash)
	uData, _ := uDb.Get(uKey)
	return uData
}

func writeUncompressedBatch(schema ItemSchema, items []ChunkItem) {
	uDb := schema.uncompressedDb()
	batch := uDb.NewBatch()
	for _, item := range items {
		uKey := schema.uncompressedKey(item.Num, item.Hash)
		if err := batch.Put(uKey, item.Uncompressed); err != nil {
			logger.Crit("Failed to write uncompressed data", "schema", schema.name(), "num", item.Num, "hash", item.Hash, "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to write uncompressed data", "schema", schema.name(), "err", err)
	}
}

func deleteUncompressedBatch(schema ItemSchema, items []ChunkItem) {
	uDb := schema.uncompressedDb()
	batch := uDb.NewBatch()
	for _, item := range items {
		uKey := schema.uncompressedKey(item.Num, item.Hash)
		if err := batch.Delete(uKey); err != nil {
			logger.Crit("Failed to delete uncompressed data", "schema", schema.name(), "num", item.Num, "hash", item.Hash, "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to write uncompressed data", "schema", schema.name(), "err", err)
	}
}

// Read the chunk that contains the block number `num`.
// Returns compressed data, from, to, and succeess flag.
func readCompressed(schema ItemSchema, num uint64) ([]byte, uint64, uint64, bool) {
	cDb := schema.compressedDb()
	prefix := schema.compressedKeyPrefix()
	bNum := make([]byte, 8)
	binary.BigEndian.PutUint64(bNum, num)

	// Find the first key `prefix || to || from` where `num <= to`. Then check that `from <= num`.
	it := cDb.NewIterator(prefix, bNum)
	if !it.Next() {
		return nil, 0, 0, false
	}
	key := it.Key()
	if from, to, ok := parseChunkKey(prefix, key); ok && from <= num && to >= num {
		return it.Value(), from, to, true
	} else {
		return nil, 0, 0, false
	}
}

func writeCompressed(schema ItemSchema, from, to uint64, data []byte) {
	cDb := schema.compressedDb()
	cKey := chunkKey(schema.compressedKeyPrefix(), from, to)
	if err := cDb.Put(cKey, data); err != nil {
		logger.Crit("Failed to write compressed data", "schema", schema.name(), "err", err)
	}
}

func deleteCompressed(schema ItemSchema, from, to uint64) {
	cDb := schema.compressedDb()
	cKey := chunkKey(schema.compressedKeyPrefix(), from, to)
	if err := cDb.Delete(cKey); err != nil {
		logger.Crit("Failed to delete compressed data", "schema", schema.name(), "err", err)
	}
}

func readNextNum(schema ItemSchema) *uint64 {
	cDb := schema.compressedDb()
	nKey := schema.nextNumKey()
	bNum, err := cDb.Get(nKey)
	if err != nil || len(bNum) == 0 {
		return nil
	}
	if len(bNum) != 8 {
		logger.Crit("Malformed next compression number", "schema", schema.name(), "length", len(bNum))
	}
	num := binary.BigEndian.Uint64(bNum)
	return &num
}

func writeNextNum(schema ItemSchema, num uint64) {
	cDb := schema.compressedDb()
	cKey := schema.nextNumKey()
	bNum := make([]byte, 8)
	binary.BigEndian.PutUint64(bNum, num)
	if err := cDb.Put(cKey, bNum); err != nil {
		logger.Crit("Failed to write next compression number", "schema", schema.name(), "err", err)
	}
}

// HeaderSchema is wrapper for the HeaderRLP schema.
type HeaderSchema struct {
	uncompressedKv database.Database
	compressedKv   database.Database
}

func NewHeaderSchema(uncompressedKv, compressedKv database.Database) *HeaderSchema {
	return &HeaderSchema{
		uncompressedKv: uncompressedKv,
		compressedKv:   compressedKv,
	}
}

func (h *HeaderSchema) name() string                      { return schemaNameHeader }
func (h *HeaderSchema) uncompressedDb() database.Database { return h.uncompressedKv }
func (h *HeaderSchema) compressedDb() database.Database   { return h.compressedKv }
func (h *HeaderSchema) uncompressedKey(num uint64, hash common.Hash) []byte {
	return database.HeaderKey(num, hash)
}
func (h *HeaderSchema) compressedKeyPrefix() []byte { return compressedHeaderKeyPrefix }
func (h *HeaderSchema) nextNumKey() []byte          { return compressedHeaderNextNumKey }

// BodySchema is wrapper for the block body RLP schema.
type BodySchema struct {
	uncompressedKv database.Database
	compressedKv   database.Database
}

func NewBodySchema(uncompressedKv, compressedKv database.Database) *BodySchema {
	return &BodySchema{
		uncompressedKv: uncompressedKv,
		compressedKv:   compressedKv,
	}
}

func (b *BodySchema) name() string                      { return schemaNameBody }
func (b *BodySchema) uncompressedDb() database.Database { return b.uncompressedKv }
func (b *BodySchema) compressedDb() database.Database   { return b.compressedKv }
func (b *BodySchema) uncompressedKey(num uint64, hash common.Hash) []byte {
	return database.BlockBodyKey(num, hash)
}
func (b *BodySchema) compressedKeyPrefix() []byte { return compressedBodyKeyPrefix }
func (b *BodySchema) nextNumKey() []byte          { return compressedBodyNextNumKey }

// ReceiptSchema is wrapper for the block receipts RLP schema.
type ReceiptSchema struct {
	uncompressedKv database.Database
	compressedKv   database.Database
}

func NewReceiptSchema(uncompressedKv, compressedKv database.Database) *ReceiptSchema {
	return &ReceiptSchema{
		uncompressedKv: uncompressedKv,
		compressedKv:   compressedKv,
	}
}

func (r *ReceiptSchema) name() string                      { return schemaNameReceipts }
func (r *ReceiptSchema) uncompressedDb() database.Database { return r.uncompressedKv }
func (r *ReceiptSchema) compressedDb() database.Database   { return r.compressedKv }
func (r *ReceiptSchema) uncompressedKey(num uint64, hash common.Hash) []byte {
	return database.BlockReceiptsKey(num, hash)
}
func (r *ReceiptSchema) compressedKeyPrefix() []byte { return compressedReceiptsKeyPrefix }
func (r *ReceiptSchema) nextNumKey() []byte          { return compressedReceiptsNextNumKey }
