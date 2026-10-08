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
// This file is derived from quorum/consensus/istanbul/core/roundchange.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/rlp"
)

// maxRoundChangeRoundsAhead retains the current round plus this many future
// rounds while a sequence is stalled. The window bounds memory while keeping the
// near-round buckets the node needs to catch up. It deliberately trades catch-up
// beyond this window for bounded local memory.
const maxRoundChangeRoundsAhead = 128

// sendNextRoundChange sends the ROUND CHANGE message with current round + 1
func (c *core) sendNextRoundChange(loc string) {
	if c.backend.NodeType() != common.CONSENSUSNODE {
		return
	}
	logger.Warn("[RC] sendNextRoundChange happened", "where", loc)
	c.sendRoundChange(new(big.Int).Add(c.currentView().Round, common.Big1))
}

// sendRoundChange sends the ROUND CHANGE message with the given round
func (c *core) sendRoundChange(round *big.Int) {
	logger := c.logger.NewWith("state", c.state)

	cv := c.currentView()
	if cv.Round.Cmp(round) >= 0 {
		logger.Warn("[RC] Skip sending out the round change message", "current round", cv.Round,
			"target round", round)
		return
	}

	logger.Warn("[RC] Prepare messages received before catchUpRound",
		"len(prepares)", c.current.Prepares.Size(), "messages", c.current.Prepares.GetMessages())
	logger.Warn("[RC] Commit messages received before catchUpRound",
		"len(commits)", c.current.Commits.Size(), "messages", c.current.Commits.GetMessages())

	preparedCertificate := c.current.PreparedCertificate()
	c.catchUpRound(&bft.View{
		// The round number we'd like to transfer to.
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(cv.Sequence),
	})

	lastProposal, _ := c.backend.LastProposal()

	// Now we have the new round number and sequence number
	cv = c.currentView()
	rc := &bft.RoundChange{
		View: cv,
	}
	// Sign only the prepared round and digest. The certificate itself rides as
	// the unsigned Evidence so a later PRE-PREPARE can carry this signed
	// message without repeating the prepared block.
	var evidence []byte
	if c.backend.IsPermissionlessAt(cv.Sequence.Uint64()) && preparedCertificate != nil {
		var err error
		if evidence, err = bft.Encode(preparedCertificate); err != nil {
			logger.Error("Failed to encode prepared certificate", "err", err)
			return
		}
		rc.Prepared = &bft.PreparedClaim{
			Round:  new(big.Int).Set(preparedCertificate.View.Round),
			Digest: preparedCertificate.Proposal.Hash(),
		}
	}

	payload, err := bft.Encode(rc)
	if err != nil {
		logger.Error("Failed to encode ROUND CHANGE", "rc", rc, "err", err)
		return
	}

	c.broadcast(&bft.Message{
		PrevHash: lastProposal.Hash(),
		Code:     bft.MsgRoundChange,
		Msg:      payload,
		Evidence: evidence,
	})
}

