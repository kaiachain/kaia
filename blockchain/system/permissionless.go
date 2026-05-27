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
	abv2data "github.com/kaiachain/kaia/contracts/bindings/abv2data"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	beaconcontract "github.com/kaiachain/kaia/contracts/bindings/beacon"
	cnstakingv4 "github.com/kaiachain/kaia/contracts/bindings/cnstakingv4"
	cnstakingv4factory "github.com/kaiachain/kaia/contracts/bindings/cnstakingv4factory"
	registrycontract "github.com/kaiachain/kaia/contracts/bindings/kip149"
	pdcontract "github.com/kaiachain/kaia/contracts/bindings/publicdelegation"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

// DefaultEpochBlockInterval is the default number of blocks per epoch for ABv2.
const DefaultEpochBlockInterval = int64(params.DefaultVRankEpoch)

// AllocPermissionlessConfig holds parameters for genesis permissionless allocation.
type AllocPermissionlessConfig struct {
	Owner              common.Address
	NodeIds            []common.Address
	NodeInfos          []addressbookv2contract.NodeInfo
	StakeAmts          []*big.Int
	DataConfig         abv2data.IABv2DataContractInitData
	EpochBlockInterval int64
}

type allocPermissionlessResult struct {
	cnStakingBeacon common.Address
	pdBeacon        common.Address
	factory         common.Address
	abv2Impl        common.Address
	abv2Data        common.Address
}

// AllocPermissionless deploys all permissionless system contracts into an
// in-memory statedb using the EVM runtime, and returns the resulting genesis
// alloc.
func AllocPermissionless(config *AllocPermissionlessConfig) (map[common.Address]blockchain.GenesisAccount, error) {
	alloc, _, err := allocPermissionless(config, true)
	return alloc, err
}

// AllocPermissionlessPrerequisites deploys the contracts needed before a
// delayed permissionless HF. It does not install ABv2 at 0x400; Finalize(HF-1)
// installs and initializes ABv2 from ABv2DataContract.
func AllocPermissionlessPrerequisites(config *AllocPermissionlessConfig) (map[common.Address]blockchain.GenesisAccount, map[string]common.Address, error) {
	return allocPermissionless(config, false)
}

func allocPermissionless(config *AllocPermissionlessConfig, installABv2 bool) (map[common.Address]blockchain.GenesisAccount, map[string]common.Address, error) {
	if config == nil {
		return nil, nil, fmt.Errorf("nil permissionless config")
	}
	if len(config.NodeIds) != len(config.NodeInfos) || len(config.NodeIds) != len(config.StakeAmts) {
		return nil, nil, fmt.Errorf("mismatched lengths: nodeIds=%d, infos=%d, stakeAmts=%d",
			len(config.NodeIds), len(config.NodeInfos), len(config.StakeAmts))
	}

	memDB := database.NewMemoryDBManager()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(memDB), nil, nil)

	registryConfig := &params.RegistryConfig{
		Records: make(map[string]common.Address),
		Owner:   config.Owner,
	}
	if err := InstallRegistry(statedb, registryConfig); err != nil {
		return nil, nil, fmt.Errorf("install registry: %w", err)
	}

	deployer := config.Owner
	for i, amt := range config.StakeAmts {
		statedb.AddBalance(config.NodeInfos[i].Manager, amt)
	}

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
	if err := deployBeaconInfra(cfg, deployer, config.EpochBlockInterval, result); err != nil {
		return nil, nil, err
	}
	if err := deployCnStakingFactory(cfg, big.NewInt(1), result); err != nil {
		return nil, nil, err
	}
	if err := deployCnStakingPerValidator(cfg, config, result); err != nil {
		return nil, nil, err
	}
	cfg.Origin = deployer
	if err := deployABv2DataContract(cfg, big.NewInt(1), config, result); err != nil {
		return nil, nil, err
	}

	// The registry entries are initially registered at activation=1, so the
	// initialization calls can resolve them while running against the genesis
	// staging state. The activation slots are patched back to zero below.
	if installABv2 {
		cfg.BlockNumber = big.NewInt(1)
		if err := installAndInitABv2(cfg, statedb, result.abv2Impl); err != nil {
			return nil, nil, err
		}
		if err := applyInitialNodeStateOverrides(cfg, config); err != nil {
			return nil, nil, err
		}
	}

	patchRegistryActivations(statedb, "CnStakingFactory", "ABv2DataContract")

	for _, info := range config.NodeInfos {
		statedb.SetBalance(info.Manager, new(big.Int))
	}

	if _, err := statedb.Commit(false); err != nil {
		return nil, nil, fmt.Errorf("statedb commit: %w", err)
	}

	excludeAddrs := map[common.Address]bool{}
	for _, info := range config.NodeInfos {
		excludeAddrs[info.Manager] = true
	}
	alloc, err := extractAlloc(statedb, excludeAddrs)
	if err != nil {
		return nil, nil, fmt.Errorf("extract alloc: %w", err)
	}
	records := map[string]common.Address{
		"CnStakingFactory": result.factory,
		"ABv2DataContract": result.abv2Data,
	}
	return alloc, records, nil
}

