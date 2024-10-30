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

package supply

import (
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ supply.SupplyModule = &SupplyModule{}

	logger = log.NewModuleLogger(log.KaiaxSupply)
)

type InitOpts struct {
	ChainKv      database.Database
	Chain        backends.BlockChainForCaller
	RewardModule reward.RewardModule
}

type SupplyModule struct {
	InitOpts
}

func NewSupplyModule() *SupplyModule {
	return &SupplyModule{}
}

func (s *SupplyModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainKv == nil || opts.Chain == nil || opts.RewardModule == nil {
		return supply.ErrInitUnexpectedNil
	}
	s.InitOpts = *opts
	return nil
}

func (s *SupplyModule) Start() error {
	return nil
}

func (s *SupplyModule) Stop() {
}
