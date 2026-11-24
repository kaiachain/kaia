// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// These reasons apply only to a matched, adopted execution. The remaining
// entries can include bundle generators and need not be executable transactions.
const (
	timeLimitWithRemaining = candidateOutcomeCount + iota
	timeLimitWithoutRemaining
	adoptedOutcomeCount
)

// CandidateObservation belongs to one preparation or execution attempt. It is
// not concurrency-safe: the owner finishes recording before handing off a
// snapshot. Fixed arrays retain neither transaction nor block identities.
type CandidateObservation struct {
	queued    uint64
	hasQueued bool
	recorder  *recorder
	phases    [candidatePhaseCount]int
	seen      [candidatePhaseCount]bool
	outcomes  [adoptedOutcomeCount]uint64
}

// NewCandidateObservation returns nil without allocating when disabled.
func NewCandidateObservation() *CandidateObservation {
	if r := current.Load(); r != nil {
		return &CandidateObservation{recorder: r}
	}
	return nil
}

// NewCandidatePreparation attaches the queued population from the same
// snapshot as the pending inputs, including a valid zero.
func NewCandidatePreparation(queued uint64) *CandidateObservation {
	o := NewCandidateObservation()
	if o != nil {
		o.queued, o.hasQueued = queued, true
	}
	return o
}

func (o *CandidateObservation) SetPhase(phase CandidatePhase, n int) {
	if o != nil && phase < candidatePhaseCount && n >= 0 {
		o.phases[phase], o.seen[phase] = n, true
	}
}

func (o *CandidateObservation) CountOutcome(outcome CandidateOutcome) {
	if o != nil && outcome < candidateOutcomeCount {
		o.outcomes[outcome]++
	}
}

// RecordTimeLimit records the state of the worklist when the worker observes its deadline.
func (o *CandidateObservation) RecordTimeLimit(hasRemaining bool) {
	if o == nil {
		return
	}
	if hasRemaining {
		o.outcomes[timeLimitWithRemaining]++
	} else {
		o.outcomes[timeLimitWithoutRemaining]++
	}
}

// Snapshot returns an independent copy for the completed execution result.
func (o *CandidateObservation) Snapshot() *CandidateObservation {
	if o == nil {
		return nil
	}
	copy := *o
	return &copy
}

func (r *recorder) registerAdoptedMetrics() {
	for i, result := range []string{"observed", "unobserved", "not_adopted"} {
		c := prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_adopted_candidate_blocks_total", Help: "Worker-consumed sealed results: observed local-author canonical candidates with complete matching observations, unobserved local-author canonical results without them, or not_adopted results. Unconsumed result channels are not counted.", ConstLabels: prometheus.Labels{"result": result}})
		r.adoptedBlocks[i] = c
		r.registry.MustRegister(c)
	}
	for i, phase := range phaseNames {
		labels := prometheus.Labels{"phase": phase}
		m := candidateMetrics{
			transactions: prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_adopted_candidate_transactions_total", Help: "Sum of transaction counts at each phase of observed local-author canonical candidates consumed by the worker.", ConstLabels: labels}),
			empty:        prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_adopted_candidate_empty_total", Help: "Zero-transaction phases of observed local-author canonical candidates consumed by the worker.", ConstLabels: labels}),
		}
		r.adoptedCandidates[i] = m
		r.registry.MustRegister(m.transactions, m.empty)
	}
	for i, reason := range append(outcomeNames[:], "time_limit_with_remaining", "time_limit_without_remaining") {
		c := prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_adopted_candidate_outcomes_total", Help: "Execution branch decisions of observed local-author canonical candidates consumed by the worker, not excluded transaction counts.", ConstLabels: prometheus.Labels{"reason": reason}})
		r.adoptedOutcomes[i] = c
		r.registry.MustRegister(c)
	}
	r.adoptedExecute = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "kaia_poolobs_adopted_candidate_execute_seconds", Help: "Executor call duration, including its lock wait, of observed local-author canonical candidates consumed by the worker; excludes finalization and sealing.", Buckets: []float64{0.01, 0.1, 0.25, 1}})
	r.adoptedQueued = prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_adopted_candidate_queued_transactions_total", Help: "Sum of queued transaction counts at the pending snapshots of observed local-author canonical candidates consumed by the worker."})
	r.registry.MustRegister(r.adoptedExecute, r.adoptedQueued)
}

// ObserveAdopted is called once for a worker-consumed local-author canonical
// result. The caller establishes canonical membership and attempt matching.
// Missing phases or recorder replacements must not appear as zero-transaction
// phases in the detailed population.
func ObserveAdopted(prep, exec *CandidateObservation, duration time.Duration, matched bool) {
	r := current.Load()
	if r == nil {
		return
	}
	valid := matched && prep != nil && exec != nil && duration >= 0
	if valid {
		valid = prep.hasQueued && prep.recorder == r && exec.recorder == r
		for phase := CandidatePoolInput; phase <= CandidateAfterModules; phase++ {
			valid = valid && prep.seen[phase]
		}
		for phase := CandidateExecutionInput; phase <= CandidateIncluded; phase++ {
			valid = valid && exec.seen[phase]
		}
	}
	if !valid {
		r.adoptedBlocks[1].Inc()
		return
	}
	r.adoptedBlocks[0].Inc()
	r.adoptedQueued.Add(float64(prep.queued))
	for phase := CandidatePoolInput; phase < candidatePhaseCount; phase++ {
		n := prep.phases[phase]
		if phase >= CandidateExecutionInput {
			n = exec.phases[phase]
		}
		m := &r.adoptedCandidates[phase]
		m.transactions.Add(float64(n))
		if n == 0 {
			m.empty.Inc()
		}
	}
	for outcome, n := range exec.outcomes {
		r.adoptedOutcomes[outcome].Add(float64(n))
	}
	r.adoptedExecute.Observe(duration.Seconds())
}

// CountNotAdopted records a consumed sealed result that did not meet the local
// canonical and local-author conditions. Nil/nonproposer results are excluded.
func CountNotAdopted() {
	if r := current.Load(); r != nil {
		r.adoptedBlocks[2].Inc()
	}
}
