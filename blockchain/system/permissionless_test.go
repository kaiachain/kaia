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
package system

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeTestPermissionlessConfig delegates to the exported MakeTestPermissionlessConfig.
func makeTestPermissionlessConfig(t *testing.T, n int) *AllocPermissionlessConfig {
	return MakeTestPermissionlessConfig(t, n)
}

// verifyPermissionlessAlloc checks alloc-level properties and ABv2 contract state.
func verifyPermissionlessAlloc(t *testing.T, config *AllocPermissionlessConfig, alloc map[common.Address]blockchain.GenesisAccount) {
	t.Helper()
	numNodes := len(config.NodeIds)

	// Alloc should contain Registry and ABv2 at known addresses
	assert.Contains(t, alloc, RegistryAddr, "Registry should be in alloc")
	assert.Contains(t, alloc, AddressBookAddr, "ABv2 should be in alloc")
	assert.NotEmpty(t, alloc[AddressBookAddr].Code, "ABv2 should have code")

	// Manager EOAs should be excluded from alloc
	for _, info := range config.NodeInfos {
		assert.NotContains(t, alloc, info.Manager, "manager EOA should be excluded from alloc")
	}

	// Alloc entries = 10 (fixed) + numNodes (CnStaking proxies):
	//   1. CnStakingV4Impl         — CREATE by owner
	//   2. CnStakingBeacon         — UpgradeableBeacon, CREATE by owner
	//   3. PDImpl                   — CREATE by owner
	//   4. PDBeacon                 — UpgradeableBeacon, CREATE by owner
	//   5. CnStakingV4Factory      — CREATE by owner
	//   6. ABv2DataContract        — CREATE by owner
	//   7. ABv2 proxy (0x400)      — ERC1967Proxy, SetCode
	//   8. ABv2 logic              — SetCode at deriveAddressBookV2LogicAddr
	//   9. Registry (0x401)        — SetCode
	//  10. Owner EOA               — deployer, nonce > 0 but no code
	//  11+. CnStaking proxy * N    — BeaconProxy, CREATE2 by Factory
	assert.Len(t, alloc, 10+numNodes)

	// Registry activation slots should be patched to 0
	registryAlloc := alloc[RegistryAddr]
	cnFactoryActivationSlot := calcArraySlot(calcMappingSlot(0, "CnStakingFactory", 0), 2, 0, 1)
	abv2DataActivationSlot := calcArraySlot(calcMappingSlot(0, "ABv2DataContract", 0), 2, 0, 1)
	assert.Equal(t, common.Hash{}, registryAlloc.Storage[cnFactoryActivationSlot], "CnStakingFactory activation should be 0")
	assert.Equal(t, common.Hash{}, registryAlloc.Storage[abv2DataActivationSlot], "ABv2DataContract activation should be 0")

	// StakingContract should have been filled in NodeInfos (side effect)
	// Each CnStaking proxy should hold the staked KAIA balance
	for i, info := range config.NodeInfos {
		assert.NotEqual(t, common.Address{}, info.StakingContract, "StakingContract should be set")
		require.Contains(t, alloc, info.StakingContract, "StakingContract should be in alloc")
		assert.NotEmpty(t, alloc[info.StakingContract].Code, "StakingContract should have code")
		assert.Equal(t, config.StakeAmts[i], alloc[info.StakingContract].Balance, "StakingContract[%d] balance should equal stakeAmt", i)
	}

	// Owner EOA: nonce > 0 (deployed contracts), no code, zero balance
	ownerAlloc := alloc[config.Owner]
	assert.Empty(t, ownerAlloc.Code, "owner should have no code")
	assert.True(t, ownerAlloc.Nonce > 0, "owner nonce should be > 0")
	assert.Equal(t, new(big.Int).SetUint64(0), ownerAlloc.Balance, "owner balance should be 0")

	// ABv2 logic contract: read address from ERC1967 implementation slot, verify it's in alloc with code
	implSlotValue := alloc[AddressBookAddr].Storage[common.BytesToHash(ImplementationSlot)]
	abv2LogicAddr := common.BytesToAddress(implSlotValue.Bytes())
	assert.NotEqual(t, common.Address{}, abv2LogicAddr, "ABv2 implementation slot should be set")
	require.Contains(t, alloc, abv2LogicAddr, "ABv2 logic should be in alloc")
	assert.NotEmpty(t, alloc[abv2LogicAddr].Code, "ABv2 logic should have code")
	// ABv2 impl has storage from constructor (_disableInitializers sets initializable slot)

	// All contract entries (except owner) should have code
	for addr, ga := range alloc {
		if addr == config.Owner {
			continue
		}
		assert.NotEmpty(t, ga.Code, "contract %s should have code", addr.Hex())
	}

	// Verify ABv2 contract state via getAllProfiles on a SimulatedBackend built from alloc
	// This verification leverages EVM execution rather the generated `alloc`, thus auxiliary validation
	backend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	caller, err := addressbookv2contract.NewAddressBookV2Caller(AddressBookAddr, backend)
	require.NoError(t, err)

	profiles, err := caller.GetAllProfiles(&bind.CallOpts{})
	require.NoError(t, err)
	assert.Len(t, profiles, numNodes)

	for i, p := range profiles {
		assert.Equal(t, config.NodeIds[i], p.NodeId, "profile[%d] NodeId", i)
		assert.Equal(t, config.NodeInfos[i].StakingContract, p.StakingContract, "profile[%d] StakingContract", i)
		assert.Equal(t, config.NodeInfos[i].RewardAddress, p.RewardAddress, "profile[%d] RewardAddress", i)
		assert.Equal(t, valset.ValActive.ToUint8(), p.State, "profile[%d] State should be ValActive", i)
	}
}

func TestAllocPermissionless(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config := makeTestPermissionlessConfig(t, 4)

	alloc, err := AllocPermissionless(config)
	require.NoError(t, err)
	verifyPermissionlessAlloc(t, config, alloc)
}

func TestAllocPermissionless_SingleValidator(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config := makeTestPermissionlessConfig(t, 1)

	alloc, err := AllocPermissionless(config)
	require.NoError(t, err)
	verifyPermissionlessAlloc(t, config, alloc)
}

func TestAllocPermissionless_MismatchedLengths(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	config := makeTestPermissionlessConfig(t, 2)

	// Remove one stakeAmt to create mismatch
	config.StakeAmts = config.StakeAmts[:1]

	_, err := AllocPermissionless(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mismatched lengths")
}
