// Copyright 2015 The go-ethereum Authors
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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

type Service struct{}

type Args struct {
	S string
}

func (s *Service) NoArgsRets() {
}

type Result struct {
	String string
	Int    int
	Args   *Args
}

func (s *Service) Echo(str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *Service) EchoWithCtx(ctx context.Context, str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *Service) Sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
	case <-ctx.Done():
	}
}

func (s *Service) Rets() (string, error) {
	return "", nil
}

func (s *Service) InvalidRets1() (error, string) {
	return nil, ""
}

func (s *Service) InvalidRets2() (string, string) {
	return "", ""
}

func (s *Service) InvalidRets3() (string, string, error) {
	return "", "", nil
}

func (s *Service) Subscription(ctx context.Context) (*Subscription, error) {
	return nil, nil
}

func TestServerRegisterName(t *testing.T) {
	server := NewServer()
	service := new(Service)

	if err := server.RegisterName("calc", service); err != nil {
		t.Fatalf("%v", err)
	}

	if len(server.services.services) != 2 {
		t.Fatalf("Expected 2 service entries, got %d", len(server.services.services))
	}

	svc, ok := server.services.services["calc"]
	if !ok {
		t.Fatalf("Expected service calc to be registered")
	}

	if len(svc.callbacks) != 5 {
		t.Errorf("Expected 5 callbacks for service 'calc', got %d", len(svc.callbacks))
	}

	if len(svc.subscriptions) != 1 {
		t.Errorf("Expected 1 subscription for service 'calc', got %d", len(svc.subscriptions))
	}
}

func testServerMethodExecution(t *testing.T, method string) {
	server := NewServer()
	service := new(Service)

	if err := server.RegisterName("test", service); err != nil {
		t.Fatalf("%v", err)
	}

	stringArg := "string arg"
	intArg := 1122
	argsArg := &Args{"abcde"}
	params := []interface{}{stringArg, intArg, argsArg}

	request := map[string]interface{}{
		"id":      12345,
		"method":  "test_" + method,
		"version": "2.0",
		"params":  params,
	}

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go server.ServeCodec(NewCodec(serverConn), 0)

	out := json.NewEncoder(clientConn)
	in := json.NewDecoder(clientConn)

	if err := out.Encode(request); err != nil {
		t.Fatal(err)
	}

	var msg jsonrpcMessage
	if err := in.Decode(&msg); err != nil {
		t.Fatal(err)
	}

	if !msg.isResponse() {
		t.Fatal("message is not response")
	}
}

func TestServerMethodExecution(t *testing.T) {
	testServerMethodExecution(t, "echo")
}

func TestServerMethodWithCtx(t *testing.T) {
	testServerMethodExecution(t, "echoWithCtx")
}

// withBatchLimits temporarily overrides the package-level batch limit vars and
// restores them on test cleanup.
func withBatchLimits(t *testing.T, requestLimit, responseMax int) {
	t.Helper()
	origReq, origResp := BatchRequestLimit, BatchResponseMaxSize
	BatchRequestLimit = requestLimit
	BatchResponseMaxSize = responseMax
	t.Cleanup(func() {
		BatchRequestLimit = origReq
		BatchResponseMaxSize = origResp
	})
}

// runBatchRequest sends a batch via ServeSingleRequest and decodes the response.
// Returns nil if the server wrote nothing (e.g. notifications-only batch).
func runBatchRequest(t *testing.T, server *Server, batch []map[string]interface{}) []*jsonrpcMessage {
	t.Helper()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		server.ServeSingleRequest(context.Background(), NewCodec(serverConn))
	}()

	if err := json.NewEncoder(clientConn).Encode(batch); err != nil {
		t.Fatalf("encode batch: %v", err)
	}
	clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Server writes either a single message (for over-limit empty/notification batches)
	// or an array. Decode generically.
	dec := json.NewDecoder(clientConn)
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	if raw[0] == '[' {
		var resps []*jsonrpcMessage
		if err := json.Unmarshal(raw, &resps); err != nil {
			t.Fatalf("decode response array: %v", err)
		}
		return resps
	}
	var single jsonrpcMessage
	if err := json.Unmarshal(raw, &single); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return []*jsonrpcMessage{&single}
}

