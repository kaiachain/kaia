package blockchain

import (
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/internal/poolobserve"
	kaiamock "github.com/kaiachain/kaia/kaiax/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/prometheus/client_golang/prometheus"
)

func propagationEvent(t *testing.T, r *prometheus.Registry, event string) float64 {
	t.Helper()
	fs, err := r.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if f.GetName() == "kaia_poolobs_events_total" {
			for _, m := range f.Metric {
				if m.Label[0].GetValue() == event {
					return m.Counter.GetValue()
				}
			}
		}
	}
	t.Fatalf("missing event %s", event)
	return 0
}

func TestPoolObservabilityRemoteEnqueueBlocked(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	pool := &TxPool{txMsgCh: make(chan types.Transactions, 1)}
	pool.txMsgCh <- nil
	done := make(chan struct{})
	go func() { pool.HandleTxMsg(make([]*types.Transaction, 2)); close(done) }()
	defer func() {
		select {
		case <-pool.txMsgCh:
		default:
		}
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Error("enqueue did not drain")
		}
	}()
	deadline := time.Now().Add(3 * time.Second)
	for poolObservationValue(t, r, "kaia_poolobs_operations_inflight", "remote_enqueue") != 1 {
		if time.Now().After(deadline) {
			t.Fatal("enqueue did not block")
		}
		time.Sleep(time.Millisecond)
	}
	if n := propagationEvent(t, r, "remote_received"); n != 2 {
		t.Fatalf("received %g", n)
	}
	if n := propagationEvent(t, r, "admitted"); n != 0 {
		t.Fatalf("enqueue incorrectly counted as admission: %g", n)
	}
	<-pool.txMsgCh
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("enqueue did not complete")
	}
	if n := poolObservationValue(t, r, "kaia_poolobs_duration_seconds", "remote_enqueue"); n != 1 {
		t.Fatalf("completed %g", n)
	}
}

func TestPoolObservabilityPublicAdmissionCounts(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	pool.Stop()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	txs := make([]*types.Transaction, 4)
	for i := range txs {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
		txs[i] = pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	}
	if err := pool.AddLocal(txs[0]); err != nil {
		t.Fatal(err)
	}
	if err := pool.AddRemote(txs[1]); err != nil {
		t.Fatal(err)
	}
	errs := pool.AddLocals([]*types.Transaction{txs[0], txs[2]})
	if errs[0] == nil || errs[1] != nil {
		t.Fatalf("local batch: %v", errs)
	}
	errs = pool.AddRemotes([]*types.Transaction{txs[1], txs[3]})
	if errs[0] == nil || errs[1] != nil {
		t.Fatalf("remote batch: %v", errs)
	}
	if n := propagationEvent(t, r, "admitted"); n != 4 {
		t.Fatalf("successful attempts %g", n)
	}
	if n := propagationEvent(t, r, "admission_rejected"); n != 2 {
		t.Fatalf("rejected attempts %g", n)
	}
	if n := propagationEvent(t, r, "promoted"); n != 4 {
		t.Fatalf("promoted %g", n)
	}
	if n := propagationEvent(t, r, "remote_admitted"); n != 2 {
		t.Fatalf("remote successful attempts %g", n)
	}
	if n := propagationEvent(t, r, "remote_rejected"); n != 1 {
		t.Fatalf("remote rejected attempts %g", n)
	}
	if n := propagationEvent(t, r, "capacity_rejected"); n != 0 {
		t.Fatalf("duplicates misclassified as capacity %g", n)
	}
	if n := poolObservationValue(t, r, "kaia_poolobs_duration_seconds", "admission_hold"); n != 4 {
		t.Fatalf("lock holds %g", n)
	}
}

