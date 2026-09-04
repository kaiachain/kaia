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
	// The limits below bound both the number and payload size of retained future
	// messages. A PREPREPARE carries an entire block, so a count-only limit is
	// not sufficient to bound backlog memory.
	maxBacklogMessagesPerSender     = 128
	maxBacklogPayloadBytesPerSender = 16 * 1024 * 1024
	// Global limits allow up to eight senders to reach their per-sender limits.
	maxBacklogMessages     = 1024
	maxBacklogPayloadBytes = 128 * 1024 * 1024
	// Keep only a small future-sequence window so far-future messages cannot
	// occupy the global backlog budget until the node catches up. A node further
	// behind than this window catches up through block synchronization rather
	// than through retained consensus messages.
	maxBacklogSequencesAhead = 8
)

// checkMessage checks the message state
// return bft.ErrInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *bft.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return bft.ErrInvalidMessage
	}

	if msgCode == bft.MsgRoundChange {
		// Round-change buckets are keyed by uint64 round in roundChangeSet.
		// Reject out-of-range round values early to avoid truncation collisions.
		if !view.Round.IsUint64() {
			return bft.ErrInvalidMessage
		}
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

func (c *core) storeBacklog(msg *bft.Message, src common.Address) {
	logger := c.logger.NewWith("from", src, "state", c.state)

	if src == c.Address() {
		logger.Warn("Backlog from self")
		return
	}

	logger.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	backlog := c.backlogs[src]
	messageBytes := backlogMessageBytes(msg)
	if c.backlogLimitReached(src, backlog, messageBytes) {
		// A full backlog is expected under load; avoid a warning for every
		// dropped message.
		logger.Trace("Discarding future message: backlog limit reached")
		return
	}

	view, err := msg.GetView()
	if err != nil || view == nil || view.Sequence == nil || view.Round == nil {
		logger.Trace("Discarding future message: cannot decode view", "err", err)
		return
	}
	if c.isBacklogSequenceTooFar(view.Sequence) {
		logger.Trace("Discarding future message: sequence is too far ahead", "sequence", view.Sequence)
		return
	}
	if backlog == nil {
		backlog = prque.New()
		c.backlogs[src] = backlog
	}
	// toPriority truncates the sequence, so it runs only after
	// isBacklogSequenceTooFar has rejected sequences that do not fit in uint64.
	backlog.Push(msg, toPriority(msg.Code, view))
	c.backlogMessageCount++
	c.backlogSenderBytes[src] += messageBytes
	c.backlogPayloadBytes += messageBytes
}

func (c *core) isBacklogSequenceTooFar(sequence *big.Int) bool {
	if !sequence.IsUint64() {
		return true
	}
	maxSequence := new(big.Int).Add(c.currentView().Sequence, big.NewInt(maxBacklogSequencesAhead))
	return sequence.Cmp(maxSequence) > 0
}

func backlogMessageBytes(msg *bft.Message) uint64 {
	return uint64(len(msg.Msg)) + uint64(len(msg.Signature)) + uint64(len(msg.CommittedSeal))
}

func exceedsBacklogLimit(used, additional, limit uint64) bool {
	return additional > limit || used > limit-additional
}

// backlogLimitReached reports whether retaining a message would exceed a
// per-sender or global backlog limit. backlogsMu must be held by the caller.
func (c *core) backlogLimitReached(src common.Address, backlog *prque.Prque, messageBytes uint64) bool {
	if backlog != nil && backlog.Size() >= maxBacklogMessagesPerSender {
		return true
	}
	if c.backlogMessageCount >= maxBacklogMessages {
		return true
	}
	if exceedsBacklogLimit(c.backlogSenderBytes[src], messageBytes, maxBacklogPayloadBytesPerSender) {
		return true
	}
	return exceedsBacklogLimit(c.backlogPayloadBytes, messageBytes, maxBacklogPayloadBytes)
}

// removeBacklogMessage updates accounting while backlogsMu is held. It also
// drops the sender's byte entry once its last retained message is removed, so
// no explicit cleanup is needed when the sender's queue becomes empty.
func (c *core) removeBacklogMessage(src common.Address, msg *bft.Message) {
	messageBytes := backlogMessageBytes(msg)
	c.backlogMessageCount--
	c.backlogPayloadBytes -= messageBytes
	if senderBytes := c.backlogSenderBytes[src]; senderBytes > messageBytes {
		c.backlogSenderBytes[src] = senderBytes - messageBytes
	} else {
		delete(c.backlogSenderBytes, src)
	}
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}

		logger := c.logger.NewWith("from", src, "state", c.state)

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !backlog.Empty() {
			m, prio := backlog.Pop()
			msg := m.(*bft.Message)
			var view *bft.View
			var prevHash common.Hash
			switch msg.Code {
			case bft.MsgPreprepare:
				var m *bft.Preprepare
				if err := msg.Decode(&m); err == nil && m != nil {
					view = m.View
					if m.Proposal != nil {
						prevHash = m.Proposal.ParentHash()
					}
				}
			default:
				var sub *bft.Subject
				if err := msg.Decode(&sub); err == nil && sub != nil {
					view = sub.View
					prevHash = sub.PrevHash
				}
			}
			if view == nil {
				logger.Debug("Nil view", "msg", msg)
				c.removeBacklogMessage(src, msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, view)
			if err != nil {
				if err == errFutureMessage {
					logger.Trace("Stop processing backlog", "msg", msg)
					backlog.Push(msg, prio)
					break
				}
				logger.Trace("Skip the backlog event", "msg", msg, "err", err)
				c.removeBacklogMessage(src, msg)
				continue
			}
			logger.Trace("Post backlog event", "msg", msg)
			c.removeBacklogMessage(src, msg)

			go c.sendEvent(backlogEvent{
				src:  src,
				msg:  msg,
				Hash: prevHash,
			})
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
