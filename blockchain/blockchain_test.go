// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/blockchain_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/compiler"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// So we can deterministically seed different blockchains
var (
	canonicalSeed = 1
	forkSeed      = 2
)

// newCanonical creates a chain database, and injects a deterministic canonical
// chain. Depending on the full flag, if creates either a full block chain or a
// header only chain.
func newCanonical(engine consensus.Engine, n int, full bool) (database.DBManager, *BlockChain, error) {
	var (
		db      = database.NewMemoryDBManager()
		genesis = new(Genesis).MustCommit(db)
	)

	// Initialize a fresh chain with only a genesis block
	blockchain, _ := NewBlockChain(db, nil, params.AllGxhashProtocolChanges, engine, vm.Config{})
	// Create and inject the requested chain
	if n == 0 {
		return db, blockchain, nil
	}
	if full {
		// Full block-chain requested
		blocks := MakeBlockChain(genesis, n, engine, db, canonicalSeed)
		_, err := blockchain.InsertChain(blocks)
		return db, blockchain, err
	}
	// Header-only chain requested
	headers := MakeHeaderChain(genesis.Header(), n, engine, db, canonicalSeed)
	_, err := blockchain.InsertHeaderChain(headers, 1)
	return db, blockchain, err
}

// Test fork of length N starting from block i
func testFork(t *testing.T, blockchain *BlockChain, i, n int, full bool, comparator func(td1, td2 *big.Int)) {
	// Copy old chain up to #i into a new db
	db, blockchain2, err := newCanonical(gxhash.NewFaker(), i, full)
	if err != nil {
		t.Fatal("could not make new canonical in testFork", err)
	}
	defer blockchain2.Stop()

	// Assert the chains have the same header/block at #i
	var hash1, hash2 common.Hash
	if full {
		hash1 = blockchain.GetBlockByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetBlockByNumber(uint64(i)).Hash()
	} else {
		hash1 = blockchain.GetHeaderByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetHeaderByNumber(uint64(i)).Hash()
	}
	if hash1 != hash2 {
		t.Errorf("chain content mismatch at %d: have hash %v, want hash %v", i, hash2, hash1)
	}
	// Extend the newly created chain
	var (
		blockChainB  []*types.Block
		headerChainB []*types.Header
	)
	if full {
		blockChainB = MakeBlockChain(blockchain2.CurrentBlock(), n, gxhash.NewFaker(), db, forkSeed)
		if _, err := blockchain2.InsertChain(blockChainB); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	} else {
		headerChainB = MakeHeaderChain(blockchain2.CurrentHeader(), n, gxhash.NewFaker(), db, forkSeed)
		if _, err := blockchain2.InsertHeaderChain(headerChainB, 1); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	}
	// Sanity check that the forked chain can be imported into the original
	var tdPre, tdPost *big.Int

	if full {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentBlock().Hash())
		if err := testBlockChainImport(blockChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked block chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(blockChainB[len(blockChainB)-1].Hash())
	} else {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentHeader().Hash())
		if err := testHeaderChainImport(headerChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked header chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(headerChainB[len(headerChainB)-1].Hash())
	}
	// Compare the total difficulties of the chains
	comparator(tdPre, tdPost)
}

func printChain(bc *BlockChain) {
	for i := bc.CurrentBlock().Number().Uint64(); i > 0; i-- {
		b := bc.GetBlockByNumber(uint64(i))
		fmt.Printf("\t%x %v\n", b.Hash(), b.BlockScore())
	}
}

// testBlockChainImport tries to process a chain of blocks, writing them into
// the database if successful.
func testBlockChainImport(chain types.Blocks, blockchain *BlockChain) error {
	for _, block := range chain {
		// Try and process the block
		err := blockchain.engine.VerifyHeader(blockchain, block.Header(), true)
		if err == nil {
			err = blockchain.validator.ValidateBody(block)
		}
		if err != nil {
			if err == ErrKnownBlock {
				continue
			}
			return err
		}
		statedb, err := state.New(blockchain.GetBlockByHash(block.ParentHash()).Root(), blockchain.stateCache, nil, nil)
		if err != nil {
			return err
		}
		receipts, _, usedGas, _, _, err := blockchain.Processor().Process(block, statedb, vm.Config{})
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		err = blockchain.validator.ValidateState(block, blockchain.GetBlockByHash(block.ParentHash()), statedb, receipts, usedGas)
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		blockchain.mu.Lock()
		blockchain.db.WriteTd(block.Hash(), block.NumberU64(), new(big.Int).Add(block.BlockScore(), blockchain.GetTdByHash(block.ParentHash())))
		blockchain.db.WriteBlock(block)
		statedb.Commit(false)
		blockchain.mu.Unlock()
	}
	return nil
}

// testHeaderChainImport tries to process a chain of header, writing them into
// the database if successful.
func testHeaderChainImport(chain []*types.Header, blockchain *BlockChain) error {
	for _, header := range chain {
		// Try and validate the header
		if err := blockchain.engine.VerifyHeader(blockchain, header, false); err != nil {
			return err
		}
		// Manually insert the header into the database, but don't reorganise (allows subsequent testing)
		blockchain.mu.Lock()
		blockchain.db.WriteTd(header.Hash(), header.Number.Uint64(), new(big.Int).Add(header.BlockScore, blockchain.GetTdByHash(header.ParentHash)))
		blockchain.db.WriteHeader(header)
		blockchain.mu.Unlock()
	}
	return nil
}

func insertChain(done chan bool, blockchain *BlockChain, chain types.Blocks, t *testing.T) {
	_, err := blockchain.InsertChain(chain)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	done <- true
}

func TestLastBlock(t *testing.T) {
	_, blockchain, err := newCanonical(gxhash.NewFaker(), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	blocks := MakeBlockChain(blockchain.CurrentBlock(), 1, gxhash.NewFullFaker(), blockchain.db, 0)
	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}
	if blocks[len(blocks)-1].Hash() != blockchain.db.ReadHeadBlockHash() {
		t.Fatalf("Write/Get HeadBlockHash failed")
	}
}

// Tests that given a starting canonical chain of a given size, it can be extended
// with various length chains.
func TestExtendCanonicalHeaders(t *testing.T) { testExtendCanonical(t, false) }
func TestExtendCanonicalBlocks(t *testing.T)  { testExtendCanonical(t, true) }

func testExtendCanonical(t *testing.T, full bool) {
	length := 5

	// Make first chain starting from genesis
	_, processor, err := newCanonical(gxhash.NewFaker(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the blockscore comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total blockscore mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Start fork from current height
	testFork(t, processor, length, 1, full, better)
	testFork(t, processor, length, 2, full, better)
	testFork(t, processor, length, 5, full, better)
	testFork(t, processor, length, 10, full, better)
}

// Tests that given a starting canonical chain of a given size, creating shorter
// forks do not take canonical ownership.
func TestShorterForkHeaders(t *testing.T) { testShorterFork(t, false) }
func TestShorterForkBlocks(t *testing.T)  { testShorterFork(t, true) }

func testShorterFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(gxhash.NewFaker(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the blockscore comparator
	worse := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) >= 0 {
			t.Errorf("total blockscore mismatch: have %v, expected less than %v", td2, td1)
		}
	}
	// Sum of numbers must be less than `length` for this to be a shorter fork
	testFork(t, processor, 0, 3, full, worse)
	testFork(t, processor, 0, 7, full, worse)
	testFork(t, processor, 1, 1, full, worse)
	testFork(t, processor, 1, 7, full, worse)
	testFork(t, processor, 5, 3, full, worse)
	testFork(t, processor, 5, 4, full, worse)
}

// Tests that given a starting canonical chain of a given size, creating longer
// forks do take canonical ownership.
func TestLongerForkHeaders(t *testing.T) { testLongerFork(t, false) }
func TestLongerForkBlocks(t *testing.T)  { testLongerFork(t, true) }

func testLongerFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(gxhash.NewFaker(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the blockscore comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total blockscore mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Sum of numbers must be greater than `length` for this to be a longer fork
	testFork(t, processor, 0, 11, full, better)
	testFork(t, processor, 0, 15, full, better)
	testFork(t, processor, 1, 10, full, better)
	testFork(t, processor, 1, 12, full, better)
	testFork(t, processor, 5, 6, full, better)
	testFork(t, processor, 5, 8, full, better)
}

// Tests that given a starting canonical chain of a given size, creating equal
// forks do take canonical ownership.
func TestEqualForkHeaders(t *testing.T) { testEqualFork(t, false) }
func TestEqualForkBlocks(t *testing.T)  { testEqualFork(t, true) }

func testEqualFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical(gxhash.NewFaker(), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the blockscore comparator
	equal := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) != 0 {
			t.Errorf("total blockscore mismatch: have %v, want %v", td2, td1)
		}
	}
	// Sum of numbers must be equal to `length` for this to be an equal fork
	testFork(t, processor, 0, 10, full, equal)
	testFork(t, processor, 1, 9, full, equal)
	testFork(t, processor, 2, 8, full, equal)
	testFork(t, processor, 5, 5, full, equal)
	testFork(t, processor, 6, 4, full, equal)
	testFork(t, processor, 9, 1, full, equal)
}

// Tests that chains missing links do not get accepted by the processor.
func TestBrokenHeaderChain(t *testing.T) { testBrokenChain(t, false) }
func TestBrokenBlockChain(t *testing.T)  { testBrokenChain(t, true) }

