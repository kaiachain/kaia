// Copyright 2025 The Kaia Authors
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

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestActiveSystemContracts(t *testing.T) {
	mainnetConfig := params.MainnetChainConfig
	mainnetConfig.Kip160ContractAddress = common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	mainnetConfig.Kip103ContractAddress = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	cases := []struct {
		name        string
		head        *big.Int
		chainConfig *params.ChainConfig
		genesisHash common.Hash
		expected    map[string]common.Address
	}{
		{
			name:        "Genesis block - testnet",
			head:        big.NewInt(0),
			chainConfig: params.KairosChainConfig,
			genesisHash: params.KairosGenesisHash,
			expected: map[string]common.Address{
				"ADDRESS_BOOK": AddressBookAddr,
			},
		},
		{
			name:        "Genesis block - mainnet",
			head:        big.NewInt(0),
			chainConfig: mainnetConfig,
			genesisHash: params.MainnetGenesisHash,
			expected: map[string]common.Address{
				"MAINNET_CREDIT": MainnetCreditAddr,
				"ADDRESS_BOOK":   AddressBookAddr,
			},
		},
		{
			name:        "Kip103 head - mainnet",
			head:        big.NewInt(119750400),
			chainConfig: mainnetConfig,
			genesisHash: params.MainnetGenesisHash,
			expected: map[string]common.Address{
				"KIP103":         mainnetConfig.Kip103ContractAddress,
				"MAINNET_CREDIT": MainnetCreditAddr,
				"ADDRESS_BOOK":   AddressBookAddr,
			},
		},
		{
			name:        "Randao head - mainnet",
			head:        big.NewInt(147534000),
			chainConfig: mainnetConfig,
			genesisHash: params.MainnetGenesisHash,
			expected: map[string]common.Address{
				"REGISTRY":       RegistryAddr,
				"KIP103":         mainnetConfig.Kip103ContractAddress,
				"MAINNET_CREDIT": MainnetCreditAddr,
				"ADDRESS_BOOK":   AddressBookAddr,
			},
		},
		{
			name:        "Kip160 head - mainnet",
			head:        big.NewInt(162900480),
			chainConfig: mainnetConfig,
			genesisHash: params.MainnetGenesisHash,
			expected: map[string]common.Address{
				"KIP160":         mainnetConfig.Kip160ContractAddress,
				"REGISTRY":       RegistryAddr,
				"KIP103":         mainnetConfig.Kip103ContractAddress,
				"MAINNET_CREDIT": MainnetCreditAddr,
				"ADDRESS_BOOK":   AddressBookAddr,
			},
		},
		{
			name:        "Prague head - mainnet",
			head:        big.NewInt(190670000),
			chainConfig: mainnetConfig,
			genesisHash: params.MainnetGenesisHash,
			expected: map[string]common.Address{
				"HISTORY_STORAGE_ADDRESS": params.HistoryStorageAddress,
				"KIP160":                  mainnetConfig.Kip160ContractAddress,
				"REGISTRY":                RegistryAddr,
				"KIP103":                  mainnetConfig.Kip103ContractAddress,
				"MAINNET_CREDIT":          MainnetCreditAddr,
				"ADDRESS_BOOK":            AddressBookAddr,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := ActiveSystemContracts(tt.chainConfig, tt.genesisHash, tt.head)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
