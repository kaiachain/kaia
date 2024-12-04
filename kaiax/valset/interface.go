package valset

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=kaiax/valset/mock/valsetmodule_mock.go github.com/kaiachain/kaia/kaiax/valset ValsetModule
type ValsetModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ExecutionModule

	GetCouncil(num uint64) (AddressList, error)
	GetCommittee(num uint64, round uint64) (AddressList, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
}
