package impl

import (
	"fmt"
	"math/big"
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

// GetCouncil returns the whole validator list of block N.
// If this network haven't voted since genesis, return genesis council which is stored at Block 0.
func (v *ValsetModule) GetCouncil(num uint64) (valset.AddressList, error) {
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
func (v *ValsetModule) GetCommittee(num uint64, round uint64) (valset.AddressList, error) {
	// if the block number is genesis, directly return genesisCouncil as committee.
	if num == 0 {
		genesisCouncil, err := v.GetCouncil(0)
		if err != nil {
			return nil, err
		}
		return genesisCouncil, nil
	}

	committeeCtx, err := newCommitteeContext(v, num)
	if err != nil {
		return nil, err
	}

	return committeeCtx.getCommittee(round)
}

// GetProposer calculates a proposer for the round of the given block.
// Note that, to skip the calculation, derive the author from header if the block exists
// - if the policy defaultSet, the proposer is picked from the council. it's defaultSet, so there's no demoted.
// - if the policy is weightedrandom and before Randao hf, the proposer is picked from the proposers
// - if the pilicy is weightedrandom and after Randao hf, the proposer is picked from the committee
func (v *ValsetModule) GetProposer(num uint64, round uint64) (common.Address, error) {
	// if the block exists, derive the author from header
	header := v.chain.GetHeaderByNumber(num)
	if header != nil && header.Number.Uint64() == num && uint64(header.Round()) == round {
		return v.chain.Engine().Author(header)
	}

	// if the block number is genesis, directly return proposer as the first element of genesisCouncil.
	if num == 0 {
		genesisCouncil, err := v.GetCouncil(0)
		if err != nil {
			return common.Address{}, err
		}
		return genesisCouncil[0], nil
	}

	committeeCtx, err := newCommitteeContext(v, num)
	if err != nil {
		return common.Address{}, err
	}

	return committeeCtx.getProposer(round)
}

// getQualifiedValidators returns a list of validators who are qualified to be a member of the committee or proposer
func (v *ValsetModule) getQualifiedValidators(num uint64) (valset.AddressList, error) {
	council, err := v.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	demoted, err := v.getDemotedValidators(num)
	if err != nil {
		return nil, err
	}

	qualified := make(valset.AddressList, 0)
	for _, v := range council {
		if demoted.GetIdxByAddress(v) == -1 { // if cannot find in demoted,
			qualified = append(qualified, v) // add it to qualified
		}
	}

	sort.Sort(qualified)
	return qualified, nil
}

// getDemotedValidators are subtract of qualified from council(N)
func (v *ValsetModule) getDemotedValidators(num uint64) (valset.AddressList, error) {
	if num == 0 {
		return valset.AddressList{}, nil
	}

	pSet, proposerPolicy, err := v.getPSetWithProposerPolicy(num)
	if err != nil {
		return nil, err
	}

	var (
		isSingleMode = pSet.GovernanceMode == "single"
		govNode      = pSet.GoverningNode
		minStaking   = pSet.MinimumStake.Uint64()
		rules        = v.chain.Config().Rules(big.NewInt(int64(num)))
	)

	// either proposer-policy is not weighted random or before istanbul HF,
	// do not filter out the demoted validators
	if !proposerPolicy.IsWeightedRandom() || !rules.IsIstanbul {
		return valset.AddressList{}, nil
	}

	council, err := v.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	// Split the council(N) into qualified and demoted.
	// Qualified is a subset of the council who are qualified to be a committee member.
	// (1) If governance mode is single, always include the governing node.
	// (2) If no council members has enough KAIA, all members become qualified.
	_, stakingAmounts, err := v.getStakingInfoWithStakingAmounts(num, council)
	if err != nil {
		return nil, err
	}

	demoted := make(valset.AddressList, 0)
	for _, addr := range council {
		staking, ok := stakingAmounts[addr]
		if !ok {
			return nil, fmt.Errorf("cannot find staking amount for %s", addr)
		}
		if staking < minStaking && !(isSingleMode && addr == govNode) {
			demoted = append(demoted, addr)
		}
	}

	// include all council members if case1 or case2
	//   case1. not a single mode && no qualified
	//   case2. single mode && len(qualified) is 1 && govNode is not qualified
	if len(demoted) == len(council) || (isSingleMode && len(demoted) == len(council)-1 && stakingAmounts[govNode] < minStaking) {
		return valset.AddressList{}, nil
	}

	sort.Sort(demoted)
	return demoted, nil
}

// getPSetWithProposerPolicy returns the govParam & proposer policy after processing block num - 1
func (v *ValsetModule) getPSetWithProposerPolicy(num uint64) (gov.ParamSet, ProposerPolicy, error) {
	pSet := v.governance.EffectiveParamSet(num)
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
