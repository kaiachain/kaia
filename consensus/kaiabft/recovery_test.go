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
	"math/big"
	"testing"
	"testing/synctest"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/require"
)

// splitLock reproduces the rc.3 split-lock counterexample at height 1: the
// round-0 proposer A (node 0) equivocates, sending Y to the next proposer B
// (node 1) and X to C/D (nodes 2, 3). A votes only for X and withholds every
// COMMIT, so C/D lock X while it stays undecided and B never locks.
func splitLock(t *testing.T, s *recoveryNet) (x *types.Block) {
	a := s.nodes[0]
	x = a.proposal(1)
	y := a.proposal(2)
	require.NotEqual(t, x.Hash(), y.Hash())
	s.equivocate(0, []int{1}, viewAt(1, 0), y)
	s.drop(bft.MsgPrepare, []int{0}, []int{1})
	s.drop(bft.MsgCommit, []int{0}, s.all())
	// A builds X itself from the request; the others' requests are not proposed.
	s.enqueue(netDelivery{kind: deliverRequest, from: 0, to: 0, request: x})
	for _, id := range []int{1, 2, 3} {
		s.enqueue(netDelivery{kind: deliverRequest, from: id, to: id, request: s.nodes[id].proposal(1)})
	}
	s.drain()
	for _, n := range s.nodes {
		n.requireUncommitted(1)
	}
	return x
}

// rc.3 TestConsensusSplitLockPreparedCertificate on the kaiabft machine: B must
// learn X from the certificates C/D attach to their ROUND CHANGEs and re-propose
// it although B never locked it.
func TestRecoverySplitLockPreparedCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		x := splitLock(t, s)
		b := s.nodes[1]
		for _, id := range []int{2, 3} {
			n := s.nodes[id]
			require.Equal(t, x.Hash(), n.m.lockedHash)
			cert := n.m.preparedCertificate
			require.NotNil(t, cert)
			require.Equal(t, int64(0), cert.View.Round.Int64())
			require.GreaterOrEqual(t, len(cert.Messages), 3, "the certificate carries a quorum of signed votes")
			require.NoError(t, b.m.verifyPreparedCertificate(cert, viewAt(1, 1)))
		}
		require.True(t, common.EmptyHash(b.m.lockedHash))

		// A's ROUND CHANGEs are lost; B, C and D time out.
		s.drop(bft.MsgRoundChange, []int{0}, s.all())
		s.timeout(1, 2, 3)

		// B re-proposed X with a justification proving X@0.
		pp := s.sentPreprepare(1, viewAt(1, 1))
		require.NotNil(t, pp)
		require.Equal(t, x.Hash(), pp.Proposal.Hash())
		require.GreaterOrEqual(t, len(pp.RoundChangeCertificate), 3)
		require.GreaterOrEqual(t, len(pp.PreparedMessages), 3)
		claims := 0
		for _, msg := range pp.RoundChangeCertificate {
			require.Empty(t, msg.Evidence, "embedded ROUND CHANGEs are stripped")
			var rc *bft.RoundChange
			require.NoError(t, msg.Decode(&rc))
			if rc.Prepared != nil {
				require.Equal(t, x.Hash(), rc.Prepared.Digest)
				require.Equal(t, int64(0), rc.Prepared.Round.Int64())
				claims++
			}
		}
		require.Equal(t, 2, claims, "C and D claim X@0")
		for _, id := range []int{1, 2, 3} {
			require.Equal(t, x.Hash(), s.nodes[id].requireCommitted(1, 1).Hash())
		}
	})
}

// rc.3 TestConsensusSplitLockReproposalRefreshesPreparedCertificate: a locked
// node re-establishes its certificate in the recovery round before its COMMIT,
// and B, which adopted X@0 from the justification, refreshes to X@1 as well.
func TestRecoverySplitLockReproposalRefreshesPreparedCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		x := splitLock(t, s)
		s.drop(bft.MsgRoundChange, []int{0}, s.all())
		// Withhold every round-1 COMMIT so the evidence is inspected before H1.
		s.drop(bft.MsgCommit, s.all(), s.all(), 1)
		s.timeout(1, 2, 3)
		for _, id := range []int{1, 2, 3} {
			n := s.nodes[id]
			n.requireUncommitted(1)
			require.Equal(t, statePrepared, n.m.state)
			cert := n.m.preparedCertificate
			require.NotNil(t, cert, "node %d", id)
			require.Equal(t, int64(1), cert.View.Round.Int64(), "node %d", id)
			require.Equal(t, x.Hash(), cert.Proposal.Hash())
			require.NoError(t, s.nodes[0].m.verifyPreparedCertificate(cert, viewAt(1, 2)))
		}
	})
}