func (c *core) handleRoundChange(msg *bft.Message, src common.Address) error {
	logger := c.logger.NewWith("state", c.state, "from", src.Hex())

	// Decode ROUND CHANGE message
	var rc *bft.RoundChange
	if err := msg.Decode(&rc); err != nil {
		logger.Error("Failed to decode message", "code", msg.Code, "err", err)
		return bft.ErrInvalidMessage
	}
	if rc == nil || rc.View == nil || rc.View.Sequence == nil {
		return bft.ErrInvalidMessage
	}
	permissionless := c.backend != nil && c.backend.IsPermissionlessAt(rc.View.Sequence.Uint64())
	// Reject a claim or Evidence at a height that does not accept them. The
	// decoder cannot produce one there, so this is defense in depth.
	if c.backend != nil && !permissionless && (len(msg.Evidence) != 0 || rc.Prepared != nil) {
		return bft.ErrInvalidMessage
	}
	// The full prepared certificate is expensive to verify, but its presence
	// must match the signed claim before checkMessage can return errFutureMessage
	// and retain the envelope. This rejects a large, unsigned attachment with no
	// claim without doing certificate decoding or signature recovery.
	if permissionless && !roundChangeEvidenceShapeValid(msg, rc) {
		return bft.ErrInvalidMessage
	}

	if err := c.checkMessage(msg, rc.View); err != nil {
		return err
	}
	// Some focused unit tests exercise round-change-set admission with a bare
	// core and no backend. Production cores always have one; the extension
	// validation itself requires it and therefore applies only when present.
	if permissionless {
		if err := c.verifyRoundChangeEvidence(msg, rc); err != nil {
			logger.Warn("Invalid prepared certificate in ROUND CHANGE", "err", err)
			return bft.ErrInvalidMessage
		}
	}

	// Restrict round-change admission to committee members, mirroring handleCommit.
	// The thresholds (requiredMessageCount, f) are committee-sized, but messageSet
	// admits any qualified validator; gating on the committee keeps admission
	// consistent with the thresholds. The committee equals the qualified set (now and
	// post-permissionless), so this is round-agnostic.
	if !c.current.committee.Contains(src) {
		logger.Warn("received an istanbul round change message from non-committee",
			"currentSequence", c.current.sequence.Uint64(), "sender", src.Hex(), "msgView", rc.View.String())
		return errNotFromCommittee
	}

	cv := c.currentView()
	roundView := rc.View

	// Add the ROUND CHANGE message to its message set and return how many
	// messages we've got with the same round number and sequence number.
	num, err := c.roundChangeSet.Add(cv.Round, roundView.Round, msg)
	if err != nil {
		logger.Trace("Discarding ROUND CHANGE", "from", src, "round", roundView.Round, "err", err)
		return err
	}

	var numCatchUp, numStartNewRound int
	n := c.current.requiredMessageCount
	f := c.current.f
	// N ROUND CHANGE messages can start new round.
	numStartNewRound = n
	// F + 1 ROUND CHANGE messages can start catch up the round.
	numCatchUp = f + 1

	if num == numStartNewRound && (c.waitingForRoundChange || cv.Round.Cmp(roundView.Round) < 0) {
		// We've received enough ROUND CHANGE messages, start a new round immediately.
		logger.Warn("[RC] Prepare messages received before startNewRound", "round", cv.Round.String(),
			"len(prepares)", c.current.Prepares.Size(), "messages", c.current.Prepares.GetMessages())
		logger.Warn("[RC] Commit messages received before startNewRound", "round", cv.Round.String(),
			"len(commits)", c.current.Commits.Size(), "messages", c.current.Commits.GetMessages())
		logger.Warn("[RC] Received 2f+1 Round Change Messages. Starting new round",
			"currentRound", cv.Round.String(), "newRound", roundView.Round.String())
		if c.backend.IsPermissionlessAt(roundView.Sequence.Uint64()) {
			c.roundChangeCertificate = c.roundChangeSet.Values(roundView.Round)
		}
		c.startNewRound(roundView.Round)
		return nil
	} else if c.waitingForRoundChange && num == numCatchUp {
		// Once we received enough ROUND CHANGE messages, those messages form a weak certificate.
		// If our round number is smaller than the certificate's round number, we would
		// try to catch up the round number.
		if cv.Round.Cmp(roundView.Round) < 0 {
			logger.Warn("[RC] Send round change because we have f+1 round change messages",
				"currentRound", cv.Round.String(), "newRound", roundView.Round.String())
			c.sendRoundChange(roundView.Round)
		}
		return nil
	} else if cv.Round.Cmp(roundView.Round) < 0 {
		// Only gossip the message with current round to other validators.
		logger.Trace("[RC] Received round is bigger but not enough number of messages. Message ignored",
			"currentRound", cv.Round.String(), "newRound", roundView.Round.String(), "numRC", num)
		return errIgnored
	}
	return nil
}

// verifyRoundChangeEvidence checks that a ROUND CHANGE's unsigned Evidence
// proves exactly the PreparedClaim its sender signed.
func (c *core) verifyRoundChangeEvidence(msg *bft.Message, rc *bft.RoundChange) error {
	if !roundChangeEvidenceShapeValid(msg, rc) {
		return errors.New("prepared claim and evidence must appear together")
	}
	if rc.Prepared == nil {
		return nil
	}
	var cert *bft.PreparedCertificate
	if err := rlp.DecodeBytes(msg.Evidence, &cert); err != nil {
		return err
	}
	// Bind the attachment to the signed claim before the expensive checks.
	if cert == nil || cert.View == nil || cert.View.Round == nil || cert.Proposal == nil || rc.Prepared.Round == nil ||
		cert.View.Round.Cmp(rc.Prepared.Round) != 0 || cert.Proposal.Hash() != rc.Prepared.Digest {
		return errors.New("evidence does not match the prepared claim")
	}
	return c.verifyPreparedCertificate(cert, rc.View)
}

