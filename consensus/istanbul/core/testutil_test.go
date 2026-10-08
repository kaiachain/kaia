// Copyright 2024 The Kaia Authors
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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/derivesha"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kaiax/gov"
	mock_gov "github.com/kaiachain/kaia/kaiax/gov/mock"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/require"
)

// Real event loops and timers run inside synctest; only worker input and network delivery are controlled.
const maxScenarioSteps = 10000

type scenarioNet struct {
	mu                sync.Mutex // Serializes network output and cross-node commit checks from real core goroutines.
	t                 *testing.T
	validators        []*validator
	config            scenarioConfig
	pendingDeliveries []scenarioEvent
	rules             []*messageRule
	sent              []scenarioEvent
	proposal          *types.Block // last PREPREPARE broadcast by any node
	failures          []error      // callback failures surfaced by the scenario-driving goroutine
}
type scenarioConfig struct {
	committeeSize  int
	chainConfig    *params.ChainConfig
	istanbulConfig *istanbul.Config
}

func newScenarioNet(t *testing.T, count, committee int, chain *params.ChainConfig) *scenarioNet {
	t.Helper()
	require.Positive(t, committee)
	require.GreaterOrEqual(t, count, committee)
	config := istanbul.DefaultConfig.Copy()
	config.ProposerPolicy = istanbul.RoundRobin
	s := &scenarioNet{t: t, config: scenarioConfig{committee, chain, config}}
	previousHash := types.HeaderHashFn
	previousDeriveSha, previousEmptyRoot := types.DeriveShaImpl, types.GetEmptyRootHash
	t.Cleanup(func() {
		for _, n := range s.validators {
			require.NoError(t, n.core.Stop())
			n.backend.mux.Stop()
		}
		fork.ClearHardForkBlockNumberConfig()
		types.SetHeaderHashFn(previousHash)
		types.DeriveShaImpl, types.GetEmptyRootHash = previousDeriveSha, previousEmptyRoot
	})
	// Prepared certificates bind a proposal's body to its TxHash.
	derivesha.InitDeriveSha(chain, nil)
	fork.ClearHardForkBlockNumberConfig()
	require.NoError(t, fork.SetHardForkBlockNumberConfig(chain))
	keys := make([]*ecdsa.PrivateKey, count)
	addresses := make([]common.Address, count)
	for i := range keys {
		var err error
		keys[i], err = crypto.ToECDSA(common.LeftPadBytes(big.NewInt(int64(i+1)).Bytes(), 32))
		require.NoError(t, err)
		addresses[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
	}
	sealer := istanbul.NewSealerImpl(keys[0])
	types.SetHeaderHashFn(sealer.HeaderHash)
	header := &types.Header{Number: new(big.Int), Time: big.NewInt(1), BlockScore: big.NewInt(1)}
	require.NoError(t, sealer.WriteValidators(header, addresses))
	for i, key := range keys {
		s.validators = append(s.validators, newValidator(s, i, key, types.NewBlockWithHeader(header), addresses))
	}
	synctest.Wait()
	s.checkFailures()
	return s
}

func (s *scenarioNet) recordFailure(err error) {
	if err == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures = append(s.failures, err)
}

func (s *scenarioNet) checkFailures() {
	s.mu.Lock()
	failures := append([]error(nil), s.failures...)
	s.mu.Unlock()
	require.NoError(s.t, errors.Join(failures...), "scenario backend invariant")
}

func (s *scenarioNet) nodes(ids ...int) []*validator {
	s.t.Helper()
	require.NotEmpty(s.t, ids)
	nodes := make([]*validator, len(ids))
	for i, id := range ids {
		require.GreaterOrEqual(s.t, id, 0)
		require.Less(s.t, id, len(s.validators))
		nodes[i] = s.validators[id]
	}
	return nodes
}

func (s *scenarioNet) ids() []int {
	ids := make([]int, len(s.validators))
	for i := range ids {
		ids[i] = i
	}
	return ids
}

func (net *scenarioNet) validatorAddresses() []common.Address {
	addresses := make([]common.Address, len(net.validators))
	for i, validator := range net.validators {
		addresses[i] = validator.backend.Address()
	}
	return addresses
}

func (net *scenarioNet) requireNode(node *validator) {
	net.t.Helper()
	require.NotNil(net.t, node)
	require.Same(net.t, net, node.backend.net, "node[%d] belongs to another scenario", node.id)
}

// advanceConsensus supplies this height's workers and drains events; excluded nodes skip only the final commit check.
func (s *scenarioNet) advanceConsensus(height uint64, skip ...[]*validator) {
	s.t.Helper()
	require.NotZero(s.t, height)
	require.LessOrEqual(s.t, len(skip), 1, "pass at most one excluded-node list")
	excluded := make(map[*validator]bool)
	if len(skip) == 1 {
		for _, node := range skip[0] {
			s.requireNode(node)
			require.False(s.t, excluded[node], "duplicate excluded node %d", node.id)
			excluded[node] = true
		}
	}
	for steps := 0; ; steps++ {
		require.Less(s.t, steps, maxScenarioSteps, "consensus did not quiesce")
		if len(s.pendingDeliveries) == 0 && !s.requestProposal(height) {
			break
		}
		s.step()
	}
	for _, n := range s.validators {
		if !excluded[n] {
			n.assertCommitted(height)
		}
	}
}

// Each ready worker submits once per sequence, including non-proposers; waiting nodes receive no new request.
func (s *scenarioNet) requestProposal(height uint64) bool {
	queued := false
	for _, n := range s.validators {
		c := n.core
		if c.current.Sequence().Uint64() == height && !c.waitingForRoundChange && c.current.pendingRequest == nil {
			proposal := n.proposal(1)
			s.pendingDeliveries = append(s.pendingDeliveries, scenarioEvent{n.id, n.id, istanbul.RequestEvent{Proposal: proposal}})
			queued = true
		}
	}
	return queued
}

func (s *scenarioNet) step() {
	s.t.Helper()
	require.NotEmpty(s.t, s.pendingDeliveries)
	ev := s.pendingDeliveries[0]
	s.pendingDeliveries = s.pendingDeliveries[1:]
	node := s.validators[ev.to]
	before := node.backend.view
	node.post(ev.data)
	after := node.backend.view
	if before.Sequence.Cmp(after.Sequence) == 0 && before.Round.Cmp(after.Round) != 0 {
		require.Positive(s.t, after.Round.Cmp(before.Round), "round regression")
		proposer := s.validators[(after.Sequence.Uint64()-1+after.Round.Uint64())%uint64(s.config.committeeSize)]
		require.Equal(s.t, proposer.backend.Address(), node.core.current.proposer)
		require.Zero(s.t, after.Cmp(node.core.currentView()))
		require.False(s.t, node.core.waitingForRoundChange)
	}
}

func (s *scenarioNet) drain() {
	s.t.Helper()
	for steps := 0; len(s.pendingDeliveries) > 0; steps++ {
		require.Less(s.t, steps, maxScenarioSteps, "message loop")
		s.step()
	}
}

// until is reserved for the worker-request/head-event ordering scenario.
func (s *scenarioNet) until(height uint64, reached func() bool) {
	for steps := 0; !reached(); steps++ {
		require.Less(s.t, steps, maxScenarioSteps)
		if len(s.pendingDeliveries) == 0 {
			require.True(s.t, s.requestProposal(height))
		}
		s.step()
	}
}

// timeout shortens selected deadlines and lets the real timer callbacks post their events.
func (s *scenarioNet) timeout(nodes []*validator) {
	s.t.Helper()
	for _, n := range nodes {
		s.requireNode(n)
		timer := n.core.roundChangeTimer.Load().(*time.Timer)
		require.True(s.t, timer.Reset(time.Nanosecond), "node %d has no armed timer", n.id)
	}
	time.Sleep(time.Nanosecond)
	synctest.Wait()
	s.checkFailures()
}

// delay holds directed routes, optionally for one round; holding self-delivery also postpones peer forwarding.
func (s *scenarioNet) delay(code uint64, from, to []*validator, round ...uint64) {
	s.t.Helper()
	require.NotEmpty(s.t, from)
	require.NotEmpty(s.t, to)
	for _, sender := range from {
		for _, recipient := range to {
			newMessageRule(s, code, sender, []*validator{recipient}, round...).hold = true
		}
	}
}

// release disables matching delays and drains stored packets without creating worker requests or timeouts.
func (s *scenarioNet) release(code uint64, from, to []*validator, round ...uint64) {
	s.t.Helper()
	require.LessOrEqual(s.t, len(round), 1)
	count := 0
	for _, rule := range s.rules {
		if rule.code == code && rule.hold && slices.Contains(from, rule.from) && slices.Contains(to, rule.recipients[0]) &&
			(len(round) == 0 || rule.round != nil && *rule.round == round[0]) {
			s.pendingDeliveries = append(s.pendingDeliveries, rule.held...)
			count += len(rule.held)
			rule.held = nil
			rule.hold = false
		}
	}
	require.Positive(s.t, count, "no matching held messages")
	s.drain()
}

// isolate drops peer traffic in both directions; local processing continues and no core is stopped.
func (s *scenarioNet) isolate(nodes []*validator) {
	for _, from := range s.validators {
		for _, to := range s.validators {
			if from != to && (slices.Contains(nodes, from) || slices.Contains(nodes, to)) {
				for _, code := range []uint64{bft.MsgPreprepare, bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange} {
					newMessageRule(s, code, from, []*validator{to}).drop = true
				}
			}
		}
	}
}

func (net *scenarioNet) modify(code uint64, from *validator, recipients []*validator, proposal *types.Block) {
	net.t.Helper()
	require.NotNil(net.t, proposal)
	require.True(net.t, code != bft.MsgRoundChange, "ROUND CHANGE has no proposal digest")
	rule := newMessageRule(net, code, from, recipients)
	rule.proposal = proposal
	net.t.Cleanup(func() {
		require.True(net.t, rule.applied, "modify did not match code %d from node %d", code, from.id)
	})
}

// drop discards messages on the selected directed routes for the scenario.
func (s *scenarioNet) drop(code uint64, from, to []*validator, round ...uint64) {
	s.t.Helper()
	require.NotEmpty(s.t, from)
	require.NotEmpty(s.t, to)
	for _, sender := range from {
		for _, recipient := range to {
			newMessageRule(s, code, sender, []*validator{recipient}, round...).drop = true
		}
	}
}

// message returns bytes actually broadcast by this node, not fabricated honest votes.
func (s *scenarioNet) message(from *validator, code, height, round uint64) istanbul.MessageEvent {
	s.t.Helper()
	for _, sent := range s.sent {
		if sent.from != from.id {
			continue
		}
		ev := sent.data.(istanbul.MessageEvent)
		msg, err := from.backend.decodeMessage(ev.Payload, nil)
		require.NoError(s.t, err)
		view, err := msg.GetView()
		require.NoError(s.t, err)
		if msg.Code == code && view.Sequence.Uint64() == height && view.Round.Uint64() == round {
			return ev
		}
	}
	s.t.Fatalf("node %d did not broadcast code %d at %d/%d", from.id, code, height, round)
	return istanbul.MessageEvent{}
}

// Validator assertions read the actual core or stored block, never a parallel model.
type validator struct {
	id      int
	core    *core
	backend *scenarioBackend
}

func newValidator(net *scenarioNet, id int, key *ecdsa.PrivateKey, genesis *types.Block, council []common.Address) *validator {
	net.t.Helper()
	backend := &scenarioBackend{
		net: net, id: id, key: key, sealer: istanbul.NewSealerImpl(key),
		mux: new(event.TypeMux), head: genesis, blocks: map[uint64]*types.Block{0: genesis},
		chainConfig: net.config.chainConfig, badProposals: make(map[common.Hash]bool),
	}

	c := New(backend, net.config.istanbulConfig.Copy()).(*core)
	validator := &validator{id: id, core: c, backend: backend}
	// Fixed committee with round-robin proposers; no validator-set churn.
	committeeSize := net.config.committeeSize
	ctrl := gomock.NewController(net.t)
	net.t.Cleanup(ctrl.Finish)
	validators := valset_mock.NewMockValsetModule(ctrl)
	governance := mock_gov.NewMockGovModule(ctrl)
	committee := council[:committeeSize]
	validators.EXPECT().GetCouncil(gomock.Any()).Return(council, nil).AnyTimes()
	validators.EXPECT().GetDemotedValidators(gomock.Any()).Return([]common.Address{}, nil).AnyTimes()
	validators.EXPECT().GetCommittee(gomock.Any(), gomock.Any()).Return(committee, nil).AnyTimes()
	validators.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(func(height, round uint64) (common.Address, error) {
		return committee[(height-1+round)%uint64(committeeSize)], nil
	}).AnyTimes()
	governance.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{CommitteeSize: uint64(committeeSize)}).AnyTimes()
	c.RegisterKaiaxModules(validators, governance)
	require.NoError(net.t, c.Start())
	synctest.Wait()
	require.NotNil(net.t, c.current)

	return validator
}

