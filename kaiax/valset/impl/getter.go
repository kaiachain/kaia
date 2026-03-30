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
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		nodes, err := v.GetNodeByState(num, valset.CouncilStates)
		if err != nil {
			return nil, err
		}
		return nodes.Addresses(), nil
	}
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	return council.List(), nil
}

// GetDemotedValidators returns the demoted validators at block `num`.
// In permissionless: demoted = council - committee = ValReady + ValPaused
// In permissioned: demoted = council members with staking < minimum
func (v *ValsetModule) GetDemotedValidators(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		m, err := v.GetNodeByState(num, []valset.State{valset.ValReady, valset.ValPaused})
		if err != nil {
			return nil, err
		}
		return m.Addresses(), nil
	}
	_, demoted, err := v.getCouncilAndDemoted(num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		// qualified = council - demoted = ValActive
		m, err := v.GetNodeByState(num, valset.CommitteeStates)
		if err != nil {
			return nil, err
		}
		return valset.NewAddressSet(m.Addresses()), nil
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

// GetNodeByState returns validators filtered by the given states at block num.
// If no states are provided, returns all nodes.
func (v *ValsetModule) GetNodeByState(num uint64, states []valset.State) (valset.NodeStateMap, error) {
	if !v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return nil, errors.New("permissionless fork is not enabled")
	}
	var nodes valset.NodeStateMap
	if cached, ok := v.nodeStatesCache.Get(num); ok {
		nodes = cached.(valset.NodeStateMap)
	} else if num == 0 {
		// Block 0 has no parent; read ABv2 directly from genesis state.
		genesisHeader := v.Chain.GetHeaderByNumber(0)
		if genesisHeader == nil {
			return nil, errParentHeaderNotFound(num)
		}
		statedb, err := v.Chain.StateAt(genesisHeader.Root)
		if err != nil {
			return nil, err
		}
		nodes, _, _, _, _, _, err = system.ReadNodeStates(statedb, v.Chain, genesisHeader)
		if err != nil {
			return nil, err
		}
	} else {
		parentHeader := v.Chain.GetHeaderByNumber(num - 1)
		if parentHeader == nil {
			return nil, errParentHeaderNotFound(num)
		}
		parentStatedb, err := v.Chain.StateAt(parentHeader.Root)
		if err != nil {
			return nil, err
		}
		nodes, err = v.getOrComputeNodeStates(num, parentStatedb)
		if err != nil {
			return nil, err
		}
	}

	// empty states means return all
	if len(states) == 0 {
		return nodes.Copy(), nil
	}

	desiredStates := make(map[valset.State]struct{}, len(states))
	for _, s := range states {
		desiredStates[s] = struct{}{}
	}
	filtered := make(valset.NodeStateMap)
	for addr, val := range nodes {
		if _, ok := desiredStates[val.State]; ok {
			filtered[addr] = val
		}
	}
	return filtered.Copy(), nil
}

func (v *ValsetModule) GetCandidates(num uint64) ([]common.Address, error) {
	candTestings, err := v.GetNodeByState(num, valset.CandidateStates)
	if err != nil {
		return nil, err
	}
	return candTestings.Addresses(), nil
}
