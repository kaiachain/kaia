package impl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/reward"
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

type subsetCouncilSlice []common.Address

func (sc subsetCouncilSlice) Len() int {
	return len(sc)
}

func (sc subsetCouncilSlice) Less(i, j int) bool {
	return strings.Compare(sc[i].String(), sc[j].String()) < 0
}

func (sc subsetCouncilSlice) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

func (sc subsetCouncilSlice) AddressStringList() []string {
	var stringAddrs []string
	for _, val := range sc {
		stringAddrs = append(stringAddrs, val.String())
	}
	return stringAddrs
}

func (sc subsetCouncilSlice) getIdxByAddress(addr common.Address) int {
	for i, val := range sc {
		if addr == val {
			return i
		}
	}
	// TODO-Kaia-Istanbul: Enable this log when non-committee nodes don't call `core.startNewRound()`
	// logger.Warn("failed to find an address in the validator list",
	// 	"address", addr, "validatorAddrs", valSet.validators.AddressStringList())
	return -1
}

// SortedAddressList retrieves the sorted address list of ValidatorSet in "ascending order".
// if public is false, sort it using bytes.Compare. It's for public purpose.
// - public-false usage: (getValidators/getDemotedValidators, defaultSet snap store, prepareExtra.validators)
// if public is true, sort it using strings.Compare. It's used for internal consensus purpose, especially for the source of committee.
// - public-true usage: (snap read/store/apply except defaultSet snap store, vrank log)
// TODO-kaia-valset: unify sorting.
func (sc subsetCouncilSlice) sortedAddressList(public bool) []common.Address {
	copiedSlice := make(subsetCouncilSlice, len(sc))
	copy(copiedSlice, sc)

	if public {
		// want reverse-sort: ascending order - bytes.Compare(ValidatorSet[i][:], ValidatorSet[j][:]) > 0
		sort.Slice(copiedSlice, func(i, j int) bool {
			return bytes.Compare(copiedSlice[i].Bytes(), copiedSlice[j].Bytes()) >= 0
		})
		sort.Sort(sort.Reverse(copiedSlice))
	} else {
		// want sort: descending order - strings.Compare(ValidatorSet[i].String(), ValidatorSet[j].String()) < 0
		sort.Sort(copiedSlice)
	}
	return copiedSlice
}

type Council struct {
	blockNumber    uint64
	round          uint64
	rules          params.Rules
	proposerPolicy ProposerPolicy // prevBlockResult.pSet.proposerPolicy

	// To calculate committee(num), we need council,prevHash,stakingInfo of lastProposal/prevBlock
	// which blocknumber is num - 1.
	prevBlockResult *blockResult
	qualified       subsetCouncilSlice // qualified is a subset of prev block's council
	demoted         subsetCouncilSlice // demoted is a subset of prev block's council which doesn't fulfill the minimum staking amount

	// latest proposer update block's information for calculating the current block's proposers, however it is deprecated since kaia KF
	// if Council.UseProposers is false, do not use and do not calculate the proposers. see the condition at Council.UseProposers method
	// if it uses cached proposers, do not calculate the proposers
	proposers []common.Address
}

