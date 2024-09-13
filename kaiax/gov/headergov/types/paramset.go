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
	for _, param := range Params {
		err := p.Set(param.Name, param.DefaultValue)
		if err != nil {
			return nil
		}
	}

	return p
}

// Set the canonical value in the ParamSet for the corresponding parameter name.
func (p *ParamSet) Set(name string, cv interface{}) error {
	param, err := GetParamByName(name)
	if err != nil {
		return err
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

func (p *ParamSet) SetFromGovernanceData(g GovData) error {
	for name, value := range g.Items() {
		err := p.Set(name, value)
		if err != nil {
			continue
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

// ToStrMap is used for test and getParams API.
func (p *ParamSet) ToStrMap() (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	// Iterate through all params in Params and ensure they're in the result
	for _, param := range Params {
		field := reflect.ValueOf(p).Elem().FieldByName(param.ParamSetFieldName)
		if field.IsValid() {
			// Convert big.Int to string for JSON compatibility
			if bigIntValue, ok := field.Interface().(*big.Int); ok {
				ret[param.Name] = bigIntValue.String()
			} else {
				ret[param.Name] = field.Interface()
			}
		}
	}

	return ret, nil
}
