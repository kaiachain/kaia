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
	ParamSet  = govtypes.ParamSet

	GovData         = headergov_types.GovData
	HeaderCache     = headergov_types.HeaderCache
	HeaderGovModule = headergov_types.HeaderGovModule
	History         = headergov_types.History
	VoteData        = headergov_types.VoteData
	VotesInEpoch    = headergov_types.VotesInEpoch
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
