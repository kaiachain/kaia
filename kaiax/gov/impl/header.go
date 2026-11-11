package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
)

func (g *GovModule) VerifyHeader(header *types.Header, parent *types.Header) error {
	return g.Hgm.VerifyHeader(header, parent)
}

func (g *GovModule) PrepareHeader(header *types.Header) error {
	return g.Hgm.PrepareHeader(header)
}
