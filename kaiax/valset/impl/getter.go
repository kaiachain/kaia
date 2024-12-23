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
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// GetCouncil returns the whole validator list for validating the block `num`.
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	} else {
		return council.List(), nil
	}
}

// GetDemotedValidators are subtract of qualified from council(N)
func (v *ValsetModule) GetDemotedValidators(num uint64) ([]common.Address, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	demoted, err := v.getDemotedValidators(council, num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	demoted, err := v.getDemotedValidators(council, num)
	if err != nil {
		return nil, err
	}
	return council.Subtract(demoted), nil
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

	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return common.Address{}, err
	}
	return v.getProposer(c, round)
}
