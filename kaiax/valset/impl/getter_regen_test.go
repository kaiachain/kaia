// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.

package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	consensus_mock "github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/require"
)

// TestGetNodeStates_MissingParentState verifies that getNodeStates(K) succeeds
// when S(K-1) is unavailable (non-archive node after graceful shutdown + restart preserves
// only the head state). getNodeStates(K) must read NodeStates(K) directly from S(K),
// not depend on S(K-1).
//
// Setup: head=4, S(4) available (real ABv2 state), S(3) returns MissingNodeError.
func TestGetNodeStates_MissingParentState(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)

	// Build a real statedb that contains an initialized ABv2.
	config, _ := system.MakeTestPermissionlessConfig(4)
	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	chainConfig := params.TestChainConfig.Copy()
	chainConfig.PermissionlessCompatibleBlock = big.NewInt(0)
	chainConfig.VRankEpoch = 5

	db := database.NewMemoryDBManager()
	genesisBlock := (&blockchain.Genesis{Config: chainConfig, Alloc: blockchain.GenesisAlloc(alloc)}).MustCommit(db)

	// Fake headers: pretend block 4 (head) has the genesis state root.
	// Block 3 has a different (fake) root that we will deliberately fail.
	headHeader := &types.Header{Number: big.NewInt(4), Root: genesisBlock.Root(), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
	parentHeader := &types.Header{Number: big.NewInt(3), Root: common.HexToHash("0xdeadbeef"), Time: big.NewInt(0), BlockScore: big.NewInt(0)}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(chainConfig).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(4)).Return(headHeader).AnyTimes()
	mockChain.EXPECT().GetHeaderByNumber(uint64(3)).Return(parentHeader).AnyTimes()
	// Return real statedb for head (S(4)), MissingNodeError for parent (S(3)).
	stateCache := state.NewDatabaseWithNewCache(db, statedb.GetEmptyTrieNodeCacheConfig())
	headState, err := state.New(genesisBlock.Root(), stateCache, nil, nil)
	require.NoError(t, err)
	mockChain.EXPECT().StateAt(genesisBlock.Root()).Return(headState, nil).AnyTimes()
	mockChain.EXPECT().StateAt(parentHeader.Root).Return(nil, &statedb.MissingNodeError{}).AnyTimes()
	// MultiCall path needs CurrentBlock for EVM context.
	mockChain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(headHeader)).AnyTimes()
	// EVM block context resolves Coinbase via chain.Engine().Author(header).
	mockEngine := consensus_mock.NewMockEngine(ctrl)
	mockEngine.EXPECT().Author(gomock.Any()).Return(common.Address{}, nil).AnyTimes()
	mockChain.EXPECT().Engine().Return(mockEngine).AnyTimes()

	mockGov := gov_mock.NewMockGovModule(ctrl)
	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: new(big.Int).SetUint64(5_000_000),
	}).AnyTimes()
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()
	mockVRank.EXPECT().GetPFS(gomock.Any()).Return(nil, nil).AnyTimes()
	mockVRank.EXPECT().GetCFSWithEpochVACount(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	v := NewValsetModule()
	v.Chain = mockChain
	v.GovModule = mockGov
	v.VRankModule = mockVRank

	// Without shortcut: getNodeStates(4) reads S(3) → ErrPrunedAncestor.
	// With shortcut: reads S(4) directly → success.
	nodes, err := v.getNodeStates(4)
	require.NoError(t, err, "getNodeStates(4) should succeed: ABv2 lives in S(4) (head); S(3) missing must NOT block this")
	require.Len(t, nodes, len(config.NodeIds))
	for _, id := range config.NodeIds {
		_, ok := nodes[id]
		require.True(t, ok, "nodeId %s should be in NodeStates(4)", id.Hex())
	}
}
