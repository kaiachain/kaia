package blockchain

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
	"github.com/prometheus/client_golang/prometheus"
)

// This component test fills the actual promotion channel while all production
// pool methods remain unchanged. It proves downstream feed pressure can extend
// pool.mu ownership across independently funded sender accounts.
func TestPoolObservabilityFeedBacklogBlocksIndependentAccounts(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	pool.Stop()
	pool.txFeedCh = make(chan types.Transactions, 1)
	pool.txFeedCh <- nil
	txs := make([]*types.Transaction, 3)
	for i := range txs {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
		txs[i] = pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	}
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	results := make(chan error, 3)
	var workers sync.WaitGroup
	launch := func(f func() error) { workers.Add(1); go func() { defer workers.Done(); results <- f() }() }
	// Start draining on both success and failure, so no test failure strands a
	// pool method holding the global mutex.
	var drainOnce sync.Once
	drainDone := make(chan struct{})
	drain := func() {
		drainOnce.Do(func() {
			go func() {
				for {
					select {
					case <-pool.txFeedCh:
					case <-drainDone:
						return
					}
				}
			}()
		})
	}
	defer func() {
		drain()
		done := make(chan struct{})
		go func() { workers.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Error("feed admissions did not drain")
		}
		close(drainDone)
	}()
	launch(func() error { return pool.AddLocal(txs[0]) })
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "feed_enqueue") == 1 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "feed_enqueue") != 1 {
		t.Fatal("promotion enqueue not observed")
	}
	launch(func() error { return pool.AddLocal(txs[1]) })
	launch(func() error { return pool.AddRemote(txs[2]) })
	deadline = time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "admission_write_single") == 2 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "admission_write_single") != 2 {
		t.Fatal("independent accounts did not wait on shared lock")
	}
	select {
	case e := <-results:
		t.Fatalf("admission completed before draining feed: %v", e)
	default:
	}
	drain()
	for i := 0; i < 3; i++ {
		select {
		case e := <-results:
			if e != nil {
				t.Fatal(e)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("feed recovery timeout")
		}
	}
	pending, e := pool.Pending()
	if e != nil {
		t.Fatal(e)
	}
	count := 0
	for _, list := range pending {
		count += len(list)
	}
	if count != 3 {
		t.Fatalf("pending=%d want3", count)
	}
	t.Log("actual promotion channel full -> first account holds pool.mu -> other local and remote write Lock wait -> channel drain -> all three pending")
}

func waitPoolObservation(t *testing.T, r *prometheus.Registry, metric, stage string, want float64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if poolObservationValue(t, r, metric, stage) == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("%s{%s}=%g want%g", metric, stage, poolObservationValue(t, r, metric, stage), want)
}

func TestPoolObservabilityFeedSendAndControl(t *testing.T) {
	for _, capacity := range []int{0, 1} {
		t.Run(fmt.Sprintf("subscriber_capacity_%d", capacity), func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			pool, key := setupTxPool()
			defer pool.Stop()
			ch := make(chan NewTxsEvent, capacity)
			sub := pool.SubscribeNewTxsEvent(ch)
			defer sub.Unsubscribe()
			registry := prometheus.NewRegistry()
			poolobserve.Register(registry)
			tx := pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
			if err := pool.AddLocal(tx); err != nil {
				t.Fatal(err)
			}
			if capacity == 0 {
				waitPoolObservation(t, registry, "kaia_poolobs_operations_inflight", "feed_send", 1)
				if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "feed_enqueue") != 0 {
					t.Fatal("enqueue wait confused with subscriber wait")
				}
			} else {
				waitPoolObservation(t, registry, "kaia_poolobs_duration_seconds", "feed_send", 1)
			}
			select {
			case event := <-ch:
				if len(event.Txs) != 1 || event.Txs[0].Hash() != tx.Hash() {
					t.Fatal("wrong delivered transaction")
				}
			case <-time.After(3 * time.Second):
				t.Fatal("feed did not deliver")
			}
			waitPoolObservation(t, registry, "kaia_poolobs_operations_inflight", "feed_send", 0)
			waitPoolObservation(t, registry, "kaia_poolobs_duration_seconds", "feed_send", 1)
		})
	}
}

