package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/kaiax/gov"
)

// EffectiveParamSet returns default parameter set in case of the following errors:
// (1) contractgov is disabled (i.e., pre-Kore or GovParam address is zero)
// (2) GovParam address is not set
// (3) Contract call to GovParam failed
// Invalid parameters in the contract (i.e., invalid parameter name or non-canonical value) are ignored.
func (c *contractGovModule) EffectiveParamSet(blockNum uint64) gov.ParamSet {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return *gov.GetDefaultGovernanceParamSet()
	}

	ret := *gov.GetDefaultGovernanceParamSet()
	for k, v := range m {
		err = ret.Set(k, v)
		if err != nil {
			return *gov.GetDefaultGovernanceParamSet()
		}
	}

	return ret
}

func (c *contractGovModule) EffectiveParamsPartial(blockNum uint64) gov.PartialParamSet {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return nil
	}
	return m
}

// TODO: add comments
func (c *contractGovModule) contractGetAllParamsAt(blockNum uint64) (gov.PartialParamSet, error) {
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
	contract, err := govcontract.NewGovParamCaller(addr, caller)
	if err != nil {
		return nil, err
	}

	names, values, err := contract.GetAllParamsAt(nil, new(big.Int).SetUint64(blockNum))
	if err != nil {
		logger.Warn("ContractEngine disabled: getAllParams call failed", "err", err)
		return nil, nil
	}

	if len(names) != len(values) {
		logger.Warn("ContractEngine disabled: getAllParams result invalid", "len(names)", len(names), "len(values)", len(values))
		return nil, nil
	}

	ret := ParseContractCall(names, values)
	return ret, nil
}

func (c *contractGovModule) contractAddrAt(blockNum uint64) (common.Address, error) {
	headerParams := c.hgm.EffectiveParamSet(blockNum)
	return headerParams.GovParamContract, nil
}

func ParseContractCall(names []string, values [][]byte) gov.PartialParamSet {
	ret := make(gov.PartialParamSet)
	for i := 0; i < len(names); i++ {
		ret.Add(names[i], values[i])
	}

	return ret
}
