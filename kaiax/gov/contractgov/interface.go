package contractgov

import (
	"github.com/kaiachain/kaia/v2/kaiax"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
)

//go:generate mockgen -destination=./mock/contractgov_mock.go -package=mock_contractgov github.com/kaiachain/kaia/v2/kaiax/gov/contractgov ContractGovModule
type ContractGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule

	GetParamSet(blockNum uint64) gov.ParamSet
	GetPartialParamSet(blockNum uint64) gov.PartialParamSet
}
