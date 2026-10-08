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
	"slices"
	"testing"
	"testing/synctest"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
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
		{"BadProposalInvalidatesExistingLock", 4, 4, func(s *scenarioNet) {
			// Hold V0 below the commit quorum after it has locked the proposal.
			s.delay(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
			s.advanceConsensus(1, s.nodes(0))
			target := s.validators[0]
			locked := s.proposal.Hash()
			target.assertHashLocked(locked)
			// A backend-reported bad proposal keeps the hash for diagnostics but invalidates the lock.
			target.backend.badProposals[locked] = true
			target.assertBadHashLocked(locked)
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

// TestConsensusRejectsVotesForAnotherParent checks that PREPARE, COMMIT and
// ROUND CHANGE messages of the current sequence must name the current chain
// head as PrevHash, in both wire formats. The legacy format signs it twice, in
// the envelope and the Subject, and normalization keeps them equal.
func TestConsensusRejectsVotesForAnotherParent(t *testing.T) {
	for _, tc := range []struct {
		name   string
		config *params.ChainConfig
	}{
		{"PrePermissionless", params.TestChainConfig.Copy()},
		{"Permissionless", params.TestKaiaConfig("permissionless")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := newScenarioNet(t, 4, 4, tc.config)
				target, sender := s.validators[0], s.validators[1]
				// Hold the target below the commit quorum so the current view stays open.
				s.delay(bft.MsgCommit, s.nodes(1, 2, 3), s.nodes(0))
				s.advanceConsensus(1, s.nodes(0))
				target.assertUncommitted(1)

				view := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}
				digest := s.proposal.Hash()
				for code, payload := range map[uint64]any{
					bft.MsgPrepare:     &bft.Prepare{View: view, Digest: digest},
					bft.MsgCommit:      &bft.Commit{View: view, Digest: digest},
					bft.MsgRoundChange: &bft.RoundChange{View: view},
				} {
					encoded, err := bft.Encode(payload)
					require.NoError(t, err)
					wrongParent := common.HexToHash("0xbad")
					signed, err := sender.core.finalizeMessage(&bft.Message{PrevHash: wrongParent, Code: code, Msg: encoded})
					require.NoError(t, err)
					target.reject(istanbul.MessageEvent{Hash: wrongParent, Payload: signed}, errInconsistentPrevHash)
				}

				s.release(bft.MsgCommit, s.nodes(1, 2, 3), s.nodes(0))
				target.assertCommitted(1)
			})
		})
	}
}

// TestConsensusJustifiedPreprepareOverridesOlderLock checks that a node locked
// on X@0 adopts a different value Y proven prepared at a higher round, and
// PREPAREs it, so the committee can finish on Y.
func TestConsensusJustifiedPreprepareOverridesOlderLock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		all := s.nodes(0, 1, 2, 3)
		a, b, d := s.validators[0], s.validators[1], s.validators[3]
		x := a.proposal(1)
		y := b.alternative(x)

		// Round 0: only D reaches the PREPARE quorum for X, and its COMMIT is lost.
		dropPeers(s, bft.MsgPrepare, all, s.nodes(0, 1, 2), 0)
		s.drop(bft.MsgCommit, s.nodes(3), all, 0)
		s.advanceConsensus(1, all)
		d.assertHashLocked(x.Hash())
		for _, n := range s.nodes(0, 1, 2) {
			n.assertHashLocked(common.Hash{})
		}

		// Round 1: B does not hear D's claim, proposes Y, and A/B/C lock Y@1.
		s.drop(bft.MsgRoundChange, s.nodes(3), s.nodes(1), 1)
		s.modify(bft.MsgPreprepare, b, all, y)
		s.drop(bft.MsgCommit, s.nodes(0, 1, 2), all, 1)
		s.timeout(all)
		settleConsensus(s, 1)
		for _, n := range s.nodes(0, 1, 2) {
			n.assertHashLocked(y.Hash())
			require.Equal(t, int64(1), n.core.current.LockedRound().Int64())
		}
		d.assertHashLocked(x.Hash())
		require.Zero(t, d.core.current.LockedRound().Sign())

		// Round 2: C re-proposes Y with the r1 certificate. Hold the votes to D
		// so its adoption is observable before the height completes.
		s.delay(bft.MsgPrepare, s.nodes(0, 1, 2), s.nodes(3), 2)
		s.delay(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3), 2)
		s.timeout(s.nodes(0, 1, 2))
		settleConsensus(s, 1)
		require.Equal(t, y.Hash(), sentPreprepare(t, s, s.validators[2], 1, 2).Proposal.Hash())
		d.assertHashLocked(y.Hash())
		require.Equal(t, int64(1), d.core.current.LockedRound().Int64())
		s.message(d, bft.MsgPrepare, 1, 2)

		s.release(bft.MsgPrepare, s.nodes(0, 1, 2), s.nodes(3), 2)
		s.release(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3), 2)
		s.advanceConsensus(1)
		for _, n := range all {
			require.Equal(t, y.Hash(), n.assertCommitted(1, 2).Hash())
		}
	})
}

