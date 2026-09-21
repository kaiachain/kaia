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
	"math/big"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
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
	t.Cleanup(func() {
		for _, n := range s.validators {
			require.NoError(t, n.core.Stop())
			n.backend.mux.Stop()
		}
		fork.ClearHardForkBlockNumberConfig()
		types.SetHeaderHashFn(previousHash)
	})
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
	return s
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
		var msg bft.Message
		require.NoError(s.t, msg.FromPayload(ev.Payload, nil))
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
		chainConfig: net.config.chainConfig,
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
	n.backend.net.t.Helper()
	require.NoError(n.backend.net.t, n.backend.mux.Post(ev))
	synctest.Wait()
}

// receive checks the message handler's return value while the event loop is idle.
func (n *validator) receive(ev istanbul.MessageEvent, want error) {
	t := n.backend.net.t
	t.Helper()
	synctest.Wait()
	err := n.core.handleMsg(ev.Payload)
	synctest.Wait()
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
	view := &bft.View{Sequence: proposal.Number(), Round: new(big.Int).SetUint64(round)}
	var subject interface{}
	switch code {
	case bft.MsgPreprepare:
		subject = &bft.Preprepare{View: view, Proposal: proposal}
	case bft.MsgPrepare, bft.MsgCommit:
		subject = &bft.Subject{View: view, Digest: proposal.Hash(), PrevHash: proposal.ParentHash()}
	case bft.MsgRoundChange:
		subject = &bft.Subject{View: view, PrevHash: proposal.ParentHash()}
	default:
		net.t.Fatalf("unknown consensus message code %d", code)
	}
	encoded, err := bft.Encode(subject)
	require.NoError(net.t, err)
	payload, err := node.core.finalizeMessage(&bft.Message{Hash: proposal.ParentHash(), Code: code, Msg: encoded})
	require.NoError(net.t, err)
	return istanbul.MessageEvent{Hash: proposal.ParentHash(), Payload: payload}
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
	commitErrors []error
	commitCalls  int
}

var _ istanbul.Backend = (*scenarioBackend)(nil)

func (b *scenarioBackend) Address() common.Address          { return crypto.PubkeyToAddress(b.key.PublicKey) }
func (b *scenarioBackend) Sealer() *istanbul.IstanbulSealer { return b.sealer }
func (b *scenarioBackend) EventMux() *event.TypeMux         { return b.mux }
func (b *scenarioBackend) NodeType() common.ConnType        { return common.CONSENSUSNODE }
func (b *scenarioBackend) HasBadProposal(common.Hash) bool  { return false }

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
	return b.fanout(hash, payload, true)
}

func (b *scenarioBackend) Gossip([]byte) error {
	b.net.t.Fatal("unexpected Gossip: core broadcasts must pass through self-delivery")
	return nil
}

// Only successful self-processing forwards to direct peers; peer relays and P2P caches are not modeled.
func (b *scenarioBackend) GossipSubPeer(hash common.Hash, payload []byte) {
	var msg bft.Message
	require.NoError(b.net.t, msg.FromPayload(payload, nil))
	if msg.Address == b.Address() {
		require.NoError(b.net.t, b.fanout(hash, payload, false))
	}
}

