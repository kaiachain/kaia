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

//go:generate mockgen -destination=mock/staking.go -package=mock github.com/kaiachain/kaia/kaiax/staking StakingModule
type StakingModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.RewindableModule

	// GetStakingInfo returns the staking info to be used for the given block number.
	// This is the most commonly used getter.
	GetStakingInfo(num uint64) (*StakingInfo, error)

	// GetStakingInfoFromDB returns the staking info from the database.
	// The given number indicates the number that the staking info is measured from, not when the staking info is used.
	// This is useful when syncing the staking info database over p2p.
	GetStakingInfoFromDB(sourceNum uint64) (*StakingInfo, error)

	// SideState management
	AllocSideStateRef() uint64
	FreeSideStateRef(refId uint64)
	AddSideState(refId uint64, header *types.Header, statedb *state.StateDB) error
}

type StakingModuleHost interface {
	RegisterStakingModule(module StakingModule)
}