// A quorum locks X at round 0 without committing. Before the fork the round-1
// proposer re-proposes its hash lock and the ROUND CHANGEs stay in the legacy
// form; after it the proposer re-proposes through the certificates its ROUND
// CHANGE quorum claims. Both commit X at round 1.
func TestRecoveryLockedProposalPreAndPostFork(t *testing.T) {
	for _, tc := range []struct {
		name           string
		cfg            *params.ChainConfig
		permissionless bool
	}{
		{"legacy", legacyConfig(), false},
		{"permissionless", permissionlessConfig(), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newRecoveryNet(t, 4, tc.cfg)
				s.drop(bft.MsgCommit, s.all(), s.all(), 0)
				s.advance(1)
				x := s.nodes[0].m.preprepare.Proposal
				for _, n := range s.nodes {
					n.requireUncommitted(1)
					require.Equal(t, x.Hash(), n.m.lockedHash)
					require.Equal(t, tc.permissionless, n.m.preparedCertificate != nil)
				}
				s.timeout(0, 1, 2, 3)
				rcs := s.sentRoundChanges()
				require.Len(t, rcs, 4)
				for _, msg := range rcs {
					var rc *bft.RoundChange
					require.NoError(t, msg.Decode(&rc))
					if tc.permissionless {
						require.NotNil(t, rc.Prepared, "every locked node claims its lock")
						require.Equal(t, x.Hash(), rc.Prepared.Digest)
						require.NotEmpty(t, msg.Evidence)
					} else {
						require.Nil(t, rc.Prepared)
						require.Empty(t, msg.Evidence)
					}
				}
				for _, n := range s.nodes {
					require.Equal(t, x.Hash(), n.requireCommitted(1, 1).Hash())
				}
			})
		})
	}
}

// Across the fork boundary: height 1 is legacy and height 2 permissionless.
// The recovery round at height 2 carries prepared claims and commits with
// round-bound seals; height 1 is untouched.
func TestRecoveryForkBoundary(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cfg := &params.ChainConfig{PermissionlessCompatibleBlock: big.NewInt(2)}
		s := newRecoveryNet(t, 4, cfg)
		s.advance(1)
		for _, n := range s.nodes {
			n.requireCommitted(1, 0)
			require.Equal(t, uint64(2), n.m.sequence.Uint64())
		}
		s.drop(bft.MsgCommit, s.all(), s.all(), 0)
		s.advance(2)
		for _, n := range s.nodes {
			n.requireUncommitted(2)
			require.NotNil(t, n.m.preparedCertificate)
		}
		s.timeout(0, 1, 2, 3)
		for _, n := range s.nodes {
			n.requireCommitted(2, 1)
		}
	})
}

// Round 0 of a permissionless height commits as before; no justification is
// attached and no recovery state is created.
func TestRecoveryRoundZeroUnchanged(t *testing.T) {
	for _, cfg := range []*params.ChainConfig{legacyConfig(), permissionlessConfig()} {
		synctest.Test(t, func(t *testing.T) {
			s := newRecoveryNet(t, 4, cfg)
			s.advance(1)
			for _, n := range s.nodes {
				n.requireCommitted(1, 0)
				require.Nil(t, n.m.roundChangeCertificate)
			}
			s.advance(2)
			for _, n := range s.nodes {
				n.requireCommitted(2, 0)
			}
		})
	}
}

