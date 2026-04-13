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
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/contracts_permissionless/contracts/multicall"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

// setupMultiCallMock creates a SimulatedBackend with MultiCallMockCode injected,
// and returns the caller, call opts, state, and a cleanup function.
func setupMultiCallMock(t *testing.T) (*multicall.MultiCallContractCaller, *bind.CallOpts, *state.StateDB, func()) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	backend := backends.NewSimulatedBackend(nil)
	originCode := MultiCallCode
	MultiCallCode = MultiCallMockCode

	header := backend.BlockChain().CurrentHeader()
	chain := backend.BlockChain()
	st, _ := backend.BlockChain().StateAt(header.Root)
	caller, _ := NewMultiCallContractCaller(st, chain, header)
	callOpts := &bind.CallOpts{BlockNumber: header.Number}

	cleanup := func() {
		MultiCallCode = originCode
		backend.Close()
	}
	return caller, callOpts, st, cleanup
}

func TestContractCallerForMultiCall(t *testing.T) {
	caller, callOpts, state, cleanup := setupMultiCallMock(t)
	defer cleanup()

	ret, err := caller.MultiCallStakingInfo(callOpts)
	assert.Nil(t, err)

	// Does not affect the original state
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
	for i := range 5 {
		assert.Equal(t, uint8(i), ret.TypeList[i])
		assert.Equal(t, expectedAddress[i], ret.AddressList[i])
	}
	assert.Equal(t, new(big.Int).Mul(big.NewInt(7_000_000), big.NewInt(params.KAIA)), ret.StakingAmounts[0])
}

func TestMultiCallStakingInfoPermissionless(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	backend := backends.NewSimulatedBackend(nil)
	originCode := MultiCallCode
	MultiCallCode = MultiCallPermlessMockCode
	defer func() {
		MultiCallCode = originCode
		backend.Close()
	}()

	header := backend.BlockChain().CurrentHeader()
	chain := backend.BlockChain()
	state, _ := backend.BlockChain().StateAt(header.Root)
	caller, _ := NewMultiCallContractCaller(state, chain, header)
	callOpts := &bind.CallOpts{BlockNumber: header.Number}

	ret, err := caller.MultiCallStakingInfoPermissionless(callOpts)
	assert.Nil(t, err)

	// Does not affect the original state
	assert.Equal(t, []byte(nil), state.GetCode(MultiCallAddr))

	// Verify profiles (6 validators with different states)
	assert.Equal(t, 6, len(ret.Profiles))
	assert.Equal(t, common.HexToAddress("0xF00"), ret.Profiles[0].NodeId)
	assert.Equal(t, uint8(valset.ValActive), ret.Profiles[0].State)
	assert.Equal(t, common.HexToAddress("0xF03"), ret.Profiles[1].NodeId)
	assert.Equal(t, uint8(valset.ValPaused), ret.Profiles[1].State)
	assert.Equal(t, common.HexToAddress("0x1000"), ret.Profiles[2].NodeId)
	assert.Equal(t, uint8(valset.ValReady), ret.Profiles[2].State)
	assert.Equal(t, common.HexToAddress("0x2000"), ret.Profiles[3].NodeId)
	assert.Equal(t, uint8(valset.ValExiting), ret.Profiles[3].State)
	assert.Equal(t, common.HexToAddress("0x3000"), ret.Profiles[4].NodeId)
	assert.Equal(t, uint8(valset.CandReady), ret.Profiles[4].State)
	assert.Equal(t, common.HexToAddress("0x4000"), ret.Profiles[5].NodeId)
	assert.Equal(t, uint8(valset.ValActive), ret.Profiles[5].State)

	// Verify staking amounts
	assert.Equal(t, 6, len(ret.StakingAmounts))
	assert.Equal(t, new(big.Int).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA)), ret.StakingAmounts[0])
	assert.Equal(t, new(big.Int).Mul(big.NewInt(10_000_000), big.NewInt(params.KAIA)), ret.StakingAmounts[1])

	// Verify fund addresses
	assert.Equal(t, common.HexToAddress("0x0a01"), ret.KefAddr)
	assert.Equal(t, common.HexToAddress("0x0a02"), ret.KifAddr)
	assert.Equal(t, common.HexToAddress("0x0a03"), ret.KpfAddr)
}
