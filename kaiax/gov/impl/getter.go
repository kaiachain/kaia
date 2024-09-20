package impl

import "github.com/kaiachain/kaia/kaiax/gov"

func (m *GovModule) EffectiveParamSet(blockNum uint64) (gov.ParamSet, error) {
	ret := gov.GetDefaultGovernanceParamSet()

	p1, err := m.hgm.EffectiveParamsPartial(blockNum)
	if err != nil {
		return gov.ParamSet{}, err
	}
	for k, v := range p1 {
		ret.Set(k, v)
	}

	p2, err := m.cgm.EffectiveParamsPartial(blockNum)
	if err != nil {
		return gov.ParamSet{}, err
	}
	for k, v := range p2 {
		ret.Set(k, v)
	}
	return *ret, nil
}