// ROUND CHANGE admission: Evidence must be bound to a signed claim and prove
// it. Every rejection leaves the sender out of the round-change set.
func TestRecoveryRoundChangeEvidenceValidation(t *testing.T) {
	type mutate func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte
	resign := func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
		claim, evidence := claimFor(t, cert)
		return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, claim, evidence)
	}
	for _, tc := range []struct {
		name string
		make mutate
		want error
	}{
		{"accepts the genuine certificate", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			return resign(t, s, cert)
		}, nil},
		{"evidence other than the signed hash", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			claim, _ := claimFor(t, cert)
			other := *cert
			other.Messages = cert.Messages[:len(cert.Messages)-1]
			return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, claim, mustEncode(t, &other))
		}, bft.ErrInvalidMessage},
		{"evidence without a claim", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, nil, mustEncode(t, cert))
		}, bft.ErrInvalidMessage},
		{"claim without evidence", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			claim, _ := claimFor(t, cert)
			return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, claim, nil)
		}, bft.ErrInvalidMessage},
		{"claim of another digest", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			claim, evidence := claimFor(t, cert)
			claim.Digest = common.Hash{1}
			return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, claim, evidence)
		}, bft.ErrInvalidMessage},
		{"claim at the round-change round", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			claim, evidence := claimFor(t, cert)
			claim.Round = big.NewInt(1)
			return s.nodes[2].signRoundChange(viewAt(1, 1), s.nodes[2].m.parentHash, claim, evidence)
		}, bft.ErrInvalidMessage},
		{"votes below quorum", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			short := *cert
			short.Messages = cert.Messages[:2]
			return resign(t, s, &short)
		}, bft.ErrInvalidMessage},
		{"duplicate voter", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			dup := *cert
			dup.Messages = []*bft.Message{cert.Messages[0], cert.Messages[0], cert.Messages[1]}
			return resign(t, s, &dup)
		}, bft.ErrInvalidMessage},
		{"forged vote signature", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			forged := *cert
			forged.Messages = append([]*bft.Message(nil), cert.Messages...)
			vote := *forged.Messages[0]
			vote.Signature = forged.Messages[1].Signature
			forged.Messages[0] = &vote
			return resign(t, s, &forged)
		}, bft.ErrInvalidMessage},
		{"vote for another view", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			moved := *cert
			moved.Messages = append([]*bft.Message(nil), cert.Messages...)
			moved.Messages[0] = s.nodes[0].signVote(bft.MsgPrepare, viewAt(1, 1), cert.Proposal)
			return resign(t, s, &moved)
		}, bft.ErrInvalidMessage},
		{"proposal padded with committed seals", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			header := cert.Proposal.Header()
			require.NoError(t, s.nodes[0].b.sealer.WriteCommittedSeals(header, [][]byte{make([]byte, 65)}))
			padded := *cert
			padded.Proposal = cert.Proposal.WithSeal(header)
			require.Equal(t, cert.Proposal.Hash(), padded.Proposal.Hash(), "the block hash excludes committed seals")
			return resign(t, s, &padded)
		}, bft.ErrInvalidMessage},
		{"another parent", func(t *testing.T, s *recoveryNet, cert *bft.PreparedCertificate) []byte {
			claim, evidence := claimFor(t, cert)
			return s.nodes[2].signRoundChange(viewAt(1, 1), common.Hash{0xba, 0xd}, claim, evidence)
		}, errInconsistentPrevHash},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newRecoveryNet(t, 4, permissionlessConfig())
				splitLock(t, s)
				cert := s.nodes[2].m.preparedCertificate
				require.NotNil(t, cert)
				b := s.nodes[1]
				payload := tc.make(t, s, cert)
				err := b.m.handleMsg(payload)
				if tc.want == nil {
					// A retained ROUND CHANGE for a round above the current one
					// without a quorum reports errIgnored, as before.
					require.ErrorIs(t, err, errIgnored)
					require.NotNil(t, b.m.roundChangeSets[1].Get(s.addrs[2]))
					key, kerr := preparedKey(big.NewInt(1), &bft.PreparedClaim{Round: big.NewInt(0), Digest: cert.Proposal.Hash()})
					require.NoError(t, kerr)
					require.NotNil(t, b.m.preparedCertificates[key], "a verified certificate is cached by claim")
					require.Nil(t, b.m.roundChangeSets[1].Get(s.addrs[2]).Evidence, "only the stripped envelope is retained")
					return
				}
				require.ErrorIs(t, err, tc.want)
				if set := b.m.roundChangeSets[1]; set != nil {
					require.Nil(t, set.Get(s.addrs[2]))
				}
			})
		})
	}
}