func roundChangeEvidenceShapeValid(msg *bft.Message, rc *bft.RoundChange) bool {
	return (rc.Prepared == nil) == (len(msg.Evidence) == 0)
}

// verifyPreparedCertificate authenticates a quorum of votes for one proposal
// in a round strictly before the target ROUND-CHANGE view.
func (c *core) verifyPreparedCertificate(cert *bft.PreparedCertificate, target *bft.View) error {
	if cert == nil || cert.View == nil || cert.Proposal == nil || target == nil ||
		cert.View.Sequence == nil || cert.View.Round == nil || target.Sequence == nil || target.Round == nil {
		return errors.New("incomplete prepared certificate")
	}
	if cert.View.Sequence.Cmp(target.Sequence) != 0 || cert.View.Round.Cmp(target.Round) >= 0 ||
		cert.Proposal.Number().Cmp(cert.View.Sequence) != 0 {
		return errors.New("prepared certificate has invalid view")
	}
	// Votes bind only the header hash. Bind the body as well, or a peer could
	// pair copied votes with different transactions and make the next proposer
	// re-propose a block that every validator rejects.
	header := cert.Proposal.Header()
	if types.DeriveTransactionsRoot(cert.Proposal.Transactions(), header.Number) != header.TxHash {
		return errors.New("prepared certificate proposal body does not match its header")
	}
	return c.verifyPreparedCertificateVotes(cert)
}

// verifyPreparedCertificateVotes performs the expensive committee, signature,
// subject, seal, and quorum checks after the certificate view has been checked.
func (c *core) verifyPreparedCertificateVotes(cert *bft.PreparedCertificate) error {
	_, committee, _, _, quorum, _, err := getRoundCommitteeState(c, cert.View.Sequence.Uint64(), cert.View.Round.Uint64())
	if err != nil {
		return err
	}
	// A certificate has at most one meaningful vote per committee member. Reject
	// excess entries before signature recovery so a Byzantine peer cannot turn a
	// syntactically valid but oversized certificate into avoidable CPU work.
	if len(cert.Messages) > committee.Len() {
		return fmt.Errorf("prepared certificate has %d votes, maximum is %d", len(cert.Messages), committee.Len())
	}
	seen := make(map[common.Address]struct{}, len(cert.Messages))
	for _, vote := range cert.Messages {
		if vote == nil || (vote.Code != bft.MsgPrepare && vote.Code != bft.MsgCommit) ||
			len(vote.Evidence) != 0 || !committee.Contains(vote.Address) {
			return errors.New("prepared certificate contains an ineligible vote")
		}
		if _, duplicate := seen[vote.Address]; duplicate {
			return fmt.Errorf("prepared certificate contains duplicate voter %s", vote.Address)
		}
		unsigned, err := vote.PayloadNoSig()
		if err != nil {
			return err
		}
		signer, err := c.validateFn(unsigned, vote.Signature)
		if err != nil || signer != vote.Address {
			return errors.New("prepared certificate contains invalid signature")
		}
		var view *bft.View
		var digest common.Hash
		var committedSeal []byte
		switch vote.Code {
		case bft.MsgPrepare:
			var prepare *bft.Prepare
			if err := vote.Decode(&prepare); err != nil || prepare == nil {
				return errors.New("prepared certificate contains invalid prepare")
			}
			view, digest = prepare.View, prepare.Digest
		case bft.MsgCommit:
			var commit *bft.Commit
			if err := vote.Decode(&commit); err != nil || commit == nil {
				return errors.New("prepared certificate contains invalid commit")
			}
			view, digest, committedSeal = commit.View, commit.Digest, commit.CommittedSeal
		}
		if view == nil || view.Cmp(cert.View) != 0 || digest != cert.Proposal.Hash() || vote.PrevHash != cert.Proposal.ParentHash() {
			return errors.New("prepared certificate vote has inconsistent payload")
		}
		if vote.Code == bft.MsgCommit {
			preimage := istanbul.PrepareCommittedSeal(digest)
			if c.backend.IsPermissionlessAt(view.Sequence.Uint64()) {
				preimage = istanbul.PrepareCommittedSealWithRound(digest, byte(view.Round.Uint64()))
			}
			committer, err := istanbul.GetSignatureAddress(preimage, committedSeal)
			if err != nil || committer != vote.Address {
				return errors.New("prepared certificate contains invalid committed seal")
			}
		}
		seen[vote.Address] = struct{}{}
	}
	if len(seen) < quorum {
		return fmt.Errorf("prepared certificate has %d votes, need %d", len(seen), quorum)
	}
	return nil
}

