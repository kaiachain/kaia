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
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
)

// Round-change recovery after the permissionless fork follows istanbul/core
// (roundchange.go) rule for rule:
//  1. A node that saw a PREPARE/COMMIT quorum for X at round r locks (r, X) and
//     keeps the quorum votes as a PreparedCertificate (lockHash).
//  2. Its ROUND CHANGE signs PreparedClaim{r, X, keccak(Evidence)} and carries
//     the certificate as unsigned Evidence; the hash binds the two.
//  3. Receivers verify the Evidence against the claim and cache the
//     certificate by (sequence, r, X); the round-change set keeps only the
//     stripped envelope.
//  4. The proposer of round r' embeds the stripped ROUND CHANGE quorum and the
//     votes of the highest claim; its proposal must be that claim's value.
//  5. A receiver accepts a justified PRE-PREPARE regardless of its own lock,
//     adopts a strictly newer certificate, and keeps its lock until a PREPARE
//     quorum at the new round replaces it. A receiver without a lock also
//     accepts a fresh proposal below a claim; a locked one first executes its
//     lock once (verifyLockedProposalExecution).
//  6. Proposals and certificate bodies carrying committed seals are rejected.
//  7. A claim for a block known to be bad is void, like a local lock on it.

// istanbulExtraVanity is the vanity prefix of a header's Extra. The value is
// shared with the Istanbul header format (istanbul.IstanbulExtraVanity).
const istanbulExtraVanity = 32

// maxConsensusP2PMessageBytes matches node/cn.ProtocolMaxMsgSize. The P2P
// message wraps the signed consensus payload in ConsensusMsg, so the payload
// itself must be smaller than this cap.
const maxConsensusP2PMessageBytes = 12 * 1024 * 1024

const (
	// A single sender may retain at most one transport-sized allocation across
	// its future ROUND CHANGE Evidence. Each Evidence must match the
	// EvidenceHash its sender signed, so only the sender can spend it.
	maxBacklogEvidenceBytesPerSender uint64 = maxConsensusP2PMessageBytes
	// Bound aggregate retained Evidence across all senders. Messages without
	// Evidence remain admissible when it is full.
	maxBacklogEvidenceBytes uint64 = 64 * 1024 * 1024
)

var (
	errInconsistentPrevHash    = errors.New("message is for another parent")
	errFreshProposalBelowClaim = errors.New("fresh proposal below a prepared claim")
	errProposalNotHighestClaim = errors.New("proposal is not the highest prepared value")
)

// consensusP2PMessageSize is the RLP size of ConsensusMsg{PrevHash, Payload}.
func consensusP2PMessageSize(payload []byte) uint64 {
	return rlp.ListSize(33 + rlp.BytesSize(payload))
}

type preparedEvidenceKey struct {
	sequence uint64
	round    uint64
	digest   common.Hash
}

// verifiedEvidenceKey identifies Evidence bytes already verified for a claim.
type verifiedEvidenceKey struct {
	claim preparedEvidenceKey
	hash  common.Hash
}

func preparedKey(sequence *big.Int, claim *bft.PreparedClaim) (preparedEvidenceKey, error) {
	if sequence == nil || !sequence.IsUint64() || claim == nil || claim.Round == nil || !claim.Round.IsUint64() {
		return preparedEvidenceKey{}, errors.New("prepared claim has invalid view")
	}
	return preparedEvidenceKey{sequence: sequence.Uint64(), round: claim.Round.Uint64(), digest: claim.Digest}, nil
}

// resetPreparedEvidence drops every certificate cache; called on a new height.
func (m *machine) resetPreparedEvidence() {
	m.roundChangeCertificate = nil
	m.preparedBlocks = make(map[common.Hash]*types.Block)
	m.preparedCertificates = make(map[preparedEvidenceKey]*bft.PreparedCertificate)
	m.verifiedEvidence = make(map[verifiedEvidenceKey]struct{})
	m.executedProposals = make(map[common.Hash]struct{})
}

