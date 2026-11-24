// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func completeCandidateObservations() (*CandidateObservation, *CandidateObservation) {
	prep, exec := NewCandidatePreparation(11), NewCandidateObservation()
	prep.SetPhase(CandidatePoolInput, 10)
	prep.SetPhase(CandidateAfterFee, 8)
	prep.SetPhase(CandidateAfterModules, 7)
	exec.SetPhase(CandidateExecutionInput, 7)
	exec.SetPhase(CandidateIncluded, 0)
	return prep, exec
}

func TestAdoptedCandidatePopulation(t *testing.T) {
	defer EnableForTesting()()
	r := current.Load()
	prep, exec := completeCandidateObservations()
	exec.CountOutcome(CandidateNonceHigh)
	exec.CountOutcome(CandidateNonceHigh)
	exec.CountOutcome(candidateOutcomeCount)
	prep.CountOutcome(CandidateSizeLimit) // Preparation outcomes must not leak into execution.
	ObserveAdopted(prep, exec.Snapshot(), 180*time.Millisecond, true)
	if got := metricValue(r.adoptedBlocks[0]); got != 1 {
		t.Fatalf("observed=%v", got)
	}
	for phase, want := range []float64{10, 8, 7, 7, 0} {
		if got := metricValue(r.adoptedCandidates[phase].transactions); got != want {
			t.Fatalf("phase %d=%v want%v", phase, got, want)
		}
		wantEmpty := float64(0)
		if want == 0 {
			wantEmpty = 1
		}
		if got := metricValue(r.adoptedCandidates[phase].empty); got != wantEmpty {
			t.Fatalf("empty phase %d=%v", phase, got)
		}
	}
	if metricValue(r.adoptedOutcomes[CandidateNonceHigh]) != 2 || metricValue(r.adoptedOutcomes[CandidateSizeLimit]) != 0 {
		t.Fatal("wrong branch population")
	}
	if got := metricValue(r.adoptedQueued); got != 11 {
		t.Fatalf("queued snapshot = %v", got)
	}
	m := &dto.Metric{}
	if err := r.adoptedExecute.Write(m); err != nil {
		t.Fatal(err)
	}
	if m.Histogram.GetSampleCount() != 1 || m.Histogram.GetSampleSum() != .180 {
		t.Fatal("wrong duration population")
	}
	// Existing all-candidate metrics must not be incremented by publication.
	if metricValue(r.candidates[CandidatePoolInput].transactions) != 0 {
		t.Fatal("adopted publication changed all-candidate metrics")
	}
}

func TestAdoptedMissingObservation(t *testing.T) {
	for _, name := range []string{"nil_prep", "nil_exec", "missing_queued", "missing_prep", "missing_exec", "mismatch", "negative_duration", "recorder_replaced"} {
		t.Run(name, func(t *testing.T) {
			defer EnableForTesting()()
			prep, exec := completeCandidateObservations()
			matched, duration := true, time.Second
			switch name {
			case "nil_prep":
				prep = nil
			case "missing_queued":
				prep.hasQueued = false
			case "nil_exec":
				exec = nil
			case "missing_prep":
				prep.seen[CandidateAfterFee] = false
			case "missing_exec":
				exec.seen[CandidateIncluded] = false
			case "mismatch":
				matched = false
			case "negative_duration":
				duration = -1
			case "recorder_replaced":
				current.Store(newRecorder())
			}
			r := current.Load()
			ObserveAdopted(prep, exec, duration, matched)
			if metricValue(r.adoptedBlocks[1]) != 1 || metricValue(r.adoptedBlocks[0]) != 0 {
				t.Fatal("unavailable observation counted as complete")
			}
			if metricValue(r.adoptedQueued) != 0 || metricValue(r.adoptedCandidates[CandidateIncluded].empty) != 0 {
				t.Fatal("missing observation counted as empty candidate")
			}
		})
	}
}

func TestCandidateSnapshotIndependent(t *testing.T) {
	defer EnableForTesting()()
	_, exec := completeCandidateObservations()
	exec.CountOutcome(CandidateNonceHigh)
	snapshot := exec.Snapshot()
	exec.CountOutcome(CandidateNonceHigh)
	exec.SetPhase(CandidateIncluded, 4)
	if snapshot.outcomes[CandidateNonceHigh] != 1 || snapshot.phases[CandidateIncluded] != 0 {
		t.Fatal("snapshot changed with source")
	}
	snapshot.SetPhase(candidatePhaseCount, 5)
	snapshot.SetPhase(CandidateIncluded, -1)
	if snapshot.phases[CandidateIncluded] != 0 {
		t.Fatal("invalid phase changed observation")
	}
}

func TestAdoptedDisabled(t *testing.T) {
	previous := current.Swap(nil)
	defer current.Store(previous)
	if allocs := testing.AllocsPerRun(100, func() {
		o := NewCandidateObservation()
		if o != nil || NewCandidatePreparation(12) != nil {
			t.Fatal("disabled allocated observation")
		}
		o.SetPhase(CandidateIncluded, 0)
		o.CountOutcome(CandidateNonceHigh)
		ObserveAdopted(nil, o.Snapshot(), time.Second, true)
		CountNotAdopted()
	}); allocs != 0 {
		t.Fatalf("disabled allocations=%g", allocs)
	}
	reg := prometheus.NewRegistry()
	Register(reg)
	fs, err := reg.Gather()
	if err != nil || len(fs) != 0 {
		t.Fatal("disabled exported observations", err)
	}
}

func TestNotAdoptedPopulation(t *testing.T) {
	defer EnableForTesting()()
	CountNotAdopted()
	r := current.Load()
	if metricValue(r.adoptedBlocks[2]) != 1 || metricValue(r.adoptedBlocks[0]) != 0 || metricValue(r.adoptedBlocks[1]) != 0 {
		t.Fatal("wrong result population")
	}
	if metricValue(r.adoptedCandidates[CandidateIncluded].empty) != 0 {
		t.Fatal("nonadopted result entered detail population")
	}
}
