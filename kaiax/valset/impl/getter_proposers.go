package impl

import (
	"math"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
)

func (v *ValsetModule) getProposerList(c *blockContext) ([]common.Address, uint64, error) {
	var (
		useGini   = c.pset.UseGiniCoeff // because UseGiniCoeff is immutable, it's safe to pass ParamSet(num).UseGiniCoeff to calculate proposer list for updateNum.
		updateNum = roundDown(c.num-1, uint64(c.pset.ProposerUpdateInterval))
	)

	// calculate base pList at updateNum
	pList, err := v.calcBaseProposers(c, updateNum, useGini)
	if err != nil {
		return nil, 0, err
	}

	// apply the remove validator votes on the proposers
	removeVoteLists := v.getRemoveVotesInInterval(updateNum, c.pset.ProposerUpdateInterval)
	pList = removeVoteLists.filteredProposerList(c.num, pList)

	return pList, updateNum, nil
}

func (v *ValsetModule) calcBaseProposers(c *blockContext, updateNum uint64, useGini bool) ([]common.Address, error) {
	// Lookup proposers at updateNum from cache
	if list, ok := v.proposerListCache.Get(updateNum); ok {
		return list.([]common.Address), nil
	}

	// Generate here
	updateHeader := v.Chain.GetHeaderByNumber(updateNum)
	if updateHeader == nil {
		return nil, errNoHeader
	}
	var list []common.Address
	if c.rules.IsKore {
		list = generateProposerListUniform(c.qualified, updateHeader.Hash())
	} else {
		si, err := v.StakingModule.GetStakingInfo(updateNum)
		if err != nil {
			return nil, err
		}
		list = generateProposerListWeighted(c.qualified, si, useGini, updateHeader.Hash())
	}
	logger.Debug("GetProposerList", "number", c.num, "list", valset.NewAddressSet(list).String())

	// Store to cache
	v.proposerListCache.Add(updateNum, list)
	return list, nil
}

func (v *ValsetModule) getRemoveVotesInInterval(pUpdateNum, pUpdateInterval uint64) removeVoteList {
	// return early if removeVoteLists is found in cache
	if list, ok := v.removeVotesCache.Get(pUpdateNum); ok {
		return list.(removeVoteList)
	}

	// collect remove validator votes by scanning blocks between [pUpdateNum, pUpdateNum+pUpdateInterval)
	scanBlocks := v.scanBlocks(pUpdateNum, pUpdateInterval)
	removeVotes := make(removeVoteList, 0)
	for _, num := range scanBlocks {
		header := v.Chain.GetHeaderByNumber(num)
		if header == nil || len(header.Vote) == 0 {
			continue // skip nil header or vote
		}
		voteKey, addresses, ok := parseValidatorVote(header)
		if !ok {
			logger.Error("Failed to parse validator vote", "block", num)
			return nil
		}
		if len(addresses) > 0 && voteKey == gov.RemoveValidator {
			// demoted validators are not considered
			qualified, err := v.getQualifiedValidators(num)
			if err != nil {
				logger.Error("Failed to get qualified validators", "block", num, "err", err)
				return nil
			}
			if removeVotes.add(num, addresses, qualified) == false {
				logger.Error("remove vote is not added sequentially ascending order", "block", num)
				return nil
			}
		}
	}

	// Store to cache
	if len(removeVotes) > 0 {
		v.removeVotesCache.Add(pUpdateNum, removeVotes)
	}
	return removeVotes
}

// scanBlocks returns the block numbers to be scanned for remove validator votes
func (v *ValsetModule) scanBlocks(pUpdateNum, pUpdateInterval uint64) []uint64 {
	scanBlocks := make([]uint64, 0)

	// if migrated, scanBlocks is set as voteBlockNums between [pUpdateNum, pUpdateNum+pUpdateInterval)
	pMinVoteNum := v.readLowestScannedVoteNumCached()
	if pMinVoteNum != nil && *pMinVoteNum <= pUpdateNum {
		if v.validatorVoteBlockNumsCache != nil {
			scanBlocks = make([]uint64, len(v.validatorVoteBlockNumsCache))
			copy(scanBlocks, v.validatorVoteBlockNumsCache)
		} else {
			scanBlocks = ReadValidatorVoteBlockNums(v.ChainKv)
		}
		scanBlocks = slices.DeleteFunc(scanBlocks, func(n uint64) bool {
			return !(n >= pUpdateNum && n < pUpdateNum+pUpdateInterval)
		})
		return scanBlocks
	}

	// if not migrated, scanBlocks is set as blocknums between [updateNum, pUpdateNum+pUpdateInterval)
	for i := pUpdateNum; i < pUpdateNum+pUpdateInterval; i++ {
		scanBlocks = append(scanBlocks, i)
	}
	return scanBlocks
}

