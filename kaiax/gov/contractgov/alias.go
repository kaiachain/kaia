package contractgov

import (
	contractgov_types "github.com/kaiachain/kaia/kaiax/gov/contractgov/types"
	headergov_types "github.com/kaiachain/kaia/kaiax/gov/headergov/types"
	govtypes "github.com/kaiachain/kaia/kaiax/gov/types"
)

type (
	ParamEnum         = govtypes.ParamEnum
	ParamSet          = govtypes.ParamSet
	ContractGovModule = contractgov_types.ContractGovModule
	HeaderGovModule   = headergov_types.HeaderGovModule
)

// Enums
var (
	GovernanceUnitPrice = govtypes.GovernanceUnitPrice
)

// Vars
var (
	Params          = govtypes.Params
	ParamNameToEnum = govtypes.ParamNameToEnum
)

// Functions
var (
	GetParamByName = govtypes.GetParamByName
)
