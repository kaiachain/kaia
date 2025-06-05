package gov

import (
	"github.com/kaiachain/kaia/v2/kaiax"
)

//go:generate mockgen -destination=./mock/govmodule_mock.go -package=mock_gov github.com/kaiachain/kaia/v2/kaiax/gov GovModule
type GovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	GetParamSet(blockNum uint64) ParamSet
}