type removeVote struct {
	voteBlkNum  uint64
	removeAddrs []common.Address
}

// removeVoteList has a list of remove validator votes contained within one proposer update interval
// The block numbers in the list must be sorted in ascending order with no duplicates.
type removeVoteList []removeVote

// isUniquelySorted checks if the block numbers are sorted in ascending order and contain no duplicates.
func (r *removeVoteList) isUniquelySorted() bool {
	list := *r
	if len(list) <= 1 {
		return true
	}

	for i := 1; i < len(list); i++ {
		// Ensure no duplicate block numbers and maintain ascending order
		if list[i-1].voteBlkNum >= list[i].voteBlkNum {
			return false
		}
	}
	return true
}

// add appends a new removeVote to the list
// The new vote block number must be added sequentially in ascending order
func (r *removeVoteList) add(voteBlkNum uint64, addresses []common.Address, qualified *valset.AddressSet) bool {
	// last vote block number should be less than the new vote block number
	if len(*r) > 0 && (*r)[len(*r)-1].voteBlkNum >= voteBlkNum {
		return false
	}

	removeAddrs := valset.NewAddressSet([]common.Address{})
	for _, addr := range addresses {
		if qualified.Contains(addr) {
			removeAddrs.Add(addr)
		}
	}
	if removeAddrs.Len() > 0 {
		*r = append(*r, removeVote{voteBlkNum, removeAddrs.List()})
	}
	return true
}

// filteredProposerList removes the removeVoteAddrs from the base proposers
func (r *removeVoteList) filteredProposerList(currentNum uint64, proposers []common.Address) []common.Address {
	if !(*r).isUniquelySorted() {
		logger.Error("removeVoteList are not sorted. create the list in ascending order.")
		return nil
	}

	removeAddrs := new(valset.AddressSet)
	for _, vote := range *r {
		if currentNum <= vote.voteBlkNum {
			break
		}
		for _, addr := range vote.removeAddrs {
			removeAddrs.Add(addr)
		}
	}

	if removeAddrs.Len() == 0 {
		return proposers
	}

	newProposers := make([]common.Address, 0)
	copy(newProposers, proposers)
	for _, removeAddr := range proposers {
		if !removeAddrs.Contains(removeAddr) {
			newProposers = append(newProposers, removeAddr)
		}
	}
	return newProposers
}

func computeGini(amounts map[common.Address]float64) float64 {
	list := make([]float64, 0, len(amounts))
	for _, amount := range amounts {
		list = append(list, amount)
	}
	return staking.ComputeGini(list)
}

// generateProposerList generates a proposer list at a certain block (i.e. proposer update interval),
// whose qualified validators are `qualified`, staking amounts are in `si`, parameters are `pset` and its block hash is `blockHash`.
func generateProposerListWeighted(qualified *valset.AddressSet, si *staking.StakingInfo, useGini bool, blockHash common.Hash) []common.Address {
	var (
		addrs          = qualified.List()
		stakingAmounts = collectStakingAmounts(addrs, si)
		gini           = computeGini(stakingAmounts) // si.Gini is computed over every CN. But we want gini among validators, so we calculate here.
		exponent       = 1.0 / (1 + gini)
		totalStakes    = float64(0)
	)

	// Adjust staking amounts and calculate the sum
	if useGini {
		for addr, amount := range stakingAmounts {
			stakingAmounts[addr] = math.Round(math.Pow(float64(amount), exponent))
			totalStakes += stakingAmounts[addr]
		}
	} else {
		for _, amount := range stakingAmounts {
			totalStakes += amount
		}
	}

	// Calculate percentile weights
	weights := make(map[common.Address]uint64)
	if totalStakes > 0 {
		for _, addr := range addrs {
			weight := uint64(math.Round(stakingAmounts[addr] * 100 / totalStakes))
			if weight <= 0 {
				weight = 1
			}
			weights[addr] = weight
		}
	} else {
		for _, addr := range addrs {
			weights[addr] = 0
		}
	}

	// Generate weighted repeated list
	proposerList := make([]common.Address, 0)
	for _, addr := range addrs {
		for i := uint64(0); i < weights[addr]; i++ {
			proposerList = append(proposerList, addr)
		}
	}
	// If the list is empty (i.e. all weights are zero), list each validator once.
	if len(proposerList) == 0 {
		for _, addr := range addrs {
			proposerList = append(proposerList, addr)
		}
	}

	seed := valset.HashToSeedLegacy(blockHash)
	return valset.NewAddressSet(proposerList).ShuffledListLegacy(seed)
}

func generateProposerListUniform(qualified *valset.AddressSet, blockHash common.Hash) []common.Address {
	proposerList := qualified.List()
	seed := valset.HashToSeedLegacy(blockHash)
	return valset.NewAddressSet(proposerList).ShuffledListLegacy(seed)
}
