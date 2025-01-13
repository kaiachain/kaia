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
	"bytes"
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

func readSubsequentCompressionBlkNumber(dbm database.DBManager, compressTyp CompressionType) uint64 {
	var (
		prefix  = getCompressPrefix(compressTyp)
		db      = dbm.GetMiscDB()
		data, _ = db.Get(prefix)
	)
	if len(data) == 0 {
		return 0
	}
	var nextCompressionNumber uint64
	if err := rlp.Decode(bytes.NewReader(data), &nextCompressionNumber); err != nil {
		logger.Error("Invalid subsequent block number RLP", "err", err)
		return 0
	}
	return nextCompressionNumber
}

func writeSubsequentCompressionBlkNumber(dbm database.DBManager, compressTyp CompressionType, number uint64) {
	prefix := getCompressPrefix(compressTyp)
	data, err := rlp.EncodeToBytes(&number)
	if err != nil {
		logger.Crit("Failed to RLP encode block number of subsequent compression", "err", err)
	}
	db := dbm.GetMiscDB()
	if err := db.Put(prefix, data); err != nil {
		logger.Crit("Failed to store block number of subsequent compression number", "err", err)
	}
}

func sanitizeRange(dbm database.DBManager, compressTyp CompressionType, from, to, headNumber, compressChunk uint64) (uint64, uint64, error) {
	var err error
	if from == 0 {
		from = readSubsequentCompressionBlkNumber(dbm, compressTyp)
	}
	if to == 0 {
		nextRange := from + compressChunk
		if headNumber < nextRange {
			to = headNumber
		} else {
			to = nextRange
		}
	}
	// input validation
	if from > to {
		err = fmt.Errorf("[%s] from(%d) must be greater than to(%d)", compressTyp.String(), from, to)
	}
	return from, to, err
}

func removeOriginAfterCompress(dbm database.DBManager, compressions []any) {
	for _, compressed := range compressions {
		switch c := compressed.(type) {
		case *HeaderCompression:
			dbm.DeleteHeaderOnly(c.BlkHash, c.BlkNumber)
		case *BodyCompression:
			dbm.DeleteBody(c.BlkHash, c.BlkNumber)
		case *ReceiptCompression:
			dbm.DeleteReceipts(c.BlkHash, c.BlkNumber)
		}
	}
}

func findDataFromChunk(dbm database.DBManager, compressTyp CompressionType, finder Finder, number uint64, blkHash common.Hash, notFoundErr error) (any, error) {
	// Badger DB does not support `NewIterator()`
	if dbm.GetDBConfig().DBType == database.BadgerDB {
		return nil, nil
	}
	var (
		prefix       = getCompressPrefix(compressTyp)
		it           = getCompressDB(dbm, compressTyp).NewIterator(prefix, toBinary(number))
		decompressed any
		err          error
	)
	// 1. Find a chunk through range search
	defer it.Release()
	for it.Next() {
		from, to := parseCompressKey(compressTyp, it.Key())
		if from <= number && number <= to {
			decompressed, err = finder(dbm, from, to, number, blkHash)
			if err != nil {
				return nil, err
			}
			if decompressed != nil {
				return decompressed, nil
			} else {
				return nil, notFoundErr
			}
		}
	}
	return nil, notFoundErr
}

func decompress(dbm database.DBManager, compressTyp CompressionType, from, to uint64) (any, error) {
	var (
		compressKey = getCompressKey(compressTyp, from, to)
		db          = getCompressDB(dbm, compressTyp)
		data, _     = db.Get(compressKey)
	)
	decompressed, err := Decompress(data)
	if err != nil {
		return nil, err
	}
	switch compressTyp {
	case HeaderCompressType:
		var headerCompressions []*HeaderCompression
		if err := rlp.DecodeBytes(decompressed, &headerCompressions); err != nil {
			return nil, err
		}
		return headerCompressions, nil
	case BodyCompressType:
		var bodyCompressions []*BodyCompression
		if err := rlp.DecodeBytes(decompressed, &bodyCompressions); err != nil {
			return nil, err
		}
		return bodyCompressions, nil
	case ReceiptCompressType:
		var receiptCompressions []*ReceiptCompression
		if err := rlp.DecodeBytes(decompressed, &receiptCompressions); err != nil {
			return nil, err
		}
		return receiptCompressions, nil
	default:
		panic("unreachable")
	}
}

