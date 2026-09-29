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

package kaiabft

import (
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/mocks"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #912: peer type validation follows the CN peer allowlist, not the council.
func TestDevParity_ValidatePeerTypeUsesCNPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	head := newEmptyBlock(5, common.Hash{})
	b, _ := newTestBackend(t, &fakeChain{cfg: permissionlessConfig(), current: head})
	b.signalPeerRegistrable()
	member := common.HexToAddress("0x00000000000000000000000000000000000000a1")
	councilOnly := common.HexToAddress("0x00000000000000000000000000000000000000a2")

	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCouncil(gomock.Any()).Return([]common.Address{member, councilOnly}, nil).AnyTimes()
	mValset.EXPECT().GetCNPeers(uint64(6)).Return([]common.Address{member}, nil).Times(2)
	b.valsetModule = mValset
	assert.NoError(t, b.ValidatePeerType(member))
	assert.Error(t, b.ValidatePeerType(councilOnly), "a council member outside the CN peer list is rejected")

	mValset.EXPECT().GetCNPeers(uint64(6)).Return(nil, nil)
	assert.NoError(t, b.ValidatePeerType(councilOnly), "a nil CN peer list disables the filter")
}

// dev #936: a validator outside the qualified set does not execute or seal.
func TestDevParity_SubmitTransactionsSkipsNonQualified(t *testing.T) {
	ctrl := gomock.NewController(t)
	b, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	other := common.HexToAddress("0x00000000000000000000000000000000000000ee")
	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetQualifiedValidators(uint64(3)).Return([]common.Address{other}, nil)
	b.valsetModule = mValset
	b.executor = mocks.NewMockExecutor(ctrl) // no calls expected
	b.currentView.Store(&bft.View{Sequence: big.NewInt(99), Round: big.NewInt(0)})

	resultCh := b.SubmitTransactions(nil, nil, &types.Header{Number: big.NewInt(3)}, nil)
	select {
	case res := <-resultCh:
		assert.Nil(t, res)
	case <-time.After(2 * time.Second):
		t.Fatal("SubmitTransactions did not return")
	}
}

// dev #944: NewChainHead does not take coreMu, so a coreMu holder waiting on
// the worker loop cannot deadlock with it.
func TestDevParity_NewChainHeadDoesNotTakeCoreMu(t *testing.T) {
	b, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	b.coreStarted.Store(true)
	b.coreMu.Lock()
	defer b.coreMu.Unlock()
	done := make(chan error, 1)
	go func() { done <- b.NewChainHead() }()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("NewChainHead blocked on coreMu")
	}
}

// dev #1029: an inbound consensus message is posted synchronously and outside
// coreMu, so a stalled consumer applies backpressure to the peer.
func TestDevParity_HandleMsgAppliesBackpressureOutsideCoreMu(t *testing.T) {
	b, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	b.coreStarted.Store(true)
	sub := b.eventMux.Subscribe(messageEvent{})
	defer sub.Unsubscribe()

	data := &bft.ConsensusMsg{Payload: []byte("backpressure test")}
	size, payload, err := rlp.EncodeToReader(data)
	require.NoError(t, err)
	done := make(chan error, 1)
	go func() {
		_, err := b.HandleMsg(common.HexToAddress("0x01"), p2p.Msg{Code: consensus.ConsensusMsgCode, Size: uint32(size), Payload: payload})
		done <- err
	}()

	hash := bft.RLPHash(data.Payload)
	deadline := time.After(time.Second)
	for {
		b.coreMu.Lock()
		_, known := b.knownMessages.Get(hash)
		b.coreMu.Unlock()
		if known {
			break
		}
		select {
		case <-deadline:
			t.Fatal("HandleMsg did not prepare the event")
		case <-time.After(time.Millisecond):
		}
	}
	select {
	case err := <-done:
		t.Fatalf("HandleMsg returned before the event was consumed: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	require.True(t, b.coreMu.TryLock(), "HandleMsg must not hold coreMu while blocked")
	b.coreMu.Unlock()

	<-sub.Chan()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("HandleMsg did not return after the event was consumed")
	}
}