// assertCommitted checks the latest height after local consensus, not a head copied by follow.
func (n *validator) assertCommitted(height uint64, round ...uint64) *types.Block {
	t := n.backend.net.t
	t.Helper()
	require.NotZero(t, height)
	require.Len(t, n.backend.committed, int(height), "node %d commit count", n.id)
	block := n.backend.committed[height-1]
	require.Equal(t, height, block.NumberU64())
	require.Equal(t, block.Hash(), n.backend.head.Hash())
	n.assertView(height+1, 0, false)
	require.Equal(t, StateAcceptRequest, n.core.state)
	if len(round) > 0 {
		require.Len(t, round, 1)
		require.Equal(t, round[0], n.committedRound(height))
	}
	return block
}

func (n *validator) assertUncommitted(height uint64) {
	n.backend.net.t.Helper()
	require.Less(n.backend.net.t, n.backend.head.NumberU64(), height, "node %d unexpectedly committed", n.id)
}

func (n *validator) assertView(height, round uint64, waiting bool) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, height, n.core.current.Sequence().Uint64(), "node %d sequence", n.id)
	require.Equal(t, round, n.core.current.Round().Uint64(), "node %d round", n.id)
	require.Equal(t, waiting, n.core.waitingForRoundChange, "node %d waiting", n.id)
}

