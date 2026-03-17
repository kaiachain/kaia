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
	"fmt"
	"math"
	"math/big"

	kaiaABI "github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/vm/runtime"
	"github.com/kaiachain/kaia/common"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2"
	abv2data "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2/abv2data"
	cnstakingv4 "github.com/kaiachain/kaia/contracts/contracts/system_contracts/CnStaking/CnStakingV4"
	cnstakingv4factory "github.com/kaiachain/kaia/contracts/contracts/system_contracts/CnStaking/CnStakingV4Factory"
	beaconcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/Proxy/beacon"
	pdcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/PublicDelegation"
	registrycontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip149"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

// DefaultEpochBlockInterval is the default number of blocks per epoch for ABv2.
const DefaultEpochBlockInterval = 86400

// AllocPermissionlessConfig holds parameters for genesis permissionless allocation.
type AllocPermissionlessConfig struct {
	Owner      common.Address                                  // Owner of beacons, Registry registrant
	NodeIds    []common.Address                                // Validator node IDs
	NodeInfos  []addressbookv2contract.NodeInfo                // Validator info — caller fills all fields except StakingContract (set after deployCnStaking)
	StakeAmts  []*big.Int                                      // Stake amounts per validator
	DataConfig addressbookv2contract.IABv2DataContractInitData // ABv2DataContract constructor data
}

// allocPermissionlessResult holds intermediate deployed addresses passed between internal steps.
type allocPermissionlessResult struct {
	cnStakingBeacon common.Address
	pdBeacon        common.Address
	factory         common.Address
	abv2Impl        common.Address
}

// AllocPermissionless deploys all permissionless system contracts into an in-memory
// statedb using the EVM runtime, and returns the resulting genesis alloc.
func AllocPermissionless(config *AllocPermissionlessConfig) (map[common.Address]blockchain.GenesisAccount, error) {
	if len(config.NodeIds) != len(config.NodeInfos) || len(config.NodeIds) != len(config.StakeAmts) {
		return nil, fmt.Errorf("mismatched lengths: nodeIds=%d, infos=%d, stakeAmts=%d",
			len(config.NodeIds), len(config.NodeInfos), len(config.StakeAmts))
	}

	// Create in-memory statedb
	memDB := database.NewMemoryDBManager()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(memDB), nil, nil)

	// Install Registry for EVM execution — included in final alloc with patched activations
	registryConfig := &params.RegistryConfig{
		Records: make(map[string]common.Address),
		Owner:   config.Owner,
	}
	if err := InstallRegistry(statedb, registryConfig); err != nil {
		return nil, fmt.Errorf("install registry: %w", err)
	}

	// Fund each validator's manager with enough KAIA for staking
	deployer := config.Owner
	for i, amt := range config.StakeAmts {
		statedb.AddBalance(config.NodeInfos[i].Manager, amt)
	}

	// Shared runtime config
	cfg := &runtime.Config{
		ChainConfig: params.TestChainConfig,
		Origin:      deployer,
		State:       statedb,
		GasLimit:    math.MaxUint64,
		GasPrice:    new(big.Int),
		Value:       new(big.Int),
		BlockNumber: new(big.Int),
		Time:        new(big.Int),
	}

	result := &allocPermissionlessResult{}

	// Step 1: Deploy implementation contracts and their UpgradeableBeacons
	if err := deployBeaconInfra(cfg, deployer, result); err != nil {
		return nil, err
	}

	// Step 2: Deploy CnStakingV4Factory and register in Registry (activation=1, block=0)
	if err := deployCnStakingFactory(cfg, big.NewInt(1), result); err != nil {
		return nil, err
	}

	// Step 3: Deploy CnStaking per validator and stake
	if err := deployCnStakingPerValidator(cfg, config, result); err != nil {
		return nil, err
	}
	cfg.Origin = deployer // restore after per-validator origin switching

	// Step 4: Deploy ABv2DataContract and register in Registry (activation=1, block=0)
	if err := deployABv2DataContract(cfg, big.NewInt(1), config, result); err != nil {
		return nil, err
	}

	// Advance block number so getActiveAddr finds records with activation=1
	cfg.BlockNumber = big.NewInt(1)

	// Step 5: Install ABv2 at 0x400 and call initialize()
	if err := installAndInitABv2(cfg, statedb, result.abv2Impl); err != nil {
		return nil, err
	}

	// Patch all Registry activation slots to 0 so contracts are active from genesis
	patchRegistryActivations(statedb, "CnStakingFactory", "ABv2DataContract")

	// Clean up balances for all managers
	for _, info := range config.NodeInfos {
		statedb.SetBalance(info.Manager, new(big.Int))
	}

	// Commit dirty state to trie so ForEachAccount/ForEachStorage can iterate
	if _, err := statedb.Commit(false); err != nil {
		return nil, fmt.Errorf("statedb commit: %w", err)
	}

	// Extract all changed accounts into genesis alloc.
	// Exclude manager EOAs (temporary funding) — they are not system contracts.
	// Registry (0x401) is included — this replaces allocateRegistry when permissionless is active.
	excludeAddrs := map[common.Address]bool{}
	for _, info := range config.NodeInfos {
		excludeAddrs[info.Manager] = true
	}
	alloc, err := extractAlloc(statedb, excludeAddrs)
	if err != nil {
		return nil, fmt.Errorf("extract alloc: %w", err)
	}

	return alloc, nil
}

