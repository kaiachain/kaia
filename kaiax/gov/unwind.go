package gov

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (m *govModule) RewindTo(newBlock *types.Block) {
	m.hgm.RewindTo(newBlock)
}

func (m *govModule) RewindDelete(hash common.Hash, num uint64) {
	m.hgm.RewindDelete(hash, num)
}
