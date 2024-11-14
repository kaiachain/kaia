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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	if block.Header() == nil {
		return errNilHeader
	}
	if len(block.Header().Vote) == 0 {
		return nil // nothing to do
	}

	var vb headergov.VoteBytes = block.Header().Vote
	vote, err := vb.ToVoteData()
	if err != nil {
		return err
	}

	if err = v.HandleValidatorVote(block.Header(), vote); err != nil {
		return err
	}
	return nil
}

// HandleValidatorVote handles addvalidator or removevalidator votes and remove them from MyVotes.
// If succeed, the voteBlk and councilAddressList db is updated.
func (v *ValsetModule) HandleValidatorVote(header *types.Header, vote headergov.VoteData) error {
	// if vote.key is in gov.Params, do nothing
	_, ok := gov.Params[vote.Name()]
	if ok {
		return nil
	}

	councilAddressList, err := ReadCouncilAddressListFromDb(v.ChainKv, header.Number.Uint64()-1)
	if err != nil {
		return err
	}

	var (
		blockNumber = header.Number.Uint64()
		c           = valset.AddressList(councilAddressList)
		name        = vote.Name()
		govNode     = v.headerGov.EffectiveParamSet(blockNumber).GoverningNode
	)

	// GovernanceAddValidator:    appends new validators only if it does not exist in current valSet
	// GovernanceRemoveValidator: delete the existing validator only if it exists in current valSet
	var addresses []common.Address
	_, ok = vote.Value().(common.Address)
	if ok {
		addresses = []common.Address{vote.Value().(common.Address)}
	} else {
		addresses = vote.Value().([]common.Address)
	}
	for _, address := range addresses {
		if address == govNode {
			return errInvalidVoteValue
		}
		idx := c.GetIdxByAddress(address)
		switch name {
		case gov.GovernanceAddValidator:
			if idx == -1 {
				c = append(c, address)
			} else {
				return errInvalidVoteValue
			}
		case gov.GovernanceRemoveValidator:
			if idx != -1 {
				c = append(c[:idx], c[idx+1:]...)
			} else {
				return errInvalidVoteValue
			}
		}
	}

	// update valSet db
	if err = WriteCouncilAddressListToDb(v.ChainKv, blockNumber, c); err != nil {
		return err
	}
	return nil
}

func (v *ValsetModule) RewindTo(block *types.Block) {
	// TODO-kaiax-valset: delete
	logger.Info("NoopModule RewindTo", "blockNum", block.Header().Number.Uint64())
}

func (v *ValsetModule) RewindDelete(hash common.Hash, num uint64) {
	logger.Info("NoopModule RewindDelete", "num", num)
}
