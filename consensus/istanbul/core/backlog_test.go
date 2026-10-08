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
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
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
		address:              common.HexToAddress("0xdead"),
		state:                StateAcceptRequest,
		logger:               logger.NewWith(),
		backlogs:             make(map[common.Address]*prque.Prque),
		backlogsMu:           new(sync.Mutex),
		backlogCounts:        make(map[common.Address]int),
		backlogEvidenceBytes: make(map[common.Address]uint64),
		backlogPreprepares:   make(map[common.Address]backlogPreprepare),
		current:              newRoundState(&bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}, qualified, common.Hash{}, nil, nil, nil, nil),
	}
}

// backlogSender returns a sender address that never collides with the address of
// the core under test, so senders can be numbered from one.
func backlogSender(i int) common.Address {
	return common.BigToAddress(big.NewInt(int64(i)))
}

func newTestBacklogMessage(t *testing.T, sequence int64) *bft.Message {
	t.Helper()
	payload, err := bft.Encode(&bft.Prepare{View: &bft.View{
		Sequence: big.NewInt(sequence),
		Round:    big.NewInt(0),
	}})
	require.NoError(t, err)
	return &bft.Message{Code: bft.MsgPrepare, Msg: payload}
}

func newTestBacklogRoundChange(t *testing.T, sequence int64) *bft.Message {
	t.Helper()
	payload, err := bft.Encode(&bft.RoundChange{View: &bft.View{
		Sequence: big.NewInt(sequence),
		Round:    big.NewInt(0),
	}})
	require.NoError(t, err)
	return &bft.Message{Code: bft.MsgRoundChange, Msg: payload}
}

func newTestBacklogPreprepare(t *testing.T, sequence, round int64) *bft.Message {
	t.Helper()
	payload, err := bft.Encode(&bft.Preprepare{
		View:     &bft.View{Sequence: big.NewInt(sequence), Round: big.NewInt(round)},
		Proposal: types.NewBlockWithHeader(&types.Header{Number: big.NewInt(sequence)}),
	})
	require.NoError(t, err)
	return &bft.Message{Code: bft.MsgPreprepare, Msg: payload}
}

func setTestBacklogView(c *core, sequence, round int64) {
	c.current = newRoundState(&bft.View{Sequence: big.NewInt(sequence), Round: big.NewInt(round)},
		valset.NewAddressSet(nil), common.Hash{}, nil, nil, nil, nil)
}

func TestStoreBacklogBoundsMessagesPerSender(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	msg := newTestBacklogMessage(t, 2)

	for range maxBacklogMessagesPerSender + 1 {
		c.storeBacklog(msg, src)
	}

	assert.Equal(t, maxBacklogMessagesPerSender, c.backlogs[src].Size())
	assert.Equal(t, maxBacklogMessagesPerSender, c.backlogCounts[src])
}

func TestStoreBacklogRetainsRoundChangeWithPreparedClaim(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	payload, err := bft.Encode(&bft.RoundChange{
		View: &bft.View{Sequence: big.NewInt(2), Round: big.NewInt(1)},
		Prepared: &bft.PreparedClaim{
			Round:  big.NewInt(0),
			Digest: common.HexToHash("0x1234"),
		},
	})
	require.NoError(t, err)
	msg := &bft.Message{Code: bft.MsgRoundChange, Msg: payload, Evidence: []byte{0xc0}}

	c.storeBacklog(msg, src)

	require.Contains(t, c.backlogs, src)
	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogCounts[src])
	assert.Equal(t, uint64(len(msg.Evidence)), c.backlogEvidenceBytes[src])
}

func TestStoreBacklogBoundsJustificationBytesPerSenderWithoutBlockingSmallMessages(t *testing.T) {
	src := backlogSender(1)
	c := newTestBacklogCore()
	msg := newTestBacklogRoundChange(t, 2)
	msg.Evidence = make([]byte, maxBacklogEvidenceBytesPerSender/2+1)
	justificationBytes := uint64(len(msg.Evidence))

	c.storeBacklog(msg, src)
	c.storeBacklog(msg, src)
	c.storeBacklog(newTestBacklogMessage(t, 2), src)

	assert.Equal(t, 2, c.backlogs[src].Size())
	assert.Equal(t, 2, c.backlogCounts[src])
	assert.Equal(t, justificationBytes, c.backlogEvidenceBytes[src])
	assert.Equal(t, justificationBytes, c.backlogTotalEvidenceBytes)
}

