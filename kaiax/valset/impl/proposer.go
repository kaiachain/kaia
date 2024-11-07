package impl

import (
	"math"
	"math/rand"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
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

func (v *ValsetModule) getProposers(num uint64) ([]common.Address, error) {
	valCtx, err := newValSetContext(v, num)
	proposersUpdateBlock := calcProposerBlockNumber(num, valCtx.prevBlockResult.pSet.ProposerUpdateInterval)
	proposers, ok := v.proposers.Get(proposersUpdateBlock)
	if ok {
		return proposers.([]common.Address), nil
	}

	// there's no stored proposers. calculate the proposers and store at cache
	valCtx, err = newValSetContext(v, proposersUpdateBlock)
	if err != nil {
		return nil, err
	}
	c, err := newCouncil(v.ChainKv, valCtx)
	if err != nil {
		return nil, err
	}

	proposersIndexes := calsSlotsInProposers(c.qualifiedValidators, valCtx)
	proposers = shuffleProposers(c.qualifiedValidators, proposersIndexes, valCtx.prevBlockResult.header.Hash())

	// store the calculated proposers
	v.proposers.Add(num, proposers)

	return proposers.([]common.Address), nil
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
func calsSlotsInProposers(qualified subsetCouncilSlice, valCtx *valSetContext) []int {
	var (
		sInfo                     = valCtx.prevBlockResult.staking
		pSet                      = valCtx.prevBlockResult.pSet
		consolidatedStakingAmount = valCtx.prevBlockResult.consolidatedStakingAmount()
		rules                     = valCtx.rules
	)

	// is calculated among all CNs (i.e. AddressBook.cnStakingContracts)
	// stakingInfo.Gini calculates the gini among the qualified subset of the council (i.e. validators)
	gini := staking.EmptyGini
	if pSet.UseGiniCoeff {
		gini = sInfo.Gini(pSet.MinimumStake.Uint64())
	}

	// calc again for totalStaking amount among qualified subset of the council.
	totalStaking := float64(0)
	for _, st := range consolidatedStakingAmount {
		if st > pSet.MinimumStake.Uint64() {
			stake := float64(st)
			if pSet.UseGiniCoeff {
				stake = math.Round(math.Pow(float64(st), 1.0/(1+gini)))
			}
			totalStaking += stake
		}
	}
	logger.Debug("calculate totalStaking", "UseGini", pSet.UseGiniCoeff, "Gini", gini, "totalStaking", totalStaking, "stakingAmounts", consolidatedStakingAmount)

	var (
		candidateValsIdx []int
		weights          = make([]uint64, len(qualified))
	)
	// weight is meaningful at next case. calculate the weight.
	if !rules.IsKore || totalStaking != 0 {
		for idx, addr := range qualified {
			weight := uint64(math.Round(float64(consolidatedStakingAmount[addr]) * 100 / totalStaking))
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

func shuffleProposers(qualifiedVals subsetCouncilSlice, candidateValsIdx []int, prevHash common.Hash) []common.Address {
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