// A claim or Evidence cannot enter a legacy height, through the codec or the
// machine.
func TestRecoveryRoundChangeRejectsEvidenceBeforeFork(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, legacyConfig())
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: common.Hash{1}, EvidenceHash: crypto.Keccak256Hash([]byte{1})}
		msg := &wireMessage{Hash: s.nodes[2].m.parentHash, Code: bft.MsgRoundChange,
			Msg: mustEncode(t, &bft.Subject{View: viewAt(1, 1), PrevHash: s.nodes[2].m.parentHash}), Prepared: claim, Evidence: []byte{1}}
		require.Nil(t, s.nodes[2].m.signPayload(msg), "the legacy codec cannot carry a claim")
		msg.Address = s.addrs[2]
		require.ErrorIs(t, s.nodes[1].m.handleRoundChange(msg, s.addrs[2]), bft.ErrInvalidMessage)
	})
}

// A ROUND CHANGE in a non-canonical encoding would be relayed, or embedded in a
// justification, with bytes its sender never signed.
func TestRecoveryRejectsNonCanonicalRoundChange(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		c := s.nodes[2]
		// RoundChange{View, Prepared: nil} with the optional claim written as an
		// explicit empty list, which decodes to the same value.
		view := mustEncode(t, viewAt(1, 1))
		inner, err := rlp.EncodeToBytes([]rlp.RawValue{view, {0xc0}})
		require.NoError(t, err)
		var decoded bft.RoundChange
		require.NoError(t, rlp.DecodeBytes(inner, &decoded))
		require.Nil(t, decoded.Prepared)
		msg := &bft.Message{PrevHash: c.m.parentHash, Code: bft.MsgRoundChange, Msg: inner, Address: c.b.address}
		unsigned, err := msg.PayloadNoSig()
		require.NoError(t, err)
		msg.Signature, err = c.b.sign(unsigned)
		require.NoError(t, err)
		payload, err := msg.Payload()
		require.NoError(t, err)
		var local wireMessage
		require.ErrorIs(t, local.FromPayloadForFork(payload, s.isPermissionlessAt, bft.GetSignatureAddress), errNonCanonicalMessage)
	})
}

// recoveryRound puts B, C and D in round 1 with B's justified PRE-PREPARE
// withheld from C and D, and returns X and B's justification.
func recoveryRound(t *testing.T, s *recoveryNet) (x *types.Block, rcc []*bft.Message, votes []*bft.Message) {
	x = splitLock(t, s)
	s.drop(bft.MsgRoundChange, []int{0}, s.all())
	s.drop(bft.MsgPreprepare, []int{1}, []int{0, 2, 3}, 1)
	s.timeout(1, 2, 3)
	b := s.nodes[1]
	require.Equal(t, int64(1), b.m.round.Int64())
	require.NotNil(t, b.m.preprepare)
	for _, id := range []int{2, 3} {
		n := s.nodes[id]
		require.Equal(t, int64(1), n.m.round.Int64())
		require.Equal(t, stateAcceptRequest, n.m.state)
		require.False(t, n.m.waitingForRoundChange)
	}
	return x, b.m.preprepare.RoundChangeCertificate, b.m.preprepare.PreparedMessages
}

// PRE-PREPARE justification: a locked receiver accepts only the value of the
// highest claim with its proof; an unjustified or malformed proposal starts the
// next round.
func TestRecoveryPreprepareJustification(t *testing.T) {
	for _, tc := range []struct {
		name string
		make func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare
		ok   bool
	}{
		{"justified re-proposal", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: x, RoundChangeCertificate: rcc, PreparedMessages: votes}
		}, true},
		{"no justification", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: x}
		}, false},
		{"round-change quorum below quorum", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: x, RoundChangeCertificate: rcc[:2], PreparedMessages: votes}
		}, false},
		{"round-change message with evidence", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			with := append([]*bft.Message(nil), rcc...)
			m := *with[0]
			m.Evidence = []byte{1}
			with[0] = &m
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: x, RoundChangeCertificate: with, PreparedMessages: votes}
		}, false},
		{"claimed value without its votes", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: x, RoundChangeCertificate: rcc}
		}, false},
		{"fresh value below the claim", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: s.nodes[1].proposal(3), RoundChangeCertificate: rcc}
		}, false},
		{"other value with the claim's votes", func(t *testing.T, s *recoveryNet, x *types.Block, rcc, votes []*bft.Message) *bft.Preprepare {
			return &bft.Preprepare{View: viewAt(1, 1), Proposal: s.nodes[1].proposal(3), RoundChangeCertificate: rcc, PreparedMessages: votes}
		}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newRecoveryNet(t, 4, permissionlessConfig())
				x, rcc, votes := recoveryRound(t, s)
				require.GreaterOrEqual(t, len(rcc), 3)
				require.GreaterOrEqual(t, len(votes), 3)
				c := s.nodes[2]
				payload := s.nodes[1].signPreprepare(tc.make(t, s, x, rcc, votes))
				err := c.m.handleMsg(payload)
				if tc.ok {
					require.NoError(t, err)
					require.Equal(t, statePreprepared, c.m.state)
					return
				}
				require.ErrorIs(t, err, bft.ErrInvalidMessage)
				require.Equal(t, x.Hash(), c.m.lockedHash, "the lock survives a refused proposal")
				require.Equal(t, int64(2), c.m.round.Int64(), "a refused proposal starts the next round")
				require.True(t, c.m.waitingForRoundChange)
			})
		})
	}
}

