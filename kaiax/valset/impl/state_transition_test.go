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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	addressbookv2contract "github.com/kaiachain/kaia/contracts_permissionless/contracts/AddressBookV2"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInstallAndInitializeABv2 tests the runtime HF migration path:
// Pre-HF state has ABv2DataContract deployed. At the HF block,
// installAndInitializeABv2 reads impl, installs proxy, and calls initialize().
func TestInstallAndInitializeABv2(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// Build a full alloc via AllocPermissionless, then remove ABv2 proxy to simulate pre-HF state
	config, _ := system.MakeTestPermissionlessConfig(3)
	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)
	delete(alloc, system.AddressBookAddr) // pre-HF: no ABv2 proxy at 0x400

	// Create SimulatedBackend with PermissionlessForkBlock=2 so InstallABv2 triggers at block 1 (HF-1)
	chainConfig := params.TestChainConfig.Copy()
	chainConfig.PermissionlessCompatibleBlock = big.NewInt(2)
	simBackend := backends.NewSimulatedBackendWithChainConfig(blockchain.GenesisAlloc(alloc), chainConfig)
	defer simBackend.Close()
	chain := simBackend.BlockChain()

	// Create ValsetModule with the chain
	v := NewValsetModule()
	v.Chain = chain

	// Get statedb at current block (block 0 = pre-HF)
	header := chain.CurrentHeader()
	statedb, err := chain.StateAt(header.Root)
	require.NoError(t, err)
	assert.Empty(t, statedb.GetCode(system.AddressBookAddr), "pre-HF: 0x400 should have no code")

	// Create EVM for the HF block (block 1)
	hfHeader := &types.Header{
		Number:     big.NewInt(1),
		Time:       big.NewInt(0),
		BlockScore: big.NewInt(1),
	}
	blockCtx := blockchain.NewEVMBlockContext(hfHeader, chain, nil)
	vmenv := vm.NewEVM(blockCtx, vm.TxContext{}, statedb, chain.Config(), &vm.Config{})

	// Call the public entry point — same as Finalize calls at HF-1 block
	err = v.InstallABv2(vmenv, hfHeader, statedb)
	require.NoError(t, err)

	// Verify: proxy code is set at 0x400
	assert.Equal(t, system.ERC1967ProxyV5Code, statedb.GetCode(system.AddressBookAddr))

	// Verify: ABv2 is initialized by querying getAllProfiles
	// Need to create a backend that can read from the modified statedb
	readBackend, err := backends.NewStateBlockchainContractBackend(chain, statedb)
	require.NoError(t, err)
	caller, err := addressbookv2contract.NewAddressBookV2Caller(system.AddressBookAddr, readBackend)
	require.NoError(t, err)

	profiles, err := caller.GetAllProfiles(&bind.CallOpts{})
	require.NoError(t, err)
	assert.Len(t, profiles, len(config.NodeIds))

	for _, p := range profiles {
		assert.Equal(t, valset.ValActive.ToUint8(), p.State)
		assert.NotEqual(t, common.Address{}, p.StakingContract)
	}
}