// verifyRoundChangeCertificate checks the signed ROUND-CHANGE quorum carried by
// a higher-round PRE-PREPARE and returns the highest prepared round claimed in
// it, or nil when no sender claims one. Only claims are checked here; the
// PRE-PREPARE must separately prove the highest claim for its proposal.
func (c *core) verifyRoundChangeCertificate(messages []*bft.Message, target *bft.View) (*big.Int, []*bft.PreparedClaim, error) {
	_, committee, _, _, quorum, _, err := getRoundCommitteeState(c, target.Sequence.Uint64(), target.Round.Uint64())
	if err != nil {
		return nil, nil, err
	}
	// A round-change certificate has at most one meaningful message from each
	// committee member.
	if len(messages) > committee.Len() {
		return nil, nil, fmt.Errorf("round-change certificate has %d messages, maximum is %d", len(messages), committee.Len())
	}
	lastProposal, _ := c.backend.LastProposal()
	if lastProposal == nil {
		return nil, nil, errors.New("last proposal unavailable")
	}
	seen := make(map[common.Address]struct{}, len(messages))
	var highest *big.Int
	claims := make([]*bft.PreparedClaim, 0, len(messages))
	for _, message := range messages {
		// Embedded messages must be stripped; their evidence is not needed and
		// would let the PRE-PREPARE grow with every claim again.
		if message == nil || message.Code != bft.MsgRoundChange || len(message.Evidence) != 0 ||
			!committee.Contains(message.Address) {
			return nil, nil, errors.New("round-change certificate contains an ineligible message")
		}
		if _, duplicate := seen[message.Address]; duplicate {
			return nil, nil, fmt.Errorf("round-change certificate contains duplicate sender %s", message.Address)
		}
		unsigned, err := message.PayloadNoSig()
		if err != nil {
			return nil, nil, err
		}
		signer, err := c.validateFn(unsigned, message.Signature)
		if err != nil || signer != message.Address {
			return nil, nil, errors.New("round-change certificate contains invalid signature")
		}
		var roundChange *bft.RoundChange
		if err := message.Decode(&roundChange); err != nil {
			return nil, nil, err
		}
		if roundChange.View == nil || roundChange.View.Cmp(target) != 0 || message.PrevHash != lastProposal.Hash() {
			return nil, nil, errors.New("round-change certificate has inconsistent view")
		}
		if claim := roundChange.Prepared; claim != nil {
			if claim.Round == nil || claim.Round.Cmp(target.Round) >= 0 || common.EmptyHash(claim.Digest) {
				return nil, nil, errors.New("round-change certificate has an invalid prepared claim")
			}
			if highest == nil || claim.Round.Cmp(highest) > 0 {
				highest = claim.Round
			}
			claims = append(claims, claim)
		}
		seen[message.Address] = struct{}{}
	}
	if len(seen) < quorum {
		return nil, nil, fmt.Errorf("round-change certificate has %d messages, need %d", len(seen), quorum)
	}
	return highest, claims, nil
}

// verifyPreprepareJustification checks the round-change justification of a
// PRE-PREPARE above round 0. When the quorum claims a
// prepared value, the proposal must be one claimed at the highest round and
// PreparedMessages must prove it; the resulting certificate is returned.
func (c *core) verifyPreprepareJustification(preprepare *bft.Preprepare) (*bft.PreparedCertificate, error) {
	highest, claims, err := c.verifyRoundChangeCertificate(preprepare.RoundChangeCertificate, preprepare.View)
	if err != nil {
		return nil, err
	}
	if highest == nil {
		if len(preprepare.PreparedMessages) != 0 {
			return nil, errors.New("prepared votes without a prepared claim")
		}
		return nil, nil
	}
	block, ok := preprepare.Proposal.(*types.Block)
	if !ok {
		return nil, errors.New("proposal is not a block")
	}
	// Two different values cannot both be prepared in one round, so only a
	// proposal claimed at the highest round can be justified. An unprovable
	// claim for another value fails to match the votes below.
	claimed := slices.ContainsFunc(claims, func(claim *bft.PreparedClaim) bool {
		return claim.Round.Cmp(highest) == 0 && claim.Digest == block.Hash()
	})
	if !claimed {
		return nil, errors.New("proposal is not the highest prepared value")
	}
	cert := &bft.PreparedCertificate{
		View:     &bft.View{Sequence: new(big.Int).Set(preprepare.View.Sequence), Round: new(big.Int).Set(highest)},
		Proposal: block,
		Messages: preprepare.PreparedMessages,
	}
	if err := c.verifyPreparedCertificateVotes(cert); err != nil {
		return nil, err
	}
	return cert, nil
}