// assertHashLocked checks the locked hash and flag; a zero hash requires no lock.
func (n *validator) assertHashLocked(hash common.Hash) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, hash, n.core.current.GetLockedHash(), "node %d lock", n.id)
	require.Equal(t, hash != (common.Hash{}), n.core.current.IsHashLocked())
}

func (n *validator) assertBadHashLocked(hash common.Hash) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, hash, n.core.current.GetLockedHash(), "node %d lock", n.id)
	require.False(t, n.core.current.IsHashLocked(), "node %d bad proposal must invalidate the lock", n.id)
}

// assertVoteCounts checks PREPAREs, COMMITs and distinct signers across both sets in the current view.
func (n *validator) assertVoteCounts(prepares, commits, distinct int) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, prepares, n.core.current.Prepares.Size())
	require.Equal(t, commits, n.core.current.Commits.Size())
	require.Equal(t, distinct, n.core.current.GetPrepareOrCommitSize())
}

func (n *validator) assertRoundChangeCount(round uint64, count int) {
	t := n.backend.net.t
	t.Helper()
	set := n.core.roundChangeSet.roundChanges[round]
	if count == 0 {
		require.Nil(t, set)
		return
	}
	require.NotNil(t, set)
	require.Equal(t, count, set.Size())
}

// post delivers through the production event loop and waits until its work is blocked again.
func (n *validator) post(ev interface{}) {
	net := n.backend.net
	net.t.Helper()
	require.NoError(net.t, n.backend.mux.Post(ev))
	synctest.Wait()
	net.checkFailures()
}

// receive checks the message handler's return value while the event loop is idle.
func (n *validator) receive(ev istanbul.MessageEvent, want error) {
	t := n.backend.net.t
	t.Helper()
	synctest.Wait()
	err := n.core.handleMsg(ev.Payload)
	synctest.Wait()
	n.backend.net.checkFailures()
	require.ErrorIs(t, err, want)
}

// reject requires an error and unchanged state; omit want only for decoder errors without a fixed core contract.
func (n *validator) reject(event istanbul.MessageEvent, want ...error) {
	t := n.backend.net.t
	t.Helper()
	require.LessOrEqual(t, len(want), 1)
	n.assertStateUnchanged(func() {
		err := n.core.handleMsg(event.Payload)
		require.Error(t, err)
		require.NotEqual(t, errFutureMessage, err, "deferral is not rejection")
		if len(want) == 1 {
			require.ErrorIs(t, err, want[0])
		}
	})
}

func (n *validator) assertQuorum(messages, faults int) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, messages, n.core.current.requiredMessageCount)
	require.Equal(t, faults, n.core.current.f)
}

func (n *validator) assertPendingRequestCount(count int) {
	n.backend.net.t.Helper()
	require.Equal(n.backend.net.t, count, n.core.pendingRequests.Size())
}

func (n *validator) assertBacklogCount(count int) {
	n.backend.net.t.Helper()
	queued := len(n.core.backlogPreprepares)
	for _, backlog := range n.core.backlogs {
		queued += backlog.Size()
	}
	require.Equal(n.backend.net.t, count, queued, "node %d backlog messages", n.id)
}
func (n *validator) failNextCommit(err error) { n.backend.commitErrors = []error{err} }
func (n *validator) commitAttempts(count int) {
	t := n.backend.net.t
	t.Helper()
	require.Equal(t, count, n.backend.commitCalls)
	require.Empty(t, n.backend.commitErrors, "configured fault must have been consumed")
}
func (n *validator) head() *types.Block { return n.backend.head }
func (n *validator) block(height uint64) *types.Block {
	n.backend.net.t.Helper()
	block := n.backend.blocks[height]
	require.NotNil(n.backend.net.t, block, "node %d missing height %d", n.id, height)
	return block
}

