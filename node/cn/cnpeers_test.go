// Copyright 2026 The Kaia Authors
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

package cn

import (
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cnPeerReadResult struct {
	addrs []common.Address
	err   error
}

type fakeCNPeerReader struct {
	results []cnPeerReadResult
	calls   []uint64
}

func (r *fakeCNPeerReader) GetCNPeers(num uint64) ([]common.Address, error) {
	r.calls = append(r.calls, num)
	if len(r.results) == 0 {
		return nil, nil
	}
	result := r.results[0]
	if len(r.results) > 1 {
		r.results = r.results[1:]
	}
	return result.addrs, result.err
}

type recordCNPeerSink struct {
	calls [][]common.Address
}

func (s *recordCNPeerSink) SetCNPeers(addrs []common.Address) {
	if addrs == nil {
		s.calls = append(s.calls, nil)
		return
	}
	cp := make([]common.Address, len(addrs))
	copy(cp, addrs)
	s.calls = append(s.calls, cp)
}

type cnPeerUpdaterOpt func(*cnPeerUpdaterOpts)

type cnPeerUpdaterOpts struct {
	results []cnPeerReadResult
	fork    uint64
}

func withCNPeerReadResults(results ...cnPeerReadResult) cnPeerUpdaterOpt {
	return func(o *cnPeerUpdaterOpts) { o.results = results }
}

func withPermissionlessFork(fork uint64) cnPeerUpdaterOpt {
	return func(o *cnPeerUpdaterOpts) { o.fork = fork }
}

type cnPeerUpdaterTest struct {
	updater *cnPeerUpdater
	reader  *fakeCNPeerReader
	sink    *recordCNPeerSink
}

func newCNPeerUpdaterTest(t *testing.T, options ...cnPeerUpdaterOpt) *cnPeerUpdaterTest {
	t.Helper()

	opts := &cnPeerUpdaterOpts{fork: 1}
	for _, opt := range options {
		opt(opts)
	}

	reader := &fakeCNPeerReader{results: opts.results}
	sink := &recordCNPeerSink{}
	return &cnPeerUpdaterTest{
		updater: newCNPeerUpdater(reader, sink, &params.ChainConfig{
			PermissionlessCompatibleBlock: new(big.Int).SetUint64(opts.fork),
		}),
		reader: reader,
		sink:   sink,
	}
}

func TestCNPeerUpdaterSyncHeadUsesNextBlockNumber(t *testing.T) {
	cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(cnPeerReadResult{addrs: cnTestAddrs(1)}))

	cn.updater.syncHead(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(10)}))

	require.Equal(t, []uint64{11}, cn.reader.calls)
	require.Len(t, cn.sink.calls, 1)
	assert.Equal(t, cnTestAddrs(1), cn.sink.calls[0])
}

func TestCNPeerUpdaterSyncHeadWaitsUntilPermissionlessPlusOne(t *testing.T) {
	cn := newCNPeerUpdaterTest(t,
		withPermissionlessFork(30),
		withCNPeerReadResults(cnPeerReadResult{addrs: cnTestAddrs(1)}),
	)

	cn.updater.syncHead(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(28)})) // target 29
	cn.updater.syncHead(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(29)})) // target 30
	cn.updater.syncHead(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(30)})) // target 31

	require.Equal(t, []uint64{31}, cn.reader.calls)
	require.Len(t, cn.sink.calls, 2)
	assert.Nil(t, cn.sink.calls[0])
	assert.Equal(t, cnTestAddrs(1), cn.sink.calls[1])
}

func TestCNPeerUpdaterSkipsDuplicatePushes(t *testing.T) {
	cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(cnPeerReadResult{addrs: cnTestAddrs(1, 2)}))

	cn.updater.sync(11)
	cn.updater.sync(12)

	require.Len(t, cn.sink.calls, 1)
	assert.Equal(t, cnTestAddrs(1, 2), cn.sink.calls[0])
}

func TestCNPeerUpdaterNilDisablesFiltering(t *testing.T) {
	cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(cnPeerReadResult{addrs: nil}))

	cn.updater.sync(11)
	cn.updater.sync(12)

	require.Len(t, cn.sink.calls, 1)
	assert.Nil(t, cn.sink.calls[0])
}

func TestCNPeerUpdaterNilAndEmptyInputsAreDistinct(t *testing.T) {
	cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(
		cnPeerReadResult{addrs: nil},
		cnPeerReadResult{addrs: []common.Address{}},
	))

	cn.updater.sync(11)
	cn.updater.sync(12)

	require.Len(t, cn.sink.calls, 2)
	assert.Nil(t, cn.sink.calls[0])
	assert.NotNil(t, cn.sink.calls[1])
	assert.Empty(t, cn.sink.calls[1])
}

func TestCNPeerUpdaterFailureDisablesBeforeFirstPush(t *testing.T) {
	cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(cnPeerReadResult{err: errors.New("missing cn peers")}))

	cn.updater.sync(11)

	require.Len(t, cn.sink.calls, 1)
	assert.Nil(t, cn.sink.calls[0])
}

func TestCNPeerUpdaterFailureKeepsLastAllowlist(t *testing.T) {
	for _, tc := range []struct {
		name       string
		failedRead cnPeerReadResult
	}{
		{"read failure", cnPeerReadResult{err: errors.New("temporary read failure")}},
		// GetCNPeers reports a failed read as a nil list, not as an error.
		{"no addresses", cnPeerReadResult{addrs: nil}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cn := newCNPeerUpdaterTest(t, withCNPeerReadResults(
				cnPeerReadResult{addrs: cnTestAddrs(1)},
				tc.failedRead,
			))

			cn.updater.sync(11)
			cn.updater.sync(12)

			require.Len(t, cn.sink.calls, 1)
			assert.Equal(t, cnTestAddrs(1), cn.sink.calls[0])
		})
	}
}

// ChainHeadEvents are posted outside the chain lock, so one can arrive after a newer
// one. The loop must read the canonical head rather than the block in the event.
func TestCNPeerSyncLoopReadsCanonicalHead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	head := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(100)})
	chain := mocks.NewMockBlockChain(ctrl)
	chain.EXPECT().CurrentBlock().Return(head).AnyTimes()

	cn := newCNPeerUpdaterTest(t)
	s := &CN{blockchain: chain, cnPeerUpdater: cn.updater}

	var feed event.Feed
	ch := make(chan blockchain.ChainHeadEvent) // unbuffered, so Send returns once the loop has taken it
	sub := feed.Subscribe(ch)
	s.cnPeerSyncWg.Add(1)
	go s.cnPeerSyncLoop(ch, sub)

	// A stale event for a block far below the canonical head.
	feed.Send(blockchain.ChainHeadEvent{Block: types.NewBlockWithHeader(&types.Header{Number: big.NewInt(50)})})
	sub.Unsubscribe()
	s.cnPeerSyncWg.Wait() // the loop finished the event before exiting

	assert.Equal(t, []uint64{101}, cn.reader.calls, "must look up the head's next block, not the event's")
}

func cnTestAddrs(n ...int) []common.Address {
	addrs := make([]common.Address, len(n))
	for i, num := range n {
		addrs[i] = common.BigToAddress(big.NewInt(int64(num)))
	}
	return addrs
}
