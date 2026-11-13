package blockchain

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/kaiachain/kaia/params"
	"github.com/prometheus/client_golang/prometheus"
)

type poolObservationChain struct {
	blockChain
	entered, release chan struct{}
}

func (c *poolObservationChain) StateAt(root common.Hash) (*state.StateDB, error) {
	close(c.entered)
	<-c.release
	return c.blockChain.StateAt(root)
}

func poolObservationValue(t *testing.T, r *prometheus.Registry, name, stage string) float64 {
	t.Helper()
	fs, err := r.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if f.GetName() == name {
			for _, m := range f.Metric {
				for _, l := range m.Label {
					if l.GetName() == "stage" && l.GetValue() == stage {
						if m.Counter != nil {
							return m.Counter.GetValue()
						}
						if m.Histogram != nil {
							return float64(m.Histogram.GetSampleCount())
						}
						return m.Gauge.GetValue()
					}
				}
			}
		}
	}
	t.Fatalf("missing %s stage=%s", name, stage)
	return 0
}

// The barrier is in the test chain provider, not a production fault hook. The
// real lockedReset holds pool.mu while its actual StateAt call is blocked.
func TestPoolObservabilityResetBlocksIndependentAccounts(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	// Stop background consumers before replacing the chain provider. The real
	// admission/reset/Pending methods below still execute on this pool.
	pool.Stop()
	txs := make([]*types.Transaction, 5)
	for i := range txs {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
		txs[i] = pricedTransaction(0, 21000, new(big.Int).Set(pool.gasPrice), key)
	}
	chain := &poolObservationChain{blockChain: pool.chain, entered: make(chan struct{}), release: make(chan struct{})}
	pool.chain = chain
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	var wg sync.WaitGroup
	var once sync.Once
	release := func() { once.Do(func() { close(chain.release) }) }
	defer func() {
		release()
		drained := make(chan struct{})
		go func() { wg.Wait(); close(drained) }()
		select {
		case <-drained:
		case <-time.After(3 * time.Second):
			t.Error("pool operations did not drain during cleanup")
		}
	}()
	wg.Add(1)
	go func() { defer wg.Done(); pool.lockedReset(nil, chain.CurrentBlock().Header()) }()
	select {
	case <-chain.entered:
	case <-time.After(time.Second):
		t.Fatal("reset did not reach StateAt")
	}
	results := make(chan error, 5)
	launch := func(f func() error) { wg.Add(1); go func() { defer wg.Done(); results <- f() }() }
	launch(func() error { return pool.AddLocal(txs[0]) })
	launch(func() error { return pool.AddLocal(txs[1]) })
	launch(func() error { return pool.AddRemote(txs[2]) })
	launch(func() error {
		for _, err := range pool.AddLocals(txs[3:]) {
			if err != nil {
				return err
			}
		}
		return nil
	})
	launch(func() error { _, err := pool.Pending(); return err })
	expected := map[string]float64{"admission_write_single": 3, "admission_write_batch": 1, "pending_wait": 1, "reset_body": 1, "reset_hold": 1}
	// Gathering must finish while the pool mutex is held. Otherwise metrics would
	// hide precisely the stall they are intended to diagnose.
	observed := make(chan bool, 1)
	go func() {
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			fs, err := registry.Gather()
			if err != nil {
				observed <- false
				return
			}
			got := map[string]float64{}
			for _, f := range fs {
				if f.GetName() == "kaia_poolobs_operations_inflight" {
					for _, m := range f.Metric {
						got[m.Label[0].GetValue()] = m.GetGauge().GetValue()
					}
				}
			}
			match := true
			for k, v := range expected {
				if got[k] != v {
					match = false
				}
			}
			if match {
				observed <- true
				return
			}
			time.Sleep(time.Millisecond)
		}
		observed <- false
	}()
	select {
	case ok := <-observed:
		if !ok {
			t.Fatal("blocked stage counts not observed")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("metrics blocked on pool lock")
	}
	if got := poolObservationValue(t, registry, "kaia_poolobs_operations_started_total", "admission_write_single"); got != 3 {
		t.Fatalf("single admission count: %g", got)
	}
	if got := poolObservationValue(t, registry, "kaia_poolobs_operations_started_total", "admission_write_batch"); got != 1 {
		t.Fatalf("batch admission count: %g", got)
	}
	select {
	case err := <-results:
		t.Fatalf("operation completed before release: %v", err)
	default:
	}
	// While the barrier is held, these stage counts are stable. A general
	// Prometheus scrape is not an atomic snapshot of changing stage counters.
	const hold = 30 * time.Millisecond
	timer := time.NewTimer(hold)
	<-timer.C
	release()
	for i := 0; i < 5; i++ {
		select {
		case err := <-results:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("operation failed to recover")
		}
	}
	wg.Wait()
	fs, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if f.GetName() == "kaia_poolobs_duration_seconds" {
			for _, m := range f.Metric {
				if _, ok := expected[m.Label[0].GetValue()]; ok && m.Histogram.GetSampleSum() < hold.Seconds() {
					t.Fatalf("%s duration omitted barrier: %g", m.Label[0].GetValue(), m.Histogram.GetSampleSum())
				}
			}
		}
	}
	for stage := range expected {
		if got := poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", stage); got != 0 {
			t.Fatalf("%s remains active: %g", stage, got)
		}
	}
	pending, err := pool.Pending()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, list := range pending {
		count += len(list)
	}
	if count != len(txs) {
		t.Fatalf("pending after recovery=%d want %d", count, len(txs))
	}
	t.Log("two independent AddLocal accounts and AddRemote blocked at shared write lock; batch and Pending blocked; gather independent; all five transactions pending after release")
}

// The no-barrier control confirms the same independent-account input is ready
// without the injected reset hold; a low nonce or fee is not the reason it waits.
func TestPoolObservabilityIndependentAccountsControl(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	pool, _ := setupTxPool()
	defer pool.Stop()
	for i := 0; i < 2; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		pool.mu.Lock()
		pool.currentState.SetBalance(crypto.PubkeyToAddress(key.PublicKey), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)))
		price := new(big.Int).Set(pool.gasPrice)
		pool.mu.Unlock()
		if err := pool.AddLocal(pricedTransaction(0, 21000, price, key)); err != nil {
			t.Fatal(err)
		}
	}
	pending, err := pool.Pending()
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 2 {
		t.Fatalf("ready independent accounts=%d", len(pending))
	}
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	for _, stage := range []string{"admission_write_single", "pending_wait"} {
		if poolObservationValue(t, registry, "kaia_poolobs_operations_inflight", stage) != 0 {
			t.Fatalf("control left %s active", stage)
		}
	}
}
