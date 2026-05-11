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

package valset

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/valset ValsetModule
type ValsetModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ExecutionModule
	kaiax.RewindableModule
	kaiax.BlockStateModule
	kaiax.HeaderModule

	GetCouncil(num uint64) ([]common.Address, error)
	GetCommittee(num uint64, round uint64) ([]common.Address, error)
	GetQualifiedValidators(num uint64) ([]common.Address, error)
	GetDemotedValidators(num uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)

	// GetCandTesting returns nodes in the CandTesting state at block `num`.
	// Stub until permissionless system contracts (AddressBookV2) land; currently returns an empty slice.
	GetCandTesting(num uint64) ([]common.Address, error)
	GetCNPeers(num uint64) ([]common.Address, error)
	GetHeaderGovVoters(num uint64) ([]common.Address, error)
	GetNodeByState(num uint64, states []NodeState) (NodeMap, error)

	// Permissionless additions (post-fork; pre-fork these are no-ops).
	WriteTransitionToABv2(vmenv *vm.EVM, header *types.Header, state *state.StateDB) error
	InstallABv2(vmenv *vm.EVM, header *types.Header, state *state.StateDB) error
}

// ValsetModuleHost is implemented by the miner. The worker reads CandTesting
// from it to fill header.VRank at epoch-start blocks.
type ValsetModuleHost interface {
	RegisterValsetModule(module ValsetModule)
}