func testBrokenChain(t *testing.T, full bool) {
	// Make chain starting from genesis
	db, blockchain, err := newCanonical(gxhash.NewFaker(), 10, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a forked chain, and try to insert with a missing link
	if full {
		chain := MakeBlockChain(blockchain.CurrentBlock(), 5, gxhash.NewFaker(), db, forkSeed)[1:]
		if err := testBlockChainImport(chain, blockchain); err == nil {
			t.Errorf("broken block chain not reported")
		}
	} else {
		chain := MakeHeaderChain(blockchain.CurrentHeader(), 5, gxhash.NewFaker(), db, forkSeed)[1:]
		if err := testHeaderChainImport(chain, blockchain); err == nil {
			t.Errorf("broken header chain not reported")
		}
	}
}

// Tests that reorganising a long difficult chain after a short easy one
// overwrites the canonical numbers and links in the database.
func TestReorgLongHeaders(t *testing.T) { testReorgLong(t, false) }
func TestReorgLongBlocks(t *testing.T)  { testReorgLong(t, true) }

func testReorgLong(t *testing.T, full bool) {
	testReorg(t, []int64{0, 0, -9}, []int64{0, 0, 0, -9}, 393280, full)
}

// Tests that reorganising a short difficult chain after a long easy one
// overwrites the canonical numbers and links in the database.
func TestReorgShortHeaders(t *testing.T) { testReorgShort(t, false) }
func TestReorgShortBlocks(t *testing.T)  { testReorgShort(t, true) }

func testReorgShort(t *testing.T, full bool) {
	// Create a long easy chain vs. a short heavy one. Due to blockscore adjustment
	// we need a fairly long chain of blocks with different difficulties for a short
	// one to become heavyer than a long one. The 96 is an empirical value.
	easy := make([]int64, 96)
	for i := 0; i < len(easy); i++ {
		easy[i] = 60
	}
	diff := make([]int64, len(easy)-1)
	for i := 0; i < len(diff); i++ {
		diff[i] = -9
	}
	testReorg(t, easy, diff, 12615120, full)
}

func testReorg(t *testing.T, first, second []int64, td int64, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(gxhash.NewFaker(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Insert an easy and a difficult chain afterwards
	easyBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.CurrentBlock(), gxhash.NewFaker(), db, len(first), func(i int, b *BlockGen) {
		b.OffsetTime(first[i])
	})
	diffBlocks, _ := GenerateChain(params.TestChainConfig, blockchain.CurrentBlock(), gxhash.NewFaker(), db, len(second), func(i int, b *BlockGen) {
		b.OffsetTime(second[i])
	})
	if full {
		if _, err := blockchain.InsertChain(easyBlocks); err != nil {
			t.Fatalf("failed to insert easy chain: %v", err)
		}
		if _, err := blockchain.InsertChain(diffBlocks); err != nil {
			t.Fatalf("failed to insert difficult chain: %v", err)
		}
	} else {
		easyHeaders := make([]*types.Header, len(easyBlocks))
		for i, block := range easyBlocks {
			easyHeaders[i] = block.Header()
		}
		diffHeaders := make([]*types.Header, len(diffBlocks))
		for i, block := range diffBlocks {
			diffHeaders[i] = block.Header()
		}
		if _, err := blockchain.InsertHeaderChain(easyHeaders, 1); err != nil {
			t.Fatalf("failed to insert easy chain: %v", err)
		}
		if _, err := blockchain.InsertHeaderChain(diffHeaders, 1); err != nil {
			t.Fatalf("failed to insert difficult chain: %v", err)
		}
	}
	// Check that the chain is valid number and link wise
	if full {
		prev := blockchain.CurrentBlock()
		for block := blockchain.GetBlockByNumber(blockchain.CurrentBlock().NumberU64() - 1); block.NumberU64() != 0; prev, block = block, blockchain.GetBlockByNumber(block.NumberU64()-1) {
			if prev.ParentHash() != block.Hash() {
				t.Errorf("parent block hash mismatch: have %x, want %x", prev.ParentHash(), block.Hash())
			}
		}
	} else {
		prev := blockchain.CurrentHeader()
		for header := blockchain.GetHeaderByNumber(blockchain.CurrentHeader().Number.Uint64() - 1); header.Number.Uint64() != 0; prev, header = header, blockchain.GetHeaderByNumber(header.Number.Uint64()-1) {
			if prev.ParentHash != header.Hash() {
				t.Errorf("parent header hash mismatch: have %x, want %x", prev.ParentHash, header.Hash())
			}
		}
	}
	// Make sure the chain total blockscore is the correct one
	want := new(big.Int).Add(blockchain.genesisBlock.BlockScore(), big.NewInt(td))
	if full {
		if have := blockchain.GetTdByHash(blockchain.CurrentBlock().Hash()); have.Cmp(want) != 0 {
			t.Errorf("total blockscore mismatch: have %v, want %v", have, want)
		}
	} else {
		if have := blockchain.GetTdByHash(blockchain.CurrentHeader().Hash()); have.Cmp(want) != 0 {
			t.Errorf("total blockscore mismatch: have %v, want %v", have, want)
		}
	}
}

// Tests that the insertion functions detect banned hashes.
func TestBadHeaderHashes(t *testing.T) { testBadHashes(t, false) }
func TestBadBlockHashes(t *testing.T)  { testBadHashes(t, true) }

func testBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(gxhash.NewFaker(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a chain, ban a hash and try to import
	if full {
		blocks := MakeBlockChain(blockchain.CurrentBlock(), 3, gxhash.NewFaker(), db, 10)

		BadHashes[blocks[2].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[2].Header().Hash()) }()

		_, err = blockchain.InsertChain(blocks)
	} else {
		headers := MakeHeaderChain(blockchain.CurrentHeader(), 3, gxhash.NewFaker(), db, 10)

		BadHashes[headers[2].Hash()] = true
		defer func() { delete(BadHashes, headers[2].Hash()) }()

		_, err = blockchain.InsertHeaderChain(headers, 1)
	}
	if err != ErrBlacklistedHash {
		t.Errorf("error mismatch: have: %v, want: %v", err, ErrBlacklistedHash)
	}
}

// Tests that bad hashes are detected on boot, and the chain rolled back to a
// good state prior to the bad hash.
func TestReorgBadHeaderHashes(t *testing.T) { testReorgBadHashes(t, false) }
func TestReorgBadBlockHashes(t *testing.T)  { testReorgBadHashes(t, true) }

func testReorgBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical(gxhash.NewFaker(), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	blockchain.Config().Istanbul = params.GetDefaultIstanbulConfig()
	// Create a chain, import and ban afterwards
	headers := MakeHeaderChain(blockchain.CurrentHeader(), 4, gxhash.NewFaker(), db, 10)
	blocks := MakeBlockChain(blockchain.CurrentBlock(), 4, gxhash.NewFaker(), db, 10)

	if full {
		if _, err = blockchain.InsertChain(blocks); err != nil {
			t.Errorf("failed to import blocks: %v", err)
		}
		if blockchain.CurrentBlock().Hash() != blocks[3].Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", blockchain.CurrentBlock().Hash(), blocks[3].Header().Hash())
		}
		BadHashes[blocks[3].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[3].Header().Hash()) }()
	} else {
		if _, err = blockchain.InsertHeaderChain(headers, 1); err != nil {
			t.Errorf("failed to import headers: %v", err)
		}
		if blockchain.CurrentHeader().Hash() != headers[3].Hash() {
			t.Errorf("last header hash mismatch: have: %x, want %x", blockchain.CurrentHeader().Hash(), headers[3].Hash())
		}
		BadHashes[headers[3].Hash()] = true
		defer func() { delete(BadHashes, headers[3].Hash()) }()
	}
	blockchain.Stop()

	// Create a new BlockChain and check that it rolled back the state.
	ncm, err := NewBlockChain(blockchain.db, nil, blockchain.chainConfig, gxhash.NewFaker(), vm.Config{})
	if err != nil {
		t.Fatalf("failed to create new chain manager: %v", err)
	}
	if full {
		if ncm.CurrentBlock().Hash() != blocks[2].Header().Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", ncm.CurrentBlock().Hash(), blocks[2].Header().Hash())
		}
	} else {
		if ncm.CurrentHeader().Hash() != headers[0].ParentHash {
			t.Errorf("last header hash mismatch: have: %x, want %x", ncm.CurrentHeader().Hash(), headers[0].ParentHash)
		}
	}
	ncm.Stop()
}

// Tests chain insertions in the face of one entity containing an invalid nonce.
func TestHeadersInsertNonceError(t *testing.T) { testInsertNonceError(t, false) }
func TestBlocksInsertNonceError(t *testing.T)  { testInsertNonceError(t, true) }

func testInsertNonceError(t *testing.T, full bool) {
	for i := 1; i < 25 && !t.Failed(); i++ {
		// Create a pristine chain and database
		db, blockchain, err := newCanonical(gxhash.NewFaker(), 0, full)
		if err != nil {
			t.Fatalf("failed to create pristine chain: %v", err)
		}
		defer blockchain.Stop()

		// Create and insert a chain with a failing nonce
		var (
			failAt  int
			failRes int
			failNum uint64
		)
		if full {
			blocks := MakeBlockChain(blockchain.CurrentBlock(), i, gxhash.NewFaker(), db, 0)

			failAt = rand.Int() % len(blocks)
			failNum = blocks[failAt].NumberU64()

			blockchain.engine = gxhash.NewFakeFailer(failNum)
			failRes, err = blockchain.InsertChain(blocks)
		} else {
			headers := MakeHeaderChain(blockchain.CurrentHeader(), i, gxhash.NewFaker(), db, 0)

			failAt = rand.Int() % len(headers)
			failNum = headers[failAt].Number.Uint64()

			blockchain.engine = gxhash.NewFakeFailer(failNum)
			blockchain.hc.engine = blockchain.engine
			failRes, err = blockchain.InsertHeaderChain(headers, 1)
		}
		// Check that the returned error indicates the failure.
		if failRes != failAt {
			t.Errorf("test %d: failure index mismatch: have %d, want %d", i, failRes, failAt)
		}
		// Check that all no blocks after the failing block have been inserted.
		for j := 0; j < i-failAt; j++ {
			if full {
				if block := blockchain.GetBlockByNumber(failNum + uint64(j)); block != nil {
					t.Errorf("test %d: invalid block in chain: %v", i, block)
				}
			} else {
				if header := blockchain.GetHeaderByNumber(failNum + uint64(j)); header != nil {
					t.Errorf("test %d: invalid header in chain: %v", i, header)
				}
			}
		}
	}
}

