package impl

import (
	"errors"
	"fmt"
)

var (
	ErrZeroEpoch                      = errors.New("epoch cannot be zero")
	ErrLowestVoteScannedBlockNotFound = errors.New("lowest vote scanned block not found")

	ErrVotePermissionDenied = errors.New("you don't have the right to vote")
	ErrInvalidKeyValue      = errors.New("your vote couldn't be placed. Please check your vote's key and value")

	ErrGovInNonEpochBlock = errors.New("governance is not allowed in non-epoch block")
	ErrNilVote            = errors.New("vote is nil")
	ErrGovVerification    = errors.New("header.Governance does not match the vote in previous epoch")

	ErrGovParamNotAccount       = errors.New("govparamcontract is not an account")
	ErrGovParamNotContract      = errors.New("govparamcontract is not an contract account")
	ErrLowerBoundBaseFee        = errors.New("lowerboundbasefee is greater than upperboundbasefee")
	ErrUpperBoundBaseFee        = errors.New("upperboundbasefee is less than lowerboundbasefee")
	ErrGovNodeInValSetVoteValue = errors.New("gov node is found in the valset vote value")
	ErrGovNodeNotInValSetList   = errors.New("gov node is not found in the valset list")
	ErrInvalidVoter             = errors.New("invalid voter")
)

func errInitNil(msg string) error {
	return fmt.Errorf("cannot init headergov module because of nil: %s", msg)
}
