// Copyright 2024 The Kaia Authors
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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/reward"
)

// VerifyHeader checks header.Rewardbase against the reward address registered for the block's
// proposer under WeightedRandom.
//   - Since Permissionless: every committee member has an AddressBookV2 profile, so the profile's
//     reward address is always enforced. A staking lookup failure is returned as is, not as an
//     invalid Rewardbase.
//   - Before: enforced only when the staking snapshot has a non-zero entry for the proposer.
func (r *RewardModule) VerifyHeader(header *types.Header, _ *types.Header) error {
	if header.Number.Sign() == 0 {
		return nil
	}
	if r.GovModule.GetParamSet(header.Number.Uint64()).ProposerPolicy != uint64(istanbul.WeightedRandom) {
		return nil
	}
	proposer, err := r.Chain.Sealer().Author(header)
	if err != nil {
		return err
	}
	expected, err := r.lookupRewardAddress(header.Number.Uint64(), proposer)
	if r.ChainConfig.IsPermissionlessForkEnabled(header.Number) {
		if err != nil {
			return fmt.Errorf("reward address lookup for block %d: %w", header.Number.Uint64(), err)
		}
		if expected == (common.Address{}) || header.Rewardbase != expected {
			return reward.ErrInvalidRewardbase
		}
		return nil
	}
	// Before Permissionless, a lookup failure or a proposer without an entry is not checked.
	if err != nil || expected == (common.Address{}) {
		return nil
	}
	if header.Rewardbase != expected {
		return reward.ErrInvalidRewardbase
	}
	return nil
}

// PrepareHeader sets header.Rewardbase to the address that receives this node's block reward,
// before any transaction is executed. The value must be final at this point: with
// reward.deferredtxfee=false the transaction fee is credited to header.Rewardbase while the
// transactions run (blockchain.StateTransition.TransitionDb, since Magma), so deciding it later
// would make the proposer execute against one address and seal another, and every other node,
// which re-executes the block with the sealed header, would reject it with a state root mismatch.
func (r *RewardModule) PrepareHeader(header *types.Header) error {
	rewardbase, err := r.effectiveRewardbase(header.Number.Uint64())
	if err != nil {
		return err
	}
	header.Rewardbase = rewardbase
	return nil
}

// effectiveRewardbase is the reward address of a block this node builds at num. It always
// passes VerifyHeader:
//   - Not WeightedRandom: the configured Rewardbase (--rewardbase).
//   - Since Permissionless: the reward address of this node's AddressBookV2 profile. Without one
//     the node cannot be in the committee, and a block with any other address would be rejected,
//     so an error is returned instead of falling back to --rewardbase.
//   - Before: the registered reward address if the staking snapshot has a non-zero one, otherwise
//     the configured Rewardbase.
//
// The lookup is keyed by NodeAddress, which must be the address the node seals with (the
// author VerifyHeader recovers from the seal); node/cn wires both from the same node key.
// A staking lookup failure is returned rather than treated as "no entry".
func (r *RewardModule) effectiveRewardbase(num uint64) (common.Address, error) {
	if r.GovModule.GetParamSet(num).ProposerPolicy != uint64(istanbul.WeightedRandom) {
		return r.Rewardbase, nil
	}
	registered, err := r.lookupRewardAddress(num, r.NodeAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("reward address lookup for block %d: %w", num, err)
	}
	if registered == (common.Address{}) {
		if r.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
			return common.Address{}, fmt.Errorf("no AddressBookV2 reward address for node %s at block %d", r.NodeAddress.Hex(), num)
		}
		logger.Trace("No reward address in the staking snapshot. Use node's rewardbase.", "number", num, "nodeAddress", r.NodeAddress, "rewardbase", r.Rewardbase)
		return r.Rewardbase, nil
	}
	logger.Trace("Use reward address for nodeValidator", "number", num, "nodeAddress", r.NodeAddress, "rewardbase", registered)
	return registered, nil
}

// lookupRewardAddress returns the reward address registered for nodeId in the staking info used
// for block num (the zero address when there is none), or the staking lookup error.
// Since Permissionless, the staking info holds every AddressBookV2 profile regardless of its state.
func (r *RewardModule) lookupRewardAddress(num uint64, nodeId common.Address) (common.Address, error) {
	sInfo, err := r.StakingModule.GetStakingInfo(num)
	if err != nil {
		return common.Address{}, err
	}
	for idx, id := range sInfo.NodeIds {
		if id != nodeId {
			continue
		}
		if idx >= len(sInfo.RewardAddrs) {
			// Only a malformed legacy staking-info record can get here; refuse rather than index past the end.
			return common.Address{}, fmt.Errorf("staking info for block %d has %d node ids but %d reward addresses", num, len(sInfo.NodeIds), len(sInfo.RewardAddrs))
		}
		return sInfo.RewardAddrs[idx], nil
	}
	return common.Address{}, nil
}
