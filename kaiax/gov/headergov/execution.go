package headergov

import (
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
)

func (h *headerGovModule) PostInsertBlock(b *types.Block) error {
	if len(b.Header().Vote) > 0 {
		vote, err := DeserializeHeaderVote(b.Header().Vote)
		if err != nil {
			logger.Error("DeserializeHeaderVote error", "vote", b.Header().Vote, "err", err)
			return err
		}
		err = h.HandleVote(b.NumberU64(), vote)
		if err != nil {
			logger.Error("HandleVote error", "vote", b.Header().Vote, "err", err)
			return err
		}
	}

	if len(b.Header().Governance) > 0 {
		gov, err := DeserializeHeaderGov(b.Header().Governance)
		if err != nil {
			logger.Error("DeserializeHeaderGov error", "governance", b.Header().Governance, "err", err)
			return err
		}
		err = h.HandleGov(b.NumberU64(), gov)
		if err != nil {
			logger.Error("HandleGov error", "governance", b.Header().Governance, "err", err)
			return err
		}
	}

	return nil
}

func (h *headerGovModule) HandleVote(blockNum uint64, vote VoteData) error {
	h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)

	var data StoredUint64Array = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &data)

	for i, myvote := range h.myVotes {
		if reflect.DeepEqual(&myvote, vote) {
			h.PopMyVotes(i)
			break
		}
	}

	return nil
}

func (h *headerGovModule) HandleGov(blockNum uint64, gov GovData) error {
	h.cache.AddGov(blockNum, gov)

	// merge gov based on latest effective params.
	gp, err := h.EffectiveParamSet(blockNum)
	if err != nil {
		logger.Error("kaiax.HandleGov error fetching EffectiveParams", "blockNum", blockNum, "gov", gov, "err", err)
		return err
	}

	err = gp.SetFromEnumMap(gov.Items())
	if err != nil {
		logger.Error("kaiax.HandleGov error setting paramset", "blockNum", blockNum, "gov", gov, "err", err, "gp", gp)
		return err
	}

	var data StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}
