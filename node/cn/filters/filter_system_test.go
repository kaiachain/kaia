// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/filters/filter_system_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package filters

import (
	"context"
	"math/big"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/bloombits"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

type testBackend struct {
	mux             *event.TypeMux
	db              database.DBManager
	sections        uint64
	txFeed          *event.Feed
	rmLogsFeed      *event.Feed
	logsFeed        *event.Feed
	chainFeed       *event.Feed
	chainConfig     *params.ChainConfig
	pendingBlock    *types.Block
	pendingReceipts types.Receipts
}

/*
head := rawdb.ReadCanonicalHash(klay.chainDb, (section+1)*params.BloomBitsBlocks-1)
head := klay.chainDB.ReadCanonicalHash((section+1)*params.BloomBitsBlocks-1)

rawdb.ReadBloomBits(klay.chainDb, task.Bit, section, head)
klay.chainDB.ReadBloomBits(database.BloomBitsKey(task.Bit, section, head))
*/

func (b *testBackend) ChainDB() database.DBManager {
	return b.db
}

func (b *testBackend) EventMux() *event.TypeMux {
	return b.mux
}

func (b *testBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	var (
		hash common.Hash
		num  uint64
	)
	if blockNr == rpc.LatestBlockNumber {
		hash = b.db.ReadHeadBlockHash()
		number := b.db.ReadHeaderNumber(hash)
		if number == nil {
			return nil, nil
		}
		num = *number
	} else {
		num = uint64(blockNr)
		hash = b.db.ReadCanonicalHash(num)
	}
	return b.db.ReadHeader(hash, num), nil
}

func (b *testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	number := b.db.ReadHeaderNumber(hash)
	if number == nil {
		return nil, nil
	}
	return b.db.ReadHeader(hash, *number), nil
}

func (b *testBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	if number := b.db.ReadHeaderNumber(hash); number != nil {
		return b.db.ReadReceipts(hash, *number)
	}
	return nil
}

func (b *testBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := b.db.ReadHeaderNumber(hash)
	if number == nil {
		return nil, nil
	}
	receipts := b.db.ReadReceipts(hash, *number)

	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *testBackend) Pending() (*types.Block, types.Receipts, *state.StateDB) {
	return b.pendingBlock, b.pendingReceipts, nil
}

func (b *testBackend) SubscribeNewTxsEvent(ch chan<- blockchain.NewTxsEvent) event.Subscription {
	return b.txFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return b.rmLogsFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.logsFeed.Subscribe(ch)
}

func (b *testBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return b.chainFeed.Subscribe(ch)
}

func (b *testBackend) BloomStatus() (uint64, uint64) {
	return params.BloomBitsBlocks, b.sections
}

func (b *testBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	requests := make(chan chan *bloombits.Retrieval)

	go session.Multiplex(16, 0, requests)
	go func() {
		for {
			// Wait for a service request or a shutdown
			select {
			case <-ctx.Done():
				return

			case request := <-requests:
				task := <-request

				task.Bitsets = make([][]byte, len(task.Sections))
				for i, section := range task.Sections {
					if rand.Int()%4 != 0 { // Handle occasional missing deliveries
						head := b.db.ReadCanonicalHash((section+1)*params.BloomBitsBlocks - 1)
						task.Bitsets[i], _ = b.db.ReadBloomBits(database.BloomBitsKey(task.Bit, section, head))
					}
				}
				request <- task
			}
		}
	}()
}

func (b *testBackend) setPending(block *types.Block, receipts types.Receipts) {
	b.pendingBlock = block
	b.pendingReceipts = receipts
}

func (b *testBackend) notifyPending(logs []*types.Log) {
	genesis := &blockchain.Genesis{
		Config: params.TestChainConfig,
	}
	db := database.NewMemoryDBManager()
	genesisBlock := genesis.MustCommit(db)
	blocks, _ := blockchain.GenerateChain(genesis.Config, genesisBlock, gxhash.NewFaker(), db, 2, func(i int, b *blockchain.BlockGen) {})
	b.setPending(blocks[1], []*types.Receipt{{Logs: logs}})
	b.chainFeed.Send(blockchain.ChainEvent{Block: blocks[0]})
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	return b.chainConfig
}

// TestBlockSubscription tests if a block subscription returns block hashes for posted chain events.
// It creates multiple subscriptions:
// - one at the start and should receive all posted chain events and a second (blockHashes)
// - one that is created after a cutoff moment and uninstalled after a second cutoff moment (blockHashes[cutoff1:cutoff2])
// - one that is created after the second cutoff moment (blockHashes[cutoff2:])
func TestBlockSubscription(t *testing.T) {
	t.Parallel()

	var (
		mux         = new(event.TypeMux)
		db          = database.NewMemoryDBManager()
		txFeed      = new(event.Feed)
		rmLogsFeed  = new(event.Feed)
		logsFeed    = new(event.Feed)
		chainFeed   = new(event.Feed)
		backend     = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api         = NewPublicFilterAPI(backend)
		genesis     = new(blockchain.Genesis).MustCommit(db)
		chain, _    = blockchain.GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 10, func(i int, gen *blockchain.BlockGen) {})
		chainEvents []blockchain.ChainEvent
	)

	for _, blk := range chain {
		chainEvents = append(chainEvents, blockchain.ChainEvent{Hash: blk.Hash(), Block: blk})
	}

	chan0 := make(chan *types.Header)
	sub0 := api.events.SubscribeNewHeads(chan0)
	chan1 := make(chan *types.Header)
	sub1 := api.events.SubscribeNewHeads(chan1)

	go func() { // simulate client
		i1, i2 := 0, 0
		for i1 != len(chainEvents) || i2 != len(chainEvents) {
			select {
			case header := <-chan0:
				if chainEvents[i1].Hash != header.Hash() {
					t.Errorf("sub0 received invalid hash on index %d, want %x, got %x", i1, chainEvents[i1].Hash, header.Hash())
				}
				i1++
			case header := <-chan1:
				if chainEvents[i2].Hash != header.Hash() {
					t.Errorf("sub1 received invalid hash on index %d, want %x, got %x", i2, chainEvents[i2].Hash, header.Hash())
				}
				i2++
			}
		}

		sub0.Unsubscribe()
		sub1.Unsubscribe()
	}()

	time.Sleep(1 * time.Second)
	for _, e := range chainEvents {
		chainFeed.Send(e)
	}

	<-sub0.Err()
	<-sub1.Err()
}