// Tests that fast importing a block chain produces the same chain data as the
// classical full block processing.
func TestFastVsFullChains(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = database.NewMemoryDBManager()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(gendb)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
	)
	blocks, receipts := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), gendb, 1024, func(i int, block *BlockGen) {
		// If the block number is multiple of 3, send a few bonus transactions to the miner
		if i%3 == 2 {
			for j := 0; j < i%4+1; j++ {
				tx, err := types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
				if err != nil {
					panic(err)
				}
				block.AddTx(tx)
			}
		}
	})
	// Import the chain as an archive node for the comparison baseline
	archiveDb := database.NewMemoryDBManager()
	gspec.MustCommit(archiveDb)
	archive, _ := NewBlockChain(archiveDb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer archive.Stop()

	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	// Fast import the chain as a non-archive node to test
	fastDb := database.NewMemoryDBManager()
	gspec.MustCommit(fastDb)
	fast, _ := NewBlockChain(fastDb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer fast.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	// Iterate over all chain data components, and cross reference
	for i := 0; i < len(blocks); i++ {
		bnum, num, hash := blocks[i].Number(), blocks[i].NumberU64(), blocks[i].Hash()

		if ftd, atd := fast.GetTdByHash(hash), archive.GetTdByHash(hash); ftd.Cmp(atd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: have %v, want %v", num, hash, ftd, atd)
		}
		if fheader, aheader := fast.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); fheader.Hash() != aheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: have %v, want %v", num, hash, fheader, aheader)
		}
		if fblock, ablock := fast.GetBlockByHash(hash), archive.GetBlockByHash(hash); fblock.Hash() != ablock.Hash() {
			t.Errorf("block #%d [%x]: block mismatch: have %v, want %v", num, hash, fblock, ablock)
		} else if types.DeriveSha(fblock.Transactions(), bnum) != types.DeriveSha(ablock.Transactions(), bnum) {
			t.Errorf("block #%d [%x]: transactions mismatch: have %v, want %v", num, hash, fblock.Transactions(), ablock.Transactions())
		}
		freceipts := fastDb.ReadReceipts(hash, *fastDb.ReadHeaderNumber(hash))
		areceipts := archiveDb.ReadReceipts(hash, *archiveDb.ReadHeaderNumber(hash))
		if types.DeriveSha(freceipts, bnum) != types.DeriveSha(areceipts, bnum) {
			t.Errorf("block #%d [%x]: receipts mismatch: have %v, want %v", num, hash, freceipts, areceipts)
		}
	}
	// Check that the canonical chains are the same between the databases
	for i := 0; i < len(blocks)+1; i++ {
		if fhash, ahash := fastDb.ReadCanonicalHash(uint64(i)), archiveDb.ReadCanonicalHash(uint64(i)); fhash != ahash {
			t.Errorf("block #%d: canonical hash mismatch: have %v, want %v", i, fhash, ahash)
		}
	}
}

// Tests that various import methods move the chain head pointers to the correct
// positions.
func TestLightVsFastVsFullChainHeads(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = database.NewMemoryDBManager()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), gendb, int(height), nil)

	// Configure a subchain to roll back
	remove := []common.Hash{}
	for _, block := range blocks[height/2:] {
		remove = append(remove, block.Hash())
	}
	// Create a small assertion method to check the three heads
	assert := func(t *testing.T, kind string, chain *BlockChain, header uint64, fast uint64, block uint64) {
		if num := chain.CurrentBlock().NumberU64(); num != block {
			t.Errorf("%s head block mismatch: have #%v, want #%v", kind, num, block)
		}
		if num := chain.CurrentFastBlock().NumberU64(); num != fast {
			t.Errorf("%s head fast-block mismatch: have #%v, want #%v", kind, num, fast)
		}
		if num := chain.CurrentHeader().Number.Uint64(); num != header {
			t.Errorf("%s head header mismatch: have #%v, want #%v", kind, num, header)
		}
	}
	// Import the chain as an archive node and ensure all pointers are updated
	archiveDb := database.NewMemoryDBManager()
	gspec.MustCommit(archiveDb)

	archive, _ := NewBlockChain(archiveDb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	defer archive.Stop()

	assert(t, "archive", archive, height, height, height)
	archive.Rollback(remove)
	assert(t, "archive", archive, height/2, height/2, height/2)

	// Import the chain as a non-archive node and ensure all pointers are updated
	fastDb := database.NewMemoryDBManager()
	gspec.MustCommit(fastDb)
	fast, _ := NewBlockChain(fastDb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer fast.Stop()

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	assert(t, "fast", fast, height, height, 0)
	fast.Rollback(remove)
	assert(t, "fast", fast, height/2, height/2, 0)

	// Import the chain as a light node and ensure all pointers are updated
	lightDb := database.NewMemoryDBManager()
	gspec.MustCommit(lightDb)

	light, _ := NewBlockChain(lightDb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	if n, err := light.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	defer light.Stop()

	assert(t, "light", light, height, 0, 0)
	light.Rollback(remove)
	assert(t, "light", light, height/2, 0, 0)
}

// Tests that chain reorganisations handle transaction removals and reinsertions.
func TestChainTxReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		addr3   = crypto.PubkeyToAddress(key3.PublicKey)
		db      = database.NewMemoryDBManager()
		gspec   = &Genesis{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				addr1: {Balance: big.NewInt(1000000)},
				addr2: {Balance: big.NewInt(1000000)},
				addr3: {Balance: big.NewInt(1000000)},
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
	)

	// Create two transactions shared between the chains:
	//  - postponed: transaction included at a later block in the forked chain
	//  - swapped: transaction included at the same block number in the forked chain
	postponed, _ := types.SignTx(types.NewTransaction(0, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)
	swapped, _ := types.SignTx(types.NewTransaction(1, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)

	// Create two transactions that will be dropped by the forked chain:
	//  - pastDrop: transaction dropped retroactively from a past block
	//  - freshDrop: transaction dropped exactly at the block where the reorg is detected
	var pastDrop, freshDrop *types.Transaction

	// Create three transactions that will be added in the forked chain:
	//  - pastAdd:   transaction added before the reorganization is detected
	//  - freshAdd:  transaction added at the exact block the reorg is detected
	//  - futureAdd: transaction added after the reorg has already finished
	var pastAdd, freshAdd, futureAdd *types.Transaction

	chain, _ := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, 3, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			pastDrop, _ = types.SignTx(types.NewTransaction(gen.TxNonce(addr2), addr2, big.NewInt(1000), params.TxGas, nil, nil), signer, key2)

			gen.AddTx(pastDrop)  // This transaction will be dropped in the fork from below the split point
			gen.AddTx(postponed) // This transaction will be postponed till block #3 in the fork

		case 2:
			freshDrop, _ = types.SignTx(types.NewTransaction(gen.TxNonce(addr2), addr2, big.NewInt(1000), params.TxGas, nil, nil), signer, key2)

			gen.AddTx(freshDrop) // This transaction will be dropped in the fork from exactly at the split point
			gen.AddTx(swapped)   // This transaction will be swapped out at the exact height

			gen.OffsetTime(9) // Lower the block blockscore to simulate a weaker chain
		}
	})
	// Import the chain. This runs all block validation rules.
	blockchain, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	if i, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert original chain[%d]: %v", i, err)
	}
	defer blockchain.Stop()

	// overwrite the old chain
	chain, _ = GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, 5, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			pastAdd, _ = types.SignTx(types.NewTransaction(gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(pastAdd) // This transaction needs to be injected during reorg

		case 2:
			gen.AddTx(postponed) // This transaction was postponed from block #1 in the original chain
			gen.AddTx(swapped)   // This transaction was swapped from the exact current spot in the original chain

			freshAdd, _ = types.SignTx(types.NewTransaction(gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(freshAdd) // This transaction will be added exactly at reorg time

		case 3:
			futureAdd, _ = types.SignTx(types.NewTransaction(gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(futureAdd) // This transaction will be added after a full reorg
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	// removed tx
	for i, tx := range (types.Transactions{pastDrop, freshDrop}) {
		if txn, _, _, _ := db.ReadTxAndLookupInfo(tx.Hash()); txn != nil {
			t.Errorf("drop %d: tx %v found while shouldn't have been", i, txn)
		}
		if rcpt, _, _, _ := db.ReadReceipt(tx.Hash()); rcpt != nil {
			t.Errorf("drop %d: receipt %v found while shouldn't have been", i, rcpt)
		}
	}
	// added tx
	for i, tx := range (types.Transactions{pastAdd, freshAdd, futureAdd}) {
		if txn, _, _, _ := db.ReadTxAndLookupInfo(tx.Hash()); txn == nil {
			t.Errorf("add %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := db.ReadReceipt(tx.Hash()); rcpt == nil {
			t.Errorf("add %d: expected receipt to be found", i)
		}
	}
	// shared tx
	for i, tx := range (types.Transactions{postponed, swapped}) {
		if txn, _, _, _ := db.ReadTxAndLookupInfo(tx.Hash()); txn == nil {
			t.Errorf("share %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := db.ReadReceipt(tx.Hash()); rcpt == nil {
			t.Errorf("share %d: expected receipt to be found", i)
		}
	}
}

func TestLogReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		db      = database.NewMemoryDBManager()
		// this code generates a log
		code    = common.Hex2Bytes("60606040525b7f24ec1d3ff24c2f6ff210738839dbc339cd45a5294d85c79361016243157aae7b60405180905060405180910390a15b600a8060416000396000f360606040526008565b00")
		gspec   = &Genesis{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}}}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
	)

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer blockchain.Stop()

	rmLogsCh := make(chan RemovedLogsEvent)
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)
	chain, _ := GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 2, func(i int, gen *BlockGen) {
		if i == 1 {
			tx, err := types.SignTx(types.NewContractCreation(gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), code), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	chain, _ = GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 3, func(i int, gen *BlockGen) {})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	timeout := time.NewTimer(1 * time.Second)
	select {
	case ev := <-rmLogsCh:
		if len(ev.Logs) == 0 {
			t.Error("expected logs")
		}
	case <-timeout.C:
		t.Fatal("Timeout. There is no RemovedLogsEvent has been sent.")
	}
}

func TestReorgSideEvent(t *testing.T) {
	var (
		db      = database.NewMemoryDBManager()
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		gspec   = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
	)

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer blockchain.Stop()

	chain, _ := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, 3, func(i int, gen *BlockGen) {})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	replacementBlocks, _ := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, 4, func(i int, gen *BlockGen) {
		tx, err := types.SignTx(types.NewContractCreation(gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), nil), signer, key1)
		if i == 2 {
			gen.OffsetTime(-9)
		}
		if err != nil {
			t.Fatalf("failed to create tx: %v", err)
		}
		gen.AddTx(tx)
	})
	chainSideCh := make(chan ChainSideEvent, 64)
	blockchain.SubscribeChainSideEvent(chainSideCh)
	if _, err := blockchain.InsertChain(replacementBlocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// first two block of the secondary chain are for a brief moment considered
	// side chains because up to that point the first one is considered the
	// heavier chain.
	expectedSideHashes := map[common.Hash]bool{
		replacementBlocks[0].Hash(): true,
		replacementBlocks[1].Hash(): true,
		chain[0].Hash():             true,
		chain[1].Hash():             true,
		chain[2].Hash():             true,
	}

	i := 0

	const timeoutDura = 10 * time.Second
	timeout := time.NewTimer(timeoutDura)
done:
	for {
		select {
		case ev := <-chainSideCh:
			block := ev.Block
			if _, ok := expectedSideHashes[block.Hash()]; !ok {
				t.Errorf("%d: didn't expect %x to be in side chain", i, block.Hash())
			}
			i++

			if i == len(expectedSideHashes) {
				timeout.Stop()

				break done
			}
			timeout.Reset(timeoutDura)

		case <-timeout.C:
			t.Fatal("Timeout. Possibly not all blocks were triggered for sideevent")
		}
	}

	// make sure no more events are fired
	select {
	case e := <-chainSideCh:
		t.Errorf("unexpected event fired: %v", e)
	case <-time.After(250 * time.Millisecond):
	}
}

// Tests if the canonical block can be fetched from the database during chain insertion.
func TestCanonicalBlockRetrieval(t *testing.T) {
	_, blockchain, err := newCanonical(gxhash.NewFaker(), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	chain, _ := GenerateChain(blockchain.chainConfig, blockchain.genesisBlock, gxhash.NewFaker(), blockchain.db, 10, func(i int, gen *BlockGen) {})

	var pend sync.WaitGroup
	pend.Add(len(chain))

	for i := range chain {
		go func(block *types.Block) {
			defer pend.Done()

			// try to retrieve a block by its canonical hash and see if the block data can be retrieved.
			for {
				ch := blockchain.db.ReadCanonicalHash(block.NumberU64())
				if ch == (common.Hash{}) {
					continue // busy wait for canonical hash to be written
				}
				if ch != block.Hash() {
					t.Errorf("unknown canonical hash, want %s, got %s", block.Hash().Hex(), ch.Hex())
					return
				}
				fb := blockchain.db.ReadBlock(ch, block.NumberU64())
				if fb == nil {
					t.Errorf("unable to retrieve block %d for canonical hash: %s", block.NumberU64(), ch.Hex())
					return
				}
				if fb.Hash() != block.Hash() {
					t.Errorf("invalid block hash for block %d, want %s, got %s", block.NumberU64(), block.Hash().Hex(), fb.Hash().Hex())
					return
				}
				return
			}
		}(chain[i])

		if _, err := blockchain.InsertChain(types.Blocks{chain[i]}); err != nil {
			t.Fatalf("failed to insert block %d: %v", i, err)
		}
	}
	pend.Wait()
}

func TestEIP155Transition(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db         = database.NewMemoryDBManager()
		key, _     = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address    = crypto.PubkeyToAddress(key.PublicKey)
		funds      = big.NewInt(1000000000)
		deleteAddr = common.Address{1}
		gspec      = &Genesis{
			Config: &params.ChainConfig{ChainID: big.NewInt(1)},
			Alloc:  GenesisAlloc{address: {Balance: funds}, deleteAddr: {Balance: new(big.Int)}},
		}
		genesis = gspec.MustCommit(db)
	)

	blockchain, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer blockchain.Stop()

	// generate an invalid chain id transaction
	config := &params.ChainConfig{ChainID: big.NewInt(2)}
	blocks, _ := GenerateChain(config, genesis, gxhash.NewFaker(), db, 4, func(i int, block *BlockGen) {
		var (
			tx      *types.Transaction
			err     error
			basicTx = func(signer types.Signer) (*types.Transaction, error) {
				return types.SignTx(types.NewTransaction(block.TxNonce(address), common.Address{}, new(big.Int), 21000, new(big.Int), nil), signer, key)
			}
		)
		if i == 0 {
			tx, err = basicTx(types.NewEIP155Signer(big.NewInt(2)))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		}
	})
	_, err := blockchain.InsertChain(blocks)
	assert.Equal(t, types.ErrSender(types.ErrInvalidChainId), err)
}

// TODO-Kaia-FailedTest Failed test. Enable this later.
/*
func TestEIP161AccountRemoval(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db      = database.NewMemoryDBManager()
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		theAddr = common.Address{1}
		gspec   = &Genesis{
			Config: &params.ChainConfig{
				ChainID:        big.NewInt(1),
			},
			Alloc: GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(db)
	)
	blockchain, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer blockchain.Stop()

	blocks, _ := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, 3, func(i int, block *BlockGen) {
		var (
			tx     *types.Transaction
			err    error
			signer = types.NewEIP155Signer(gspec.Config.ChainID)
		)
		switch i {
		case 0:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 1:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 2:
			tx, err = types.SignTx(types.NewTransaction(block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		}
		if err != nil {
			t.Fatal(err)
		}
		block.AddTx(tx)
	})
	// account must exist pre eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[0]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); !st.Exist(theAddr) {
		t.Error("expected account to exist")
	}

	// account needs to be deleted post eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[1]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}

	// account musn't be created post eip 161
	if _, err := blockchain.InsertChain(types.Blocks{blocks[2]}); err != nil {
		t.Fatal(err)
	}
	if st, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}
}
*/

// This is a regression test (i.e. as weird as it is, don't delete it ever), which
// tests that under weird reorg conditions the blockchain and its internal header-
// chain return the same latest block/header.
//
// https://github.com/ethereum/go-ethereum/pull/15941
func TestBlockchainHeaderchainReorgConsistency(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine := gxhash.NewFaker()

	db := database.NewMemoryDBManager()
	genesis := new(Genesis).MustCommit(db)
	blocks, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, 64, func(i int, b *BlockGen) { b.SetRewardbase(common.Address{1}) })

	// Generate a bunch of fork blocks, each side forking from the canonical chain
	forks := make([]*types.Block, len(blocks))
	for i := 0; i < len(forks); i++ {
		parent := genesis
		if i > 0 {
			parent = blocks[i-1]
		}
		fork, _ := GenerateChain(params.TestChainConfig, parent, engine, db, 1, func(i int, b *BlockGen) { b.SetRewardbase(common.Address{2}) })
		forks[i] = fork[0]
	}
	// Import the canonical and fork chain side by side, verifying the current block
	// and current header consistency
	diskdb := database.NewMemoryDBManager()
	new(Genesis).MustCommit(diskdb)

	chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	for i := 0; i < len(blocks); i++ {
		if _, err := chain.InsertChain(blocks[i : i+1]); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", i, err)
		}
		if chain.CurrentBlock().Hash() != chain.CurrentHeader().Hash() {
			t.Errorf("block %d: current block/header mismatch: block #%d [%x…], header #%d [%x…]", i, chain.CurrentBlock().Number(), chain.CurrentBlock().Hash().Bytes()[:4], chain.CurrentHeader().Number, chain.CurrentHeader().Hash().Bytes()[:4])
		}
		if _, err := chain.InsertChain(forks[i : i+1]); err != nil {
			t.Fatalf(" fork %d: failed to insert into chain: %v", i, err)
		}
		if chain.CurrentBlock().Hash() != chain.CurrentHeader().Hash() {
			t.Errorf(" fork %d: current block/header mismatch: block #%d [%x…], header #%d [%x…]", i, chain.CurrentBlock().Number(), chain.CurrentBlock().Hash().Bytes()[:4], chain.CurrentHeader().Number, chain.CurrentHeader().Hash().Bytes()[:4])
		}
	}
}