// buildPreparedCertificate collects the PREPARE/COMMIT quorum of the current
// view for the accepted proposal, one vote per sender, as signed.
func (m *machine) buildPreparedCertificate() *bft.PreparedCertificate {
	if m.preprepare == nil {
		return nil
	}
	block, ok := m.preprepare.Proposal.(*types.Block)
	if !ok {
		return nil
	}
	messages := make([]*bft.Message, 0, m.prepares.Size()+m.commits.Size())
	seen := make(map[common.Address]struct{})
	for _, set := range []*messageSet{m.prepares, m.commits} {
		for _, msg := range set.Values() {
			if _, exists := seen[msg.Address]; exists {
				continue
			}
			vote, err := msg.normalized()
			if err != nil {
				return nil
			}
			seen[msg.Address] = struct{}{}
			messages = append(messages, vote)
		}
	}
	return &bft.PreparedCertificate{
		View:     m.currentView(),
		Proposal: block,
		Messages: messages,
	}
}

// adoptPreparedCertificate moves the lock to a strictly newer, independently
// verified certificate. Keeping the comparison with the mutation prevents any
// caller from downgrading a newer local lock.
func (m *machine) adoptPreparedCertificate(cert *bft.PreparedCertificate) {
	if cert == nil || cert.View == nil || cert.View.Round == nil || cert.Proposal == nil {
		return
	}
	if own := m.preparedCertificate; own != nil && own.View != nil && own.View.Round != nil &&
		cert.View.Round.Cmp(own.View.Round) <= 0 {
		return
	}
	m.lockedHash = cert.Proposal.Hash()
	m.preparedCertificate = cert
}

// roundChangeEvidenceMatchesClaim reports whether msg carries Evidence exactly
// when it claims a prepared value, and whether that Evidence hashes to the
// EvidenceHash the sender signed.
func roundChangeEvidenceMatchesClaim(msg *wireMessage) bool {
	if msg.Prepared == nil {
		return len(msg.Evidence) == 0
	}
	return len(msg.Evidence) != 0 && crypto.Keccak256Hash(msg.Evidence) == msg.Prepared.EvidenceHash
}

// verifyRoundChangeEvidence checks that a ROUND CHANGE's unsigned Evidence
// proves exactly the PreparedClaim its sender signed. The caller has already
// checked that the claim and Evidence appear together.
func (m *machine) verifyRoundChangeEvidence(msg *wireMessage, view *bft.View) (*bft.PreparedCertificate, error) {
	claim := msg.Prepared
	if claim == nil {
		return nil, nil
	}
	key, err := preparedKey(view.Sequence, claim)
	if err != nil {
		return nil, err
	}
	// A cached certificate skips verifyPreparedCertificate, so enforce the claim
	// rules of verifyRoundChangeCertificate here.
	if view.Round == nil || claim.Round.Cmp(view.Round) >= 0 || common.EmptyHash(claim.Digest) {
		return nil, errors.New("prepared claim is not before the round-change view")
	}
	verified := verifiedEvidenceKey{claim: key, hash: claim.EvidenceHash}
	if cert := m.preparedCertificates[key]; cert != nil {
		if _, ok := m.verifiedEvidence[verified]; ok {
			return cert, nil
		}
	}
	var cert *bft.PreparedCertificate
	if err := rlp.DecodeBytes(msg.Evidence, &cert); err != nil {
		return nil, err
	}
	if cert == nil || cert.View == nil || cert.View.Round == nil || cert.Proposal == nil ||
		cert.View.Round.Cmp(claim.Round) != 0 || cert.Proposal.Hash() != claim.Digest {
		return nil, errors.New("evidence does not match the prepared claim")
	}
	if err := m.verifyPreparedCertificate(cert, view); err != nil {
		return nil, err
	}
	m.verifiedEvidence[verified] = struct{}{}
	return m.canonicalPreparedCertificate(cert), nil
}