// TestPendingTxFilter tests whether pending tx filters retrieve all pending transactions that are posted to the event mux.
func TestPendingTxFilter(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api        = NewPublicFilterAPI(backend)

		transactions = []*types.Transaction{
			types.NewTransaction(0, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil),
			types.NewTransaction(1, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil),
			types.NewTransaction(2, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil),
			types.NewTransaction(3, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil),
			types.NewTransaction(4, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil),
		}

		hashes []common.Hash
	)

	fid0 := api.NewPendingTransactionFilter()

	time.Sleep(1 * time.Second)
	txFeed.Send(blockchain.NewTxsEvent{Txs: transactions})

	timeout := time.Now().Add(1 * time.Second)
	for {
		results, err := api.GetFilterChanges(fid0)
		if err != nil {
			t.Fatalf("Unable to retrieve logs: %v", err)
		}

		h := results.([]common.Hash)
		hashes = append(hashes, h...)
		if len(hashes) >= len(transactions) {
			break
		}
		// check timeout
		if time.Now().After(timeout) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if len(hashes) != len(transactions) {
		t.Errorf("invalid number of transactions, want %d transactions(s), got %d", len(transactions), len(hashes))
		return
	}
	for i := range hashes {
		if hashes[i] != transactions[i].Hash() {
			t.Errorf("hashes[%d] invalid, want %x, got %x", i, transactions[i].Hash(), hashes[i])
		}
	}
}

// TestLogFilterCreation test whether a given filter criteria makes sense.
// If not it must return an error.
func TestLogFilterCreation(t *testing.T) {
	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api        = NewPublicFilterAPI(backend)

		testCases = []struct {
			crit    FilterCriteria
			success bool
		}{
			// defaults
			{FilterCriteria{}, true},
			// valid block number range
			{FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(2)}, true},
			// "mined" block range to pending
			{FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, true},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(2), ToBlock: big.NewInt(1)}, false},
			// from block "higher" than to block
			{FilterCriteria{FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(100)}, false},
			// errPendingLogsUnsupported
			{FilterCriteria{FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(100)}, false},
			// errPendingLogsUnsupported
			{FilterCriteria{FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, false},
			// errPendingLogsUnsupported
			{FilterCriteria{FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(rpc.PendingBlockNumber.Int64())}, false},
			// topics more than 4 // NOTE: Kaia doesn't support errExceedMaxTopics
			{FilterCriteria{Topics: [][]common.Hash{{}, {}, {}, {}, {}}}, true},
		}
	)

	for i, test := range testCases {
		id, err := api.NewFilter(test.crit)
		if test.success && err != nil {
			t.Errorf("expected filter creation for case %d to success, got %v", i, err)
		}
		if err == nil {
			api.UninstallFilter(id)
			if !test.success {
				t.Errorf("expected testcase %d to fail with an error", i)
			}
		}
	}
}

