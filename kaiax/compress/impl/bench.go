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
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCopyOriginData(c *CompressModule, compressTyp CompressionType, copyTestDB database.Batch, from, to uint64) {
	// Copy origin
	for i := from; i < to; i++ {
		hash := c.Dbm.ReadCanonicalHash(i)
		switch compressTyp {
		case HeaderCompressType:
			originHeader := c.Dbm.ReadHeader(hash, i)
			c.Dbm.PutHeaderToBatch(copyTestDB, originHeader)
		case BodyCompressType:
			originBody := c.Dbm.ReadBody(hash, i)
			c.Dbm.PutBodyToBatch(copyTestDB, hash, i, originBody)
		case ReceiptCompressType:
			originReceipts := c.Dbm.ReadReceipts(hash, i)
			c.Dbm.PutReceiptsToBatch(copyTestDB, hash, i, originReceipts)
		}
	}
}

func testVerifyCompressionIntegrity(c *CompressModule, compressTyp CompressionType, from, to uint64) error {
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

func testCompress(c *CompressModule, t *testing.T, compressTyp CompressionType, from, to uint64, tempDir string, wg *sync.WaitGroup) {
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
		testCopyOriginData(c, compressTyp, copyTestDB, from, subsequentBlkNumber-1)
		_, err = database.WriteBatches(copyTestDB)
		assert.Nil(t, err)
		// Compression integrity test
		err = testVerifyCompressionIntegrity(c, compressTyp, from, subsequentBlkNumber-1)
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
	go testCompress(mCompress, t, HeaderCompressType, from, to, "copy_header", &wg)
	go testCompress(mCompress, t, BodyCompressType, from, to, "copy_body", &wg)
	go testCompress(mCompress, t, ReceiptCompressType, from, to, "copy_receipts", &wg)
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

func TestCompressQuery(t *testing.T, setup func(t *testing.T) (*blockchain.BlockChain, database.DBManager, error)) {
	url := os.Getenv("TEST_COMPRESS_ENDPOINT_URL")
	if url == "" {
		fmt.Println("Endpoint URL is empty (Set thorugh `TEST_COMPRESS_ENDPOINT_URL`")
		return
	}
	mCompress, chainDB := testInit(t, setup)
	if mCompress == nil || chainDB == nil {
		return
	}
	defer chainDB.Close()

	rand.Seed(time.Now().UnixNano())

	var (
		c, err          = client.Dial(url)
		nextNum         = uint64(0)
		dist            = uint64(10000)
		latestBlockHash = chainDB.ReadHeadBlockHash()
		latestBlockNum  = *chainDB.ReadHeaderNumber(latestBlockHash)
		maxNum          = latestBlockNum - (dist * 10)
		queryCnt        = 0
		headerElapsed   time.Duration
		bodyElapsed     time.Duration
		receiptsElapsed time.Duration
	)
	require.Nil(t, err)

	for {
		clearCache()
		r := nextNum + (rand.Uint64() % dist)
		if r > maxNum {
			break
		}
		// header verification
		chStart := time.Now()
		oh, err := c.HeaderByNumber(context.Background(), big.NewInt(int64(r)))
		chEnd := time.Since(chStart)
		require.Nil(t, err)
		ohHash := oh.Hash()
		headerStart := time.Now()
		ch := chainDB.ReadHeader(ohHash, oh.Number.Uint64())
		headerElapsed = time.Since(headerStart)
		require.Equal(t, ohHash, ch.Hash())

		// body and receipts verification
		cbStart := time.Now()
		ob, err := c.BlockByNumber(context.Background(), big.NewInt(int64(r)))
		cbEnd := time.Since(cbStart)
		require.Nil(t, err)
		txs := ob.Transactions()
		if len(txs) > 0 {
			bodyStart := time.Now()
			compressedBody := chainDB.ReadBody(ohHash, r)
			bodyElapsed = time.Since(bodyStart)
			receiptsStart := time.Now()
			compressedReceipts := chainDB.ReadReceipts(ohHash, r)
			receiptsElapsed = time.Since(receiptsStart)
			for idx, originTx := range txs {
				require.True(t, originTx.Equal(compressedBody.Transactions[idx]))

				or, err := c.TransactionReceipt(context.Background(), originTx.Hash())
				require.Nil(t, err)
				if compressedReceipts[idx].Status != types.ReceiptStatusSuccessful {
					compressedReceipts[idx].Status = types.ReceiptStatusFailed
				}
				require.True(t, reflect.DeepEqual(or, compressedReceipts[idx]))
			}
		}
		fmt.Printf("[%10d] queryNum=%10d, headerElapsed=%s, bodyElapsed=%s, receiptsElapsed=%s, queryHeaderElapsed=%s, queryBlockElapsed=%s, progress=%5f%%, txs=%d\n",
			queryCnt, r, headerElapsed.String(), bodyElapsed.String(), receiptsElapsed.String(), chEnd.String(), cbEnd.String(), float64(r)/float64(maxNum)*100.0, len(txs))
		nextNum += dist
		queryCnt++
	}
	fmt.Println("DONE")
}
