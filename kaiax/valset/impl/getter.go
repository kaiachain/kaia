package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
)

type blockResult struct {
	councilAddrList []common.Address
	staking         *staking.StakingInfo
	header          *types.Header
	author          common.Address
	pSet            gov.ParamSet
}

// consolidatedStakingAmounts get total staking amounts per staking contracts by nodeIds
func (br *blockResult) consolidatedStakingAmount() map[common.Address]uint64 {
	consolidatedStakingAmounts := make(map[common.Address]uint64, len(br.staking.NodeIds))
	for idx, nAddr := range br.staking.NodeIds {
		consolidatedStakingAmounts[nAddr] = br.staking.ConsolidatedNodes()[idx].StakingAmount
	}
	return consolidatedStakingAmounts
}

func (v *ValsetModule) getBlockResultsByNumber(num uint64) (*blockResult, error) {
	councilAddrList, err := v.GetCouncilAddressList(num)
	if err != nil {
		return nil, err
	}
	sInfo, err := v.stakingInfo.GetStakingInfo(num)
	if err != nil {
		return nil, err
	}
	header := v.chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errNilHeader
	}
	author, err := v.chain.Engine().Author(header)
	if err != nil {
		return nil, err
	}
	pSet := v.headerGov.EffectiveParamSet(num)

	return &blockResult{councilAddrList, sInfo, header, author, pSet}, nil
}

// GetCouncilAddressList returns the whole validator list of block N.
// If this network haven't voted since genesis, return genesis council which is stored at Block 0.
func (v *ValsetModule) GetCouncilAddressList(num uint64) ([]common.Address, error) {
	closestValidateVoteBlk, _ := v.headerGov.GetLatestValidatorVote(num)

	// The committee of genesis block can not be calculated because it requires a previous block.
	if closestValidateVoteBlk == 0 {
		header := v.chain.GetHeaderByNumber(num)
		if header != nil {
			return nil, errNilHeader
		}

		istanbulExtra, err := types.ExtractIstanbulExtra(header)
		if err != nil {
			return nil, errExtractIstanbulExtra
		}
		return istanbulExtra.Validators, nil
	}

	councilAddresses, err := ReadCouncilAddressListFromDb(v.ChainKv, closestValidateVoteBlk)
	if err != nil {
		return nil, err
	}
	return councilAddresses, nil
}

// GetCommitteeAddressList returns the current round or block's committee.
func (v *ValsetModule) GetCommitteeAddressList(num uint64, round uint64) ([]common.Address, error) {
	// if the block number is genesis, directly return council as committee.
	if num == 0 {
		committee, err := v.GetCouncilAddressList(0)
		if err != nil {
			return nil, err
		}
		return committee, nil
	}
	// prepare council
	council, err := v.NewCouncil(num)
	if err != nil {
		return nil, err
	}
	return council.selectCommittee(round)
}

func (v *ValsetModule) GetProposer(num uint64, round uint64) (common.Address, error) {
	// if the block number is genesis, directly return council as committee.
	if num == 0 {
		committee, err := v.GetCouncilAddressList(0)
		if err != nil {
			return common.Address{}, err
		}
		return committee[0], nil
	}
	// prepare council
	council, err := v.NewCouncil(num)
	if err != nil {
		return common.Address{}, err
	}
	proposer, idx := council.proposer(round)
	if idx == -1 {
		return common.Address{}, errUnknownProposer
	}
	return proposer, nil
}
