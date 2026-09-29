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
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #923: an accepted PRE-PREPARE is relayed to the VRank module, both when
// a validator accepts a peer's proposal and when the proposer self-accepts.
func TestDevParity_AcceptedPreprepareReachesVRank(t *testing.T) {
	ctrl := gomock.NewController(t)
	parent := newEmptyBlock(9, common.Hash{})
	chain := &fakeChain{cfg: legacyConfig(), current: parent, headers: map[common.Hash]*types.Header{parent.Hash(): parent.Header()}}
	b, key := newTestBackend(t, chain)
	self := crypto.PubkeyToAddress(key.PublicKey)
	proposerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)

	got := make(chan *bft.View, 2)
	mVRank := vrank_mock.NewMockVRankModule(ctrl)
	mVRank.EXPECT().HandleIstanbulPreprepare(gomock.Any(), gomock.Any()).Do(func(_ *types.Block, view *bft.View) { got <- view }).Times(2)
	b.RegisterVRankModule(mVRank)
	b.startPrepreparedRelay()
	defer b.stopPrepreparedRelay()

	// A validator accepts the proposer's PRE-PREPARE.
	block := newEmptyBlock(10, parent.Hash())
	m := newTestMachine(b, 10, 0, nil, []common.Address{self, proposer})
	m.state = stateAcceptRequest
	m.proposer = proposer
	pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(0)}, Proposal: block}
	require.NoError(t, m.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp)}, proposer))

	// The proposer accepts its own proposal when it comes back through the
	// self-loop (backend.broadcast posts every message to the local machine).
	m2 := newTestMachine(b, 10, 1, nil, []common.Address{self, proposer})
	m2.state = stateAcceptRequest
	m2.proposer = self
	own := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(1)}, Proposal: newEmptyBlock(10, parent.Hash())}
	require.NoError(t, m2.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, own)}, self))

	for _, want := range []int64{0, 1} {
		select {
		case view := <-got:
			assert.Equal(t, int64(10), view.Sequence.Int64())
			assert.Equal(t, want, view.Round.Int64())
		case <-time.After(2 * time.Second):
			t.Fatalf("no VRank relay for round %d", want)
		}
	}
}

// dev #923: the relay runs exactly while the engine is started.
func TestDevParity_VRankRelayFollowsLifecycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	genesis := newEmptyBlock(0, common.Hash{})
	b, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig(), current: genesis})
	b.RegisterVRankModule(vrank_mock.NewMockVRankModule(ctrl))

	require.NoError(t, b.Start(b.chain, nil))
	assert.NotNil(t, b.prepreparedSub, "Start must start the VRank relay")
	require.NoError(t, b.Stop())
	assert.Nil(t, b.prepreparedSub, "Stop must stop the VRank relay")
}
