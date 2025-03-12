package impl

import (
	"github.com/kaiachain/kaia/kaiax/gov"
)

func (m *GovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	ret := gov.GetDefaultGovernanceParamSet()

	p0 := m.Fallback
	for k, v := range p0 {
		err := ret.Set(k, v)
		if err != nil {
			logger.CritWithStack("Failed to add param from Fallback", "name", k, "value", v, "error", err)
		}
	}

	p1 := m.Hgm.GetPartialParamSet(blockNum)
	for k, v := range p1 {
		err := ret.Set(k, v)
		if err != nil {
			logger.CritWithStack("Failed to add param from HeaderGov", "name", k, "value", v, "error", err)
		}
	}

	if m.isKoreHF(blockNum) {
		p2 := m.Cgm.GetPartialParamSet(blockNum)
		for k, v := range p2 {
			err := ret.Set(k, v)
			if err != nil {
				logger.CritWithStack("Failed to add param from ContractGov", "name", k, "value", v, "error", err)
			}
		}
	}

	return *ret
}
