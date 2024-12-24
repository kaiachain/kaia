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
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestCompressStorage(t *testing.T) {
	for _, compressTyp := range allCompressTypes {
		switch compressTyp {
		case HeaderCompressType:
			testHeaderCompress(t)
		case BodyCompressType:
			testBodyCompress(t)
		case ReceiptCompressType:
			testReceiptsCompress(t)
		}
	}
}

func TestDecompressStorage(t *testing.T) {
	for _, compressTyp := range allCompressTypes {
		switch compressTyp {
		case HeaderCompressType:
			testHeaderDecompress(t)
		case BodyCompressType:
			testBodyDecompress(t)
		case ReceiptCompressType:
			testReceiptsDecompress(t)
		}
	}
}

func TestDeleteStorage(t *testing.T) {
	for _, compressTyp := range allCompressTypes {
		switch compressTyp {
		case HeaderCompressType:
			testCompressedHeaderDelete(t)
		case BodyCompressType:
			testCompressedBodyDelete(t)
		case ReceiptCompressType:
			testCompressedReceiptsDelete(t)
		}
	}
}

func TestCompressModule(t *testing.T) {
	var (
		nBlocks                           = 100
		_, dbm, headers, bodies, receipts = runCompress(t, nBlocks)
	)
	checkCompressedIntegrity(t, dbm, 0, nBlocks-1, headers, bodies, receipts, false)
}

func TestRewind(t *testing.T) {
	// Scenario Description:
	// 1. `setHead` invoked targeting block number 55
	//    - Removed origin data from 100 to 55
	//    - Removed compressed data from 100 to 50 (because chunk size is set to 10)
	// 2. Generate new blocks from 55 to 100
	// 3. Restart the compression module
	// 	  - start compression from block number 50 to 100

	/*
		[Phase1: Setup]
			Compression completed
			0 ------------ 50 ------------ 99
											^
								next compression block number
			Chunks = C1|C2|C3|C4|C5|C6|C7|C8|C9|C10

		[Phase2: Rewind] (`Stop` called)
			Once `setHead` is invoked,
			0 ------------ 50 ---- 55
							^
				next compression block number
			Chunks = C1|C2|C3|C4|C5

		[Phase3: Compress again] (`Start` called)
			compressed data range 50-59 is restored and Sync is started from 55. Finally,
			0 ------------ 50 ------------ 99
			Chunks = C1|C2|C3|C4|C5|C6|C7|C8|C9|C10
	*/

	var (
		nBlocks                                   = 100
		setHeadTo                                 = 55
		mCompress, dbm, headers, bodies, receipts = runCompress(t, nBlocks)
	)

	mCompress.Stop()
	for i := nBlocks - 1; i >= setHeadTo; i-- {
		num := uint64(i)
		hash := dbm.ReadCanonicalHash(num)
		// delete origin data
		dbm.DeleteBody(hash, num)
		dbm.DeleteReceipts(hash, num)
		dbm.DeleteHeader(hash, num)
		// delete compression data
		mCompress.RewindDelete(hash, num)
	}
	checkCompressedIntegrity(t, dbm, setHeadTo, nBlocks-1, headers, bodies, receipts, true)

	var (
		newHeaders  []*types.Header
		newBodies   []*types.Body
		newReceipts []types.Receipts
	)
	// expected next compression block number should be equal or less than `setHeadTo`.
	// Given the value of `setHeadTo` is 55 and chunk size is 10,
	// The rewinded next compression block number should be 50.
	for _, compressTyp := range allCompressTypes {
		nextCompressionNumber := readSubsequentCompressionBlkNumber(dbm, compressTyp)
		assert.Equal(t, int(nextCompressionNumber), setHeadTo-(setHeadTo%int(mCompress.getCompressChunk())))
	}

	start := setHeadTo + setHeadTo%int(mCompress.getCompressChunk())
	for i := start; i < nBlocks; i++ {
		h := genHeader()
		h.Number = big.NewInt(int64(i))
		hn, hh := h.Number.Uint64(), h.Hash()
		dbm.WriteCanonicalHash(hh, hn)
		dbm.WriteHeader(h)

		b, r := genBody(100), genReceipts(100)
		dbm.WriteBody(hh, hn, b)
		dbm.WriteReceipts(hh, hn, *r)

		newHeaders = append(newHeaders, h)
		newBodies = append(newBodies, b)
		newReceipts = append(newReceipts, *r)
	}
	canonicalHeaders := append(headers[:start], newHeaders[:]...)
	canonicalBodies := append(bodies[:start], newBodies[:]...)
	canonicalReceipts := append(receipts[:start], newReceipts[:]...)

	mCompress.Start() // fragment restore invoked before starting compression
	waitCompression(mCompress)
	checkCompressedIntegrity(t, dbm, 0, nBlocks-1, canonicalHeaders, canonicalBodies, canonicalReceipts, false)

	// Once completed the compression, next compression block number reaches to `nBlocks - 1`
	for _, compressTyp := range allCompressTypes {
		nextCompressionNumber := readSubsequentCompressionBlkNumber(dbm, compressTyp)
		assert.Equal(t, int(nextCompressionNumber), nBlocks-1)
	}
}

