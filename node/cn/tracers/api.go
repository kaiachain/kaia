// Modifications Copyright 2022 The klaytn Authors
// Copyright 2021 The go-ethereum Authors
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
// This file is derived from eth/tracers/api.go (2022/08/08).
// Modified and improved for the klaytn development.

package tracers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	kaiaapi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultLoggerTimeout is the amount of time a logger can aggregate trace logs
	defaultLoggerTimeout = 1 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)

	// defaultTracechainMemLimit is the size of the triedb, at which traceChain
	// switches over and tries to use a disk-backed database instead of building
	// on top of memory.
	// For non-archive nodes, this limit _will_ be overblown, as disk-backed tries
	// will only be found every ~15K blocks or so.
	// For Kaia, this value is set to a value 4 times larger compared to the ethereum setting.
	defaultTracechainMemLimit = common.StorageSize(4 * 500 * 1024 * 1024)

	// fastCallTracer is the go-version callTracer which is lighter and faster than
	// Javascript version.
	fastCallTracer = "fastCallTracer"
)

var (
	HeavyAPIRequestLimit int32 = 500 // WARN: changing this value will have no effect. This value is for test. See HeavyDebugRequestLimitFlag
	heavyAPIRequestCount int32 = 0
)

// StateReleaseFunc is used to deallocate resources held by constructing a
// historical state for tracing purposes.
type StateReleaseFunc func()

// Backend interface provides the common API services with access to necessary functions.
type Backend interface {
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	RPCGasCap() *big.Int
	ChainConfig() *params.ChainConfig
	ChainDB() database.DBManager
	Engine() consensus.Engine
	// StateAtBlock returns the state corresponding to the stateroot of the block.
	// N.B: For executing transactions on block N, the required stateRoot is block N-1,
	// so this method should be called with the parent.
	StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (*state.StateDB, StateReleaseFunc, error)
	StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (blockchain.Message, vm.BlockContext, vm.TxContext, *state.StateDB, StateReleaseFunc, error)
}

// CommonAPI contains
// - public methods that change behavior depending on `.unsafeTrace` flag.
// For instance, TraceTransaction and TraceCall may or may not support custom tracers.
// - private helper methods such as traceTx
type CommonAPI struct {
	backend     Backend
	unsafeTrace bool
}

// API contains public methods that are considered "safe" to expose in public RPC.
type API struct {
	CommonAPI
}

// UnsafeAPI contains public methods that are considered "unsafe" to expose in public RPC.
type UnsafeAPI struct {
	CommonAPI
}

// NewUnsafeAPI creates a new UnsafeAPI definition
func NewUnsafeAPI(backend Backend) *UnsafeAPI {
	return &UnsafeAPI{
		CommonAPI{backend: backend, unsafeTrace: true},
	}
}

// NewAPI creates a new API definition
func NewAPI(backend Backend) *API {
	return &API{
		CommonAPI{backend: backend, unsafeTrace: false},
	}
}

type chainContext struct {
	backend Backend
	ctx     context.Context
}

func (context *chainContext) Engine() consensus.Engine {
	return context.backend.Engine()
}

func (context *chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	header, err := context.backend.HeaderByNumber(context.ctx, rpc.BlockNumber(number))
	if err != nil {
		return nil
	}
	if header.Hash() == hash {
		return header
	}
	header, err = context.backend.HeaderByHash(context.ctx, hash)
	if err != nil {
		return nil
	}
	return header
}

// chainContext constructs the context reader which is used by the evm for reading
// the necessary chain context.
func newChainContext(ctx context.Context, backend Backend) blockchain.ChainContext {
	return &chainContext{backend: backend, ctx: ctx}
}

// blockByNumber is the wrapper of the chain access function offered by the backend.
// It will return an error if the block is not found.
func (api *CommonAPI) blockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	return api.backend.BlockByNumber(ctx, number)
}

// blockByHash is the wrapper of the chain access function offered by the backend.
// It will return an error if the block is not found.
func (api *CommonAPI) blockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return api.backend.BlockByHash(ctx, hash)
}

// blockByNumberAndHash is the wrapper of the chain access function offered by
// the backend. It will return an error if the block is not found.
func (api *CommonAPI) blockByNumberAndHash(ctx context.Context, number rpc.BlockNumber, hash common.Hash) (*types.Block, error) {
	block, err := api.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if block.Hash() == hash {
		return block, nil
	}
	return api.blockByHash(ctx, hash)
}

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*vm.LogConfig
	Tracer        *string
	Timeout       *string
	LoggerTimeout *string
	Reexec        *uint64
}

