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
	"errors"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/prque"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// PBFT state constants. Values and semantics match istanbul/core.
const (
	stateAcceptRequest uint64 = iota
	statePreprepared
	statePrepared
	stateCommitted
)

var stateNames = map[uint64]string{
	stateAcceptRequest: "AcceptRequest",
	statePreprepared:   "Preprepared",
	statePrepared:      "Prepared",
	stateCommitted:     "Committed",
}

var (
	errFutureMessage        = errors.New("future message")
	errOldMessage           = errors.New("old message")
	errNotFromProposer      = errors.New("message does not come from proposer")
	errNotFromCommittee     = errors.New("message does not come from committee")
	errInconsistentSubject  = errors.New("inconsistent subjects")
	errInvalidCommittedSeal = errors.New("invalid committed seal")
	errRoundChangeTooFar    = errors.New("round change is too far in the future")
	errRoundChangeLimit     = errors.New("round change message limit reached")
	errMessageTooLarge      = errors.New("message is too large")
	errUnauthorizedAddress  = errors.New("unauthorized address")
	errIgnored              = errors.New("ignored")
)

// machine implements the single-runloop BFT state machine.
// All state mutations happen inside runLoop; external goroutines communicate
// exclusively via the event mux channels.
type machine struct {
	b *backend

	// Subscriptions
	events     *event.TypeMuxSubscription
	timeoutSub *event.TypeMuxSubscription
	chainHead  *event.TypeMuxSubscription
	wg         sync.WaitGroup

	// Current round state — only accessed by runLoop
	sequence *big.Int
	round    *big.Int
	state    uint64

	// Proposal and message tracking
	preprepare *bft.Preprepare
	prepares   *messageSet
	commits    *messageSet
	lockedHash common.Hash

	// Committee state for current round
	qualified            *valset.AddressSet
	committee            *valset.AddressSet
	proposer             common.Address
	committeeSize        uint64
	requiredMessageCount int
	f                    int

	pendingRequest        *bft.Request
	waitingForRoundChange bool

	// Backlog
	backlogs           map[common.Address]*prque.Prque
	backlogCounts      map[common.Address]int               // queued PREPARE, COMMIT and ROUND CHANGE per sender
	backlogPreprepares map[common.Address]backlogPreprepare // the one retained PREPREPARE per sender
	backlogsMu         sync.Mutex

	// Round change
	roundChangeSets map[uint64]*messageSet
	roundChangesMu  sync.Mutex

	// Timer
	roundChangeTimer atomic.Value // *time.Timer

	// Future preprepare retry timer
	futurePreprepareTimer *time.Timer

	// Consensus timing
	consensusTimestamp time.Time

	metrics *bft.CoreMetrics
}

func newMachine(b *backend) *machine {
	return &machine{
		b:                  b,
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogCounts:      make(map[common.Address]int),
		backlogPreprepares: make(map[common.Address]backlogPreprepare),
		roundChangeSets:    make(map[uint64]*messageSet),
		metrics:            bft.NewCoreMetrics("consensus/istanbul/core"),
	}
}

// ---------------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------------

func (m *machine) start() error {
	m.events = m.b.eventMux.Subscribe(
		requestEvent{},
		messageEvent{},
		backlogEvent{},
	)
	m.timeoutSub = m.b.eventMux.Subscribe(timeoutEvent{})
	m.chainHead = m.b.eventMux.Subscribe(chainHeadEvent{})

	m.startNewRound(common.Big0)

	m.wg.Add(1)
	go m.runLoop()
	return nil
}

func (m *machine) stop() {
	m.stopTimer()
	m.events.Unsubscribe()
	m.timeoutSub.Unsubscribe()
	m.chainHead.Unsubscribe()
	m.wg.Wait()
}

// ---------------------------------------------------------------------------
// Event loop — single goroutine owns all state
// ---------------------------------------------------------------------------

