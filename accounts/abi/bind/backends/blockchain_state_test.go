// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package backends

import (
	"context"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain/types"
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

func TestStateHeaderCallContractUsesHeader(t *testing.T) {
	bc := newTestBlockchain()
	statedb, err := bc.State()
	assert.NoError(t, err)

	// A header that is not in the chain yet, as during block processing.
	header := types.CopyHeader(bc.CurrentHeader())
	header.Number = new(big.Int).Add(header.Number, common.Big1)
	header.ParentHash = bc.CurrentHeader().Hash()

	c, err := NewStateHeaderContractBackend(bc, statedb, header)
	assert.NoError(t, err)

	data, _ := parsedAbi1.Pack("receive", []byte("X"))
	call := kaia.CallMsg{From: common.Address{}, To: &code1Addr, Gas: 1000000, Data: data}

	for _, num := range []*big.Int{nil, header.Number} {
		ret, err := c.CallContract(context.Background(), call, num)
		assert.NoError(t, err)
		assert.Equal(t, expectedReturn, ret)
	}

	_, err = c.CallContract(context.Background(), call, bc.CurrentHeader().Number)
	assert.ErrorIs(t, err, errBlockDoesNotExist)

	_, err = NewStateHeaderContractBackend(bc, statedb, nil)
	assert.Error(t, err)
}