// canonicalPreparedCertificate returns the stored certificate for cert's
// (sequence, round, digest), storing a copy of cert when none is. Each block is
// retained once per height and each vote set once per exact claim.
func (m *machine) canonicalPreparedCertificate(cert *bft.PreparedCertificate) *bft.PreparedCertificate {
	if cert == nil || cert.View == nil || cert.View.Sequence == nil || cert.View.Round == nil || cert.Proposal == nil {
		return cert
	}
	key := preparedEvidenceKey{sequence: cert.View.Sequence.Uint64(), round: cert.View.Round.Uint64(), digest: cert.Proposal.Hash()}
	if existing := m.preparedCertificates[key]; existing != nil {
		return existing
	}
	block := m.preparedBlocks[key.digest]
	if block == nil {
		block = m.localPreparedBlock(key.digest)
	}
	if block == nil {
		block = cert.Proposal
	}
	m.preparedBlocks[key.digest] = block
	canonical := &bft.PreparedCertificate{
		View: &bft.View{
			Sequence: new(big.Int).Set(cert.View.Sequence),
			Round:    new(big.Int).Set(cert.View.Round),
		},
		Proposal: block,
		Messages: slices.Clone(cert.Messages),
	}
	m.preparedCertificates[key] = canonical
	return canonical
}

// localPreparedBlock prefers a body already verified through the normal
// PRE-PREPARE path over an equivalent peer-supplied certificate body.
func (m *machine) localPreparedBlock(digest common.Hash) *types.Block {
	if cert := m.preparedCertificate; cert != nil && cert.Proposal != nil && cert.Proposal.Hash() == digest {
		return cert.Proposal
	}
	if m.preprepare != nil {
		if block, ok := m.preprepare.Proposal.(*types.Block); ok && block.Hash() == digest {
			return block
		}
	}
	return nil
}

// verifyPreparedCertificate authenticates a quorum of votes for one proposal
// in a round strictly before the target ROUND-CHANGE view.
func (m *machine) verifyPreparedCertificate(cert *bft.PreparedCertificate, target *bft.View) error {
	if cert == nil || cert.View == nil || cert.Proposal == nil || target == nil ||
		cert.View.Sequence == nil || cert.View.Round == nil || target.Sequence == nil || target.Round == nil {
		return errors.New("incomplete prepared certificate")
	}
	if cert.View.Sequence.Cmp(target.Sequence) != 0 || cert.View.Round.Cmp(target.Round) >= 0 ||
		cert.Proposal.Number().Cmp(cert.View.Sequence) != 0 {
		return errors.New("prepared certificate has invalid view")
	}
	if err := m.verifyPreparedCertificateVotes(cert); err != nil {
		return err
	}
	// The block hash excludes the committed-seal list, so bodies that differ
	// only there share a digest. A proposal carries none until it is committed.
	padded, err := proposalCarriesCommittedSeals(cert.Proposal.Header())
	if err != nil {
		return fmt.Errorf("prepared certificate proposal has invalid extra: %w", err)
	}
	if padded {
		return errors.New("prepared certificate proposal carries committed seals")
	}
	// Votes bind the header hash, not every body component.
	return m.b.verifyBody(cert.Proposal)
}

// proposalCarriesCommittedSeals reports whether a peer-supplied proposal header
// has a committed-seal list. It decodes the extra directly rather than through
// the sealer, whose parser caches every distinct extra it sees. The layout is
// istanbul.IstanbulExtra.
func proposalCarriesCommittedSeals(header *types.Header) (bool, error) {
	if len(header.Extra) < istanbulExtraVanity {
		return false, errors.New("invalid istanbul header extra-data")
	}
	var extra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := rlp.DecodeBytes(header.Extra[istanbulExtraVanity:], &extra); err != nil {
		return false, err
	}
	return len(extra.CommittedSeal) != 0, nil
}

