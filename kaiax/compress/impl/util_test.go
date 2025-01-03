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
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	mrand "math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/crypto"
	blockchain_mock "github.com/kaiachain/kaia/kaiax/compress/impl/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

type rewindTest struct {
	canonicalBlocks int     // Number of blocks to generate for the canonical chain (heavier)
	sidechainBlocks int     // Number of blocks to generate for the side chain (lighter)
	commitBlock     uint64  // Block number for which to commit the state to disk
	pivotBlock      *uint64 // Pivot block number in case of fast sync

	setheadBlock       uint64 // Block number to set head back to
	expCanonicalBlocks int    // Number of canonical blocks expected to remain in the database (excl. genesis)
	expSidechainBlocks int    // Number of sidechain blocks expected to remain in the database (excl. genesis)
	expHeadHeader      uint64 // Block number of the expected head header
	expHeadFastBlock   uint64 // Block number of the expected head fast sync block
	expHeadBlock       uint64 // Block number of the expected head full block
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func generateBigInt(nBits int) *big.Int {
	v, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(nBits)))
	return v
}

func generateUint8() uint8 {
	return generateRandomBytes(64)[0]
}

func generateUint64() uint64 {
	return binary.BigEndian.Uint64(generateRandomBytes(64))
}

func generateAddress() common.Address {
	randomBytes := generateRandomBytes(common.AddressLength)
	return common.Address(randomBytes)
}

func generateBloom() types.Bloom {
	randomBytes := generateRandomBytes(types.BloomByteLength)
	return types.Bloom(randomBytes)
}

func generateHash() common.Hash {
	randomBytes := generateRandomBytes(common.HashLength)
	hash := sha256.Sum256(randomBytes)
	return common.Hash(hash)
}

func genHeader() *types.Header {
	return &types.Header{
		ParentHash:   generateHash(),
		Rewardbase:   generateAddress(),
		Root:         generateHash(),
		TxHash:       generateHash(),
		ReceiptHash:  generateHash(),
		Bloom:        generateBloom(),
		BlockScore:   generateBigInt(64),
		Number:       generateBigInt(64),
		GasUsed:      generateUint64(),
		Time:         generateBigInt(64),
		TimeFoS:      generateUint8(),
		Extra:        generateRandomBytes(8096),
		Governance:   generateRandomBytes(512),
		Vote:         generateRandomBytes(512),
		BaseFee:      generateBigInt(64),
		RandomReveal: generateRandomBytes(192),
		MixHash:      generateHash().Bytes(),
	}
}

func genHeaders() ([]*types.Header, []*types.Header) {
	h1, h2, h3 := genHeader(), genHeader(), genHeader()
	copyH2, copyH3 := types.CopyHeader(h1), types.CopyHeader(h1)
	headers := []*types.Header{h1, h2, h3}
	copyHeaders := []*types.Header{h1, copyH2, copyH3}

	for idx := range headers {
		headers[idx].Number = big.NewInt(int64(idx) + 1)
		copyHeaders[idx].Number = big.NewInt(int64(idx) + 1)
	}
	return headers, copyHeaders
}

func genBody(n int) *types.Body {
	testKey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	txs := make([]*types.Transaction, n)
	for i := range n {
		tx := types.NewTransaction(
			generateUint64(),
			common.StringToAddress("ASD"),
			generateBigInt(64),
			generateUint64(),
			generateBigInt(64),
			generateRandomBytes(128),
		)
		signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(big.NewInt(1)), testKey)
		if err != nil {
			panic(err)
		}
		txs[i] = signedTx
	}
	return &types.Body{
		Transactions: txs,
	}
}

func genBodies() ([]*types.Body, []*types.Body, int) {
	b1, b2, b3 := genBody(100), genBody(100), genBody(100)
	copyB2, copyB3 := *b1, *b1
	bodies := []*types.Body{b1, b2, b3}
	copyBodies := []*types.Body{b1, &copyB2, &copyB3}
	return bodies, copyBodies, len(bodies)
}

