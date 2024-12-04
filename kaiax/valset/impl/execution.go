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

	if err := v.HandleValidatorVote(header.Number.Uint64(), header.Vote); err != nil {
		return err
	}
	return nil
}

// HandleValidatorVote handles addvalidator or removevalidator votes.
// If succeeded, the voteBlk and council db is updated.
func (v *ValsetModule) HandleValidatorVote(blockNumber uint64, voteByte []byte) error {
	council, err := v.GetCouncil(blockNumber)
	if err != nil {
		return err
	}
	govNode := v.governance.EffectiveParamSet(blockNumber).GoverningNode
	council, err = applyValSetVote(voteByte, council, govNode)
	if council == nil {
		// if err is nil, it means there's no valSet vote to handle. otherwise, there's issue during handling
		return err
	}

	// write new record at council db and update the valSet vote block db
	if err = writeCouncil(v.ChainKv, blockNumber, council); err != nil {
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

	// AddValidator:    appends new validators only if it does not exist in current valSet
	// RemoveValidator: delete the existing validator only if it exists in current valSet
	for _, address := range addresses {
		if address == govNode {
			return nil, errInvalidVoteValue
		}
		idx := c.GetIdxByAddress(address)

		//nolint:exhaustive
		switch vote.Name() {
		case gov.AddValidator:
			if idx == -1 {
				c = append(c, address)
			} else {
				return nil, errInvalidVoteValue
			}
		case gov.RemoveValidator:
			if idx != -1 {
				c = append(c[:idx], c[idx+1:]...)
			} else {
				return nil, errInvalidVoteValue
			}
		default:
			return nil, errInvalidVoteKey
		}
	}

	return c, nil
}
