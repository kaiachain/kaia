package types

import (
	"github.com/kaiachain/kaia/kaiax"
)

type StakingModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.UnwindableModule

	GetStakingInfo(num uint64) (*StakingInfo, error)
}
