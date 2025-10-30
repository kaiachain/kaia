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

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// ActiveSystemContracts returns the currently active system contracts at the
// given block number.
func ActiveSystemContracts(c *params.ChainConfig, genesis common.Hash, head *big.Int) map[string]common.Address {
	active := make(map[string]common.Address)

	if c.IsPragueForkEnabled(head) {
		active["HISTORY_STORAGE_ADDRESS"] = params.HistoryStorageAddress
	}
	if c.IsKip160ForkEnabled(head) {
		active["KIP160"] = c.Kip160ContractAddress
	}
	if c.IsRandaoForkEnabled(head) {
		active["REGISTRY"] = RegistryAddr
	}
	if c.IsKip103ForkEnabled(head) {
		active["KIP103"] = c.Kip103ContractAddress
	}

	// These contracts are active from genesis.
	if genesis == params.MainnetGenesisHash {
		active["MAINNET_CREDIT"] = MainnetCreditAddr
	}
	active["ADDRESS_BOOK"] = AddressBookAddr

	return active
}
