package gov

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	for name, param := range Params {
		assert.NotEmpty(t, param.Canonicalizer, name)
		assert.NotEmpty(t, param.FormatChecker, name)
		assert.NotEmpty(t, param.ChainConfigValue, name)
	}
}

func TestAddressCanonicalizer(t *testing.T) {
	tcs := []struct {
		desc          string
		input         any
		expected      common.Address
		expectedError error
	}{
		{desc: "Valid byte slice", input: common.FromHex("0x0102030405060708090a0b0c0d0e0f1011121314"), expected: common.HexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314")},
		{desc: "Valid hex string", input: "0x0102030405060708090a0b0c0d0e0f1011121314", expected: common.HexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314")},
		{desc: "Valid hex string", input: "0102030405060708090a0b0c0d0e0f1011121314", expected: common.HexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314")},
		{desc: "Valid common.Address", input: common.HexToAddress("0x1234567890123456789012345678901234567890"), expected: common.HexToAddress("0x1234567890123456789012345678901234567890")},
		{desc: "Invalid byte slice length", input: []byte{1, 2, 3}, expectedError: ErrCanonicalizeByteToAddress},
		{desc: "Invalid hex string - non-hexdigits", input: "0xinvalid", expectedError: ErrCanonicalizeStringToAddress},
		{desc: "Invalid hex string - length 2", input: "01", expectedError: ErrCanonicalizeStringToAddress},
		{desc: "Invalid hex string - length 38", input: "0102030405060708090a0b0c0d0e0f10111213", expectedError: ErrCanonicalizeStringToAddress},
		{desc: "Invalid type", input: 123, expectedError: ErrCanonicalizeToAddress},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := addressCanonicalizer(tc.input)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result.(common.Address))
			}
		})
	}
}

func TestValidatorAddressCanonicalizer(t *testing.T) {
	tcs := []struct {
		desc          string
		input         any
		expected      []common.Address
		expectedError error
	}{
		{desc: "Valid string, single address", input: "0x1234567890123456789012345678901234567890", expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890")}},
		{desc: "Valid string, multiple addresses", input: "0x1234567890123456789012345678901234567890,0x0987654321098765432109876543210987654321", expected: []common.Address{common.HexToAddress("0x1234567890123456789012345678901234567890"), common.HexToAddress("0x0987654321098765432109876543210987654321")}},
		{desc: "Invalid string", input: "0xinvalid", expectedError: ErrCanonicalizeStringToAddress},
		{desc: "Valid bytes, one address", input: hexutil.MustDecode("0x1212121212121212121212121212121212121212"), expected: []common.Address{common.HexToAddress("0x1212121212121212121212121212121212121212")}},
		{desc: "Valid bytes, hex-encoded one address", input: hexutil.MustDecode("0x307831366331393235383561306162323462353532373833623462663764386463396636383535633335"), expected: []common.Address{common.HexToAddress("0x3833623462663764386463396636383535633335")}},
		{desc: "Valid bytes, two addresses", input: hexutil.MustDecode("0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc90f79bf6eb2c4f870365e785982e1f101e93b906"), expected: []common.Address{common.HexToAddress("0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc"), common.HexToAddress("0x90f79bf6eb2c4f870365e785982e1f101e93b906")}},
		{desc: "Invalid type", input: 123, expectedError: ErrCanonicalizeToAddressList},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := validatorAddressListCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.([]common.Address))
			}
		})
	}
}

func TestBigIntCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         any
		expected      *big.Int
		expectedError error
	}{
		{name: "Valid big.Int", input: big.NewInt(12345), expected: big.NewInt(12345), expectedError: nil},
		{name: "Valid string", input: "67890", expected: big.NewInt(67890), expectedError: nil},
		{name: "Valid byte slice", input: []byte("100000"), expected: big.NewInt(100000), expectedError: nil},
		{name: "Invalid string", input: "invalid", expected: nil, expectedError: ErrCanonicalizeStringToBigInt},
		{name: "Invalid byte slice", input: []byte("invalid"), expected: nil, expectedError: ErrCanonicalizeByteToBigInt},
		{name: "Invalid type (int)", input: 12345, expected: nil, expectedError: ErrCanonicalizeBigInt},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := bigIntCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(*big.Int))
			}
		})
	}
}

func TestBoolCanonicalizer(t *testing.T) {
	tcs := []struct {
		name          string
		input         any
		expected      bool
		expectedError error
	}{
		{name: "Valid bool true", input: true, expected: true, expectedError: nil},
		{name: "Valid bool false", input: false, expected: false, expectedError: nil},
		{name: "Valid byte slice true", input: []byte{0x01}, expected: true, expectedError: nil},
		{name: "Valid byte slice false", input: []byte{0x00}, expected: false, expectedError: nil},
		{name: "Invalid byte slice", input: []byte{0x02}, expected: false, expectedError: ErrCanonicalizeByteToBool},
		{name: "Invalid type string", input: "true", expected: false, expectedError: ErrCanonicalizeBool},
		{name: "Invalid type int", input: 1, expected: false, expectedError: ErrCanonicalizeBool},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := boolCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(bool))
			}
		})
	}
}

func TestStringCanonicalizer(t *testing.T) {
	tcs := []struct {
		desc          string
		input         any
		expected      string
		expectedError error
	}{
		{desc: "Valid string", input: "hello", expected: "hello", expectedError: nil},
		{desc: "Valid byte slice", input: []byte("world"), expected: "world", expectedError: nil},
		{desc: "Valid empty string", input: "", expected: "", expectedError: nil},
		{desc: "Valid empty byte slice", input: []byte{}, expected: "", expectedError: nil},
		{desc: "Invalid type int", input: 123, expected: "", expectedError: ErrCanonicalizeString},
		{desc: "Invalid type bool", input: true, expected: "", expectedError: ErrCanonicalizeString},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := stringCanonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(string))
			}
		})
	}
}

func TestUint64Canonicalizer(t *testing.T) {
	tcs := []struct {
		desc          string
		input         any
		expected      uint64
		expectedError error
	}{
		{desc: "Valid uint64", input: uint64(12345), expected: 12345, expectedError: nil},
		{desc: "Valid float64", input: float64(67890), expected: 67890, expectedError: nil},
		{desc: "Valid byte slice", input: []byte{0, 0, 0, 0, 0, 0x10, 0, 0}, expected: 0x100000, expectedError: nil},
		{desc: "Invalid float64 (not an integer)", input: float64(123.45), expected: 0, expectedError: ErrCanonicalizeFloatToUint64},
		{desc: "Invalid type (string)", input: "12345", expected: 0, expectedError: ErrCanonicalizeUint64},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := uint64Canonicalizer(tc.input)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result.(uint64))
			}
		})
	}
}
