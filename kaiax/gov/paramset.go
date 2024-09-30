package gov

import (
	"encoding/json"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
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
	param, ok := Params[name]
	if !ok {
		return ErrInvalidParamName
	}

	field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
	if !field.IsValid() || !field.CanSet() {
		return ErrCannotSet
	}

	fieldValue := reflect.ValueOf(cv)
	if !fieldValue.Type().AssignableTo(field.Type()) {
		return ErrInvalidParamValue
	}

	field.Set(fieldValue)
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

func (p *ParamSet) ToMap() map[ParamName]any {
	ret := make(map[ParamName]any)

	// Iterate through all params in Params and ensure they're in the result
	for name, param := range Params {
		field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
		if field.IsValid() {
			// Convert big.Int to string for JSON compatibility at API
			if bigIntValue, ok := field.Interface().(*big.Int); ok {
				ret[name] = bigIntValue.String()
			} else {
				ret[name] = field.Interface()
			}
		}
	}

	return ret
}

// TODO: remove this. Currently it's used for GetRewards API.
func (p *ParamSet) ToGovParamSet() *params.GovParamSet {
	m := make(map[string]any)
	for name := range Params {
		param := Params[name]
		fieldValue := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
		if fieldValue.IsValid() {
			m[string(name)] = fieldValue.Interface()
		}
	}

	ps, _ := params.NewGovParamSetStrMap(m)
	return ps
}

func (p PartialParamSet) Add(name string, value any) {
	param, ok := Params[ParamName(name)]
	if !ok {
		return
	}

	cv, err := param.Canonicalizer(value)
	if err != nil {
		return
	}

	if !param.FormatChecker(cv) {
		return
	}

	p[ParamName(name)] = cv
}
