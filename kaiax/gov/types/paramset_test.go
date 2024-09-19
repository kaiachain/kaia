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
		enum  ParamEnum
		value interface{}
	}{
		{enum: GovernanceDeriveShaImpl, value: uint64(2)},
		{enum: GovernanceGovernanceMode, value: "none"},
		{enum: GovernanceGoverningNode, value: common.HexToAddress("0x000000000000000000000000000abcd000000000")},
		{enum: GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{enum: GovernanceUnitPrice, value: uint64(25e9)},
		{enum: IstanbulCommitteeSize, value: uint64(7)},
		{enum: IstanbulEpoch, value: uint64(406800)},
		{enum: IstanbulPolicy, value: uint64(2)},
		{enum: Kip71BaseFeeDenominator, value: uint64(64)},
		{enum: Kip71GasTarget, value: uint64(15000000)},
		{enum: Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{enum: Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{enum: Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{enum: RewardDeferredTxFee, value: true},
		{enum: RewardKip82Ratio, value: "10/90"},
		{enum: RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{enum: RewardMinimumStake, value: new(big.Int).SetUint64(2000000)},
		{enum: RewardProposerUpdateInterval, value: uint64(3600)},
		{enum: RewardRatio, value: "100/0/0"},
		{enum: RewardStakingUpdateInterval, value: uint64(86400)},
		{enum: RewardUseGiniCoeff, value: true},
	}

	ps := ParamSet{}
	for _, tc := range tcs {
		t.Run(Params[tc.enum].Name, func(t *testing.T) {
			err := ps.Set(tc.enum, tc.value)
			assert.NoError(t, err)
		})
	}
}

func TestGetDefaultGovernanceParamSet(t *testing.T) {
	ps := GetDefaultGovernanceParamSet()
	assert.NotNil(t, ps)
	for _, param := range Params {
		t.Run(param.Name, func(t *testing.T) {
			fieldName := param.ParamSetFieldName
			fieldValue := reflect.ValueOf(ps).Elem().FieldByName(fieldName).Interface()
			assert.Equal(t, param.DefaultValue, fieldValue, "Mismatch for %s", param.Name)
		})
	}
}