func (m *machine) runLoop() {
	defer func() {
		m.sequence = nil
		m.wg.Done()
	}()

	for {
		select {
		case ev, ok := <-m.events.Chan():
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case requestEvent:
				if err := m.handleRequest(&bft.Request{Proposal: e.Proposal}); err == errFutureMessage {
					m.storePendingRequest(&bft.Request{Proposal: e.Proposal})
				}
			case messageEvent:
				if err := m.handleMsg(e.Payload); err == nil {
					m.b.gossipSubPeer(e.Hash, e.Payload)
				}
			case backlogEvent:
				if m.qualified != nil && !m.qualified.Contains(e.src) {
					continue
				}
				if err := m.handleCheckedMsg(e.msg, e.src); err == nil {
					p, err := e.msg.Payload()
					if err != nil {
						continue
					}
					m.b.gossipSubPeer(e.Hash, p)
				}
			}

		case ev, ok := <-m.timeoutSub.Chan():
			if !ok || ev.Data == nil {
				return
			}
			data, ok := ev.Data.(timeoutEvent)
			if !ok || data.nextView == nil {
				return
			}
			m.handleTimeout(data.nextView)

		case ev, ok := <-m.chainHead.Chan():
			if !ok {
				return
			}
			if _, ok := ev.Data.(chainHeadEvent); ok {
				m.handleChainHead()
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Message decoding and dispatch
// ---------------------------------------------------------------------------

func (m *machine) handleMsg(payload []byte) error {
	msg := new(bft.Message)
	if err := msg.FromPayload(payload, m.checkValidatorSignature); err != nil {
		if m.b.nodetype == common.CONSENSUSNODE && err != errUnauthorizedAddress {
			logger.Error("Failed to decode message from payload", "err", err)
		}
		return err
	}
	if m.qualified == nil || !m.qualified.Contains(msg.Address) {
		return errUnauthorizedAddress
	}
	return m.handleCheckedMsg(msg, msg.Address)
}

func (m *machine) handleCheckedMsg(msg *bft.Message, src common.Address) error {
	if err := checkMessageSize(msg); err != nil {
		return err
	}

	testBacklog := func(err error) error {
		if err == errFutureMessage {
			m.storeBacklog(msg, src)
			lastProposal, _ := m.b.lastProposal()
			if lastProposal != nil && lastProposal.Number().Cmp(m.sequence) >= 0 {
				m.startNewRound(common.Big0)
			}
		}
		return err
	}

	switch msg.Code {
	case bft.MsgPreprepare:
		return testBacklog(m.handlePreprepare(msg, src))
	case bft.MsgPrepare:
		return testBacklog(m.handlePrepare(msg, src))
	case bft.MsgCommit:
		return testBacklog(m.handleCommit(msg, src))
	case bft.MsgRoundChange:
		return testBacklog(m.handleRoundChange(msg, src))
	default:
		return bft.ErrInvalidMessage
	}
}

// ---------------------------------------------------------------------------
// PBFT handlers
// ---------------------------------------------------------------------------

func (m *machine) handleRequest(request *bft.Request) error {
	if request == nil || request.Proposal == nil {
		return bft.ErrInvalidMessage
	}
	cmp := m.sequence.Cmp(request.Proposal.Number())
	if cmp > 0 {
		return errOldMessage
	}
	if cmp < 0 {
		return errFutureMessage
	}
	m.pendingRequest = request
	if m.state == stateAcceptRequest {
		m.sendPreprepare(request)
	}
	return nil
}

func (m *machine) handlePreprepare(msg *bft.Message, src common.Address) error {
	var pp *bft.Preprepare
	if err := msg.Decode(&pp); err != nil {
		return bft.ErrInvalidMessage
	}

	// The proposal number must equal the view sequence, else the proposal
	// could be aliased to another height.
	if !proposalNumberMatchesView(pp) {
		return bft.ErrInvalidMessage
	}

	if err := m.checkMessage(bft.MsgPreprepare, pp.View); err != nil {
		if err == errOldMessage {
			// This PRE-PREPARE targets an already-finalized height. Reply with a
			// COMMIT to help the sender finish that block, only if:
			// 1. The sender is the scheduled proposer of the given (sequence, round).
			// 2. This node has the given block in its chain.
			// 3. Post-permissionless: the given round equals the round this node
			//    stored. The committed seal signs (hash, round), so answering
			//    another round would sign a commit at a round this node never
			//    committed at; callers could collect such seals into a quorum
			//    certificate for a round that never reached quorum.
			proposer, getErr := m.b.valsetModule.GetProposer(pp.View.Sequence.Uint64(), pp.View.Round.Uint64())
			if getErr != nil {
				return getErr
			}
			storedRound, hasProposal := m.b.proposalRound(pp.Proposal.Hash(), pp.Proposal.Number())
			roundMatches := !m.b.isPermissionlessAt(pp.View.Sequence.Uint64()) ||
				uint64(storedRound) == pp.View.Round.Uint64()
			if proposer == src && hasProposal && roundMatches {
				m.sendCommitForOldBlock(pp.View, pp.Proposal.Hash(), pp.Proposal.ParentHash())
				return nil
			}
		}
		return err
	}

	if m.proposer != src {
		return errNotFromProposer
	}

	if duration, err := m.b.verify(pp.Proposal); err != nil {
		if err == consensus.ErrFutureBlock {
			m.stopFuturePreprepareTimer()
			m.futurePreprepareTimer = time.AfterFunc(duration, func() {
				m.b.eventMux.Post(backlogEvent{src: src, msg: msg, Hash: msg.Hash})
			})
		} else {
			m.sendNextRoundChange("handlePreprepare: verification failure")
		}
		return err
	}

	if m.state == stateAcceptRequest {
		if m.isHashLocked() {
			header := m.preprepare.Proposal.Header()
			m.b.sealer.WriteRound(header, m.currentView().Round.Int64())
			m.preprepare.Proposal = m.preprepare.Proposal.WithSeal(header)

			if pp.Proposal.Hash() == m.lockedHash {
				logger.Info("Received preprepare matching hash-locked proposal, moving to Prepared",
					"seq", m.sequence, "round", m.round, "hash", pp.Proposal.Hash())
				m.acceptPreprepare(pp)
				m.setState(statePrepared)
				m.sendCommit()
			} else {
				logger.Warn("[RC] Hash locked but received different proposal hash",
					"seq", m.sequence, "round", m.round,
					"locked", m.lockedHash, "received", pp.Proposal.Hash())
				m.sendNextRoundChange("handlePreprepare: hash locked but different proposal")
			}
		} else {
			logger.Debug("Accepted preprepare, moving to Preprepared",
				"seq", m.sequence, "round", m.round, "from", src, "hash", pp.Proposal.Hash())
			// Speculative execution (kaiabft-specific).
			if !m.isProposer() {
				m.b.startSpeculativeExecution(pp.Proposal)
			}

			// Accept preprepare and move to Preprepared state
			m.acceptPreprepare(pp)
			m.setState(statePreprepared)
			m.sendPrepare()
		}
	}

	return nil
}

func (m *machine) handlePrepare(msg *bft.Message, src common.Address) error {
	var prepare *bft.Subject
	if err := msg.Decode(&prepare); err != nil {
		return bft.ErrInvalidMessage
	}
	if err := m.checkMessage(bft.MsgPrepare, prepare.View); err != nil {
		return err
	}
	if err := m.verifySubject(prepare); err != nil {
		return err
	}
	if !m.committee.Contains(src) {
		return errNotFromCommittee
	}

	m.prepares.Add(msg)

	if m.state < statePrepared {
		if m.isHashLocked() && prepare.Digest == m.lockedHash {
			logger.Info("Received prepare matching hash-locked proposal, moving to Prepared",
				"seq", m.sequence, "round", m.round, "from", src)
			m.setState(statePrepared)
			m.sendCommit()
		} else if m.getPrepareOrCommitSize() >= m.requiredMessageCount {
			logger.Info("Received quorum of prepare/commit messages, moving to Prepared",
				"seq", m.sequence, "round", m.round,
				"prepares", m.prepares.Size(), "commits", m.commits.Size(),
				"required", m.requiredMessageCount)
			m.lockHash()
			m.setState(statePrepared)
			m.sendCommit()
		}
	}
	return nil
}

func (m *machine) handleCommit(msg *bft.Message, src common.Address) error {
	var commit *bft.Subject
	if err := msg.Decode(&commit); err != nil {
		return bft.ErrInvalidMessage
	}
	if err := m.checkMessage(bft.MsgCommit, commit.View); err != nil {
		return err
	}
	if err := m.verifySubject(commit); err != nil {
		return err
	}
	if !m.committee.Contains(src) {
		return errNotFromCommittee
	}

	// Verify msg.CommittedSeal is the sender's signature over the proposal's
	// committed-seal preimage. Without this, an arbitrary seal would be copied
	// verbatim into the sealed block. commit.Digest is the proposal hash, already
	// validated by verifySubject above.
	committer, err := bft.GetSignatureAddress(m.b.committedSealPreimage(commit.View, commit.Digest), msg.CommittedSeal)
	if err != nil || committer != src {
		logger.Warn("Invalid committed seal in commit message", "sender", src, "recovered", committer, "err", err)
		return errInvalidCommittedSeal
	}

	m.commits.Add(msg)

	// Both PREPARE and COMMIT messages count for the prepared quorum.
	if m.state < statePrepared {
		if m.isHashLocked() && commit.Digest == m.lockedHash {
			logger.Info("Received commit matching hash-locked proposal, moving to Prepared",
				"seq", m.sequence, "round", m.round, "from", src)
			m.setState(statePrepared)
			m.sendCommit()
		} else if m.getPrepareOrCommitSize() >= m.requiredMessageCount {
			logger.Info("Received quorum of prepare/commit messages (via commit), moving to Prepared",
				"seq", m.sequence, "round", m.round,
				"prepares", m.prepares.Size(), "commits", m.commits.Size(),
				"required", m.requiredMessageCount)
			m.lockHash()
			m.setState(statePrepared)
			m.sendCommit()
		}
	}

	if m.state < stateCommitted && m.commits.Size() >= m.requiredMessageCount {
		logger.Info("Received quorum of commit messages, committing",
			"seq", m.sequence, "round", m.round,
			"commits", m.commits.Size(), "required", m.requiredMessageCount)
		m.lockHash()
		m.doCommit()
	}
	return nil
}

func (m *machine) handleRoundChange(msg *bft.Message, src common.Address) error {
	var rc *bft.Subject
	if err := msg.Decode(&rc); err != nil {
		return bft.ErrInvalidMessage
	}
	if err := m.checkMessage(bft.MsgRoundChange, rc.View); err != nil {
		return err
	}
	// Restrict round-change admission to committee members, mirroring
	// handleCommit. The thresholds (requiredMessageCount, f) are
	// committee-sized, so gating on the committee keeps admission consistent
	// with the thresholds.
	if !m.committee.Contains(src) {
		return errNotFromCommittee
	}

	cv := m.currentView()
	num, err := m.addRoundChange(rc.View.Round, msg)
	if err != nil {
		// Rejections are expected under load (round window, per-round limit),
		// so they are not logged above Trace, as in istanbul.
		logger.Trace("Discarding ROUND CHANGE", "from", msg.Address, "round", rc.View.Round, "err", err)
		return err
	}

	numStartNewRound := m.requiredMessageCount
	numCatchUp := m.f + 1

	logger.Debug("[RC] Round change received",
		"from", msg.Address, "rcRound", rc.View.Round, "currentRound", cv.Round,
		"seq", cv.Sequence, "numRC", num, "need", numStartNewRound,
		"waiting", m.waitingForRoundChange, "state", stateNames[m.state])

	if num == numStartNewRound && (m.waitingForRoundChange || cv.Round.Cmp(rc.View.Round) < 0) {
		logger.Warn("[RC] Received 2f+1 round change messages, starting new round",
			"seq", cv.Sequence, "currentRound", cv.Round, "newRound", rc.View.Round,
			"proposer", m.proposer, "prepares", m.prepares.Size(), "commits", m.commits.Size())
		m.startNewRound(rc.View.Round)
		return nil
	}
	if m.waitingForRoundChange && num == numCatchUp {
		if cv.Round.Cmp(rc.View.Round) < 0 {
			logger.Warn("[RC] Sending round change: received f+1 round change messages",
				"currentRound", cv.Round, "newRound", rc.View.Round, "seq", cv.Sequence)
			m.sendRoundChange(rc.View.Round)
		}
		return nil
	}
	if cv.Round.Cmp(rc.View.Round) < 0 {
		logger.Debug("[RC] Received higher round but not enough messages, ignored",
			"currentRound", cv.Round, "rcRound", rc.View.Round, "numRC", num, "need", numStartNewRound)
		return errIgnored
	}
	return nil
}

func (m *machine) handleTimeout(nextView *bft.View) {
	if m.b.nodetype != common.CONSENSUSNODE {
		return
	}
	lastProposal, _ := m.b.lastProposal()
	if lastProposal == nil {
		logger.Error("[RC] Timeout but last proposal is nil", "nextView", nextView)
		return
	}
	if lastProposal.Number().Cmp(nextView.Sequence) >= 0 {
		logger.Debug("[RC] Timeout outdated, chain already advanced",
			"blockNumber", lastProposal.Number().Uint64(), "timeoutSeq", nextView.Sequence, "timeoutRound", nextView.Round)
		return
	}

	if !m.waitingForRoundChange {
		maxRound := m.maxRoundChangeRound(m.f + 1)
		if maxRound != nil && maxRound.Cmp(m.round) > 0 {
			logger.Warn("[RC] Sending round change on timeout: catching up to f+1 max round",
				"currentRound", m.round, "maxRound", maxRound, "seq", m.sequence)
			m.sendRoundChange(maxRound)
			return
		}
	}

	if lastProposal.Number().Cmp(m.sequence) >= 0 {
		logger.Debug("[RC] Timeout: chain caught up, starting new sequence",
			"blockNumber", lastProposal.Number().Uint64(), "seq", m.sequence)
		m.startNewRound(common.Big0)
	} else {
		logger.Warn("[RC] Sending round change on timeout",
			"seq", m.sequence, "currentRound", m.round, "nextRound", nextView.Round,
			"waiting", m.waitingForRoundChange, "state", stateNames[m.state])
		m.sendRoundChange(nextView.Round)
	}
}

func (m *machine) handleChainHead() {
	lastProposal, _ := m.b.lastProposal()
	if lastProposal == nil || m.sequence == nil {
		return
	}
	if lastProposal.Number().Cmp(m.sequence) >= 0 {
		m.startNewRound(common.Big0)
	}
}

// ---------------------------------------------------------------------------
// Round management
// ---------------------------------------------------------------------------

func (m *machine) startNewRound(round *big.Int) {
	roundChange := false
	lastProposal, _ := m.b.lastProposal()

	if lastProposal == nil {
		logger.Error("Cannot start new round: last proposal is nil")
		return
	}

	if m.sequence == nil {
		logger.Debug("Starting initial round")
	} else if lastProposal.Number().Cmp(m.sequence) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), m.sequence)
		m.metrics.Sequence.Mark(new(big.Int).Add(diff, common.Big1).Int64())
		if !m.consensusTimestamp.IsZero() {
			m.metrics.ConsensusTime.Update(int64(time.Since(m.consensusTimestamp)))
			m.consensusTimestamp = time.Time{}
		}
		logger.Debug("Chain advanced past our sequence, starting new sequence",
			"lastBlock", lastProposal.Number().Uint64(), "currentSeq", m.sequence)
	} else if lastProposal.Number().Cmp(big.NewInt(m.sequence.Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			return // same seq, round 0 — nothing to do
		}
		if round.Cmp(m.round) < 0 {
			return // can't go backwards
		}
		roundChange = true
	} else {
		return
	}

	var newView *bft.View
	if roundChange {
		newView = &bft.View{
			Sequence: new(big.Int).Set(m.sequence),
			Round:    new(big.Int).Set(round),
		}
	} else {
		newView = &bft.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
	}

	seq, r := newView.Sequence.Uint64(), newView.Round.Uint64()
	qualified, committeeSet, proposer, committeeSize, required, fNum, err := m.getRoundCommitteeState(seq, r)
	if err != nil {
		logger.Error("Failed to get round committee state", "err", err)
		return
	}

	m.b.currentView.Store(newView)

	// Cancel any in-flight speculative execution from the previous round.
	m.b.cancelSpeculativeExecution()

	// Reset round state.
	m.sequence = newView.Sequence
	m.round = newView.Round
	m.qualified = qualified
	m.committee = committeeSet
	m.proposer = proposer
	m.committeeSize = committeeSize
	m.requiredMessageCount = required
	m.f = fNum

	if !roundChange {
		councilSize := int64(qualified.Len())
		cs := int64(committeeSize)
		if cs > councilSize {
			cs = councilSize
		}
		m.metrics.CouncilSize.Update(councilSize)
		m.metrics.CommitteeSize.Update(cs)
	}
	m.metrics.CurrentRound.Update(m.round.Int64())
	if m.isHashLocked() {
		m.metrics.HashLock.Update(1)
	} else {
		m.metrics.HashLock.Update(0)
	}

	if roundChange && m.isHashLocked() {
		// Keep the locked preprepare and proposal.
	} else if !roundChange {
		m.preprepare = nil
		m.lockedHash = common.Hash{}
		m.pendingRequest = nil
	} else {
		m.preprepare = nil
		m.lockedHash = common.Hash{}
	}
	m.prepares = newMessageSet(qualified)
	m.commits = newMessageSet(qualified)
	m.resetRoundChangeSets()
	m.waitingForRoundChange = false
	m.state = stateAcceptRequest

	m.processPendingRequests()
	m.processBacklog()

	if roundChange && m.isProposer() {
		if m.isHashLocked() && m.preprepare != nil {
			m.sendPreprepare(&bft.Request{Proposal: m.preprepare.Proposal})
		} else if m.pendingRequest != nil {
			m.sendPreprepare(m.pendingRequest)
		} else {
			// Round-change proposer with no local pendingRequest: ask worker to build lazily.
			logger.Info("Requesting local proposal build for round-change proposer",
				"seq", newView.Sequence.Uint64(), "round", newView.Round.Uint64(),
				"proposer", m.proposer, "self", m.b.address)
			m.b.eventMux.Post(consensus.NewSequenceEvent{IsProposer: true})
		}
	}

	if !roundChange {
		m.b.eventMux.Post(consensus.NewSequenceEvent{IsProposer: m.isProposer()})
	}

	m.newRoundChangeTimer()

	logger.Debug("New round", "round", newView.Round, "seq", newView.Sequence,
		"proposer", m.proposer, "isProposer", m.isProposer(),
		"roundChange", roundChange, "quorum", m.requiredMessageCount, "f", m.f)
	logger.Trace("New round committee", "round", newView.Round, "seq", newView.Sequence,
		"qualifiedSize", m.qualified.Len(), "committeeSize", m.committeeSize)
}

