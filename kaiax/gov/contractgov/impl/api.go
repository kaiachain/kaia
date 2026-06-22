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
	blockNum, err := api.blockNumber(num)
	if err != nil {
		return nil, err
	}
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

func (api *contractGovAPI) blockNumber(num rpc.BlockNumber) (uint64, error) {
	head := api.c.Chain.CurrentBlock().NumberU64()
	if num == rpc.LatestBlockNumber || num == rpc.PendingBlockNumber {
		return head, nil
	}
	if num.Int64() < 0 || num.Uint64() > head {
		return 0, gov.ErrUnknownBlock
	}
	return num.Uint64(), nil
}