func TestPoolObservabilitySingletonBatch(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, key := setupTxPool()
	defer pool.Stop()
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	tx := pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	errs := pool.AddLocals([]*types.Transaction{tx})
	if len(errs) != 1 || errs[0] != nil {
		t.Fatal(errs)
	}
	if poolObservationValue(t, registry, "kaia_poolobs_operations_started_total", "admission_write_single") != 1 || poolObservationValue(t, registry, "kaia_poolobs_operations_started_total", "admission_write_batch") != 0 {
		t.Fatal("singleton fast path misclassified as batch lock")
	}
}

type observationFeeGov struct{}

func (observationFeeGov) GetParamSet(uint64) gov.ParamSet {
	return gov.ParamSet{UnitPrice: 25000000000, LowerBoundBaseFee: 25000000000, UpperBoundBaseFee: 750000000000, GasTarget: 15000000, MaxBlockGasUsedForBaseFee: 30000000, BaseFeeDenominator: 20}
}

// This preserves, rather than fixes, the existing reset fee-update order. The
// nonce-gap control prevents a later promotion from being attributed to fee alone.
func TestPoolObservabilityResetFeeOrder(t *testing.T) {
	for _, nonce := range []uint64{0, 1} {
		t.Run(fmt.Sprintf("nonce_%d", nonce), func(t *testing.T) {
			cfg := params.TestChainConfig.Copy()
			cfg.MagmaCompatibleBlock = big.NewInt(0)
			pool, key := setupTxPoolWithConfig(cfg)
			pool.Stop()
			pool.govModule = observationFeeGov{}
			tx := pricedTransaction(nonce, 100000, big.NewInt(96000000000), key)
			if err := pool.AddLocal(tx); err != nil {
				t.Fatal(err)
			}
			pool.lockedReset(nil, &types.Header{Number: big.NewInt(1), BaseFee: big.NewInt(100000000000), GasUsed: 15000000})
			pool.lockedReset(nil, &types.Header{Number: big.NewInt(2), BaseFee: big.NewInt(100000000000), GasUsed: 0})
			pending, queued := pool.Stats()
			if pending != 0 || queued != 1 || pool.gasPrice.Int64() != 95000000000 {
				t.Fatalf("after fee update: p=%d q=%d fee=%v", pending, queued, pool.gasPrice)
			}
			pool.lockedReset(nil, &types.Header{Number: big.NewInt(3), BaseFee: big.NewInt(95000000000), GasUsed: 0})
			pending, queued = pool.Stats()
			if nonce == 0 && (pending != 1 || queued != 0) {
				t.Fatal("eligible nonce0 not promoted next reset")
			}
			if nonce == 1 && (pending != 0 || queued != 1) {
				t.Fatal("nonce-gap control promoted")
			}
			t.Logf("nonce=%d old100G/new95G left queued; next reset pending=%d queued=%d", nonce, pending, queued)
		})
	}
}

func TestPoolObservabilityReplacementFeedWait(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	cfg := params.TestChainConfig.Copy()
	cfg.MagmaCompatibleBlock = big.NewInt(0)
	pool, key := setupTxPoolWithConfig(cfg)
	pool.Stop()
	pool.txFeedCh = make(chan types.Transactions, 1)
	old := pricedTransaction(0, 21000, big.NewInt(30000000000), key)
	replacement := pricedTransaction(0, 21000, big.NewInt(40000000000), key)
	if err := pool.AddLocal(old); err != nil {
		t.Fatal(err)
	}
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	done := make(chan error, 1)
	finished := make(chan struct{})
	go func() { defer close(finished); done <- pool.AddLocal(replacement) }()
	// Release the actual channel producer even if an assertion fails.
	defer func() {
		select {
		case <-pool.txFeedCh:
		default:
		}
		select {
		case <-finished:
		case <-time.After(3 * time.Second):
			t.Error("replacement producer did not stop")
		}
	}()
	waitPoolObservation(t, registry, "kaia_poolobs_operations_inflight", "feed_enqueue", 1)
	if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", "admission_write_single") != 0 {
		t.Fatal("replacement still waiting to acquire mutex")
	}
	<-pool.txFeedCh
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("replacement did not recover")
	}
	select {
	case txs := <-pool.txFeedCh:
		if len(txs) != 1 || txs[0].Hash() != replacement.Hash() {
			t.Fatal("wrong replacement event")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("replacement event missing")
	}
	waitPoolObservation(t, registry, "kaia_poolobs_operations_inflight", "feed_enqueue", 0)
	if poolObservationValue(t, registry, "kaia_poolobs_duration_seconds", "feed_enqueue") != 2 {
		t.Fatal("replacement enqueue was not counted")
	}
}
