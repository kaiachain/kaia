// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PeerQueue tracks only aggregate observations for an existing bounded peer
// queue. No peer or transaction identity is retained. Its lock never covers a
// channel operation or network write.
type PeerQueue struct {
	mu       sync.Mutex
	recorder *recorder
	pending  int
	closed   bool
}

func (r *recorder) registerPeerQueueMetrics() {
	r.peerQueued = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_peer_queue_batches", Help: "Observed async transaction batches awaiting dequeue, including provisional nonblocking enqueue attempts. Not transactions."})
	r.peerAbandoned = prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_peer_queue_abandoned_batches_total", Help: "Pending async batch observations detached when their peer closes. Not proof of transaction loss."})
	r.peerDelay = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "kaia_poolobs_peer_queue_delay_seconds", Help: "Monotonic time from async enqueue attempt to dequeue. Excludes actual peer write; closed queue observations are abandoned instead.", Buckets: durationBuckets})
	r.registry.MustRegister(r.peerQueued, r.peerAbandoned, r.peerDelay)
}

// NewPeerQueue is nil when observation is disabled; all methods accept nil.
func NewPeerQueue() *PeerQueue {
	if r := current.Load(); r != nil {
		return &PeerQueue{recorder: r}
	}
	return nil
}

// Begin runs just before a nonblocking enqueue attempt. The caller must use
// Cancel on rejection, or pass the timestamp with the batch to Dequeue. The
// provisional increment precedes publication so a fast consumer cannot produce
// a negative count. Each timestamp must be completed at most once.
func (q *PeerQueue) Begin() time.Time {
	if q == nil {
		return time.Time{}
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return time.Time{}
	}
	at := time.Now()
	q.pending++
	q.recorder.peerQueued.Inc()
	return at
}
func (q *PeerQueue) Cancel(at time.Time)  { q.finish(at, false) }
func (q *PeerQueue) Dequeue(at time.Time) { q.finish(at, true) }
func (q *PeerQueue) finish(at time.Time, observe bool) {
	if q == nil || at.IsZero() {
		return
	}
	q.mu.Lock()
	if q.closed || q.pending == 0 {
		q.mu.Unlock()
		return
	}
	q.pending--
	q.recorder.peerQueued.Dec()
	q.mu.Unlock()
	if observe {
		q.recorder.peerDelay.Observe(time.Since(at).Seconds())
	}
}

// Close detaches pending observations without draining or changing the channel.
// Enqueues racing after Close preserve existing peer behavior but are untracked.
func (q *PeerQueue) Close() {
	if q == nil {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	q.recorder.peerQueued.Sub(float64(q.pending))
	q.recorder.peerAbandoned.Add(float64(q.pending))
	q.pending = 0
}
