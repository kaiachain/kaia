package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
)

func (g *GovModule) PostInsertBlock(b *types.Block) error {
	return g.Hgm.PostInsertBlock(b)
}
