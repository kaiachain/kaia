// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from eth/filters/filter_system.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package filters

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
)

// Type determines the kind of filter and is used to put the filter in to
// the correct bucket when added.
type Type byte

const (
	// UnknownSubscription indicates an unknown subscription type
	UnknownSubscription Type = iota
	// LogsSubscription queries for new or removed (chain reorg) logs
	LogsSubscription
	// PendingTransactionsSubscription queries tx hashes for pending
	// transactions entering the pending state
	PendingTransactionsSubscription
	// BlocksSubscription queries hashes for blocks that are imported
	BlocksSubscription
	// LastSubscription keeps track of the last index
	LastIndexSubscription
)

const (

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// rmLogsChanSize is the size of channel listening to RemovedLogsEvent.
	rmLogsChanSize = 10
	// logsChanSize is the size of channel listening to LogsEvent.
	logsChanSize = 10
	// chainEvChanSize is the size of channel listening to ChainEvent.
	chainEvChanSize = 10
)

var (
	ErrInvalidSubscriptionID = errors.New("invalid id")
	logger                   = log.NewModuleLogger(log.NodeCNFilters)
)

type subscription struct {
	id        rpc.ID
	typ       Type
	created   time.Time
	logsCrit  kaia.FilterQuery
	logs      chan []*types.Log
	txs       chan []*types.Transaction
	headers   chan *types.Header
	installed chan struct{} // closed when the filter is installed
	err       chan error    // closed when the filter is uninstalled
}

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria.
type EventSystem struct {
	mux     *event.TypeMux
	backend Backend

	// Subscriptions
	txsSub    event.Subscription // Subscription for new transaction event
	logsSub   event.Subscription // Subscription for new log event
	rmLogsSub event.Subscription // Subscription for removed log event
	chainSub  event.Subscription // Subscription for new chain event

	// Channels
	install   chan *subscription               // install filter for event notification
	uninstall chan *subscription               // remove filter for event notification
	txsCh     chan blockchain.NewTxsEvent      // Channel to receive new transactions event
	logsCh    chan []*types.Log                // Channel to receive new log event
	rmLogsCh  chan blockchain.RemovedLogsEvent // Channel to receive removed log event
	chainCh   chan blockchain.ChainEvent       // Channel to receive new chain event
}

// NewEventSystem creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewEventSystem(mux *event.TypeMux, backend Backend) *EventSystem {
	m := &EventSystem{
		mux:       mux,
		backend:   backend,
		install:   make(chan *subscription),
		uninstall: make(chan *subscription),
		txsCh:     make(chan blockchain.NewTxsEvent, txChanSize),
		logsCh:    make(chan []*types.Log, logsChanSize),
		rmLogsCh:  make(chan blockchain.RemovedLogsEvent, rmLogsChanSize),
		chainCh:   make(chan blockchain.ChainEvent, chainEvChanSize),
	}

	// Subscribe events
	m.txsSub = m.backend.SubscribeNewTxsEvent(m.txsCh)
	m.logsSub = m.backend.SubscribeLogsEvent(m.logsCh)
	m.rmLogsSub = m.backend.SubscribeRemovedLogsEvent(m.rmLogsCh)
	m.chainSub = m.backend.SubscribeChainEvent(m.chainCh)

	// Make sure none of the subscriptions are empty
	if m.txsSub == nil || m.logsSub == nil || m.rmLogsSub == nil || m.chainSub == nil {
		logger.Crit("Subscribe for event system failed")
	}

	go m.eventLoop()
	return m
}

// Subscription is created when the client registers itself for a particular event.
type Subscription struct {
	ID        rpc.ID
	f         *subscription
	es        *EventSystem
	unsubOnce sync.Once
}

// Err returns a channel that is closed when unsubscribed.
func (sub *Subscription) Err() <-chan error {
	return sub.f.err
}

// Unsubscribe uninstalls the subscription from the event broadcast loop.
func (sub *Subscription) Unsubscribe() {
	sub.unsubOnce.Do(func() {
	uninstallLoop:
		for {
			// write uninstall request and consume logs/hashes. This prevents
			// the eventLoop broadcast method to deadlock when writing to the
			// filter event channel while the subscription loop is waiting for
			// this method to return (and thus not reading these events).
			select {
			case sub.es.uninstall <- sub.f:
				break uninstallLoop
			case <-sub.f.logs:
			case <-sub.f.txs:
			case <-sub.f.headers:
			}
		}

		// wait for filter to be uninstalled in work loop before returning
		// this ensures that the manager won't use the event channel which
		// will probably be closed by the client asap after this method returns.
		<-sub.Err()
	})
}

