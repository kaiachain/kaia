// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestFlowCountersAreBoundedAndConcurrent(t *testing.T) {
	defer EnableForTesting()()
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Count(Admitted, 3)
			ObserveCandidate(CandidatePoolInput, 3)
			ObserveCandidate(CandidateIncluded, 0)
			CountCandidateOutcome(CandidateSizeLimit)
		}()
	}
	wg.Wait()
	r := current.Load()
	if metricValue(r.outcomes[CandidateSizeLimit]) != 40 {
		t.Fatal("lost candidate decisions")
	}
	CountCandidateOutcome(candidateOutcomeCount)
	if metricValue(r.events[Admitted]) != 120 || metricValue(r.candidates[CandidatePoolInput].transactions) != 120 || metricValue(r.candidates[CandidateIncluded].empty) != 40 {
		t.Fatal("incorrect concurrent counts")
	}
	Count(Admitted, -1)
	Count(eventCount, 5)
	ObserveCandidate(candidatePhaseCount, 4)
	ObserveCandidate(CandidateIncluded, -1)
	if metricValue(r.events[Admitted]) != 120 || metricValue(r.preparations) != 40 {
		t.Fatal("invalid input changed counters")
	}
}

func TestHeadStagesRemainDistinct(t *testing.T) {
	defer EnableForTesting()()
	var head atomic.Uint64
	head.Store(12)
	stop := TrackChainHead(head.Load)
	SetPoolStateHead(11)
	SetPoolResetHead(10)
	r := current.Load()
	if metricValue(r.poolState) != 11 || metricValue(r.poolReset) != 10 {
		t.Fatal("lost head boundary")
	}
	if metricValue(r.poolStateObserved) != 1 || metricValue(r.poolResetObserved) != 1 {
		t.Fatal("missing observed state")
	}
	second := TrackChainHead(func() uint64 { return 13 })
	stop()
	if r.chainHead.Load().number() != 13 {
		t.Fatal("old cleanup detached new registration")
	}
	second()
	if r.chainHead.Load() != nil {
		t.Fatal("head source retained after stop")
	}
	// A restart/genesis observation must be distinguishable from never observed.
	SetPoolStateHead(0)
	SetPoolResetHead(0)
	if metricValue(r.poolStateObserved) != 1 {
		t.Fatal("genesis treated as unavailable")
	}
}

func TestFlowDisabledDoesNotReadSource(t *testing.T) {
	old := current.Swap(nil)
	defer current.Store(old)
	TrackChainHead(func() uint64 { t.Fatal("disabled head callback called"); return 0 })()
	Count(Admitted, 1)
	CountCandidateOutcome(CandidateSizeLimit)
	ObserveCandidate(CandidateIncluded, 1)
	SetPoolStateHead(1)
	SetPoolResetHead(1)
	reg := prometheus.NewRegistry()
	Register(reg)
	fs, err := reg.Gather()
	if err != nil || len(fs) != 0 {
		t.Fatal("disabled metrics exported")
	}
}

func BenchmarkFlowCounters(b *testing.B) {
	defer EnableForTesting()()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Count(Admitted, 1)
		ObserveCandidate(CandidateIncluded, 1000)
	}
}

func metricValue(m prometheus.Metric) float64 {
	v := &dto.Metric{}
	if e := m.Write(v); e != nil {
		panic(e)
	}
	if v.Counter != nil {
		return v.Counter.GetValue()
	}
	return v.Gauge.GetValue()
}

func BenchmarkObservationConcurrentScrape(b *testing.B) {
	defer EnableForTesting()()
	r := current.Load()
	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		for {
			select {
			case <-done:
				return
			default:
				_, _ = r.registry.Gather()
			}
		}
	}()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Begin(PeerSend)()
		Count(PeerSendTransactions, 1)
	}
	b.StopTimer()
	close(done)
	<-stopped
}
