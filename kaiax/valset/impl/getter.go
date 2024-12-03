package impl

import (
	"sort"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// GetCouncilAddressList returns the whole validator list of block N.
// If this network haven't voted since genesis, return genesis council which is stored at Block 0.
func (v *ValsetModule) GetCouncilAddressList(num uint64) ([]common.Address, error) {
	scannedNum, err := readLowestScannedCheckpointIntervalNum(v.ChainKv)
	if err != nil {
		return nil, err
	}

	// council(N-1)+header[N-1]=council(N), so council(N) is finalized at block N.
	// That is why the key of the db is the finalized block num: snap[N-1], valSet[N-1]
	finalizedBlockNum := uint64(0)
	if num > 0 {
		finalizedBlockNum = num - 1
	}

	// if valSet council db is not prepared (target block number is less than lowest scanned number)
	// read the council address list from Istanbul snapshot
	// then, regenerate the target block's council address list by iterating the votes in header within snapshot interval
	if finalizedBlockNum < scannedNum {
		if finalizedBlockNum == 0 {
			header := v.chain.GetHeaderByNumber(num)
			if header == nil {
				return nil, errNilHeader
			}
			istanbulExtra, err := types.ExtractIstanbulExtra(header)
			if err != nil {
				return nil, errExtractIstanbulExtra
			}
			return istanbulExtra.Validators, nil
		}
		// snap[checkpoint]+header[checkpoint+1]+header[checkpoint+2] => snap[checkpoint+2] => council(checkpoint+3)
		return v.replayValSetVotes(finalizedBlockNum-finalizedBlockNum%params.CheckpointInterval, finalizedBlockNum, false)
	}

	// make sure that the num doesn't exceed current block
	// valSet council db is prepared. let's read council address list directly from the db
	return readCouncilAddressListFromValSetCouncilDB(v.ChainKv, finalizedBlockNum)
}

// GetCommittee returns the current block's committee.
func (v *ValsetModule) GetCommittee(num uint64, round uint64) ([]common.Address, error) {
	c, err := newCouncil(v, num)
	if err != nil {
		return nil, err
	}

	// if the block number is genesis, directly return councilAddressList as committee.
	if num == 0 {
		return c.councilAddressList, nil
	}

	pSet, proposerPolicy, err := v.getPSetWithProposerPolicy(num)
	if err != nil {
		return nil, err
	}

	if pSet.CommitteeSize >= uint64(len(c.qualifiedValidators)) {
		return c.qualifiedValidators, nil
	}

	proposer, err := v.GetProposer(num, round)
	if err != nil {
		return nil, err
	}

	// return early if the committee size is 1
	if pSet.CommitteeSize == 1 {
		return []common.Address{proposer}, nil
	}

	prevHeader := v.chain.GetHeaderByNumber(num - 1)
	if prevHeader == nil {
		return nil, errNilHeader
	}

	if !proposerPolicy.IsWeightedRandom() || !c.rules.IsRandao {
		// closest next proposer who has different address with the proposer
		// pick current round's proposer and next proposer which address is different from current proposer
		var nextDistinctProposer common.Address
		for i := uint64(1); i < pSet.ProposerUpdateInterval; i++ {
			nextDistinctProposer, err = v.GetProposer(num, round+i)
			if err != nil {
				return nil, err
			}
			if proposer != nextDistinctProposer {
				break
			}
		}

		prevAuthor, err := v.chain.Engine().Author(prevHeader) // author of block N-1
		if err != nil {
			return nil, err
		}

		return c.selectRandomCommittee(round, pSet, proposer, nextDistinctProposer, prevHeader, prevAuthor)
	}
	return c.selectRandaoCommittee(prevHeader, pSet)
}

// GetProposer calculates a proposer for the round of the given block.
// Note that, to skip the calculation, derive the author from header if the block exists
// - if the policy defaultSet, the proposer is picked from the council. it's defaultSet, so there's no demoted.
// - if the policy is weightedrandom and before Randao hf, the proposer is picked from the proposers
// - if the pilicy is weightedrandom and after Randao hf, the proposer is picked from the committee
func (v *ValsetModule) GetProposer(num uint64, round uint64) (common.Address, error) {
	header := v.chain.GetHeaderByNumber(num)
	if header != nil && header.Number.Uint64() == num && uint64(header.Round()) == round {
		return v.chain.Engine().Author(header)
	}

	c, err := newCouncil(v, num)
	if err != nil {
		return common.Address{}, err
	}

	// if the block number is genesis, directly return proposer as the first element of councilAddressList.
	if num == 0 {
		return c.councilAddressList[0], nil
	}

	pSet, proposerPolicy, err := v.getPSetWithProposerPolicy(num)
	if err != nil {
		return common.Address{}, err
	}

	prevHeader := v.chain.GetHeaderByNumber(num - 1)
	if prevHeader == nil {
		return common.Address{}, errNilHeader
	}
	prevAuthor, err := v.chain.Engine().Author(prevHeader) // author of block N-1
	if err != nil {
		return common.Address{}, err
	}

	if proposerPolicy.IsDefaultSet() {
		// if the policy is round-robin or sticky, all the council members are qualified.
		// be cautious that the proposer may not be included in the committee list.
		copied := make(valset.AddressList, len(c.councilAddressList))
		copy(copied, c.councilAddressList)

		// sorting on council address list
		sort.Sort(copied)

		proposer, _ := pickRoundRobinProposer(c.councilAddressList, proposerPolicy, prevAuthor, round)
		return proposer, nil
	}

	// before Randao, weightedrandom uses proposers to pick the proposer.
	if !c.rules.IsRandao {
		pUpdateBlock := calcProposerBlockNumber(num, pSet.ProposerUpdateInterval)
		proposers, err := v.getProposers(pUpdateBlock)
		if err != nil {
			return common.Address{}, err
		}
		proposer, _ := pickWeightedRandomProposer(proposers, pUpdateBlock, num, round, c.qualifiedValidators, prevAuthor)
		return proposer, nil
	}

	// after Randao, pick proposer from randao committee
	committee, err := c.selectRandaoCommittee(prevHeader, pSet)
	if err != nil {
		return common.Address{}, err
	}
	return committee[int(round)%len(c.qualifiedValidators)], nil
}

// getPSetWithProposerPolicy returns the govParam & proposer policy after processing block num - 1
func (v *ValsetModule) getPSetWithProposerPolicy(num uint64) (gov.ParamSet, ProposerPolicy, error) {
	pSet := v.headerGov.EffectiveParamSet(num)
	if pSet.CommitteeSize == 0 {
		return gov.ParamSet{}, 0, errInvalidCommitteeSize // it cannot happen. just to make sure
	}
	// if config.Istanbul is nil, it means the consensus is not 'istanbul' so use default proposer policy( = RoundRobin).
	proposerPolicy := ProposerPolicy(params.DefaultProposerPolicy)
	if v.chain.Config().Istanbul != nil {
		proposerPolicy = ProposerPolicy(pSet.ProposerPolicy)
	}
	return pSet, proposerPolicy, nil
}

// getStakingInfoWithStakingAmounts returns the stakingInfo & parsed staking amounts after processing block num - 1
func (v *ValsetModule) getStakingInfoWithStakingAmounts(num uint64, cList []common.Address) (*staking.StakingInfo, map[common.Address]uint64, error) {
	sInfo, err := v.stakingInfo.GetStakingInfo(num)
	if err != nil {
		return nil, nil, err
	}
	stakingAmounts := make(map[common.Address]uint64, len(cList))
	for _, node := range cList {
		stakingAmounts[node] = uint64(0)
	}
	for _, consolidated := range sInfo.ConsolidatedNodes() {
		for _, nAddr := range consolidated.NodeIds {
			if _, ok := stakingAmounts[nAddr]; ok {
				stakingAmounts[nAddr] = consolidated.StakingAmount
			}
		}
	}
	return sInfo, stakingAmounts, nil
}

// convertHashToSeed takes the first 8 bytes (64 bits) and convert to int64
func convertHashToSeed(hash common.Hash) (int64, error) {
	hashstring := strings.TrimPrefix(hash.Hex(), "0x")
	if len(hashstring) > 15 {
		hashstring = hashstring[:15]
	}

	seed, err := strconv.ParseInt(hashstring, 16, 64)
	if err != nil {
		logger.Error("fail to make sub-list of validators", "hash", hash.Hex(), "seed", seed, "err", err)
		return 0, err
	}
	return seed, nil
}
