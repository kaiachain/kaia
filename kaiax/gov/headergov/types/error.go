package types

import "errors"

var (
	ErrInvalidGovData  = errors.New("invalid gov data")
	ErrInvalidVoteData = errors.New("invalid vote data")
	ErrNoHistory       = errors.New("history search failed")
)
