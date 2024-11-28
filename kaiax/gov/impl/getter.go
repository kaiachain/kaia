package impl

import (
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
