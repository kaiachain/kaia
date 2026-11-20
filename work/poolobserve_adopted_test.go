package work

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	consensusmocks "github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/prometheus/client_golang/prometheus"
)

type adoptedTestSealer struct {
	consensus.Sealer
	author common.Address
	err    error
}

func (s adoptedTestSealer) Author(*types.Header) (common.Address, error) { return s.author, s.err }

func adoptedTestSnapshots() (*poolobserve.CandidateObservation, *poolobserve.CandidateObservation) {
	p, e := poolobserve.NewCandidateObservation(), poolobserve.NewCandidateObservation()
	p.SetPhase(poolobserve.CandidatePoolInput, 4)
	p.SetPhase(poolobserve.CandidateAfterFee, 3)
	p.SetPhase(poolobserve.CandidateAfterModules, 2)
	e.SetPhase(poolobserve.CandidateExecutionInput, 2)
	e.SetPhase(poolobserve.CandidateIncluded, 0)
	e.CountOutcome(poolobserve.CandidateNonceHigh)
	return p, e
}

func TestAdoptedCandidateResultClassification(t *testing.T) {
	for _, name := range []string{"canonical", "known_canonical", "known_side", "side", "write_error", "foreign_author", "author_error", "missing_execution", "missing_preparation", "number_mismatch", "parent_mismatch", "tx_count_mismatch"} {
		t.Run(name, func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			reg := prometheus.NewRegistry()
			poolobserve.Register(reg)
			chain := mocks.NewMockBlockChain(gomock.NewController(t))
			h := &types.Header{Number: big.NewInt(5), Time: big.NewInt(5), ParentHash: common.Hash{3}}
			b := types.NewBlockWithHeader(h)
			prep, exec := adoptedTestSnapshots()
			exec.RecordTimeLimit(true)
			task := &Task{header: h, candidateObservation: prep}
			result := &consensus.ExecutionResult{Block: b, CandidateObservation: exec, ExecuteTime: 180 * time.Millisecond}
			w := &worker{chain: chain, nodeAddr: common.Address{7}}
			status := blockchain.WriteResult{Status: blockchain.CanonStatTy}
			var err error
			want := "observed"
			switch name {
			case "known_canonical", "known_side":
				err = blockchain.ErrKnownBlock
				canonical := h
				if name == "known_side" {
					canonical = &types.Header{Number: big.NewInt(5), Time: big.NewInt(6)}
					want = "not_adopted"
				}
				chain.EXPECT().GetHeaderByNumber(uint64(5)).Return(canonical)
			case "side":
				status.Status = blockchain.SideStatTy
				want = "not_adopted"
			case "write_error":
				err = errors.New("write failed")
				want = "not_adopted"
			case "missing_execution":
				result.CandidateObservation = nil
				want = "unobserved"
			case "missing_preparation":
				task.candidateObservation = nil
				want = "unobserved"
			case "number_mismatch":
				task.header = &types.Header{Number: big.NewInt(4), ParentHash: h.ParentHash}
				want = "unobserved"
			case "tx_count_mismatch":
				result.Txs = []*types.Transaction{types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)}
				want = "unobserved"
			case "parent_mismatch":
				task.header = &types.Header{Number: big.NewInt(5), ParentHash: common.Hash{8}}
				want = "unobserved"
			}
			if want != "not_adopted" {
				sealer := adoptedTestSealer{author: w.nodeAddr}
				if name == "foreign_author" {
					sealer.author = common.Address{9}
					want = "not_adopted"
				}
				if name == "author_error" {
					sealer.err = errors.New("invalid author seal")
					want = "not_adopted"
				}
				chain.EXPECT().Sealer().Return(sealer)
			}
			w.observeAdoptedCandidate(task, result, status, err)
			for _, kind := range []string{"observed", "unobserved", "not_adopted"} {
				expected := float64(0)
				if kind == want {
					expected = 1
				}
				if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_blocks_total", "result", kind); got != expected {
					t.Fatalf("%s=%v, want %v", kind, got, expected)
				}
			}
			expected := float64(0)
			if want == "observed" {
				expected = 4
			}
			wantTimeout := float64(0)
			if want == "observed" {
				wantTimeout = 1
			}
			if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_outcomes_total", "reason", "time_limit_with_remaining"); got != wantTimeout {
				t.Fatalf("published time limits = %v want %v", got, wantTimeout)
			}
			if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_transactions_total", "phase", "pool_input"); got != expected {
				t.Fatalf("detail=%v, want %v", got, expected)
			}
		})
	}
}

