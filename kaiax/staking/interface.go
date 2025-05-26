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

package staking

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=./mock/staking.go -package=mock github.com/kaiachain/kaia/kaiax/staking StakingModule
type StakingModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	// GetStakingInfo returns the staking info to be used for the given block number.
	// This is the most commonly used getter.
	GetStakingInfo(num uint64) (*StakingInfo, error)

	// Directly access the database.
	// Note that db access is only effective before Kaia hardfork.
	// The given number indicates the number that the staking info is measured from, not when the staking info is used.
	// This is useful when syncing the staking info database over p2p.
	GetStakingInfoFromDB(sourceNum uint64) *StakingInfo
	PutStakingInfoToDB(sourceNum uint64, stakingInfo *StakingInfo)

	// Preload features allow staking info to be preloaded into memory
	// which helps the situation where the state is not available in the database
	// but the state is only in the memory temporarily (e.g., state regen).
	AllocPreloadRef() uint64
	FreePreloadRef(refId uint64)
	PreloadFromState(refId uint64, header *types.Header, statedb *state.StateDB) error
}

type StakingModuleHost interface {
	RegisterStakingModule(module StakingModule)
}
