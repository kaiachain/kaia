// Copyright 2024 The Kaia Authors
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
package system

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestContractCallerForMultiCall(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	backend := backends.NewSimulatedBackend(nil)
	originCode := MultiCallCode
	defer func() {
		MultiCallCode = originCode
		backend.Close()
	}()

	// Temporary code injection
	MultiCallCode = MultiCallMockCode

	header := backend.BlockChain().CurrentHeader()
	chain := backend.BlockChain()

	tempState, _ := backend.BlockChain().StateAt(header.Root)
	caller, _ := NewMultiCallContractCaller(tempState, chain, header)
	ret, err := caller.MultiCallStakingInfo(&bind.CallOpts{BlockNumber: header.Number})
	assert.Nil(t, err)

	// Does not affect the original state
	state, _ := backend.BlockChain().StateAt(header.Root)
	assert.Equal(t, []byte(nil), state.GetCode(MultiCallAddr))

	// Mock data
	assert.Equal(t, 5, len(ret.TypeList))
	assert.Equal(t, 5, len(ret.AddressList))
	assert.Equal(t, 1, len(ret.StakingAmounts))

	expectedAddress := []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000F00"),
		common.HexToAddress("0x0000000000000000000000000000000000000F01"),
		common.HexToAddress("0x0000000000000000000000000000000000000F02"),
		common.HexToAddress("0x0000000000000000000000000000000000000F03"),
		common.HexToAddress("0x0000000000000000000000000000000000000F04"),
	}
	for i := 0; i < 5; i++ {
		assert.Equal(t, uint8(i), ret.TypeList[i])
		assert.Equal(t, expectedAddress[i], ret.AddressList[i])
	}
	assert.Equal(t, new(big.Int).Mul(big.NewInt(7_000_000), big.NewInt(params.KAIA)), ret.StakingAmounts[0])
}
