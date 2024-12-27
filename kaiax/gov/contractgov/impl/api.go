package impl

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/networks/rpc"
)

func (c *contractGovModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   NewContractGovAPI(c),
			Public:    true,
		},
	}
}

type contractGovAPI struct {
	c *contractGovModule
}

func NewContractGovAPI(c *contractGovModule) *contractGovAPI {
	return &contractGovAPI{c}
}

func (api *contractGovAPI) GetContractParams(num rpc.BlockNumber, govParam *common.Address) (gov.PartialParamSet, error) {
	blockNum := num.Uint64()
	var govParamAddr common.Address
	if govParam != nil {
		govParamAddr = *govParam
	} else {
		// Use default GovParam address from headergov
		govParamAddr = api.c.Hgm.GetParamSet(blockNum).GovParamContract
	}

	params, err := api.c.contractGetAllParamsAtFromAddr(blockNum, govParamAddr)
	if err != nil {
		return nil, err
	}

	return params, nil
}