func (m *machine) catchUpRound(view *bft.View) {
	oldRound := new(big.Int).Set(m.round)
	oldSeq := new(big.Int).Set(m.sequence)
	oldProposer := m.proposer // capture before the update below so the log is accurate

	if diff := new(big.Int).Sub(view.Round, m.round); diff.Sign() > 0 {
		m.metrics.Round.Mark(diff.Int64())
	}

	m.waitingForRoundChange = true
	// Preserve lock state through round catch-up.
	oldPreprepare := m.preprepare
	oldLockedHash := m.lockedHash
	oldPendingReq := m.pendingRequest

	m.sequence = view.Sequence
	m.round = view.Round
	m.prepares = newMessageSet(m.qualified)
	m.commits = newMessageSet(m.qualified)

	if !common.EmptyHash(oldLockedHash) {
		m.preprepare = oldPreprepare
		m.lockedHash = oldLockedHash
	} else {
		m.preprepare = nil
		m.lockedHash = common.Hash{}
	}
	m.pendingRequest = oldPendingReq
	m.clearRoundChangeSets(view.Round)
	m.newRoundChangeTimer()

	m.metrics.CurrentRound.Update(m.round.Int64())
	if m.isHashLocked() {
		m.metrics.HashLock.Update(1)
	} else {
		m.metrics.HashLock.Update(0)
	}

	newProposer, err := m.b.valsetModule.GetProposer(view.Sequence.Uint64(), view.Round.Uint64())
	if err != nil {
		logger.Warn("Failed to get proposer for catch-up round", "err", err)
	} else {
		m.proposer = newProposer
	}
	// Keep the shared view in sync with the round-change view so isProposer()
	// and the SubmitTransactions proposer check don't act on the pre-RC round.
	m.b.currentView.Store(&bft.View{
		Sequence: new(big.Int).Set(view.Sequence),
		Round:    new(big.Int).Set(view.Round),
	})
	logger.Warn("[RC] Catch up round",
		"oldRound", oldRound, "oldSeq", oldSeq, "oldProposer", oldProposer,
		"newRound", view.Round, "newSeq", view.Sequence, "newProposer", newProposer,
		"hashLocked", !common.EmptyHash(oldLockedHash))
}