func genReceipts(n int) *types.Receipts {
	receipts := make([]*types.Receipt, n)
	for i := range n {
		receipt := &types.Receipt{
			Status:          uint(generateUint64()),
			Bloom:           types.Bloom{},
			Logs:            []*types.Log{},
			TxHash:          generateHash(),
			ContractAddress: generateAddress(),
			GasUsed:         generateUint64(),
		}
		receipts[i] = receipt
	}
	ret := types.Receipts(receipts)
	return &ret
}

func genReceiptsSlice() ([]*types.Receipts, []*types.Receipts, int) {
	r1, r2, r3 := genReceipts(100), genReceipts(100), genReceipts(100)
	copyR2, copyR3 := *r1, *r1
	receipts := []*types.Receipts{r1, r2, r3}
	copyBodies := []*types.Receipts{r1, &copyR2, &copyR3}
	return receipts, copyBodies, len(receipts)
}

func headerCompress(from, to, headerNumber uint64, headers []*types.Header) (database.DBManager, uint64, common.StorageSize, common.StorageSize, error) {
	dbm := database.NewMemoryDBManager()
	var originSize common.StorageSize
	for _, h := range headers {
		dbm.WriteCanonicalHash(h.Hash(), h.Number.Uint64())
		dbm.WriteHeader(h)
		originSize += h.Size()
	}
	nextCompressNum, _, compressedSize, err := compressHeader(dbm, from, to, headerNumber, blockchain.DefaultChunkBlockSize, blockchain.DefaultCompressChunkCap, true)
	return dbm, nextCompressNum, originSize, common.StorageSize(float64(compressedSize)), err
}

func bodyCompress(from, to, headerNumber uint64, blkHashes []common.Hash, bodies []*types.Body) (database.DBManager, uint64, common.StorageSize, common.StorageSize, error) {
	originSize := CompressionSize(0)
	dbm := database.NewMemoryDBManager()
	for idx, body := range bodies {
		blkNum := uint64(idx) + 1
		dbm.WriteCanonicalHash(blkHashes[idx], blkNum)
		dbm.WriteBody(blkHashes[idx], blkNum, body)
		rlp.Encode(&originSize, body)
	}
	nextCompressNum, _, compressedSize, err := compressBody(dbm, from, to, headerNumber, blockchain.DefaultChunkBlockSize, blockchain.DefaultCompressChunkCap, true)
	return dbm, nextCompressNum, common.StorageSize(originSize), common.StorageSize(float64(compressedSize)), err
}

func receiptsCompress(from, to, headerNumber uint64, blkHashes []common.Hash, receiptsSlice []*types.Receipts) (database.DBManager, uint64, common.StorageSize, common.StorageSize, error) {
	originSize := CompressionSize(0)
	dbm := database.NewMemoryDBManager()
	for idx, receipts := range receiptsSlice {
		blkNum := uint64(idx) + 1
		dbm.WriteCanonicalHash(blkHashes[idx], blkNum)
		dbm.WriteReceipts(blkHashes[idx], blkNum, *receipts)
		rlp.Encode(&originSize, receipts)
	}
	nextCompressNum, _, compressedSize, err := compressReceipts(dbm, from, to, headerNumber, blockchain.DefaultChunkBlockSize, blockchain.DefaultCompressChunkCap, true)
	return dbm, nextCompressNum, common.StorageSize(originSize), common.StorageSize(float64(compressedSize)), err
}

func testHeaderCompress(t *testing.T) {
	headers, copyHeaders := genHeaders()
	from, to, headerNumber := headers[0].Number.Uint64(), headers[len(headers)-1].Number.Uint64()+1, headers[0].Number.Uint64()

	_, nextCompressNum, originSize, compressedSize1, err := headerCompress(from, to, headerNumber, headers)
	_, _, _, compressedSize2, _ := headerCompress(from, to, headerNumber, copyHeaders)
	assert.Nil(t, err)
	assert.Equal(t, nextCompressNum, to)
	assert.True(t, originSize > compressedSize1)
	// Since copied header compression has a higher entropy, the compressed size should be effective more.
	assert.True(t, compressedSize1 > compressedSize2)
}

