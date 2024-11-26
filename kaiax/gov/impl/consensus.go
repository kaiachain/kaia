package impl

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
)

func (g *GovModule) VerifyHeader(header *types.Header) error {
	return g.Hgm.VerifyHeader(header)
}

func (g *GovModule) PrepareHeader(header *types.Header) error {
	return g.Hgm.PrepareHeader(header)
}

func (g *GovModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	return g.Hgm.FinalizeHeader(header, state, txs, receipts)
}
