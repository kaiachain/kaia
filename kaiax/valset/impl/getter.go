package impl

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/params"
)

type blockResult struct {
	proposerPolicy ProposerPolicy // prevBlockResult.pSet.proposerPolicy
	staking        *staking.StakingInfo
	header         *types.Header
	author         common.Address
	pSet           gov.ParamSet
}

// consolidatedStakingAmounts get total staking amounts per staking contracts by nodeIds
func (br *blockResult) consolidatedStakingAmount() map[common.Address]uint64 {
	consolidatedStakingAmounts := make(map[common.Address]uint64, len(br.staking.NodeIds))
	for idx, nAddr := range br.staking.NodeIds {
		consolidatedStakingAmounts[nAddr] = br.staking.ConsolidatedNodes()[idx].StakingAmount
	}
	return consolidatedStakingAmounts
}

type valSetContext struct {
	blockNumber     uint64 // id of valSetContext
	rules           params.Rules
	prevBlockResult *blockResult // previous block's results - states
}

func newValSetContext(v *ValsetModule, bn uint64) (*valSetContext, error) {
	if bn == 0 {
		bn = bn + 1
	}
	// 1. get previous block's results
	sInfo, err := v.stakingInfo.GetStakingInfo(bn - 1) // stakingInfo after processing block N-1
	if err != nil {
		return nil, err
	}
	header := v.chain.GetHeaderByNumber(bn - 1)
	if header == nil {
		return nil, errNilHeader
	}
	author, err := v.chain.Engine().Author(header) // author of block N-1
	if err != nil {
		return nil, err
	}
	pSet := v.headerGov.EffectiveParamSet(bn - 1) // govParam after processing block N-1
	if pSet.CommitteeSize == 0 {
		return nil, errInvalidCommitteeSize // it cannot happen. just to make sure
	}
	// if config.Istanbul is nil, it means the consensus is not 'istanbul' so use defualt proposer policy( = RoundRobin).
	proposerPolicy := ProposerPolicy(params.DefaultProposerPolicy)
	if v.chain.Config().Istanbul != nil {
		proposerPolicy = ProposerPolicy(pSet.ProposerPolicy)
	}
	prevBlockResult := &blockResult{proposerPolicy, sInfo, header, author, pSet}

	// 2. get network config - Rules
	rules := v.chain.Config().Rules(big.NewInt(int64(bn)))
	return &valSetContext{bn, rules, prevBlockResult}, nil
}

// GetCouncilAddressList returns the whole validator list of block N.
// If this network haven't voted since genesis, return genesis council which is stored at Block 0.
func (v *ValsetModule) GetCouncilAddressList(num uint64) ([]common.Address, error) {
	// make sure that the num doesn't exceed current block
	councilAddresses, err := ReadCouncilAddressListFromDb(v.ChainKv, num)
	if err != nil {
		return nil, err
	}
	return councilAddresses, nil
}

// GetCommitteeAddressList returns the current round or block's committee.
func (v *ValsetModule) GetCommitteeAddressList(num uint64, round uint64) ([]common.Address, error) {
	// if the block number is genesis, directly return councilAddressList as committee.
	if num == 0 {
		councilAddressList, err := ReadCouncilAddressListFromDb(v.ChainKv, 0)
		if err != nil {
			return nil, err
		}
		committeeSize := v.headerGov.EffectiveParamSet(0).CommitteeSize
		return councilAddressList[:committeeSize], nil
	}

	valCtx, err := newValSetContext(v, num)
	if err != nil {
		return nil, err
	}

	c, err := newCouncil(v.ChainKv, valCtx)
	if err != nil {
		return nil, err
	}

	var (
		proposerPolicy = valCtx.prevBlockResult.proposerPolicy
		rules          = valCtx.rules
	)
	if valCtx.prevBlockResult.pSet.CommitteeSize >= uint64(len(c.qualifiedValidators)) {
		return c.qualifiedValidators, nil
	}
	if !proposerPolicy.IsWeightedRandom() || !rules.IsRandao {
		pUpdateBlock := calcProposerBlockNumber(num, valCtx.prevBlockResult.pSet.ProposerUpdateInterval)
		proposers, err := v.getProposers(pUpdateBlock)
		if err != nil {
			return nil, err
		}
		return c.selectRandomCommittee(valCtx, round, proposers)
	}
	return c.selectRandaoCommittee(valCtx.prevBlockResult)
}

// GetProposer calculates a proposer for the (N, Round) view.
// There's two way to derive the proposer:
// - if the view (N, Round) is same as existing header, derive it from header.
// - otherwise, calc it.
//   - if the policy defaultSet, the proposer is picked from the council. it's defaultSet, so there's no demoted.
//   - if the policy is weightedrandom and before Randao hf, the proposer is picked from the proposers
//   - if the pilicy is weightedrandom and after Randao hf, the proposer is picked from the committee
func (v *ValsetModule) GetProposer(num uint64, round uint64) (common.Address, error) {
	header := v.chain.GetHeaderByNumber(num)
	if header.Number.Uint64() == num && uint64(header.Round()) == round {
		return v.chain.Engine().Author(header)
	}

	// if the block number is genesis, directly return proposer as the first element of councilAddressList.
	if num == 0 {
		councilAddressList, err := ReadCouncilAddressListFromDb(v.ChainKv, 0)
		if err != nil {
			return common.Address{}, err
		}
		return councilAddressList[0], nil
	}

	valCtx, err := newValSetContext(v, num)
	if err != nil {
		return common.Address{}, err
	}
	c, err := newCouncil(v.ChainKv, valCtx)
	if err != nil {
		return common.Address{}, err
	}

	var (
		proposerPolicy = valCtx.prevBlockResult.proposerPolicy
		author         = valCtx.prevBlockResult.author
	)

	if proposerPolicy.IsDefaultSet() {
		proposer, _ := pickRoundRobinProposer(c.councilAddressList, proposerPolicy, author, round)
		return proposer, nil
	}

	// before Randao, weightedrandom uses proposers to pick the proposer.
	if !valCtx.rules.IsRandao {
		pUpdateBlock := calcProposerBlockNumber(num, valCtx.prevBlockResult.pSet.ProposerUpdateInterval)
		proposers, err := v.getProposers(pUpdateBlock)
		if err != nil {
			return common.Address{}, err
		}
		proposer, _ := pickWeightedRandomProposer(proposers, pUpdateBlock, num, round, c.qualifiedValidators, author)
		return proposer, nil
	}

	// after Randao, pick proposer from randao committee.
	committee, err := c.selectRandaoCommittee(valCtx.prevBlockResult)
	if err != nil {
		return common.Address{}, err
	}
	return committee[int(round)%len(c.qualifiedValidators)], nil
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
