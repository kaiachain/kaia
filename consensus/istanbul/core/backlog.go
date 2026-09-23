// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
//
// This file is derived from quorum/consensus/istanbul/core/backlog.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/prque"
	"github.com/kaiachain/kaia/consensus/bft"
)

// msgPriority is defined for calculating processing priority to speedup consensus
// bft.MsgPreprepare > bft.MsgCommit > bft.MsgPrepare
var msgPriority = map[uint64]int{
	bft.MsgPreprepare: 1,
	bft.MsgCommit:     2,
	bft.MsgPrepare:    3,
}

const (
	// The backlog retains future messages per sender only, with no budget shared
	// between senders: a shared budget lets a few senders fill it ahead of an
	// honest one, and during a round change the message dropped that way is the
	// next round's PREPREPARE. The per-sender limits bound the total on their
	// own, since handleMsg admits only the qualified council, the council is
	// bounded by ABv2 (DefaultMaxValActivePausedCount), and every message is
	// bounded by checkMessageSize or the block size.

	// Keep only a small future-sequence window so far-future messages cannot
	// occupy a sender's budget until the node catches up. A node further behind
	// than this window catches up through block synchronization rather than
	// through retained consensus messages.
	maxBacklogSequencesAhead = 8

	// PREPARE, COMMIT and ROUND CHANGE are bounded by checkMessageSize, so a
	// count limit bounds the memory they occupy. A PREPREPARE carries an entire
	// block and takes the sender's single PREPREPARE slot instead.
	maxBacklogMessagesPerSender = 128
)

// checkMessage checks the message state
// return bft.ErrInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *bft.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return bft.ErrInvalidMessage
	}

	// The round is consumed as a uint64 (View.Round.Uint64()); reject out-of-range values.
	if !view.Round.IsUint64() {
		return bft.ErrInvalidMessage
	}

	if msgCode == bft.MsgRoundChange {
		if view.Sequence.Cmp(c.currentView().Sequence) > 0 {
			return errFutureMessage
		} else if view.Cmp(c.currentView()) < 0 {
			return errOldMessage
		}
		return nil
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}

	if c.waitingForRoundChange {
		return errFutureMessage
	}

	// StateAcceptRequest only accepts bft.MsgPreprepare
	// other messages are future messages
	if c.state == StateAcceptRequest {
		if msgCode > bft.MsgPreprepare {
			return errFutureMessage
		}
		return nil
	}

	// For states(StatePreprepared, StatePrepared, StateCommitted),
	// can accept all message types if processing with same view
	return nil
}

// backlogPreprepare is the one PREPREPARE retained per sender. The view and
// parent hash are decoded once at retention, so that neither a replacement nor
// processBacklog decodes the block again.
type backlogPreprepare struct {
	msg      *bft.Message
	view     *bft.View
	prevHash common.Hash
}

// storeBacklog retains a future message for processBacklog. PREPARE, COMMIT
// and ROUND CHANGE queue per sender up to maxBacklogMessagesPerSender. A
// PREPREPARE takes the sender's single slot, replacing a retained one only when
// its view is higher: a sender proposes at most once per round, and only its
// proposal for the highest view can still be handled, while a delayed older one
// must not evict it. The slot is never shared, so no sender can take the slot
// of the proposer the node is waiting for.
func (c *core) storeBacklog(msg *bft.Message, src common.Address) {
	logger := c.logger.NewWith("from", src, "state", c.state)

	if src == c.Address() {
		logger.Warn("Backlog from self")
		return
	}

	logger.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	view, prevHash := backlogMessageView(msg)
	if view == nil {
		logger.Trace("Discarding future message: cannot decode view")
		return
	}
	if c.isBacklogSequenceTooFar(view.Sequence) {
		logger.Trace("Discarding future message: sequence is too far ahead", "sequence", view.Sequence)
		return
	}

	if msg.Code == bft.MsgPreprepare {
		if held, ok := c.backlogPreprepares[src]; ok && view.Cmp(held.view) <= 0 {
			logger.Trace("Discarding future PREPREPARE: not newer than the retained one", "view", view, "retained", held.view)
			return
		}
		c.backlogPreprepares[src] = backlogPreprepare{msg: msg, view: view, prevHash: prevHash}
		return
	}

	if c.backlogCounts[src] >= maxBacklogMessagesPerSender {
		// A full backlog is expected under load; avoid a warning for every
		// dropped message.
		logger.Trace("Discarding future message: sender backlog limit reached")
		return
	}
	backlog := c.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
		c.backlogs[src] = backlog
	}
	// toPriority truncates the sequence, so it runs only after
	// isBacklogSequenceTooFar has rejected sequences that do not fit in uint64.
	backlog.Push(msg, toPriority(msg.Code, view))
	c.backlogCounts[src]++
}

