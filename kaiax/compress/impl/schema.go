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
	"encoding/binary"
	"io"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/klauspost/compress/zstd"
)

type CompressionType uint8

const (
	HeaderCompressType = iota
	BodyCompressType
	ReceiptCompressType
	TotalCompressTypeSize
)

var allCompressTypes = [TotalCompressTypeSize]CompressionType{
	HeaderCompressType,
	BodyCompressType,
	ReceiptCompressType,
}

var (
	compressHeaderPrefix  = []byte("CompressdHeader-")
	compressReceiptPrefix = []byte("CompressdReceipt-")
	compressBodyPrefix    = []byte("CompressdBody-")

	lastHeaderCompressionDeleteKey     = []byte("LastHeaderCompressionDeleteKey")
	lastHeaderCompressionDeleteValue   = []byte("LastHeaderCompressionDeleteValue")
	lastBodyCompressionDeleteKey       = []byte("LastBodyCompressionDeleteKey")
	lastBodyCompressionDeleteValue     = []byte("LastBodyCompressionDeleteValue")
	lastReceiptsCompressionDeleteKey   = []byte("LastReceiptsCompressionDeleteKey")
	lastReceiptsCompressionDeleteValue = []byte("LastReceiptsCompressionDeleteValue")
	lastCompressionCleared             = []byte("NoLastCompressionkeyReserved")

	// Create a writer that caches compressors.
	// For this operation type we supply a nil Reader.
	encoder, _ = zstd.NewWriter(nil)
	decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
)

func (typ CompressionType) String() string {
	switch typ {
	case HeaderCompressType:
		return "Header Compression"
	case BodyCompressType:
		return "Body Compression"
	case ReceiptCompressType:
		return "Receipts Compression"
	default:
		return ""
	}
}

func getLsatCompressionDeleteKeyPrefix(typ CompressionType) []byte {
	switch typ {
	case HeaderCompressType:
		return lastHeaderCompressionDeleteKey
	case BodyCompressType:
		return lastBodyCompressionDeleteKey
	case ReceiptCompressType:
		return lastReceiptsCompressionDeleteKey
	default:
		panic("unreacahble")
	}
}

func getLsatCompressionDeleteValuePrefix(typ CompressionType) []byte {
	switch typ {
	case HeaderCompressType:
		return lastHeaderCompressionDeleteValue
	case BodyCompressType:
		return lastBodyCompressionDeleteValue
	case ReceiptCompressType:
		return lastReceiptsCompressionDeleteValue
	default:
		panic("unreacahble")
	}
}

// Compressed range is represented as `to-from`
func getCompressKey(typ CompressionType, from, to uint64) []byte {
	bFrom, bTo := make([]byte, 8), make([]byte, 8)
	binary.BigEndian.PutUint64(bFrom, from)
	binary.BigEndian.PutUint64(bTo, to)

	var prefix []byte
	switch typ {
	case HeaderCompressType:
		prefix = compressHeaderPrefix
	case BodyCompressType:
		prefix = compressBodyPrefix
	case ReceiptCompressType:
		prefix = compressReceiptPrefix
	}
	return append(append(prefix, bTo...), bFrom...)
}

// Returned compressed range is represented as `from-to`
func parseCompressKey(typ CompressionType, key []byte) (uint64, uint64) {
	var prefixLen int
	switch typ {
	case HeaderCompressType:
		prefixLen = len(compressHeaderPrefix)
	case BodyCompressType:
		prefixLen = len(compressBodyPrefix)
	case ReceiptCompressType:
		prefixLen = len(compressReceiptPrefix)
	}

	to := binary.BigEndian.Uint64(key[prefixLen : prefixLen+8])
	from := binary.BigEndian.Uint64(key[prefixLen+8:])
	return from, to
}

func toBinary(number uint64) []byte {
	bstart := make([]byte, 8)
	binary.BigEndian.PutUint64(bstart, number)
	return bstart
}

func getCompressDB(dbm database.DBManager, compressTyp CompressionType) database.Database {
	if dbm.GetDBConfig().DBType == database.MemoryDB {
		return dbm.GetMemDB()
	}
	switch compressTyp {
	case HeaderCompressType:
		return dbm.GetCompressHeaderDB()
	case BodyCompressType:
		return dbm.GetCompressBodyDB()
	case ReceiptCompressType:
		return dbm.GetCompressReceiptsDB()
	default:
		panic("unreachable")
	}
}

