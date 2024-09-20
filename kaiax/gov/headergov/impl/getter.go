package impl

import (
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (h *headerGovModule) EffectiveParamSet(blockNum uint64) (gov.ParamSet, error) {
	// TODO: only return when num <= head + 1
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.GetGovernanceHistory()
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		logger.Error("EffectiveParams error", "prevEpochStart", prevEpochStart, "blockNum", blockNum, "err", err,
			"govHistory", gh, "govs", h.cache.Govs())
		return gov.ParamSet{}, err
	} else {
		return gp, nil
	}
}

func (h *headerGovModule) EffectiveParamsPartial(blockNum uint64) (map[gov.ParamEnum]interface{}, error) {
	ret := make(map[gov.ParamEnum]interface{})
	for num, gov := range h.cache.Govs() {
		if num > blockNum {
			continue
		}
		for enum, value := range gov.Items() {
			ret[enum] = value
		}
	}

	return ret, nil
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
