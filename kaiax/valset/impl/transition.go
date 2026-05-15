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
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

const (
	DefaultValPausedTimeout        = time.Hour * 8
	DefaultValIdleTimeout          = 30 * 24 * time.Hour
	DefaultMaxNodeCount            = 100
	DefaultMaxValActivePausedCount = 50
)

// #region getter

// getNodes returns the transitioned node map for producing block `num`.
// Resolves parent state from chain (StateAt(parent.Root)). Use this for
// post-commit reads (RPC, peer set queries).
func (v *ValsetModule) getNodes(num uint64) (valset.NodeMap, error) {
	if cached, ok := v.transitionResultCache.Get(num); ok {
		return cached.(*TransitionResult).Nodes, nil
	}
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return nil, errParentHeaderNotFound(num)
	}
	parentStatedb, err := v.Chain.StateAt(parentHeader.Root)
	if err != nil {
		return nil, fmt.Errorf("StateAt(%d) failed: %w", num-1, err)
	}
	result, err := v.getTransitionResult(num, parentStatedb)
	if err != nil {
		return nil, err
	}
	return result.Nodes, nil
}

func (v *ValsetModule) getTransitionResult(num uint64, parentStatedb *state.StateDB) (*TransitionResult, error) {
	if cached, ok := v.transitionResultCache.Get(num); ok {
		return cached.(*TransitionResult), nil
	}
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return nil, errParentHeaderNotFound(num)
	}
	result, err := v.applyTransition(parentHeader, parentStatedb)
	if err != nil {
		return nil, err
	}
	v.transitionResultCache.Add(num, result)
	return result, nil
}

// applyTransition does all chain/contract/gov/vrank reads and invokes applyTr(N).
// The argument header number `N` is used as an input to generate the next block.
// That is, it produces NodeStates(N+1).
// applyTransition(header(N), state(N)) = ABv2(N) + applyTr(N)
//
//	= ABv2(N) + Block(N) + Gov(N+1) + isVrankEpoch(N+1) + GetCFS(N) + GetPFS(N) + GetPfReport(N)
//
// Caller must ensure statedb == StateAt(header.Root).
func (v *ValsetModule) applyTransition(header *types.Header, statedb *state.StateDB) (*TransitionResult, error) {
	headNum := header.Number.Uint64()
	nextNum := headNum + 1

	// ABv2 read from state `headNum`
	abv2result, err := system.ReadNodeStates(statedb, v.Chain, header)
	if err != nil {
		return nil, err
	}
	abv2result.Nodes.MarkSuspended(abv2result.SuspendedValidators)

	// Gov read from `nextNum`
	pset := v.GovModule.GetParamSet(nextNum)

	// VRank reads `headNum` (fail-open via nil maps on error).
	cfs, pfs, pfReport := v.fetchVRankCtx(headNum)

	// Build the transition context.
	ctx := NewTransitionContext()
	ctx.SetBlockCtx(header, v.isVrankEpoch(nextNum))
	ctx.SetGovCtx(pset)
	ctx.SetABv2TransitionParam(abv2result.TransitionParam())
	ctx.SetSlotsCtx(slotLimitsFor(abv2result.EpochVACount))
	ctx.SetVRankCtx(cfs, pfs, pfReport)

	return ctx.ApplyAllTransitions(abv2result.Nodes), nil
}

// fetchVRankCtx pulls the three vrank scores transitions need at block num.
// Errors are swallowed here so the orchestrator stays robust to transient
// vrank failures; the resulting nil/empty maps cause transitions to fail open.
func (v *ValsetModule) fetchVRankCtx(num uint64) (cfs, pfs map[common.Address]uint64, pfReport []common.Address) {
	if v.VRankModule == nil {
		logger.Error("VRankModule is nil")
		return nil, nil, nil
	}
	if nextNum := num + 1; v.isVrankEpoch(nextNum) {
		c, err := v.VRankModule.GetCFS(num)
		if err != nil {
			logger.Warn("fetchVRankCtx: GetCFS failed", "num", num, "err", err)
		} else {
			cfs = c
		}
	}
	r, err := v.VRankModule.GetPfReport(num)
	if err != nil {
		logger.Warn("fetchVRankCtx: GetPfReport failed", "num", num, "err", err)
	} else {
		pfReport = r
	}
	if len(pfReport) > 0 {
		p, err := v.VRankModule.GetPFS(num)
		if err != nil {
			logger.Warn("fetchVRankCtx: GetPFS failed", "num", num, "err", err)
		} else {
			pfs = p
		}
	}
	return cfs, pfs, pfReport
}

