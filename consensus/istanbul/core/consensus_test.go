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
	"errors"
	"fmt"
	"math/big"
	"testing"
	"testing/synctest"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

// This file tests consensus safety and progress through scenarios involving multiple validators.
// It checks that successful validators agree on a block, insufficient or invalid votes cannot authorize a commit,
// and locks survive round changes. Recovery scenarios check that consensus resumes once enough validators
// can exchange valid messages, including across sequence changes and the Permissionless fork boundary.
//
// Each scenario runs the production event loops and timers inside synctest; scenarioNet controls message delivery.
// advanceConsensus supplies worker requests and processes queued events, checking commits except for excluded nodes.
// Validator methods check intermediate states, rejected inputs and recovery results where the scenario requires them.
//
// The scenarios are grouped by the consensus behavior they exercise:
// TestConsensusNetworkAndVoting: agreement and progress under reordered, missing, duplicate or conflicting votes.
// TestConsensusRoundChange: lock preservation and recovery through round changes, timeouts and commit errors.
// TestConsensusInputValidation: rejection of invalid messages, unauthorized voters and invalid committed seals.
// TestConsensusSequence: future messages, historical proposal replies and worker requests ahead of head notifications.
// TestConsensusForkBoundary: height-dependent seal rules before, at and after Permissionless activation.

