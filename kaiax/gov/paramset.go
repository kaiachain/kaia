package gov

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/params"
)

type PartialParamSet map[ParamName]any

type ParamSet struct {
	// governance
	GovernanceMode                  string
	GoverningNode, GovParamContract common.Address

	// istanbul
	CommitteeSize, ProposerPolicy, Epoch uint64

	// reward
	Ratio, Kip82Ratio                             string
	StakingUpdateInterval, ProposerUpdateInterval uint64
	MintingAmount, MinimumStake                   *big.Int
	UseGiniCoeff, DeferredTxFee                   bool

	// KIP-71
	LowerBoundBaseFee, UpperBoundBaseFee, GasTarget, MaxBlockGasUsedForBaseFee, BaseFeeDenominator uint64

	// etc.
	DeriveShaImpl uint64
	UnitPrice     uint64
}

// GetDefaultGovernanceParamSet must not return nil, which is unit-tested.
func GetDefaultGovernanceParamSet() *ParamSet {
	ps := &ParamSet{}
	for name, param := range Params {
		err := ps.Set(name, param.DefaultValue)
		if err != nil {
			return nil
		}
	}

	return ps
}

// Set the canonical value in the ParamSet for the corresponding parameter name.
func (p *ParamSet) Set(name ParamName, cv any) error {
	var ok bool
	switch name {
	case GovernanceGovernanceMode:
		p.GovernanceMode, ok = cv.(string)
	case GovernanceGoverningNode:
		p.GoverningNode, ok = cv.(common.Address)
	case GovernanceGovParamContract:
		p.GovParamContract, ok = cv.(common.Address)
	case GovernanceUnitPrice:
		p.UnitPrice, ok = cv.(uint64)
	case IstanbulCommitteeSize:
		p.CommitteeSize, ok = cv.(uint64)
	case IstanbulEpoch:
		p.Epoch, ok = cv.(uint64)
	case IstanbulPolicy:
		p.ProposerPolicy, ok = cv.(uint64)
	case Kip71BaseFeeDenominator:
		p.BaseFeeDenominator, ok = cv.(uint64)
	case Kip71GasTarget:
		p.GasTarget, ok = cv.(uint64)
	case Kip71LowerBoundBaseFee:
		p.LowerBoundBaseFee, ok = cv.(uint64)
	case Kip71MaxBlockGasUsedForBaseFee:
		p.MaxBlockGasUsedForBaseFee, ok = cv.(uint64)
	case Kip71UpperBoundBaseFee:
		p.UpperBoundBaseFee, ok = cv.(uint64)
	case RewardDeferredTxFee:
		p.DeferredTxFee, ok = cv.(bool)
	case RewardKip82Ratio:
		p.Kip82Ratio, ok = cv.(string)
	case RewardMintingAmount:
		p.MintingAmount, ok = cv.(*big.Int)
	case RewardMinimumStake:
		p.MinimumStake, ok = cv.(*big.Int)
	case RewardProposerUpdateInterval:
		p.ProposerUpdateInterval, ok = cv.(uint64)
	case RewardRatio:
		p.Ratio, ok = cv.(string)
	case RewardStakingUpdateInterval:
		p.StakingUpdateInterval, ok = cv.(uint64)
	case RewardUseGiniCoeff:
		p.UseGiniCoeff, ok = cv.(bool)
	case GovernanceDeriveShaImpl:
		p.DeriveShaImpl, ok = cv.(uint64)
	default:
		return ErrInvalidParamName
	}

	if !ok {
		return ErrInvalidParamValue
	}
	return nil
}

func (p *ParamSet) SetFromMap(m PartialParamSet) error {
	for name, value := range m {
		err := p.Set(name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ParamSet) ToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// TODO: remove this. Currently it's used for kaia_getParams API.
func (p *ParamSet) ToMap() map[ParamName]any {
	ret := make(map[ParamName]any)

	// Iterate through all params in Params and ensure they're in the result
	for name := range Params {
		switch name {
		case GovernanceGovernanceMode:
			ret[name] = p.GovernanceMode
		case GovernanceGoverningNode:
			ret[name] = p.GoverningNode
		case GovernanceGovParamContract:
			ret[name] = p.GovParamContract
		case GovernanceUnitPrice:
			ret[name] = p.UnitPrice
		case IstanbulCommitteeSize:
			ret[name] = p.CommitteeSize
		case IstanbulEpoch:
			ret[name] = p.Epoch
		case IstanbulPolicy:
			ret[name] = p.ProposerPolicy
		case Kip71BaseFeeDenominator:
			ret[name] = p.BaseFeeDenominator
		case Kip71GasTarget:
			ret[name] = p.GasTarget
		case Kip71LowerBoundBaseFee:
			ret[name] = p.LowerBoundBaseFee
		case Kip71MaxBlockGasUsedForBaseFee:
			ret[name] = p.MaxBlockGasUsedForBaseFee
		case Kip71UpperBoundBaseFee:
			ret[name] = p.UpperBoundBaseFee
		case RewardDeferredTxFee:
			ret[name] = p.DeferredTxFee
		case RewardKip82Ratio:
			ret[name] = p.Kip82Ratio
		case RewardMintingAmount:
			ret[name] = p.MintingAmount.String()
		case RewardMinimumStake:
			ret[name] = p.MinimumStake.String()
		case RewardProposerUpdateInterval:
			ret[name] = p.ProposerUpdateInterval
		case RewardRatio:
			ret[name] = p.Ratio
		case RewardStakingUpdateInterval:
			ret[name] = p.StakingUpdateInterval
		case RewardUseGiniCoeff:
			ret[name] = p.UseGiniCoeff
		case GovernanceDeriveShaImpl:
			ret[name] = p.DeriveShaImpl
		}
	}

	return ret
}

// TODO: remove this. Currently it's used for GetRewards API.
func (p *ParamSet) ToGovParamSet() *params.GovParamSet {
	m := make(map[string]any)
	for name, val := range p.ToMap() {
		m[string(name)] = val
	}

	ps, _ := params.NewGovParamSetStrMap(m)
	return ps
}

func (p *ParamSet) ToKip71Config() *params.KIP71Config {
	return &params.KIP71Config{
		LowerBoundBaseFee:         p.LowerBoundBaseFee,
		UpperBoundBaseFee:         p.UpperBoundBaseFee,
		GasTarget:                 p.GasTarget,
		MaxBlockGasUsedForBaseFee: p.MaxBlockGasUsedForBaseFee,
		BaseFeeDenominator:        p.BaseFeeDenominator,
	}
}

func (p PartialParamSet) Add(name string, value any) error {
	param, ok := Params[ParamName(name)]
	if !ok {
		return ErrInvalidParamName
	}

	cv, err := param.Canonicalizer(value)
	if err != nil {
		return err
	}

	if !param.FormatChecker(cv) {
		return ErrInvalidParamValue
	}

	p[ParamName(name)] = cv
	return nil
}
