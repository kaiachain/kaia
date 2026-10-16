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
	"crypto/ecdsa"
	"math/big"
	"slices"
	"sort"
	"sync"
	"testing"
	"testing/synctest"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

// recoveryNet runs several kaiabft machines in one synctest bubble. Each
// machine is driven by the test goroutine instead of runLoop, but through the
// same entry points runLoop uses: handleMsg for wire payloads, handleCheckedMsg
// for backlog events, handleRequest, handleTimeout and handleChainHead. A
// node's own broadcast is self-delivered first and reaches peers only through
// gossipSubPeer after the node accepted it, as in production. Relays of peer
// messages are not modelled.
type recoveryNet struct {
	t          *testing.T
	cfg        *params.ChainConfig
	nodes      []*recoveryNode
	addrs      []common.Address
	mu         sync.Mutex
	pending    []netDelivery
	rules      []*routeRule
	seq        int
	sent       [][]byte // every self-delivered broadcast, before routing
	committees map[uint64][]common.Address
	sentTo     [][]byte // every payload sent to a peer, before routing
}

type recoveryNode struct {
	net        *recoveryNet
	id         int
	key        *ecdsa.PrivateKey
	b          *backend
	m          *machine
	chain      *recoveryChain
	sub        *event.TypeMuxSubscription
	committed  map[uint64]*types.Block
	executions int
}

type deliveryKind int

const (
	deliverSelf deliveryKind = iota // own broadcast, self-delivered
	deliverPeer                     // payload sent to a peer
	deliverBacklog
	deliverRequest
	deliverChainHead
)

type netDelivery struct {
	kind    deliveryKind
	from    int
	to      int
	seq     int
	hash    common.Hash
	payload []byte
	backlog backlogEvent
	request bft.Proposal
}

// routeRule drops, or replaces the payload of, messages of one code from one
// signer to one recipient, optionally only in one round.
type routeRule struct {
	code    uint64
	from    int
	to      int
	round   int64 // -1 for every round
	replace []byte
}

// recoveryChain is the per-node chain reader the machine reaches.
type recoveryChain struct {
	consensus.ChainReader
	node           *recoveryNode
	cfg            *params.ChainConfig
	headers        map[common.Hash]*types.Header
	current        *types.Block
	bad            map[common.Hash]bool
	badExec        map[common.Hash]bool
	executionError error
}

func (c *recoveryChain) Config() *params.ChainConfig                { return c.cfg }
func (c *recoveryChain) HasBadBlock(h common.Hash) bool             { return c.bad[h] }
func (c *recoveryChain) ValidateProposalHeader(*types.Header) error { return nil }
func (c *recoveryChain) GetHeader(h common.Hash, _ uint64) *types.Header {
	return c.headers[h]
}
func (c *recoveryChain) CurrentBlock() *types.Block   { return c.current }
func (c *recoveryChain) CurrentHeader() *types.Header { return c.current.Header() }
func (c *recoveryChain) ValidateProposalExecution(block *types.Block) error {
	c.node.executions++
	if c.executionError != nil {
		return c.executionError
	}
	if c.badExec[block.Hash()] {
		c.bad[block.Hash()] = true
		return errInvalidCommittedSeal // any error; the record is what matters
	}
	return nil
}

type netBroadcaster struct{ node *recoveryNode }

func (nb *netBroadcaster) Enqueue(_ string, block *types.Block) { nb.node.onCommit(block) }
func (nb *netBroadcaster) FindPeers(map[common.Address]bool) map[common.Address]consensus.Peer {
	return nil
}
func (nb *netBroadcaster) FindCNPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	peers := make(map[common.Address]consensus.Peer)
	for _, n := range nb.node.net.nodes {
		if targets[n.b.address] {
			peers[n.b.address] = &netPeer{from: nb.node, to: n}
		}
	}
	return peers
}
func (nb *netBroadcaster) GetCNPeers() map[common.Address]consensus.Peer            { return nil }
func (nb *netBroadcaster) GetENPeers() map[common.Address]consensus.Peer            { return nil }
func (nb *netBroadcaster) RegisterValidator(common.ConnType, p2p.PeerTypeValidator) {}

