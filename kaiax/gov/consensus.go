package gov

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
)

func (g *govModule) VerifyHeader(header *types.Header) error {
	return g.hgm.VerifyHeader(header)
}

func (g *govModule) PrepareHeader(header *types.Header) error {
	return g.hgm.PrepareHeader(header)
}

func (g *govModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	return g.hgm.FinalizeHeader(header, state, txs, receipts)
}