func TestAdoptedCandidateDisabledDoesNotReadChain(t *testing.T) {
	if poolobserve.Enabled() {
		t.Skip("requires KAIA_POOLOBSERVE to be disabled at process start")
	}
	// A canonical result would require an author lookup without the disabled guard.
	// Strict mocks reject that lookup, rather than taking an unrelated early return.
	chain := mocks.NewMockBlockChain(gomock.NewController(t))
	w := &worker{chain: chain}
	header := &types.Header{Number: big.NewInt(1), Time: big.NewInt(1)}
	result := &consensus.ExecutionResult{Block: types.NewBlockWithHeader(header)}
	w.observeAdoptedCandidate(&Task{header: header}, result, blockchain.WriteResult{Status: blockchain.CanonStatTy}, nil)
}

func TestAdoptedCandidateFinalizationCountsOnce(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	reg := prometheus.NewRegistry()
	poolobserve.Register(reg)
	chain := mocks.NewMockBlockChain(gomock.NewController(t))
	db, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	h := &types.Header{Number: big.NewInt(1), Time: big.NewInt(1)}
	b := types.NewBlockWithHeader(h)
	prep, exec := adoptedTestSnapshots()
	w := &worker{chain: chain, nodeAddr: common.Address{1}, mux: new(event.TypeMux), pendingWork: &Task{header: h, state: db, candidateObservation: prep}}
	result := &consensus.ExecutionResult{Block: b, State: db, CandidateObservation: exec, ExecuteTime: time.Millisecond}
	chain.EXPECT().WriteBlockWithState(b, gomock.Any(), db).Return(blockchain.WriteResult{Status: blockchain.CanonStatTy}, nil).Times(1)
	chain.EXPECT().Sealer().Return(adoptedTestSealer{author: w.nodeAddr}).Times(1)
	chain.EXPECT().PostChainEvents(gomock.Any(), gomock.Any()).Times(1)
	w.handleFinalizedBlock(nil)
	w.handleFinalizedBlock(&consensus.ExecutionResult{})
	w.handleFinalizedBlock(result)
	w.handleFinalizedBlock(result)
	if w.pendingWork != nil || w.finalizeCh != nil {
		t.Fatal("completed channel pair retained")
	}
	if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_blocks_total", "result", "observed"); got != 1 {
		t.Fatalf("observed=%v", got)
	}
}

func TestAdoptedExecutionSnapshotIsPerCall(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	reg := prometheus.NewRegistry()
	poolobserve.Register(reg)
	db, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	config := &params.ChainConfig{ChainID: big.NewInt(1)}
	signer := types.LatestSignerForChainID(config.ChainID)
	e := &DefaultExecutor{config: config, chain: &candidateMetricsChain{}, state: db, header: &types.Header{Number: big.NewInt(1), Time: big.NewInt(1), BlockScore: big.NewInt(1)}, signer: signer, initialized: true}
	execute := func() *consensus.ExecutionResult {
		r, err := e.Execute(types.NewTransactionsByPriceAndNonce(signer, map[common.Address]types.Transactions{}, nil), new(event.TypeMux))
		if err != nil {
			t.Fatal(err)
		}
		if r.CandidateObservation == nil {
			t.Fatal("missing execution snapshot")
		}
		return r
	}
	first := execute()
	first.CandidateObservation.CountOutcome(poolobserve.CandidateNonceLow)
	second := execute()
	if first.CandidateObservation == second.CandidateObservation {
		t.Fatal("snapshot shared between calls")
	}
	prep, _ := adoptedTestSnapshots()
	poolobserve.ObserveAdopted(prep, second.CandidateObservation, time.Millisecond, true)
	if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_outcomes_total", "reason", "nonce_low"); got != 0 {
		t.Fatalf("previous result leaked: %v", got)
	}
}

type adoptedTestProcessor struct{ blockchain.Processor }

func (adoptedTestProcessor) InitializeState(*types.Header, *state.StateDB) error { return nil }

