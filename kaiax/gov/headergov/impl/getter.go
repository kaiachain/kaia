package impl

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (h *headerGovModule) EffectiveParamSet(blockNum uint64) gov.ParamSet {
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.GetGovernanceHistory()
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		return *gov.GetDefaultGovernanceParamSet()
	}
	return gp
}

func (h *headerGovModule) EffectiveParamsPartial(blockNum uint64) map[gov.ParamName]any {
	ret := make(map[gov.ParamName]any)
	for num, gov := range h.cache.Govs() {
		if num > blockNum {
			continue
		}
		for enum, value := range gov.Items() {
			ret[enum] = value
		}
	}

	return ret
}

func (h *headerGovModule) NodeAddress() common.Address {
	return h.nodeAddress
}

func (h *headerGovModule) GetGovernanceHistory() headergov.History {
	return h.cache.History()
}

func PrevEpochStart(blockNum, epoch uint64, isKore bool) uint64 {
	if blockNum <= epoch {
		return 0
	}
	if !isKore {
		blockNum -= 1
	}
	return blockNum - blockNum%epoch - epoch
}