// NewCouncil returns preprocessed council(N-1). It is useful to calculate the committee(N,R) or proposer(N,R).
// defaultSet - do nothing, weightedrandom - filter out by minimum staking amount after istanbul HF and sort it
func (v *ValsetModule) NewCouncil(blockNumber uint64) (*Council, error) {
	prevBlockResult, err := v.getBlockResultsByNumber(blockNumber - 1)
	if err != nil {
		return nil, err
	}

	council := &Council{
		blockNumber:     blockNumber,
		rules:           v.chain.Config().Rules(big.NewInt(int64(blockNumber))),
		prevBlockResult: prevBlockResult,
	}
	// if config.Istanbul is nil, it means the consensus is not 'istanbul'.
	// use default proposer policy (= RoundRobin). If not, use paramSet.proposerPolicy.
	if v.chain.Config().Istanbul == nil {
		council.proposerPolicy = ProposerPolicy(params.DefaultProposerPolicy)
	} else {
		council.proposerPolicy = ProposerPolicy(prevBlockResult.pSet.ProposerPolicy)
	}

	// defaultSet does not filter out demoted validators or calculate proposers since it is PoA
	if council.proposerPolicy.IsDefaultSet() {
		council.qualified = make([]common.Address, len(council.prevBlockResult.councilAddrList))
		copy(council.qualified, council.prevBlockResult.councilAddrList)
		return council, nil
	}

	// weighted random filter out under-staked nodes since istanbul HF
	if council.rules.IsIstanbul {
		council.qualified, council.demoted = splitByMinimumStakingAmount(council.prevBlockResult)
	}

	// latest proposer update block's information for calculating the current block's proposers, however it is deprecated since kaia KF
	// in some cases, skip the proposers calculation
	//   case1. already cached at proposers
	//   case2. after Randao HF, proposers is deprecated
	//   case3. if proposer policy is not weighted random, proposers is not used
	if council.UseProposers() {
		proposerUpdateBlock := params.CalcProposerBlockNumber(blockNumber)
		cachedProposers, ok := v.proposers.Get(proposerUpdateBlock)
		if !ok {
			proposerUpdateBlockRules := v.chain.Config().Rules(big.NewInt(int64(proposerUpdateBlock)))
			proposerUpdatePrevBlockResult, err := v.getBlockResultsByNumber(proposerUpdateBlock - 1)
			if err != nil {
				return nil, err
			}
			proposerUpdateBlockQualified, _ := splitByMinimumStakingAmount(proposerUpdatePrevBlockResult)
			council.proposers = calculateProposers(proposerUpdateBlockQualified, proposerUpdatePrevBlockResult, proposerUpdateBlockRules)
			v.proposers.Add(blockNumber, council.proposers)
		} else {
			council.proposers = cachedProposers.([]common.Address)
		}
	}

	return council, nil
}

// calcWeight updates each validator's weight based on the ratio of its staking amount vs. the total staking amount.
func calcWeight(qualified subsetCouncilSlice, prevBnResult *blockResult) []uint64 {
	var (
		sInfo                     = prevBnResult.staking
		pSet                      = prevBnResult.pSet
		consolidatedStakingAmount = prevBnResult.consolidatedStakingAmount()
	)
	// stakingInfo.Gini is calculated among all CNs (i.e. AddressBook.cnStakingContracts)
	// But we want the gini calculated among the subset of CNs (i.e. validators)
	totalStaking, gini := float64(0), reward.DefaultGiniCoefficient
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

func calculateProposers(qualified subsetCouncilSlice, prevBnResult *blockResult, rules params.Rules) []common.Address {
	// Although this is for selecting proposer, update it
	// otherwise, all parameters should be re-calculated at `RefreshProposers` method.
	var candidateValsIdx []int
	if !rules.IsKore {
		weights := calcWeight(qualified, prevBnResult)
		for index := range qualified {
			for i := uint64(0); i < weights[index]; i++ {
				candidateValsIdx = append(candidateValsIdx, index)
			}
		}
	}

	// All validators has zero weight. Let's use all validators as candidate proposers.
	if len(candidateValsIdx) == 0 {
		for index := 0; index < len(qualified); index++ {
			candidateValsIdx = append(candidateValsIdx, index)
		}
		logger.Trace("Refresh uses all validators as candidate proposers, because all weight is zero.", "candidateValsIdx", candidateValsIdx)
	}

	// shuffle it
	proposers := make([]common.Address, len(candidateValsIdx))
	seed, err := convertHashToSeed(prevBnResult.header.Hash())
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
		proposers[i] = qualified[candidateValsIdx[i]]
	}
	return proposers
}

func (c Council) UseProposers() bool {
	return c.proposerPolicy.IsWeightedRandom() && !c.rules.IsRandao
}

