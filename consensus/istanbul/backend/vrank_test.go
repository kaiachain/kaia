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

package backend

import (
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBackend_startPrepreparedRelay asserts that, with a VRank module registered, a
// PrepreparedEvent posted on the Istanbul mux is forwarded once to the module's
// HandleIstanbulPreprepare with the same (block, view).
func TestBackend_startPrepreparedRelay(t *testing.T) {
	sb := newTestBackend()
	defer sb.Stop()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockVRank := vrank_mock.NewMockVRankModule(mockCtrl)

	// Capture the forwarded (block, view); the channel send/receive synchronizes the
	// relay goroutine's writes with the assertions below.
	var (
		gotBlock *types.Block
		gotView  *bft.View
		got      = make(chan struct{}, 1)
	)
	mockVRank.EXPECT().HandleIstanbulPreprepare(gomock.Any(), gomock.Any()).
		DoAndReturn(func(b *types.Block, v *bft.View) {
			gotBlock, gotView = b, v
			got <- struct{}{}
		}).Times(1)

	sb.RegisterVRankModule(mockVRank)
	sb.startPrepreparedRelay()
	defer sb.stopPrepreparedRelay()

	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(7)})
	view := &bft.View{Sequence: big.NewInt(7), Round: big.NewInt(0)}
	sb.istanbulEventMux.Post(istanbul.PrepreparedEvent{Block: block, View: view})

	select {
	case <-got:
		require.NotNil(t, gotBlock)
		assert.Equal(t, block.Hash(), gotBlock.Hash())
		require.NotNil(t, gotView)
		assert.Equal(t, view.Sequence.Uint64(), gotView.Sequence.Uint64())
		assert.Equal(t, view.Round.Uint64(), gotView.Round.Uint64())
	case <-time.After(3 * time.Second):
		t.Fatal("relay did not forward PrepreparedEvent to the VRank module")
	}
}

// TestBackend_startPrepreparedRelay_NoModule ensures the relay is a no-op when no
// VRank module is registered (e.g. on non-CN nodes), rather than panicking.
func TestBackend_startPrepreparedRelay_NoModule(t *testing.T) {
	sb := newTestBackend()
	defer sb.Stop()

	sb.startPrepreparedRelay()
	defer sb.stopPrepreparedRelay()

	// Posting an event must not panic or block even though nothing consumes it.
	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	sb.istanbulEventMux.Post(istanbul.PrepreparedEvent{Block: block, View: &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}})

	assert.Nil(t, sb.prepreparedSub, "relay must not subscribe without a VRank module")
}
