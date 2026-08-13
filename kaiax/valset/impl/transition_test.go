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
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/contracts/bindings/abv2data"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	"github.com/kaiachain/kaia/contracts/bindings/cnstakingv4"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	vrank_mock "github.com/kaiachain/kaia/kaiax/vrank/mock"
	"github.com/kaiachain/kaia/params"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPermissionlessConfig(permissionlessBlock int64, vrankEpoch uint64) *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	config.PermissionlessCompatibleBlock = big.NewInt(permissionlessBlock)
	config.VRankEpoch = vrankEpoch
	return config
}

func TestGetNodesCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10))

		v := NewValsetModule()
		v.Chain = mockChain
		cached := NodeMap{
			addr1: {State: ValActive, StakingAmount: aboveMinStake},
		}
		v.transitionResultCache.Add(uint64(10), &TransitionResult{Nodes: cached})

		nodes, err := v.getNodes(10)
		require.NoError(t, err)
		require.Same(t, cached[addr1], nodes[addr1])
		assert.Equal(t, cached, nodes)
	})

	t.Run("genesis reads genesis state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		root := common.HexToHash("0x1234")
		genesisHeader := &types.Header{Number: big.NewInt(0), Root: root}
		stateErr := errors.New("state unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10))
		mockChain.EXPECT().GetHeaderByNumber(uint64(0)).Return(genesisHeader)
		mockChain.EXPECT().StateAt(root).Return(nil, stateErr)

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(0)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, stateErr)
	})

	t.Run("no cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		root := common.HexToHash("0x1234")
		parentHeader := &types.Header{Number: big.NewInt(9), Root: root}
		stateErr := errors.New("state unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10))
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(parentHeader)
		mockChain.EXPECT().StateAt(root).Return(nil, stateErr)

		v := NewValsetModule()
		v.Chain = mockChain

		_, inCache := v.transitionResultCache.Get(uint64(10))
		require.False(t, inCache, "cache must be cold before first call")

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, stateErr)
	})

	t.Run("permissionless fork not enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, errPermissionlessDisabled)
	})
}

func TestGetTransitionResultCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		v := NewValsetModule()
		cached := &TransitionResult{
			Nodes: NodeMap{
				addr1: {State: ValActive, StakingAmount: aboveMinStake},
			},
		}
		v.transitionResultCache.Add(uint64(10), cached)

		result, err := v.getTransitionResult(10, nil)
		require.NoError(t, err)
		require.Same(t, cached, result)
	})

	t.Run("genesis has no parent transition", func(t *testing.T) {
		v := NewValsetModule()

		result, err := v.getTransitionResult(0, nil)
		require.Nil(t, result)
		require.ErrorContains(t, err, "parent header not found for block 0")
	})

	t.Run("no cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(nil)

		v := NewValsetModule()
		v.Chain = mockChain

		_, inCache := v.transitionResultCache.Get(uint64(10))
		require.False(t, inCache, "cache must be cold before first call")

		result, err := v.getTransitionResult(10, nil)
		require.Nil(t, result)
		require.ErrorContains(t, err, "parent header not found for block 10")
	})
}

func TestGetNodesParentReadErrors(t *testing.T) {
	t.Run("missing parent header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10))
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(nil)

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorContains(t, err, "parent header not found for block 10")
	})

	t.Run("StateAt error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		root := common.HexToHash("0x1234")
		parentHeader := &types.Header{Number: big.NewInt(9), Root: root}
		stateErr := errors.New("state unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10))
		mockChain.EXPECT().GetHeaderByNumber(uint64(9)).Return(parentHeader)
		mockChain.EXPECT().StateAt(root).Return(nil, stateErr)

		v := NewValsetModule()
		v.Chain = mockChain

		nodes, err := v.getNodes(10)
		require.Nil(t, nodes)
		require.ErrorIs(t, err, stateErr)
	})
}