// StdTraceConfig holds extra parameters to standard-json trace functions.
type StdTraceConfig struct {
	*vm.LogConfig
	Reexec *uint64
	TxHash common.Hash
}

// txTraceResult is the result of a single transaction trace.
type txTraceResult struct {
	TxHash common.Hash `json:"txHash,omitempty"` // transaction hash
	Result interface{} `json:"result,omitempty"` // Trace results produced by the tracer
	Error  string      `json:"error,omitempty"`  // Trace failure produced by the tracer
}

// blockTraceTask represents a single block trace task when an entire chain is
// being traced.
type blockTraceTask struct {
	statedb *state.StateDB   // Intermediate state prepped for tracing
	block   *types.Block     // Block to trace the transactions from
	release StateReleaseFunc // The function to release the held resource for this task
	results []*txTraceResult // Trace results procudes by the task
}

// blockTraceResult represets the results of tracing a single block when an entire
// chain is being traced.
type blockTraceResult struct {
	Block  hexutil.Uint64   `json:"block"`  // Block number corresponding to this trace
	Hash   common.Hash      `json:"hash"`   // Block hash corresponding to this trace
	Traces []*txTraceResult `json:"traces"` // Trace results produced by the task
}

// txTraceTask represents a single transaction trace task when an entire block
// is being traced.
type txTraceTask struct {
	statedb *state.StateDB // Intermediate state prepped for tracing
	index   int            // Transaction offset in the block
}

func checkRangeAndReturnBlock(api *CommonAPI, ctx context.Context, start, end rpc.BlockNumber) (*types.Block, *types.Block, error) {
	// Fetch the block interval that we want to trace
	from, err := api.blockByNumber(ctx, start)
	if err != nil {
		return nil, nil, err
	}
	to, err := api.blockByNumber(ctx, end)
	if err != nil {
		return nil, nil, err
	}

	// Trace the chain if we've found all our blocks
	if from == nil {
		return nil, nil, fmt.Errorf("starting block #%d not found", start)
	}
	if to == nil {
		return nil, nil, fmt.Errorf("end block #%d not found", end)
	}
	if from.Number().Cmp(to.Number()) >= 0 {
		return nil, nil, fmt.Errorf("end block #%d needs to come after start block #%d", end, start)
	}
	return from, to, nil
}

// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.
func (api *UnsafeAPI) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	from, to, err := checkRangeAndReturnBlock(&api.CommonAPI, ctx, start, end)
	if err != nil {
		return nil, err
	}
	// Tracing a chain is a **long** operation, only do with subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	sub := notifier.CreateSubscription()
	_, err = api.traceChain(from, to, config, notifier, sub)
	return sub, err
}

// releaser is a helper tool responsible for caching the release
// callbacks of tracing state.
type releaser struct {
	releases []StateReleaseFunc
	lock     sync.Mutex
}

func (r *releaser) add(release StateReleaseFunc) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.releases = append(r.releases, release)
}

func (r *releaser) call() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, release := range r.releases {
		release()
	}
	r.releases = r.releases[:0]
}