// Tests that importing small side forks doesn't leave junk in the trie database
// cache (which would eventually cause memory issues).
func TestTrieForkGC(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine := gxhash.NewFaker()

	db := database.NewMemoryDBManager()
	genesis := new(Genesis).MustCommit(db)
	blocks, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, 2*DefaultTriesInMemory, func(i int, b *BlockGen) { b.SetRewardbase(common.Address{1}) })

	// Generate a bunch of fork blocks, each side forking from the canonical chain
	forks := make([]*types.Block, len(blocks))
	for i := 0; i < len(forks); i++ {
		parent := genesis
		if i > 0 {
			parent = blocks[i-1]
		}
		fork, _ := GenerateChain(params.TestChainConfig, parent, engine, db, 1, func(i int, b *BlockGen) { b.SetRewardbase(common.Address{2}) })
		forks[i] = fork[0]
	}
	// Import the canonical and fork chain side by side, forcing the trie cache to cache both
	diskdb := database.NewMemoryDBManager()
	new(Genesis).MustCommit(diskdb)

	chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	for i := 0; i < len(blocks); i++ {
		if _, err := chain.InsertChain(blocks[i : i+1]); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", i, err)
		}
		if _, err := chain.InsertChain(forks[i : i+1]); err != nil {
			t.Fatalf("fork %d: failed to insert into chain: %v", i, err)
		}
	}
	// Dereference all the recent tries and ensure no past trie is left in
	for i := 0; i < DefaultTriesInMemory; i++ {
		chain.stateCache.TrieDB().Dereference(blocks[len(blocks)-1-i].Root())
		chain.stateCache.TrieDB().Dereference(forks[len(blocks)-1-i].Root())
	}
	if len(chain.stateCache.TrieDB().Nodes()) > 0 {
		t.Fatalf("stale tries still alive after garbase collection")
	}
}

// Tests that State pruning indeed deletes obsolete trie nodes.
func TestStatePruning(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)
	var (
		db      = database.NewMemoryDBManager()
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = common.HexToAddress("0xaaaa")

		gspec = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
		engine  = gxhash.NewFaker()

		// Latest `retention` blocks survive.
		// Blocks 1..7 are pruned, blocks 8..10 are kept.
		retention = uint64(3)
		numBlocks = 10
		pruneNum  = uint64(numBlocks) - retention
	)

	db.WritePruningEnabled()
	cacheConfig := &CacheConfig{
		ArchiveMode:          false,
		CacheSize:            512,
		BlockInterval:        2, // Write frequently to test pruning
		TriesInMemory:        DefaultTriesInMemory,
		LivePruningRetention: retention,
		TrieNodeCacheConfig:  statedb.GetEmptyTrieNodeCacheConfig(),
	}
	blockchain, _ := NewBlockChain(db, cacheConfig, gspec.Config, engine, vm.Config{})

	chain, _ := GenerateChain(gspec.Config, genesis, engine, db, numBlocks, func(i int, gen *BlockGen) {
		tx, _ := types.SignTx(types.NewTransaction(
			gen.TxNonce(addr1), addr2, common.Big1, 21000, common.Big1, nil), signer, key1)
		gen.AddTx(tx)
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}
	assert.Equal(t, uint64(numBlocks), blockchain.CurrentBlock().NumberU64())

	// Give some time for pruning loop to run
	time.Sleep(100 * time.Millisecond)

	// Note that even if trie nodes are deleted from disk (DiskDB),
	// they may still be cached in memory (TrieDB).
	//
	// Therefore reopen the blockchain from the DiskDB with a clean TrieDB.
	// This simulates the node program restart.
	blockchain.Stop()
	blockchain, _ = NewBlockChain(db, cacheConfig, gspec.Config, engine, vm.Config{})

	// Genesis block always survives
	state, err := blockchain.StateAt(genesis.Root())
	assert.Nil(t, err)
	assert.NotZero(t, state.GetBalance(addr1).Uint64())

	// Pruned blocks should be inaccessible.
	for num := uint64(1); num <= pruneNum; num++ {
		_, err := blockchain.StateAt(blockchain.GetBlockByNumber(num).Root())
		assert.IsType(t, &statedb.MissingNodeError{}, err, num)
	}

	// Recent unpruned blocks should be accessible.
	for num := pruneNum + 1; num < uint64(numBlocks); num++ {
		state, err := blockchain.StateAt(blockchain.GetBlockByNumber(num).Root())
		require.Nil(t, err, num)
		assert.NotZero(t, state.GetBalance(addr1).Uint64())
		assert.NotZero(t, state.GetBalance(addr2).Uint64())
	}
	blockchain.Stop()
}