func TestRetention(t *testing.T) {
	var (
		nBlocks    = 100
		chain, dbm = initMock(t, nBlocks)
		mCompress  = NewCompression()
		err        = mCompress.Init(&InitOpts{
			ChunkBlockSize: blockchain.DefaultChunkBlockSize,
			ChunkCap:       blockchain.DefaultCompressChunkCap,
			Chain:          chain,
			Dbm:            dbm,
		})
	)
	assert.Nil(t, err)
	dbm.SetCompressModule(mCompress)
	mCompress.setCompressChunk(10)
	mCompress.setCompressRetention(uint64(nBlocks))
	mCompress.Start()
	waitCompression(mCompress)

	// compression is not performed by retention
	for _, compressTyp := range allCompressTypes {
		nextCompressionNumber := readSubsequentCompressionBlkNumber(dbm, compressTyp)
		assert.Equal(t, nextCompressionNumber, uint64(0))
	}
	mCompress.Stop()

	// compress work by reset retention
	mCompress.setCompressRetention(0)
	mCompress.Start()
	waitCompression(mCompress)
	for _, compressTyp := range allCompressTypes {
		nextCompressionNumber := readSubsequentCompressionBlkNumber(dbm, compressTyp)
		assert.Equal(t, int(nextCompressionNumber), nBlocks-1)
	}
}

func TestCache(t *testing.T) {
	var (
		nBlocks        = 100
		chunkBlockSize = uint64(10)
		chain, dbm     = initMock(t, nBlocks)
		mCompress      = NewCompression()
		err            = mCompress.Init(&InitOpts{
			ChunkBlockSize: chunkBlockSize,
			ChunkCap:       blockchain.DefaultCompressChunkCap,
			Chain:          chain,
			Dbm:            dbm,
		})
	)
	assert.Nil(t, err)
	mCompress.Start()
	waitCompression(mCompress)
	targetNum := uint64(30)
	hn, hh := uint64(targetNum), dbm.ReadCanonicalHash(targetNum)

	// header compression cache
	_, ok := getFromCache(HeaderCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.False(t, ok)
	decompressedH, err := findHeaderFromChunkWithBlkHash(dbm, hn, hh)
	assert.Nil(t, err)
	assert.NotNil(t, decompressedH)
	_, ok = getFromCache(HeaderCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.True(t, ok)

	// body compression cache
	_, ok = getFromCache(BodyCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.False(t, ok)
	decompressedB, err := findBodyFromChunkWithBlkHash(dbm, hn, hh)
	assert.Nil(t, err)
	assert.NotNil(t, decompressedB)
	_, ok = getFromCache(BodyCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.True(t, ok)

	// receipts compression cache
	_, ok = getFromCache(ReceiptCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.False(t, ok)
	decompressedR, err := findReceiptsFromChunkWithBlkHash(dbm, hn, hh)
	assert.Nil(t, err)
	assert.NotNil(t, decompressedR)
	_, ok = getFromCache(ReceiptCompressType, targetNum, uint64(targetNum+chunkBlockSize-1))
	assert.True(t, ok)
}
