package backend

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
)

func (sb *backend) GetValidatorSet(num uint64) (*istanbul.BlockValSet, error) {
	council, err := sb.valsetModule.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	demoted, err := sb.valsetModule.GetDemotedValidators(num)
	if err != nil {
		return nil, err
	}

	return istanbul.NewBlockValSet(council, demoted), nil
}

func (sb *backend) GetCommitteeState(num uint64) (*istanbul.RoundCommitteeState, error) {
	header := sb.chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errUnknownBlock
	}

	return sb.GetCommitteeStateByRound(num, uint64(header.Round()))
}

func (sb *backend) GetCommitteeStateByRound(num uint64, round uint64) (*istanbul.RoundCommitteeState, error) {
	blockValSet, err := sb.GetValidatorSet(num)
	if err != nil {
		return nil, err
	}

	committee, err := sb.valsetModule.GetCommittee(num, round)
	if err != nil {
		return nil, err
	}

	proposer, err := sb.valsetModule.GetProposer(num, round)
	if err != nil {
		return nil, err
	}

	committeeSize := sb.govModule.EffectiveParamSet(num).CommitteeSize
	return istanbul.NewRoundCommitteeState(blockValSet, committeeSize, committee, proposer), nil
}

// GetProposer implements istanbul.Backend.GetProposer
func (sb *backend) GetProposer(number uint64) common.Address {
	if h := sb.chain.GetHeaderByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

func (sb *backend) GetRewardAddress(num uint64, nodeId common.Address) common.Address {
	sInfo, err := sb.stakingModule.GetStakingInfo(num)
	if err != nil {
		return common.Address{}
	}

	for idx, id := range sInfo.NodeIds {
		if id == nodeId {
			return sInfo.RewardAddrs[idx]
		}
	}
	return common.Address{}
}