// fanout applies recipient rules to signed core outputs and queues them for the normal message handler.
func (b *scenarioBackend) fanout(hash common.Hash, payload []byte, self bool) error {
	b.net.mu.Lock()
	defer b.net.mu.Unlock()
	var msg bft.Message
	require.NoError(b.net.t, msg.FromPayload(payload, nil))
	if self {
		b.net.sent = append(b.net.sent, scenarioEvent{b.id, b.id, istanbul.MessageEvent{Hash: hash, Payload: slices.Clone(payload)}})
		if msg.Code == bft.MsgPreprepare {
			var pp bft.Preprepare
			require.NoError(b.net.t, msg.Decode(&pp))
			b.net.proposal = pp.Proposal.(*types.Block)
		}
	}
	if self && msg.Code == bft.MsgCommit {
		var subject bft.Subject
		require.NoError(b.net.t, msg.Decode(&subject))
		current := b.net.validators[b.id].core.current
		if subject.View.Sequence.Cmp(current.Sequence()) == 0 {
			require.True(b.net.t, current.IsHashLocked(), "node %d: honest COMMIT without a lock", b.id)
			require.Equal(b.net.t, current.GetLockedHash(), subject.Digest, "node %d: COMMIT conflicts with lock", b.id)
		}
	}
	view, err := msg.GetView()
	require.NoError(b.net.t, err)
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
					require.NoError(b.net.t, msg.Decode(&pp))
					require.Equal(b.net.t, pp.View.Sequence, rule.proposal.Number(), "modified proposal must keep the message's sequence")
					// An equivocating proposer changes only the block. Preserve a real
					// higher-round ROUND-CHANGE certificate so receivers can test that
					// it binds the proposal rather than merely being present.
					pp.Proposal = rule.proposal
					encoded, err := bft.Encode(pp)
					require.NoError(b.net.t, err)
					payload, err := rule.from.core.finalizeMessage(&bft.Message{Hash: rule.proposal.ParentHash(), Code: msg.Code, Msg: encoded})
					require.NoError(b.net.t, err)
					ev.data = istanbul.MessageEvent{Hash: rule.proposal.ParentHash(), Payload: payload}
				} else {
					var subject *bft.Subject
					require.NoError(b.net.t, msg.Decode(&subject))
					require.Equal(b.net.t, subject.View.Sequence, rule.proposal.Number(), "modified proposal must keep the message's sequence")
					ev.data = rule.from.message(msg.Code, rule.proposal, subject.View.Round.Uint64())
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
	author, err := b.sealer.Author(block.Header())
	if err != nil {
		return 0, err
	}
	if !slices.Contains(b.net.validatorAddresses(), author) {
		return 0, istanbul.ErrUnauthorizedAddress
	}
	return 0, nil
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
	b.committed = append(b.committed, block)
	b.net.validators[b.id].checkCommit()
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
	pending              common.Hash
	futureRequests       int
	sequence, round      uint64
	state                State
	waiting              bool
	lock, proposal, head common.Hash
	prepares, commits    map[common.Address]common.Hash
	roundChanges         map[uint64]map[common.Address]common.Hash
	sent                 int
	backlogs             map[common.Address]int
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
		sequence: c.current.Sequence().Uint64(), round: c.current.Round().Uint64(),
		state: c.state, waiting: c.waitingForRoundChange, lock: c.current.GetLockedHash(), head: node.backend.head.Hash(),
		prepares: messages(c.current.Prepares), commits: messages(c.current.Commits),
		roundChanges: make(map[uint64]map[common.Address]common.Hash),
	}
	s.backlogs = make(map[common.Address]int)
	for from, q := range c.backlogs {
		s.backlogs[from] = q.Size()
	}
	for from := range c.backlogPreprepares {
		s.backlogs[from]++
	}
	s.futureRequests = c.pendingRequests.Size()
	if c.current.pendingRequest != nil {
		s.pending = c.current.pendingRequest.Proposal.Hash()
	}
	if proposal := c.current.Proposal(); proposal != nil {
		s.proposal = proposal.Hash()
	}
	for round, set := range c.roundChangeSet.roundChanges {
		s.roundChanges[round] = messages(set)
	}
	for _, sent := range net.sent {
		if sent.from == node.id {
			s.sent++
		}
	}
	return s
}

func (node *validator) checkCommit() {
	net, id := node.backend.net, node.id
	net.t.Helper()
	backend := node.backend
	block := backend.committed[len(backend.committed)-1]
	height := block.NumberU64()
	require.Equal(net.t, backend.head.NumberU64()+1, height, "node %d: non-consecutive commit", id)
	require.Equal(net.t, backend.head.Hash(), block.ParentHash(), "node %d: wrong committed parent", id)
	require.Len(net.t, backend.committed, int(height), "node %d: exactly one commit per height", id)
	node.assertCommittedSeals(block)
	author, err := backend.sealer.Author(block.Header())
	require.NoError(net.t, err)
	require.Contains(net.t, net.validatorAddresses(), author, "unauthorized block author")

	// All equivocation scenarios stay within f faulty signers, preserving the agreement assumption.
	for _, other := range net.validators {
		if other.id == id {
			continue
		}
		if committed := other.backend.blocks[height]; committed != nil {
			require.Equal(net.t, committed.Hash(), block.Hash(),
				"nodes %d and %d committed different blocks at height %d", other.id, id, height)
		}
	}
}

type scenarioSealFormat struct {
	round      uint64
	roundBound bool
}