// ---------------------------------------------------------------------------
// Send helpers
// ---------------------------------------------------------------------------

func (m *machine) sendPreprepare(request *bft.Request) {
	header := request.Proposal.Header()
	m.b.sealer.WriteRound(header, m.currentView().Round.Int64())
	request.Proposal = request.Proposal.WithSeal(header)

	if m.sequence.Cmp(request.Proposal.Number()) != 0 || !m.isProposer() {
		return
	}
	curView := m.currentView()
	pp := &bft.Preprepare{View: curView, Proposal: request.Proposal}
	encoded, err := bft.Encode(pp)
	if err != nil {
		return
	}
	msg := &bft.Message{
		Hash: request.Proposal.ParentHash(),
		Code: bft.MsgPreprepare,
		Msg:  encoded,
	}
	// Self-accept and gossip directly to peers; skip the self-loop (otherwise
	// the proposer pays the decode + verify on its own block). Skipping the
	// self-loop also skips checkMessage, so re-check waitingForRoundChange here:
	// a late request during round change must not gossip a stale-round proposal.
	if m.state == stateAcceptRequest && !m.isHashLocked() && !m.waitingForRoundChange {
		payload := m.signPayload(msg)
		if payload == nil {
			return
		}
		m.b.gossipSubPeer(msg.Hash, payload)
		m.acceptPreprepare(pp)
		m.setState(statePreprepared)
		m.sendPrepare()
		return
	}
	// Hash-locked round-change path keeps the legacy self-loop.
	m.broadcastMsg(msg)
}

