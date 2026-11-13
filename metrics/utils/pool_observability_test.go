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
	if r.Code != 200 || !strings.Contains(r.Body.String(), "kaia_poolobs_schema_version 5") || !strings.Contains(r.Body.String(), "stage=\"feed_send\"") {
		t.Fatal("native schema missing from private exporter")
	}
	for _, path := range []string{"/debug/pprof/", "/debug/pprof/goroutine"} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
		if r.Code != 404 {
			t.Fatal("private metrics mux exposed global debug handler")
		}
	}
}
