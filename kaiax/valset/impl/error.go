package impl

import (
	"errors"
)

var (
	errGenesisNotCalculable = errors.New("genesis block committee or proposer is not calculable")
	errInitUnexpectedNil    = errors.New("unexpected nil during module init")
	errExtractIstanbulExtra = errors.New("extract Istanbul Extra from block header of the given block number")
	errNilHeader            = errors.New("nil block header")
	errInvalidCommitteeSize = errors.New("invalid committee size")

	errPendingNotAllowed = errors.New("pending is not allowed")
	errUnknownBlock      = errors.New("unknown block")
	errUnknownProposer   = errors.New("unknown proposer")
)
