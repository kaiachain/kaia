// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/common"
)

type countingBody struct {
	read int64
}

func (b *countingBody) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}
	b.read += int64(len(p))
	return len(p), nil
}

func (b *countingBody) Close() error { return nil }

func TestHTTPErrorResponseWithDelete(t *testing.T) {
	testHTTPErrorResponse(t, http.MethodDelete, contentType, "", http.StatusMethodNotAllowed)
}

func TestHTTPErrorResponseWithPut(t *testing.T) {
	testHTTPErrorResponse(t, http.MethodPut, contentType, "", http.StatusMethodNotAllowed)
}

func TestHTTPErrorResponseWithMaxContentLength(t *testing.T) {
	body := make([]rune, common.MaxRequestContentLength+1)
	testHTTPErrorResponse(t,
		http.MethodPost, contentType, string(body), http.StatusRequestEntityTooLarge)
}

func TestGetRPCRequestsBodyLimit(t *testing.T) {
	t.Run("chunked", func(t *testing.T) {
		body := new(countingBody)
		req := &http.Request{Body: body, ContentLength: -1}
		_, _, err := getRPCRequests(req)
		if !errors.Is(err, errAPMRequestTooLarge) {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(common.MaxRequestContentLength + 1)
		if body.read != want {
			t.Fatalf("read %d bytes, want %d", body.read, want)
		}
		replayed, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		if int64(len(replayed)) != want {
			t.Fatalf("replayed %d bytes, want %d", len(replayed), want)
		}
	})

	t.Run("content length", func(t *testing.T) {
		body := new(countingBody)
		req := &http.Request{Body: body, ContentLength: int64(common.MaxRequestContentLength + 1)}
		_, _, err := getRPCRequests(req)
		if !errors.Is(err, errAPMRequestTooLarge) {
			t.Fatalf("unexpected error: %v", err)
		}
		if body.read != 0 {
			t.Fatalf("read %d bytes", body.read)
		}
	})
}

func TestHTTPErrorResponseWithEmptyContentType(t *testing.T) {
	testHTTPErrorResponse(t, http.MethodPost, "", "", http.StatusUnsupportedMediaType)
}

func TestHTTPErrorResponseWithValidRequest(t *testing.T) {
	testHTTPErrorResponse(t, http.MethodPost, contentType, "", 0)
}

func testHTTPErrorResponse(t *testing.T, method, contentType, body string, expected int) {
	request := httptest.NewRequest(method, "http://url.com", strings.NewReader(body))
	request.Header.Set("content-type", contentType)
	if code, _ := validateRequest(request); code != expected {
		t.Fatalf("response code should be %d not %d", expected, code)
	}
}
