// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package metricutils

import (
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof" // registers /debug/pprof/* on http.DefaultServeMux, like the node binaries
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPrometheusExporterDoesNotExposePprof verifies that newPrometheusMux
// serves /metrics and does not expose the pprof routes registered on
// http.DefaultServeMux.
func TestPrometheusExporterDoesNotExposePprof(t *testing.T) {
	server := httptest.NewServer(newPrometheusMux())
	t.Cleanup(server.Close)

	cases := []struct {
		path string
		want int
	}{
		{"/metrics", http.StatusOK},
		{"/debug/pprof/", http.StatusNotFound},
		{"/debug/pprof/cmdline", http.StatusNotFound},
		{"/debug/pprof/goroutine", http.StatusNotFound},
	}
	for _, tc := range cases {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+tc.path, nil)
		require.NoError(t, err, tc.path)
		resp, err := server.Client().Do(req)
		require.NoError(t, err, tc.path)
		resp.Body.Close()
		assert.Equal(t, tc.want, resp.StatusCode, "unexpected status for %s", tc.path)
	}
}