func TestAdoptedReplayDoesNotCarryExecutionSnapshot(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	chain := mocks.NewMockBlockChain(gomock.NewController(t))
	chain.EXPECT().Processor().Return(adoptedTestProcessor{})
	chain.EXPECT().Sealer().Return(adoptedTestSealer{})
	e := &DefaultExecutor{chain: chain, header: &types.Header{Number: big.NewInt(1)}, initialized: true}
	result, err := e.ExecuteTransactions(nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.CandidateObservation != nil {
		t.Fatal("validator replay returned local candidate snapshot")
	}
}

func TestAdoptedReplacedResultChannelIsNotConsumed(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	reg := prometheus.NewRegistry()
	poolobserve.Register(reg)
	ctrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(ctrl)
	engine := consensusmocks.NewMockEngine(ctrl)
	pool := mocks.NewMockTxPool(ctrl)
	db, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	parent := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0), Time: big.NewInt(0), BlockScore: big.NewInt(1)})
	chain.EXPECT().CurrentBlock().Return(parent).Times(2)
	chain.EXPECT().Validator().Return(&fakeValidator{}).Times(2)
	pool.EXPECT().Pending().Return(map[common.Address]types.Transactions{}, nil).Times(2)
	chain.EXPECT().PrepareHeader(gomock.Any()).DoAndReturn(func(h *types.Header) error { h.BlockScore = big.NewInt(1); return nil }).Times(2)
	chain.EXPECT().StateAt(parent.Root()).Return(db, nil).Times(2)
	chain.EXPECT().Processor().Return(adoptedTestProcessor{}).Times(2)
	var channels []chan *consensus.ExecutionResult
	var results []*consensus.ExecutionResult
	engine.EXPECT().SubmitTransactions(gomock.Any(), db, gomock.Any(), gomock.Any()).DoAndReturn(func(_ *types.TransactionsByPriceAndNonce, _ *state.StateDB, h *types.Header, _ *event.TypeMux) <-chan *consensus.ExecutionResult {
		ch := make(chan *consensus.ExecutionResult, 1)
		exec := poolobserve.NewCandidateObservation()
		exec.SetPhase(poolobserve.CandidateExecutionInput, 0)
		exec.SetPhase(poolobserve.CandidateIncluded, 0)
		channels = append(channels, ch)
		results = append(results, &consensus.ExecutionResult{State: db, Block: types.NewBlockWithHeader(h), CandidateObservation: exec, ExecuteTime: time.Millisecond})
		return ch
	}).Times(2)
	var txFeed, headFeed, sideFeed event.Feed
	mux := new(event.TypeMux)
	w := &worker{config: &params.ChainConfig{ChainID: big.NewInt(1)}, chain: chain, engine: engine, backend: &testBackend{chain: chain, txPool: pool}, mux: mux, nodeAddr: common.Address{1}, txsCh: make(chan blockchain.NewTxsEvent), chainHeadCh: make(chan blockchain.ChainHeadEvent), chainSideCh: make(chan blockchain.ChainSideEvent), newSequenceSub: mux.Subscribe(newSequenceEvent{})}
	w.txsSub = txFeed.Subscribe(w.txsCh)
	w.chainHeadSub = headFeed.Subscribe(w.chainHeadCh)
	w.chainSideSub = sideFeed.Subscribe(w.chainSideCh)
	// Both real preparations use the same parent and height. Only their owned
	// result channels differ, so a height-only selection would not distinguish them.
	w.commitNewWork()
	w.commitNewWork()
	channels[0] <- results[0]
	channels[1] <- results[1]
	chain.EXPECT().WriteBlockWithState(results[1].Block, gomock.Any(), db).Return(blockchain.WriteResult{Status: blockchain.CanonStatTy}, nil).Times(1)
	chain.EXPECT().Sealer().Return(adoptedTestSealer{author: w.nodeAddr}).Times(1)
	posted := make(chan struct{})
	chain.EXPECT().PostChainEvents(gomock.Any(), gomock.Any()).DoAndReturn(func([]interface{}, []*types.Log) { close(posted) }).Times(1)
	stopped := make(chan struct{})
	go func() { w.update(); close(stopped) }()
	select {
	case <-posted:
	case <-time.After(time.Second):
		t.Fatal("current result was not consumed")
	}
	w.chainHeadSub.Unsubscribe()
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("worker did not stop")
	}
	if len(channels[0]) != 1 {
		t.Fatal("replaced result channel was consumed")
	}
	if got := candidateCounter(t, reg, "kaia_poolobs_adopted_candidate_blocks_total", "result", "observed"); got != 1 {
		t.Fatalf("observed=%v, want 1", got)
	}
}
