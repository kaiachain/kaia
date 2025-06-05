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

package impl

import (
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/consensus"
	"github.com/kaiachain/kaia/v2/kaiax/reward"
	"github.com/kaiachain/kaia/v2/kaiax/staking"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
)

var (
	_ reward.RewardModule = &RewardModule{}

	logger = log.NewModuleLogger(log.KaiaxReward)
)

type blockChain interface {
	CurrentBlock() *types.Block
	GetHeaderByNumber(number uint64) *types.Header
	GetBlockByNumber(number uint64) *types.Block
	GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts
	Engine() consensus.Engine
}

type InitOpts struct {
	ChainConfig   *params.ChainConfig
	Chain         blockChain
	GovModule     reward.GovModule // TODO-kaiax: Restore to gov.GovModule after introducing kaiax/gov
	StakingModule staking.StakingModule
}

type RewardModule struct {
	InitOpts
}

func NewRewardModule() *RewardModule {
	return &RewardModule{}
}

func (r *RewardModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.Chain == nil || opts.GovModule == nil || opts.StakingModule == nil {
		return reward.ErrInitUnexpectedNil
	}
	r.InitOpts = *opts
	return nil
}

func (r *RewardModule) Start() error {
	return nil
}

func (r *RewardModule) Stop() {
}
