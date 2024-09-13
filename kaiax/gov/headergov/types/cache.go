package types

import (
	"sort"
)

type VotesInEpoch map[uint64]VoteData

// in-memory cache of data from DB
type HeaderCache struct {
	groupedVotes map[uint64]VotesInEpoch
	governances  map[uint64]GovData
	history      History
}

func NewHeaderGovCache() *HeaderCache {
	return &HeaderCache{
		groupedVotes: make(map[uint64]VotesInEpoch),
		governances:  make(map[uint64]GovData),
		history:      History{},
	}
}

func (h *HeaderCache) GroupedVotes() map[uint64]VotesInEpoch {
	votes := make(map[uint64]VotesInEpoch)
	for epochIdx, votesInEpoch := range h.groupedVotes {
		votes[epochIdx] = make(VotesInEpoch)
		for blockNum, vote := range votesInEpoch {
			votes[epochIdx][blockNum] = vote
		}
	}
	return votes
}

func (h *HeaderCache) Govs() map[uint64]GovData {
	govs := make(map[uint64]GovData)
	for blockNum, gov := range h.governances {
		govs[blockNum] = gov
	}
	return govs
}

func (h *HeaderCache) History() History {
	return h.history
}

func (h *HeaderCache) VoteBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.groupedVotes {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *HeaderCache) GovBlockNums() []uint64 {
	blockNums := make([]uint64, 0)
	for num := range h.governances {
		blockNums = append(blockNums, num)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})
	return blockNums
}

func (h *HeaderCache) AddVote(epochIdx, blockNum uint64, vote VoteData) {
	if _, ok := h.groupedVotes[epochIdx]; !ok {
		h.groupedVotes[epochIdx] = make(VotesInEpoch)
	}
	h.groupedVotes[epochIdx][blockNum] = vote
}

func (h *HeaderCache) AddGov(blockNum uint64, gov GovData) {
	h.governances[blockNum] = gov

	h.history = GetHistory(h.governances)
}

func (h *HeaderCache) RemoveVotesAfter(blockNum uint64) {
	for epochIdxIter, votes := range h.groupedVotes {
		for blockNumIter := range votes {
			if blockNumIter > blockNum {
				delete(h.groupedVotes[epochIdxIter], blockNumIter)
			}
		}
		// If all votes for this epoch have been removed, delete the epoch entry
		if len(h.groupedVotes[epochIdxIter]) == 0 {
			delete(h.groupedVotes, epochIdxIter)
		}
	}
}

func (h *HeaderCache) RemoveGovernanceAfter(blockNum uint64) {
	for blockNumIter := range h.governances {
		if blockNumIter > blockNum {
			delete(h.governances, blockNumIter)
		}
	}

	// Regenerate the governance history after removing entries
	h.history = GetHistory(h.governances)
}