// subscribe installs the subscription in the event broadcast loop.
func (es *EventSystem) subscribe(sub *subscription) *Subscription {
	es.install <- sub
	<-sub.installed
	return &Subscription{ID: sub.id, f: sub, es: es}
}

// SubscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel. Default value for the from and to
// block is "latest". If the fromBlock > toBlock an error is returned.
func (es *EventSystem) SubscribeLogs(crit kaia.FilterQuery, logs chan []*types.Log) (*Subscription, error) {
	var from, to rpc.BlockNumber
	if crit.FromBlock == nil {
		from = rpc.LatestBlockNumber
	} else {
		from = rpc.BlockNumber(crit.FromBlock.Int64())
	}
	if crit.ToBlock == nil {
		to = rpc.LatestBlockNumber
	} else {
		to = rpc.BlockNumber(crit.ToBlock.Int64())
	}

	// Pending logs are not supported anymore.
	if from == rpc.PendingBlockNumber || to == rpc.PendingBlockNumber {
		return nil, errPendingLogsUnsupported
	}

	// only interested in new mined logs
	if from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber {
		return es.subscribeLogs(crit, logs), nil
	}
	// only interested in mined logs within a specific block range
	if from >= 0 && to >= 0 && to >= from {
		return es.subscribeLogs(crit, logs), nil
	}
	// interested in logs from a specific block number to new mined blocks
	if from >= 0 && to == rpc.LatestBlockNumber {
		return es.subscribeLogs(crit, logs), nil
	}
	return nil, fmt.Errorf("invalid from and to block combination: from > to")
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (es *EventSystem) subscribeLogs(crit kaia.FilterQuery, logs chan []*types.Log) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       LogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      logs,
		txs:       make(chan []*types.Transaction),
		headers:   make(chan *types.Header),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

// SubscribeNewHeads creates a subscription that writes the header of a block that is
// imported in the chain.
func (es *EventSystem) SubscribeNewHeads(headers chan *types.Header) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       BlocksSubscription,
		created:   time.Now(),
		logs:      make(chan []*types.Log),
		txs:       make(chan []*types.Transaction),
		headers:   headers,
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

// SubscribePendingTxs creates a subscription that writes transaction for
// transactions that enter the transaction pool.
func (es *EventSystem) SubscribePendingTxs(txs chan []*types.Transaction) *Subscription {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       PendingTransactionsSubscription,
		created:   time.Now(),
		logs:      make(chan []*types.Log),
		txs:       txs,
		headers:   make(chan *types.Header),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return es.subscribe(sub)
}

type filterIndex map[Type]map[rpc.ID]*subscription

// broadcast event to filters that match criteria.
func (es *EventSystem) broadcast(filters filterIndex, ev interface{}) {
	if ev == nil {
		return
	}

	switch e := ev.(type) {
	case []*types.Log:
		if len(e) > 0 {
			for _, f := range filters[LogsSubscription] {
				if matchedLogs := filterLogs(e, f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics); len(matchedLogs) > 0 {
					f.logs <- matchedLogs
				}
			}
		}
	case blockchain.RemovedLogsEvent:
		for _, f := range filters[LogsSubscription] {
			if matchedLogs := filterLogs(e.Logs, f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics); len(matchedLogs) > 0 {
				f.logs <- matchedLogs
			}
		}
	case blockchain.NewTxsEvent:
		for _, f := range filters[PendingTransactionsSubscription] {
			f.txs <- e.Txs
		}
	case blockchain.ChainEvent:
		for _, f := range filters[BlocksSubscription] {
			f.headers <- e.Block.Header()
		}
	}
}

// eventLoop (un)installs filters and processes mux events.
func (es *EventSystem) eventLoop() {
	// Ensure all subscriptions get cleaned up
	defer func() {
		es.txsSub.Unsubscribe()
		es.logsSub.Unsubscribe()
		es.rmLogsSub.Unsubscribe()
		es.chainSub.Unsubscribe()
	}()

	index := make(filterIndex)
	for i := UnknownSubscription; i < LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*subscription)
	}

	for {
		select {
		// Handle subscribed events
		case ev := <-es.txsCh:
			es.broadcast(index, ev)
		case ev := <-es.logsCh:
			es.broadcast(index, ev)
		case ev := <-es.rmLogsCh:
			es.broadcast(index, ev)
		case ev := <-es.chainCh:
			es.broadcast(index, ev)

		case f := <-es.install:
			index[f.typ][f.id] = f
			close(f.installed)

		case f := <-es.uninstall:
			delete(index[f.typ], f.id)
			close(f.err)

			// System stopped
		case <-es.txsSub.Err():
			return
		case <-es.logsSub.Err():
			return
		case <-es.rmLogsSub.Err():
			return
		case <-es.chainSub.Err():
			return
		}
	}
}
