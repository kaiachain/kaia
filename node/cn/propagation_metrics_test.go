package cn

import (
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/prometheus/client_golang/prometheus"
)

func propagationMetric(t *testing.T, r *prometheus.Registry, name, value string) float64 {
	t.Helper()
	fs, err := r.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if f.GetName() == name {
			for _, m := range f.Metric {
				if value == "" && len(m.Label) == 0 {
					if m.Counter != nil {
						return m.Counter.GetValue()
					}
					if m.Histogram != nil {
						return float64(m.Histogram.GetSampleCount())
					}
					return m.Gauge.GetValue()
				}
				for _, l := range m.Label {
					if l.GetValue() == value {
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
	t.Fatalf("missing %s %s", name, value)
	return 0
}

func awaitPropagationMetric(t *testing.T, r *prometheus.Registry, name, value string, expected float64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for propagationMetric(t, r, name, value) != expected {
		if time.Now().After(deadline) {
			t.Fatalf("%s %s did not reach %g", name, value, expected)
		}
		time.Sleep(time.Millisecond)
	}
}

// MsgPipe blocks the actual encoded write until its payload is consumed. This
// exercises the CN async consumer rather than timing a mocked helper call.
func TestPoolObservabilityAsyncPeerWrite(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	peer, local, remote := newBasePeer()
	done := make(chan struct{})
	go func() { peer.Broadcast(); close(done) }()
	defer func() {
		peer.Close()
		local.Close()
		remote.Close()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Error("peer consumer did not stop")
		}
	}()
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
	peer.AsyncSendTransactions(types.Transactions{tx})
	awaitPropagationMetric(t, r, "kaia_poolobs_operations_inflight", "peer_send", 1)
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_transactions"); n != 1 {
		t.Fatalf("attempted copies %g", n)
	}
	if n := propagationMetric(t, r, "kaia_poolobs_duration_seconds", "peer_send"); n != 0 {
		t.Fatalf("blocked write already completed %g", n)
	}
	msg, err := remote.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}
	if err = msg.Discard(); err != nil {
		t.Fatal(err)
	}
	awaitPropagationMetric(t, r, "kaia_poolobs_operations_inflight", "peer_send", 0)
	// Closed transport exercises the real resend error path.
	remote.Close()
	if err := peer.ReSendTransactions(types.Transactions{tx, tx}); err == nil {
		t.Fatal("closed pipe write succeeded")
	}
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_errors"); n != 1 {
		t.Fatalf("failed calls %g", n)
	}
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_transactions"); n != 3 {
		t.Fatalf("attempted copies %g", n)
	}
}

func TestPoolObservabilityAsyncQueueDrop(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	peer, local, remote := newBasePeer()
	defer local.Close()
	defer remote.Close()
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
	for i := 0; i < maxQueuedTxs; i++ {
		peer.AsyncSendTransactions(types.Transactions{tx})
	}
	peer.AsyncSendTransactions(types.Transactions{tx, tx})
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_dropped"); n != 2 {
		t.Fatalf("dropped copies %g", n)
	}
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_transactions"); n != 0 {
		t.Fatalf("queued copies reported as writes %g", n)
	}
}

func TestPoolObservabilityFailedBatchExcludesUnsentTail(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	writeErr := errors.New("first batch write failed")
	calls := 0
	err := sendTransactionBatches(make(types.Transactions, maxTxMsgItems+1), func(batch types.Transactions) error {
		calls++
		if len(batch) != maxTxMsgItems {
			t.Fatalf("first batch size %d", len(batch))
		}
		return writeErr
	})
	if !errors.Is(err, writeErr) || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_transactions"); n != float64(maxTxMsgItems) {
		t.Fatalf("attempted copies %g, want %d; unsent tail must be excluded", n, maxTxMsgItems)
	}
	if n := propagationMetric(t, r, "kaia_poolobs_events_total", "peer_send_errors"); n != 1 {
		t.Fatalf("failed calls %g", n)
	}
}

func TestPoolObservabilityAsyncQueueWait(t *testing.T) {
	for _, multi := range []bool{false, true} {
		name := "single"
		if multi {
			name = "multi"
		}
		t.Run(name, func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			r := prometheus.NewRegistry()
			poolobserve.Register(r)
			var peer Peer
			local, remote := p2p.MsgPipe()
			if multi {
				var err error
				peer, err = newPeerWithRWs(version, p2pPeers[0], []p2p.MsgReadWriter{local, local})
				if err != nil {
					t.Fatal(err)
				}
			} else {
				peer = newPeer(version, p2pPeers[0], local)
			}
			done := make(chan struct{})
			go func() { peer.Broadcast(); close(done) }()
			defer func() {
				peer.Close()
				local.Close()
				remote.Close()
				select {
				case <-done:
				case <-time.After(3 * time.Second):
					t.Error("consumer did not stop")
				}
			}()
			tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
			peer.AsyncSendTransactions(types.Transactions{tx})
			awaitPropagationMetric(t, r, "kaia_poolobs_operations_inflight", "peer_send", 1)
			peer.AsyncSendTransactions(types.Transactions{tx})
			if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_batches", ""); got != 1 {
				t.Fatalf("blocked write hides queued batch: %v", got)
			}
			if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_delay_seconds", ""); got != 1 {
				t.Fatalf("only dequeued batch should have duration: %v", got)
			}
			// Unblock the first write. The second batch then leaves the queue but remains
			// in peer_send until its encoded payload is consumed.
			msg, err := remote.ReadMsg()
			if err != nil {
				t.Fatal(err)
			}
			if err = msg.Discard(); err != nil {
				t.Fatal(err)
			}
			awaitPropagationMetric(t, r, "kaia_poolobs_peer_queue_batches", "", 0)
			awaitPropagationMetric(t, r, "kaia_poolobs_peer_queue_delay_seconds", "", 2)
			msg, err = remote.ReadMsg()
			if err != nil {
				t.Fatal(err)
			}
			if err = msg.Discard(); err != nil {
				t.Fatal(err)
			}
			awaitPropagationMetric(t, r, "kaia_poolobs_operations_inflight", "peer_send", 0)
		})
	}
}

func TestPoolObservabilityAsyncQueueClose(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	peer, local, remote := newBasePeer()
	defer local.Close()
	defer remote.Close()
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
	for i := 0; i < maxQueuedTxs; i++ {
		peer.AsyncSendTransactions(types.Transactions{tx})
	}
	peer.AsyncSendTransactions(types.Transactions{tx})
	if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_batches", ""); got != maxQueuedTxs {
		t.Fatalf("drop must roll back provisional count: %v", got)
	}
	peer.Close()
	if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_batches", ""); got != 0 {
		t.Fatalf("close leaked queue count: %v", got)
	}
	if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_abandoned_batches_total", ""); got != maxQueuedTxs {
		t.Fatalf("detached observations: %v", got)
	}
	if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_delay_seconds", ""); got != 0 {
		t.Fatalf("close must not masquerade as dequeue: %v", got)
	}
}

func TestPoolObservabilityAsyncQueueCloseRace(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	r := prometheus.NewRegistry()
	poolobserve.Register(r)
	peer, local, remote := newBasePeer()
	defer local.Close()
	defer remote.Close()
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(25_000_000_000), nil)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 32; j++ {
				peer.AsyncSendTransactions(types.Transactions{tx})
			}
		}()
	}
	close(start)
	peer.Close()
	wg.Wait()
	if got := propagationMetric(t, r, "kaia_poolobs_peer_queue_batches", ""); got != 0 {
		t.Fatalf("close race leaked queue count: %v", got)
	}
}