// TestConsensusOlderCertificateDoesNotOverrideNewerLock checks the inverse: a
// node locked on Y@1 keeps that lock when the proposer re-proposes a different
// value X justified only by an older X@0 certificate, and asks for a new round.
func TestConsensusOlderCertificateDoesNotOverrideNewerLock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		all := s.nodes(0, 1, 2, 3)
		a, b, c, d := s.validators[0], s.validators[1], s.validators[2], s.validators[3]
		x := a.proposal(1)
		y := b.alternative(x)

		// Round 0: only A locks X.
		dropPeers(s, bft.MsgPrepare, all, s.nodes(1, 2, 3), 0)
		s.drop(bft.MsgCommit, s.nodes(0), all, 0)
		s.advanceConsensus(1, all)
		a.assertHashLocked(x.Hash())

		// Round 1: B does not hear A's claim and proposes Y; only D locks Y@1.
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(1), 1)
		s.modify(bft.MsgPreprepare, b, all, y)
		dropPeers(s, bft.MsgPrepare, all, s.nodes(0, 1, 2), 1)
		s.drop(bft.MsgCommit, s.nodes(3), all, 1)
		s.timeout(all)
		settleConsensus(s, 1)
		d.assertHashLocked(y.Hash())
		require.Equal(t, int64(1), d.core.current.LockedRound().Int64())

		// Round 2: C's quorum {A, B, C} carries only A's X@0 claim, so C
		// re-proposes X. D must keep its newer lock and request round 3.
		s.drop(bft.MsgRoundChange, s.nodes(3), s.nodes(2), 2)
		s.timeout(s.nodes(1, 2, 3))
		settleConsensus(s, 1)
		preprepare := sentPreprepare(t, s, c, 1, 2)
		require.Equal(t, x.Hash(), preprepare.Proposal.Hash())
		d.assertHashLocked(y.Hash())
		require.Equal(t, int64(1), d.core.current.LockedRound().Int64())
		s.message(d, bft.MsgRoundChange, 1, 3)
	})
}

// TestConsensusRoundChangeCarriesCommitVoteCertificate checks a certificate in
// which a COMMIT counts toward the prepared quorum, end to end: the claimant's
// ROUND CHANGE is admitted by a peer, a round-bound seal for another round is
// rejected, and the recovery PRE-PREPARE carries the COMMIT vote.
func TestConsensusRoundChangeCarriesCommitVoteCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		all := s.nodes(0, 1, 2, 3)
		a, b, c, d := s.validators[0], s.validators[1], s.validators[2], s.validators[3]
		x := a.proposal(1)

		// Round 0: nobody reaches a PREPARE quorum. D holds PREPAREs from C and
		// itself, then a COMMIT from A completes its lock. D's COMMIT is lost,
		// so D is the only claimant.
		dropPeers(s, bft.MsgPrepare, all, s.nodes(0, 1, 2), 0)
		s.drop(bft.MsgPrepare, s.nodes(0, 1), s.nodes(3), 0)
		s.drop(bft.MsgCommit, s.nodes(3), all, 0)
		s.advanceConsensus(1, all)
		d.assertHashLocked(common.Hash{})
		d.receive(a.message(bft.MsgCommit, x, 0), nil)
		d.assertHashLocked(x.Hash())
		cert := d.core.current.PreparedCertificate()
		require.NotNil(t, cert)
		require.True(t, slices.ContainsFunc(cert.Messages, func(m *bft.Message) bool { return m.Code == bft.MsgCommit }))

		// A seal over another round byte, re-signed by A, must fail verification.
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: x.Hash()}
		forged := *cert
		forged.Messages = slices.Clone(cert.Messages)
		for i, vote := range forged.Messages {
			if vote.Code != bft.MsgCommit {
				continue
			}
			var commit bft.Commit
			require.NoError(t, vote.Decode(&commit))
			var err error
			commit.CommittedSeal, err = a.backend.sealer.MakeCommittedSealFromHashWithRound(x.Hash(), 1)
			require.NoError(t, err)
			reSealed := *vote
			reSealed.Msg, err = bft.Encode(&commit)
			require.NoError(t, err)
			unsigned, err := reSealed.PayloadNoSig()
			require.NoError(t, err)
			reSealed.Signature, err = a.backend.Sign(unsigned)
			require.NoError(t, err)
			forged.Messages[i] = &reSealed
		}
		b.reject(signedRoundChange(t, d, 1, claim, &forged), bft.ErrInvalidMessage)

		// A peer admits the genuine certificate.
		c.receive(signedRoundChange(t, d, 1, claim, cert), errIgnored)
		c.assertRoundChangeCount(1, 1)

		// Round 1: B's quorum {B, C, D} carries D's claim, so B re-proposes X
		// with D's mixed vote set.
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(1), 1)
		s.timeout(all)
		settleConsensus(s, 1)
		preprepare := sentPreprepare(t, s, b, 1, 1)
		require.Equal(t, x.Hash(), preprepare.Proposal.Hash())
		require.True(t, slices.ContainsFunc(preprepare.PreparedMessages, func(m *bft.Message) bool {
			return m.Code == bft.MsgCommit && m.Address == a.backend.Address()
		}))
		s.advanceConsensus(1)
		for _, n := range all {
			require.Equal(t, x.Hash(), n.assertCommitted(1, 1).Hash())
		}
	})
}

// dropPeers drops code on every route from a sender in from to a different
// node in to. Self-delivery is kept: a node forwards a message to its peers
// only after handling it locally.
func dropPeers(s *scenarioNet, code uint64, from, to []*validator, round uint64) {
	s.t.Helper()
	for _, sender := range from {
		for _, recipient := range to {
			if sender != recipient {
				s.drop(code, []*validator{sender}, []*validator{recipient}, round)
			}
		}
	}
}

// settleConsensus delivers queued messages and worker requests for height
// without requiring the height to commit.
func settleConsensus(s *scenarioNet, height uint64) {
	s.t.Helper()
	for steps := 0; len(s.pendingDeliveries) > 0 || s.requestProposal(height); steps++ {
		require.Less(s.t, steps, maxScenarioSteps, "consensus did not quiesce")
		s.step()
	}
}