type netPeer struct{ from, to *recoveryNode }

func (p *netPeer) Send(_ uint64, data any) error {
	msg := data.(*bft.ConsensusMsg)
	p.from.net.enqueue(netDelivery{kind: deliverPeer, from: p.from.id, to: p.to.id, hash: msg.PrevHash, payload: slices.Clone(msg.Payload)})
	return nil
}
func (p *netPeer) RegisterConsensusMsgCode(uint64) error { return nil }

// newRecoveryNet must run inside synctest.Test.
func newRecoveryNet(t *testing.T, count int, cfg *params.ChainConfig) *recoveryNet {
	s := &recoveryNet{t: t, cfg: cfg, committees: make(map[uint64][]common.Address)}
	keys := make([]*ecdsa.PrivateKey, count)
	for i := range keys {
		var err error
		keys[i], err = crypto.ToECDSA(common.LeftPadBytes(big.NewInt(int64(i+1)).Bytes(), 32))
		require.NoError(t, err)
		s.addrs = append(s.addrs, crypto.PubkeyToAddress(keys[i].PublicKey))
	}
	sealer := istanbul.NewSealerImpl(keys[0])
	previous := types.HeaderHashFn
	types.SetHeaderHashFn(sealer.HeaderHash)
	header := &types.Header{Number: new(big.Int), Time: big.NewInt(1), BlockScore: big.NewInt(1)}
	require.NoError(t, sealer.WriteValidators(header, s.addrs))
	genesis := types.NewBlockWithHeader(header)

	for i, key := range keys {
		n := &recoveryNode{net: s, id: i, key: key, committed: make(map[uint64]*types.Block)}
		n.chain = &recoveryChain{node: n, cfg: cfg, headers: map[common.Hash]*types.Header{genesis.Hash(): genesis.Header()},
			current: genesis, bad: make(map[common.Hash]bool), badExec: make(map[common.Hash]bool)}
		n.b = New(&Opts{Timeout: 10000, PrivateKey: key, NodeType: common.CONSENSUSNODE, Sealer: istanbul.NewSealerImpl(key)}).(*backend)
		n.b.chain = n.chain
		n.b.broadcaster = &netBroadcaster{node: n}
		n.b.commitCh = nil // as Start leaves it: no local seal is waiting

		ctrl := gomock.NewController(t)
		mValset := valset_mock.NewMockValsetModule(ctrl)
		addrs := s.addrs
		mValset.EXPECT().GetCouncil(gomock.Any()).Return(addrs, nil).AnyTimes()
		mValset.EXPECT().GetDemotedValidators(gomock.Any()).Return([]common.Address{}, nil).AnyTimes()
		mValset.EXPECT().GetCommittee(gomock.Any(), gomock.Any()).DoAndReturn(func(height, round uint64) ([]common.Address, error) {
			if committee, ok := s.committees[round]; ok {
				return committee, nil
			}
			return addrs, nil
		}).AnyTimes()
		mValset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(func(height, round uint64) (common.Address, error) {
			return addrs[(height-1+round)%uint64(len(addrs))], nil
		}).AnyTimes()
		n.b.RegisterKaiaxModules(gov_mock.NewMockGovModule(ctrl), mValset)

		n.m = newMachine(n.b)
		n.b.machine = n.m
		n.sub = n.b.eventMux.Subscribe(messageEvent{}, backlogEvent{}, requestEvent{})
		go n.collect()
		s.nodes = append(s.nodes, n)
	}
	for _, n := range s.nodes {
		n.m.startNewRound(common.Big0)
	}
	t.Cleanup(func() {
		for _, n := range s.nodes {
			n.m.stopTimer()
			n.sub.Unsubscribe()
		}
		types.SetHeaderHashFn(previous)
	})
	synctest.Wait()
	return s
}