func TestStoreBacklogBoundsTotalJustificationBytesWithoutBlockingSmallMessages(t *testing.T) {
	c := newTestBacklogCore()
	msg := newTestBacklogRoundChange(t, 2)
	msg.Evidence = make([]byte, maxBacklogEvidenceBytesPerSender)

	for sender := 1; sender <= int(maxBacklogEvidenceBytes/maxBacklogEvidenceBytesPerSender)+1; sender++ {
		c.storeBacklog(msg, backlogSender(sender))
	}

	accepted := int(maxBacklogEvidenceBytes / maxBacklogEvidenceBytesPerSender)
	assert.Len(t, c.backlogs, accepted)
	assert.Equal(t, uint64(accepted)*maxBacklogEvidenceBytesPerSender, c.backlogTotalEvidenceBytes)
	assert.NotContains(t, c.backlogs, backlogSender(accepted+1))

	// A full attachment budget must not block ordinary consensus messages,
	// including from a sender whose own attachment budget is also full.
	c.storeBacklog(newTestBacklogMessage(t, 2), backlogSender(1))
	c.storeBacklog(newTestBacklogMessage(t, 2), backlogSender(accepted+1))
	assert.Equal(t, 2, c.backlogs[backlogSender(1)].Size())
	assert.Equal(t, 1, c.backlogs[backlogSender(accepted+1)].Size())
	assert.Equal(t, uint64(accepted)*maxBacklogEvidenceBytesPerSender, c.backlogTotalEvidenceBytes)
}

// The message budget is per sender, so a sender that has filled its own budget
// leaves every other sender's untouched.
func TestStoreBacklogMessageBudgetIsPerSender(t *testing.T) {
	c := newTestBacklogCore()
	msg := newTestBacklogMessage(t, 2)
	for range maxBacklogMessagesPerSender {
		c.storeBacklog(msg, backlogSender(1))
	}

	c.storeBacklog(msg, backlogSender(2))

	assert.Equal(t, 1, c.backlogs[backlogSender(2)].Size())
}

func TestStoreBacklogBoundsFutureSequence(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()

	c.storeBacklog(newTestBacklogMessage(t, 1+maxBacklogSequencesAhead), src)
	c.storeBacklog(newTestBacklogMessage(t, 2+maxBacklogSequencesAhead), src)
	c.storeBacklog(newTestBacklogPreprepare(t, 2+maxBacklogSequencesAhead, 0), src)

	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogCounts[src])
	assert.Empty(t, c.backlogPreprepares)
}

func TestBacklogRejectsSequenceOutsideUint64(t *testing.T) {
	c := newTestBacklogCore()
	// Put the local sequence at the top of uint64 so the window alone would
	// admit the message; only the uint64 check may reject it.
	c.current = newRoundState(
		&bft.View{Sequence: new(big.Int).SetUint64(math.MaxUint64), Round: big.NewInt(0)},
		valset.NewAddressSet(nil), common.Hash{}, nil, nil, nil, nil)
	tooLarge := new(big.Int).Add(new(big.Int).SetUint64(math.MaxUint64), big.NewInt(1))

	assert.True(t, c.isBacklogSequenceTooFar(tooLarge))
}

func TestStoreBacklogSkipsUndecodableMessage(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()

	c.storeBacklog(&bft.Message{Code: bft.MsgPrepare, Msg: []byte{0xff}}, src)
	c.storeBacklog(&bft.Message{Code: bft.MsgPreprepare, Msg: []byte{0xff}}, src)

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogCounts)
	assert.Empty(t, c.backlogPreprepares)
}

func TestProcessBacklogFreesCapacityForLaterMessages(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	msg := newTestBacklogRoundChange(t, 2)
	msg.Evidence = make([]byte, 1024)

	for range maxBacklogMessagesPerSender {
		c.storeBacklog(msg, src)
	}
	setTestBacklogView(c, 3, 0)
	c.processBacklog()

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogCounts)
	assert.Empty(t, c.backlogEvidenceBytes)
	assert.Zero(t, c.backlogTotalEvidenceBytes)

	c.storeBacklog(msg, src)
	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogCounts[src])
	assert.Equal(t, uint64(len(msg.Evidence)), c.backlogEvidenceBytes[src])
	assert.Equal(t, uint64(len(msg.Evidence)), c.backlogTotalEvidenceBytes)
}

func TestProcessBacklogKeepsExactJustificationAccountingForRemainingMessages(t *testing.T) {
	src := backlogSender(1)
	c := newTestBacklogCore()
	old := newTestBacklogRoundChange(t, 2)
	old.Evidence = make([]byte, 100)
	future := newTestBacklogRoundChange(t, 4)
	future.Evidence = make([]byte, 200)
	c.storeBacklog(old, src)
	c.storeBacklog(future, src)

	// The sequence-2 message becomes old, while sequence 4 is still future and
	// is pushed back into the queue with its accounting intact.
	setTestBacklogView(c, 3, 0)
	c.processBacklog()

	assert.Equal(t, 1, c.backlogs[src].Size())
	assert.Equal(t, 1, c.backlogCounts[src])
	assert.Equal(t, uint64(len(future.Evidence)), c.backlogEvidenceBytes[src])
	assert.Equal(t, uint64(len(future.Evidence)), c.backlogTotalEvidenceBytes)
}

