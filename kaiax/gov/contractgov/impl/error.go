package impl

import "errors"

var (
	ErrInitNil = errors.New("cannot init contractgov module because of nil")

	ErrNotReady      = errors.New("ContractEngine is not ready")
	ErrHeaderGovFail = errors.New("headerGov EffectiveParams() failed")
)
