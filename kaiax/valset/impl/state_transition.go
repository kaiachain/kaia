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
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

const (
	DefaultValPausedTimeout     = time.Hour * 8
	DefaultValIdleTimeout       = 30 * 24 * time.Hour
	DefaultActiveValidatorCount = 50
)

// getOrComputeNodeStates returns the node states for block `num`.
// If parentStatedb is provided, it is used as the parent state (avoids StateAt which may fail on pruned state).
// 1. Cache hit → return cached value.
// 2. Block N committed (header exists) → read ABv2(N) directly.
// 3. Block N not committed or state pruned → read ABv2(N-1) + apply transitions.
func (v *ValsetModule) getOrComputeNodeStates(num uint64, parentStatedb *state.StateDB) (valset.NodeStateMap, error) {
	// 1. Cache hit
	if cached, ok := v.nodeStatesCache.Get(num); ok {
		return cached.(valset.NodeStateMap), nil
	}

	// 2. Block N committed → read ABv2(N) directly (optimization: skip transition computation)
	// Falls through to case 3 if state is pruned.
	if parentStatedb == nil {
		if header := v.Chain.GetHeaderByNumber(num); header != nil {
			if statedb, err := v.Chain.StateAt(header.Root); err == nil {
				validators, err := system.ReadGetAllValidators(statedb, v.Chain, header)
				if err == nil {
					v.nodeStatesCache.Add(num, validators)
					return validators, nil
				}
			}
		}
	}

	// 3. Block N not committed or state pruned → read ABv2(N-1) + apply transitions
	if num == 0 {
		return nil, errors.New("block 0 has no committed state for permissionless")
	}
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return nil, fmt.Errorf("parent header not found for block %d", num)
	}
	if parentStatedb == nil {
		var err error
		parentStatedb, err = v.Chain.StateAt(parentHeader.Root)
		if err != nil {
			return nil, err
		}
	}
	parentValidators, err := system.ReadGetAllValidators(parentStatedb, v.Chain, parentHeader)
	if err != nil {
		return nil, err
	}

	newValidators, err := v.applyAllTransitions(parentValidators, num, parentStatedb, parentHeader)
	if err != nil {
		return nil, err
	}
	v.nodeStatesCache.Add(num, newValidators)
	return newValidators, nil
}

// applyAllTransitions applies VrankViolation → Timeout → Epoch transitions.
func (v *ValsetModule) applyAllTransitions(
	validators valset.NodeStateMap,
	num uint64,
	statedb *state.StateDB,
	header *types.Header,
) (valset.NodeStateMap, error) {
	newValidators := v.getViolationTransition(num, validators)

	backend, err := backends.NewStateBlockchainContractBackend(v.Chain, statedb.Copy())
	if err != nil {
		return nil, err
	}

	pauseTimeout, idleTimeout, err := system.ReadABv2Timeouts(backend, header.Number)
	if err != nil {
		logger.Error("Failed to read ABv2 timeouts, using defaults", "number", num, "err", err)
		pauseTimeout, idleTimeout = DefaultValPausedTimeout, DefaultValIdleTimeout
	}
	blockTime := time.Unix(header.Time.Int64(), 0)
	newValidators = v.getTimeoutTransition(newValidators, pauseTimeout, idleTimeout, blockTime)

	maxValCount, _, err := system.ReadABv2MaxCounts(backend, header.Number)
	if err != nil {
		logger.Error("Failed to read ABv2 max counts, using defaults", "number", num, "err", err)
		maxValCount = DefaultActiveValidatorCount
	}
	newValidators = v.getEpochTransition(num, newValidators, idleTimeout, int(maxValCount), blockTime)
	return newValidators, nil
}

// TODO-Permissionless: Replace with KIP-227 implementation
func (v *ValsetModule) isPassVrankTest() bool {
	return true
}

func (v *ValsetModule) isVrankEpoch(num uint64) bool {
	return num%v.vrankEpoch == 0
}

// WriteStatesToContract computes state transitions and writes updated validator states to ABv2 contract.
func (v *ValsetModule) WriteStatesToContract(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	config := v.Chain.Config()

	if config.IsPermissionlessForkEnabled(header.Number) {
		if err := v.writeNodesToContract(vmenv, header, state); err != nil {
			return err
		}
	}
	return nil
}