// TestConsensusNetworkAndVoting checks agreement and progress with delayed, missing or conflicting inputs.
// It exercises quorum boundaries, signer deduplication and proposer equivocation using actual core transitions.
// Successful nodes must commit the same block; insufficient or conflicting votes must not create a certificate.
func TestConsensusNetworkAndVoting(t *testing.T) {
	for _, tc := range []struct {
		name             string
		count, committee int
		run              func(*scenarioNet)
	}{
		{"ConsecutiveBlocks", 4, 4, func(s *scenarioNet) {
			for height := uint64(1); height <= 4; height++ {
				s.advanceConsensus(height)
			}
		}},
		{"CommitsWithReorderedMessages", 4, 4, func(s *scenarioNet) {
			// V1 receives votes before the proposal; releasing it must drain those deferred votes.
			s.delay(bft.MsgPreprepare, s.nodes(0), s.nodes(1))
			s.advanceConsensus(1, s.nodes(1))
			s.release(bft.MsgPreprepare, s.nodes(0), s.nodes(1))
			s.validators[1].assertCommitted(1)
		}},
		{"CommitsAfterSilentProposer", 4, 4, func(s *scenarioNet) {
			s.isolate(s.nodes(0))
			s.timeout(s.nodes(1, 2, 3))
			s.advanceConsensus(1, s.nodes(0))
			s.advanceConsensus(2, s.nodes(0))
			for _, n := range s.nodes(1, 2, 3) {
				require.Equal(s.t, uint64(1), n.committedRound(1))
			}
			s.validators[0].assertUncommitted(1)
			// Peer isolation leaves V0 running: its own PREPARE still counts locally.
			s.validators[0].assertVoteCounts(1, 0, 1)
		}},
		{"CommitsWithOneConflictingVoter", 4, 4, func(s *scenarioNet) {
			// Keep only V0 at Q−1 COMMITs while the other nodes finish normally.
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			s.advanceConsensus(1, s.nodes(0))
			target, voter := s.validators[0], s.validators[3]
			b := target.alternative(s.proposal)
			for _, code := range []uint64{bft.MsgPrepare, bft.MsgCommit} {
				target.reject(voter.message(code, b, 0), errInconsistentSubject)
			}
			s.release(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			target.assertCommitted(1)
		}},
		{"DoesNotCommitAfterLosingQuorum", 4, 4, func(s *scenarioNet) {
			s.advanceConsensus(1)
			// Two late PREPREPARE deliveries leave only V0/V1 able to vote, below Q=3 at the next height.
			s.delay(bft.MsgPreprepare, s.nodes(1), s.nodes(2, 3), 0)
			s.advanceConsensus(2, s.nodes(0, 1, 2, 3))
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertUncommitted(2)
			}
			for _, n := range s.nodes(0, 1) {
				n.assertVoteCounts(2, 0, 2)
				n.assertHashLocked(common.Hash{})
			}
		}},
		{"DuplicateAndMixedVotesDoNotReachQuorum", 4, 4, func(s *scenarioNet) {
			s.delay(bft.MsgPrepare, s.nodes(2, 3), s.nodes(0))
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			s.advanceConsensus(1, s.nodes(0))
			target := s.validators[0]
			// Replay at the core boundary checks signer deduplication independently of the P2P cache.
			for _, code := range []uint64{bft.MsgPrepare, bft.MsgCommit} {
				target.assertStateUnchanged(func() { target.receive(s.message(s.validators[1], code, 1, 0), nil) })
			}
			target.assertVoteCounts(2, 1, 2)
			target.assertHashLocked(common.Hash{})
			target.assertUncommitted(1)
			// A third identity supplies the prepare certificate; the normal self COMMIT finalizes.
			s.release(bft.MsgCommit, s.nodes(2), s.nodes(0))
			target.assertCommitted(1)
		}},
		{"CommitsWithRemoteCommitsBeforePrepares", 4, 4, func(s *scenarioNet) {
			s.delay(bft.MsgPrepare, s.nodes(1, 2, 3), s.nodes(0))
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			s.advanceConsensus(1, s.nodes(0))
			target := s.validators[0]
			// Self PREPARE plus two remote COMMITs must suffice even when remote PREPAREs arrive late.
			target.assertVoteCounts(1, 1, 2)
			target.assertHashLocked(common.Hash{})
			s.release(bft.MsgCommit, s.nodes(2), s.nodes(0))
			target.assertCommitted(1)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, tc.count, tc.committee, params.TestChainConfig.Copy())
				tc.run(s)
			})
		})
	}

	// Hold the final distinct vote at each phase, then deliver it to cross Q−1 → Q.
	for _, tc := range []struct {
		name                             string
		count, committee, quorum, faults int
	}{
		{"1Validators", 1, 1, 1, 0},
		{"2Validators", 2, 2, 2, 0},
		{"3Validators", 3, 3, 3, 0},
		{"4Validators", 4, 4, 3, 1},
		{"5Validators", 5, 5, 4, 1},
		{"6Validators", 6, 6, 4, 1},
		{"7Validators", 7, 7, 5, 2},
		{"Qualified7Committee4", 7, 4, 3, 1},
	} {
		for _, phase := range []struct {
			name string
			code uint64
		}{{"Prepare", bft.MsgPrepare}, {"Commit", bft.MsgCommit}} {
			t.Run(tc.name+"/"+phase.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					s := newScenarioNet(t, tc.count, tc.committee, params.TestChainConfig.Copy())
					ids := s.ids()
					committee := s.nodes(ids[:tc.committee]...)
					target := s.validators[0]
					target.assertQuorum(tc.quorum, tc.faults)
					// Delay only unseen signers' COMMITs; already-counted identities cannot increase the mixed quorum.
					if phase.code == bft.MsgPrepare {
						s.delay(bft.MsgCommit, committee[tc.quorum-1:], s.nodes(0))
					}
					s.delay(phase.code, committee[tc.quorum-1:], s.nodes(0))
					s.advanceConsensus(1, s.nodes(ids...))
					target.assertUncommitted(1)
					if phase.code == bft.MsgPrepare {
						target.assertVoteCounts(tc.quorum-1, max(0, tc.quorum-2), tc.quorum-1)
						target.assertHashLocked(common.Hash{})
					} else {
						target.assertVoteCounts(tc.committee, tc.quorum-1, tc.committee)
						target.assertHashLocked(s.proposal.Hash())
					}
					s.release(phase.code, committee[tc.quorum-1:tc.quorum], s.nodes(0))
					if phase.code == bft.MsgPrepare {
						target.assertHashLocked(s.proposal.Hash())
						target.assertUncommitted(1)
						s.release(bft.MsgCommit, committee[tc.quorum-1:], s.nodes(0))
					}
					for _, n := range committee {
						n.assertCommitted(1, 0)
					}
				})
			})
		}
	}

	// One Byzantine proposer may sign both A and B; honest recipients must not combine their votes.
	for _, tc := range []struct {
		name string
		run  func(*scenarioNet)
	}{
		{"IgnoresSecondValidPreprepare", func(s *scenarioNet) {
			sender, target := s.validators[0], s.validators[1]
			a := sender.proposal(1)
			b := sender.alternative(a)
			// For C=5/7, delaying the latter half's votes leaves V1 at Q−1 while its peers finish.
			late := s.nodes(s.ids()[(len(s.validators)+1)/2:]...)
			s.delay(bft.MsgPrepare, late, s.nodes(1))
			s.delay(bft.MsgCommit, late, s.nodes(1))
			s.advanceConsensus(1, s.nodes(1))
			target.assertStateUnchanged(func() { target.receive(sender.message(bft.MsgPreprepare, b, 0), nil) })
			s.release(bft.MsgPrepare, late, s.nodes(1))
			s.release(bft.MsgCommit, late, s.nodes(1))
			require.Equal(s.t, a.Hash(), target.assertCommitted(1).Hash())
		}},
		{"SplitVotesDoNotCommit", func(s *scenarioNet) {
			ids := s.ids()
			// Each half gets its own proposal plus V0's vote: exactly Q−1 distinct identities.
			middle := 1 + (len(ids)-1)/2
			sender := s.validators[0]
			a, b := sender.proposal(1), sender.proposal(2)
			for _, code := range []uint64{bft.MsgPreprepare, bft.MsgPrepare} {
				s.modify(code, sender, s.nodes(ids[middle:]...), b)
			}
			s.advanceConsensus(1, s.nodes(ids...))
			// A COMMIT from the same proposer must not count as an additional identity.
			for _, group := range []struct {
				ids      []int
				proposal *types.Block
			}{{ids[1:middle], a}, {ids[middle:], b}} {
				for _, n := range s.nodes(group.ids...) {
					n.receive(sender.message(bft.MsgCommit, group.proposal, 0), nil)
					n.assertVoteCounts(len(group.ids)+1, 1, len(group.ids)+1)
					n.assertHashLocked(common.Hash{})
				}
			}
			s.drain()
			for _, n := range s.nodes(ids[1:]...) {
				n.assertUncommitted(1)
			}
		}},
	} {
		for _, count := range []int{5, 7} {
			t.Run(fmt.Sprintf("%s/%dValidators", tc.name, count), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					s := newScenarioNet(t, count, count, params.TestChainConfig.Copy())
					tc.run(s)
				})
			})
		}
	}
}

