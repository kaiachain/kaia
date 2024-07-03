// Copyright 2018 The klaytn Authors
// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package genesis

import (
	"math/big"
	"strings"

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/cmd/homi/extra"
	"github.com/kaiachain/kaia/consensus/clique"
	testcontract "github.com/kaiachain/kaia/contracts/contracts/testing/reward"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
)

type Option func(*blockchain.Genesis)

var logger = log.NewModuleLogger(log.CMDIstanbul)

func Validators(addrs ...common.Address) Option {
	return func(genesis *blockchain.Genesis) {
		extraData, err := extra.Encode("0x00", addrs)
		if err != nil {
			logger.Error("Failed to encode extra data", "err", err)
			return
		}
		genesis.ExtraData = hexutil.MustDecode(extraData)
	}
}

func ValidatorsOfClique(signers ...common.Address) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.ExtraData = make([]byte, clique.ExtraVanity+len(signers)*common.AddressLength+clique.ExtraSeal)
		for i, signer := range signers {
			copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
		}
	}
}

func makeGenesisAccount(addrs []common.Address, balance *big.Int) map[common.Address]blockchain.GenesisAccount {
	alloc := make(map[common.Address]blockchain.GenesisAccount)
	for _, addr := range addrs {
		alloc[addr] = blockchain.GenesisAccount{Balance: balance}
	}
	return alloc
}

func Alloc(addrs []common.Address, balance *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		alloc := makeGenesisAccount(addrs, balance)
		genesis.Alloc = alloc
	}
}

func AllocWithPremainnetContract(addrs []common.Address, balance *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		alloc := makeGenesisAccount(addrs, balance)
		alloc[system.MainnetCreditAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(MainnetCreditBin),
			Balance: big.NewInt(0),
		}
		alloc[system.AddressBookAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(PremainnetAddressBookBin),
			Balance: big.NewInt(0),
		}
		genesis.Alloc = alloc
	}
}

func AllocWithMainnetContract(addrs []common.Address, balance *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		alloc := makeGenesisAccount(addrs, balance)
		alloc[system.MainnetCreditAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(MainnetCreditBin),
			Balance: big.NewInt(0),
		}
		alloc[system.AddressBookAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(MainnetAddressBookBin),
			Balance: big.NewInt(0),
		}
		genesis.Alloc = alloc
	}
}

func AllocWithPretestnetContract(addrs []common.Address, balance *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		alloc := makeGenesisAccount(addrs, balance)
		alloc[system.AddressBookAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(PretestnetAddressBookBin),
			Balance: big.NewInt(0),
		}
		genesis.Alloc = alloc
	}
}

func AllocWithTestnetContract(addrs []common.Address, balance *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		alloc := makeGenesisAccount(addrs, balance)
		alloc[system.AddressBookAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(TestnetAddressBookBin),
			Balance: big.NewInt(0),
		}
		genesis.Alloc = alloc
	}
}

// Patch the hardcoded line in AddressBook.sol:constructContract().
func PatchAddressBook(addr common.Address) Option {
	return func(genesis *blockchain.Genesis) {
		contractAddr := system.AddressBookAddr
		contractAccount, ok := genesis.Alloc[contractAddr]
		if !ok {
			log.Fatalf("No AddressBook to patch")
		}

		codeHex := hexutil.Encode(contractAccount.Code)
		var oldAddr string
		switch codeHex {
		case MainnetAddressBookBin:
			oldAddr = "854ca8508c8be2bb1f3c244045786410cb7d5d0a"
		case TestnetAddressBookBin:
			oldAddr = "88bb3838aa0a140acb73eeb3d4b25a8d3afd58d4"
		case PremainnetAddressBookBin, PretestnetAddressBookBin:
			oldAddr = "fe1ffd5293fc94857a33dcd284fe82bc106be4c7"
		}

		// The hardcoded address appears exactly once, hence Replace(.., 1)
		newAddr := strings.ToLower(addr.Hex()[2:])
		codeHex = strings.Replace(codeHex, oldAddr, newAddr, 1)

		genesis.Alloc[contractAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(codeHex),
			Balance: contractAccount.Balance,
		}
	}
}