// TODO-Kaia-FailedTest Failed test. Enable this later.
/*
// Tests that doing large reorgs works even if the state associated with the
// forking point is not available any more.
func TestLargeReorgTrieGC(t *testing.T) {
	// Generate the original common chain segment and the two competing forks
	engine := gxhash.NewFaker()

	db := database.NewMemoryDBManager()
	genesis := new(Genesis).MustCommit(db)

	shared, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, 64, func(i int, b *BlockGen) { b.SetCoinbase(common.Address{1}) })
	original, _ := GenerateChain(params.TestChainConfig, shared[len(shared)-1], engine, db, 2*DefaultTriesInMemory, func(i int, b *BlockGen) { b.SetCoinbase(common.Address{2}) })
	competitor, _ := GenerateChain(params.TestChainConfig, shared[len(shared)-1], engine, db, 2*DefaultTriesInMemory+1, func(i int, b *BlockGen) { b.SetCoinbase(common.Address{3}) })

	// Import the shared chain and the original canonical one
	diskdb := database.NewMemoryDBManager()
	new(Genesis).MustCommit(diskdb)

	chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if _, err := chain.InsertChain(shared); err != nil {
		t.Fatalf("failed to insert shared chain: %v", err)
	}
	if _, err := chain.InsertChain(original); err != nil {
		t.Fatalf("failed to insert shared chain: %v", err)
	}
	// Ensure that the state associated with the forking point is pruned away
	if node, _ := chain.stateCache.TrieDB().Node(shared[len(shared)-1].Root()); node != nil {
		t.Fatalf("common-but-old ancestor still cache")
	}
	// Import the competitor chain without exceeding the canonical's TD and ensure
	// we have not processed any of the blocks (protection against malicious blocks)
	if _, err := chain.InsertChain(competitor[:len(competitor)-2]); err != nil {
		t.Fatalf("failed to insert competitor chain: %v", err)
	}
	for i, block := range competitor[:len(competitor)-2] {
		if node, _ := chain.stateCache.TrieDB().Node(block.Root()); node != nil {
			t.Fatalf("competitor %d: low TD chain became processed", i)
		}
	}
	// Import the head of the competitor chain, triggering the reorg and ensure we
	// successfully reprocess all the stashed away blocks.
	if _, err := chain.InsertChain(competitor[len(competitor)-2:]); err != nil {
		t.Fatalf("failed to finalize competitor chain: %v", err)
	}
	for i, block := range competitor[:len(competitor)-DefaultTriesInMemory] {
		if node, _ := chain.stateCache.TrieDB().Node(block.Root()); node != nil {
			t.Fatalf("competitor %d: competing chain state missing", i)
		}
	}
}
*/

// Taken from go-ethereum core.TestEIP2718Transition
// https://github.com/ethereum/go-ethereum/blob/v1.12.2/core/blockchain_test.go#L3441
func TestAccessListTx(t *testing.T) {
	config := params.TestChainConfig.Copy()
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0
	config.ShanghaiCompatibleBlock = common.Big0
	config.CancunCompatibleBlock = common.Big0
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Governance.KIP71.LowerBoundBaseFee = 0
	var (
		contractAddr = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		engine       = gxhash.NewFaker()
		signer       = types.LatestSigner(config)

		// A sender who makes transactions, has some funds
		senderKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		senderAddr   = crypto.PubkeyToAddress(senderKey.PublicKey)
		senderNonce  = uint64(0)
		gspec        = &Genesis{
			Config: config,
			Alloc: GenesisAlloc{
				senderAddr: {Balance: big.NewInt(params.KAIA)},
				contractAddr: { // SLOAD 0x00 and 0x01
					Code: []byte{
						byte(vm.PC),
						byte(vm.PC),
						byte(vm.SLOAD),
						byte(vm.SLOAD),
					},
					Nonce:   0,
					Balance: big.NewInt(0),
				},
			},
		}
		db       = database.NewMemoryDBManager()
		block    = gspec.MustCommit(db)
		accesses = types.AccessList{{
			Address:     contractAddr,
			StorageKeys: []common.Hash{{0}},
		}}
	)

	// Import the canonical chain
	chain, err := NewBlockChain(db, nil, gspec.Config, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	defer chain.Stop()

	// helper function to insert a block with a transaction
	insertBlockWithTx := func(list types.AccessList) *types.Block {
		// Generate blocks
		blocks, _ := GenerateChain(gspec.Config, block, engine, db, 1, func(i int, b *BlockGen) {
			b.SetRewardbase(common.Address{1})

			// One transaction to 0xAAAA
			intrinsicGas, _ := types.IntrinsicGas([]byte{}, list, nil, false, gspec.Config.Rules(block.Number()))
			tx, _ := types.SignTx(types.NewMessage(senderAddr, &contractAddr, senderNonce, big.NewInt(0), 30000, big.NewInt(1), nil, nil, []byte{}, false, intrinsicGas, list, nil, nil), signer, senderKey)
			b.AddTx(tx)
		})
		if n, err := chain.InsertChain(blocks); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", n, err)
		}
		senderNonce++
		return chain.CurrentBlock()
	}

	block = insertBlockWithTx(accesses) // with AccessList

	// Expected gas is intrinsic + 2 * pc + hot load + cold load, since only one load is in the access list
	expected := params.TxGas + params.TxAccessListAddressGas + params.TxAccessListStorageKeyGas +
		vm.GasQuickStep*2 + params.WarmStorageReadCostEIP2929 + params.ColdSloadCostEIP2929
	if block.GasUsed() != expected {
		t.Fatalf("incorrect amount of gas spent: expected %d, got %d", expected, block.GasUsed())
	}

	block = insertBlockWithTx(nil) // without AccessList

	// Expected gas is intrinsic + 2 * pc + 2 * cold load, since no access list provided
	expected = params.TxGas + vm.GasQuickStep*2 + params.ColdSloadCostEIP2929*2
	if block.GasUsed() != expected {
		t.Fatalf("incorrect amount of gas spent: expected %d, got %d", expected, block.GasUsed())
	}
}

func TestEIP3651(t *testing.T) {
	var (
		aa     = params.AuthorAddressForTesting
		bb     = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		engine = gxhash.NewFaker()
		db     = database.NewMemoryDBManager()

		// A sender who makes transactions, has some funds
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		funds   = new(big.Int).Mul(common.Big1, big.NewInt(params.KAIA))
		gspec   = &Genesis{
			Config: params.MainnetChainConfig.Copy(),
			Alloc: GenesisAlloc{
				addr1: {Balance: funds},
				addr2: {Balance: funds},
				// The address 0xAAAA sloads 0x00 and 0x01
				aa: {
					Code: []byte{
						byte(vm.PC),
						byte(vm.PC),
						byte(vm.SLOAD),
						byte(vm.SLOAD),
					},
					Nonce:   0,
					Balance: big.NewInt(0),
				},
				// The address 0xBBBB calls 0xAAAA
				// delegatecall(gas, address, in_offset(argsOffset), in_size(argsSize), out_offset(retOffset))
				// bb.Code: execute delegatecall to the contract which address is same as coinbase(producer)
				bb: {
					Code: []byte{
						byte(vm.PUSH1), 0, // out size
						byte(vm.DUP1),   // out offset
						byte(vm.DUP1),   // out insize
						byte(vm.DUP1),   // in offset
						byte(vm.PUSH20), // address
						byte(0xc0), byte(0xea), byte(0x08), byte(0xa2),
						byte(0xd4), byte(0x04), byte(0xd3), byte(0x17),
						byte(0x2d), byte(0x2a), byte(0xdd), byte(0x29),
						byte(0xa4), byte(0x5b), byte(0xe5), byte(0x6d),
						byte(0xa4), byte(0x0e), byte(0x29), byte(0x49),
						byte(vm.GAS), // gas
						byte(vm.DELEGATECALL),
					},
					Nonce:   0,
					Balance: big.NewInt(0),
				},
			},
		}
	)
	gspec.Config.SetDefaults()
	gspec.Config.IstanbulCompatibleBlock = common.Big0
	gspec.Config.LondonCompatibleBlock = common.Big0
	gspec.Config.EthTxTypeCompatibleBlock = common.Big0
	gspec.Config.MagmaCompatibleBlock = common.Big0
	gspec.Config.KoreCompatibleBlock = common.Big0
	gspec.Config.ShanghaiCompatibleBlock = common.Big0

	signer := types.LatestSigner(gspec.Config)
	genesis := gspec.MustCommit(db)

	blocks, _ := GenerateChain(gspec.Config, genesis, engine, db, 1, func(i int, b *BlockGen) {
		// One transaction to Coinbase
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyTo:       bb,
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(20000000),
			types.TxValueKeyGasPrice: big.NewInt(750 * params.Gkei),
			types.TxValueKeyData:     []byte{},
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, values)
		assert.Equal(t, nil, err)

		tx, err = types.SignTx(tx, signer, key1)
		assert.Equal(t, nil, err)

		b.AddTx(tx)
	})
	chain, err := NewBlockChain(db, nil, gspec.Config, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	block := chain.GetBlockByNumber(1)

	// 1+2: Ensure access lists are accounted for via gas usage.
	innerGas := vm.GasQuickStep*2 + params.ColdSloadCostEIP2929*2
	expectedGas := params.TxGas + 5*vm.GasFastestStep + vm.GasQuickStep + 100 + innerGas // 100 because 0xaaaa is in access list
	if block.GasUsed() != expectedGas {
		t.Fatalf("incorrect amount of gas spent: expected %d, got %d", expectedGas, block.GasUsed())
	}

	state, _ := chain.State()

	// 3: Ensure that miner received only the mining fee (consensus is gxHash, so 3 KAIA is the total reward)
	actual := state.GetBalance(params.AuthorAddressForTesting)
	expected := gxhash.ByzantiumBlockReward
	if actual.Cmp(expected) != 0 {
		t.Fatalf("miner balance incorrect: expected %d, got %d", expected, actual)
	}

	// 4: Ensure the tx sender paid for the gasUsed * (block baseFee).
	actual = new(big.Int).Sub(funds, state.GetBalance(addr1))
	expected = new(big.Int).Mul(new(big.Int).SetUint64(block.GasUsed()), block.Header().BaseFee)
	if actual.Cmp(expected) != 0 {
		t.Fatalf("sender balance incorrect: expected %d, got %d", expected, actual)
	}
}

