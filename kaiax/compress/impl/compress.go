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
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

const (
	SEC_TEN           = time.Second * 10
	COMPACTION_PERIOD = 86400 * 14 // 2 weeks
)

func (c *CompressModule) stopCompress() {
	for range allCompressTypes {
		c.terminateCompress <- struct{}{}
	}
	for _, compressTyp := range allCompressTypes {
		c.wg.Wait()
		c.setIdleState(compressTyp, nil)
	}
}

func (c *CompressModule) Compress() {
	c.wg.Add(TotalCompressTypeSize)
	go c.compressHeader()
	go c.compressBody()
	go c.compressReceipts()
}

func (c *CompressModule) compressHeader() {
	c.compress(HeaderCompressType, compressHeader)
}

func (c *CompressModule) compressBody() {
	c.compress(BodyCompressType, compressBody)
}

func (c *CompressModule) compressReceipts() {
	c.compress(ReceiptCompressType, compressReceipts)
}

func (c *CompressModule) idle(compressTyp CompressionType, nBlocks, curBlkNum uint64) bool {
	idealIdleTime := time.Second * time.Duration(nBlocks)
	c.setIdleState(compressTyp, &IdleState{true, idealIdleTime})
	logger.Info("[Compression] Enter idle state", "type", compressTyp.String(), "idle", SEC_TEN, "ideal idle time", idealIdleTime, "curBlkNum", curBlkNum, "watingBlocks", nBlocks)
	for {
		timer := time.NewTimer(SEC_TEN)
		select {
		case <-c.terminateCompress:
			logger.Info("[Compression] Stop signal received", "type", compressTyp.String())
			return true
		case <-timer.C:
			return false
		}
	}
}

func (c *CompressModule) compress(compressTyp CompressionType, compressFn CompressFn) {
	defer c.wg.Done()
	log := logger.Trace
	if c.ForceLogging {
		log = logger.Info
	}
	var (
		totalChunks    = 0
		compactBaseNum = c.Chain.CurrentBlock().NumberU64()
	)
	for {
		select {
		case <-c.terminateCompress:
			logger.Info("[Compression] Stop signal received", "type", compressTyp.String())
			return
		default:
		}

		var (
			curBlkNum               = c.Chain.CurrentBlock().NumberU64()
			residualBlkCnt          = curBlkNum % c.getCompressChunk()
			nextCompressionBlkNum   = readSubsequentCompressionBlkNumber(c.Dbm, compressTyp)
			nextCompressionDistance = curBlkNum - nextCompressionBlkNum
			// Do not wait if next compression block number is far awway. Start migration right now
			noWait     = curBlkNum > nextCompressionBlkNum && nextCompressionDistance > c.getCompressChunk()
			originFrom = readSubsequentCompressionBlkNumber(c.Dbm, compressTyp)
			from       = originFrom
		)
		// 1. Idle check
		if curBlkNum < c.getCompressChunk() || (residualBlkCnt != 0 && !noWait) {
			if c.idle(compressTyp, c.getCompressChunk()-residualBlkCnt, curBlkNum) {
				return
			}
			continue
		}
		if nextCompressionDistance == 0 || c.getCompressRetention() > nextCompressionDistance {
			if c.idle(compressTyp, c.getCompressRetention()-nextCompressionDistance, curBlkNum) {
				return
			}
			continue
		}
		c.setIdleState(compressTyp, nil)
		// 2. Main loop (compression)
		logger.Info("[Compression] Start compression loop", "type", compressTyp.String(), "from", originFrom, "curBlknum", curBlkNum, "totalChunks", totalChunks)
		for {
			select {
			case <-c.terminateCompress:
				logger.Info("[Compression] Stop signal received", "type", compressTyp.String())
				return
			default:
			}
			subsequentBlkNumber, originSize, compressedSize, err := compressFn(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
			if err != nil {
				logger.Warn("[Compression] failed to compress chunk", "type", compressTyp.String(), "err", err)
				break
			}
			if compressedSize != 0 {
				log("[Compression] chunk compressed", "type", compressTyp.String(), "originFrom", originFrom, "startFrom", from, "subsequentBlkNumber", subsequentBlkNumber, "curBlknum", curBlkNum, "originSize", common.StorageSize(originSize), "compressedSize", common.StorageSize(compressedSize), "totalChunks", totalChunks)
				totalChunks++
			}
			if subsequentBlkNumber >= curBlkNum {
				break
			}
			from = subsequentBlkNumber
			if curBlkNum-compactBaseNum > COMPACTION_PERIOD {
				go c.compact(false)
				compactBaseNum = curBlkNum
			}
			time.Sleep(c.loopIdleTime) // unconditional 50ms idle
		}
	}
}

func (c *CompressModule) deleteHeader(hash common.Hash, num uint64) {
	newFromAfterRewind, shouldUpdate, err := deleteHeaderFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Header Compression] Failed to delete header", "blockNum", num, "blockHash", hash.String())
	}
	c.headerChunkDeleted <- ChunkDelete{shouldUpdate, newFromAfterRewind}
}