// TestInvalidLogFilterCreation tests whether invalid filter log criteria results in an error
// when the filter is created.
func TestInvalidLogFilterCreation(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api        = NewPublicFilterAPI(backend)
	)

	// different situations where log filter creation should fail.
	// Reason: fromBlock > toBlock
	testCases := []FilterCriteria{
		0: {FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())},
		1: {FromBlock: big.NewInt(rpc.PendingBlockNumber.Int64()), ToBlock: big.NewInt(100)},
		2: {FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()), ToBlock: big.NewInt(100)},
	}

	for i, test := range testCases {
		if _, err := api.NewFilter(test); err == nil {
			t.Errorf("Expected NewFilter for case #%d to fail", i)
		}
	}
}

func TestInvalidGetLogsRequest(t *testing.T) {
	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api        = NewPublicFilterAPI(backend)
		blockHash  = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	)

	// Reason: Cannot specify both BlockHash and FromBlock/ToBlock)
	testCases := []FilterCriteria{
		0: {BlockHash: &blockHash, FromBlock: big.NewInt(100)},
		1: {BlockHash: &blockHash, ToBlock: big.NewInt(500)},
		2: {BlockHash: &blockHash, FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64())},
	}

	for i, test := range testCases {
		if _, err := api.GetLogs(context.Background(), test); err == nil {
			t.Errorf("Expected Logs for case #%d to fail", i)
		}
	}
}

