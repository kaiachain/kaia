package impl

import (
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