func getDBType(compressTyp CompressionType) database.DBEntryType {
	switch compressTyp {
	case HeaderCompressType:
		return database.CompressHeaderDB
	case BodyCompressType:
		return database.CompressBodyDB
	case ReceiptCompressType:
		return database.CompressReceiptsDB
	default:
		panic("unreachable")
	}
}

// getDBType returns db type,  and compress key
func getDBTypeWithCompressKey(compressTyp CompressionType, from, to uint64) (database.DBEntryType, []byte) {
	switch compressTyp {
	case HeaderCompressType:
		return database.CompressHeaderDB, getCompressKey(compressTyp, from, to)
	case BodyCompressType:
		return database.CompressBodyDB, getCompressKey(compressTyp, from, to)
	case ReceiptCompressType:
		return database.CompressReceiptsDB, getCompressKey(compressTyp, from, to)
	default:
		panic("unreachable")
	}
}

func getCompressPrefix(compressTyp CompressionType) []byte {
	switch compressTyp {
	case HeaderCompressType:
		return compressHeaderPrefix
	case BodyCompressType:
		return compressBodyPrefix
	case ReceiptCompressType:
		return compressReceiptPrefix
	default:
		panic("unreachable")
	}
}

type CompressStructTyp interface {
	GetBlkHash() common.Hash
}

type (
	Finder       func(dbm database.DBManager, from, to, headNumber uint64, blkOrTxHash common.Hash) (any, error)
	CompressFn   func(dbm database.DBManager, from, to, headNumber, compressChunk, maxSize uint64, migrationMode bool) (uint64, int, int, error)
	DecompressFn func(dbm database.DBManager, compressTyp CompressionType, from, to uint64) ([]CompressStructTyp, error)
)

// Compress a buffer.
// If you have a destination buffer, the allocation in the call can also be eliminated.
func Compress(src []byte) []byte {
	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}

type ReceiptCompression struct {
	BlkNumber       uint64
	BlkHash         common.Hash
	StorageReceipts []*types.ReceiptForStorage
}

type ReceiptCompressionRLP struct {
	BlkNumber       uint64
	BlkHash         common.Hash
	StorageReceipts []*types.ReceiptForStorage
}

func (r *ReceiptCompression) GetBlkHash() common.Hash {
	return r.BlkHash
}

func (r *ReceiptCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &ReceiptCompressionRLP{BlkNumber: r.BlkNumber, BlkHash: r.BlkHash, StorageReceipts: r.StorageReceipts})
}

func (r *ReceiptCompression) DecodeRLP(s *rlp.Stream) error {
	var dec ReceiptCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.BlkNumber = dec.BlkNumber
	r.BlkHash = dec.BlkHash
	r.StorageReceipts = dec.StorageReceipts
	return nil
}

type BodyCompression struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Body      *types.Body
}

type BodyCompressionRLP struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Body      *types.Body
}

func (b *BodyCompression) GetBlkHash() common.Hash {
	return b.BlkHash
}

func (b *BodyCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &BodyCompressionRLP{BlkNumber: b.BlkNumber, BlkHash: b.BlkHash, Body: b.Body})
}

func (b *BodyCompression) DecodeRLP(s *rlp.Stream) error {
	var dec BodyCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	b.BlkNumber = dec.BlkNumber
	b.BlkHash = dec.BlkHash
	b.Body = dec.Body
	return nil
}

type HeaderCompression struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Header    *types.Header
}

type HeaderCompressionRLP struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Header    *types.Header
}

func (h *HeaderCompression) GetBlkNumber() uint64 {
	return h.BlkNumber
}

func (h *HeaderCompression) GetBlkHash() common.Hash {
	return h.BlkHash
}

func (h *HeaderCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &HeaderCompressionRLP{BlkNumber: h.BlkNumber, BlkHash: h.BlkHash, Header: h.Header})
}

func (h *HeaderCompression) DecodeRLP(s *rlp.Stream) error {
	var dec HeaderCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	h.BlkNumber = dec.BlkNumber
	h.BlkHash = dec.BlkHash
	h.Header = dec.Header
	return nil
}

type CompressionSize common.StorageSize

func (c *CompressionSize) Write(b []byte) (int, error) {
	*c += CompressionSize(len(b))
	return len(b), nil
}
