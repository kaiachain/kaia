package types

import (
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=kaiax/gov/contractgov/mocks/contractgov_mock.go github.com/kaiachain/kaia/kaiax/gov/contractgov/types ContractGovModule
type ContractGovModule interface {
	kaiax.BaseModule

	EffectiveParamSet(blockNum uint64) (ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[ParamEnum]interface{}, error)
}