// collect turns the events the machine posts to itself into deliveries.
func (n *recoveryNode) collect() {
	for ev := range n.sub.Chan() {
		switch e := ev.Data.(type) {
		case messageEvent:
			n.net.enqueue(netDelivery{kind: deliverSelf, from: n.id, to: n.id, hash: e.Hash, payload: slices.Clone(e.Payload)})
		case backlogEvent:
			n.net.enqueue(netDelivery{kind: deliverBacklog, from: n.id, to: n.id, backlog: e})
		case requestEvent:
			n.net.enqueue(netDelivery{kind: deliverRequest, from: n.id, to: n.id, request: e.Proposal})
		}
	}
}

func (s *recoveryNet) enqueue(d netDelivery) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	d.seq = s.seq
	s.pending = append(s.pending, d)
}

// take returns the queued deliveries in a deterministic order: deliveries
// posted from goroutines arrive in scheduling order, so they are sorted by
// content within a batch.
func (s *recoveryNet) take() []netDelivery {
	synctest.Wait()
	s.mu.Lock()
	batch := s.pending
	s.pending = nil
	s.mu.Unlock()
	sort.SliceStable(batch, func(i, j int) bool {
		a, b := batch[i], batch[j]
		if a.to != b.to {
			return a.to < b.to
		}
		if a.kind != b.kind {
			return a.kind < b.kind
		}
		if a.from != b.from {
			return a.from < b.from
		}
		return bytes.Compare(crypto.Keccak256(a.payload), crypto.Keccak256(b.payload)) < 0
	})
	return batch
}

// drain delivers until no node has anything left to say.
func (s *recoveryNet) drain() {
	for steps := 0; ; steps++ {
		require.Less(s.t, steps, 1000, "message loop")
		batch := s.take()
		if len(batch) == 0 {
			return
		}
		for _, d := range batch {
			s.deliver(d)
		}
	}
}

// advance hands each node not waiting for a round change its own proposal for
// height, then drains.
func (s *recoveryNet) advance(height uint64) {
	for _, n := range s.nodes {
		if n.m.sequence.Uint64() == height && !n.m.waitingForRoundChange && n.m.pendingRequest == nil {
			s.enqueue(netDelivery{kind: deliverRequest, from: n.id, to: n.id, request: n.proposal(1)})
		}
	}
	s.drain()
}

// timeout fires the round-change timeout of the given nodes, as their timer would.
func (s *recoveryNet) timeout(ids ...int) {
	for _, id := range ids {
		m := s.nodes[id].m
		m.handleTimeout(&bft.View{Sequence: new(big.Int).Set(m.sequence), Round: new(big.Int).Add(m.round, common.Big1)})
	}
	s.drain()
}

func (s *recoveryNet) deliver(d netDelivery) {
	n := s.nodes[d.to]
	switch d.kind {
	case deliverSelf, deliverPeer:
		if d.kind == deliverSelf {
			s.sent = append(s.sent, d.payload)
		} else {
			s.sentTo = append(s.sentTo, d.payload)
		}
		payload, drop := s.route(d)
		if drop {
			return
		}
		err := n.m.handleMsg(payload)
		if d.kind == deliverSelf && err == nil {
			n.b.gossipSubPeer(d.hash, payload)
		}
	case deliverBacklog:
		if n.m.qualified != nil && !n.m.qualified.Contains(d.backlog.src) {
			return
		}
		n.m.handleCheckedMsg(d.backlog.msg, d.backlog.src)
	case deliverRequest:
		if err := n.m.handleRequest(&bft.Request{Proposal: d.request}); err == errFutureMessage {
			n.m.storePendingRequest(&bft.Request{Proposal: d.request})
		}
	case deliverChainHead:
		n.m.handleChainHead()
	}
}

