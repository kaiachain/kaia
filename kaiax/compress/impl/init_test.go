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

package impl

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialCompression(t *testing.T) {
	compressLogInterval = 500
	compressCompactionPeriod = 800
	defer func() {
		compressLogInterval = 102400
		compressCompactionPeriod = 86400 * 30
	}()
	var (
		c, answers = makeTestModule(t)
		currNum    = uint64(1000)
		endNum     = uint64(872) // last block that can be compressed (= currNum - retention)
	)

	for _, schema := range c.schemas {
		// 1. Schemas must have been initialized during module Init.
		assert.Equal(t, uint64(1), *readNextNum(schema))
	}

	// Start compression threads and wait for completion.
	sizeBefore := dirSize(t, c.DBM.(database.DBManager).GetDBConfig().Dir)
	c.Start()
	waitCompletion(t, c, endNum)
	c.Stop()
	sizeAfter := dirSize(t, c.DBM.(database.DBManager).GetDBConfig().Dir)
	t.Logf("sizeBefore: %d, sizeAfter: %d", sizeBefore, sizeAfter)
	// assert.Less(t, sizeAfter, sizeBefore) // compression & compaction must reduce the size -> not necessarily true with small data size

	for _, schema := range c.schemas {
		// 2. The nextNum must be updated.
		// Because recent blocks might not be enough to form a chunk, nextNum can be less than retention.
		//                                          <-retention->
		// [-chunk-][-chunk-][-chunk-][-chunk-]------------------
		// |                                   |    |           |
		// genesis                        nextNum endNum    currNum
		nextNum := *readNextNum(schema)
		assert.LessOrEqual(t, nextNum, endNum, schema.name())

		// 3. Uncompressed and compressed data must be correct.
		validateSchema(t, c, schema, nextNum, currNum, answers[schema.name()])
	}
}

func waitCompletion(t *testing.T, c *CompressModule, wantedEndNum uint64) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if c.compressEndNum() == wantedEndNum {
			return // success
		}
	}
	t.Fail()
}

func validateSchema(t *testing.T, c *CompressModule, schema ItemSchema, nextNum, currNum uint64, answers [][]byte) {
	// blocks [0] and [nextNum,currNum] must survive in	uncompressed database.
	data, _ := schema.uncompressedDb().Get(schema.uncompressedKey(0, c.DBM.ReadCanonicalHash(0)))
	assert.Equal(t, answers[0], data)

	for i := nextNum; i < currNum; i++ {
		data, _ := schema.uncompressedDb().Get(schema.uncompressedKey(i, c.DBM.ReadCanonicalHash(i)))
		assert.Equal(t, answers[i], data)
	}

	// blocks [1,nextNum) must be deleted and exist in compressed database.
	for i := uint64(1); i < nextNum; i++ {
		ok, _ := schema.compressedDb().Has(schema.uncompressedKey(i, c.DBM.ReadCanonicalHash(i)))
		assert.False(t, ok, "%s %d", schema.name(), i)
	}
	for i := uint64(1); i < nextNum; i++ {
		uncompressed, ok := c.findFromCompressed(schema, i, c.DBM.ReadCanonicalHash(i))
		assert.True(t, ok, "%s %d", schema.name(), i)
		assert.Equal(t, answers[i], uncompressed)
	}
}

// Create a test module with sample blockchain and database.
// The chain contains 1000 blocks, each with a random number of transactions.
func makeTestModule(t *testing.T) (*CompressModule, map[string][][]byte) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError) // silence GenerateChain
	defer log.EnableLogForTest(log.LvlCrit, log.LvlInfo)
	var (
		currNum = uint64(1000)
		rng     = rand.New(rand.NewSource(0x12345678)) // deterministic seed for reproducible tests

		dir = t.TempDir()
		dbc = &database.DBConfig{
			Dir:               dir,
			DBType:            database.LevelDB,
			LevelDBCacheSize:  32,
			PebbleDBCacheSize: 32,
			OpenFilesLimit:    32,
		}
		dbm = database.NewDBManager(dbc)

		engine = gxhash.NewFaker()
		config = &params.ChainConfig{
			ChainID:       big.NewInt(31337),
			DeriveShaImpl: 2,
		}
		signer  = types.LatestSignerForChainID(config.ChainID)
		key1, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		addr1   = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
		nonce1  = uint64(0)
		addr2   = common.HexToAddress("0xdead")
		alloc   = blockchain.GenesisAlloc{
			addr1: {Balance: big.NewInt(10000000000)},
			addr2: {Balance: big.NewInt(0)},
		}
	)

	genesis := blockchain.Genesis{Config: config, Alloc: alloc}
	genesis.MustCommit(dbm)
	chain, err := blockchain.NewBlockChain(dbm, nil, config, engine, vm.Config{})
	require.NoError(t, err)

	block0 := chain.CurrentBlock()
	blocks, _ := blockchain.GenerateChain(config, block0, engine, dbm, int(currNum), func(i int, b *blockchain.BlockGen) {
		count := rng.Intn(10)
		for j := 0; j < count; j++ {
			unsignedTx := types.NewTransaction(nonce1, addr2, common.Big1, 21000, common.Big0, nil)
			nonce1++
			signedTx, err := types.SignTx(unsignedTx, signer, key1)
			require.NoError(t, err)
			b.AddTx(signedTx)
		}
	})
	chain.InsertChain(blocks)
	require.Equal(t, uint64(1000), chain.CurrentBlock().NumberU64())

	c := NewCompressModule()
	c.Init(&InitOpts{
		Chain: chain,
		DBM:   dbm,
		CompressConfig: compress.CompressConfig{
			Retention:    128,
			ChunkItemCap: 100,
			ChunkByteCap: 10 * 1024,
		},
	})

	answers := map[string][][]byte{}
	for _, schema := range c.schemas {
		// Remember correct data.
		data := [][]byte{}
		for i := uint64(0); i <= currNum; i++ {
			hash := dbm.ReadCanonicalHash(i)
			uncompressed, _ := schema.uncompressedDb().Get(schema.uncompressedKey(i, hash))
			data = append(data, uncompressed)
		}
		answers[schema.name()] = data
	}
	return c, answers
}