// roundChangeJustification selects what a proposer attaches to a PRE-PREPARE
// for the given view: the stripped signed ROUND CHANGE quorum, plus the
// certificate behind its highest prepared claim. Every message was verified
// together with its Evidence when it entered the round-change set.
func (c *core) roundChangeJustification(messages []*bft.Message, target *bft.View) ([]*bft.Message, *bft.PreparedCertificate, error) {
	certificate := make([]*bft.Message, 0, len(messages))
	var highest *bft.PreparedCertificate
	for _, message := range messages {
		certificate = append(certificate, message.WithoutEvidence())
		var roundChange *bft.RoundChange
		if err := message.Decode(&roundChange); err != nil {
			return nil, nil, err
		}
		if roundChange.Prepared == nil || (highest != nil && roundChange.Prepared.Round.Cmp(highest.View.Round) <= 0) {
			continue
		}
		var cert *bft.PreparedCertificate
		if err := rlp.DecodeBytes(message.Evidence, &cert); err != nil {
			return nil, nil, err
		}
		highest = cert
	}
	if _, _, err := c.verifyRoundChangeCertificate(certificate, target); err != nil {
		return nil, nil, err
	}
	return certificate, highest, nil
}

// ----------------------------------------------------------------------------

func newRoundChangeSet(qualified *valset.AddressSet, maxMessagesPerRound int) *roundChangeSet {
	return &roundChangeSet{
		qualified:           qualified,
		maxMessagesPerRound: maxMessagesPerRound,
		roundChanges:        make(map[uint64]*messageSet),
		mu:                  new(sync.Mutex),
	}
}

type roundChangeSet struct {
	qualified           *valset.AddressSet
	maxMessagesPerRound int
	roundChanges        map[uint64]*messageSet
	mu                  *sync.Mutex
}

// Add retains a ROUND CHANGE only when its round is within the current window.
// A retained round needs no more than a quorum's worth of messages: once that
// limit is reached, any further distinct sender cannot change its outcome.
// Add holds rcs.mu and is the only writer of a retained messageSet, so reading
// that set's size before adding to it cannot race.
func (rcs *roundChangeSet) Add(currentRound, messageRound *big.Int, msg *bft.Message) (int, error) {
	// roundChanges is keyed by uint64, so reject values that would truncate to
	// an existing bucket even if callers bypass checkMessage.
	if !messageRound.IsUint64() {
		return 0, bft.ErrInvalidMessage
	}

	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	maxRound := new(big.Int).Add(currentRound, big.NewInt(maxRoundChangeRoundsAhead))
	if messageRound.Cmp(maxRound) > 0 {
		return 0, errRoundChangeTooFar
	}

	// A bucket is stored only after a message is accepted, so a rejected message
	// never leaves an empty bucket behind.
	round := messageRound.Uint64()
	messages := rcs.roundChanges[round]
	if messages == nil {
		messages = newMessageSet(rcs.qualified)
	}
	if messages.Get(msg.Address) == nil && messages.Size() >= rcs.maxMessagesPerRound {
		return 0, errRoundChangeMessageLimit
	}
	if err := messages.Add(msg); err != nil {
		return 0, err
	}
	rcs.roundChanges[round] = messages
	return messages.Size(), nil
}

// Clear deletes the messages with smaller round
func (rcs *roundChangeSet) Clear(round *big.Int) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	for k, rms := range rcs.roundChanges {
		if len(rms.Values()) == 0 || k < round.Uint64() {
			delete(rcs.roundChanges, k)
		}
	}
}

// MaxRound returns the max round which the number of messages is equal or larger than num
func (rcs *roundChangeSet) MaxRound(num int) *big.Int {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	var maxRound *big.Int
	for k, rms := range rcs.roundChanges {
		if rms.Size() < num {
			continue
		}
		r := big.NewInt(int64(k))
		if maxRound == nil || maxRound.Cmp(r) < 0 {
			maxRound = r
		}
	}
	return maxRound
}

// Values returns the messages collected for one target round.
func (rcs *roundChangeSet) Values(round *big.Int) []*bft.Message {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	set := rcs.roundChanges[round.Uint64()]
	if set == nil {
		return nil
	}
	return set.Values()
}
