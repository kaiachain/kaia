package metricutils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/prometheus/client_golang/prometheus"
)

func TestPoolObservationExporter(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	old := prometheus.DefaultGatherer
	prometheus.DefaultGatherer = registry
	defer func() { prometheus.DefaultGatherer = old }()
	mux := newPrometheusMux()
	r := httptest.NewRecorder()
	mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if r.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d", r.Code, http.StatusOK)
	}
	for _, sample := range []string{
		"kaia_poolobs_schema_version 10",
		`kaia_poolobs_tx_messages_total{result="forwarded"} 0`,
		`kaia_poolobs_promotion_stops_total{reason="below_base_fee"} 0`,
		`kaia_poolobs_adopted_candidate_outcomes_total{reason="time_limit_with_remaining"} 0`,
		`kaia_poolobs_adopted_candidate_outcomes_total{reason="time_limit_without_remaining"} 0`,
		"stage=\"feed_send\"",
		`kaia_poolobs_remote_handoff_seconds_bucket{le="0.1"} 0`,
		"kaia_poolobs_adopted_candidate_queued_transactions_total 0",
		`kaia_poolobs_adopted_candidate_outcomes_total{reason="bundle_incorporation_failed"} 0`,
		`kaia_poolobs_candidate_outcomes_total{reason="other_execution_error"} 0`,
		"kaia_poolobs_adopted_candidate_blocks_total{result=\"observed\"} 0",
		"kaia_poolobs_adopted_candidate_execute_seconds_count 0",
	} {
		if !strings.Contains(r.Body.String(), sample) {
			t.Errorf("private exporter missing %q", sample)
		}
	}
	for _, removed := range []string{"kaia_poolobs_queue_observed", "kaia_poolobs_candidate_preparations_total"} {
		if strings.Contains(r.Body.String(), removed) {
			t.Errorf("removed metric still exported: %s", removed)
		}
	}
	for _, path := range []string{"/debug/pprof/", "/debug/pprof/goroutine"} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
		if r.Code != 404 {
			t.Fatal("private metrics mux exposed global debug handler")
		}
	}
}
