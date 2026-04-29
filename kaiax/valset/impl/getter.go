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
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// GetCouncil returns the whole validator list for validating the block `num`.
// In permissionless: Council = {ValActive, ValPaused}. VR excluded (voluntary standby).
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.getNodeStates(num)
		if err != nil {
			return nil, err
		}
		return nodes.Council().Addresses(), nil
	}
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	return council.List(), nil
}

// GetDemotedValidators returns the demoted validators at block `num`.
// In permissionless: demoted = council - committee = {ValPaused, suspended ValActive}
// In permissioned: demoted = council members with staking < minimum
func (v *ValsetModule) GetDemotedValidators(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.getNodeStates(num)
		if err != nil {
			return nil, err
		}
		council := nodes.Council()
		qualified, err := v.getQualifiedValidators(num)
		if err != nil {
			return nil, err
		}
		return valset.NewAddressSet(council.Addresses()).Subtract(qualified).List(), nil
	}
	_, demoted, err := v.getCouncilAndDemoted(num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.getNodeStates(num)
		if err != nil {
			return nil, err
		}
		committee := nodes.Committee()
		// Safety fallback: if all ValActive are suspended, ignore the suspended set
		// to prevent consensus halt.
		allActive := nodes.FilterByState(valset.ValActive)
		if len(committee) == 0 && len(allActive) > 0 {
			if v.lastSuspendFallbackLog != num {
				logger.Warn("all ValActive are suspended, ignoring suspended set for committee", "num", num)
				v.lastSuspendFallbackLog = num
			}
			return valset.NewAddressSet(allActive.Addresses()), nil
		}
		return valset.NewAddressSet(committee.Addresses()), nil
	}
	council, demoted, err := v.getCouncilAndDemoted(num)
	if err != nil {
		return nil, err
	}
	return council.Subtract(demoted), nil
}

func (v *ValsetModule) getCouncilAndDemoted(num uint64) (*valset.AddressSet, *valset.AddressSet, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, nil, err
	}
	demoted, err := v.getDemotedValidators(council, num)
	if err != nil {
		return nil, nil, err
	}
	return council, demoted, nil
}

// GetCommittee returns the current block's committee.
func (v *ValsetModule) GetCommittee(num uint64, round uint64) ([]common.Address, error) {
	if num == 0 {
		return v.GetCouncil(0)
	}

	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return nil, err
	}
	return v.getCommittee(c, round)
}

func (v *ValsetModule) GetProposer(num, round uint64) (common.Address, error) {
	if num == 0 {
		return common.Address{}, nil
	}
	if header := v.Chain.GetHeaderByNumber(num); header != nil {
		if uint64(header.Round()) == round {
			return v.Chain.Engine().Author(header)
		}
	}
	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return common.Address{}, err
	}
	return v.getProposer(c, round)
}

// getNodeStates returns all node states at block num. Internal use only — no Copy.
func (v *ValsetModule) getNodeStates(num uint64) (valset.NodeStateMap, error) {
	if !v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return nil, errors.New("permissionless fork is not enabled")
	}
	if cached, ok := v.nodeStatesCache.Get(num); ok {
		return cached.(nodeStatesCacheEntry).validators, nil
	}
	if num == 0 {
		// Block 0 has no parent; read ABv2 directly from genesis state.
		genesisHeader := v.Chain.GetHeaderByNumber(0)
		if genesisHeader == nil {
			return nil, errParentHeaderNotFound(num)
		}
		statedb, err := v.Chain.StateAt(genesisHeader.Root)
		if err != nil {
			return nil, err
		}
		res, err := system.ReadNodeStates(statedb, v.Chain, genesisHeader)
		if err != nil {
			return nil, err
		}
		return res.Validators, nil
	}
	// Shortcut: if block num is already committed, read NodeStates(num) from S(num) directly.
	// Initialize(num) writes NodeStates(num) to ABv2 via system tx, so ABv2(num) == NodeStates(num)
	// after commit — no applyTr needed.
	// Without this shortcut, getNodeStates(K) reads S(K-1); callers like TallyCfReport(N-1)
	// invoke getNodeStates(N-1), which would read S(N-2) — absent on a non-archive node
	// after graceful shutdown/restart (only head state is preserved on Stop).
	if header := v.Chain.GetHeaderByNumber(num); header != nil {
		if statedb, err := v.Chain.StateAt(header.Root); err == nil {
			if res, err := system.ReadNodeStates(statedb, v.Chain, header); err == nil {
				res.Validators.MarkSuspended(res.SuspendedValidators)
				v.nodeStatesCache.Add(num, nodeStatesCacheEntry{validators: res.Validators})
				return res.Validators, nil
			}
		}
	}
	parentHeader := v.Chain.GetHeaderByNumber(num - 1)
	if parentHeader == nil {
		return nil, errParentHeaderNotFound(num)
	}
	parentStatedb, err := v.Chain.StateAt(parentHeader.Root)
	if err != nil {
		return nil, err
	}
	nodes, _, err := v.getOrComputeNodeStates(num, parentStatedb)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNodeByState returns validators filtered by the given states at block num.
// If no states are provided, returns all nodes. Used by RPC API.
func (v *ValsetModule) GetNodeByState(num uint64, states []valset.State) (valset.NodeStateMap, error) {
	nodes, err := v.getNodeStates(num)
	if err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nodes.Copy(), nil
	}
	return nodes.FilterByState(states...).Copy(), nil
}

func (v *ValsetModule) GetCandTesting(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodeStates(num)
	if err != nil {
		return nil, err
	}
	return nodes.FilterByState(valset.CandTesting).Addresses(), nil
}

// GetCNPeers returns addresses of nodes that should maintain CN-CN P2P connections.
// In permissionless: CNPeers = {VA, VR, VP, CR, CT}
// In permissioned: equivalent to council (all council members are CN peers)
func (v *ValsetModule) GetCNPeers(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.getNodeStates(num)
		if err != nil {
			return nil, err
		}
		return nodes.CNPeers().Addresses(), nil
	}
	return v.GetCouncil(num)
}

// GetHeaderGovVoters returns addresses of validators eligible for governance header votes.
// In permissionless: HeaderGovVoters = {VA} excluding suspended
// In permissioned: equivalent to council
func (v *ValsetModule) GetHeaderGovVoters(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.getNodeStates(num)
		if err != nil {
			return nil, err
		}
		return nodes.HeaderGovVoters().Addresses(), nil
	}
	return v.GetCouncil(num)
}
