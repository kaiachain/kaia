package contractgov

import (
	contractgov_types "github.com/kaiachain/kaia/kaiax/gov/contractgov/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/gov/headergov/types"
	gov_types "github.com/kaiachain/kaia/kaiax/gov/types"
)

type (
	ParamEnum         = gov_types.ParamEnum
	ParamSet          = gov_types.ParamSet
	ContractGovModule = contractgov_types.ContractGovModule
	HeaderGovModule   = headergov_types.HeaderGovModule
)

// Enums
var (
	GovernanceUnitPrice = gov_types.GovernanceUnitPrice
)

// Vars
var (
	Params          = gov_types.Params
	ParamNameToEnum = gov_types.ParamNameToEnum
)

// Functions
var (
	GetParamByName = gov_types.GetParamByName
)
