package gov

import (
	contractgov_types "github.com/kaiachain/kaia/kaiax/gov/contractgov/types"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	headergov_types "github.com/kaiachain/kaia/kaiax/gov/headergov/types"
	govtypes "github.com/kaiachain/kaia/kaiax/gov/types"
)

type (
	ContractGovModule = contractgov_types.ContractGovModule

	GovModule = govtypes.GovModule
	ParamEnum = govtypes.ParamEnum
	ParamSet  = govtypes.ParamSet

	GovData         = headergov_types.GovData
	HeaderCache     = headergov_types.HeaderCache
	HeaderGovModule = headergov_types.HeaderGovModule
	History         = headergov_types.History
	VoteData        = headergov_types.VoteData
	VotesInEpoch    = headergov_types.VotesInEpoch
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
	DeserializeHeaderGov  = headergov_types.DeserializeHeaderGov
	DeserializeHeaderVote = headergov_types.DeserializeHeaderVote

	NewHeaderGovAPI = headergov.NewHeaderGovAPI

	NewGovData        = headergov_types.NewGovData
	NewHeaderGovCache = headergov_types.NewHeaderGovCache
	NewVoteData       = headergov_types.NewVoteData

	GetDefaultGovernanceParamSet = govtypes.GetDefaultGovernanceParamSet
)
