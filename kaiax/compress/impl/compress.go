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
	"math/rand"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

const SEC_TEN = time.Second * 10

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

func (c *CompressModule) idle(compressTyp CompressionType, nBlocks uint64) bool {
	idealIdleTime := time.Second * time.Duration(nBlocks)
	c.setIdleState(compressTyp, &IdleState{true, idealIdleTime})
	logger.Info("[Compression] Enter idle state", "type", compressTyp.String(), "idle", SEC_TEN, "ideal idle time", idealIdleTime)
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
	totalChunks := 0
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
			if c.idle(compressTyp, c.getCompressChunk()-residualBlkCnt) {
				return
			}
			continue
		}
		if c.getCompressRetention() > nextCompressionDistance {
			if c.idle(compressTyp, c.getCompressRetention()-nextCompressionDistance) {
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
				logger.Info("[Compression] chunk compressed", "type", compressTyp.String(), "from", originFrom, "subsequentBlkNumber", subsequentBlkNumber, "curBlknum", curBlkNum, "originSize", common.StorageSize(originSize), "compressedSize", common.StorageSize(compressedSize), "totalChunks", totalChunks)
				totalChunks++
			}
			if subsequentBlkNumber >= curBlkNum {
				break
			}
			from = subsequentBlkNumber
		}
	}
}

func (c *CompressModule) deleteHeader(hash common.Hash, num uint64) uint64 {
	newFromAfterRewind, err := deleteHeaderFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Header Compression] Failed to delete header", "blockNum", num, "blockHash", hash.String())
	}
	return newFromAfterRewind
}

func (c *CompressModule) deleteBody(hash common.Hash, num uint64) uint64 {
	newFromAfterRewind, err := deleteBodyFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Body Compression] Failed to delete body", "blockNum", num, "blockHash", hash.String())
	}
	return newFromAfterRewind
}

