package valset

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=kaiax/valset/mock/valsetmodule_mock.go github.com/kaiachain/kaia/kaiax/valset ValsetModule
type ValsetModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule

	GetCouncilAddressList(num uint64) ([]common.Address, error)
	GetCommitteeAddressList(num uint64, round uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
	Vote(blockNumber uint64, voter common.Address, name string, value any) (string, error)
}
