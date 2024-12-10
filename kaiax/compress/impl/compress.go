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
			residualBlkCnt        = curBlkNum % CompressMigrationThreshold
			nextCompressionBlkNum = readSubsequentCompressionBlkNumber(c.Dbm, compressTyp)
			// Do not wait if next compression block number is far awway. Start migration right now
			noWait = curBlkNum > nextCompressionBlkNum && curBlkNum-nextCompressionBlkNum > CompressMigrationThreshold
		)

		if residualBlkCnt != 0 && !noWait {
			time.Sleep(time.Second * time.Duration(CompressMigrationThreshold-residualBlkCnt))
			continue
		}
		from, to := uint64(0), uint64(0)
		for {
			subsequentBlkNumber, err := compressFn(c.Dbm, from, to, curBlkNum, true)
			if err != nil {
				logger.Warn("[Compression] failed to compress chunk", "err", err)
				break
			}
			if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
				logger.Info("[Compression] compression is completed", "from", from, "to", to, "subsequentBlkNumber", subsequentBlkNumber)
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
		switch compressTyp {
		case HeaderCompressType:
			originHeader := c.Dbm.ReadHeader(c.Dbm.ReadCanonicalHash(i), i)
			compressedHeader, err := findHeaderFromChunkWithBlkHash(c.Dbm, i, originHeader.Hash())
			if err != nil {
				return err
			}
			if originHeader.Hash() != compressedHeader.Hash() {
				return fmt.Errorf("[Header Compression Test] Compressed header is not the same data with origin header data (number=%d, txHash=%s)", i, originHeader.Hash().String())
			}
		case BodyCompressType:
			originBody := c.Dbm.ReadBody(c.Dbm.ReadCanonicalHash(i), i)
			for _, originTx := range originBody.Transactions {
				compressedTx, err := findTxFromChunkWithTxHash(c.Dbm, i, originTx.Hash())
				if err != nil {
					return err
				}
				if !originTx.Equal(compressedTx) {
					return fmt.Errorf("[Body Compression Test] Compressed body is not the same data with origin body data (number=%d, txHash=%s)", i, originTx.Hash().String())
				}
			}
		case ReceiptCompressType:
			for _, originReceipt := range c.Dbm.ReadReceipts(c.Dbm.ReadCanonicalHash(i), i) {
				compressedReceipt, err := findReceiptFromChunkWithTxHash(c.Dbm, i, originReceipt.TxHash)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(originReceipt, compressedReceipt) {
					return fmt.Errorf("[Receipt Compression Test] Compressed receipt is not the same data with origin receipt data (number=%d, txHash=%s)", i, originReceipt.TxHash.String())
				}
			}
		}
	}
	return nil
}

func (c *CompressModule) TestCompress(compressTyp CompressionType, from, to uint64, startNum *uint64, tempDir string) error {
	dbConfig := c.Dbm.GetDBConfig()
	copyTestDB, err := database.TestCreateNewDB(dbConfig, filepath.Join(dbConfig.Dir, tempDir))
	if err != nil {
		return err
	}
	defer copyTestDB.Release()
	if startNum != nil {
		writeSubsequentCompressionBlkNumber(c.Dbm, compressTyp, *startNum)
	}
	curBlkNum := c.Chain.CurrentBlock().NumberU64()
	for {
		// Copy origin receipts
		c.testCopyOriginData(compressTyp, copyTestDB, from, to)

		var (
			subsequentBlkNumber uint64
			err                 error
		)
		switch compressTyp {
		case HeaderCompressType:
			subsequentBlkNumber, err = compressHeader(c.Dbm, from, to, curBlkNum, true)
		case BodyCompressType:
			subsequentBlkNumber, err = compressBody(c.Dbm, from, to, curBlkNum, true)
		case ReceiptCompressType:
			subsequentBlkNumber, err = compressReceipts(c.Dbm, from, to, curBlkNum, true)
		}
		if err != nil {
			return err
		}
		if err = c.testVerifyCompressionIntegrity(compressTyp, from, to); err != nil {
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

func TestCompress(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) {
	var (
		copyTempDirHC         = "copy_header"
		copyTempDirBC         = "copy_body"
		copyTempDirRC         = "copy_receipts"
		bc, chainDB, setupErr = setup(t)
		initOpts              = &InitOpts{
			Chain: bc,
			Dbm:   chainDB,
		}
		mCompress     = NewCompression()
		moduleInitErr = mCompress.Init(initOpts)
	)

	if err := errors.Join(setupErr, moduleInitErr); err != nil {
		if errors.Is(err, compress.ErrInitNil) {
			// If no environment varaible set, do not execute compression test
			// TODO-hyunsooda: Change this test to functional test and remove temp storage directory
			return
		} else {
			t.Fatal(err)
		}
	}
	chainDB.SetCompressModule(mCompress)
	defer chainDB.Close()

	from, to := uint64(0), uint64(3300)
	// receipt compression test
	assert.Nil(t, mCompress.TestCompress(HeaderCompressType, from, to, &from, copyTempDirHC))
	assert.Nil(t, mCompress.TestCompress(BodyCompressType, from, to, &from, copyTempDirBC))
	assert.Nil(t, mCompress.TestCompress(ReceiptCompressType, from, to, &from, copyTempDirRC))
}
