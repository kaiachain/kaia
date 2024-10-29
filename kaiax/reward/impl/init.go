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
	"math/big"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
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
	StateAt(root common.Hash) (*state.StateDB, error)
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

// TODO-kaiax: remove FromLegacy after introducing kaiax/gov
type LegacyGovEngine interface {
	EffectiveParams(num uint64) (*params.GovParamSet, error)
}

type LegacyGovModule struct {
	engine LegacyGovEngine
}

// Wraps LegacyGovEngine to implement govModule
func FromLegacy(engine LegacyGovEngine) reward.GovModule {
	return &LegacyGovModule{engine: engine}
}

func (l *LegacyGovModule) EffectiveParamSet(num uint64) gov.ParamSet {
	pset, err := l.engine.EffectiveParams(num)
	if err != nil {
		logger.Crit("Failed to get effective params", "num", num, "err", err)
		return gov.ParamSet{}
	}
	// Only implement the ones needed in RewardConfig.
	return gov.ParamSet{
		ProposerPolicy: pset.Policy(),
		UnitPrice:      pset.UnitPrice(),
		MintingAmount:  new(big.Int).Set(pset.MintingAmountBig()),
		MinimumStake:   new(big.Int).Set(pset.MinimumStakeBig()),
		DeferredTxFee:  pset.DeferredTxFee(),
		Ratio:          pset.Ratio(),
		Kip82Ratio:     pset.Kip82Ratio(),
	}
}