// splitLock reproduces the round-0 split lock used by the tests above: C/D
// lock X, B (the round-1 proposer) saw only Y, and A's COMMIT is withheld.
func splitLock(t *testing.T) (s *scenarioNet, x *types.Block) {
	s = newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
	attacker := s.validators[0]
	x = attacker.proposal(1)
	s.modify(bft.MsgPreprepare, attacker, s.nodes(1), attacker.alternative(x))
	s.drop(bft.MsgPrepare, s.nodes(0), s.nodes(1))
	s.drop(bft.MsgCommit, s.nodes(0), s.nodes(0, 1, 2, 3))
	s.advanceConsensus(1, s.nodes(0, 1, 2, 3))
	for _, n := range s.nodes(2, 3) {
		n.assertHashLocked(x.Hash())
	}
	return s, x
}

// sentPreprepare decodes the PRE-PREPARE a node actually broadcast.
func sentPreprepare(t *testing.T, s *scenarioNet, from *validator, height, round uint64) *bft.Preprepare {
	msg, err := from.backend.decodeMessage(s.message(from, bft.MsgPreprepare, height, round).Payload, nil)
	require.NoError(t, err)
	var preprepare *bft.Preprepare
	require.NoError(t, msg.Decode(&preprepare))
	return preprepare
}

// signedRoundChange signs a ROUND CHANGE for the sender's head with an
// arbitrary claim and attached certificate.
func signedRoundChange(t *testing.T, sender *validator, round uint64, claim *bft.PreparedClaim, cert *bft.PreparedCertificate) istanbul.MessageEvent {
	var evidence []byte
	if cert != nil {
		var err error
		evidence, err = bft.Encode(cert)
		require.NoError(t, err)
	}
	return signedRoundChangeAt(t, sender, sender.head().NumberU64()+1, round, claim, evidence)
}

// signedRoundChangeAt signs a ROUND CHANGE for an arbitrary sequence with raw
// Evidence bytes, which the sender's signature does not cover.
func signedRoundChangeAt(t *testing.T, sender *validator, sequence, round uint64, claim *bft.PreparedClaim, evidence []byte) istanbul.MessageEvent {
	payload, err := sender.core.finalizeMessage(roundChangeMessageAt(t, sender, sequence, round, claim, evidence))
	require.NoError(t, err)
	return istanbul.MessageEvent{Hash: sender.head().Hash(), Payload: payload}
}

// postPermissionlessRoundChange signs a ROUND CHANGE in the post-Permissionless
// wire format regardless of the fork active at sequence, as a peer could.
func postPermissionlessRoundChange(t *testing.T, sender *validator, sequence, round uint64, claim *bft.PreparedClaim, evidence []byte) istanbul.MessageEvent {
	msg := roundChangeMessageAt(t, sender, sequence, round, claim, evidence)
	msg.Address = sender.backend.Address()
	unsigned, err := msg.PayloadNoSig()
	require.NoError(t, err)
	msg.Signature, err = sender.backend.Sign(unsigned)
	require.NoError(t, err)
	payload, err := msg.Payload()
	require.NoError(t, err)
	return istanbul.MessageEvent{Hash: msg.PrevHash, Payload: payload}
}

// roundChangeMessageAt builds an unsigned ROUND CHANGE. A claim without an
// EvidenceHash is bound to evidence, as an honest sender would sign it; set the
// hash explicitly to sign a claim for other bytes.
func roundChangeMessageAt(t *testing.T, sender *validator, sequence, round uint64, claim *bft.PreparedClaim, evidence []byte) *bft.Message {
	if claim != nil && claim.EvidenceHash == (common.Hash{}) && len(evidence) != 0 {
		bound := *claim
		bound.EvidenceHash = crypto.Keccak256Hash(evidence)
		claim = &bound
	}
	encoded, err := bft.Encode(&bft.RoundChange{
		View:     &bft.View{Sequence: new(big.Int).SetUint64(sequence), Round: new(big.Int).SetUint64(round)},
		Prepared: claim,
	})
	require.NoError(t, err)
	return &bft.Message{PrevHash: sender.head().Hash(), Code: bft.MsgRoundChange, Msg: encoded, Evidence: evidence}
}

// TestConsensusRecoveryPreprepareCarriesProposalOnce checks the QBFT-style
// justification: the recovery PRE-PREPARE embeds signed ROUND CHANGE claims
// stripped of their certificates, plus only the winning certificate's votes.
// The prepared block travels once, as the proposal, however many senders
// hold it.
func TestConsensusRecoveryPreprepareCarriesProposalOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, x := splitLock(t)
		nextProposer := s.validators[1]
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.timeout(s.nodes(1, 2, 3))
		s.advanceConsensus(1, s.nodes(0))

		preprepare := sentPreprepare(t, s, nextProposer, 1, 1)
		require.Equal(t, x.Hash(), preprepare.Proposal.Hash())
		require.Len(t, preprepare.RoundChangeCertificate, nextProposer.core.current.requiredMessageCount)
		claims := 0
		for _, rc := range preprepare.RoundChangeCertificate {
			require.Empty(t, rc.Evidence, "embedded ROUND CHANGE must not repeat the prepared block")
			require.LessOrEqual(t, retainedMessageBytes(rc), uint64(maxSubjectMessageBytes))
			var roundChange *bft.RoundChange
			require.NoError(t, rc.Decode(&roundChange))
			if roundChange.Prepared != nil {
				claims++
				require.Equal(t, x.Hash(), roundChange.Prepared.Digest)
				require.Zero(t, roundChange.Prepared.Round.Sign())
			}
		}
		require.Equal(t, 2, claims, "C and D both claim X")
		require.GreaterOrEqual(t, len(preprepare.PreparedMessages), nextProposer.core.current.requiredMessageCount)
	})
}