func (node *validator) assertCommittedSeals(block *types.Block, expected ...scenarioSealFormat) {
	net, id, backend := node.backend.net, node.id, node.backend
	net.t.Helper()
	var (
		committers []common.Address
		err        error
	)
	if len(expected) != 0 {
		require.Len(net.t, expected, 1)
		format := expected[0]
		require.LessOrEqual(net.t, format.round, uint64(255))
		preimage := append(block.Hash().Bytes(), byte(bft.MsgCommit))
		if format.roundBound {
			preimage = append(preimage, byte(format.round))
		}
		_, seals, err := backend.sealer.RawSeals(block.Header())
		require.NoError(net.t, err)
		for _, seal := range seals {
			signer, err := istanbul.GetSignatureAddress(preimage, seal)
			require.NoError(net.t, err)
			committers = append(committers, signer)
		}
	} else if backend.IsPermissionlessAt(block.NumberU64()) {
		committers, err = backend.sealer.CommittersWithRound(block.Header())
	} else {
		committers, err = backend.sealer.Committers(block.Header())
	}
	require.NoError(net.t, err)

	committee := net.validatorAddresses()[:net.config.committeeSize]
	unique := make(map[common.Address]struct{}, len(committers))
	for _, committer := range committers {
		require.NotContains(net.t, unique, committer, "node %d: duplicate committed seals", id)
		require.Contains(net.t, committee, committer, "node %d: sealer outside the committee", id)
		unique[committer] = struct{}{}
	}
	// Restate the tiny-committee exception independently of calcQuorumSize.
	quorum := net.config.committeeSize
	if quorum >= 4 {
		quorum = (2*quorum + 2) / 3
	}
	require.GreaterOrEqual(net.t, len(unique), quorum, "node %d: seals below quorum", id)
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
	var msg bft.Message
	require.NoError(s.t, msg.FromPayload(ev.Payload, nil))
	switch invalid {
	case consensusMalformedRLP:
		ev.Payload = []byte{0xff}
	case consensusUnknownCode:
		msg.Code = 99
	case consensusInvalidSignature:
		msg.Signature = []byte{1}
	case consensusMissingView, consensusOverflowRound:
		var subject bft.Subject
		require.NoError(s.t, msg.Decode(&subject))
		if invalid == consensusMissingView {
			subject.View = nil
		} else {
			msg.Code = bft.MsgRoundChange
			subject.View.Round = new(big.Int).Lsh(big.NewInt(1), 64)
		}
		var err error
		msg.Msg, err = bft.Encode(&subject)
		require.NoError(s.t, err)
	default:
		s.t.Fatalf("unknown invalid input %s", invalid)
	}
	if invalid != consensusMalformedRLP {
		if invalid != consensusInvalidSignature {
			unsigned, err := msg.PayloadNoSig()
			require.NoError(s.t, err)
			msg.Signature, err = crypto.Sign(crypto.Keccak256(unsigned), sender.backend.key)
			require.NoError(s.t, err)
		}
		var err error
		ev.Payload, err = msg.Payload()
		require.NoError(s.t, err)
	}
	return ev
}

func (sender *validator) corruptCommit(corruption consensusSealMutation) istanbul.MessageEvent {
	s := sender.backend.net
	ev := s.message(sender, bft.MsgCommit, 1, 0)
	var msg bft.Message
	require.NoError(s.t, msg.FromPayload(ev.Payload, nil))
	var subject bft.Subject
	require.NoError(s.t, msg.Decode(&subject))
	if corruption == consensusMalformedSeal {
		msg.CommittedSeal = []byte{1}
	} else {
		digest, round, key := subject.Digest, byte(subject.View.Round.Uint64()), sender.backend.key
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
		if sender.backend.IsPermissionlessAt(subject.View.Sequence.Uint64()) {
			preimage = append(preimage, round)
		}
		var err error
		msg.CommittedSeal, err = crypto.Sign(crypto.Keccak256(preimage), key)
		require.NoError(s.t, err)
	}
	unsigned, err := msg.PayloadNoSig()
	require.NoError(s.t, err)
	msg.Signature, err = crypto.Sign(crypto.Keccak256(unsigned), sender.backend.key)
	require.NoError(s.t, err)
	ev.Payload, err = msg.Payload()
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
	var msg bft.Message
	require.NoError(s.t, msg.FromPayload(sent.data.(istanbul.MessageEvent).Payload, istanbul.GetSignatureAddress))
	require.Equal(s.t, bft.MsgCommit, msg.Code)
	var subject bft.Subject
	require.NoError(s.t, msg.Decode(&subject))
	require.Equal(s.t, height, subject.View.Sequence.Uint64())
	require.Equal(s.t, round, subject.View.Round.Uint64())
	require.Equal(s.t, n.backend.blocks[height].Hash(), subject.Digest)
	preimage := append(subject.Digest.Bytes(), byte(bft.MsgCommit))
	if roundBound {
		preimage = append(preimage, byte(round))
	}
	signer, err := istanbul.GetSignatureAddress(preimage, msg.CommittedSeal)
	require.NoError(s.t, err)
	require.Equal(s.t, n.backend.Address(), signer)
}