// TestConsensusRoundChange checks how locked proposals and round-change evidence survive interrupted rounds.
// It distinguishes waiting for enough votes from starting a round and rejects effects from stale timeouts.
// Recovery must commit a permitted proposal without discarding a valid lock or inventing a successful commit.
func TestConsensusRoundChange(t *testing.T) {
	for _, tc := range []struct {
		name             string
		count, committee int
		run              func(*scenarioNet)
	}{
		{"IgnoresTimeoutFromCommittedHeight", 4, 4, func(s *scenarioNet) {
			target := s.validators[0]
			// A timeout event posted before cancellation may arrive after the view has advanced.
			stale := timeoutEvent{nextView: &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}}
			s.advanceConsensus(1)
			target.assertStateUnchanged(func() { target.post(stale); s.drain() })
			s.advanceConsensus(2)
		}},
		{"CommitsLockedProposalAfterRoundChange", 4, 4, func(s *scenarioNet) {
			// Round-0 COMMITs are delayed between two groups; PREPAREs and self-delivery work.
			s.delay(bft.MsgCommit, s.nodes(0, 1), s.nodes(2, 3), 0)
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0, 1), 0)
			s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
			a := s.proposal
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertUncommitted(1)
				n.assertHashLocked(a.Hash())
				n.assertVoteCounts(4, 2, 4)
			}
			// New-round traffic is healthy; V1 re-proposes A and collects fresh round-1 votes.
			s.timeout(s.nodes(0, 1, 2))
			s.advanceConsensus(1)
			// Late round-0 packets cannot cause a second commit.
			s.release(bft.MsgCommit, s.nodes(0, 1, 2, 3), s.nodes(0, 1, 2, 3), 0)
			for _, n := range s.nodes(0, 1, 2, 3) {
				require.Equal(s.t, a.Hash(), n.assertCommitted(1, 1).Hash())
			}
		}},
		{"CommitsAfterTwoSilentProposers", 7, 7, func(s *scenarioNet) {
			s.isolate(s.nodes(0, 1))
			s.timeout(s.nodes(2, 3, 4, 5, 6))
			s.advanceConsensus(1, s.nodes(0, 1, 2, 3, 4, 5, 6))
			for _, n := range s.nodes(2, 3, 4, 5, 6) {
				n.assertUncommitted(1)
				n.assertView(1, 1, false)
			}
			s.timeout(s.nodes(2, 3, 4, 5, 6))
			s.advanceConsensus(1, s.nodes(0, 1))
			for _, n := range s.nodes(2, 3, 4, 5, 6) {
				n.assertCommitted(1, 2)
			}
		}},
		{"RoundChangeQuorumBoundaries", 7, 7, func(s *scenarioNet) {
			// With C=7, a non-waiting node stores f+1=3 votes but cannot start a round before Q=5.
			s.delay(bft.MsgRoundChange, s.nodes(3, 4), s.nodes(0))
			s.timeout(s.nodes(1, 2, 3, 4))
			s.drain()
			target := s.validators[0]
			target.assertRoundChangeCount(1, 2)
			target.assertView(1, 0, false)
			s.release(bft.MsgRoundChange, s.nodes(3), s.nodes(0))
			target.assertStateUnchanged(func() { target.receive(s.message(s.validators[3], bft.MsgRoundChange, 1, 1), errIgnored) })
			target.assertRoundChangeCount(1, 3)
			target.assertView(1, 0, false)
			// Local timeout adds V0's own fourth vote; the last peer vote is still needed to start.
			target.timeout()
			s.drain()
			target.assertRoundChangeCount(1, 4)
			target.assertView(1, 1, true)
			s.release(bft.MsgRoundChange, s.nodes(4), s.nodes(0))
			target.assertView(1, 1, false)
			target.assertRoundChangeCount(1, 0)
			s.advanceConsensus(1)
		}},
		{"TimeoutFollowsHigherWeakCertificate", 7, 7, func(s *scenarioNet) {
			// Three peers time out twice while V0 stays in round 0; f+1 votes alone cannot start round 2.
			for range 2 {
				s.timeout(s.nodes(1, 2, 3))
				s.drain()
			}
			target := s.validators[0]
			target.assertView(1, 0, false)
			target.assertRoundChangeCount(2, 3)
			// V0's timeout must follow the round-2 evidence rather than request only round 1.
			target.timeout()
			s.drain()
			target.assertView(1, 2, true)
			target.assertRoundChangeCount(2, 4)
		}},
		{"CommitsAfterPartialLockAndConflictingProposal", 4, 4, func(s *scenarioNet) {
			// Only V0/V2 lock A; V0's COMMIT is safe for V1/V3 because they already counted V0's PREPARE.
			s.delay(bft.MsgPrepare, s.nodes(2, 3), s.nodes(1), 0)
			s.delay(bft.MsgPrepare, s.nodes(1, 2), s.nodes(3), 0)
			s.delay(bft.MsgCommit, s.nodes(2), s.nodes(1, 3), 0)
			s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
			a := s.proposal
			for _, n := range s.nodes(0, 2) {
				n.assertHashLocked(a.Hash())
				n.assertUncommitted(1)
			}
			for _, n := range s.nodes(1, 3) {
				n.assertHashLocked(common.Hash{})
				n.assertVoteCounts(2, 1, 2)
			}
			// Unlocked proposer V1 proposes B; the locked nodes refuse it and request round 2.
			s.timeout(s.nodes(0, 1, 2, 3))
			s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
			require.NotEqual(s.t, a.Hash(), s.proposal.Hash())
			for _, n := range s.nodes(0, 2) {
				n.assertView(1, 2, true)
				n.assertHashLocked(a.Hash())
				n.assertUncommitted(1)
			}
			for _, n := range s.nodes(1, 3) {
				n.assertView(1, 1, false)
				n.assertHashLocked(common.Hash{})
				n.assertUncommitted(1)
			}
			// V1/V3 also time out; locked proposer V2 can now recover A in round 2.
			s.timeout(s.nodes(1, 3))
			s.advanceConsensus(1)
			for _, n := range s.nodes(0, 1, 2, 3) {
				require.Equal(s.t, a.Hash(), n.assertCommitted(1, 2).Hash())
			}
		}},
		{"SingleLockedNodeFollowsFinalizedHead", 4, 4, func(s *scenarioNet) {
			// Only V0 locks A; its COMMIT adds no new identity to peers that already counted its PREPARE.
			s.delay(bft.MsgPrepare, s.nodes(1), s.nodes(2, 3), 0)
			s.delay(bft.MsgPrepare, s.nodes(2), s.nodes(1, 3), 0)
			s.delay(bft.MsgPrepare, s.nodes(3), s.nodes(1, 2), 0)
			s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
			target := s.validators[0]
			a := s.proposal
			target.assertHashLocked(a.Hash())
			target.assertUncommitted(1)
			for _, n := range s.nodes(1, 2, 3) {
				n.assertHashLocked(common.Hash{})
			}
			s.timeout(s.nodes(0, 1, 2, 3))
			s.advanceConsensus(1, s.nodes(0))
			b := s.validators[1].assertCommitted(1, 1)
			require.NotEqual(s.t, a.Hash(), b.Hash())
			target.assertHashLocked(a.Hash())
			target.assertUncommitted(1)
			// Supply an already finalized head to the core; block import itself is outside this fixture.
			target.follow(s.validators[1])
			target.assertView(2, 0, false)
			target.assertHashLocked(common.Hash{})
		}},
		{"FollowsFinalizedHeadAfterCommitError", 4, 4, func(s *scenarioNet) {
			target := s.validators[0]
			target.failNextCommit(errors.New("temporary Commit failure"))
			s.advanceConsensus(1, s.nodes(0))
			target.assertUncommitted(1)
			target.assertHashLocked(common.Hash{})
			target.assertView(1, 1, true)
			target.assertRoundChangeCount(1, 1)
			target.commitAttempts(1)
			// Peers already finalized H1; recovery adopts that head instead of forcing them to retry H1.
			target.follow(s.validators[1])
			target.assertView(2, 0, false)
			target.assertHashLocked(common.Hash{})
		}},
		{"WaitingNodeFollowsWeakCertificate", 7, 7, func(s *scenarioNet) {
			all := s.nodes(0, 1, 2, 3, 4, 5, 6)
			s.delay(bft.MsgRoundChange, s.nodes(4), s.nodes(0), 1)
			s.delay(bft.MsgRoundChange, s.nodes(3, 4), s.nodes(0), 2)
			target := s.validators[0]
			target.timeout()
			s.drain()
			for range 2 {
				s.timeout(s.nodes(1, 2, 3, 4))
				s.drain()
			}
			// Two R2 votes cannot move V0; the third triggers catch-up and its own fourth vote before Q=5 starts.
			target.assertRoundChangeCount(2, 2)
			target.assertView(1, 1, true)
			s.release(bft.MsgRoundChange, s.nodes(3), s.nodes(0), 2)
			target.assertRoundChangeCount(2, 4)
			target.assertView(1, 2, true)
			s.release(bft.MsgRoundChange, s.nodes(4), s.nodes(0), 2)
			target.assertView(1, 2, false)
			target.assertRoundChangeCount(2, 0)
			s.advanceConsensus(1)
			for _, n := range all {
				n.assertCommitted(1, 2)
			}
		}},
		{"IgnoresTimeoutFromEarlierRound", 4, 4, func(s *scenarioNet) {
			target := s.validators[0]
			// A timeout event posted before cancellation may arrive after the view has advanced.
			stale := timeoutEvent{nextView: &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}}
			for range 2 {
				s.timeout(s.nodes(0, 1, 2, 3))
				s.drain()
			}
			target.assertView(1, 2, false)
			target.assertStateUnchanged(func() { target.post(stale); s.drain() })
			s.advanceConsensus(1)
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertCommitted(1, 2)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				// These cases assert legacy hash-lock recovery. Permissionless
				// prepared-certificate behavior is covered separately below.
				s := newScenarioNet(t, tc.count, tc.committee, params.TestKaiaConfig("osaka"))
				tc.run(s)
			})
		})
	}
}

