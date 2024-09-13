package gov

import "github.com/kaiachain/kaia/blockchain/types"

func (g *govModule) PostInsertBlock(b *types.Block) error {
	return g.hgm.PostInsertBlock(b)
}
