// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"sync"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

func TestPeerQueueLifecycleAndSubscrapeCompletion(t *testing.T) {
	defer EnableForTesting()()
	r := current.Load()
	q := NewPeerQueue()
	other := NewPeerQueue()
	first := q.Begin()
	drop := q.Begin()
	other.Begin()
	if metricValue(r.peerQueued) != 3 {
		t.Fatal("missing provisional observations")
	}
	q.Cancel(drop)
	q.Dequeue(first)
	var h dto.Metric
	if err := r.peerDelay.Write(&h); err != nil {
		t.Fatal(err)
	}
	if h.GetHistogram().GetSampleCount() != 1 || metricValue(r.peerQueued) != 1 {
		t.Fatal("completed event not retained independently of current queue")
	}
	q.Begin()
	q.Close()
	q.Close()
	if metricValue(r.peerQueued) != 1 || metricValue(r.peerAbandoned) != 1 {
		t.Fatal("peer close affected other peer or double counted")
	}
	if !q.Begin().IsZero() {
		t.Fatal("closed queue observed new work")
	}
	q.Dequeue(first)
	other.Close()
	if metricValue(r.peerQueued) != 0 || metricValue(r.peerAbandoned) != 2 {
		t.Fatal("remaining observation leak")
	}
}

func TestPeerQueueConcurrentClose(t *testing.T) {
	defer EnableForTesting()()
	q := NewPeerQueue()
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); at := q.Begin(); q.Dequeue(at) }()
	}
	q.Close()
	wg.Wait()
	if metricValue(current.Load().peerQueued) != 0 {
		t.Fatal("queue observations leaked")
	}
}

func TestPeerQueueDisabled(t *testing.T) {
	old := current.Swap(nil)
	defer current.Store(old)
	q := NewPeerQueue()
	if q != nil {
		t.Fatal("disabled allocation")
	}
	at := q.Begin()
	q.Dequeue(at)
	q.Cancel(at)
	q.Close()
}

func BenchmarkPeerQueue(b *testing.B) {
	defer EnableForTesting()()
	q := NewPeerQueue()
	defer q.Close()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at := q.Begin()
		q.Dequeue(at)
	}
}