// TestConsensusPreprepareJustificationValidation mutates a real recovery
// PRE-PREPARE and checks that a receiver rejects every unjustified variant.
func TestConsensusPreprepareJustificationValidation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, x := splitLock(t)
		nextProposer, receiver := s.validators[1], s.validators[2]
		// Keep the receivers at (1, 1) so the captured PRE-PREPARE stays current.
		s.drop(bft.MsgRoundChange, s.nodes(0), s.nodes(0, 1, 2, 3))
		s.drop(bft.MsgPreprepare, s.nodes(1), s.nodes(0, 2, 3), 1)
		s.timeout(s.nodes(1, 2, 3))
		s.drain()
		receiver.assertView(1, 1, false)
		quorum := receiver.core.current.requiredMessageCount

		cert, err := receiver.core.verifyPreprepareJustification(sentPreprepare(t, s, nextProposer, 1, 1))
		require.NoError(t, err)
		require.Equal(t, x.Hash(), cert.Proposal.Hash())
		require.Zero(t, cert.View.Round.Sign())

		// Each case asserts its own error, so that a check that is removed but
		// masked by a later one still fails the test.
		for _, tc := range []struct {
			name   string
			mutate func(*bft.Preprepare)
			want   string
		}{
			{
				"proposal is not the highest prepared value", func(pp *bft.Preprepare) { pp.Proposal = nextProposer.alternative(x) },
				"proposal is not the highest prepared value",
			},
			{
				"missing prepared votes", func(pp *bft.Preprepare) { pp.PreparedMessages = nil },
				"prepared certificate has 0 votes",
			},
			{
				"prepared votes below quorum", func(pp *bft.Preprepare) { pp.PreparedMessages = pp.PreparedMessages[:quorum-1] },
				"prepared certificate has",
			},
			{"round-change certificate below quorum", func(pp *bft.Preprepare) {
				pp.RoundChangeCertificate = pp.RoundChangeCertificate[:quorum-1]
			}, "round-change certificate has"},
			// Appending keeps a quorum of distinct senders, so only the duplicate
			// check can reject it.
			{"duplicate round-change sender", func(pp *bft.Preprepare) {
				pp.RoundChangeCertificate = append(pp.RoundChangeCertificate, pp.RoundChangeCertificate[0])
			}, "duplicate sender"},
			{"embedded round change keeps its evidence", func(pp *bft.Preprepare) {
				embedded := *pp.RoundChangeCertificate[0]
				embedded.Evidence = []byte{0xc0}
				pp.RoundChangeCertificate[0] = &embedded
			}, "ineligible message"},
		} {
			pp := sentPreprepare(t, s, nextProposer, 1, 1)
			tc.mutate(pp)
			_, err := receiver.core.verifyPreprepareJustification(pp)
			require.ErrorContains(t, err, tc.want, tc.name)
		}
	})
}

// TestConsensusRoundChangeRejectsUnboundEvidence checks that a ROUND
// CHANGE certificate must prove exactly the signed claim, including the
// prepared block's body, which the votes alone do not bind.
func TestConsensusRoundChangeRejectsUnboundEvidence(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, x := splitLock(t)
		attacker, receiver := s.validators[0], s.validators[1]
		cert := s.validators[2].core.current.PreparedCertificate()
		require.NotNil(t, cert)
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: x.Hash()}

		// Same header and votes, different transactions: the block hash is
		// unchanged, so only the body check can catch it. The transaction is
		// signed so that the Evidence decodes and reaches that check.
		txKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		tx, err := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(0), nil),
			types.LatestSignerForChainID(s.config.chainConfig.ChainID), txKey)
		require.NoError(t, err)
		forged := &bft.PreparedCertificate{View: cert.View, Proposal: cert.Proposal.WithBody(types.Transactions{tx}), Messages: cert.Messages}
		require.Equal(t, x.Hash(), forged.Proposal.Hash())
		forgedMsg, err := receiver.backend.decodeMessage(signedRoundChange(t, attacker, 2, claim, forged).Payload, nil)
		require.NoError(t, err)
		var forgedRC *bft.RoundChange
		require.NoError(t, forgedMsg.Decode(&forgedRC))
		_, err = receiver.core.verifyRoundChangeEvidence(forgedMsg, forgedRC)
		require.ErrorIs(t, err, istanbul.ErrMismatchTxhashes, "the forged body reaches the body check")

		for _, tc := range []struct {
			name  string
			claim *bft.PreparedClaim
			cert  *bft.PreparedCertificate
		}{
			{"forged proposal body", claim, forged},
			{"claim without evidence", claim, nil},
			{"evidence without claim", nil, cert},
			{"claim for another value", &bft.PreparedClaim{Round: big.NewInt(0), Digest: attacker.alternative(x).Hash()}, cert},
			{"claim for another round", &bft.PreparedClaim{Round: big.NewInt(1), Digest: x.Hash()}, cert},
		} {
			t.Log(tc.name)
			receiver.reject(signedRoundChange(t, attacker, 2, tc.claim, tc.cert), bft.ErrInvalidMessage)
		}
		// The same construction with a matching claim and certificate is admitted;
		// a lone future-round message is counted but reported as ignored.
		receiver.receive(signedRoundChange(t, attacker, 2, claim, cert), errIgnored)
		receiver.assertRoundChangeCount(2, 1)
	})
}