func TestPoolObservabilityFeedBlocksAdmissionCompletion(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	pool.Stop()
	pool.txFeedCh = make(chan types.Transactions, 1)
	pool.txFeedCh <- nil
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
	tx := pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	done := make(chan error, 1)
	finished := make(chan struct{})
	go func() { done <- pool.AddLocal(tx); close(finished) }()
	defer func() {
		select {
		case <-pool.txFeedCh:
		default:
		}
		select {
		case <-finished:
		case <-time.After(3 * time.Second):
			t.Error("admission did not stop")
		}
	}()
	deadline := time.Now().Add(3 * time.Second)
	for poolObservationValue(t, r, "kaia_poolobs_operations_inflight", "feed_enqueue") != 1 {
		if time.Now().After(deadline) {
			t.Fatal("feed enqueue did not block")
		}
		time.Sleep(time.Millisecond)
	}
	if n := poolObservationValue(t, r, "kaia_poolobs_operations_inflight", "admission_hold"); n != 1 {
		t.Fatalf("admission hold %g", n)
	}
	if n := propagationEvent(t, r, "promoted"); n != 1 {
		t.Fatalf("promoted %g", n)
	}
	if n := propagationEvent(t, r, "admitted"); n != 0 {
		t.Fatalf("unfinished public admission counted %g", n)
	}
	<-pool.txFeedCh
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("admission did not drain")
	}
	if n := propagationEvent(t, r, "admitted"); n != 1 {
		t.Fatalf("completed public admission %g", n)
	}
}

// Use the real pool loop and reset, with barriers only in its existing chain
// provider and module interfaces. The chain head read remains concurrency-safe.
type stagedPoolHeadChain struct {
	*testBlockChain
	head                       atomic.Pointer[types.Block]
	blockState                 atomic.Bool
	stateEntered, stateRelease chan struct{}
}

func (c *stagedPoolHeadChain) CurrentBlock() *types.Block { return c.head.Load() }
func (c *stagedPoolHeadChain) StateAt(root common.Hash) (*state.StateDB, error) {
	if c.blockState.Load() {
		close(c.stateEntered)
		<-c.stateRelease
	}
	return c.testBlockChain.StateAt(root)
}

func TestPoolObservabilityHeadResetStages(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	oldPool, _ := setupTxPool()
	oldPool.Stop()
	chain := &stagedPoolHeadChain{testBlockChain: oldPool.chain.(*testBlockChain), stateEntered: make(chan struct{}), stateRelease: make(chan struct{})}
	oldHead := chain.testBlockChain.CurrentBlock()
	chain.head.Store(oldHead)
	pool := NewTxPool(testTxPoolConfig, params.TestChainConfig, chain, oldPool.govModule)
	postEntered, postRelease := make(chan struct{}), make(chan struct{})
	var stateOnce, postOnce sync.Once
	releaseState := func() { stateOnce.Do(func() { close(chain.stateRelease) }) }
	releasePost := func() { postOnce.Do(func() { close(postRelease) }) }
	defer func() { releaseState(); releasePost(); pool.Stop() }()
	module := kaiamock.NewMockTxPoolModule(gomock.NewController(t))
	module.EXPECT().PreReset(gomock.Any(), gomock.Any()).Return(nil)
	module.EXPECT().PostReset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(*types.Header, *types.Header, map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
		close(postEntered)
		<-postRelease
	})
	pool.RegisterTxPoolModule(module)
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	value := func(name string) float64 {
		fs, err := registry.Gather()
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range fs {
			if f.GetName() == name {
				return f.Metric[0].Gauge.GetValue()
			}
		}
		t.Fatalf("missing %s", name)
		return 0
	}
	wait := func(name string, want float64) {
		deadline := time.Now().Add(3 * time.Second)
		for value(name) != want {
			if time.Now().After(deadline) {
				t.Fatalf("%s did not reach %g", name, want)
			}
			time.Sleep(time.Millisecond)
		}
	}
	wait("kaia_poolobs_chain_head_observed", 1)
	head := types.NewBlock(&types.Header{Number: big.NewInt(1), ParentHash: oldHead.Hash()}, nil, nil)
	chain.blockState.Store(true)
	chain.head.Store(head)
	chain.chainHeadFeed.Send(ChainHeadEvent{Block: head})
	select {
	case <-chain.stateEntered:
	case <-time.After(3 * time.Second):
		t.Fatal("reset did not request new state")
	}
	if value("kaia_poolobs_chain_head_number") != 1 || value("kaia_poolobs_pool_state_head_number") != 0 || value("kaia_poolobs_pool_reset_head_number") != 0 {
		t.Fatal("state loading did not preserve old pool head")
	}
	releaseState()
	select {
	case <-postEntered:
	case <-time.After(3 * time.Second):
		t.Fatal("reset did not reach PostReset")
	}
	if value("kaia_poolobs_pool_state_head_number") != 1 || value("kaia_poolobs_pool_reset_head_number") != 0 {
		t.Fatal("state installation and complete reset were not distinguished")
	}
	releasePost()
	wait("kaia_poolobs_pool_reset_head_number", 1)
}