// traceChain configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The tracing chain range includes
// the end block but excludes the start one. The return value will be one item per
// transaction, dependent on the requested tracer.
// The tracing procedure should be aborted in case the closed signal is received.
//
// The traceChain operates in two modes: subscription mode and rpc mode
//   - if notifier and sub is not nil, it works as a subscription mode and returns nothing
//   - if those parameters are nil, it works as a rpc mode and returns the block trace results, so it can pass the result through rpc-call
func (api *CommonAPI) traceChain(start, end *types.Block, config *TraceConfig, notifier *rpc.Notifier, sub *rpc.Subscription) (map[uint64]*blockTraceResult, error) {
	// Prepare all the states for tracing. Note this procedure can take very
	// long time. Timeout mechanism is necessary.
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	// Execute all the transaction contained within the chain concurrently for each block
	blocks := int(end.NumberU64() - start.NumberU64())
	threads := runtime.NumCPU()
	if threads > blocks {
		threads = blocks
	}
	var (
		pend     = new(sync.WaitGroup)
		tasks    = make(chan *blockTraceTask, threads)
		results  = make(chan *blockTraceTask, threads)
		localctx = context.Background()
		reler    = new(releaser)
	)
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the block trace tasks
			for task := range tasks {
				signer := types.MakeSigner(api.backend.ChainConfig(), task.block.Number())
				blockCtx := blockchain.NewEVMBlockContext(task.block.Header(), newChainContext(localctx, api.backend), nil)

				// Trace all the transactions contained within
				for i, tx := range task.block.Transactions() {
					msg, err := tx.AsMessageWithAccountKeyPicker(signer, task.statedb, task.block.NumberU64())
					if err != nil {
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						break
					}

					txCtx := blockchain.NewEVMTxContext(msg, task.block.Header(), api.backend.ChainConfig())

					res, err := api.traceTx(localctx, msg, blockCtx, txCtx, task.statedb, config)
					if err != nil {
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						break
					}
					task.statedb.Finalise(true, true)
					task.results[i] = &txTraceResult{TxHash: tx.Hash(), Result: res}
				}
				// Tracing state is used up, queue it for de-referencing
				reler.add(task.release)

				// Stream the result back to the result catcher or abort on teardown
				if notifier != nil {
					// Stream the result back to the user or abort on teardown
					select {
					case results <- task:
					case <-notifier.Closed():
						return
					}
				} else {
					results <- task
				}
			}
		}()
	}
	// Start a goroutine to feed all the blocks into the tracers

	go func() {
		var (
			logged  time.Time
			begin   = time.Now()
			number  uint64
			traced  uint64
			failed  error
			statedb *state.StateDB
			release StateReleaseFunc
		)
		// Ensure everything is properly cleaned up on any exit path
		defer func() {
			close(tasks)
			pend.Wait()

			// Clean out any pending derefs.
			reler.call()

			// Log the chain result
			switch {
			case failed != nil:
				logger.Warn("Chain tracing failed", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin), "err", failed)
			case number < end.NumberU64():
				logger.Warn("Chain tracing aborted", "start", start.NumberU64(), "end", end.NumberU64(), "abort", number, "transactions", traced, "elapsed", time.Since(begin))
			default:
				logger.Info("Chain tracing finished", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin))
			}
			close(results)
		}()

		// Feed all the blocks both into the tracer, as well as fast process concurrently
		for number = start.NumberU64(); number < end.NumberU64(); number++ {
			if notifier != nil {
				// Stop tracing if interruption was requested
				select {
				case <-notifier.Closed():
					return
				default:
				}
			}
			// Print progress logs if long enough time elapsed
			if time.Since(logged) > log.StatsReportLimit {
				logged = time.Now()
				logger.Info("Tracing chain segment", "start", start.NumberU64(), "end", end.NumberU64(), "current", number, "transactions", traced, "elapsed", time.Since(begin))
			}
			next, err := api.blockByNumber(localctx, rpc.BlockNumber(number+1))
			if err != nil {
				failed = err
				break
			}

			// Prepare the statedb for tracing. Don't use the live database for
			// tracing to avoid persisting state junks into the database. Switch
			// over to `preferDisk` mode only if the memory usage exceeds the
			// limit, the trie database will be reconstructed from scratch only
			// if the relevant state is available in disk.
			var preferDisk bool
			if statedb != nil {
				s1, s2, s3 := statedb.Database().TrieDB().Size()
				preferDisk = s1+s2+s3 > defaultTracechainMemLimit
			}
			_, _, _, statedb, release, err = api.backend.StateAtTransaction(localctx, next, 0, reexec, statedb, false, preferDisk)
			if err != nil {
				failed = err
				break
			}

			// Clean out any pending derefs. Note this step must be done after
			// constructing tracing state, because the tracing state of block
			// next depends on the parent state and construction may fail if
			// we release too early.
			reler.call()

			// Send the block over to the concurrent tracers (if not in the fast-forward phase)
			txs := next.Transactions()
			if notifier != nil {
				select {
				case tasks <- &blockTraceTask{statedb: statedb.Copy(), block: next, release: release, results: make([]*txTraceResult, len(txs))}:
				case <-notifier.Closed():
					return
				}
			} else {
				tasks <- &blockTraceTask{statedb: statedb.Copy(), block: next, release: release, results: make([]*txTraceResult, len(txs))}
			}
			traced += uint64(len(txs))
		}
	}()

	waitForResult := func() map[uint64]*blockTraceResult {
		// Keep reading the trace results and stream the to the user
		var (
			done = make(map[uint64]*blockTraceResult)
			next = start.NumberU64() + 1
		)
		for res := range results {
			// Queue up next received result
			result := &blockTraceResult{
				Block:  hexutil.Uint64(res.block.NumberU64()),
				Hash:   res.block.Hash(),
				Traces: res.results,
			}
			done[uint64(result.Block)] = result

			if notifier != nil {
				// Stream completed traces to the user, aborting on the first error
				for result, ok := done[next]; ok; result, ok = done[next] {
					if len(result.Traces) > 0 || next == end.NumberU64() {
						notifier.Notify(sub.ID, result)
					}
					delete(done, next)
					next++
				}
			} else {
				if len(done) == blocks {
					return done
				}
			}
		}
		return nil
	}

	if notifier != nil {
		// Keep reading the trace results and stream them to result channel.
		go waitForResult()
		return nil, nil
	}

	return waitForResult(), nil
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *API) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
	block, err := api.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockByNumberRange returns the ranged blocks tracing results
