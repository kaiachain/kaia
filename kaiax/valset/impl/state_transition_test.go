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

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInstallAndInitializeABv2 tests the runtime HF migration path:
// Pre-HF state has ABv2DataContract deployed. At the HF block,
// installAndInitializeABv2 reads impl, installs proxy, and calls initialize().
func TestInstallAndInitializeABv2(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// Build a full alloc via AllocPermissionless, then remove ABv2 proxy to simulate pre-HF state
	config := system.MakeTestPermissionlessConfig(t, 3)
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
	assert.Equal(t, system.ERC1967ProxyCode, statedb.GetCode(system.AddressBookAddr))

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

// TestReadGetAllValidators tests ReadGetAllValidators before and after a state transition.
func TestReadGetAllValidators(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config := system.MakeTestPermissionlessConfig(t, 3)

	alloc, err := system.AllocPermissionless(config)
	require.NoError(t, err)

	simBackend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	defer simBackend.Close()

	chain := simBackend.BlockChain()
	header := chain.CurrentHeader()
	statedb, err := chain.StateAt(header.Root)
	require.NoError(t, err)

	// Before transition(initial state): all validators are ValActive
	validators, err := system.ReadGetAllValidators(statedb, chain, header)
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

	// After transition: nodeIds[0] should be ValPaused
	validators, err = system.ReadGetAllValidators(statedb, chain, header)
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