// TestConsensusRoundChangeCachesPreparedEvidenceOnce ensures repeated claims
// retain only compact signed ROUND CHANGEs. The prepared block and quorum of
// votes are cached once per exact prepared claim, even across target rounds,
// and changing only the unsigned attachment cannot trigger verification again.
func TestConsensusRoundChangeCachesPreparedEvidenceOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, x := splitLock(t)
		sender, receiver := s.validators[0], s.validators[1]
		cert := s.validators[2].core.current.PreparedCertificate()
		require.NotNil(t, cert)
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: x.Hash()}

		first := signedRoundChange(t, sender, 2, claim, cert)
		receiver.receive(first, errIgnored)
		require.Len(t, receiver.core.preparedBlocks, 1)
		require.Len(t, receiver.core.preparedCertificates, 1)
		require.Equal(t, x.Hash(), receiver.core.preparedBlocks[x.Hash()].Hash())
		for _, cached := range receiver.core.preparedCertificates {
			require.Same(t, cached.Proposal, receiver.core.preparedBlocks[x.Hash()])
		}
		require.Len(t, receiver.core.roundChangeSet.Values(big.NewInt(2)), 1)
		require.Empty(t, receiver.core.roundChangeSet.Values(big.NewInt(2))[0].Evidence)

		// Count signature recoveries from here on: a cached claim must cost only
		// the envelope signature, not the certificate's votes.
		recoveries := 0
		validate := receiver.core.validateFn
		receiver.core.validateFn = func(data, sig []byte) (common.Address, error) {
			recoveries++
			return validate(data, sig)
		}

		// Evidence is excluded from the ROUND CHANGE signature but bound by the
		// signed EvidenceHash. A relay that replaces it is rejected before the
		// replacement is decoded or any certificate signature is recovered.
		var duplicate bft.Message
		require.NoError(t, duplicate.FromPayload(first.Payload, nil))
		duplicate.Evidence = []byte{0xc0}
		payload, err := duplicate.Payload()
		require.NoError(t, err)
		before := recoveries
		receiver.reject(istanbul.MessageEvent{Hash: first.Hash, Payload: payload}, bft.ErrInvalidMessage)
		require.Equal(t, before+1, recoveries, "only the envelope signature is recovered")

		// New target rounds may carry the same prepared claim, but must reuse its
		// single cached block and vote set rather than multiplying the evidence.
		for round := uint64(3); round <= 5; round++ {
			before := recoveries
			receiver.receive(signedRoundChange(t, sender, round, claim, cert), errIgnored)
			require.Equal(t, before+1, recoveries, "round %d: only the envelope signature is recovered", round)
			stored := receiver.core.roundChangeSet.Values(new(big.Int).SetUint64(round))
			require.Len(t, stored, 1)
			require.Empty(t, stored[0].Evidence)
		}
		require.Len(t, receiver.core.preparedBlocks, 1)
		require.Len(t, receiver.core.preparedCertificates, 1)
	})
}

// TestConsensusRoundChangeRejectsCachedClaimNotBeforeItsRound checks that a
// cached certificate does not bypass the claim rules. A claim must precede the
// ROUND CHANGE view, or the proposer that selects it cannot justify its round.
func TestConsensusRoundChangeRejectsCachedClaimNotBeforeItsRound(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		x := s.validators[0].proposal(1)
		receiver := s.validators[3]
		claim := &bft.PreparedClaim{Round: big.NewInt(1), Digest: x.Hash()}
		cert := preparedCertificate(t, s.nodes(0, 1, 2), x, 1)
		receiver.receive(signedRoundChange(t, s.validators[0], 2, claim, cert), errIgnored)
		require.Len(t, receiver.core.preparedCertificates, 1)

		// The same verified Evidence bytes hit the cache and skip certificate
		// verification, so only the claim rules can reject these.
		for _, round := range []uint64{0, 1} {
			receiver.reject(signedRoundChange(t, s.validators[1], round, claim, cert), bft.ErrInvalidMessage)
		}
	})
}

// TestConsensusRoundChangeRejectsUnverifiedEvidenceForCachedClaim checks that
// a cached certificate does not vouch for other Evidence bytes. A relay can
// replace the unsigned attachment of another sender's ROUND CHANGE, and an
// accepted message is relayed as received, so only bytes already verified for
// the claim may skip verification.
func TestConsensusRoundChangeRejectsUnverifiedEvidenceForCachedClaim(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		x := s.validators[0].proposal(1)
		receiver := s.validators[3]
		cert := preparedCertificate(t, s.nodes(0, 1, 2), x, 0)
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: x.Hash()}
		receiver.receive(signedRoundChange(t, s.validators[0], 1, claim, cert), errIgnored)
		require.Len(t, receiver.core.preparedCertificates, 1)

		receiver.reject(signedRoundChangeAt(t, s.validators[1], 1, 1, claim, []byte{0xc0}), bft.ErrInvalidMessage)
		receiver.receive(signedRoundChange(t, s.validators[1], 1, claim, cert), errIgnored)
		receiver.assertRoundChangeCount(1, 2)
	})
}

