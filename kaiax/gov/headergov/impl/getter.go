package impl

import (
	"slices"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"golang.org/x/exp/maps" // TODO: use "maps"
)

func (h *headerGovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	h.mu.RLock()
	defer h.mu.RUnlock()

	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	gh := h.history
	gp, err := gh.Search(prevEpochStart)
	if err != nil {
		logger.Warn("No param set", "blockNum", blockNum, "prevEpochStart", prevEpochStart)
		return *gov.GetDefaultGovernanceParamSet()
	}
	return gp
}

func (h *headerGovModule) GetPartialParamSet(blockNum uint64) gov.PartialParamSet {
	prevEpochStart := PrevEpochStart(blockNum, h.epoch, h.isKoreHF(blockNum))
	ret := make(gov.PartialParamSet)

	// merge all governance sets before num's prevEpochStart.
	for _, num := range h.GovBlockNums() {
		if num <= prevEpochStart {
			for name, value := range h.governances[num].Items() {
				ret.Add(string(name), value)
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

	blockNums := []uint64{}
	for _, votes := range h.groupedVotes {
		blockNums = append(blockNums, maps.Keys(votes)...)
	}

	slices.Sort(blockNums)
	return blockNums
}

func (h *headerGovModule) GovBlockNums() []uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	blockNums := maps.Keys(h.governances)
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
