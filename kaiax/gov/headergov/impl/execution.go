package impl

import (
	"bytes"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax/gov"
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
	// if governance vote (i.e., not validator vote), add to vote
	if _, ok := gov.Params[vote.Name()]; ok {
		h.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)
		InsertVoteDataBlockNum(h.ChainKv, blockNum)
	}

	// if the vote was mine, remove it.
	for i, myvote := range h.myVotes {
		if bytes.Equal(myvote.Voter().Bytes(), vote.Voter().Bytes()) &&
			myvote.Name() == vote.Name() &&
			reflect.DeepEqual(myvote.Value(), vote.Value()) {
			h.PopMyVotes(i)
			break
		}
	}

	// if vote.key is in ValSetVoteKeyMap, do not store the voteBlock.
	if _, ok := gov.ValSetVoteKeyMap[vote.Name()]; ok {
		return nil
	}
	InsertVoteDataBlockNum(h.ChainKv, blockNum)

	return nil
}

func (h *headerGovModule) HandleGov(blockNum uint64, gov headergov.GovData) error {
	h.cache.AddGov(blockNum, gov)

	var data StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &data)
	return nil
}
