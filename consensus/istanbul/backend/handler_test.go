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
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackend_HandleMsg(t *testing.T) {
	_, backend := newBlockChain(1)
	defer backend.Stop()
	eventSub := backend.istanbulEventMux.Subscribe(istanbul.MessageEvent{})

	addr := common.StringToAddress("test addr")
	data := &istanbul.ConsensusMsg{
		PrevHash: common.HexToHash("0x1234"),
		Payload:  []byte("test data"),
	}
	hash := istanbul.RLPHash(data.Payload)
	size, payload, _ := rlp.EncodeToReader(data)

	// Success case
	{
		msg := p2p.Msg{
			Code:    IstanbulMsg,
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
			Code:    IstanbulMsg,
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
			Code:    IstanbulMsg,
			Size:    uint32(size),
			Payload: payload,
		}
		_ = backend.Stop()
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, istanbul.ErrStoppedEngine, err)
		assert.True(t, isHandled)
	}
}

func TestBackend_Protocol(t *testing.T) {
	backend := newTestBackend()
	defer backend.Stop()

	assert.Equal(t, IstanbulProtocol, backend.Protocol())
}

func TestBackend_ValidatePeerType(t *testing.T) {
	t.Run("permissioned", func(t *testing.T) {
		_, b := newBlockChain(1)
		defer b.Stop()

		assert.NoError(t, b.ValidatePeerType(b.address))
		assert.Equal(t, errInvalidPeerAddress, b.ValidatePeerType(common.HexToAddress("0xdead")))
	})

	t.Run("permissionless", func(t *testing.T) {
		b, cleanup := newPermissionlessBlockChain(t, 1)
		defer cleanup()

		cnPeers, err := b.valsetModule.GetCNPeers(1)
		require.NoError(t, err)
		require.NotEmpty(t, cnPeers)
		assert.NoError(t, b.ValidatePeerType(cnPeers[0]))
		assert.Equal(t, errInvalidPeerAddress, b.ValidatePeerType(common.HexToAddress("0xdead")))
	})
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

func TestGetHeaderGovVoters(t *testing.T) {
	t.Run("permissioned", func(t *testing.T) {
		_, b := newBlockChain(1)
		defer b.Stop()

		voters, err := b.valsetModule.GetHeaderGovVoters(1)
		require.NoError(t, err)
		assert.Contains(t, voters, b.address)
	})

	t.Run("permissionless", func(t *testing.T) {
		b, cleanup := newPermissionlessBlockChain(t, 1)
		defer cleanup()

		voters, err := b.valsetModule.GetHeaderGovVoters(1)
		require.NoError(t, err)
		// MakeTestPermissionlessConfig(1) registers 1 ValActive node → 1 gov voter
		assert.Len(t, voters, 1)
	})
}

func newPermissionlessBlockChain(t *testing.T, n int) (*backend, func()) {
	t.Helper()
	config := params.TestChainConfig.Copy()
	config.PermissionlessCompatibleBlock = big.NewInt(0)

	ctrl := gomock.NewController(t)
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()

	_, b := newBlockChain(n, config, mockVRank)
	return b, func() {
		b.Stop()
		ctrl.Finish()
	}
}
