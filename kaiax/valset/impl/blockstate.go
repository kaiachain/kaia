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

package impl

import (
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
)

// InitializeState writes the post-fork node-state diff into the
// AddressBookV2 contract via a system tx. No-op pre-fork. Runs before user
// transactions.
func (v *ValsetModule) InitializeState(header *types.Header, statedb *state.StateDB) {
	config := v.Chain.Config()
	if !config.IsPermissionlessForkEnabled(header.Number) {
		return
	}
	context := blockchain.NewEVMBlockContext(header, v.Chain, nil)
	vmenv := vm.NewEVM(context, vm.TxContext{}, statedb, config, &vm.Config{})
	if err := v.WriteTransitionToABv2(vmenv, header, statedb); err != nil {
		logger.Error("Failed to apply node transition to ABv2", "number", header.Number.Uint64(), "err", err)
	}
}

// FinalizeState installs and initializes the AddressBookV2 proxy at the HF-1 block.
func (v *ValsetModule) FinalizeState(header *types.Header, statedb *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	config := v.Chain.Config()
	if !config.IsPermissionlessForkBlockParent(header.Number) {
		return nil
	}
	context := blockchain.NewEVMBlockContext(header, v.Chain, nil)
	vmenv := vm.NewEVM(context, vm.TxContext{}, statedb, config, &vm.Config{})
	return v.InstallABv2(vmenv, header, statedb)
}
