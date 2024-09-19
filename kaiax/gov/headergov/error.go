package headergov

import "errors"

var (
	ErrZeroEpoch     = errors.New("epoch cannot be zero")
	ErrNoChainConfig = errors.New("ChainConfig or Istanbul is not set")

	ErrVotePermissionDenied = errors.New("you don't have the right to vote")
	ErrInvalidKeyValue      = errors.New("your vote couldn't be placed. Please check your vote's key and value")

	ErrGovInNonEpochBlock = errors.New("governance is not allowed in non-epoch block")
	ErrNilVote            = errors.New("vote is nil")
	ErrGovVerification    = errors.New("header.Governance does not match the vote in previous epoch")

	ErrGovParamNotAccount  = errors.New("govparamcontract is not an account")
	ErrGovParamNotContract = errors.New("govparamcontract is not an contract account")
	ErrLowerBoundBaseFee   = errors.New("lowerboundbasefee is greater than upperboundbasefee")
	ErrUpperBoundBaseFee   = errors.New("upperboundbasefee is less than lowerboundbasefee")
)