// route applies the drop and replace rules to a wire payload.
func (s *recoveryNet) route(d netDelivery) ([]byte, bool) {
	var msg bft.Message
	if err := msg.FromPayloadForFork(d.payload, s.isPermissionlessAt, nil); err != nil {
		return d.payload, false
	}
	view, err := msg.GetView()
	if err != nil {
		return d.payload, false
	}
	from := slices.Index(s.addrs, msg.Address)
	for _, r := range s.rules {
		if r.code == msg.Code && r.from == from && r.to == d.to && (r.round < 0 || view.Round.Int64() == r.round) {
			if r.replace != nil {
				return r.replace, false
			}
			return nil, true
		}
	}
	return d.payload, false
}

func (s *recoveryNet) isPermissionlessAt(num uint64) bool {
	return s.cfg.IsPermissionlessForkEnabled(new(big.Int).SetUint64(num))
}

// drop discards messages of code from every sender in from to every recipient
// in to, in the given round or, without one, in every round.
func (s *recoveryNet) drop(code uint64, from, to []int, round ...int64) {
	r := int64(-1)
	if len(round) == 1 {
		r = round[0]
	}
	for _, f := range from {
		for _, t := range to {
			s.rules = append(s.rules, &routeRule{code: code, from: f, to: t, round: r})
		}
	}
}

// equivocate makes the proposer from deliver replacement, a PRE-PREPARE it
// signs for the same view, to the recipients in to.
func (s *recoveryNet) equivocate(from int, to []int, view *bft.View, replacement *types.Block) {
	proposer := s.nodes[from]
	header := replacement.Header()
	proposer.b.sealer.WriteRound(header, view.Round.Int64())
	pp := &bft.Preprepare{View: view, Proposal: replacement.WithSeal(header)}
	payload := proposer.m.signPayload(&wireMessage{Hash: replacement.ParentHash(), Code: bft.MsgPreprepare, Msg: mustEncode(s.t, pp)})
	require.NotNil(s.t, payload)
	for _, t := range to {
		s.rules = append(s.rules, &routeRule{code: bft.MsgPreprepare, from: from, to: t, round: view.Round.Int64(), replace: payload})
	}
}

// sentRoundChanges decodes every ROUND CHANGE broadcast so far.
func (s *recoveryNet) sentRoundChanges() []*bft.Message {
	var out []*bft.Message
	for _, payload := range s.sent {
		var msg bft.Message
		require.NoError(s.t, msg.FromPayloadForFork(payload, s.isPermissionlessAt, nil))
		if msg.Code == bft.MsgRoundChange {
			out = append(out, &msg)
		}
	}
	return out
}

// sentPreprepare returns the PRE-PREPARE from signs for view, as sent to peers.
func (s *recoveryNet) sentPreprepare(from int, view *bft.View) *bft.Preprepare {
	for _, payload := range s.sentTo {
		var msg bft.Message
		require.NoError(s.t, msg.FromPayloadForFork(payload, s.isPermissionlessAt, nil))
		if msg.Code != bft.MsgPreprepare || msg.Address != s.addrs[from] {
			continue
		}
		var pp *bft.Preprepare
		require.NoError(s.t, msg.Decode(&pp))
		if pp.View.Cmp(view) == 0 {
			return pp
		}
	}
	return nil
}

func (s *recoveryNet) all() []int {
	ids := make([]int, len(s.nodes))
	for i := range ids {
		ids[i] = i
	}
	return ids
}

// proposal builds this node's proposal on its head; variant changes the hash.
func (n *recoveryNode) proposal(variant int64) *types.Block {
	t := n.net.t
	head := n.chain.current
	header := &types.Header{
		ParentHash: head.Hash(), Number: new(big.Int).Add(head.Number(), common.Big1),
		Time: new(big.Int).Add(head.Time(), big.NewInt(variant)), BlockScore: big.NewInt(1),
	}
	header.TxHash = types.DeriveTransactionsRoot(nil, header.Number)
	require.NoError(t, n.b.sealer.WriteValidators(header, n.net.addrs))
	seal, err := n.b.sealer.MakeAuthorSeal(header)
	require.NoError(t, err)
	require.NoError(t, n.b.sealer.WriteAuthorSeal(header, seal))
	return types.NewBlockWithHeader(header)
}