func TestWriteTransitionToABv2PreForkNoOp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

	v := NewValsetModule()
	v.Chain = mockChain

	err := v.WriteTransitionToABv2(nil, &types.Header{Number: big.NewInt(99)}, nil)
	require.NoError(t, err)
}

func TestInstallABv2OutsideForkParentNoOp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(testPermissionlessConfig(100, 10))

	v := NewValsetModule()
	v.Chain = mockChain

	err := v.InstallABv2(nil, &types.Header{Number: big.NewInt(98)}, nil)
	require.NoError(t, err)
}

func TestFetchVRankCtx(t *testing.T) {
	t.Run("epoch with pf report reads all vrank inputs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		cfs := map[common.Address]uint64{addr1: 1}
		pfs := map[common.Address]uint64{addr1: 2}
		pfReport := []common.Address{addr1}

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetCFS(uint64(9)).Return(cfs, nil)
		mockVRank.EXPECT().GetPfReport(uint64(9)).Return(pfReport, nil)
		mockVRank.EXPECT().GetPFS(uint64(9)).Return(pfs, nil)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(9)
		assert.Equal(t, cfs, gotCFS)
		assert.Equal(t, pfs, gotPFS)
		assert.Equal(t, pfReport, gotPfReport)
	})

	t.Run("non-epoch without pf report skips CFS and PFS", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetPfReport(uint64(8)).Return(nil, nil)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(8)
		assert.Nil(t, gotCFS)
		assert.Nil(t, gotPFS)
		assert.Nil(t, gotPfReport)
	})

	t.Run("errors fail open", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		pfReport := []common.Address{addr1}
		vrankErr := errors.New("vrank unavailable")

		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(0, 10)).AnyTimes()

		mockVRank := vrank_mock.NewMockVRankModule(ctrl)
		mockVRank.EXPECT().GetCFS(uint64(9)).Return(nil, vrankErr)
		mockVRank.EXPECT().GetPfReport(uint64(9)).Return(pfReport, nil)
		mockVRank.EXPECT().GetPFS(uint64(9)).Return(nil, vrankErr)

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = mockVRank

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(9)
		assert.Nil(t, gotCFS)
		assert.Nil(t, gotPFS)
		assert.Equal(t, pfReport, gotPfReport)
	})

	t.Run("pre-fork block returns empty context without vrank reads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockChain := chain_mock.NewMockBlockChain(ctrl)
		mockChain.EXPECT().Config().Return(testPermissionlessConfig(30, 10)).AnyTimes()

		v := NewValsetModule()
		v.Chain = mockChain
		v.VRankModule = vrank_mock.NewMockVRankModule(ctrl)

		gotCFS, gotPFS, gotPfReport := v.fetchVRankCtx(29)
		assert.Nil(t, gotCFS)
		assert.Nil(t, gotPFS)
		assert.Nil(t, gotPfReport)
	})
}

func TestDiffNodeStates(t *testing.T) {
	oldIdleTimeout := time.Unix(100, 0)
	newIdleTimeout := time.Unix(200, 0)

	parent := NodeMap{
		addr1: {State: ValActive, StakingAmount: highStake},
		addr2: {State: ValPaused, PausedTimeout: oldIdleTimeout},
		addr3: {State: ValReady, IdleTimeout: oldIdleTimeout},
		addr4: {State: ValActive},
	}
	current := NodeMap{
		addr1: {State: ValActive, StakingAmount: lowStake}, // staking amount is not written by this path
		addr2: {State: ValActive},
		addr3: {State: ValReady, IdleTimeout: newIdleTimeout},
		addr4: {State: ValActive, Suspended: true}, // suspended is derived, not written
		addr5: {State: CandReady},
	}

	diff := diffNodeStates(parent, current)
	assert.Equal(t, []common.Address{addr2, addr3, addr5}, diff.Addresses())
	assert.Equal(t, ValActive, diff[addr2].State)
	assert.Equal(t, newIdleTimeout, diff[addr3].IdleTimeout)
	assert.Equal(t, CandReady, diff[addr5].State)

	diff[addr2].State = Registered
	assert.Equal(t, ValActive, current[addr2].State, "diff must not alias current nodes")
}