// TestConsensusSplitLockPreparedCertificate covers the post-Permissionless
// counterexample where the next proposer did not observe the prepared value.
func TestConsensusSplitLockPreparedCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		require.True(t, s.validators[0].backend.IsPermissionlessAt(1))
		attacker, nextProposer := s.validators[0], s.validators[1]
		locked := s.nodes(2, 3)
		x := attacker.proposal(1)
		y := attacker.alternative(x)

		// A equivocates: B, the next proposer, sees Y while C/D see X. A votes
		// only for X and withholds its COMMIT. C/D lock X, but it is undecided.
		s.modify(bft.MsgPreprepare, attacker, s.nodes(1), y)
		s.drop(bft.MsgPrepare, s.nodes(0), s.nodes(1))
		s.drop(bft.MsgCommit, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
		for _, n := range locked {
			n.assertHashLocked(x.Hash())
			n.assertUncommitted(1)
			require.NoError(t, nextProposer.core.verifyPreparedCertificate(
				n.core.current.PreparedCertificate(), &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}))
		}
		nextProposer.assertHashLocked(common.Hash{})

		// The honest round-change quorum carries X's prepared proof. In round 1,
		// B must learn and re-propose X even though B was not itself locked on X.
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.timeout(s.nodes(1, 2, 3))
		s.advanceConsensus(1, s.nodes(0))
		for _, n := range s.nodes(1, 2, 3) {
			require.Equal(t, x.Hash(), n.assertCommitted(1, 1).Hash())
		}
	})
}

