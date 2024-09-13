package contractgov

import "errors"

var (
	errNoChainConfig          = errors.New("ChainConfig or Istanbul is not set")
	errContractEngineNotReady = errors.New("ContractEngine is not ready")
	errParamsAtFail           = errors.New("headerGov EffectiveParams() failed")
	errGovParamNotExist       = errors.New("GovParam does not exist")
	errInvalidGovParam        = errors.New("GovParam conversion failed")
)
