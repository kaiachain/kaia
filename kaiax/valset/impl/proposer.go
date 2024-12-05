package impl

import (
	"math"
	"math/big"
	"math/rand"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
	WeightedRandom
)

func (p ProposerPolicy) IsDefaultSet() bool {
	return p == RoundRobin || p == Sticky
}

func (p ProposerPolicy) IsWeightedRandom() bool {
	return p == WeightedRandom
}

func (v *ValsetModule) getLegacyProposersList(pUpdateBlock uint64) ([]common.Address, error) {
	if cachedProposers, ok := v.proposers.Get(pUpdateBlock); ok {
		if proposers, ok := cachedProposers.([]common.Address); ok {
			return proposers, nil
		}
		return nil, errInvalidProposersType
	}

	qualified, err := v.getQualifiedValidators(pUpdateBlock)
	if err != nil {
		return nil, err
	}

	pSet, _, err := v.getPSetWithProposerPolicy(pUpdateBlock)
	if err != nil {
		return nil, err
	}

	sInfo, stakingAmounts, err := v.getStakingInfoWithStakingAmounts(pUpdateBlock, qualified)
	if err != nil {
		return nil, err
	}

	prevHeader := v.chain.GetHeaderByNumber(pUpdateBlock - 1)
	if prevHeader == nil {
		return nil, errNilHeader
	}

	rules := v.chain.Config().Rules(big.NewInt(int64(pUpdateBlock)))
	proposersIndexes := calsSlotsInProposers(qualified, rules, pSet, sInfo, stakingAmounts)
	proposers := shuffleProposers(qualified, proposersIndexes, prevHeader.Hash())

	// store the calculated proposers
	v.proposers.Add(pUpdateBlock, proposers)

	// return the calculated proposers
	return proposers, nil
}

func pickRoundRobinProposer(cList valset.AddressList, policy ProposerPolicy, prevAuthor common.Address, round uint64) (common.Address, int) {
	if common.EmptyAddress(prevAuthor) {
		// if prevAuthor is empty, return -1 as an error
		return common.Address{}, -1
	}
	var (
		lastProposerIdx = cList.GetIdxByAddress(prevAuthor)
		seed            = round
	)

	if lastProposerIdx > -1 {
		seed += uint64(lastProposerIdx)
	}
	if policy == RoundRobin {
		seed += 1
	}

	idx := int(seed) % len(cList)
	return cList[idx], idx
}

func pickWeightedRandomProposer(proposers []common.Address, pUpdateBlock, num, round uint64, qualified valset.AddressList, author common.Address) (common.Address, int) {
	proposer := proposers[(int(num+round)-int(pUpdateBlock))%len(proposers)]
	idx := qualified.GetIdxByAddress(proposer)
	if idx != -1 {
		return proposer, idx
	}

	// fall-back to roundrobin proposer if proposer cannot found at qualified validators
	logger.Warn("Failed to select a new proposer, thus fall back to roundRobinProposer")
	return pickRoundRobinProposer(qualified, Sticky, author, round)
}

// CalcProposerBlockNumber returns number of block where list of proposers is updated for block blockNum
func calcProposerBlockNumber(blockNum uint64, proposerUpdateInterval uint64) uint64 {
	if blockNum == 0 {
		return 0
	}
	number := blockNum - (blockNum % proposerUpdateInterval)
	if blockNum%proposerUpdateInterval == 0 {
		number = blockNum - proposerUpdateInterval
	}
	return number
}

// proposersIndexes updates each validator's weight based on the ratio of its staking amount vs. the total staking amount.
func calsSlotsInProposers(qualified valset.AddressList, rules params.Rules, pSet gov.ParamSet, sInfo *staking.StakingInfo, stakingAmounts map[common.Address]uint64) []int {
	// is calculated among all CNs (i.e. AddressBook.cnStakingContracts)
	// stakingInfo.Gini calculates the gini among the qualified subset of the council (i.e. validators)
	gini := staking.EmptyGini
	if pSet.UseGiniCoeff {
		gini = sInfo.Gini(pSet.MinimumStake.Uint64())
	}

	// calc again for totalStaking amount among qualified subset of the council.
	totalStaking := float64(0)
	for _, st := range stakingAmounts {
		if st >= pSet.MinimumStake.Uint64() {
			stake := float64(st)
			if pSet.UseGiniCoeff {
				stake = math.Round(math.Pow(float64(st), 1.0/(1+gini)))
			}
			totalStaking += stake
		}
	}
	logger.Debug("calculate totalStaking", "UseGini", pSet.UseGiniCoeff, "Gini", gini, "totalStaking", totalStaking, "stakingAmounts", stakingAmounts)

	var (
		candidateValsIdx []int
		weights          = make([]uint64, len(qualified))
	)
	// weight is meaningful at next case. calculate the weight.
	if totalStaking != 0 && !rules.IsKore {
		for idx, addr := range qualified {
			weight := uint64(math.Round(float64(stakingAmounts[addr]) * 100 / totalStaking))
			if weight <= 0 {
				// A validator, who holds zero or small stake, has minimum weight, 1.
				weight = 1
			}
			weights[idx] = weight
		}
	} else {
		for idx := range qualified {
			weights[idx] = 1
		}
	}

	// allocate the validator slot per weight in proposers
	for index := 0; index < len(qualified); index++ {
		for i := uint64(0); i < weights[index]; i++ {
			candidateValsIdx = append(candidateValsIdx, index)
		}
	}

	return candidateValsIdx
}

func shuffleProposers(qualifiedVals valset.AddressList, candidateValsIdx []int, prevHash common.Hash) []common.Address {
	// shuffle it
	proposers := make([]common.Address, len(candidateValsIdx))
	seed, err := convertHashToSeed(prevHash)
	if err != nil {
		return nil
	}
	picker := rand.New(rand.NewSource(seed))
	for i := 0; i < len(candidateValsIdx); i++ {
		randIndex := picker.Intn(len(candidateValsIdx))
		candidateValsIdx[i], candidateValsIdx[randIndex] = candidateValsIdx[randIndex], candidateValsIdx[i]
	}

	// copy it
	for i := 0; i < len(candidateValsIdx); i++ {
		proposers[i] = qualifiedVals[candidateValsIdx[i]]
	}
	return proposers
}