func (n *recoveryNode) onCommit(block *types.Block) {
	n.committed[block.NumberU64()] = block
	n.chain.current = block
	n.chain.headers[block.Hash()] = block.Header()
	n.net.enqueue(netDelivery{kind: deliverChainHead, from: n.id, to: n.id})
}

// requireCommitted checks the committed block of height, its round, and that
// its committed seals are a quorum of committee signatures bound to the round.
func (n *recoveryNode) requireCommitted(height uint64, round byte) *types.Block {
	t := n.net.t
	block := n.committed[height]
	require.NotNil(t, block, "node %d did not commit height %d", n.id, height)
	got, err := n.b.sealer.Round(block.Header())
	require.NoError(t, err)
	require.Equal(t, round, got, "node %d committed height %d at another round", n.id, height)
	committers, err := n.b.sealer.CommittersWithRound(block.Header())
	if !n.net.isPermissionlessAt(height) {
		committers, err = n.b.sealer.Committers(block.Header())
	}
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(committers), calcQuorumSize(len(n.net.addrs), uint64(len(n.net.addrs))))
	for _, c := range committers {
		require.Contains(t, n.net.addrs, c)
	}
	return block
}

func (n *recoveryNode) requireUncommitted(height uint64) {
	require.Nil(n.net.t, n.committed[height], "node %d committed height %d", n.id, height)
}

// signRoundChange signs a ROUND CHANGE from n for view with the given claim and
// Evidence, which a test may set inconsistently.
func (n *recoveryNode) signRoundChange(view *bft.View, prevHash common.Hash, claim *bft.PreparedClaim, evidence []byte) []byte {
	msg := &wireMessage{Hash: prevHash, Code: bft.MsgRoundChange, Msg: mustEncode(n.net.t, &bft.Subject{View: view, PrevHash: prevHash}), Prepared: claim, Evidence: evidence}
	payload := n.m.signPayload(msg)
	require.NotNil(n.net.t, payload)
	return payload
}

// signVote signs a PREPARE or COMMIT from n for block at view, in the shared
// message form a certificate embeds.
func (n *recoveryNode) signVote(code uint64, view *bft.View, block *types.Block) *bft.Message {
	sub := &bft.Subject{View: view, Digest: block.Hash(), PrevHash: block.ParentHash()}
	msg := &wireMessage{Hash: block.ParentHash(), Code: code, Msg: mustEncode(n.net.t, sub)}
	require.NotNil(n.net.t, n.m.signPayload(msg))
	out, err := msg.normalized()
	require.NoError(n.net.t, err)
	return out
}

// signPreprepare signs a PRE-PREPARE from n carrying pp.
func (n *recoveryNode) signPreprepare(pp *bft.Preprepare) []byte {
	payload := n.m.signPayload(&wireMessage{Hash: pp.Proposal.ParentHash(), Code: bft.MsgPreprepare, Msg: mustEncode(n.net.t, pp)})
	require.NotNil(n.net.t, payload)
	return payload
}

func viewAt(seq, round int64) *bft.View {
	return &bft.View{Sequence: big.NewInt(seq), Round: big.NewInt(round)}
}

func claimFor(t *testing.T, cert *bft.PreparedCertificate) (*bft.PreparedClaim, []byte) {
	evidence := mustEncode(t, cert)
	return &bft.PreparedClaim{Round: new(big.Int).Set(cert.View.Round), Digest: cert.Proposal.Hash(), EvidenceHash: crypto.Keccak256Hash(evidence)}, evidence
}