func (c *CompressModule) deleteReceipts(hash common.Hash, num uint64) uint64 {
	newFromAfterRewind, err := deleteReceiptsFromChunk(c.Dbm, num, hash)
	if err != nil {
		logger.Warn("[Receipt Compression] Failed to delete receipt", "blockNum", num, "blockHash", hash.String())
	}
	return newFromAfterRewind
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

func (c *CompressModule) restoreFragmentByRewind() {
	for _, compressTyp := range allCompressTypes {
		var (
			lastCompressDeleteKeyPrefix, lastCompressDeleteValuePrefix = getLsatCompressionDeleteKeyPrefix(compressTyp), getLsatCompressionDeleteValuePrefix(compressTyp)
			miscDB                                                     = c.Dbm.GetMiscDB()
		)
		key, err := miscDB.Get(lastCompressDeleteKeyPrefix)
		if err != nil {
			return
		}
		value, err := miscDB.Get(lastCompressDeleteValuePrefix)
		if err != nil {
			return
		}
		if bytes.Equal(key, lastCompressionCleared) && bytes.Equal(value, lastCompressionCleared) {
			// No reserved last compression data. Noting to do
			return
		}

		var (
			db       = getCompressDB(c.Dbm, compressTyp)
			from, to = parseCompressKey(compressTyp, key)
		)
		// 1. restore last compression data
		if err := db.Put(key, value); err != nil {
			logger.Crit(fmt.Sprintf("Failed to restore last compressed data. err(%s) type(%s), from=%d, to=%d", err.Error(), compressTyp.String(), from, to))
		}
		// 2. clear last compression data
		if err := miscDB.Put(lastCompressDeleteKeyPrefix, lastCompressionCleared); err != nil {
			logger.Crit("Failed to clear last decompression key")
		}
		if err := miscDB.Put(lastCompressDeleteValuePrefix, lastCompressionCleared); err != nil {
			logger.Crit("Failed to clear last decompression value")
		}
	}
}

// TODO-hyunsooda: Move to `compress_test.go`
func (c *CompressModule) testCopyOriginData(compressTyp CompressionType, copyTestDB database.Batch, from, to uint64) {
	// Copy origin
	for i := from; i < to; i++ {
		hash := c.Dbm.ReadCanonicalHash(i)
		switch compressTyp {
		case HeaderCompressType:
			originHeader := c.Dbm.ReadHeader(hash, i)
			c.Dbm.PutHeaderToBatch(copyTestDB, hash, i, originHeader)
		case BodyCompressType:
			originBody := c.Dbm.ReadBody(hash, i)
			c.Dbm.PutBodyToBatch(copyTestDB, hash, i, originBody)
		case ReceiptCompressType:
			originReceipts := c.Dbm.ReadReceipts(hash, i)
			c.Dbm.PutReceiptsToBatch(copyTestDB, hash, i, originReceipts)
		}
	}
}

func (c *CompressModule) testVerifyCompressionIntegrity(compressTyp CompressionType, from, to uint64) error {
	for i := from; i < to; i++ {
		blkHash := c.Dbm.ReadCanonicalHash(i)
		switch compressTyp {
		case HeaderCompressType:
			originHeader := c.Dbm.ReadHeader(blkHash, i)
			compressedHeader, err := findHeaderFromChunkWithBlkHash(c.Dbm, i, originHeader.Hash())
			if err != nil {
				return err
			}
			if originHeader.Hash() != compressedHeader.Hash() {
				return fmt.Errorf("[Header Compression Test] Compressed header is not the same data with origin header data (number=%d, txHash=%s)", i, originHeader.Hash().String())
			}
		case BodyCompressType:
			originBody := c.Dbm.ReadBody(blkHash, i)
			if len(originBody.Transactions) > 0 {
				compressedBody, err := findBodyFromChunkWithBlkHash(c.Dbm, i, blkHash)
				if err != nil {
					return err
				}
				for idx, originTx := range originBody.Transactions {
					if !originTx.Equal(compressedBody.Transactions[idx]) {
						return fmt.Errorf("[Body Compression Test] Compressed body is not the same data with origin body data (number=%d, txHash=%s)", i, originTx.Hash().String())
					}
				}
			}
		case ReceiptCompressType:
			originReceipts := c.Dbm.ReadReceipts(blkHash, i)
			if len(originReceipts) > 0 {
				compressedReceipts, err := findReceiptsFromChunkWithBlkHash(c.Dbm, i, blkHash)
				if err != nil {
					return err
				}
				for idx, originReceipt := range originReceipts {
					if !reflect.DeepEqual(originReceipt, compressedReceipts[idx]) {
						return fmt.Errorf("[Receipt Compression Test] Compressed receipt is not the same data with origin receipt data (number=%d, txHash=%s)", i, originReceipt.TxHash.String())
					}
				}
			}
		}
	}
	return nil
}

func (c *CompressModule) testFinder(compressTyp CompressionType, from, to uint64) {
	for i := from; i < to; i++ {
		blkHash := c.Dbm.ReadCanonicalHash(i)
		switch compressTyp {
		case HeaderCompressType:
			findHeaderFromChunkWithBlkHash(c.Dbm, i, blkHash)
		case BodyCompressType:
			findBodyFromChunkWithBlkHash(c.Dbm, i, blkHash)
		case ReceiptCompressType:
			findReceiptsFromChunkWithBlkHash(c.Dbm, i, blkHash)
		}
	}
}

func (c *CompressModule) testCompress(t *testing.T, compressTyp CompressionType, from, to uint64, tempDir string, wg *sync.WaitGroup) {
	dbConfig := c.Dbm.GetDBConfig()
	copyTestDB, err := database.TestCreateNewDB(dbConfig, filepath.Join(dbConfig.Dir, tempDir))
	assert.Nil(t, err)
	defer func() {
		copyTestDB.Release()
		wg.Done()
	}()
	writeSubsequentCompressionBlkNumber(c.Dbm, compressTyp, from)
	curBlkNum := c.Chain.CurrentBlock().NumberU64()

	var (
		totalHeaderCompressedElapsedTime   time.Duration
		totalBodyCompressedElapsedTime     time.Duration
		totalReceiptsCompressedElapsedTime time.Duration
		nCompressedHeader                  = uint64(0)
		nCompressedBody                    = uint64(0)
		nCompressedReceipts                = uint64(0)
		nHeaderChunks                      = 0
		nBodyChunkcs                       = 0
		nReceiptsChunkcs                   = 0
		loopCnt                            = 0

		originFrom = from
	)

	for {
		var (
			subsequentBlkNumber uint64
			err                 error
			originSize          int
			compressedSize      int
		)
		startCompress := time.Now()
		switch compressTyp {
		case HeaderCompressType:
			subsequentBlkNumber, originSize, compressedSize, err = compressHeader(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
			totalHeaderCompressedElapsedTime += time.Since(startCompress)
			nCompressedHeader = subsequentBlkNumber - 1
			nHeaderChunks++
		case BodyCompressType:
			subsequentBlkNumber, originSize, compressedSize, err = compressBody(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
			totalBodyCompressedElapsedTime += time.Since(startCompress)
			nCompressedBody = subsequentBlkNumber - 1
			nBodyChunkcs++
		case ReceiptCompressType:
			subsequentBlkNumber, originSize, compressedSize, err = compressReceipts(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
			totalReceiptsCompressedElapsedTime += time.Since(startCompress)
			nCompressedReceipts = subsequentBlkNumber - 1
			nReceiptsChunkcs++
		}
		assert.Nil(t, err)
		// Copy origin receipts
		c.testCopyOriginData(compressTyp, copyTestDB, from, subsequentBlkNumber-1)
		_, err = database.WriteBatches(copyTestDB)
		assert.Nil(t, err)
		// Compression integrity test
		err = c.testVerifyCompressionIntegrity(compressTyp, from, subsequentBlkNumber-1)
		assert.Nil(t, err)
		if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
			logger.Info("[Compression] chunk compressed", "type", compressTyp.String(), "from", originFrom, "subsequentBlkNumber", subsequentBlkNumber, "curBlknum", curBlkNum, "originSize", common.StorageSize(originSize), "compressedSize", common.StorageSize(compressedSize))
			break
		}
		from = subsequentBlkNumber
		loopCnt++
		if loopCnt%100 == 0 {
			fmt.Printf("%s ... remaining: %d\n", compressTyp.String(), to-subsequentBlkNumber)
		}
	}
	switch compressTyp {
	case HeaderCompressType:
		fmt.Printf("number of compressed header    (from=%d) (to=%d) (# of blocks = %d) (# of chunks = %d) \n", originFrom, to, nCompressedHeader-originFrom, nHeaderChunks)
		fmt.Println("total header   compression    elapsed time: ", totalHeaderCompressedElapsedTime)
		fmt.Println("avg   header   compression    elapsed time: ", totalHeaderCompressedElapsedTime/time.Duration(nHeaderChunks))
	case BodyCompressType:
		fmt.Printf("number of compressed body      (from=%d) (to=%d) (# of blocks = %d) (# of chunks = %d) \n", originFrom, to, nCompressedBody-originFrom, nBodyChunkcs)
		fmt.Println("total body     compression    elapsed time: ", totalBodyCompressedElapsedTime)
		fmt.Println("avg   body     compression    elapsed time: ", totalBodyCompressedElapsedTime/time.Duration(nBodyChunkcs))
	case ReceiptCompressType:
		fmt.Printf("number of compressed receipts  (from=%d) (to=%d) (# of blocks = %d) (# of chunks = %d) \n", originFrom, to, nCompressedReceipts-originFrom, nReceiptsChunkcs)
		fmt.Println("total receipts compression    elapsed time: ", totalReceiptsCompressedElapsedTime)
		fmt.Println("avg   receipts compression    elapsed time: ", totalReceiptsCompressedElapsedTime/time.Duration(nReceiptsChunkcs))
	}
}

func testInit(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) (*CompressModule, database.DBManager) {
	var (
		bc, chainDB, setupErr = setup(t)
		initOpts              = &InitOpts{
			ChunkBlockSize: blockchain.DefaultChunkBlockSize,
			ChunkCap:       blockchain.DefaultCompressChunkCap,
			Chain:          bc,
			Dbm:            chainDB,
		}
		mCompress     = NewCompression()
		moduleInitErr = mCompress.Init(initOpts)
	)

	if err := errors.Join(setupErr, moduleInitErr); err != nil {
		if errors.Is(err, compress.ErrInitNil) {
			// If no environment varaible set, do not execute compression test
			return nil, nil
		} else {
			t.Fatal(err)
		}
	}
	chainDB.SetCompressModule(mCompress)
	return mCompress, chainDB
}

func TestCompressFunction(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) {
	mCompress, chainDB := testInit(t, setup)
	if mCompress == nil || chainDB == nil {
		return
	}
	defer chainDB.Close()

	from, to := uint64(0), uint64(5000)
	var wg sync.WaitGroup
	wg.Add(3)
	go mCompress.testCompress(t, HeaderCompressType, from, to, "copy_header", &wg)
	go mCompress.testCompress(t, BodyCompressType, from, to, "copy_body", &wg)
	go mCompress.testCompress(t, ReceiptCompressType, from, to, "copy_receipts", &wg)
	wg.Wait()
}

func TestCompressFinder(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) {
	mCompress, chainDB := testInit(t, setup)
	if mCompress == nil || chainDB == nil {
		return
	}
	defer chainDB.Close()

	rand.Seed(time.Now().UnixNano())

	var (
		headerElapsed         time.Duration
		bodyElapsed           time.Duration
		receiptsElapsed       time.Duration
		originHeaderElapsed   time.Duration
		originBodyElapsed     time.Duration
		originReceiptsElapsed time.Duration
		n                     = 100
	)
	for range n {
		for _, compressTyp := range allCompressTypes {
			var (
				r    = rand.Uint64() % (5000)
				hash = mCompress.Dbm.ReadCanonicalHash(r)
			)
			startOrigin := time.Now()
			switch compressTyp {
			case HeaderCompressType:
				header := mCompress.Dbm.ReadHeader(hash, r)
				assert.NotNil(t, header)
				originHeaderElapsed = time.Since(startOrigin)
			case BodyCompressType:
				body := mCompress.Dbm.ReadBody(hash, r)
				assert.NotNil(t, body)
				originBodyElapsed = time.Since(startOrigin)
			case ReceiptCompressType:
				receipts := mCompress.Dbm.ReadReceipts(hash, r)
				assert.NotNil(t, receipts)
				originReceiptsElapsed = time.Since(startOrigin)
			}

			startFinder := time.Now()
			switch compressTyp {
			case HeaderCompressType:
				header, _ := findHeaderFromChunkWithBlkHash(mCompress.Dbm, r, hash)
				assert.NotNil(t, header)
				headerElapsed = time.Since(startFinder)
			case BodyCompressType:
				body, _ := findBodyFromChunkWithBlkHash(mCompress.Dbm, r, hash)
				assert.NotNil(t, body)
				bodyElapsed = time.Since(startFinder)
			case ReceiptCompressType:
				receipts, _ := findReceiptsFromChunkWithBlkHash(mCompress.Dbm, r, hash)
				assert.NotNil(t, receipts)
				receiptsElapsed = time.Since(startFinder)
			}
		}
	}

	fmt.Println("avg header elapsed:", headerElapsed/time.Duration(n))
	fmt.Println("avg body elapsed:", bodyElapsed/time.Duration(n))
	fmt.Println("avg receipts elapsed:", receiptsElapsed/time.Duration(n))
	fmt.Println("avg origin header elapsed:", originHeaderElapsed/time.Duration(n))
	fmt.Println("avg origin body elapsed:", originBodyElapsed/time.Duration(n))
	fmt.Println("avg origin receipts elapsed:", originReceiptsElapsed/time.Duration(n))
}
