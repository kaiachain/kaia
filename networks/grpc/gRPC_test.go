// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package grpc

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	TEST_BLOCK_NUMBER = float64(123456789)
)

type APIgRPC struct{}

func (a APIgRPC) BlockNumber() float64 {
	return TEST_BLOCK_NUMBER
}

func TestGRPC(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	wg := &sync.WaitGroup{}
	wg.Add(3)

	addr := "127.0.0.1:4000"
	handler := rpc.NewServer()

	handler.RegisterName("kaia", &APIgRPC{})

	listener := &Listener{Addr: addr}
	listener.SetRPCServer(handler)
	go listener.Start()
	defer listener.Stop()

	time.Sleep(2 * time.Second)

	go testCall(t, addr, wg)
	go testBiCall(t, addr, wg)
	go testSubscribe(t, addr, wg)
	wg.Wait()
}

func testCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKaiaClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKaiaClient(timeout)
	assert.NoError(t, err)

	request, err := kclient.makeRPCRequest("kaia", "kaia_blockNumber", nil)
	assert.NoError(t, err)

	response, err := knclient.Call(kclient.ctx, request)
	assert.NoError(t, err)

	err = kclient.handleRPCResponse(response)
	assert.NoError(t, err)

	var out jsonSuccessResponse
	json.Unmarshal(response.Payload, &out)
	assert.NoError(t, err)
	assert.Equal(t, out.Result, TEST_BLOCK_NUMBER)
}

func testBiCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKaiaClient(addr)

	knclient, err := kclient.makeKaiaClient(timeout)
	assert.NoError(t, err)

	stream, _ := knclient.BiCall(kclient.ctx)
	done := make(chan struct{})
	// Close must precede <-done: closing the connection causes handleBiCall's
	// internal goroutines to unblock from their stream operations, which lets
	// handleBiCall return and signal done. Reversing the order deadlocks.
	defer func() {
		kclient.Close()
		<-done
	}()

	go func() {
		kclient.handleBiCall(stream, func() (request *RPCRequest, e error) {
			request, err := kclient.makeRPCRequest("kaia", "kaia_blockNumber", nil)
			if assert.NoError(t, err) {
				return request, err
			}
			return request, nil
		}, func(response *RPCResponse) error {
			var out jsonSuccessResponse
			err = json.Unmarshal(response.Payload, &out)
			assert.NoError(t, err)
			assert.Equal(t, out.Result, TEST_BLOCK_NUMBER)
			return nil
		})
		close(done)
	}()

	time.Sleep(3 * time.Second)
}

func testSubscribe(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKaiaClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKaiaClient(timeout)
	assert.NoError(t, err)

	request, err := kclient.makeRPCRequest("kaia", "kaia_blockNumber", nil)
	assert.NoError(t, err)

	stream, err := knclient.Subscribe(kclient.ctx, request)
	assert.NoError(t, err)

	kclient.handleSubscribe(stream, func(response *RPCResponse) error {
		var out jsonSuccessResponse
		err = json.Unmarshal(response.Payload, &out)
		assert.NoError(t, err)
		assert.Equal(t, out.Result, TEST_BLOCK_NUMBER)
		return nil
	})
}

// newNotificationTestServer starts an in-process gRPC server and returns a
// connected client, the underlying API stub, and a cleanup function.
func newNotificationTestServer(t *testing.T) (KlaytnNodeClient, *notificationAPI, func()) {
	t.Helper()
	handler := rpc.NewServer()
	api := &notificationAPI{}
	if err := handler.RegisterName("kaia", api); err != nil {
		t.Fatal(err)
	}
	lis, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := googlegrpc.NewServer()
	RegisterKlaytnNodeServer(srv, &kaiaServer{handler: handler})
	go srv.Serve(lis)

	dialCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := googlegrpc.DialContext(dialCtx, lis.Addr().String(),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
		googlegrpc.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	return NewKlaytnNodeClient(conn), api, func() { conn.Close(); srv.Stop() }
}

type notificationAPI struct{ calls int32 }

func (a *notificationAPI) BlockNumber() float64 {
	atomic.AddInt32(&a.calls, 1)
	return TEST_BLOCK_NUMBER
}

// TestNotificationRejected verifies that JSON-RPC notifications (no "id" field)
// sent over gRPC unary Call are rejected with InvalidArgument before dispatch.
func TestNotificationRejected(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	client, api, cleanup := newNotificationTestServer(t)
	defer cleanup()

	notification := []byte(`{"jsonrpc":"2.0","method":"kaia_blockNumber","params":[]}`)

	for i := 0; i < 5; i++ {
		_, err := client.Call(context.Background(), &RPCRequest{Params: notification})
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("request %d: expected InvalidArgument, got %v", i, err)
		}
	}

	if n := atomic.LoadInt32(&api.calls); n != 0 {
		t.Errorf("method was executed %d time(s); want 0 (notification must not be dispatched)", n)
	}
}

// TestSubscribeNotificationRejected verifies that JSON-RPC notifications sent
// over gRPC Subscribe are rejected with InvalidArgument before dispatch.
func TestSubscribeNotificationRejected(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	client, api, cleanup := newNotificationTestServer(t)
	defer cleanup()

	notification := []byte(`{"jsonrpc":"2.0","method":"kaia_blockNumber","params":[]}`)
	stream, err := client.Subscribe(context.Background(), &RPCRequest{Params: notification})
	if err == nil {
		_, err = stream.Recv()
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", err)
	}
	if n := atomic.LoadInt32(&api.calls); n != 0 {
		t.Errorf("method was executed %d time(s); want 0 (notification must not be dispatched)", n)
	}
}
