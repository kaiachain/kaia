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
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ (valset.ValsetModule) = &ValsetModule{}

	logger = log.NewModuleLogger(log.KaiaxValset)
)

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	CurrentBlock() *types.Block
	Config() *params.ChainConfig
	Engine() consensus.Engine
}

type headerGov interface {
	// TODO-kaiax-valset: do not read valset voting vlks from headergov. read valset voting blks from valset db.
	//   voting blks is updated when adding validator.
	GetLatestValidatorVote(num uint64) (uint64, headergov.VoteData)
	GetMyVotes() []headergov.VoteData
	EffectiveParamSet(blockNum uint64) gov.ParamSet
	PopMyVotes(idx int)
}

type stakingInfo interface {
	GetStakingInfo(num uint64) (*staking.StakingInfo, error)
}

type InitOpts struct {
	ChainKv     database.Database
	Chain       chain
	HeaderGov   headerGov
	StakingInfo stakingInfo
	NodeAddress common.Address
}

type ValsetModule struct {
	ChainKv     database.Database
	chain       chain
	headerGov   headerGov
	stakingInfo stakingInfo
	nodeAddress common.Address

	// caches
	proposers lru.Cache
	voteBlks  []uint64
}

func NewValsetModule() *ValsetModule {
	return &ValsetModule{}
}

func (v *ValsetModule) Init(opts *InitOpts) error {
	if opts == nil {
		return errInitUnexpectedNil
	}
	v.ChainKv = opts.ChainKv
	v.chain = opts.Chain
	v.headerGov = opts.HeaderGov
	v.stakingInfo = opts.StakingInfo
	v.nodeAddress = opts.NodeAddress
	return nil
}

func (v *ValsetModule) Start() error {
	// TODO-kaiax-valset: check the genesis is stored. if stored, check the valset list for genesis is stored.
	return nil
}

func (v *ValsetModule) Stop() {
	logger.Info("ValsetModule Stop")
}
