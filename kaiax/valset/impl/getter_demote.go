package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
)

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

// TODO-kaiax: move the feature into staking_info.go
func collectStakingAmounts(nodes []common.Address, si *staking.StakingInfo) map[common.Address]float64 {
	cns := si.ConsolidatedNodes()
	stakingAmounts := make(map[common.Address]float64, len(nodes))
	for _, node := range nodes {
		stakingAmounts[node] = 0
	}
	for _, cn := range cns {
		for _, node := range cn.NodeIds {
			if _, ok := stakingAmounts[node]; ok {
				stakingAmounts[node] = float64(cn.StakingAmount)
			}
		}
	}
	return stakingAmounts
}