func TestProcessBacklogRemovesMessageWithNilView(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	payload, err := bft.Encode(&bft.Prepare{})
	require.NoError(t, err)
	msg := &bft.Message{Code: bft.MsgPrepare, Msg: payload}

	c.backlogs[src] = prque.New()
	c.backlogs[src].Push(msg, 0)
	c.backlogCounts[src] = 1
	c.processBacklog()

	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogCounts)
	assert.Empty(t, c.backlogEvidenceBytes)
	assert.Zero(t, c.backlogTotalEvidenceBytes)
}

// A sender's slot keeps the PREPREPARE for the highest view: a higher view
// replaces it, while a delayed older one, or another proposal for the same
// view, cannot evict it.
func TestStoreBacklogKeepsHighestViewPrepreparePerSender(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	round1 := newTestBacklogPreprepare(t, 1, 1)
	round2 := newTestBacklogPreprepare(t, 1, 2)
	nextSequence := newTestBacklogPreprepare(t, 2, 0)

	c.storeBacklog(round1, src)
	c.storeBacklog(round2, src)
	assert.Same(t, round2, c.backlogPreprepares[src].msg)

	c.storeBacklog(round1, src)
	assert.Same(t, round2, c.backlogPreprepares[src].msg, "an older view must not replace a newer one")

	c.storeBacklog(newTestBacklogPreprepare(t, 1, 2), src)
	assert.Same(t, round2, c.backlogPreprepares[src].msg, "the same view must not replace the retained one")

	c.storeBacklog(nextSequence, src)
	assert.Same(t, nextSequence, c.backlogPreprepares[src].msg)
	assert.Empty(t, c.backlogs)
	assert.Empty(t, c.backlogCounts)
}

// The PREPREPARE slot is per sender, so senders flooding their own slots cannot
// take the slot of the proposer the node is waiting for.
func TestStoreBacklogPreprepareSlotIsPerSender(t *testing.T) {
	const flooders = 8
	c := newTestBacklogCore()
	for sender := 1; sender <= flooders; sender++ {
		for round := range int64(16) {
			c.storeBacklog(newTestBacklogPreprepare(t, 1, round+1), backlogSender(sender))
		}
	}
	proposer := backlogSender(flooders + 1)
	msg := newTestBacklogPreprepare(t, 1, 1)

	c.storeBacklog(msg, proposer)

	assert.Len(t, c.backlogPreprepares, flooders+1)
	assert.Same(t, msg, c.backlogPreprepares[proposer].msg)
}

// A retained PREPREPARE does not consume the sender's message budget, and a
// full message budget does not block the sender's PREPREPARE.
func TestStoreBacklogPreprepareSlotIsApartFromMessageBudget(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	msg := newTestBacklogMessage(t, 2)
	for range maxBacklogMessagesPerSender {
		c.storeBacklog(msg, src)
	}
	preprepare := newTestBacklogPreprepare(t, 2, 0)

	c.storeBacklog(preprepare, src)
	c.storeBacklog(msg, src)

	assert.Same(t, preprepare, c.backlogPreprepares[src].msg)
	assert.Equal(t, maxBacklogMessagesPerSender, c.backlogCounts[src])
}

// During a round change, checkMessage reports even the current view as future,
// so the next round's PREPREPARE is retained. It must stay retained until the
// round starts, then be delivered and released from the slot.
func TestProcessBacklogDeliversPreprepareAfterRoundChange(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	mockBackend, _, _, _ := newMockBackend(t, []common.Address{c.address, src, common.HexToAddress("0x3")}, false)
	c.backend = mockBackend
	events := mockBackend.EventMux().Subscribe(backlogEvent{})
	defer events.Unsubscribe()

	setTestBacklogView(c, 1, 1)
	c.waitingForRoundChange = true
	msg := newTestBacklogPreprepare(t, 1, 1)
	c.storeBacklog(msg, src)
	c.processBacklog()
	assert.Same(t, msg, c.backlogPreprepares[src].msg, "retained while the round change is pending")

	c.waitingForRoundChange = false
	c.processBacklog()

	assert.Empty(t, c.backlogPreprepares)
	select {
	case ev := <-events.Chan():
		assert.Same(t, msg, ev.Data.(backlogEvent).msg)
	case <-time.After(time.Second):
		t.Fatal("retained PREPREPARE was not delivered")
	}
}

// A retained PREPREPARE whose view has passed is released without delivery.
func TestProcessBacklogDropsStalePreprepare(t *testing.T) {
	src := common.HexToAddress("0x1")
	c := newTestBacklogCore()
	c.storeBacklog(newTestBacklogPreprepare(t, 1, 1), src)

	setTestBacklogView(c, 1, 2)
	c.processBacklog()

	assert.Empty(t, c.backlogPreprepares)
}