func testBodyCompress(t *testing.T) {
	bodies, copyBodies, bodyLen := genBodies()
	from, to, headerNumber := uint64(1), uint64(bodyLen+1), uint64(1)
	blkHashes := make([]common.Hash, bodyLen)

	for i := range bodyLen {
		blkHashes[i] = generateHash()
	}

	_, nextCompressNum, originSize, compressedSize1, err := bodyCompress(from, to, headerNumber, blkHashes, bodies)
	_, _, _, compressedSize2, _ := bodyCompress(from, to, headerNumber, blkHashes, copyBodies)
	assert.Nil(t, err)
	assert.Equal(t, nextCompressNum, to)
	assert.True(t, originSize > compressedSize1)
	// Since copied header compression has a higher entropy, the compressed size should be effective more.
	assert.True(t, compressedSize1 > compressedSize2)
}

func testReceiptsCompress(t *testing.T) {
	receipts, copyReceipts, receiptsLen := genReceiptsSlice()
	from, to, headerNumber := uint64(1), uint64(receiptsLen+1), uint64(1)
	blkHashes := make([]common.Hash, receiptsLen)

	for i := range receiptsLen {
		blkHashes[i] = generateHash()
	}

	_, nextCompressNum, originSize, compressedSize1, err := receiptsCompress(from, to, headerNumber, blkHashes, receipts)
	_, _, _, compressedSize2, _ := receiptsCompress(from, to, headerNumber, blkHashes, copyReceipts)
	assert.Nil(t, err)
	assert.Equal(t, nextCompressNum, to)
	assert.True(t, originSize > compressedSize1)
	// Since copied header compression has a higher entropy, the compressed size should be effective more.
	assert.True(t, compressedSize1 > compressedSize2)
}

func testHeaderDecompress(t *testing.T) {
	h1, h2 := genHeader(), genHeader()
	h1.Number = big.NewInt(1)
	h2.Number = big.NewInt(2)
	h1Num, h2Num := h1.Number.Uint64(), h2.Number.Uint64()
	headers := []*types.Header{h1, h2}
	dbm, _, _, _, err := headerCompress(h1Num, h2Num+1, h1Num, headers)
	assert.Nil(t, err)

	decompressedH1, err := findHeaderFromChunkWithBlkHash(dbm, h1Num, h1.Hash())
	assert.Nil(t, err)
	decompressedH2, err := findHeaderFromChunkWithBlkHash(dbm, h2Num, h2.Hash())
	assert.Nil(t, err)
	assert.Equal(t, decompressedH1.Hash(), h1.Hash())
	assert.Equal(t, decompressedH2.Hash(), h2.Hash())
}

func testBodyDecompress(t *testing.T) {
	b1, b2 := genBody(100), genBody(100)
	bodies := []*types.Body{b1, b2}
	blkHashes := []common.Hash{generateHash(), generateHash()}
	from, to, headerNumber := uint64(1), uint64(3), uint64(1)
	dbm, _, _, _, err := bodyCompress(from, to, headerNumber, blkHashes, bodies)
	assert.Nil(t, err)

	decompressedB1, err := findBodyFromChunkWithBlkHash(dbm, 1, blkHashes[0])
	assert.Nil(t, err)
	decompressedB2, err := findBodyFromChunkWithBlkHash(dbm, 2, blkHashes[1])
	assert.Nil(t, err)

	for idx, originTx := range b1.Transactions {
		assert.True(t, originTx.Equal(decompressedB1.Transactions[idx]))
	}
	for idx, originTx := range b2.Transactions {
		assert.True(t, originTx.Equal(decompressedB2.Transactions[idx]))
	}
}

