package headergov

import (
	headergov_types "github.com/kaiachain/kaia/kaiax/gov/headergov/types"
	gov_types "github.com/kaiachain/kaia/kaiax/gov/types"
)

type (
	GovData         = headergov_types.GovData
	HeaderCache     = headergov_types.HeaderCache
	HeaderGovModule = headergov_types.HeaderGovModule
	History         = headergov_types.History
	VoteData        = headergov_types.VoteData
	VotesInEpoch    = headergov_types.VotesInEpoch

	ParamEnum = gov_types.ParamEnum
	ParamSet  = gov_types.ParamSet
)

// Enums
var (
	GovernanceDeriveShaImpl        = gov_types.GovernanceDeriveShaImpl
	GovernanceGovernanceMode       = gov_types.GovernanceGovernanceMode
	GovernanceGoverningNode        = gov_types.GovernanceGoverningNode
	GovernanceGovParamContract     = gov_types.GovernanceGovParamContract
	GovernanceUnitPrice            = gov_types.GovernanceUnitPrice
	IstanbulCommitteeSize          = gov_types.IstanbulCommitteeSize
	IstanbulEpoch                  = gov_types.IstanbulEpoch
	IstanbulPolicy                 = gov_types.IstanbulPolicy
	Kip71BaseFeeDenominator        = gov_types.Kip71BaseFeeDenominator
	Kip71GasTarget                 = gov_types.Kip71GasTarget
	Kip71LowerBoundBaseFee         = gov_types.Kip71LowerBoundBaseFee
	Kip71MaxBlockGasUsedForBaseFee = gov_types.Kip71MaxBlockGasUsedForBaseFee
	Kip71UpperBoundBaseFee         = gov_types.Kip71UpperBoundBaseFee
	RewardDeferredTxFee            = gov_types.RewardDeferredTxFee
	RewardKip82Ratio               = gov_types.RewardKip82Ratio
	RewardMintingAmount            = gov_types.RewardMintingAmount
	RewardMinimumStake             = gov_types.RewardMinimumStake
	RewardProposerUpdateInterval   = gov_types.RewardProposerUpdateInterval
	RewardRatio                    = gov_types.RewardRatio
	RewardStakingUpdateInterval    = gov_types.RewardStakingUpdateInterval
	RewardUseGiniCoeff             = gov_types.RewardUseGiniCoeff
)

// Vars
var (
	Params          = gov_types.Params
	ParamNameToEnum = gov_types.ParamNameToEnum
)

// Functions
var (
	DeserializeHeaderGov  = headergov_types.DeserializeHeaderGov
	DeserializeHeaderVote = headergov_types.DeserializeHeaderVote

	NewGovData        = headergov_types.NewGovData
	NewHeaderGovCache = headergov_types.NewHeaderGovCache
	NewVoteData       = headergov_types.NewVoteData

	GetDefaultGovernanceParamSet = gov_types.GetDefaultGovernanceParamSet
	EnumMapToStrMap              = gov_types.EnumMapToStrMap
	StrMapToEnumMap              = gov_types.StrMapToEnumMap
)
