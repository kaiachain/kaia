package contractgov

import "errors"

var (
	ErrNoChainConfig = errors.New("ChainConfig or Istanbul is not set")
	ErrNotReady      = errors.New("ContractEngine is not ready")
	ErrHeaderGovFail = errors.New("headerGov EffectiveParams() failed")
)