func (m *machine) sendPrepare() {
	if !m.committee.Contains(m.b.address) {
		return
	}
	sub := m.subject()
	if sub == nil {
		return
	}
	encoded, err := bft.Encode(sub)
	if err != nil {
		return
	}
	m.broadcastMsg(&bft.Message{
		Hash: m.preprepare.Proposal.ParentHash(),
		Code: bft.MsgPrepare,
		Msg:  encoded,
	})
}

func (m *machine) sendCommit() {
	if m.preprepare == nil || !m.committee.Contains(m.b.address) {
		return
	}
	sub := m.subject()
	if sub == nil {
		return
	}
	m.broadcastCommit(sub)
}

func (m *machine) sendCommitForOldBlock(view *bft.View, digest, prevHash common.Hash) {
	m.broadcastCommit(&bft.Subject{View: view, Digest: digest, PrevHash: prevHash})
}

func (m *machine) broadcastCommit(sub *bft.Subject) {
	encoded, err := bft.Encode(sub)
	if err != nil {
		return
	}
	m.broadcastMsg(&bft.Message{
		Hash: sub.PrevHash,
		Code: bft.MsgCommit,
		Msg:  encoded,
	})
}

func (m *machine) sendNextRoundChange(reason string) {
	if m.b.nodetype != common.CONSENSUSNODE {
		return
	}
	logger.Warn("[RC] sendNextRoundChange", "where", reason)
	m.sendRoundChange(new(big.Int).Add(m.currentView().Round, common.Big1))
}

func (m *machine) sendRoundChange(round *big.Int) {
	cv := m.currentView()
	if cv.Round.Cmp(round) >= 0 {
		logger.Warn("[RC] Skip sending round change, current round >= target",
			"current", cv.Round, "target", round, "seq", cv.Sequence)
		return
	}

	logger.Warn("[RC] Sending round change",
		"seq", cv.Sequence, "currentRound", cv.Round, "targetRound", round,
		"proposer", m.proposer, "prepares", m.prepares.Size(), "commits", m.commits.Size())

	m.catchUpRound(&bft.View{
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(cv.Sequence),
	})

	lastProposal, _ := m.b.lastProposal()
	cv = m.currentView()
	rc := &bft.Subject{
		View:     cv,
		Digest:   common.Hash{},
		PrevHash: lastProposal.Hash(),
	}
	payload, err := bft.Encode(rc)
	if err != nil {
		return
	}
	m.broadcastMsg(&bft.Message{
		Hash: rc.PrevHash,
		Code: bft.MsgRoundChange,
		Msg:  payload,
	})
}

func (m *machine) broadcastMsg(msg *bft.Message) {
	payload := m.signPayload(msg)
	if payload == nil {
		return
	}
	if err := m.b.broadcast(msg.Hash, payload); err != nil {
		logger.Error("Failed to broadcast message", "err", err)
	}
}

func (m *machine) signPayload(msg *bft.Message) []byte {
	msg.Address = m.b.address
	msg.CommittedSeal = []byte{}
	// A COMMIT seal signs the digest carried in the message's own subject (the
	// block this COMMIT votes for) rather than the current proposal, so a COMMIT
	// re-sent for an old block (sendCommitForOldBlock) attests that block. After
	// the permissionless fork the preimage also binds the round.
	if msg.Code == bft.MsgCommit {
		var sub *bft.Subject
		if err := msg.Decode(&sub); err != nil {
			return nil
		}
		seal, err := m.b.sign(m.b.committedSealPreimage(sub.View, sub.Digest))
		if err != nil {
			return nil
		}
		msg.CommittedSeal = seal
	}
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil
	}
	msg.Signature, err = m.b.sign(data)
	if err != nil {
		return nil
	}
	payload, err := msg.Payload()
	if err != nil {
		return nil
	}
	return payload
}

