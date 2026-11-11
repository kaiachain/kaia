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

package rpc

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

func TestDatadogHTTPHandlerMetadata(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	for _, tc := range []struct {
		name     string
		body     string
		response string
		method   string
		spans    int
	}{
		{
			name:     "single",
			body:     `{"jsonrpc":"2.0","id":1,"method":"personal_importRawKey","params":["test-private-key","test-passphrase"]}`,
			response: `{"jsonrpc":"2.0","id":1,"result":true}`,
			method:   "personal_importRawKey",
			spans:    1,
		},
		{
			name:     "batch",
			body:     `[{"jsonrpc":"2.0","id":1,"method":"personal_importRawKey","params":["test-private-key","test-passphrase"]},{"jsonrpc":"2.0","id":2,"method":"personal_unlockAccount","params":["test-account","test-passphrase"]}]`,
			response: `[{"jsonrpc":"2.0","id":1,"result":true},{"jsonrpc":"2.0","id":2,"result":true}]`,
			method:   "personal_importRawKey_batch",
			spans:    3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mt.Reset()
			called := false
			handler := newDatadogHTTPHandler(&DatadogTracer{Service: "rpc-test"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, tc.body, string(body))
				_, err = io.WriteString(w, tc.response)
				require.NoError(t, err)
			}))
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "http://localhost/", strings.NewReader(tc.body))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
			require.True(t, called)
			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, tc.response, recorder.Body.String())

			spans := mt.FinishedSpans()
			require.Len(t, spans, tc.spans)
			var requestSpan mocktracer.Span
			for _, span := range spans {
				tags := span.Tags()
				require.NotContains(t, tags, "request.params")
				require.NotContains(t, fmt.Sprint(tags), "test-private-key")
				require.NotContains(t, fmt.Sprint(tags), "test-passphrase")
				if span.OperationName() == "http.request" {
					requestSpan = span
				}
			}
			require.NotNil(t, requestSpan)
			require.Equal(t, tc.method, requestSpan.Tag("request.method"))
			require.Equal(t, "200", requestSpan.Tag("http.status_code"))
			require.False(t, requestSpan.StartTime().IsZero())
			require.False(t, requestSpan.FinishTime().Before(requestSpan.StartTime()))
		})
	}
}