// TestLogFilter tests whether log filters match the correct logs that are posted to the event feed.
func TestLogFilter(t *testing.T) {
	t.Parallel()

	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		api        = NewPublicFilterAPI(backend)

		firstAddr      = common.HexToAddress("0x1111111111111111111111111111111111111111")
		secondAddr     = common.HexToAddress("0x2222222222222222222222222222222222222222")
		thirdAddress   = common.HexToAddress("0x3333333333333333333333333333333333333333")
		notUsedAddress = common.HexToAddress("0x9999999999999999999999999999999999999999")
		firstTopic     = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		secondTopic    = common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")
		notUsedTopic   = common.HexToHash("0x9999999999999999999999999999999999999999999999999999999999999999")

		// posted twice, once as vm.Logs and once as blockchain.PendingLogsEvent
		allLogs = []*types.Log{
			{Address: firstAddr},
			{Address: firstAddr, Topics: []common.Hash{firstTopic}, BlockNumber: 1},
			{Address: secondAddr, Topics: []common.Hash{firstTopic}, BlockNumber: 1},
			{Address: thirdAddress, Topics: []common.Hash{secondTopic}, BlockNumber: 2},
			{Address: thirdAddress, Topics: []common.Hash{secondTopic}, BlockNumber: 3},
		}

		testCases = []struct {
			crit     FilterCriteria
			expected []*types.Log
			id       rpc.ID
		}{
			// match all
			0: {FilterCriteria{}, allLogs, ""},
			// match none due to no matching addresses
			1: {FilterCriteria{Addresses: []common.Address{{}, notUsedAddress}, Topics: [][]common.Hash{nil}}, []*types.Log{}, ""},
			// match logs based on addresses, ignore topics
			2: {FilterCriteria{Addresses: []common.Address{firstAddr}}, allLogs[:2], ""},
			// match none due to no matching topics (match with address)
			3: {FilterCriteria{Addresses: []common.Address{secondAddr}, Topics: [][]common.Hash{{notUsedTopic}}}, []*types.Log{}, ""},
			// match logs based on addresses and topics
			4: {FilterCriteria{Addresses: []common.Address{thirdAddress}, Topics: [][]common.Hash{{firstTopic, secondTopic}}}, allLogs[3:5], ""},
			// match logs based on multiple addresses and "or" topics
			5: {FilterCriteria{Addresses: []common.Address{secondAddr, thirdAddress}, Topics: [][]common.Hash{{firstTopic, secondTopic}}}, allLogs[2:5], ""},
			// all "mined" logs with block num >= 2
			6: {FilterCriteria{FromBlock: big.NewInt(2), ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, allLogs[3:], ""},
			// all "mined" logs
			7: {FilterCriteria{ToBlock: big.NewInt(rpc.LatestBlockNumber.Int64())}, allLogs, ""},
			// all "mined" logs with 1>= block num <=2 and topic secondTopic
			8: {FilterCriteria{FromBlock: big.NewInt(1), ToBlock: big.NewInt(2), Topics: [][]common.Hash{{secondTopic}}}, allLogs[3:4], ""},
			// match all logs due to wildcard topic
			9: {FilterCriteria{Topics: [][]common.Hash{nil}}, allLogs[1:], ""},
		}
	)

	// create all filters
	for i := range testCases {
		testCases[i].id, _ = api.NewFilter(testCases[i].crit)
	}

	// raise events
	time.Sleep(1 * time.Second)
	if nsend := logsFeed.Send(allLogs); nsend == 0 {
		t.Fatal("Shoud have at least one subscription")
	}

	// set pending logs
	backend.notifyPending(allLogs)

	for i, tt := range testCases {
		var fetched []*types.Log
		timeout := time.Now().Add(1 * time.Second)
		for { // fetch all expected logs
			results, err := api.GetFilterChanges(tt.id)
			if err != nil {
				t.Fatalf("test %d: unable to fetch logs: %v", i, err)
			}

			fetched = append(fetched, results.([]*types.Log)...)
			if len(fetched) >= len(tt.expected) {
				break
			}
			// check timeout
			if time.Now().After(timeout) {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}

		if len(fetched) != len(tt.expected) {
			t.Errorf("invalid number of logs for case %d, want %d log(s), got %d", i, len(tt.expected), len(fetched))
			return
		}

		for l := range fetched {
			if fetched[l].Removed {
				t.Errorf("expected log not to be removed for log %d in case %d", l, i)
			}
			if !reflect.DeepEqual(fetched[l], tt.expected[l]) {
				t.Errorf("invalid log on index %d for case %d", l, i)
			}
		}
	}
}

// TestPendingTxFilterDeadlock tests if the event loop hangs when pending
// txs arrive at the same time that one of multiple filters is timing out.
func TestPendingTxFilterDeadlock(t *testing.T) {
	t.Parallel()
	timeout := 100 * time.Millisecond

	var (
		mux        = new(event.TypeMux)
		db         = database.NewMemoryDBManager()
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		done       = make(chan struct{})
	)

	// instead of using NewPublicFilterAPI, define it directly
	// It is for faster test (timeout: 5minute -> 100millisecond)
	api := &PublicFilterAPI{
		backend: backend,
		mux:     backend.EventMux(),
		chainDB: backend.ChainDB(),
		events:  NewEventSystem(backend.EventMux(), backend),
		filters: make(map[rpc.ID]*filter),
		timeout: timeout,
	}
	go api.timeoutLoop()

	go func() {
		// Bombard feed with txs until signal was received to stop
		i := uint64(0)
		for {
			select {
			case <-done:
				return
			default:
			}

			tx := types.NewTransaction(i, common.HexToAddress("0xb794f5ea0ba39494ce83a213fffba74279579268"), new(big.Int), 0, new(big.Int), nil)
			backend.txFeed.Send(blockchain.NewTxsEvent{Txs: []*types.Transaction{tx}})
			i++
		}
	}()

	// Create a bunch of filters
	fids := make([]rpc.ID, 20)
	for i := 0; i < len(fids); i++ {
		fid := api.NewPendingTransactionFilter()
		fids[i] = fid
		// Wait for at least one tx to arrive in filter
		for {
			hashes, err := api.GetFilterChanges(fid)
			if err != nil {
				t.Fatalf("Filter should exist: %v\n", err)
			}
			if len(hashes.([]common.Hash)) > 0 {
				break
			}
			runtime.Gosched()
		}
	}

	// Wait until filters have timed out
	time.Sleep(3 * timeout)

	// If tx loop doesn't consume `done` after a second
	// it's hanging.
	select {
	case done <- struct{}{}:
		// Check that all filters have been uninstalled
		for _, fid := range fids {
			if _, err := api.GetFilterChanges(fid); err == nil {
				t.Errorf("Filter %s should have been uninstalled\n", fid)
			}
		}
	case <-time.After(1 * time.Second):
		t.Error("Tx sending loop hangs")
	}
}
