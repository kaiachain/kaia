package headergov

import (
	"errors"

	"github.com/kaiachain/kaia/log"
)

var logger = log.NewModuleLogger(log.KaiaxGov)

var (
	ErrInvalidRlp      = errors.New("invalid rlp")
	ErrInvalidJson     = errors.New("invalid json")
	ErrInvalidGovData  = errors.New("invalid gov data")
	ErrInvalidVoteData = errors.New("invalid vote data")
	ErrNoHistory       = errors.New("history search failed")
)
