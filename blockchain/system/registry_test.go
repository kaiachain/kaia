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

	"github.com/kaiachain/kaia/v2/accounts/abi/bind"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/common"
	contracts "github.com/kaiachain/kaia/v2/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/stretchr/testify/assert"
)

func TestReadRegistry(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KAIA),
			},
			RegistryAddr: {
				Code:    RegistryMockCode,
				Balance: common.Big0,
			},
		}
		backend = backends.NewSimulatedBackend(alloc)

		recordName = "AcmeContract"
		recordAddr = common.HexToAddress("0xaaaa")
	)

	// Without a record
	addr, err := ReadActiveAddressFromRegistry(backend, recordName, common.Big0)
	assert.Nil(t, err)
	assert.Equal(t, common.Address{}, addr)

	// Register a record
	contract, err := contracts.NewRegistryMockTransactor(RegistryAddr, backend)
	_, err = contract.Register(sender, recordName, recordAddr, common.Big1)
	assert.Nil(t, err)
	backend.Commit()

	// With the record
	addr, err = ReadActiveAddressFromRegistry(backend, recordName, common.Big1)
	assert.Nil(t, err)
	assert.Equal(t, recordAddr, addr)
}

// Test that AllocRegistry correctly reproduces the storage state
// identical to the state after a series of register() call.
func TestAllocRegistry(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// 1. Create storage with AllocRegistry
	allocStorage := AllocRegistry(&params.RegistryConfig{
		Records: map[string]common.Address{
			"AcmeContract": common.HexToAddress("0xaaaa"),
			"TestContract": common.HexToAddress("0xcccc"),
		},
		Owner: common.HexToAddress("0xffff"),
	})

	// 2. Create storage by calling register()
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KAIA),
			},
			RegistryAddr: {
				Code:    RegistryMockCode,
				Balance: common.Big0,
			},
		}
		backend     = backends.NewSimulatedBackend(alloc)
		contract, _ = contracts.NewRegistryMockTransactor(RegistryAddr, backend)
	)

	contract.Register(sender, "AcmeContract", common.HexToAddress("0xaaaa"), common.Big0)
	contract.Register(sender, "TestContract", common.HexToAddress("0xcccc"), common.Big0)
	contract.TransferOwnership(sender, common.HexToAddress("0xffff"))
	backend.Commit()

	execStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(RegistryAddr, func(key common.Hash, value common.Hash) bool {
		execStorage[key] = value
		return true
	})

	// 3. Compare the two states
	for k, v := range allocStorage {
		assert.Equal(t, v.Hex(), execStorage[k].Hex(), k.Hex())
		t.Logf("%x %x\n", k, v)
	}
	for k, v := range execStorage {
		assert.Equal(t, v.Hex(), allocStorage[k].Hex(), k.Hex())
	}
}
