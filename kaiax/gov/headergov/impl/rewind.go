package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (h *headerGovModule) RewindTo(newBlock *types.Block) {
	// do nothing
}

func (h *headerGovModule) RewindDelete(hash common.Hash, num uint64) {
	h.RemoveVotesAfter(num)
	h.RemoveGovAfter(num)
}

func (h *headerGovModule) RemoveVotesAfter(blockNum uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	dirty := false
	for epochIdxIter, votes := range h.groupedVotes {
		for blockNumIter := range votes {
			if blockNumIter > blockNum {
				dirty = true
				delete(h.groupedVotes[epochIdxIter], blockNumIter)

				// If all votes for this epoch have been removed, delete the epoch entry
				if len(h.groupedVotes[epochIdxIter]) == 0 {
					delete(h.groupedVotes, epochIdxIter)
				}
			}
		}
	}

	if dirty {
		WriteVoteDataBlockNums(h.ChainKv, h.VoteBlockNums())
	}
}

func (h *headerGovModule) RemoveGovAfter(blockNum uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	dirty := false
	for blockNumIter := range h.governances {
		if blockNumIter > blockNum {
			dirty = true
			delete(h.governances, blockNumIter)
		}
	}

	if dirty {
		WriteGovDataBlockNums(h.ChainKv, h.GovBlockNums())
		h.history = headergov.GovsToHistory(h.governances)
	}
}
