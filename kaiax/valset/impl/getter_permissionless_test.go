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

package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/engine"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/storage/database"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testGetterBlock = uint64(10)

func TestPermissionlessPublicGetters(t *testing.T) {
	addr6 := numToAddr(6)
	addr7 := numToAddr(7)
	v := newCachedPermissionlessValset(t, NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive, Suspended: true},
		addr3: {State: ValPaused},
		addr4: {State: ValReady},
		addr5: {State: CandReady},
		addr6: {State: CandTesting},
		addr7: {State: Registered},
	})

	council, err := v.GetCouncil(testGetterBlock)
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Address{addr1, addr2, addr3}, council)

	cnPeers, err := v.GetCNPeers(testGetterBlock)
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Address{addr1, addr2, addr3, addr4, addr5, addr6}, cnPeers)

	demoted, err := v.GetDemotedValidators(testGetterBlock)
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Address{addr2, addr3}, demoted)

	committee, err := v.GetCommittee(testGetterBlock, 0)
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Address{addr1}, committee)

	candTesting, err := v.GetCandTesting(testGetterBlock)
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Address{addr6}, candTesting)
}

func TestGetNodesByState(t *testing.T) {
	addr6 := numToAddr(6)
	nodes := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValPaused},
		addr3: {State: ValInactive},
		addr4: {State: CandReady},
		addr5: {State: ValExiting},
		addr6: {State: ValActive, Suspended: true},
	}
	v := newCachedPermissionlessValset(t, nodes)

	t.Run("filter single state", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, []NodeState{ValActive})
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr6}, NodeMap(result).Addresses())
	})

	t.Run("filter multiple states", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, []NodeState{ValActive, ValPaused})
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr2, addr6}, NodeMap(result).Addresses())
	})

	t.Run("empty states returns all", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, []NodeState{})
		require.NoError(t, err)
		assert.ElementsMatch(t, nodes.Addresses(), NodeMap(result).Addresses())
	})

	t.Run("nil states returns all", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, nil)
		require.NoError(t, err)
		assert.ElementsMatch(t, nodes.Addresses(), NodeMap(result).Addresses())
	})

	t.Run("no match returns empty", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, []NodeState{ValReady})
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns defensive copy", func(t *testing.T) {
		result, err := v.GetNodesByState(testGetterBlock, nil)
		require.NoError(t, err)
		result[addr1].State = Registered

		again, err := v.GetNodesByState(testGetterBlock, []NodeState{ValActive})
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr6}, NodeMap(again).Addresses())
		assert.Equal(t, ValActive, again[addr1].State)
	})
}

func TestSuspendedFallback(t *testing.T) {
	t.Run("all ValActive suspended falls back to all ValActive", func(t *testing.T) {
		v := newCachedPermissionlessValset(t, NodeMap{
			addr1: {State: ValActive, Suspended: true},
			addr2: {State: ValActive, Suspended: true},
			addr3: {State: ValPaused, Suspended: true},
		})

		qualified, err := v.getQualifiedValidators(testGetterBlock)
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr2}, qualified.List())

		committee, err := v.GetCommittee(testGetterBlock, 0)
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr2}, committee)
	})

	t.Run("some ValActive suspended excludes only suspended validators", func(t *testing.T) {
		v := newCachedPermissionlessValset(t, NodeMap{
			addr1: {State: ValActive, Suspended: true},
			addr2: {State: ValActive},
			addr3: {State: ValActive},
		})

		qualified, err := v.getQualifiedValidators(testGetterBlock)
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr2, addr3}, qualified.List())
	})

	t.Run("demoted is council minus qualified", func(t *testing.T) {
		v := newCachedPermissionlessValset(t, NodeMap{
			addr1: {State: ValActive, Suspended: true},
			addr2: {State: ValActive},
			addr3: {State: ValPaused},
			addr4: {State: CandReady},
		})

		demoted, err := v.getDemoted(testGetterBlock)
		require.NoError(t, err)
		assert.ElementsMatch(t, []common.Address{addr1, addr3}, demoted)
	})
}

func TestGetNodesByStatePreForkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

	v := NewValsetModule()
	v.Chain = mockChain

	nodes, err := v.GetNodesByState(10, nil)
	require.Nil(t, nodes)
	require.ErrorIs(t, err, errPermissionlessDisabled)
}

