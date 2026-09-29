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

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #1053: future messages are retained per sender only, within a sequence
// window and a per-sender count, with a single PRE-PREPARE slot per sender.
func TestDevParity_BacklogIsBounded(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	m := newTestMachine(b, 10, 0, nil, []common.Address{crypto.PubkeyToAddress(key.PublicKey)})
	src := common.HexToAddress("0x00000000000000000000000000000000000000bb")

	m.storeBacklog(subjectMsg(t, bft.MsgPrepare, src, 10+maxBacklogSequencesAhead+1, 0), src)
	assert.Zero(t, m.backlogCounts[src], "a sequence past the window is dropped")

	for i := 0; i < maxBacklogMessagesPerSender+5; i++ {
		m.storeBacklog(subjectMsg(t, bft.MsgPrepare, src, 11, int64(i)), src)
	}
	assert.Equal(t, maxBacklogMessagesPerSender, m.backlogCounts[src])
	assert.Equal(t, maxBacklogMessagesPerSender, m.backlogs[src].Size())

	ppAt := func(round int64) *bft.Message {
		pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(11), Round: big.NewInt(round)}, Proposal: newEmptyBlock(11, common.Hash{1})}
		return &bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp), Address: src}
	}
	m.storeBacklog(ppAt(1), src)
	m.storeBacklog(ppAt(2), src)
	m.storeBacklog(ppAt(0), src)
	require.Contains(t, m.backlogPreprepares, src)
	assert.Equal(t, int64(2), m.backlogPreprepares[src].view.Round.Int64(), "only a higher view replaces the retained PRE-PREPARE")
}

// dev #1053: ROUND CHANGE retention is limited to a round window and a quorum's
// worth of senders per round.
func TestDevParity_RoundChangeRetentionIsBounded(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	a := crypto.PubkeyToAddress(key.PublicKey)
	c := common.HexToAddress("0x00000000000000000000000000000000000000cc")
	d := common.HexToAddress("0x00000000000000000000000000000000000000dd")
	m := newTestMachine(b, 10, 0, nil, []common.Address{a, c, d})
	m.requiredMessageCount = 2

	_, err := m.addRoundChange(big.NewInt(maxRoundChangeRoundsAhead+1), subjectMsg(t, bft.MsgRoundChange, a, 10, maxRoundChangeRoundsAhead+1))
	assert.ErrorIs(t, err, errRoundChangeTooFar)

	n, err := m.addRoundChange(big.NewInt(1), subjectMsg(t, bft.MsgRoundChange, a, 10, 1))
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	n, err = m.addRoundChange(big.NewInt(1), subjectMsg(t, bft.MsgRoundChange, c, 10, 1))
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	_, err = m.addRoundChange(big.NewInt(1), subjectMsg(t, bft.MsgRoundChange, d, 10, 1))
	assert.ErrorIs(t, err, errRoundChangeLimit)
	n, err = m.addRoundChange(big.NewInt(1), subjectMsg(t, bft.MsgRoundChange, a, 10, 1))
	require.NoError(t, err, "a retained sender may replace its own message")
	assert.Equal(t, 2, n)
}

// dev #1053: PREPARE, COMMIT and ROUND CHANGE are size-bounded before any
// retention path.
func TestDevParity_SubjectMessageSizeIsBounded(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	m := newTestMachine(b, 10, 0, nil, []common.Address{self})
	for _, code := range []uint64{bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange} {
		oversized := &bft.Message{Code: code, Msg: make([]byte, maxSubjectMessageBytes+1), Address: self}
		assert.ErrorIs(t, m.handleCheckedMsg(oversized, self), errMessageTooLarge, "code %d", code)
	}
	// A PREPREPARE carries a block and is exempt from the subject size limit.
	bigPreprepare := &bft.Message{Code: bft.MsgPreprepare, Msg: make([]byte, maxSubjectMessageBytes+1), Address: self}
	assert.NoError(t, checkMessageSize(bigPreprepare))
}

// The round-change window follows the round a node catches up to.
func TestDevParity_RoundChangeWindowFollowsCatchUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig(), current: newEmptyBlock(9, common.Hash{})})
	self := crypto.PubkeyToAddress(key.PublicKey)
	m := newTestMachine(b, 10, 0, nil, []common.Address{self})
	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(self, nil).AnyTimes()
	b.valsetModule = mValset

	m.catchUpRound(&bft.View{Sequence: big.NewInt(10), Round: big.NewInt(5)})
	defer m.stopTimer()
	_, err := m.addRoundChange(big.NewInt(5+maxRoundChangeRoundsAhead), subjectMsg(t, bft.MsgRoundChange, self, 10, 5+maxRoundChangeRoundsAhead))
	assert.NoError(t, err)
	_, err = m.addRoundChange(big.NewInt(6+maxRoundChangeRoundsAhead), subjectMsg(t, bft.MsgRoundChange, self, 10, 6+maxRoundChangeRoundsAhead))
	assert.ErrorIs(t, err, errRoundChangeTooFar)
}

// Delivering retained messages releases the sender's budget, so a sender is
// not locked out after maxBacklogMessagesPerSender messages.
func TestDevParity_ProcessBacklogReleasesSenderBudget(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	m := newTestMachine(b, 10, 0, nil, []common.Address{crypto.PubkeyToAddress(key.PublicKey)})
	src := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	for i := 0; i < maxBacklogMessagesPerSender; i++ {
		m.storeBacklog(subjectMsg(t, bft.MsgPrepare, src, 11, 0), src)
	}
	require.Equal(t, maxBacklogMessagesPerSender, m.backlogCounts[src])

	// Move to the view the messages belong to; they are no longer future.
	m.sequence, m.round, m.state = big.NewInt(11), big.NewInt(0), statePreprepared
	m.processBacklog()
	assert.Zero(t, m.backlogCounts[src])
	assert.NotContains(t, m.backlogs, src)

	m.storeBacklog(subjectMsg(t, bft.MsgPrepare, src, 12, 0), src)
	assert.Equal(t, 1, m.backlogCounts[src], "the sender can queue again")
}