type testABv2GenesisNode struct {
	NodeID          common.Address
	Manager         common.Address
	StakingContract common.Address
	RewardAddress   common.Address
	VoterAddress    common.Address
}

func TestInstallAndInitializeABv2(t *testing.T) {
	h := newABv2Harness(t, withABv2ChainConfig(testPermissionlessConfig(3, 10)))
	v, statedb, nodes := h.Valset, h.StateDB, h.Nodes
	assert.Equal(t, system.ERC1967ProxyV5Code, statedb.GetCode(system.AddressBookAddr))

	readBackend, err := backends.NewStateBlockchainContractBackend(v.Chain, statedb)
	require.NoError(t, err)
	caller, err := addressbookv2contract.NewAddressBookV2Caller(system.AddressBookAddr, readBackend)
	require.NoError(t, err)
	profiles, err := caller.GetAllProfiles(&bind.CallOpts{})
	require.NoError(t, err)
	require.Len(t, profiles, len(nodes))

	profilesByNodeID := make(map[common.Address]addressbookv2contract.Profile, len(profiles))
	for _, profile := range profiles {
		profilesByNodeID[profile.NodeId] = profile
	}
	for _, node := range nodes {
		profile, ok := profilesByNodeID[node.NodeID]
		require.True(t, ok)
		assert.Equal(t, node.StakingContract, profile.StakingContract)
		assert.Equal(t, node.RewardAddress, profile.RewardAddress)
		assert.Equal(t, ValActive.ToUint8(), profile.State)
		assert.Zero(t, profile.TimeoutAt.Sign())
	}
}

func TestReadAndWriteABv2Snapshot(t *testing.T) {
	h := newABv2Harness(t)
	v, parentHeader, statedb, nodes := h.Valset, h.Header, h.StateDB, h.Nodes

	snapshot, err := system.ReadABv2Snapshot(statedb, v.Chain, parentHeader)
	require.NoError(t, err)
	require.Len(t, snapshot.Nodes, len(nodes))
	for _, node := range nodes {
		require.Contains(t, snapshot.Nodes, node.NodeID)
		got := snapshot.Nodes[node.NodeID]
		assert.Equal(t, ValActive, got.State)
		assert.Zero(t, got.StakingAmount)
		assert.True(t, got.IdleTimeout.IsZero())
		assert.True(t, got.PausedTimeout.IsZero())
	}

	pauseTimeout := time.Unix(9999, 0)
	modified := snapshot.Nodes.Copy()
	modified[nodes[0].NodeID].State = ValPaused
	modified[nodes[0].NodeID].PausedTimeout = pauseTimeout

	writeHeader := types.CopyHeader(parentHeader)
	writeHeader.Number = new(big.Int).Add(parentHeader.Number, common.Big1)
	writeHeader.Time = new(big.Int).Add(parentHeader.Time, common.Big1)

	v.transitionResultCache.Add(writeHeader.Number.Uint64(), &TransitionResult{Nodes: modified})
	blockCtx := blockchain.NewEVMBlockContext(writeHeader, v.Chain, nil)
	vmenv := vm.NewEVM(blockCtx, vm.TxContext{}, statedb, v.Chain.Config(), &vm.Config{})
	require.NoError(t, v.WriteTransitionToABv2(vmenv, writeHeader, statedb))

	updated, err := system.ReadABv2Snapshot(statedb, v.Chain, writeHeader)
	require.NoError(t, err)
	assert.Equal(t, ValPaused, updated.Nodes[nodes[0].NodeID].State)
	assert.Equal(t, pauseTimeout, updated.Nodes[nodes[0].NodeID].PausedTimeout)
	for _, node := range nodes[1:] {
		assert.Equal(t, ValActive, updated.Nodes[node.NodeID].State)
	}
}

func TestGetNodesExcludesUserTxAtCurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	nodes, nodeKeys := testABv2GenesisNodesWithKeys(t, 4)
	chainConfig := testPermissionlessConfig(0, 10)
	h := newABv2Harness(t,
		withABv2Nodes(nodes),
		withABv2ChainConfig(chainConfig),
	)
	simBackend := h.Backend
	chain := simBackend.BlockChain()

	mockGov := gov_mock.NewMockGovModule(ctrl)
	mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
		MinimumStake: common.Big0,
	}).AnyTimes()
	mockVRank := vrank_mock.NewMockVRankModule(ctrl)
	mockVRank.EXPECT().GetPfReport(gomock.Any()).Return(nil, nil).AnyTimes()

	v := h.Valset
	v.GovModule = mockGov
	v.VRankModule = mockVRank

	target := nodes[0].NodeID
	pauseData, err := system.AddressBookV2ABI.Pack("pause", target)
	require.NoError(t, err)
	tx := types.NewTransaction(0, system.AddressBookAddr, common.Big0, 1_000_000, common.Big0, pauseData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainConfig.ChainID), nodeKeys[0])
	require.NoError(t, err)
	// userTx is sent at block 1.
	require.NoError(t, simBackend.SendTransaction(context.Background(), signedTx))
	simBackend.Commit()

	receipt, err := simBackend.TransactionReceipt(context.Background(), signedTx.Hash())
	require.NoError(t, err)
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	header1 := chain.GetHeaderByNumber(1)
	require.NotNil(t, header1)
	state1, err := chain.StateAt(header1.Root)
	require.NoError(t, err)
	snapshot1, err := system.ReadABv2Snapshot(state1, chain, header1)
	require.NoError(t, err)
	// ABv2(1): state is ValPaused.
	require.Equal(t, ValPaused, snapshot1.Nodes[target].State)

	nodesForBlock1, err := v.getNodes(1)
	require.NoError(t, err)
	// getNodes(1): state is ValActive.
	assert.Equal(t, ValActive, nodesForBlock1[target].State, "block 1 production must exclude block 1 user tx")

	nodesForBlock2, err := v.getNodes(2)
	require.NoError(t, err)
	// getNodes(2): state is ValPaused.
	assert.Equal(t, ValPaused, nodesForBlock2[target].State, "block 2 production must include committed block 1 user tx")
}

type testABv2Harness struct {
	Backend *backends.SimulatedBackend
	Valset  *ValsetModule
	Header  *types.Header
	StateDB *state.StateDB
	Nodes   []testABv2GenesisNode
}

type abv2HarnessOpt func(*abv2HarnessOpts)

type abv2HarnessOpts struct {
	nodes       []testABv2GenesisNode
	chainConfig *params.ChainConfig
}

func withABv2Nodes(nodes []testABv2GenesisNode) abv2HarnessOpt {
	return func(o *abv2HarnessOpts) { o.nodes = nodes }
}

func withABv2ChainConfig(config *params.ChainConfig) abv2HarnessOpt {
	return func(o *abv2HarnessOpts) { o.chainConfig = config.Copy() }
}

