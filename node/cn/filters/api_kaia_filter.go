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
// This file is derived from eth/filters/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package filters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/kaiachain/kaia"
	kaiaApi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/rpc"
)

var (
	defaultFilterDeadline = 5 * time.Minute // consider a filter inactive if it has not been polled for within deadline

	getLogsCxtKeyMaxItems = "maxItems"       // the value of the context key should have the type of GetLogsMaxItems
	GetLogsDeadline       = 10 * time.Second // execution deadlines for getLogs and getFilterLogs APIs
	GetLogsMaxItems       = int(10000)       // maximum allowed number of return items for getLogs and getFilterLogs APIs
)

var errPendingLogsUnsupported = errors.New("pending logs are not supported")

// filter is a helper struct that holds meta information over the filter type
// and associated subscription in the event system.
type filter struct {
	typ      Type
	deadline *time.Timer // filter is inactiv when deadline triggers
	hashes   []common.Hash
	fullTx   bool
	txs      []*types.Transaction
	crit     FilterCriteria
	logs     []*types.Log
	s        *Subscription // associated subscription in event system
}

type FilterAPI struct {
	backend   Backend
	events    *EventSystem
	filtersMu sync.Mutex
	filters   map[rpc.ID]*filter

	supportFullTx bool
	ethHead       bool

	// this field is for test. it makes the filter timeout more flexible when testing
	timeout time.Duration
}

func NewFilterAPI(backend Backend, supportFullTx, ethHead bool) *FilterAPI {
	api := &FilterAPI{
		backend:       backend,
		events:        NewEventSystem(backend.EventMux(), backend),
		filters:       make(map[rpc.ID]*filter),
		timeout:       defaultFilterDeadline,
		supportFullTx: supportFullTx,
		ethHead:       ethHead,
	}
	go api.timeoutLoop()
	return api
}

// timeoutLoop runs every 5 minutes and deletes filters that have not been recently used.
// It is started when the api is created.
func (api *FilterAPI) timeoutLoop() {
	var toUninstall []*Subscription
	ticker := time.NewTicker(api.timeout)
	defer ticker.Stop()
	for {
		<-ticker.C
		api.filtersMu.Lock()
		for id, f := range api.filters {
			select {
			case <-f.deadline.C:
				toUninstall = append(toUninstall, f.s)
				delete(api.filters, id)
			default:
				continue
			}
		}
		api.filtersMu.Unlock()

		// Unsubscribes are processed outside the lock to avoid the following scenario:
		// event loop attempts broadcasting events to still active filters while
		// Unsubscribe is waiting for it to process the uninstall request.
		for _, s := range toUninstall {
			s.Unsubscribe()
		}
		toUninstall = nil
	}
}

// Events return private field events of PublicFilterAPI.
func (api *FilterAPI) Events() *EventSystem {
	return api.events
}

func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		headers := make(chan *types.Header)
		headersSub := api.events.SubscribeNewHeads(headers)

		for {
			select {
			case h := <-headers:
				var header map[string]interface{}
				var err error
				if api.ethHead {
					header, err = kaiaApi.RpcMarshalEthHeader(h, api.backend.Engine(), api.backend.ChainConfig(), true)
					if err != nil {
						logger.Error("Failed to marshal header during newHeads subscription", "err", err)
						headersSub.Unsubscribe()
						return
					}
				} else {
					header = kaiaApi.RPCMarshalHeader(h, api.backend.ChainConfig().Rules(h.Number))
				}
				notifier.Notify(rpcSub.ID, header)
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}

// NewPendingTransactions creates a subscription that is triggered each time a transaction
// enters the transaction pool and was signed from one of the transactions this nodes manages.
func (api *FilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	wantFullTx := fullTx != nil && *fullTx
	if !api.supportFullTx && wantFullTx {
		return nil, rpc.ErrSubscriptionNotFound
	}

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		txs := make(chan []*types.Transaction, 128)
		pendingTxSub := api.events.SubscribePendingTxs(txs)
		defer pendingTxSub.Unsubscribe()

		for {
			select {
			case txs := <-txs:
				// To keep the original behaviour, send a single tx hash in one notification.
				// TODO(rjl493456442) Send a batch of tx hashes in one notification
				for _, tx := range txs {
					if wantFullTx {
						m := tx.MakeRPCOutput()
						m["time"] = tx.Time()
						notifier.Notify(rpcSub.ID, m)
					} else {
						notifier.Notify(rpcSub.ID, tx.Hash())
					}
				}
			case <-rpcSub.Err():
				return
			}
		}
	}()

	return rpcSub, nil
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (api *FilterAPI) Logs(ctx context.Context, crit FilterCriteria) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	var (
		rpcSub      = notifier.CreateSubscription()
		matchedLogs = make(chan []*types.Log)
	)

	logsSub, err := api.events.SubscribeLogs(kaia.FilterQuery(crit), matchedLogs)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case logs := <-matchedLogs:
				for _, log := range logs {
					notifier.Notify(rpcSub.ID, &log)
				}
			case <-rpcSub.Err(): // client send an unsubscribe request
				logsSub.Unsubscribe()
				return
			case <-notifier.Closed(): // connection dropped
				logsSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}