func (n *validator) alternative(base *types.Block) *types.Block {
	t := n.backend.net.t
	t.Helper()
	block := n.proposal(2)
	require.NotEqual(t, base.Hash(), block.Hash())
	_, err := n.backend.Verify(block)
	require.NoError(t, err, "conflict must be an independently valid proposal")
	return block
}

func (n *validator) assertStateUnchanged(action func()) {
	t := n.backend.net.t
	t.Helper()
	before := n.snapshot()
	action()
	require.Equal(t, before, n.snapshot(), "node %d changed", n.id)
}

func (node *validator) committedRound(height uint64) uint64 {
	t := node.backend.net.t
	t.Helper()
	block := node.backend.blocks[height]
	require.NotNil(t, block, "node %d: no block at height %d", node.id, height)
	round, err := node.backend.sealer.Round(block.Header())
	require.NoError(t, err, "node %d: cannot read round at height %d", node.id, height)
	return uint64(round)
}

func (node *validator) proposal(variant int64) *types.Block {
	net := node.backend.net
	net.t.Helper()
	header := &types.Header{
		ParentHash: node.backend.head.Hash(), Number: new(big.Int).Add(node.backend.head.Number(), common.Big1),
		Time: new(big.Int).Add(node.backend.head.Time(), big.NewInt(variant)), BlockScore: big.NewInt(1),
	}
	header.TxHash = types.DeriveTransactionsRoot(nil, header.Number)
	require.NoError(net.t, node.backend.sealer.WriteValidators(header, net.validatorAddresses()))
	seal, err := node.backend.sealer.MakeAuthorSeal(header)
	require.NoError(net.t, err)
	require.NoError(net.t, node.backend.sealer.WriteAuthorSeal(header, seal))
	return types.NewBlockWithHeader(header)
}

// message signs a synthetic faulty input; the sender's real core outputs still undergo all invariant checks.
func (node *validator) message(code uint64, proposal *types.Block, round uint64) istanbul.MessageEvent {
	net := node.backend.net
	net.t.Helper()
	event, err := node.makeMessage(code, proposal, round)
	require.NoError(net.t, err)
	return event
}

func (node *validator) makeMessage(code uint64, proposal *types.Block, round uint64) (istanbul.MessageEvent, error) {
	view := &bft.View{Sequence: proposal.Number(), Round: new(big.Int).SetUint64(round)}
	var payload any
	switch code {
	case bft.MsgPreprepare:
		payload = &bft.Preprepare{View: view, Proposal: proposal}
	case bft.MsgPrepare:
		payload = &bft.Prepare{View: view, Digest: proposal.Hash()}
	case bft.MsgCommit:
		payload = &bft.Commit{View: view, Digest: proposal.Hash()}
	case bft.MsgRoundChange:
		payload = &bft.RoundChange{View: view}
	default:
		return istanbul.MessageEvent{}, fmt.Errorf("unknown consensus message code %d", code)
	}
	encoded, err := bft.Encode(payload)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	signed, err := node.core.finalizeMessage(&bft.Message{PrevHash: proposal.ParentHash(), Code: code, Msg: encoded})
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	return istanbul.MessageEvent{Hash: proposal.ParentHash(), Payload: signed}, nil
}

// timeout expires this node's real timer early; other nodes retain their own deadlines.
func (node *validator) timeout() { node.backend.net.timeout([]*validator{node}) }

// scenarioBackend records core commits as-is; validator assertions expose invalid quorum or conflicting blocks.
type scenarioBackend struct {
	net          *scenarioNet
	id           int
	key          *ecdsa.PrivateKey
	sealer       *istanbul.IstanbulSealer
	mux          *event.TypeMux
	view         *bft.View
	head         *types.Block
	blocks       map[uint64]*types.Block
	committed    []*types.Block
	chainConfig  *params.ChainConfig
	badProposals map[common.Hash]bool
	commitErrors []error
	commitCalls  int
}

var _ istanbul.Backend = (*scenarioBackend)(nil)

func (b *scenarioBackend) Address() common.Address          { return crypto.PubkeyToAddress(b.key.PublicKey) }
func (b *scenarioBackend) Sealer() *istanbul.IstanbulSealer { return b.sealer }
func (b *scenarioBackend) EventMux() *event.TypeMux         { return b.mux }
func (b *scenarioBackend) NodeType() common.ConnType        { return common.CONSENSUSNODE }
func (b *scenarioBackend) HasBadProposal(hash common.Hash) bool {
	return b.badProposals[hash]
}

func (b *scenarioBackend) IsPermissionlessAt(number uint64) bool {
	return b.chainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(number))
}

func (b *scenarioBackend) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(data), b.key)
}

func (b *scenarioBackend) SetCurrentView(view *bft.View) {
	b.view = &bft.View{Sequence: new(big.Int).Set(view.Sequence), Round: new(big.Int).Set(view.Round)}
}

func (b *scenarioBackend) Broadcast(hash common.Hash, payload []byte) error {
	err := b.fanout(hash, payload, true)
	b.net.recordFailure(err)
	return err
}

func (b *scenarioBackend) Gossip([]byte) error {
	err := errors.New("unexpected Gossip: core broadcasts must pass through self-delivery")
	b.net.recordFailure(err)
	return err
}

