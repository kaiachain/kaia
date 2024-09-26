package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/kaiax/gov"
)

func (c *contractGovModule) EffectiveParamSet(blockNum uint64) gov.ParamSet {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return *gov.GetDefaultGovernanceParamSet()
	}

	ret := gov.ParamSet{}
	for k, v := range m {
		err = ret.Set(k, v)
		if err != nil {
			return *gov.GetDefaultGovernanceParamSet()
		}
	}

	return ret
}

func (c *contractGovModule) EffectiveParamsPartial(blockNum uint64) map[gov.ParamName]any {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return nil
	}
	return m
}

// TODO: add comments
func (c *contractGovModule) contractGetAllParamsAt(blockNum uint64) (map[gov.ParamName]any, error) {
	chain := c.Chain
	if chain == nil {
		return nil, ErrNotReady
	}

	config := c.ChainConfig
	if !config.IsKoreForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, ErrNotReady
	}

	addr, err := c.contractAddrAt(blockNum)
	if err != nil {
		return nil, err
	}
	if common.EmptyAddress(addr) {
		logger.Trace("ContractEngine disabled: GovParamContract address not set")
		return nil, nil
	}

	caller := backends.NewBlockchainContractBackend(chain, nil, nil)
	contract, _ := govcontract.NewGovParamCaller(addr, caller)

	names, values, err := contract.GetAllParamsAt(nil, new(big.Int).SetUint64(blockNum))
	if err != nil {
		logger.Warn("ContractEngine disabled: getAllParams call failed", "err", err)
		return nil, nil
	}

	if len(names) != len(values) {
		logger.Warn("ContractEngine disabled: getAllParams result invalid", "len(names)", len(names), "len(values)", len(values))
		return nil, nil
	}

	ret := make(map[gov.ParamName]any)
	for i := 0; i < len(names); i++ {
		param, ok := gov.Params[gov.ParamName(names[i])]
		if !ok {
			return nil, gov.ErrInvalidParamName
		}
		cv, err := param.Canonicalizer(values[i])
		if err != nil {
			return nil, err
		}
		ret[gov.ParamName(names[i])] = cv
	}

	return ret, nil
}

func (c *contractGovModule) contractAddrAt(blockNum uint64) (common.Address, error) {
	headerParams := c.hgm.EffectiveParamSet(blockNum)
	return headerParams.GovParamContract, nil
}