// TestConsensusRoundChangeRejectsCertificateBodyWithoutBlobSidecar checks the
// body rules that votes do not bind. TxHash excludes blob sidecars, so a
// certificate whose blob transaction lost its sidecar keeps the voted hash.
func TestConsensusRoundChangeRejectsCertificateBodyWithoutBlobSidecar(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		sender, receiver := s.validators[0], s.validators[3]
		txs := types.Transactions{types.NewTx(&types.TxInternalDataEthereumBlob{
			GasFeeCap:  uint256.NewInt(1),
			GasLimit:   21000,
			BlobFeeCap: uint256.NewInt(1),
			BlobHashes: []common.Hash{{0x01}},
			V:          big.NewInt(0),
			R:          big.NewInt(1),
			S:          big.NewInt(1),
		})}
		header := sender.proposal(1).Header()
		header.TxHash = types.DeriveTransactionsRoot(txs, header.Number)
		stripped := types.NewBlockWithHeader(header).WithBody(txs)
		cert := preparedCertificate(t, s.nodes(0, 1, 2), stripped, 0)

		target := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}
		require.ErrorIs(t, receiver.core.verifyPreparedCertificate(cert, target), istanbul.ErrNoBlobSidecarForBlobTx)
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: stripped.Hash()}
		receiver.reject(signedRoundChange(t, sender, 1, claim, cert), bft.ErrInvalidMessage)
		require.Empty(t, receiver.core.preparedCertificates)
	})
}

// TestConsensusRecoveryUsesHighestRoundVotesForSameDigest checks that claims
// for one digest at different rounds share the cached block but keep their own
// vote sets, and the proposer attaches the votes of the highest claimed round.
// Receivers verify PreparedMessages at that round, so lower-round votes would
// make the recovery PRE-PREPARE invalid.
func TestConsensusRecoveryUsesHighestRoundVotesForSameDigest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		x := s.validators[0].proposal(1)
		receiver := s.validators[3]
		voters := s.nodes(0, 1, 2)
		lower, higher := preparedCertificate(t, voters, x, 0), preparedCertificate(t, voters, x, 1)
		for _, ev := range []istanbul.MessageEvent{
			signedRoundChange(t, s.validators[0], 2, &bft.PreparedClaim{Round: big.NewInt(0), Digest: x.Hash()}, lower),
			signedRoundChange(t, s.validators[1], 2, &bft.PreparedClaim{Round: big.NewInt(1), Digest: x.Hash()}, higher),
			signedRoundChange(t, s.validators[2], 2, nil, nil),
		} {
			err := receiver.core.handleMsg(ev.Payload)
			synctest.Wait()
			s.checkFailures()
			if err != nil {
				require.ErrorIs(t, err, errIgnored)
			}
		}
		require.Len(t, receiver.core.preparedBlocks, 1)
		require.Len(t, receiver.core.preparedCertificates, 2)
		receiver.assertView(1, 2, false)

		// The quorum moved the receiver to round 2, which keeps both caches.
		target := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(2)}
		_, prepared, err := receiver.core.roundChangeJustification(receiver.core.roundChangeCertificate, target)
		require.NoError(t, err)
		require.NotNil(t, prepared)
		require.Equal(t, int64(1), prepared.View.Round.Int64())
		require.Same(t, receiver.core.preparedBlocks[x.Hash()], prepared.Proposal)
		require.Len(t, prepared.Messages, len(higher.Messages))
		for i, vote := range prepared.Messages {
			require.Equal(t, higher.Messages[i].Msg, vote.Msg)
		}
	})
}

// TestConsensusEvidenceVariantsDoNotExhaustSignerBacklog checks that a relay
// cannot vary the unsigned Evidence of one signed future-height ROUND CHANGE
// to fill its signer's backlog allowance. Every variant fails the signed
// EvidenceHash before it is retained, so the genuine copy and a later vote
// from the same signer are still retained.
func TestConsensusEvidenceVariantsDoNotExhaustSignerBacklog(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		sender, receiver := s.validators[0], s.validators[3]
		// The claim's content is irrelevant until height 2 is reached; only its
		// binding to the attached bytes is checked on receipt.
		evidence := []byte{0xaa, 0xbb}
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: common.HexToHash("0xd1"), EvidenceHash: crypto.Keccak256Hash(evidence)}
		genuine := signedRoundChangeAt(t, sender, 2, 1, claim, evidence)

		var signed bft.Message
		require.NoError(t, signed.FromPayload(genuine.Payload, nil))
		for i := range maxBacklogMessagesPerSender {
			variant := signed
			variant.Evidence = []byte{byte(i)}
			payload, err := variant.Payload()
			require.NoError(t, err)
			receiver.reject(istanbul.MessageEvent{Hash: genuine.Hash, Payload: payload}, bft.ErrInvalidMessage)
		}
		receiver.assertBacklogCount(0)

		receiver.receive(genuine, errFutureMessage)
		encoded, err := bft.Encode(&bft.Prepare{View: &bft.View{Sequence: big.NewInt(2), Round: big.NewInt(0)}, Digest: common.HexToHash("0xd2")})
		require.NoError(t, err)
		prepare, err := sender.core.finalizeMessage(&bft.Message{PrevHash: common.HexToHash("0x01"), Code: bft.MsgPrepare, Msg: encoded})
		require.NoError(t, err)
		receiver.receive(istanbul.MessageEvent{Payload: prepare}, errFutureMessage)
		receiver.assertBacklogCount(2)
	})
}

