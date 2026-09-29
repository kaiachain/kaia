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
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #962, #985: the COMMIT seal signs the digest of the message's own subject,
// round-bound after the permissionless fork.
func TestDevParity_CommitSealSignsSubjectWithRound(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  *params.ChainConfig
	}{
		{"legacy", legacyConfig()},
		{"permissionless", permissionlessConfig()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, key := newTestBackend(t, &fakeChain{cfg: tc.cfg})
			self := crypto.PubkeyToAddress(key.PublicKey)
			current := newEmptyBlock(10, common.Hash{1})
			old := newEmptyBlock(9, common.Hash{2})
			m := newTestMachine(b, 10, 3, current, []common.Address{self})

			msgs, stop := captureMessages(b)
			defer stop()

			// A COMMIT for the old block while a different block is the current
			// proposal: the seal must attest the subject's digest, not m.preprepare.
			sub := &bft.Subject{View: &bft.View{Sequence: big.NewInt(9), Round: big.NewInt(2)}, Digest: old.Hash(), PrevHash: old.ParentHash()}
			m.broadcastMsg(&bft.Message{Hash: sub.PrevHash, Code: bft.MsgCommit, Msg: mustEncode(t, sub)})
			msg := recvMessage(t, msgs)
			require.Equal(t, uint64(bft.MsgCommit), msg.Code)
			require.Len(t, msg.CommittedSeal, crypto.SignatureLength)

			want := istanbul.PrepareCommittedSeal(old.Hash())
			if tc.cfg.IsPermissionlessForkEnabled(big.NewInt(9)) {
				want = istanbul.PrepareCommittedSealWithRound(old.Hash(), 2)
			}
			assert.True(t, recovers(want, msg.CommittedSeal, self), "seal must sign the subject digest (and round post-fork)")
		})
	}
}

// dev #924, #985: a COMMIT whose seal is not the sender's signature over the
// committed-seal preimage is rejected and never retained.
func TestDevParity_HandleCommitVerifiesSeal(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	otherKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	block := newEmptyBlock(10, common.Hash{1})
	m := newTestMachine(b, 10, 1, block, []common.Address{self})
	sub := m.subject()

	bad := &bft.Message{Code: bft.MsgCommit, Msg: mustEncode(t, sub), Address: self,
		CommittedSeal: sealOver(t, otherKey, istanbul.PrepareCommittedSealWithRound(sub.Digest, 1))}
	assert.ErrorIs(t, m.handleCommit(bad, self), errInvalidCommittedSeal)
	assert.Equal(t, 0, m.commits.Size())

	legacySeal := &bft.Message{Code: bft.MsgCommit, Msg: mustEncode(t, sub), Address: self,
		CommittedSeal: sealOver(t, key, istanbul.PrepareCommittedSeal(sub.Digest))}
	assert.ErrorIs(t, m.handleCommit(legacySeal, self), errInvalidCommittedSeal, "post-fork seals must bind the round")

	good := &bft.Message{Code: bft.MsgCommit, Msg: mustEncode(t, sub), Address: self,
		CommittedSeal: sealOver(t, key, istanbul.PrepareCommittedSealWithRound(sub.Digest, 1))}
	assert.NoError(t, m.handleCommit(good, self))
	assert.Equal(t, 1, m.commits.Size())
}

// dev #1054: post-permissionless, a PRE-PREPARE for an already-committed block
// is answered with a COMMIT only for the round this node stored.
func TestDevParity_OldBlockCommitRequiresStoredRound(t *testing.T) {
	for _, tc := range []struct {
		name           string
		cfg            *params.ChainConfig
		requestedRound int64
		expectCommit   bool
	}{
		{"permissionless answers the stored round", permissionlessConfig(), 1, true},
		{"permissionless refuses another round", permissionlessConfig(), 5, false},
		{"legacy answers any round", legacyConfig(), 5, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			chain := &fakeChain{cfg: tc.cfg, headers: map[common.Hash]*types.Header{}}
			b, _ := newTestBackend(t, chain)
			proposerKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)

			header := &types.Header{Number: big.NewInt(10), ParentHash: common.Hash{9}, Time: big.NewInt(0)}
			require.NoError(t, b.sealer.WriteValidators(header, []common.Address{proposer}))
			b.sealer.WriteRound(header, 1)
			stored := types.NewBlockWithHeader(header)
			chain.headers[stored.Hash()] = stored.Header()

			mValset := valset_mock.NewMockValsetModule(ctrl)
			mValset.EXPECT().GetProposer(uint64(10), uint64(tc.requestedRound)).Return(proposer, nil)
			b.valsetModule = mValset

			m := newTestMachine(b, 11, 0, nil, []common.Address{proposer})
			m.state = stateAcceptRequest
			msgs, stop := captureMessages(b)
			defer stop()

			pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(tc.requestedRound)}, Proposal: stored}
			err = m.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp)}, proposer)
			if tc.expectCommit {
				assert.NoError(t, err)
				commit := recvMessage(t, msgs)
				assert.Equal(t, uint64(bft.MsgCommit), commit.Code)
			} else {
				assert.ErrorIs(t, err, errOldMessage)
				expectNoMessage(t, msgs)
			}
		})
	}
}

