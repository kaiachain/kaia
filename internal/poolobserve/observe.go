// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Stage identifies a fixed measurement boundary. No transaction or account
// identity is retained, so recorder state and exported cardinality are fixed.
// Each enabled in-flight interval still has its own completion closure.
type Stage uint8

const (
	AdmissionWriteSingle Stage = iota
	AdmissionWriteBatch
	HeadResetWait
	PendingWait // Only pool.mu acquisition, not the subsequent pool.txMu acquisition.
	ResetBody
	ResetHold
	FeedEnqueue // One event-batch submission, not one transaction.
	FeedSend    // One Feed.Send invocation to all subscribers.
	AdmissionHold
	RemoteEnqueue
	RemoteAdmission
	Broadcast
	PeerSend
	Rebroadcast
	BlockImport
	CandidateBuild
	CandidatePending
	stageCount
)

var names = [stageCount]string{"admission_write_single", "admission_write_batch", "head_reset_wait", "pending_wait", "reset_body", "reset_hold", "feed_enqueue", "feed_send", "admission_hold", "remote_enqueue", "remote_admission", "broadcast", "peer_send", "rebroadcast", "block_import", "candidate_build", "candidate_pending"}

// Queue identifies an event channel, whose slots can contain multiple txs.
type Queue uint8

const (
	FeedQueue Queue = iota
	BroadcastQueue
	RemoteQueue
	queueCount
)

var queueNames = [queueCount]string{"feed", "broadcast", "remote"}

type queueSource struct {
	depth    func() int
	capacity int
}

type stageMetrics struct {
	started  prometheus.Counter
	inflight prometheus.Gauge
	duration prometheus.Histogram
}

// Resolve 100-250ms waits, the one-second block interval, and longer propagation stalls.
// Completed operations already have an exact count in the histogram.
var durationBuckets = []float64{.1, .15, .25, 1, 5, 10, 20}

type recorder struct {
	peerQueued                           prometheus.Gauge
	peerAbandoned                        prometheus.Counter
	peerDelay                            prometheus.Histogram
	registry                             *prometheus.Registry
	stages                               [stageCount]stageMetrics
	queues                               [queueCount]atomic.Pointer[queueSource]
	events                               [eventCount]prometheus.Counter
	candidates                           [candidatePhaseCount]candidateMetrics
	preparations                         prometheus.Counter
	outcomes                             [candidateOutcomeCount]prometheus.Counter
	chainHead                            atomic.Pointer[headSource]
	poolState, poolReset                 prometheus.Gauge
	poolStateObserved, poolResetObserved prometheus.Gauge
}

var current atomic.Pointer[recorder]

func init() {
	if os.Getenv("KAIA_POOLOBSERVE") == "1" {
		current.Store(newRecorder())
	}
}

func newRecorder() *recorder {
	r := &recorder{registry: prometheus.NewRegistry()}
	schema := prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_schema_version", Help: "Pool observation schema version."})
	schema.Set(5)
	r.registry.MustRegister(schema)
	start := prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_process_start_time_seconds", Help: "Wall clock when the process observation recorder initialized; not a duration clock."})
	start.Set(float64(time.Now().UnixNano()) / 1e9)
	r.registry.MustRegister(start)
	for i, name := range queueNames {
		source := &r.queues[i]
		labels := prometheus.Labels{"queue": name}
		r.registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "kaia_poolobs_queue_observed", Help: "One while the local queue consumer is registered; zero means not observed.", ConstLabels: labels}, func() float64 {
			if source.Load() != nil {
				return 1
			}
			return 0
		}))
		r.registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "kaia_poolobs_queue_slots", Help: "Current event-batch slots in the local channel, not transaction count; consult queue_observed.", ConstLabels: labels}, func() float64 {
			if q := source.Load(); q != nil {
				return float64(q.depth())
			}
			return 0
		}))
		r.registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "kaia_poolobs_queue_capacity", Help: "Capacity in event-batch slots of the observed local channel.", ConstLabels: labels}, func() float64 {
			if q := source.Load(); q != nil {
				return float64(q.capacity)
			}
			return 0
		}))
	}
	for i, name := range names {
		labels := prometheus.Labels{"stage": name}
		s := &r.stages[i]
		s.started = prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_operations_started_total", Help: "Local transaction processing stage operations started.", ConstLabels: labels})
		s.inflight = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_operations_inflight", Help: "Local transaction processing stage operations not yet completed.", ConstLabels: labels})
		s.duration = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "kaia_poolobs_duration_seconds", Help: "Completed local transaction processing stage monotonic duration, including unsuccessful operations.", ConstLabels: labels, Buckets: durationBuckets})
		r.registry.MustRegister(s.started, s.inflight, s.duration)
	}
	r.registerFlowMetrics()
	r.registerPeerQueueMetrics()
	return r
}

// Register adds the native cumulative collector to the existing exporter when
// enabled at process start. Disabled observation exports no poolobs series.
func Register(reg prometheus.Registerer) {
	if r := current.Load(); r != nil {
		reg.MustRegister(r.registry)
	}
}
func noop() {}

// Begin starts a paired monotonic interval. The returned completion is idempotent.
// Wait intervals end immediately after acquiring the lock; hold intervals end
// immediately before releasing it. These intervals overlap and must not be summed.
func Begin(stage Stage) func() {
	r := current.Load()
	if r == nil || stage >= stageCount {
		return noop
	}
	s := &r.stages[stage]
	start := time.Now()
	s.started.Inc()
	s.inflight.Inc()
	var ended atomic.Bool
	return func() {
		if !ended.CompareAndSwap(false, true) {
			return
		}
		s.duration.Observe(time.Since(start).Seconds())
		s.inflight.Dec()
	}
}

// EnableForTesting installs an isolated recorder and restores the previous one.
// Tests using this process-wide recorder must not run in parallel. In-flight
// completions retain their original recorder and cannot affect its replacement.
func EnableForTesting() func() {
	previous := current.Swap(newRecorder())
	return func() { current.Store(previous) }
}

// TrackQueue binds one of the fixed channel sources for its consumer's
// lifetime. depth must be a nonblocking channel len read, never a pool snapshot
// that acquires pool.mu. Cleanup cannot detach a newer consumer's registration.
func TrackQueue(queue Queue, depth func() int, capacity int) func() {
	r := current.Load()
	if r == nil || queue >= queueCount {
		return noop
	}
	q := &queueSource{depth: depth, capacity: capacity}
	r.queues[queue].Store(q)
	return func() { r.queues[queue].CompareAndSwap(q, nil) }
}