func newABv2Harness(t *testing.T, options ...abv2HarnessOpt) *testABv2Harness {
	t.Helper()

	opts := &abv2HarnessOpts{
		nodes:       testABv2GenesisNodes(),
		chainConfig: testPermissionlessConfig(0, 10),
	}
	for _, opt := range options {
		opt(opts)
	}

	deployerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	deployer := bind.NewKeyedTransactor(deployerKey)

	alloc := testABv2Alloc(deployer.From, opts.nodes)

	preinitializedGenesis := opts.chainConfig.IsPermissionlessForkEnabled(common.Big0)
	installConfig := opts.chainConfig
	if preinitializedGenesis {
		installConfig = opts.chainConfig.Copy()
		installConfig.PermissionlessCompatibleBlock = big.NewInt(3)
	}

	simBackend := backends.NewSimulatedBackendWithChainConfig(alloc, installConfig)
	t.Cleanup(func() { _ = simBackend.Close() })

	implAddr, _, _, err := addressbookv2contract.DeployAddressBookV2(deployer, simBackend, new(big.Int).SetUint64(opts.chainConfig.VRankEpoch))
	require.NoError(t, err)
	simBackend.Commit()

	dataContractAddr, _, _, err := abv2data.DeployABv2DataContract(deployer, simBackend, implAddr, testABv2InitData(opts.nodes))
	require.NoError(t, err)
	simBackend.Commit()

	chain := simBackend.BlockChain()
	v := NewValsetModule()
	v.Chain = chain

	header := chain.CurrentHeader()
	statedb, err := chain.StateAt(header.Root)
	require.NoError(t, err)
	require.NoError(t, system.InstallRegistry(statedb, &params.RegistryConfig{
		Records: map[string]common.Address{
			"ABv2DataContract": dataContractAddr,
		},
	}))
	assert.Empty(t, statedb.GetCode(system.AddressBookAddr))

	blockCtx := blockchain.NewEVMBlockContext(header, chain, nil)
	vmenv := vm.NewEVM(blockCtx, vm.TxContext{}, statedb, chain.Config(), &vm.Config{})
	require.NoError(t, v.InstallABv2(vmenv, header, statedb))

	if !preinitializedGenesis {
		return &testABv2Harness{
			Backend: simBackend,
			Valset:  v,
			Header:  header,
			StateDB: statedb,
			Nodes:   opts.nodes,
		}
	}

	root, err := statedb.Commit(true)
	require.NoError(t, err)
	require.NoError(t, statedb.Database().TrieDB().Commit(root, false, header.Number.Uint64()))

	initializedState, err := state.New(root, statedb.Database(), nil, nil)
	require.NoError(t, err)
	initializedAlloc := genesisAllocFromState(initializedState)

	initializedBackend := backends.NewSimulatedBackendWithChainConfig(initializedAlloc, opts.chainConfig)
	t.Cleanup(func() { _ = initializedBackend.Close() })
	initializedChain := initializedBackend.BlockChain()
	initializedHeader := initializedChain.CurrentHeader()
	initializedStateDB, err := initializedChain.StateAt(initializedHeader.Root)
	require.NoError(t, err)
	initializedValset := NewValsetModule()
	initializedValset.Chain = initializedChain

	return &testABv2Harness{
		Backend: initializedBackend,
		Valset:  initializedValset,
		Header:  initializedHeader,
		StateDB: initializedStateDB,
		Nodes:   opts.nodes,
	}
}

func genesisAllocFromState(statedb *state.StateDB) blockchain.GenesisAlloc {
	alloc := make(blockchain.GenesisAlloc)
	statedb.ForEachAccount(func(addr common.Address, data account.Account) {
		storage := make(map[common.Hash]common.Hash)
		statedb.ForEachStorage(addr, func(key, value common.Hash) bool {
			if value != (common.Hash{}) {
				storage[key] = value
			}
			return true
		})
		code := statedb.GetCode(addr)
		if len(code) > 0 {
			code = append([]byte(nil), code...)
		}
		genesisAccount := blockchain.GenesisAccount{
			Code:    code,
			Balance: new(big.Int).Set(data.GetBalance()),
			Nonce:   data.GetNonce(),
		}
		if len(storage) > 0 {
			genesisAccount.Storage = storage
		}
		alloc[addr] = genesisAccount
	})
	return alloc
}

func testABv2GenesisNodes() []testABv2GenesisNode {
	return []testABv2GenesisNode{
		{
			NodeID:          common.HexToAddress("0x1001"),
			Manager:         common.HexToAddress("0x2001"),
			StakingContract: common.HexToAddress("0x3001"),
			RewardAddress:   common.HexToAddress("0x4001"),
			VoterAddress:    common.HexToAddress("0x5001"),
		},
		{
			NodeID:          common.HexToAddress("0x1002"),
			Manager:         common.HexToAddress("0x2002"),
			StakingContract: common.HexToAddress("0x3002"),
			RewardAddress:   common.HexToAddress("0x4002"),
			VoterAddress:    common.HexToAddress("0x5002"),
		},
		{
			NodeID:          common.HexToAddress("0x1003"),
			Manager:         common.HexToAddress("0x2003"),
			StakingContract: common.HexToAddress("0x3003"),
			RewardAddress:   common.HexToAddress("0x4003"),
			VoterAddress:    common.HexToAddress("0x5003"),
		},
	}
}

