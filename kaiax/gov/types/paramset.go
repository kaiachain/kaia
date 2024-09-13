package types

import (
	"encoding/json"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/common"
)

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
	p := &ParamSet{}
	for enum, param := range Params {
		err := p.Set(enum, param.DefaultValue)
		if err != nil {
			return nil
		}
	}

	return p
}

// Set the canonical value in the ParamSet for the corresponding parameter name.
func (p *ParamSet) Set(enum ParamEnum, cv interface{}) error {
	param, ok := Params[enum]
	if !ok {
		return ErrInvalidParamEnum
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

func (p *ParamSet) SetFromEnumMap(m map[ParamEnum]interface{}) error {
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

// ToEnumMap is used for test and getParams API.
func (p *ParamSet) ToEnumMap() (map[ParamEnum]interface{}, error) {
	ret := make(map[ParamEnum]interface{})

	// Iterate through all params in Params and ensure they're in the result
	for enum, param := range Params {
		field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
		if field.IsValid() {
			// Convert big.Int to string for JSON compatibility at API
			if bigIntValue, ok := field.Interface().(*big.Int); ok {
				ret[enum] = bigIntValue.String()
			} else {
				ret[enum] = field.Interface()
			}
		}
	}

	return ret, nil
}
