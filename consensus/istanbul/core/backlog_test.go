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

package core

import (
	"math"
	"math/big"
	"sync"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/prque"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestBacklogCore() *core {
	qualified := valset.NewAddressSet(nil)
	return &core{
		address:            common.HexToAddress("0xdead"),
		state:              StateAcceptRequest,
		logger:             logger.NewWith(),
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		backlogSenderBytes: make(map[common.Address]uint64),
		current:            newRoundState(&bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}, qualified, common.Hash{}, nil, nil, nil),
	}
}

// backlogSender returns a sender address that never collides with the address of
// the core under test, so senders can be numbered from one.
func backlogSender(i int) common.Address {
	return common.BigToAddress(big.NewInt(int64(i)))
}

func newTestBacklogMessage(t *testing.T, sequence int64) *bft.Message {
	t.Helper()
	payload, err := bft.Encode(&bft.Subject{View: &bft.View{
		Sequence: big.NewInt(sequence),
		Round:    big.NewInt(0),
	}})
	require.NoError(t, err)
	return &bft.Message{Code: bft.MsgPrepare, Msg: payload}
}

func TestStoreBacklogBoundsMessagesPerSender(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	msg := newTestBacklogMessage(t, 2)

	for range maxBacklogMessagesPerSender + 1 {
		c.storeBacklog(msg, src)
	}

	assert.Equal(t, maxBacklogMessagesPerSender, c.backlogs[src].Size())
	assert.Equal(t, maxBacklogMessagesPerSender, c.backlogMessageCount)
}

func TestStoreBacklogBoundsPayloadBytesPerSender(t *testing.T) {
	src := common.HexToAddress("0x1")
	msg := newTestBacklogMessage(t, 2)
	c := newTestBacklogCore()
	msg.Signature = make([]byte, int(maxBacklogPayloadBytesPerSender-backlogMessageBytes(msg)+1))

	c.storeBacklog(msg, src)

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogSenderBytes)
	assert.Zero(t, c.backlogMessageCount)
	assert.Zero(t, c.backlogPayloadBytes)
}

func TestStoreBacklogBoundsMessagesAcrossSenders(t *testing.T) {
	require.Zero(t, maxBacklogMessages%maxBacklogMessagesPerSender,
		"the test fills the global cap with senders that each reach the per-sender cap")
	senders := maxBacklogMessages / maxBacklogMessagesPerSender
	c := newTestBacklogCore()
	msg := newTestBacklogMessage(t, 2)

	for sender := 1; sender <= senders; sender++ {
		for range maxBacklogMessagesPerSender {
			c.storeBacklog(msg, backlogSender(sender))
		}
	}
	rejected := backlogSender(senders + 1)
	c.storeBacklog(msg, rejected)

	assert.Equal(t, maxBacklogMessages, c.backlogMessageCount)
	assert.Len(t, c.backlogs, senders)
	assert.NotContains(t, c.backlogs, rejected)
}

func TestStoreBacklogBoundsPayloadBytesAcrossSenders(t *testing.T) {
	require.Zero(t, maxBacklogPayloadBytes%maxBacklogPayloadBytesPerSender,
		"the test fills the global cap with senders that each reach the per-sender byte cap")
	senders := maxBacklogPayloadBytes / maxBacklogPayloadBytesPerSender
	msg := newTestBacklogMessage(t, 2)
	// One message per sender, sized to exactly the per-sender byte cap.
	msg.Signature = make([]byte, int(maxBacklogPayloadBytesPerSender-backlogMessageBytes(msg)))
	c := newTestBacklogCore()

	for sender := 1; sender <= senders; sender++ {
		c.storeBacklog(msg, backlogSender(sender))
	}
	rejected := backlogSender(senders + 1)
	c.storeBacklog(msg, rejected)

	assert.Equal(t, senders, c.backlogMessageCount)
	assert.Equal(t, uint64(maxBacklogPayloadBytes), c.backlogPayloadBytes)
	assert.NotContains(t, c.backlogs, rejected)
}

func TestStoreBacklogBoundsFutureSequence(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()

	c.storeBacklog(newTestBacklogMessage(t, 1+maxBacklogSequencesAhead), src)
	c.storeBacklog(newTestBacklogMessage(t, 2+maxBacklogSequencesAhead), src)

	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogMessageCount)
}

func TestBacklogRejectsSequenceOutsideUint64(t *testing.T) {
	c := newTestBacklogCore()
	// Put the local sequence at the top of uint64 so the window alone would
	// admit the message; only the uint64 check may reject it.
	c.current = newRoundState(
		&bft.View{Sequence: new(big.Int).SetUint64(math.MaxUint64), Round: big.NewInt(0)},
		valset.NewAddressSet(nil), common.Hash{}, nil, nil, nil)
	tooLarge := new(big.Int).Add(new(big.Int).SetUint64(math.MaxUint64), big.NewInt(1))

	assert.True(t, c.isBacklogSequenceTooFar(tooLarge))
}

func TestStoreBacklogSkipsUndecodableMessage(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()

	c.storeBacklog(&bft.Message{Code: bft.MsgPrepare, Msg: []byte{0xff}}, src)

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogSenderBytes)
	assert.Zero(t, c.backlogMessageCount)
	assert.Zero(t, c.backlogPayloadBytes)
}

func TestProcessBacklogFreesCapacityForLaterMessages(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	qualified := valset.NewAddressSet(nil)
	msg := newTestBacklogMessage(t, 2)

	for range maxBacklogMessagesPerSender {
		c.storeBacklog(msg, src)
	}
	c.current = newRoundState(&bft.View{Sequence: big.NewInt(3), Round: big.NewInt(0)}, qualified, common.Hash{}, nil, nil, nil)
	c.processBacklog()

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogSenderBytes)
	assert.Zero(t, c.backlogMessageCount)
	assert.Zero(t, c.backlogPayloadBytes)

	c.storeBacklog(msg, src)
	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogMessageCount)
}

func TestProcessBacklogRemovesMessageWithNilView(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	payload, err := bft.Encode(&bft.Subject{})
	require.NoError(t, err)
	msg := &bft.Message{Code: bft.MsgPrepare, Msg: payload}

	c.backlogs[src] = prque.New()
	c.backlogs[src].Push(msg, 0)
	c.backlogMessageCount = 1
	c.backlogSenderBytes[src] = backlogMessageBytes(msg)
	c.backlogPayloadBytes = backlogMessageBytes(msg)
	c.processBacklog()

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogSenderBytes)
	assert.Zero(t, c.backlogMessageCount)
	assert.Zero(t, c.backlogPayloadBytes)
}

func TestExceedsBacklogLimit(t *testing.T) {
	assert.False(t, exceedsBacklogLimit(15, 1, 16))
	assert.True(t, exceedsBacklogLimit(15, 2, 16))
	assert.True(t, exceedsBacklogLimit(0, 17, 16))
}
