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

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// GetCouncil returns the whole validator list for validating the block `num`.
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	var (
		council valset.CommonAddressSet
		err     error
	)
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		council, err = v.getCouncilPermissionless(num)
	} else {
		council, err = v.getCouncilPermissioned(num)
	}
	if err != nil {
		return nil, err
	} else {
		return council.Council(), nil
	}
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
	_, demoted, err := v.getCouncilAndDemotedPermissioned(num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		// qualified = council - demoted = ValActive
		m, err := v.GetNodeByState(num, []valset.State{valset.ValActive})
		if err != nil {
			return nil, err
		}
		return valset.NewAddressSet(m.Addresses()), nil
	}
	council, demoted, err := v.getCouncilAndDemotedPermissioned(num)
	if err != nil {
		return nil, err
	}
	return council.Subtract(demoted), nil
}

func (v *ValsetModule) getCouncilAndDemotedPermissioned(num uint64) (valset.CommonAddressSet, *valset.AddressSet, error) {
	council, err := v.getCouncilPermissioned(num)
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

func (v *ValsetModule) GetNodeByState(num uint64, states []valset.State) (valset.ValidatorStateMap, error) {
	if !v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return nil, errors.New("permissionless fork is not enabled")
	}
	allNodes, err := v.getAllStateNodes(num)
	if err != nil {
		return nil, err
	}
	validatorList, ok := allNodes.(*ValidatorList)
	if !ok {
		return nil, errors.New("not a permissionless council object")
	}

	validatorList.permlessMu.RLock()
	defer validatorList.permlessMu.RUnlock()

	// empty states means return all
	if len(states) == 0 {
		return validatorList.permlessVals.Copy(), nil
	}

	desiredStates := make(map[valset.State]struct{}, len(states))
	for _, s := range states {
		desiredStates[s] = struct{}{}
	}
	filtered := make(valset.ValidatorStateMap)
	for addr, val := range validatorList.permlessVals {
		if _, ok := desiredStates[val.State]; ok {
			filtered[addr] = val
		}
	}
	return filtered.Copy(), nil
}

func (v *ValsetModule) GetEpochTransition(validators valset.ValidatorStateMap, num uint64, state *state.StateDB) (valset.ValidatorStateMap, error) {
	si, err := v.StakingModule.GetStakingInfoFromState(num, state)
	if err != nil {
		return nil, err
	}
	return v.getEpochTransition(si, num, validators), nil
}

func (v *ValsetModule) GetVrankViolationTransition(validators valset.ValidatorStateMap, num uint64, state *state.StateDB) (valset.ValidatorStateMap, error) {
	si, err := v.StakingModule.GetStakingInfoFromState(num, state)
	if err != nil {
		return nil, err
	}
	return v.deactiveStakersLessMinStakingAmount(si, num, validators), nil
}

func (v *ValsetModule) GetTimeoutTransition(validators valset.ValidatorStateMap) valset.ValidatorStateMap {
	return v.getTimeoutTransition(validators)
}

func (v *ValsetModule) GetCandidates(num uint64) ([]common.Address, error) {
	candTestings, err := v.GetNodeByState(num, []valset.State{valset.CandTesting})
	if err != nil {
		return nil, err
	}
	return candTestings.Addresses(), nil
}