func (c *core) isBacklogSequenceTooFar(sequence *big.Int) bool {
	if !sequence.IsUint64() {
		return true
	}
	maxSequence := new(big.Int).Add(c.currentView().Sequence, big.NewInt(maxBacklogSequencesAhead))
	return sequence.Cmp(maxSequence) > 0
}

// retainedMessageBytes reports the memory a retained message occupies, so that
// every size limit measures a message the same way.
func retainedMessageBytes(msg *bft.Message) uint64 {
	return uint64(len(msg.Msg)) + uint64(len(msg.Signature)) + uint64(len(msg.CommittedSeal))
}

// removeBacklogMessage releases one queued message of a sender while backlogsMu
// is held. It drops the sender's counter at zero, so no explicit cleanup is
// needed when the sender's queue becomes empty.
func (c *core) removeBacklogMessage(src common.Address) {
	if c.backlogCounts[src] <= 1 {
		delete(c.backlogCounts, src)
		return
	}
	c.backlogCounts[src]--
}

// backlogMessageView decodes the view a message belongs to and the parent hash
// to post it under. The view is nil when the message cannot be retained.
func backlogMessageView(msg *bft.Message) (view *bft.View, prevHash common.Hash) {
	if msg.Code == bft.MsgPreprepare {
		var p *bft.Preprepare
		if err := msg.Decode(&p); err != nil || p == nil || p.Proposal == nil {
			return nil, common.Hash{}
		}
		view, prevHash = p.View, p.Proposal.ParentHash()
	} else {
		var sub *bft.Subject
		if err := msg.Decode(&sub); err != nil || sub == nil {
			return nil, common.Hash{}
		}
		view, prevHash = sub.View, sub.PrevHash
	}
	if view == nil || view.Sequence == nil || view.Round == nil {
		return nil, common.Hash{}
	}
	return view, prevHash
}

// postBacklogMessage posts a retained message once its view is current, or
// discards it once its view has passed. It reports whether the message is
// still in the future and must stay retained.
func (c *core) postBacklogMessage(src common.Address, msg *bft.Message, view *bft.View, prevHash common.Hash) (future bool) {
	logger := c.logger.NewWith("from", src, "state", c.state)

	if view == nil {
		logger.Debug("Nil view", "msg", msg)
		return false
	}
	err := c.checkMessage(msg.Code, view)
	if err == errFutureMessage {
		return true
	}
	if err != nil {
		logger.Trace("Skip the backlog event", "msg", msg, "err", err)
		return false
	}
	logger.Trace("Post backlog event", "msg", msg)

	go c.sendEvent(backlogEvent{
		src:  src,
		msg:  msg,
		Hash: prevHash,
	})
	return false
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, held := range c.backlogPreprepares {
		if !c.postBacklogMessage(src, held.msg, held.view, held.prevHash) {
			delete(c.backlogPreprepares, src)
		}
	}

	for src, backlog := range c.backlogs {
		// Stop at the first future message: the queue is ordered by view, so
		// everything behind it is in the future as well.
		for !backlog.Empty() {
			m, prio := backlog.Pop()
			msg := m.(*bft.Message)
			view, prevHash := backlogMessageView(msg)
			if c.postBacklogMessage(src, msg, view, prevHash) {
				backlog.Push(msg, prio)
				break
			}
			c.removeBacklogMessage(src)
		}

		// Do not retain prque's backing storage after all messages from this
		// sender have been processed or discarded.
		if backlog.Empty() {
			delete(c.backlogs, src)
		}
	}
}

func toPriority(msgCode uint64, view *bft.View) int64 {
	if msgCode == bft.MsgRoundChange {
		// For bft.MsgRoundChange, set the message priority based on its sequence
		return -int64(view.Sequence.Uint64() * 1000)
	}
	// FIXME: round will be reset as 0 while new sequence
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Sequence limits the range of round is from 0 to 99
	return -int64(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
}
