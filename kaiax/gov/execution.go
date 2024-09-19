package gov

import "github.com/kaiachain/kaia/blockchain/types"

func (g *GovModule) PostInsertBlock(b *types.Block) error {
	logger.Info("PostInsertBlock", "block number", b.NumberU64())
	return g.hgm.PostInsertBlock(b)
}