func TestServerBatchRequestLimit(t *testing.T) {
	withBatchLimits(t, 3, 0)

	server := newTestServer("test", new(Service))
	defer server.Stop()

	// 4 calls — exceeds the limit of 3.
	batch := make([]map[string]interface{}, 4)
	for i := range batch {
		batch[i] = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      i + 100,
			"method":  "test_echo",
			"params":  []interface{}{"x", 1, &Args{S: "a"}},
		}
	}

	resps := runBatchRequest(t, server, batch)
	if len(resps) != 1 {
		t.Fatalf("expected 1 error response, got %d", len(resps))
	}
	if resps[0].Error == nil {
		t.Fatalf("expected error response, got: %+v", resps[0])
	}
	if resps[0].Error.Message != errMsgBatchTooLarge {
		t.Errorf("error message = %q, want %q", resps[0].Error.Message, errMsgBatchTooLarge)
	}
	if string(resps[0].ID) != "100" {
		t.Errorf("error response id = %s, want first call id 100", string(resps[0].ID))
	}
}

func TestServerBatchRequestLimitNotificationsOnly(t *testing.T) {
	withBatchLimits(t, 2, 0)

	server := newTestServer("test", new(Service))
	defer server.Stop()

	// 3 notifications — exceeds the limit of 2. No id field.
	batch := []map[string]interface{}{
		{"jsonrpc": "2.0", "method": "test_noArgsRets"},
		{"jsonrpc": "2.0", "method": "test_noArgsRets"},
		{"jsonrpc": "2.0", "method": "test_noArgsRets"},
	}

	resps := runBatchRequest(t, server, batch)
	if len(resps) != 1 {
		t.Fatalf("expected 1 error response, got %d", len(resps))
	}
	if resps[0].Error == nil || resps[0].Error.Message != errMsgBatchTooLarge {
		t.Fatalf("expected batchTooLarge error, got: %+v", resps[0])
	}
	// errorMessage uses the package-level `null` raw JSON for the id; that's the
	// JSON-RPC null-id convention.
	if string(resps[0].ID) != "null" {
		t.Errorf("expected null id, got %q", string(resps[0].ID))
	}
}

func TestServerBatchResponseMaxSize(t *testing.T) {
	// Each test_echo response Result is ~80 bytes for the payload below.
	// Cap at 200 bytes -> first 2 succeed, 3rd is what tips us over, remaining
	// calls receive responseTooLarge.
	withBatchLimits(t, 0, 200)

	server := newTestServer("test", new(Service))
	defer server.Stop()

	const numCalls = 5
	batch := make([]map[string]interface{}, numCalls)
	for i := range batch {
		batch[i] = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      i + 1,
			"method":  "test_echo",
			"params":  []interface{}{fmt.Sprintf("payload-%d", i), i, &Args{S: "abcde"}},
		}
	}

	resps := runBatchRequest(t, server, batch)
	if len(resps) != numCalls {
		t.Fatalf("expected %d responses (mix of success + errors), got %d", numCalls, len(resps))
	}

	var sawSuccess, sawTooLarge bool
	for _, r := range resps {
		if r.Error == nil {
			sawSuccess = true
			continue
		}
		if r.Error.Code == errcodeResponseTooLarge {
			sawTooLarge = true
		}
	}
	if !sawSuccess {
		t.Errorf("expected at least one successful response before the limit was hit")
	}
	if !sawTooLarge {
		t.Errorf("expected at least one responseTooLarge error after the limit was hit")
	}
}

func TestServerBatchLimitsDisabled(t *testing.T) {
	withBatchLimits(t, 0, 0)

	server := newTestServer("test", new(Service))
	defer server.Stop()

	const numCalls = 50
	batch := make([]map[string]interface{}, numCalls)
	for i := range batch {
		batch[i] = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      i + 1,
			"method":  "test_echo",
			"params":  []interface{}{"x", i, &Args{S: "a"}},
		}
	}

	resps := runBatchRequest(t, server, batch)
	if len(resps) != numCalls {
		t.Fatalf("expected %d responses, got %d", numCalls, len(resps))
	}
	for i, r := range resps {
		if r.Error != nil {
			t.Errorf("response %d had unexpected error: %+v", i, r.Error)
		}
	}
}

// This test checks that responses are delivered for very short-lived connections that
// only carry a single request.
func TestServerShortLivedConn(t *testing.T) {
	server := newTestServer("service", new(Service))
	defer server.Stop()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener)

	var (
		request  = `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
		wantResp = `{"jsonrpc":"2.0","id":1,"result":{"rpc":"1.0","service":"1.0"}}` + "\n"
		deadline = time.Now().Add(10 * time.Second)
	)
	for range 20 {
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Fatal("can't dial:", err)
		}
		defer conn.Close()
		conn.SetDeadline(deadline)
		// Write the request, then half-close the connection so the server stops reading.
		conn.Write([]byte(request))
		conn.(*net.TCPConn).CloseWrite()
		// Now try to get the response.
		buf := make([]byte, 2000)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal("read error:", err)
		}
		if !bytes.Equal(buf[:n], []byte(wantResp)) {
			t.Fatalf("wrong response: %s", buf[:n])
		}
	}
}
