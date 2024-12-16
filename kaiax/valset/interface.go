package valset

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/valset ValsetModule
type ValsetModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	GetCouncil(num uint64) ([]common.Address, error)
	GetCommittee(num uint64, round uint64) ([]common.Address, error)
	GetDemotedValidators(num uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
}