// InstallABv2 installs and initializes ABv2 at the HF-1 block's Finalize.
func (v *ValsetModule) InstallABv2(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	config := v.Chain.Config()
	if config.IsPermissionlessForBlockParent(header.Number) {
		return v.installAndInitializeABv2(vmenv, header, statedb)
	}
	return nil
}

// writeNodesToContract computes the diff between parent(N-1) and current(N) node states,
// and writes only the changed nodes to the AddressBookV2 contract via processSystemTransition.
// At epoch blocks, the call is always made (even with empty diff) to update epochValCount.
func (v *ValsetModule) writeNodesToContract(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	num := header.Number.Uint64()
	currentNodes, err := v.getOrComputeNodeStates(num, statedb)
	if err != nil {
		return err
	}

	// Compute diff: only nodes whose state or timeout changed.
	// Error is ignored: if parent is unavailable (nil), diffNodeStates treats all current nodes as changed (fallback to full write).
	parentNodes, _ := v.getOrComputeNodeStates(num-1, nil)
	diff := diffNodeStates(parentNodes, currentNodes)

	// Skip call if no changes and not an epoch block (epoch blocks need epochValCount update)
	if len(diff) == 0 && !v.isVrankEpoch(num) {
		return nil
	}

	config := v.Chain.Config()
	msg, from, err := prepareNodeWrite(config, header.Number, diff)
	if err != nil {
		return err
	}
	if ret, err := blockchain.SystemTxCall(msg, from, header, vmenv, statedb, config.Rules(header.Number)); err != nil {
		return fmt.Errorf("processSystemTransition failed: %w (ret=%s)", err, common.Bytes2Hex(ret))
	}
	return nil
}

// diffNodeStates returns nodes whose state or timeout changed between parent and current.
func diffNodeStates(parent, current valset.NodeStateMap) valset.NodeStateMap {
	diff := make(valset.NodeStateMap)
	for addr, cur := range current {
		prev, exists := parent[addr]
		if !exists || prev.State != cur.State || prev.IdleTimeout != cur.IdleTimeout || prev.PausedTimeout != cur.PausedTimeout {
			diff[addr] = cur
		}
	}
	return diff
}

// installAndInitializeABv2 deploys and initializes ABv2 at the hardfork block.
func (v *ValsetModule) installAndInitializeABv2(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	config := v.Chain.Config()

	// Read ABv2 implementation address from ABv2DataContract (pre-deployed by governance)
	backend, err := backends.NewStateBlockchainContractBackend(v.Chain, statedb)
	if err != nil {
		return fmt.Errorf("create contract backend: %w", err)
	}
	// nil uses CurrentBlock (parent), since the current block isn't committed yet
	implAddr, err := system.ReadABv2Implementation(backend, nil)
	if err != nil {
		logger.Error("Failed to read ABv2 implementation", "number", header.Number, "err", err)
		return err
	}

	if err := system.InstallAddressBookV2(statedb, implAddr); err != nil {
		logger.Error("Failed to install AddressBookV2", "number", header.Number, "err", err.Error())
		return err
	}
	logger.Info("Installed AddressBookV2", "number", header.Number, "impl", implAddr.Hex())

	// ABv2.initialize() reads all genesis data from ABv2DataContract via Registry(0x401).
	from, msg, err := system.EncodeInitializeABv2(config.Rules(header.Number))
	if err != nil {
		logger.Error("Failed to encode initialize ABv2", "number", header.Number, "err", err.Error())
		return err
	}
	if ret, evmErr := blockchain.SystemTxCall(msg, from, header, vmenv, statedb, config.Rules(header.Number)); evmErr != nil {
		logger.Error("Failed to call initialize ABv2", "number", header.Number, "err", evmErr, "ret", common.Bytes2Hex(ret))
		return evmErr
	}

	logger.Info("Initialized AddressBookV2", "number", header.Number)
	return nil
}

// prepareNodeWrite builds the ABI-encoded input for writing changed validator states to ABv2.
func prepareNodeWrite(
	config *params.ChainConfig,
	num *big.Int,
	nodes valset.NodeStateMap,
) (*types.Transaction, common.Address, error) {
	from, msg, err := system.EncodeWriteNodes(
		config.Rules(num),
		nodes,
	)
	if err != nil {
		logger.Error("Failed to encode processSystemTransition", "number", num.Uint64(), "err", err.Error(), "nodes", nodes.String())
		return nil, common.Address{}, err
	}
	return msg, from, err
}
