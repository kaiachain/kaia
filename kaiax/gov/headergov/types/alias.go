package types

import (
	govtypes "github.com/kaiachain/kaia/kaiax/gov/types"
)

type (
	Param     = govtypes.Param
	ParamEnum = govtypes.ParamEnum
	ParamSet  = govtypes.ParamSet
)

// Enums
var (
	GovernanceDeriveShaImpl        = govtypes.GovernanceDeriveShaImpl
	GovernanceGovernanceMode       = govtypes.GovernanceGovernanceMode
	GovernanceGoverningNode        = govtypes.GovernanceGoverningNode
	GovernanceGovParamContract     = govtypes.GovernanceGovParamContract
	GovernanceUnitPrice            = govtypes.GovernanceUnitPrice
	IstanbulCommitteeSize          = govtypes.IstanbulCommitteeSize
	IstanbulEpoch                  = govtypes.IstanbulEpoch
	IstanbulPolicy                 = govtypes.IstanbulPolicy
	Kip71BaseFeeDenominator        = govtypes.Kip71BaseFeeDenominator
	Kip71GasTarget                 = govtypes.Kip71GasTarget
	Kip71LowerBoundBaseFee         = govtypes.Kip71LowerBoundBaseFee
	Kip71MaxBlockGasUsedForBaseFee = govtypes.Kip71MaxBlockGasUsedForBaseFee
	Kip71UpperBoundBaseFee         = govtypes.Kip71UpperBoundBaseFee
	RewardDeferredTxFee            = govtypes.RewardDeferredTxFee
	RewardKip82Ratio               = govtypes.RewardKip82Ratio
	RewardMintingAmount            = govtypes.RewardMintingAmount
	RewardMinimumStake             = govtypes.RewardMinimumStake
	RewardProposerUpdateInterval   = govtypes.RewardProposerUpdateInterval
	RewardRatio                    = govtypes.RewardRatio
	RewardStakingUpdateInterval    = govtypes.RewardStakingUpdateInterval
	RewardUseGiniCoeff             = govtypes.RewardUseGiniCoeff
)

// Vars
var (
	Params          = govtypes.Params
	ParamNameToEnum = govtypes.ParamNameToEnum
)

// Functions
var (
	GetParamByName  = govtypes.GetParamByName
	StrMapToEnumMap = govtypes.StrMapToEnumMap
)
