package types

import (
	"github.com/kaiachain/kaia/kaiax"
)

type GovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	EffectiveParamSet(blockNum uint64) (ParamSet, error)
}
