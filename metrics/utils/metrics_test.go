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
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* on http.DefaultServeMux, like the node binaries
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// freePort returns an ephemeral TCP port that is free at call time.
func freePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	require.NoError(t, l.Close())
	return port
}

// httpGet issues a context-aware GET and returns the response.
func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func newExporterCtx(port int) *cli.Context {
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.Int(PrometheusExporterPortFlag, port, "")
	return cli.NewContext(cli.NewApp(), set, nil)
}

// TestPrometheusExporterDoesNotExposePprof checks that the exporter serves
// /metrics but not net/http/pprof's /debug/pprof/*, guarding against a regression
// to a nil (DefaultServeMux) handler that would leak them on the Prometheus port.
func TestPrometheusExporterDoesNotExposePprof(t *testing.T) {
	// Enabled/EnabledPrometheusExport are normally set by init() from os.Args.
	prevEnabled, prevExport := Enabled, EnabledPrometheusExport
	Enabled, EnabledPrometheusExport = true, true
	t.Cleanup(func() { Enabled, EnabledPrometheusExport = prevEnabled, prevExport })

	port := freePort(t)
	StartMetricCollectionAndExport(newExporterCtx(port))

	base := fmt.Sprintf("http://127.0.0.1:%d", port)

	// Wait for the listener to come up.
	require.Eventually(t, func() bool {
		resp, err := httpGet(base + "/metrics")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 20*time.Millisecond, "/metrics never became reachable")

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
		resp, err := httpGet(base + tc.path)
		require.NoError(t, err, tc.path)
		resp.Body.Close()
		assert.Equal(t, tc.want, resp.StatusCode, "unexpected status for %s", tc.path)
	}
}
