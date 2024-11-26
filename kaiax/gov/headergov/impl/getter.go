package impl

import (
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
)

func (h *headerGovModule) EffectiveParamSet(blockNum uint64) gov.ParamSet {
	h.mu.RLock()
	defer h.mu.RUnlock()

	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.history
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		return *gov.GetDefaultGovernanceParamSet()
	}
	return gp
}

func (h *headerGovModule) EffectiveParamsPartial(blockNum uint64) gov.PartialParamSet {
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	ret := make(gov.PartialParamSet)

	// merge all governance sets before num's prevEpochStart.
	for _, num := range h.GovBlockNums() {
		if num <= prevEpochStart {
			for name, value := range h.governances[num].Items() {
				ret[name] = value
			}
		}
	}

	return ret
}

func (h *headerGovModule) NodeAddress() common.Address {
	return h.nodeAddress
}

func (h *headerGovModule) VoteBlockNums() []uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// TODO: replace with maps.Keys(h.groupedVotes)
	blockNums := make([]uint64, 0)
	for _, group := range h.groupedVotes {
		for num := range group {
			blockNums = append(blockNums, num)
		}
	}

	slices.Sort(blockNums)
	return blockNums
}

func (h *headerGovModule) GovBlockNums() []uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// TODO: replace with maps.Keys(h.Govs())
	blockNums := make([]uint64, 0)
	for num := range h.governances {
		blockNums = append(blockNums, num)
	}

	slices.Sort(blockNums)
	return blockNums
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
