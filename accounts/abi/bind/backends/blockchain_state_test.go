package backends

import (
	"context"
	"testing"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestStateBlockchainCallContractLeavesStateUnchanged(t *testing.T) {
	bc := newTestBlockchain()
	statedb, err := bc.State()
	assert.NoError(t, err)

	c, err := NewStateBlockchainContractBackend(bc, statedb)
	assert.NoError(t, err)

	data, _ := parsedAbi1.Pack("receive", []byte("X"))
	nonceBefore := statedb.GetNonce(common.Address{})
	rootBefore := statedb.IntermediateRoot(true)

	ret, err := c.CallContract(context.Background(), kaia.CallMsg{
		From: common.Address{},
		To:   &code1Addr,
		Gas:  1000000,
		Data: data,
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedReturn, ret)

	assert.Equal(t, nonceBefore, statedb.GetNonce(common.Address{}))
	assert.Equal(t, rootBefore, statedb.IntermediateRoot(true))
}