func testReceiptsDecompress(t *testing.T) {
	r1, r2 := genReceipts(100), genReceipts(100)
	receipts := []*types.Receipts{r1, r2}
	blkHashes := []common.Hash{generateHash(), generateHash()}
	from, to, headerNumber := uint64(1), uint64(3), uint64(1)
	dbm, _, _, _, err := receiptsCompress(from, to, headerNumber, blkHashes, receipts)
	assert.Nil(t, err)

	decompressedR1, err := findReceiptsFromChunkWithBlkHash(dbm, 1, blkHashes[0])
	assert.Nil(t, err)
	decompressedR2, err := findReceiptsFromChunkWithBlkHash(dbm, 2, blkHashes[1])
	assert.Nil(t, err)

	for idx, originReceipt := range *r1 {
		assert.True(t, reflect.DeepEqual(originReceipt, decompressedR1[idx]))
	}
	for idx, originReceipt := range *r2 {
		assert.True(t, reflect.DeepEqual(originReceipt, decompressedR2[idx]))
	}
}

func testCompressedHeaderDelete(t *testing.T) {
	headers, _ := genHeaders()
	h1Num, h1Hash := headers[0].Number.Uint64(), headers[0].Hash()
	from, to, headerNumber := h1Num, headers[len(headers)-1].Number.Uint64()+1, h1Num

	// compress headers
	dbm, _, _, _, err := headerCompress(from, to, headerNumber, headers)
	assert.Nil(t, err)
	// find a compressed chunk
	decompressedH1, err := findHeaderFromChunkWithBlkHash(dbm, h1Num, h1Hash)
	assert.Nil(t, err)
	assert.Equal(t, decompressedH1.Hash(), h1Hash)

	// removed a compressed chunk
	subsequentBlkNum, _, err := deleteHeaderFromChunk(dbm, h1Num, h1Hash)
	assert.Nil(t, err)
	assert.Equal(t, subsequentBlkNum, h1Num)

	// try to find a compressed chunk again
	decompressedH1, err = findHeaderFromChunkWithBlkHash(dbm, h1Num, h1Hash)
	assert.Nil(t, decompressedH1)
	assert.NotNil(t, err)
}

func testCompressedBodyDelete(t *testing.T) {
	bodies, _, bodyLen := genBodies()
	from, to, headerNumber := uint64(1), uint64(bodyLen+1), uint64(1)
	blkHashes := make([]common.Hash, bodyLen)
	for i := range bodyLen {
		blkHashes[i] = generateHash()
	}
	targetNum, targetHash := uint64(1), blkHashes[0]

	// compress bodies
	dbm, _, _, _, err := bodyCompress(from, to, headerNumber, blkHashes, bodies)
	assert.Nil(t, err)
	// find a compressed chunk
	decompressedB1, err := findBodyFromChunkWithBlkHash(dbm, targetNum, targetHash)
	assert.Nil(t, err)
	for idx, originTx := range bodies[0].Transactions {
		assert.True(t, originTx.Equal(decompressedB1.Transactions[idx]))
	}

	// removed a compressed chunk
	subsequentBlkNum, _, err := deleteBodyFromChunk(dbm, targetNum, targetHash)
	assert.Nil(t, err)
	assert.Equal(t, subsequentBlkNum, targetNum)

	// try to find a compressed chunk again
	decompressedB1, err = findBodyFromChunkWithBlkHash(dbm, targetNum, targetHash)
	assert.Nil(t, decompressedB1)
	assert.NotNil(t, err)
}

