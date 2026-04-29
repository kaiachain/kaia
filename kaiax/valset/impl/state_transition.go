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
	DefaultValPausedTimeout        = time.Hour * 8
	DefaultValIdleTimeout          = 30 * 24 * time.Hour
	DefaultMaxNodeCount            = 100
	DefaultMaxValActivePausedCount = 50
)

// nodeStatesCacheEntry holds the result of getOrComputeNodeStates.
// epochVACount is the VA count after epoch transition (before violation); 0 for non-epoch blocks.
type nodeStatesCacheEntry struct {
	validators   valset.NodeStateMap
	epochVACount uint64
}

// getOrComputeNodeStates returns the node states for block `num` and the epoch epochVACount(0 if not epoch).
// Always computes ABv2(N-1) + applyTr(N-1) to exclude user transactions at block N.
// parentStatedb must be the committed state of block N-1.
func (v *ValsetModule) getOrComputeNodeStates(num uint64, parentStatedb *state.StateDB) (valset.NodeStateMap, uint64, error) {
	// 1. Cache hit
	if cached, ok := v.nodeStatesCache.Get(num); ok {
		entry := cached.(nodeStatesCacheEntry)
		return entry.validators, entry.epochVACount, nil
	}

	// 2. Read ABv2(N-1) + apply transitions for block N
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return nil, 0, errParentHeaderNotFound(num)
	}
	// Assert parentStatedb matches block num-1
	if root := parentStatedb.IntermediateRoot(false); root != parentHeader.Root {
		return nil, 0, fmt.Errorf("parentStatedb root mismatch for block %d: got %s, want %s", num-1, root.Hex(), parentHeader.Root.Hex())
	}
	newValidators, epochVACount, err := v.computeNodeStates(parentStatedb, parentHeader)
	if err != nil {
		return nil, 0, err
	}
	v.nodeStatesCache.Add(num, nodeStatesCacheEntry{validators: newValidators, epochVACount: epochVACount})
	return newValidators, epochVACount, nil
}

// computeNodeStates reads ABv2 validators, timeouts, and max counts in a single MultiCall,
// then applies all transitions.
func (v *ValsetModule) computeNodeStates(statedb *state.StateDB, header *types.Header) (valset.NodeStateMap, uint64, error) {
	res, err := system.ReadNodeStates(statedb, v.Chain, header)
	if err != nil {
		return nil, 0, err
	}
	res.Validators.MarkSuspended(res.SuspendedValidators)
	slotLimitsFn := func(n uint64) (uint64, uint64, error) {
		return v.readSlotLimitsFor(statedb, header, n)
	}
	return v.applyAllTransitions(res, header, slotLimitsFn)
}

// applyAllTransitions computes applyTr(N-1) given ABv2(N-1) validators and block N-1 parameters.
// NodeStates(N) = ABv2(N-1) + applyTr(N-1), where applyTr(N-1) = Epoch → Violation → Timeout.
// At epoch blocks, epoch runs first so that violation uses the new epochVACount-based slot limits.
// header must be the parent header (N-1).
// slotLimitsFn computes (maxSlotAvailable, minActiveCount) for a given epochVACount.
// Returns the final NodeStateMap and the epoch epochVACount(VA count after epoch, before violation; 0 if not epoch).
func (v *ValsetModule) applyAllTransitions(
	res *system.NodeStatesResult,
	header *types.Header, // parent header (N-1)
	slotLimitsFn func(n uint64) (uint64, uint64, error),
) (valset.NodeStateMap, uint64, error) {
	var (
		nextNum   = header.Number.Uint64() + 1
		pset      = v.GovModule.GetParamSet(nextNum)
		minStake  = pset.MinimumStake.Uint64()
		blockTime = time.Unix(header.Time.Int64(), 0)

		maxSlotAvailable = res.MaxSlotAvailable
		minActiveCount   = res.MinActiveCount
		epochVACount     uint64
	)

	newValidators := res.Validators
	if v.isVrankEpoch(nextNum) {
		newValidators = v.getEpochTransition(minStake, newValidators, res.IdleTimeout, int(res.MaxValActivePausedCount), blockTime, header.Number.Uint64(), res.CfsThreshold, res.EpochVACount)
		// Recompute slot limits from new epochVACount(VA count after epoch, before violation).
		// This epochVACountdefines the epoch's canonical slot budget and is stored in the contract.
		epochVACount = newValidators.CountByState(valset.ValActive)
		var err error
		maxSlotAvailable, minActiveCount, err = slotLimitsFn(epochVACount)
		if err != nil {
			return nil, 0, fmt.Errorf("slot limits at epoch (n=%d): %w", epochVACount, err)
		}
	}
	newValidators = v.getViolationTransition(minStake, newValidators, header.Number.Uint64(), res.PfsThreshold, maxSlotAvailable, minActiveCount, res.IdleTimeout, blockTime)
	newValidators = v.getTimeoutTransition(newValidators, res.IdleTimeout, res.PauseTimeout, blockTime)
	return newValidators, epochVACount, nil
}