func (m *machine) doCommit() {
	m.setState(stateCommitted)

	proposal := m.proposal()
	if proposal == nil {
		m.unlockHash()
		m.sendNextRoundChange("commit: nil proposal")
		return
	}

	committedSeals := make([][]byte, m.commits.Size())
	for i, v := range m.commits.Values() {
		committedSeals[i] = make([]byte, len(v.CommittedSeal))
		copy(committedSeals[i], v.CommittedSeal)
	}

	if err := m.b.commit(proposal, committedSeals); err != nil {
		m.unlockHash()
		m.sendNextRoundChange("commit failure")
	}
}

// ---------------------------------------------------------------------------
// State helpers
// ---------------------------------------------------------------------------

func (m *machine) currentView() *bft.View {
	return &bft.View{
		Sequence: new(big.Int).Set(m.sequence),
		Round:    new(big.Int).Set(m.round),
	}
}

func (m *machine) isProposer() bool {
	return m.proposer == m.b.address
}

func (m *machine) subject() *bft.Subject {
	if m.preprepare == nil {
		return nil
	}
	return &bft.Subject{
		View: &bft.View{
			Round:    new(big.Int).Set(m.round),
			Sequence: new(big.Int).Set(m.sequence),
		},
		Digest:   m.preprepare.Proposal.Hash(),
		PrevHash: m.preprepare.Proposal.ParentHash(),
	}
}

func (m *machine) proposal() bft.Proposal {
	if m.preprepare != nil {
		return m.preprepare.Proposal
	}
	return nil
}

// proposalNumberMatchesView reports whether the proposal's block number equals
// the view sequence.
func proposalNumberMatchesView(pp *bft.Preprepare) bool {
	if pp == nil || pp.Proposal == nil || pp.Proposal.Number() == nil ||
		pp.View == nil || pp.View.Sequence == nil {
		return false
	}
	return pp.Proposal.Number().Cmp(pp.View.Sequence) == 0
}

func (m *machine) acceptPreprepare(pp *bft.Preprepare) {
	m.consensusTimestamp = time.Now()
	m.preprepare = pp
	m.postPrepreparedEvent(pp)
}

func (m *machine) setState(s uint64) {
	m.state = s
	if s == stateAcceptRequest {
		m.processPendingRequests()
	}
	m.processBacklog()
}

func (m *machine) lockHash() {
	if m.preprepare != nil {
		m.lockedHash = m.preprepare.Proposal.Hash()
	}
}

func (m *machine) unlockHash() {
	m.lockedHash = common.Hash{}
}

func (m *machine) isHashLocked() bool {
	if common.EmptyHash(m.lockedHash) {
		return false
	}
	return !m.b.hasBadProposal(m.lockedHash)
}

func (m *machine) getPrepareOrCommitSize() int {
	result := m.prepares.Size() + m.commits.Size()
	for _, msg := range m.prepares.Values() {
		if m.commits.Get(msg.Address) != nil {
			result--
		}
	}
	return result
}

func (m *machine) verifySubject(sub *bft.Subject) error {
	expected := m.subject()
	if expected == nil || !sub.Equal(expected) {
		return errInconsistentSubject
	}
	return nil
}

func (m *machine) checkMessage(msgCode uint64, view *bft.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return bft.ErrInvalidMessage
	}
	// The round is consumed as a uint64 (round-change buckets, backlog
	// priority, proposer lookup); reject out-of-range values for every message
	// code, as istanbul's checkMessage does.
	if !view.Round.IsUint64() {
		return bft.ErrInvalidMessage
	}
	cv := m.currentView()

	if msgCode == bft.MsgRoundChange {
		if view.Sequence.Cmp(cv.Sequence) > 0 {
			return errFutureMessage
		}
		if view.Cmp(cv) < 0 {
			return errOldMessage
		}
		return nil
	}

	if view.Cmp(cv) > 0 {
		return errFutureMessage
	}
	if view.Cmp(cv) < 0 {
		return errOldMessage
	}

	if m.waitingForRoundChange {
		return errFutureMessage
	}

	if m.state == stateAcceptRequest && msgCode > bft.MsgPreprepare {
		return errFutureMessage
	}
	return nil
}

func (m *machine) checkValidatorSignature(data, sig []byte) (common.Address, error) {
	signer, err := bft.GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	if m.qualified != nil && m.qualified.Contains(signer) {
		return signer, nil
	}
	return common.Address{}, errUnauthorizedAddress
}

// ---------------------------------------------------------------------------
// Committee helpers
// ---------------------------------------------------------------------------

func (m *machine) getRoundCommitteeState(seq, r uint64) (
	qualified *valset.AddressSet,
	committeeSet *valset.AddressSet,
	proposer common.Address,
	committeeSize uint64,
	requiredMsgCnt int,
	fNum int,
	err error,
) {
	if m.b.valsetModule == nil || m.b.govModule == nil {
		return nil, nil, common.Address{}, 0, 0, 0, errNoModule
	}
	council, err := m.b.valsetModule.GetCouncil(seq)
	if err != nil {
		return qualified, committeeSet, proposer, committeeSize, requiredMsgCnt, fNum, err
	}
	demoted, err := m.b.valsetModule.GetDemotedValidators(seq)
	if err != nil {
		return qualified, committeeSet, proposer, committeeSize, requiredMsgCnt, fNum, err
	}
	qualified = valset.NewAddressSet(council).Subtract(valset.NewAddressSet(demoted))

	committeeAddrs, err := m.b.valsetModule.GetCommittee(seq, r)
	if err != nil {
		return qualified, committeeSet, proposer, committeeSize, requiredMsgCnt, fNum, err
	}
	proposer, err = m.b.valsetModule.GetProposer(seq, r)
	if err != nil {
		return qualified, committeeSet, proposer, committeeSize, requiredMsgCnt, fNum, err
	}

	committeeSet = valset.NewAddressSet(committeeAddrs)
	// Size the quorum by the committee actually selected. Post-permissionless
	// the committee is the qualified set and the governance CommitteeSize no
	// longer bounds it.
	committeeSize = uint64(committeeSet.Len())

	qLen := qualified.Len()
	requiredMsgCnt = calcQuorumSize(qLen, committeeSize)
	fNum = calcFaultTolerance(qLen, committeeSize)
	return qualified, committeeSet, proposer, committeeSize, requiredMsgCnt, fNum, err
}

