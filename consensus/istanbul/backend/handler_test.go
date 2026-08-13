// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
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

package backend

import (
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
)

func TestBackend_HandleMsg(t *testing.T) {
	_, backend := newBlockChain(t, 1)
	defer backend.Stop()
	eventSub := backend.istanbulEventMux.Subscribe(istanbul.MessageEvent{})

	addr := common.StringToAddress("test addr")
	data := &bft.ConsensusMsg{
		PrevHash: common.HexToHash("0x1234"),
		Payload:  []byte("test data"),
	}
	hash := istanbul.RLPHash(data.Payload)
	size, payload, _ := rlp.EncodeToReader(data)

	// Success case
	{
		msg := p2p.Msg{
			Code:    consensus.ConsensusMsgCode,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Nil(t, err)
		assert.True(t, isHandled)

		if err != nil {
			t.Fatalf("handle message failed: %v", err)
		}

		recentMsg, ok := backend.recentMessages.Get(addr)
		assert.True(t, ok)

		cachedMsg, ok := recentMsg.(*lru.ARCCache)
		assert.True(t, ok)

		value, ok := cachedMsg.Get(hash)
		assert.True(t, ok)
		assert.True(t, value.(bool))

		value, ok = backend.knownMessages.Get(hash)
		assert.True(t, ok)
		assert.True(t, value.(bool))

		evTimer := time.NewTimer(3 * time.Second)
		defer evTimer.Stop()

		select {
		case event := <-eventSub.Chan():
			switch ev := event.Data.(type) {
			case istanbul.MessageEvent:
				assert.Equal(t, data.Payload, ev.Payload)
				assert.Equal(t, data.PrevHash, ev.Hash)
			default:
				t.Fatal("unexpected message type")
			}
		case <-evTimer.C:
			t.Fatal("failed to subscribe istanbul message event")
		}
	}

	// Failure case - undefined message code
	{
		msg := p2p.Msg{
			Code:    0x99,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, nil, err)
		assert.False(t, isHandled)
	}

	// Failure case - invalid message data
	{
		size, payload, _ := rlp.EncodeToReader([]byte{0x1, 0x2})
		msg := p2p.Msg{
			Code:    consensus.ConsensusMsgCode,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, errDecodeFailed, err)
		assert.True(t, isHandled)
	}

	// Failure case - stopped istanbul engine
	{
		msg := p2p.Msg{
			Code:    consensus.ConsensusMsgCode,
			Size:    uint32(size),
			Payload: payload,
		}
		_ = backend.Stop()
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, istanbul.ErrStoppedEngine, err)
		assert.True(t, isHandled)
	}
}

func TestBackend_ValidatePeerType(t *testing.T) {
	_, backend := newBlockChain(t, 1)
	defer backend.Stop()

	// Return nil if the input address is a validator
	{
		err := backend.ValidatePeerType(backend.address)
		assert.Nil(t, err)
	}

	// Return an error if the input address is invalid
	{
		err := backend.ValidatePeerType(common.Address{})
		assert.Equal(t, errInvalidPeerAddress, err)
	}
}

// startBlockedValidatePeer creates a test backend and starts a goroutine that
// blocks on ValidatePeerType. Returns the backend and a channel that closes
// when ValidatePeerType returns. Asserts that ValidatePeerType is still
// blocking after 200ms.
func startBlockedValidatePeer(t *testing.T) (*backend, chan struct{}) {
	t.Helper()
	b := newTestBackend()

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() { recover() }() // valsetModule is nil in test backend
		b.ValidatePeerType(common.Address{})
	}()

	select {
	case <-done:
		t.Fatal("ValidatePeerType should block before engine is ready")
	case <-time.After(200 * time.Millisecond):
	}
	return b, done
}

func TestValidatePeerType_BlocksUntilSignal(t *testing.T) {
	b, done := startBlockedValidatePeer(t)
	defer b.Stop()

	b.SignalPeerRegistrable()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ValidatePeerType should unblock after SignalPeerRegistrable")
	}
}

func TestValidatePeerType_UnblocksOnStop(t *testing.T) {
	b, done := startBlockedValidatePeer(t)
	defer b.Stop()

	b.Stop()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ValidatePeerType should unblock after Stop")
	}
}