// A PRE-PREPARE at round 0 carries no justification.
func TestRecoveryRoundZeroPreprepareRejectsJustification(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		a := s.nodes[0]
		x := a.proposal(1)
		rc := s.nodes[1].signRoundChange(viewAt(1, 1), s.nodes[1].m.parentHash, nil, nil)
		var embedded bft.Message
		require.NoError(t, embedded.FromPayload(rc, nil))
		payload := a.signPreprepare(&bft.Preprepare{View: viewAt(1, 0), Proposal: x, RoundChangeCertificate: []*bft.Message{&embedded}})
		require.ErrorIs(t, s.nodes[1].m.handleMsg(payload), bft.ErrInvalidMessage)
		require.Equal(t, int64(1), s.nodes[1].m.round.Int64())
	})
}

// A locked node that refuses a proposal below the highest claim first executes
// its lock once. A lock that fails execution is released, and the proposal is
// accepted as it would be by a node without a lock.
func TestRecoveryLockedNodeExecutesLockBeforeRefusing(t *testing.T) {
	for _, bad := range []bool{false, true} {
		synctest.Test(t, func(t *testing.T) {
			s := newRecoveryNet(t, 4, permissionlessConfig())
			x, rcc, _ := recoveryRound(t, s)
			c := s.nodes[2]
			if bad {
				c.chain.badExec[x.Hash()] = true
			}
			fresh := &bft.Preprepare{View: viewAt(1, 1), Proposal: s.nodes[1].proposal(3), RoundChangeCertificate: rcc}
			err := c.m.handleMsg(s.nodes[1].signPreprepare(fresh))
			require.Equal(t, 1, c.executions)
			if bad {
				require.NoError(t, err)
				require.False(t, c.m.isHashLocked())
				require.Equal(t, statePreprepared, c.m.state)
				return
			}
			require.ErrorIs(t, err, bft.ErrInvalidMessage)
			require.True(t, c.m.isHashLocked())
		})
	}
}