// deployBeaconInfra deploys CnStakingV4 and PublicDelegation implementations
// along with their UpgradeableBeacons (step 1).
func deployBeaconInfra(cfg *runtime.Config, owner common.Address, result *allocPermissionlessResult) error {
	// Deploy CnStakingV4 implementation
	cnImplAddr, err := evmCreate(cfg, common.FromHex(cnstakingv4.CnStakingV4Bin))
	if err != nil {
		return fmt.Errorf("deploy CnStakingV4 impl: %w", err)
	}

	// Deploy UpgradeableBeacon for CnStaking
	beaconABI, _ := beaconcontract.UpgradeableBeaconMetaData.GetAbi()
	cnBeaconInput, err := packConstructor(beaconABI, common.FromHex(beaconcontract.UpgradeableBeaconBin), cnImplAddr, owner)
	if err != nil {
		return fmt.Errorf("pack CnStaking beacon constructor: %w", err)
	}
	cnBeaconAddr, err := evmCreate(cfg, cnBeaconInput)
	if err != nil {
		return fmt.Errorf("deploy CnStaking beacon: %w", err)
	}
	result.cnStakingBeacon = cnBeaconAddr

	// Deploy PublicDelegation implementation + UpgradeableBeacon
	pdImplAddr, err := evmCreate(cfg, common.FromHex(pdcontract.PublicDelegationBin))
	if err != nil {
		return fmt.Errorf("deploy PD impl: %w", err)
	}

	pdBeaconInput, err := packConstructor(beaconABI, common.FromHex(beaconcontract.UpgradeableBeaconBin), pdImplAddr, owner)
	if err != nil {
		return fmt.Errorf("pack PD beacon constructor: %w", err)
	}
	pdBeaconAddr, err := evmCreate(cfg, pdBeaconInput)
	if err != nil {
		return fmt.Errorf("deploy PD beacon: %w", err)
	}
	result.pdBeacon = pdBeaconAddr

	// Deploy AddressBookV2 implementation (used by ABv2DataContract and proxy setup)
	abv2ABI, _ := addressbookv2contract.AddressBookV2MetaData.GetAbi()
	abv2ImplInput, err := packConstructor(abv2ABI, common.FromHex(addressbookv2contract.AddressBookV2Bin), big.NewInt(DefaultEpochBlockInterval))
	if err != nil {
		return fmt.Errorf("pack ABv2 impl constructor: %w", err)
	}
	abv2ImplAddr, err := evmCreate(cfg, abv2ImplInput)
	if err != nil {
		return fmt.Errorf("deploy ABv2 impl: %w", err)
	}
	result.abv2Impl = abv2ImplAddr

	return nil
}

