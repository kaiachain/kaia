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
	"math/big"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
)

func (c *core) sendPreprepare(request *bft.Request) {
	logger := c.logger.NewWith("state", c.state)

	header := request.Proposal.Header()
	c.backend.Sealer().WriteRound(header, c.currentView().Round.Int64())
	request.Proposal = request.Proposal.WithSeal(header)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() {
		curView := c.currentView()
		message := &bft.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		}
		if c.backend.IsPermissionlessAt(curView.Sequence.Uint64()) && curView.Round.Sign() > 0 {
			message.RoundChangeCertificate = c.roundChangeCertificate
		}
		preprepare, err := bft.Encode(message)
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&bft.Message{
			Hash: request.Proposal.ParentHash(),
			Code: bft.MsgPreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *bft.Message, src common.Address) error {
	logger := c.logger.NewWith("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare *bft.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		logger.Error("Failed to decode message", "code", msg.Code, "err", err)
		return bft.ErrInvalidMessage
	}

	// Proposal number must equal the view sequence, else it can be aliased to another height.
	if !proposalNumberMatchesView(preprepare) {
		logger.Warn("Ignore preprepare whose proposal number does not match view sequence")
		return bft.ErrInvalidMessage
	}

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(bft.MsgPreprepare, preprepare.View); err != nil {
		if err == errOldMessage {
			// This PRE-PREPARE targets an already-finalized height. Reply with a COMMIT
			// to help the sender finish that block, only if:
			// 1. The sender is the scheduled proposer of the given (sequence, round).
			// 2. This node has the given block in its chain.
			// 3. Post-permissionless: the given round equals the round this node stored.
			//    The committed seal signs (hash, round), so answering another round would
			//    sign a commit at a round this node never committed at; callers could
			//    collect such seals into a quorum certificate for a round that never
			//    reached quorum.
			proposer, getProposerErr := c.valsetModule.GetProposer(preprepare.View.Sequence.Uint64(), preprepare.View.Round.Uint64())
			if getProposerErr != nil {
				return getProposerErr
			}
			storedRound, hasProposal := c.backend.ProposalRound(preprepare.Proposal.Hash(), preprepare.Proposal.Number())
			roundMatches := !c.backend.IsPermissionlessAt(preprepare.View.Sequence.Uint64()) ||
				uint64(storedRound) == preprepare.View.Round.Uint64()
			if proposer == src && hasProposal && roundMatches {
				c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash(), preprepare.Proposal.ParentHash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
	if c.current.proposer != src {
		logger.Warn("Ignore preprepare messages from non-proposer")
		return errNotFromProposer
	}

	var highestPrepared *bft.PreparedCertificate
	if c.backend.IsPermissionlessAt(preprepare.View.Sequence.Uint64()) {
		if preprepare.View.Round.Sign() == 0 {
			if len(preprepare.RoundChangeCertificate) != 0 {
				return bft.ErrInvalidMessage
			}
		} else {
			highestPrepared, err = c.verifyRoundChangeCertificate(preprepare.RoundChangeCertificate, preprepare.View)
			if err != nil {
				logger.Warn("Invalid round-change certificate in PRE-PREPARE", "err", err)
				return bft.ErrInvalidMessage
			}
			if highestPrepared != nil && highestPrepared.Proposal.Hash() != preprepare.Proposal.Hash() {
				logger.Warn("PRE-PREPARE does not carry the highest prepared proposal",
					"expected", highestPrepared.Proposal.Hash(), "got", preprepare.Proposal.Hash())
				return bft.ErrInvalidMessage
			}
		}
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
		if highestPrepared != nil {
			lockedRound := c.current.LockedRound()
			if lockedRound == nil || highestPrepared.View.Round.Cmp(lockedRound) > 0 {
				c.current.AdoptPreparedCertificate(highestPrepared)
			}
		}
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			// Legacy IBFT re-seals the locally retained proposal with the new
			// round before comparing it. Keep that wire-compatible behavior on
			// the legacy path. A Permissionless node may instead have adopted a
			// certificate without a local Preprepare, so it must not dereference
			// or mutate that optional local proposal.
			if !c.backend.IsPermissionlessAt(preprepare.View.Sequence.Uint64()) && c.current.Preprepare != nil {
				header := c.current.Preprepare.Proposal.Header()
				c.backend.Sealer().WriteRound(header, c.currentView().Round.Int64())
				c.current.Preprepare.Proposal = c.current.Preprepare.Proposal.WithSeal(header)
			}
			if preprepare.Proposal.Hash() == c.current.GetLockedHash() {
				if !c.backend.IsPermissionlessAt(preprepare.View.Sequence.Uint64()) {
					logger.Warn("Received preprepare message of the hash locked proposal and change state to prepared")
					// Preserve the legacy shortcut before Permissionless activation.
					c.acceptPreprepare(preprepare)
					c.postPrepreparedEvent(preprepare)
					c.setState(StatePrepared)
					c.sendCommit()
				} else {
					// A carried certificate proves a prior round, not this one. Re-establish
					// a quorum at the current view before sending COMMIT, so every honest
					// committer can advertise a current-round certificate in a later RC.
					c.acceptPreprepare(preprepare)
					c.postPrepreparedEvent(preprepare)
					c.setState(StatePreprepared)
					c.sendPrepare()
				}
			} else {
				// Send round change
				c.sendNextRoundChange("handlePreprepare. HashLocked, but received hash is different from locked hash")
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptPreprepare(preprepare)
			c.postPrepreparedEvent(preprepare)
			c.setState(StatePreprepared)
			c.sendPrepare()
		}
	}

	return nil
}

// proposalNumberMatchesView reports whether the proposal's block number equals the view sequence.
func proposalNumberMatchesView(pp *bft.Preprepare) bool {
	if pp == nil || pp.Proposal == nil || pp.Proposal.Number() == nil ||
		pp.View == nil || pp.View.Sequence == nil {
		return false
	}
	return pp.Proposal.Number().Cmp(pp.View.Sequence) == 0
}

func (c *core) acceptPreprepare(preprepare *bft.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}

// postPrepreparedEvent notifies subscribers (VRank) that this node accepted the
// PRE-PREPARE for the given (block, view). The view is deep-copied so later
// round mutations don't alias the event payload.
func (c *core) postPrepreparedEvent(preprepare *bft.Preprepare) {
	block, ok := preprepare.Proposal.(*types.Block)
	if !ok {
		c.logger.Warn("Failed to post preprepared event due to unexpected proposal type")
		return
	}
	view := &bft.View{
		Round:    new(big.Int).Set(preprepare.View.Round),
		Sequence: new(big.Int).Set(preprepare.View.Sequence),
	}
	c.sendEvent(istanbul.PrepreparedEvent{Block: block, View: view})
}