// TestConsensusReplaysFutureRoundChangeWithCertificate covers the full
// future-height path of a prepared ROUND CHANGE: a node one height behind
// retains it in the backlog, and once it reaches that height the replayed
// message enters the round-change set with its certificate verified and cached.
func TestConsensusReplaysFutureRoundChangeWithCertificate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		sender, receiver := s.validators[0], s.validators[3]
		// Hold the other COMMITs to the receiver so it stays at height 1.
		s.delay(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3))
		s.advanceConsensus(1, s.nodes(3))
		receiver.assertUncommitted(1)
		require.Equal(t, uint64(1), sender.head().NumberU64())

		next := sender.proposal(1)
		require.Equal(t, uint64(2), next.NumberU64())
		cert := preparedCertificate(t, s.nodes(0, 1, 2), next, 0)
		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: next.Hash()}
		receiver.receive(signedRoundChange(t, sender, 1, claim, cert), errFutureMessage)
		receiver.assertBacklogCount(1)
		require.NotZero(t, receiver.core.backlogEvidenceBytes[sender.backend.Address()])

		s.release(bft.MsgCommit, s.nodes(0, 1, 2), s.nodes(3))
		synctest.Wait()
		receiver.assertCommitted(1)
		receiver.assertView(2, 0, false)
		receiver.assertBacklogCount(0)
		require.Zero(t, receiver.core.backlogTotalEvidenceBytes)

		stored := receiver.core.roundChangeSet.Values(big.NewInt(1))
		require.Len(t, stored, 1)
		require.Equal(t, sender.backend.Address(), stored[0].Address)
		require.Empty(t, stored[0].Evidence)
		key := preparedEvidenceKey{sequence: 2, round: 0, digest: next.Hash()}
		require.NotNil(t, receiver.core.preparedCertificates[key], "the replayed certificate is verified and cached")
	})
}

// TestConsensusProposerReproposesNewerOwnLock checks a proposer locked at a
// newer round than every claim in its round-change quorum. Late or dropped
// ROUND CHANGEs can leave only a lower claim in the quorum; proposing that value
// would be refused by the proposer itself and by every node locked at the newer
// round. The proposer must add its own signed claim and re-propose its lock,
// and must not send a later request with that justification.
func TestConsensusProposerReproposesNewerOwnLock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		proposer, receiver := s.validators[2], s.validators[3] // validators[2] proposes height 1, round 2
		x := s.validators[0].proposal(1)
		y := s.validators[0].alternative(x)
		voters := s.nodes(0, 1, 3)
		lower := preparedCertificate(t, voters, y, 0)
		newer := preparedCertificate(t, voters, x, 1)
		proposer.core.current.AdoptPreparedCertificate(newer)
		require.Equal(t, int64(1), proposer.core.current.LockedRound().Int64())

		// The quorum for round 2 carries only the lower claim y@0.
		for _, ev := range []istanbul.MessageEvent{
			signedRoundChange(t, s.validators[0], 2, &bft.PreparedClaim{Round: big.NewInt(0), Digest: y.Hash()}, lower),
			signedRoundChange(t, s.validators[1], 2, nil, nil),
			signedRoundChange(t, s.validators[3], 2, nil, nil),
		} {
			err := proposer.core.handleMsg(ev.Payload)
			synctest.Wait()
			s.checkFailures()
			if err != nil {
				require.ErrorIs(t, err, errIgnored)
			}
		}
		proposer.assertView(1, 2, false)
		require.Equal(t, proposer.backend.Address(), proposer.core.current.proposer)

		preprepare := sentPreprepare(t, s, proposer, 1, 2)
		require.Equal(t, x.Hash(), preprepare.Proposal.Hash(), "the proposer re-proposes its newer lock")
		ownClaim := false
		for _, rc := range preprepare.RoundChangeCertificate {
			var roundChange *bft.RoundChange
			require.NoError(t, rc.Decode(&roundChange))
			if rc.Address == proposer.backend.Address() {
				require.NotNil(t, roundChange.Prepared)
				require.Equal(t, int64(1), roundChange.Prepared.Round.Int64())
				require.Equal(t, x.Hash(), roundChange.Prepared.Digest)
				ownClaim = true
			}
		}
		require.True(t, ownClaim, "the justification carries the proposer's own claim")
		cert, err := receiver.core.verifyPreprepareJustification(preprepare)
		require.NoError(t, err, "other validators accept the justification")
		require.Equal(t, x.Hash(), cert.Proposal.Hash())
		require.Equal(t, int64(1), cert.View.Round.Int64())

		// A late request for another block must not go out with this justification.
		sentPreprepares := func() int {
			count := 0
			for _, sent := range s.sent {
				if sent.from != proposer.id {
					continue
				}
				msg, err := proposer.backend.decodeMessage(sent.data.(istanbul.MessageEvent).Payload, nil)
				require.NoError(t, err)
				if msg.Code == bft.MsgPreprepare {
					count++
				}
			}
			return count
		}
		before := sentPreprepares()
		proposer.core.sendPreprepare(&bft.Request{Proposal: y})
		synctest.Wait()
		require.Equal(t, before, sentPreprepares())
	})
}

// TestRoundChangeSetKeepsOneMessagePerCommitteeMember checks that a round keeps
// up to one ROUND CHANGE per committee member rather than a quorum. Claims carry
// block-sized Evidence and tend to arrive last, so a quorum-sized limit could
// drop the highest prepared claim.
func TestRoundChangeSetKeepsOneMessagePerCommitteeMember(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		c := s.validators[0].core
		require.Equal(t, 3, c.current.requiredMessageCount)
		require.Equal(t, c.current.committee.Len(), c.roundChangeSet.maxMessagesPerRound)
	})
}