func (c *CompressModule) deleteBody(hash common.Hash, num uint64) {
	newFromAfterRewind, shouldUpdate, err := deleteBodyFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Body Compression] Failed to delete body", "blockNum", num, "blockHash", hash.String())
	}
	c.bodyChunkDeleted <- ChunkDelete{shouldUpdate, newFromAfterRewind}
}

func (c *CompressModule) deleteReceipts(hash common.Hash, num uint64) {
	newFromAfterRewind, shouldUpdate, err := deleteReceiptsFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Receipt Compression] Failed to delete receipt", "blockNum", num, "blockHash", hash.String())
	}
	c.receiptsChunkDeleted <- ChunkDelete{shouldUpdate, newFromAfterRewind}
}

func (c *CompressModule) FindHeaderFromChunkWithBlkHash(dbm database.DBManager, number uint64, hash common.Hash) (*types.Header, error) {
	return findHeaderFromChunkWithBlkHash(dbm, number, hash)
}

func (c *CompressModule) FindBodyFromChunkWithBlkHash(dbm database.DBManager, number uint64, hash common.Hash) (*types.Body, error) {
	return findBodyFromChunkWithBlkHash(dbm, number, hash)
}

func (c *CompressModule) FindReceiptsFromChunkWithBlkHash(dbm database.DBManager, number uint64, hash common.Hash) (types.Receipts, error) {
	return findReceiptsFromChunkWithBlkHash(dbm, number, hash)
}

func (c *CompressModule) compact(compactCompression bool) {
	// if header, body, and receipts DB does not exist (i.e., unit test)
	if c.Dbm.GetDBConfig().DBType == database.BadgerDB {
		return
	}
	if atomic.LoadInt32(&c.isCompacting) != 0 {
		return
	}
	atomic.StoreInt32(&c.isCompacting, 1)
	defer atomic.StoreInt32(&c.isCompacting, 0)

	var (
		headerDB    = c.Dbm.GetHeaderDB()
		bodyDB      = c.Dbm.GetBodyDB()
		receiptsDB  = c.Dbm.GetReceiptsDB()
		cHeaderDB   = c.Dbm.GetCompressHeaderDB()
		cBodyDB     = c.Dbm.GetCompressBodyDB()
		cReceiptsDB = c.Dbm.GetCompressReceiptsDB()
		dbs         = []struct {
			db        database.Database
			str       string
			doCompact bool
		}{
			{headerDB, "headerDB", true},
			{bodyDB, "bodyDB", true},
			{receiptsDB, "receiptsDB", true},
			{cHeaderDB, "compressionHeaderDB", compactCompression},
			{cBodyDB, "compressionBodyDB", compactCompression},
			{cReceiptsDB, "compressionReceiptsDB", compactCompression},
		}
	)
	for _, db := range dbs {
		logger.Info("[Compression] start database compaction", "type", db.str)
		if db.db != nil && db.doCompact {
			err := db.db.Compact(nil, nil)
			logger.Info("[Compression] end database compaction", "type", db.str, "err", err)
		}
	}
}