// decodeMessage decodes payload with the wire codec its sequence selects on this chain.
func (b *scenarioBackend) decodeMessage(payload []byte, validateFn func([]byte, []byte) (common.Address, error)) (*bft.Message, error) {
	var msg bft.Message
	if err := msg.FromPayloadForFork(payload, b.IsPermissionlessAt, validateFn); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Only successful self-processing forwards to direct peers; peer relays and P2P caches are not modeled.
func (b *scenarioBackend) GossipSubPeer(hash common.Hash, payload []byte) {
	msg, err := b.decodeMessage(payload, nil)
	if err != nil {
		b.net.recordFailure(fmt.Errorf("node %d decode relayed message: %w", b.id, err))
		return
	}
	if msg.Address == b.Address() {
		b.net.recordFailure(b.fanout(hash, payload, false))
	}
}

// fanout applies recipient rules to signed core outputs and queues them for the normal message handler.
func (b *scenarioBackend) fanout(hash common.Hash, payload []byte, self bool) error {
	b.net.mu.Lock()
	defer b.net.mu.Unlock()
	msg, err := b.decodeMessage(payload, nil)
	if err != nil {
		return fmt.Errorf("node %d decode broadcast message: %w", b.id, err)
	}
	if self {
		b.net.sent = append(b.net.sent, scenarioEvent{b.id, b.id, istanbul.MessageEvent{Hash: hash, Payload: slices.Clone(payload)}})
		if msg.Code == bft.MsgPreprepare {
			var pp bft.Preprepare
			if err := msg.Decode(&pp); err != nil {
				return fmt.Errorf("node %d decode PREPREPARE: %w", b.id, err)
			}
			proposal, ok := pp.Proposal.(*types.Block)
			if !ok {
				return fmt.Errorf("node %d broadcast non-block PREPREPARE", b.id)
			}
			b.net.proposal = proposal
		}
	}
	if self && msg.Code == bft.MsgCommit {
		var commit bft.Commit
		if err := msg.Decode(&commit); err != nil {
			return fmt.Errorf("node %d decode COMMIT: %w", b.id, err)
		}
		current := b.net.validators[b.id].core.current
		if commit.View.Sequence.Cmp(current.Sequence()) == 0 {
			if !current.IsHashLocked() {
				return fmt.Errorf("node %d: honest COMMIT without a lock", b.id)
			}
			if current.GetLockedHash() != commit.Digest {
				return fmt.Errorf("node %d: COMMIT conflicts with lock: have %s, want %s", b.id, commit.Digest, current.GetLockedHash())
			}
		}
	}
	view, err := msg.GetView()
	if err != nil {
		return fmt.Errorf("node %d read message view: %w", b.id, err)
	}
	for to, recipient := range b.net.validators {
		if (self && to != b.id) || (!self && to == b.id) {
			continue
		}
		ev := scenarioEvent{b.id, to, istanbul.MessageEvent{Hash: hash, Payload: append([]byte(nil), payload...)}}
		deliver := true
		for _, rule := range b.net.rules {
			if rule.code != msg.Code || rule.from.id != b.id || !slices.Contains(rule.recipients, recipient) ||
				(rule.round != nil && *rule.round != view.Round.Uint64()) {
				continue
			}
			rule.applied = true
			switch {
			case rule.drop:
				deliver = false
			case rule.hold:
				rule.held = append(rule.held, ev)
				deliver = false
			case rule.proposal != nil:
				if msg.Code == bft.MsgPreprepare {
					var pp *bft.Preprepare
					if err := msg.Decode(&pp); err != nil {
						return fmt.Errorf("node %d decode modified PREPREPARE: %w", b.id, err)
					}
					if pp.View.Sequence.Cmp(rule.proposal.Number()) != 0 {
						return fmt.Errorf("modified proposal must keep the message's sequence: have %s, want %s", rule.proposal.Number(), pp.View.Sequence)
					}
					// An equivocating proposer changes only the block. Preserve a real
					// higher-round ROUND-CHANGE certificate so receivers can test that
					// it binds the proposal rather than merely being present.
					pp.Proposal = rule.proposal
					encoded, err := bft.Encode(pp)
					if err != nil {
						return fmt.Errorf("node %d encode modified PREPREPARE: %w", b.id, err)
					}
					payload, err := rule.from.core.finalizeMessage(&bft.Message{PrevHash: rule.proposal.ParentHash(), Code: msg.Code, Msg: encoded})
					if err != nil {
						return fmt.Errorf("node %d create modified PREPREPARE: %w", b.id, err)
					}
					ev.data = istanbul.MessageEvent{Hash: rule.proposal.ParentHash(), Payload: payload}
				} else {
					if view.Sequence.Cmp(rule.proposal.Number()) != 0 {
						return fmt.Errorf("modified proposal must keep the message's sequence: have %s, want %s", rule.proposal.Number(), view.Sequence)
					}
					modified, err := rule.from.makeMessage(msg.Code, rule.proposal, view.Round.Uint64())
					if err != nil {
						return fmt.Errorf("node %d create modified message: %w", b.id, err)
					}
					ev.data = modified
				}
			}
		}
		if deliver {
			b.net.pendingDeliveries = append(b.net.pendingDeliveries, ev)
		}
	}
	return nil
}

func (b *scenarioBackend) LastProposal() (bft.Proposal, common.Address) {
	author, _ := b.sealer.Author(b.head.Header()) // genesis has no proposer seal
	return b.head, author
}

func (b *scenarioBackend) ProposalRound(hash common.Hash, number *big.Int) (byte, bool) {
	block := b.blocks[number.Uint64()]
	if block == nil || block.Hash() != hash {
		return 0, false
	}
	round, err := b.sealer.Round(block.Header())
	if err != nil {
		return 0, false
	}
	return round, true
}

func (b *scenarioBackend) HasPropsal(hash common.Hash, height *big.Int) bool {
	block := b.blocks[height.Uint64()]
	return block != nil && block.Hash() == hash
}

// Verify checks empty-block parent, height and author; execution and full header validation belong to backend tests.
func (b *scenarioBackend) Verify(proposal bft.Proposal) (time.Duration, error) {
	block, ok := proposal.(*types.Block)
	if !ok || block.NumberU64() != b.head.NumberU64()+1 || block.ParentHash() != b.head.Hash() {
		return 0, istanbul.ErrInvalidProposal
	}
	if err := b.VerifyProposalBody(block); err != nil {
		return 0, err
	}
	author, err := b.sealer.Author(block.Header())
	if err != nil {
		return 0, err
	}
	if !slices.Contains(b.net.validatorAddresses(), author) {
		return 0, istanbul.ErrUnauthorizedAddress
	}
	return 0, nil
}

func (b *scenarioBackend) VerifyProposalBody(proposal bft.Proposal) error {
	block, ok := proposal.(*types.Block)
	if !ok {
		return istanbul.ErrInvalidProposal
	}
	if b.net.config.chainConfig.IsOsakaForkEnabled(block.Number()) && block.Size() > params.MaxBlockSize {
		return blockchain.ErrBlockOversized
	}
	header := block.Header()
	if types.DeriveTransactionsRoot(block.Transactions(), block.Number()) != header.TxHash {
		return istanbul.ErrMismatchTxhashes
	}
	var blobs int
	for _, tx := range block.Transactions() {
		if header.BaseFee != nil && header.BaseFee.Cmp(tx.GasPrice()) > 0 {
			return fmt.Errorf("invalid GasPrice: txHash %x", tx.Hash())
		}
		blobs += len(tx.BlobHashes())
		if tx.Type() == types.TxTypeEthereumBlob {
			sidecar := tx.BlobTxSidecar()
			if sidecar == nil {
				return istanbul.ErrNoBlobSidecarForBlobTx
			}
			if err := sidecar.ValidateWithBlobHashes(tx.BlobHashes()); err != nil {
				return istanbul.ErrInvalidBlobTxWithSidecar
			}
		}
	}
	if header.BlobGasUsed != nil {
		if want := *header.BlobGasUsed / params.BlobTxBlobGasPerBlob; uint64(blobs) != want {
			return errors.New("blob gas used mismatch")
		}
	} else if blobs > 0 {
		return errors.New("data blobs present in block body")
	}
	return nil
}

func (b *scenarioBackend) Commit(proposal bft.Proposal, seals [][]byte) error {
	b.net.mu.Lock()
	defer b.net.mu.Unlock()
	b.commitCalls++
	if len(b.commitErrors) != 0 {
		err := b.commitErrors[0]
		b.commitErrors = b.commitErrors[1:]
		return err // Core error handling only; no real backend/import is tested.
	}
	block, ok := proposal.(*types.Block)
	if !ok {
		return istanbul.ErrInvalidProposal
	}
	header := block.Header()
	b.sealer.WriteRound(header, b.view.Round.Int64())
	if err := b.sealer.WriteCommittedSeals(header, seals); err != nil {
		return err
	}
	block = block.WithSeal(header)
	if err := b.net.validators[b.id].checkCommit(block); err != nil {
		b.net.failures = append(b.net.failures, err)
		return err
	}
	b.committed = append(b.committed, block)
	b.head = block
	b.blocks[block.NumberU64()] = block
	b.net.pendingDeliveries = append(b.net.pendingDeliveries, scenarioEvent{b.id, b.id, istanbul.ChainHeadEvent{}})
	return nil
}

type scenarioEvent struct {
	from, to int
	data     interface{}
}

type messageRule struct {
	code       uint64
	from       *validator
	recipients []*validator
	drop       bool
	hold       bool
	proposal   *types.Block
	held       []scenarioEvent
	applied    bool
	round      *uint64
}

func newMessageRule(net *scenarioNet, code uint64, from *validator, recipients []*validator, rounds ...uint64) *messageRule {
	net.t.Helper()
	net.requireNode(from)
	require.Contains(net.t, []uint64{bft.MsgPreprepare, bft.MsgPrepare, bft.MsgCommit, bft.MsgRoundChange}, code)
	require.NotEmpty(net.t, recipients)
	require.LessOrEqual(net.t, len(rounds), 1)
	var round *uint64
	if len(rounds) == 1 {
		round = &rounds[0]
	}
	for _, recipient := range recipients {
		net.requireNode(recipient)
		for _, existing := range net.rules {
			if round != nil && existing.round != nil && *round != *existing.round {
				continue
			}
			if existing.code == code && existing.from == from &&
				(existing.drop || existing.hold || existing.proposal != nil) {
				require.NotContains(net.t, existing.recipients, recipient, "overlapping message rules")
			}
		}
	}
	rule := &messageRule{code: code, from: from, recipients: append([]*validator(nil), recipients...), round: round}
	net.rules = append(net.rules, rule)
	return rule
}

// Snapshot is used only around targeted rejected/ignored inputs, not every event.
type scenarioSnapshot struct {
	Pending              common.Hash
	FutureRequests       int
	Sequence, Round      uint64
	State                State
	Waiting              bool
	Lock, Proposal, Head common.Hash
	Prepares, Commits    map[common.Address]common.Hash
	RoundChanges         map[uint64]map[common.Address]common.Hash
	Sent                 int
	Backlogs             map[common.Address]int
}

func (node *validator) snapshot() scenarioSnapshot {
	net := node.backend.net
	c := node.core
	messages := func(set *messageSet) map[common.Address]common.Hash {
		result := make(map[common.Address]common.Hash)
		for _, msg := range set.Values() {
			payload, err := msg.Payload()
			require.NoError(net.t, err)
			result[msg.Address] = crypto.Keccak256Hash(payload)
		}
		return result
	}
	s := scenarioSnapshot{
		Sequence: c.current.Sequence().Uint64(), Round: c.current.Round().Uint64(),
		State: c.state, Waiting: c.waitingForRoundChange, Lock: c.current.GetLockedHash(), Head: node.backend.head.Hash(),
		Prepares: messages(c.current.Prepares), Commits: messages(c.current.Commits),
		RoundChanges: make(map[uint64]map[common.Address]common.Hash),
	}
	s.Backlogs = make(map[common.Address]int)
	for from, q := range c.backlogs {
		s.Backlogs[from] = q.Size()
	}
	for from := range c.backlogPreprepares {
		s.Backlogs[from]++
	}
	s.FutureRequests = c.pendingRequests.Size()
	if c.current.pendingRequest != nil {
		s.Pending = c.current.pendingRequest.Proposal.Hash()
	}
	if proposal := c.current.Proposal(); proposal != nil {
		s.Proposal = proposal.Hash()
	}
	for round, set := range c.roundChangeSet.roundChanges {
		s.RoundChanges[round] = messages(set)
	}
	for _, sent := range net.sent {
		if sent.from == node.id {
			s.Sent++
		}
	}
	return s
}

func (node *validator) checkCommit(block *types.Block) error {
	net, id := node.backend.net, node.id
	backend := node.backend
	height := block.NumberU64()
	if want := backend.head.NumberU64() + 1; height != want {
		return fmt.Errorf("node %d: non-consecutive commit: have height %d, want %d", id, height, want)
	}
	if want := backend.head.Hash(); block.ParentHash() != want {
		return fmt.Errorf("node %d: wrong committed parent: have %s, want %s", id, block.ParentHash(), want)
	}
	if want := int(height); len(backend.committed)+1 != want {
		return fmt.Errorf("node %d: commit count %d, want %d", id, len(backend.committed)+1, want)
	}
	if err := node.validateCommittedSeals(block); err != nil {
		return err
	}
	author, err := backend.sealer.Author(block.Header())
	if err != nil {
		return fmt.Errorf("node %d: recover committed block author: %w", id, err)
	}
	if !slices.Contains(net.validatorAddresses(), author) {
		return fmt.Errorf("node %d: unauthorized block author %s", id, author)
	}

	// All equivocation scenarios stay within f faulty signers, preserving the agreement assumption.
	for _, other := range net.validators {
		if other.id == id {
			continue
		}
		if committed := other.backend.blocks[height]; committed != nil {
			if committed.Hash() != block.Hash() {
				return fmt.Errorf("nodes %d and %d committed different blocks at height %d: %s != %s",
					other.id, id, height, committed.Hash(), block.Hash())
			}
		}
	}
	return nil
}

type scenarioSealFormat struct {
	round      uint64
	roundBound bool
}

func (node *validator) assertCommittedSeals(block *types.Block, expected ...scenarioSealFormat) {
	net := node.backend.net
	net.t.Helper()
	require.NoError(net.t, node.validateCommittedSeals(block, expected...))
}

func (node *validator) validateCommittedSeals(block *types.Block, expected ...scenarioSealFormat) error {
	net, id, backend := node.backend.net, node.id, node.backend
	var (
		committers []common.Address
		err        error
	)
	if len(expected) != 0 {
		if len(expected) != 1 {
			return fmt.Errorf("node %d: got %d expected seal formats, want at most one", id, len(expected))
		}
		format := expected[0]
		if format.round > 255 {
			return fmt.Errorf("node %d: seal round %d overflows byte", id, format.round)
		}
		preimage := append(block.Hash().Bytes(), byte(bft.MsgCommit))
		if format.roundBound {
			preimage = append(preimage, byte(format.round))
		}
		_, seals, err := backend.sealer.RawSeals(block.Header())
		if err != nil {
			return fmt.Errorf("node %d: read committed seals: %w", id, err)
		}
		for _, seal := range seals {
			signer, err := istanbul.GetSignatureAddress(preimage, seal)
			if err != nil {
				return fmt.Errorf("node %d: recover committed seal: %w", id, err)
			}
			committers = append(committers, signer)
		}
	} else if backend.IsPermissionlessAt(block.NumberU64()) {
		committers, err = backend.sealer.CommittersWithRound(block.Header())
	} else {
		committers, err = backend.sealer.Committers(block.Header())
	}
	if err != nil {
		return fmt.Errorf("node %d: recover committers: %w", id, err)
	}

	committee := net.validatorAddresses()[:net.config.committeeSize]
	unique := make(map[common.Address]struct{}, len(committers))
	for _, committer := range committers {
		if _, ok := unique[committer]; ok {
			return fmt.Errorf("node %d: duplicate committed seal from %s", id, committer)
		}
		if !slices.Contains(committee, committer) {
			return fmt.Errorf("node %d: sealer %s outside the committee", id, committer)
		}
		unique[committer] = struct{}{}
	}
	// Restate the tiny-committee exception independently of calcQuorumSize.
	quorum := net.config.committeeSize
	if quorum >= 4 {
		quorum = (2*quorum + 2) / 3
	}
	if len(unique) < quorum {
		return fmt.Errorf("node %d: %d committed seals below quorum %d", id, len(unique), quorum)
	}
	return nil
}

// follow models a head notification only; this scenario ends before further local commits.
func (n *validator) follow(source *validator) {
	s := n.backend.net
	s.t.Helper()
	commits := len(n.backend.committed)
	for height := n.backend.head.NumberU64() + 1; height <= source.backend.head.NumberU64(); height++ {
		block := source.backend.blocks[height]
		require.NotNil(s.t, block)
		require.Equal(s.t, n.backend.head.Hash(), block.ParentHash())
		n.backend.blocks[height] = block
		n.backend.head = block
	}
	n.post(istanbul.ChainHeadEvent{})
	s.drain()
	require.Len(s.t, n.backend.committed, commits, "head notification must not invent local commits")
	require.Equal(s.t, source.head().Hash(), n.head().Hash())
}

type consensusInvalidKind string

const (
	consensusMalformedRLP     consensusInvalidKind = "MalformedRLP"
	consensusUnknownCode      consensusInvalidKind = "UnknownCode"
	consensusInvalidSignature consensusInvalidKind = "InvalidSignature"
	consensusMissingView      consensusInvalidKind = "MissingView"
	consensusOverflowRound    consensusInvalidKind = "OverflowRound"
)

type consensusSealMutation string

const (
	consensusMalformedSeal consensusSealMutation = "Malformed"
	consensusWrongSigner   consensusSealMutation = "WrongSigner"
	consensusWrongDigest   consensusSealMutation = "WrongDigest"
	consensusOtherRound    consensusSealMutation = "OtherRound"
)

func (sender *validator) invalidMessage(invalid consensusInvalidKind) istanbul.MessageEvent {
	s := sender.backend.net
	ev := sender.message(bft.MsgPrepare, sender.proposal(1), 0)
	if invalid == consensusMalformedRLP {
		ev.Payload = []byte{0xff}
		return ev
	}
	// Mutate the envelope in the wire format of height 1 directly: a mutation
	// may remove the view that would otherwise select the codec.
	if !sender.backend.IsPermissionlessAt(1) {
		ev.Payload = resignLegacy(s.t, ev.Payload, sender.backend.key, invalid == consensusInvalidSignature, func(msg *bft.PrePermissionlessMessage) {
			var subject bft.Subject
			require.NoError(s.t, rlp.DecodeBytes(msg.Msg, &subject))
			switch invalid {
			case consensusUnknownCode:
				msg.Code = 99
			case consensusInvalidSignature:
				msg.Signature = []byte{1}
			case consensusMissingView:
				subject.View = nil
			case consensusOverflowRound:
				msg.Code = bft.MsgRoundChange
				subject.View.Round = new(big.Int).Lsh(big.NewInt(1), 64)
				subject.Digest = common.Hash{}
			default:
				s.t.Fatalf("unknown invalid input %s", invalid)
			}
			var err error
			msg.Msg, err = bft.Encode(&subject)
			require.NoError(s.t, err)
		})
		return ev
	}
	msg, err := sender.backend.decodeMessage(ev.Payload, nil)
	require.NoError(s.t, err)
	switch invalid {
	case consensusUnknownCode:
		msg.Code = 99
	case consensusInvalidSignature:
		msg.Signature = []byte{1}
	case consensusMissingView, consensusOverflowRound:
		var prepare bft.Prepare
		require.NoError(s.t, msg.Decode(&prepare))
		var payload any = &prepare
		if invalid == consensusMissingView {
			prepare.View = nil
		} else {
			msg.Code = bft.MsgRoundChange
			prepare.View.Round = new(big.Int).Lsh(big.NewInt(1), 64)
			payload = &bft.RoundChange{View: prepare.View}
		}
		msg.Msg, err = bft.Encode(payload)
		require.NoError(s.t, err)
	default:
		s.t.Fatalf("unknown invalid input %s", invalid)
	}
	if invalid != consensusInvalidSignature {
		unsigned, err := msg.PayloadNoSig()
		require.NoError(s.t, err)
		msg.Signature, err = crypto.Sign(crypto.Keccak256(unsigned), sender.backend.key)
		require.NoError(s.t, err)
	}
	ev.Payload, err = msg.Payload()
	require.NoError(s.t, err)
	return ev
}

// resignLegacy applies mutate to a legacy wire message and, unless keepSignature
// is set, signs the result over the legacy preimage.
func resignLegacy(t *testing.T, payload []byte, key *ecdsa.PrivateKey, keepSignature bool, mutate func(*bft.PrePermissionlessMessage)) []byte {
	t.Helper()
	var legacy bft.PrePermissionlessMessage
	require.NoError(t, rlp.DecodeBytes(payload, &legacy))
	mutate(&legacy)
	if !keepSignature {
		legacy.Signature = nil
		unsigned, err := bft.Encode(&legacy)
		require.NoError(t, err)
		legacy.Signature, err = crypto.Sign(crypto.Keccak256(unsigned), key)
		require.NoError(t, err)
	}
	resigned, err := bft.Encode(&legacy)
	require.NoError(t, err)
	return resigned
}

func (sender *validator) corruptCommit(corruption consensusSealMutation) istanbul.MessageEvent {
	s := sender.backend.net
	ev := s.message(sender, bft.MsgCommit, 1, 0)
	msg, err := sender.backend.decodeMessage(ev.Payload, nil)
	require.NoError(s.t, err)
	var commit bft.Commit
	require.NoError(s.t, msg.Decode(&commit))
	if corruption == consensusMalformedSeal {
		commit.CommittedSeal = []byte{1}
	} else {
		digest, round, key := commit.Digest, byte(commit.View.Round.Uint64()), sender.backend.key
		switch corruption {
		case consensusWrongSigner:
			key = s.validators[2].backend.key
		case consensusWrongDigest:
			digest = common.HexToHash("0xbad")
		case consensusOtherRound:
			round++
		default:
			s.t.Fatalf("unknown corruption %s", corruption)
		}
		preimage := append(digest.Bytes(), byte(bft.MsgCommit))
		if sender.backend.IsPermissionlessAt(commit.View.Sequence.Uint64()) || corruption == consensusOtherRound {
			preimage = append(preimage, round)
		}
		commit.CommittedSeal, err = crypto.Sign(crypto.Keccak256(preimage), key)
		require.NoError(s.t, err)
	}
	msg.Msg, err = bft.Encode(&commit)
	require.NoError(s.t, err)
	unsigned, err := msg.PayloadNoSigForFork(sender.backend.IsPermissionlessAt)
	require.NoError(s.t, err)
	msg.Signature, err = crypto.Sign(crypto.Keccak256(unsigned), sender.backend.key)
	require.NoError(s.t, err)
	ev.Payload, err = msg.PayloadForFork(sender.backend.IsPermissionlessAt)
	require.NoError(s.t, err)
	return ev
}

func (sender *validator) invalidParentProposal() *types.Block {
	s := sender.backend.net
	bad := sender.proposal(1).Header()
	bad.ParentHash = common.HexToHash("0xdeadbeef")
	authorSeal, err := sender.backend.sealer.MakeAuthorSeal(bad)
	require.NoError(s.t, err)
	require.NoError(s.t, sender.backend.sealer.WriteAuthorSeal(bad, authorSeal))
	proposal := sender.proposal(1).WithSeal(bad)
	return proposal
}

// assertCommitReply checks construction only; current core rejects the old COMMIT on self-delivery before forwarding.
func (n *validator) assertCommitReply(height, round uint64, roundBound bool, action func()) {
	s := n.backend.net
	s.t.Helper()
	start := len(s.sent)
	action()
	require.Len(s.t, s.sent[start:], 1, "old proposal must produce one reply")
	sent := s.sent[start]
	require.Equal(s.t, n.id, sent.from)
	msg, err := n.backend.decodeMessage(sent.data.(istanbul.MessageEvent).Payload, istanbul.GetSignatureAddress)
	require.NoError(s.t, err)
	require.Equal(s.t, bft.MsgCommit, msg.Code)
	var commit bft.Commit
	require.NoError(s.t, msg.Decode(&commit))
	require.Equal(s.t, height, commit.View.Sequence.Uint64())
	require.Equal(s.t, round, commit.View.Round.Uint64())
	require.Equal(s.t, n.backend.blocks[height].Hash(), commit.Digest)
	preimage := append(commit.Digest.Bytes(), byte(bft.MsgCommit))
	if roundBound {
		preimage = append(preimage, byte(round))
	}
	signer, err := istanbul.GetSignatureAddress(preimage, commit.CommittedSeal)
	require.NoError(s.t, err)
	require.Equal(s.t, n.backend.Address(), signer)
}