// TestConsensusSplitLockReproposalRefreshesPreparedCertificate ensures a
// locked node re-establishes its prepared proof in the recovery round before
// it can contribute a COMMIT. Retaining only the r0 proof after sending a
// r1 COMMIT would let a later RCC prefer a conflicting, newer certificate.
func TestConsensusSplitLockReproposalRefreshesPreparedCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		attacker := s.validators[0]
		locked := s.nodes(2, 3)
		x := attacker.proposal(1)
		y := attacker.alternative(x)

		// C/D prepare and lock X in round 0 while B sees an equivocated Y.
		s.modify(bft.MsgPreprepare, attacker, s.nodes(1), y)
		s.drop(bft.MsgPrepare, s.nodes(0), s.nodes(1))
		s.drop(bft.MsgCommit, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
		for _, n := range locked {
			require.Equal(t, int64(0), n.core.current.LockedRound().Int64())
		}

		// In r1, B learns X through C/D's round changes. Drop the locked nodes'
		// final COMMITs so the test can inspect their evidence before H1 commits.
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.drop(bft.MsgCommit, locked, s.nodes(0, 1, 2, 3), 1)
		s.timeout(s.nodes(1, 2, 3))
		s.drain()

		for _, n := range locked {
			cert := n.core.current.PreparedCertificate()
			require.NotNil(t, cert)
			require.Equal(t, int64(1), cert.View.Round.Int64())
			require.Equal(t, x.Hash(), cert.Proposal.Hash())
		}
	})
}