func testCompressedReceiptsDelete(t *testing.T) {
	receipts, _, receiptsLen := genReceiptsSlice()
	from, to, headerNumber := uint64(1), uint64(receiptsLen+1), uint64(1)
	blkHashes := make([]common.Hash, receiptsLen)
	for i := range receiptsLen {
		blkHashes[i] = generateHash()
	}
	targetNum, targetHash := uint64(1), blkHashes[0]

	// compress receipts
	dbm, _, _, _, err := receiptsCompress(from, to, headerNumber, blkHashes, receipts)
	assert.Nil(t, err)
	// find a compressed chunk
	decompressedR1, err := findReceiptsFromChunkWithBlkHash(dbm, targetNum, targetHash)
	assert.Nil(t, err)
	for idx, originReceipt := range *receipts[targetNum-1] {
		assert.True(t, reflect.DeepEqual(originReceipt, decompressedR1[idx]))
	}

	// removed a compressed chunk
	subsequentBlkNum, _, err := deleteReceiptsFromChunk(dbm, targetNum, targetHash)
	assert.Nil(t, err)
	assert.Equal(t, subsequentBlkNum, targetNum)

	// try to find a compressed chunk again
	decompressedR1, err = findReceiptsFromChunkWithBlkHash(dbm, targetNum, targetHash)
	assert.Nil(t, decompressedR1)
	assert.NotNil(t, err)
}

func initMock(t *testing.T, n int) (*blockchain_mock.MockBlockChain, database.DBManager) {
	var (
		chain      = blockchain_mock.NewMockBlockChain(gomock.NewController(t))
		dbm        = database.NewMemoryDBManager()
		lastHeader *types.Header
	)

	// insert `n` blocks (header, body, receipts)
	for i := range n {
		h := genHeader()
		h.Number = big.NewInt(int64(i))
		hn, hh := h.Number.Uint64(), h.Hash()
		dbm.WriteCanonicalHash(hh, hn)
		dbm.WriteHeader(h)

		b, r := genBody(100), genReceipts(100)
		dbm.WriteBody(hh, hn, b)
		dbm.WriteReceipts(hh, hn, *r)
		lastHeader = h
	}

	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(lastHeader)).AnyTimes()
	return chain, dbm
}

func waitCompression(m *CompressModule) {
	time.AfterFunc(SEC_TEN*2, func() { panic("Compression timeout") })
	var (
		headerCompressCompleted   = false
		bodyCompressCompleted     = false
		receiptsCompressCompleted = false
	)
	for {
		if is := m.getIdleState(HeaderCompressType); is.isIdle {
			headerCompressCompleted = true
		}
		if is := m.getIdleState(BodyCompressType); is.isIdle {
			bodyCompressCompleted = true
		}
		if is := m.getIdleState(ReceiptCompressType); is.isIdle {
			receiptsCompressCompleted = true
		}
		if headerCompressCompleted &&
			bodyCompressCompleted &&
			receiptsCompressCompleted {
			return
		}
	}
}

func readOriginData(t *testing.T, dbm database.DBManager, nBlocks int) ([]*types.Header, []*types.Body, []types.Receipts) {
	var (
		headers  = make([]*types.Header, nBlocks)
		bodies   = make([]*types.Body, nBlocks)
		receipts = make([]types.Receipts, nBlocks)
	)
	for i := range nBlocks {
		num := uint64(i)
		hash := dbm.ReadCanonicalHash(num)
		h := dbm.ReadHeader(hash, num)
		assert.NotNil(t, h)
		headers[i] = h
		b := dbm.ReadBody(hash, num)
		assert.NotNil(t, b)
		bodies[i] = b
		r := dbm.ReadReceipts(hash, num)
		assert.NotNil(t, r)
		receipts[i] = r
	}
	return headers, bodies, receipts
}

func assertNextCompressNum(t *testing.T, dbm database.DBManager, expected uint64) {
	for _, compressTyp := range allCompressTypes {
		nextCompressionNumber := readSubsequentCompressionBlkNumber(dbm, compressTyp)
		assert.Equal(t, nextCompressionNumber, expected)
	}
}

