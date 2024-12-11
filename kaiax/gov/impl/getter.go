package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/kaiax/gov"
)

func (m *GovModule) EffectiveParamSet(blockNum uint64) gov.ParamSet {
	ret := gov.GetDefaultGovernanceParamSet()

	p0 := m.Fallback
	for k, v := range p0 {
		ret.Set(k, v)
	}

	p1 := m.Hgm.EffectiveParamsPartial(blockNum)
	for k, v := range p1 {
		ret.Set(k, v)
	}

	if m.isKoreHF(blockNum) {
		p2 := m.Cgm.EffectiveParamsPartial(blockNum)
		for k, v := range p2 {
			ret.Set(k, v)
		}
	}

	return *ret
}

func (m *GovModule) distill(gp gov.ParamSet, blockNum uint64) gov.ParamSet {
	rule := m.Chain.Config().Rules(new(big.Int).SetUint64(blockNum))

	// To avoid confusion, override some parameters that are deprecated after hardforks.
	// e.g., stakingupdateinterval is shown as 86400 but actually irrelevant (i.e. updated every block)
	if rule.IsKore {
		// Gini option deprecated since Kore, as All committee members have an equal chance
		// of being elected block proposers.
		gp.UseGiniCoeff = false
	}
	if rule.IsRandao {
		// Block proposer is randomly elected at every block with Randao,
		// no more precalculated proposer list.
		gp.ProposerUpdateInterval = 1
	}
	if rule.IsKaia {
		// Staking information updated every block since Kaia.
		gp.StakingUpdateInterval = 1
	}
	return gp
}