// TestConsensusInputValidation checks malformed messages, sender eligibility, proposal validity and commit seals.
// Pure rejection must leave node state unchanged; invalid proposals may instead trigger a round change.
// Valid traffic must still progress after the faulty input, using the seal rules active at the block height.
func TestConsensusInputValidation(t *testing.T) {
	for _, tc := range []struct {
		name             string
		count, committee int
		run              func(*scenarioNet)
	}{
		{"SingleHighRoundVoterCannotForceCatchUp", 4, 4, func(s *scenarioNet) {
			// One attacker can claim round 2 immediately, but one signature cannot form a weak certificate.
			sender := s.validators[3]
			target := s.validators[0]
			target.receive(sender.message(bft.MsgRoundChange, sender.proposal(1), 2), errIgnored)
			target.assertRoundChangeCount(2, 1)
			target.assertView(1, 0, false)
			s.timeout(s.nodes(0, 1, 2))
			s.advanceConsensus(1)
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertCommitted(1, 1)
			}
		}},
		{"CommitsDespiteNonCommitteeTraffic", 5, 4, func(s *scenarioNet) {
			target, outsider := s.validators[0], s.validators[4]
			target.reject(outsider.message(bft.MsgPreprepare, outsider.proposal(1), 0), errNotFromProposer)
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			s.advanceConsensus(1, s.nodes(0, 4))
			// V0 has processed PREPREPARE, so votes must be rejected rather than deferred.
			for _, code := range []uint64{bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange} {
				target.reject(outsider.message(code, s.proposal, 0), errNotFromCommittee)
			}
			s.release(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			target.assertCommitted(1)
		}},
		{"ChangesRoundAfterInvalidProposal", 4, 4, func(s *scenarioNet) {
			sender := s.validators[0]
			// A faulty proposer signs a bad-parent block; validation must request the next round.
			invalid := sender.message(bft.MsgPreprepare, sender.invalidParentProposal(), 0)
			for _, n := range s.nodes(1, 2, 3) {
				n.receive(invalid, istanbul.ErrInvalidProposal)
				n.assertView(1, 1, true)
				n.assertHashLocked(common.Hash{})
			}
			s.advanceConsensus(1)
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertCommitted(1, 1)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, tc.count, tc.committee, params.TestChainConfig.Copy())
				tc.run(s)
			})
		})
	}

	// Malformed signed input must leave state unchanged before normal consensus resumes.
	for _, tc := range []struct {
		name string
		kind consensusInvalidKind
		want error
	}{
		{"MalformedRLP", consensusMalformedRLP, nil},
		{"UnknownCode", consensusUnknownCode, bft.ErrInvalidMessage},
		{"InvalidSignature", consensusInvalidSignature, nil},
		{"MissingView", consensusMissingView, bft.ErrInvalidMessage},
		{"OverflowRound", consensusOverflowRound, bft.ErrInvalidMessage},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, 4, 4, params.TestChainConfig.Copy())
				target := s.validators[0]
				ev := s.validators[3].invalidMessage(tc.kind)
				if tc.want == nil {
					target.reject(ev)
				} else {
					target.reject(ev, tc.want)
				}
				s.advanceConsensus(1)
			})
		})
	}

	// Malformed seal lengths fail envelope validation; valid lengths still require matching signatures.
	for _, tc := range []struct {
		name     string
		mutation consensusSealMutation
		want     error
	}{
		{"Malformed", consensusMalformedSeal, bft.ErrInvalidMessage},
		{"WrongSigner", consensusWrongSigner, errInvalidCommittedSeal},
		{"WrongDigest", consensusWrongDigest, errInvalidCommittedSeal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, 4, 4, params.TestChainConfig.Copy())
				target := s.validators[0]
				s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
				s.advanceConsensus(1, s.nodes(0))
				// Only V0 is short of one COMMIT; the faulty seal must not become that deciding vote.
				ev := s.validators[3].corruptCommit(tc.mutation)
				target.reject(ev, tc.want)
				target.assertUncommitted(1)
				s.release(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
				target.assertCommitted(1)
			})
		})
	}
}