func compressStorage(dbm database.DBManager, compressTyp CompressionType, readData func(common.Hash, uint64) any, from, to, headNumber, compressChunk, maxSize uint64, migrationMode bool) (uint64, int, int, error) {
	from, to, err := sanitizeRange(dbm, compressTyp, from, to, headNumber, compressChunk)
	if err != nil {
		return 0, 0, 0, err
	}

	var (
		itIdx               = uint64(0)
		compressedTo        = from
		compressions        = make([]any, to-from)
		accumulatedByteSize = CompressionSize(0)
	)
	// migration loop
	for i := from; i < to; i++ {
		blkHash := dbm.ReadCanonicalHash(i)
		if common.EmptyHash(blkHash) {
			return 0, 0, 0, fmt.Errorf("[%s Compression] Block does not exist (number=%d)", compressTyp.String(), i)
		}
		data := readData(blkHash, i)
		switch compressTyp {
		case HeaderCompressType:
			header := data.(*types.Header)
			compressions[itIdx] = &HeaderCompression{
				BlkNumber: i,
				BlkHash:   blkHash,
				Header:    header,
			}
			rlp.Encode(&accumulatedByteSize, header)
		case BodyCompressType:
			body := &types.Body{}
			if data != nil {
				body = data.(*types.Body)
			}
			compressions[itIdx] = &BodyCompression{
				BlkNumber: i,
				BlkHash:   blkHash,
				Body:      body,
			}
			rlp.Encode(&accumulatedByteSize, body)
		case ReceiptCompressType:
			receipts := []*types.ReceiptForStorage{}
			if data != nil {
				receipts = data.([]*types.ReceiptForStorage)
			}
			compressions[itIdx] = &ReceiptCompression{
				BlkNumber:       i,
				BlkHash:         blkHash,
				StorageReceipts: receipts,
			}
			rlp.Encode(&accumulatedByteSize, receipts)
		}
		itIdx++
		compressedTo = uint64(i)
		if uint64(accumulatedByteSize) > maxSize {
			break
		}
	}
	if itIdx == 0 {
		// There is no data to compress
		return to, 0, 0, nil
	}
	bytes, err := rlp.EncodeToBytes(compressions[:itIdx])
	if err != nil {
		return 0, 0, 0, err
	}

	compressedSize := writeCompression(dbm, compressTyp, bytes, from, compressedTo)
	nextCompressStart := compressedTo + 1
	if migrationMode {
		writeSubsequentCompressionBlkNumber(dbm, compressTyp, nextCompressStart)
	}

	removeOriginAfterCompress(dbm, compressions[:itIdx])
	return nextCompressStart, len(bytes), compressedSize, nil
}

func writeCompression(dbm database.DBManager, compressTyp CompressionType, compressedBytes []byte, from, to uint64) int {
	var (
		compressKey = getCompressKey(compressTyp, from, to)
		db          = getCompressDB(dbm, compressTyp)
	)

	// Store compressed receipts (range `from` to `to`)
	compressedSize, compressedData := Compress(compressedBytes)
	if err := db.Put(compressKey, compressedData); err != nil {
		logger.Crit("Failed to store compressed receipts")
	}
	return compressedSize
}

func compressHeader(dbm database.DBManager, from, to, headNumber, compressChunk, maxSize uint64, migrationMode bool) (uint64, int, int, error) {
	readData := func(blkHash common.Hash, blkNumber uint64) any {
		return dbm.ReadHeader(blkHash, blkNumber)
	}
	return compressStorage(dbm, HeaderCompressType, readData, from, to, headNumber, compressChunk, maxSize, migrationMode)
}

func compressBody(dbm database.DBManager, from, to, headNumber, compressChunk, maxSize uint64, migrationMode bool) (uint64, int, int, error) {
	readData := func(blkHash common.Hash, blkNumber uint64) any {
		body := dbm.ReadBody(blkHash, blkNumber)
		if body == nil || len(body.Transactions) == 0 {
			return nil
		}
		return body
	}
	return compressStorage(dbm, BodyCompressType, readData, from, to, headNumber, compressChunk, maxSize, migrationMode)
}

func compressReceipts(dbm database.DBManager, from, to, headNumber, compressChunk, maxSize uint64, migrationMode bool) (uint64, int, int, error) {
	readData := func(blkHash common.Hash, blkNumber uint64) any {
		receipts := dbm.ReadReceipts(blkHash, blkNumber)
		if receipts == nil || len(receipts) == 0 {
			return nil
		}
		storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
		for number, receipt := range receipts {
			storageReceipts[number] = (*types.ReceiptForStorage)(receipt)
		}
		return storageReceipts
	}
	return compressStorage(dbm, ReceiptCompressType, readData, from, to, headNumber, compressChunk, maxSize, migrationMode)
}

func decompressHeader(dbm database.DBManager, from, to uint64) ([]*HeaderCompression, error) {
	decompressed, err := decompress(dbm, HeaderCompressType, from, to)
	if err != nil {
		return nil, err
	}
	return decompressed.([]*HeaderCompression), nil
}