func deployBeaconInfra(cfg *runtime.Config, owner common.Address, epochBlockInterval int64, result *allocPermissionlessResult) error {
	cnImplAddr, err := evmCreate(cfg, common.FromHex(cnstakingv4.CnStakingV4Bin))
	if err != nil {
		return fmt.Errorf("deploy CnStakingV4 impl: %w", err)
	}

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

	abv2ABI, _ := addressbookv2contract.AddressBookV2MetaData.GetAbi()
	if epochBlockInterval == 0 {
		epochBlockInterval = DefaultEpochBlockInterval
	}
	abv2ImplInput, err := packConstructor(abv2ABI, common.FromHex(addressbookv2contract.AddressBookV2Bin), big.NewInt(epochBlockInterval))
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
	result.abv2Data = dataContractAddr

	registryABI, _ := registrycontract.RegistryMetaData.GetAbi()
	if err := evmCallABI(cfg, RegistryAddr, registryABI, "register", "ABv2DataContract", dataContractAddr, activation); err != nil {
		return fmt.Errorf("register ABv2DataContract: %w", err)
	}
	return nil
}

func installAndInitABv2(cfg *runtime.Config, statedb *state.StateDB, implAddr common.Address) error {
	if err := InstallAddressBookV2(statedb, implAddr); err != nil {
		return fmt.Errorf("install ABv2: %w", err)
	}
	if err := evmCallABI(cfg, AddressBookAddr, AddressBookV2ABI, "initialize"); err != nil {
		return fmt.Errorf("ABv2.initialize: %w", err)
	}
	statedb.SetState(AddressBookAddr, common.BigToHash(big.NewInt(12)), common.BigToHash(big.NewInt(1)))
	return nil
}

func applyInitialNodeStateOverrides(cfg *runtime.Config, config *AllocPermissionlessConfig) error {
	var (
		valActiveState = valset.ValActive.ToUint8()
		nodeIds        []common.Address
		newStates      []uint8
		timeoutAts     []*big.Int
	)
	for i, info := range config.NodeInfos {
		if info.State != valActiveState {
			nodeIds = append(nodeIds, config.NodeIds[i])
			newStates = append(newStates, info.State)
			timeoutAts = append(timeoutAts, new(big.Int))
		}
	}
	if len(nodeIds) == 0 {
		return nil
	}

	abv2ABI, _ := addressbookv2contract.AddressBookV2MetaData.GetAbi()
	origOrigin := cfg.Origin
	cfg.Origin = params.SystemAddress
	defer func() { cfg.Origin = origOrigin }()

	epochVACount := big.NewInt(int64(len(config.NodeInfos) - len(nodeIds)))
	if err := evmCallABI(cfg, AddressBookAddr, abv2ABI, "processSystemTransition", nodeIds, newStates, timeoutAts, epochVACount); err != nil {
		return fmt.Errorf("applyInitialNodeStateOverrides: %w", err)
	}
	return nil
}

func patchRegistryActivations(statedb *state.StateDB, names ...string) {
	for _, name := range names {
		arraySlot := calcMappingSlot(0, name, 0)
		activationSlot := calcArraySlot(arraySlot, 2, 0, 1)
		statedb.SetState(RegistryAddr, activationSlot, common.Hash{})
	}
}

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
			Name:     info.Name,
			Metadata: info.Metadata,
			State:    info.State,
		}
	}
	d := config.DataConfig
	return abv2data.IABv2DataContractInitData{
		InitialOwner:            d.InitialOwner,
		InitialSuspender:        d.InitialSuspender,
		InitialConfigurator:     d.InitialConfigurator,
		PfsThreshold:            d.PfsThreshold,
		CfsThreshold:            d.CfsThreshold,
		PauseTimeout:            d.PauseTimeout,
		IdleTimeout:             d.IdleTimeout,
		MaxNodeCount:            d.MaxNodeCount,
		MaxValActivePausedCount: d.MaxValActivePausedCount,
		MaxCandReadyCount:       d.MaxCandReadyCount,
		KefAddress:              d.KefAddress,
		KifAddress:              d.KifAddress,
		KpfAddress:              d.KpfAddress,
		NodeIds:                 config.NodeIds,
		Infos:                   infos,
	}
}

func evmCreate(cfg *runtime.Config, input []byte) (common.Address, error) {
	ret, addr, _, err := runtime.Create(input, cfg)
	if err != nil {
		return common.Address{}, fmt.Errorf("%w (revert: %x)", err, ret)
	}
	return addr, nil
}

func evmCallABI(cfg *runtime.Config, to common.Address, parsedABI *kaiaABI.ABI, method string, args ...interface{}) error {
	_, err := evmCallABIReturn(cfg, to, parsedABI, method, args...)
	return err
}

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

func packConstructor(parsedABI *kaiaABI.ABI, bytecode []byte, args ...interface{}) ([]byte, error) {
	packed, err := parsedABI.Pack("", args...)
	if err != nil {
		return nil, err
	}
	return append(bytecode, packed...), nil
}

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
