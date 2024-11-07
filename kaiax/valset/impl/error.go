package impl

import (
	"errors"
)

var (
	errGenesisNotCalculable = errors.New("genesis block committee or proposer is not calculable")
	errInitUnexpectedNil    = errors.New("unexpected nil during module init")
	errExtractIstanbulExtra = errors.New("extract Istanbul Extra from block header of the given block number")
	errNilHeader            = errors.New("nil block header")
	errNilMixHash           = errors.New("nil mixHash on block header")
	errInvalidCommitteeSize = errors.New("invalid committee size")

	// rpc related errors
	errPendingNotAllowed = errors.New("pending is not allowed")
	errUnknownBlock      = errors.New("unknown block")
	errUnknownProposer   = errors.New("unknown proposer")

	// voting related errors
	errEmptyVoteBlock            = errors.New("failed to read vote blocks from db")
	errInvalidVoter              = errors.New("failed to verify voter")
	errInvalidVoteKey            = errors.New("your vote failed due to the wrong key")
	errInvalidVoteValue          = errors.New("your vote failed due to the wrong value")
	errCanonicalizeToAddressList = errors.New("could not canonicalize value to address list")
)