// The highest claim wins: with X claimed at round 0 and Z at round 1 in the
// same quorum, only Z with Z@1 votes is justified.
func TestRecoveryHighestClaimSelection(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		x, _, xVotes := recoveryRound(t, s)
		z := s.nodes[1].proposal(5)
		zView := viewAt(1, 1)
		var zVotes []*bft.Message
		for _, id := range []int{0, 1, 3} {
			zVotes = append(zVotes, s.nodes[id].signVote(bft.MsgPrepare, zView, z))
		}
		zCert := &bft.PreparedCertificate{View: zView, Proposal: z, Messages: zVotes}
		xCert := &bft.PreparedCertificate{View: viewAt(1, 0), Proposal: x, Messages: xVotes}
		target := viewAt(1, 2)
		parent := s.nodes[0].m.parentHash
		var rcc []*bft.Message
		for id, cert := range map[int]*bft.PreparedCertificate{0: xCert, 2: xCert, 3: zCert} {
			claim, _ := claimFor(t, cert)
			var msg bft.Message
			require.NoError(t, msg.FromPayload(s.nodes[id].signRoundChange(target, parent, claim, nil), nil))
			rcc = append(rcc, &msg)
		}
		judge := s.nodes[1].m // B: unlocked would accept anything below; use the strict check
		highest, _, err := judge.verifyRoundChangeCertificate(rcc, target)
		require.NoError(t, err)
		require.Equal(t, int64(1), highest.Int64())

		_, err = judge.verifyPreprepareJustification(&bft.Preprepare{View: target, Proposal: x, RoundChangeCertificate: rcc, PreparedMessages: xVotes})
		require.ErrorIs(t, err, errProposalNotHighestClaim)
		cert, err := judge.verifyPreprepareJustification(&bft.Preprepare{View: target, Proposal: z, RoundChangeCertificate: rcc, PreparedMessages: zVotes})
		require.NoError(t, err)
		require.Equal(t, z.Hash(), cert.Proposal.Hash())
		require.Equal(t, int64(1), cert.View.Round.Int64())

		// Adoption only moves forward.
		judge.adoptPreparedCertificate(cert)
		require.Equal(t, z.Hash(), judge.lockedHash)
		judge.adoptPreparedCertificate(xCert)
		require.Equal(t, z.Hash(), judge.lockedHash, "an older certificate never replaces a newer lock")

		// A proposer locked newer than every claim of its quorum re-proposes its
		// lock, justified by its own signed ROUND CHANGE (rc.3
		// TestConsensusProposerReproposesNewerOwnLock).
		var xOnly []*bft.Message
		for _, id := range []int{0, 2, 3} {
			claim, _ := claimFor(t, xCert)
			var msg bft.Message
			require.NoError(t, msg.FromPayload(s.nodes[id].signRoundChange(target, parent, claim, nil), nil))
			xOnly = append(xOnly, &msg)
		}
		justification, prepared, err := judge.roundChangeJustification(xOnly, target)
		require.NoError(t, err)
		require.Equal(t, z.Hash(), prepared.Proposal.Hash())
		require.Len(t, justification, 4)
		var own *bft.RoundChange
		require.NoError(t, justification[3].Decode(&own))
		require.Equal(t, s.addrs[1], justification[3].Address)
		require.Equal(t, int64(1), own.Prepared.Round.Int64())
		_, err = judge.verifyPreprepareJustification(&bft.Preprepare{View: target, Proposal: z, RoundChangeCertificate: justification, PreparedMessages: prepared.Messages})
		require.NoError(t, err)

		// A claim for a block this node found bad is void.
		s.nodes[1].chain.bad[z.Hash()] = true
		highest, _, err = judge.verifyRoundChangeCertificate(rcc, target)
		require.NoError(t, err)
		require.Equal(t, int64(0), highest.Int64(), "the void Z@1 claim no longer wins")
	})
}

// A sender's ROUND CHANGE is replaced only by one with a strictly higher claim,
// so alternating signed variants cannot force re-verification.
func TestRecoveryRoundChangeReplacement(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		splitLock(t, s)
		b, c := s.nodes[1], s.nodes[2]
		retained := func() *wireMessage { return b.m.roundChangeSets[1].Get(c.b.address) }
		bare := c.signRoundChange(viewAt(1, 1), c.m.parentHash, nil, nil)
		// Round 1 is above B's round, so a retained message reports errIgnored.
		require.ErrorIs(t, b.m.handleMsg(bare), errIgnored)
		require.NotNil(t, retained())
		require.Nil(t, retained().Prepared)
		claim, evidence := claimFor(t, c.m.preparedCertificate)
		withClaim := c.signRoundChange(viewAt(1, 1), c.m.parentHash, claim, evidence)
		require.ErrorIs(t, b.m.handleMsg(withClaim), errIgnored)
		require.NotNil(t, retained().Prepared, "a claim replaces a bare ROUND CHANGE")
		verified := len(b.m.verifiedEvidence)
		require.ErrorIs(t, b.m.checkRoundChange(big.NewInt(1), retained()), errIgnored, "an identical copy is ignored")
		require.ErrorIs(t, b.m.handleMsg(bare), errIgnored)
		require.NotNil(t, retained().Prepared, "a bare ROUND CHANGE cannot replace a claim")
		require.ErrorIs(t, b.m.handleMsg(withClaim), errIgnored)
		require.Equal(t, verified, len(b.m.verifiedEvidence))
	})
}

