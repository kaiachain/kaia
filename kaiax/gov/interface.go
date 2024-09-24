package gov

import (
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=kaiax/gov/mock/govmodule_mock.go github.com/kaiachain/kaia/kaiax/gov GovModule
type GovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	EffectiveParamSet(blockNum uint64) (ParamSet, error)
}