func calcQuorumSize(qualifiedLen int, committeeSize uint64) int {
	size := min(qualifiedLen, int(committeeSize))
	if size < 4 {
		return size
	}
	return int(math.Ceil(float64(2*size) / 3))
}

func calcFaultTolerance(qualifiedLen int, committeeSize uint64) int {
	if qualifiedLen > int(committeeSize) {
		return int(math.Ceil(float64(committeeSize)/3)) - 1
	}
	return int(math.Ceil(float64(qualifiedLen)/3)) - 1
}

// ---------------------------------------------------------------------------
// Timer
// ---------------------------------------------------------------------------

func (m *machine) newRoundChangeTimer() {
	m.stopTimer()
	timeout := time.Duration(m.b.timeout) * time.Millisecond
	round := m.round.Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
	}

	seq := new(big.Int).Set(m.sequence)
	r := new(big.Int).Set(m.round)

	// Capture current state for logging inside the timer callback.
	preprepareNil := m.preprepare == nil
	preparesSize := m.prepares.Size()
	commitsSize := m.commits.Size()
	proposer := m.proposer
	state := m.state

	loc := "startNewRound"
	if round > 0 {
		loc = "catchUpRound"
	}

	m.roundChangeTimer.Store(time.AfterFunc(timeout, func() {
		if m.b.nodetype == common.CONSENSUSNODE {
			logger.Warn("[RC] Timeout fired",
				"setBy", loc, "seq", seq, "round", r, "proposer", proposer,
				"preprepareNil", preprepareNil, "prepares", preparesSize,
				"commits", commitsSize, "state", stateNames[state])
		}
		m.b.eventMux.Post(timeoutEvent{&bft.View{
			Sequence: seq,
			Round:    new(big.Int).Add(r, common.Big1),
		}})
	}))

	logger.Debug("[RC] Timer set", "seq", seq, "round", round, "timeout", timeout)
}

func (m *machine) stopFuturePreprepareTimer() {
	if m.futurePreprepareTimer != nil {
		m.futurePreprepareTimer.Stop()
		m.futurePreprepareTimer = nil
	}
}

func (m *machine) stopTimer() {
	m.stopFuturePreprepareTimer()
	if t := m.roundChangeTimer.Load(); t != nil {
		t.(*time.Timer).Stop()
	}
}

// ---------------------------------------------------------------------------
// Backlog
// ---------------------------------------------------------------------------

var msgPriority = map[uint64]int{
	bft.MsgPreprepare: 1,
	bft.MsgCommit:     2,
	bft.MsgPrepare:    3,
}

// Retention limits, the same as istanbul core. Future messages are kept per
// sender only, so no sender can take another's budget. A retained PREPARE,
// COMMIT or ROUND CHANGE is bounded by checkMessageSize; the one retained
// PREPREPARE per sender is stored before verify, so it is bounded by the p2p
// message size limit rather than by the block size cap.
const (
	// maxBacklogSequencesAhead keeps only a small future-sequence window; a node
	// further behind catches up through block synchronization.
	maxBacklogSequencesAhead = 8
	// maxBacklogMessagesPerSender bounds queued PREPARE, COMMIT and ROUND
	// CHANGE per sender. A PREPREPARE takes the sender's single slot instead.
	maxBacklogMessagesPerSender = 128
	// maxRoundChangeRoundsAhead retains the current round plus this many
	// future rounds of ROUND CHANGE while a sequence is stalled.
	maxRoundChangeRoundsAhead = 128
	// maxSubjectMessageBytes bounds a PREPARE, COMMIT or ROUND CHANGE, whose
	// well-formed size is a few hundred bytes.
	maxSubjectMessageBytes = 1024
)

// checkMessageSize rejects an oversized PREPARE, COMMIT or ROUND CHANGE before
// any retention path. A PREPREPARE carries a block and is exempt; the single
// backlog slot per sender bounds how many of them are retained.
func checkMessageSize(msg *bft.Message) error {
	if msg.Code == bft.MsgPreprepare {
		return nil
	}
	if retainedMessageBytes(msg) > maxSubjectMessageBytes {
		return errMessageTooLarge
	}
	return nil
}

// retainedMessageBytes reports the memory a retained message occupies.
func retainedMessageBytes(msg *bft.Message) uint64 {
	return uint64(len(msg.Msg)) + uint64(len(msg.Signature)) + uint64(len(msg.CommittedSeal))
}

// backlogPreprepare is the one PREPREPARE retained per sender, with its view
// and parent hash decoded once at retention.
type backlogPreprepare struct {
	msg      *bft.Message
	view     *bft.View
	prevHash common.Hash
}

// storeBacklog retains a future message for processBacklog. PREPARE, COMMIT
// and ROUND CHANGE queue per sender up to maxBacklogMessagesPerSender. A
// PREPREPARE takes the sender's single slot, replacing a retained one only when
// its view is higher.
func (m *machine) storeBacklog(msg *bft.Message, src common.Address) {
	if src == m.b.address {
		return
	}
	m.backlogsMu.Lock()
	defer m.backlogsMu.Unlock()

	view, prevHash := backlogMessageView(msg)
	if view == nil || m.isBacklogSequenceTooFar(view.Sequence) {
		return
	}

	if msg.Code == bft.MsgPreprepare {
		if held, ok := m.backlogPreprepares[src]; ok && view.Cmp(held.view) <= 0 {
			return
		}
		m.backlogPreprepares[src] = backlogPreprepare{msg: msg, view: view, prevHash: prevHash}
		return
	}

	if m.backlogCounts[src] >= maxBacklogMessagesPerSender {
		return
	}
	backlog := m.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
		m.backlogs[src] = backlog
	}
	// toPriority truncates the sequence, so it runs only after
	// isBacklogSequenceTooFar has rejected sequences that do not fit in uint64.
	backlog.Push(msg, toPriority(msg.Code, view))
	m.backlogCounts[src]++
}