// deployCnStakingFactory deploys CnStakingV4Factory and registers it in Registry (step 2).
func deployCnStakingFactory(cfg *runtime.Config, activation *big.Int, result *allocPermissionlessResult) error {
	factoryABI, _ := cnstakingv4factory.CnStakingV4FactoryMetaData.GetAbi()
	factoryInput, err := packConstructor(factoryABI, common.FromHex(cnstakingv4factory.CnStakingV4FactoryBin), result.cnStakingBeacon, result.pdBeacon)
	if err != nil {
		return fmt.Errorf("pack factory constructor: %w", err)
	}
	factoryAddr, err := evmCreate(cfg, factoryInput)
	if err != nil {
		return fmt.Errorf("deploy factory: %w", err)
	}
	result.factory = factoryAddr

	registryABI, _ := registrycontract.RegistryMetaData.GetAbi()
	if err := evmCallABI(cfg, RegistryAddr, registryABI, "register", "CnStakingFactory", factoryAddr, activation); err != nil {
		return fmt.Errorf("register CnStakingFactory: %w", err)
	}
	return nil
}

// deployCnStakingPerValidator deploys a CnStaking proxy per validator via Factory and
// stakes KAIA via delegate() (step 3).
// Switches cfg.Origin to each validator's manager so that Factory tracks the correct deployer.
func deployCnStakingPerValidator(cfg *runtime.Config, config *AllocPermissionlessConfig, result *allocPermissionlessResult) error {
	factoryABI, _ := cnstakingv4factory.CnStakingV4FactoryMetaData.GetAbi()
	cnStakingABI, _ := cnstakingv4.CnStakingV4MetaData.GetAbi()

	for i := range config.NodeIds {
		cfg.Origin = config.NodeInfos[i].Manager
		retData, err := evmCallABIReturn(cfg, result.factory, factoryABI, "deployCnStaking", config.NodeInfos[i].Manager)
		if err != nil {
			return fmt.Errorf("deployCnStaking[%d]: %w", i, err)
		}
		proxyAddr := common.BytesToAddress(retData[12:32])

		cfg.Value = config.StakeAmts[i]
		if err := evmCallABI(cfg, proxyAddr, cnStakingABI, "delegate"); err != nil {
			return fmt.Errorf("delegate[%d]: %w", i, err)
		}
		cfg.Value = new(big.Int)

		config.NodeInfos[i].StakingContract = proxyAddr
	}
	return nil
}

// deployABv2DataContract deploys ABv2DataContract and registers it in Registry (step 4).
func deployABv2DataContract(cfg *runtime.Config, activation *big.Int, config *AllocPermissionlessConfig, result *allocPermissionlessResult) error {
	abv2DataABI, _ := abv2data.ABv2DataContractMetaData.GetAbi()
	dataInput := convertToABv2DataInitData(config)
	dataContractInput, err := packConstructor(abv2DataABI, common.FromHex(abv2data.ABv2DataContractBin), result.abv2Impl, dataInput)
	if err != nil {
		return fmt.Errorf("pack ABv2DataContract constructor: %w", err)
	}
	dataContractAddr, err := evmCreate(cfg, dataContractInput)
	if err != nil {
		return fmt.Errorf("deploy ABv2DataContract: %w", err)
	}

	registryABI, _ := registrycontract.RegistryMetaData.GetAbi()
	if err := evmCallABI(cfg, RegistryAddr, registryABI, "register", "ABv2DataContract", dataContractAddr, activation); err != nil {
		return fmt.Errorf("register ABv2DataContract: %w", err)
	}
	return nil
}

// installAndInitABv2 installs ABv2 code at 0x400 and calls initialize() (step 5).
func installAndInitABv2(cfg *runtime.Config, statedb *state.StateDB, implAddr common.Address) error {
	if err := InstallAddressBookV2(statedb, implAddr); err != nil {
		return fmt.Errorf("install ABv2: %w", err)
	}
	if err := evmCallABI(cfg, AddressBookAddr, AddressBookV2ABI, "initialize"); err != nil {
		return fmt.Errorf("ABv2.initialize: %w", err)
	}
	return nil
}

// patchRegistryActivations overwrites activation slots to 0 in Registry storage
// so that all registered contracts are active from genesis (block 0).
func patchRegistryActivations(statedb *state.StateDB, names ...string) {
	for _, name := range names {
		// Registry storage layout: records[name][0].activation @ Hash(Hash(name, 0)) + 1
		arraySlot := calcMappingSlot(0, name, 0)
		activationSlot := calcArraySlot(arraySlot, 2, 0, 1)
		statedb.SetState(RegistryAddr, activationSlot, common.Hash{})
	}
}