// TestConsensusSequence checks messages and worker requests that arrive across a height transition.
// Future messages and requests must wait for the matching height, while stale requests must not alter completed work.
// Historical proposals are checked for locally generated reply signatures, without claiming successful peer delivery.
func TestConsensusSequence(t *testing.T) {
	for _, tc := range []struct {
		name             string
		count, committee int
		run              func(*scenarioNet)
	}{
		{"GeneratesCommitForStoredOldProposal", 4, 4, func(s *scenarioNet) {
			s.advanceConsensus(1)
			target := s.validators[1]
			// Check the reply signature only; current core rejects the old COMMIT before forwarding it.
			target.assertCommitReply(1, 0, true, func() {
				target.receive(s.message(s.validators[0], bft.MsgPreprepare, 1, 0), nil)
			})
		}},
		{"CommitsBufferedNextHeightMessages", 4, 4, func(s *scenarioNet) {
			// Only V3 misses the COMMIT quorum; its peers can finish H1 and send H2 messages ahead of it.
			s.delay(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3), 0)
			s.advanceConsensus(1, s.nodes(3))
			s.advanceConsensus(2, s.nodes(3))
			target := s.validators[3]
			target.assertUncommitted(1)
			target.assertView(1, 0, false)
			// H2's PREPREPARE and three PREPAREs wait in the backlog while H1 is still unfinished.
			target.assertBacklogCount(4)
			s.release(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3), 0)
			target.assertCommitted(2)
			target.assertBacklogCount(0)
		}},
		{"QueuesFutureWorkerRequestBeforeHeadEvent", 4, 4, func(s *scenarioNet) {
			target := s.validators[1]
			// Head storage has advanced but its notification is pending, so the H2 worker request is early.
			s.until(1, func() bool { return target.head().NumberU64() == 1 })
			target.post(istanbul.RequestEvent{Proposal: target.proposal(1)})
			target.assertPendingRequestCount(1)
			s.drain()
			for _, n := range s.nodes(0, 1, 2, 3) {
				n.assertCommitted(2)
			}
			target.assertStateUnchanged(func() { target.post(istanbul.RequestEvent{Proposal: target.block(1)}) })
			target.assertPendingRequestCount(0)
			s.advanceConsensus(3)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, tc.count, tc.committee, params.TestChainConfig.Copy())
				tc.run(s)
			})
		})
	}
}

