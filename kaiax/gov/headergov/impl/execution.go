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
		h.AddVote(blockNum, vote)
		InsertVoteDataBlockNum(h.ChainKv, blockNum)
	}

	// if the vote was mine, remove it.
	for i, myvote := range h.myVotes {
		if isEqualVotes(myvote, vote) {
			logger.Debug("Removing myvote", "vote", myvote)
			h.PopMyVotes(i)
			break
		}
	}

	return nil
}

func (h *headerGovModule) HandleGov(blockNum uint64, gov headergov.GovData) error {
	h.AddGov(blockNum, gov)

	data := h.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, data)
	return nil
}

func (h *headerGovModule) AddGov(blockNum uint64, gov headergov.GovData) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.governances[blockNum] = gov
	h.history = headergov.GovsToHistory(h.governances)
}

func (h *headerGovModule) AddVote(blockNum uint64, vote headergov.VoteData) {
	h.mu.Lock()
	defer h.mu.Unlock()

	epochIdx := calcEpochIdx(blockNum, h.epoch)

	if _, ok := h.groupedVotes[epochIdx]; !ok {
		h.groupedVotes[epochIdx] = make(headergov.VotesInEpoch)
	}
	h.groupedVotes[epochIdx][blockNum] = vote
}

func isEqualVotes(a headergov.VoteData, b headergov.VoteData) bool {
	return bytes.Equal(a.Voter().Bytes(), b.Voter().Bytes()) &&
		a.Name() == b.Name() &&
		reflect.DeepEqual(a.Value(), b.Value())
}
