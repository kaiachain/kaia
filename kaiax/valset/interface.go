package valset

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

type ValsetModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule

	GetCouncilAddressList(num uint64) ([]common.Address, error)
	GetCommitteeAddressList(num uint64, round uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
}
