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

func (h *headerGovModule) EffectiveParamsPartial(blockNum uint64) gov.PartialParamSet {
	ret := make(gov.PartialParamSet)
	for _, num := range h.cache.GovBlockNums() {
		if num > blockNum {
			break
		}
		for name, value := range h.cache.Govs()[num].Items() {
			ret[name] = value
		}
	}

	return ret
}

func (h *headerGovModule) NodeAddress() common.Address {
	return h.nodeAddress
}

// GetLatestValidatorVote returns non-zero voteBlk and latest addvalidator or removevalidator vote
// If there's no vote, return 0, nil.
func (h *headerGovModule) GetLatestValidatorVote(num uint64) (uint64, headergov.VoteData) {
	votesBlks := h.cache.VoteBlockNums()

	for i := len(votesBlks) - 1; i >= 0; i-- {
		voteBlk := votesBlks[i]
		vote := h.cache.GroupedVotes()[calcEpochIdx(voteBlk, h.epoch)][voteBlk]
		if voteBlk < num && (vote.Name() == "governance.addvalidator" || vote.Name() == "governance.removevalidator") {
			return voteBlk, vote
		}
	}
	return 0, nil
}

func (h *headerGovModule) GetMyVotes() []headergov.VoteData {
	return h.myVotes
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
