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
	"bytes"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	if len(block.Header().Vote) == 0 {
		return nil
	}

	var vb headergov.VoteBytes = block.Header().Vote
	vote, err := vb.ToVoteData()
	if err != nil {
		return err
	}
	// other votes will be handled in headergov
	if vote.Name() != "governance.addvalidator" && vote.Name() == "governance.removevalidator" {
		return nil
	}

	if err = v.HandleValidatorVote(block.NumberU64(), vote); err != nil {
		return err
	}
	return nil
}

// HandleValidatorVote handles addvalidator or removevalidator votes and remove them from MyVotes.
// If succeed, the voteBlk and councilAddressList db is updated.
func (v *ValsetModule) HandleValidatorVote(blockNumber uint64, voteData headergov.VoteData) error {
	prevBlockResult, err := v.getBlockResultsByNumber(blockNumber - 1)
	if err != nil {
		return err
	}

	var (
		valset    = subsetCouncilSlice(prevBlockResult.councilAddrList)
		value     = voteData.Value()
		name      = voteData.Name()
		addresses []common.Address
	)
	// parse the vote
	if addr, ok := value.(common.Address); ok {
		addresses = append(addresses, addr)
	} else if addrs, ok := value.([]common.Address); ok {
		addresses = addrs
	} else {
		logger.Warn("Invalid value Type", "number", blockNumber, "voter", voteData.Voter(), "key", name, "value", value)
	}

	// handle votes and remove it from myVotes
	for _, address := range addresses {
		if name == "governance.addvalidator" {
			valset = append(valset, address)
		} else if name == "governance.removevalidator" {
			idx := valset.getIdxByAddress(address)
			if idx != -1 {
				valset = append(valset[:idx], valset[idx+1:]...)
			}
		}
		if err = WriteCouncilAddressListToDb(v.ChainKv, blockNumber, valset); err != nil {
			return err
		}
	}

	// store new vote block in the db
	v.voteBlks = append(v.voteBlks, blockNumber)
	if err = WriteValidatorVoteDataBlockNums(v.ChainKv, &v.voteBlks); err != nil {
		return err
	}

	// pop my vote from voteData
	if prevBlockResult.author == v.nodeAddress {
		for i, myvote := range v.headerGov.GetMyVotes() {
			if bytes.Equal(myvote.Voter().Bytes(), voteData.Voter().Bytes()) &&
				myvote.Name() == name &&
				myvote.Value() == value {
				v.headerGov.PopMyVotes(i)
				break
			}
		}
	}
	return nil
}

// TODO-kaiax-valset: do we need to check this condition?
func checkVote(address common.Address, isKeyAddValidator bool, valset subsetCouncilSlice) bool {
	idx := valset.getIdxByAddress(address)
	return (idx != -1 && !isKeyAddValidator) || (idx == -1 && isKeyAddValidator)
}

func (v *ValsetModule) RewindTo(block *types.Block) {
	// TODO-kaiax-valset: delete
	logger.Info("NoopModule RewindTo", "blockNum", block.Header().Number.Uint64())
}

func (v *ValsetModule) RewindDelete(hash common.Hash, num uint64) {
	logger.Info("NoopModule RewindDelete", "num", num)
}