// verifyPreparedCertificateVotes performs the committee, signature, subject,
// seal and quorum checks after the certificate view has been checked.
func (m *machine) verifyPreparedCertificateVotes(cert *bft.PreparedCertificate) error {
	_, committee, _, _, quorum, _, err := m.getRoundCommitteeState(cert.View.Sequence.Uint64(), cert.View.Round.Uint64())
	if err != nil {
		return err
	}
	// At most one meaningful vote per committee member; reject excess entries
	// before signature recovery.
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
		signer, err := bft.GetSignatureAddress(unsigned, vote.Signature)
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
		if view == nil || view.Sequence == nil || view.Round == nil || view.Cmp(cert.View) != 0 ||
			digest != cert.Proposal.Hash() || vote.PrevHash != cert.Proposal.ParentHash() {
			return errors.New("prepared certificate vote has inconsistent payload")
		}
		if vote.Code == bft.MsgCommit {
			// Certificates exist only where committed seals bind the round.
			committer, err := bft.GetSignatureAddress(m.b.committedSealPreimage(view, digest), committedSeal)
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
// it, or nil when no sender claims one, with the claims themselves.
func (m *machine) verifyRoundChangeCertificate(messages []*bft.Message, target *bft.View) (*big.Int, []*bft.PreparedClaim, error) {
	_, committee, _, _, quorum, _, err := m.getRoundCommitteeState(target.Sequence.Uint64(), target.Round.Uint64())
	if err != nil {
		return nil, nil, err
	}
	if len(messages) > committee.Len() {
		return nil, nil, fmt.Errorf("round-change certificate has %d messages, maximum is %d", len(messages), committee.Len())
	}
	seen := make(map[common.Address]struct{}, len(messages))
	var highest *big.Int
	claims := make([]*bft.PreparedClaim, 0, len(messages))
	for _, message := range messages {
		// Embedded messages must be stripped of their Evidence.
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
		signer, err := bft.GetSignatureAddress(unsigned, message.Signature)
		if err != nil || signer != message.Address {
			return nil, nil, errors.New("round-change certificate contains invalid signature")
		}
		var roundChange *bft.RoundChange
		if err := message.Decode(&roundChange); err != nil {
			return nil, nil, err
		}
		if roundChange.View == nil || roundChange.View.Sequence == nil || roundChange.View.Round == nil ||
			roundChange.View.Cmp(target) != 0 || message.PrevHash != m.parentHash {
			return nil, nil, errors.New("round-change certificate has inconsistent view")
		}
		if claim := roundChange.Prepared; claim != nil {
			if claim.Round == nil || claim.Round.Cmp(target.Round) >= 0 || common.EmptyHash(claim.Digest) {
				return nil, nil, errors.New("round-change certificate has an invalid prepared claim")
			}
			// A claim for a block this node found invalid is void, as a local
			// lock on it is.
			if m.b.hasBadProposal(claim.Digest) {
				seen[message.Address] = struct{}{}
				continue
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
// PRE-PREPARE above round 0. When the quorum claims a prepared value, the
// proposal must be one claimed at the highest round and PreparedMessages must
// prove it; the resulting certificate is returned.
func (m *machine) verifyPreprepareJustification(pp *bft.Preprepare) (*bft.PreparedCertificate, error) {
	highest, claims, err := m.verifyRoundChangeCertificate(pp.RoundChangeCertificate, pp.View)
	if err != nil {
		return nil, err
	}
	if highest == nil {
		if len(pp.PreparedMessages) != 0 {
			return nil, errors.New("prepared votes without a prepared claim")
		}
		return nil, nil
	}
	block, ok := pp.Proposal.(*types.Block)
	if !ok {
		return nil, errors.New("proposal is not a block")
	}
	claimed := slices.ContainsFunc(claims, func(claim *bft.PreparedClaim) bool {
		return claim.Round.Cmp(highest) == 0 && claim.Digest == block.Hash()
	})
	if !claimed {
		if len(pp.PreparedMessages) == 0 {
			return nil, errFreshProposalBelowClaim
		}
		return nil, errProposalNotHighestClaim
	}
	cert := &bft.PreparedCertificate{
		View:     &bft.View{Sequence: new(big.Int).Set(pp.View.Sequence), Round: new(big.Int).Set(highest)},
		Proposal: block,
		Messages: pp.PreparedMessages,
	}
	if err := m.verifyPreparedCertificateVotes(cert); err != nil {
		return nil, err
	}
	return cert, nil
}

// verifyRecoveryJustification judges the justification of a PRE-PREPARE above
// round 0 from this node's point of view. A proposal below the highest claim is
// accepted by a node without a lock; a locked node first executes its lock once
// and refuses the proposal only if the lock survives.
func (m *machine) verifyRecoveryJustification(pp *bft.Preprepare) (*bft.PreparedCertificate, error) {
	highest, err := m.verifyPreprepareJustification(pp)
	if err != errFreshProposalBelowClaim && err != errProposalNotHighestClaim {
		return highest, err
	}
	if m.isHashLocked() {
		m.verifyLockedProposalExecution()
		if m.isHashLocked() {
			return nil, err
		}
		if highest, err = m.verifyPreprepareJustification(pp); err != errFreshProposalBelowClaim && err != errProposalNotHighestClaim {
			return highest, err
		}
	}
	logger.Debug("Accept proposal below a prepared claim without a lock", "view", pp.View)
	return nil, nil
}

// verifyLockedProposalExecution executes the locked block once per height, so
// that a lock on a block this node never imported cannot refuse recovery from
// a value known bad elsewhere. A failure is recorded by the chain and read
// back through isHashLocked.
func (m *machine) verifyLockedProposalExecution() {
	locked := m.lockedHash
	var block *types.Block
	if cert := m.preparedCertificate; cert != nil && cert.Proposal != nil && cert.Proposal.Hash() == locked {
		block = cert.Proposal
	} else if p, ok := m.proposal().(*types.Block); ok && p.Hash() == locked {
		block = p
	} else {
		return
	}
	if _, done := m.executedProposals[locked]; done {
		return
	}
	err := m.b.chain.ValidateProposalExecution(block)
	if err == nil || m.b.hasBadProposal(locked) {
		// A definite result; a failure to execute at all is retried later.
		m.executedProposals[locked] = struct{}{}
	}
	if err != nil {
		logger.Warn("Locked proposal fails execution", "hash", locked, "released", !m.isHashLocked(), "err", err)
	}
}

// roundChangeJustification selects what a proposer attaches to a PRE-PREPARE
// for the given view: the stripped signed ROUND CHANGE quorum, plus the
// certificate behind its highest prepared claim. Every message was verified
// with its Evidence when it entered the round-change set, which cached the
// certificate by claim.
func (m *machine) roundChangeJustification(messages []*bft.Message, target *bft.View) ([]*bft.Message, *bft.PreparedCertificate, error) {
	certificate := make([]*bft.Message, 0, len(messages)+1)
	var highest *bft.PreparedCertificate
	for _, message := range messages {
		certificate = append(certificate, message.WithoutEvidence())
		var roundChange *bft.RoundChange
		if err := message.Decode(&roundChange); err != nil {
			return nil, nil, err
		}
		if roundChange.Prepared == nil || m.b.hasBadProposal(roundChange.Prepared.Digest) ||
			(highest != nil && roundChange.Prepared.Round.Cmp(highest.View.Round) <= 0) {
			continue
		}
		key, err := preparedKey(target.Sequence, roundChange.Prepared)
		if err != nil {
			return nil, nil, err
		}
		cert := m.preparedCertificates[key]
		if cert == nil {
			return nil, nil, errors.New("prepared evidence is unavailable")
		}
		highest = cert
	}
	// A proposer locked at a newer round than every claim in the quorum must
	// still re-propose its lock. Its own signed ROUND CHANGE carrying that claim
	// justifies the lock.
	if own := m.preparedCertificate; own != nil && m.isHashLocked() && own.View.Sequence.Cmp(target.Sequence) == 0 &&
		own.View.Round.Cmp(target.Round) < 0 && (highest == nil || own.View.Round.Cmp(highest.View.Round) > 0) {
		own = m.canonicalPreparedCertificate(own)
		ownRoundChange, err := m.signedRoundChange(target, own)
		if err != nil {
			return nil, nil, err
		}
		replaced := false
		for i, message := range certificate {
			if message.Address == ownRoundChange.Address {
				certificate[i], replaced = ownRoundChange, true
				break
			}
		}
		if !replaced {
			certificate = append(certificate, ownRoundChange)
		}
		highest = own
	}
	if _, _, err := m.verifyRoundChangeCertificate(certificate, target); err != nil {
		return nil, nil, err
	}
	return certificate, highest, nil
}

// signedRoundChange signs this node's ROUND CHANGE for target claiming cert,
// without Evidence, for embedding in a PRE-PREPARE justification.
func (m *machine) signedRoundChange(target *bft.View, cert *bft.PreparedCertificate) (*bft.Message, error) {
	evidence, err := bft.Encode(cert)
	if err != nil {
		return nil, err
	}
	payload, err := bft.Encode(&bft.RoundChange{
		View: &bft.View{Sequence: new(big.Int).Set(target.Sequence), Round: new(big.Int).Set(target.Round)},
		Prepared: &bft.PreparedClaim{
			Round:        new(big.Int).Set(cert.View.Round),
			Digest:       cert.Proposal.Hash(),
			EvidenceHash: crypto.Keccak256Hash(evidence),
		},
	})
	if err != nil {
		return nil, err
	}
	msg := &bft.Message{PrevHash: m.parentHash, Code: bft.MsgRoundChange, Msg: payload, Address: m.b.address}
	unsigned, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	if msg.Signature, err = m.b.sign(unsigned); err != nil {
		return nil, err
	}
	return msg, nil
}

// checkRoundChange performs the cheap, non-mutating admission checks of
// addRoundChange before the attached certificate is decoded and verified.
func (m *machine) checkRoundChange(round *big.Int, msg *wireMessage) error {
	if !round.IsUint64() {
		return bft.ErrInvalidMessage
	}
	m.roundChangesMu.Lock()
	defer m.roundChangesMu.Unlock()
	if round.Cmp(new(big.Int).Add(m.round, big.NewInt(maxRoundChangeRoundsAhead))) > 0 {
		return errRoundChangeTooFar
	}
	messages := m.roundChangeSets[round.Uint64()]
	if messages == nil {
		return nil
	}
	if existing := messages.Get(msg.Address); existing != nil {
		if !replacesRoundChange(existing, msg) {
			return errIgnored
		}
		return nil
	}
	if messages.Size() >= m.requiredMessageCount {
		return errRoundChangeLimit
	}
	return nil
}

// replacesRoundChange reports whether next may replace prev, both from the
// same sender for the same round. Only a strictly higher prepared claim may,
// so a sender cannot force its Evidence to be verified again by alternating
// signed messages.
func replacesRoundChange(prev, next *wireMessage) bool {
	if next.Prepared == nil || next.Prepared.Round == nil {
		return false
	}
	return prev.Prepared == nil || prev.Prepared.Round == nil || next.Prepared.Round.Cmp(prev.Prepared.Round) > 0
}

// sameRoundChange reports whether a and b carry the same signed content.
func sameRoundChange(a, b *wireMessage) bool {
	if !bytes.Equal(a.Msg, b.Msg) {
		return false
	}
	if a.Prepared == nil || b.Prepared == nil {
		return a.Prepared == nil && b.Prepared == nil
	}
	return a.Prepared.Round != nil && b.Prepared.Round != nil && a.Prepared.Round.Cmp(b.Prepared.Round) == 0 &&
		a.Prepared.Digest == b.Prepared.Digest && a.Prepared.EvidenceHash == b.Prepared.EvidenceHash
}

// roundChangeMessages returns the signed ROUND CHANGEs retained for a round,
// in the shared message form.
func (m *machine) roundChangeMessages(round *big.Int) []*bft.Message {
	m.roundChangesMu.Lock()
	set := m.roundChangeSets[round.Uint64()]
	m.roundChangesMu.Unlock()
	if set == nil {
		return nil
	}
	values := set.Values()
	messages := make([]*bft.Message, 0, len(values))
	for _, v := range values {
		msg, err := v.normalized()
		if err != nil {
			continue
		}
		messages = append(messages, msg.WithoutEvidence())
	}
	return messages
}