// convertToABv2DataInitData converts addressbookv2contract types to abv2data types.
func convertToABv2DataInitData(config *AllocPermissionlessConfig) abv2data.IABv2DataContractInitData {
	infos := make([]abv2data.NodeInfo, len(config.NodeInfos))
	for i, info := range config.NodeInfos {
		infos[i] = abv2data.NodeInfo{
			Manager:         info.Manager,
			StakingContract: info.StakingContract,
			RewardAddress:   info.RewardAddress,
			VoterAddress:    info.VoterAddress,
			TimeoutAt:       info.TimeoutAt,
			GcId:            info.GcId,
			BlsInfo: abv2data.BlsPublicKeyInfo{
				PublicKey: info.BlsInfo.PublicKey,
				Pop:       info.BlsInfo.Pop,
			},
			Metadata: info.Metadata,
			State:    info.State,
		}
	}
	d := config.DataConfig
	return abv2data.IABv2DataContractInitData{
		InitialOwner:           d.InitialOwner,
		ExitThreshold:          d.ExitThreshold,
		PauseTimeout:           d.PauseTimeout,
		IdleTimeout:            d.IdleTimeout,
		MaxValidatorCount:      d.MaxValidatorCount,
		MaxReadyCandidateCount: d.MaxReadyCandidateCount,
		KefAddress:             d.KefAddress,
		KifAddress:             d.KifAddress,
		KpfAddress:             d.KpfAddress,
		NodeIds:                config.NodeIds,
		Infos:                  infos,
	}
}

// evmCreate deploys a contract using the EVM runtime and returns the deployed address.
func evmCreate(cfg *runtime.Config, input []byte) (common.Address, error) {
	ret, addr, _, err := runtime.Create(input, cfg)
	if err != nil {
		return common.Address{}, fmt.Errorf("%w (revert: %x)", err, ret)
	}
	return addr, nil
}

// evmCallABI encodes a call using ABI and executes it. Returns error if the call fails.
func evmCallABI(cfg *runtime.Config, to common.Address, parsedABI *kaiaABI.ABI, method string, args ...interface{}) error {
	_, err := evmCallABIReturn(cfg, to, parsedABI, method, args...)
	return err
}

// evmCallABIReturn encodes a call using ABI, executes it, and returns the raw return data.
func evmCallABIReturn(cfg *runtime.Config, to common.Address, parsedABI *kaiaABI.ABI, method string, args ...interface{}) ([]byte, error) {
	input, err := parsedABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("abi pack %s: %w", method, err)
	}
	ret, _, err := runtime.Call(to, input, cfg)
	if err != nil {
		return nil, fmt.Errorf("call %s: %w (revert: %x)", method, err, ret)
	}
	return ret, nil
}

// packConstructor appends ABI-encoded constructor arguments to the bytecode.
func packConstructor(parsedABI *kaiaABI.ABI, bytecode []byte, args ...interface{}) ([]byte, error) {
	packed, err := parsedABI.Pack("", args...)
	if err != nil {
		return nil, err
	}
	return append(bytecode, packed...), nil
}

// extractAlloc iterates over all committed accounts in the statedb and builds a genesis alloc map,
// excluding the specified addresses. statedb.Commit must be called before this function.
func extractAlloc(statedb *state.StateDB, exclude map[common.Address]bool) (map[common.Address]blockchain.GenesisAccount, error) {
	alloc := make(map[common.Address]blockchain.GenesisAccount)

	statedb.ForEachAccount(func(addr common.Address, _ account.Account) {
		if exclude[addr] {
			return
		}

		ga := blockchain.GenesisAccount{
			Balance: statedb.GetBalance(addr),
			Nonce:   statedb.GetNonce(addr),
			Code:    statedb.GetCode(addr),
		}

		storage := make(map[common.Hash]common.Hash)
		statedb.ForEachStorage(addr, func(key, value common.Hash) bool {
			storage[key] = value
			return true
		})
		if len(storage) > 0 {
			ga.Storage = storage
		}

		alloc[addr] = ga
	})

	return alloc, nil
}
