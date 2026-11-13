// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import "github.com/prometheus/client_golang/prometheus"

// Event has a fixed vocabulary. Counts are local observations, not network-wide
// unique transactions: retries, rebroadcast and per-peer copies can repeat them.
type Event uint8

const (
	Admitted Event = iota
	AdmissionRejected
	Promoted
	BroadcastTransactions
	PeerSendTransactions
	PeerSendErrors
	RemoteReceived
	RebroadcastTransactions
	PeerSendDropped
	AdmissionCapacityRejected
	RemoteAdmitted
	RemoteRejected
	RemoteDenied
	RemoteThrottled
	BroadcastNoRecipients
	RebroadcastNoRecipients
	eventCount
)

var eventNames = [eventCount]string{"admitted", "admission_rejected", "promoted", "broadcast_transactions", "peer_send_transactions", "peer_send_errors", "remote_received", "rebroadcast_transactions", "peer_send_dropped", "capacity_rejected", "remote_admitted", "remote_rejected", "remote_denied", "remote_throttled", "broadcast_no_recipients", "rebroadcast_no_recipients"}

type CandidatePhase uint8

const (
	CandidatePoolInput CandidatePhase = iota
	CandidateAfterFee
	CandidateAfterModules
	CandidateExecutionInput
	CandidateIncluded
	candidatePhaseCount
)

var phaseNames = [candidatePhaseCount]string{"pool_input", "after_fee", "after_modules", "execution_input", "included"}

// CandidateOutcome counts branch decisions, not affected transaction counts.
// One skipped bundle can contain several transactions; repeated candidates can
// encounter the same decision again.
type CandidateOutcome uint8

const (
	CandidateBundleTargetRejected CandidateOutcome = iota
	CandidateGeneratorFailed
	CandidateSizeLimit
	CandidateBlobLimit
	CandidateUnsupportedType
	CandidateBundleReverted
	CandidateNonceLow
	CandidateNonceHigh
	candidateOutcomeCount
)

var outcomeNames = [candidateOutcomeCount]string{"bundle_target_rejected", "generator_failed", "size_limit", "blob_space_rejected", "unsupported_type", "bundle_reverted", "nonce_low", "nonce_high"}

type (
	candidateMetrics struct{ transactions, empty prometheus.Counter }
	headSource       struct{ number func() uint64 }
)

func Enabled() bool { return current.Load() != nil }

func (r *recorder) registerFlowMetrics() {
	for i, name := range outcomeNames {
		c := prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_candidate_outcomes_total", Help: "Local candidate branch decisions, not excluded transaction counts or canonical block outcomes.", ConstLabels: prometheus.Labels{"reason": name}})
		r.outcomes[i] = c
		r.registry.MustRegister(c)
	}

	for i, name := range eventNames {
		c := prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_events_total", Help: "Local event units: transaction copies/attempts except peer_send_errors (failed peer calls) and *_no_recipients (nonempty calls without recipients). Capacity rejection includes internal attempts. Not unique hashes or remote acknowledgements.", ConstLabels: prometheus.Labels{"event": name}})
		r.events[i] = c
		r.registry.MustRegister(c)
	}
	r.preparations = prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_candidate_preparations_total", Help: "Successful pending snapshots starting local candidate preparation. Fee and module phases share this population."})
	r.registry.MustRegister(r.preparations)
	for i, name := range phaseNames {
		labels := prometheus.Labels{"phase": name}
		m := candidateMetrics{
			transactions: prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_candidate_transactions_total", Help: "Sum of concrete transaction counts at local candidate phases; repeated candidates can count transactions again.", ConstLabels: labels}),
			empty:        prometheus.NewCounter(prometheus.CounterOpts{Name: "kaia_poolobs_candidate_empty_total", Help: "Local candidate phase observations with zero concrete transactions; not canonical empty blocks.", ConstLabels: labels}),
		}
		r.candidates[i] = m
		r.registry.MustRegister(m.transactions, m.empty)
	}
	r.registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "kaia_poolobs_chain_head_observed", Help: "One while a canonical head reader is registered."}, func() float64 {
		if r.chainHead.Load() != nil {
			return 1
		}
		return 0
	}))
	r.registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "kaia_poolobs_chain_head_number", Help: "Current local canonical head at scrape time; consult chain_head_observed. Not atomically sampled with pool head."}, func() float64 {
		if h := r.chainHead.Load(); h != nil {
			return float64(h.number())
		}
		return 0
	}))
	r.poolState = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_pool_state_head_number", Help: "Head whose state was most recently installed by pool reset; promotion may still be running."})
	r.poolReset = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_pool_reset_head_number", Help: "Head at the most recent fully completed pool reset; not proof any specific transaction is ready."})
	r.poolStateObserved = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_pool_state_head_observed", Help: "One after pool state head was observed, including genesis."})
	r.poolResetObserved = prometheus.NewGauge(prometheus.GaugeOpts{Name: "kaia_poolobs_pool_reset_head_observed", Help: "One after a complete pool reset was observed, including genesis."})
	r.registry.MustRegister(r.poolState, r.poolReset, r.poolStateObserved, r.poolResetObserved)
}

// Count ignores invalid selectors and nonpositive increments. Disabled mode is
// a single atomic read and retains neither transaction nor peer identities.
func Count(event Event, n int) {
	if r := current.Load(); r != nil && event < eventCount && n > 0 {
		r.events[event].Add(float64(n))
	}
}

func ObserveCandidate(phase CandidatePhase, n int) {
	if r := current.Load(); r != nil && phase < candidatePhaseCount && n >= 0 {
		m := &r.candidates[phase]
		if phase == CandidatePoolInput {
			r.preparations.Inc()
		}
		m.transactions.Add(float64(n))
		if n == 0 {
			m.empty.Inc()
		}
	}
}

// TrackChainHead requires a nonblocking, concurrency-safe canonical-head read.
// The source is detached with the pool consumer, never by a newer registration.
func TrackChainHead(number func() uint64) func() {
	r := current.Load()
	if r == nil || number == nil {
		return noop
	}
	s := &headSource{number: number}
	r.chainHead.Store(s)
	return func() { r.chainHead.CompareAndSwap(s, nil) }
}

func SetPoolStateHead(n uint64) {
	if r := current.Load(); r != nil {
		r.poolState.Set(float64(n))
		r.poolStateObserved.Set(1)
	}
}

func SetPoolResetHead(n uint64) {
	if r := current.Load(); r != nil {
		r.poolReset.Set(float64(n))
		r.poolResetObserved.Set(1)
	}
}

// CountCandidateOutcome records one execution branch decision.
func CountCandidateOutcome(outcome CandidateOutcome) {
	if r := current.Load(); r != nil && outcome < candidateOutcomeCount {
		r.outcomes[outcome].Inc()
	}
}
