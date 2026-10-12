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

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	num := block.Header().Number.Uint64()
	if num == 0 {
		return nil
	}
	if !v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		return v.postInsertBlockPermissioned(block)
	}
	return v.postInsertBlockPermissionless(block)
}

func (v *ValsetModule) postInsertBlockPermissionless(block *types.Block) error {
	header := block.Header()
	nextNum := header.Number.Uint64() + 1
	parentStatedb, err := v.Chain.StateAt(header.Root)
	if err != nil {
		return err
	}
	_, err = v.getTransitionResult(nextNum, parentStatedb) // to cache the transition result
	return err
}

// postInsertBlockPermissioned ingests the validator vote (if any) carried in
// the block's header and persists the updated council to the DB.
func (v *ValsetModule) postInsertBlockPermissioned(block *types.Block) error {
	header := block.Header()
	num := header.Number.Uint64()
	council, err := v.getCouncilPermissioned(num)
	if err != nil {
		return err
	}
	governingNode := v.GovModule.GetParamSet(num).GoverningNode
	if applyVote(header, council, governingNode) {
		insertValidatorVoteBlockNums(v.ChainKv, num)
		writeCouncil(v.ChainKv, num, council.List())
		v.validatorVoteBlockNumsCache = nil
	}
	return nil
}

func (v *ValsetModule) RewindTo(block *types.Block) {
	trimValidatorVoteBlockNums(v.ChainKv, block.Header().Number.Uint64())
	v.validatorVoteBlockNumsCache = nil
}

func (v *ValsetModule) RewindDelete(hash common.Hash, num uint64) {
	deleteCouncil(v.ChainKv, num)
}