// TODO-tracer: limit the result by the size of the return
func (api *UnsafeAPI) TraceBlockByNumberRange(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (map[uint64]*blockTraceResult, error) {
	// When the block range is [start,end], the actual tracing block would be [start+1,end]
	// this is the reason why we change the block range to [start-1, end] so that we can trace [start,end] blocks
	from, to, err := checkRangeAndReturnBlock(&api.CommonAPI, ctx, start-1, end)
	if err != nil {
		return nil, err
	}
	return api.traceChain(from, to, config, nil, nil)
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *API) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *CommonAPI) TraceBlock(ctx context.Context, blob hexutil.Bytes, config *TraceConfig) ([]*txTraceResult, error) {
	block := new(types.Block)
	if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
		return nil, fmt.Errorf("could not decode block: %v", err)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *UnsafeAPI) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([]*txTraceResult, error) {
	blob, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, common.Hex2Bytes(string(blob)), config)
}

// TraceBadBlock returns the structured logs created during the execution of
// EVM against a block pulled from the pool of bad ones and returns them as a JSON
// object.
func (api *API) TraceBadBlock(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	blocks, err := api.backend.ChainDB().ReadAllBadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.traceBlock(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// StandardTraceBlockToFile dumps the structured logs created during the
// execution of EVM to the local file system and returns a list of files
// to the caller.
func (api *UnsafeAPI) StandardTraceBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("block %#x not found", hash)
	}
	return api.standardTraceBlockToFile(ctx, block, config)
}

// StandardTraceBadBlockToFile dumps the structured logs created during the
// execution of EVM against a block pulled from the pool of bad ones to the
// local file system and returns a list of files to the caller.
func (api *UnsafeAPI) StandardTraceBadBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	blocks, err := api.backend.ChainDB().ReadAllBadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.standardTraceBlockToFile(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// traceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *CommonAPI) traceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	if block.NumberU64() == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	_, _, _, statedb, release, err := api.backend.StateAtTransaction(ctx, block, 0, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	// Execute all the transaction contained within the block concurrently
	var (
		signer  = types.MakeSigner(api.backend.ChainConfig(), block.Number())
		txs     = block.Transactions()
		results = make([]*txTraceResult, len(txs))

		pend = new(sync.WaitGroup)
		jobs = make(chan *txTraceTask, len(txs))

		header   = block.Header()
		blockCtx = blockchain.NewEVMBlockContext(header, newChainContext(ctx, api.backend), nil)
	)

	threads := runtime.NumCPU()
	if threads > len(txs) {
		threads = len(txs)
	}
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next transaction trace tasks
			for task := range jobs {
				msg, err := txs[task.index].AsMessageWithAccountKeyPicker(signer, task.statedb, block.NumberU64())
				if err != nil {
					logger.Warn("Tracing failed", "tx idx", task.index, "block", block.NumberU64(), "err", err)
					results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Error: err.Error()}
					continue
				}

				txCtx := blockchain.NewEVMTxContext(msg, block.Header(), api.backend.ChainConfig())
				res, err := api.traceTx(ctx, msg, blockCtx, txCtx, task.statedb, config)
				if err != nil {
					results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Error: err.Error()}
					continue
				}
				results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Result: res}
			}
		}()
	}
	// Feed the transactions into the tracers and return
	var failed error
	for i, tx := range txs {
		// Send the trace task over for execution
		jobs <- &txTraceTask{statedb: statedb.Copy(), index: i}

		// Generate the next state snapshot fast without tracing
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("Tracing failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			failed = err
			break
		}

		txCtx := blockchain.NewEVMTxContext(msg, block.Header(), api.backend.ChainConfig())
		blockCtx := blockchain.NewEVMBlockContext(block.Header(), newChainContext(ctx, api.backend), nil)
		vmenv := vm.NewEVM(blockCtx, txCtx, statedb, api.backend.ChainConfig(), &vm.Config{})
		if _, err = blockchain.ApplyMessage(vmenv, msg); err != nil {
			failed = err
			break
		}
		// Finalize the state so any modifications are written to the trie
		statedb.Finalise(true, true)
	}
	close(jobs)
	pend.Wait()

	// If execution failed in between, abort
	if failed != nil {
		return nil, failed
	}
	return results, nil
}

