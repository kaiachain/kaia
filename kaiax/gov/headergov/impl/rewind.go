package impl

import (
	"slices"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (h *headerGovModule) RewindTo(newBlock *types.Block) {
	// do nothing
}

func (h *headerGovModule) RewindDelete(hash common.Hash, num uint64) {
	votesOld := h.cache.VoteBlockNums()
	govOld := h.cache.GovBlockNums()

	h.cache.RemoveVotesAfter(num)
	h.cache.RemoveGovAfter(num)

	// Update stored block numbers for votes and governance
	if votesNew := h.cache.VoteBlockNums(); !slices.Equal(votesOld, votesNew) {
		WriteVoteDataBlockNums(h.ChainKv, votesNew)
	}

	if govNew := h.cache.GovBlockNums(); !slices.Equal(govOld, govNew) {
		WriteGovDataBlockNums(h.ChainKv, govNew)
	}
}
