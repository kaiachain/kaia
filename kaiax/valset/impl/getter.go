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
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// GetCouncil returns the whole validator list for validating the block `num`.
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.getCouncil(num)
	}
	council, err := v.getCouncilPermissioned(num)
	if err != nil {
		return nil, err
	}
	return council.List(), nil
}

// GetDemotedValidators returns council − qualified at block `num`.
func (v *ValsetModule) GetDemotedValidators(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.getDemoted(num)
	}
	demoted, err := v.getDemotedPermissioned(num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) GetQualifiedValidators(num uint64) ([]common.Address, error) {
	var (
		qualified *valset.AddressSet
		err       error
	)
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		qualified, err = v.getQualifiedValidators(num)
	} else {
		qualified, err = v.getQualifiedPermissioned(num)
	}
	if err != nil {
		return nil, err
	}
	return qualified.List(), nil
}

// GetCommittee returns the current block's committee.
func (v *ValsetModule) GetCommittee(num uint64, round uint64) ([]common.Address, error) {
	if num == 0 {
		return v.GetCouncil(0)
	}
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.getCommittee(num)
	}

	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return nil, err
	}
	return v.getCommitteePermissioned(c, round)
}

func (v *ValsetModule) GetProposer(num, round uint64) (common.Address, error) {
	if num == 0 {
		return common.Address{}, nil
	}
	if header := v.Chain.GetHeaderByNumber(num); header != nil {
		headerRound, err := v.Chain.Sealer().Round(header)
		if err == nil && uint64(headerRound) == round {
			return v.Chain.Sealer().Author(header)
		}
	}
	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return common.Address{}, err
	}
	return v.getProposer(c, round)
}

// #region Permissionless-only public getters.
//
// Pre-fork they return permissioned equivalents for peer/voter sets, empty for
// CandTesting, and an explicit fork-disabled error for NodesByState.

// GetCandTesting returns nodes in the CandTesting state. Empty pre-fork.
func (v *ValsetModule) GetCandTesting(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.getCandTesting(num)
	}
	return []common.Address{}, nil
}

// GetCNPeers returns the CN-CN P2P peer set. Pre-fork: council.
func (v *ValsetModule) GetCNPeers(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		cnPeers, err := v.getCNPeers(num)
		if len(cnPeers) > 0 {
			return cnPeers, nil
		}
		if err != nil {
			logger.Warn("GetCNPeers: disabling CN peer filter after post-fork read failed", "num", num, "err", err)
		} else {
			logger.Warn("GetCNPeers: disabling CN peer filter after post-fork read returned empty", "num", num)
		}
		return nil, nil // allow all CN peers as fallback
	}
	cnPeers, err := v.GetCouncil(num)
	if err != nil {
		logger.Warn("GetCNPeers: disabling CN peer filter after pre-fork GetCouncil failed", "num", num, "err", err)
		return nil, nil
	}
	return cnPeers, nil
}

// GetHeaderGovVoters returns validators eligible for header governance votes. Pre-fork: council.
func (v *ValsetModule) GetHeaderGovVoters(num uint64) ([]common.Address, error) {
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.getHeaderGovVoters(num)
	}
	return v.GetCouncil(num)
}

// GetNodesByState filters NodeMap by state. Empty states means "all".
// Returns an explicit fork-disabled error pre-fork.
func (v *ValsetModule) GetNodesByState(num uint64, states []valset.NodeState) (map[common.Address]*valset.Node, error) {
	nodes, err := v.getNodesByState(num, states)
	if err != nil {
		return nil, err
	}
	return map[common.Address]*valset.Node(nodes), nil
}