// KaiaFilterAPI offers support to create and manage filters. This will allow external clients to retrieve various
// information related to the Kaia protocol such as blocks, transactions and logs.
type KaiaFilterAPI struct {
	*FilterAPI
}

// NewKaiaFilterAPI returns a new FilterAPI instance for kaia namespace.
func NewKaiaFilterAPI(backend Backend) *KaiaFilterAPI {
	return &KaiaFilterAPI{NewFilterAPI(backend, false, false)}
}

// NewEthFilterAPI returns a new FilterAPI instance for eth namespace.
func NewEthFilterAPI(backend Backend) *KaiaFilterAPI {
	return &KaiaFilterAPI{NewFilterAPI(backend, false, true)}
}

// NewAuctionFilterAPI returns a new FilterAPI instance for auction namespace.
func NewAuctionFilterAPI(backend Backend) *KaiaFilterAPI {
	return &KaiaFilterAPI{NewFilterAPI(backend, true, false)}
}

// NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
// as transactions enter the pending state.
//
// It is part of the filter package because this filter can be used through the
// `kaia_getFilterChanges` polling method that is also used for log filters.
func (api *KaiaFilterAPI) NewPendingTransactionFilter() rpc.ID {
	var (
		pendingTxs   = make(chan []*types.Transaction)
		pendingTxSub = api.events.SubscribePendingTxs(pendingTxs)
	)

	api.filtersMu.Lock()
	api.filters[pendingTxSub.ID] = &filter{typ: PendingTransactionsSubscription, deadline: time.NewTimer(api.timeout), txs: make([]*types.Transaction, 0), s: pendingTxSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case pTx := <-pendingTxs:
				api.filtersMu.Lock()
				if f, found := api.filters[pendingTxSub.ID]; found {
					f.txs = append(f.txs, pTx...)
				}
				api.filtersMu.Unlock()
			case <-pendingTxSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, pendingTxSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return pendingTxSub.ID
}

func (api *KaiaFilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	return api.FilterAPI.NewPendingTransactions(ctx, fullTx)
}

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
func (api *KaiaFilterAPI) NewBlockFilter() rpc.ID {
	var (
		headers   = make(chan *types.Header)
		headerSub = api.events.SubscribeNewHeads(headers)
	)

	api.filtersMu.Lock()
	api.filters[headerSub.ID] = &filter{typ: BlocksSubscription, deadline: time.NewTimer(api.timeout), hashes: make([]common.Hash, 0), s: headerSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case h := <-headers:
				api.filtersMu.Lock()
				if f, found := api.filters[headerSub.ID]; found {
					f.hashes = append(f.hashes, h.Hash())
				}
				api.filtersMu.Unlock()
			case <-headerSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, headerSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return headerSub.ID
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *KaiaFilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	return api.FilterAPI.NewHeads(ctx)
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (api *KaiaFilterAPI) Logs(ctx context.Context, crit FilterCriteria) (*rpc.Subscription, error) {
	return api.FilterAPI.Logs(ctx, crit)
}

// FilterCriteria represents a request to create a new filter.
// Same as Kaia.FilterQuery but with UnmarshalJSON() method.
type FilterCriteria kaia.FilterQuery

// NewFilter creates a new filter and returns the filter id. It can be
// used to retrieve logs when the state changes. This method cannot be
// used to fetch logs that are already stored in the state.
//
// Default criteria for the from and to block are "latest".
// Using "latest" as block number will return logs for mined blocks.
// Using "pending" as block number returns logs for not yet mined (pending) blocks.
// In case logs are removed (chain reorg) previously returned logs are returned
// again but with the removed property set to true.
//
// In case "fromBlock" > "toBlock" an error is returned.
func (api *KaiaFilterAPI) NewFilter(crit FilterCriteria) (rpc.ID, error) {
	logs := make(chan []*types.Log)
	logsSub, err := api.events.SubscribeLogs(kaia.FilterQuery(crit), logs)
	if err != nil {
		return rpc.ID(""), err
	}

	api.filtersMu.Lock()
	api.filters[logsSub.ID] = &filter{typ: LogsSubscription, crit: crit, deadline: time.NewTimer(api.timeout), logs: make([]*types.Log, 0), s: logsSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case l := <-logs:
				api.filtersMu.Lock()
				if f, found := api.filters[logsSub.ID]; found {
					f.logs = append(f.logs, l...)
				}
				api.filtersMu.Unlock()
			case <-logsSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, logsSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()

	return logsSub.ID, nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
func (api *KaiaFilterAPI) GetLogs(ctx context.Context, crit FilterCriteria) ([]*types.Log, error) {
	ctx = context.WithValue(ctx, getLogsCxtKeyMaxItems, GetLogsMaxItems)
	ctx, cancelFnc := context.WithTimeout(ctx, GetLogsDeadline)
	defer cancelFnc()

	var filter *Filter
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.backend, *crit.BlockHash, crit.Addresses, crit.Topics)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewRangeFilter(api.backend, begin, end, crit.Addresses, crit.Topics)
	}

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), err
}

// UninstallFilter removes the filter with the given filter id.
func (api *KaiaFilterAPI) UninstallFilter(id rpc.ID) bool {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	if found {
		delete(api.filters, id)
	}
	api.filtersMu.Unlock()
	if found {
		f.s.Unsubscribe()
	}

	return found
}

// GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
func (api *KaiaFilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*types.Log, error) {
	ctx = context.WithValue(ctx, getLogsCxtKeyMaxItems, GetLogsMaxItems)
	ctx, cancelFnc := context.WithTimeout(ctx, GetLogsDeadline)
	defer cancelFnc()

	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()

	if !found || f.typ != LogsSubscription {
		return nil, fmt.Errorf("filter not found")
	}

	begin := rpc.LatestBlockNumber.Int64()
	if f.crit.FromBlock != nil {
		begin = f.crit.FromBlock.Int64()
	}
	end := rpc.LatestBlockNumber.Int64()
	if f.crit.ToBlock != nil {
		end = f.crit.ToBlock.Int64()
	}
	// Create and run the filter to get all the logs
	filter := NewRangeFilter(api.backend, begin, end, f.crit.Addresses, f.crit.Topics)

	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil
}

// GetFilterChanges returns the logs for the filter with the given id since
// last time it was called. This can be used for polling.
//
// For pending transaction and block filters the result is []common.Hash.
// (pending)Log filters return []Log.
func (api *KaiaFilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if f, found := api.filters[id]; found {
		if !f.deadline.Stop() {
			// timer expired but filter is not yet removed in timeout loop
			// receive timer value and reset timer
			<-f.deadline.C
		}
		f.deadline.Reset(api.timeout)

		switch f.typ {
		case BlocksSubscription:
			hashes := f.hashes
			f.hashes = nil
			return returnHashes(hashes), nil
		case PendingTransactionsSubscription:
			if f.fullTx {
				txs := make([]map[string]interface{}, 0, len(f.txs))
				for _, tx := range f.txs {
					txs = append(txs, tx.MakeRPCOutput())
				}
				f.txs = nil
				return txs, nil
			} else {
				hashes := make([]common.Hash, 0, len(f.txs))
				for _, tx := range f.txs {
					hashes = append(hashes, tx.Hash())
				}
				f.txs = nil
				return hashes, nil
			}
		case LogsSubscription:
			logs := f.logs
			f.logs = nil
			return returnLogs(logs), nil
		}
	}

	return []interface{}{}, fmt.Errorf("filter not found")
}

// returnHashes is a helper that will return an empty hash array case the given hash array is nil,
// otherwise the given hashes array is returned.
func returnHashes(hashes []common.Hash) []common.Hash {
	if hashes == nil {
		return []common.Hash{}
	}
	return hashes
}

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*types.Log) []*types.Log {
	if logs == nil {
		return []*types.Log{}
	}
	return logs
}

// UnmarshalJSON sets *args fields with given data.
func (args *FilterCriteria) UnmarshalJSON(data []byte) error {
	type input struct {
		BlockHash *common.Hash     `json:"blockHash"`
		FromBlock *rpc.BlockNumber `json:"fromBlock"`
		ToBlock   *rpc.BlockNumber `json:"toBlock"`
		Addresses interface{}      `json:"address"`
		Topics    []interface{}    `json:"topics"`
	}

	var raw input
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.BlockHash != nil {
		if raw.FromBlock != nil || raw.ToBlock != nil {
			// BlockHash is mutually exclusive with FromBlock/ToBlock criteria
			return fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")
		}
		args.BlockHash = raw.BlockHash
	} else {
		if raw.FromBlock != nil {
			args.FromBlock = big.NewInt(raw.FromBlock.Int64())
		}
		if raw.ToBlock != nil {
			args.ToBlock = big.NewInt(raw.ToBlock.Int64())
		}
	}

	args.Addresses = []common.Address{}

	if raw.Addresses != nil {
		// raw.Address can contain a single address or an array of addresses
		switch rawAddr := raw.Addresses.(type) {
		case []interface{}:
			for i, addr := range rawAddr {
				if strAddr, ok := addr.(string); ok {
					addr, err := decodeAddress(strAddr)
					if err != nil {
						return fmt.Errorf("invalid address at index %d: %v", i, err)
					}
					args.Addresses = append(args.Addresses, addr)
				} else {
					return fmt.Errorf("non-string address at index %d", i)
				}
			}
		case string:
			addr, err := decodeAddress(rawAddr)
			if err != nil {
				return fmt.Errorf("invalid address: %v", err)
			}
			args.Addresses = []common.Address{addr}
		default:
			return errors.New("invalid addresses in query")
		}
	}

	// topics is an array consisting of strings and/or arrays of strings.
	// JSON null values are converted to common.Hash{} and ignored by the filter manager.
	if len(raw.Topics) > 0 {
		args.Topics = make([][]common.Hash, len(raw.Topics))
		for i, t := range raw.Topics {
			switch topic := t.(type) {
			case nil:
				// ignore topic when matching logs

			case string:
				// match specific topic
				top, err := decodeTopic(topic)
				if err != nil {
					return err
				}
				args.Topics[i] = []common.Hash{top}

			case []interface{}:
				// or case e.g. [null, "topic0", "topic1"]
				for _, rawTopic := range topic {
					if rawTopic == nil {
						// null component, match all
						args.Topics[i] = nil
						break
					}
					if topic, ok := rawTopic.(string); ok {
						parsed, err := decodeTopic(topic)
						if err != nil {
							return err
						}
						args.Topics[i] = append(args.Topics[i], parsed)
					} else {
						return fmt.Errorf("invalid topic(s)")
					}
				}
			default:
				return fmt.Errorf("invalid topic(s)")
			}
		}
	}

	return nil
}

func decodeAddress(s string) (common.Address, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != common.AddressLength {
		err = fmt.Errorf("hex has invalid length %d after decoding", len(b))
	}
	return common.BytesToAddress(b), err
}

func decodeTopic(s string) (common.Hash, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != common.HashLength {
		err = fmt.Errorf("hex has invalid length %d after decoding", len(b))
	}
	return common.BytesToHash(b), err
}