func (m *machine) isBacklogSequenceTooFar(sequence *big.Int) bool {
	if !sequence.IsUint64() {
		return true
	}
	maxSequence := new(big.Int).Add(m.sequence, big.NewInt(maxBacklogSequencesAhead))
	return sequence.Cmp(maxSequence) > 0
}

// removeBacklogMessage releases one queued message of a sender while
// backlogsMu is held.
func (m *machine) removeBacklogMessage(src common.Address) {
	if m.backlogCounts[src] <= 1 {
		delete(m.backlogCounts, src)
		return
	}
	m.backlogCounts[src]--
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
func (m *machine) postBacklogMessage(src common.Address, msg *bft.Message, view *bft.View, prevHash common.Hash) (future bool) {
	if view == nil {
		return false
	}
	err := m.checkMessage(msg.Code, view)
	if err == errFutureMessage {
		return true
	}
	if err != nil {
		return false
	}
	go m.b.eventMux.Post(backlogEvent{src: src, msg: msg, Hash: prevHash})
	return false
}

func (m *machine) processBacklog() {
	m.backlogsMu.Lock()
	defer m.backlogsMu.Unlock()

	for src, held := range m.backlogPreprepares {
		if !m.postBacklogMessage(src, held.msg, held.view, held.prevHash) {
			delete(m.backlogPreprepares, src)
		}
	}

	for src, backlog := range m.backlogs {
		// The queue is ordered by view, so stop at the first future message.
		for !backlog.Empty() {
			item, prio := backlog.Pop()
			msg := item.(*bft.Message)
			view, prevHash := backlogMessageView(msg)
			if m.postBacklogMessage(src, msg, view, prevHash) {
				backlog.Push(msg, prio)
				break
			}
			m.removeBacklogMessage(src)
		}
		if backlog.Empty() {
			delete(m.backlogs, src)
		}
	}
}

func (m *machine) storePendingRequest(request *bft.Request) {
	// For simplicity, store only one pending request (latest).
	m.pendingRequest = request
}

func (m *machine) processPendingRequests() {
	if m.pendingRequest == nil {
		return
	}
	cmp := m.sequence.Cmp(m.pendingRequest.Proposal.Number())
	if cmp == 0 {
		go m.b.eventMux.Post(requestEvent{Proposal: m.pendingRequest.Proposal})
	} else if cmp > 0 {
		m.pendingRequest = nil
	}
	// If cmp < 0 (future), keep for later.
}

func toPriority(msgCode uint64, view *bft.View) int64 {
	if msgCode == bft.MsgRoundChange {
		return -int64(view.Sequence.Uint64() * 1000)
	}
	return -int64(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
}

// ---------------------------------------------------------------------------
// Round change set
// ---------------------------------------------------------------------------

// addRoundChange retains a ROUND CHANGE only when its round is within
// maxRoundChangeRoundsAhead of the current round, and keeps at most a quorum's
// worth of senders per round: a further distinct sender cannot change the
// outcome. A rejected message never leaves an empty bucket behind.
func (m *machine) addRoundChange(round *big.Int, msg *bft.Message) (int, error) {
	if !round.IsUint64() {
		return 0, bft.ErrInvalidMessage
	}
	m.roundChangesMu.Lock()
	defer m.roundChangesMu.Unlock()

	maxRound := new(big.Int).Add(m.round, big.NewInt(maxRoundChangeRoundsAhead))
	if round.Cmp(maxRound) > 0 {
		return 0, errRoundChangeTooFar
	}
	r := round.Uint64()
	messages := m.roundChangeSets[r]
	if messages == nil {
		messages = newMessageSet(m.qualified)
	}
	if messages.Get(msg.Address) == nil && messages.Size() >= m.requiredMessageCount {
		return 0, errRoundChangeLimit
	}
	if err := messages.Add(msg); err != nil {
		return 0, err
	}
	m.roundChangeSets[r] = messages
	return messages.Size(), nil
}

// resetRoundChangeSets creates a fresh (empty) round change set, matching
// Istanbul's behaviour in startNewRound.  This ensures stale entries from a
// previous sequence (which may carry an outdated qualified-validator set)
// are never reused.
func (m *machine) resetRoundChangeSets() {
	m.roundChangesMu.Lock()
	defer m.roundChangesMu.Unlock()
	m.roundChangeSets = make(map[uint64]*messageSet)
}

// clearRoundChangeSets removes entries for rounds strictly below the given
// round.  Used in catchUpRound where we want to preserve higher-round entries
// within the same sequence (same qualified set).
func (m *machine) clearRoundChangeSets(round *big.Int) {
	m.roundChangesMu.Lock()
	defer m.roundChangesMu.Unlock()
	for k, rms := range m.roundChangeSets {
		if rms.Size() == 0 || k < round.Uint64() {
			delete(m.roundChangeSets, k)
		}
	}
}

func (m *machine) maxRoundChangeRound(minCount int) *big.Int {
	m.roundChangesMu.Lock()
	defer m.roundChangesMu.Unlock()
	var maxRound *big.Int
	for k, rms := range m.roundChangeSets {
		if rms.Size() < minCount {
			continue
		}
		r := big.NewInt(int64(k))
		if maxRound == nil || maxRound.Cmp(r) < 0 {
			maxRound = r
		}
	}
	return maxRound
}

// ---------------------------------------------------------------------------
// Message set — deduplicated per validator address
// ---------------------------------------------------------------------------

type messageSet struct {
	mu       sync.Mutex
	messages map[common.Address]*bft.Message
	allowed  *valset.AddressSet
}

func newMessageSet(allowed *valset.AddressSet) *messageSet {
	return &messageSet{
		messages: make(map[common.Address]*bft.Message),
		allowed:  allowed,
	}
}

func (ms *messageSet) Add(msg *bft.Message) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if ms.allowed != nil && !ms.allowed.Contains(msg.Address) {
		return errUnauthorizedAddress
	}
	ms.messages[msg.Address] = msg
	return nil
}

func (ms *messageSet) Values() []*bft.Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	result := make([]*bft.Message, 0, len(ms.messages))
	for _, v := range ms.messages {
		result = append(result, v)
	}
	return result
}

func (ms *messageSet) Size() int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return len(ms.messages)
}

func (ms *messageSet) Get(addr common.Address) *bft.Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.messages[addr]
}
