package impl

import (
	"fmt"
	"testing"

	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func newHeaderGovAPI(t *testing.T) *headerGovAPI {
	h := newHeaderGovModule(t, &params.ChainConfig{
		Istanbul: &params.IstanbulConfig{
			Epoch: 1000,
		},
	})
	return NewHeaderGovAPI(h)
}

func TestUpperBoundBaseFeeSet(t *testing.T) {
	api := newHeaderGovAPI(t)
	s, err := api.Vote("kip71.upperboundbasefee", uint64(1))
	assert.Equal(t, ErrUpperBoundBaseFee, err)
	assert.Equal(t, "", s)
}

func TestLowerBoundBaseFeeSet(t *testing.T) {
	api := newHeaderGovAPI(t)
	s, err := api.Vote("kip71.lowerboundbasefee", uint64(1e18))
	assert.Equal(t, ErrLowerBoundBaseFee, err)
	assert.Equal(t, "", s)
}

func TestValidatorVotes(t *testing.T) {
	testCases := []struct {
		key         string
		value       string
		expectedErr error
	}{
		{key: "governance.addvalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{key: "governance.addvalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{key: "governance.addvalidator", value: ",0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", expectedErr: ErrInvalidKeyValue},
		{key: "governance.removevalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{key: "governance.removevalidator", value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{key: "governance.removevalidator", value: ",0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", expectedErr: ErrInvalidKeyValue},
	}

	api := newHeaderGovAPI(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_, err := api.Vote(tc.key, tc.value)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