// preparedCertificate builds a certificate for proposal at round from PREPARE
// votes signed by voters.
func preparedCertificate(t *testing.T, voters []*validator, proposal *types.Block, round uint64) *bft.PreparedCertificate {
	messages := make([]*bft.Message, 0, len(voters))
	for _, voter := range voters {
		var message bft.Message
		require.NoError(t, message.FromPayload(voter.message(bft.MsgPrepare, proposal, round).Payload, nil))
		messages = append(messages, &message)
	}
	return &bft.PreparedCertificate{
		View:     &bft.View{Sequence: proposal.Number(), Round: new(big.Int).SetUint64(round)},
		Proposal: proposal,
		Messages: messages,
	}
}

// TestPreparedCertificateAcceptsCommitVotes covers the certificate branch in
// which some quorum members contribute COMMIT rather than PREPARE envelopes.
// The outer vote signature and the round-bound committed seal must both match.
func TestPreparedCertificateAcceptsCommitVotes(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newScenarioNet(t, 4, 4, params.TestKaiaConfig("permissionless"))
		proposal := s.validators[0].proposal(1)
		messages := make([]*bft.Message, 0, 3)
		for i, code := range []uint64{bft.MsgPrepare, bft.MsgPrepare, bft.MsgCommit} {
			event := s.validators[i].message(code, proposal, 0)
			var message bft.Message
			require.NoError(t, message.FromPayload(event.Payload, nil))
			messages = append(messages, &message)
		}
		cert := &bft.PreparedCertificate{
			View:     &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)},
			Proposal: proposal,
			Messages: messages,
		}
		target := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}
		require.NoError(t, s.validators[3].core.verifyPreparedCertificate(cert, target))

		wrongCommit := *messages[2]
		var commit bft.Commit
		require.NoError(t, wrongCommit.Decode(&commit))
		var err error
		commit.CommittedSeal, err = s.validators[2].backend.sealer.MakeCommittedSealFromHashWithRound(proposal.Hash(), 1)
		require.NoError(t, err)
		wrongCommit.Msg, err = bft.Encode(&commit)
		require.NoError(t, err)
		unsigned, err := wrongCommit.PayloadNoSig()
		require.NoError(t, err)
		wrongCommit.Signature, err = s.validators[2].backend.Sign(unsigned)
		require.NoError(t, err)
		badCert := *cert
		badCert.Messages = []*bft.Message{messages[0], messages[1], &wrongCommit}
		require.ErrorContains(t, s.validators[3].core.verifyPreparedCertificate(&badCert, target), "invalid committed seal")
	})
}

// TestConsensusRoundChangeRejectsEvidenceBeforeFork checks that before
// Permissionless a ROUND CHANGE keeps the legacy shape. The legacy codec cannot
// express a PreparedClaim or Evidence, so local senders cannot produce one, and
// a peer's post-Permissionless envelope for a pre-fork height is rejected
// before it can enter the round-change set or the backlog. From the
// activation height a signed PreparedClaim and its attachment must appear
// together before a future message can enter the backlog.
func TestConsensusRoundChangeRejectsEvidenceBeforeFork(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		config := params.TestKaiaConfig("permissionless")
		config.PermissionlessCompatibleBlock = big.NewInt(3)
		s := newScenarioNet(t, 4, 4, config)
		sender, receiver := s.validators[0], s.validators[1]
		require.False(t, receiver.backend.IsPermissionlessAt(1))

		claim := &bft.PreparedClaim{Round: big.NewInt(0), Digest: common.HexToHash("0x01")}
		junk := []byte{0xc0}
		for _, tc := range []struct {
			name     string
			sequence uint64
			claim    *bft.PreparedClaim
			evidence []byte
		}{
			{"plain", 1, nil, nil},
			{"evidence without claim", 1, nil, junk},
			{"claim without evidence", 1, claim, nil},
			{"claim with evidence", 1, claim, junk},
			{"future-height evidence", 2, nil, junk},
		} {
			t.Log(tc.name)
			if tc.claim != nil || tc.evidence != nil {
				_, err := sender.core.finalizeMessage(roundChangeMessageAt(t, sender, tc.sequence, 1, tc.claim, tc.evidence))
				require.ErrorIs(t, err, bft.ErrInvalidMessage, "legacy codec must not drop a claim or evidence")
			}
			receiver.reject(postPermissionlessRoundChange(t, sender, tc.sequence, 1, tc.claim, tc.evidence))
		}
		// The legacy ROUND CHANGE from the same sender is still admitted.
		receiver.receive(signedRoundChangeAt(t, sender, 1, 1, nil, nil), errIgnored)
		receiver.assertRoundChangeCount(1, 1)
		// At the activation height the malformed envelope is rejected before the
		// future-message path can retain its unsigned attachment. A claim-less
		// ROUND CHANGE without an attachment remains compatible and deferrable.
		require.True(t, receiver.backend.IsPermissionlessAt(3))
		receiver.reject(signedRoundChangeAt(t, sender, 3, 1, nil, junk), bft.ErrInvalidMessage)
		receiver.receive(signedRoundChangeAt(t, sender, 3, 1, nil, nil), errFutureMessage)
		receiver.assertBacklogCount(1)
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
	}{
		{"OtherRoundSealBeforeActivation", 2},
		{"OtherRoundSealAfterActivation", 1},
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
				// A post-fork round-bound seal is invalid before activation, and a seal for another round is invalid after it.
				target.reject(ev, errInvalidCommittedSeal)
				target.assertUncommitted(1)
				s.release(bft.MsgCommit, s.nodes(2, 3), s.nodes(0))
				target.assertCommitted(1)
			})
		})
	}
}
