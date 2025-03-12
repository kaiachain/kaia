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
// This file is derived from quorum/consensus/istanbul/core/request.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"github.com/kaiachain/kaia/consensus/istanbul"
)

func (c *core) handleRequest(request *istanbul.Request) error {
	logger := c.logger.NewWith("state", c.state, "seq", c.current.sequence)

	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return err
		}
		logger.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.Hash())
		return err
	}

	logger.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.current.pendingRequest = request
	if c.state == StateAcceptRequest {
		c.sendPreprepare(request)
	}
	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *istanbul.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	cmp := c.current.sequence.Cmp(request.Proposal.Number())
	if cmp > 0 { // request.Proposal.Number < c.current.sequence
		return errOldMessage
	} else if cmp < 0 { // request.Proposer.Number > c.current.Sequence
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *istanbul.Request) {
	logger := c.logger.NewWith("state", c.state)

	logger.Trace("Store future request", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, -request.Proposal.Number().Int64())
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*istanbul.Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())

		go c.sendEvent(istanbul.RequestEvent{
			Proposal: r.Proposal,
		})
	}
}
