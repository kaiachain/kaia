package impl

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	contractgov_mock "github.com/kaiachain/kaia/v2/kaiax/gov/contractgov/mock"
	headergov_mock "github.com/kaiachain/kaia/v2/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/work/mocks"
	"github.com/stretchr/testify/assert"
)

func newGovModule(t *testing.T, chainConfig *params.ChainConfig) (*GovModule, *mocks.MockBlockChain, *headergov_mock.MockHeaderGovModule, *contractgov_mock.MockContractGovModule) {
	govModule := NewGovModule()
	mockChain := mocks.NewMockBlockChain(gomock.NewController(t))
	mockHgm := headergov_mock.NewMockHeaderGovModule(gomock.NewController(t))
	mockCgm := contractgov_mock.NewMockContractGovModule(gomock.NewController(t))

	govModule.Chain = mockChain
	govModule.Hgm = mockHgm
	govModule.Cgm = mockCgm
	return govModule, mockChain, mockHgm, mockCgm
}

func TestChainConfigFallback(t *testing.T) {
	tests := []struct {
		desc     string
		config   *params.ChainConfig
		expected gov.PartialParamSet
	}{
		{
			desc:     "nil config",
			config:   nil,
			expected: gov.PartialParamSet{},
		},
		{
			desc:     "empty config",
			config:   &params.ChainConfig{},
			expected: gov.PartialParamSet{gov.GovernanceUnitPrice: uint64(0)},
		},
		{
			desc: "config with non-default values",
			config: &params.ChainConfig{
				DeriveShaImpl: 1,
				UnitPrice:     uint64(123456789),
				Istanbul: &params.IstanbulConfig{
					Epoch:          uint64(100),
					ProposerPolicy: uint64(2),
					SubGroupSize:   uint64(99),
				},
				Governance: &params.GovernanceConfig{
					GovernanceMode:   "single",
					GoverningNode:    common.HexToAddress("0x1234567890123456789012345678901234567890"),
					GovParamContract: common.HexToAddress("0x2345678901234567890123456789012345678901"),
					Reward: &params.RewardConfig{
						MintingAmount:          big.NewInt(1e18),
						Ratio:                  "34/33/33",
						Kip82Ratio:             "50/50",
						UseGiniCoeff:           true,
						DeferredTxFee:          true,
						StakingUpdateInterval:  uint64(12345),
						ProposerUpdateInterval: uint64(6789),
						MinimumStake:           big.NewInt(5000000),
					},
					KIP71: &params.KIP71Config{
						LowerBoundBaseFee:         uint64(50000000000),
						UpperBoundBaseFee:         uint64(500000000000),
						GasTarget:                 uint64(40000000),
						MaxBlockGasUsedForBaseFee: uint64(80000000),
						BaseFeeDenominator:        uint64(25),
					},
				},
			},
			expected: gov.PartialParamSet{
				gov.GovernanceDeriveShaImpl:        uint64(1),
				gov.GovernanceUnitPrice:            uint64(123456789),
				gov.IstanbulEpoch:                  uint64(100),
				gov.IstanbulPolicy:                 uint64(2),
				gov.IstanbulCommitteeSize:          uint64(99),
				gov.GovernanceGovernanceMode:       "single",
				gov.GovernanceGoverningNode:        common.HexToAddress("0x1234567890123456789012345678901234567890"),
				gov.GovernanceGovParamContract:     common.HexToAddress("0x2345678901234567890123456789012345678901"),
				gov.RewardMintingAmount:            big.NewInt(1e18),
				gov.RewardRatio:                    "34/33/33",
				gov.RewardKip82Ratio:               "50/50",
				gov.RewardUseGiniCoeff:             true,
				gov.RewardDeferredTxFee:            true,
				gov.RewardStakingUpdateInterval:    uint64(12345),
				gov.RewardProposerUpdateInterval:   uint64(6789),
				gov.RewardMinimumStake:             big.NewInt(5000000),
				gov.Kip71LowerBoundBaseFee:         uint64(50000000000),
				gov.Kip71UpperBoundBaseFee:         uint64(500000000000),
				gov.Kip71GasTarget:                 uint64(40000000),
				gov.Kip71MaxBlockGasUsedForBaseFee: uint64(80000000),
				gov.Kip71BaseFeeDenominator:        uint64(25),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			fallback := ChainConfigFallback(tc.config)
			assert.Equal(t, tc.expected, fallback)
		})
	}
}
