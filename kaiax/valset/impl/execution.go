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
	"github.com/kaiachain/kaia/kaiax/valset"
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	header := block.Header()
	num := header.Number.Uint64()
	if num == 0 {
		return nil
	}

	// Ingest validator vote
	var (
		council valset.CommonAddressSet
		err     error
	)
	if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
		council, err = v.getCouncilPermissionless(num)
	} else {
		council, err = v.getCouncilPermissioned(num)
	}
	if err != nil {
		return err
	}
	// TODO-Permissionless: Revisit me and refine
	governingNode := v.GovModule.GetParamSet(num).GoverningNode
	if applyVote(header, council, governingNode) {
		insertValidatorVoteBlockNums(v.ChainKv, num)
		if v.Chain.Config().IsPermissionlessForkEnabled(new(big.Int).SetUint64(num)) {
			writeCouncilPermissionless(v.ChainKv, num, council)
		} else {
			writeCouncilPermissioned(v.ChainKv, num, council)
		}
		v.validatorVoteBlockNumsCache = nil
	}

	return nil
}

func (v *ValsetModule) RewindTo(block *types.Block) {
	trimValidatorVoteBlockNums(v.ChainKv, block.Header().Number.Uint64())
	trimValidatorStateChangeBlockNums(v.ChainKv, block.Header().Number.Uint64())
	v.validatorVoteBlockNumsCache = nil
	v.validatorStateChangeBlockNumsCache = nil
}

func (v *ValsetModule) RewindDelete(hash common.Hash, num uint64) {
	deleteCouncil(v.ChainKv, num)
}