func TestPoolObservabilityCapacityRejections(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	pool.Stop()
	pool.config.ExecSlotsAll = 1
	pool.config.NonExecSlotsAll = 1
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	txs := make([]*types.Transaction, 3)
	for i := range txs {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
		txs[i] = pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	}
	if err := pool.AddLocal(txs[0]); err != nil {
		t.Fatal(err)
	}
	errs := pool.AddRemotes(txs[1:])
	if len(errs) != 2 || errs[0] != nil || errs[1] != ErrTxPoolOverflow {
		t.Fatalf("partial batch: %v", errs)
	}
	if n := propagationEvent(t, r, "capacity_rejected"); n != 1 {
		t.Fatalf("partial capacity count %g", n)
	}
	if err := pool.AddLocal(txs[2]); err == nil {
		t.Fatal("full local admission succeeded")
	}
	if n := propagationEvent(t, r, "capacity_rejected"); n != 2 {
		t.Fatalf("local precheck count %g", n)
	}
	if err := pool.AddRemote(txs[2]); err == nil {
		t.Fatal("full remote admission succeeded")
	}
	if n := propagationEvent(t, r, "capacity_rejected"); n != 3 {
		t.Fatalf("remote capacity count %g", n)
	}
	if n := propagationEvent(t, r, "remote_admitted"); n != 1 {
		t.Fatalf("remote admitted %g", n)
	}
	if n := propagationEvent(t, r, "remote_rejected"); n != 2 {
		t.Fatalf("remote rejected %g", n)
	}
	if n := propagationEvent(t, r, "admitted"); n != 2 {
		t.Fatalf("all admitted %g", n)
	}
	if n := propagationEvent(t, r, "admission_rejected"); n != 3 {
		t.Fatalf("all rejected %g", n)
	}
}

func TestPoolObservabilityRemotePolicyPaths(t *testing.T) {
	t.Run("denied", func(t *testing.T) {
		defer poolobserve.EnableForTesting()()
		r := prometheus.NewRegistry()
		poolobserve.Register(r)
		pool := &TxPool{config: TxPoolConfig{DenyRemoteTx: true}}
		pool.HandleTxMsg(make(types.Transactions, 2))
		if n := propagationEvent(t, r, "remote_received"); n != 2 {
			t.Fatalf("received %g", n)
		}
		if n := propagationEvent(t, r, "remote_denied"); n != 2 {
			t.Fatalf("denied %g", n)
		}
		if n := propagationEvent(t, r, "remote_rejected"); n != 0 {
			t.Fatalf("policy denial is not attempted admission: %g", n)
		}
	})
	t.Run("throttled", func(t *testing.T) {
		defer poolobserve.EnableForTesting()()
		pool, key := setupTxPool()
		pool.Stop()
		if err := pool.AddLocal(pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)); err != nil {
			t.Fatal(err)
		}
		config := *DefaultSpamThrottlerConfig
		config.ActivateTxPoolSize = 0
		th := newTestThrottler(&config)
		th.throttleCh = make(chan *types.Transaction, 1)
		target := common.HexToAddress("0xbeef")
		th.throttled[target] = 10
		spamThrottlerMu.Lock()
		previous := spamThrottler
		spamThrottler = th
		spamThrottlerMu.Unlock()
		defer func() { spamThrottlerMu.Lock(); spamThrottler = previous; spamThrottlerMu.Unlock() }()
		r := prometheus.NewRegistry()
		poolobserve.Register(r)
		dropsBefore := throttlerDropCount.Count()
		tx := types.NewTransaction(0, target, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
		pool.HandleTxMsg(types.Transactions{tx, tx})
		if n := propagationEvent(t, r, "remote_received"); n != 2 {
			t.Fatalf("received %g", n)
		}
		if n := propagationEvent(t, r, "remote_throttled"); n != 1 {
			t.Fatalf("throttle transfers %g", n)
		}
		if n := propagationEvent(t, r, "remote_admitted"); n != 0 {
			t.Fatalf("transfer is not admission: %g", n)
		}
		if len(th.throttleCh) != 1 || throttlerDropCount.Count()-dropsBefore != 1 {
			t.Fatal("throttle queue transfer and overflow mismatch")
		}
	})
}
