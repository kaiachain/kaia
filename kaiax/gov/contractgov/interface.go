package contractgov

import (
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
)

//go:generate mockgen -destination=mock/contractgov_mock.go github.com/kaiachain/kaia/kaiax/gov/contractgov ContractGovModule
type ContractGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule

	GetParamSet(blockNum uint64) gov.ParamSet
	GetPartialParamSet(blockNum uint64) gov.PartialParamSet
}
