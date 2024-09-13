package types

import "errors"

var (
	ErrInvalidParamName  = errors.New("invalid param name")
	ErrInvalidParamValue = errors.New("invalid param value")
	ErrCannotSet         = errors.New("invalid field or cannot set the value")

	ErrInvalidVoteData   = errors.New("invalid vote data")
	ErrNotFoundInHistory = errors.New("blockNum not found from governance history")
)