func AddressBookMock() Option {
	return func(genesis *blockchain.Genesis) {
		contractAddr := system.AddressBookAddr
		contractAccount, ok := genesis.Alloc[contractAddr]
		if !ok {
			log.Fatalf("No AddressBook to patch")
		}

		code := testcontract.AddressBookMockBinRuntime
		genesis.Alloc[contractAddr] = blockchain.GenesisAccount{
			Code:    common.FromHex(code),
			Balance: contractAccount.Balance,
		}
	}
}

func AllocateRegistry(storage map[common.Hash]common.Hash) Option {
	return func(genesis *blockchain.Genesis) {
		registryCode := system.RegistryCode
		registryAddr := system.RegistryAddr

		genesis.Alloc[registryAddr] = blockchain.GenesisAccount{
			Code:    registryCode,
			Storage: storage,
			Balance: big.NewInt(0),
		}
	}
}

func RegistryMock() Option {
	return func(genesis *blockchain.Genesis) {
		registryMockCode := system.RegistryMockCode
		registryAddr := system.RegistryAddr

		registry, ok := genesis.Alloc[registryAddr]
		if !ok {
			log.Fatalf("No registry to patch")
		}

		genesis.Alloc[registryAddr] = blockchain.GenesisAccount{
			Code:    registryMockCode,
			Storage: registry.Storage,
			Balance: big.NewInt(0),
		}
	}
}

func AllocateKip113(kip113ProxyAddr, kip113LogicAddr common.Address, owner common.Address, proxyStorage, logicStorage map[common.Hash]common.Hash) Option {
	return func(genesis *blockchain.Genesis) {
		proxyCode := system.ERC1967ProxyCode
		logicCode := system.Kip113Code

		genesis.Alloc[kip113ProxyAddr] = blockchain.GenesisAccount{
			Code:    proxyCode,
			Storage: proxyStorage,
			Balance: big.NewInt(0),
		}
		genesis.Alloc[kip113LogicAddr] = blockchain.GenesisAccount{
			Code:    logicCode,
			Storage: logicStorage,
			Balance: big.NewInt(0),
		}
		genesis.Config.RandaoRegistry = &params.RegistryConfig{
			Records: map[string]common.Address{
				system.Kip113Name: kip113ProxyAddr,
			},
			Owner: owner,
		}
	}
}

func Kip113Mock(kip113LogicAddr common.Address) Option {
	return func(genesis *blockchain.Genesis) {
		kip113MockCode := system.Kip113MockCode

		_, ok := genesis.Alloc[kip113LogicAddr]
		if !ok {
			log.Fatalf("No kip113 to patch")
		}
		_, ok = genesis.Config.RandaoRegistry.Records[system.Kip113Name]
		if !ok {
			log.Fatalf("No kip113 record to patch")
		}

		genesis.Alloc[kip113LogicAddr] = blockchain.GenesisAccount{
			Code:    kip113MockCode,
			Balance: big.NewInt(0),
		}
		genesis.Config.RandaoRegistry.Records[system.Kip113Name] = kip113LogicAddr
	}
}

func ChainID(chainID *big.Int) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.ChainID = chainID
	}
}

func UnitPrice(price uint64) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.UnitPrice = price
	}
}

func Istanbul(config *params.IstanbulConfig) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.Istanbul = config
	}
}

func DeriveShaImpl(impl int) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.DeriveShaImpl = impl
	}
}

func Governance(config *params.GovernanceConfig) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.Governance = config
	}
}

func Clique(config *params.CliqueConfig) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.Clique = config
	}
}

func StakingInterval(interval uint64) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.Governance.Reward.StakingUpdateInterval = interval
	}
}

func ProposerInterval(interval uint64) Option {
	return func(genesis *blockchain.Genesis) {
		genesis.Config.Governance.Reward.ProposerUpdateInterval = interval
	}
}