func runCompress(t *testing.T, nBlocks int) (*CompressModule, database.DBManager, []*types.Header, []*types.Body, []types.Receipts) {
	var (
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
	dbm.SetCompressModule(mCompress)
	mCompress.loopIdleTime = 0
	headers, bodies, receipts := readOriginData(t, dbm, nBlocks)
	mCompress.setCompressChunk(10)
	mCompress.Start()
	waitCompression(mCompress)
	return mCompress, dbm, headers, bodies, receipts
}

func checkCompressedIntegrity(t *testing.T, dbm database.DBManager, from, to int, originHeaders []*types.Header, originBodies []*types.Body, originReceipts []types.Receipts, mustErr bool) {
	for i := from; i < to; i++ {
		num := uint64(i)
		hash := dbm.ReadCanonicalHash(num)
		// compressed header integrity verification
		{
			decompressedH, err := findHeaderFromChunkWithBlkHash(dbm, num, hash)
			if mustErr {
				assert.NotNil(t, err)
				assert.Nil(t, dbm.ReadHeader(hash, num))
			} else {
				assert.Equal(t, decompressedH.Hash(), originHeaders[i].Hash())
			}
		}
		// compressed body integrity verification
		{
			decompressedB, err := findBodyFromChunkWithBlkHash(dbm, num, hash)
			if mustErr {
				assert.Nil(t, dbm.ReadBody(hash, num))
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for idx, originTx := range originBodies[i].Transactions {
					assert.True(t, originTx.Equal(decompressedB.Transactions[idx]))
				}
			}
		}
		// compressed receipts integrity verification
		{
			decompressedR, err := findReceiptsFromChunkWithBlkHash(dbm, num, hash)
			if mustErr {
				assert.Nil(t, dbm.ReadReceipts(hash, num))
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for idx, originReceipt := range originReceipts[i] {
					assert.True(t, reflect.DeepEqual(originReceipt, decompressedR[idx]))
				}
			}
		}
	}
}

// verifyNoGaps checks that there are no gaps after the initial set of blocks in
// the database and errors if found.
func verifyNoGaps(t *testing.T, chain *blockchain.BlockChain, canonical bool, inserted types.Blocks) {
	t.Helper()

	var end uint64
	for i := uint64(0); i <= uint64(len(inserted)); i++ {
		header := chain.GetHeaderByNumber(i)
		if header == nil && end == 0 {
			end = i
		}
		if header != nil && end > 0 {
			if canonical {
				t.Errorf("Canonical header gap between #%d-#%d", end, i-1)
			} else {
				t.Errorf("Sidechain header gap between #%d-#%d", end, i-1)
			}
			end = 0 // Reset for further gap detection
		}
	}
	end = 0
	for i := uint64(0); i <= uint64(len(inserted)); i++ {
		block := chain.GetBlockByNumber(i)
		if block == nil && end == 0 {
			end = i
		}
		if block != nil && end > 0 {
			if canonical {
				t.Errorf("Canonical block gap between #%d-#%d", end, i-1)
			} else {
				t.Errorf("Sidechain block gap between #%d-#%d", end, i-1)
			}
			end = 0 // Reset for further gap detection
		}
	}
	end = 0
	for i := uint64(1); i <= uint64(len(inserted)); i++ {
		receipts := chain.GetReceiptsByBlockHash(inserted[i-1].Hash())
		if receipts == nil && end == 0 {
			end = i
		}
		if receipts != nil && end > 0 {
			if canonical {
				t.Errorf("Canonical receipt gap between #%d-#%d", end, i-1)
			} else {
				t.Errorf("Sidechain receipt gap between #%d-#%d", end, i-1)
			}
			end = 0 // Reset for further gap detection
		}
	}
}

// verifyCutoff checks that there are no chain data available in the chain after
// the specified limit, but that it is available before.
func verifyCutoff(t *testing.T, chain *blockchain.BlockChain, canonical bool, inserted types.Blocks, head int) {
	t.Helper()

	for i := 1; i <= len(inserted); i++ {
		if i <= head {
			if header := chain.GetHeader(inserted[i-1].Hash(), uint64(i)); header == nil {
				if canonical {
					t.Errorf("Canonical header   #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain header   #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if block := chain.GetBlock(inserted[i-1].Hash(), uint64(i)); block == nil {
				if canonical {
					t.Errorf("Canonical block    #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain block    #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if body := chain.GetReceiptsByBlockHash(inserted[i-1].Hash()); body == nil {
				if canonical {
					t.Errorf("Canonical body #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain body #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if receipts := chain.GetReceiptsByBlockHash(inserted[i-1].Hash()); receipts == nil {
				if canonical {
					t.Errorf("Canonical receipts #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain receipts #%2d [%x...] missing before cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
		} else {
			if header := chain.GetHeader(inserted[i-1].Hash(), uint64(i)); header != nil {
				if canonical {
					t.Errorf("Canonical header   #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain header   #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if block := chain.GetBlock(inserted[i-1].Hash(), uint64(i)); block != nil {
				if canonical {
					t.Errorf("Canonical block    #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain block    #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if body := chain.GetReceiptsByBlockHash(inserted[i-1].Hash()); body != nil {
				if canonical {
					t.Errorf("Canonical body #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain body #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
			if receipts := chain.GetReceiptsByBlockHash(inserted[i-1].Hash()); receipts != nil {
				if canonical {
					t.Errorf("Canonical receipts #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				} else {
					t.Errorf("Sidechain receipts #%2d [%x...] present after cap %d", inserted[i-1].Number(), inserted[i-1].Hash().Bytes()[:3], head)
				}
			}
		}
	}
}

func testRewind(t *testing.T, tt *rewindTest) {
	var (
		db      = database.NewMemoryDBManager()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &blockchain.Genesis{
			Config: params.TestChainConfig,
			Alloc:  blockchain.GenesisAlloc{address: {Balance: funds}},
		}
		signer    = types.LatestSignerForChainID(gspec.Config.ChainID)
		genesis   = gspec.MustCommit(db)
		chain, _  = blockchain.NewBlockChain(db, nil, gspec.Config, gxhash.NewFullFaker(), vm.Config{})
		mCompress = NewCompression()
		err       = mCompress.Init(&InitOpts{
			ChunkBlockSize: blockchain.DefaultChunkBlockSize,
			ChunkCap:       blockchain.DefaultCompressChunkCap,
			Chain:          chain,
			Dbm:            db,
			Enable:         true,
		})
	)
	assert.Nil(t, err)
	mCompress.loopIdleTime = 0

	chain.RegisterRewindableModule(mCompress)

	// Initialize a fresh chain with only a genesis block
	defer chain.Stop()
	chain.Config().Istanbul = params.GetDefaultIstanbulConfig()
	mrand.Seed(time.Now().UnixNano())

	canonblocks, _ := blockchain.GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, tt.canonicalBlocks, func(i int, b *blockchain.BlockGen) {
		for range 3 {
			tx, err := types.SignTx(types.NewTransaction(b.TxNonce(address), generateAddress(), generateBigInt(10), params.TxGas, nil, nil), signer, key)
			assert.Nil(t, err)
			b.AddTx(tx)
		}
	})
	if _, err := chain.InsertChain(canonblocks[:tt.commitBlock]); err != nil {
		t.Fatalf("Failed to import canonical chain start: %v", err)
	}
	if tt.commitBlock > 0 {
		chain.StateCache().TrieDB().Commit(canonblocks[tt.commitBlock-1].Root(), true, tt.commitBlock)
	}
	if _, err := chain.InsertChain(canonblocks[tt.commitBlock:]); err != nil {
		t.Fatalf("Failed to import canonical chain tail: %v", err)
	}

	// Manually dereference anything not committed to not have to work with 128+ tries
	for _, block := range canonblocks {
		chain.StateCache().TrieDB().Dereference(block.Root())
	}

	db.SetCompressModule(mCompress)
	blockChunks := uint64(5)
	mCompress.setCompressChunk(blockChunks)
	mCompress.Start()
	waitCompression(mCompress)

	// Set the head of the chain back to the requested number
	chain.SetHead(tt.setheadBlock)

	// Iterate over all the remaining blocks and ensure there are no gaps
	verifyNoGaps(t, chain, true, canonblocks)
	verifyCutoff(t, chain, true, canonblocks, tt.expCanonicalBlocks)

	if head := chain.CurrentHeader(); head.Number.Uint64() != tt.expHeadHeader {
		t.Errorf("Head header mismatch!: have %d, want %d", head.Number, tt.expHeadHeader)
	}
	if head := chain.CurrentFastBlock(); head.NumberU64() != tt.expHeadFastBlock {
		t.Errorf("Head fast block mismatch: have %d, want %d", head.NumberU64(), tt.expHeadFastBlock)
	}
	if head := chain.CurrentBlock(); head.NumberU64() != tt.expHeadBlock {
		t.Errorf("Head block mismatch!!: have %d, want %d", head.NumberU64(), tt.expHeadBlock)
	}

	// check next compression block number
	assertNextCompressNum(t, db, blockChunks)

	// check canonical blocks after rewinding
	for i := 0; i < len(canonblocks); i++ {
		block := chain.GetHeader(canonblocks[i].Hash(), canonblocks[i].NumberU64())
		body := chain.GetBodyRLP(canonblocks[i].Hash())
		receipts := chain.GetReceiptsByBlockHash(canonblocks[i].Hash())
		if canonblocks[i].Number().Uint64() <= tt.setheadBlock {
			assert.False(t, block == nil)
			assert.False(t, body == nil)
			assert.False(t, receipts == nil)
		} else {
			assert.Nil(t, block)
			assert.Nil(t, body)
			assert.Nil(t, receipts)
		}
	}

	// generate new blocks
	curBlock := chain.CurrentBlock()
	expectedNBlocks := uint64(30)
	nBlocks := int(expectedNBlocks - curBlock.NumberU64())
	nTxs := 3
	newBlocks, _ := blockchain.GenerateChain(gspec.Config, curBlock, gxhash.NewFaker(), db, nBlocks, func(i int, b *blockchain.BlockGen) {
		for range nTxs {
			tx, err := types.SignTx(types.NewTransaction(b.TxNonce(address), generateAddress(), generateBigInt(10), params.TxGas, nil, nil), signer, key)
			assert.Nil(t, err)
			b.AddTx(tx)
		}
	})
	if _, err := chain.InsertChain(newBlocks); err != nil {
		t.Fatalf("Failed to import canonical chain start: %v", err)
	}

	// immediately restart to avoid idle time
	mCompress.Stop()
	mCompress.Start()
	waitCompression(mCompress)

	// check next compression block number
	assertNextCompressNum(t, db, expectedNBlocks)

	// check caonnical blocks after inserting some blocks
	allBlocks := append(canonblocks[:tt.setheadBlock], newBlocks...)
	assert.Equal(t, len(allBlocks), int(expectedNBlocks))
	for i := 0; i < len(allBlocks); i++ {
		block := chain.GetHeader(allBlocks[i].Hash(), allBlocks[i].NumberU64())
		body := chain.GetBodyRLP(allBlocks[i].Hash())
		receipts := chain.GetReceiptsByBlockHash(allBlocks[i].Hash())
		assert.NotNil(t, block)
		assert.NotNil(t, body)
		assert.NotNil(t, receipts)
	}
}

func genBlocks(dbm database.DBManager, startNum, nBlocks int) *types.Header {
	var lastHeader *types.Header
	for i := range nBlocks {
		n := int64(startNum + i)
		h := genHeader()
		h.Number = big.NewInt(n)

		hn, hh := h.Number.Uint64(), h.Hash()
		dbm.WriteCanonicalHash(hh, hn)
		dbm.WriteHeader(h)
		b, r := genBody(100), genReceipts(100)
		dbm.WriteBody(hh, hn, b)
		dbm.WriteReceipts(hh, hn, *r)
		lastHeader = h
	}
	return lastHeader
}
