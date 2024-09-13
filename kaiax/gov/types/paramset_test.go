package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestParamSet_Set(t *testing.T) {
	tcs := []struct {
		ty    ParamEnum
		value interface{}
	}{
		{ty: GovernanceDeriveShaImpl, value: uint64(2)},
		{ty: GovernanceGovernanceMode, value: "none"},
		{ty: GovernanceGoverningNode, value: common.HexToAddress("0x000000000000000000000000000abcd000000000")},
		{ty: GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{ty: GovernanceUnitPrice, value: uint64(25e9)},
		{ty: IstanbulCommitteeSize, value: uint64(7)},
		{ty: IstanbulEpoch, value: uint64(406800)},
		{ty: IstanbulPolicy, value: uint64(2)},
		{ty: Kip71BaseFeeDenominator, value: uint64(64)},
		{ty: Kip71GasTarget, value: uint64(15000000)},
		{ty: Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{ty: Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{ty: Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{ty: RewardDeferredTxFee, value: true},
		{ty: RewardKip82Ratio, value: "10/90"},
		{ty: RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{ty: RewardMinimumStake, value: new(big.Int).SetUint64(2000000)},
		{ty: RewardProposerUpdateInterval, value: uint64(3600)},
		{ty: RewardRatio, value: "100/0/0"},
		{ty: RewardStakingUpdateInterval, value: uint64(86400)},
		{ty: RewardUseGiniCoeff, value: true},
	}

	p := ParamSet{}
	for _, tc := range tcs {
		param := Params[tc.ty]
		t.Run(param.Name, func(t *testing.T) {
			err := p.Set(param.Name, tc.value)
			assert.NoError(t, err)
		})
	}
}

func TestGetDefaultGovernanceParamSet(t *testing.T) {
	p := GetDefaultGovernanceParamSet()
	assert.NotNil(t, p)
	for _, param := range Params {
		t.Run(param.Name, func(t *testing.T) {
			fieldName := param.ParamSetFieldName
			fieldValue := reflect.ValueOf(p).Elem().FieldByName(fieldName).Interface()
			assert.Equal(t, param.DefaultValue, fieldValue, "Mismatch for %s", param.Name)
		})
	}
}
