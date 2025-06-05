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
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/accounts/abi/bind"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/common"
	contracts "github.com/kaiachain/kaia/v2/contracts/contracts/system_contracts/kip113"
	proxycontracts "github.com/kaiachain/kaia/v2/contracts/contracts/system_contracts/proxy"
	testcontracts "github.com/kaiachain/kaia/v2/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/crypto/bls"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/stretchr/testify/assert"
)

func TestReadKip113(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KAIA),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)

		nodeId        = common.HexToAddress("0xaaaa")
		_, pub1, pop1 = makeBlsKey()
		_, pub2, _    = makeBlsKey()
	)

	// Deploy Proxy contract
	transactor, contractAddr := deployKip113Mock(t, sender, backend)

	// With a valid record
	transactor.Register(sender, nodeId, pub1, pop1)
	backend.Commit()

	caller, _ := contracts.NewSimpleBlsRegistryCaller(contractAddr, backend)

	opts := &bind.CallOpts{BlockNumber: nil}
	owner, _ := caller.Owner(opts)
	assert.Equal(t, sender.From, owner)
	t.Logf("owner: %x", owner)

	infos, err := ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, pub1, infos[nodeId].PublicKey)
	assert.Equal(t, pop1, infos[nodeId].Pop)
	assert.Nil(t, infos[nodeId].VerifyErr) // valid PublicKeyInfo

	// With an invalid record
	// Another register() call for the same nodeId overwrites the existing info.
	transactor.Register(sender, nodeId, pub2, pop1) // pub vs. pop mismatch
	backend.Commit()

	// Returns one record with VerifyErr.
	infos, err = ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.ErrorIs(t, infos[nodeId].VerifyErr, ErrKip113BadPop) // invalid PublicKeyInfo
}

func TestAllocKip113(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		nodeId1       = common.HexToAddress("0xaaaa")
		nodeId2       = common.HexToAddress("0xbbbb")
		_, pub1, pop1 = makeBlsKey()
		_, pub2, pop2 = makeBlsKey()

		abi, _   = testcontracts.KIP113MockMetaData.GetAbi()
		input, _ = abi.Pack("initialize")

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KAIA),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)
	)
	kip113Addr, _, _, _ := testcontracts.DeployKIP113Mock(sender, backend)
	backend.Commit()
	contractAddr, _, _, _ := proxycontracts.DeployERC1967Proxy(sender, backend, kip113Addr, input)
	backend.Commit()
	var (
		allocProxyStorage  = AllocProxy(kip113Addr)
		allocKip113Storage = AllocKip113Proxy(AllocKip113Init{
			Infos: BlsPublicKeyInfos{
				nodeId1: {PublicKey: pub1, Pop: pop1},
				nodeId2: {PublicKey: pub2, Pop: pop2},
			},
			Owner: sender.From,
		})
		allocLogicStorage = AllocKip113Logic()
	)

	// 1. Merge two storage maps
	allocStorage := MergeStorage(allocProxyStorage, allocKip113Storage)

	// 2. Create storage by calling register()
	contract, _ := testcontracts.NewKIP113MockTransactor(contractAddr, backend)

	contract.Register(sender, nodeId1, pub1, pop1)
	contract.Register(sender, nodeId2, pub2, pop2)
	backend.Commit()

	// 3. Compare the two states
	compareStorage(t, backend, contractAddr, allocStorage)
	compareStorage(t, backend, kip113Addr, allocLogicStorage)
}

func deployKip113Mock(t *testing.T, sender *bind.TransactOpts, backend *backends.SimulatedBackend, params ...interface{}) (*testcontracts.KIP113MockTransactor, common.Address) {
	// Prepare input data for ERC1967Proxy constructor
	abi, err := testcontracts.KIP113MockMetaData.GetAbi()
	assert.Nil(t, err)
	data, err := abi.Pack("initialize")
	assert.Nil(t, err)

	// Deploy Proxy contract
	// 1. Deploy KIP113Mock implementation contract
	implAddr, _, _, err := testcontracts.DeployKIP113Mock(sender, backend)
	backend.Commit()
	assert.Nil(t, err)
	t.Logf("KIP113Mock impl at %x", implAddr)

	// 2. Deploy ERC1967Proxy(KIP113Mock.address, _data)
	contractAddr, _, _, err := proxycontracts.DeployERC1967Proxy(sender, backend, implAddr, data)
	backend.Commit()
	assert.Nil(t, err)
	t.Logf("ERC1967Proxy at %x", contractAddr)

	// 3. Attach KIP113Mock contract to the proxy
	transactor, _ := testcontracts.NewKIP113MockTransactor(contractAddr, backend)

	return transactor, contractAddr
}

func compareStorage(t *testing.T, backend *backends.SimulatedBackend, contractAddr common.Address, allocStorage map[common.Hash]common.Hash) {
	execStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(contractAddr, func(key common.Hash, value common.Hash) bool {
		execStorage[key] = value
		return true
	})

	for k, v := range allocStorage {
		assert.Equal(t, v.Hex(), execStorage[k].Hex(), k.Hex())
		t.Logf("%x %x\n", k, v)
	}
	for k, v := range execStorage {
		assert.Equal(t, v.Hex(), allocStorage[k].Hex(), k.Hex())
	}
}

func makeBlsKey() (priv, pub, pop []byte) {
	ikm := make([]byte, 32)
	rand.Read(ikm)

	sk, _ := bls.GenerateKey(ikm)
	pk := sk.PublicKey()
	sig := bls.PopProve(sk)

	priv = sk.Marshal()
	pub = pk.Marshal()
	pop = sig.Marshal()
	if len(priv) != 32 || len(pub) != 48 || len(pop) != 96 {
		panic("bad bls key")
	}
	return
}
