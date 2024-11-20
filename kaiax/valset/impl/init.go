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
	EffectiveParamSet(blockNum uint64) gov.ParamSet
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
	proposers *lru.Cache
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
	cache, err := lru.New(128)
	if err != nil {
		return err
	}
	v.proposers = cache
	return nil
}

func (v *ValsetModule) Start() error {
	var (
		currentBlockNum  = v.chain.CurrentBlock().NumberU64()
		intervalBlockNum = currentBlockNum - currentBlockNum%params.CheckpointInterval
	)

	// valSet initialization at genesis block
	voteBlks := ReadValidatorVoteDataBlockNums(v.ChainKv)
	if voteBlks == nil {
		header := v.chain.GetHeaderByNumber(0)
		if header == nil {
			return errNilHeader
		}
		istanbulExtra, err := types.ExtractIstanbulExtra(header)
		if err != nil {
			return errExtractIstanbulExtra
		}
		if err = WriteCouncilAddressListToDb(v.ChainKv, 0, istanbulExtra.Validators); err != nil {
			return err
		}
		if currentBlockNum == 0 {
			if err = writeLowestScannedCheckpointIntervalNum(v.ChainKv, intervalBlockNum); err != nil {
				return err
			}
		}
	}

	// generate valset db between [lastIntervalBlock, currentBlock)
	lowestSciNum, err := readLowestScannedCheckpointIntervalNum(v.ChainKv)
	if err != nil {
		if err = writeLowestScannedCheckpointIntervalNum(v.ChainKv, intervalBlockNum); err != nil {
			return err
		}
		_, err = v.replayValSetVotes(lowestSciNum, currentBlockNum, true)
		if err != nil {
			return err
		}

		// update lowestScannedNum to figure out if migration is needed or not
		lowestSciNum = intervalBlockNum
	}

	if lowestSciNum == 0 {
		return nil
	}

	go v.migrate()

	return nil
}

func (v *ValsetModule) migrate() {
	lowestSciNum, err := readLowestScannedCheckpointIntervalNum(v.ChainKv)
	if err != nil {
		logger.Error(err.Error())
	}
	for lowestSciNum >= 0 {
		_, err = v.replayValSetVotes(lowestSciNum, lowestSciNum+params.CheckpointInterval, true)
		if err != nil {
			logger.Error(err.Error())
		}
		lowestSciNum = lowestSciNum - 1024
	}
}

// replayValSetVotes replays the valset module handleVote.
// lowestSciNum: the lowest istanbul snapshot checkpoint interval. replay the valSet votes from there.
// targetBlockNum: the target block number where replay ends. it should not exceed current block number.
// writeValSetDb: if it is set, it writes the valSet db
func (v *ValsetModule) replayValSetVotes(intervalBlockNum uint64, targetBlockNum uint64, writeValSetDb bool) ([]common.Address, error) {
	intervalHeader := v.chain.GetHeaderByNumber(intervalBlockNum)
	if intervalHeader == nil {
		return nil, errNilHeader
	}

	c, err := readCouncilAddressListFromIstanbulSnapshot(v.ChainKv, intervalHeader.Hash())
	if err != nil {
		return nil, err
	}

	// replay the header votes until next checkpoint interval
	for num := intervalBlockNum + 1; num <= targetBlockNum; num++ {
		header := v.chain.GetHeaderByNumber(num)
		if header == nil {
			return nil, errNilHeader
		}

		// apply addvalidator/removevalidator vote to council
		govNode := v.headerGov.EffectiveParamSet(num).GoverningNode
		cList, err := applyValSetVote(header.Vote, c, govNode)
		if err != nil {
			return nil, err
		}
		if cList == nil {
			continue // nothing to do
		}

		c = cList

		if writeValSetDb == false {
			continue
		}
		// update to valSet db (council list, voteBlk)
		if err = WriteCouncilAddressListToDb(v.ChainKv, num, c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (v *ValsetModule) Stop() {
	logger.Info("ValsetModule Stop")
}
