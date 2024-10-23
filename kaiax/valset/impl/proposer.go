package impl

import (
	"math"
	"math/rand"

	"github.com/kaiachain/kaia/common"
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

func (v *ValsetModule) getProposers(num uint64) ([]common.Address, error) {
	proposersUpdateBlock := params.CalcProposerBlockNumber(num)
	proposers, ok := v.proposers.Get(proposersUpdateBlock)
	if ok {
		return proposers.([]common.Address), nil
	}

	// there's no stored proposers. calculate the proposers and store at cache
	valCtx, err := newValSetContext(v, proposersUpdateBlock)
	if err != nil {
		return nil, err
	}
	council, err := newCouncil(v.ChainKv, valCtx)
	if err != nil {
		return nil, err
	}
	proposers = calculateProposers(council, valCtx)
	v.proposers.Add(num, proposers)

	return proposers.([]common.Address), nil
}

func calculateProposers(c *council, valCtx *valSetContext) []common.Address {
	// Although this is for selecting proposer, update it
	// otherwise, all parameters should be re-calculated at `RefreshProposers` method.
	var candidateValsIdx []int
	if !valCtx.rules.IsKore {
		weights := calcWeight(c.qualifiedValidators, valCtx)
		for index := range c.qualifiedValidators {
			for i := uint64(0); i < weights[index]; i++ {
				candidateValsIdx = append(candidateValsIdx, index)
			}
		}
	}

	// All validators has zero weight. Let's use all validators as candidate proposers.
	if len(candidateValsIdx) == 0 {
		for index := 0; index < len(c.qualifiedValidators); index++ {
			candidateValsIdx = append(candidateValsIdx, index)
		}
		logger.Trace("Refresh uses all validators as candidate proposers, because all weight is zero.", "candidateValsIdx", candidateValsIdx)
	}

	// shuffle it
	proposers := make([]common.Address, len(candidateValsIdx))
	seed, err := convertHashToSeed(valCtx.prevBlockResult.header.Hash())
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
		proposers[i] = c.qualifiedValidators[candidateValsIdx[i]]
	}
	return proposers
}

// calcWeight updates each validator's weight based on the ratio of its staking amount vs. the total staking amount.
func calcWeight(qualified subsetCouncilSlice, valCtx *valSetContext) []uint64 {
	var (
		sInfo                     = valCtx.prevBlockResult.staking
		pSet                      = valCtx.prevBlockResult.pSet
		consolidatedStakingAmount = valCtx.prevBlockResult.consolidatedStakingAmount()
	)
	// stakingInfo.Gini is calculated among all CNs (i.e. AddressBook.cnStakingContracts)
	// But we want the gini calculated among the subset of CNs (i.e. validators)
	totalStaking, gini := float64(0), float64(-1)
	if pSet.UseGiniCoeff {
		gini = sInfo.Gini(pSet.MinimumStake.Uint64())
		for _, st := range consolidatedStakingAmount {
			if st > pSet.MinimumStake.Uint64() {
				totalStaking += math.Round(math.Pow(float64(st), 1.0/(1+gini)))
			}
		}
	} else {
		for _, st := range consolidatedStakingAmount {
			if st > pSet.MinimumStake.Uint64() {
				totalStaking += float64(st)
			}
		}
	}
	logger.Debug("calculate totalStaking", "UseGini", pSet.UseGiniCoeff, "Gini", gini, "totalStaking", totalStaking, "stakingAmounts", consolidatedStakingAmount)

	// calculate and store each weight
	weights := make([]uint64, 0, len(qualified))
	if totalStaking == 0 {
		return weights
	}
	for idx, addr := range qualified {
		weight := uint64(math.Round(float64(consolidatedStakingAmount[addr]) * 100 / totalStaking))
		if weight <= 0 {
			// A validator, who holds zero or small stake, has minimum weight, 1.
			weight = 1
		}
		weights[idx] = weight
	}
	return weights
}
