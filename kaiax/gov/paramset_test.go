package gov

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/stretchr/testify/assert"
)

func TestParametersSetAndAdd(t *testing.T) {
	tcs := []struct {
		name  ParamName
		value any
	}{
		{name: GovernanceDeriveShaImpl, value: uint64(2)},
		{name: GovernanceGovernanceMode, value: "none"},
		{name: GovernanceGoverningNode, value: common.HexToAddress("0x000000000000000000000000000abcd000000000")},
		{name: GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{name: GovernanceUnitPrice, value: uint64(25e9)},
		{name: IstanbulCommitteeSize, value: uint64(7)},
		{name: IstanbulEpoch, value: uint64(406800)},
		{name: IstanbulPolicy, value: uint64(2)},
		{name: Kip71BaseFeeDenominator, value: uint64(64)},
		{name: Kip71GasTarget, value: uint64(15000000)},
		{name: Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{name: Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{name: Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{name: RewardDeferredTxFee, value: true},
		{name: RewardKip82Ratio, value: "10/90"},
		{name: RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{name: RewardMinimumStake, value: new(big.Int).SetUint64(2000000)},
		{name: RewardProposerUpdateInterval, value: uint64(3600)},
		{name: RewardRatio, value: "100/0/0"},
		{name: RewardStakingUpdateInterval, value: uint64(86400)},
		{name: RewardUseGiniCoeff, value: true},
	}

	ps := ParamSet{}
	for _, tc := range tcs {
		t.Run(string(tc.name), func(t *testing.T) {
			err := ps.Set(tc.name, tc.value)
			assert.NoError(t, err)

			pps := PartialParamSet{}
			pps.Add(string(tc.name), tc.value)
			assert.Equal(t, 1, len(pps))
		})
	}
}

func TestGetDefaultGovernanceParamSet(t *testing.T) {
	ps := GetDefaultGovernanceParamSet()
	assert.NotNil(t, ps)
}
