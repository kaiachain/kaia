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
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	if block.Header() == nil {
		return errNilHeader
	}
	if len(block.Header().Vote) == 0 {
		return nil // nothing to do
	}

	name, vote, err := newVoteDataFromBytes(block.Header().Vote)

	// if vote.key is in gov.Params, do nothing
	_, ok := gov.Params[gov.ParamName(name)]
	if !ok {
		return nil
	}

	// otherwise, handle the votebyte
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
func (v *ValsetModule) HandleValidatorVote(header *types.Header, voteData *voteData) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	councilAddressList, err := ReadCouncilAddressListFromDb(v.ChainKv, header.Number.Uint64()-1)
	if err != nil {
		return err
	}
	var (
		blockNumber = header.Number.Uint64()
		valset      = subsetCouncilSlice(councilAddressList)
		name        = voteData.Name()
	)

	if err = v.checkConsistency(blockNumber, voteData); err != nil {
		return err
	}

	author, err := v.chain.Engine().Author(header)
	if err != nil {
		return err
	}
	if voteData.voter != author {
		return errInvalidVoter
	}

	// GovernanceAddValidator:    appends new validators only if it does not exist in current valSet
	// GovernanceRemoveValidator: delete the existing validator only if it exists in current valSet
	for _, address := range voteData.Value() {
		if author == address {
			return errInvalidVoteValue
		}
		idx := valset.getIdxByAddress(address)
		switch gov.ValSetVoteKeyMap[name] {
		case gov.GovernanceAddValidator:
			if idx == -1 {
				valset = append(valset, address)
			} else {
				return errInvalidVoteValue
			}
		case gov.GovernanceRemoveValidator:
			if idx != -1 {
				valset = append(valset[:idx], valset[idx+1:]...)
			} else {
				return errInvalidVoteValue
			}
		}
	}

	// update valSet db
	if err = WriteCouncilAddressListToDb(v.ChainKv, blockNumber, valset); err != nil {
		return err
	}

	// remove it from myVotes only if the block is proposed by me
	if author == v.nodeAddress {
		for idx, myVote := range v.myVotes {
			if voteData.Equal(myVote) {
				v.myVotes = append(v.myVotes[:idx], v.myVotes[idx+1:]...)
				break
			}
		}
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
