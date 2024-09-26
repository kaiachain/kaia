package contractgov

import (
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
)

//go:generate mockgen -destination=kaiax/gov/contractgov/mock/contractgov_mock.go github.com/kaiachain/kaia/kaiax/gov/contractgov ContractGovModule
type ContractGovModule interface {
	kaiax.BaseModule

	EffectiveParamSet(blockNum uint64) (gov.ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[gov.ParamEnum]any, error)
}
