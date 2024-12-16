package impl

import (
	"errors"
	"fmt"
)

var (
	errInitUnexpectedNil     = errors.New("unexpected nil during module init")
	errNoHeader              = errors.New("no header found")
	errNoBlock               = errors.New("no block found")
	errInvalidProposerPolicy = errors.New("invalid proposer policy")
	errNoLowestScannedNum    = errors.New("no lowest scanned checkpoint interval")
	errEmptyVoteBlock        = errors.New("failed to read vote blocks from db")

	// rpc related errors
	errPendingNotAllowed = errors.New("pending is not allowed")
	errUnknownBlock      = errors.New("unknown block")
	errUnknownProposer   = errors.New("unknown proposer")
)

func ErrNoIstanbulSnapshot(num uint64) error {
	return fmt.Errorf("no istanbul snapshot at block %d", num)
}