// proposer picks a proposer for the (N, Round) view. it is picked from different sources.
// - if the policy defaultSet, the proposer is picked from the council. it's defaultSet, so there's no demoted.
// - if the policy is weightedrandom and before Randao hf, the proposer is picked from the proposers
// - if the pilicy is weightedrandom and after Randao hf, the proposer is picked from the committee
func (c Council) proposer(round uint64) (common.Address, int) {
	if c.proposerPolicy.IsDefaultSet() {
		lastProposerIdx := c.qualified.getIdxByAddress(c.prevBlockResult.author)
		seed := defaultSetNextProposerSeed(c.proposerPolicy, c.prevBlockResult.author, lastProposerIdx, round)
		idx := int(seed) % len(c.qualified)
		return c.qualified[idx], idx
	}

	var (
		proposer    common.Address
		proposerIdx int
	)

	// before Randao, weightedrandom uses proposers to pick the proposer.
	// the proposers have always been calculated before entering proposer
	if !c.rules.IsRandao {
		proposer = c.proposers[int(c.blockNumber+round)%3600%len(c.proposers)]
		proposerIdx = c.qualified.getIdxByAddress(proposer)
	}

	// after Randao, proposers is deprecated. pick proposer from qualified council member list.
	committee, err := c.selectRandaoCommittee()
	if err != nil {
		return common.Address{}, -1
	}

	// for non-default, fall-back to roundrobin proposer
	idx := int(round) % len(c.qualified)
	proposer, proposerIdx = committee[idx], idx
	if proposerIdx == -1 {
		logger.Warn("Failed to select a new proposer, thus fall back to roundRobinProposer")
		idx = c.qualified.getIdxByAddress(proposer)
		seed := defaultSetNextProposerSeed(params.Sticky, c.prevBlockResult.author, idx, round)
		proposerIdx = int(seed) % len(c.qualified)
		proposer = c.qualified[proposerIdx]
	}
	return proposer, proposerIdx
}

func defaultSetNextProposerSeed(policy ProposerPolicy, proposer common.Address, proposerIdx int, round uint64) uint64 {
	seed := round
	if proposerIdx > -1 {
		seed += uint64(proposerIdx)
	}
	if policy == params.RoundRobin && !common.EmptyAddress(proposer) {
		seed += 1
	}
	return seed
}

func (c Council) selectCommittee(round uint64) (subsetCouncilSlice, error) {
	var (
		committeeSize = c.prevBlockResult.pSet.CommitteeSize
		policy        = c.proposerPolicy
		rules         = c.rules
		validators    = c.qualified
	)
	if committeeSize >= uint64(len(validators)) {
		committee := make([]common.Address, len(validators))
		copy(committee, validators)
		return committee, nil
	}

	// it cannot be happened. just to make sure
	if committeeSize == 0 {
		return nil, errors.New("invalid committee size")
	}

	if policy.IsDefaultSet() || (policy.IsWeightedRandom() && !rules.IsRandao) {
		return c.selectRandomCommittee(round)
	}
	return c.selectRandaoCommittee()
}