// readSlotLimitsFor computes slot limits for a given epochVACountvia contract call.
func (v *ValsetModule) readSlotLimitsFor(statedb *state.StateDB, header *types.Header, n uint64) (uint64, uint64, error) {
	backend, err := backends.NewStateBlockchainContractBackend(v.Chain, statedb)
	if err != nil {
		return 0, 0, fmt.Errorf("create contract backend for slot limits: %w", err)
	}
	maxSlot, minActive, err := system.ReadSlotLimitsFor(backend, header.Number, n)
	if err != nil {
		return 0, 0, fmt.Errorf("read slot limits for n=%d: %w", n, err)
	}
	return maxSlot, minActive, nil
}

// isPassVrankTest returns true if the candidate's CFS is below the threshold.
// CFS < cfsThreshold → pass (CandTesting → ValActive promotion eligible).
func (v *ValsetModule) isPassVrankTest(addr common.Address, num, cfsThreshold, epochVACount uint64) bool {
	cfsScores, err := v.VRankModule.GetCFSWithEpochVACount(num, epochVACount)
	if err != nil { // CFS unavailable → pass to avoid blocking promotion
		logger.Warn("isPassVrankTest: GetCFSWithEpochVACount failed", "num", num, "err", err)
		return true
	}
	cfs, ok := cfsScores[addr]
	if !ok {
		return true // no score recorded → pass
	}
	if cfs >= cfsThreshold {
		logger.Info("CFS test failed: demoting to Registered", "addr", addr, "cfs", cfs, "cfsThreshold", cfsThreshold, "num", num)
		return false
	}
	return true
}

func (v *ValsetModule) isVrankEpoch(num uint64) bool {
	return num%v.Chain.Config().VRankEpoch == 0
}

// WriteStatesToContract computes state transitions and writes updated validator states to ABv2 contract.
func (v *ValsetModule) WriteStatesToContract(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	config := v.Chain.Config()

	if config.IsPermissionlessForkEnabled(header.Number) {
		if err := v.writeNodeStateUpdateToContract(vmenv, header, state); err != nil {
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

// writeNodeStateUpdateToContract computes the diff between parent(N-1) and current(N) node states,
// and writes only the changed nodes to the AddressBookV2 contract via processSystemTransition.
// At epoch blocks, the call is always made (even with empty diff) to update the epochVACount snapshot.
func (v *ValsetModule) writeNodeStateUpdateToContract(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	num := header.Number.Uint64()
	currentNodes, epochVACount, err := v.getOrComputeNodeStates(num, statedb)
	if err != nil {
		return err
	}

	// Compute diff against committed ABv2(N-1) to get only applyTr(N-1) changes.
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return errParentHeaderNotFound(num)
	}
	parentRes, err := system.ReadNodeStates(statedb, v.Chain, parentHeader)
	if err != nil {
		return fmt.Errorf("failed to read ABv2(N-1): %w", err)
	}
	diff := diffNodeStates(parentRes.Validators, currentNodes)

	// Skip call if no changes and not an epoch block (epoch blocks need epochVACount update)
	if len(diff) == 0 && !v.isVrankEpoch(num) {
		return nil
	}

	config := v.Chain.Config()
	msg, from, err := prepareNodeStateUpdate(config, header.Number, diff, epochVACount)
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
			copied := *cur
			diff[addr] = &copied
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

// prepareNodeStateUpdate builds the ABI-encoded input for writing changed validator states to ABv2.
// epochVACount is the VA count after epoch transition; 0 for non-epoch blocks.
func prepareNodeStateUpdate(
	config *params.ChainConfig,
	num *big.Int,
	nodes valset.NodeStateMap,
	epochVACount uint64,
) (*types.Transaction, common.Address, error) {
	from, msg, err := system.EncodeNodeStateUpdate(
		config.Rules(num),
		nodes,
		epochVACount,
	)
	if err != nil {
		logger.Error("Failed to encode processSystemTransition", "number", num.Uint64(), "err", err.Error(), "nodes", nodes.String())
		return nil, common.Address{}, err
	}
	return msg, from, err
}
