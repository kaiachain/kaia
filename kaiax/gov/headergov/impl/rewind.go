package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (h *headerGovModule) RewindTo(newBlock *types.Block) {
	// Remove entries from h.cache that are larger than num
	h.cache.RemoveVotesAfter(newBlock.NumberU64())
	h.cache.RemoveGovAfter(newBlock.NumberU64())
}

func (h *headerGovModule) RewindDelete(hash common.Hash, num uint64) {
	h.cache.RemoveVotesAfter(num)
	h.cache.RemoveGovAfter(num)

	// Update stored block numbers for votes and governance
	var voteBlockNums StoredUint64Array = h.cache.VoteBlockNums()
	WriteVoteDataBlockNums(h.ChainKv, &voteBlockNums)

	var govBlockNums StoredUint64Array = h.cache.GovBlockNums()
	WriteGovDataBlockNums(h.ChainKv, &govBlockNums)
}
