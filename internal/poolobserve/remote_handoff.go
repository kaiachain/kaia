// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (r *recorder) registerRemoteHandoff() {
	r.remoteHandoff = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "kaia_poolobs_remote_handoff_seconds",
		Help:    "Monotonic time from a normal remote batch's channel send attempt to dequeue, including a blocked send and excluding pool admission.",
		Buckets: []float64{0.1, 1},
	})
	r.registry.MustRegister(r.remoteHandoff)
}

// RemoteHandoff travels with one batch on the normal remote intake channel.
type RemoteHandoff struct {
	recorder *recorder
	start    time.Time
}

// StartRemoteHandoff runs immediately before the producer attempts the send.
func StartRemoteHandoff() RemoteHandoff {
	if r := current.Load(); r != nil {
		return RemoteHandoff{recorder: r, start: time.Now()}
	}
	return RemoteHandoff{}
}

// ObserveDequeue runs once when the consumer receives the batch.
func (h RemoteHandoff) ObserveDequeue() {
	if h.recorder != nil {
		h.recorder.remoteHandoff.Observe(time.Since(h.start).Seconds())
	}
}
