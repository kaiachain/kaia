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
// This file is derived from quorum/consensus/istanbul/core/preprepare.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
)

func (c *core) sendPreprepare(request *istanbul.Request) {
	logger := c.logger.NewWith("state", c.state)

	header := types.SetRoundToHeader(request.Proposal.Header(), c.currentView().Round.Int64())
	request.Proposal = request.Proposal.WithSeal(header)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() {
		curView := c.currentView()
		preprepare, err := Encode(&istanbul.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&message{
			Hash: request.Proposal.ParentHash(),
			Code: msgPreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *message, src common.Address) error {
	logger := c.logger.NewWith("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare *istanbul.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		logger.Error("Failed to decode message", "code", msg.Code, "err", err)
		return errInvalidMessage
	}

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgPreprepare, preprepare.View); err != nil {
		if err == errOldMessage {
			// Get validator set for the given proposal
			councilState, getCouncilError := c.backend.GetCommitteeStateByRound(preprepare.View.Sequence.Uint64(), preprepare.View.Round.Uint64())
			if getCouncilError != nil {
				return getCouncilError
			}
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if councilState.IsProposer(src) && c.backend.HasPropsal(preprepare.Proposal.Hash(), preprepare.Proposal.Number()) {
				c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash(), preprepare.Proposal.ParentHash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
	if !c.currentCommittee.IsProposer(src) {
		logger.Warn("Ignore preprepare messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(preprepare.Proposal); err != nil {
		logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
		// if it's a future block, we will handle it again after the duration
		if err == consensus.ErrFutureBlock {
			c.stopFuturePreprepareTimer()
			c.futurePreprepareTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src:  src,
					msg:  msg,
					Hash: msg.Hash,
				})
			})
		} else {
			c.sendNextRoundChange("handlePreprepare. Proposal verification failure. Not ErrFutureBlock")
		}
		return err
	}

	// Here is about to accept the PRE-PREPARE
	if c.state == StateAcceptRequest {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			header := types.SetRoundToHeader(c.current.Preprepare.Proposal.Header(), c.currentView().Round.Int64())
			c.current.Preprepare.Proposal = c.current.Preprepare.Proposal.WithSeal(header)

			if preprepare.Proposal.Hash() == c.current.GetLockedHash() {
				logger.Warn("Received preprepare message of the hash locked proposal and change state to prepared")
				// Broadcast COMMIT and enters Prepared state directly
				c.acceptPreprepare(preprepare)
				c.setState(StatePrepared)
				c.sendCommit()

				if vrank != nil {
					vrank.Log()
				}
				vrank = NewVrank(*c.currentView(), c.currentCommittee.Committee().List())
			} else {
				// Send round change
				c.sendNextRoundChange("handlePreprepare. HashLocked, but received hash is different from locked hash")
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptPreprepare(preprepare)
			c.setState(StatePreprepared)
			c.sendPrepare()

			if vrank != nil {
				vrank.Log()
			}
			vrank = NewVrank(*c.currentView(), c.currentCommittee.Committee().List())
		}
	}

	return nil
}

func (c *core) acceptPreprepare(preprepare *istanbul.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}
