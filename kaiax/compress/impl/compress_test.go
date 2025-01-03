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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	blockchain_mock "github.com/kaiachain/kaia/kaiax/compress/impl/mock"
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
			Enable:         true,
		})
	)
	assert.Nil(t, err)
	mCompress.loopIdleTime = 0
	dbm.SetCompressModule(mCompress)
	mCompress.setCompressChunk(10)
	mCompress.setCompressRetention(uint64(nBlocks))
	mCompress.Start()
	waitCompression(mCompress)

	assertNextCompressNum(t, dbm, uint64(0))
	mCompress.Stop()

	// compress work by reset retention
	mCompress.setCompressRetention(0)
	mCompress.Start()
	waitCompression(mCompress)
	assertNextCompressNum(t, dbm, uint64(nBlocks-1))
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
			Enable:         true,
		})
	)
	assert.Nil(t, err)
	mCompress.loopIdleTime = 0
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

func TestRewind(t *testing.T) {
	setHeadBlock1 := uint64(6)
	testRewind(t, &rewindTest{
		canonicalBlocks:    24,
		commitBlock:        4,
		pivotBlock:         nil,
		setheadBlock:       setHeadBlock1,
		expCanonicalBlocks: int(setHeadBlock1),
		expSidechainBlocks: 0,
		expHeadHeader:      setHeadBlock1,
		expHeadFastBlock:   setHeadBlock1,
		expHeadBlock:       setHeadBlock1,
	})
}

func TestEnableFlag(t *testing.T) {
	var (
		nBlocks                                   = 100
		mCompress, dbm, headers, bodies, receipts = runCompress(t, nBlocks)
	)
	checkCompressedIntegrity(t, dbm, 0, nBlocks-1, headers, bodies, receipts, false)

	// further insertion 100 blocks
	lastHeader1 := genBlocks(dbm, nBlocks, nBlocks)

	chain := blockchain_mock.NewMockBlockChain(gomock.NewController(t))
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(lastHeader1)).AnyTimes()
	mCompress.Chain = chain
	assertNextCompressNum(t, dbm, uint64(nBlocks-1))

	mCompress.Stop()
	mCompress.Enable = false
	mCompress.Start()
	time.Sleep(time.Second * 2)
	assertNextCompressNum(t, dbm, uint64(nBlocks-1))

	mCompress.Stop()
	mCompress.Enable = true
	mCompress.Start()
	time.Sleep(time.Second * 2)
	assertNextCompressNum(t, dbm, lastHeader1.Number.Uint64())

	// further insertion 100 blocks
	lastHeader2 := genBlocks(dbm, nBlocks*2, nBlocks)
	chain = blockchain_mock.NewMockBlockChain(gomock.NewController(t))
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(lastHeader2)).AnyTimes()
	mCompress.Chain = chain

	mCompress.Stop()
	mCompress.Enable = false
	mCompress.Start()
	time.Sleep(time.Second * 2)
	assertNextCompressNum(t, dbm, lastHeader1.Number.Uint64())

	mCompress.Stop()
	mCompress.Enable = true
	mCompress.Start()
	time.Sleep(time.Second * 2)
	assertNextCompressNum(t, dbm, lastHeader2.Number.Uint64())

	// The query remains accessible even if the module is stopped or disabled.
	mCompress.Stop()
	mCompress.Enable = false
	clearCache()
	for i := 0; i < int(lastHeader2.Number.Uint64()); i++ {
		hn := uint64(i)
		hh := dbm.ReadCanonicalHash(hn)
		dbm.ReadHeader(hh, hn)
		dbm.ReadBody(hh, hn)
		dbm.ReadReceipts(hh, hn)
	}
}