// TestGetOrComputeNodeStates tests getOrComputeNodeStates which should always
// return ABv2(N-1) + applyTr(N-1), excluding user transactions at block N.
// Uses GenerateChain + InsertChain to commit a real pause tx for realistic chain state.
func TestGetOrComputeNodeStates(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	config, nodeKeys := system.MakeTestPermissionlessConfig(3)
	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	// Build chain with genesis — fund config.NodeIds[0] for pause tx
	db := database.NewMemoryDBManager()
	genesisAlloc := blockchain.GenesisAlloc(alloc)
	genesisAlloc[config.NodeIds[0]] = blockchain.GenesisAccount{Balance: big.NewInt(1_000_000_000)}
	gspec := &blockchain.Genesis{
		Config: params.TestChainConfig,
		Alloc:  genesisAlloc,
	}
	genesis := gspec.MustCommit(db)

	// Generate block 1 with a pause tx for config.NodeIds[0]
	pauseData, err := system.AddressBookV2ABI.Pack("pause", config.NodeIds[0])
	require.NoError(t, err)

	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	blocks, _ := blockchain.GenerateChain(params.TestChainConfig, genesis, faker.NewFaker(), db, 1, func(i int, gen *blockchain.BlockGen) {
		tx, err := types.SignTx(
			types.NewTransaction(gen.TxNonce(config.NodeIds[0]), system.AddressBookAddr, common.Big0, 1_000_000, common.Big0, pauseData),
			signer, nodeKeys[0],
		)
		require.NoError(t, err)
		gen.AddTx(tx)
	})

	// Insert block 1
	chain, err := blockchain.NewBlockChain(db, nil, params.TestChainConfig, faker.NewFaker(), vm.Config{})
	require.NoError(t, err)
	defer chain.Stop()
	_, err = chain.InsertChain(blocks)
	require.NoError(t, err)

	// Verify ABv2(1) has ValPaused from pause tx
	header1 := chain.GetHeaderByNumber(1)
	statedb1, err := chain.StateAt(header1.Root)
	require.NoError(t, err)
	validators, _, _, _, _, _, err := system.ReadNodeStates(statedb1, chain, header1)
	require.NoError(t, err)
	require.Equal(t, valset.ValPaused, validators[config.NodeIds[0]].State, "ABv2(1) should have ValPaused after pause tx")

	newModule := func(t *testing.T) *ValsetModule {
		ctrl := gomock.NewController(t)
		mockGov := gov_mock.NewMockGovModule(ctrl)
		mockGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{
			MinimumStake: new(big.Int).SetUint64(5_000_000),
		}).AnyTimes()

		v := NewValsetModule()
		v.Chain = chain
		v.GovModule = mockGov
		v.VRankModule = &stubVRankModule{}
		chain.Config().VRankEpoch = testVRankEpoch
		return v
	}

	genesisStatedb := func() *state.StateDB {
		s, err := chain.StateAt(chain.GetHeaderByNumber(0).Root)
		require.NoError(t, err)
		return s
	}
	block1Statedb := func() *state.StateDB {
		s, err := chain.StateAt(chain.GetHeaderByNumber(1).Root)
		require.NoError(t, err)
		return s
	}

	t.Run("cache hit returns cached value", func(t *testing.T) {
		v := newModule(t)
		cached := valset.NodeStateMap{
			config.NodeIds[0]: {State: valset.ValPaused, StakingAmount: 999},
		}
		v.nodeStatesCache.Add(uint64(1), cached)

		result, err := v.getOrComputeNodeStates(1, genesisStatedb())
		require.NoError(t, err)
		assert.Equal(t, valset.ValPaused, result[config.NodeIds[0]].State)
		assert.Equal(t, uint64(999), result[config.NodeIds[0]].StakingAmount)
	})

	t.Run("excludes userTx at block N from getOrComputeNodeStates(N)", func(t *testing.T) {
		v := newModule(t)
		// Block 1 is committed with pause tx → ABv2(1) has ValPaused.
		// But getOrComputeNodeStates(1) should read ABv2(0) + applyTr(0) = all ValActive.
		result, err := v.getOrComputeNodeStates(1, genesisStatedb())
		require.NoError(t, err)
		assert.Equal(t, valset.ValActive, result[config.NodeIds[0]].State,
			"should return ValActive from ABv2(0), not ValPaused from committed ABv2(1)")
	})

	t.Run("reads committed ABv2(N-1) including userTx from previous block", func(t *testing.T) {
		v := newModule(t)
		// pause() is an anytime transition (user tx, any block). ABv2(1) has ValPaused written by that tx.
		// Block 2: reads ABv2(1) with ValPaused + applyTr(1) (no timeout expiry, not epoch block → state unchanged)
		result, err := v.getOrComputeNodeStates(2, block1Statedb())
		require.NoError(t, err)
		assert.Equal(t, valset.ValPaused, result[config.NodeIds[0]].State,
			"should read ValPaused from committed ABv2(1)")
		for _, nodeId := range config.NodeIds[1:] {
			assert.Equal(t, valset.ValActive, result[nodeId].State)
		}
	})
}

// TestReadGetAllValidators tests ReadGetAllValidators before and after a state transition.
func TestReadGetAllValidators(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config, _ := system.MakeTestPermissionlessConfig(3)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	simBackend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	defer simBackend.Close()

	chain := simBackend.BlockChain()
	header := chain.CurrentHeader()
	statedb, err := chain.StateAt(header.Root)
	require.NoError(t, err)

	// Before transition(initial state): all validators are ValActive
	validators, _, _, _, _, _, err := system.ReadNodeStates(statedb, chain, header)
	require.NoError(t, err)
	assert.Len(t, validators, 3)

	for _, nodeId := range config.NodeIds {
		vs, ok := validators[nodeId]
		require.True(t, ok, "nodeId %s should be in validators", nodeId.Hex())
		assert.Equal(t, valset.ValActive, vs.State)
		assert.Equal(t, uint64(5_000_000), vs.StakingAmount)
		assert.True(t, vs.IdleTimeout.IsZero())
		assert.True(t, vs.PausedTimeout.IsZero())
	}

	// Apply a state transition via writeNodesToContract (core path):
	// Seed cache with parent(block 0) = current validators, current(block 1) = modified validators
	v := NewValsetModule()
	v.Chain = chain

	pauseTimeout := time.Unix(9999, 0)
	modified := validators.Copy()
	modified[config.NodeIds[0]] = &valset.ValidatorState{
		State:         valset.ValPaused,
		StakingAmount: 5_000_000,
		PausedTimeout: pauseTimeout,
	}
	v.nodeStatesCache.Add(uint64(0), validators) // parent
	v.nodeStatesCache.Add(uint64(1), modified)   // current

	hfHeader := &types.Header{
		Number:     big.NewInt(1),
		Time:       big.NewInt(0),
		BlockScore: big.NewInt(1),
	}
	blockCtx := blockchain.NewEVMBlockContext(hfHeader, chain, nil)
	vmenv := vm.NewEVM(blockCtx, vm.TxContext{}, statedb, chain.Config(), &vm.Config{})
	err = v.writeNodesToContract(vmenv, hfHeader, statedb)
	require.NoError(t, err)

	// After transition: config.NodeIds[0] should be ValPaused
	validators, _, _, _, _, _, err = system.ReadNodeStates(statedb, chain, header)
	require.NoError(t, err)
	assert.Len(t, validators, 3)

	vs0 := validators[config.NodeIds[0]]
	assert.Equal(t, valset.ValPaused, vs0.State)
	assert.Equal(t, pauseTimeout, vs0.PausedTimeout)

	// Other validators unchanged
	for _, nodeId := range config.NodeIds[1:] {
		vs := validators[nodeId]
		assert.Equal(t, valset.ValActive, vs.State)
		assert.Equal(t, uint64(5_000_000), vs.StakingAmount)
		assert.True(t, vs.IdleTimeout.IsZero())
		assert.True(t, vs.PausedTimeout.IsZero())
	}
}