// Benchmarks large blocks with value transfers to non-existing accounts
func benchmarkLargeNumberOfValueToNonexisting(b *testing.B, numTxs, numBlocks int, recipientFn func(uint64) common.Address, dataFn func(uint64) []byte) {
	var (
		testBankKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
		bankFunds       = big.NewInt(100000000000000000)
		gspec           = Genesis{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				testBankAddress: {Balance: bankFunds},
				common.HexToAddress("0xc0de"): {
					Code:    []byte{0x60, 0x01, 0x50},
					Balance: big.NewInt(0),
				}, // push 1, pop
			},
		}
		signer = types.LatestSignerForChainID(gspec.Config.ChainID)
	)
	// Generate the original common chain segment and the two competing forks
	engine := gxhash.NewFaker()
	db := database.NewMemoryDBManager()
	genesis := gspec.MustCommit(db)

	blockGenerator := func(i int, block *BlockGen) {
		for txi := 0; txi < numTxs; txi++ {
			uniq := uint64(i*numTxs + txi)
			recipient := recipientFn(uniq)
			// recipient := common.BigToAddress(big.NewInt(0).SetUint64(1337 + uniq))
			tx, err := types.SignTx(types.NewTransaction(uniq, recipient, big.NewInt(1), params.TxGas, big.NewInt(1), nil), signer, testBankKey)
			if err != nil {
				b.Error(err)
			}
			block.AddTx(tx)
		}
	}

	shared, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, numBlocks, blockGenerator)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Import the shared chain and the original canonical one
		diskdb := database.NewMemoryDBManager()
		gspec.MustCommit(diskdb)

		chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
		if err != nil {
			b.Fatalf("failed to create tester chain: %v", err)
		}
		b.StartTimer()
		if _, err := chain.InsertChain(shared); err != nil {
			b.Fatalf("failed to insert shared chain: %v", err)
		}
		b.StopTimer()
		if got := chain.CurrentBlock().Transactions().Len(); got != numTxs*numBlocks {
			b.Fatalf("Transactions were not included, expected %d, got %d", (numTxs * numBlocks), got)
		}
	}
}

func BenchmarkBlockChain_1x1000ValueTransferToNonexisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337 + nonce))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

func BenchmarkBlockChain_1x1000ValueTransferToExisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

func BenchmarkBlockChain_1x1000Executions(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(0xc0de))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}

	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

// TestCheckBlockChainVersion tests the functionality of CheckBlockChainVersion function.
func TestCheckBlockChainVersion(t *testing.T) {
	memDB := database.NewMemoryDBManager()

	// 1. If DatabaseVersion is not stored yet,
	// calling CheckBlockChainVersion stores BlockChainVersion to DatabaseVersion.
	assert.Nil(t, memDB.ReadDatabaseVersion())
	assert.NoError(t, CheckBlockChainVersion(memDB))
	assert.Equal(t, uint64(BlockChainVersion), *memDB.ReadDatabaseVersion())

	// 2. If DatabaseVersion is stored but less than BlockChainVersion,
	// calling CheckBlockChainVersion stores BlockChainVersion to DatabaseVersion.
	memDB.WriteDatabaseVersion(BlockChainVersion - 1)
	assert.NoError(t, CheckBlockChainVersion(memDB))
	assert.Equal(t, uint64(BlockChainVersion), *memDB.ReadDatabaseVersion())

	// 3. If DatabaseVersion is stored but greater than BlockChainVersion,
	// calling CheckBlockChainVersion returns an error and does not change the value.
	memDB.WriteDatabaseVersion(BlockChainVersion + 1)
	assert.Error(t, CheckBlockChainVersion(memDB))
	assert.Equal(t, uint64(BlockChainVersion+1), *memDB.ReadDatabaseVersion())
}

var (
	internalTxContractCode string
	internalTxContractAbi  string
)

func genInternalTxTransaction(t *testing.T, block *BlockGen, address common.Address, signer types.Signer, key *ecdsa.PrivateKey) {
	// 1. Deploy internal transaction sample contract
	nonce := block.TxNonce(address)
	amount := new(big.Int).SetUint64(100000000)
	gasLimit := big.NewInt(500000).Uint64()
	gasPrice := big.NewInt(25000)

	// It has to be cached because TestBlockChain_SetCanonicalBlock calls it many times
	if len(internalTxContractCode) == 0 {
		contracts, err := compiler.CompileSolidityOrLoad("",
			"../contracts/contracts/testing/internal_tx_contract/internal_tx_contract.sol")
		assert.Nil(t, err)

		var contract compiler.Contract
		for _, v := range contracts { // take the first one
			contract = *v
			break
		}

		internalTxContractCode = strings.TrimPrefix(contract.Code, "0x")
		abi, err := json.Marshal(contract.Info.AbiDefinition)
		assert.Nil(t, err)
		internalTxContractAbi = string(abi)
	}

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyAmount:        amount,
		types.TxValueKeyGasLimit:      gasLimit,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyFrom:          address,
		types.TxValueKeyData:          common.Hex2Bytes(internalTxContractCode),
	}
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
	assert.Nil(t, err)

	err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key})
	assert.Nil(t, err)

	block.AddTx(tx)

	contractAddr := crypto.CreateAddress(address, nonce)

	// 2. Contract Execution
	abii, err := abi.JSON(strings.NewReader(internalTxContractAbi))
	assert.Equal(t, nil, err)

	// the contract method "sendKlay" send 1 kei to address 3 times
	data, err := abii.Pack("sendKlay", uint32(3), address)
	assert.Equal(t, nil, err)

	values = map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    block.TxNonce(address),
		types.TxValueKeyFrom:     address,
		types.TxValueKeyTo:       contractAddr,
		types.TxValueKeyAmount:   big.NewInt(0),
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
	}
	tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key})
	assert.NoError(t, err)

	block.AddTx(tx)
}

// TestCallTraceChainEventSubscription tests if the method insertChain posts a chain event correctly.
// Scenario:
//  1. Deploy a contract
//     sendKlay(n uint32, receiver address): send 1 kei to `receiver` address `n` times.
//  2. Send a smart contract execution transaction
func TestCallTraceChainEventSubscription(t *testing.T) {
	// configure and generate a sample block chain
	var (
		gendb       = database.NewMemoryDBManager()
		key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address     = crypto.PubkeyToAddress(key.PublicKey)
		funds       = big.NewInt(100000000000000000)
		testGenesis = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = testGenesis.MustCommit(gendb)
		signer  = types.LatestSignerForChainID(testGenesis.Config.ChainID)
	)
	db := database.NewMemoryDBManager()
	testGenesis.MustCommit(db)

	// create new blockchain with enabled internal tx tracing option
	blockchain, _ := NewBlockChain(db, nil, testGenesis.Config, gxhash.NewFaker(), vm.Config{Debug: true, EnableInternalTxTracing: true})
	defer blockchain.Stop()

	// subscribe a new chain event channel
	chainEventCh := make(chan ChainEvent, 1)
	subscription := blockchain.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()

	// generate blocks
	blocks, _ := GenerateChain(testGenesis.Config, genesis, gxhash.NewFaker(), gendb, 1, func(i int, block *BlockGen) {
		// Deploy a contract which can trigger internal transactions
		genInternalTxTransaction(t, block, address, signer, key)
	})

	// insert the generated blocks into the test chain
	if n, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()
	// compare the published chain event with the expected test data
	select {
	case <-timer.C:
		t.Fatal("Timeout. There is no chain event posted for 1 second")
	case ev := <-chainEventCh:
		// a contract deploy tx and a contract execution tx
		assert.Equal(t, 2, len(ev.InternalTxTraces))

		// compare contract deploy result
		assert.Equal(t, address, *ev.InternalTxTraces[0].From)
		assert.Equal(t, 0, len(ev.InternalTxTraces[0].Calls))
		assert.Equal(t, "0x"+internalTxContractCode, ev.InternalTxTraces[0].Input)
		assert.Equal(t, fmt.Sprintf("0x%x", 100000000), ev.InternalTxTraces[0].Value)

		// compare contract execution result
		assert.Equal(t, address, *ev.InternalTxTraces[1].From)
		assert.Equal(t, 3, len(ev.InternalTxTraces[1].Calls))
		assert.Equal(t, fmt.Sprintf("0x%x", 0), ev.InternalTxTraces[1].Value)
	}
}

// TestBlockChain_SetCanonicalBlock tests SetCanonicalBlock.
// It first generates the chain and then call SetCanonicalBlock to change CurrentBlock.
func TestBlockChain_SetCanonicalBlock(t *testing.T) {
	// configure and generate a sample block chain
	var (
		gendb       = database.NewMemoryDBManager()
		key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address     = crypto.PubkeyToAddress(key.PublicKey)
		funds       = big.NewInt(100000000000000000)
		testGenesis = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = testGenesis.MustCommit(gendb)
		signer  = types.LatestSignerForChainID(testGenesis.Config.ChainID)
	)
	db := database.NewMemoryDBManager()
	testGenesis.MustCommit(db)

	// Archive mode is given to avoid mismatching between the given starting block number and
	// the actual block number where the blockchain has been rolled back to due to 128 blocks interval commit.
	cacheConfig := &CacheConfig{
		ArchiveMode:         true,
		CacheSize:           512,
		BlockInterval:       DefaultBlockInterval,
		TriesInMemory:       DefaultTriesInMemory,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:   512,
	}
	// create new blockchain with enabled internal tx tracing option
	blockchain, _ := NewBlockChain(db, cacheConfig, testGenesis.Config, gxhash.NewFaker(), vm.Config{Debug: true, EnableInternalTxTracing: true})
	defer blockchain.Stop()

	chainLength := rand.Int63n(500) + 100

	// generate blocks
	blocks, _ := GenerateChain(testGenesis.Config, genesis, gxhash.NewFaker(), gendb, int(chainLength), func(i int, block *BlockGen) {
		// Deploy a contract which can trigger internal transactions
		genInternalTxTransaction(t, block, address, signer, key)
	})

	// insert the generated blocks into the test chain
	if n, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}

	// target block number is 1/2 of the original chain length
	targetBlockNum := uint64(chainLength / 2)
	targetBlock := blockchain.db.ReadBlockByNumber(targetBlockNum)

	// set the canonical block with the target block number
	blockchain.SetCanonicalBlock(targetBlockNum)

	// compare the current block to the target block
	newHeadBlock := blockchain.CurrentBlock()
	assert.Equal(t, targetBlock.Hash(), newHeadBlock.Hash())
	assert.EqualValues(t, targetBlock, newHeadBlock)
}