// Future ROUND CHANGE Evidence is charged to its sender's backlog budget and
// released on delivery.
func TestRecoveryBacklogEvidenceAccounting(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		splitLock(t, s)
		b, c := s.nodes[1], s.nodes[2]
		claim, evidence := claimFor(t, c.m.preparedCertificate)
		msg := &wireMessage{Hash: c.m.parentHash, Code: bft.MsgRoundChange,
			Msg: mustEncode(t, &bft.Subject{View: viewAt(2, 1), PrevHash: c.m.parentHash}), Address: c.b.address, Prepared: claim, Evidence: evidence}
		b.m.storeBacklog(msg, c.b.address)
		require.Equal(t, uint64(len(evidence)), b.m.backlogEvidenceBytes[c.b.address])
		require.Equal(t, uint64(len(evidence)), b.m.backlogTotalEvidenceBytes)

		huge := *msg
		huge.Evidence = make([]byte, maxBacklogEvidenceBytesPerSender)
		b.m.storeBacklog(&huge, c.b.address)
		require.Equal(t, 1, b.m.backlogCounts[c.b.address], "the sender's evidence budget is exhausted")

		b.m.sequence = big.NewInt(3) // the retained message is now old
		b.m.processBacklog()
		require.Zero(t, b.m.backlogTotalEvidenceBytes)
		require.NotContains(t, b.m.backlogEvidenceBytes, c.b.address)
		b.m.sequence = big.NewInt(1)
	})
}

// The node's own claim survives its relay unchanged: the payload a peer
// re-encodes from the internal form is the payload the sender signed.
func TestRecoveryRoundChangeRelayPreservesSignedBytes(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		splitLock(t, s)
		c := s.nodes[2]
		claim, evidence := claimFor(t, c.m.preparedCertificate)
		payload := c.signRoundChange(viewAt(1, 1), c.m.parentHash, claim, evidence)
		var local wireMessage
		require.NoError(t, local.FromPayloadForFork(payload, s.isPermissionlessAt, bft.GetSignatureAddress))
		relayed, err := local.PayloadForFork(s.isPermissionlessAt)
		require.NoError(t, err)
		require.Equal(t, payload, relayed)
		var shared bft.Message
		require.NoError(t, shared.FromPayloadForFork(relayed, s.isPermissionlessAt, bft.GetSignatureAddress))
		require.Equal(t, evidence, shared.Evidence)
	})
}

// Unavailable parent state is not evidence that a certified block is invalid.
// Preserve the lock on transient execution failure and retry on the next proposal.
func TestRecoveryTransientExecutionFailureRetries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		x, rcc, _ := recoveryRound(t, s)
		c := s.nodes[2]
		fresh := &bft.Preprepare{View: viewAt(1, 1), Proposal: s.nodes[1].proposal(3), RoundChangeCertificate: rcc}
		c.chain.executionError = errors.New("parent state temporarily unavailable")
		_, err := c.m.verifyRecoveryJustification(fresh)
		require.ErrorIs(t, err, errFreshProposalBelowClaim)
		require.Equal(t, 1, c.executions)
		require.True(t, c.m.isHashLocked())
		require.NotContains(t, c.m.executedProposals, x.Hash())
		require.False(t, c.chain.HasBadBlock(x.Hash()))
		c.chain.executionError = nil
		_, err = c.m.verifyRecoveryJustification(fresh)
		require.ErrorIs(t, err, errFreshProposalBelowClaim)
		require.Equal(t, 2, c.executions)
		require.Contains(t, c.m.executedProposals, x.Hash())
		require.True(t, c.m.isHashLocked())
		_, err = c.m.verifyRecoveryJustification(fresh)
		require.ErrorIs(t, err, errFreshProposalBelowClaim)
		require.Equal(t, 2, c.executions, "a definite successful execution is reused")
	})
}

// Votes are checked against their own round, not the receiving recovery round.
func TestRecoveryCertificateUsesVotingRoundCommittee(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newRecoveryNet(t, 4, permissionlessConfig())
		splitLock(t, s)
		cert := s.nodes[2].m.preparedCertificate
		judge := s.nodes[1].m
		s.committees[1] = []common.Address{s.addrs[1]}
		require.NoError(t, judge.verifyPreparedCertificate(cert, viewAt(1, 1)))
		s.committees[0] = []common.Address{s.addrs[1]}
		require.Error(t, judge.verifyPreparedCertificate(cert, viewAt(1, 1)))
	})
}