func decompressBody(dbm database.DBManager, from, to uint64) ([]*BodyCompression, error) {
	decompressed, err := decompress(dbm, BodyCompressType, from, to)
	if err != nil {
		return nil, err
	}
	return decompressed.([]*BodyCompression), nil
}

func decompressReceipts(dbm database.DBManager, from, to uint64) ([]*ReceiptCompression, error) {
	decompressed, err := decompress(dbm, ReceiptCompressType, from, to)
	if err != nil {
		return nil, err
	}
	return decompressed.([]*ReceiptCompression), nil
}

func migrateHeaderChunkToOriginDB(dbm database.DBManager, compressions any) error {
	headerCompressions := compressions.([]*HeaderCompression)
	batch := dbm.NewBatch(dbm.GetHeaderDBEntryType())
	for _, c := range headerCompressions {
		if !dbm.HasHeaderInDisk(c.BlkHash, c.BlkNumber) {
			dbm.PutHeaderToBatch(batch, c.Header)
		}
	}
	_, err := database.WriteBatches(batch)
	return err
}

func migrateBodyChunkToOriginDB(dbm database.DBManager, compressions any) error {
	bodyCompressions := compressions.([]*BodyCompression)
	batch := dbm.NewBatch(dbm.GetBodyDBEntryType())
	for _, c := range bodyCompressions {
		if !dbm.HasBody(c.BlkHash, c.BlkNumber) {
			dbm.PutBodyToBatch(batch, c.BlkHash, c.BlkNumber, c.Body)
		}
	}
	_, err := database.WriteBatches(batch)
	return err
}

func migrateReceiptChunkToOriginDB(dbm database.DBManager, compressions any) error {
	receiptsCompressions := compressions.([]*ReceiptCompression)
	batch := dbm.NewBatch(dbm.GetReceiptsDBEntryType())
	for _, c := range receiptsCompressions {
		if !dbm.HasReceipts(c.BlkHash, c.BlkNumber) {
			receipts := make([]*types.Receipt, len(c.StorageReceipts))
			for i, r := range c.StorageReceipts {
				receipts[i] = (*types.Receipt)(r)
			}
			dbm.PutReceiptsToBatch(batch, c.BlkHash, c.BlkNumber, types.Receipts(receipts))
		}
	}
	_, err := database.WriteBatches(batch)
	return err
}

func deleteDataFromChunk(dbm database.DBManager, compressTyp CompressionType, number uint64, blkHash common.Hash) (uint64, bool, error) {
	// Badger DB does not support `NewIterator()`
	if dbm.GetDBConfig().DBType == database.BadgerDB {
		return 0, false, nil
	}

	var (
		prefix = getCompressPrefix(compressTyp)
		db     = getCompressDB(dbm, compressTyp)
	)
	// 1. Find a chunk through range search
	it := db.NewIterator(prefix, toBinary(number))
	defer it.Release()
	for it.Next() { // ascending order iteration
		from, to := parseCompressKey(compressTyp, it.Key())
		if from > number { // early exit if `from` is larger than target `number`
			return 0, false, nil
		}
		if from <= number && number <= to {
			// 1. Once a find chunk, decompress it
			compressions, err := decompress(dbm, compressTyp, from, to)
			if err != nil {
				logger.Crit("[Compression] Failed to decompress a chunk", "compressTyp", compressTyp.String(), "from", from, "to", to)
			}
			// 2. Store decompressed chunk to origin database
			switch compressTyp {
			case HeaderCompressType:
				err = migrateHeaderChunkToOriginDB(dbm, compressions)
			case BodyCompressType:
				err = migrateBodyChunkToOriginDB(dbm, compressions)
			case ReceiptCompressType:
				err = migrateReceiptChunkToOriginDB(dbm, compressions)
			}
			if err != nil {
				logger.Crit("[Compression] failed to reverse migrate", "compressTyp", compressTyp.String(), "from", from, "to", to)
			}
			// delete compression and return the starting number so that the compression moduel can start work from there
			if err := db.Delete(it.Key()); err != nil {
				logger.Crit(fmt.Sprintf("Failed to delete compressed data. err(%s) type(%s), from=%d, to=%d", err.Error(), compressTyp.String(), from, to))
			}
			logger.Info("[Compression] chunk deleted", "type", compressTyp.String(), "from", from, "to", to)
			return from, true, nil
		}
	}
	return 0, false, nil
}

func deleteHeaderFromChunk(dbm database.DBManager, number uint64, blkHash common.Hash) (uint64, bool, error) {
	return deleteDataFromChunk(dbm, HeaderCompressType, number, blkHash)
}