func (v *ValsetModule) isVrankEpoch(num uint64) bool {
	return num%v.Chain.Config().VRankEpoch == 0
}

// #region writer

// WriteTransitionToABv2 computes the node-state transition for the current
// block and writes any diffs to AddressBookV2 via SystemTx (ABv2.processSystemTransition).
// No-op pre-fork.
// WriteTransitionToABv2(n) is for generating the next block.
// WriteTransitionToABv2(n) = Write NodeStates(N) - NodeStates(N-1) at Initialize(N).
func (v *ValsetModule) WriteTransitionToABv2(
	vmenv *vm.EVM,
	header *types.Header,
	state *state.StateDB,
) error {
	config := v.Chain.Config()
	if !config.IsPermissionlessForkEnabled(header.Number) {
		return nil
	}
	return v.writeTransitionToABv2(vmenv, header, state)
}

// writeTransitionToABv2 computes the diff between NodeStates(N-1) and NodeStates(N),
// and writes only the changed nodes to AddressBookV2 via SystemTx (ABv2.processSystemTransition).
// At epoch blocks, the call is always made (even with empty diff) to update the epochVACount snapshot.
func (v *ValsetModule) writeTransitionToABv2(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	num := header.Number.Uint64()
	tr, err := v.getTransitionResult(num, statedb)
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
	diff := diffNodeStates(parentRes.Nodes, tr.Nodes)

	// Skip the call if no changes and not an epoch block (epoch blocks need
	// the epochVACount snapshot update regardless)
	if len(diff) == 0 && !v.isVrankEpoch(num) {
		return nil
	}

	config := v.Chain.Config()
	from, msg, err := system.EncodeNodeStateUpdate(config.Rules(header.Number), diff, tr.epochVACountForWrite)
	if err != nil {
		logger.Error("Failed to encode processSystemTransition", "number", header.Number.Uint64(), "err", err.Error(), "nodes", diff.String())
		return err
	}
	if ret, err := blockchain.SystemTxCall(msg, from, header, vmenv, statedb, config.Rules(header.Number)); err != nil {
		return fmt.Errorf("processSystemTransition failed: %w (ret=%s)", err, common.Bytes2Hex(ret))
	}
	return nil
}

// diffNodeStates returns nodes whose state or timeout changed between parent and current.
// Other Node fields are intentionally omitted: processSystemTransition writes
// only lifecycle state and timeoutAt. StakingAmount comes from staking reads,
// and Suspended is derived from ABv2's suspended validator list.
func diffNodeStates(parent, current valset.NodeMap) valset.NodeMap {
	diff := make(valset.NodeMap)
	for addr, cur := range current {
		prev, exists := parent[addr]
		if !exists || prev.State != cur.State || prev.IdleTimeout != cur.IdleTimeout || prev.PausedTimeout != cur.PausedTimeout {
			copied := *cur
			diff[addr] = &copied
		}
	}
	return diff
}

// InstallABv2 installs and initializes the AddressBookV2 proxy at the
// HF-1 block's Finalize step. No-op outside that single block.
func (v *ValsetModule) InstallABv2(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	config := v.Chain.Config()
	if !config.IsPermissionlessForkBlockParent(header.Number) {
		return nil
	}
	return v.installAndInitializeABv2(vmenv, header, statedb)
}

// installAndInitializeABv2 deploys and initializes ABv2 at the hardfork block.
func (v *ValsetModule) installAndInitializeABv2(
	vmenv *vm.EVM,
	header *types.Header,
	statedb *state.StateDB,
) error {
	config := v.Chain.Config()

	// Read ABv2 implementation address from ABv2DataContract (pre-deployed by governance).
	backend, err := backends.NewStateBlockchainContractBackend(v.Chain, statedb)
	if err != nil {
		return fmt.Errorf("create contract backend: %w", err)
	}
	// nil block number → caller uses CurrentBlock (parent), since the current
	// block isn't committed yet.
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
