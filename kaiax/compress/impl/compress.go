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
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func (c *CompressModule) Compress() {
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

func (c *CompressModule) compress(compressTyp CompressionType, compressFn CompressFn) {
	for {
		var (
			curBlkNum             = c.Chain.CurrentBlock().NumberU64()
			residualBlkCnt        = curBlkNum % c.getCompressChunk()
			nextCompressionBlkNum = readSubsequentCompressionBlkNumber(c.Dbm, compressTyp)
			// Do not wait if next compression block number is far awway. Start migration right now
			noWait     = curBlkNum > nextCompressionBlkNum && curBlkNum-nextCompressionBlkNum > c.getCompressChunk()
			originFrom = readSubsequentCompressionBlkNumber(c.Dbm, compressTyp)
			from       = originFrom
		)
		if residualBlkCnt != 0 && !noWait {
			idleTime := time.Second * time.Duration(c.getCompressChunk()-residualBlkCnt)
			c.setIdleState(compressTyp, &IdleState{true, idleTime})
			time.Sleep(idleTime)
			continue
		}
		c.setIdleState(compressTyp, nil)
		for {
			subsequentBlkNumber, compressedSize, err := compressFn(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
			if err != nil {
				logger.Warn("[Compression] failed to compress chunk", "type", compressTyp.String(), "err", err)
				break
			}
			if subsequentBlkNumber >= curBlkNum {
				logger.Info("[Compression] chunk compressed", "type", compressTyp.String(), "from", originFrom, "subsequentBlkNumber", subsequentBlkNumber, "curBlknum", curBlkNum, "compressedSize", common.StorageSize(compressedSize))
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

func (c *CompressModule) testCompress(compressTyp CompressionType, from, to uint64, tempDir string) error {
	dbConfig := c.Dbm.GetDBConfig()
	copyTestDB, err := database.TestCreateNewDB(dbConfig, filepath.Join(dbConfig.Dir, tempDir))
	if err != nil {
		return err
	}
	defer copyTestDB.Release()
	writeSubsequentCompressionBlkNumber(c.Dbm, compressTyp, from)
	curBlkNum := c.Chain.CurrentBlock().NumberU64()
	for {

		var (
			subsequentBlkNumber uint64
			err                 error
		)
		switch compressTyp {
		case HeaderCompressType:
			subsequentBlkNumber, _, err = compressHeader(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
		case BodyCompressType:
			subsequentBlkNumber, _, err = compressBody(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
		case ReceiptCompressType:
			subsequentBlkNumber, _, err = compressReceipts(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
		}
		if err != nil {
			return err
		}
		// Copy origin receipts
		c.testCopyOriginData(compressTyp, copyTestDB, from, subsequentBlkNumber-1)
		// Compression integrity test
		if err = c.testVerifyCompressionIntegrity(compressTyp, from, subsequentBlkNumber-1); err != nil {
			return err
		}
		if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
			logger.Info("[Compression] compression is completed", "from", from, "to", to, "subsequentBlkNumber", subsequentBlkNumber)
			break
		}
		from = subsequentBlkNumber
	}
	if _, err := database.WriteBatches(copyTestDB); err != nil {
		return err
	}
	return nil
}

func (c *CompressModule) testCompressPerformance(from, to uint64) error {
	var (
		originFrom = from
		curBlkNum  = c.Chain.CurrentBlock().NumberU64()

		totalHeaderCompressedElapsedTime   time.Duration
		totalBodyCompressedElapsedTime     time.Duration
		totalReceiptsCompressedElapsedTime time.Duration
		totalHeaderFinderElapsedTime       time.Duration
		totalBodyFinderElapsedTime         time.Duration
		totalReceiptsFinderElapsedTime     time.Duration
	)
	for _, compressTyp := range []CompressionType{HeaderCompressType, BodyCompressType, ReceiptCompressType} {
		from = originFrom
		writeSubsequentCompressionBlkNumber(c.Dbm, compressTyp, from)
		for {
			var (
				subsequentBlkNumber uint64
				err                 error
			)
			startCompress := time.Now()
			switch compressTyp {
			case HeaderCompressType:
				subsequentBlkNumber, _, err = compressHeader(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
				totalHeaderCompressedElapsedTime += time.Since(startCompress)
			case BodyCompressType:
				subsequentBlkNumber, _, err = compressBody(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
				totalBodyCompressedElapsedTime += time.Since(startCompress)
			case ReceiptCompressType:
				subsequentBlkNumber, _, err = compressReceipts(c.Dbm, from, 0, curBlkNum, c.getCompressChunk(), c.getChunkCap(), true)
				totalReceiptsCompressedElapsedTime += time.Since(startCompress)
			}
			if err != nil {
				return err
			}
			startFinder := time.Now()
			c.testFinder(compressTyp, from, subsequentBlkNumber-1)
			switch compressTyp {
			case HeaderCompressType:
				totalHeaderFinderElapsedTime += time.Since(startFinder)
			case BodyCompressType:
				totalBodyFinderElapsedTime += time.Since(startFinder)
			case ReceiptCompressType:
				totalReceiptsFinderElapsedTime += time.Since(startFinder)
			}
			if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
				logger.Info("[Compression] compression is completed", "from", from, "to", to, "subsequentBlkNumber", subsequentBlkNumber)
				break
			}
			from = subsequentBlkNumber
		}
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("block range (from=%d) (to=%d) (# of blocks = %d)\n", originFrom, to, to-originFrom)
	fmt.Println("total header   compression elpased time: ", totalHeaderCompressedElapsedTime)
	fmt.Println("total body     compression elpased time: ", totalBodyCompressedElapsedTime)
	fmt.Println("total receipts compression elpased time: ", totalReceiptsCompressedElapsedTime)
	fmt.Println("total header   finder      elpased time: ", totalHeaderFinderElapsedTime)
	fmt.Println("total body     finder      elpased time: ", totalBodyFinderElapsedTime)
	fmt.Println("total receipts finder      elpased time: ", totalReceiptsFinderElapsedTime)
	fmt.Println("--------------------------------------------------")
	return nil
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

	from, to := uint64(0), uint64(3300)
	assert.Nil(t, mCompress.testCompress(HeaderCompressType, from, to, "copy_header"))
	assert.Nil(t, mCompress.testCompress(BodyCompressType, from, to, "copy_body"))
	assert.Nil(t, mCompress.testCompress(ReceiptCompressType, from, to, "copy_receipts"))
}

func TestCompressPerformance(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) {
	mCompress, chainDB := testInit(t, setup)
	if mCompress == nil || chainDB == nil {
		return
	}
	defer chainDB.Close()

	from, to := uint64(0), uint64(3300)
	assert.Nil(t, mCompress.testCompressPerformance(from, to))
}
