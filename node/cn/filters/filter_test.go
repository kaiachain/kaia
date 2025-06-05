// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/consensus/gxhash"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/event"
	"github.com/kaiachain/kaia/v2/networks/rpc"
	cn "github.com/kaiachain/kaia/v2/node/cn/filters/mock"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	addr1 = common.HexToAddress("111")
	addr2 = common.HexToAddress("222")
	addrs []common.Address
)

var (
	topic1 common.Hash
	topic2 common.Hash
	topics [][]common.Hash
)

var (
	begin = int64(12345)
	end   = int64(12345)
)

var header *types.Header

var someErr = errors.New("some error")

func init() {
	addrs = []common.Address{addr1, addr2}
	topics = [][]common.Hash{{topic1}, {topic2}}
	header = &types.Header{
		Number:     big.NewInt(int64(123)),
		BlockScore: big.NewInt(int64(1)),
		Extra:      addrs[0][:],
		Governance: addrs[0][:],
		Vote:       addrs[0][:],
	}
}

func genFilter(t *testing.T) (*gomock.Controller, *cn.MockBackend, *Filter) {
	mockCtrl := gomock.NewController(t)
	mockBackend := cn.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().BloomStatus().Return(uint64(123), uint64(321)).Times(1)
	newFilter := NewRangeFilter(mockBackend, begin, end, addrs, topics)
	return mockCtrl, mockBackend, newFilter
}

func TestFilter_New(t *testing.T) {
	mockCtrl, mockBackend, newFilter := genFilter(t)
	defer mockCtrl.Finish()

	assert.NotNil(t, newFilter)
	assert.Equal(t, mockBackend, newFilter.backend)
	assert.Equal(t, begin, newFilter.begin)
	assert.Equal(t, end, newFilter.end)
	assert.Equal(t, topics, newFilter.topics)
	assert.Equal(t, addrs, newFilter.addresses)
	assert.NotNil(t, newFilter.matcher)
}