func TestGetCNPeersFilterFallback(t *testing.T) {
	t.Run("pre-fork returns council", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10)).AnyTimes()

		db := database.NewMemDB()
		writeCouncil(db, 0, numsToAddrs(2, 1))
		writeValidatorVoteBlockNums(db, []uint64{0})
		writeLowestScannedVoteNum(db, 0)

		v := NewValsetModule()
		v.Chain = mockChain
		v.ChainKv = db

		cnPeers, err := v.GetCNPeers(11)
		require.NoError(t, err)
		assert.Equal(t, numsToAddrs(1, 2), cnPeers)
	})

	t.Run("pre-fork GetCouncil failure returns nil to disable filtering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10)).AnyTimes()

		v := NewValsetModule()
		v.Chain = mockChain
		v.ChainKv = database.NewMemDB()

		cnPeers, err := v.GetCNPeers(11)
		require.NoError(t, err)
		assert.Nil(t, cnPeers)
	})

	t.Run("post-fork read failure returns nil to disable filtering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()
		mockChain.EXPECT().GetHeaderByNumber(uint64(10)).Return(nil)

		v := NewValsetModule()
		v.Chain = mockChain

		cnPeers, err := v.GetCNPeers(11)
		require.NoError(t, err)
		assert.Nil(t, cnPeers)
	})

	t.Run("post-fork empty peer set returns nil to disable filtering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		v := NewValsetModule()
		v.Chain = mockChain
		v.transitionResultCache.Add(uint64(11), &TransitionResult{Nodes: NodeMap{}})

		cnPeers, err := v.GetCNPeers(11)
		require.NoError(t, err)
		assert.Nil(t, cnPeers)
	})
}

func newCachedPermissionlessValset(t *testing.T, nodes NodeMap) *ValsetModule {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()
	// The committee/qualified getters prefer the canonical header; return nil so
	// they fall back to the pre-cached transition result (the state path under test).
	mockChain.EXPECT().GetHeaderByNumber(gomock.Any()).Return(nil).AnyTimes()

	v := NewValsetModule()
	v.Chain = mockChain
	v.transitionResultCache.Add(testGetterBlock, &TransitionResult{Nodes: nodes})
	return v
}

// TestPermissionlessCommitteeFromHeader checks that GetCommittee and
// GetQualifiedValidators serve the set from the canonical header: no
// transitionResultCache is set up, so any state read would fail the test.
func TestPermissionlessCommitteeFromHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	sealer := engine.NewSealer(nil, nil)

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()
	mockChain.EXPECT().Sealer().Return(sealer).AnyTimes()

	header := &types.Header{Number: new(big.Int).SetUint64(testGetterBlock)}
	require.NoError(t, sealer.WriteValidators(header, numsToAddrs(1, 2)))
	mockChain.EXPECT().GetHeaderByNumber(testGetterBlock).Return(header).AnyTimes()

	v := NewValsetModule()
	v.Chain = mockChain

	committee, err := v.GetCommittee(testGetterBlock, 0)
	require.NoError(t, err)
	assert.Equal(t, numsToAddrs(1, 2), committee)

	qualified, err := v.GetQualifiedValidators(testGetterBlock)
	require.NoError(t, err)
	assert.Equal(t, numsToAddrs(1, 2), qualified)
}

// GetProposer answers "who was scheduled for this round", never "who signed this block". A
// hash-locked block commits at a later round while still carrying the seal of the validator that
// built it, so reading the author here would name a node that did not propose the committing round.
// Re-sealing the same header by a different validator must not move the answer.
func TestGetProposerIgnoresHeaderAuthor(t *testing.T) {
	var (
		ctrl      = gomock.NewController(t)
		sealer    = engine.NewSealer(nil, nil)
		mockChain = chain_mock.NewMockBlockChain(ctrl)
		mockGov   = gov_mock.NewMockGovModule(ctrl)
		keyA, _   = crypto.GenerateKey()
		keyB, _   = crypto.GenerateKey()
		addrA     = crypto.PubkeyToAddress(keyA.PublicKey)
		addrB     = crypto.PubkeyToAddress(keyB.PublicKey)
		qualified = []common.Address{addrA, addrB, numToAddr(3), numToAddr(4)}
	)

	header := &types.Header{Number: big.NewInt(1)}
	require.NoError(t, sealer.WriteValidators(header, qualified))
	sealer.WriteRound(header, 2) // committed at round 2

	mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()
	mockChain.EXPECT().Sealer().Return(sealer).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(1)).Return(header).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(&types.Header{Number: big.NewInt(0)}).AnyTimes()
	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		ProposerPolicy: uint64(istanbul.WeightedRandom),
		MinimumStake:   big.NewInt(0),
		GovernanceMode: "none",
	}).AnyTimes()

	v := NewValsetModule()
	v.Chain = mockChain
	v.GovModule = mockGov

	// The schedule's answer for round 2, computed without going through GetProposer.
	c, err := v.getBlockContext(1)
	require.NoError(t, err)
	roundProposer, err := v.getProposer(c, 2)
	require.NoError(t, err)

	// Seal the block as a validator that is NOT the round-2 proposer: that is the hash lock
	// shape, where the block still carries the seal of the round it was originally built at.
	authorKey, author := keyA, addrA
	if roundProposer == addrA {
		authorKey, author = keyB, addrB
	}
	require.NotEqual(t, roundProposer, author)
	seal, err := istanbul.NewSealerImpl(authorKey).MakeAuthorSeal(header)
	require.NoError(t, err)
	require.NoError(t, sealer.WriteAuthorSeal(header, seal))
	sealed, err := sealer.Author(header)
	require.NoError(t, err)
	require.Equal(t, author, sealed, "the header really is sealed by `author`")

	got, err := v.GetProposer(1, 2)
	require.NoError(t, err)
	assert.Equal(t, roundProposer, got, "the round's proposer answers")
	assert.NotEqual(t, author, got, "the header's author does not")
}