// TestConsensusForkBoundary advances one chain through the height where Permissionless seal rules activate.
// Nonzero-round commits before, at and after activation must use the rule for their own block height.
// Wrong-round seals are rejected after activation, while historical replies retain their block height's format.
func TestConsensusForkBoundary(t *testing.T) {
	config := params.TestKaiaConfig("permissionless")
	config.PermissionlessCompatibleBlock = big.NewInt(2)
	t.Run("PermissionlessActivationBoundary", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			s := newScenarioNet(t, 4, 4, config)
			for height := uint64(1); height <= 3; height++ {
				// Let Q nodes time out before worker submission to exercise a nonzero round across the fork.
				s.timeout(s.nodes(0, 1, 2))
				s.advanceConsensus(height)
				for _, n := range s.nodes(0, 1, 2, 3) {
					block := n.assertCommitted(height, 1)
					n.assertCommittedSeals(block, scenarioSealFormat{1, height >= 2})
				}
			}
			// Even after activation, the H1 reply must use the pre-fork seal format.
			target := s.validators[2]
			target.assertCommitReply(1, 1, false, func() { target.receive(s.message(s.validators[1], bft.MsgPreprepare, 1, 1), nil) })
		})
	})

	for _, tc := range []struct {
		name             string
		activationHeight int64
		want             error
	}{
		{"OtherRoundSealBeforeActivation", 2, nil},
		{"OtherRoundSealAfterActivation", 1, errInvalidCommittedSeal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				config := params.TestKaiaConfig("permissionless")
				config.PermissionlessCompatibleBlock = big.NewInt(tc.activationHeight)
				s := newScenarioNet(t, 4, 4, config)
				target := s.validators[0]
				s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
				s.advanceConsensus(1, s.nodes(0))
				ev := s.validators[3].corruptCommit(consensusOtherRound)
				// Before activation the seal omits the round; afterward the same mutation must be rejected.
				if tc.want == nil {
					target.receive(ev, nil)
					s.drain()
					target.assertCommitted(1)
				} else {
					target.reject(ev, tc.want)
					target.assertUncommitted(1)
				}
				s.release(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
				target.assertCommitted(1)
			})
		})
	}
}
