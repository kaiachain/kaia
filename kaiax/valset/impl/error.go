package impl

import (
	"errors"
	"fmt"
)

var (
	errInitUnexpectedNil     = errors.New("unexpected nil during module init")
	errInvalidProposerPolicy = errors.New("invalid proposer policy")
	errNoHeader              = errors.New("no header found")
	errNoBlock               = errors.New("no block found")
	errNoLowestScannedNum    = errors.New("no lowest scanned validator vote num")
	errNoVoteBlockNums       = errors.New("no validator vote block nums")

	// rpc related errors
	errPendingNotAllowed = errors.New("pending is not allowed")
	errUnknownBlock      = errors.New("unknown block")
	errUnknownProposer   = errors.New("unknown proposer")
)

func ErrNoIstanbulSnapshot(num uint64) error {
	return fmt.Errorf("no istanbul snapshot at block %d", num)
}
