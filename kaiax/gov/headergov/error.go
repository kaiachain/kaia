package headergov

import "errors"

var (
	errZeroEpoch     = errors.New("epoch cannot be zero")
	errNoChainConfig = errors.New("ChainConfig or Istanbul is not set")
)