func testABv2GenesisNodesWithKeys(t *testing.T, count int) ([]testABv2GenesisNode, []*ecdsa.PrivateKey) {
	t.Helper()

	nodes := make([]testABv2GenesisNode, count)
	keys := make([]*ecdsa.PrivateKey, count)
	for i := range count {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		n := int64(i + 1)
		keys[i] = key
		nodes[i] = testABv2GenesisNode{
			NodeID:          crypto.PubkeyToAddress(key.PublicKey),
			Manager:         common.BigToAddress(big.NewInt(0x2000 + n)),
			StakingContract: common.BigToAddress(big.NewInt(0x3000 + n)),
			RewardAddress:   common.BigToAddress(big.NewInt(0x4000 + n)),
			VoterAddress:    common.BigToAddress(big.NewInt(0x5000 + n)),
		}
	}
	return nodes, keys
}

func testABv2Alloc(deployer common.Address, nodes []testABv2GenesisNode) blockchain.GenesisAlloc {
	alloc := blockchain.GenesisAlloc{
		deployer: {
			Balance: new(big.Int).Mul(big.NewInt(100), big.NewInt(params.KAIA)),
		},
	}

	cnStakingCode := common.FromHex("0x" + cnstakingv4.CnStakingV4BinRuntime)
	for _, node := range nodes {
		alloc[node.NodeID] = blockchain.GenesisAccount{
			Balance: new(big.Int).Mul(big.NewInt(10), big.NewInt(params.KAIA)),
		}
		alloc[node.StakingContract] = blockchain.GenesisAccount{
			Code:    cnStakingCode,
			Balance: common.Big0,
		}
	}
	return alloc
}

func testABv2InitData(nodes []testABv2GenesisNode) abv2data.IABv2DataContractInitData {
	nodeIDs := make([]common.Address, len(nodes))
	infos := make([]abv2data.NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.NodeID
		infos[i] = abv2data.NodeInfo{
			Manager:         node.Manager,
			StakingContract: node.StakingContract,
			RewardAddress:   node.RewardAddress,
			VoterAddress:    node.VoterAddress,
			TimeoutAt:       common.Big0,
			GcId:            big.NewInt(int64(i + 1)),
			BlsInfo:         testABv2BlsInfo(byte(i + 1)),
			Name:            "genesis-node",
			Metadata:        "{}",
			State:           ValActive.ToUint8(),
		}
	}
	return abv2data.IABv2DataContractInitData{
		InitialOwner:            common.HexToAddress("0x9001"),
		InitialSuspender:        common.HexToAddress("0x9002"),
		InitialConfigurator:     common.HexToAddress("0x9003"),
		PfsThreshold:            big.NewInt(10),
		CfsThreshold:            big.NewInt(10),
		PauseTimeout:            big.NewInt(3600),
		IdleTimeout:             big.NewInt(7200),
		MaxNodeCount:            big.NewInt(100),
		MaxValActivePausedCount: big.NewInt(100),
		MaxCandReadyCount:       big.NewInt(100),
		KefAddress:              common.HexToAddress("0x9004"),
		KifAddress:              common.HexToAddress("0x9005"),
		KpfAddress:              common.HexToAddress("0x9006"),
		NodeIds:                 nodeIDs,
		Infos:                   infos,
	}
}

func testABv2BlsInfo(seed byte) abv2data.BlsPublicKeyInfo {
	return abv2data.BlsPublicKeyInfo{
		PublicKey: bytes.Repeat([]byte{seed}, 48),
		Pop:       bytes.Repeat([]byte{seed + 10}, 96),
	}
}
