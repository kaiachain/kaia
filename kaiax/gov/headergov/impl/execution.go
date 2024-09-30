package impl

import (
	"bytes"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (h *headerGovModule) PostInsertBlock(b *types.Block) error {
	if len(b.Header().Vote) > 0 {
		var vb headergov.VoteBytes = b.Header().Vote
		vote, err := vb.ToVoteData()
		if err != nil {
			logger.Error("ToVoteData error", "vote", vb, "err", err)
			return err
		}
		err = h.HandleVote(b.NumberU64(), vote)
		if err != nil {
			logger.Error("HandleVote error", "vote", vb, "err", err)
			return err
		}
	}

	if len(b.Header().Governance) > 0 {
		var gb headergov.GovBytes = b.Header().Governance
		gov, err := gb.ToGovData()
		if err != nil {
			logger.Error("DeserializeHeaderGov error", "governance", gb, "err", err)
			return err
		}
		err = h.HandleGov(b.NumberU64(), gov)
		if err != nil {
			logger.Error("HandleGov error", "governance", gb, "err", err)
			return err
		}
	}

	return nil
}

func (h *headerGovModule) HandleVote(blockNum uint64, vote headergov.VoteData) error {
	h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)

	var data StoredUint64Array = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &data)

	// if the vote was mine, remove it.
	for i, myvote := range h.myVotes {
		if bytes.Equal(myvote.Voter().Bytes(), vote.Voter().Bytes()) &&
			myvote.Name() == vote.Name() &&
			myvote.Value() == vote.Value() {
			h.PopMyVotes(i)
			break
		}
	}

	return nil
}

func (h *headerGovModule) HandleGov(blockNum uint64, gov headergov.GovData) error {
	h.cache.AddGov(blockNum, gov)

	// merge gov based on latest effective params.
	gp := h.EffectiveParamSet(blockNum)
	err := gp.SetFromMap(gov.Items())
	if err != nil {
		logger.Error("kaiax.HandleGov error setting paramset", "blockNum", blockNum, "gov", gov, "err", err, "gp", gp)
		return err
	}

	var data StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}