func TestBlockChain_writeBlockLogsToRemoteCache(t *testing.T) {
	storage.SkipLocalTest(t)

	// prepare blockchain
	blockchain := &BlockChain{
		stateCache: state.NewDatabaseWithNewCache(database.NewMemoryDBManager(), &statedb.TrieNodeCacheConfig{
			CacheType:          statedb.CacheTypeHybrid,
			LocalCacheSizeMiB:  100,
			RedisEndpoints:     []string{"localhost:6379"},
			RedisClusterEnable: false,
		}),
	}

	// prepare test data to be written in the cache
	key := []byte{1, 2, 3, 4}
	log := &types.Log{
		Address:     common.Address{},
		Topics:      []common.Hash{common.BytesToHash(hexutil.MustDecode("0x123456789abcdef123456789abcdefffffffffff"))},
		Data:        []uint8{0x11, 0x22, 0x33, 0x44},
		BlockNumber: uint64(1000),
	}
	receipt := &types.Receipt{
		TxHash:  common.Hash{},
		GasUsed: uint64(999999),
		Status:  types.ReceiptStatusSuccessful,
		Logs:    []*types.Log{log},
	}

	// write log to cache
	blockchain.writeBlockLogsToRemoteCache(key, []*types.Receipt{receipt})

	// get log from cache
	ret := blockchain.stateCache.TrieDB().TrieNodeCache().Get(key)
	if ret == nil {
		t.Fatal("no cache")
	}

	// decode return data to the original log format
	storageLog := []*types.LogForStorage{}
	if err := rlp.DecodeBytes(ret, &storageLog); err != nil {
		t.Fatal(err)
	}
	logs := make([]*types.Log, len(storageLog))
	for i, log := range storageLog {
		logs[i] = (*types.Log)(log)
	}

	assert.Equal(t, log, logs[0])
}

// TestDeleteCreateRevert tests a weird state transition corner case that we hit
// while changing the internals of statedb. The workflow is that a contract is
// self destructed, then in a followup transaction (but same block) it's created
// again and the transaction reverted.
func TestDeleteCreateRevert(t *testing.T) {
	var (
		aa = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		bb = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		// Generate a canonical chain to act as the main dataset
		engine = gxhash.NewFaker()
		db     = database.NewMemoryDBManager()

		// A sender who makes transactions, has some funds
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				address: {Balance: funds},
				// The address 0xAAAAA selfdestructs if called
				aa: {
					// Code needs to just selfdestruct
					Code:    []byte{byte(vm.PC), 0xFF},
					Nonce:   1,
					Balance: big.NewInt(0),
				},
				// The address 0xBBBB send 1 kei to 0xAAAA, then reverts
				bb: {
					Code: []byte{
						byte(vm.PC),          // [0]
						byte(vm.DUP1),        // [0,0]
						byte(vm.DUP1),        // [0,0,0]
						byte(vm.DUP1),        // [0,0,0,0]
						byte(vm.PUSH1), 0x01, // [0,0,0,0,1] (value)
						byte(vm.PUSH2), 0xaa, 0xaa, // [0,0,0,0,1, 0xaaaa]
						byte(vm.GAS),
						byte(vm.CALL),
						byte(vm.REVERT),
					},
					Balance: big.NewInt(1),
				},
			},
		}
		genesis = gspec.MustCommit(db)
	)

	blocks, _ := GenerateChain(params.TestChainConfig, genesis, engine, db, 1, func(i int, b *BlockGen) {
		b.SetRewardbase(common.Address{1})
		signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
		// One transaction to AAAA
		tx, _ := types.SignTx(types.NewTransaction(0, aa,
			big.NewInt(0), 50000, big.NewInt(1), nil), signer, key)
		b.AddTx(tx)
		// One transaction to BBBB
		tx, _ = types.SignTx(types.NewTransaction(1, bb,
			big.NewInt(0), 100000, big.NewInt(1), nil), signer, key)
		b.AddTx(tx)
	})
	// Import the canonical chain
	diskdb := database.NewMemoryDBManager()
	gspec.MustCommit(diskdb)

	chain, err := NewBlockChain(diskdb, nil, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
}

// TestBlockChain_InsertChain_InsertFutureBlocks inserts future blocks that have a missing ancestor.
// It should return an expected error, but not panic.
func TestBlockChain_InsertChain_InsertFutureBlocks(t *testing.T) {
	// configure and generate a sample blockchain
	var (
		db          = database.NewMemoryDBManager()
		key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address     = crypto.PubkeyToAddress(key.PublicKey)
		funds       = big.NewInt(100000000000000000)
		testGenesis = &Genesis{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = testGenesis.MustCommit(db)
	)

	// Archive mode is given to avoid mismatching between the given starting block number and
	// the actual block number where the blockchain has been rolled back to due to 128 blocks interval commit.
	cacheConfig := &CacheConfig{
		ArchiveMode:         true,
		CacheSize:           512,
		BlockInterval:       DefaultBlockInterval,
		TriesInMemory:       DefaultTriesInMemory,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:   512,
	}
	cacheConfig.TrieNodeCacheConfig.NumFetcherPrefetchWorker = 3

	// create new blockchain with enabled internal tx tracing option
	blockchain, _ := NewBlockChain(db, cacheConfig, testGenesis.Config, gxhash.NewFaker(), vm.Config{})
	defer blockchain.Stop()

	// generate blocks
	blocks, _ := GenerateChain(testGenesis.Config, genesis, gxhash.NewFaker(), db, 10, func(i int, block *BlockGen) {})

	// insert the generated blocks into the test chain
	if n, err := blockchain.InsertChain(blocks[:2]); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}

	// insert future blocks
	_, err := blockchain.InsertChain(blocks[4:])
	if err == nil {
		t.Fatal("should be failed")
	}

	assert.Equal(t, consensus.ErrUnknownAncestor, err)
}

// TestTransientStorageReset ensures the transient storage is wiped correctly
// between transactions.
func TestTransientStorageReset(t *testing.T) {
	var (
		key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address     = crypto.PubkeyToAddress(key.PublicKey)
		destAddress = crypto.CreateAddress(address, 0)
		funds       = big.NewInt(1000000000000000000)

		testEngine = gxhash.NewFaker()
	)
	code := append([]byte{
		// TLoad value with location 1
		byte(vm.PUSH1), 0x1,
		byte(vm.TLOAD),

		// PUSH location
		byte(vm.PUSH1), 0x1,

		// SStore location:value
		byte(vm.SSTORE),
	}, make([]byte, 32-6)...)
	initCode := []byte{
		// TSTORE 1:1
		byte(vm.PUSH1), 0x1,
		byte(vm.PUSH1), 0x1,
		byte(vm.TSTORE),

		// Get the runtime-code on the stack
		byte(vm.PUSH32),
	}
	initCode = append(initCode, code...)
	initCode = append(initCode, []byte{
		byte(vm.PUSH1), 0x0, // offset
		byte(vm.MSTORE),
		byte(vm.PUSH1), 0x6, // size
		byte(vm.PUSH1), 0x0, // offset
		byte(vm.RETURN), // return 6 bytes of zero-code
	}...)

	gspec := &Genesis{
		Config: params.TestChainConfig.Copy(),
		Alloc: GenesisAlloc{
			address: {Balance: funds},
		},
	}
	gspec.Config.SetDefaults()
	gspec.Config.IstanbulCompatibleBlock = common.Big0
	gspec.Config.LondonCompatibleBlock = common.Big0
	gspec.Config.EthTxTypeCompatibleBlock = common.Big0
	gspec.Config.MagmaCompatibleBlock = common.Big0
	gspec.Config.KoreCompatibleBlock = common.Big0
	gspec.Config.ShanghaiCompatibleBlock = common.Big0
	gspec.Config.CancunCompatibleBlock = common.Big0
	gspec.Config.RandaoCompatibleBlock = common.Big0

	testdb := database.NewMemoryDBManager()
	genesis := gspec.MustCommit(testdb)
	blocks, _ := GenerateChain(gspec.Config, genesis, testEngine, testdb, 10, func(i int, gen *BlockGen) {
		fee := big.NewInt(1)
		if gen.header.BaseFee != nil {
			fee = gen.header.BaseFee
		}
		gen.SetRewardbase(common.Address{1})
		signer := types.LatestSigner(gen.config)
		tx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(address), common.Address{}, nil, 100000, fee, initCode), signer, key)
		gen.AddTx(tx)
		tx, _ = types.SignTx(types.NewTransaction(gen.TxNonce(address), destAddress, nil, 100000, fee, nil), signer, key)
		gen.AddTx(tx)
	})

	// Initialize the blockchain with 1153 enabled.
	testdb = database.NewMemoryDBManager()
	gspec.MustCommit(testdb)
	chain, err := NewBlockChain(testdb, nil, gspec.Config, testEngine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	defer chain.Stop()

	// Import the blocks
	if _, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to insert into chain: %v", err)
	}
	// Check the storage
	state, err := chain.StateAt(chain.CurrentHeader().Root)
	if err != nil {
		t.Fatalf("Failed to load state %v", err)
	}
	loc := common.BytesToHash([]byte{1})
	slot := state.GetState(destAddress, loc)
	if slot != (common.Hash{}) {
		t.Fatalf("Unexpected dirty storage slot")
	}
}

