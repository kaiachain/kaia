package impl

import (
	"math/big"
	"sort"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
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

func (v *ValsetModule) getCouncil(num uint64) (*valset.AddressSet, error) {
	if num == 0 {
		return getCouncilGenesis(v.Chain.GetHeaderByNumber(0))
	}

	pBorder := ReadLowestScannedSnapshotNum(v.ChainKv)
	if pBorder == nil || *pBorder > 0 { // migration not started or migration not completed.
		council, _, err := v.replayFromIstanbulSnapshot(num, false)
		return council, err
	} else {
		return v.getCouncilDB(num)
	}
}

// getCouncilGenesis parses the genesis council from the header's extraData.
func getCouncilGenesis(header *types.Header) (*valset.AddressSet, error) {
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	return valset.NewAddressSet(istanbulExtra.Validators), nil
}

func (v *ValsetModule) getCouncilDB(num uint64) (*valset.AddressSet, error) {
	if v.validatorVoteBlockNumsCache == nil {
		v.validatorVoteBlockNumsCache = ReadValidatorVoteBlockNums(v.ChainKv)
	}
	nums := v.validatorVoteBlockNumsCache
	if nums == nil {
		return nil, errEmptyVoteBlock
	}
	voteNum := lastVoteBlockNum(nums, num)
	council := valset.NewAddressSet(ReadCouncil(v.ChainKv, voteNum))
	return council, nil
}

// lastVoteBlockNum returns the last block number in the list that is less than the given block number.
// For instance, if nums = [0, 10, 20, 30] and num = 25, the result is 20.
func lastVoteBlockNum(nums []uint64, num uint64) uint64 {
	// idx is the smallest index that is greater than or equal to `num`.
	// idx-1 is the largest index that is less than `num`.
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= num
	})
	if idx > 0 && nums[idx-1] < num {
		return nums[idx-1]
	} else {
		return 0
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

// getDemotedValidators returns the demoted validators at the given block number.
func (v *ValsetModule) getDemotedValidators(council *valset.AddressSet, num uint64) (*valset.AddressSet, error) {
	if num == 0 {
		return valset.NewAddressSet(nil), nil
	}

	pset := v.GovModule.EffectiveParamSet(num)
	rules := v.Chain.Config().Rules(new(big.Int).SetUint64(num))

	switch istanbul.ProposerPolicy(pset.ProposerPolicy) {
	case istanbul.RoundRobin, istanbul.Sticky:
		// All council members are qualified for both RoundRobin and Sticky.
		return valset.NewAddressSet(nil), nil
	case istanbul.WeightedRandom:
		// All council members are qualified for WeightedRandom before Istanbul hardfork.
		if !rules.IsIstanbul {
			return valset.NewAddressSet(nil), nil
		}
		// Otherwise, filter out based on staking amounts.
		si, err := v.StakingModule.GetStakingInfo(num)
		if err != nil {
			return nil, err
		}
		return filterValidatorsIstanbul(council, si, pset), nil
	default:
		return nil, errInvalidProposerPolicy
	}
}

func filterValidatorsIstanbul(council *valset.AddressSet, si *staking.StakingInfo, pset gov.ParamSet) *valset.AddressSet {
	var (
		demoted        = valset.NewAddressSet(nil)
		singleMode     = pset.GovernanceMode == "single"
		governingNode  = pset.GoverningNode
		minStake       = pset.MinimumStake.Uint64() // in KAIA
		stakingAmounts = collectStakingAmounts(council.List(), si)
	)

	// First filter by staking amounts.
	for _, node := range council.List() {
		if uint64(stakingAmounts[node]) < minStake {
			demoted.Add(node)
		}
	}

	// If all validators are demoted, then no one is demoted.
	if demoted.Len() == len(council.List()) {
		demoted = valset.NewAddressSet(nil)
	}

	// Under single governnace mode, governing node cannot be demoted.
	if singleMode && demoted.Contains(governingNode) {
		demoted.Remove(governingNode)
	}
	return demoted
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