func deleteBodyFromChunk(dbm database.DBManager, number uint64, blkHash common.Hash) (uint64, bool, error) {
	return deleteDataFromChunk(dbm, BodyCompressType, number, blkHash)
}

func deleteReceiptsFromChunk(dbm database.DBManager, number uint64, blkHash common.Hash) (uint64, bool, error) {
	return deleteDataFromChunk(dbm, ReceiptCompressType, number, blkHash)
}

func compressedHeaderFinder(dbm database.DBManager, from, to, number uint64, blkHash common.Hash) (any, error) {
	var headerCompressions []*HeaderCompression
	if hc, ok := getFromCache(HeaderCompressType, from, to); ok {
		// Found from cache
		headerCompressions = hc.([]*HeaderCompression)
	} else {
		// Find a chunk and decompress
		if hc, err := decompressHeader(dbm, from, to); err == nil {
			headerCompressions = hc
			addCache(HeaderCompressType, from, to, headerCompressions)
		} else {
			return nil, err
		}
	}
	// Make a `types.Header` struct and returns it`
	for _, hc := range headerCompressions {
		if hc.BlkHash == blkHash {
			return hc.Header, nil
		}
	}
	return nil, nil
}

func compressedBodyFinder(dbm database.DBManager, from, to, number uint64, blkHash common.Hash) (any, error) {
	var bodyCompressions []*BodyCompression
	if bc, ok := getFromCache(BodyCompressType, from, to); ok {
		bodyCompressions = bc.([]*BodyCompression)
	} else {
		// Find a chunk and decompress
		if bc, err := decompressBody(dbm, from, to); err == nil {
			bodyCompressions = bc
			addCache(BodyCompressType, from, to, bodyCompressions)
		} else {
			return nil, err
		}
	}
	// Make a `types.Body` struct and returns it`
	for _, bc := range bodyCompressions {
		if bc.BlkHash == blkHash {
			return bc.Body, nil
		}
	}
	return nil, nil
}

func compressedReceiptFinder(dbm database.DBManager, from, to, number uint64, blkHash common.Hash) (any, error) {
	var receiptsCompressions []*ReceiptCompression
	if rc, ok := getFromCache(ReceiptCompressType, from, to); ok {
		receiptsCompressions = rc.([]*ReceiptCompression)
	} else {
		// Find a chunk and decompress
		if rc, err := decompressReceipts(dbm, from, to); err == nil {
			receiptsCompressions = rc
			addCache(ReceiptCompressType, from, to, receiptsCompressions)
		} else {
			return nil, err
		}
	}
	// Make a `types.Receipt` struct and returns it`
	for _, rc := range receiptsCompressions {
		if rc.BlkHash == blkHash {
			receipts := make(types.Receipts, len(rc.StorageReceipts))
			for idx, receipt := range rc.StorageReceipts {
				receipts[idx] = (*types.Receipt)(receipt)
			}
			return receipts, nil
		}
	}
	return nil, nil
}

func findHeaderFromChunkWithBlkHash(dbm database.DBManager, number uint64, blkHash common.Hash) (*types.Header, error) {
	notFoundErr := fmt.Errorf("[Header Compression] Failed to find a header (blkNumber= %d, blkHash=%s)", number, blkHash.String())
	decompressed, err := findDataFromChunk(dbm, HeaderCompressType, compressedHeaderFinder, number, blkHash, notFoundErr)
	if err != nil {
		return nil, err
	}
	if decompressed == nil {
		return nil, errors.New("[Header Compression] header not found")
	}
	return decompressed.(*types.Header), nil
}

func findBodyFromChunkWithBlkHash(dbm database.DBManager, number uint64, blkHash common.Hash) (*types.Body, error) {
	notFoundErr := fmt.Errorf("[Body Compression] Failed to find transactions (blkNumber= %d, blkHash=%s)", number, blkHash.String())
	decompressed, err := findDataFromChunk(dbm, BodyCompressType, compressedBodyFinder, number, blkHash, notFoundErr)
	if err != nil {
		return nil, err
	}
	if decompressed == nil {
		return nil, errors.New("[Body Compression] body not found")
	}
	return decompressed.(*types.Body), nil
}

func findReceiptsFromChunkWithBlkHash(dbm database.DBManager, number uint64, blkHash common.Hash) (types.Receipts, error) {
	notFoundErr := fmt.Errorf("[Receipt Compression] Failed to find receipts (blkNumber= %d, blkHash=%s)", number, blkHash.String())
	decompressed, err := findDataFromChunk(dbm, ReceiptCompressType, compressedReceiptFinder, number, blkHash, notFoundErr)
	if err != nil {
		return nil, err
	}
	if decompressed == nil {
		return nil, errors.New("[Receipt Compression] receipt not found")
	}
	return decompressed.(types.Receipts), nil
}
