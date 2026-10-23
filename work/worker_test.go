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
package work

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	consensusmocks "github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
)

// newSequenceEvent stands in for istanbul.NewSequenceEvent.
type newSequenceEvent struct{}

type testBackend struct {
	chain  BlockChain
	txPool TxPool
}

func (b *testBackend) AccountManager() accounts.AccountManager { return nil }
func (b *testBackend) BlockChain() BlockChain                  { return b.chain }
func (b *testBackend) TxPool() TxPool                          { return b.txPool }
func (b *testBackend) ChainDB() database.DBManager             { return nil }
func (b *testBackend) ReBroadcastTxs(types.Transactions)       {}

// newTestWorker runs a worker's update loop on finalizeCh and mux. Each new sequence it takes is sent on taken.
func newTestWorker(t *testing.T, ctrl *gomock.Controller, engine consensus.Engine, finalizeCh <-chan *consensus.ExecutionResult, mux *event.TypeMux) (w *worker, taken <-chan struct{}) {
	chain := mocks.NewMockBlockChain(ctrl)
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0), Time: big.NewInt(0)})).AnyTimes()
	// commitNewWork stops at Pending.
	takenCh := make(chan struct{}, 1)
	txPool := mocks.NewMockTxPool(ctrl)
	txPool.EXPECT().Pending().DoAndReturn(func() (map[common.Address]types.Transactions, error) {
		takenCh <- struct{}{}
		return nil, errors.New("no pending txs")
	}).AnyTimes()

	var txsFeed, headFeed, sideFeed event.Feed
	w = &worker{
		config:         params.TestChainConfig,
		engine:         engine,
		backend:        &testBackend{chain: chain, txPool: txPool},
		chain:          chain,
		txsCh:          make(chan blockchain.NewTxsEvent, txChanSize),
		chainHeadCh:    make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		chainSideCh:    make(chan blockchain.ChainSideEvent, chainSideChanSize),
		finalizeCh:     finalizeCh,
		newSequenceSub: mux.Subscribe(newSequenceEvent{}),
	}
	w.txsSub = txsFeed.Subscribe(w.txsCh)
	w.chainHeadSub = headFeed.Subscribe(w.chainHeadCh)
	w.chainSideSub = sideFeed.Subscribe(w.chainSideCh)
	go w.update()
	t.Cleanup(w.chainHeadSub.Unsubscribe)
	return w, takenCh
}

// Tests that starting or stopping the engine doesn't deadlock when the update loop gets a result meanwhile.
func TestWorkerStartStopWithPendingResult(t *testing.T) {
	for _, tc := range []struct {
		name   string
		expect func(engine *consensusmocks.MockEngine, engineCall func() error)
		run    func(w *worker)
	}{
		{
			name: "start",
			expect: func(engine *consensusmocks.MockEngine, engineCall func() error) {
				engine.EXPECT().Start(gomock.Any(), gomock.Any()).DoAndReturn(func(consensus.ChainReader, consensus.Executor) error { return engineCall() })
			},
			run: func(w *worker) { w.start() },
		},
		{
			name: "stop",
			expect: func(engine *consensusmocks.MockEngine, engineCall func() error) {
				engine.EXPECT().Stop().DoAndReturn(engineCall)
			},
			run: func(w *worker) { w.stop() },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			engine := consensusmocks.NewMockEngine(ctrl)
			finalizeCh := make(chan *consensus.ExecutionResult)
			mux := new(event.TypeMux)
			w, taken := newTestWorker(t, ctrl, engine, finalizeCh, mux)

			// Like the engine: a Seal left from before the stop returns nil, then a new sequence waits for the update loop.
			tc.expect(engine, func() error {
				finalizeCh <- nil
				return mux.Post(newSequenceEvent{})
			})
			done := make(chan struct{})
			go func() {
				tc.run(w)
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatalf("worker %s deadlocked", tc.name)
			}
			select {
			case <-taken:
			case <-time.After(time.Second):
				t.Fatal("update loop didn't take the new sequence")
			}
		})
	}
}

// Tests that the update loop doesn't deadlock posting its own ChainHeadEvent when chainHeadCh is full.
// The loop is the only reader of chainHeadCh, so it must not block sending to it.
func TestWorkerPostOwnChainHeadWithFullChannel(t *testing.T) {
	for _, tc := range []struct {
		name   string
		queued int // ChainHeadEvents imported while the loop writes the mined block
	}{
		{name: "not full", queued: chainHeadChanSize - 2},
		{name: "full", queued: chainHeadChanSize},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			engine := consensusmocks.NewMockEngine(ctrl)
			chain := mocks.NewMockBlockChain(ctrl)
			var txsFeed, headFeed, sideFeed event.Feed

			// Like an importer holding bc.mu, WriteBlockWithState waits for release.
			writing, release := make(chan struct{}), make(chan struct{})
			chain.EXPECT().WriteBlockWithState(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(*types.Block, []*types.Receipt, *state.StateDB) (blockchain.WriteResult, error) {
					close(writing)
					<-release
					return blockchain.WriteResult{Status: blockchain.CanonStatTy}, nil
				})
			// Like BlockChain.PostChainEvents, send the ChainHeadEvent on the head feed.
			posted := make(chan struct{})
			chain.EXPECT().PostChainEvents(gomock.Any(), gomock.Any()).DoAndReturn(func(events []interface{}, _ []*types.Log) {
				for _, ev := range events {
					if head, ok := ev.(blockchain.ChainHeadEvent); ok {
						headFeed.Send(head)
					}
				}
				close(posted)
			})

			finalizeCh := make(chan *consensus.ExecutionResult)
			mux := new(event.TypeMux)
			statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
			header := &types.Header{Number: big.NewInt(1), Time: big.NewInt(1)}
			w := &worker{
				config:         params.TestChainConfig,
				engine:         engine,
				chain:          chain,
				mux:            mux,
				txsCh:          make(chan blockchain.NewTxsEvent, txChanSize),
				chainHeadCh:    make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
				chainSideCh:    make(chan blockchain.ChainSideEvent, chainSideChanSize),
				finalizeCh:     finalizeCh,
				newSequenceSub: mux.Subscribe(newSequenceEvent{}),
				pendingWork:    &Task{state: statedb, header: header},
			}
			w.txsSub = txsFeed.Subscribe(w.txsCh)
			w.chainHeadSub = headFeed.Subscribe(w.chainHeadCh)
			w.chainSideSub = sideFeed.Subscribe(w.chainSideCh)
			go w.update()
			t.Cleanup(w.chainHeadSub.Unsubscribe)

			block := types.NewBlockWithHeader(header)
			finalizeCh <- &consensus.ExecutionResult{State: statedb, Block: block}
			<-writing
			for i := 0; i < tc.queued; i++ {
				headFeed.Send(blockchain.ChainHeadEvent{Block: block})
			}
			close(release)

			// The next import must be able to post its ChainHeadEvent.
			done := make(chan struct{})
			go func() {
				headFeed.Send(blockchain.ChainHeadEvent{Block: block})
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fatalf("update loop deadlocked posting its ChainHeadEvent, chainHeadCh len=%d", len(w.chainHeadCh))
			}
			select {
			case <-posted:
			case <-time.After(time.Second):
				t.Fatal("mined block's chain events weren't posted")
			}
		})
	}
}