// standardTraceBlockToFile configures a new tracer which uses standard JSON output,
// and traces either a full block or an individual transaction. The return value will
// be one filename per transaction traced.
func (api *CommonAPI) standardTraceBlockToFile(ctx context.Context, block *types.Block, config *StdTraceConfig) ([]string, error) {
	// If we're tracing a single transaction, make sure it's present
	if config != nil && !common.EmptyHash(config.TxHash) {
		if !containsTx(block, config.TxHash) {
			return nil, fmt.Errorf("transaction %#x not found in block", config.TxHash)
		}
	}
	if block.NumberU64() == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	parent, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(block.NumberU64()-1), block.ParentHash())
	if err != nil {
		return nil, err
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	_, _, _, statedb, release, err := api.backend.StateAtTransaction(ctx, parent, 0, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	// Retrieve the tracing configurations, or use default values
	var (
		logConfig vm.LogConfig
		txHash    common.Hash
	)
	if config != nil {
		if config.LogConfig != nil {
			logConfig = *config.LogConfig
		}
		txHash = config.TxHash
	}
	logConfig.Debug = true

	header := block.Header()
	blockCtx := blockchain.NewEVMBlockContext(header, newChainContext(ctx, api.backend), nil)

	// Execute transaction, either tracing all or just the requested one
	var (
		signer = types.MakeSigner(api.backend.ChainConfig(), block.Number())
		dumps  []string
	)
	for i, tx := range block.Transactions() {
		// Prepare the transaction for un-traced execution
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("Tracing failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			return nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}

		var (
			txCtx = blockchain.NewEVMTxContext(msg, block.Header(), api.backend.ChainConfig())

			vmConf vm.Config
			dump   *os.File
		)

		// If the transaction needs tracing, swap out the configs
		if tx.Hash() == txHash || common.EmptyHash(txHash) {
			// Generate a unique temporary file to dump it into
			prefix := fmt.Sprintf("block_%#x-%d-%#x-", block.Hash().Bytes()[:4], i, tx.Hash().Bytes()[:4])

			dump, err = os.CreateTemp(os.TempDir(), prefix)
			if err != nil {
				return nil, err
			}
			dumps = append(dumps, dump.Name())

			// Swap out the noop logger to the standard tracer
			vmConf = vm.Config{
				Debug:                   true,
				Tracer:                  vm.NewJSONLogger(&logConfig, bufio.NewWriter(dump)),
				EnablePreimageRecording: true,
			}
		}
		// Execute the transaction and flush any traces to disk
		vmenv := vm.NewEVM(blockCtx, txCtx, statedb, api.backend.ChainConfig(), &vmConf)
		_, err = blockchain.ApplyMessage(vmenv, msg)

		if dump != nil {
			dump.Close()
			logger.Info("Wrote standard trace", "file", dump.Name())
		}
		if err != nil {
			return dumps, err
		}
		// Finalize the state so any modifications are written to the trie
		statedb.Finalise(true, true)

		// If we've traced the transaction we were looking for, abort
		if tx.Hash() == txHash {
			break
		}
	}
	return dumps, nil
}

// containsTx reports whether the transaction with a certain hash
// is contained within the specified block.
func containsTx(block *types.Block, hash common.Hash) bool {
	for _, tx := range block.Transactions() {
		if tx.Hash() == hash {
			return true
		}
	}
	return false
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *CommonAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	// Retrieve the transaction and assemble its EVM context
	tx, blockHash, blockNumber, index := api.backend.GetTxAndLookupInfo(hash)
	if tx == nil {
		return nil, fmt.Errorf("transaction %#x not found", hash)
	}
	// It shouldn't happen in practice.
	if blockNumber == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	block, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(blockNumber), blockHash)
	if err != nil {
		return nil, err
	}
	msg, blockCtx, txCtx, statedb, release, err := api.backend.StateAtTransaction(ctx, block, int(index), reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	// Trace the transaction and return
	return api.traceTx(ctx, msg, blockCtx, txCtx, statedb, config)
}

// TraceCall lets you trace a given kaia_call. It collects the structured logs
// created during the execution of EVM if the given transaction was added on
// top of the provided block and returns them as a JSON object.
func (api *CommonAPI) TraceCall(ctx context.Context, args kaiaapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	// Try to retrieve the specified block
	var (
		err   error
		block *types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = api.blockByHash(ctx, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		block, err = api.blockByNumber(ctx, number)
	} else {
		return nil, errors.New("invalid arguments; neither block nor hash specified")
	}
	if err != nil {
		return nil, err
	}
	// try to recompute the state
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	statedb, release, err := api.backend.StateAtBlock(ctx, block, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	// Execute the trace
	intrinsicGas, err := types.IntrinsicGas(args.InputData(), args.GetAccessList(), args.To == nil, api.backend.ChainConfig().Rules(block.Number()))
	if err != nil {
		return nil, err
	}
	basefee := new(big.Int).SetUint64(params.ZeroBaseFee)
	if block.Header().BaseFee != nil {
		basefee = block.Header().BaseFee
	}
	gasCap := uint64(0)
	if rpcGasCap := api.backend.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	msg, err := args.ToMessage(gasCap, basefee, intrinsicGas)
	if err != nil {
		return nil, err
	}

	// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
	statedb.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), basefee))

	txCtx := blockchain.NewEVMTxContext(msg, block.Header(), api.backend.ChainConfig())
	blockCtx := blockchain.NewEVMBlockContext(block.Header(), newChainContext(ctx, api.backend), nil)

	return api.traceTx(ctx, msg, blockCtx, txCtx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (api *CommonAPI) traceTx(ctx context.Context, message blockchain.Message, blockCtx vm.BlockContext, txCtx vm.TxContext, statedb *state.StateDB, config *TraceConfig) (interface{}, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer vm.Tracer
		err    error
	)
	switch {
	case config != nil && config.Tracer != nil:
		// Define a meaningful timeout of a single transaction trace
		timeout := defaultTraceTimeout
		if config.Timeout != nil {
			if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
				return nil, err
			}
		}

		if *config.Tracer == "fastCallTracer" || *config.Tracer == "callTracer" {
			tracer = vm.NewCallTracer()
		} else {
			// Construct the JavaScript tracer to execute with
			if tracer, err = New(*config.Tracer, new(Context), api.unsafeTrace); err != nil {
				return nil, err
			}
		}
		// Handle timeouts and RPC cancellations
		deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-deadlineCtx.Done()
			if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
				switch t := tracer.(type) {
				case *Tracer:
					t.Stop(errors.New("execution timeout"))
				case *vm.InternalTxTracer:
					t.Stop(errors.New("execution timeout"))
				case *vm.CallTracer:
					t.Stop(errors.New("execution timeout"))
				default:
					logger.Warn("unknown tracer type", "type", reflect.TypeOf(t).String())
				}
			}
		}()
		defer cancel()

	case config == nil:
		tracer = vm.NewStructLogger(nil)

	default:
		tracer = vm.NewStructLogger(config.LogConfig)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(blockCtx, txCtx, statedb, api.backend.ChainConfig(), &vm.Config{Debug: true, Tracer: tracer})

	ret, err := blockchain.ApplyMessage(vmenv, message)
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %v", err)
	}
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		loggerTimeout := defaultLoggerTimeout
		if config != nil && config.LoggerTimeout != nil {
			if loggerTimeout, err = time.ParseDuration(*config.LoggerTimeout); err != nil {
				return nil, err
			}
		}
		if logs, err := kaiaapi.FormatLogs(loggerTimeout, tracer.StructLogs()); err == nil {
			return &kaiaapi.ExecutionResult{
				Gas:         ret.UsedGas,
				Failed:      ret.Failed(),
				ReturnValue: fmt.Sprintf("%x", ret.Return()),
				StructLogs:  logs,
			}, nil
		} else {
			return nil, err
		}

	case *Tracer:
		return tracer.GetResult()
	case *vm.InternalTxTracer:
		return tracer.GetResult()
	case *vm.CallTracer:
		return tracer.GetResult()

	default:
		panic(fmt.Sprintf("bad tracer type %T", tracer))
	}
}
