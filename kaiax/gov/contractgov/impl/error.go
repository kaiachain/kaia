package impl

import (
	"errors"
	"fmt"
)

var (
	ErrNotReady      = errors.New("ContractEngine is not ready")
	ErrHeaderGovFail = errors.New("headerGov GetParamSet() failed")
)

func errInitNil(msg string) error {
	return fmt.Errorf("cannot init contractgov module because of nil: %s", msg)
}
