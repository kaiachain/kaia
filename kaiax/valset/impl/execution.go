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
	header := block.Header()
	if header == nil {
		return errNilHeader
	}

	// read prevBlock's council address list
	councilAddressList, err := v.GetCouncilAddressList(header.Number.Uint64() - 1)
	if err != nil {
		return err
	}

	if err = v.HandleValidatorVote(header.Number.Uint64(), header.Vote, councilAddressList); err != nil {
		return err
	}
	return nil
}

// HandleValidatorVote handles addvalidator or removevalidator votes and remove them from MyVotes.
// If succeed, the voteBlk and councilAddressList db is updated.
func (v *ValsetModule) HandleValidatorVote(blockNumber uint64, voteByte []byte, c valset.AddressList) error {
	govNode := v.headerGov.EffectiveParamSet(blockNumber).GoverningNode
	cList, err := applyValSetVote(voteByte, c, govNode)
	if cList == nil {
		// if err is nil, it means there's no valSet vote to handle. otherwise, there's issue during handling
		return err
	}

	// update valSet db (council list, voteBlk)
	if err = WriteCouncilAddressListToDb(v.ChainKv, blockNumber, cList); err != nil {
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

func applyValSetVote(vb headergov.VoteBytes, c valset.AddressList, govNode common.Address) ([]common.Address, error) {
	if len(vb) == 0 {
		return nil, nil // nothing to do
	}

	vote, err := vb.ToVoteData()
	if err != nil {
		return nil, err
	}

	// if vote.key is in gov.Params, do nothing
	_, ok := gov.Params[vote.Name()]
	if ok {
		return nil, nil
	}

	var addresses []common.Address
	_, ok = vote.Value().(common.Address)
	if ok {
		addresses = []common.Address{vote.Value().(common.Address)}
	} else {
		addresses = vote.Value().([]common.Address)
	}

	// GovernanceAddValidator:    appends new validators only if it does not exist in current valSet
	// GovernanceRemoveValidator: delete the existing validator only if it exists in current valSet
	for _, address := range addresses {
		if address == govNode {
			return nil, errInvalidVoteValue
		}
		idx := c.GetIdxByAddress(address)
		switch vote.Name() {
		case gov.GovernanceAddValidator:
			if idx == -1 {
				c = append(c, address)
			} else {
				return nil, errInvalidVoteValue
			}
		case gov.GovernanceRemoveValidator:
			if idx != -1 {
				c = append(c[:idx], c[idx+1:]...)
			} else {
				return nil, errInvalidVoteValue
			}
		}
	}

	return c, nil
}

func (v *ValsetModule) simulateValSetVotes() {
}