func TestFilter_Logs(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.LatestBlockNumber).Times(1).Return(nil, errors.New("latest header not found"))
		logs, err := newFilter.Logs(ctx)
		assert.Nil(t, logs)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_unindexedLogs(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.BlockNumber(newFilter.begin)).Times(1).Return(nil, nil)
		logs, err := newFilter.unindexedLogs(ctx, uint64(newFilter.end))
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.BlockNumber(newFilter.begin)).Times(1).Return(header, nil)
		logs, err := newFilter.unindexedLogs(ctx, uint64(newFilter.end))
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_checkMatches(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().GetLogs(ctx, header.Hash()).Times(1).Return(nil, someErr)
		logs, err := newFilter.checkMatches(ctx, header)
		assert.Nil(t, logs)
		assert.Equal(t, someErr, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().GetLogs(ctx, header.Hash()).Times(1).Return(nil, nil)
		logs, err := newFilter.checkMatches(ctx, header)
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_bloomFilter(t *testing.T) {
	{
		assert.True(t, bloomFilter(types.Bloom{}, nil, nil))
	}
	{
		assert.False(t, bloomFilter(types.Bloom{}, nil, [][]common.Hash{{topic1}}))
	}
	{
		assert.False(t, bloomFilter(types.Bloom{}, []common.Address{addr1}, nil))
	}
}

func makeReceipt(addr common.Address) *types.Receipt {
	receipt := genReceipt(false, 0)
	receipt.Logs = []*types.Log{
		{Address: addr},
	}
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt
}

func BenchmarkFilters(b *testing.B) {
	var (
		db         = database.NewMemoryDBManager()
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1      = crypto.PubkeyToAddress(key1.PublicKey)
		addr2      = common.BytesToAddress([]byte("jeff"))
		addr3      = common.BytesToAddress([]byte("ethereum"))
		addr4      = common.BytesToAddress([]byte("random addresses please"))
	)
	defer db.Close()

	genesis := blockchain.GenesisBlockForTesting(db, addr1, big.NewInt(1000000))
	chain, receipts := blockchain.GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 100010, func(i int, gen *blockchain.BlockGen) {
		switch i {
		case 2403:
			receipt := makeReceipt(addr1)
			gen.AddUncheckedReceipt(receipt)
		case 1034:
			receipt := makeReceipt(addr2)
			gen.AddUncheckedReceipt(receipt)
		case 34:
			receipt := makeReceipt(addr3)
			gen.AddUncheckedReceipt(receipt)
		case 99999:
			receipt := makeReceipt(addr4)
			gen.AddUncheckedReceipt(receipt)

		}
	})
	for i, block := range chain {
		db.WriteBlock(block)
		db.WriteCanonicalHash(block.Hash(), block.NumberU64())
		db.WriteHeadBlockHash(block.Hash())
		db.WriteReceipts(block.Hash(), block.NumberU64(), receipts[i])
	}
	b.ResetTimer()

	filter := NewRangeFilter(backend, 0, -1, []common.Address{addr1, addr2, addr3, addr4}, nil)

	for i := 0; i < b.N; i++ {
		logs, _ := filter.Logs(context.Background())
		if len(logs) != 4 {
			b.Fatal("expected 4 logs, got", len(logs))
		}
	}
}

func genReceipt(failed bool, cumulativeGasUsed uint64) *types.Receipt {
	r := &types.Receipt{GasUsed: cumulativeGasUsed}
	if failed {
		r.Status = types.ReceiptStatusFailed
	} else {
		r.Status = types.ReceiptStatusSuccessful
	}
	return r
}

func TestFilters(t *testing.T) {
	var (
		db         = database.NewMemoryDBManager()
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig, nil, nil}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr       = crypto.PubkeyToAddress(key1.PublicKey)

		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
	)
	defer db.Close()

	genesis := blockchain.GenesisBlockForTesting(db, addr, big.NewInt(1000000))
	chain, receipts := blockchain.GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 1000, func(i int, gen *blockchain.BlockGen) {
		switch i {
		case 1:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash1},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), 1, big.NewInt(1), nil))
		case 2:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash2},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(2, common.HexToAddress("0x2"), big.NewInt(2), 2, big.NewInt(2), nil))

		case 998:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash3},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(998, common.HexToAddress("0x998"), big.NewInt(998), 998, big.NewInt(998), nil))
		case 999:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash4},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, big.NewInt(999), nil))
		}
	})
	for i, block := range chain {
		db.WriteBlock(block)
		db.WriteCanonicalHash(block.Hash(), block.NumberU64())
		db.WriteHeadBlockHash(block.Hash())
		db.WriteReceipts(block.Hash(), block.NumberU64(), receipts[i])
	}

	for i, tc := range []struct {
		f    *Filter
		want string
		err  string
	}{
		{
			f:    NewRangeFilter(backend, 0, int64(rpc.LatestBlockNumber), []common.Address{addr}, [][]common.Hash{{hash1, hash2, hash3, hash4}}),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696331"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696332"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696333"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696334"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f:    NewRangeFilter(backend, 900, 999, []common.Address{addr}, [][]common.Hash{{hash3}}),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696333"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f:    NewRangeFilter(backend, 990, int64(rpc.LatestBlockNumber), []common.Address{addr}, [][]common.Hash{{hash3}}),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696333"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f:    NewRangeFilter(backend, 1, 10, []common.Address{addr}, [][]common.Hash{{hash1, hash2}}),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696331"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696332"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f:    NewRangeFilter(backend, 1, 10, nil, [][]common.Hash{{hash1, hash2}}),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696331"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696332"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f: NewRangeFilter(backend, 0, int64(rpc.LatestBlockNumber), nil, [][]common.Hash{{common.BytesToHash([]byte("fail"))}}),
		},
		{
			f: NewRangeFilter(backend, 0, int64(rpc.LatestBlockNumber), []common.Address{common.BytesToAddress([]byte("failmenow"))}, nil),
		},
		{
			f: NewRangeFilter(backend, 0, int64(rpc.LatestBlockNumber), nil, [][]common.Hash{{common.BytesToHash([]byte("fail"))}, {hash1}}),
		},
		{
			f:    NewRangeFilter(backend, int64(rpc.LatestBlockNumber), int64(rpc.LatestBlockNumber), nil, nil),
			want: `[{"address":"0x71562b71999873db5b286df957af199ec94617f7","topics":["0x0000000000000000000000000000000000000000000000000000746f70696334"],"data":"0x","blockNumber":"0x0","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false}]`,
		},
		{
			f:   NewRangeFilter(backend, int64(rpc.PendingBlockNumber), int64(rpc.PendingBlockNumber), nil, nil),
			err: errPendingLogsUnsupported.Error(),
		},
		{
			f:   NewRangeFilter(backend, int64(rpc.LatestBlockNumber), int64(rpc.PendingBlockNumber), nil, nil),
			err: errPendingLogsUnsupported.Error(),
		},
		{
			f:   NewRangeFilter(backend, int64(rpc.PendingBlockNumber), int64(rpc.LatestBlockNumber), nil, nil),
			err: errPendingLogsUnsupported.Error(),
		},
	} {
		logs, err := tc.f.Logs(context.Background())
		if err == nil && tc.err != "" {
			t.Fatalf("test %d, expected error %q, got nil", i, tc.err)
		} else if err != nil && err.Error() != tc.err {
			t.Fatalf("test %d, expected error %q, got %q", i, tc.err, err.Error())
		}
		if tc.want == "" && len(logs) == 0 {
			continue
		}
		have, err := json.Marshal(logs)
		if err != nil {
			t.Fatal(err)
		}
		if string(have) != tc.want {
			t.Fatalf("test %d, have:\n%s\nwant:\n%s", i, have, tc.want)
		}
	}

	t.Run("timeout", func(t *testing.T) {
		f := NewRangeFilter(backend, 0, rpc.LatestBlockNumber.Int64(), nil, nil)
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
		defer cancel()
		_, err := f.Logs(ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != errors.New("query timeout exceeded").Error() {
			t.Fatalf("expected context.DeadlineExceeded, got %v", err)
		}
	})
}
