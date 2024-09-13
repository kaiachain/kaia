package types

import "errors"

var (
	ErrInvalidParamName  = errors.New("invalid param name")
	ErrInvalidParamValue = errors.New("invalid param value")
	ErrCannotSet         = errors.New("invalid field or cannot set the value")

	ErrInvalidVoteData   = errors.New("invalid vote data")
	ErrNotFoundInHistory = errors.New("blockNum not found from governance history")

	ErrCanonicalizeUint64        = errors.New("could not canonicalize value to uint64")
	ErrCanonicalizeString        = errors.New("could not canonicalize value to string")
	ErrCanonicalizeToAddress     = errors.New("could not canonicalize value to address")
	ErrCanonicalizeBigInt        = errors.New("could not canonicalize value to big.Int")
	ErrCanonicalizeBool          = errors.New("could not canonicalize value to bool")
	ErrCanonicalizeToAddressList = errors.New("could not canonicalize value to address list")

	ErrCanonicalizeByteToAddress   = errors.New("could not canonicalize []byte to address")
	ErrCanonicalizeByteToUint64    = errors.New("could not canonicalize []byte to uint64")
	ErrCanonicalizeFloatToUint64   = errors.New("could not canonicalize float64 to uint64")
	ErrCanonicalizeStringToAddress = errors.New("could not canonicalize string to address")
	ErrCanonicalizeByteToBigInt    = errors.New("could not canonicalize []byte to big.Int")
	ErrCanonicalizeStringToBigInt  = errors.New("could not canonicalize string to big.Int")
	ErrCanonicalizeByteToBool      = errors.New("could not canonicalize []byte to bool")
)
