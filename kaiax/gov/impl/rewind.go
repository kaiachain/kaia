package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (m *GovModule) RewindTo(newBlock *types.Block) {
	m.Hgm.RewindTo(newBlock)
}

func (m *GovModule) RewindDelete(hash common.Hash, num uint64) {
	m.Hgm.RewindDelete(hash, num)
}