// selectRandomCommittee composes a committee selecting validators randomly based on the seed value.
// It returns nil if the given committeeSize is bigger than validatorSize or proposer indexes are invalid.
func (c Council) selectRandomCommittee(round uint64) ([]common.Address, error) {
	var (
		validatorSize                                         = len(c.qualified)
		validator                                             = c.qualified
		committeeSize                                         = c.prevBlockResult.pSet.CommitteeSize
		proposer, proposerIdx                                 = c.proposer(round)
		closestDifferentProposer, closestDifferentProposerIdx = c.proposer(round + 1)
	)

	// return early if the committee size is 1
	if committeeSize == 1 {
		return []common.Address{proposer}, nil
	}

	// closest next proposer who has different address with the proposer
	for i := uint64(1); i < params.ProposerUpdateInterval(); i++ {
		if proposer != closestDifferentProposer {
			break
		}
		closestDifferentProposer, closestDifferentProposerIdx = c.proposer(round + i)
	}

	// ensure validator indexes are valid
	if proposerIdx < 0 || closestDifferentProposerIdx < 0 || proposerIdx == closestDifferentProposerIdx ||
		validatorSize <= proposerIdx || validatorSize <= closestDifferentProposerIdx {
		return nil, fmt.Errorf("invalid indexes of validators. validatorSize: %d, proposerIdx:%d, nextProposerIdx:%d",
			validatorSize, proposerIdx, closestDifferentProposerIdx)
	}

	seed, err := convertHashToSeed(c.prevBlockResult.header.Hash())
	if err != nil {
		return nil, err
	}
	// shuffle the qualified validators except two proposers
	if c.rules.IsIstanbul {
		seed += int64(round)
	}
	committee := make([]common.Address, committeeSize)
	picker := rand.New(rand.NewSource(seed))
	pickSize := validatorSize - 2
	indexs := make([]int, pickSize)
	idx := 0
	for i := 0; i < validatorSize; i++ {
		if i != proposerIdx && i != closestDifferentProposerIdx {
			indexs[idx] = i
			idx++
		}
	}
	for i := 0; i < pickSize; i++ {
		randIndex := picker.Intn(pickSize)
		indexs[i], indexs[randIndex] = indexs[randIndex], indexs[i]
	}

	// first committee is the proposer and the second committee is the next proposer
	committee[0], committee[1] = proposer, closestDifferentProposer
	for i := uint64(0); i < committeeSize-2; i++ {
		committee[i+2] = validator[indexs[i]]
	}

	return committee, nil
}

// SelectRandaoCommittee composes a committee selecting validators randomly based on the mixHash.
// It is guaranteed that the committee include proposers for all rounds because
// the proposer is picked from the this committee. See weightedRandomProposer().
//
// def select_committee_KIP146(validators, committee_size, seed):
//
//	shuffled = shuffle_validators_KIP146(validators, seed)
//	return shuffled[:min(committee_size, len(validators))]
func (c Council) selectRandaoCommittee() ([]common.Address, error) {
	if c.prevBlockResult.header.MixHash == nil {
		return nil, fmt.Errorf("nil mixHash")
	}

	copied := make(subsetCouncilSlice, len(c.qualified))
	copy(copied, c.qualified)

	seed := int64(binary.BigEndian.Uint64(c.prevBlockResult.header.MixHash[:8]))
	rand.New(rand.NewSource(seed)).Shuffle(len(copied), copied.Swap)
	return copied[:c.prevBlockResult.pSet.CommitteeSize], nil
}

// splitByMinimumStakingAmount split the council members into qualified, demoted.
// Qualified is a subset of the council who have staked more than minimum staking amount. Demoted stakes less than minimum.
// There's two rules.
// (1) If governance mode is single, always include the governing node.
// (2) If no council members has enough KAIA, all members become qualified.
func splitByMinimumStakingAmount(prevBlockResult *blockResult) (subsetCouncilSlice, subsetCouncilSlice) {
	var (
		qualified subsetCouncilSlice
		demoted   subsetCouncilSlice

		// get params
		isSingleMode  = prevBlockResult.pSet.GovernanceMode == "single"
		govNode       = prevBlockResult.pSet.GoverningNode
		minStaking    = prevBlockResult.pSet.MinimumStake.Uint64()
		stakingAmount = prevBlockResult.consolidatedStakingAmount()
	)

	// filter out the demoted members who have staked less than minimum staking amounts
	for addr, val := range stakingAmount {
		if val >= minStaking || (isSingleMode && addr == govNode) {
			qualified = append(qualified, addr)
		} else {
			demoted = append(qualified, addr)
		}
	}

	// include all council members if case1 or case2
	//   case1. not a single mode && no qualified
	//   case2. single mode && len(qualified) is 1 && govNode is not qualified
	if len(qualified) == 0 || (isSingleMode && len(qualified) == 1 && stakingAmount[govNode] < minStaking) {
		demoted = subsetCouncilSlice{} // ensure demoted is empty
		qualified = make(subsetCouncilSlice, len(prevBlockResult.councilAddrList))
		copy(qualified, prevBlockResult.councilAddrList)
	}

	sort.Sort(qualified)
	sort.Sort(demoted)
	return qualified, demoted
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
