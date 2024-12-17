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