func TestProcessParentBlockHash(t *testing.T) {
	var (
		chainConfig = &params.ChainConfig{
			ShanghaiCompatibleBlock: common.Big0, // Shanghai fork is necesasry because `params.HistoryStorageCode` contains `PUSH0(0x5f)` instruction
		}
		hashA    = common.Hash{0x01}
		hashB    = common.Hash{0x02}
		header   = &types.Header{ParentHash: hashA, Number: big.NewInt(2), Time: common.Big0, BlockScore: common.Big0}
		parent   = &types.Header{ParentHash: hashB, Number: big.NewInt(1), Time: common.Big0, BlockScore: common.Big0}
		coinbase = common.Address{}
		rules    = params.Rules{}
		db       = state.NewDatabase(database.NewMemoryDBManager())
	)
	test := func(statedb *state.StateDB) {
		if err := statedb.SetCode(params.HistoryStorageAddress, params.HistoryStorageCode); err != nil {
			t.Error(err)
		}
		statedb.SetNonce(params.HistoryStorageAddress, 1)
		statedb.IntermediateRoot(true)

		vmContext := NewEVMBlockContext(header, nil, &coinbase)
		evm := vm.NewEVM(vmContext, vm.TxContext{}, statedb, chainConfig, &vm.Config{})
		if err := ProcessParentBlockHash(header, evm, statedb, rules); err != nil {
			t.Error(err)
		}

		vmContext = NewEVMBlockContext(parent, nil, &coinbase)
		evm = vm.NewEVM(vmContext, vm.TxContext{}, statedb, chainConfig, &vm.Config{})
		if err := ProcessParentBlockHash(parent, evm, statedb, rules); err != nil {
			t.Error(err)
		}

		// make sure that the state is correct
		if have := getParentBlockHash(statedb, 1); have != hashA {
			t.Errorf("want parent hash %v, have %v", hashA, have)
		}
		if have := getParentBlockHash(statedb, 0); have != hashB {
			t.Errorf("want parent hash %v, have %v", hashB, have)
		}
	}
	t.Run("MPT", func(t *testing.T) {
		statedb, _ := state.New(types.EmptyRootHash, db, nil, nil)
		test(statedb)
	})
}

func getParentBlockHash(statedb *state.StateDB, number uint64) common.Hash {
	ringIndex := number % params.HistoryServeWindow
	var key common.Hash
	binary.BigEndian.PutUint64(key[24:], ringIndex)
	return statedb.GetState(params.HistoryStorageAddress, key)
}

func newGkei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Gkei))
}

// TestEIP7702 deploys two delegation designations and calls them. It writes one
// value to storage which is verified after.
func TestEIP7702(t *testing.T) {
	var (
		// Generate a canonical chain to act as the main dataset
		engine  = gxhash.NewFaker()
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		aa      = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		bb      = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		funds   = big.NewInt(100000000000000000)
	)
	gspec := &Genesis{
		Config: params.TestChainConfig.Copy(),
		Alloc: GenesisAlloc{
			addr1: {Balance: funds},
			addr2: {Balance: funds},
			// The address 0xAAAA sstores 1 into slot 2.
			aa: {
				Code: []byte{
					byte(vm.PC),          // [0]
					byte(vm.DUP1),        // [0,0]
					byte(vm.DUP1),        // [0,0,0]
					byte(vm.DUP1),        // [0,0,0,0]
					byte(vm.PUSH1), 0x01, // [0,0,0,0,1] (value)
					byte(vm.PUSH20), addr2[0], addr2[1], addr2[2], addr2[3], addr2[4], addr2[5], addr2[6], addr2[7], addr2[8], addr2[9], addr2[10], addr2[11], addr2[12], addr2[13], addr2[14], addr2[15], addr2[16], addr2[17], addr2[18], addr2[19],
					byte(vm.GAS),
					byte(vm.CALL),
					byte(vm.STOP),
				},
				Nonce:   0,
				Balance: big.NewInt(0),
			},
			// The address 0xBBBB sstores 42 into slot 42.
			bb: {
				Code: []byte{
					byte(vm.PUSH1), 0x42,
					byte(vm.DUP1),
					byte(vm.SSTORE),
					byte(vm.STOP),
				},
				Nonce:   0,
				Balance: big.NewInt(0),
			},
		},
	}
	gspec.Config.SetDefaults()
	gspec.Config.IstanbulCompatibleBlock = common.Big0
	gspec.Config.LondonCompatibleBlock = common.Big0
	gspec.Config.EthTxTypeCompatibleBlock = common.Big0
	gspec.Config.MagmaCompatibleBlock = common.Big0
	gspec.Config.KoreCompatibleBlock = common.Big0
	gspec.Config.ShanghaiCompatibleBlock = common.Big0
	gspec.Config.CancunCompatibleBlock = common.Big0
	gspec.Config.KaiaCompatibleBlock = common.Big0
	gspec.Config.PragueCompatibleBlock = common.Big0

	// Sign authorization tuples.
	// The way the auths are combined, it becomes
	// 1. tx -> addr1 which is delegated to 0xaaaa
	// 2. addr1:0xaaaa calls into addr2:0xbbbb
	// 3. addr2:0xbbbb  writes to storage
	auth1, _ := types.SignSetCode(key1, types.SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(gspec.Config.ChainID),
		Address: aa,
		Nonce:   1,
	})

	auth2, _ := types.SignSetCode(key2, types.SetCodeAuthorization{
		ChainID: *uint256.NewInt(0),
		Address: bb,
		Nonce:   0,
	})

	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)

	testdb := database.NewMemoryDBManager()
	genesis := gspec.MustCommit(testdb)
	blocks, _ := GenerateChain(gspec.Config, genesis, engine, testdb, 1, func(i int, b *BlockGen) {
		b.SetRewardbase(common.Address{1})

		authorizationList := []types.SetCodeAuthorization{auth1, auth2}
		intrinsicGas, err := types.IntrinsicGas(nil, nil, authorizationList, false, params.TestRules)
		if err != nil {
			t.Fatalf("failed to run intrinsic gas: %v", err)
		}

		tx, err := types.SignTx(types.NewMessage(addr1, &addr1, uint64(0), nil, 500000, nil, newGkei(50),
			big.NewInt(20), nil, false, intrinsicGas, nil, nil, authorizationList), signer, key1)
		if err != nil {
			t.Fatalf("failed to sign tx: %v", err)
		}
		b.AddTx(tx)
	})

	chain, err := NewBlockChain(testdb, nil, gspec.Config, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	defer chain.Stop()
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	// Verify delegation designations were deployed.
	state, _ := chain.State()
	code, want := state.GetCode(addr1), types.AddressToDelegation(auth1.Address)
	if !bytes.Equal(code, want) {
		t.Fatalf("addr1 code incorrect: got %s, want %s", common.Bytes2Hex(code), common.Bytes2Hex(want))
	}
	if vmVersion, ok := state.GetVmVersion(addr1); !(vmVersion == params.VmVersion1 && ok) {
		t.Fatalf("addr1 code info incorrect: got %v, want %v", vmVersion, params.VmVersion1)
	}

	code, want = state.GetCode(addr2), types.AddressToDelegation(auth2.Address)
	if !bytes.Equal(code, want) {
		t.Fatalf("addr2 code incorrect: got %s, want %s", common.Bytes2Hex(code), common.Bytes2Hex(want))
	}
	if vmVersion, ok := state.GetVmVersion(addr2); !(vmVersion == params.VmVersion1 && ok) {
		t.Fatalf("addr2 code info incorrect: got %v, want %v", vmVersion, params.VmVersion1)
	}

	// Verify delegation executed the correct code.
	var (
		fortyTwo = common.BytesToHash([]byte{0x42})
		actual   = state.GetState(addr2, fortyTwo)
	)
	if !bytes.Equal(actual[:], fortyTwo[:]) {
		t.Fatalf("addr2 storage wrong: expected %d, got %d", fortyTwo, actual)
	}

	// Check if EOA with code can be called with kaia's execution type
	execTxTypes := []types.TxType{
		types.TxTypeSmartContractExecution,
		types.TxTypeFeeDelegatedSmartContractExecution,
		types.TxTypeFeeDelegatedSmartContractExecutionWithRatio,
	}
	for _, execTxType := range execTxTypes {
		tx, err := genTxForExecutionType(execTxType, addr1, addr1, state.GetNonce(addr1))
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key1})
		assert.Equal(t, nil, err)

		if execTxType.IsFeeDelegatedTransaction() {
			err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{key1})
			assert.Equal(t, nil, err)
		}

		// It only checks whether the transaction is successful.
		err = applyTransaction(chain, state, tx, addr1)
		assert.Equal(t, nil, err)
	}

	// Set 0x0000000000000000000000000000000000000000000 test
	{
		authForEmpty, _ := types.SignSetCode(key1, types.SetCodeAuthorization{
			ChainID: *uint256.MustFromBig(gspec.Config.ChainID),
			Address: common.Address{},
			Nonce:   state.GetNonce(addr1) + 1,
		})

		authorizationList := []types.SetCodeAuthorization{authForEmpty}
		intrinsicGas, err := types.IntrinsicGas(nil, nil, authorizationList, false, params.TestRules)
		if err != nil {
			t.Fatalf("failed to run intrinsic gas: %v", err)
		}

		tx, err := types.SignTx(types.NewMessage(addr1, &addr1, state.GetNonce(addr1), nil, 500000, nil, newGkei(50),
			big.NewInt(20), nil, false, intrinsicGas, nil, nil, authorizationList), signer, key1)
		if err != nil {
			t.Fatalf("failed to sign tx: %v", err)
		}
		err = applyTransaction(chain, state, tx, addr1)
		assert.Equal(t, nil, err)

		state.GetCode(addr1)
	}

	// Checks whether code and codeInfo are initialized.
	if !bytes.Equal(state.GetCode(addr1), nil) {
		t.Fatalf("addr1 code incorrect: got %s, want %s", common.Bytes2Hex(code), common.Bytes2Hex(want))
	}
	if vmVersion, ok := state.GetVmVersion(addr1); !(vmVersion == params.VmVersion0 && !ok) {
		t.Fatalf("addr1 code info incorrect: got %v, want %v", vmVersion, params.VmVersion1)
	}
}

func genTxForExecutionType(txType types.TxType, from, to common.Address, nonce uint64) (*types.Transaction, error) {
	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     from,
		types.TxValueKeyTo:       to,
		types.TxValueKeyAmount:   big.NewInt(0),
		types.TxValueKeyGasLimit: uint64(500000),
		types.TxValueKeyGasPrice: big.NewInt(25 * params.Gkei),
		types.TxValueKeyData:     []byte{},
	}

	if txType.IsFeeDelegatedTransaction() {
		values[types.TxValueKeyFeePayer] = from
	}

	if txType.IsFeeDelegatedWithRatioTransaction() {
		values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)
	}

	return types.NewTransactionWithMap(txType, values)
}

func applyTransaction(chain *BlockChain, state *state.StateDB, tx *types.Transaction, author common.Address) error {
	parent := chain.CurrentBlock()
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), common.Big1),
		BlockScore: big.NewInt(0),
	}
	usedGas := uint64(0)
	_, _, err := chain.ApplyTransaction(chain.Config(), &author, state, header, tx, &usedGas, &vm.Config{})
	return err
}