// The seal preimage kaiabft builds must be byte-identical to istanbul's, which
// verifySeals and istanbul validators use to recover committers.
func TestDevParity_CommittedSealPreimageMatchesIstanbul(t *testing.T) {
	digest := common.HexToHash("0x1234")
	view := &bft.View{Sequence: big.NewInt(7), Round: big.NewInt(3)}

	legacy, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	assert.Equal(t, istanbul.PrepareCommittedSeal(digest), legacy.committedSealPreimage(view, digest))

	permissionless, _ := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	assert.Equal(t, istanbul.PrepareCommittedSealWithRound(digest, 3), permissionless.committedSealPreimage(view, digest))
}

// A node's own COMMIT round-trips through handleMsg (envelope shape, signature
// and seal checks), and its seal is byte-identical to istanbul's for the same
// key and view.
func TestDevParity_SelfCommitRoundTrip(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	block := newEmptyBlock(10, common.Hash{1})
	m := newTestMachine(b, 10, 2, block, []common.Address{self})
	msgs, stop := captureMessages(b)
	defer stop()

	m.sendCommit()
	commit := recvMessage(t, msgs)
	want, err := istanbul.NewSealerImpl(key).MakeCommittedSealFromHashWithRound(block.Hash(), 2)
	require.NoError(t, err)
	assert.Equal(t, want, commit.CommittedSeal)

	payload, err := commit.Payload()
	require.NoError(t, err)
	require.NoError(t, m.handleMsg(payload))
	assert.Equal(t, 1, m.commits.Size())
}

// A hash-locked block re-proposed in a later round is committed with seals
// bound to the new round; the round byte stays outside the block hash.
func TestDevParity_HashLockedReproposalCommitBindsNewRound(t *testing.T) {
	sealer := istanbul.NewSealerImpl(nil)
	prev := types.HeaderHashFn
	types.SetHeaderHashFn(sealer.HeaderHash)
	defer types.SetHeaderHashFn(prev)

	b, key := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	proposerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)

	header := newEmptyBlock(10, common.Hash{1}).Header()
	require.NoError(t, b.sealer.WriteValidators(header, []common.Address{self, proposer}))
	b.sealer.WriteRound(header, 2)
	locked := types.NewBlockWithHeader(header)

	m := newTestMachine(b, 10, 3, locked, []common.Address{self, proposer})
	m.preprepare.View = &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(2)}
	m.lockedHash = locked.Hash()
	m.state = stateAcceptRequest
	m.proposer = proposer

	reHeader := locked.Header()
	b.sealer.WriteRound(reHeader, 3)
	reproposed := locked.WithSeal(reHeader)
	require.Equal(t, locked.Hash(), reproposed.Hash())

	msgs, stop := captureMessages(b)
	defer stop()
	pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(3)}, Proposal: reproposed}
	require.NoError(t, m.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp)}, proposer))
	require.Equal(t, statePrepared, m.state)

	commit := recvMessage(t, msgs)
	require.Equal(t, uint64(bft.MsgCommit), commit.Code)
	assert.True(t, recovers(istanbul.PrepareCommittedSealWithRound(locked.Hash(), 3), commit.CommittedSeal, self))
	payload, err := commit.Payload()
	require.NoError(t, err)
	require.NoError(t, m.handleMsg(payload))
	assert.Equal(t, 1, m.commits.Size())
}

// Before the permissionless fork the COMMIT seal is the legacy one: a legacy
// seal is accepted and a round-bound seal is not.
func TestDevParity_HandleCommitUsesLegacySealBeforeFork(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	m := newTestMachine(b, 10, 1, newEmptyBlock(10, common.Hash{1}), []common.Address{self})
	sub := m.subject()

	roundBound := &bft.Message{Code: bft.MsgCommit, Msg: mustEncode(t, sub), Address: self,
		CommittedSeal: sealOver(t, key, istanbul.PrepareCommittedSealWithRound(sub.Digest, 1))}
	assert.ErrorIs(t, m.handleCommit(roundBound, self), errInvalidCommittedSeal)

	legacy := &bft.Message{Code: bft.MsgCommit, Msg: mustEncode(t, sub), Address: self,
		CommittedSeal: sealOver(t, key, istanbul.PrepareCommittedSeal(sub.Digest))}
	assert.NoError(t, m.handleCommit(legacy, self))
	assert.Equal(t, 1, m.commits.Size())
}
